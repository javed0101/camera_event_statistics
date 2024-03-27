package redis

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/javed0101/cameraevents/config"
	"github.com/javed0101/cameraevents/helper"
	"github.com/javed0101/cameraevents/internal/core/models"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

var redisClient *Redis
var totalEvent string

func GetRedisClient() *Redis {
	if redisClient != nil {
		return redisClient
	}
	return nil
}

func InitRedis() *Redis {
	conf := config.GetConfig()
	opts := redis.Options{
		Addr:       *conf.Redis.HostName,
		MaxRetries: *conf.Redis.MaxRetries,
		DB:         *conf.Redis.DB,
	}
	redisClient = &Redis{
		client: redis.NewClient(&opts),
	}
	log.Info("Initializing redis client: ", *conf.Redis.HostName)
	return redisClient
}

func (r *Redis) AddEventToRedis(ctx context.Context, key *string, redisEvent *models.RedisEvent, pulsarEvent *models.PulsarEvent) error {

	getAllEvents := r.client.HGetAll(ctx, *key)
	result, err := getAllEvents.Result()
	if err != nil {
		log.Error("Error getting key from redis. Error: ", err)
		return err
	}
	if len(result) == 0 {
		redisEvent.StartTime = pulsarEvent.TimeStamp
		log.Info("Inserting camera event into redis with cameraID: ", *pulsarEvent.CamersID)
	} else {
		log.Info("Updating event into redis with cameraID: ", *pulsarEvent.CamersID)
		redisEvent.EndTime = pulsarEvent.TimeStamp
	}
	*redisEvent.Count, _ = strconv.Atoi(result["count"])
	*redisEvent.Count += 1
	upsertCmd := r.client.HSet(ctx, *key, redisEvent)
	_, err = upsertCmd.Result()
	if err != nil {
		log.Error("Failed to insert camera event into redis. Error: ", err)
	} else {
		log.Infof("Successfully inserted camera event into redis with eventID: %s and redis key: %s", *pulsarEvent.EventID, *key)
	}
	r.client.Incr(context.Background(), totalEvent).Result()
	return err
}

func (r *Redis) GetEventFromRedis(ctx context.Context, key *string) (*models.RedisEvent, error) {
	// if err := redisClient.client.Ping(ctx); err != nil {
	// 	return nil, nil
	// }
	getEvent := r.client.HGetAll(ctx, *key)
	result, err := getEvent.Result()
	if err != nil || result == nil || len(result) == 0 {
		log.Error("Error getting key from redis. Error: ", err)
		return nil, err
	}
	redisEvent := new(models.RedisEvent)
	count, _ := strconv.Atoi(result["count"])
	redisEvent.Count = &count
	redisEvent.StartTime = helper.StringPointer(result["startTime"])
	redisEvent.EndTime = helper.StringPointer(result["endTime"])
	return redisEvent, nil
}

func (r *Redis) CloseDBConnection() error {
	if r != nil && r.client != nil {
		log.Info("Closing the redis connection")
		if err := r.client.Close(); err != nil {
			log.Infof("Closing the redis connection failed with error: [%v]", err)
			return err
		}
	}
	return nil
}

func GetTotalEventCount(rc *Redis) *int {
	counterValue, _ := rc.client.Get(context.Background(), totalEvent).Result()
	totalEvent, _ := strconv.Atoi(counterValue)
	return &totalEvent
}

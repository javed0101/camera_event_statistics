package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/gofiber/fiber/v2/log"
	configmanager "github.com/javed0101/cameraevents/config"
	"github.com/javed0101/cameraevents/helper"
	"github.com/javed0101/cameraevents/internal/core/models"
	"github.com/javed0101/cameraevents/internal/sources/redis"
)

type CameraEvent struct {
	Topic   *string
	Channel chan CameraEvent
}

func getCameraEventProcessor(topic string, jobChan chan CameraEvent) IJobProcessor {
	return &CameraEvent{Topic: &topic, Channel: jobChan}
}

func (cam *CameraEvent) ProcessJob() {
	cam.Channel <- *cam
	log.Info("Initializing workers for consuming events from topic: ", *cam.Topic)
	config := configmanager.GetConfig()
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               *config.Pulsar.HostName,
		ConnectionTimeout: time.Second * 50,
	})
	if err != nil {
		fmt.Println("Error creating Pulsar client:", err)
		return
	}
	defer client.Close()

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            *cam.Topic,
		SubscriptionName: *config.Pulsar.SubscriptionName,
		Type:             pulsar.Shared,
	})
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}

	defer consumer.Close()

	for {
		message, err := consumer.Receive(context.Background())
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}
		fmt.Printf("Received message: '%s' with ID: %s at: '%s'\n", string(message.Payload()), message.ID(), message.PublishTime())
		var pulsarEvent *models.PulsarEvent
		err = json.Unmarshal(message.Payload(), &pulsarEvent)
		if err != nil {
			log.Info("Error unmarshalling pulsar event. Error: ", err.Error())
			continue
		}
		fmt.Println("CameraID:", *pulsarEvent.Info.Event.CameraID, "******************", "Event Type:", *pulsarEvent.Info.Event.EventType)
		redisEvent := new(models.RedisEvent)
		redisEvent.Count = helper.IntPointer(0)
		redisEvent.EndTime = pulsarEvent.Info.Event.Timestamp
		ctx := context.Background()
		PushEvents(ctx, redisEvent, pulsarEvent)
		if err != nil {
			log.Fatalf("Error inserting event into redis. Error: ", err)
			continue
		}
		consumer.Ack(message)
	}
}

func PushEvents(ctx context.Context, redisEvent *models.RedisEvent, pulsarEvent *models.PulsarEvent) {
	redisClient := redis.GetRedisClient()
	redisKey := helper.StringPointer(*pulsarEvent.Info.Event.CameraID + ":" + *pulsarEvent.Info.Event.EventType)
	log.Info("Redis Key: ", *redisKey)
	_ = redisClient.AddEventToRedis(ctx, redisKey, redisEvent, pulsarEvent)
}

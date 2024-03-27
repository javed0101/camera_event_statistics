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
	Topic *string
}

func getCameraEventProcessor(topic string) IJobProcessor {
	return &CameraEvent{Topic: &topic}
}

func (cam *CameraEvent) ProcessJob() {
	log.Info("Initializing workers for consuming events from topic: ", cam.Topic)
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
			log.Fatalf("Error unmarshalling pulsar event. Error: ", err.Error())
		}
		redisEvent := new(models.RedisEvent)
		redisEvent.Count = helper.IntPointer(0)
		redisEvent.EndTime = pulsarEvent.TimeStamp
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
	key := helper.StringPointer(*pulsarEvent.CamersID + *pulsarEvent.EventType)
	_ = redisClient.AddEventToRedis(ctx, key, redisEvent, pulsarEvent)
}

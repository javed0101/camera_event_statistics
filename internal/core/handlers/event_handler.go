package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/javed0101/cameraevents/helper"
	"github.com/javed0101/cameraevents/internal/sources/redis"
	cameraevents "github.com/javed0101/cameraevents/pkg/contracts/camera_events"
	eenerror "github.com/javed0101/cameraevents/pkg/utils"
)

func CameraEventHandler(c *fiber.Ctx) error {
	cameraEventRequest := new(cameraevents.CameraEventRequest)
	cameraEventRequest.ExtractCameraEvent(c)
	if _, validate := cameraEventRequest.ValidateCameraEvent(c); validate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": eenerror.ErrorInvalidParam,
		})
	}
	rc := redis.GetRedisClient()
	key := helper.StringPointer(*cameraEventRequest.CameraID + *cameraEventRequest.EventType)
	redisData, err := rc.GetEventFromRedis(context.Background(), key)
	// if redisData == nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": eenerror.ErrorRedisDown,
	// 	})
	// }

	defer StatusCodeMetrics(c)

	if err != nil || redisData == nil {
		if err != nil || redisData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": eenerror.ErrorNoEventFound,
			})
		}
	}
	cameraEventResponse := new(cameraevents.CameraEventResponse)
	newUUID := uuid.New()
	cameraEventResponse.RequestID = helper.StringPointer(newUUID.String())
	cameraEventResponse.EventCount = redisData.Count
	cameraEventResponse.TotalEventCount = redis.GetTotalEventCount(rc)
	timestamp := new(cameraevents.TimeStamp)
	timestamp.FirstEvent = redisData.StartTime
	timestamp.LastEvent = redisData.EndTime
	cameraEventResponse.TimeStamp = timestamp

	if err != nil || redisData == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": eenerror.ErrorInvalidParam,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"statusCode": fiber.StatusOK,
		"event":      cameraEventResponse,
	})
}

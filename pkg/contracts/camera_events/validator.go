package cameraevents

import (
	"github.com/eencloud/goeen/api.v3"
	"github.com/gofiber/fiber/v2"
	eenerror "github.com/javed0101/cameraevents/pkg/utils"
)

func (req *CameraEventRequest) ValidateCameraEvent(c *fiber.Ctx) (*CameraEventRequest, *api.Reason) {
	if req.CameraID == nil || req.EventType == nil {
		return req, &eenerror.ErrorInvalidParam
	}
	return req, nil
}

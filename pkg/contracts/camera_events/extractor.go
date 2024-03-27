package cameraevents

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javed0101/cameraevents/pkg/contracts"
	enum "github.com/javed0101/cameraevents/pkg/types"
)

func (req *CameraEventRequest) ExtractCameraEvent(c *fiber.Ctx) (*CameraEventRequest, *error) {
	request := c.Request()
	contentType := request.Header.ContentType()
	if request.Header.ContentLength() > 0 {
		_, err := contracts.ContentDecoder(string(contentType))
		if err != nil {
			return nil, err
		}
		// if err := decoder(request.Body).Decode(req); err != nil {
		// 	if _, ok := err.(Error); ok {
		// 		return ErrBadRequestInvalidParameter("request format is invalid")
		// 	}
		// 	return ErrBadRequestInvalidBody(err.Error())
		// }
	}
	if req == nil {
		req = &CameraEventRequest{}
	}
	if req.CameraID == nil {
		cameraID := c.Query(enum.CAMERA_ID)
		if cameraID != "" {
			req.CameraID = &cameraID
		}
	}
	if req.EventType == nil {
		eventType := c.Query(enum.EVENT_TYPE)
		if eventType != "" {
			req.EventType = &eventType
		}
	}
	return req, nil
}

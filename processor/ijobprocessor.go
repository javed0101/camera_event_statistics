package processor

import configmanager "github.com/javed0101/cameraevents/config"

type IJobProcessor interface {
	ProcessJob()
}

func JobProcessor(topicName string) IJobProcessor {
	config := configmanager.GetConfig()
	switch topicName {
	case *config.Pulsar.Topic.CameraEvent:
		return getCameraEventProcessor(*config.Pulsar.Topic.CameraEvent)
	}
	return nil
}

package models

import "time"

type Alert struct {
	CameraName  string    `json:"cameraName"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

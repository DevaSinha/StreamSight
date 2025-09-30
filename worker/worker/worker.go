package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DevaSinha/StreamSight/worker/config"
	"github.com/DevaSinha/StreamSight/worker/models"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

func RunWorker(cfg config.Config, alertWSURL string) {
	alertCh := make(chan models.Alert, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		alertSender(alertWSURL, alertCh)
	}()

	for {
		cameras, err := fetchCameras(cfg.ApiEndpoint)
		if err != nil {
			log.Printf("Failed to fetch cameras: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, camera := range cameras {
			go processCamera(camera, alertCh)
		}

		time.Sleep(5 * time.Minute)
	}
}

func alertSender(wsURL string, alertCh <-chan models.Alert) {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			log.Printf("Failed to connect to alert websocket: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Connected to alert WebSocket server")

		for alert := range alertCh {
			msg, err := json.Marshal(alert)
			if err != nil {
				log.Printf("Failed to marshal alert: %v", err)
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Failed to send alert, reconnecting websocket: %v", err)
				conn.Close()
				break
			}
		}

		conn.Close()
		time.Sleep(10 * time.Second)
	}
}

func fetchCameras(apiEndpoint string) ([]models.Camera, error) {
	req, err := http.NewRequest("GET", apiEndpoint+"/cameras", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTkyOTQxNzF9.VwW__6wN4SoO8pOwxv58UeNXO1SVA7DduBUX-qD_4P8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cameras []models.Camera
	if err := json.NewDecoder(resp.Body).Decode(&cameras); err != nil {
		return nil, err
	}
	return cameras, nil
}

func processCamera(camera models.Camera, alertCh chan<- models.Alert) {
	log.Printf("Processing camera %s with URL: %s", camera.Name, camera.URL)

	// Open video capture from RTSP stream
	webcam, err := gocv.VideoCaptureFile(camera.URL)
	if err != nil {
		log.Printf("Error opening video capture for camera %s: %v", camera.Name, err)
		return
	}
	defer webcam.Close()

	// Load face detection classifier
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	// Load the Haar cascade file for face detection
	if !classifier.Load("haarcascade_frontalface_default.xml") {
		log.Printf("Error reading cascade file for camera %s", camera.Name)
		return
	}

	// Prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	log.Printf("Starting face detection for camera: %s", camera.Name)

	frameCount := 0
	for {
		if ok := webcam.Read(&img); !ok {
			log.Printf("Cannot read frame from camera %s", camera.Name)
			break
		}

		if img.Empty() {
			continue
		}

		frameCount++

		// Detect faces every 30 frames (roughly once per second at 30fps)
		if frameCount%30 == 0 {
			// Detect faces
			rects := classifier.DetectMultiScale(img)

			if len(rects) > 0 {
				log.Printf("Found %d face(s) in camera %s", len(rects), camera.Name)

				alert := models.Alert{
					CameraName:  camera.Name,
					Timestamp:   time.Now(),
					Description: "Face detected",
				}

				// Send alert to channel
				select {
				case alertCh <- alert:
					log.Printf("Alert sent for camera %s", camera.Name)
				default:
					log.Printf("Alert channel full, dropping alert for camera %s", camera.Name)
				}
			}
		}

		// Small delay to prevent excessive CPU usage
		time.Sleep(33 * time.Millisecond) // ~30fps
	}

	log.Printf("Camera %s processing stopped", camera.Name)
}

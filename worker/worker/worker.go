package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/DevaSinha/StreamSight/worker/config"
	"github.com/DevaSinha/StreamSight/worker/models"
	"github.com/aler9/gortsplib"
	url2 "github.com/aler9/gortsplib/pkg/url"
	"github.com/gorilla/websocket"
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
			if camera.Active {
				go processCamera(camera, alertCh)
			}
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

		// Connection closed or error; reconnect after short delay
		conn.Close()
		time.Sleep(10 * time.Second)
	}
}

func fetchCameras(apiEndpoint string) ([]models.Camera, error) {
	req, err := http.NewRequest("GET", apiEndpoint+"/cameras", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

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
	u, err := url.Parse(camera.RtspURL)
	if err != nil {
		log.Printf("Invalid RTSP URL %s: %v", camera.RtspURL, err)
		return
	}
	client := gortsplib.Client{}

	err = client.Start(u.Scheme, u.Host)
	if err != nil {
		log.Printf("Failed to start RTSP client for camera %s: %v", camera.Name, err)
		return
	}
	defer client.Close()

	var frameCounter int64

	client.OnPacketRTP = func(ctx *gortsplib.ClientOnPacketRTPCtx) {
		frameCounter++
		if frameCounter%200 == 0 {
			alert := models.Alert{
				CameraName:  camera.Name,
				Timestamp:   time.Now(),
				Description: "Face detected",
			}
			alertCh <- alert
		}
	}
	medias, baseURL, _, err := client.Describe((*url2.URL)(u))
	if err != nil {
		log.Printf("Describe error: %v", err)
		return
	}

	if err = client.SetupAll(medias, baseURL); err != nil {
		log.Printf("Failed to setup tracks for camera %s: %v", camera.Name, err)
		return
	}

	if _, err := client.Play(nil); err != nil {
		log.Printf("Failed to play stream for camera %s: %v", camera.Name, err)
		return
	}

	select {} // keep this goroutine alive
}

package models

type Camera struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	RtspURL  string `json:"rtsp_url"`
	Location string `json:"location"`
	Active   bool   `json:"active"`
}

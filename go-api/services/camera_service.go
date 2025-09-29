package services

import (
	"database/sql"

	"github.com/DevaSinha/StreamSight/go-api/config"
	"github.com/DevaSinha/StreamSight/go-api/models"
)

func GetAllCameras() ([]models.Camera, error) {
	rows, err := config.DB.Query("SELECT id, user_id, name, url, location, created_at FROM cameras")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cameras []models.Camera
	for rows.Next() {
		var cam models.Camera
		err := rows.Scan(&cam.ID, &cam.UserID, &cam.Name, &cam.URL, &cam.Location, &cam.CreatedAt)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, cam)
	}
	return cameras, nil
}

func CreateCamera(userID int, name, url, location string) (*models.Camera, error) {
	var cam models.Camera
	err := config.DB.QueryRow(
		"INSERT INTO cameras (user_id, name, url, location) VALUES ($1,$2,$3,$4) RETURNING id, user_id, name, url, location, created_at",
		userID, name, url, location,
	).Scan(&cam.ID, &cam.UserID, &cam.Name, &cam.URL, &cam.Location, &cam.CreatedAt)

	return &cam, err
}

func GetCameraByID(id string) (*models.Camera, error) {
	var cam models.Camera
	err := config.DB.QueryRow(
		"SELECT id, user_id, name, url, location, created_at FROM cameras WHERE id=$1", id,
	).Scan(&cam.ID, &cam.UserID, &cam.Name, &cam.URL, &cam.Location, &cam.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Camera not found
	}
	return &cam, err
}

func UpdateCamera(id, name, url, location string) error {
	_, err := config.DB.Exec(
		"UPDATE cameras SET name=$1, url=$2, location=$3 WHERE id=$4",
		name, url, location, id,
	)
	return err
}

func DeleteCamera(id string) error {
	_, err := config.DB.Exec("DELETE FROM cameras WHERE id=$1", id)
	return err
}

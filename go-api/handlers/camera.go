package handlers

import (
	"net/http"

	"github.com/DevaSinha/StreamSight/go-api/services"
	"github.com/gin-gonic/gin"
)

func ListCameras(c *gin.Context) {
	cameras, err := services.GetAllCameras()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cameras"})
		return
	}
	c.JSON(http.StatusOK, cameras)
}

func CreateCamera(c *gin.Context) {
	var req struct {
		UserID   int    `json:"user_id"`
		Name     string `json:"name"`
		URL      string `json:"url"`
		Location string `json:"location"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	camera, err := services.CreateCamera(req.UserID, req.Name, req.URL, req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create camera"})
		return
	}

	c.JSON(http.StatusCreated, camera)
}

func GetCamera(c *gin.Context) {
	id := c.Param("id")
	camera, err := services.GetCameraByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	if camera == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "camera not found"})
		return
	}
	c.JSON(http.StatusOK, camera)
}

func UpdateCamera(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Location string `json:"location"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := services.UpdateCamera(id, req.Name, req.URL, req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "camera updated"})
}

func DeleteCamera(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteCamera(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "camera deleted"})
}

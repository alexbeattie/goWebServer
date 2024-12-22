package handlers

import (
	"net/http"
	"strconv"
	"time"
	"github.com/alexbeattie/golangone/models"
	"github.com/alexbeattie/golangone/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPreferences(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	preferences, err := h.service.GetPreferences(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preferences not found"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func (h *Handler) UpdatePreferences(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var preferences models.UserPreferences
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	preferences.UserID = uint(userID)
	if err := h.service.UpdatePreferences(&preferences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

func (h *Handler) GetDevices(c *gin.Context) {
	devices, err := h.service.FetchDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}
func (h *Handler) GetDeviceOdometer(c *gin.Context) {
    deviceID := c.Param("id")
    
    odometerData, err := h.service.FetchDeviceOdometer(deviceID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, odometerData)
}
func (h *Handler) GetDeviceInfo(c *gin.Context) {
    // You can add query params handling if needed
    deviceInfo, err := h.service.FetchDeviceInfo(nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, deviceInfo)
}
func (h *Handler) GetDriveStopRoute(c *gin.Context) {
    deviceID := c.Query("device_id")
    if deviceID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
        return
    }

    fromStr := c.Query("dt_tracker_from")
    toStr := c.Query("dt_tracker_to")
    stopDuration := c.DefaultQuery("stop_duration", "5m0s")

    from, err := time.Parse(time.RFC3339, fromStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format"})
        return
    }

    to, err := time.Parse(time.RFC3339, toStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format"})
        return
    }

    routeData, err := h.service.FetchDriveStopRoute(deviceID, from, to, stopDuration)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, routeData)
}
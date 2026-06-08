package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/store"
)

// DeviceHandler handles device endpoints
type DeviceHandler struct {
	store *store.Store
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		store: globalStore,
	}
}

// List returns all devices
func (h *DeviceHandler) List(c *gin.Context) {
	devices := h.store.ListDevices()
	c.JSON(http.StatusOK, gin.H{
		"count":   len(devices),
		"results": devices,
	})
}

// Create creates a new device
func (h *DeviceHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Site     int    `json:"site" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Status   string `json:"status"`
		Serial   string `json:"serial"`
		AssetTag string `json:"asset_tag"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	device := &store.Device{
		Name:     req.Name,
		Site:     req.Site,
		Type:     req.Type,
		Status:   req.Status,
		Serial:   req.Serial,
		AssetTag: req.AssetTag,
	}
	if device.Status == "" {
		device.Status = "active"
	}

	device = h.store.CreateDevice(device)
	c.JSON(http.StatusCreated, device)
}

// Get returns a device by ID
func (h *DeviceHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid device ID",
		})
		return
	}

	device := h.store.GetDevice(id)
	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "device not found",
		})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Update updates a device
func (h *DeviceHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid device ID",
		})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Status   string `json:"status"`
		Serial   string `json:"serial"`
		AssetTag string `json:"asset_tag"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	device := h.store.UpdateDevice(id, &store.Device{
		Name:     req.Name,
		Type:     req.Type,
		Status:   req.Status,
		Serial:   req.Serial,
		AssetTag: req.AssetTag,
	})

	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "device not found",
		})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Delete deletes a device
func (h *DeviceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid device ID",
		})
		return
	}

	if !h.store.DeleteDevice(id) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "device not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

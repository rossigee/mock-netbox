package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/store"
)

// InterfaceHandler handles interface endpoints
type InterfaceHandler struct {
	store *store.Store
}

// NewInterfaceHandler creates a new interface handler
func NewInterfaceHandler() *InterfaceHandler {
	return &InterfaceHandler{
		store: globalStore,
	}
}

// List returns all interfaces
func (h *InterfaceHandler) List(c *gin.Context) {
	interfaces := h.store.ListInterfaces()
	c.JSON(http.StatusOK, gin.H{
		"count":    len(interfaces),
		"results": interfaces,
	})
}

// Create creates a new interface
func (h *InterfaceHandler) Create(c *gin.Context) {
	var req struct {
		Device  int    `json:"device" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Status  string `json:"status"`
		Enabled bool   `json:"enabled"`
		MTU     int    `json:"mtu"`
		Mode    string `json:"mode"`
		MACAddr string `json:"mac_address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	iface := &store.Interface{
		Device:  req.Device,
		Name:    req.Name,
		Type:    req.Type,
		Status:  req.Status,
		Enabled: req.Enabled,
		MTU:     req.MTU,
		Mode:    req.Mode,
		MACAddr: req.MACAddr,
	}
	if iface.Status == "" {
		iface.Status = "active"
	}
	if !iface.Enabled {
		iface.Enabled = true
	}

	iface = h.store.CreateInterface(iface)
	c.JSON(http.StatusCreated, iface)
}

// Get returns an interface by ID
func (h *InterfaceHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid interface ID",
		})
		return
	}

	iface := h.store.GetInterface(id)
	if iface == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "interface not found",
		})
		return
	}

	c.JSON(http.StatusOK, iface)
}

// Delete deletes an interface
func (h *InterfaceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid interface ID",
		})
		return
	}

	if !h.store.DeleteInterface(id) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "interface not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

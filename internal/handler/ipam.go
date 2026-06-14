package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/store"
)

// IPAMHandler handles IPAM endpoints
type IPAMHandler struct {
	store *store.Store
}

// NewIPAMHandler creates a new IPAM handler
func NewIPAMHandler() *IPAMHandler {
	return &IPAMHandler{
		store: globalStore,
	}
}

// ListPrefixes returns all prefixes
func (h *IPAMHandler) ListPrefixes(c *gin.Context) {
	prefixes := h.store.ListPrefixes()
	c.JSON(http.StatusOK, gin.H{
		"count":   len(prefixes),
		"results": prefixes,
	})
}

// GetPrefix returns a prefix by ID
func (h *IPAMHandler) GetPrefix(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prefix ID"})
		return
	}

	prefix := h.store.GetPrefix(id)
	if prefix == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prefix not found"})
		return
	}

	c.JSON(http.StatusOK, prefix)
}

// ListIPAddresses returns all IP addresses
func (h *IPAMHandler) ListIPAddresses(c *gin.Context) {
	addresses := h.store.ListIPAddresses()
	c.JSON(http.StatusOK, gin.H{
		"count":   len(addresses),
		"results": addresses,
	})
}

// GetIPAddress returns an IP address by ID
func (h *IPAMHandler) GetIPAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address ID"})
		return
	}

	addr := h.store.GetIPAddress(id)
	if addr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IP address not found"})
		return
	}

	c.JSON(http.StatusOK, addr)
}

// AllocateIP allocates an available IP address
func (h *IPAMHandler) AllocateIP(c *gin.Context) {
	var req struct {
		Prefix      string `json:"prefix"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := h.store.AllocateIP(req.Prefix, req.Description)
	if ip == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no available IPs in prefix"})
		return
	}

	c.JSON(http.StatusCreated, ip)
}

// ReleaseIP releases an IP address
func (h *IPAMHandler) ReleaseIP(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address ID"})
		return
	}

	if !h.store.ReleaseIP(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "IP address not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/store"
)

// SiteHandler handles site endpoints
type SiteHandler struct {
	store *store.Store
}

// NewSiteHandler creates a new site handler
func NewSiteHandler() *SiteHandler {
	return &SiteHandler{
		store: globalStore,
	}
}

// List returns all sites
func (h *SiteHandler) List(c *gin.Context) {
	sites := h.store.ListSites()
	c.JSON(http.StatusOK, gin.H{
		"count":    len(sites),
		"results": sites,
	})
}

// Create creates a new site
func (h *SiteHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug"`
		Region      string `json:"region"`
		Facility    string `json:"facility"`
		ASN         int    `json:"asn"`
		TimeZone    string `json:"time_zone"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	site := &store.Site{
		Name:        req.Name,
		Slug:        req.Slug,
		Region:      req.Region,
		Facility:    req.Facility,
		ASN:         req.ASN,
		TimeZone:    req.TimeZone,
		Description: req.Description,
	}

	site = h.store.CreateSite(site)
	c.JSON(http.StatusCreated, site)
}

// Get returns a site by ID
func (h *SiteHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid site ID",
		})
		return
	}

	site := h.store.GetSite(id)
	if site == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "site not found",
		})
		return
	}

	c.JSON(http.StatusOK, site)
}

// Delete deletes a site
func (h *SiteHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid site ID",
		})
		return
	}

	if !h.store.DeleteSite(id) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "site not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

package store

import (
	"sync"

	"github.com/google/uuid"
)

// Device represents a Netbox device
type Device struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Site     int    `json:"site"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Serial   string `json:"serial,omitempty"`
	AssetTag string `json:"asset_tag,omitempty"`
}

// Interface represents a network interface on a device
type Interface struct {
	ID       int    `json:"id"`
	Device   int    `json:"device"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Enabled  bool   `json:"enabled"`
	MTU      int    `json:"mtu,omitempty"`
	Mode     string `json:"mode,omitempty"`
	MACAddr  string `json:"mac_address,omitempty"`
}

// Site represents a Netbox site
type Site struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Region      string `json:"region,omitempty"`
	Facility    string `json:"facility,omitempty"`
	ASN         int    `json:"asn,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
	Description string `json:"description,omitempty"`
}

// Store manages in-memory data for Netbox mock
type Store struct {
	mu         sync.RWMutex
	devices    map[int]*Device
	interfaces map[int]*Interface
	sites      map[int]*Site
	nextID     int
}

// NewStore creates a new store instance
func NewStore() *Store {
	return &Store{
		devices:    make(map[int]*Device),
		interfaces: make(map[int]*Interface),
		sites:      make(map[int]*Site),
		nextID:     1,
	}
}

// Device operations

// CreateDevice adds a new device to the store
func (s *Store) CreateDevice(d *Device) *Device {
	s.mu.Lock()
	defer s.mu.Unlock()
	d.ID = s.nextID
	s.nextID++
	s.devices[d.ID] = d
	return d
}

// GetDevice retrieves a device by ID
func (s *Store) GetDevice(id int) *Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.devices[id]
}

// ListDevices returns all devices
func (s *Store) ListDevices() []*Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Device, 0, len(s.devices))
	for _, d := range s.devices {
		result = append(result, d)
	}
	return result
}

// UpdateDevice updates an existing device
func (s *Store) UpdateDevice(id int, d *Device) *Device {
	s.mu.Lock()
	defer s.mu.Unlock()
	if device, exists := s.devices[id]; exists {
		if d.Name != "" {
			device.Name = d.Name
		}
		if d.Type != "" {
			device.Type = d.Type
		}
		if d.Status != "" {
			device.Status = d.Status
		}
		if d.Serial != "" {
			device.Serial = d.Serial
		}
		if d.AssetTag != "" {
			device.AssetTag = d.AssetTag
		}
		return device
	}
	return nil
}

// DeleteDevice removes a device from the store
func (s *Store) DeleteDevice(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[id]; exists {
		delete(s.devices, id)
		return true
	}
	return false
}

// Interface operations

// CreateInterface adds a new interface to the store
func (s *Store) CreateInterface(i *Interface) *Interface {
	s.mu.Lock()
	defer s.mu.Unlock()
	i.ID = s.nextID
	s.nextID++
	s.interfaces[i.ID] = i
	return i
}

// GetInterface retrieves an interface by ID
func (s *Store) GetInterface(id int) *Interface {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.interfaces[id]
}

// ListInterfaces returns all interfaces
func (s *Store) ListInterfaces() []*Interface {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Interface, 0, len(s.interfaces))
	for _, i := range s.interfaces {
		result = append(result, i)
	}
	return result
}

// DeleteInterface removes an interface from the store
func (s *Store) DeleteInterface(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.interfaces[id]; exists {
		delete(s.interfaces, id)
		return true
	}
	return false
}

// Site operations

// CreateSite adds a new site to the store
func (s *Store) CreateSite(site *Site) *Site {
	s.mu.Lock()
	defer s.mu.Unlock()
	site.ID = s.nextID
	s.nextID++
	if site.Slug == "" {
		site.Slug = uuid.New().String()
	}
	s.sites[site.ID] = site
	return site
}

// GetSite retrieves a site by ID
func (s *Store) GetSite(id int) *Site {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sites[id]
}

// ListSites returns all sites
func (s *Store) ListSites() []*Site {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Site, 0, len(s.sites))
	for _, s := range s.sites {
		result = append(result, s)
	}
	return result
}

// DeleteSite removes a site from the store
func (s *Store) DeleteSite(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.sites[id]; exists {
		delete(s.sites, id)
		return true
	}
	return false
}

package cause_test

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world example: IoT Device Management System
type IoTDeviceConfig struct {
	DeviceID        string         `json:"device_id"`
	DeviceName      string         `json:"device_name"`
	DeviceType      string         `json:"device_type"`
	Model           string         `json:"model"`
	Manufacturer    string         `json:"manufacturer"`
	FirmwareVersion string         `json:"firmware_version"`
	NetworkConfig   NetworkConfig  `json:"network_config"`
	SensorConfig    []SensorConfig `json:"sensor_config"`
	SecurityConfig  SecurityConfig `json:"security_config"`
	PowerConfig     PowerConfig    `json:"power_config"`
	LocationInfo    LocationInfo   `json:"location_info"`
	Metadata        DeviceMetadata `json:"metadata"`
	Status          string         `json:"status"`
	LastSeen        time.Time      `json:"last_seen"`
	RegisteredAt    time.Time      `json:"registered_at"`
}

type NetworkConfig struct {
	WiFiSSID       string   `json:"wifi_ssid,omitempty"`
	WiFiPassword   string   `json:"wifi_password,omitempty"`
	IPAddress      string   `json:"ip_address,omitempty"`
	SubnetMask     string   `json:"subnet_mask,omitempty"`
	Gateway        string   `json:"gateway,omitempty"`
	DNSServers     []string `json:"dns_servers,omitempty"`
	MACAddress     string   `json:"mac_address"`
	NetworkMode    string   `json:"network_mode"`
	SignalStrength int      `json:"signal_strength,omitempty"`
	Bandwidth      int      `json:"bandwidth_mbps,omitempty"`
}

type SensorConfig struct {
	SensorID    string            `json:"sensor_id"`
	SensorType  string            `json:"sensor_type"`
	SensorName  string            `json:"sensor_name"`
	Unit        string            `json:"unit"`
	SampleRate  int               `json:"sample_rate_seconds"`
	Threshold   SensorThreshold   `json:"threshold"`
	Calibration SensorCalibration `json:"calibration"`
	Enabled     bool              `json:"enabled"`
	Location    string            `json:"location,omitempty"`
}

type SensorThreshold struct {
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	CriticalMin  float64 `json:"critical_min"`
	CriticalMax  float64 `json:"critical_max"`
	AlertEnabled bool    `json:"alert_enabled"`
}

type SensorCalibration struct {
	Offset          float64   `json:"offset"`
	Scale           float64   `json:"scale"`
	LastCalibrated  time.Time `json:"last_calibrated"`
	CalibratedBy    string    `json:"calibrated_by"`
	CalibrationCert string    `json:"calibration_cert,omitempty"`
}

type SecurityConfig struct {
	EncryptionEnabled bool      `json:"encryption_enabled"`
	EncryptionType    string    `json:"encryption_type"`
	CertificateID     string    `json:"certificate_id,omitempty"`
	AuthMethod        string    `json:"auth_method"`
	APIKeys           []APIKey  `json:"api_keys,omitempty"`
	FirewallEnabled   bool      `json:"firewall_enabled"`
	VPNEnabled        bool      `json:"vpn_enabled"`
	LastSecurityScan  time.Time `json:"last_security_scan,omitempty"`
}

type APIKey struct {
	KeyID       string    `json:"key_id"`
	KeyName     string    `json:"key_name"`
	Permissions []string  `json:"permissions"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
}

type PowerConfig struct {
	PowerSource      string  `json:"power_source"`
	BatteryLevel     int     `json:"battery_level,omitempty"`
	VoltageRating    float64 `json:"voltage_rating"`
	PowerConsumption float64 `json:"power_consumption_watts"`
	SleepMode        bool    `json:"sleep_mode_enabled"`
	SleepSchedule    string  `json:"sleep_schedule,omitempty"`
	LowPowerAlert    bool    `json:"low_power_alert"`
	PowerThreshold   int     `json:"power_threshold_percent,omitempty"`
}

type LocationInfo struct {
	Building    string  `json:"building,omitempty"`
	Room        string  `json:"room,omitempty"`
	Floor       int     `json:"floor,omitempty"`
	Zone        string  `json:"zone,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	Altitude    float64 `json:"altitude,omitempty"`
	Address     string  `json:"address,omitempty"`
	Description string  `json:"description,omitempty"`
}

type DeviceMetadata struct {
	Tags              []string          `json:"tags,omitempty"`
	Environment       string            `json:"environment"`
	DeploymentType    string            `json:"deployment_type"`
	MaintenanceWindow string            `json:"maintenance_window,omitempty"`
	SupportContact    string            `json:"support_contact,omitempty"`
	PurchaseDate      time.Time         `json:"purchase_date,omitempty"`
	WarrantyExpiry    time.Time         `json:"warranty_expiry,omitempty"`
	CustomFields      map[string]string `json:"custom_fields,omitempty"`
}

func (idc *IoTDeviceConfig) Validate() error {
	return cause.Map{
		"device_id": cause.Required(idc.DeviceID).
			Select(map[string]bool{
				"invalid device ID format": !isValidDeviceID(idc.DeviceID),
				"device ID too long":       len(idc.DeviceID) > 50,
			}),

		"device_name": cause.Required(idc.DeviceName).
			Select(map[string]bool{
				"device name too short":                   len(idc.DeviceName) < 3,
				"device name too long":                    len(idc.DeviceName) > 100,
				"device name contains invalid characters": containsSpecialChars(idc.DeviceName),
			}),

		"device_type": cause.Required(idc.DeviceType).
			Select(map[string]bool{
				"invalid device type": !isValidDeviceType(idc.DeviceType),
			}),

		"model": cause.Required(idc.Model).
			Select(map[string]bool{
				"model name too short": len(idc.Model) < 2,
				"model name too long":  len(idc.Model) > 50,
			}),

		"manufacturer": cause.Required(idc.Manufacturer).
			Select(map[string]bool{
				"manufacturer name too short": len(idc.Manufacturer) < 2,
				"manufacturer name too long":  len(idc.Manufacturer) > 50,
			}),

		"firmware_version": cause.Required(idc.FirmwareVersion).
			Select(map[string]bool{
				"invalid firmware version format": !isValidVersion(idc.FirmwareVersion),
			}),

		"network_config":  cause.Required(idc.NetworkConfig).Err(),
		"sensor_config":   cause.Optional(idc.SensorConfig).Err(),
		"security_config": cause.Required(idc.SecurityConfig).Err(),
		"power_config":    cause.Required(idc.PowerConfig).Err(),
		"location_info":   cause.Optional(idc.LocationInfo).Err(),
		"metadata":        cause.Required(idc.Metadata).Err(),

		"status": cause.Required(idc.Status).
			Select(map[string]bool{
				"invalid device status": !isValidDeviceStatus(idc.Status),
			}),

		"last_seen": cause.Required(idc.LastSeen).
			Select(map[string]bool{
				"last seen cannot be in the future": idc.LastSeen.After(time.Now()),
				"last seen date too old":            idc.LastSeen.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),

		"registered_at": cause.Required(idc.RegisteredAt).
			Select(map[string]bool{
				"registration date cannot be in the future": idc.RegisteredAt.After(time.Now()),
				"registration date too old":                 idc.RegisteredAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			}),
	}.Err()
}

func (nc *NetworkConfig) Validate() error {
	return cause.Map{
		"wifi_ssid": cause.Optional(nc.WiFiSSID).
			Select(map[string]bool{
				"WiFi SSID too long":                    len(nc.WiFiSSID) > 32,
				"WiFi SSID contains invalid characters": containsInvalidSSIDChars(nc.WiFiSSID),
			}),

		"wifi_password": cause.Optional(nc.WiFiPassword).
			Select(map[string]bool{
				"WiFi password too short": len(nc.WiFiPassword) < 8 && len(nc.WiFiPassword) > 0,
				"WiFi password too long":  len(nc.WiFiPassword) > 63,
			}),

		"ip_address": cause.Optional(nc.IPAddress).
			Select(map[string]bool{
				"invalid IP address format": !isValidIPAddress(nc.IPAddress),
			}),

		"subnet_mask": cause.Optional(nc.SubnetMask).
			Select(map[string]bool{
				"invalid subnet mask format": !isValidIPAddress(nc.SubnetMask),
			}),

		"gateway": cause.Optional(nc.Gateway).
			Select(map[string]bool{
				"invalid gateway IP address": !isValidIPAddress(nc.Gateway),
			}),

		"dns_servers": cause.Optional(nc.DNSServers).
			Select(map[string]bool{
				"too many DNS servers":         len(nc.DNSServers) > 5,
				"invalid DNS server addresses": !areValidDNSServers(nc.DNSServers),
			}),

		"mac_address": cause.Required(nc.MACAddress).
			Select(map[string]bool{
				"invalid MAC address format": !isValidMACAddress(nc.MACAddress),
			}),

		"network_mode": cause.Required(nc.NetworkMode).
			Select(map[string]bool{
				"invalid network mode": !isValidNetworkMode(nc.NetworkMode),
			}),

		"signal_strength": cause.Optional(nc.SignalStrength).
			Select(map[string]bool{
				"signal strength outside valid range": nc.SignalStrength < -100 || nc.SignalStrength > 0,
			}),

		"bandwidth_mbps": cause.Optional(nc.Bandwidth).
			Select(map[string]bool{
				"bandwidth cannot be negative":       nc.Bandwidth < 0,
				"bandwidth exceeds reasonable limit": nc.Bandwidth > 10000,
			}),
	}.Err()
}

func (sc *SensorConfig) Validate() error {
	return cause.Map{
		"sensor_id": cause.Required(sc.SensorID).
			Select(map[string]bool{
				"invalid sensor ID format": !isValidSensorID(sc.SensorID),
				"sensor ID too long":       len(sc.SensorID) > 30,
			}),

		"sensor_type": cause.Required(sc.SensorType).
			Select(map[string]bool{
				"invalid sensor type": !isValidSensorType(sc.SensorType),
			}),

		"sensor_name": cause.Required(sc.SensorName).
			Select(map[string]bool{
				"sensor name too short": len(sc.SensorName) < 2,
				"sensor name too long":  len(sc.SensorName) > 50,
			}),

		"unit": cause.Required(sc.Unit).
			Select(map[string]bool{
				"invalid measurement unit": !isValidUnit(sc.Unit),
			}),

		"sample_rate_seconds": cause.Required(sc.SampleRate).
			Select(map[string]bool{
				"sample rate must be at least 1 second": sc.SampleRate < 1,
				"sample rate cannot exceed 24 hours":    sc.SampleRate > 86400,
			}),

		"threshold":   cause.Required(sc.Threshold).Err(),
		"calibration": cause.Required(sc.Calibration).Err(),

		"location": cause.Optional(sc.Location).
			Select(map[string]bool{
				"sensor location description too long": len(sc.Location) > 100,
			}),
	}.Err()
}

func (st *SensorThreshold) Validate() error {
	return cause.Map{
		"min_value": cause.Required(st.MinValue).
			When(st.MinValue >= st.MaxValue, "minimum value must be less than maximum value").Err(),

		"max_value": cause.Required(st.MaxValue).Err(),

		"critical_min": cause.Required(st.CriticalMin).
			When(st.CriticalMin > st.MinValue, "critical minimum must be less than or equal to minimum value").Err(),

		"critical_max": cause.Required(st.CriticalMax).
			When(st.CriticalMax < st.MaxValue, "critical maximum must be greater than or equal to maximum value").Err(),
	}.Err()
}

func (sc *SensorCalibration) Validate() error {
	return cause.Map{
		"scale": cause.Required(sc.Scale).
			When(sc.Scale <= 0, "calibration scale must be positive").Err(),

		"last_calibrated": cause.Required(sc.LastCalibrated).
			When(sc.LastCalibrated.After(time.Now()), "calibration date cannot be in the future").
			When(sc.LastCalibrated.Before(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)), "calibration date too old").Err(),

		"calibrated_by": cause.Required(sc.CalibratedBy).
			When(len(sc.CalibratedBy) < 2, "calibrated by field too short").
			When(len(sc.CalibratedBy) > 50, "calibrated by field too long").Err(),

		"calibration_cert": cause.Optional(sc.CalibrationCert).
			When(len(sc.CalibrationCert) > 100, "calibration certificate ID too long").Err(),
	}.Err()
}

func (sec *SecurityConfig) Validate() error {
	return cause.Map{
		"encryption_type": cause.RequiredMessage(sec.EncryptionType, "encryption type required when encryption enabled").
			When(sec.EncryptionEnabled && !isValidEncryptionType(sec.EncryptionType), "invalid encryption type").Err(),

		"certificate_id": cause.Optional(sec.CertificateID).
			When(sec.EncryptionEnabled && len(sec.CertificateID) == 0, "certificate ID required when encryption enabled").
			When(len(sec.CertificateID) > 50, "certificate ID too long").Err(),

		"auth_method": cause.Required(sec.AuthMethod).
			When(!isValidAuthMethod(sec.AuthMethod), "invalid authentication method").Err(),

		"api_keys": cause.Optional(sec.APIKeys).
			When(len(sec.APIKeys) > 10, "too many API keys").Err(),

		"last_security_scan": cause.Optional(sec.LastSecurityScan).
			When(!sec.LastSecurityScan.IsZero() && sec.LastSecurityScan.After(time.Now()), "security scan date cannot be in the future").
			When(!sec.LastSecurityScan.IsZero() && sec.LastSecurityScan.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)), "security scan date too old").Err(),
	}.Err()
}

func (ak *APIKey) Validate() error {
	return cause.Map{
		"key_id": cause.Required(ak.KeyID).
			When(!isValidAPIKeyID(ak.KeyID), "invalid API key ID format").
			When(len(ak.KeyID) > 50, "API key ID too long").Err(),

		"key_name": cause.Required(ak.KeyName).
			When(len(ak.KeyName) < 3, "API key name too short").
			When(len(ak.KeyName) > 50, "API key name too long").Err(),

		"permissions": cause.Required(ak.Permissions).
			When(len(ak.Permissions) == 0, "at least one permission required").
			When(len(ak.Permissions) > 20, "too many permissions").
			When(!areValidPermissions(ak.Permissions), "invalid permissions").Err(),

		"expires_at": cause.Required(ak.ExpiresAt).
			When(ak.ExpiresAt.Before(ak.CreatedAt), "expiry date cannot be before creation date").
			When(ak.ExpiresAt.Before(time.Now()), "API key has expired").
			When(ak.ExpiresAt.After(time.Now().AddDate(5, 0, 0)), "expiry date too far in future").Err(),

		"created_at": cause.Required(ak.CreatedAt).
			When(ak.CreatedAt.After(time.Now()), "creation date cannot be in the future").
			When(ak.CreatedAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)), "creation date too old").Err(),
	}.Err()
}

func (pc *PowerConfig) Validate() error {
	return cause.Map{
		"power_source": cause.Required(pc.PowerSource).
			When(!isValidPowerSource(pc.PowerSource), "invalid power source").Err(),

		"battery_level": cause.Optional(pc.BatteryLevel).
			When(pc.PowerSource == "Battery" && pc.BatteryLevel < 0, "battery level cannot be negative").
			When(pc.PowerSource == "Battery" && pc.BatteryLevel > 100, "battery level cannot exceed 100%").Err(),

		"voltage_rating": cause.Required(pc.VoltageRating).
			When(pc.VoltageRating <= 0, "voltage rating must be positive").
			When(pc.VoltageRating > 1000, "voltage rating exceeds safe limits").Err(),

		"power_consumption_watts": cause.Required(pc.PowerConsumption).
			When(pc.PowerConsumption < 0, "power consumption cannot be negative").
			When(pc.PowerConsumption > 10000, "power consumption exceeds reasonable limits").Err(),

		"sleep_schedule": cause.Optional(pc.SleepSchedule).
			When(pc.SleepMode && len(pc.SleepSchedule) == 0, "sleep schedule required when sleep mode enabled").
			When(!isValidCronSchedule(pc.SleepSchedule), "invalid sleep schedule format").Err(),

		"power_threshold_percent": cause.Optional(pc.PowerThreshold).
			When(pc.LowPowerAlert && pc.PowerThreshold <= 0, "power threshold required when low power alert enabled").
			When(pc.PowerThreshold < 0 || pc.PowerThreshold > 100, "power threshold must be between 0 and 100").Err(),
	}.Err()
}

func (li *LocationInfo) Validate() error {
	return cause.Map{
		"building": cause.Optional(li.Building).
			When(len(li.Building) > 100, "building name too long").Err(),

		"room": cause.Optional(li.Room).
			When(len(li.Room) > 50, "room identifier too long").Err(),

		"floor": cause.Optional(li.Floor).
			When(li.Floor < -10 || li.Floor > 200, "floor number outside reasonable range").Err(),

		"zone": cause.Optional(li.Zone).
			When(len(li.Zone) > 50, "zone identifier too long").Err(),

		"latitude": cause.Optional(li.Latitude).
			When(li.Latitude < -90 || li.Latitude > 90, "latitude outside valid range").Err(),

		"longitude": cause.Optional(li.Longitude).
			When(li.Longitude < -180 || li.Longitude > 180, "longitude outside valid range").Err(),

		"altitude": cause.Optional(li.Altitude).
			When(li.Altitude < -1000 || li.Altitude > 10000, "altitude outside reasonable range").Err(),

		"address": cause.Optional(li.Address).
			When(len(li.Address) > 200, "address too long").Err(),

		"description": cause.Optional(li.Description).
			When(len(li.Description) > 500, "location description too long").Err(),
	}.Err()
}

func (dm *DeviceMetadata) Validate() error {
	return cause.Map{
		"tags": cause.Optional(dm.Tags).
			When(len(dm.Tags) > 20, "too many tags").
			When(!areValidTags(dm.Tags), "invalid tag format").Err(),

		"environment": cause.Required(dm.Environment).
			When(!isValidEnvironment(dm.Environment), "invalid environment").Err(),

		"deployment_type": cause.Required(dm.DeploymentType).
			When(!isValidDeploymentType(dm.DeploymentType), "invalid deployment type").Err(),

		"maintenance_window": cause.Optional(dm.MaintenanceWindow).
			When(!isValidCronSchedule(dm.MaintenanceWindow), "invalid maintenance window format").Err(),

		"support_contact": cause.Optional(dm.SupportContact).
			When(!isValidEmail(dm.SupportContact), "invalid support contact email").Err(),

		"purchase_date": cause.Optional(dm.PurchaseDate).
			When(!dm.PurchaseDate.IsZero() && dm.PurchaseDate.After(time.Now()), "purchase date cannot be in the future").
			When(!dm.PurchaseDate.IsZero() && dm.PurchaseDate.Before(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)), "purchase date too old").Err(),

		"warranty_expiry": cause.Optional(dm.WarrantyExpiry).
			When(!dm.WarrantyExpiry.IsZero() && !dm.PurchaseDate.IsZero() && dm.WarrantyExpiry.Before(dm.PurchaseDate), "warranty expiry cannot be before purchase date").
			When(!dm.WarrantyExpiry.IsZero() && dm.WarrantyExpiry.After(time.Now().AddDate(10, 0, 0)), "warranty period too long").Err(),

		"custom_fields": cause.Optional(dm.CustomFields).
			When(len(dm.CustomFields) > 50, "too many custom fields").Err(),
	}.Err()
}

// IoT-specific validation helper functions
func isValidDeviceID(id string) bool {
	matched, _ := regexp.MatchString(`^[A-Z0-9]{8,32}$`, id)
	return matched
}

func containsSpecialChars(s string) bool {
	// Device names should only contain alphanumeric characters, spaces, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9 _-]+$`, s)
	return !matched
}

func isValidVersion(version string) bool {
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$`, version)
	return matched
}

func isValidDeviceStatus(status string) bool {
	validStatuses := []string{"Online", "Offline", "Maintenance", "Error", "Initializing", "Updating"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func containsInvalidSSIDChars(ssid string) bool {
	// SSID cannot contain certain control characters
	for _, char := range ssid {
		if char < 32 || char > 126 {
			return true
		}
	}
	return false
}

func isValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func areValidDNSServers(servers []string) bool {
	for _, server := range servers {
		if !isValidIPAddress(server) {
			return false
		}
	}
	return true
}

func isValidMACAddress(mac string) bool {
	matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, mac)
	return matched
}

func isValidNetworkMode(mode string) bool {
	validModes := []string{"WiFi", "Ethernet", "Cellular", "LoRa", "Zigbee", "Bluetooth"}
	for _, valid := range validModes {
		if mode == valid {
			return true
		}
	}
	return false
}

func isValidSensorID(id string) bool {
	matched, _ := regexp.MatchString(`^S[0-9A-Z]{4,20}$`, id)
	return matched
}

func isValidSensorType(sensorType string) bool {
	validTypes := []string{
		"Temperature", "Humidity", "Pressure", "Motion", "Light",
		"Sound", "Gas", "Proximity", "Accelerometer", "Gyroscope",
		"GPS", "Current", "Voltage", "pH", "Flow",
	}
	for _, valid := range validTypes {
		if sensorType == valid {
			return true
		}
	}
	return false
}

func isValidUnit(unit string) bool {
	validUnits := []string{
		"°C", "°F", "K", "%", "Pa", "hPa", "bar", "psi",
		"m/s", "km/h", "mph", "lux", "dB", "ppm", "V", "A",
		"W", "Hz", "rpm", "m", "cm", "mm", "kg", "g",
	}
	for _, valid := range validUnits {
		if unit == valid {
			return true
		}
	}
	return false
}

func isValidEncryptionType(encType string) bool {
	validTypes := []string{"AES256", "AES128", "TLS1.2", "TLS1.3", "WPA2", "WPA3"}
	for _, valid := range validTypes {
		if encType == valid {
			return true
		}
	}
	return false
}

func isValidAuthMethod(method string) bool {
	validMethods := []string{"Certificate", "Key", "Username/Password", "OAuth2", "JWT", "HMAC"}
	for _, valid := range validMethods {
		if method == valid {
			return true
		}
	}
	return false
}

func isValidAPIKeyID(id string) bool {
	matched, _ := regexp.MatchString(`^ak_[0-9a-zA-Z]{16,32}$`, id)
	return matched
}

func areValidPermissions(permissions []string) bool {
	validPerms := []string{"read", "write", "delete", "admin", "config", "monitor"}
	for _, perm := range permissions {
		found := false
		for _, valid := range validPerms {
			if perm == valid {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func isValidPowerSource(source string) bool {
	validSources := []string{"AC", "Battery", "Solar", "PoE", "USB", "DC"}
	for _, valid := range validSources {
		if source == valid {
			return true
		}
	}
	return false
}

func isValidCronSchedule(schedule string) bool {
	if schedule == "" {
		return true // Optional field
	}
	// Simplified cron validation - real implementation would be more thorough
	parts := strings.Fields(schedule)
	return len(parts) == 5 || len(parts) == 6
}

func areValidTags(tags []string) bool {
	for _, tag := range tags {
		if len(tag) == 0 || len(tag) > 50 {
			return false
		}
		// Tags should be alphanumeric with hyphens and underscores
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_]+$`, tag)
		if !matched {
			return false
		}
	}
	return true
}

func isValidEnvironment(env string) bool {
	validEnvs := []string{"Production", "Staging", "Development", "Testing", "Lab"}
	for _, valid := range validEnvs {
		if env == valid {
			return true
		}
	}
	return false
}

func isValidDeploymentType(depType string) bool {
	validTypes := []string{"Indoor", "Outdoor", "Industrial", "Marine", "Automotive", "Aerospace"}
	for _, valid := range validTypes {
		if depType == valid {
			return true
		}
	}
	return false
}

// Example test function
func ExampleIoTDeviceConfig_Validate() {
	// Valid IoT device configuration
	validDevice := &IoTDeviceConfig{
		DeviceID:        "DEV12345ABCD",
		DeviceName:      "Temperature Sensor 01",
		DeviceType:      "Sensor",
		Model:           "TS-2000",
		Manufacturer:    "SensorTech",
		FirmwareVersion: "2.1.0",
		NetworkConfig: NetworkConfig{
			WiFiSSID:    "IoT_Network",
			MACAddress:  "00:1B:63:84:45:E6",
			NetworkMode: "WiFi",
		},
		SecurityConfig: SecurityConfig{
			EncryptionEnabled: true,
			EncryptionType:    "AES256",
			AuthMethod:        "Certificate",
		},
		PowerConfig: PowerConfig{
			PowerSource:      "Battery",
			BatteryLevel:     85,
			VoltageRating:    3.3,
			PowerConsumption: 0.5,
		},
		Metadata: DeviceMetadata{
			Environment:    "Production",
			DeploymentType: "Indoor",
		},
		Status:       "Online",
		LastSeen:     time.Now().Add(-5 * time.Minute),
		RegisteredAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	err := validDevice.Validate()
	fmt.Printf("Valid IoT device error: %v\n", err)

	// Invalid IoT device configuration
	invalidDevice := &IoTDeviceConfig{
		DeviceID:        "INVALID_ID!",
		DeviceName:      "AB", // Too short
		DeviceType:      "Invalid Type",
		FirmwareVersion: "invalid.version",
		NetworkConfig: NetworkConfig{
			MACAddress:  "invalid-mac",
			NetworkMode: "Invalid Mode",
		},
		SecurityConfig: SecurityConfig{
			EncryptionEnabled: true,
			EncryptionType:    "", // Required when encryption enabled
			AuthMethod:        "Invalid Method",
		},
		PowerConfig: PowerConfig{
			PowerSource:      "Invalid Source",
			VoltageRating:    -5.0, // Negative voltage
			PowerConsumption: -1.0, // Negative consumption
		},
		Metadata: DeviceMetadata{
			Environment:    "Invalid Env",
			DeploymentType: "Invalid Type",
		},
		Status:       "Invalid Status",
		LastSeen:     time.Now().Add(time.Hour), // Future time
		RegisteredAt: time.Now().Add(time.Hour), // Future time
	}

	err = invalidDevice.Validate()
	if err != nil {
		fmt.Printf("Invalid IoT device has errors: %v\n", err != nil)
	}

	// Output:
	// Valid IoT device error: <nil>
	// Invalid IoT device has errors: true
}

func isValidDeviceType(deviceType string) bool {
	validTypes := []string{"Sensor", "Actuator", "Controller", "Gateway", "Monitor", "Alarm", "Display"}
	for _, valid := range validTypes {
		if deviceType == valid {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	if email == "" {
		return true // Optional field
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

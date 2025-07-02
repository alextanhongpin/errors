package cause_test

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world example: REST API Request Validation
type CreateUserAPIRequest struct {
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Profile  UserProfile       `json:"profile"`
	Settings UserSettings      `json:"settings"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type UserProfile struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Location    string    `json:"location,omitempty"`
	Website     string    `json:"website,omitempty"`
}

type UserSettings struct {
	Theme         string               `json:"theme"`
	Language      string               `json:"language"`
	Timezone      string               `json:"timezone"`
	Notifications NotificationSettings `json:"notifications"`
	Privacy       PrivacySettings      `json:"privacy"`
}

type NotificationSettings struct {
	Email     bool `json:"email"`
	Push      bool `json:"push"`
	SMS       bool `json:"sms"`
	Marketing bool `json:"marketing"`
}

type PrivacySettings struct {
	ProfilePublic  bool `json:"profile_public"`
	ShowEmail      bool `json:"show_email"`
	ShowLastSeen   bool `json:"show_last_seen"`
	AllowDiscovery bool `json:"allow_discovery"`
}

func (req *CreateUserAPIRequest) Validate() error {
	return cause.Map{
		"username": cause.Required(req.Username).
			When(len(req.Username) < 3, "username must be at least 3 characters").
			When(len(req.Username) > 30, "username must not exceed 30 characters").
			When(!isValidUsername(req.Username), "username contains invalid characters").
			When(isReservedUsername(req.Username), "username is reserved"),

		"email": cause.Required(req.Email).
			When(!isValidEmail(req.Email), "invalid email format").
			When(isDisposableEmail(req.Email), "disposable email addresses not allowed").
			When(isDomainBlacklisted(req.Email), "email domain is blacklisted"),

		"password": cause.Required(req.Password).
			When(len(req.Password) < 8, "password must be at least 8 characters").
			When(len(req.Password) > 128, "password must not exceed 128 characters").
			When(!hasRequiredPasswordComplexity(req.Password), "password does not meet complexity requirements").
			When(isCommonPassword(req.Password), "password is too common"),

		"profile":  cause.Required(req.Profile),
		"settings": cause.Required(req.Settings),

		"metadata": cause.Optional(req.Metadata).
			When(len(req.Metadata) > 20, "too many metadata fields").
			When(hasInvalidMetadataKeys(req.Metadata), "metadata contains invalid keys"),
	}.Err()
}

func (profile *UserProfile) Validate() error {
	return cause.Map{
		"first_name": cause.Required(profile.FirstName).
			When(len(profile.FirstName) < 1, "first name is required").
			When(len(profile.FirstName) > 50, "first name too long").
			When(containsInvalidNameChars(profile.FirstName), "first name contains invalid characters"),

		"last_name": cause.Required(profile.LastName).
			When(len(profile.LastName) < 1, "last name is required").
			When(len(profile.LastName) > 50, "last name too long").
			When(containsInvalidNameChars(profile.LastName), "last name contains invalid characters"),

		"display_name": cause.Optional(profile.DisplayName).
			When(len(profile.DisplayName) > 100, "display name too long").
			When(containsProfanity(profile.DisplayName), "display name contains inappropriate content"),

		"bio": cause.Optional(profile.Bio).
			When(len(profile.Bio) > 500, "bio too long").
			When(containsProfanity(profile.Bio), "bio contains inappropriate content").
			When(hasExcessiveLinks(profile.Bio), "bio contains too many links"),

		"avatar": cause.Optional(profile.Avatar).
			When(!isValidImageURL(profile.Avatar), "invalid avatar image URL").
			When(!isSupportedImageFormat(profile.Avatar), "unsupported image format"),

		"date_of_birth": cause.Optional(profile.DateOfBirth).
			When(!profile.DateOfBirth.IsZero() && profile.DateOfBirth.After(time.Now()), "date of birth cannot be in the future").
			When(!profile.DateOfBirth.IsZero() && getAge(profile.DateOfBirth) < 13, "must be at least 13 years old").
			When(!profile.DateOfBirth.IsZero() && getAge(profile.DateOfBirth) > 120, "invalid date of birth"),

		"location": cause.Optional(profile.Location).
			When(len(profile.Location) > 100, "location too long"),

		"website": cause.Optional(profile.Website).
			When(!isValidURL(profile.Website), "invalid website URL").
			When(isMaliciousURL(profile.Website), "website URL appears malicious"),
	}.Err()
}

func (settings *UserSettings) Validate() error {
	return cause.Map{
		"theme": cause.Required(settings.Theme).
			When(!isValidTheme(settings.Theme), "invalid theme"),

		"language": cause.Required(settings.Language).
			When(!isValidLanguageCode(settings.Language), "invalid language code"),

		"timezone": cause.Required(settings.Timezone).
			When(!isValidTimezone(settings.Timezone), "invalid timezone"),

		"notifications": cause.Required(settings.Notifications),
		"privacy":       cause.Required(settings.Privacy),
	}.Err()
}

func (ns *NotificationSettings) Validate() error {
	// Notification settings are just booleans, so basic validation
	return nil
}

func (ps *PrivacySettings) Validate() error {
	// Privacy settings are just booleans, so basic validation
	return nil
}

// Real-world example: Search API Request
type SearchRequest struct {
	Query      string            `json:"query"`
	Filters    SearchFilters     `json:"filters,omitempty"`
	Sort       SortOptions       `json:"sort,omitempty"`
	Pagination PaginationOptions `json:"pagination,omitempty"`
	Facets     []string          `json:"facets,omitempty"`
	Highlight  HighlightOptions  `json:"highlight,omitempty"`
	DebugMode  bool              `json:"debug_mode,omitempty"`
}

type SearchFilters struct {
	Categories    []string          `json:"categories,omitempty"`
	PriceRange    *PriceRange       `json:"price_range,omitempty"`
	DateRange     *DateRange        `json:"date_range,omitempty"`
	Location      *LocationFilter   `json:"location,omitempty"`
	CustomFilters map[string]string `json:"custom_filters,omitempty"`
}

type PriceRange struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type LocationFilter struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Unit      string  `json:"unit"`
}

type SortOptions struct {
	Field     string      `json:"field"`
	Direction string      `json:"direction"`
	Secondary []SortField `json:"secondary,omitempty"`
}

type SortField struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type PaginationOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset,omitempty"`
}

type HighlightOptions struct {
	Enabled     bool     `json:"enabled"`
	Fields      []string `json:"fields,omitempty"`
	PreTag      string   `json:"pre_tag,omitempty"`
	PostTag     string   `json:"post_tag,omitempty"`
	MaxSnippets int      `json:"max_snippets,omitempty"`
}

func (req *SearchRequest) Validate() error {
	return cause.Map{
		"query": cause.Required(req.Query).
			When(len(req.Query) < 1, "search query cannot be empty").
			When(len(req.Query) > 500, "search query too long").
			When(containsMaliciousPatterns(req.Query), "search query contains potentially malicious patterns").
			When(isOnlyWhitespace(req.Query), "search query cannot be only whitespace"),

		"filters":    cause.Optional(req.Filters),
		"sort":       cause.Optional(req.Sort),
		"pagination": cause.Optional(req.Pagination),

		"facets": cause.Optional(req.Facets).
			When(len(req.Facets) > 50, "too many facets requested").
			When(hasInvalidFacets(req.Facets), "contains invalid facet names"),

		"highlight": cause.Optional(req.Highlight),
	}.Err()
}

func (filters *SearchFilters) Validate() error {
	return cause.Map{
		"categories": cause.Optional(filters.Categories).
			When(len(filters.Categories) > 20, "too many categories").
			When(hasInvalidCategories(filters.Categories), "contains invalid categories"),

		"price_range": cause.Optional(filters.PriceRange),
		"date_range":  cause.Optional(filters.DateRange),
		"location":    cause.Optional(filters.Location),

		"custom_filters": cause.Optional(filters.CustomFilters).
			When(len(filters.CustomFilters) > 10, "too many custom filters").
			When(hasInvalidCustomFilterKeys(filters.CustomFilters), "contains invalid filter keys"),
	}.Err()
}

func (pr *PriceRange) Validate() error {
	return cause.Map{
		"min": cause.Required(pr.Min).
			When(pr.Min < 0, "minimum price cannot be negative"),

		"max": cause.Required(pr.Max).
			When(pr.Max < 0, "maximum price cannot be negative").
			When(pr.Max <= pr.Min, "maximum price must be greater than minimum price").
			When(pr.Max > 1000000, "maximum price too high"),

		"currency": cause.Required(pr.Currency).
			When(!isValidCurrencyCode(pr.Currency), "invalid currency code"),
	}.Err()
}

func (dr *DateRange) Validate() error {
	return cause.Map{
		"start": cause.Required(dr.Start).
			When(dr.Start.After(time.Now().AddDate(10, 0, 0)), "start date too far in future"),

		"end": cause.Required(dr.End).
			When(dr.End.Before(dr.Start), "end date must be after start date").
			When(dr.End.After(time.Now().AddDate(10, 0, 0)), "end date too far in future").
			When(dr.End.Sub(dr.Start) > 365*24*time.Hour, "date range too large (max 1 year)"),
	}.Err()
}

func (lf *LocationFilter) Validate() error {
	return cause.Map{
		"latitude": cause.Required(lf.Latitude).
			When(lf.Latitude < -90 || lf.Latitude > 90, "latitude must be between -90 and 90"),

		"longitude": cause.Required(lf.Longitude).
			When(lf.Longitude < -180 || lf.Longitude > 180, "longitude must be between -180 and 180"),

		"radius": cause.Required(lf.Radius).
			When(lf.Radius <= 0, "radius must be positive").
			When(lf.Radius > 1000, "radius too large (max 1000)"),

		"unit": cause.Required(lf.Unit).
			When(!isValidDistanceUnit(lf.Unit), "invalid distance unit"),
	}.Err()
}

func (so *SortOptions) Validate() error {
	return cause.Map{
		"field": cause.Required(so.Field).
			When(!isValidSortField(so.Field), "invalid sort field"),

		"direction": cause.Required(so.Direction).
			When(!isValidSortDirection(so.Direction), "invalid sort direction"),

		"secondary": cause.Optional(so.Secondary).
			When(len(so.Secondary) > 5, "too many secondary sort fields"),
	}.Err()
}

func (sf *SortField) Validate() error {
	return cause.Map{
		"field": cause.Required(sf.Field).
			When(!isValidSortField(sf.Field), "invalid sort field"),

		"direction": cause.Required(sf.Direction).
			When(!isValidSortDirection(sf.Direction), "invalid sort direction"),
	}.Err()
}

func (po *PaginationOptions) Validate() error {
	return cause.Map{
		"page": cause.Required(po.Page).
			When(po.Page < 1, "page must be at least 1").
			When(po.Page > 10000, "page number too high"),

		"page_size": cause.Required(po.PageSize).
			When(po.PageSize < 1, "page size must be at least 1").
			When(po.PageSize > 1000, "page size too large (max 1000)"),

		"offset": cause.Optional(po.Offset).
			When(po.Offset < 0, "offset cannot be negative").
			When(po.Offset > 100000, "offset too large"),
	}.Err()
}

func (ho *HighlightOptions) Validate() error {
	if !ho.Enabled {
		return nil
	}

	return cause.Map{
		"fields": cause.Optional(ho.Fields).
			When(len(ho.Fields) > 20, "too many highlight fields").
			When(hasInvalidHighlightFields(ho.Fields), "contains invalid highlight fields"),

		"pre_tag": cause.Optional(ho.PreTag).
			When(len(ho.PreTag) > 50, "pre tag too long"),

		"post_tag": cause.Optional(ho.PostTag).
			When(len(ho.PostTag) > 50, "post tag too long"),

		"max_snippets": cause.Optional(ho.MaxSnippets).
			When(ho.MaxSnippets < 0, "max snippets cannot be negative").
			When(ho.MaxSnippets > 100, "max snippets too high"),
	}.Err()
}

// Example functions demonstrating the validation results
func ExampleCreateUserAPIRequest_validation() {
	// Valid user creation request
	validRequest := &CreateUserAPIRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
		Profile: UserProfile{
			FirstName:   "John",
			LastName:    "Doe",
			DisplayName: "Johnny",
			Bio:         "Software developer passionate about clean code",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Website:     "https://johndoe.dev",
		},
		Settings: UserSettings{
			Theme:    "dark",
			Language: "en",
			Timezone: "UTC",
			Notifications: NotificationSettings{
				Email:     true,
				Push:      true,
				SMS:       false,
				Marketing: false,
			},
			Privacy: PrivacySettings{
				ProfilePublic:  true,
				ShowEmail:      false,
				ShowLastSeen:   true,
				AllowDiscovery: true,
			},
		},
	}

	fmt.Println("Valid user creation request:")
	if err := validRequest.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Invalid request
	invalidRequest := &CreateUserAPIRequest{
		Username: "ab", // Too short
		Email:    "invalid-email",
		Password: "weak", // Too weak
		Profile: UserProfile{
			FirstName: "",                       // Empty required field
			LastName:  "123Invalid",             // Contains numbers
			Bio:       strings.Repeat("x", 600), // Too long
		},
		Settings: UserSettings{
			Theme:         "invalid-theme",
			Language:      "invalid-lang",
			Timezone:      "invalid-tz",
			Notifications: NotificationSettings{},
			Privacy:       PrivacySettings{},
		},
	}

	fmt.Println("\nInvalid user creation request:")
	if err := invalidRequest.Validate(); err != nil {
		if validationErr, ok := err.(interface{ Map() map[string]any }); ok {
			fieldErrors := validationErr.Map()
			for field, fieldErr := range fieldErrors {
				fmt.Printf("  %s: %v\n", field, fieldErr)
			}
		}
	}

	// Note: Output field order may vary due to Go map iteration order
	// Output:
	// Valid user creation request:
	// Validation passed!
}

func ExampleSearchRequest_validation() {
	// Valid search request
	validSearch := &SearchRequest{
		Query: "laptop computers",
		Filters: SearchFilters{
			Categories: []string{"electronics", "computers"},
			PriceRange: &PriceRange{
				Min:      100.00,
				Max:      2000.00,
				Currency: "USD",
			},
			DateRange: &DateRange{
				Start: time.Now().AddDate(0, -1, 0),
				End:   time.Now(),
			},
		},
		Sort: SortOptions{
			Field:     "price",
			Direction: "asc",
		},
		Pagination: PaginationOptions{
			Page:     1,
			PageSize: 20,
		},
		Facets: []string{"brand", "category", "price_range"},
		Highlight: HighlightOptions{
			Enabled:     true,
			Fields:      []string{"title", "description"},
			PreTag:      "<em>",
			PostTag:     "</em>",
			MaxSnippets: 3,
		},
	}

	fmt.Println("Valid search request:")
	if err := validSearch.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Output:
	// Valid search request:
	// Validation passed!
}

// Helper functions for API validation
func isValidUsername(username string) bool {
	// Username can contain letters, numbers, underscores, and hyphens
	pattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

func isReservedUsername(username string) bool {
	reserved := []string{"admin", "root", "api", "www", "support", "help", "about", "contact"}
	usernameLower := strings.ToLower(username)
	for _, res := range reserved {
		if usernameLower == res {
			return true
		}
	}
	return false
}

func isDomainBlacklisted(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	blacklistedDomains := []string{"malicious.com", "spam.com", "phishing.net"}
	domain := strings.ToLower(parts[1])
	for _, blacklisted := range blacklistedDomains {
		if domain == blacklisted {
			return true
		}
	}
	return false
}

func hasRequiredPasswordComplexity(password string) bool {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func isCommonPassword(password string) bool {
	commonPasswords := []string{"password", "123456", "password123", "admin", "qwerty"}
	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if passwordLower == common {
			return true
		}
	}
	return false
}

func hasInvalidMetadataKeys(metadata map[string]string) bool {
	invalidKeys := []string{"password", "secret", "private", "internal"}
	for key := range metadata {
		keyLower := strings.ToLower(key)
		for _, invalid := range invalidKeys {
			if keyLower == invalid {
				return true
			}
		}
	}
	return false
}

func containsInvalidNameChars(name string) bool {
	// Names should only contain letters, spaces, hyphens, and apostrophes
	pattern := `^[a-zA-Z\s\-']+$`
	matched, _ := regexp.MatchString(pattern, name)
	return !matched
}

func containsProfanity(text string) bool {
	if text == "" {
		return false
	}

	profanity := []string{"badword1", "badword2", "inappropriate"}
	textLower := strings.ToLower(text)
	for _, word := range profanity {
		if strings.Contains(textLower, word) {
			return true
		}
	}
	return false
}

func hasExcessiveLinks(text string) bool {
	// Count HTTP/HTTPS links
	linkPattern := `https?://[^\s]+`
	re := regexp.MustCompile(linkPattern)
	matches := re.FindAllString(text, -1)
	return len(matches) > 3
}

func isSupportedImageFormat(url string) bool {
	if url == "" {
		return true // Optional field
	}

	supportedFormats := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	urlLower := strings.ToLower(url)
	for _, format := range supportedFormats {
		if strings.HasSuffix(urlLower, format) {
			return true
		}
	}
	return false
}

func isMaliciousURL(url string) bool {
	if url == "" {
		return false
	}

	maliciousDomains := []string{"malicious.com", "phishing.net", "scam.org"}
	urlLower := strings.ToLower(url)
	for _, domain := range maliciousDomains {
		if strings.Contains(urlLower, domain) {
			return true
		}
	}
	return false
}

func isValidTheme(theme string) bool {
	validThemes := []string{"light", "dark", "auto"}
	for _, valid := range validThemes {
		if theme == valid {
			return true
		}
	}
	return false
}

func isValidLanguageCode(lang string) bool {
	validLanguages := []string{"en", "es", "fr", "de", "it", "ja", "ko", "zh"}
	for _, valid := range validLanguages {
		if lang == valid {
			return true
		}
	}
	return false
}

func isValidTimezone(tz string) bool {
	validTimezones := []string{"UTC", "America/New_York", "Europe/London", "Asia/Tokyo", "Australia/Sydney"}
	for _, valid := range validTimezones {
		if tz == valid {
			return true
		}
	}
	return false
}

func containsMaliciousPatterns(query string) bool {
	maliciousPatterns := []string{"<script", "javascript:", "eval(", "union select"}
	queryLower := strings.ToLower(query)
	for _, pattern := range maliciousPatterns {
		if strings.Contains(queryLower, pattern) {
			return true
		}
	}
	return false
}

func isOnlyWhitespace(text string) bool {
	return strings.TrimSpace(text) == ""
}

func hasInvalidFacets(facets []string) bool {
	validFacets := []string{"brand", "category", "price_range", "rating", "availability", "location"}
	for _, facet := range facets {
		valid := false
		for _, validFacet := range validFacets {
			if facet == validFacet {
				valid = true
				break
			}
		}
		if !valid {
			return true
		}
	}
	return false
}

func hasInvalidCustomFilterKeys(filters map[string]string) bool {
	validKeys := []string{"color", "size", "material", "condition", "shipping"}
	for key := range filters {
		valid := false
		for _, validKey := range validKeys {
			if key == validKey {
				valid = true
				break
			}
		}
		if !valid {
			return true
		}
	}
	return false
}

func isValidDistanceUnit(unit string) bool {
	validUnits := []string{"km", "mi", "m", "ft"}
	for _, valid := range validUnits {
		if unit == valid {
			return true
		}
	}
	return false
}

func isValidSortField(field string) bool {
	validFields := []string{"relevance", "price", "date", "rating", "popularity", "name"}
	for _, valid := range validFields {
		if field == valid {
			return true
		}
	}
	return false
}

func isValidSortDirection(direction string) bool {
	validDirections := []string{"asc", "desc"}
	for _, valid := range validDirections {
		if direction == valid {
			return true
		}
	}
	return false
}

func hasInvalidHighlightFields(fields []string) bool {
	validFields := []string{"title", "description", "content", "tags", "category"}
	for _, field := range fields {
		valid := false
		for _, validField := range validFields {
			if field == validField {
				valid = true
				break
			}
		}
		if !valid {
			return true
		}
	}
	return false
}

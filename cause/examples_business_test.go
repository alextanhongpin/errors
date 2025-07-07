package cause_test

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world example: Financial Transaction System
type TransferRequest struct {
	FromAccountID string            `json:"from_account_id"`
	ToAccountID   string            `json:"to_account_id"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Reference     string            `json:"reference"`
	Description   string            `json:"description,omitempty"`
	ScheduledDate *time.Time        `json:"scheduled_date,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func (tr *TransferRequest) Validate() error {
	return cause.Map{
		"from_account_id": cause.Required(tr.FromAccountID).
			Select(map[string]bool{
				"invalid from account ID format":  !isValidAccountID(tr.FromAccountID),
				"cannot transfer to same account": tr.FromAccountID == tr.ToAccountID,
			}),

		"to_account_id": cause.Required(tr.ToAccountID).
			Select(map[string]bool{
				"invalid to account ID format": !isValidAccountID(tr.ToAccountID),
			}),

		"amount": cause.Required(tr.Amount).
			Select(map[string]bool{
				"amount must be positive":            tr.Amount <= 0,
				"amount exceeds transfer limit":      tr.Amount > 1000000,
				"amount has too many decimal places": !hasValidDecimalPlaces(tr.Amount),
			}),

		"currency": cause.Required(tr.Currency).
			Select(map[string]bool{
				"invalid currency code":  !isValidCurrencyCode(tr.Currency),
				"currency not supported": !isSupportedCurrency(tr.Currency),
			}),

		"reference": cause.Required(tr.Reference).
			Select(map[string]bool{
				"reference too short":            len(tr.Reference) < 3,
				"reference too long":             len(tr.Reference) > 50,
				"reference must be alphanumeric": !isAlphaNumeric(tr.Reference),
			}),

		"description": cause.Optional(tr.Description).
			Select(map[string]bool{
				"description too long": len(tr.Description) > 200,
			}),

		"scheduled_date": cause.Optional(tr.ScheduledDate).
			Select(map[string]bool{
				"scheduled date cannot be in the past": tr.ScheduledDate != nil && tr.ScheduledDate.Before(time.Now()),
				"scheduled date too far in future":     tr.ScheduledDate != nil && tr.ScheduledDate.After(time.Now().AddDate(1, 0, 0)),
			}),

		"metadata": cause.Optional(tr.Metadata).
			Select(map[string]bool{
				"too many metadata entries":                len(tr.Metadata) > 10,
				"metadata contains invalid keys or values": hasInvalidMetadata(tr.Metadata),
			}),
	}.Err()
}

// Real-world example: Content Management System
type BlogPost struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      string    `json:"author_id"`
	Categories    []string  `json:"categories"`
	Tags          []string  `json:"tags"`
	PublishDate   time.Time `json:"publish_date"`
	Status        string    `json:"status"`
	SEOData       SEOData   `json:"seo_data"`
	Featured      bool      `json:"featured"`
	AllowComments bool      `json:"allow_comments"`
}

type SEOData struct {
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Keywords        []string `json:"keywords"`
	CanonicalURL    string   `json:"canonical_url,omitempty"`
	OpenGraphImage  string   `json:"og_image,omitempty"`
}

func (bp *BlogPost) Validate() error {
	return cause.Map{
		"title": cause.Required(bp.Title).
			When(len(bp.Title) < 5, "title too short").
			When(len(bp.Title) > 100, "title too long").
			When(hasRepeatedWords(bp.Title), "title contains repeated words").Err(),

		"content": cause.Required(bp.Content).
			When(len(bp.Content) < 100, "content too short for publication").
			When(len(bp.Content) > 50000, "content exceeds maximum length").
			When(hasProhibitedContent(bp.Content), "content contains prohibited material").Err(),

		"author_id": cause.Required(bp.AuthorID).
			When(!isValidUserID(bp.AuthorID), "invalid author ID").Err(),

		"categories": cause.Required(bp.Categories).
			When(len(bp.Categories) == 0, "at least one category required").
			When(len(bp.Categories) > 5, "too many categories").
			When(hasInvalidCategories(bp.Categories), "contains invalid categories").Err(),

		"tags": cause.Optional(bp.Tags).
			When(len(bp.Tags) > 20, "too many tags").
			When(hasDuplicates(bp.Tags), "duplicate tags not allowed").
			When(hasInvalidTags(bp.Tags), "contains invalid tags").Err(),

		"publish_date": cause.Required(bp.PublishDate).
			When(bp.Status == "published" && bp.PublishDate.After(time.Now()), "cannot publish in future").
			When(bp.PublishDate.Before(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)), "publish date too old").Err(),

		"status": cause.Required(bp.Status).
			When(!isValidPostStatus(bp.Status), "invalid post status").Err(),

		"seo_data": cause.Required(bp.SEOData).Err(),
	}.Err()
}

func (seo *SEOData) Validate() error {
	return cause.Map{
		"meta_title": cause.Required(seo.MetaTitle).
			When(len(seo.MetaTitle) < 10, "meta title too short").
			When(len(seo.MetaTitle) > 60, "meta title too long for SEO").Err(),

		"meta_description": cause.Required(seo.MetaDescription).
			When(len(seo.MetaDescription) < 50, "meta description too short").
			When(len(seo.MetaDescription) > 160, "meta description too long for SEO").Err(),

		"keywords": cause.Optional(seo.Keywords).
			When(len(seo.Keywords) > 10, "too many keywords").
			When(hasDuplicates(seo.Keywords), "duplicate keywords").Err(),

		"canonical_url": cause.Optional(seo.CanonicalURL).
			When(!isValidURL(seo.CanonicalURL), "invalid canonical URL").Err(),

		"og_image": cause.Optional(seo.OpenGraphImage).
			When(!isValidImageURL(seo.OpenGraphImage), "invalid Open Graph image URL").Err(),
	}.Err()
}

// Real-world example: Event Management System
type EventRegistration struct {
	EventID          string           `json:"event_id"`
	AttendeeInfo     AttendeeInfo     `json:"attendee_info"`
	TicketType       string           `json:"ticket_type"`
	Quantity         int              `json:"quantity"`
	PaymentInfo      PaymentInfo      `json:"payment_info"`
	DietaryReqs      []string         `json:"dietary_requirements,omitempty"`
	Accessibility    []string         `json:"accessibility_needs,omitempty"`
	EmergencyContact EmergencyContact `json:"emergency_contact"`
	CustomFields     map[string]any   `json:"custom_fields,omitempty"`
}

type AttendeeInfo struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Company     string    `json:"company,omitempty"`
	JobTitle    string    `json:"job_title,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitzero"`
}

type PaymentInfo struct {
	Method         string  `json:"method"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PromoCode      string  `json:"promo_code,omitempty"`
	BillingAddress Address `json:"billing_address"`
}

type EmergencyContact struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
	Email        string `json:"email,omitempty"`
}

func (er *EventRegistration) Validate() error {
	return cause.Map{
		"event_id": cause.Required(er.EventID).
			When(!isValidEventID(er.EventID), "invalid event ID").
			When(!isEventOpen(er.EventID), "event registration is closed").Err(),

		"attendee_info": cause.Required(er.AttendeeInfo).Err(),

		"ticket_type": cause.Required(er.TicketType).
			When(!isValidTicketType(er.EventID, er.TicketType), "invalid ticket type for event").Err(),

		"quantity": cause.Required(er.Quantity).
			When(er.Quantity <= 0, "quantity must be positive").
			When(er.Quantity > getMaxTicketsPerRegistration(er.EventID), "exceeds maximum tickets per registration").
			When(!isTicketAvailable(er.EventID, er.TicketType, er.Quantity), "not enough tickets available").Err(),

		"payment_info": cause.Required(er.PaymentInfo).Err(),

		"dietary_requirements": cause.Optional(er.DietaryReqs).
			When(len(er.DietaryReqs) > 10, "too many dietary requirements").
			When(hasInvalidDietaryReqs(er.DietaryReqs), "contains invalid dietary requirements").Err(),

		"accessibility_needs": cause.Optional(er.Accessibility).
			When(len(er.Accessibility) > 10, "too many accessibility needs").
			When(hasInvalidAccessibilityNeeds(er.Accessibility), "contains invalid accessibility needs").Err(),

		"emergency_contact": cause.Required(er.EmergencyContact).Err(),

		"custom_fields": cause.Optional(er.CustomFields).
			When(!validateCustomFields(er.EventID, er.CustomFields), "invalid custom field values").Err(),
	}.Err()
}

func (ai *AttendeeInfo) Validate() error {
	return cause.Map{
		"first_name": cause.Required(ai.FirstName).
			When(len(ai.FirstName) < 2, "first name too short").
			When(len(ai.FirstName) > 50, "first name too long").
			When(containsNumbers(ai.FirstName), "first name cannot contain numbers").Err(),

		"last_name": cause.Required(ai.LastName).
			When(len(ai.LastName) < 2, "last name too short").
			When(len(ai.LastName) > 50, "last name too long").
			When(containsNumbers(ai.LastName), "last name cannot contain numbers").Err(),

		"email": cause.Required(ai.Email).
			When(!isValidEmail(ai.Email), "invalid email format").
			When(isDisposableEmail(ai.Email), "disposable emails not allowed").Err(),

		"phone": cause.Required(ai.Phone).
			When(!isValidPhoneNumber(ai.Phone), "invalid phone number").Err(),

		"company": cause.Optional(ai.Company).
			When(len(ai.Company) > 100, "company name too long").Err(),

		"job_title": cause.Optional(ai.JobTitle).
			When(len(ai.JobTitle) > 100, "job title too long").Err(),

		"date_of_birth": cause.Optional(ai.DateOfBirth).
			When(!ai.DateOfBirth.IsZero() && getAge(ai.DateOfBirth) < 18, "must be at least 18 years old").
			When(!ai.DateOfBirth.IsZero() && ai.DateOfBirth.After(time.Now()), "date of birth cannot be in future").Err(),
	}.Err()
}

func (pi *PaymentInfo) Validate() error {
	return cause.Map{
		"method": cause.Required(pi.Method).
			When(!isValidPaymentMethod(pi.Method), "invalid payment method").Err(),

		"amount": cause.Required(pi.Amount).
			When(pi.Amount <= 0, "amount must be positive").
			When(!hasValidDecimalPlaces(pi.Amount), "amount has too many decimal places").Err(),

		"currency": cause.Required(pi.Currency).
			When(!isValidCurrencyCode(pi.Currency), "invalid currency code").Err(),

		"promo_code": cause.Optional(pi.PromoCode).
			When(!isValidPromoCode(pi.PromoCode), "invalid promo code format").
			When(!isPromoCodeActive(pi.PromoCode), "promo code is expired or inactive").Err(),

		"billing_address": cause.Required(pi.BillingAddress).Err(),
	}.Err()
}

func (ec *EmergencyContact) Validate() error {
	return cause.Map{
		"name": cause.Required(ec.Name).
			When(len(ec.Name) < 2, "emergency contact name too short").
			When(len(ec.Name) > 100, "emergency contact name too long").Err(),

		"relationship": cause.Required(ec.Relationship).
			When(!isValidRelationship(ec.Relationship), "invalid relationship type").Err(),

		"phone": cause.Required(ec.Phone).
			When(!isValidPhoneNumber(ec.Phone), "invalid emergency contact phone number").Err(),

		"email": cause.Optional(ec.Email).
			When(!isValidEmail(ec.Email), "invalid emergency contact email").Err(),
	}.Err()
}

// Example demonstrating complex business validation
func ExampleTransferRequest_validation() {
	// Valid transfer
	validTransfer := &TransferRequest{
		FromAccountID: "ACC-123456789",
		ToAccountID:   "ACC-987654321",
		Amount:        1500.50,
		Currency:      "USD",
		Reference:     "RENT2024JAN",
		Description:   "Monthly rent payment",
	}

	fmt.Println("Valid transfer:")
	if err := validTransfer.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Invalid transfer
	invalidTransfer := &TransferRequest{
		FromAccountID: "invalid-account",
		ToAccountID:   "invalid-account",        // Same as from account
		Amount:        -100.00,                  // Negative amount
		Currency:      "XYZ",                    // Invalid currency
		Reference:     "$$",                     // Too short and invalid chars
		Description:   strings.Repeat("x", 250), // Too long
	}

	fmt.Println("\nInvalid transfer:")
	if err := invalidTransfer.Validate(); err != nil {
		if validationErr, ok := err.(interface{ Map() map[string]any }); ok {
			fieldErrors := validationErr.Map()
			fields := slices.Sorted(maps.Keys(fieldErrors))
			for _, field := range fields {
				fieldErr := fieldErrors[field]
				fmt.Printf("  %s: %v\n", field, fieldErr)
			}
		}
	}

	// Output:
	// Valid transfer:
	// Validation passed!
	//
	// Invalid transfer:
	//   amount: amount must be positive
	//   currency: currency not supported, invalid currency code
	//   description: description too long
	//   from_account_id: cannot transfer to same account, invalid from account ID format
	//   reference: reference must be alphanumeric, reference too short
	//   to_account_id: invalid to account ID format
}

func ExampleEventRegistration_validation() {
	// Valid registration
	validRegistration := &EventRegistration{
		EventID: "EVT-CONF2024-001",
		AttendeeInfo: AttendeeInfo{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice@company.com",
			Phone:     "+1-555-123-4567",
			Company:   "Tech Corp",
			JobTitle:  "Software Engineer",
		},
		TicketType: "standard",
		Quantity:   1,
		PaymentInfo: PaymentInfo{
			Method:   "credit_card",
			Amount:   299.99,
			Currency: "USD",
			BillingAddress: Address{
				Street:  "123 Business Ave",
				City:    "San Francisco",
				State:   "CA",
				Country: "US",
				ZipCode: "94105",
			},
		},
		EmergencyContact: EmergencyContact{
			Name:         "Bob Johnson",
			Relationship: "spouse",
			Phone:        "+1-555-987-6543",
			Email:        "bob@example.com",
		},
	}

	fmt.Println("Valid event registration:")
	if err := validRegistration.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Output:
	// Valid event registration:
	// Validation passed!
}

// Helper functions for financial validation
func isValidAccountID(accountID string) bool {
	return len(accountID) >= 10 && strings.HasPrefix(accountID, "ACC-")
}

func hasValidDecimalPlaces(amount float64) bool {
	// Check if amount has more than 2 decimal places
	formatted := fmt.Sprintf("%.2f", amount)
	original := fmt.Sprintf("%g", amount)
	return formatted == original || len(strings.Split(original, ".")) <= 1 ||
		len(strings.Split(original, ".")[1]) <= 2
}

func isValidCurrencyCode(currency string) bool {
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF"}
	return slices.Contains(validCurrencies, currency)
}

func isSupportedCurrency(currency string) bool {
	// For demo, assume all valid currencies are supported
	return isValidCurrencyCode(currency)
}

func isAlphaNumeric(s string) bool {
	for _, char := range s {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

func hasInvalidMetadata(metadata map[string]string) bool {
	for key, value := range metadata {
		if len(key) > 50 || len(value) > 200 {
			return true
		}
	}
	return false
}

// Helper functions for content validation
func hasRepeatedWords(title string) bool {
	words := strings.Fields(strings.ToLower(title))
	seen := make(map[string]int)
	for _, word := range words {
		seen[word]++
		if seen[word] > 2 {
			return true
		}
	}
	return false
}

func hasProhibitedContent(content string) bool {
	prohibited := []string{"spam", "illegal", "violence"}
	contentLower := strings.ToLower(content)
	for _, word := range prohibited {
		if strings.Contains(contentLower, word) {
			return true
		}
	}
	return false
}

func isValidUserID(userID string) bool {
	return len(userID) >= 5 && isAlphaNumeric(userID)
}

func hasInvalidCategories(categories []string) bool {
	validCategories := []string{"technology", "business", "lifestyle", "health", "education", "entertainment"}
	for _, category := range categories {
		valid := slices.Contains(validCategories, category)
		if !valid {
			return true
		}
	}
	return false
}

func hasDuplicates(slice []string) bool {
	seen := make(map[string]bool)
	for _, item := range slice {
		if seen[item] {
			return true
		}
		seen[item] = true
	}
	return false
}

func hasInvalidTags(tags []string) bool {
	for _, tag := range tags {
		if len(tag) < 2 || len(tag) > 20 || !isAlphaNumeric(tag) {
			return true
		}
	}
	return false
}

func isValidPostStatus(status string) bool {
	validStatuses := []string{"draft", "published", "archived", "scheduled"}
	return slices.Contains(validStatuses, status)
}

func isValidURL(url string) bool {
	if url == "" {
		return true // Optional field
	}
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidImageURL(url string) bool {
	if url == "" {
		return true // Optional field
	}
	return isValidURL(url) && (strings.HasSuffix(url, ".jpg") ||
		strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") ||
		strings.HasSuffix(url, ".gif"))
}

// Helper functions for event validation
func isValidEventID(eventID string) bool {
	return len(eventID) >= 10 && strings.HasPrefix(eventID, "EVT-")
}

func isEventOpen(eventID string) bool {
	// Simplified: assume event is open for demo
	return true
}

func isValidTicketType(eventID, ticketType string) bool {
	validTypes := []string{"standard", "premium", "vip", "student"}
	return slices.Contains(validTypes, ticketType)
}

func getMaxTicketsPerRegistration(eventID string) int {
	// Simplified: return fixed limit
	return 10
}

func isTicketAvailable(eventID, ticketType string, quantity int) bool {
	// Simplified: assume tickets are available
	return true
}

func hasInvalidDietaryReqs(reqs []string) bool {
	validReqs := []string{"vegetarian", "vegan", "gluten-free", "dairy-free", "nut-free", "halal", "kosher"}
	for _, req := range reqs {
		valid := slices.Contains(validReqs, req)
		if !valid {
			return true
		}
	}
	return false
}

func hasInvalidAccessibilityNeeds(needs []string) bool {
	validNeeds := []string{"wheelchair", "hearing-impaired", "visually-impaired", "mobility-assistance"}
	for _, need := range needs {
		valid := slices.Contains(validNeeds, need)
		if !valid {
			return true
		}
	}
	return false
}

func validateCustomFields(eventID string, fields map[string]any) bool {
	// Simplified validation
	return len(fields) <= 5
}

func isValidPaymentMethod(method string) bool {
	validMethods := []string{"credit_card", "debit_card", "paypal", "bank_transfer"}
	return slices.Contains(validMethods, method)
}

func isValidPromoCode(code string) bool {
	if code == "" {
		return true // Optional field
	}
	return len(code) >= 5 && len(code) <= 15 && isAlphaNumeric(code)
}

func isPromoCodeActive(code string) bool {
	if code == "" {
		return true // Optional field
	}
	// Simplified: assume promo code is active
	return true
}

func isValidRelationship(relationship string) bool {
	validRelationships := []string{"spouse", "parent", "child", "sibling", "friend", "colleague", "other"}
	return slices.Contains(validRelationships, relationship)
}

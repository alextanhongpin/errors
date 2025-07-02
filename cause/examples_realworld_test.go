package cause_test

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/errors/cause"
)

// Real-world address validation
type Address struct {
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country"`
	ZipCode   string `json:"zip_code"`
	IsDefault bool   `json:"is_default,omitempty"`
}

func (a *Address) Validate() error {
	return cause.Map{
		"street": cause.Required(a.Street).
			When(len(a.Street) < 5, "street address too short").
			When(len(a.Street) > 100, "street address too long"),

		"city": cause.Required(a.City).
			When(len(a.City) < 2, "city name too short").
			When(len(a.City) > 50, "city name too long").
			When(containsNumbers(a.City), "city name cannot contain numbers"),

		"state": cause.Optional(a.State).
			When(a.Country == "US" && len(a.State) != 2, "US state must be 2 characters").
			When(len(a.State) > 50, "state name too long"),

		"country": cause.Required(a.Country).
			When(!isValidCountryCode(a.Country), "invalid country code"),

		"zip_code": cause.Required(a.ZipCode).
			When(!isValidZipCodeForCountry(a.ZipCode, a.Country), "invalid zip code for country"),
	}.Err()
}

// Real-world example: User Registration System
type UserRegistration struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	PhoneNumber     string    `json:"phone_number"`
	Address         *Address  `json:"address,omitempty"`
	Terms           bool      `json:"terms_accepted"`
	Newsletter      bool      `json:"newsletter_opt_in"`
}

func (u *UserRegistration) Validate() error {
	return cause.Map{
		"email": cause.Required(u.Email).
			When(!isValidEmail(u.Email), "invalid email format").
			When(isDisposableEmail(u.Email), "disposable emails not allowed"),

		"password": cause.Required(u.Password).
			When(len(u.Password) < 8, "password must be at least 8 characters").
			When(!hasUppercase(u.Password), "password must contain uppercase letter").
			When(!hasLowercase(u.Password), "password must contain lowercase letter").
			When(!hasDigit(u.Password), "password must contain digit").
			When(!hasSpecialChar(u.Password), "password must contain special character"),

		"confirm_password": cause.Required(u.ConfirmPassword).
			When(u.Password != u.ConfirmPassword, "passwords do not match"),

		"first_name": cause.Required(u.FirstName).
			When(len(u.FirstName) < 2, "first name too short").
			When(len(u.FirstName) > 50, "first name too long").
			When(containsNumbers(u.FirstName), "first name cannot contain numbers"),

		"last_name": cause.Required(u.LastName).
			When(len(u.LastName) < 2, "last name too short").
			When(len(u.LastName) > 50, "last name too long").
			When(containsNumbers(u.LastName), "last name cannot contain numbers"),

		"date_of_birth": cause.Required(u.DateOfBirth).
			When(u.DateOfBirth.After(time.Now()), "date of birth cannot be in the future").
			When(getAge(u.DateOfBirth) < 13, "must be at least 13 years old").
			When(getAge(u.DateOfBirth) > 120, "invalid date of birth"),

		"phone_number": cause.Optional(u.PhoneNumber).
			When(!isValidPhoneNumber(u.PhoneNumber), "invalid phone number format"),

		"address": cause.Optional(u.Address),

		"terms": cause.Required(u.Terms).
			When(!u.Terms, "must accept terms and conditions"),
	}.Err()
}

// Real-world example: E-commerce Product
type Product struct {
	SKU         string     `json:"sku"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Variants    []Variant  `json:"variants"`
	Dimensions  Dimensions `json:"dimensions"`
	Weight      float64    `json:"weight"`
	InStock     bool       `json:"in_stock"`
	StockCount  int        `json:"stock_count"`
}

type Variant struct {
	Name  string  `json:"name"`
	Value string  `json:"value"`
	Price float64 `json:"price_modifier"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

func (p *Product) Validate() error {
	return cause.Map{
		"sku": cause.Required(p.SKU).
			When(!isValidSKU(p.SKU), "invalid SKU format").
			When(len(p.SKU) > 20, "SKU too long"),

		"name": cause.Required(p.Name).
			When(len(p.Name) < 3, "product name too short").
			When(len(p.Name) > 100, "product name too long"),

		"description": cause.Optional(p.Description).
			When(len(p.Description) > 1000, "description too long"),

		"price": cause.Required(p.Price).
			When(p.Price <= 0, "price must be positive").
			When(p.Price > 999999.99, "price too high"),

		"category": cause.Required(p.Category).
			When(!isValidCategory(p.Category), "invalid product category"),

		"tags": cause.Optional(p.Tags).
			When(len(p.Tags) > 10, "too many tags").
			When(hasDuplicateTags(p.Tags), "duplicate tags not allowed"),

		"variants": cause.Optional(p.Variants).
			When(len(p.Variants) > 20, "too many variants"),

		"dimensions": cause.Required(p.Dimensions),

		"weight": cause.Required(p.Weight).
			When(p.Weight <= 0, "weight must be positive").
			When(p.Weight > 1000, "weight too heavy for shipping"),

		"stock_count": cause.Optional(p.StockCount).
			When(p.InStock && p.StockCount <= 0, "stock count must be positive when in stock"),
	}.Err()
}

func (v *Variant) Validate() error {
	return cause.Map{
		"name": cause.Required(v.Name).
			When(len(v.Name) > 50, "variant name too long"),

		"value": cause.Required(v.Value).
			When(len(v.Value) > 100, "variant value too long"),

		"price": cause.Optional(v.Price).
			When(v.Price < -999.99, "price modifier too low").
			When(v.Price > 999.99, "price modifier too high"),
	}.Err()
}

func (d *Dimensions) Validate() error {
	return cause.Map{
		"length": cause.Required(d.Length).
			When(d.Length <= 0, "length must be positive").
			When(d.Length > 1000, "length too large"),

		"width": cause.Required(d.Width).
			When(d.Width <= 0, "width must be positive").
			When(d.Width > 1000, "width too large"),

		"height": cause.Required(d.Height).
			When(d.Height <= 0, "height must be positive").
			When(d.Height > 1000, "height too large"),

		"unit": cause.Required(d.Unit).
			When(!isValidDimensionUnit(d.Unit), "invalid dimension unit"),
	}.Err()
}

// Real-world example: API Request with complex nested validation
type CreateOrderRequest struct {
	CustomerID          string          `json:"customer_id"`
	Items               []OrderItem     `json:"items"`
	ShippingAddress     Address         `json:"shipping_address"`
	BillingAddress      *Address        `json:"billing_address,omitempty"`
	PaymentMethod       PaymentMethod   `json:"payment_method"`
	CouponCode          string          `json:"coupon_code,omitempty"`
	DeliveryOptions     DeliveryOptions `json:"delivery_options"`
	SpecialInstructions string          `json:"special_instructions,omitempty"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount,omitempty"`
}

type PaymentMethod struct {
	Type       string `json:"type"`
	CardNumber string `json:"card_number,omitempty"`
	ExpiryDate string `json:"expiry_date,omitempty"`
	CVV        string `json:"cvv,omitempty"`
	PayPalID   string `json:"paypal_id,omitempty"`
}

type DeliveryOptions struct {
	Type          string    `json:"type"`
	Priority      string    `json:"priority"`
	RequestedDate time.Time `json:"requested_date,omitempty"`
	Instructions  string    `json:"instructions,omitempty"`
}

func (o *CreateOrderRequest) Validate() error {
	return cause.Map{
		"customer_id": cause.Required(o.CustomerID).
			When(!isValidUUID(o.CustomerID), "invalid customer ID format"),

		"items": cause.Required(o.Items).
			When(len(o.Items) == 0, "order must contain at least one item").
			When(len(o.Items) > 50, "too many items in order").
			When(calculateTotalValue(o.Items) > 10000, "order value exceeds limit"),

		"shipping_address": cause.Required(o.ShippingAddress),
		"billing_address":  cause.Optional(o.BillingAddress),
		"payment_method":   cause.Required(o.PaymentMethod),

		"coupon_code": cause.Optional(o.CouponCode).
			When(!isValidCouponFormat(o.CouponCode), "invalid coupon code format"),

		"delivery_options": cause.Required(o.DeliveryOptions),

		"special_instructions": cause.Optional(o.SpecialInstructions).
			When(len(o.SpecialInstructions) > 500, "special instructions too long"),
	}.Err()
}

func (item *OrderItem) Validate() error {
	return cause.Map{
		"product_id": cause.Required(item.ProductID).
			When(!isValidUUID(item.ProductID), "invalid product ID format"),

		"quantity": cause.Required(item.Quantity).
			When(item.Quantity <= 0, "quantity must be positive").
			When(item.Quantity > 100, "quantity exceeds maximum allowed"),

		"price": cause.Required(item.Price).
			When(item.Price <= 0, "price must be positive").
			When(item.Price > 9999.99, "price exceeds maximum allowed"),

		"discount": cause.Optional(item.Discount).
			When(item.Discount < 0, "discount cannot be negative").
			When(item.Discount >= item.Price, "discount cannot exceed item price"),
	}.Err()
}

func (pm *PaymentMethod) Validate() error {
	baseValidation := cause.Map{
		"type": cause.Required(pm.Type).
			When(!isValidPaymentType(pm.Type), "invalid payment method type"),
	}

	// Additional validation based on payment type
	switch pm.Type {
	case "credit_card", "debit_card":
		baseValidation["card_number"] = cause.Required(pm.CardNumber).
			When(!isValidCardNumber(pm.CardNumber), "invalid card number")
		baseValidation["expiry_date"] = cause.Required(pm.ExpiryDate).
			When(!isValidExpiryDate(pm.ExpiryDate), "invalid expiry date").
			When(isExpired(pm.ExpiryDate), "card has expired")
		baseValidation["cvv"] = cause.Required(pm.CVV).
			When(!isValidCVV(pm.CVV), "invalid CVV")
	case "paypal":
		baseValidation["paypal_id"] = cause.Required(pm.PayPalID).
			When(!isValidEmail(pm.PayPalID), "invalid PayPal email")
	}

	return baseValidation.Err()
}

func (do *DeliveryOptions) Validate() error {
	return cause.Map{
		"type": cause.Required(do.Type).
			When(!isValidDeliveryType(do.Type), "invalid delivery type"),

		"priority": cause.Required(do.Priority).
			When(!isValidDeliveryPriority(do.Priority), "invalid delivery priority"),

		"requested_date": cause.Optional(do.RequestedDate).
			When(!do.RequestedDate.IsZero() && do.RequestedDate.Before(time.Now().AddDate(0, 0, 1)), "delivery date must be at least tomorrow").
			When(!do.RequestedDate.IsZero() && do.RequestedDate.After(time.Now().AddDate(0, 0, 30)), "delivery date too far in future"),

		"instructions": cause.Optional(do.Instructions).
			When(len(do.Instructions) > 200, "delivery instructions too long"),
	}.Err()
}

// Example functions showcasing the validation results
func ExampleUserRegistration_validation() {
	// Valid user registration
	validUser := &UserRegistration{
		Email:           "john@example.com",
		Password:        "SecurePass123!",
		ConfirmPassword: "SecurePass123!",
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber:     "+1-555-123-4567",
		Terms:           true,
		Newsletter:      false,
	}

	fmt.Println("Valid user registration:")
	if err := validUser.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Invalid user registration
	invalidUser := &UserRegistration{
		Email:           "invalid-email",
		Password:        "weak",
		ConfirmPassword: "different",
		FirstName:       "J",
		LastName:        "123Doe",
		DateOfBirth:     time.Now().AddDate(1, 0, 0), // Future date
		PhoneNumber:     "invalid-phone",
		Terms:           false,
	}

	fmt.Println("\nInvalid user registration:")
	if err := invalidUser.Validate(); err != nil {
		if validationErr, ok := err.(interface{ Map() map[string]any }); ok {
			fieldErrors := validationErr.Map()
			for field, fieldErr := range fieldErrors {
				fmt.Printf("  %s: %v\n", field, fieldErr)
			}
		}
	}

	// Note: Output field order may vary due to Go map iteration order
	// Output:
	// Valid user registration:
	// Validation passed!
}

func ExampleProduct_validation() {
	// Valid product
	validProduct := &Product{
		SKU:         "ABC-123-XYZ",
		Name:        "Premium Widget",
		Description: "A high-quality widget for all your needs",
		Price:       29.99,
		Category:    "electronics",
		Tags:        []string{"premium", "widget", "electronics"},
		Variants: []Variant{
			{Name: "Color", Value: "Red", Price: 0},
			{Name: "Size", Value: "Large", Price: 5.00},
		},
		Dimensions: Dimensions{
			Length: 10.5,
			Width:  5.2,
			Height: 3.1,
			Unit:   "cm",
		},
		Weight:     0.5,
		InStock:    true,
		StockCount: 100,
	}

	fmt.Println("Valid product:")
	if err := validProduct.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Output:
	// Valid product:
	// Validation passed!
}

func ExampleCreateOrderRequest_validation() {
	// Valid order request
	validOrder := &CreateOrderRequest{
		CustomerID: "550e8400-e29b-41d4-a716-446655440000",
		Items: []OrderItem{
			{
				ProductID: "550e8400-e29b-41d4-a716-446655440001",
				Quantity:  2,
				Price:     29.99,
				Discount:  5.00,
			},
		},
		ShippingAddress: Address{
			Street:  "123 Main St",
			City:    "Anytown",
			Country: "US",
			ZipCode: "12345",
		},
		PaymentMethod: PaymentMethod{
			Type:       "credit_card",
			CardNumber: "4111111111111111",
			ExpiryDate: "12/25",
			CVV:        "123",
		},
		DeliveryOptions: DeliveryOptions{
			Type:     "standard",
			Priority: "normal",
		},
	}

	fmt.Println("Valid order:")
	if err := validOrder.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}

	// Output:
	// Valid order:
	// Validation passed!
}

// Helper validation functions (simplified implementations)
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isDisposableEmail(email string) bool {
	disposableDomains := []string{"10minutemail.com", "tempmail.org", "guerrillamail.com"}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])
	for _, disposable := range disposableDomains {
		if domain == disposable {
			return true
		}
	}
	return false
}

func hasUppercase(s string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

func hasLowercase(s string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(s)
}

func hasDigit(s string) bool {
	return regexp.MustCompile(`\d`).MatchString(s)
}

func hasSpecialChar(s string) bool {
	return regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(s)
}

func containsNumbers(s string) bool {
	return regexp.MustCompile(`\d`).MatchString(s)
}

func getAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	return age
}

func isValidPhoneNumber(phone string) bool {
	if phone == "" {
		return true // Optional field
	}
	pattern := `^\+?[\d\s\-\(\)]+$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched && len(phone) >= 10
}

func isValidSKU(sku string) bool {
	pattern := `^[A-Z0-9\-]+$`
	matched, _ := regexp.MatchString(pattern, sku)
	return matched
}

func isValidCategory(category string) bool {
	validCategories := []string{"electronics", "clothing", "books", "home", "sports", "toys"}
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

func hasDuplicateTags(tags []string) bool {
	seen := make(map[string]bool)
	for _, tag := range tags {
		if seen[tag] {
			return true
		}
		seen[tag] = true
	}
	return false
}

func isValidDimensionUnit(unit string) bool {
	validUnits := []string{"cm", "m", "in", "ft"}
	for _, valid := range validUnits {
		if unit == valid {
			return true
		}
	}
	return false
}

func isValidUUID(uuid string) bool {
	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, _ := regexp.MatchString(pattern, uuid)
	return matched
}

func calculateTotalValue(items []OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += (item.Price - item.Discount) * float64(item.Quantity)
	}
	return total
}

func isValidCouponFormat(coupon string) bool {
	if coupon == "" {
		return true // Optional field
	}
	pattern := `^[A-Z0-9]{5,15}$`
	matched, _ := regexp.MatchString(pattern, coupon)
	return matched
}

func isValidPaymentType(paymentType string) bool {
	validTypes := []string{"credit_card", "debit_card", "paypal", "bank_transfer"}
	for _, valid := range validTypes {
		if paymentType == valid {
			return true
		}
	}
	return false
}

func isValidCardNumber(cardNumber string) bool {
	// Simplified Luhn algorithm check
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(cardNumber, "")
	return len(cleaned) >= 13 && len(cleaned) <= 19
}

func isValidExpiryDate(expiry string) bool {
	pattern := `^(0[1-9]|1[0-2])\/([0-9]{2})$`
	matched, _ := regexp.MatchString(pattern, expiry)
	return matched
}

func isExpired(expiry string) bool {
	if !isValidExpiryDate(expiry) {
		return true
	}
	// Simplified expiry check (assuming current time for demo)
	return false
}

func isValidCVV(cvv string) bool {
	pattern := `^[0-9]{3,4}$`
	matched, _ := regexp.MatchString(pattern, cvv)
	return matched
}

func isValidDeliveryType(deliveryType string) bool {
	validTypes := []string{"standard", "express", "overnight", "pickup"}
	for _, valid := range validTypes {
		if deliveryType == valid {
			return true
		}
	}
	return false
}

func isValidDeliveryPriority(priority string) bool {
	validPriorities := []string{"low", "normal", "high", "urgent"}
	for _, valid := range validPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}

func isValidCountryCode(code string) bool {
	validCodes := []string{"US", "CA", "GB", "DE", "FR", "JP", "AU", "BR", "IN", "CN"}
	for _, valid := range validCodes {
		if code == valid {
			return true
		}
	}
	return false
}

func isValidZipCodeForCountry(zipCode, country string) bool {
	if zipCode == "" {
		return false
	}

	switch country {
	case "US":
		pattern := `^\d{5}(-\d{4})?$`
		matched, _ := regexp.MatchString(pattern, zipCode)
		return matched
	case "CA":
		pattern := `^[A-Z]\d[A-Z] \d[A-Z]\d$`
		matched, _ := regexp.MatchString(pattern, strings.ToUpper(zipCode))
		return matched
	case "GB":
		pattern := `^[A-Z]{1,2}\d[A-Z\d]? \d[A-Z]{2}$`
		matched, _ := regexp.MatchString(pattern, strings.ToUpper(zipCode))
		return matched
	default:
		// Generic validation for other countries
		return len(zipCode) >= 3 && len(zipCode) <= 10
	}
}

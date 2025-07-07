// Package cause provides examples of common usage patterns for the structured error handling package.
// This file demonstrates best practices and common patterns that developers should follow.
package cause

import (
	"log/slog"
	"slices"

	"github.com/alextanhongpin/errors/codes"
)

// Common error definitions - define these at package level for reusability
var (
	// User-related errors
	ErrUserNotFound      = New(codes.NotFound, "UserNotFound", "User not found")
	ErrUserAlreadyExists = New(codes.Exists, "UserExists", "User already exists")
	ErrInvalidUserData   = New(codes.BadRequest, "InvalidUser", "Invalid user data")

	// Authentication errors
	ErrUnauthorized     = New(codes.Unauthorized, "Unauthorized", "Authentication required")
	ErrInvalidToken     = New(codes.BadRequest, "InvalidToken", "Invalid or expired token")
	ErrPermissionDenied = New(codes.Forbidden, "PermissionDenied", "Insufficient permissions")

	// System errors
	ErrDatabaseConnection = New(codes.Internal, "DatabaseError", "Database connection failed")
	ErrServiceUnavailable = New(codes.Unavailable, "ServiceUnavailable", "Service temporarily unavailable")
)

// Example: User domain model with validation
type User struct {
	ID      string   `json:"id"`
	Email   string   `json:"email"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address *Address `json:"address,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

// Validate implements validation for User
func (u *User) Validate() error {
	return Map{
		"id":      Required(u.ID).Err(),
		"email":   Required(u.Email).When(!isValidEmail(u.Email), "invalid email format").Err(),
		"name":    Required(u.Name).When(len(u.Name) < 2, "name too short").Err(),
		"age":     Optional(u.Age).When(u.Age < 0, "age cannot be negative").When(u.Age > 150, "unrealistic age").Err(),
		"address": Optional(u.Address).Err(),
		"tags":    Optional(u.Tags).When(len(u.Tags) > 10, "too many tags").Err(),
	}.Err()
}

// Validate implements validation for Address
func (a *Address) Validate() error {
	return Map{
		"street":   Required(a.Street).Err(),
		"city":     Required(a.City).Err(),
		"country":  Required(a.Country).When(!isValidCountryCode(a.Country), "invalid country code").Err(),
		"zip_code": Optional(a.ZipCode).When(!isValidZipCode(a.ZipCode), "invalid zip code format").Err(),
	}.Err()
}

// Example service methods demonstrating error handling patterns

// UserService demonstrates service-level error handling
type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}

// GetUser demonstrates error wrapping and context addition
func (s *UserService) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidUserData.WithDetails(map[string]any{
			"field": "id",
			"value": id,
		})
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		// Wrap the repository error with additional context
		return nil, ErrUserNotFound.WithDetails(map[string]any{
			"user_id": id,
			"action":  "get_user",
		}).Wrap(err).WithStack()
	}

	return user, nil
}

// CreateUser demonstrates validation and error enrichment
func (s *UserService) CreateUser(user *User) error {
	// Validate the user input
	if err := user.Validate(); err != nil {
		return ErrInvalidUserData.WithDetails(map[string]any{
			"action":  "create_user",
			"user_id": user.ID,
		}).Wrap(err)
	}

	// Check if user already exists
	if existing, _ := s.repo.GetByID(user.ID); existing != nil {
		return ErrUserAlreadyExists.WithDetails(map[string]any{
			"user_id":        user.ID,
			"existing_email": existing.Email,
		}).WithAttrs(
			slog.String("operation", "create_user"),
			slog.String("conflict_field", "id"),
		)
	}

	// Attempt to create user
	if err := s.repo.Create(user); err != nil {
		return ErrDatabaseConnection.WithDetails(map[string]any{
			"action":  "create_user",
			"user_id": user.ID,
		}).Wrap(err).WithStack()
	}

	return nil
}

// ValidateUserBatch demonstrates slice validation
func (s *UserService) ValidateUserBatch(users []User) error {
	return SliceFunc(users, func(u User) error {
		return u.Validate()
	}).Validate()
}

// Helper functions for validation (would be implemented elsewhere)
func isValidEmail(email string) bool {
	// Implementation would use proper email validation
	return email != "" && len(email) > 5
}

func isValidCountryCode(code string) bool {
	// Implementation would check against ISO country codes
	return len(code) == 2
}

func isValidZipCode(zip string) bool {
	// Implementation would validate zip code format
	return len(zip) >= 5
}

// Example: Error handling middleware pattern
func ErrorHandlingMiddleware(next func() error) error {
	defer func() {
		if r := recover(); r != nil {
			// Convert panics to structured errors
			err := New(codes.Internal, "PanicRecovered", "Unexpected panic occurred").
				WithDetails(map[string]any{
					"panic_value": r,
				}).
				WithStack()

			// Log the error
			slog.Error("Panic recovered", "error", err)
		}
	}()

	return next()
}

// Example: Complex validation scenario
type OrderRequest struct {
	UserID          string      `json:"user_id"`
	Items           []OrderItem `json:"items"`
	ShippingAddress *Address    `json:"shipping_address"`
	PaymentMethod   string      `json:"payment_method"`
	Discount        *float64    `json:"discount,omitempty"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func (o *OrderRequest) Validate() error {
	return Map{
		"user_id": Required(o.UserID).Err(),
		"items": Required(o.Items).
			When(len(o.Items) == 0, "at least one item required").
			When(len(o.Items) > 100, "too many items").Err(),
		"shipping_address": Required(o.ShippingAddress).Err(),
		"payment_method": Required(o.PaymentMethod).
			When(!isValidPaymentMethod(o.PaymentMethod), "invalid payment method").Err(),
		"discount": Optional(o.Discount).
			When(o.Discount != nil && *o.Discount < 0, "discount cannot be negative").
			When(o.Discount != nil && *o.Discount > 1, "discount cannot exceed 100%").Err(),
	}.Err()
}

func (item *OrderItem) Validate() error {
	return Map{
		"product_id": Required(item.ProductID).Err(),
		"quantity": Required(item.Quantity).
			When(item.Quantity <= 0, "quantity must be positive").
			When(item.Quantity > 1000, "quantity too large").Err(),
		"price": Required(item.Price).
			When(item.Price < 0, "price cannot be negative").Err(),
	}.Err()
}

func isValidPaymentMethod(method string) bool {
	validMethods := []string{"credit_card", "debit_card", "paypal", "bank_transfer"}
	return slices.Contains(validMethods, method)
}

package cause

import (
	"testing"

	"github.com/alextanhongpin/errors/codes"
)

// BenchmarkNewError benchmarks the creation of new errors
func BenchmarkNewError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(codes.Invalid, "TestError", "Test message")
	}
}

// BenchmarkErrorWithDetails benchmarks adding details to errors
func BenchmarkErrorWithDetails(b *testing.B) {
	err := New(codes.Invalid, "TestError", "Test message")
	details := map[string]any{
		"field": "test",
		"value": "invalid",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.WithDetails(details)
	}
}

// BenchmarkErrorWithStack benchmarks stack trace capture
func BenchmarkErrorWithStack(b *testing.B) {
	err := New(codes.Invalid, "TestError", "Test message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.WithStack()
	}
}

// BenchmarkValidationSimple benchmarks simple validation
func BenchmarkValidationSimple(b *testing.B) {
	user := &User{
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
		Age:   25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

// BenchmarkValidationComplex benchmarks complex nested validation
func BenchmarkValidationComplex(b *testing.B) {
	order := &OrderRequest{
		UserID: "123",
		Items: []OrderItem{
			{ProductID: "p1", Quantity: 2, Price: 10.0},
			{ProductID: "p2", Quantity: 1, Price: 15.0},
		},
		ShippingAddress: &Address{
			Street:  "123 Main St",
			City:    "Anytown",
			Country: "US",
			ZipCode: "12345",
		},
		PaymentMethod: "credit_card",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = order.Validate()
	}
}

// BenchmarkErrorClone benchmarks error cloning
func BenchmarkErrorClone(b *testing.B) {
	err := New(codes.Invalid, "TestError", "Test message").
		WithDetails(map[string]any{"field": "test"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Clone()
	}
}

// BenchmarkSliceValidation benchmarks slice validation
func BenchmarkSliceValidation(b *testing.B) {
	users := make([]User, 10)
	for i := range users {
		users[i] = User{
			ID:    "user" + string(rune(i)),
			Email: "user" + string(rune(i)) + "@example.com",
			Name:  "User " + string(rune(i)),
			Age:   20 + i,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SliceFunc(users, func(u User) error {
			return u.Validate()
		}).Validate()
	}
}

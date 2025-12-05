package external

import (
	"context"
	"fmt"
)

// PaymentGateway defines the interface for payment processing.
type PaymentGateway interface {
	AuthorizePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	CapturePayment(ctx context.Context, paymentID string, amount float64) error
	RefundPayment(ctx context.Context, paymentID string, amount float64) error
}

// PaymentRequest represents a payment authorization request.
type PaymentRequest struct {
	OrderID       string
	Amount        float64
	Currency      string
	PaymentMethod string
	CardToken     string // For card payments
}

// PaymentResponse represents the response from payment gateway.
type PaymentResponse struct {
	PaymentID     string
	Status        string // authorized, failed, pending
	TransactionID string
}

// MockPaymentGateway is a mock implementation for development/testing.
type MockPaymentGateway struct {
	shouldSucceed bool
}

// NewMockPaymentGateway creates a new mock payment gateway.
func NewMockPaymentGateway(shouldSucceed bool) *MockPaymentGateway {
	return &MockPaymentGateway{shouldSucceed: shouldSucceed}
}

// AuthorizePayment simulates payment authorization.
func (g *MockPaymentGateway) AuthorizePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	if !g.shouldSucceed {
		return nil, fmt.Errorf("payment authorization failed")
	}

	return &PaymentResponse{
		PaymentID:     fmt.Sprintf("pay_%s", req.OrderID),
		Status:        "authorized",
		TransactionID: fmt.Sprintf("txn_%s", req.OrderID),
	}, nil
}

// CapturePayment simulates payment capture.
func (g *MockPaymentGateway) CapturePayment(ctx context.Context, paymentID string, amount float64) error {
	if !g.shouldSucceed {
		return fmt.Errorf("payment capture failed")
	}
	return nil
}

// RefundPayment simulates payment refund.
func (g *MockPaymentGateway) RefundPayment(ctx context.Context, paymentID string, amount float64) error {
	if !g.shouldSucceed {
		return fmt.Errorf("payment refund failed")
	}
	return nil
}

// StripeGateway is an example adapter for Stripe payment gateway.
// TODO: Implement actual Stripe integration
type StripeGateway struct {
	apiKey string
}

// NewStripeGateway creates a new Stripe payment gateway adapter.
func NewStripeGateway(apiKey string) *StripeGateway {
	return &StripeGateway{apiKey: apiKey}
}

// AuthorizePayment authorizes payment via Stripe API.
func (g *StripeGateway) AuthorizePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement Stripe API call
	return nil, fmt.Errorf("not implemented")
}

// CapturePayment captures payment via Stripe API.
func (g *StripeGateway) CapturePayment(ctx context.Context, paymentID string, amount float64) error {
	// TODO: Implement Stripe API call
	return fmt.Errorf("not implemented")
}

// RefundPayment refunds payment via Stripe API.
func (g *StripeGateway) RefundPayment(ctx context.Context, paymentID string, amount float64) error {
	// TODO: Implement Stripe API call
	return fmt.Errorf("not implemented")
}

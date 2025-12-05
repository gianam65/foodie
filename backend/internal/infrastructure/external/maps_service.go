package external

import (
	"context"
	"fmt"
)

// MapsService provides geocoding and routing functionality.
type MapsService interface {
	Geocode(ctx context.Context, address string) (*Location, error)
	CalculateDistance(ctx context.Context, from, to Location) (float64, error)
	EstimateDeliveryTime(ctx context.Context, from, to Location) (int, error) // minutes
}

// Location represents a geographic location.
type Location struct {
	Latitude  float64
	Longitude float64
	Address   string
}

// MockMapsService is a mock implementation for development/testing.
type MockMapsService struct {
}

// NewMockMapsService creates a new mock maps service.
func NewMockMapsService() *MockMapsService {
	return &MockMapsService{}
}

// Geocode simulates geocoding an address.
func (s *MockMapsService) Geocode(ctx context.Context, address string) (*Location, error) {
	// Mock: return fixed coordinates
	return &Location{
		Latitude:  10.762622,
		Longitude: 106.660172,
		Address:   address,
	}, nil
}

// CalculateDistance calculates distance between two locations in kilometers.
func (s *MockMapsService) CalculateDistance(ctx context.Context, from, to Location) (float64, error) {
	// Mock: return fixed distance
	return 5.5, nil
}

// EstimateDeliveryTime estimates delivery time in minutes.
func (s *MockMapsService) EstimateDeliveryTime(ctx context.Context, from, to Location) (int, error) {
	// Mock: return fixed time
	return 30, nil
}

// GoogleMapsService is an example adapter for Google Maps API.
// TODO: Implement actual Google Maps integration
type GoogleMapsService struct {
	apiKey string
}

// NewGoogleMapsService creates a new Google Maps service adapter.
func NewGoogleMapsService(apiKey string) *GoogleMapsService {
	return &GoogleMapsService{apiKey: apiKey}
}

// Geocode geocodes an address using Google Maps API.
func (s *GoogleMapsService) Geocode(ctx context.Context, address string) (*Location, error) {
	// TODO: Implement Google Maps Geocoding API call
	return nil, fmt.Errorf("not implemented")
}

// CalculateDistance calculates distance using Google Maps Distance Matrix API.
func (s *GoogleMapsService) CalculateDistance(ctx context.Context, from, to Location) (float64, error) {
	// TODO: Implement Google Maps Distance Matrix API call
	return 0, fmt.Errorf("not implemented")
}

// EstimateDeliveryTime estimates delivery time using Google Maps API.
func (s *GoogleMapsService) EstimateDeliveryTime(ctx context.Context, from, to Location) (int, error) {
	// TODO: Implement Google Maps Directions API call
	return 0, fmt.Errorf("not implemented")
}

package utils

import (
	"testing"
)

func TestCelciusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celcius  float64
		expected float64
	}{
		{"zero", 0, 32},
		{"boiling", 100, 212},
		{"negative", -40, -40},
		{"fractional", 28.5, 83.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelciusToFahrenheit(tt.celcius)
			if result != tt.expected {
				t.Errorf("CelciusToFahrenheit(%v) = %v, want %v", tt.celcius, result, tt.expected)
			}
		})
	}
}

func TestCelciusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celcius  float64
		expected float64
	}{
		{"zero", 0, 273},
		{"boiling", 100, 373},
		{"negative", -273, 0},
		{"fractional", 28.5, 301.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelciusToKelvin(tt.celcius)
			if result != tt.expected {
				t.Errorf("CelciusToKelvin(%v) = %v, want %v", tt.celcius, result, tt.expected)
			}
		})
	}
}
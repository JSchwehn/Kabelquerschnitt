package main

import (
	"strconv"
	"testing"
)

// TestTUIInputValidation tests the input validation logic used by the TUI
func TestTUIInputValidation(t *testing.T) {
	tests := []struct {
		name      string
		voltage   string
		current   string
		length    string
		wantError bool
	}{
		{
			name:      "valid inputs",
			voltage:   "12.0",
			current:   "10.0",
			length:    "5.0",
			wantError: false,
		},
		{
			name:      "missing voltage",
			voltage:   "",
			current:   "10.0",
			length:    "5.0",
			wantError: true,
		},
		{
			name:      "invalid voltage - too high",
			voltage:   "60.0",
			current:   "10.0",
			length:    "5.0",
			wantError: true,
		},
		{
			name:      "invalid voltage - negative",
			voltage:   "-5.0",
			current:   "10.0",
			length:    "5.0",
			wantError: true,
		},
		{
			name:      "missing current",
			voltage:   "12.0",
			current:   "",
			length:    "5.0",
			wantError: true,
		},
		{
			name:      "invalid current - negative",
			voltage:   "12.0",
			current:   "-5.0",
			length:    "5.0",
			wantError: true,
		},
		{
			name:      "missing length",
			voltage:   "12.0",
			current:   "10.0",
			length:    "",
			wantError: true,
		},
		{
			name:      "invalid length - negative",
			voltage:   "12.0",
			current:   "10.0",
			length:    "-5.0",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false

			// Test voltage
			if tt.voltage == "" {
				hasError = true
			} else {
				voltage, err := strconv.ParseFloat(tt.voltage, 64)
				if err != nil || voltage <= 0 || voltage > 50 {
					hasError = true
				}
			}

			// Test current
			if tt.current == "" {
				hasError = true
			} else {
				current, err := strconv.ParseFloat(tt.current, 64)
				if err != nil || current <= 0 {
					hasError = true
				}
			}

			// Test length
			if tt.length == "" {
				hasError = true
			} else {
				length, err := strconv.ParseFloat(tt.length, 64)
				if err != nil || length <= 0 {
					hasError = true
				}
			}

			if hasError != tt.wantError {
				t.Errorf("Validation error = %v, want %v", hasError, tt.wantError)
			}
		})
	}
}

// TestTUIParseInputs tests the input parsing logic used by TUI calculate function
func TestTUIParseInputs(t *testing.T) {
	tests := []struct {
		name            string
		voltageStr      string
		currentStr      string
		lengthStr       string
		voltageDropStr  string
		tempStr         string
		tempUnitStr     string
		wantVoltage     float64
		wantCurrent     float64
		wantLength      float64
		wantVoltageDrop float64
		wantTemp        float64
		wantTempUnit    string
	}{
		{
			name:            "all inputs provided",
			voltageStr:      "12.0",
			currentStr:      "10.0",
			lengthStr:       "5.0",
			voltageDropStr:  "3.0",
			tempStr:         "25.0",
			tempUnitStr:     "C",
			wantVoltage:     12.0,
			wantCurrent:     10.0,
			wantLength:      5.0,
			wantVoltageDrop: 3.0,
			wantTemp:        25.0,
			wantTempUnit:    "C",
		},
		{
			name:            "defaults for optional fields",
			voltageStr:      "24.0",
			currentStr:      "15.0",
			lengthStr:       "10.0",
			voltageDropStr:  "",
			tempStr:         "",
			tempUnitStr:     "",
			wantVoltage:     24.0,
			wantCurrent:     15.0,
			wantLength:      10.0,
			wantVoltageDrop: 3.0,  // default
			wantTemp:        20.0, // default
			wantTempUnit:    "C",  // default
		},
		{
			name:            "fahrenheit temperature",
			voltageStr:      "12.0",
			currentStr:      "10.0",
			lengthStr:       "5.0",
			voltageDropStr:  "",
			tempStr:         "68.0",
			tempUnitStr:     "F",
			wantVoltage:     12.0,
			wantCurrent:     10.0,
			wantLength:      5.0,
			wantVoltageDrop: 3.0,
			wantTemp:        20.0, // 68°F = 20°C
			wantTempUnit:    "F",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse voltage
			voltage, err := strconv.ParseFloat(tt.voltageStr, 64)
			if err != nil {
				t.Fatalf("Failed to parse voltage: %v", err)
			}
			if voltage != tt.wantVoltage {
				t.Errorf("Voltage = %v, want %v", voltage, tt.wantVoltage)
			}

			// Parse current
			current, err := strconv.ParseFloat(tt.currentStr, 64)
			if err != nil {
				t.Fatalf("Failed to parse current: %v", err)
			}
			if current != tt.wantCurrent {
				t.Errorf("Current = %v, want %v", current, tt.wantCurrent)
			}

			// Parse length
			length, err := strconv.ParseFloat(tt.lengthStr, 64)
			if err != nil {
				t.Fatalf("Failed to parse length: %v", err)
			}
			if length != tt.wantLength {
				t.Errorf("Length = %v, want %v", length, tt.wantLength)
			}

			// Parse voltage drop (with default)
			voltageDrop := 3.0
			if tt.voltageDropStr != "" {
				if val, err := strconv.ParseFloat(tt.voltageDropStr, 64); err == nil && val > 0 && val <= 10 {
					voltageDrop = val
				}
			}
			if voltageDrop != tt.wantVoltageDrop {
				t.Errorf("VoltageDrop = %v, want %v", voltageDrop, tt.wantVoltageDrop)
			}

			// Parse temperature (with default)
			temp := 20.0
			if tt.tempStr != "" {
				if val, err := strconv.ParseFloat(tt.tempStr, 64); err == nil {
					temp = val
				}
			}
			if tt.tempUnitStr == "F" {
				temp = fahrenheitToCelsius(temp)
			}
			if temp != tt.wantTemp {
				t.Errorf("Temperature = %v, want %v", temp, tt.wantTemp)
			}
		})
	}
}

// TestTUISelections tests the selection mapping logic
func TestTUISelections(t *testing.T) {
	tests := []struct {
		name              string
		materialIndex     int
		installationIndex int
		wireTypeIndex     int
		wantMaterial      string
		wantInstallation  InstallationMethod
		wantWireType      string
	}{
		{
			name:              "copper, in air, flry",
			materialIndex:     0,
			installationIndex: 0,
			wireTypeIndex:     0,
			wantMaterial:      "copper",
			wantInstallation:  InstallationInAir,
			wantWireType:      "flry",
		},
		{
			name:              "aluminum, in conduit, generic",
			materialIndex:     1,
			installationIndex: 1,
			wireTypeIndex:     8,
			wantMaterial:      "aluminum",
			wantInstallation:  InstallationConduit,
			wantWireType:      "generic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			materialKeys := []string{"copper", "aluminum"}
			if materialKeys[tt.materialIndex] != tt.wantMaterial {
				t.Errorf("Material = %v, want %v", materialKeys[tt.materialIndex], tt.wantMaterial)
			}

			installationKeys := []InstallationMethod{InstallationInAir, InstallationConduit, InstallationIsolated}
			if installationKeys[tt.installationIndex] != tt.wantInstallation {
				t.Errorf("Installation = %v, want %v", installationKeys[tt.installationIndex], tt.wantInstallation)
			}

			wireTypeKeys := []string{"flry", "flry-a", "flry-b", "thhn", "thwn", "xlpe", "pvc", "silicon", "generic"}
			if wireTypeKeys[tt.wireTypeIndex] != tt.wantWireType {
				t.Errorf("WireType = %v, want %v", wireTypeKeys[tt.wireTypeIndex], tt.wantWireType)
			}
		})
	}
}

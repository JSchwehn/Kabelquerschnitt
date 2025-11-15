package main

import (
	"math"
	"testing"
)

func TestCalculateCableArea(t *testing.T) {
	tests := []struct {
		name                  string
		voltage               float64
		current               float64
		length                float64
		maxVoltageDropPercent float64
		material              CableMaterial
		roundTrip             bool
		ambientTempCelsius    float64
		installation          InstallationMethod
		want                  float64
		tolerance             float64
	}{
		{
			name:                  "12V system, 10A, 5m, 3% drop, copper, one-way, 20°C, in air",
			voltage:               12.0,
			current:               10.0,
			length:                5.0,
			maxVoltageDropPercent: 3.0,
			material:              materials["copper"],
			roundTrip:             false,
			ambientTempCelsius:    20.0,
			installation:          InstallationInAir,
			want:                  2.4305555555555554, // (10 * 0.0175 * 5 * 1) / (12 * 0.03)
			tolerance:             0.01,
		},
		{
			name:                  "12V system, 10A, 5m, 3% drop, copper, round trip, 20°C, in air",
			voltage:               12.0,
			current:               10.0,
			length:                5.0,
			maxVoltageDropPercent: 3.0,
			material:              materials["copper"],
			roundTrip:             true,
			ambientTempCelsius:    20.0,
			installation:          InstallationInAir,
			want:                  4.861111111111111, // (10 * 0.0175 * 5 * 2) / (12 * 0.03)
			tolerance:             0.01,
		},
		{
			name:                  "24V system, 20A, 10m, 5% drop, copper, one-way, 20°C, in air",
			voltage:               24.0,
			current:               20.0,
			length:                10.0,
			maxVoltageDropPercent: 5.0,
			material:              materials["copper"],
			roundTrip:             false,
			ambientTempCelsius:    20.0,
			installation:          InstallationInAir,
			want:                  2.9166666666666665, // (20 * 0.0175 * 10 * 1) / (24 * 0.05)
			tolerance:             0.01,
		},
		{
			name:                  "48V system, 15A, 20m, 3% drop, aluminum, one-way, 20°C, in air",
			voltage:               48.0,
			current:               15.0,
			length:                20.0,
			maxVoltageDropPercent: 3.0,
			material:              materials["aluminum"],
			roundTrip:             false,
			ambientTempCelsius:    20.0,
			installation:          InstallationInAir,
			want:                  5.902777777777778, // (15 * 0.0283 * 20 * 1) / (48 * 0.03)
			tolerance:             0.01,
		},
		{
			name:                  "50V system, 5A, 15m, 2% drop, copper, round trip, 20°C, in air",
			voltage:               50.0,
			current:               5.0,
			length:                15.0,
			maxVoltageDropPercent: 2.0,
			material:              materials["copper"],
			roundTrip:             true,
			ambientTempCelsius:    20.0,
			installation:          InstallationInAir,
			want:                  2.625, // (5 * 0.0175 * 15 * 2) / (50 * 0.02) = 2.625 / 1.0 = 2.625
			tolerance:             0.01,
		},
		{
			name:                  "12V system, 10A, 5m, 3% drop, copper, one-way, 40°C, in conduit",
			voltage:               12.0,
			current:               10.0,
			length:                5.0,
			maxVoltageDropPercent: 3.0,
			material:              materials["copper"],
			roundTrip:             false,
			ambientTempCelsius:    40.0,
			installation:          InstallationConduit,
			// Effective temp: 40 + 10 = 50°C
			// Resistivity at 50°C: 0.0175 * (1 + 0.00393 * (50-20)) = 0.0175 * 1.1179 = 0.01956325
			// Area: (10 * 0.01956325 * 5 * 1) / (12 * 0.03) = 0.9781625 / 0.36 = 2.717
			want:      2.717,
			tolerance: 0.01,
		},
		{
			name:                  "12V system, 10A, 5m, 3% drop, copper, one-way, 0°C, isolated",
			voltage:               12.0,
			current:               10.0,
			length:                5.0,
			maxVoltageDropPercent: 3.0,
			material:              materials["copper"],
			roundTrip:             false,
			ambientTempCelsius:    0.0,
			installation:          InstallationIsolated,
			// Effective temp: 0 + 20 = 20°C
			// Resistivity at 20°C: 0.0175 (same as reference)
			// Area: (10 * 0.0175 * 5 * 1) / (12 * 0.03) = 2.431
			want:      2.431,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCableArea(tt.voltage, tt.current, tt.length, tt.maxVoltageDropPercent, tt.material, tt.roundTrip, tt.ambientTempCelsius, tt.installation)
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("calculateCableArea() = %v, want %v (tolerance: %v)", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestAreaToDiameter(t *testing.T) {
	tests := []struct {
		name      string
		area      float64
		want      float64
		tolerance float64
	}{
		{
			name:      "1 mm² area",
			area:      1.0,
			want:      1.1283791670955126, // 2 * sqrt(1/π)
			tolerance: 0.0001,
		},
		{
			name:      "2.5 mm² area",
			area:      2.5,
			want:      1.7841241161527712, // 2 * sqrt(2.5/π)
			tolerance: 0.0001,
		},
		{
			name:      "10 mm² area",
			area:      10.0,
			want:      3.5682482323055424, // 2 * sqrt(10/π)
			tolerance: 0.0001,
		},
		{
			name:      "25 mm² area",
			area:      25.0,
			want:      5.641895835477563, // 2 * sqrt(25/π)
			tolerance: 0.0001,
		},
		{
			name:      "zero area",
			area:      0.0,
			want:      0.0,
			tolerance: 0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := areaToDiameter(tt.area)
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("areaToDiameter() = %v, want %v (tolerance: %v)", got, tt.want, tt.tolerance)
			}
		})
	}

	// Verify the formula is correct: area = π * (diameter/2)²
	t.Run("formula verification", func(t *testing.T) {
		testArea := 10.0
		diameter := areaToDiameter(testArea)
		calculatedArea := math.Pi * math.Pow(diameter/2, 2)
		if math.Abs(calculatedArea-testArea) > 0.0001 {
			t.Errorf("Formula verification failed: area %v -> diameter %v -> area %v", testArea, diameter, calculatedArea)
		}
	})
}

func TestFindClosestMetricSize(t *testing.T) {
	tests := []struct {
		name         string
		requiredArea float64
		wantSize     float64
		wantDiff     float64
		tolerance    float64
	}{
		{
			name:         "exact match - 2.5 mm²",
			requiredArea: 2.5,
			wantSize:     2.5,
			wantDiff:     0.0,
			tolerance:    0.0001,
		},
		{
			name:         "close to 1.5 mm²",
			requiredArea: 1.6,
			wantSize:     1.5,
			wantDiff:     0.1,
			tolerance:    0.0001,
		},
		{
			name:         "close to 4.0 mm²",
			requiredArea: 3.8,
			wantSize:     4.0,
			wantDiff:     0.2,
			tolerance:    0.0001,
		},
		{
			name:         "very small area",
			requiredArea: 0.3,
			wantSize:     0.5,
			wantDiff:     0.2,
			tolerance:    0.0001,
		},
		{
			name:         "large area",
			requiredArea: 200.0,
			wantSize:     185.0, // 200 is closer to 185 than 240
			wantDiff:     15.0,
			tolerance:    0.0001,
		},
		{
			name:         "between 6.0 and 10.0",
			requiredArea: 8.0,
			wantSize:     6.0, // 8.0 is closer to 6.0 (diff=2.0) than 10.0 (diff=2.0), but 6.0 comes first
			wantDiff:     2.0,
			tolerance:    0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSize, gotDiff := findClosestMetricSize(tt.requiredArea)
			if math.Abs(gotSize-tt.wantSize) > tt.tolerance {
				t.Errorf("findClosestMetricSize() size = %v, want %v", gotSize, tt.wantSize)
			}
			if math.Abs(gotDiff-tt.wantDiff) > tt.tolerance {
				t.Errorf("findClosestMetricSize() diff = %v, want %v", gotDiff, tt.wantDiff)
			}
		})
	}
}

func TestFindClosestAWG(t *testing.T) {
	tests := []struct {
		name         string
		requiredArea float64
		wantLabel    string
		wantArea     float64
		wantDiff     float64
		tolerance    float64
	}{
		{
			name:         "exact match - AWG 12",
			requiredArea: 3.309,
			wantLabel:    "12",
			wantArea:     3.309,
			wantDiff:     0.0,
			tolerance:    0.0001,
		},
		{
			name:         "close to AWG 14",
			requiredArea: 2.0,
			wantLabel:    "14",
			wantArea:     2.081,
			wantDiff:     0.081,
			tolerance:    0.0001,
		},
		{
			name:         "close to AWG 10",
			requiredArea: 5.5,
			wantLabel:    "10",
			wantArea:     5.261,
			wantDiff:     0.239,
			tolerance:    0.0001,
		},
		{
			name:         "very small area - AWG 18",
			requiredArea: 0.5,
			wantLabel:    "18",
			wantArea:     0.823,
			wantDiff:     0.323,
			tolerance:    0.0001,
		},
		{
			name:         "large area - AWG 4/0",
			requiredArea: 100.0,
			wantLabel:    "4/0",
			wantArea:     107.2,
			wantDiff:     7.2,
			tolerance:    0.0001,
		},
		{
			name:         "between AWG 1 and 1/0",
			requiredArea: 48.0,
			wantLabel:    "1/0",
			wantArea:     53.49,
			wantDiff:     5.49,
			tolerance:    0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLabel, gotArea, gotDiff := findClosestAWG(tt.requiredArea)
			if gotLabel != tt.wantLabel {
				t.Errorf("findClosestAWG() label = %v, want %v", gotLabel, tt.wantLabel)
			}
			if math.Abs(gotArea-tt.wantArea) > tt.tolerance {
				t.Errorf("findClosestAWG() area = %v, want %v", gotArea, tt.wantArea)
			}
			if math.Abs(gotDiff-tt.wantDiff) > tt.tolerance {
				t.Errorf("findClosestAWG() diff = %v, want %v", gotDiff, tt.wantDiff)
			}
		})
	}
}

func TestMaterials(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantName string
		wantRes  float64
	}{
		{
			name:     "copper material",
			key:      "copper",
			wantName: "Copper",
			wantRes:  copperResistivity20C,
		},
		{
			name:     "aluminum material",
			key:      "aluminum",
			wantName: "Aluminum",
			wantRes:  aluminumResistivity20C,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			material, ok := materials[tt.key]
			if !ok {
				t.Fatalf("Material %s not found", tt.key)
			}
			if material.Name != tt.wantName {
				t.Errorf("Material name = %v, want %v", material.Name, tt.wantName)
			}
			if material.Resistivity20C != tt.wantRes {
				t.Errorf("Material resistivity = %v, want %v", material.Resistivity20C, tt.wantRes)
			}
		})
	}
}

func TestTemperatureConversion(t *testing.T) {
	tests := []struct {
		name       string
		celsius    float64
		fahrenheit float64
		tolerance  float64
	}{
		{
			name:       "0°C = 32°F",
			celsius:    0.0,
			fahrenheit: 32.0,
			tolerance:  0.01,
		},
		{
			name:       "20°C = 68°F",
			celsius:    20.0,
			fahrenheit: 68.0,
			tolerance:  0.01,
		},
		{
			name:       "100°C = 212°F",
			celsius:    100.0,
			fahrenheit: 212.0,
			tolerance:  0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Celsius to Fahrenheit
			gotF := celsiusToFahrenheit(tt.celsius)
			if math.Abs(gotF-tt.fahrenheit) > tt.tolerance {
				t.Errorf("celsiusToFahrenheit(%v) = %v, want %v", tt.celsius, gotF, tt.fahrenheit)
			}

			// Test Fahrenheit to Celsius
			gotC := fahrenheitToCelsius(tt.fahrenheit)
			if math.Abs(gotC-tt.celsius) > tt.tolerance {
				t.Errorf("fahrenheitToCelsius(%v) = %v, want %v", tt.fahrenheit, gotC, tt.celsius)
			}
		})
	}
}

func TestResistivityAtTemp(t *testing.T) {
	tests := []struct {
		name      string
		material  CableMaterial
		tempC     float64
		want      float64
		tolerance float64
	}{
		{
			name:      "Copper at 20°C",
			material:  materials["copper"],
			tempC:     20.0,
			want:      0.0175,
			tolerance: 0.0001,
		},
		{
			name:      "Copper at 50°C",
			material:  materials["copper"],
			tempC:     50.0,
			want:      0.0175 * (1 + 0.00393*(50-20)), // 0.01956325
			tolerance: 0.0001,
		},
		{
			name:      "Aluminum at 20°C",
			material:  materials["aluminum"],
			tempC:     20.0,
			want:      0.0283,
			tolerance: 0.0001,
		},
		{
			name:      "Aluminum at 0°C",
			material:  materials["aluminum"],
			tempC:     0.0,
			want:      0.0283 * (1 + 0.00403*(0-20)), // Lower resistivity at lower temp
			tolerance: 0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateResistivityAtTemp(tt.material, tt.tempC)
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("calculateResistivityAtTemp() = %v, want %v (tolerance: %v)", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestEffectiveTemp(t *testing.T) {
	tests := []struct {
		name         string
		ambientTempC float64
		installation InstallationMethod
		want         float64
		tolerance    float64
	}{
		{
			name:         "20°C in air",
			ambientTempC: 20.0,
			installation: InstallationInAir,
			want:         20.0,
			tolerance:    0.01,
		},
		{
			name:         "20°C in conduit",
			ambientTempC: 20.0,
			installation: InstallationConduit,
			want:         30.0, // 20 + 10
			tolerance:    0.01,
		},
		{
			name:         "30°C isolated",
			ambientTempC: 30.0,
			installation: InstallationIsolated,
			want:         50.0, // 30 + 20
			tolerance:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateEffectiveTemp(tt.ambientTempC, tt.installation)
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("calculateEffectiveTemp() = %v, want %v (tolerance: %v)", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// Test a complete scenario: 12V system, 10A, 10m, 3% drop, copper, round trip
	voltage := 12.0
	current := 10.0
	length := 10.0
	maxVoltageDropPercent := 3.0
	material := materials["copper"]
	roundTrip := true

	// Calculate required area (using 20°C and in-air for compatibility)
	requiredArea := calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, roundTrip, 20.0, InstallationInAir)

	// Should be positive
	if requiredArea <= 0 {
		t.Errorf("Required area should be positive, got %v", requiredArea)
	}

	// Calculate diameter
	diameter := areaToDiameter(requiredArea)
	if diameter <= 0 {
		t.Errorf("Diameter should be positive, got %v", diameter)
	}

	// Find closest sizes
	closestMetric, metricDiff := findClosestMetricSize(requiredArea)
	closestAWG, awgArea, awgDiff := findClosestAWG(requiredArea)

	// Closest sizes should be positive
	if closestMetric <= 0 {
		t.Errorf("Closest metric size should be positive, got %v", closestMetric)
	}
	if closestAWG == "" {
		t.Errorf("Closest AWG label should not be empty")
	}
	if awgArea <= 0 {
		t.Errorf("Closest AWG area should be positive, got %v", awgArea)
	}

	// Differences should be reasonable
	if metricDiff < 0 {
		t.Errorf("Metric difference should be non-negative, got %v", metricDiff)
	}
	if awgDiff < 0 {
		t.Errorf("AWG difference should be non-negative, got %v", awgDiff)
	}

	// Verify voltage drop calculation with recommended sizes
	// For round trip: V_drop = I × ρ(T) × L × 2 / A
	distanceFactor := 2.0
	effectiveTemp := calculateEffectiveTemp(20.0, InstallationInAir)
	resistivity := calculateResistivityAtTemp(material, effectiveTemp)
	actualDropMetric := (current * resistivity * length * distanceFactor) / closestMetric
	actualDropAWG := (current * resistivity * length * distanceFactor) / awgArea

	// Voltage drops: if closest size is >= required, drop should be <= max
	// If closest size is < required (rounded down), drop may exceed max
	maxVoltageDrop := voltage * (maxVoltageDropPercent / 100.0)
	if closestMetric >= requiredArea {
		if actualDropMetric > maxVoltageDrop*1.05 { // Allow 5% tolerance
			t.Errorf("Actual voltage drop with metric size (%.2f V) exceeds maximum (%.2f V) for size >= required", actualDropMetric, maxVoltageDrop)
		}
	}
	if awgArea >= requiredArea {
		if actualDropAWG > maxVoltageDrop*1.05 { // Allow 5% tolerance
			t.Errorf("Actual voltage drop with AWG size (%.2f V) exceeds maximum (%.2f V) for size >= required", actualDropAWG, maxVoltageDrop)
		}
	}
	// If sizes are smaller than required, voltage drop exceeding max is expected
}

func TestEdgeCases(t *testing.T) {
	t.Run("very small current", func(t *testing.T) {
		area := calculateCableArea(12.0, 0.1, 5.0, 3.0, materials["copper"], false, 20.0, InstallationInAir)
		if area <= 0 {
			t.Errorf("Area should be positive for small current, got %v", area)
		}
	})

	t.Run("very short length", func(t *testing.T) {
		area := calculateCableArea(12.0, 10.0, 0.1, 3.0, materials["copper"], false, 20.0, InstallationInAir)
		if area <= 0 {
			t.Errorf("Area should be positive for short length, got %v", area)
		}
	})

	t.Run("maximum voltage", func(t *testing.T) {
		area := calculateCableArea(50.0, 10.0, 10.0, 3.0, materials["copper"], false, 20.0, InstallationInAir)
		if area <= 0 {
			t.Errorf("Area should be positive for 50V, got %v", area)
		}
	})

	t.Run("very large required area", func(t *testing.T) {
		// This would require a very large cable
		area := calculateCableArea(12.0, 100.0, 100.0, 1.0, materials["copper"], true, 20.0, InstallationInAir)
		closestMetric, _ := findClosestMetricSize(area)
		closestAWG, _, _ := findClosestAWG(area)

		// Should return the largest available sizes
		if closestMetric < 100.0 {
			t.Logf("Large area %v -> closest metric %v (may be limited by available sizes)", area, closestMetric)
		}
		if closestAWG == "" {
			t.Errorf("Should find an AWG size, got empty string")
		}
	})
}

func TestValidateWireTemperature(t *testing.T) {
	tests := []struct {
		name          string
		effectiveTemp float64
		wireType      WireType
		wantValid     bool
		wantWarning   bool // true if should have warning message
	}{
		{
			name:          "FLRY at safe temperature",
			effectiveTemp: 80.0,
			wireType:      wireTypes["flry"],
			wantValid:     true,
			wantWarning:   false,
		},
		{
			name:          "FLRY at warning temperature (90% of max)",
			effectiveTemp: 95.0, // 90% of 105°C
			wireType:      wireTypes["flry"],
			wantValid:     true,
			wantWarning:   true,
		},
		{
			name:          "FLRY exceeds maximum",
			effectiveTemp: 110.0,
			wireType:      wireTypes["flry"],
			wantValid:     false,
			wantWarning:   true,
		},
		{
			name:          "PVC at safe temperature",
			effectiveTemp: 50.0,
			wireType:      wireTypes["pvc"],
			wantValid:     true,
			wantWarning:   false,
		},
		{
			name:          "PVC exceeds maximum",
			effectiveTemp: 75.0,
			wireType:      wireTypes["pvc"],
			wantValid:     false,
			wantWarning:   true,
		},
		{
			name:          "Silicone at high but safe temperature",
			effectiveTemp: 180.0,
			wireType:      wireTypes["silicon"],
			wantValid:     true,
			wantWarning:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotMsg := ValidateWireTemperature(tt.effectiveTemp, tt.wireType)
			if gotValid != tt.wantValid {
				t.Errorf("ValidateWireTemperature() valid = %v, want %v", gotValid, tt.wantValid)
			}
			hasWarning := gotMsg != ""
			if hasWarning != tt.wantWarning {
				t.Errorf("ValidateWireTemperature() has warning = %v, want %v (message: %s)", hasWarning, tt.wantWarning, gotMsg)
			}
		})
	}
}

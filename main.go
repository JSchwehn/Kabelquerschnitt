package main

// DC Cable Diameter Calculator
//
// This code has been developed with AI assistance.
// While calculations and logic have been reviewed and tested,
// users should verify results for critical applications.
//
// See README.md for full disclaimer and safety warnings.

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	// Resistivity of copper at 20°C (Ω·mm²/m)
	// Value: 0.0175 Ω·mm²/m
	// See DEVELOPER.md for temperature compensation details
	copperResistivity20C = 0.0175

	// Resistivity of aluminum at 20°C (Ω·mm²/m)
	// Value: 0.0283 Ω·mm²/m
	aluminumResistivity20C = 0.0283

	// Temperature coefficient for copper (per °C)
	// Value: 0.00393 per °C (approximately 0.004/°C)
	copperTempCoefficient = 0.00393

	// Temperature coefficient for aluminum (per °C)
	// Value: 0.00403 per °C (approximately 0.004/°C)
	aluminumTempCoefficient = 0.00403

	// Reference temperature for resistivity values (°C)
	referenceTemp = 20.0
)

type CableMaterial struct {
	Name            string
	Resistivity20C  float64
	TempCoefficient float64
}

var materials = map[string]CableMaterial{
	"copper":   {"Copper", copperResistivity20C, copperTempCoefficient},
	"aluminum": {"Aluminum", aluminumResistivity20C, aluminumTempCoefficient},
}

// InstallationMethod represents how the cable is installed
type InstallationMethod string

const (
	InstallationInAir    InstallationMethod = "air"
	InstallationConduit  InstallationMethod = "conduit"
	InstallationIsolated InstallationMethod = "isolated"
)

// Temperature adjustment factors for installation methods
// These represent the temperature rise above ambient due to installation method
// Values are approximate temperature increases in °C
var installationTempAdjustments = map[InstallationMethod]float64{
	InstallationInAir:    0.0,  // Good cooling, minimal temperature rise
	InstallationConduit:  10.0, // Reduced cooling, moderate temperature rise
	InstallationIsolated: 20.0, // Poor cooling, significant temperature rise
}

// WireType represents the wire/cable type with its maximum temperature rating
type WireType struct {
	Name           string
	MaxTempCelsius float64
	Description    string
}

// Common wire types with their maximum operating temperatures
var wireTypes = map[string]WireType{
	"flry": {
		Name:           "FLRY",
		MaxTempCelsius: 105.0,
		Description:    "Automotive thin-wall PVC (FLRY-A/B), stranded copper",
	},
	"flry-a": {
		Name:           "FLRY-A",
		MaxTempCelsius: 105.0,
		Description:    "Automotive thin-wall PVC, flexible stranded",
	},
	"flry-b": {
		Name:           "FLRY-B",
		MaxTempCelsius: 105.0,
		Description:    "Automotive thin-wall PVC, symmetrical stranded",
	},
	"thhn": {
		Name:           "THHN",
		MaxTempCelsius: 90.0,
		Description:    "Thermoplastic, high heat, nylon coated",
	},
	"thwn": {
		Name:           "THWN",
		MaxTempCelsius: 75.0,
		Description:    "Thermoplastic, heat/water resistant, nylon coated",
	},
	"xlpe": {
		Name:           "XLPE",
		MaxTempCelsius: 90.0,
		Description:    "Cross-linked polyethylene insulation",
	},
	"pvc": {
		Name:           "PVC",
		MaxTempCelsius: 70.0,
		Description:    "Standard PVC insulation",
	},
	"silicon": {
		Name:           "Silicone",
		MaxTempCelsius: 200.0,
		Description:    "Silicone rubber insulation, high temperature",
	},
	"generic": {
		Name:           "Generic",
		MaxTempCelsius: 90.0,
		Description:    "Generic wire type (assumes 90°C rating)",
	},
}

// AWG size structure
type AWGSize struct {
	Label string
	Area  float64
}

// Standard AWG to mm² conversion (common sizes)
var awgSizes = []AWGSize{
	{Label: "18", Area: 0.823},
	{Label: "16", Area: 1.309},
	{Label: "14", Area: 2.081},
	{Label: "12", Area: 3.309},
	{Label: "10", Area: 5.261},
	{Label: "8", Area: 8.367},
	{Label: "6", Area: 13.30},
	{Label: "4", Area: 21.15},
	{Label: "2", Area: 33.62},
	{Label: "1", Area: 42.41},
	{Label: "1/0", Area: 53.49},
	{Label: "2/0", Area: 67.43},
	{Label: "3/0", Area: 85.01},
	{Label: "4/0", Area: 107.2},
}

// Standard metric cable sizes (mm²)
var standardMetricSizes = []float64{
	0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, 16.0, 25.0, 35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0,
}

// Calculate resistivity at given temperature.
//
// Formula: ρ(T) = ρ(20°C) × [1 + α × (T - 20)]
// Where:
//   - ρ(T) = resistivity at temperature T
//   - ρ(20°C) = resistivity at 20°C
//   - α = temperature coefficient (per °C)
//   - T = temperature in Celsius
//
// See DEVELOPER.md for detailed calculation methodology.
func calculateResistivityAtTemp(material CableMaterial, tempCelsius float64) float64 {
	return material.Resistivity20C * (1 + material.TempCoefficient*(tempCelsius-referenceTemp))
}

// Convert Fahrenheit to Celsius
func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

// Convert Celsius to Fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

// Calculate effective operating temperature considering installation method.
//
// The effective temperature accounts for ambient temperature plus
// temperature rise due to installation method (poor cooling in conduits/isolated).
func calculateEffectiveTemp(ambientTempCelsius float64, installation InstallationMethod) float64 {
	adjustment := installationTempAdjustments[installation]
	return ambientTempCelsius + adjustment
}

// ValidateWireTemperature checks if the effective operating temperature
// exceeds the wire type's maximum temperature rating.
//
// Returns true if temperature is within limits, false if exceeded.
// Also returns a warning message if temperature is close to the limit (>90% of max).
func ValidateWireTemperature(effectiveTempCelsius float64, wireType WireType) (bool, string) {
	if effectiveTempCelsius > wireType.MaxTempCelsius {
		return false, fmt.Sprintf("WARNING: Effective operating temperature (%.1f°C) exceeds %s maximum rating (%.0f°C)! Wire insulation may fail.", effectiveTempCelsius, wireType.Name, wireType.MaxTempCelsius)
	}

	// Warn if within 10% of maximum
	if effectiveTempCelsius > wireType.MaxTempCelsius*0.9 {
		return true, fmt.Sprintf("CAUTION: Effective operating temperature (%.1f°C) is close to %s maximum rating (%.0f°C). Consider using a higher temperature rated wire.", effectiveTempCelsius, wireType.Name, wireType.MaxTempCelsius)
	}

	return true, ""
}

// Calculate required cross-sectional area based on voltage drop.
//
// Formula: A = (I × ρ(T) × L × distanceFactor) / V_drop_max
// Where:
//   - A = cross-sectional area (mm²)
//   - I = current (A)
//   - ρ(T) = material resistivity at operating temperature (Ω·mm²/m)
//   - L = cable length (m)
//   - distanceFactor = 1.0 for one-way, 2.0 for round trip
//   - V_drop_max = V_system × (maxVoltageDropPercent / 100)
//   - T = effective operating temperature (ambient + installation adjustment)
//
// For round trip, the factor is 2 because current flows through both
// the positive and return conductors.
//
// See DEVELOPER.md for detailed calculation methodology.
func calculateCableArea(voltage, current, length, maxVoltageDropPercent float64, material CableMaterial, roundTrip bool, ambientTempCelsius float64, installation InstallationMethod) float64 {
	maxVoltageDrop := voltage * (maxVoltageDropPercent / 100.0)

	distanceFactor := 1.0
	if roundTrip {
		distanceFactor = 2.0
	}

	// Calculate effective operating temperature
	effectiveTemp := calculateEffectiveTemp(ambientTempCelsius, installation)

	// Calculate resistivity at operating temperature
	resistivity := calculateResistivityAtTemp(material, effectiveTemp)

	area := (current * resistivity * length * distanceFactor) / maxVoltageDrop

	return area
}

// Calculate diameter from cross-sectional area.
//
// Formula: diameter = 2 × √(area / π)
// Assumes circular cross-section.
//
// See DEVELOPER.md for calculation details.
func areaToDiameter(area float64) float64 {
	return 2 * math.Sqrt(area/math.Pi)
}

// Find closest standard metric cable size.
//
// Returns the standard metric size (mm²) closest to the required area
// and the absolute difference between them.
//
// Standard sizes: 0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, 16.0, 25.0,
// 35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0 mm²
func findClosestMetricSize(requiredArea float64) (float64, float64) {
	var closestSize float64
	minDiff := math.MaxFloat64

	for _, size := range standardMetricSizes {
		diff := math.Abs(size - requiredArea)
		if diff < minDiff {
			minDiff = diff
			closestSize = size
		}
	}

	return closestSize, minDiff
}

// Find closest AWG (American Wire Gauge) size.
//
// Returns the AWG label (e.g., "12", "1/0", "2/0"), the cross-sectional
// area of that AWG size, and the absolute difference from the required area.
//
// Supported AWG sizes: 18, 16, 14, 12, 10, 8, 6, 4, 2, 1, 1/0, 2/0, 3/0, 4/0
func findClosestAWG(requiredArea float64) (string, float64, float64) {
	var closestLabel string
	var closestArea float64
	minDiff := math.MaxFloat64

	for _, awg := range awgSizes {
		diff := math.Abs(awg.Area - requiredArea)
		if diff < minDiff {
			minDiff = diff
			closestLabel = awg.Label
			closestArea = awg.Area
		}
	}

	return closestLabel, closestArea, minDiff
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== DC Cable Diameter Calculator ===")
	fmt.Println("Supports 12V, 24V, 48V, 50V DC systems")
	fmt.Println()

	// Get system voltage
	fmt.Print("Enter system voltage (V): ")
	voltageStr, _ := reader.ReadString('\n')
	voltageStr = strings.TrimSpace(voltageStr)
	voltage, err := strconv.ParseFloat(voltageStr, 64)
	if err != nil || voltage <= 0 || voltage > 50 {
		fmt.Println("Error: Invalid voltage. Please enter a value between 0 and 50V (inclusive).")
		return
	}

	// Get current
	fmt.Print("Enter current (A): ")
	currentStr, _ := reader.ReadString('\n')
	currentStr = strings.TrimSpace(currentStr)
	current, err := strconv.ParseFloat(currentStr, 64)
	if err != nil || current <= 0 {
		fmt.Println("Error: Invalid current. Please enter a positive value.")
		return
	}

	// Get cable length
	fmt.Print("Enter cable length (m): ")
	lengthStr, _ := reader.ReadString('\n')
	lengthStr = strings.TrimSpace(lengthStr)
	length, err := strconv.ParseFloat(lengthStr, 64)
	if err != nil || length <= 0 {
		fmt.Println("Error: Invalid length. Please enter a positive value.")
		return
	}

	// Get voltage drop percentage
	fmt.Print("Enter maximum voltage drop percentage (default 3%): ")
	dropStr, _ := reader.ReadString('\n')
	dropStr = strings.TrimSpace(dropStr)
	maxVoltageDropPercent := 3.0
	if dropStr != "" {
		maxVoltageDropPercent, err = strconv.ParseFloat(dropStr, 64)
		if err != nil || maxVoltageDropPercent <= 0 || maxVoltageDropPercent > 10 {
			fmt.Println("Warning: Invalid voltage drop percentage. Using default 3%.")
			maxVoltageDropPercent = 3.0
		}
	}

	// Get round trip option
	fmt.Print("Is this round trip length? (y/n, default: n): ")
	roundTripStr, _ := reader.ReadString('\n')
	roundTripStr = strings.TrimSpace(strings.ToLower(roundTripStr))
	roundTrip := roundTripStr == "y" || roundTripStr == "yes"

	// Get material
	fmt.Print("Cable material (copper/aluminum, default: copper): ")
	materialStr, _ := reader.ReadString('\n')
	materialStr = strings.TrimSpace(strings.ToLower(materialStr))
	material, ok := materials[materialStr]
	if !ok {
		material = materials["copper"]
		fmt.Println("Using default: Copper")
	}

	// Get temperature
	fmt.Print("Temperature unit (C/F, default: C): ")
	tempUnitStr, _ := reader.ReadString('\n')
	tempUnitStr = strings.TrimSpace(strings.ToUpper(tempUnitStr))
	if tempUnitStr != "F" && tempUnitStr != "C" && tempUnitStr != "" {
		tempUnitStr = "C"
	}
	if tempUnitStr == "" {
		tempUnitStr = "C"
	}

	fmt.Print("Enter ambient temperature: ")
	tempStr, _ := reader.ReadString('\n')
	tempStr = strings.TrimSpace(tempStr)
	ambientTemp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		fmt.Println("Error: Invalid temperature. Using default 20°C.")
		ambientTemp = 20.0
		tempUnitStr = "C"
	}

	// Convert to Celsius if needed
	ambientTempCelsius := ambientTemp
	if tempUnitStr == "F" {
		ambientTempCelsius = fahrenheitToCelsius(ambientTemp)
	}

	// Get installation method
	fmt.Print("Installation method (air/conduit/isolated, default: air): ")
	installStr, _ := reader.ReadString('\n')
	installStr = strings.TrimSpace(strings.ToLower(installStr))
	var installation InstallationMethod
	switch installStr {
	case "conduit":
		installation = InstallationConduit
	case "isolated":
		installation = InstallationIsolated
	case "air", "":
		installation = InstallationInAir
	default:
		installation = InstallationInAir
		fmt.Println("Using default: In air")
	}

	// Get wire type
	fmt.Print("Wire type (flry/flry-a/flry-b/thhn/thwn/xlpe/pvc/silicon/generic, default: generic): ")
	wireTypeStr, _ := reader.ReadString('\n')
	wireTypeStr = strings.TrimSpace(strings.ToLower(wireTypeStr))
	wireType, ok := wireTypes[wireTypeStr]
	if !ok {
		wireType = wireTypes["generic"]
		fmt.Println("Using default: Generic (90°C)")
	}

	fmt.Println()
	fmt.Println("=== Calculation Results ===")
	fmt.Printf("System Voltage: %.1f V\n", voltage)
	fmt.Printf("Current: %.2f A\n", current)
	fmt.Printf("Cable Length: %.2f m (%s)\n", length, map[bool]string{true: "round trip", false: "one-way"}[roundTrip])
	fmt.Printf("Material: %s\n", material.Name)
	fmt.Printf("Wire Type: %s (Max: %.0f°C) - %s\n", wireType.Name, wireType.MaxTempCelsius, wireType.Description)
	fmt.Printf("Ambient Temperature: %.1f°%s (%.1f°C)\n", ambientTemp, tempUnitStr, ambientTempCelsius)
	fmt.Printf("Installation Method: %s\n", map[InstallationMethod]string{
		InstallationInAir:    "In air",
		InstallationConduit:  "In conduit",
		InstallationIsolated: "Isolated/Insulated",
	}[installation])

	effectiveTemp := calculateEffectiveTemp(ambientTempCelsius, installation)
	fmt.Printf("Effective Operating Temperature: %.1f°C\n", effectiveTemp)

	// Validate wire temperature rating
	isValid, warningMsg := ValidateWireTemperature(effectiveTemp, wireType)
	if !isValid {
		fmt.Println()
		fmt.Println("⚠️  " + warningMsg)
		fmt.Println("   The calculated cable size may not be safe for this wire type!")
		fmt.Println("   Consider: using a higher temperature rated wire, reducing ambient temperature,")
		fmt.Println("   improving cooling, or increasing cable size to reduce heat generation.")
		fmt.Println()
	} else if warningMsg != "" {
		fmt.Println()
		fmt.Println("⚠️  " + warningMsg)
		fmt.Println()
	}

	fmt.Printf("Maximum Voltage Drop: %.2f%% (%.2f V)\n", maxVoltageDropPercent, voltage*maxVoltageDropPercent/100)
	fmt.Println()

	// Calculate required area
	requiredArea := calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, roundTrip, ambientTempCelsius, installation)
	requiredDiameter := areaToDiameter(requiredArea)

	fmt.Printf("Required Cross-Sectional Area: %.2f mm²\n", requiredArea)
	fmt.Printf("Required Diameter: %.2f mm\n", requiredDiameter)
	fmt.Println()

	// Find standard sizes
	closestMetric, metricDiff := findClosestMetricSize(requiredArea)
	closestAWG, awgArea, awgDiff := findClosestAWG(requiredArea)

	fmt.Println("=== Recommended Standard Sizes ===")
	fmt.Printf("Metric: %.2f mm² (difference: %.2f mm²)\n", closestMetric, metricDiff)
	fmt.Printf("AWG: %s (%.2f mm², difference: %.2f mm²)\n", closestAWG, awgArea, awgDiff)
	fmt.Println()

	// Calculate actual voltage drop with recommended sizes
	fmt.Println("=== Voltage Drop with Recommended Sizes ===")

	// Calculate resistivity at operating temperature for voltage drop calculations
	resistivity := calculateResistivityAtTemp(material, effectiveTemp)
	distanceFactor := map[bool]float64{true: 2.0, false: 1.0}[roundTrip]

	// Metric size
	actualDropMetric := (current * resistivity * length * distanceFactor) / closestMetric
	actualDropPercentMetric := (actualDropMetric / voltage) * 100
	fmt.Printf("With %.2f mm²: %.2f V (%.2f%%)\n", closestMetric, actualDropMetric, actualDropPercentMetric)

	// AWG size
	actualDropAWG := (current * resistivity * length * distanceFactor) / awgArea
	actualDropPercentAWG := (actualDropAWG / voltage) * 100
	fmt.Printf("With AWG %s (%.2f mm²): %.2f V (%.2f%%)\n", closestAWG, awgArea, actualDropAWG, actualDropPercentAWG)
}

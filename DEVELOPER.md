# Developer Documentation

This document provides technical details for developers working on the DC Cable Diameter Calculator.

## Table of Contents

- [Architecture](#architecture)
- [Calculation Methodology](#calculation-methodology)
- [Code Structure](#code-structure)
- [Key Functions](#key-functions)
- [Testing](#testing)
- [Adding New Features](#adding-new-features)
- [Maintaining Documentation](#maintaining-documentation)

## Architecture

The application is a simple command-line tool written in Go. It follows a straightforward structure:

```
Kabelquerschnitt/
├── main.go          # Main application code
├── main_test.go     # Test suite
├── go.mod           # Go module definition
├── README.md        # User documentation
└── DEVELOPER.md     # This file
```

## Calculation Methodology

### Voltage Drop Formula

The core calculation is based on Ohm's law and the resistance formula for conductors.

#### Basic Principles

1. **Resistance of a conductor:**
   ```
   R = ρ × L / A
   ```
   Where:
   - `R` = Resistance (Ω)
   - `ρ` = Resistivity of the material (Ω·mm²/m)
   - `L` = Length (m)
   - `A` = Cross-sectional area (mm²)

2. **Voltage drop (Ohm's law):**
   ```
   V_drop = I × R
   ```
   Where:
   - `V_drop` = Voltage drop (V)
   - `I` = Current (A)
   - `R` = Resistance (Ω)

#### Combined Formula

For a DC system, the voltage drop over a cable is:

**One-way (single conductor):**
```
V_drop = I × ρ(T) × L / A
```

**Round trip (power + return conductor):**
```
V_drop = I × ρ(T) × L × 2 / A
```

Where `ρ(T)` is the resistivity at the effective operating temperature.

#### Solving for Cross-Sectional Area

To find the required cross-sectional area `A` for a given maximum voltage drop:

**One-way:**
```
A = (I × ρ(T) × L) / V_drop_max
```

**Round trip:**
```
A = (I × ρ(T) × L × 2) / V_drop_max
```

Where:
- `V_drop_max = V_system × (voltage_drop_percent / 100)`
- `ρ(T) = ρ(20°C) × [1 + α × (T_effective - 20)]`
- `T_effective = T_ambient + installation_adjustment`

### Implementation

The calculation is implemented in `calculateCableArea()`:

```go
func calculateCableArea(voltage, current, length, maxVoltageDropPercent float64, 
                       material CableMaterial, roundTrip bool, 
                       ambientTempCelsius float64, installation InstallationMethod) float64 {
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
```

#### Supporting Functions

**Temperature conversion:**
```go
func fahrenheitToCelsius(f float64) float64 {
    return (f - 32) * 5 / 9
}

func celsiusToFahrenheit(c float64) float64 {
    return c*9/5 + 32
}
```

**Resistivity at temperature:**
```go
func calculateResistivityAtTemp(material CableMaterial, tempCelsius float64) float64 {
    return material.Resistivity20C * (1 + material.TempCoefficient*(tempCelsius-referenceTemp))
}
```

**Effective temperature:**
```go
func calculateEffectiveTemp(ambientTempCelsius float64, installation InstallationMethod) float64 {
    adjustment := installationTempAdjustments[installation]
    return ambientTempCelsius + adjustment
}
```

### Diameter Calculation

From cross-sectional area to diameter (assuming circular cross-section):

```
diameter = 2 × √(area / π)
```

Implemented in `areaToDiameter()`:

```go
func areaToDiameter(area float64) float64 {
    return 2 * math.Sqrt(area/math.Pi)
}
```

### Material Resistivity

The program uses standard resistivity values at 20°C:

- **Copper**: 0.0175 Ω·mm²/m (temperature coefficient: 0.00393 per °C)
- **Aluminum**: 0.0283 Ω·mm²/m (temperature coefficient: 0.00403 per °C)

#### Temperature Compensation

Resistivity changes with temperature according to:

```
ρ(T) = ρ(20°C) × [1 + α × (T - 20)]
```

Where:
- `ρ(T)` = resistivity at temperature T
- `ρ(20°C)` = resistivity at 20°C
- `α` = temperature coefficient (per °C)
- `T` = temperature in Celsius

#### Effective Operating Temperature

The effective operating temperature accounts for:
1. **Ambient temperature**: The environmental temperature
2. **Installation method adjustment**: Additional temperature rise due to reduced cooling
   - **In air**: +0°C (good cooling)
   - **In conduit**: +10°C (reduced cooling)
   - **Isolated/Insulated**: +20°C (poor cooling)

```
T_effective = T_ambient + installation_adjustment
```

The resistivity is then calculated at the effective operating temperature.

## Code Structure

### Constants

```go
const (
    copperResistivity20C   = 0.0175  // Ω·mm²/m at 20°C
    aluminumResistivity20C = 0.0283  // Ω·mm²/m at 20°C
    copperTempCoefficient   = 0.00393 // per °C
    aluminumTempCoefficient = 0.00403 // per °C
    referenceTemp           = 20.0    // °C
)
```

### Types

#### CableMaterial
```go
type CableMaterial struct {
    Name            string
    Resistivity20C  float64  // Resistivity at 20°C
    TempCoefficient float64  // Temperature coefficient per °C
}
```

#### InstallationMethod
```go
type InstallationMethod string

const (
    InstallationInAir     InstallationMethod = "air"
    InstallationConduit   InstallationMethod = "conduit"
    InstallationIsolated  InstallationMethod = "isolated"
)

var installationTempAdjustments = map[InstallationMethod]float64{
    InstallationInAir:     0.0,  // Good cooling
    InstallationConduit:   10.0, // Reduced cooling
    InstallationIsolated:  20.0, // Poor cooling
}
```

#### AWGSize
```go
type AWGSize struct {
    Label string
    Area  float64
}
```

### Data Structures

#### Standard Cable Sizes

**Metric sizes (mm²):**
```go
var standardMetricSizes = []float64{
    0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, 16.0, 25.0, 
    35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0,
}
```

**AWG sizes:**
```go
var awgSizes = []AWGSize{
    {Label: "18", Area: 0.823},
    {Label: "16", Area: 1.309},
    // ... up to 4/0
}
```

## Key Functions

### calculateCableArea()

Calculates the required cross-sectional area based on voltage drop requirements, accounting for temperature effects.

**Parameters:**
- `voltage`: System voltage in volts
- `current`: Current in amperes
- `length`: Cable length in meters
- `maxVoltageDropPercent`: Maximum allowed voltage drop as percentage
- `material`: Cable material (copper or aluminum)
- `roundTrip`: Whether length is round trip (true) or one-way (false)
- `ambientTempCelsius`: Ambient temperature in Celsius
- `installation`: Installation method (air, conduit, or isolated)

**Returns:**
- Cross-sectional area in mm²

**Formula:**
```
A = (I × ρ(T_effective) × L × distanceFactor) / (V × maxDropPercent / 100)
```

Where:
- `T_effective = T_ambient + installation_adjustment`
- `ρ(T) = ρ(20°C) × [1 + α × (T - 20)]`

### calculateResistivityAtTemp()

Calculates material resistivity at a given temperature.

**Parameters:**
- `material`: Cable material with resistivity and temperature coefficient
- `tempCelsius`: Temperature in Celsius

**Returns:**
- Resistivity at the specified temperature (Ω·mm²/m)

**Formula:**
```
ρ(T) = ρ(20°C) × [1 + α × (T - 20)]
```

### calculateEffectiveTemp()

Calculates the effective operating temperature considering installation method.

**Parameters:**
- `ambientTempCelsius`: Ambient temperature in Celsius
- `installation`: Installation method (air, conduit, or isolated)

**Returns:**
- Effective operating temperature in Celsius

**Formula:**
```
T_effective = T_ambient + installation_adjustment
```

### fahrenheitToCelsius() / celsiusToFahrenheit()

Temperature conversion utilities.

**Parameters:**
- `f`: Temperature in Fahrenheit (for fahrenheitToCelsius)
- `c`: Temperature in Celsius (for celsiusToFahrenheit)

**Returns:**
- Converted temperature

**Formulas:**
```
°C = (°F - 32) × 5/9
°F = °C × 9/5 + 32
```

### areaToDiameter()

Converts cross-sectional area to diameter.

**Parameters:**
- `area`: Cross-sectional area in mm²

**Returns:**
- Diameter in mm

**Formula:**
```
diameter = 2 × √(area / π)
```

### findClosestMetricSize()

Finds the closest standard metric cable size.

**Parameters:**
- `requiredArea`: Required cross-sectional area in mm²

**Returns:**
- `closestSize`: Closest standard size in mm²
- `diff`: Difference between required and closest size

**Algorithm:**
Iterates through all standard sizes and finds the one with minimum absolute difference.

### findClosestAWG()

Finds the closest AWG cable size.

**Parameters:**
- `requiredArea`: Required cross-sectional area in mm²

**Returns:**
- `label`: AWG label (e.g., "12", "1/0", "2/0")
- `area`: Cross-sectional area of closest AWG size
- `diff`: Difference between required and closest size

## Testing

### Running Tests

```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Run with coverage
go test -cover

# Run with race detector
go test -race
```

### Test Coverage

The test suite (`main_test.go`) includes:

1. **TestCalculateCableArea**: Tests cable area calculations with various scenarios
2. **TestAreaToDiameter**: Tests diameter calculations and formula verification
3. **TestFindClosestMetricSize**: Tests metric size selection
4. **TestFindClosestAWG**: Tests AWG size selection
5. **TestMaterials**: Tests material properties
6. **TestIntegration**: End-to-end integration test
7. **TestEdgeCases**: Edge case handling

### Test Data

Tests use realistic scenarios:
- 12V, 24V, 48V, 50V systems
- Various current levels (5A to 20A)
- Different cable lengths (5m to 20m)
- One-way and round trip scenarios
- Copper and aluminum materials

## Adding New Features

### Adding a New Material

1. Add resistivity constant:
```go
const newMaterialResistivity = 0.XXXX  // Ω·mm²/m
```

2. Add to materials map:
```go
var materials = map[string]CableMaterial{
    "copper":   {"Copper", copperResistivity},
    "aluminum": {"Aluminum", aluminumResistivity},
    "newmaterial": {"New Material", newMaterialResistivity},
}
```

3. Update user documentation (README.md)

### Adding New Cable Sizes

**For metric sizes:**
Add to `standardMetricSizes` array in ascending order.

**For AWG sizes:**
Add to `awgSizes` array with proper label and area.

### Adding Temperature Compensation

To add temperature compensation:

1. Add temperature parameter to `calculateCableArea()`
2. Apply temperature coefficient:
   ```
   ρ(T) = ρ(20°C) × [1 + α × (T - 20)]
   ```
   Where α ≈ 0.004/°C for copper

3. Update function signature and documentation

### Adding Current-Carrying Capacity

To add ampacity (current-carrying capacity) checks:

1. Create ampacity lookup tables (depends on installation method, ambient temperature)
2. Add validation after calculating required area
3. Warn if selected cable cannot handle the current

## Maintaining Documentation

### When to Update Documentation

Update documentation when:

1. **Adding new features:**
   - Update README.md with new usage instructions
   - Update DEVELOPER.md with implementation details

2. **Changing calculations:**
   - Update DEVELOPER.md calculation methodology section
   - Update any affected examples

3. **Adding new materials or cable sizes:**
   - Update README.md usage section
   - Update DEVELOPER.md data structures section

4. **Changing function signatures:**
   - Update DEVELOPER.md key functions section
   - Update code comments

### Documentation Standards

- **README.md**: User-facing, non-technical language, examples
- **DEVELOPER.md**: Technical details, formulas, implementation notes
- **Code comments**: Inline documentation for complex logic

### Keeping Documentation in Sync

1. **Code changes trigger doc updates:**
   - Function signature changes → Update DEVELOPER.md
   - New features → Update both README.md and DEVELOPER.md
   - UI/UX changes → Update README.md examples

2. **Regular review:**
   - Review documentation when reviewing code
   - Verify examples still work
   - Check that formulas match implementation

3. **Version control:**
   - Commit documentation changes with code changes
   - Use meaningful commit messages

## Code Quality

### Linting

```bash
# Run go vet
go vet

# Run golangci-lint (if installed)
golangci-lint run
```

### Formatting

```bash
# Format code
go fmt ./...

# Check formatting
go fmt -d .
```

## Future Enhancements

Potential improvements:

1. **Command-line arguments**: Support non-interactive mode
2. **Temperature compensation**: Account for operating temperature
3. **Ampacity checking**: Verify current-carrying capacity
4. **Multiple calculations**: Batch processing capability
5. **Export results**: Save results to file (CSV, JSON)
6. **GUI version**: Web or desktop interface
7. **More materials**: Add silver, gold, other conductors
8. **Installation method**: Account for cable installation (free air, conduit, etc.)

## References

- IEC 60228: International standard for conductor sizes
- AWG (American Wire Gauge) standard
- Electrical engineering handbooks for resistivity values
- National Electrical Code (NEC) for ampacity tables

## Contributing

When contributing:

1. Follow Go best practices
2. Add tests for new features
3. Update documentation
4. Ensure all tests pass
5. Format code with `go fmt`

## Questions or Issues?

For technical questions or to report issues, please refer to the project's issue tracker or contact the maintainers.


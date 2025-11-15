# DC Cable Diameter Calculator

A command-line tool for calculating the required cable cross-sectional area and diameter for DC electrical systems (12V, 24V, 48V, up to 50V). The calculator determines the appropriate cable size based on voltage drop requirements.

> **Note:** This code has been developed with AI assistance. While the calculations and logic have been reviewed, users should verify results for critical applications.

## Features

- ✅ Supports DC systems from 12V to 50V
- ✅ Calculates required cable cross-sectional area and diameter
- ✅ Supports both copper and aluminum cables
- ✅ Handles one-way and round-trip cable lengths
- ✅ Provides recommendations in both metric (mm²) and AWG sizes
- ✅ Shows actual voltage drop with recommended cable sizes
- ✅ **Interactive TUI mode** with step-by-step form interface
- ✅ Temperature compensation (Celsius/Fahrenheit)
- ✅ Installation method support (air, conduit, isolated)
- ✅ Wire type selection with temperature validation

## Installation

### Prerequisites

- Go 1.16 or later

### Build from Source

```bash
git clone <repository-url>
cd Kabelquerschnitt
go build -o cablecalc .
```

Or run directly:

```bash
go run .
```

**Note:** The application consists of `main.go` and `main_tui.go`. Both files must be included when building.

## Usage

The application supports two modes: **CLI mode** (default) and **TUI mode** (interactive terminal interface).

### TUI Mode (Recommended)

Run with the `--tui` or `-t` flag for an interactive terminal user interface:

```bash
./cablecalc --tui
# or
./cablecalc -t
```

The TUI provides:
- **Step-by-step form** with field navigation (Tab/Enter to move forward, Shift+Tab/Up to go back)
- **Interactive lists** for material, installation method, and wire type selection
- **Real-time validation** with error messages
- **Color-coded results** with warnings and recommendations
- **Keyboard shortcuts**: `q` to quit, `Esc` to go back, `r` to restart after viewing results

### CLI Mode

Run without flags for the traditional command-line interface:

```bash
./cablecalc
```

The program will prompt you for the following information:

1. **System Voltage (V)**: Enter the DC system voltage (e.g., 12, 24, 48, 50)
2. **Current (A)**: Enter the current in amperes
3. **Cable Length (m)**: Enter the cable length in meters
4. **Maximum Voltage Drop Percentage**: Enter the maximum allowed voltage drop (default: 3%)
5. **Round Trip Length**: Answer 'y' if the length is round trip (power + return), 'n' for one-way
6. **Cable Material**: Enter 'copper' or 'aluminum' (default: copper)
7. **Temperature Unit**: Enter 'C' for Celsius or 'F' for Fahrenheit (default: C)
8. **Ambient Temperature**: Enter the ambient/environment temperature
9. **Installation Method**: Enter 'air', 'conduit', or 'isolated' (default: air)
   - **air**: Cable installed in open air (best cooling)
   - **conduit**: Cable in conduit (reduced cooling, +10°C adjustment)
   - **isolated**: Cable insulated/isolated (poor cooling, +20°C adjustment)
10. **Wire Type**: Enter the wire type (default: generic)
    - **flry/flry-a/flry-b**: Automotive thin-wall PVC (105°C max)
    - **thhn**: Thermoplastic, high heat, nylon (90°C max)
    - **thwn**: Thermoplastic, heat/water resistant (75°C max)
    - **xlpe**: Cross-linked polyethylene (90°C max)
    - **pvc**: Standard PVC (70°C max)
    - **silicon**: Silicone rubber (200°C max)
    - **generic**: Generic wire (90°C max)

### Example Session

```
=== DC Cable Diameter Calculator ===
Supports 12V, 24V, 48V, 50V DC systems

Enter system voltage (V): 12
Enter current (A): 10
Enter cable length (m): 5
Enter maximum voltage drop percentage (default 3%): 3
Is this round trip length? (y/n, default: n): y
Cable material (copper/aluminum, default: copper): copper
Temperature unit (C/F, default: C): C
Enter ambient temperature: 25
Installation method (air/conduit/isolated, default: air): conduit
Wire type (flry/flry-a/flry-b/thhn/thwn/xlpe/pvc/silicon/generic, default: generic): flry

=== Calculation Results ===
System Voltage: 12.0 V
Current: 10.00 A
Cable Length: 5.00 m (round trip)
Material: Copper
Wire Type: FLRY (Max: 105°C) - Automotive thin-wall PVC (FLRY-A/B), stranded copper
Ambient Temperature: 25.0°C (25.0°C)
Installation Method: In conduit
Effective Operating Temperature: 35.0°C
Maximum Voltage Drop: 3.00% (0.36 V)

Required Cross-Sectional Area: 4.86 mm²
Required Diameter: 2.49 mm

=== Recommended Standard Sizes ===
Metric: 4.00 mm² (difference: 0.86 mm²)
AWG: 10 (5.26 mm², difference: 0.40 mm²)

=== Voltage Drop with Recommended Sizes ===
With 4.00 mm²: 0.44 V (3.65%)
With AWG 10 (5.26 mm²): 0.33 V (2.77%)
```

## Understanding the Results

### Required Cross-Sectional Area
The minimum cable cross-sectional area needed to meet your voltage drop requirements.

### Required Diameter
The diameter of a circular cable with the required cross-sectional area.

### Recommended Standard Sizes
The program suggests the closest standard cable sizes available:
- **Metric**: Standard metric sizes in mm² (e.g., 0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, etc.)
- **AWG**: American Wire Gauge sizes (e.g., 18, 16, 14, 12, 10, 8, 6, 4, 2, 1, 1/0, 2/0, 3/0, 4/0)

### Voltage Drop with Recommended Sizes
Shows the actual voltage drop you'll experience with the recommended cable sizes. This helps you verify that the selected cable meets your requirements.

## Important Notes

### Voltage Drop
- **3%** is the default and recommended for most applications
- **5%** may be acceptable for some applications
- Lower percentages provide better voltage regulation but require larger cables

### Round Trip vs One-Way
- **One-way**: Length from power source to load (single conductor)
- **Round trip**: Total length including both positive and negative/ground conductors (common in DC systems)

### Material Selection
- **Copper**: Lower resistance, better conductivity, more expensive
- **Aluminum**: Higher resistance, lighter weight, less expensive

### Temperature and Installation Method
The calculator accounts for temperature effects on cable resistance:
- **Resistivity increases with temperature**: Higher temperatures increase cable resistance, requiring larger cables
- **Installation method affects operating temperature**:
  - **In air**: Best cooling, minimal temperature rise above ambient
  - **In conduit**: Reduced cooling, approximately +10°C above ambient
  - **Isolated/Insulated**: Poor cooling, approximately +20°C above ambient

The program calculates the effective operating temperature (ambient + installation adjustment) and adjusts resistivity accordingly. This ensures accurate calculations for real-world conditions.

### Wire Type Selection

Different wire types have different maximum operating temperatures based on their insulation material:
- **FLRY/FLRY-A/FLRY-B**: Automotive wires, typically 105°C maximum
- **THHN**: 90°C maximum
- **THWN**: 75°C maximum  
- **XLPE**: 90°C maximum
- **PVC**: 70°C maximum
- **Silicone**: 200°C maximum (high temperature applications)

The program validates that the calculated effective operating temperature does not exceed the wire type's maximum rating. If it does, a warning is displayed recommending:
- Using a higher temperature rated wire
- Reducing ambient temperature
- Improving cooling (better installation method)
- Increasing cable size to reduce heat generation

## Common Use Cases

### 12V DC Systems
- Automotive applications
- RV and marine electrical systems
- Solar panel installations
- LED lighting systems

### 24V DC Systems
- Truck and commercial vehicle systems
- Industrial control systems
- Telecommunications equipment

### 48V DC Systems
- Data center power distribution
- Telecommunications systems
- Electric vehicle charging

## Limitations

- Maximum system voltage: 50V DC
- Calculations assume standard temperature (20°C)
- Does not account for temperature derating
- Does not consider current-carrying capacity (ampacity) - always verify cables can handle the current
- Standard cable sizes are limited to common sizes

## Safety Warning

⚠️ **Always consult with a qualified electrician or electrical engineer for critical applications.** This tool provides calculations for voltage drop only and does not replace professional engineering judgment. Always verify:
- Current-carrying capacity (ampacity) of the selected cable
- Temperature derating factors
- Local electrical codes and regulations
- Safety margins for your specific application

## AI-Assisted Development

This software has been developed with AI assistance. While the code has been reviewed and tested, users are advised to:
- Review the source code before using in production environments
- Verify calculations independently for critical applications
- Report any issues or inaccuracies found
- Understand that AI-generated code may contain errors and should be validated

## Troubleshooting

### "Invalid voltage" Error
- Ensure voltage is between 0 and 50V
- Use decimal notation (e.g., 12.0, 24.5)

### "Invalid current" Error
- Enter a positive number
- Use decimal notation if needed (e.g., 10.5)

### Unexpectedly Large Cable Sizes
- Check if you selected "round trip" when you meant "one-way"
- Consider increasing the allowed voltage drop percentage
- Verify your current and length inputs

## Documentation

- **[README.md](README.md)**: User guide (this file)
- **[DEVELOPER.md](DEVELOPER.md)**: Technical documentation for developers
- **[DOCUMENTATION.md](DOCUMENTATION.md)**: Documentation maintenance guide

## Contributing

See [DEVELOPER.md](DEVELOPER.md) for development guidelines and technical details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


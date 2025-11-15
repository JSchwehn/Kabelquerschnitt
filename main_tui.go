package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A49FA5")).
			MarginRight(2)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	resultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#50FA7B")).
			Padding(1, 2).
			MarginTop(1)
)

type model struct {
	inputs               []textinput.Model
	focused              int
	selectedMaterial     int
	selectedInstallation int
	selectedWireType     int
	roundTrip            bool
	materialList         list.Model
	installationList     list.Model
	wireTypeList         list.Model
	showResults          bool
	results              calculationResults
	err                  string
	warning              string
	step                 int // 0: inputs, 1: material, 2: installation, 3: wire type, 4: round trip, 5: results
}

type calculationResults struct {
	voltage               float64
	current               float64
	length                float64
	maxVoltageDropPercent float64
	roundTrip             bool
	material              CableMaterial
	installation          InstallationMethod
	wireType              WireType
	ambientTemp           float64
	ambientTempDisplay    float64
	tempUnit              string
	effectiveTemp         float64
	requiredArea          float64
	requiredDiameter      float64
	closestMetric         float64
	closestAWG            string
	awgArea               float64
	actualDropMetric      float64
	actualDropAWG         float64
}

func initialModel() model {
	inputs := make([]textinput.Model, 6)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "12, 24, 48, 50"
	inputs[0].Focus()
	inputs[0].CharLimit = 10
	inputs[0].Width = 20

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "10.0"
	inputs[1].CharLimit = 10
	inputs[1].Width = 20

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "5.0"
	inputs[2].CharLimit = 10
	inputs[2].Width = 20

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "3.0"
	inputs[3].CharLimit = 10
	inputs[3].Width = 20

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "20.0"
	inputs[4].CharLimit = 10
	inputs[4].Width = 20

	inputs[5] = textinput.New()
	inputs[5].Placeholder = "C or F"
	inputs[5].CharLimit = 1
	inputs[5].Width = 5

	// Material list
	materialItems := []list.Item{
		item{title: "Copper", desc: "Lower resistance, better conductivity"},
		item{title: "Aluminum", desc: "Higher resistance, lighter weight"},
	}
	materialList := list.New(materialItems, itemDelegate{}, 40, 5)
	materialList.Title = "Select Cable Material"
	materialList.SetShowStatusBar(false)
	materialList.SetFilteringEnabled(false)

	// Installation list
	installationItems := []list.Item{
		item{title: "In Air", desc: "Best cooling, minimal temperature rise"},
		item{title: "In Conduit", desc: "Reduced cooling, +10°C adjustment"},
		item{title: "Isolated", desc: "Poor cooling, +20°C adjustment"},
	}
	installationList := list.New(installationItems, itemDelegate{}, 40, 5)
	installationList.Title = "Select Installation Method"
	installationList.SetShowStatusBar(false)
	installationList.SetFilteringEnabled(false)

	// Wire type list
	wireTypeItems := []list.Item{
		item{title: "FLRY", desc: "Automotive thin-wall PVC (105°C max)"},
		item{title: "FLRY-A", desc: "Automotive flexible stranded (105°C max)"},
		item{title: "FLRY-B", desc: "Automotive symmetrical stranded (105°C max)"},
		item{title: "THHN", desc: "Thermoplastic, high heat, nylon (90°C max)"},
		item{title: "THWN", desc: "Thermoplastic, heat/water resistant (75°C max)"},
		item{title: "XLPE", desc: "Cross-linked polyethylene (90°C max)"},
		item{title: "PVC", desc: "Standard PVC (70°C max)"},
		item{title: "Silicone", desc: "Silicone rubber (200°C max)"},
		item{title: "Generic", desc: "Generic wire (90°C max)"},
	}
	wireTypeList := list.New(wireTypeItems, itemDelegate{}, 50, 10)
	wireTypeList.Title = "Select Wire Type"
	wireTypeList.SetShowStatusBar(false)
	wireTypeList.SetFilteringEnabled(true)

	return model{
		inputs:               inputs,
		focused:              0,
		selectedMaterial:     0,
		selectedInstallation: 0,
		selectedWireType:     0,
		materialList:         materialList,
		installationList:     installationList,
		wireTypeList:         wireTypeList,
		showResults:          false,
		step:                 0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case 0: // Input step
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				if m.focused < len(m.inputs)-1 {
					m.focused++
					m.inputs[m.focused].Focus()
					m.inputs[m.focused-1].Blur()
					return m, textinput.Blink
				} else {
					// Validate inputs and move to material selection
					if err := m.validateInputs(); err != nil {
						m.err = err.Error()
						return m, nil
					}
					m.err = ""
					m.step = 1
					m.materialList.Select(0)
					return m, nil
				}
			case "tab", "down":
				if m.focused < len(m.inputs)-1 {
					m.focused++
					m.inputs[m.focused].Focus()
					m.inputs[m.focused-1].Blur()
					return m, textinput.Blink
				}
			case "shift+tab", "up":
				if m.focused > 0 {
					m.focused--
					m.inputs[m.focused].Focus()
					m.inputs[m.focused+1].Blur()
					return m, textinput.Blink
				}
			}
			// Update focused input
			m.inputs[m.focused], _ = m.inputs[m.focused].Update(msg)

		case 1: // Material selection
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.selectedMaterial = m.materialList.Index()
				m.step = 2
				m.installationList.Select(0)
				return m, nil
			case "esc":
				m.step = 0
				return m, nil
			}
			var cmd tea.Cmd
			m.materialList, cmd = m.materialList.Update(msg)
			return m, cmd

		case 2: // Installation selection
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.selectedInstallation = m.installationList.Index()
				m.step = 3
				m.wireTypeList.Select(8) // Default to Generic
				return m, nil
			case "esc":
				m.step = 1
				return m, nil
			}
			var cmd tea.Cmd
			m.installationList, cmd = m.installationList.Update(msg)
			return m, cmd

		case 3: // Wire type selection
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.selectedWireType = m.wireTypeList.Index()
				m.step = 4
				return m, nil
			case "esc":
				m.step = 2
				return m, nil
			}
			var cmd tea.Cmd
			m.wireTypeList, cmd = m.wireTypeList.Update(msg)
			return m, cmd

		case 4: // Round trip selection
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "y", "Y":
				m.roundTrip = true
				m.calculate()
				m.step = 5
				return m, nil
			case "n", "N", "enter":
				m.roundTrip = false
				m.calculate()
				m.step = 5
				return m, nil
			case "esc":
				m.step = 3
				return m, nil
			}

		case 5: // Results
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "r":
				// Reset to start
				m = initialModel()
				return m, textinput.Blink
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) validateInputs() error {
	// Voltage
	if m.inputs[0].Value() == "" {
		return fmt.Errorf("voltage is required")
	}
	voltage, err := strconv.ParseFloat(m.inputs[0].Value(), 64)
	if err != nil || voltage <= 0 || voltage > 50 {
		return fmt.Errorf("voltage must be between 0 and 50V")
	}

	// Current
	if m.inputs[1].Value() == "" {
		return fmt.Errorf("current is required")
	}
	current, err := strconv.ParseFloat(m.inputs[1].Value(), 64)
	if err != nil || current <= 0 {
		return fmt.Errorf("current must be positive")
	}

	// Length
	if m.inputs[2].Value() == "" {
		return fmt.Errorf("length is required")
	}
	length, err := strconv.ParseFloat(m.inputs[2].Value(), 64)
	if err != nil || length <= 0 {
		return fmt.Errorf("length must be positive")
	}

	// Voltage drop (optional, defaults to 3%)
	// Temperature (optional, defaults to 20°C)
	// Temp unit (optional, defaults to C)

	return nil
}

func (m *model) calculate() {
	// Parse inputs
	voltage, _ := strconv.ParseFloat(m.inputs[0].Value(), 64)
	current, _ := strconv.ParseFloat(m.inputs[1].Value(), 64)
	length, _ := strconv.ParseFloat(m.inputs[2].Value(), 64)

	maxVoltageDropPercent := 3.0
	if m.inputs[3].Value() != "" {
		if val, err := strconv.ParseFloat(m.inputs[3].Value(), 64); err == nil && val > 0 && val <= 10 {
			maxVoltageDropPercent = val
		}
	}

	ambientTemp := 20.0
	if m.inputs[4].Value() != "" {
		if val, err := strconv.ParseFloat(m.inputs[4].Value(), 64); err == nil {
			ambientTemp = val
		}
	}

	tempUnit := "C"
	if m.inputs[5].Value() != "" {
		unit := strings.ToUpper(m.inputs[5].Value())
		if unit == "F" || unit == "C" {
			tempUnit = unit
		}
	}

	ambientTempCelsius := ambientTemp
	if tempUnit == "F" {
		ambientTempCelsius = fahrenheitToCelsius(ambientTemp)
	}

	// Get selections
	materialKeys := []string{"copper", "aluminum"}
	material := materials[materialKeys[m.selectedMaterial]]

	installationKeys := []InstallationMethod{InstallationInAir, InstallationConduit, InstallationIsolated}
	installation := installationKeys[m.selectedInstallation]

	wireTypeKeys := []string{"flry", "flry-a", "flry-b", "thhn", "thwn", "xlpe", "pvc", "silicon", "generic"}
	wireType := wireTypes[wireTypeKeys[m.selectedWireType]]

	// Calculate
	effectiveTemp := calculateEffectiveTemp(ambientTempCelsius, installation)
	requiredArea := calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, m.roundTrip, ambientTempCelsius, installation)
	requiredDiameter := areaToDiameter(requiredArea)

	closestMetric, _ := findClosestMetricSize(requiredArea)
	closestAWG, awgArea, _ := findClosestAWG(requiredArea)

	resistivity := calculateResistivityAtTemp(material, effectiveTemp)
	distanceFactor := map[bool]float64{true: 2.0, false: 1.0}[m.roundTrip]
	actualDropMetric := (current * resistivity * length * distanceFactor) / closestMetric
	actualDropAWG := (current * resistivity * length * distanceFactor) / awgArea

	// Validate temperature
	isValid, warningMsg := ValidateWireTemperature(effectiveTemp, wireType)
	m.warning = ""
	if !isValid || warningMsg != "" {
		m.warning = warningMsg
	}

	m.results = calculationResults{
		voltage:               voltage,
		current:               current,
		length:                length,
		maxVoltageDropPercent: maxVoltageDropPercent,
		roundTrip:             m.roundTrip,
		material:              material,
		installation:          installation,
		wireType:              wireType,
		ambientTemp:           ambientTempCelsius,
		ambientTempDisplay:    ambientTemp,
		tempUnit:              tempUnit,
		effectiveTemp:         effectiveTemp,
		requiredArea:          requiredArea,
		requiredDiameter:      requiredDiameter,
		closestMetric:         closestMetric,
		closestAWG:            closestAWG,
		awgArea:               awgArea,
		actualDropMetric:      actualDropMetric,
		actualDropAWG:         actualDropAWG,
	}
}

func (m model) View() string {
	if m.step == 0 {
		return m.inputView()
	} else if m.step == 1 {
		return m.materialView()
	} else if m.step == 2 {
		return m.installationView()
	} else if m.step == 3 {
		return m.wireTypeView()
	} else if m.step == 4 {
		return m.roundTripView()
	} else {
		return m.resultsView()
	}
}

func (m model) inputView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" DC Cable Diameter Calculator "))
	b.WriteString("\n\n")

	if m.err != "" {
		b.WriteString(errorStyle.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	labels := []string{
		"System Voltage (V):",
		"Current (A):",
		"Cable Length (m):",
		"Max Voltage Drop (%):",
		"Ambient Temperature:",
		"Temp Unit (C/F):",
	}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focused {
			b.WriteString(labelStyle.Render("> " + label))
		} else {
			b.WriteString(labelStyle.Render("  " + label))
		}
		b.WriteString(input.View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Tab/Enter to move to next field, Shift+Tab/Up to go back, Enter on last field to continue, q to quit"))

	return borderStyle.Render(b.String())
}

func (m model) materialView() string {
	return borderStyle.Render(
		titleStyle.Render(" Select Cable Material ") + "\n\n" +
			m.materialList.View() + "\n\n" +
			helpStyle.Render("↑/↓: Navigate, Enter: Select, Esc: Back, q: Quit"),
	)
}

func (m model) installationView() string {
	return borderStyle.Render(
		titleStyle.Render(" Select Installation Method ") + "\n\n" +
			m.installationList.View() + "\n\n" +
			helpStyle.Render("↑/↓: Navigate, Enter: Select, Esc: Back, q: Quit"),
	)
}

func (m model) wireTypeView() string {
	return borderStyle.Render(
		titleStyle.Render(" Select Wire Type ") + "\n\n" +
			m.wireTypeList.View() + "\n\n" +
			helpStyle.Render("↑/↓: Navigate, Enter: Select, Esc: Back, /: Filter, q: Quit"),
	)
}

func (m model) roundTripView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(" Round Trip Length "))
	b.WriteString("\n\n")
	b.WriteString("Is this round trip length (power + return)?\n\n")
	b.WriteString(helpStyle.Render("Press 'y' for yes, 'n' for no, Esc to go back, q to quit"))
	return borderStyle.Render(b.String())
}

func (m model) resultsView() string {
	var b strings.Builder
	r := m.results

	b.WriteString(titleStyle.Render(" Calculation Results "))
	b.WriteString("\n\n")

	// Input summary
	b.WriteString(valueStyle.Render("Input Parameters:\n"))
	b.WriteString(fmt.Sprintf("  System Voltage: %.1f V\n", r.voltage))
	b.WriteString(fmt.Sprintf("  Current: %.2f A\n", r.current))
	b.WriteString(fmt.Sprintf("  Cable Length: %.2f m (%s)\n", r.length, map[bool]string{true: "round trip", false: "one-way"}[r.roundTrip]))
	b.WriteString(fmt.Sprintf("  Max Voltage Drop: %.2f%% (%.2f V)\n", r.maxVoltageDropPercent, r.voltage*r.maxVoltageDropPercent/100))
	b.WriteString(fmt.Sprintf("  Material: %s\n", r.material.Name))
	b.WriteString(fmt.Sprintf("  Wire Type: %s (Max: %.0f°C)\n", r.wireType.Name, r.wireType.MaxTempCelsius))
	b.WriteString(fmt.Sprintf("  Installation: %s\n", map[InstallationMethod]string{
		InstallationInAir:    "In Air",
		InstallationConduit:  "In Conduit",
		InstallationIsolated: "Isolated",
	}[r.installation]))
	b.WriteString(fmt.Sprintf("  Ambient Temp: %.1f°%s (%.1f°C)\n", r.ambientTempDisplay, r.tempUnit, r.ambientTemp))
	b.WriteString(fmt.Sprintf("  Effective Temp: %.1f°C\n", r.effectiveTemp))

	if m.warning != "" {
		b.WriteString("\n")
		if strings.Contains(m.warning, "WARNING") {
			b.WriteString(errorStyle.Render("⚠️  " + m.warning + "\n"))
		} else {
			b.WriteString(warningStyle.Render("⚠️  " + m.warning + "\n"))
		}
	}

	b.WriteString("\n")
	b.WriteString(valueStyle.Render("Required Cable Size:\n"))
	b.WriteString(fmt.Sprintf("  Cross-Sectional Area: %.2f mm²\n", r.requiredArea))
	b.WriteString(fmt.Sprintf("  Diameter: %.2f mm\n", r.requiredDiameter))

	b.WriteString("\n")
	b.WriteString(valueStyle.Render("Recommended Standard Sizes:\n"))
	b.WriteString(fmt.Sprintf("  Metric: %.2f mm²\n", r.closestMetric))
	b.WriteString(fmt.Sprintf("  AWG: %s (%.2f mm²)\n", r.closestAWG, r.awgArea))

	b.WriteString("\n")
	b.WriteString(valueStyle.Render("Voltage Drop with Recommended Sizes:\n"))
	b.WriteString(fmt.Sprintf("  %.2f mm²: %.2f V (%.2f%%)\n", r.closestMetric, r.actualDropMetric, (r.actualDropMetric/r.voltage)*100))
	b.WriteString(fmt.Sprintf("  AWG %s: %.2f V (%.2f%%)\n", r.closestAWG, r.actualDropAWG, (r.actualDropAWG/r.voltage)*100))

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press 'r' to restart, 'q' to quit"))

	return resultBoxStyle.Render(b.String())
}

// Item for lists
type item struct {
	title, desc string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)
	str := fmt.Sprintf("%d. %s", index+1, i.title)

	if len(i.desc) > 0 {
		str += fmt.Sprintf(" - %s", i.desc)
	}

	// Highlight selected item
	if index == m.Index() {
		str = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Render("> " + str)
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

// Helper key struct (not used but required by bubbles)
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Esc   key.Binding
	Quit  key.Binding
}

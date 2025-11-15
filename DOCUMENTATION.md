# Documentation Maintenance Guide

This guide helps keep documentation synchronized with code changes.

## Documentation Files

- **README.md**: User-facing documentation (installation, usage, examples)
- **DEVELOPER.md**: Technical documentation (formulas, architecture, implementation)
- **DOCUMENTATION.md**: This file - maintenance guidelines
- **Code comments**: Inline documentation in `main.go`

## When to Update Documentation

### Code Changes Checklist

When making code changes, check if documentation needs updates:

- [ ] **Function signature changes** → Update DEVELOPER.md "Key Functions" section
- [ ] **New functions** → Add to DEVELOPER.md and add code comments
- [ ] **Changed calculations** → Update DEVELOPER.md "Calculation Methodology" section
- [ ] **New features** → Update both README.md and DEVELOPER.md
- [ ] **UI/UX changes** → Update README.md examples
- [ ] **New materials/sizes** → Update README.md usage and DEVELOPER.md data structures
- [ ] **Changed constants** → Update code comments and DEVELOPER.md

### Specific Scenarios

#### Adding a New Material

1. Update `main.go`:
   - Add resistivity constant
   - Add to `materials` map

2. Update `README.md`:
   - Add to "Material Selection" section
   - Update examples if needed

3. Update `DEVELOPER.md`:
   - Add to "Material Resistivity" section
   - Update "Adding New Features" section

4. Update tests:
   - Add test cases for new material

#### Changing Calculation Formula

1. Update `main.go`:
   - Update function implementation
   - Update function comments with new formula

2. Update `DEVELOPER.md`:
   - Update "Calculation Methodology" section
   - Update "Key Functions" section
   - Verify formula examples match code

3. Update tests:
   - Update or add test cases with new expected values

#### Adding Command-Line Arguments

1. Update `main.go`:
   - Implement argument parsing
   - Update prompts/help text

2. Update `README.md`:
   - Add "Command-Line Options" section
   - Update usage examples
   - Add examples for new arguments

3. Update `DEVELOPER.md`:
   - Add to "Future Enhancements" if planned, or move to "Architecture" if implemented

## Documentation Standards

### README.md Style

- **Target audience**: End users, non-technical
- **Language**: Clear, simple, avoid jargon
- **Structure**: Installation → Usage → Examples → Troubleshooting
- **Examples**: Real-world scenarios with expected output

### DEVELOPER.md Style

- **Target audience**: Developers, contributors
- **Language**: Technical, precise
- **Structure**: Architecture → Methodology → Code Structure → Functions
- **Formulas**: Include mathematical notation and explanations
- **Code examples**: Show actual implementation, not pseudocode

### Code Comments Style

- **Function comments**: Use Go doc comment format (`// FunctionName ...`)
- **Include formulas**: Document mathematical formulas in comments
- **Reference docs**: Link to DEVELOPER.md for detailed explanations
- **Parameter documentation**: Explain non-obvious parameters

## Verification Steps

Before committing documentation changes:

1. **Readability**: Read through the documentation as if you're a new user/developer
2. **Accuracy**: Verify formulas match implementation
3. **Examples**: Test all code examples and command-line examples
4. **Links**: Check all internal links work
5. **Consistency**: Ensure terminology is consistent across all docs
6. **Completeness**: Ensure all new features are documented

## Quick Reference

### Formula Locations

- **Voltage drop calculation**: `DEVELOPER.md` → "Calculation Methodology"
- **Diameter calculation**: `DEVELOPER.md` → "Calculation Methodology" → "Diameter Calculation"
- **Code implementation**: `main.go` → `calculateCableArea()` and `areaToDiameter()`

### Standard Values

- **Copper resistivity**: 0.0175 Ω·mm²/m (documented in `main.go` constants and `DEVELOPER.md`)
- **Aluminum resistivity**: 0.0283 Ω·mm²/m (documented in `main.go` constants and `DEVELOPER.md`)
- **Standard metric sizes**: Listed in `main.go` and `DEVELOPER.md`
- **AWG sizes**: Listed in `main.go` and `DEVELOPER.md`

### Common Updates

**When voltage limit changes:**
- Update validation message in `main.go`
- Update README.md "Features" and "Limitations"
- Update DEVELOPER.md if relevant

**When adding cable sizes:**
- Update `standardMetricSizes` or `awgSizes` in `main.go`
- Update function comments
- Update DEVELOPER.md "Data Structures" section
- Update README.md if it lists specific sizes

**When changing default values:**
- Update `main.go` default value
- Update README.md usage instructions
- Update examples if they use defaults

## Testing Documentation

After updating documentation:

```bash
# Verify code still compiles
go build

# Run tests to ensure examples are still valid
go test

# Check for broken links (if using markdown link checker)
# markdown-link-check README.md DEVELOPER.md
```

## Version History

Document significant documentation changes here:

- **2024-XX-XX**: Initial documentation created
  - README.md: User guide with examples
  - DEVELOPER.md: Technical documentation with formulas
  - Enhanced code comments with formula references

## Questions?

If unsure about documentation updates:

1. Check this guide first
2. Review similar changes in git history
3. Ask maintainers or create an issue

Remember: **Better to over-document than under-document!**


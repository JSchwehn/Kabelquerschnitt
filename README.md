# Web Version - DC Cable Diameter Calculator

This is a static web version of the DC Cable Diameter Calculator that can be deployed to GitHub Pages.

## Features

- ✅ Full calculation functionality matching the Go CLI/TUI versions
- ✅ Modern, responsive UI with Tailwind CSS
- ✅ Dark mode support
- ✅ **Multi-language support** (English, German, French, Swedish)
- ✅ Language preference persistence (saved in browser)
- ✅ **Cable weight calculation** (for copper and aluminum)
- ✅ **Accessibility features** (WCAG 2.1 compliant, screen reader support, keyboard navigation)
- ✅ No backend required - runs entirely in the browser
- ✅ Works offline after initial load
- ✅ Rounds up to nearest standard cable size for safety

## Deployment to GitHub Pages

### Option 1: Deploy from `web` branch

1. Go to your repository settings on GitHub
2. Navigate to "Pages" in the left sidebar
3. Under "Source", select "Deploy from a branch"
4. Select branch: `web`
5. Select folder: `/ (root)`
6. Click "Save"

Your site will be available at: `https://[username].github.io/Kabelquerschnitt/`

### Option 2: Deploy from `gh-pages` branch

1. Create a new branch called `gh-pages`:
   ```bash
   git checkout -b gh-pages
   git add index.html
   git commit -m "Add web version"
   git push origin gh-pages
   ```

2. Go to repository settings → Pages
3. Select `gh-pages` branch as source
4. Save

### Option 3: Use GitHub Actions (Recommended)

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches:
      - web

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./
          publish_branch: gh-pages
```

## Local Testing

Simply open `index.html` in your web browser, or use a local server:

```bash
# Python 3
python3 -m http.server 8000

# Node.js (with http-server)
npx http-server

# PHP
php -S localhost:8000
```

Then visit `http://localhost:8000`

## File Structure

```
.
├── index.html          # Main web application (single file)
├── test.html           # Test suite for calculation functions
├── README.md           # This file
└── ...                 # Other project files
```

## Testing

A comprehensive test suite is included in `test.html`. To run the tests:

1. Open `test.html` in your web browser
2. The tests will run automatically and display results
3. All calculation functions are tested including:
   - Temperature conversions (Fahrenheit ↔ Celsius)
   - Resistivity calculations at different temperatures
   - Effective temperature calculations
   - Wire temperature validation
   - Cable area calculations
   - Diameter calculations
   - Standard size lookups (metric and AWG)

The test suite includes:
- ✅ 30+ test cases
- ✅ Edge case testing
- ✅ Round-trip verification
- ✅ Comparison tests (copper vs aluminum, different temperatures, etc.)
- ✅ Visual test results with pass/fail indicators

## Cable Weight Calculation

The calculator now includes cable weight estimation based on:
- Cross-sectional area (mm²)
- Cable length (meters)
- Material (copper or aluminum)
- Round trip vs one-way configuration

Weight is displayed in grams (g) for smaller cables and kilograms (kg) for larger cables. The calculation accounts for both conductors in round-trip configurations.

**Material densities:**
- Copper: 8.96 g/m per mm²
- Aluminum: 2.70 g/m per mm²

## Browser Compatibility

- Chrome/Edge (latest) - Full support
- Firefox (latest) - Full support
- Safari (latest) - Full support
- Mobile browsers (iOS Safari, Chrome Mobile) - Full support

All modern browsers with JavaScript enabled are supported. The application uses standard web APIs and should work in any browser from the last 5 years.

## Internationalization

The web version supports multiple languages:

- **English** (en) - Default
- **German** (de) - Deutsch
- **French** (fr) - Français
- **Swedish** (sv) - Svenska

### Language Selection

- Use the language selector in the top-right corner of the page
- Your language preference is automatically saved in your browser's localStorage
- The page will remember your choice on future visits
- All UI elements, error messages, warnings, and results are translated

### Adding New Languages

To add a new language, edit `index.html` and add a new translation object to the `translations` object with all required keys. See the existing language objects for reference.

## Accessibility

The web version follows WCAG 2.1 accessibility guidelines:

- ✅ **Semantic HTML** - Proper use of `<header>`, `<main>`, `<section>`, and definition lists
- ✅ **ARIA labels** - All interactive elements have proper ARIA attributes
- ✅ **Screen reader support** - ARIA live regions announce calculation results
- ✅ **Keyboard navigation** - Full keyboard support with skip links
- ✅ **Focus management** - Automatic focus on results after calculation
- ✅ **Form labels** - All form fields are properly labeled and associated
- ✅ **Error announcements** - Validation errors are announced to screen readers

The application is fully accessible to users with disabilities and works well with screen readers like NVDA, JAWS, and VoiceOver.

## Notes

- The web version uses Tailwind CSS via CDN (requires internet connection for initial load)
- All calculations are performed client-side using vanilla JavaScript
- No data is sent to any server - completely private
- The calculations match the Go implementation exactly
- Language preferences are stored locally in your browser (localStorage)
- Cable sizes are rounded up to the nearest standard size for safety
- Maximum system voltage is configurable via `MAX_VOLTAGE` constant (default: 50V)

## License

Same as the main project - MIT License


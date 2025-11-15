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
├── index.html          # Main HTML structure and UI
├── translations.js     # All translation strings (i18n) for all supported languages
├── script.js           # JavaScript application logic and calculations
├── test.html           # Test suite for calculation functions
├── README.md           # This file
└── ...                 # Other project files
```

The application follows a clean separation of concerns:
- **HTML** (structure) and **JavaScript** (behavior) are separated
- **Translations** are isolated in their own file for easy maintenance
- This structure makes it easy to add new languages or modify translations without touching the main application logic

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
- **German (Simple Language)** (de-simple) - Deutsch (Einfache Sprache) - Easy-to-read German for accessibility
- **French** (fr) - Français
- **Swedish** (sv) - Svenska

### Language Selection

- Use the language selector in the top-right corner of the page
- Your language preference is automatically saved in your browser's localStorage
- The page will remember your choice on future visits
- All UI elements, error messages, warnings, and results are translated

### Simple Language (Einfache Sprache)

The German Simple Language version (de-simple) follows accessibility guidelines for easy-to-read text:
- Short sentences (max 15 words)
- Simple words and clear explanations
- No technical terms without explanation
- Clear structure and active voice
- Concrete examples

This makes the application more accessible to people with reading difficulties, learning disabilities, or limited German language skills.

### Adding New Languages

To add a new language, edit `translations.js` and add a new translation object to the `translations` object with all required keys. See the existing language objects (en, de, de-simple, fr, sv) for reference.

Example:
```javascript
const translations = {
    // ... existing languages ...
    "new-lang": {
        title: "Translation",
        subtitle: "Translation",
        // ... all other keys ...
    }
};
```

After adding the translation, also add the language option to the language selector in `index.html`.

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

## Code Structure

The application follows a clean separation of concerns:

- **`index.html`**: Contains the HTML structure, semantic markup, and accessibility features
- **`translations.js`**: Contains all translation strings for all supported languages (en, de, de-simple, fr, sv)
- **`script.js`**: Contains all JavaScript application logic including:
  - Internationalization (i18n) system and translation function
  - Calculation functions (cable area, diameter, weight, etc.)
  - Form validation and user interaction handling
  - UI updates and screen reader announcements

This separation provides:
- Better code organization and maintainability
- Easier debugging and testing
- Improved browser caching (each file can be cached separately)
- Clear separation between structure (HTML), translations (i18n), and behavior (JavaScript)
- Easy to add new languages without modifying application logic

## Notes

- The web version uses Tailwind CSS via CDN (requires internet connection for initial load)
- All calculations are performed client-side using vanilla JavaScript
- No data is sent to any server - completely private
- The calculations match the Go implementation exactly
- Language preferences are stored locally in your browser (localStorage)
- Cable sizes are rounded up to the nearest standard size for safety
- Maximum system voltage is configurable via `MAX_VOLTAGE` constant in `script.js` (default: 60V)
- Translations are managed in `translations.js` - easy to add new languages or modify existing translations

## License

Same as the main project - MIT License


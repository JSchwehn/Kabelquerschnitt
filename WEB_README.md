# Web Version - DC Cable Diameter Calculator

This is a static web version of the DC Cable Diameter Calculator that can be deployed to GitHub Pages.

## Features

- ✅ Full calculation functionality matching the Go CLI/TUI versions
- ✅ Modern, responsive UI with Tailwind CSS
- ✅ Dark mode support
- ✅ No backend required - runs entirely in the browser
- ✅ Works offline after initial load

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
├── WEB_README.md       # This file
└── ...                 # Other project files
```

## Browser Compatibility

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Mobile browsers (iOS Safari, Chrome Mobile)

## Notes

- The web version uses Tailwind CSS via CDN (requires internet connection for initial load)
- All calculations are performed client-side using vanilla JavaScript
- No data is sent to any server - completely private
- The calculations match the Go implementation exactly

## License

Same as the main project - MIT License


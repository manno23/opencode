# TUI to Ogen Type Flow Visualization

This directory contains an interactive Sankey diagram visualizing the transformation flow from TUI/Internal types to Ogen-generated types through intermediate steps.

## Files

- `index.html`: The main HTML page with Bootstrap styling and CDN libraries.
- `data.json`: JSON data defining nodes (types and intermediates) and links (transformation steps).
- `app.js`: ES6 module script that loads data, builds the Sankey diagram, and handles interactions.
- `styles.css`: Custom CSS for styling the diagram and page.
- `README.md`: This file.

## Opening the Visualization

### Option 1: Direct File Open

Open `index.html` directly in a web browser. Note that some browsers may restrict fetching `data.json` from `file://` protocol, but the app includes a fallback to inline data.

### Option 2: Simple Static Server

For full functionality, serve the directory with a static server:

```bash
cd packages/tui/api/test/visualization
python -m http.server 8000
# or
npx http-server
```

Then open `http://localhost:8000/index.html`.

## Editing Data

To add more types or modify flows:

1. Edit `data.json`:
   - Add nodes to the `nodes` array with unique `id` and display `name`.
   - Add links to the `links` array with `source`, `target`, `value` (usually 1), and `label` describing the step.

2. Reload the page to see changes.

The diagram uses real type names as specified and labels links with transformation descriptions.

# Assets

This directory holds non-code project assets for the triagegeist project. Only the logo is described here; if further assets are added (e.g. diagrams, screenshots), they should be listed and described in this file.

---

## Logo

The project logo is used in the root [README](../README.md) and, if needed, in documentation or release assets. It must be provided by the maintainers or contributors; the repository may ship without it (the README will show a broken image until the file is added).

### Specifications

| Property | Value |
|----------|--------|
| Filename | **logo.svg** (preferred) or **triagegeist-logo.svg** |
| Format | SVG (scalable vector graphics) |
| Background | None (transparent). The SVG should not include an opaque background so that it displays correctly on light and dark themes. |
| Canvas size | **8000 x 2000 px** (4:1 aspect ratio). This is the design canvas; the logo can be scaled down for display. |
| README display | The root README embeds the logo with `width="320"`, so the displayed size is 320 px wide (height scales to 80 px for 4:1). This keeps the page readable. |

### Instructions

1. Save your logo file as **assets/logo.svg** (or **assets/triagegeist-logo.svg**). If you use the latter, update the root README image source from `assets/logo.svg` to `assets/triagegeist-logo.svg`.
2. Ensure the SVG has no background (transparent). Remove any opaque background rectangle or fill if present.
3. Do not commit large raster images (PNG, JPG) for the main logo; SVG is preferred for scalability and repository size. Raster exports may be added separately for release artefacts if needed.
4. The logo is not covered by the same licence as the code unless explicitly stated; if you contribute a logo, please confirm that you have the right to license it for use in the project (e.g. original work or permissive licence).

### File location and README

The root README contains an image tag that points to this directory:

```html
<img src="assets/logo.svg" width="320" alt="triagegeist logo" />
```

If the file is missing, the image will not render until **assets/logo.svg** (or the chosen filename) is present. See [README](../README.md) for the exact tag and any alternative filename note.

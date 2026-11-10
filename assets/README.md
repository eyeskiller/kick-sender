# Kick Sender — app icon assets

Ready-to-use application icons derived from the approved Kick Sender icon.

- `KickSender.png` — 1024×1024 RGBA PNG (transparent outside the rounded app tile). Recommended source image for packaging and icon generation.
- `KickSender.icns` — native macOS app icon (multi-resolution).
- `KickSender.ico` — native Windows app icon (multi-resolution).
- `KickSender-{size}.png` — UI/favicon-sized versions in 16, 32, 64, 128, 256 and 512 pixels.

## Instructions for the coding agent
Use these existing assets directly; do not redraw, recreate or generate a mockup. Put the macOS ICNS file into the build's app icon resources, and use the ICO file for Windows. Use the 1024px PNG as the fallback/source for other platforms. Ensure the generated executable/app bundle actually uses the icon in Finder/Dock/taskbar.

# Kick Interval GUI

Cross-platform graphical interface for Kick Interval, built with Electron.

## Features

- 🖥️ **Modern GUI** - Clean, intuitive interface with real-time status updates
- 🌍 **Cross-platform** - Works on macOS, Windows, and Linux
- 🔌 **Integrated Backend** - Automatically starts the Go backend process
- 📊 **Live Logs** - Stream activity logs in real-time
- 💾 **Persistent Config** - Settings are automatically saved

## Quick Start

### Prerequisites

1. **Build the Go backend** (if not already done):
   ```bash
   cd /Users/lubosberan/Dev/kick-sender
   make build
   ```

2. **Install Electron dependencies**:
   ```bash
   cd electron
   npm install
   ```

### Running the GUI

```bash
cd electron
npm start
```

The application will:
- Start the Go backend server on port 18765
- Open the GUI window
- Load your saved configuration (if any)

## Usage

### 1. Authentication
- Click the **Authenticate** button
- Your browser will open to the Kick OAuth flow
- Authorize the application
- The GUI will show "Authenticated" when complete

### 2. Configuration
Fill in the three required fields:
- **Channel**: Enter a Kick channel slug (e.g., `examplechannel`) or URL
- **Message**: Enter the message to send (max 500 characters)
- **Interval**: Enter how often to send (e.g., `10m`, `1h`, `60s`)

Click the **Set** button next to each field to save it.

### 3. Start Sending
- Click **Start** to begin automated message sending
- The status will change to "Running"
- You'll see real-time logs of each message sent
- Click **Stop** to stop the scheduler

### 4. Activity Log
- View all activity in the logs section at the bottom
- Logs show timestamps, status (INFO/ERROR), and messages
- Successful sends show "Message sent successfully"
- Failed sends show the error details

### 5. Settings & Reset
Click the gear icon (⚙️) in the top-right corner to access settings:

- **Clear Authentication** — Remove OAuth access/refresh tokens only
- **Clear Settings** — Remove channel, message, and interval (keeps authentication)
- **Reset Everything** — Clear all data including Kick Developer App credentials (client ID, client secret)

**"Reset Everything"** is designed for testing the first-run experience. It deletes:
- OAuth tokens from the system keychain
- `oauth.json` (Kick Developer App credentials)
- `config.json` (channel, message, interval)

After resetting, the onboarding wizard will appear again as if the app was never configured.

## Building for Distribution

### Build for Current Platform

```bash
cd electron
npm run build
```

### Build for Specific Platforms

```bash
# macOS (Universal binary)
npm run build:mac

# Windows (NSIS installer)
npm run build:win

# Linux (AppImage)
npm run build:linux

# All platforms
npm run build:all
```

The built applications will be in `electron/dist/`.

## Architecture

```
electron/
├── src/
│   ├── main.js           # Electron main process
│   ├── preload.js        # Secure bridge between main/renderer
│   └── renderer/
│       ├── index.html    # GUI layout
│       ├── style.css     # Styling
│       └── app.js        # Frontend logic
├── package.json          # Electron dependencies
└── dist/                 # Built applications (after build)
```

### How It Works

1. **Main Process** (`main.js`):
   - Spawns the Go backend as a child process
   - Creates and manages the Electron window
   - Handles app lifecycle (start/stop)

2. **Backend Communication**:
   - Go backend runs HTTP server on port 18765
   - Electron renderer makes API calls via fetch
   - Real-time logs via Server-Sent Events (SSE)

3. **Security**:
   - Context isolation enabled
   - No Node.js access in renderer
   - CORS enabled for local communication

## Troubleshooting

### Backend Not Starting

**Problem**: "Failed to start backend"

**Solution**:
1. Make sure you built the Go backend: `make build`
2. Check that `build/kick-interval-arm64` (or `-amd64`) exists
3. Check the Electron console for error details

### GUI Not Connecting

**Problem**: Status shows "Not authenticated" even after auth

**Solution**:
1. Check backend logs in Electron console
2. Verify API is accessible: `curl http://localhost:18765/api/status`
3. Restart the application

### Authentication Not Working

**Problem**: OAuth flow doesn't complete

**Solution**:
1. Make sure you have a Kick Developer App configured
2. Verify redirect URI is `http://localhost:18765/callback`
3. Check backend logs for OAuth errors

### Port Already in Use

**Problem**: "Port 18765 already in use"

**Solution**:
1. Kill the existing backend process: `pkill -f kick-interval`
2. Or change the port in `electron/src/main.js`

## Development

### Run in Development Mode

```bash
cd electron
npm start
```

### View Logs

Open DevTools in Electron: `Cmd+Option+I` (Mac) or `Ctrl+Shift+I` (Windows/Linux)

### Hot Reload

For development, you can use `electron-reload`:

```bash
npm install --save-dev electron-reload
```

## Notes

- The GUI uses port **18765** (CLI uses 8765) to avoid conflicts
- Backend process is automatically killed when GUI closes
- Configuration is shared between CLI and GUI
- OAuth credentials are stored in macOS Keychain (or equivalent on other platforms)

## Platform-Specific Notes

### macOS
- Uses native Keychain for credential storage
- Universal binary supports both Intel and Apple Silicon

### Windows
- Uses Windows Credential Manager
- NSIS installer for easy installation

### Linux
- Uses libsecret for credential storage
- AppImage for portable installation

## License

Same as the main Kick Interval project (MIT)

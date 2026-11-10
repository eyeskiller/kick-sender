# Kick Interval Chat Sender

A cross-platform application that automatically sends configurable messages to Kick.com livestream chat at user-defined intervals. Available as both a CLI and a modern GUI.

## Features

- 🖥️ **Dual Interface**: Choose between a modern GUI (Electron) or a lightweight CLI
- 🌍 **Cross-Platform**: GUI works on macOS, Windows, and Linux
- 🎯 **Simple CLI**: Interactive terminal interface with colored output
- 🔐 **Secure OAuth 2.1**: Uses official Kick API with PKCE authentication
- 💾 **Persistent Configuration**: Settings survive application restarts
- 🔑 **Platform Keychain**: Credentials stored securely (macOS Keychain, Windows Credential Manager, Linux libsecret)
- ⏱️ **Flexible Intervals**: Send messages every 60s, 5m, 10m, 1h, or custom duration
- 🛡️ **Robust Error Handling**: Automatic token refresh, rate limiting, and graceful shutdown
- 📊 **Live Activity Logs**: Real-time monitoring of message delivery status

## Prerequisites

- A Kick.com account
- A Kick Developer App (see setup instructions below)

**For CLI only:**
- macOS (Apple Silicon or Intel), Windows, or Linux

**For GUI:**
- Node.js 18+ (for building/running the Electron app)
- Works on macOS, Windows, and Linux

## Setup

### 1. Create a Kick Developer App

1. Sign up for a [Kick.com](https://kick.com) account if you haven't already
2. Enable Two-Factor Authentication (2FA) in your [Kick Account Settings](https://kick.com/settings/security)
3. Navigate to the [Developer Tab](https://kick.com/settings/developer)
4. Click "Create App" and fill in the details:
   - **App Name**: `Kick Interval Sender` (or any name you prefer)
   - **Redirect URL**: `http://localhost:8765/callback`
   - **Description**: Brief description of your app
5. After creating the app, you'll receive:
   - **Client ID**
   - **Client Secret**
   - **Redirect URL** (should match what you entered)

**Important**: The redirect URL must exactly match `http://localhost:8765/callback` for this application to work.

### 2. Build or Install

#### Option A: Build from source

```bash
# Clone or download the repository
cd kick-sender

# Build for your platform
make build-local

# The binary will be in build/kick-interval
./build/kick-interval --help
```

#### Option B: Install system-wide

```bash
make install
```

This installs the binary to `/usr/local/bin/kick-interval`.

## Usage

### First Run

When you run the application for the first time:

```bash
kick-interval
```

The application will guide you through:

1. **Authentication Setup**: Enter your Client ID and Client Secret from your Kick Developer App
2. **Browser Authorization**: Your browser opens to authenticate with Kick
3. **Channel Selection**: Enter the target channel (slug or URL)
4. **Message Configuration**: Enter the message to send
5. **Interval Selection**: Choose how often to send (e.g., `10m` for every 10 minutes)
6. **Confirmation**: Review settings and start sending

### Interactive Mode

Simply run without arguments:

```bash
kick-interval
```

Follow the prompts to configure everything interactively.

### Command-Line Mode

Provide all settings via flags:

```bash
kick-interval --channel examplechannel \
  --message "Follow the stream and join our Discord!" \
  --interval 10m
```

### Send Immediately

Use `--send-now` to send the first message immediately instead of waiting for the first interval:

```bash
kick-interval --send-now
```

### Dry Run

Validate your configuration without actually sending messages:

```bash
kick-interval --dry-run
```

### Saved Configuration

After the first run, your settings are saved. On subsequent launches:

```bash
$ kick-interval

Saved configuration found:

Channel:  examplechannel
Message:  Follow the stream and join our Discord!
Interval: 10m

[1] Start with saved settings
[2] Change channel
[3] Change message
[4] Change interval
[5] Change Kick account
[6] Exit

Select:
>
```

### Logout

Clear stored credentials:

```bash
kick-interval logout
```

### Help

```bash
kick-interval --help
```

## GUI (Graphical Interface)

The Electron GUI provides a user-friendly interface for all platforms.

### Quick Start

```bash
# Build the Go backend (if not already done)
make build

# Install Electron dependencies
cd electron
npm install

# Run the GUI
npm start
```

The GUI will:
- Start the backend server automatically
- Open a modern window with all controls
- Load your saved configuration (if any)
- Show real-time activity logs

### GUI Features

- **Visual Configuration**: Set channel, message, and interval with a clean interface
- **Authentication Button**: One-click OAuth flow
- **Start/Stop Controls**: Easy scheduler management
- **Live Status**: See message count, next send time, and running status
- **Activity Log**: Real-time logs showing every send attempt
- **Settings & Reset**: Clear authentication, configuration, or everything (including Kick Developer App credentials)
- **Cross-Platform**: Works on macOS, Windows, and Linux

### Building GUI for Distribution

```bash
cd electron

# Build for current platform
npm run build

# Build for specific platforms
npm run build:mac     # macOS (Universal binary)
npm run build:win     # Windows (NSIS installer)
npm run build:linux   # Linux (AppImage)
npm run build:all     # All platforms
```

Built applications will be in `electron/dist/`.

For more GUI details, see [electron/README.md](electron/README.md).

## Interval Format

Intervals are specified as a number followed by a unit:

- `s` - seconds (e.g., `60s`)
- `m` - minutes (e.g., `5m`, `10m`, `30m`)
- `h` - hours (e.g., `1h`, `2h`)

**Minimum interval**: 60 seconds (to respect rate limits)

**Default interval**: 10 minutes

## Examples

### Basic Usage

```bash
# Send "Hello!" every 10 minutes
kick-interval --channel mystreamer --message "Hello!" --interval 10m
```

### Promotional Message

```bash
# Send promotional message every 30 minutes
kick-interval --channel gamer123 \
  --message "Check out our Discord: discord.gg/example" \
  --interval 30m
```

### Quick Start

```bash
# Start immediately with saved settings
kick-interval --send-now
```

### Testing

```bash
# Validate configuration without sending
kick-interval --dry-run
```

## Troubleshooting

### Authentication Issues

**Problem**: "No tokens found" or authentication fails

**Solution**:
1. Verify your Kick Developer App is configured correctly
2. Ensure the redirect URL matches exactly: `http://localhost:8765/callback`
3. Try logging out and re-authenticating:
   ```bash
   kick-interval logout
   kick-interval
   ```

### Channel Not Found

**Problem**: "Channel validation failed"

**Solution**:
1. Verify the channel exists on Kick.com
2. Use the exact channel slug from the URL (e.g., `kick.com/username` → `username`)
3. Check that you're authenticated and have permission to chat in that channel

### Rate Limiting

**Problem**: Messages show "RETRY" or rate limit errors

**Solution**:
1. Increase your interval (minimum is 60 seconds)
2. The application automatically handles rate limits with exponential backoff
3. Wait before retrying

### Token Expiration

**Problem**: "Unauthorized" errors during sending

**Solution**:
- The application automatically refreshes tokens
- If refresh fails, re-authenticate:
  ```bash
  kick-interval logout
  kick-interval
  ```

### Application Won't Start

**Problem**: Binary not found or permission denied

**Solution**:
```bash
# Rebuild the application
make build-local

# Or install to system path
make install
```

## Development

### Prerequisites

- Go 1.21 or later
- macOS (for building and testing)

### Build

```bash
# Build for current platform
make build-local

# Build for all macOS architectures
make build

# Build for ARM64 (Apple Silicon) only
make build-arm64

# Build for AMD64 (Intel) only
make build-amd64
```

### Test

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run all linters
make lint
```

### Clean

```bash
# Remove build artifacts
make clean
```

## Architecture

```
kick-sender/
├── cmd/
│   └── kick-interval/
│       └── main.go           # CLI entry point
├── internal/
│   ├── auth/
│   │   ├── oauth.go          # OAuth utilities (PKCE, state)
│   │   ├── flow.go           # OAuth flow implementation
│   │   └── keychain.go       # macOS Keychain integration
│   ├── kick/
│   │   ├── client.go         # Kick API client
│   │   └── channels.go       # Channel/message validation
│   ├── scheduler/
│   │   └── scheduler.go      # Message scheduler
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── server/
│   │   └── server.go         # HTTP API server (for GUI)
│   └── ui/
│       └── terminal.go       # Terminal UI utilities
├── electron/                 # Electron GUI
│   ├── src/
│   │   ├── main.js           # Electron main process
│   │   ├── preload.js        # Security bridge
│   │   └── renderer/         # GUI frontend
│   └── package.json
├── Makefile
├── README.md
└── .gitignore
```

## Security

- **OAuth 2.1 with PKCE**: Secure authentication flow
- **Platform Keychain**: Credentials stored in system keychain (macOS Keychain, Windows Credential Manager, Linux libsecret), never in plaintext
- **No Password Storage**: Never stores or transmits your Kick password
- **Token Refresh**: Automatic token refresh without re-authentication
- **Minimal Scopes**: Only requests necessary permissions (`chat:write`, `channel:read`, `user:read`)
- **Context Isolation**: GUI uses Electron security best practices

## Limitations

- Single account, single channel, single message
- Minimum 60-second interval
- Respects Kick's rate limits and anti-spam policies

## Manual Integration Testing

To perform a live test with your own Kick account:

1. Create a Kick Developer App as described above
2. Build the application: `make build-local`
3. Run: `./build/kick-interval`
4. Follow the prompts to authenticate
5. Select your own channel or a test channel you control
6. Enter a test message
7. Set a short interval (e.g., `60s`)
8. Verify messages appear in the chat

**Note**: Always test on channels you own or have permission to use.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues related to:

- **This application**: Check the troubleshooting section or open an issue
- **Kick API**: Refer to the [Kick Developer Documentation](https://docs.kick.com)
- **Kick Developer Apps**: Visit [kick.com/settings/developer](https://kick.com/settings/developer)

## Acknowledgments

- Built with the official [Kick Developer API](https://docs.kick.com/apis)
- Uses [go-keyring](https://github.com/zalando/go-keyring) for secure credential storage

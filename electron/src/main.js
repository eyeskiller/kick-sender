const { app, BrowserWindow, Menu } = require('electron');
const path = require('path');
const { spawn } = require('child_process');

let mainWindow;
let backendProcess;

function getBackendPath() {
  // In production, the backend binary is in extraResources
  if (app.isPackaged) {
    const basePath = process.resourcesPath;
    const platform = process.platform;
    const arch = process.arch;

    let binaryName = 'kick-interval';
    if (platform === 'darwin') {
      binaryName = arch === 'arm64' ? 'kick-interval-arm64' : 'kick-interval-amd64';
    } else if (platform === 'win32') {
      binaryName = 'kick-interval.exe';
    } else if (platform === 'linux') {
      binaryName = 'kick-interval-linux';
    }

    return path.join(basePath, 'bin', binaryName);
  }

  // In development, look in the parent build directory
  const basePath = path.join(__dirname, '..', '..', 'build');
  const platform = process.platform;
  const arch = process.arch;

  let binaryName = 'kick-interval';
  if (platform === 'darwin') {
    binaryName = arch === 'arm64' ? 'kick-interval-arm64' : 'kick-interval-amd64';
  } else if (platform === 'win32') {
    binaryName = 'kick-interval.exe';
  } else if (platform === 'linux') {
    binaryName = 'kick-interval-linux';
  }

  return path.join(basePath, binaryName);
}

function startBackend() {
  const backendPath = getBackendPath();
  const port = '18765'; // Use different port than CLI to avoid conflicts

  console.log('Starting backend:', backendPath, 'on port', port);

  backendProcess = spawn(backendPath, ['serve', port], {
    stdio: ['pipe', 'pipe', 'pipe'],
  });

  backendProcess.stdout.on('data', (data) => {
    console.log(`Backend: ${data}`);
  });

  backendProcess.stderr.on('data', (data) => {
    console.error(`Backend error: ${data}`);
  });

  backendProcess.on('close', (code) => {
    console.log(`Backend process exited with code ${code}`);
  });

  backendProcess.on('error', (err) => {
    console.error('Failed to start backend:', err);
  });

  return port;
}

function stopBackend() {
  if (backendProcess) {
    backendProcess.kill('SIGTERM');
    backendProcess = null;
  }
}

function createWindow(port) {
  mainWindow = new BrowserWindow({
    width: 1100,
    height: 720,
    minWidth: 960,
    minHeight: 600,
    title: 'Kick Sender',
    backgroundColor: '#0e0e10',
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
    },
  });

  mainWindow.loadFile(path.join(__dirname, 'renderer', 'index.html'));

  // Pass the API port to the renderer
  mainWindow.webContents.once('did-finish-load', () => {
    mainWindow.webContents.send('set-api-port', port);
  });

  // Remove default menu
  Menu.setApplicationMenu(null);

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

app.whenReady().then(() => {
  const port = startBackend();

  // Give the backend a moment to start
  setTimeout(() => {
    createWindow(port);
  }, 1000);

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      const port = process.env.BACKEND_PORT || '18765';
      createWindow(port);
    }
  });
});

app.on('window-all-closed', () => {
  stopBackend();
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  stopBackend();
});

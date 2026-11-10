const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  onSetApiPort: (callback) => ipcRenderer.on('set-api-port', (event, port) => callback(port)),
});

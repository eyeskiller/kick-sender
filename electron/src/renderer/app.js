let API_BASE = 'http://localhost:18765';
let onboarded = false;
let nextSendUnix = 0;
let isRunning = false;

// Initialize language before anything else
initLanguage();

// Listen for API port from main process
window.electronAPI.onSetApiPort((port) => {
  API_BASE = `http://localhost:${port}`;
  initApp();
});

document.addEventListener('DOMContentLoaded', () => {
  setTimeout(() => {
    if (!window.appInitialized) initApp();
  }, 1500);
});

async function initApp() {
  window.appInitialized = true;
  await checkAndRender();
  setInterval(pollStatus, 2000);
}

async function checkAndRender() {
  try {
    const status = await fetchStatus();
    const needsOnboarding = !status.authenticated && !status.channel && !status.message && !status.interval;

    if (needsOnboarding && !onboarded) {
      showOnboarding();
      initOnboarding(status);
    } else {
      showMainApp();
      initMainApp(status);
    }
  } catch (err) {
    console.error('Failed to check status:', err);
    showOnboarding();
    initOnboarding({ authenticated: false, has_oauth_config: false });
  }
}

async function fetchStatus() {
  const res = await fetch(`${API_BASE}/api/status`);
  return res.json();
}

function showOnboarding() {
  document.getElementById('onboarding').style.display = 'flex';
  document.getElementById('mainApp').style.display = 'none';
}

function showMainApp() {
  document.getElementById('onboarding').style.display = 'none';
  document.getElementById('mainApp').style.display = 'flex';
}

// ============================================================
// ONBOARDING
// ============================================================

function initOnboarding(status) {
  // Re-apply translations after showing (preset buttons have dynamic labels)
  applyTranslations();

  let startStep = 1;
  if (status.authenticated) startStep = 4;
  else if (status.has_oauth_config) startStep = 3;
  showObStep(startStep);

  document.querySelectorAll('.ob-next[data-next]').forEach(btn => {
    btn.addEventListener('click', () => showObStep(parseInt(btn.dataset.next)));
  });

  document.querySelectorAll('.ob-back[data-back]').forEach(btn => {
    btn.addEventListener('click', () => showObStep(parseInt(btn.dataset.back)));
  });

  document.getElementById('obOpenPortal').addEventListener('click', () => {
    window.open('https://kick.com/settings/developer', '_blank');
  });
  const openDevLink = document.getElementById('openDevPortal');
  if (openDevLink) {
    openDevLink.addEventListener('click', (e) => {
      e.preventDefault();
      window.open('https://kick.com/settings/developer', '_blank');
    });
  }

  document.getElementById('obSaveCreds').addEventListener('click', saveObCredentials);
  document.getElementById('obStartAuth').addEventListener('click', startObAuth);

  document.querySelectorAll('#onboarding .preset-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('#onboarding .preset-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      document.getElementById('obInterval').value = btn.dataset.value;
    });
  });

  document.getElementById('obFinish').addEventListener('click', finishObSetup);
}

function showObStep(step) {
  document.querySelectorAll('.ob-panel').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.ob-step').forEach(s => {
    const sNum = parseInt(s.dataset.step);
    s.classList.remove('active', 'done');
    if (sNum === step) s.classList.add('active');
    else if (sNum < step) s.classList.add('done');
  });
  const panel = document.querySelector(`.ob-panel[data-panel="${step}"]`);
  if (panel) panel.classList.add('active');
}

async function saveObCredentials() {
  const clientId = document.getElementById('obClientId').value.trim();
  const clientSecret = document.getElementById('obClientSecret').value.trim();
  const redirectUri = document.getElementById('obRedirectUri').value.trim() || 'http://localhost:8765/callback';
  const errEl = document.getElementById('obCredError');

  if (!clientId || !clientSecret) {
    errEl.textContent = t('ob.s2.error_empty');
    errEl.classList.add('show');
    return;
  }
  errEl.classList.remove('show');

  try {
    const res = await fetch(`${API_BASE}/api/auth/config`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        client_id: clientId,
        client_secret: clientSecret,
        redirect_uri: redirectUri,
      }),
    });

    if (!res.ok) {
      const err = await res.text();
      errEl.textContent = t('ob.s2.error_save', err);
      errEl.classList.add('show');
      return;
    }

    showObStep(3);
  } catch (e) {
    errEl.textContent = t('ob.s2.error_conn', e.message);
    errEl.classList.add('show');
  }
}

async function startObAuth() {
  const statusText = document.getElementById('obAuthText');
  const indicator = document.getElementById('obAuthIndicator');
  const waitEl = document.getElementById('obAuthWait');

  statusText.textContent = t('ob.s3.opening');
  indicator.className = 'status-indicator';
  waitEl.style.display = 'flex';

  try {
    const res = await fetch(`${API_BASE}/api/auth/start`, { method: 'POST' });
    const data = await res.json();

    if (res.ok) {
      statusText.textContent = t('ob.s3.waiting');
      const checkAuth = setInterval(async () => {
        try {
          const s = await fetchStatus();
          if (s.authenticated) {
            clearInterval(checkAuth);
            statusText.textContent = t('ob.s3.success');
            indicator.className = 'status-indicator authenticated';
            waitEl.style.display = 'none';
            setTimeout(() => showObStep(4), 1000);
          }
        } catch (_) { /* ignore */ }
      }, 1500);
    } else {
      statusText.textContent = t('ob.s3.error', data.message || 'Failed');
      indicator.className = 'status-indicator not-authenticated';
      waitEl.style.display = 'none';
    }
  } catch (e) {
    statusText.textContent = t('ob.s3.error', e.message);
    indicator.className = 'status-indicator not-authenticated';
    waitEl.style.display = 'none';
  }
}

async function finishObSetup() {
  const channel = document.getElementById('obChannel').value.trim();
  const message = document.getElementById('obMessage').value.trim();
  const interval = document.getElementById('obInterval').value.trim() || '10m';
  const errEl = document.getElementById('obConfigError');

  if (!channel) { errEl.textContent = t('ob.s4.error_channel'); errEl.classList.add('show'); return; }
  if (!message) { errEl.textContent = t('ob.s4.error_message'); errEl.classList.add('show'); return; }
  errEl.classList.remove('show');

  try {
    let r = await fetch(`${API_BASE}/api/channel`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ channel }),
    });
    if (!r.ok) { errEl.textContent = t('ob.s4.error_field', t('main.channel'), await r.text()); errEl.classList.add('show'); return; }
  } catch (e) { errEl.textContent = t('status.conn_error', e.message); errEl.classList.add('show'); return; }

  try {
    let r = await fetch(`${API_BASE}/api/message`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ message }),
    });
    if (!r.ok) { errEl.textContent = t('ob.s4.error_field', t('main.message'), await r.text()); errEl.classList.add('show'); return; }
  } catch (e) { errEl.textContent = t('status.conn_error', e.message); errEl.classList.add('show'); return; }

  try {
    let r = await fetch(`${API_BASE}/api/interval`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ interval }),
    });
    if (!r.ok) { errEl.textContent = t('ob.s4.error_field', t('main.interval'), await r.text()); errEl.classList.add('show'); return; }
  } catch (e) { errEl.textContent = t('status.conn_error', e.message); errEl.classList.add('show'); return; }

  onboarded = true;
  showMainApp();
  initMainApp(await fetchStatus());
}

// ============================================================
// MAIN APP
// ============================================================

let mainAppInited = false;

function initMainApp(status) {
  if (mainAppInited) {
    updateMainStatus(status);
    return;
  }
  mainAppInited = true;

  // Apply translations to the now-visible main UI
  applyTranslations();

  // Language selector
  const langSelect = document.getElementById('langSelect');
  if (langSelect) {
    langSelect.value = currentLang;
    langSelect.addEventListener('change', (e) => setLanguage(e.target.value));
  }

  document.getElementById('setChannelBtn').addEventListener('click', setChannel);
  document.getElementById('setMessageBtn').addEventListener('click', setMessage);
  document.getElementById('setIntervalBtn').addEventListener('click', setInterval_);
  document.getElementById('startBtn').addEventListener('click', startScheduler);
  document.getElementById('stopBtn').addEventListener('click', stopScheduler);

  document.querySelectorAll('#mainApp .preset-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
      document.querySelectorAll('#mainApp .preset-btn').forEach(b => b.classList.remove('active'));
      e.target.classList.add('active');
      document.getElementById('interval').value = e.target.dataset.value;
      setInterval_();
    });
  });

  const settingsBtn = document.getElementById('settingsBtn');
  const settingsModal = document.getElementById('settingsModal');
  const modalClose = settingsModal.querySelector('.modal-close');

  settingsBtn.addEventListener('click', () => {
    settingsModal.classList.add('show');
    document.getElementById('backendPort').textContent = API_BASE.split(':')[2];
  });

  modalClose.addEventListener('click', () => settingsModal.classList.remove('show'));
  settingsModal.addEventListener('click', (e) => {
    if (e.target === settingsModal) settingsModal.classList.remove('show');
  });

  document.getElementById('clearAuthBtn').addEventListener('click', () => resetData(false, true));
  document.getElementById('clearConfigBtn').addEventListener('click', () => resetData(true, false));
  document.getElementById('resetAllBtn').addEventListener('click', () => resetData(true, true));

  loadConfig(status);
  streamLogs();
  updateMainStatus(status);
}

function loadConfig(status) {
  if (status.channel) {
    document.getElementById('channel').value = status.channel;
    showStatus('channelStatus', t('status.loaded'), 'success');
  }
  if (status.message) {
    document.getElementById('message').value = status.message;
    showStatus('messageStatus', t('status.loaded'), 'success');
  }
  if (status.interval) {
    document.getElementById('interval').value = status.interval;
    showStatus('intervalStatus', t('status.loaded'), 'success');
    document.querySelectorAll('#mainApp .preset-btn').forEach(btn => {
      if (btn.dataset.value === status.interval) btn.classList.add('active');
    });
  }
}

async function pollStatus() {
  try {
    const status = await fetchStatus();

    if (document.getElementById('mainApp').style.display !== 'none') {
      updateMainStatus(status);
    }

    const obAuthPanel = document.querySelector('.ob-panel[data-panel="3"]');
    if (obAuthPanel && obAuthPanel.classList.contains('active') && status.authenticated) {
      const statusText = document.getElementById('obAuthText');
      const indicator = document.getElementById('obAuthIndicator');
      const waitEl = document.getElementById('obAuthWait');
      if (statusText) {
        statusText.textContent = t('ob.s3.success');
        indicator.className = 'status-indicator authenticated';
        waitEl.style.display = 'none';
        setTimeout(() => showObStep(4), 800);
      }
    }
  } catch (_) {}
}

function updateMainStatus(status) {
  const headerAuthText = document.getElementById('headerAuthText');
  const headerAuthIndicator = document.getElementById('headerAuthIndicator');
  if (headerAuthText && headerAuthIndicator) {
    if (status.authenticated) {
      headerAuthText.textContent = t('auth.authenticated');
      headerAuthIndicator.className = 'status-indicator sm authenticated';
    } else {
      headerAuthText.textContent = t('auth.not_authenticated');
      headerAuthIndicator.className = 'status-indicator sm not-authenticated';
    }
  }

  const runningStatus = document.getElementById('runningStatus');
  if (runningStatus) {
    if (status.running) {
      runningStatus.textContent = t('main.running');
      runningStatus.classList.add('running');
      document.getElementById('startBtn').disabled = true;
      document.getElementById('stopBtn').disabled = false;
    } else {
      runningStatus.textContent = t('main.stopped');
      runningStatus.classList.remove('running');
      document.getElementById('startBtn').disabled = false;
      document.getElementById('stopBtn').disabled = true;
    }
  }

  const sentEl = document.getElementById('messagesSent');
  const failedEl = document.getElementById('failedCount');
  const nextEl = document.getElementById('nextSend');
  if (sentEl) sentEl.textContent = status.messages_sent || 0;

  if (status.next_send_unix) {
    nextSendUnix = status.next_send_unix;
  } else {
    nextSendUnix = 0;
  }
  isRunning = status.running;
  tickCountdown();

  if (failedEl) {
    failedEl.textContent = status.failed_count || 0;
    failedEl.style.color = status.failed_count ? 'var(--error)' : 'var(--kick-green)';
  }
}

function showStatus(elementId, message, type) {
  const el = document.getElementById(elementId);
  if (!el) return;
  el.textContent = message;
  el.className = `status-message show ${type}`;
  setTimeout(() => el.classList.remove('show'), 5000);
}

async function setChannel() {
  const channel = document.getElementById('channel').value.trim();
  if (!channel) { showStatus('channelStatus', t('status.enter_channel'), 'error'); return; }
  showStatus('channelStatus', t('status.validating'), 'success');
  try {
    const res = await fetch(`${API_BASE}/api/channel`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ channel }),
    });
    if (res.ok) {
      const data = await res.json();
      showStatus('channelStatus', t('status.channel_set', data.channel, data.id), 'success');
    } else {
      showStatus('channelStatus', t('status.err_prefix') + await res.text(), 'error');
    }
  } catch (e) {
    showStatus('channelStatus', t('status.conn_error', e.message), 'error');
  }
}

async function setMessage() {
  const message = document.getElementById('message').value.trim();
  if (!message) { showStatus('messageStatus', t('status.enter_message'), 'error'); return; }
  showStatus('messageStatus', t('status.saving'), 'success');
  try {
    const res = await fetch(`${API_BASE}/api/message`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ message }),
    });
    if (res.ok) {
      showStatus('messageStatus', t('status.message_set', message.length), 'success');
    } else {
      showStatus('messageStatus', t('status.err_prefix') + await res.text(), 'error');
    }
  } catch (e) {
    showStatus('messageStatus', t('status.conn_error', e.message), 'error');
  }
}

async function setInterval_() {
  const interval = document.getElementById('interval').value.trim();
  if (!interval) { showStatus('intervalStatus', t('status.enter_interval'), 'error'); return; }
  showStatus('intervalStatus', t('status.saving'), 'success');
  try {
    const res = await fetch(`${API_BASE}/api/interval`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ interval }),
    });
    if (res.ok) {
      showStatus('intervalStatus', t('status.interval_set', interval), 'success');
    } else {
      showStatus('intervalStatus', t('status.err_prefix') + await res.text(), 'error');
    }
  } catch (e) {
    showStatus('intervalStatus', t('status.conn_error', e.message), 'error');
  }
}

async function startScheduler() {
  try {
    const res = await fetch(`${API_BASE}/api/start`, { method: 'POST' });
    if (res.ok) {
      document.getElementById('startBtn').disabled = true;
      document.getElementById('stopBtn').disabled = false;
      addLog('INFO', 'Scheduler started');
    } else {
      addLog('ERROR', 'Failed to start: ' + await res.text());
    }
  } catch (e) {
    addLog('ERROR', 'Failed to start: ' + e.message);
  }
}

async function stopScheduler() {
  try {
    const res = await fetch(`${API_BASE}/api/stop`, { method: 'POST' });
    if (res.ok) {
      document.getElementById('startBtn').disabled = false;
      document.getElementById('stopBtn').disabled = true;
      addLog('INFO', 'Scheduler stopped');
    } else {
      addLog('ERROR', 'Failed to stop: ' + await res.text());
    }
  } catch (e) {
    addLog('ERROR', 'Failed to stop: ' + e.message);
  }
}

async function streamLogs() {
  const logsContainer = document.getElementById('logs');
  if (!logsContainer) return;

  try {
    const res = await fetch(`${API_BASE}/api/logs`, { headers: { 'Accept': 'text/event-stream' } });
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';
      for (const line of lines) {
        if (line.startsWith('data: ')) {
          try { addLog(JSON.parse(line.substring(6)).level, JSON.parse(line.substring(6)).message, JSON.parse(line.substring(6)).timestamp); } catch (_) {}
        }
      }
    }
  } catch (_) { setTimeout(streamLogs, 5000); }
}

function addLog(level, message, timestamp) {
  const logsContainer = document.getElementById('logs');
  if (!logsContainer) return;
  const emptyState = logsContainer.querySelector('.log-empty');
  if (emptyState) emptyState.remove();
  const entry = document.createElement('div');
  entry.className = 'log-entry';
  const time = timestamp ? new Date(timestamp).toLocaleTimeString() : new Date().toLocaleTimeString();
  entry.innerHTML = `<span class="timestamp">${time}</span><span class="level ${level}">${level}</span><span class="message">${escapeHtml(message)}</span>`;
  logsContainer.insertBefore(entry, logsContainer.firstChild);
  while (logsContainer.children.length > 100) logsContainer.removeChild(logsContainer.lastChild);
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

function tickCountdown() {
  const nextEl = document.getElementById('nextSend');
  if (!nextEl) return;
  if (!isRunning || !nextSendUnix) { nextEl.textContent = '-'; return; }
  const now = Math.floor(Date.now() / 1000);
  let remaining = nextSendUnix - now;
  if (remaining < 0) remaining = 0;
  const hours = Math.floor(remaining / 3600);
  const mins = Math.floor((remaining % 3600) / 60);
  const secs = remaining % 60;
  if (hours > 0) nextEl.textContent = `${hours}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
  else nextEl.textContent = `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

async function resetData(clearConfig, clearAuth) {
  let msg = '';
  if (clearConfig && clearAuth) msg = t('confirm.reset_all');
  else if (clearAuth) msg = t('confirm.clear_auth');
  else if (clearConfig) msg = t('confirm.clear_config');

  if (!confirm(msg)) return;

  try {
    const res = await fetch(`${API_BASE}/api/reset`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ clear_config: clearConfig, clear_auth: clearAuth }),
    });

    if (res.ok) {
      addLog('INFO', t('reset.success'));

      if (clearConfig) {
        document.getElementById('channel').value = '';
        document.getElementById('message').value = '';
        document.getElementById('interval').value = '';
        document.querySelectorAll('#mainApp .preset-btn').forEach(b => b.classList.remove('active'));
        ['channelStatus', 'messageStatus', 'intervalStatus'].forEach(id => {
          const el = document.getElementById(id);
          if (el) el.classList.remove('show');
        });
      }

      document.getElementById('settingsModal').classList.remove('show');

      if (clearAuth || (clearConfig && clearAuth)) {
        mainAppInited = false;
        onboarded = false;
        await checkAndRender();
      } else {
        updateMainStatus(await fetchStatus());
      }
    } else {
      alert(t('reset.failed', await res.text()));
    }
  } catch (e) {
    alert(t('reset.failed', e.message));
  }
}

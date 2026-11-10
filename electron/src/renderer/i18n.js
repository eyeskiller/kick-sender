// i18n: English, Czech, Slovak
// Usage: t('key'), t('key_with_arg', value), setLanguage('cs')
// HTML: <span data-i18n="key"></span>, <input data-i18n-placeholder="key">

const translations = {
  en: {
    // Onboarding tagline
    'ob.tagline': 'Let\u2019s get you set up',

    // Steps
    'ob.step.app': 'Kick App',
    'ob.step.creds': 'Credentials',
    'ob.step.auth': 'Auth',
    'ob.step.config': 'Configure',

    // Step 1
    'ob.s1.title': '1. Create a Kick Developer App',
    'ob.s1.desc': 'This app needs API access to your Kick account. You\u2019ll create a developer app to get that access.',
    'ob.s1.i1': 'Go to <a href="https://kick.com/settings/developer" target="_blank">kick.com/settings/developer</a>',
    'ob.s1.i2': 'Enable <strong>2FA</strong> on your Kick account (required for developer access)',
    'ob.s1.i3': 'Click <strong>"Create App"</strong> and give it any name',
    'ob.s1.i4': 'Set the <strong>Redirect URL</strong> to:',
    'ob.s1.i4.code': 'http://localhost:8765/callback',
    'ob.s1.i5': 'Copy your <strong>Client ID</strong> and <strong>Client Secret</strong> from the app page',
    'ob.s1.open_portal': 'Open Developer Portal',
    'ob.s1.next': 'I\u2019ve created my app',

    // Step 2
    'ob.s2.title': '2. Enter your credentials',
    'ob.s2.desc': 'Paste the Client ID and Client Secret from your Kick Developer App. These are stored locally on your computer only.',
    'ob.s2.client_id': 'Client ID',
    'ob.s2.client_id_ph': 'Paste your Client ID here',
    'ob.s2.client_secret': 'Client Secret',
    'ob.s2.client_secret_ph': 'Paste your Client Secret here',
    'ob.s2.redirect': 'Redirect URL',
    'ob.s2.redirect_opt': '(optional, defaults to localhost)',
    'ob.s2.redirect_ph': 'http://localhost:8765/callback',
    'ob.s2.save': 'Save & Continue',
    'ob.s2.error_empty': 'Please enter both Client ID and Client Secret.',
    'ob.s2.error_save': 'Failed to save: {0}',
    'ob.s2.error_conn': 'Connection error: {0}',

    // Step 3
    'ob.s3.title': '3. Authenticate with Kick',
    'ob.s3.desc': 'Click below to open your browser. Log in to Kick and authorize the app. You\u2019ll be redirected back automatically.',
    'ob.s3.ready': 'Ready to authenticate',
    'ob.s3.auth_btn': 'Open Browser & Authenticate',
    'ob.s3.opening': 'Opening browser...',
    'ob.s3.waiting': 'Waiting for authentication... Check your browser.',
    'ob.s3.wait_text': 'Waiting for authentication to complete...',
    'ob.s3.success': 'Authenticated!',
    'ob.s3.error': 'Error: {0}',

    // Step 4
    'ob.s4.title': '4. Configure your sender',
    'ob.s4.desc': 'Set up your channel, message, and interval. You can always change these later.',
    'ob.s4.channel': 'Target channel',
    'ob.s4.channel_ph': 'username or https://kick.com/username',
    'ob.s4.message': 'Message to send',
    'ob.s4.message_ph': 'e.g., Follow the stream and join our Discord!',
    'ob.s4.interval': 'Send interval',
    'ob.s4.interval_ph': 'or type custom: 60s, 5m, 10m, 1h',
    'ob.s4.error_channel': 'Please enter a channel.',
    'ob.s4.error_message': 'Please enter a message.',
    'ob.s4.error_field': '{0}: {1}',
    'ob.s4.finish': 'Start Using Kick Sender',

    // Presets
    'preset.1m': '1 min',
    'preset.5m': '5 min',
    'preset.10m': '10 min',
    'preset.15m': '15 min',
    'preset.30m': '30 min',
    'preset.1h': '1 hour',

    // Main app card headers
    'main.channel': 'Channel',
    'main.message': 'Message',
    'main.interval': 'Interval',
    'main.status': 'Status',
    'main.logs': 'Activity Log',

    // Main app placeholders
    'main.channel_ph': 'username or kick.com URL',
    'main.message_ph': 'Message to send...',
    'main.interval_ph': 'Custom: 60s, 5m, 10m...',

    // Help texts
    'main.channel_help': 'Kick channel to send messages to. Paste a URL or just the username.',
    'main.message_help': 'Repeated message. Max 500 characters, 2048 bytes.',
    'main.interval_help': 'How often to send. Min 60s.',

    // Buttons
    'btn.set': 'Set',

    // Status grid
    'main.state': 'State',
    'main.sent': 'Sent',
    'main.next_in': 'Next in',
    'main.failed': 'Failed',
    'main.running': 'Running',
    'main.stopped': 'Stopped',

    // Start/stop
    'btn.start': 'Start Sending',
    'btn.stop': 'Stop',

    // Log
    'logs.empty': 'No activity yet',

    // Footer
    'footer': 'Kick Sender v0.1.0 \u2022 Uses official Kick API',

    // Header auth
    'auth.authenticated': 'Authenticated',
    'auth.not_authenticated': 'Not authenticated',

    // Settings
    'settings.title': 'Settings',
    'settings.reset': 'Reset Options',
    'settings.clear_auth': 'Clear Authentication',
    'settings.clear_auth_desc': 'Remove OAuth tokens',
    'settings.clear_config': 'Clear Settings',
    'settings.clear_config_desc': 'Remove channel, message, interval',
    'settings.reset_all': 'Reset Everything',
    'settings.reset_all_desc': 'Clear all data including Kick Developer App credentials',
    'settings.about': 'About',
    'settings.version': 'Version',
    'settings.backend': 'Backend',
    'btn.clear': 'Clear',
    'btn.reset_all': 'Reset All',

    // Status messages (dynamic)
    'status.loaded': 'Loaded from saved configuration',
    'status.validating': 'Validating...',
    'status.saving': 'Saving...',
    'status.channel_set': 'Channel set: {0} (ID: {1})',
    'status.message_set': 'Set ({0} chars)',
    'status.interval_set': 'Interval: {0}',
    'status.enter_channel': 'Please enter a channel',
    'status.enter_message': 'Please enter a message',
    'status.enter_interval': 'Please enter an interval',
    'status.conn_error': 'Connection error: {0}',
    'status.err_prefix': 'Error: ',
    'btn.back_short': 'Back',

    // Confirm dialogs
    'confirm.reset_all': 'Reset everything? This clears all data and Kick Developer App credentials.',
    'confirm.clear_auth': 'Clear authentication? You will need to re-authenticate.',
    'confirm.clear_config': 'Clear saved settings?',
    'confirm.ok': 'OK',
    'confirm.cancel': 'Cancel',

    // Reset feedback
    'reset.success': 'Reset successful',
    'reset.failed': 'Reset failed: {0}',

    // Settings language
    'settings.language': 'Language',
  },

  cs: {
    'ob.tagline': 'Poj\u010fme to nastavit',

    'ob.step.app': 'Kick App',
    'ob.step.creds': 'P\u0159\u00EDstupov\u00E9 \u00FAdaje',
    'ob.step.auth': 'P\u0159ihl\u00E1\u0161en\u00ED',
    'ob.step.config': 'Nastaven\u00ED',

    'ob.s1.title': '1. Vytvo\u0159te Kick Developer App',
    'ob.s1.desc': 'Tato aplikace pot\u0159ebuje API p\u0159\u00EDstup k va\u0161emu Kick \u00FA\u010Dtu. Vytvo\u0159\u00EDte si developer aplikaci pro z\u00EDsk\u00E1n\u00ED tohoto p\u0159\u00EDstupu.',
    'ob.s1.i1': 'Jd\u011Bte na <a href="https://kick.com/settings/developer" target="_blank">kick.com/settings/developer</a>',
    'ob.s1.i2': 'Povolte <strong>2FA</strong> na sv\u00E9m Kick \u00FA\u010Dtu (vy\u017Eadov\u00E1no pro developer p\u0159\u00EDstup)',
    'ob.s1.i3': 'Klikn\u011Bte na <strong>"Create App"</strong> a zadejte libovoln\u00FD n\u00E1zev',
    'ob.s1.i4': 'Nastavte <strong>Redirect URL</strong> na:',
    'ob.s1.i4.code': 'http://localhost:8765/callback',
    'ob.s1.i5': 'Zkop\u00EDrujte sv\u00E9 <strong>Client ID</strong> a <strong>Client Secret</strong> ze str\u00E1nky aplikace',
    'ob.s1.open_portal': 'Otev\u0159\u00EDt Developer Portal',
    'ob.s1.next': 'Vytvo\u0159il jsem aplikaci',

    'ob.s2.title': '2. Zadejte p\u0159\u00EDstupov\u00E9 \u00FAdaje',
    'ob.s2.desc': 'Vlo\u017Ete Client ID a Client Secret z va\u0161\u00ED Kick Developer aplikace. Tyto \u00FAdaje jsou ulo\u017Eeny pouze lok\u00E1ln\u011B na va\u0161em po\u010D\u00EDta\u010Di.',
    'ob.s2.client_id': 'Client ID',
    'ob.s2.client_id_ph': 'Vlo\u017Ete sv\u00E9 Client ID',
    'ob.s2.client_secret': 'Client Secret',
    'ob.s2.client_secret_ph': 'Vlo\u017Ete sv\u00E9 Client Secret',
    'ob.s2.redirect': 'Redirect URL',
    'ob.s2.redirect_opt': '(voliteln\u00E9, v\u00FDchoz\u00ED je localhost)',
    'ob.s2.redirect_ph': 'http://localhost:8765/callback',
    'ob.s2.save': 'Ulo\u017Eit a pokra\u010Dovat',
    'ob.s2.error_empty': 'Vypl\u0148te pros\u00EDm Client ID i Client Secret.',
    'ob.s2.error_save': 'Ulo\u017Een\u00ED se nezda\u0159ilo: {0}',
    'ob.s2.error_conn': 'Chyba p\u0159ipojen\u00ED: {0}',

    'ob.s3.title': '3. P\u0159ihl\u00E1sit se ke Kick',
    'ob.s3.desc': 'Klikn\u011Bte n\u00ED\u017Ee pro otev\u0159en\u00ED prohl\u00ED\u017Ee\u010De. P\u0159ihl\u00E1ste se ke Kick a autorizujte aplikaci. Budete automaticky p\u0159esm\u011Brov\u00E1ni zp\u011Bt.',
    'ob.s3.ready': 'P\u0159ipraveno k p\u0159ihl\u00E1\u0161en\u00ED',
    'ob.s3.auth_btn': 'Otev\u0159\u00EDt prohl\u00ED\u017Ee\u010D a p\u0159ihl\u00E1sit se',
    'ob.s3.opening': 'Otev\u00EDr\u00E1m prohl\u00ED\u017Ee\u010D...',
    'ob.s3.waiting': '\u010Cek\u00E1m na p\u0159ihl\u00E1\u0161en\u00ED... Zkontrolujte prohl\u00ED\u017Ee\u010D.',
    'ob.s3.wait_text': '\u010Cek\u00E1m na dokon\u010Den\u00ED p\u0159ihl\u00E1\u0161en\u00ED...',
    'ob.s3.success': 'P\u0159ihl\u00E1\u0161eno!',
    'ob.s3.error': 'Chyba: {0}',

    'ob.s4.title': '4. Nastavte odes\u00EDl\u00E1n\u00ED',
    'ob.s4.desc': 'Nastavte kan\u00E1l, zpr\u00E1vu a interval. Tyto hodnoty m\u016F\u017Eete kdykoliv zm\u011Bnit.',
    'ob.s4.channel': 'C\u00EDlov\u00FD kan\u00E1l',
    'ob.s4.channel_ph': 'u\u017Eivatelsk\u00E9_jm\u00E9no nebo https://kick.com/u\u017Eivatel',
    'ob.s4.message': 'Zpr\u00E1va k odesl\u00E1n\u00ED',
    'ob.s4.message_ph': 'nap\u0159. Sleduj stream a p\u0159idej se na Discord!',
    'ob.s4.interval': 'Interval odes\u00EDl\u00E1n\u00ED',
    'ob.s4.interval_ph': 'nebo zadejte vlastn\u00ED: 60s, 5m, 10m, 1h',
    'ob.s4.error_channel': 'Zadejte pros\u00EDm kan\u00E1l.',
    'ob.s4.error_message': 'Zadejte pros\u00EDm zpr\u00E1vu.',
    'ob.s4.error_field': '{0}: {1}',
    'ob.s4.finish': 'Spustit Kick Sender',

    'preset.1m': '1 min',
    'preset.5m': '5 min',
    'preset.10m': '10 min',
    'preset.15m': '15 min',
    'preset.30m': '30 min',
    'preset.1h': '1 hodina',

    'main.channel': 'Kan\u00E1l',
    'main.message': 'Zpr\u00E1va',
    'main.interval': 'Interval',
    'main.status': 'Stav',
    'main.logs': 'Protokol aktivit',

    'main.channel_ph': 'u\u017Eivatelsk\u00E9_jm\u00E9no nebo kick.com URL',
    'main.message_ph': 'Zpr\u00E1va k odesl\u00E1n\u00ED...',
    'main.interval_ph': 'Vlastn\u00ED: 60s, 5m, 10m...',

    'main.channel_help': 'Kick kan\u00E1l pro odes\u00EDl\u00E1n\u00ED zpr\u00E1v. Vlo\u017Ete URL nebo jen u\u017Eivatelsk\u00E9 jm\u00E9no.',
    'main.message_help': 'Opakovan\u00E1 zpr\u00E1va. Max 500 znak\u016F, 2048 bajt\u016F.',
    'main.interval_help': 'Jak \u010Dasto odes\u00EDlat. Min 60s.',

    'btn.set': 'Nastavit',

    'main.state': 'Stav',
    'main.sent': 'Odesl\u00E1no',
    'main.next_in': 'Dal\u0161\u00ED za',
    'main.failed': 'Chyby',
    'main.running': 'B\u011B\u017E\u00ED',
    'main.stopped': 'Zastaveno',

    'btn.start': 'Spustit odes\u00EDl\u00E1n\u00ED',
    'btn.stop': 'Zastavit',

    'logs.empty': 'Zat\u00EDm \u017E\u00E1dn\u00E1 aktivita',

    'footer': 'Kick Sender v0.1.0 \u2022 Pou\u017E\u00EDv\u00E1 ofici\u00E1ln\u00ED Kick API',

    'auth.authenticated': 'P\u0159ihl\u00E1\u0161eno',
    'auth.not_authenticated': 'Nep\u0159ihl\u00E1\u0161eno',

    'settings.title': 'Nastaven\u00ED',
    'settings.reset': 'Mo\u017Enosti resetu',
    'settings.clear_auth': 'Smazat p\u0159ihl\u00E1\u0161en\u00ED',
    'settings.clear_auth_desc': 'Odstranit OAuth tokeny',
    'settings.clear_config': 'Smazat nastaven\u00ED',
    'settings.clear_config_desc': 'Odstranit kan\u00E1l, zpr\u00E1vu, interval',
    'settings.reset_all': 'Resetovat v\u0161e',
    'settings.reset_all_desc': 'Smazat v\u0161echna data v\u010Detn\u011B p\u0159\u00EDstupov\u00FDch \u00FAdaj\u016F Kick Developer App',
    'settings.about': 'O aplikaci',
    'settings.version': 'Verze',
    'settings.backend': 'Backend',
    'btn.clear': 'Smazat',
    'btn.reset_all': 'Resetovat',

    'status.loaded': 'Na\u010Dteno z ulo\u017Een\u00E9ho nastaven\u00ED',
    'status.validating': 'Ov\u011B\u0159uji...',
    'status.saving': 'Ukl\u00E1d\u00E1m...',
    'status.channel_set': 'Kan\u00E1l nastaven: {0} (ID: {1})',
    'status.message_set': 'Nastaveno ({0} znak\u016F)',
    'status.interval_set': 'Interval: {0}',
    'status.enter_channel': 'Zadejte pros\u00EDm kan\u00E1l',
    'status.enter_message': 'Zadejte pros\u00EDm zpr\u00E1vu',
    'status.enter_interval': 'Zadejte pros\u00EDm interval',
    'status.conn_error': 'Chyba p\u0159ipojen\u00ED: {0}',
    'status.err_prefix': 'Chyba: ',

    'confirm.reset_all': 'Resetovat v\u0161e? T\u00EDmto se sma\u017Eou v\u0161echna data a p\u0159\u00EDstupov\u00E9 \u00FAdaje Kick Developer App.',
    'confirm.clear_auth': 'Smazat p\u0159ihl\u00E1\u0161en\u00ED? Budete se muset znovu p\u0159ihl\u00E1sit.',
    'confirm.clear_config': 'Smazat ulo\u017Een\u00E9 nastaven\u00ED?',

    'confirm.ok': 'OK',
    'confirm.cancel': 'Zru\u0161it',

    'reset.success': 'Reset \u00FAsp\u011B\u0161n\u00FD',
    'reset.failed': 'Reset se nezda\u0159il: {0}',

    'settings.language': 'Jazyk',
    'btn.back_short': 'Zpět',
  },

  sk: {
    'ob.tagline': 'Po\u010fme to nastavi\u0165',

    'ob.step.app': 'Kick App',
    'ob.step.creds': 'Pr\u00EDstupov\u00E9 \u00FAdaje',
    'ob.step.auth': 'Prihl\u00E1senie',
    'ob.step.config': 'Nastavenie',

    'ob.s1.title': '1. Vytvorte Kick Developer App',
    'ob.s1.desc': 'T\u00E1to aplik\u00E1cia potrebuje API pr\u00EDstup k va\u0161emu Kick \u00FA\u010Dtu. Vytvor\u00EDte si developer aplik\u00E1ciu na z\u00EDskanie tohto pr\u00EDstupu.',
    'ob.s1.i1': 'Cho\u010Fte na <a href="https://kick.com/settings/developer" target="_blank">kick.com/settings/developer</a>',
    'ob.s1.i2': 'Povo\u013Ete <strong>2FA</strong> na svojom Kick \u00FA\u010Dte (vy\u017Eadovan\u00E9 pre developer pr\u00EDstup)',
    'ob.s1.i3': 'Kliknite na <strong>"Create App"</strong> a zadajte \u013Eubovoln\u00FD n\u00E1zov',
    'ob.s1.i4': 'Nastavte <strong>Redirect URL</strong> na:',
    'ob.s1.i4.code': 'http://localhost:8765/callback',
    'ob.s1.i5': 'Skop\u00EDrujte svoj <strong>Client ID</strong> a <strong>Client Secret</strong> zo str\u00E1nky aplik\u00E1cie',
    'ob.s1.open_portal': 'Otvori\u0165 Developer Portal',
    'ob.s1.next': 'Vytvoril som aplik\u00E1ciu',

    'ob.s2.title': '2. Zadajte pr\u00EDstupov\u00E9 \u00FAdaje',
    'ob.s2.desc': 'Vlo\u017Ete Client ID a Client Secret z va\u0161ej Kick Developer aplik\u00E1cie. Tieto \u00FAdaje s\u00FA ulo\u017Een\u00E9 len lok\u00E1lne na va\u0161om po\u010D\u00EDta\u010Di.',
    'ob.s2.client_id': 'Client ID',
    'ob.s2.client_id_ph': 'Vlo\u017Ete svoj Client ID',
    'ob.s2.client_secret': 'Client Secret',
    'ob.s2.client_secret_ph': 'Vlo\u017Ete svoj Client Secret',
    'ob.s2.redirect': 'Redirect URL',
    'ob.s2.redirect_opt': '(volite\u013En\u00E9, v\u00FDchodzie je localhost)',
    'ob.s2.redirect_ph': 'http://localhost:8765/callback',
    'ob.s2.save': 'Ulo\u017Ei\u0165 a pokra\u010Dova\u0165',
    'ob.s2.error_empty': 'Vypl\u0148te pros\u00EDm Client ID aj Client Secret.',
    'ob.s2.error_save': 'Ulo\u017Eenie zlyhalo: {0}',
    'ob.s2.error_conn': 'Chyba pripojenia: {0}',

    'ob.s3.title': '3. Prihl\u00E1si\u0165 sa ku Kick',
    'ob.s3.desc': 'Kliknite ni\u017E\u0161ie pre otvorenie prehliada\u010Da. Prihl\u00E1ste sa ku Kick a autorizujte aplik\u00E1ciu. Budete automaticky presmerovan\u00ED sp\u00E4\u0165.',
    'ob.s3.ready': 'Pripraven\u00E9 na prihl\u00E1senie',
    'ob.s3.auth_btn': 'Otvori\u0165 prehliada\u010D a prihl\u00E1si\u0165 sa',
    'ob.s3.opening': 'Otv\u00E1ram prehliada\u010D...',
    'ob.s3.waiting': '\u010Cak\u00E1m na prihl\u00E1senie... Skontrolujte prehliada\u010D.',
    'ob.s3.wait_text': '\u010Cak\u00E1m na dokon\u010Denie prihl\u00E1senia...',
    'ob.s3.success': 'Prihl\u00E1sen\u00E9!',
    'ob.s3.error': 'Chyba: {0}',

    'ob.s4.title': '4. Nastavte odosielanie',
    'ob.s4.desc': 'Nastavte kan\u00E1l, spr\u00E1vu a interval. Tieto hodnoty m\u00F4\u017Eete kedyko\u013Evek zmeni\u0165.',
    'ob.s4.channel': 'Cie\u013Eov\u00FD kan\u00E1l',
    'ob.s4.channel_ph': 'pou\u017E\u00EDvate\u013Esk\u00E9_meno alebo https://kick.com/pou\u017E\u00EDvate\u013E',
    'ob.s4.message': 'Spr\u00E1va na odoslanie',
    'ob.s4.message_ph': 'napr. Sleduj stream a pridaj sa na Discord!',
    'ob.s4.interval': 'Interval odosielania',
    'ob.s4.interval_ph': 'alebo zadajte vlastn\u00FD: 60s, 5m, 10m, 1h',
    'ob.s4.error_channel': 'Zadajte pros\u00EDm kan\u00E1l.',
    'ob.s4.error_message': 'Zadajte pros\u00EDm spr\u00E1vu.',
    'ob.s4.error_field': '{0}: {1}',
    'ob.s4.finish': 'Spusti\u0165 Kick Sender',

    'preset.1m': '1 min',
    'preset.5m': '5 min',
    'preset.10m': '10 min',
    'preset.15m': '15 min',
    'preset.30m': '30 min',
    'preset.1h': '1 hodina',

    'main.channel': 'Kan\u00E1l',
    'main.message': 'Spr\u00E1va',
    'main.interval': 'Interval',
    'main.status': 'Stav',
    'main.logs': 'Protokol aktiv\u00EDt',

    'main.channel_ph': 'pou\u017E\u00EDvate\u013Esk\u00E9_meno alebo kick.com URL',
    'main.message_ph': 'Spr\u00E1va na odoslanie...',
    'main.interval_ph': 'Vlastn\u00FD: 60s, 5m, 10m...',

    'main.channel_help': 'Kick kan\u00E1l na odosielanie spr\u00E1v. Vlo\u017Ete URL alebo len pou\u017E\u00EDvate\u013Esk\u00E9 meno.',
    'main.message_help': 'Opakovan\u00E1 spr\u00E1va. Max 500 znakov, 2048 bajtov.',
    'main.interval_help': 'Ako \u010Dasto odosiela\u0165. Min 60s.',

    'btn.set': 'Nastavi\u0165',

    'main.state': 'Stav',
    'main.sent': 'Odoslan\u00E9',
    'main.next_in': '\u010Eal\u0161ie za',
    'main.failed': 'Chyby',
    'main.running': 'Be\u017E\u00ED',
    'main.stopped': 'Zastaven\u00E9',

    'btn.start': 'Spusti\u0165 odosielanie',
    'btn.stop': 'Zastavi\u0165',

    'logs.empty': 'Zatia\u013E \u017Eiadna aktivita',

    'footer': 'Kick Sender v0.1.0 \u2022 Pou\u017E\u00EDva ofici\u00E1lne Kick API',

    'auth.authenticated': 'Prihl\u00E1sen\u00E9',
    'auth.not_authenticated': 'Neprihl\u00E1sen\u00E9',

    'settings.title': 'Nastavenia',
    'settings.reset': 'Mo\u017Enosti resetu',
    'settings.clear_auth': 'Vymaza\u0165 prihl\u00E1senie',
    'settings.clear_auth_desc': 'Odstr\u00E1ni\u0165 OAuth tokeny',
    'settings.clear_config': 'Vymaza\u0165 nastavenia',
    'settings.clear_config_desc': 'Odstr\u00E1ni\u0165 kan\u00E1l, spr\u00E1vu, interval',
    'settings.reset_all': 'Resetova\u0165 v\u0161etko',
    'settings.reset_all_desc': 'Vymaza\u0165 v\u0161etky d\u00E1ta vr\u00E1tane pr\u00EDstupov\u00FDch \u00FAdajov Kick Developer App',
    'settings.about': 'O aplik\u00E1cii',
    'settings.version': 'Verzia',
    'settings.backend': 'Backend',
    'btn.clear': 'Vymaza\u0165',
    'btn.reset_all': 'Resetova\u0165',

    'status.loaded': 'Na\u010D\u00EDtan\u00E9 z ulo\u017Een\u00E9ho nastavenia',
    'status.validating': 'Overujem...',
    'status.saving': 'Uklad\u00E1m...',
    'status.channel_set': 'Kan\u00E1l nastaven\u00FD: {0} (ID: {1})',
    'status.message_set': 'Nastaven\u00E9 ({0} znakov)',
    'status.interval_set': 'Interval: {0}',
    'status.enter_channel': 'Zadajte pros\u00EDm kan\u00E1l',
    'status.enter_message': 'Zadajte pros\u00EDm spr\u00E1vu',
    'status.enter_interval': 'Zadajte pros\u00EDm interval',
    'status.conn_error': 'Chyba pripojenia: {0}',
    'status.err_prefix': 'Chyba: ',

    'confirm.reset_all': 'Resetova\u0165 v\u0161etko? T\u00FDmto sa vyma\u017E\u00FA v\u0161etky d\u00E1ta a pr\u00EDstupov\u00E9 \u00FAdaje Kick Developer App.',
    'confirm.clear_auth': 'Vymaza\u0165 prihl\u00E1senie? Budete sa musie\u0165 znova prihl\u00E1si\u0165.',
    'confirm.clear_config': 'Vymaza\u0165 ulo\u017Een\u00E9 nastavenia?',

    'confirm.ok': 'OK',
    'confirm.cancel': 'Zru\u0161i\u0165',

    'reset.success': 'Reset \u00FAspe\u0161n\u00FD',
    'reset.failed': 'Reset zlyhal: {0}',

    'settings.language': 'Jazyk',
    'btn.back_short': 'Späť',
  },
};

let currentLang = 'en';

function detectLanguage() {
  const nav = navigator.language || 'en';
  const code = nav.split('-')[0].toLowerCase();
  if (code === 'cs') return 'cs';
  if (code === 'sk') return 'sk';
  return 'en';
}

// Format a key with {0}, {1}... placeholders
function t(key, ...args) {
  const raw = translations[currentLang]?.[key] ?? translations.en[key] ?? key;
  if (args.length === 0) return raw;
  return raw.replace(/\{(\d+)\}/g, (_, i) => args[parseInt(i, 10)] ?? '');
}

function setLanguage(lang) {
  if (!(lang in translations)) return;
  currentLang = lang;
  try { localStorage.setItem('kick-sender-lang', lang); } catch (_) {}
  applyTranslations();
}

function savedLanguage() {
  try { return localStorage.getItem('kick-sender-lang') || null; } catch (_) { return null; }
}

function initLanguage() {
  currentLang = savedLanguage() || detectLanguage();
  applyTranslations();
}

// Walk the DOM and apply translations to any element with data-i18n* attributes.
function applyTranslations() {
  document.querySelectorAll('[data-i18n]').forEach(el => {
    el.textContent = t(el.dataset.i18n);
  });
  document.querySelectorAll('[data-i18n-html]').forEach(el => {
    el.innerHTML = t(el.dataset.i18nHtml);
  });
  document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
    el.placeholder = t(el.dataset.i18nPlaceholder);
  });
  document.querySelectorAll('[data-i18n-title]').forEach(el => {
    el.title = t(el.dataset.i18nTitle);
  });

  // Update preset button text (their visible label comes from translations)
  document.querySelectorAll('.preset-btn[data-value]').forEach(btn => {
    const key = 'preset.' + btn.dataset.value;
    const val = t(key);
    if (val !== key) btn.textContent = val;
  });

  // Sync the language selector
  const sel = document.getElementById('langSelect');
  if (sel) sel.value = currentLang;
}

// Kick off as soon as this script loads (it runs before app.js, which calls applyTranslations
// again once it has initialized the main UI).
document.addEventListener('DOMContentLoaded', initLanguage);

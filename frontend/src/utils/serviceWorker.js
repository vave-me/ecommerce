/**
 * Service Worker Registration and Management
 * Mobile-optimized PWA utilities with intelligent update handling
 */
const isLocalhost = Boolean(
  typeof window !== 'undefined' &&
  (window.location.hostname === 'localhost' ||
  window.location.hostname === '[::1]' ||
  window.location.hostname.match(
    /^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/
  ))
);
/**
 * Register service worker with mobile-optimized configuration
 */
export function registerSW() {
  if (typeof window === 'undefined' || !('serviceWorker' in navigator)) {
    return;
  }
  const publicUrl = new URL(process.env.PUBLIC_URL || '', window.location.href);
  if (publicUrl.origin !== window.location.origin) {
    return;
  }
  window.addEventListener('load', () => {
    const swUrl = `${process.env.PUBLIC_URL || ''}/sw.js`;
    if (isLocalhost) {
      checkValidServiceWorker(swUrl);
      navigator.serviceWorker.ready.then(() => {
        if (process.env.NODE_ENV === 'development') {
          }
      });
    } else {
      registerValidSW(swUrl);
    }
  });
}
/**
 * Register valid service worker
 */
function registerValidSW(swUrl) {
  navigator.serviceWorker
    .register(swUrl)
    .then((registration) => {
      // Handle updates
      registration.addEventListener('updatefound', () => {
        const installingWorker = registration.installing;
        if (installingWorker == null) {
          return;
        }
        installingWorker.addEventListener('statechange', () => {
          if (installingWorker.state === 'installed') {
            if (navigator.serviceWorker.controller) {
              // New content is available; please refresh
              // Show update notification to user
              showUpdateNotification(registration);
            } else {
              // Content is cached for offline use
              showOfflineReadyNotification();
            }
          }
        });
      });
      // Check for updates periodically (mobile-friendly)
      setInterval(() => {
        registration.update();
      }, 60000); // Check every minute
    })
    .catch((error) => {
    });
}
/**
 * Check if service worker is valid
 */
function checkValidServiceWorker(swUrl) {
  fetch(swUrl, {
    headers: { 'Service-Worker': 'script' },
  })
    .then((response) => {
      const contentType = response.headers.get('content-type');
      if (
        response.status === 404 ||
        (contentType != null && contentType.indexOf('javascript') === -1)
      ) {
        navigator.serviceWorker.ready.then((registration) => {
          registration.unregister().then(() => {
            window.location.reload();
          });
        });
      } else {
        registerValidSW(swUrl);
      }
    })
    .catch(() => {
      });
}
/**
 * Unregister service worker
 */
export function unregisterSW() {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.ready
      .then((registration) => {
        registration.unregister();
      })
      .catch((error) => {
      });
  }
}
/**
 * Show update notification with mobile-friendly UI
 * SECURITY: Uses secure DOM creation instead of innerHTML
 */
function showUpdateNotification(registration) {
  // Import secure DOM utility
  import('./secureDom.js').then(({ createSecureNotification }) => {
    const notification = createSecureNotification({
      type: 'update',
      title: 'Update Available',
      message: 'New features and improvements are ready!',
      onAction: () => {
        if (registration.waiting) {
          registration.waiting.postMessage({ type: 'SKIP_WAITING' });
          window.location.reload();
        }
      },
      onDismiss: () => {
        notification.remove();
      }
    });
    document.body.appendChild(notification);
    // Auto-dismiss after 10 seconds
    setTimeout(() => {
      if (notification.parentNode) {
        notification.remove();
      }
    }, 10000);
  }).catch(error => {
  });
}
/**
 * Show offline ready notification
 */
function showOfflineReadyNotification() {
  const notification = document.createElement('div');
  notification.innerHTML = `
    <div style="
      position: fixed;
      bottom: 20px;
      left: 20px;
      right: 20px;
      background: #10b981;
      color: white;
      padding: 16px;
      border-radius: 12px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      z-index: 10000;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      animation: slideUp 0.3s ease-out;
    ">
      <div style="flex: 1; margin-right: 12px;">
        <div style="font-weight: 600; margin-bottom: 4px;">Ready for Offline</div>
        <div style="opacity: 0.9;">App is now available offline!</div>
      </div>
      <button onclick="this.parentElement.parentElement.remove()" style="
        background: transparent;
        color: white;
        border: none;
        padding: 8px;
        cursor: pointer;
        font-size: 18px;
        line-height: 1;
      ">
        ×
      </button>
    </div>
  `;
  document.body.appendChild(notification);
  // Auto-dismiss after 5 seconds
  setTimeout(() => {
    if (notification.parentNode) {
      notification.remove();
    }
  }, 5000);
}
/**
 * Check if app is running as PWA
 */
export function isPWA() {
  if (typeof window === 'undefined') {
    return false;
  }
  return window.matchMedia('(display-mode: standalone)').matches ||
         window.navigator.standalone === true ||
         document.referrer.includes('android-app://');
}
/**
 * Get PWA install prompt
 */
let deferredPrompt;
export function setupPWAInstallPrompt() {
  if (typeof window === 'undefined') {
    return;
  }
  window.addEventListener('beforeinstallprompt', (e) => {
    // Prevent Chrome 67 and earlier from automatically showing the prompt
    e.preventDefault();
    // Stash the event so it can be triggered later
    deferredPrompt = e;
    // Show custom install button/banner
    showInstallPrompt();
  });
  // Handle successful installation
  window.addEventListener('appinstalled', () => {
    deferredPrompt = null;
    hideInstallPrompt();
  });
}
/**
 * Show PWA install prompt
 */
function showInstallPrompt() {
  // Only show on mobile devices
  if (!isMobileDevice()) {
    return;
  }
  const installBanner = document.createElement('div');
  installBanner.id = 'pwa-install-banner';
  installBanner.innerHTML = `
    <div style="
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      background: linear-gradient(135deg, #2980b9, #7c3aed);
      color: white;
      padding: 12px 20px;
      z-index: 10000;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    ">
      <div style="flex: 1; margin-right: 12px;">
        <div style="font-weight: 600; margin-bottom: 2px;">Install sfx markt</div>
        <div style="opacity: 0.9; font-size: 12px;">Get the full app experience</div>
      </div>
      <button id="pwa-install-btn" style="
        background: white;
        color: #2980b9;
        border: none;
        padding: 8px 16px;
        border-radius: 8px;
        font-weight: 600;
        font-size: 14px;
        cursor: pointer;
        margin-right: 8px;
      ">
        Install
      </button>
      <button id="pwa-install-dismiss" style="
        background: transparent;
        color: white;
        border: none;
        padding: 8px;
        cursor: pointer;
        font-size: 18px;
        line-height: 1;
      ">
        ×
      </button>
    </div>
  `;
  document.body.appendChild(installBanner);
  // Handle install button click
  document.getElementById('pwa-install-btn').addEventListener('click', async () => {
    if (deferredPrompt) {
      deferredPrompt.prompt();
      const { outcome } = await deferredPrompt.userChoice;
      if (outcome === 'accepted') {
        } else {
        }
      deferredPrompt = null;
      hideInstallPrompt();
    }
  });
  // Handle dismiss button click
  document.getElementById('pwa-install-dismiss').addEventListener('click', () => {
    hideInstallPrompt();
    // Don't show again for 7 days
    localStorage.setItem('pwa-install-dismissed', Date.now().toString());
  });
  // Check if user previously dismissed
  const dismissed = localStorage.getItem('pwa-install-dismissed');
  if (dismissed && Date.now() - parseInt(dismissed) < 7 * 24 * 60 * 60 * 1000) {
    hideInstallPrompt();
  }
}
/**
 * Hide PWA install prompt
 */
function hideInstallPrompt() {
  const banner = document.getElementById('pwa-install-banner');
  if (banner) {
    banner.remove();
  }
}
/**
 * Check if device is mobile
 */
function isMobileDevice() {
  if (typeof window === 'undefined' || typeof navigator === 'undefined') {
    return false;
  }
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ||
         window.innerWidth <= 768;
}
/**
 * Get network status
 */
export function getNetworkStatus() {
  if (typeof window === 'undefined' || typeof navigator === 'undefined') {
    return {
      online: true,
      connection: null,
      effectiveType: 'unknown'
    };
  }
  return {
    online: navigator.onLine,
    connection: navigator.connection || navigator.mozConnection || navigator.webkitConnection,
    effectiveType: navigator.connection?.effectiveType || 'unknown'
  };
}
/**
 * Setup network status monitoring
 */
export function setupNetworkMonitoring() {
  if (typeof window === 'undefined') {
    return;
  }
  window.addEventListener('online', () => {
    showNetworkStatus('online');
  });
  window.addEventListener('offline', () => {
    showNetworkStatus('offline');
  });
}
/**
 * Show network status notification
 */
function showNetworkStatus(status) {
  const notification = document.createElement('div');
  const isOnline = status === 'online';
  notification.innerHTML = `
    <div style="
      position: fixed;
      top: 20px;
      left: 20px;
      right: 20px;
      background: ${isOnline ? '#10b981' : '#ef4444'};
      color: white;
      padding: 12px 16px;
      border-radius: 8px;
      z-index: 10000;
      text-align: center;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      font-weight: 600;
      animation: slideDown 0.3s ease-out;
    ">
      ${isOnline ? '🌐 Back Online' : '📱 Offline Mode'}
    </div>
    <style>
      @keyframes slideDown {
        from { transform: translateY(-100%); opacity: 0; }
        to { transform: translateY(0); opacity: 1; }
      }
    </style>
  `;
  document.body.appendChild(notification);
  // Auto-dismiss after 3 seconds
  setTimeout(() => {
    if (notification.parentNode) {
      notification.style.animation = 'slideDown 0.3s ease-out reverse';
      setTimeout(() => notification.remove(), 300);
    }
  }, 3000);
} 
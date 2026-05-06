/**
 * Service Worker for vave me PWA
 * Provides offline functionality, caching, and background sync
 */

const CACHE_NAME = 'vaveme-v1.0.0';
const OFFLINE_URL = '/offline.html';

// Critical assets to cache for immediate use
const CRITICAL_ASSETS = [
  '/',
  '/offline.html',
  '/manifest.json',
  '/favicon.ico'
];

// Assets that should be cached on first request
const RUNTIME_CACHE_PATTERNS = [
  // Images
  /\.(?:png|jpg|jpeg|svg|gif|webp|ico)$/,
  // Fonts
  /\.(?:woff|woff2|ttf|otf|eot)$/,
  // JavaScript and CSS
  /\.(?:js|css)$/,
  // API responses (selective)
  /^\/api\/(posts|products|properties|services|jobs)\//
];

// Network-first patterns (always try network first)
const NETWORK_FIRST_PATTERNS = [
  /^\/api\/auth\//,
  /^\/api\/users\//,
  /^\/api\/messages\//,
  /^\/api\/notifications\//
];

// Install event - cache critical assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(CRITICAL_ASSETS))
      .then(() => self.skipWaiting())
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames
            .filter((cacheName) => cacheName !== CACHE_NAME)
            .map((cacheName) => caches.delete(cacheName))
        );
      })
      .then(() => self.clients.claim())
  );
});

// Fetch event - handle all network requests
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return;
  }

  // Skip chrome-extension and other non-http(s) requests
  if (!url.protocol.startsWith('http')) {
    return;
  }

  // Handle different types of requests
  if (NETWORK_FIRST_PATTERNS.some(pattern => pattern.test(url.pathname))) {
    // Network first for real-time data
    event.respondWith(networkFirst(request));
  } else if (RUNTIME_CACHE_PATTERNS.some(pattern => pattern.test(url.pathname))) {
    // Cache first for static assets
    event.respondWith(cacheFirst(request));
  } else if (url.pathname.startsWith('/api/')) {
    // Stale while revalidate for API responses
    event.respondWith(staleWhileRevalidate(request));
  } else {
    // Default navigation handling
    event.respondWith(handleNavigation(request));
  }
});

// Handle background sync for offline actions
self.addEventListener('sync', (event) => {
  if (event.tag === 'background-sync') {
    event.waitUntil(backgroundSync());
  }
});

// Handle push notifications
self.addEventListener('push', (event) => {
  if (!event.data) return;

  const data = event.data.json();
  const options = {
    body: data.body,
    icon: '/logo192.png',
    badge: '/favicon.ico',
    data: data.data || {},
    actions: data.actions || []
  };

  event.waitUntil(
    self.registration.showNotification(data.title, options)
  );
});

// Handle notification clicks
self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  const action = event.action;
  const data = event.notification.data;

  let url = '/';
  if (action === 'view' && data.url) {
    url = data.url;
  } else if (data.url) {
    url = data.url;
  }

  event.waitUntil(
    clients.matchAll({ type: 'window' }).then((clientList) => {
      // Try to focus existing window
      for (const client of clientList) {
        if (client.url === url && 'focus' in client) {
          return client.focus();
        }
      }
      // Open new window
      if (clients.openWindow) {
        return clients.openWindow(url);
      }
    })
  );
});

// Handle messages from main thread
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});

// Caching strategies

/**
 * Safe response cloning utility
 */
function safeCloneResponse(response) {
  try {
    // Check if response body has been consumed
    if (response.bodyUsed) {
      console.warn('[SW] Response body already used, cannot clone');
      return null;
    }
    return response.clone();
  } catch (error) {
    console.warn('[SW] Failed to clone response:', error);
    return null;
  }
}

/**
 * Network first strategy - try network, fallback to cache
 */
async function networkFirst(request) {
  try {
    const response = await fetch(request);
    if (response.ok) {
      const responseClone = safeCloneResponse(response);
      if (responseClone) {
        const cache = await caches.open(CACHE_NAME);
        cache.put(request, responseClone);
      }
    }
    return response;
  } catch (error) {
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    throw error;
  }
}

/**
 * Cache first strategy - try cache, fallback to network
 */
async function cacheFirst(request) {
  const cachedResponse = await caches.match(request);
  if (cachedResponse) {
    return cachedResponse;
  }

  try {
    const response = await fetch(request);
    if (response.ok) {
      const responseClone = safeCloneResponse(response);
      if (responseClone) {
        const cache = await caches.open(CACHE_NAME);
        cache.put(request, responseClone);
      }
    }
    return response;
  } catch (error) {
    // Return offline fallback for images
    if (request.destination === 'image') {
      return new Response(
        '<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200"><rect width="200" height="200" fill="#f3f4f6"/><text x="100" y="100" text-anchor="middle" dy=".3em" font-family="Arial" font-size="14" fill="#9ca3af">Image offline</text></svg>',
        { headers: { 'Content-Type': 'image/svg+xml' } }
      );
    }
    throw error;
  }
}

/**
 * Stale while revalidate - return cache immediately, update in background
 */
async function staleWhileRevalidate(request) {
  const cachedResponse = await caches.match(request);
  
  const fetchPromise = fetch(request).then((response) => {
    if (response.ok) {
      const responseClone = safeCloneResponse(response);
      if (responseClone) {
        const cache = caches.open(CACHE_NAME);
        cache.then(c => c.put(request, responseClone));
      }
    }
    return response;
  }).catch(() => {
    // Network failed, return cached response if available
    return cachedResponse;
  });

  return cachedResponse || fetchPromise;
}

/**
 * Handle navigation requests
 */
async function handleNavigation(request) {
  try {
    // Try network first for navigation
    const response = await fetch(request);
    if (response.ok) {
      const responseClone = safeCloneResponse(response);
      if (responseClone) {
        const cache = await caches.open(CACHE_NAME);
        cache.put(request, responseClone);
      }
    }
    return response;
  } catch (error) {
    // Fallback to cache
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // Ultimate fallback to offline page
    const offlineResponse = await caches.match(OFFLINE_URL);
    if (offlineResponse) {
      return offlineResponse;
    }
    
    // Return a basic offline message
    return new Response(
      `<!DOCTYPE html>
      <html>
        <head>
          <title>Offline - sfx markt</title>
          <meta charset="utf-8">
          <meta name="viewport" content="width=device-width, initial-scale=1">
          <style>
            body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; text-align: center; padding: 50px; }
            h1 { color: #4f46e5; }
          </style>
        </head>
        <body>
          <h1>You're offline</h1>
          <p>Please check your internet connection and try again.</p>
          <button onclick="window.location.reload()">Retry</button>
        </body>
      </html>`,
      { 
        status: 200,
        headers: { 'Content-Type': 'text/html' }
      }
    );
  }
}

/**
 * Background sync for offline actions
 */
async function backgroundSync() {
  // Sync any pending offline actions
  const offlineActions = await getOfflineActions();
  
  for (const action of offlineActions) {
    try {
      await syncAction(action);
      await removeOfflineAction(action.id);
    } catch (error) {
      console.error('Failed to sync action:', action, error);
    }
  }
}

/**
 * Get offline actions from IndexedDB
 */
async function getOfflineActions() {
  // Simplified - in real app, use IndexedDB
  return [];
}

/**
 * Sync individual action
 */
async function syncAction(action) {
  // Implementation for syncing offline actions
  return Promise.resolve();
}

/**
 * Remove offline action after sync
 */
async function removeOfflineAction(actionId) {
  // Implementation for removing synced actions
  return Promise.resolve();
} 
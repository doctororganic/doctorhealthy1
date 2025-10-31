const CACHE_NAME = 'nutrition-platform-v2';
const STATIC_CACHE = 'static-v2';
const DYNAMIC_CACHE = 'dynamic-v2';
const API_CACHE = 'api-v2';

const urlsToCache = [
  '/',
  '/index.html',
  '/diet-planning.html',
  '/workout-generator.html',
  '/manifest.json',
  '/favicon.svg',
  '/icons/icon-72x72.svg',
  '/icons/icon-96x96.svg',
  '/icons/icon-128x128.svg',
  '/icons/icon-192x192.svg',
  '/icons/icon-512x512.svg',
  '/src/js/app.js',
  '/src/js/auth.js',
  '/src/js/language.js',
  '/src/js/diet-planning.js',
  '/src/js/workout-generator.js',
  '/src/js/validation.js',
  '/src/css/style.css',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.rtl.min.css',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js',
  'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css',
  'https://fonts.googleapis.com/css2?family=Noto+Sans+Arabic:wght@300;400;500;600;700&display=swap'
];

const API_ENDPOINTS = [
  '/api/v1/users',
  '/api/v1/foods',
  '/api/v1/health'
];

// Install event - cache resources
self.addEventListener('install', event => {
  console.log('Service Worker installing...');
  event.waitUntil(
    caches.open(STATIC_CACHE)
      .then(cache => {
        console.log('Opened static cache');
        return cache.addAll(urlsToCache);
      })
      .then(() => {
        console.log('Static resources cached successfully');
        return self.skipWaiting();
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', event => {
  console.log('Service Worker activating...');
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames.map(cacheName => {
          if (cacheName !== STATIC_CACHE && cacheName !== DYNAMIC_CACHE && cacheName !== API_CACHE) {
            console.log('Deleting old cache:', cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    }).then(() => {
      console.log('Service Worker activated');
      return self.clients.claim();
    })
  );
});

// Fetch event - serve from cache when offline
self.addEventListener('fetch', event => {
  const { request } = event;
  const url = new URL(request.url);

  // Handle API requests
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(handleApiRequest(request));
    return;
  }

  // Handle static assets
  if (request.destination === 'image' || request.destination === 'font' || 
      request.destination === 'style' || request.destination === 'script') {
    event.respondWith(handleStaticAssets(request));
    return;
  }

  // Handle navigation requests
  if (request.mode === 'navigate') {
    event.respondWith(handleNavigation(request));
    return;
  }

  // Default cache-first strategy
   event.respondWith(
     caches.match(request)
       .then(response => response || fetch(request))
       .catch(() => {
         // Return offline fallback
         if (request.destination === 'document') {
           return caches.match('/index.html');
         }
       })
   );
 });

 // Handle API requests with network-first strategy
 async function handleApiRequest(request) {
   try {
     const networkResponse = await fetch(request);
     if (networkResponse.ok) {
       const cache = await caches.open(API_CACHE);
       cache.put(request, networkResponse.clone());
     }
     return networkResponse;
   } catch (error) {
     console.log('Network failed, trying cache for API request');
     const cachedResponse = await caches.match(request);
     if (cachedResponse) {
       return cachedResponse;
     }
     // Return offline API response
     return new Response(JSON.stringify({
       error: true,
       message: 'Offline - data not available',
       offline: true
     }), {
       status: 503,
       headers: { 'Content-Type': 'application/json' }
     });
   }
 }

 // Handle static assets with cache-first strategy
 async function handleStaticAssets(request) {
   const cachedResponse = await caches.match(request);
   if (cachedResponse) {
     return cachedResponse;
   }

   try {
     const networkResponse = await fetch(request);
     if (networkResponse.ok) {
       const cache = await caches.open(DYNAMIC_CACHE);
       cache.put(request, networkResponse.clone());
     }
     return networkResponse;
   } catch (error) {
     console.log('Failed to fetch static asset:', request.url);
     throw error;
   }
 }

 // Handle navigation requests
 async function handleNavigation(request) {
   try {
     const networkResponse = await fetch(request);
     return networkResponse;
   } catch (error) {
     console.log('Navigation request failed, serving cached index.html');
     const cachedResponse = await caches.match('/index.html');
     return cachedResponse || new Response('Offline', { status: 503 });
   }
 }
self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames.map(cacheName => {
          if (cacheName !== CACHE_NAME) {
            console.log('Deleting old cache:', cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
});

// Background sync for offline actions
self.addEventListener('sync', event => {
  if (event.tag === 'background-sync') {
    event.waitUntil(doBackgroundSync());
  }
});

function doBackgroundSync() {
  // Handle offline actions when connection is restored
  return new Promise((resolve) => {
    // Implement sync logic here
    console.log('Background sync triggered');
    resolve();
  });
}

// Push notifications
self.addEventListener('push', event => {
  const options = {
    body: event.data ? event.data.text() : 'New notification from Nutrition Platform',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/icon-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    },
    actions: [
      {
        action: 'explore',
        title: 'View',
        icon: '/icons/icon-96x96.png'
      },
      {
        action: 'close',
        title: 'Close',
        icon: '/icons/icon-96x96.png'
      }
    ]
  };

  event.waitUntil(
    self.registration.showNotification('Nutrition Platform', options)
  );
});

// Notification click handling
self.addEventListener('notificationclick', event => {
  event.notification.close();

  if (event.action === 'explore') {
    event.waitUntil(
      clients.openWindow('/')
    );
  }
});
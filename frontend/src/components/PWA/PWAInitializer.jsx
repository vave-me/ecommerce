"use client";
import { useEffect, memo } from 'react';
import { 
  registerSW, 
  setupPWAInstallPrompt, 
  setupNetworkMonitoring,
  isPWA,
  getNetworkStatus 
} from '../../utils/serviceWorker';
import { initMobilePerformanceMonitoring } from '../../utils/mobilePerformance';
/**
 * PWA Initializer Component
 * Handles service worker registration and PWA features setup
 * Optimized for mobile-first experience with performance monitoring
 */
const PWAInitializer = memo(function PWAInitializer() {
  useEffect(() => {
    // Only initialize PWA features in browser environment
    if (typeof window === 'undefined') {
      return;
    }
    let mounted = true;
    const initializePWA = async () => {
      try {
        // 1. Register service worker
        if ('serviceWorker' in navigator) {
          registerSW();
        } else {
        }
        // 2. Initialize mobile performance monitoring
        initMobilePerformanceMonitoring();
        // 3. Setup PWA install prompt (mobile-optimized)
        setupPWAInstallPrompt();
        // 4. Setup network monitoring
        setupNetworkMonitoring();
        // 5. Log PWA status
        const pwaStatus = isPWA();
        const networkStatus = getNetworkStatus();
        // 6. Setup PWA-specific optimizations
        if (pwaStatus) {
          // Hide browser UI elements when running as PWA
          document.documentElement.style.setProperty('--pwa-mode', '1');
          // Add PWA-specific CSS class
          document.body.classList.add('pwa-mode');
          // Prevent zoom on double tap (iOS Safari)
          let lastTouchEnd = 0;
          document.addEventListener('touchend', (event) => {
            const now = (new Date()).getTime();
            if (now - lastTouchEnd <= 300) {
              event.preventDefault();
            }
            lastTouchEnd = now;
          }, false);
        }
        // 7. Setup viewport height fix for mobile browsers
        const setViewportHeight = () => {
          const vh = window.innerHeight * 0.01;
          document.documentElement.style.setProperty('--vh', `${vh}px`);
        };
        setViewportHeight();
        window.addEventListener('resize', setViewportHeight);
        window.addEventListener('orientationchange', () => {
          setTimeout(setViewportHeight, 100);
        });
        // 8. Setup performance monitoring for PWA
        if (process.env.NODE_ENV === 'development') {
          // Monitor PWA-specific metrics
          if ('performance' in window && 'PerformanceObserver' in window) {
            try {
              const observer = new PerformanceObserver((list) => {
                list.getEntries().forEach((entry) => {
                  if (entry.entryType === 'navigation') {
                    }
                });
              });
              observer.observe({ entryTypes: ['navigation'] });
            } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
          }
        }
        // 9. Setup app shortcuts (if supported)
        if ('getInstalledRelatedApps' in navigator) {
          try {
            const relatedApps = await navigator.getInstalledRelatedApps();
          } catch (e) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', e);
        }
        // Continue with default behavior
    }
        }
        // 10. Setup background sync registration (if supported)
        if ('serviceWorker' in navigator && 'sync' in window.ServiceWorkerRegistration.prototype) {
          navigator.serviceWorker.ready.then((registration) => {
            // Register for background sync
            return registration.sync.register('background-sync');
          }).catch((err) => {
            });
        }
        // 11. Setup Web Share API (mobile-optimized)
        if ('share' in navigator) {
          window.__PWA_SHARE_AVAILABLE__ = true;
          }
        // 12. Setup screen wake lock (for certain features)
        if ('wakeLock' in navigator) {
          window.__PWA_WAKE_LOCK_AVAILABLE__ = true;
          }
        // 13. Setup device orientation (mobile-specific)
        if ('DeviceOrientationEvent' in window) {
          // Request permission for iOS 13+
          if (typeof DeviceOrientationEvent.requestPermission === 'function') {
            // Will be requested when needed by specific features
            window.__PWA_ORIENTATION_PERMISSION_REQUIRED__ = true;
          }
        }
        // 14. Setup mobile-specific touch optimizations
        if ('ontouchstart' in window) {
          // Add touch-optimized class for CSS targeting
          document.body.classList.add('touch-device');
          // Optimize touch scrolling
          document.body.style.webkitOverflowScrolling = 'touch';
          // Prevent pull-to-refresh on mobile browsers (when not needed)
          document.body.addEventListener('touchstart', (e) => {
            if (e.touches.length > 1) {
              e.preventDefault();
            }
          }, { passive: false });
          let lastTouchEnd = 0;
          document.body.addEventListener('touchend', (e) => {
            const now = (new Date()).getTime();
            if (now - lastTouchEnd <= 300) {
              e.preventDefault();
            }
            lastTouchEnd = now;
          }, false);
        }
        // 15. Setup mobile keyboard handling
        if (/iPhone|iPad|iPod|Android/i.test(navigator.userAgent)) {
          let initialViewportHeight = window.innerHeight;
          window.addEventListener('resize', () => {
            const currentHeight = window.innerHeight;
            const heightDifference = initialViewportHeight - currentHeight;
            // Detect virtual keyboard
            if (heightDifference > 150) {
              document.body.classList.add('keyboard-open');
            } else {
              document.body.classList.remove('keyboard-open');
            }
          });
        }
        // 16. Setup offline data synchronization preparation
        if ('serviceWorker' in navigator) {
          navigator.serviceWorker.addEventListener('message', (event) => {
            if (event.data && event.data.type === 'CACHE_UPDATED') {
              // Optionally notify user of updated content
            }
          });
        }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    };
    // Initialize with a small delay to ensure DOM is ready
    const timeoutId = setTimeout(() => {
      if (mounted) {
        initializePWA();
      }
    }, 100);
    // Cleanup function
    return () => {
      mounted = false;
      clearTimeout(timeoutId);
    };
  }, []);
  // This component doesn't render anything
  return null;
});
export default PWAInitializer; 
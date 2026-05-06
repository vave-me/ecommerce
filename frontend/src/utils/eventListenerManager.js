/**
 * PRODUCTION EVENT LISTENER MANAGER
 * Prevents memory leaks from unhandled event listeners
 * Centralized cleanup and monitoring system
 */
class EventListenerManager {
  constructor() {
    this.listeners = new Map();
    this.debugMode = process.env.NODE_ENV === 'development';
    this.maxListeners = 100; // Safety limit
  }
  /**
   * Safe addEventListener with automatic cleanup tracking
   */
  addEventListener(target, event, handler, options = {}) {
    if (!target || typeof handler !== 'function') {
      if (this.debugMode) {
      }
      return null;
    }
    // Create unique key for this listener
    const key = this.generateKey(target, event, handler);
    // Check for duplicates
    if (this.listeners.has(key)) {
      if (this.debugMode) {
      }
      return key;
    }
    // Safety check for too many listeners
    if (this.listeners.size >= this.maxListeners) {
      if (this.debugMode) {
      }
      return null;
    }
    // Add the listener
    try {
      target.addEventListener(event, handler, options);
      // Store listener info for cleanup
      this.listeners.set(key, {
        target,
        event,
        handler,
        options,
        addedAt: Date.now(),
        component: this.getCurrentComponent()
      });
      if (this.debugMode) {
        }
      return key;
    } catch (error) {
      if (this.debugMode) {
      }
      return null;
    }
  }
  /**
   * Remove specific event listener
   */
  removeEventListener(key) {
    if (!this.listeners.has(key)) {
      if (this.debugMode) {
      }
      return false;
    }
    const listener = this.listeners.get(key);
    try {
      listener.target.removeEventListener(listener.event, listener.handler, listener.options);
      this.listeners.delete(key);
      if (this.debugMode) {
        }
      return true;
    } catch (error) {
      if (this.debugMode) {
      }
      return false;
    }
  }
  /**
   * Remove all listeners for a specific component
   */
  removeComponentListeners(componentName) {
    let removed = 0;
    for (const [key, listener] of this.listeners.entries()) {
      if (listener.component === componentName) {
        if (this.removeEventListener(key)) {
          removed++;
        }
      }
    }
    if (this.debugMode && removed > 0) {
      }
    return removed;
  }
  /**
   * Remove all listeners (emergency cleanup)
   */
  removeAllListeners() {
    const count = this.listeners.size;
    for (const key of this.listeners.keys()) {
      this.removeEventListener(key);
    }
    if (this.debugMode) {
      }
    return count;
  }
  /**
   * Get current active listeners count
   */
  getListenerCount() {
    return this.listeners.size;
  }
  /**
   * Get listeners by component
   */
  getComponentListeners(componentName) {
    const componentListeners = [];
    for (const [key, listener] of this.listeners.entries()) {
      if (listener.component === componentName) {
        componentListeners.push({ key, ...listener });
      }
    }
    return componentListeners;
  }
  /**
   * Performance monitoring
   */
  getPerformanceReport() {
    const now = Date.now();
    const report = {
      totalListeners: this.listeners.size,
      byComponent: {},
      byEvent: {},
      oldListeners: []
    };
    for (const [key, listener] of this.listeners.entries()) {
      // By component
      if (!report.byComponent[listener.component]) {
        report.byComponent[listener.component] = 0;
      }
      report.byComponent[listener.component]++;
      // By event type
      if (!report.byEvent[listener.event]) {
        report.byEvent[listener.event] = 0;
      }
      report.byEvent[listener.event]++;
      // Old listeners (potential leaks)
      const age = now - listener.addedAt;
      if (age > 300000) { // 5 minutes
        report.oldListeners.push({
          key,
          component: listener.component,
          event: listener.event,
          ageMinutes: Math.round(age / 60000)
        });
      }
    }
    return report;
  }
  /**
   * Generate unique key for listener
   */
  generateKey(target, event, handler) {
    const targetId = target === window ? 'window' : 
                    target === document ? 'document' :
                    target.id || target.tagName || 'unknown';
    const handlerName = handler.name || 'anonymous';
    const timestamp = Date.now();
    return `${targetId}-${event}-${handlerName}-${timestamp}`;
  }
  /**
   * Get current component name from call stack
   */
  getCurrentComponent() {
    try {
      const stack = new Error().stack;
      const lines = stack.split('\n');
      // Look for React component names in stack
      for (const line of lines) {
        const match = line.match(/at\s+([A-Z][A-Za-z0-9]*)/);
        if (match) {
          return match[1];
        }
      }
      return 'Unknown';
    } catch {
      return 'Unknown';
    }
  }
  /**
   * Periodic cleanup of old listeners
   */
  startPeriodicCleanup(intervalMs = 60000) { // 1 minute
    if (this.cleanupInterval) {
      clearInterval(this.cleanupInterval);
    }
    this.cleanupInterval = setInterval(() => {
      const report = this.getPerformanceReport();
      if (this.debugMode && report.oldListeners.length > 0) {
      }
      // Auto-cleanup listeners older than 10 minutes (safety measure)
      const tenMinutesAgo = Date.now() - 600000;
      let cleaned = 0;
      for (const [key, listener] of this.listeners.entries()) {
        if (listener.addedAt < tenMinutesAgo) {
          if (this.removeEventListener(key)) {
            cleaned++;
          }
        }
      }
      if (this.debugMode && cleaned > 0) {
        }
    }, intervalMs);
    return this.cleanupInterval;
  }
  /**
   * Stop periodic cleanup
   */
  stopPeriodicCleanup() {
    if (this.cleanupInterval) {
      clearInterval(this.cleanupInterval);
      this.cleanupInterval = null;
    }
  }
}
// Create singleton instance
const eventListenerManager = new EventListenerManager();
// Auto-start periodic cleanup in development
if (process.env.NODE_ENV === 'development') {
  eventListenerManager.startPeriodicCleanup();
}
// Global cleanup on page unload
if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', () => {
    eventListenerManager.removeAllListeners();
    eventListenerManager.stopPeriodicCleanup();
  });
}
export default eventListenerManager; 
export class MessageDeduplicator<T> {
  private messageCache: Map<string, { timestamp: number; data: T }>;
  private windowMs: number;
  private maxSize: number;
  private getKey: (message: T) => string;

  constructor(options: {
    windowMs?: number;
    maxSize?: number;
    getKey?: (message: T) => string;
  } = {}) {
    this.messageCache = new Map();
    this.windowMs = options.windowMs || 5000; // 5 seconds default
    this.maxSize = options.maxSize || 1000;
    this.getKey = options.getKey || ((msg: T) => {
      if (typeof msg === 'object' && msg !== null && 'id' in msg) {
        return String((msg as any).id);
      }
      return JSON.stringify(msg);
    });
  }

  isDuplicate(message: T): boolean {
    const key = this.getKey(message);
    const now = Date.now();
    
    // Clean up old entries
    this.cleanup(now);
    
    // Check if we've seen this message
    if (this.messageCache.has(key)) {
      const cached = this.messageCache.get(key)!;
      if (now - cached.timestamp < this.windowMs) {
        return true;
      }
    }
    
    // Add to cache
    this.messageCache.set(key, { timestamp: now, data: message });
    
    // Enforce size limit
    if (this.messageCache.size > this.maxSize) {
      const oldestKey = this.messageCache.keys().next().value;
      if (oldestKey !== undefined) {
        this.messageCache.delete(oldestKey);
      }
    }
    
    return false;
  }

  private cleanup(now: number) {
    for (const [key, value] of this.messageCache.entries()) {
      if (now - value.timestamp > this.windowMs) {
        this.messageCache.delete(key);
      } else {
        // Map maintains insertion order, so once we find a recent entry,
        // all subsequent entries are also recent
        break;
      }
    }
  }

  clear() {
    this.messageCache.clear();
  }

  size(): number {
    return this.messageCache.size;
  }
}
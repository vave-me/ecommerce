export interface RetryOptions {
  maxAttempts?: number;
  initialDelayMs?: number;
  maxDelayMs?: number;
  backoffMultiplier?: number;
  shouldRetry?: (error: Error, attempt: number) => boolean;
  onRetry?: (error: Error, attempt: number) => void;
}

export async function withRetry<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> {
  const {
    maxAttempts = 3,
    initialDelayMs = 1000,
    maxDelayMs = 30000,
    backoffMultiplier = 2,
    shouldRetry = () => true,
    onRetry
  } = options;

  let lastError: Error;
  
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error as Error;
      
      if (attempt === maxAttempts || !shouldRetry(lastError, attempt)) {
        throw lastError;
      }
      
      onRetry?.(lastError, attempt);
      
      const delay = Math.min(
        initialDelayMs * Math.pow(backoffMultiplier, attempt - 1),
        maxDelayMs
      );
      
      await sleep(delay);
    }
  }
  
  throw lastError!;
}

export function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

export class ExponentialBackoff {
  private attempt: number = 0;
  private readonly initialDelayMs: number;
  private readonly maxDelayMs: number;
  private readonly multiplier: number;
  private readonly jitterMs: number;

  constructor(options: {
    initialDelayMs?: number;
    maxDelayMs?: number;
    multiplier?: number;
    jitterMs?: number;
  } = {}) {
    this.initialDelayMs = options.initialDelayMs || 1000;
    this.maxDelayMs = options.maxDelayMs || 30000;
    this.multiplier = options.multiplier || 2;
    this.jitterMs = options.jitterMs || 0;
  }

  nextDelay(): number {
    const baseDelay = Math.min(
      this.initialDelayMs * Math.pow(this.multiplier, this.attempt),
      this.maxDelayMs
    );
    
    const jitter = this.jitterMs > 0
      ? Math.random() * this.jitterMs * 2 - this.jitterMs
      : 0;
    
    this.attempt++;
    
    return Math.max(0, baseDelay + jitter);
  }

  reset() {
    this.attempt = 0;
  }

  async wait() {
    const delay = this.nextDelay();
    await sleep(delay);
  }
}
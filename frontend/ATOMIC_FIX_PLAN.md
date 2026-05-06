# Atomic Fix Plan for Production Readiness

## Phase 1: Critical Security Fixes (Day 1-2)

### 1.1 Remove Console Statements
```bash
# Create utility to remove console statements
npm install -D @babel/plugin-transform-remove-console

# Update babel config
# babel.config.js
module.exports = {
  presets: ['next/babel'],
  plugins: [
    ['transform-remove-console', { exclude: ['error', 'warn'] }]
  ]
};
```

### 1.2 Fix dangerouslySetInnerHTML Usage
```javascript
// Create safe HTML renderer
// src/utils/safeHtml.js
import DOMPurify from 'dompurify';

export const sanitizeHtml = (html) => {
  if (typeof window === 'undefined') return '';
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'a', 'ul', 'ol', 'li'],
    ALLOWED_ATTR: ['href', 'target', 'rel']
  });
};

// Replace all instances
// Before: dangerouslySetInnerHTML={{ __html: content }}
// After: dangerouslySetInnerHTML={{ __html: sanitizeHtml(content) }}
```

### 1.3 Fix Authentication Token Storage
```javascript
// src/utils/secureAuth.js
import Cookies from 'js-cookie';

export const secureTokenStorage = {
  setToken: (token) => {
    Cookies.set('auth_token', token, {
      secure: true,
      httpOnly: false, // Can't use httpOnly in client
      sameSite: 'strict',
      expires: 7
    });
  },
  
  getToken: () => {
    return Cookies.get('auth_token');
  },
  
  removeToken: () => {
    Cookies.remove('auth_token');
  }
};
```

## Phase 2: Error Handling (Day 3-4)

### 2.1 Create Global Error Handler
```javascript
// src/utils/errorHandler.js
import * as Sentry from '@sentry/nextjs';

export class AppError extends Error {
  constructor(message, code, statusCode) {
    super(message);
    this.code = code;
    this.statusCode = statusCode;
    this.isOperational = true;
  }
}

export const handleError = (error, context = {}) => {
  // Log to Sentry
  Sentry.captureException(error, {
    extra: context
  });
  
  // Log in development
  if (process.env.NODE_ENV === 'development') {
    console.error('Error:', error);
  }
  
  // Return user-friendly message
  if (error.isOperational) {
    return {
      message: error.message,
      code: error.code
    };
  }
  
  return {
    message: 'An unexpected error occurred',
    code: 'INTERNAL_ERROR'
  };
};
```

### 2.2 Fix Empty Catch Blocks
```javascript
// Create standard error handling pattern
// Before:
try {
  // code
} catch (error) {
  // empty
}

// After:
try {
  // code
} catch (error) {
  handleError(error, { component: 'ComponentName' });
  // Show user notification if needed
  showErrorNotification('Operation failed. Please try again.');
}
```

### 2.3 Add Error Boundaries
```javascript
// src/components/ErrorBoundary/GlobalErrorBoundary.jsx
import React from 'react';
import * as Sentry from '@sentry/nextjs';

class GlobalErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }
  
  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }
  
  componentDidCatch(error, errorInfo) {
    Sentry.captureException(error, {
      contexts: {
        react: {
          componentStack: errorInfo.componentStack
        }
      }
    });
  }
  
  render() {
    if (this.state.hasError) {
      return (
        <div className="error-fallback">
          <h1>Something went wrong</h1>
          <button onClick={() => window.location.reload()}>
            Reload Page
          </button>
        </div>
      );
    }
    
    return this.props.children;
  }
}
```

## Phase 3: Memory Leak Fixes (Day 5-6)

### 3.1 Create Cleanup Manager
```javascript
// src/utils/cleanupManager.js
export class CleanupManager {
  constructor() {
    this.cleanups = new Set();
  }
  
  add(cleanup) {
    this.cleanups.add(cleanup);
    return () => this.cleanups.delete(cleanup);
  }
  
  cleanup() {
    this.cleanups.forEach(fn => fn());
    this.cleanups.clear();
  }
}

// Usage in components
const cleanup = new CleanupManager();

useEffect(() => {
  const timer = setTimeout(() => {}, 1000);
  cleanup.add(() => clearTimeout(timer));
  
  const listener = () => {};
  window.addEventListener('resize', listener);
  cleanup.add(() => window.removeEventListener('resize', listener));
  
  return () => cleanup.cleanup();
}, []);
```

### 3.2 Fix Event Listeners
```javascript
// src/hooks/useEventListener.js
import { useEffect, useRef } from 'react';

export function useEventListener(eventName, handler, element = window) {
  const savedHandler = useRef();
  
  useEffect(() => {
    savedHandler.current = handler;
  }, [handler]);
  
  useEffect(() => {
    const isSupported = element && element.addEventListener;
    if (!isSupported) return;
    
    const eventListener = (event) => savedHandler.current(event);
    element.addEventListener(eventName, eventListener);
    
    return () => {
      element.removeEventListener(eventName, eventListener);
    };
  }, [eventName, element]);
}
```

## Phase 4: Performance Optimization (Day 7-8)

### 4.1 Implement Code Splitting
```javascript
// src/utils/lazyLoad.js
import { lazy, Suspense } from 'react';
import LoadingSpinner from '@/components/LoadingSpinner';

export const lazyLoadComponent = (importFunc, fallback = <LoadingSpinner />) => {
  const LazyComponent = lazy(importFunc);
  
  return (props) => (
    <Suspense fallback={fallback}>
      <LazyComponent {...props} />
    </Suspense>
  );
};

// Usage
const LazyModal = lazyLoadComponent(
  () => import('@/components/Modal/UnifiedModal')
);
```

### 4.2 Optimize Images
```javascript
// src/components/OptimizedImage.jsx
import Image from 'next/image';
import { useState } from 'react';

export function OptimizedImage({ src, alt, ...props }) {
  const [isLoading, setLoading] = useState(true);
  
  return (
    <div className="image-container">
      <Image
        src={src}
        alt={alt}
        loading="lazy"
        placeholder="blur"
        blurDataURL="data:image/jpeg;base64,/9j/4AAQSkZJRg..."
        onLoadingComplete={() => setLoading(false)}
        {...props}
      />
      {isLoading && <div className="image-skeleton" />}
    </div>
  );
}
```

### 4.3 Add React.memo to Heavy Components
```javascript
// src/components/Feed/FeedItem.jsx
import React, { memo } from 'react';

const FeedItem = memo(({ item, onInteraction }) => {
  // Component implementation
}, (prevProps, nextProps) => {
  // Custom comparison
  return prevProps.item.id === nextProps.item.id &&
         prevProps.item.updatedAt === nextProps.item.updatedAt;
});
```

## Phase 5: Type Safety (Day 9-10)

### 5.1 Create TypeScript Migration Plan
```typescript
// tsconfig.json
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true
  }
}
```

### 5.2 Add Type Definitions
```typescript
// src/types/index.ts
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'user' | 'admin' | 'vendor';
}

export interface Product {
  id: string;
  title: string;
  price: number;
  description: string;
  images: string[];
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
  status: number;
}
```

## Phase 6: Testing Implementation (Day 11-12)

### 6.1 Setup Testing Framework
```javascript
// jest.config.js
module.exports = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/src/tests/setup.js'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
  },
  collectCoverageFrom: [
    'src/**/*.{js,jsx}',
    '!src/**/*.test.{js,jsx}',
    '!src/tests/**/*',
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
};
```

### 6.2 Add Critical Tests
```javascript
// src/components/__tests__/Header.test.jsx
import { render, screen, fireEvent } from '@testing-library/react';
import Header from '../Header/Header';

describe('Header', () => {
  it('renders navigation items', () => {
    render(<Header />);
    expect(screen.getByText('Home')).toBeInTheDocument();
  });
  
  it('handles search submission', () => {
    const onSearch = jest.fn();
    render(<Header onSearch={onSearch} />);
    
    const input = screen.getByPlaceholderText('Search...');
    fireEvent.change(input, { target: { value: 'test' } });
    fireEvent.submit(input.closest('form'));
    
    expect(onSearch).toHaveBeenCalledWith('test');
  });
});
```

## Phase 7: Monitoring Setup (Day 13-14)

### 7.1 Implement Sentry
```javascript
// sentry.client.config.js
import * as Sentry from '@sentry/nextjs';

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  tracesSampleRate: 0.1,
  debug: false,
  replaysOnErrorSampleRate: 1.0,
  replaysSessionSampleRate: 0.1,
  integrations: [
    new Sentry.Replay({
      maskAllText: true,
      blockAllMedia: true,
    }),
  ],
});
```

### 7.2 Add Performance Monitoring
```javascript
// src/utils/performance.js
export const measurePerformance = (metricName, fn) => {
  const startTime = performance.now();
  const result = fn();
  const endTime = performance.now();
  
  if (window.gtag) {
    window.gtag('event', 'timing_complete', {
      name: metricName,
      value: Math.round(endTime - startTime),
    });
  }
  
  return result;
};
```

## Implementation Schedule

### Week 1
- Day 1-2: Security fixes
- Day 3-4: Error handling
- Day 5-6: Memory leaks

### Week 2  
- Day 7-8: Performance optimization
- Day 9-10: Type safety
- Day 11-12: Testing

### Week 3
- Day 13-14: Monitoring setup
- Day 15-16: SEO improvements
- Day 17-18: Final testing

### Week 4
- Day 19-20: Bug fixes
- Day 21: Production deployment

## Success Metrics

1. **Security**: 0 console logs in production, all inputs sanitized
2. **Performance**: Core Web Vitals all green
3. **Reliability**: < 0.1% error rate
4. **Testing**: > 80% code coverage
5. **Monitoring**: Real-time error tracking active

## Rollback Plan

1. Keep current version tagged
2. Blue-green deployment strategy
3. Feature flags for new changes
4. Quick rollback procedure documented
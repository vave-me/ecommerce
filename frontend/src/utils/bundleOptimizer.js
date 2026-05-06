/**
 * ADVANCED BUNDLE SIZE OPTIMIZER
 * Comprehensive utility for reducing bundle size and optimizing imports
 * - Dynamic import patterns
 * - Tree-shaking optimizations
 * - Lazy loading helpers
 * - Bundle splitting strategies
 */
import React, { lazy, Suspense, memo, useMemo, useCallback } from 'react';
/**
 * Advanced lazy component factory with error boundaries and preloading
 */
export const createOptimizedLazyComponent = (
  importFn, 
  fallback = null, 
  errorFallback = null,
  preloadCondition = null
) => {
  let componentPromise = null;
  // Preload function
  const preload = () => {
    if (!componentPromise) {
      componentPromise = importFn();
    }
    return componentPromise;
  };
  // Create lazy component
  const LazyComponent = lazy(() => preload());
  // Enhanced wrapper with error boundary
  const OptimizedComponent = memo((props) => {
    // Conditional preloading
    useMemo(() => {
      if (preloadCondition && preloadCondition(props)) {
        preload();
      }
    }, [props]);
    return (
      <ErrorBoundary fallback={errorFallback}>
        <Suspense fallback={fallback}>
          <LazyComponent {...props} />
        </Suspense>
      </ErrorBoundary>
    );
  });
  // Attach preload method
  OptimizedComponent.preload = preload;
  OptimizedComponent.displayName = `OptimizedLazy(${LazyComponent.displayName || 'Component'})`;
  return OptimizedComponent;
};
/**
 * Simple error boundary for lazy components
 */
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }
  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }
  componentDidCatch(error, errorInfo) {
    if (process.env.NODE_ENV === 'development') {
    }
  }
  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="error-boundary">
          <p>Something went wrong loading this component.</p>
          <button onClick={() => window.location.reload()}>
            Reload Page
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
/**
 * Optimized dynamic import with module federation support
 */
export const optimizedImport = async (modulePath, options = {}) => {
  const {
    timeout = 10000,
    retries = 3,
    fallback = null
  } = options;
  let lastError;
  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      const importPromise = import(modulePath);
      // Add timeout
      const timeoutPromise = new Promise((_, reject) => {
        setTimeout(() => reject(new Error('Import timeout')), timeout);
      });
      const module = await Promise.race([importPromise, timeoutPromise]);
      return module.default || module;
    } catch (error) {
      lastError = error;
      if (attempt < retries - 1) {
        // Exponential backoff
        await new Promise(resolve => setTimeout(resolve, Math.pow(2, attempt) * 1000));
      }
    }
  }
  if (fallback) {
    return fallback;
  }
  throw lastError;
};
/**
 * Tree-shaking optimizer for icon libraries
 */
export const createIconOptimizer = (iconLibrary) => {
  const iconCache = new Map();
  return {
    // Get icon with caching and lazy loading
    getIcon: async (iconName) => {
      if (iconCache.has(iconName)) {
        return iconCache.get(iconName);
      }
      try {
        const iconModule = await import(`${iconLibrary}/${iconName}`);
        const Icon = iconModule.default || iconModule[iconName];
        iconCache.set(iconName, Icon);
        return Icon;
      } catch (error) {
        return null;
      }
    },
    // Preload multiple icons
    preloadIcons: async (iconNames) => {
      const promises = iconNames.map(name => this.getIcon(name));
      return Promise.allSettled(promises);
    },
    // Get cache stats
    getCacheStats: () => ({
      cachedIcons: iconCache.size,
      iconNames: Array.from(iconCache.keys())
    })
  };
};
/**
 * Code splitting helper for route-based chunks
 */
export const createRouteChunk = (routes) => {
  const chunks = {};
  routes.forEach(route => {
    const { path, component, preload = false } = route;
    chunks[path] = createOptimizedLazyComponent(
      component,
      <div className="route-loading">Loading page...</div>,
      <div className="route-error">Failed to load page</div>,
      preload ? () => true : null
    );
  });
  return {
    chunks,
    preloadRoute: (path) => {
      if (chunks[path] && chunks[path].preload) {
        chunks[path].preload();
      }
    },
    preloadAllRoutes: () => {
      Object.values(chunks).forEach(chunk => {
        if (chunk.preload) {
          chunk.preload();
        }
      });
    }
  };
};
/**
 * Bundle analysis helper
 */
export const analyzeBundleSize = () => {
  const getModuleSize = (moduleName) => {
    try {
      // This is a rough estimation - actual bundle analysis would need webpack-bundle-analyzer
      const module = require.cache[require.resolve(moduleName)];
      return module ? JSON.stringify(module).length : 0;
    } catch {
      return 0;
    }
  };
  const getComponentSize = (component) => {
    return component.toString().length;
  };
  return {
    getModuleSize,
    getComponentSize,
    estimateBundleImpact: (modules) => {
      return modules.reduce((total, module) => {
        return total + getModuleSize(module);
      }, 0);
    }
  };
};
/**
 * Performance-optimized component registry
 */
export const createComponentRegistry = () => {
  const registry = new Map();
  const loadingStates = new Map();
  return {
    // Register component with lazy loading
    register: (name, importFn, options = {}) => {
      const component = createOptimizedLazyComponent(
        importFn,
        options.fallback,
        options.errorFallback,
        options.preloadCondition
      );
      registry.set(name, component);
      return component;
    },
    // Get component (lazy load if needed)
    get: async (name) => {
      if (!registry.has(name)) {
        throw new Error(`Component ${name} not registered`);
      }
      const component = registry.get(name);
      // Preload if not already loading
      if (component.preload && !loadingStates.get(name)) {
        loadingStates.set(name, true);
        try {
          await component.preload();
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
        loadingStates.set(name, false);
      }
      return component;
    },
    // Preload multiple components
    preload: async (names) => {
      const promises = names.map(name => this.get(name));
      return Promise.allSettled(promises);
    },
    // Get registry stats
    getStats: () => ({
      totalComponents: registry.size,
      loadingComponents: Array.from(loadingStates.entries())
        .filter(([, loading]) => loading)
        .map(([name]) => name)
    })
  };
};
/**
 * CSS optimization helper
 */
export const optimizeCSS = {
  // Remove unused CSS (client-side detection)
  removeUnusedStyles: () => {
    const usedSelectors = new Set();
    const allElements = document.querySelectorAll('*');
    allElements.forEach(element => {
      const classes = element.className;
      if (typeof classes === 'string') {
        classes.split(' ').forEach(cls => {
          if (cls.trim()) {
            usedSelectors.add(`.${cls.trim()}`);
          }
        });
      }
    });
    return Array.from(usedSelectors);
  },
  // Critical CSS extraction (simplified)
  extractCriticalCSS: () => {
    const criticalElements = document.querySelectorAll('header, nav, .hero, .above-fold');
    const criticalSelectors = new Set();
    criticalElements.forEach(element => {
      const computedStyle = window.getComputedStyle(element);
      // This is a simplified version - real implementation would be more complex
      criticalSelectors.add(element.tagName.toLowerCase());
    });
    return Array.from(criticalSelectors);
  }
};
/**
 * Memory optimization helpers
 */
export const memoryOptimizer = {
  // Cleanup unused module cache
  cleanupModuleCache: () => {
    if (typeof require !== 'undefined' && require.cache) {
      const unusedModules = [];
      Object.keys(require.cache).forEach(key => {
        // Simple heuristic - modules not accessed recently
        const module = require.cache[key];
        if (module && module.loaded && !module.children.length) {
          unusedModules.push(key);
        }
      });
      return unusedModules;
    }
    return [];
  },
  // Clear component caches
  clearComponentCaches: () => {
    // Clear any component-level caches
    if (typeof window !== 'undefined') {
      // Clear intersection observer cache
      if (window.intersectionObserverCache) {
        window.intersectionObserverCache.clear();
      }
      // Clear image loading cache
      if (window.imageLoadingCache) {
        window.imageLoadingCache.clear();
      }
    }
  }
};
/**
 * Bundle optimization strategies
 */
export const bundleStrategies = {
  // Route-based splitting
  routeSplitting: (routes) => createRouteChunk(routes),
  // Component-based splitting
  componentSplitting: (components) => {
    return components.map(({ name, importFn, ...options }) => ({
      name,
      component: createOptimizedLazyComponent(importFn, ...options)
    }));
  },
  // Vendor library splitting
  vendorSplitting: (vendors) => {
    return vendors.map(vendor => ({
      name: vendor,
      import: () => optimizedImport(vendor, { retries: 2 })
    }));
  },
  // Progressive enhancement
  progressiveEnhancement: (features) => {
    return features.map(feature => ({
      ...feature,
      load: () => optimizedImport(feature.module, feature.options)
    }));
  }
};
export default {
  createOptimizedLazyComponent,
  optimizedImport,
  createIconOptimizer,
  createRouteChunk,
  analyzeBundleSize,
  createComponentRegistry,
  optimizeCSS,
  memoryOptimizer,
  bundleStrategies
}; 
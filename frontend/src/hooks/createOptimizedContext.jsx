import React, { createContext, useContext, useMemo } from 'react';
/**
 * Creates an optimized context with selector support to prevent unnecessary re-renders
 * 
 * @param {Object} options - Configuration options
 * @returns {Object} - Context provider and hooks
 */
export function createOptimizedContext(options = {}) {
  const {
    name = 'OptimizedContext',
    initialState = {},
    displayName = null,
  } = options;
  // Create base context
  const Context = createContext(initialState);
  // Set display name for dev tools
  if (displayName) {
    Context.displayName = displayName;
  }
  // Create selector tracker to optimize rerenders
  const SelectorContext = createContext(new Set());
  // Provider component with performance optimizations
  const Provider = ({ children, value }) => {
    // Create memoized selector set
    const selectorSet = useMemo(() => new Set(), []);
    // Memoize the value to avoid unnecessary rerenders
    const memoizedValue = useMemo(() => value, [
      // JSON stringify with a stable sort to ensure consistent key order
      // Only recompute when the actual data changes, not just the object reference
      JSON.stringify(value, Object.keys(value).sort())
    ]);
    return (
      <SelectorContext.Provider value={selectorSet}>
        <Context.Provider value={memoizedValue}>
          {children}
        </Context.Provider>
      </SelectorContext.Provider>
    );
  };
  /**
   * Hook to use the context with selector for performance
   * @param {Function} selector - Optional selector function
   * @returns {any} - Selected state or entire context
   */
  const useContextSelector = (selector) => {
    const value = useContext(Context);
    const selectorSet = useContext(SelectorContext);
    // If no selector is provided, return the entire context value
    if (!selector) {
      return value;
    }
    // Track this selector for potential optimization 
    if (!selectorSet.has(selector)) {
      selectorSet.add(selector);
    }
    // Apply the selector and memoize the result
    return useMemo(() => selector(value), [
      selector,
      // Only recompute when the selected value changes
      JSON.stringify(selector(value))
    ]);
  };
  // Create a simple hook for backward compatibility
  const useContextValue = () => useContext(Context);
  return {
    Provider,
    useContextSelector,
    useContextValue,
    Context, // Export raw context for advanced use cases
  };
} 
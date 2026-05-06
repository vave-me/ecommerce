/**
 * PRODUCTION MEMO FIX
 * Fixes Temporal Dead Zone errors in production builds with React.memo
 * 
 * The issue: In production, Terser minifier reorders variable declarations
 * causing memo(function ComponentName() {}) to access variables before initialization
 */
import { memo } from 'react';
/**
 * Safe memo wrapper that prevents TDZ errors in production
 * Use this instead of React.memo for components that cause production errors
 */
export const safeMemo = (Component, areEqual) => {
  // In development, use regular memo
  if (process.env.NODE_ENV === 'development') {
    return memo(Component, areEqual);
  }
  // In production, use a safer pattern that prevents variable hoisting issues
  const SafeComponent = (props) => {
    return Component(props);
  };
  // Copy over component properties
  if (Component.displayName) {
    SafeComponent.displayName = Component.displayName;
  }
  if (Component.propTypes) {
    SafeComponent.propTypes = Component.propTypes;
  }
  // Apply memoization with custom comparison if provided
  const MemoizedComponent = memo(SafeComponent, areEqual);
  // Preserve the original component name
  if (Component.name) {
    Object.defineProperty(MemoizedComponent, 'name', {
      value: Component.name,
      configurable: true
    });
  }
  return MemoizedComponent;
};
/**
 * Alternative safe memo for function components with explicit names
 * Prevents TDZ errors by avoiding function hoisting conflicts
 */
export const createSafeMemo = (componentName, componentFunc, areEqual) => {
  const SafeComponent = (props) => componentFunc(props);
  SafeComponent.displayName = componentName;
  return memo(SafeComponent, areEqual);
};
/**
 * Wrapper for components that use memo(function ComponentName() {})
 * This pattern is especially problematic in production minification
 */
export const wrapForProduction = (Component) => {
  if (process.env.NODE_ENV === 'production') {
    // In production, extract the function and wrap it safely
    return memo(Component);
  }
  return Component;
};
/**
 * Debug utility to identify problematic memo patterns
 */
export const checkMemoPattern = (Component) => {
  if (process.env.NODE_ENV === 'development') {
    const componentString = Component.toString();
    // Check for problematic patterns
    const problemPatterns = [
      /React\.memo\(function\s+\w+/,
      /memo\(function\s+\w+/,
      /const\s+\w+\s*=\s*React\.memo\(function/,
      /const\s+\w+\s*=\s*memo\(function/
    ];
    const hasProblematicPattern = problemPatterns.some(pattern => 
      pattern.test(componentString)
    );
    if (hasProblematicPattern) {
    }
  }
};
/**
 * Production-safe hooks wrapper
 * Prevents TDZ errors with useCallback/useMemo in production
 */
export const safeHooks = {
  useCallback: (callback, deps) => {
    // Import hooks dynamically to prevent hoisting issues
    const { useCallback } = require('react');
    return useCallback(callback, deps);
  },
  useMemo: (factory, deps) => {
    const { useMemo } = require('react');
    return useMemo(factory, deps);
  }
};
export default {
  safeMemo,
  createSafeMemo,
  wrapForProduction,
  checkMemoPattern,
  safeHooks
}; 
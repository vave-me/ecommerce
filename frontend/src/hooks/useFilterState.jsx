import { useState, useCallback, useEffect, useMemo } from 'react';
/**
 * Custom hook for efficient filter state management
 * - Handles local state with debouncing
 * - Provides validation
 * - Offers optimized change handlers
 * 
 * @param {Object} options - Hook configuration options
 * @returns {Object} Filter state and handlers
 */
export function useFilterState(options = {}) {
  const {
    initialFilters = {},
    debounceTime = 300,
    validateFilters = null,
    onFilterChange = null
  } = options;
  // Local state for filters
  const [filters, setFilters] = useState(initialFilters);
  // Track which filters are being debounced
  const [debouncedValues, setDebouncedValues] = useState({});
  const [errorKey, setErrorKey] = useState('');
  // Use flags to track which values have changed
  const [pendingChanges, setPendingChanges] = useState({});
  // Create debounced filter update
  useEffect(() => {
    // Skip if there are no pending changes
    if (Object.keys(pendingChanges).length === 0) return;
    // Create timers for each pending change
    const timers = Object.keys(pendingChanges).map(key => {
      return setTimeout(() => {
        setDebouncedValues(prev => ({
          ...prev,
          [key]: filters[key]
        }));
        // Clear this specific pending change
        setPendingChanges(prev => {
          const next = { ...prev };
          delete next[key];
          return next;
        });
      }, debounceTime);
    });
    // Cleanup timers on unmount or when pendingChanges changes
    return () => {
      timers.forEach(timer => clearTimeout(timer));
    };
  }, [pendingChanges, filters, debounceTime]);
  // Call onFilterChange when debounced values change (optimized)
  useEffect(() => {
    if (onFilterChange && Object.keys(debouncedValues).length > 0) {
      const mergedFilters = { ...filters };
      onFilterChange(mergedFilters, Object.keys(debouncedValues));
    }
  }, [debouncedValues, onFilterChange]);
  // Handle individual filter changes
  const handleChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
    // Mark this key as having a pending change
    setPendingChanges(prev => ({ ...prev, [key]: true }));
    // Clear any previous errors
    setErrorKey('');
  }, []);
  // Handle toggle filters (boolean values)
  const handleToggle = useCallback((key) => {
    setFilters(prev => ({ ...prev, [key]: !prev[key] }));
    // Mark this key as having a pending change
    setPendingChanges(prev => ({ ...prev, [key]: true }));
    // Clear any previous errors
    setErrorKey('');
  }, []);
  // Handle array values like tags
  const handleArrayChange = useCallback((key, valueStr, delimiter = ',') => {
    const arrayValue = valueStr
      .split(delimiter)
      .map(item => item.trim())
      .filter(Boolean);
    setFilters(prev => ({ ...prev, [key]: arrayValue }));
    // Mark this key as having a pending change
    setPendingChanges(prev => ({ ...prev, [key]: true }));
    // Clear any previous errors
    setErrorKey('');
  }, []);
  // Reset filters to initial state
  const resetFilters = useCallback(() => {
    setFilters(initialFilters);
    setDebouncedValues({});
    setPendingChanges({});
    setErrorKey('');
  }, [initialFilters]);
  // Update entire filter state at once
  const updateFilters = useCallback((newFilters) => {
    setFilters(newFilters);
    setDebouncedValues({});
    setPendingChanges({});
    setErrorKey('');
  }, []);
  // Validate filters and return error key if invalid
  const validate = useCallback(() => {
    if (!validateFilters) return '';
    const errorKey = validateFilters(filters);
    setErrorKey(errorKey);
    return errorKey;
  }, [filters, validateFilters]);
  // Check if filters have changed from initial state
  const hasChanges = useMemo(() => {
    // Quick length check first
    if (Object.keys(filters).length !== Object.keys(initialFilters).length) {
      return true;
    }
    // Deep comparison of values
    for (const key in filters) {
      const initialValue = initialFilters[key];
      const currentValue = filters[key];
      // Handle array comparison
      if (Array.isArray(currentValue) && Array.isArray(initialValue)) {
        if (currentValue.length !== initialValue.length) {
          return true;
        }
        for (let i = 0; i < currentValue.length; i++) {
          if (currentValue[i] !== initialValue[i]) {
            return true;
          }
        }
      } 
      // Handle object comparison
      else if (
        typeof currentValue === 'object' && 
        currentValue !== null &&
        typeof initialValue === 'object' && 
        initialValue !== null
      ) {
        if (JSON.stringify(currentValue) !== JSON.stringify(initialValue)) {
          return true;
        }
      }
      // Simple value comparison
      else if (currentValue !== initialValue) {
        return true;
      }
    }
    return false;
  }, [filters, initialFilters]);
  return {
    filters,
    errorKey,
    handleChange,
    handleToggle,
    handleArrayChange,
    resetFilters,
    updateFilters,
    validate,
    hasChanges,
    hasPendingChanges: Object.keys(pendingChanges).length > 0
  };
} 
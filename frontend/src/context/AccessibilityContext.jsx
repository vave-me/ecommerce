"use client";

import React, { createContext, useContext, useState, useCallback, useRef } from 'react';

// Context for managing screen reader announcements
const AccessibilityContext = createContext({
  announce: () => {},
  announcePolite: () => {},
  announceAssertive: () => {},
  clearAnnouncement: () => {}
});

export const useAccessibility = () => {
  const context = useContext(AccessibilityContext);
  if (!context) {
    throw new Error('useAccessibility must be used within an AccessibilityProvider');
  }
  return context;
};

export const AccessibilityProvider = ({ children }) => {
  const [announcement, setAnnouncement] = useState('');
  const [politeness, setPoliteness] = useState('polite');
  const announcementTimeoutRef = useRef(null);

  // Clear any existing timeout
  const clearTimeout = useCallback(() => {
    if (announcementTimeoutRef.current) {
      window.clearTimeout(announcementTimeoutRef.current);
      announcementTimeoutRef.current = null;
    }
  }, []);

  // Generic announce function
  const announce = useCallback((message, politenessLevel = 'polite', clearAfter = 5000) => {
    clearTimeout();
    setAnnouncement(message);
    setPoliteness(politenessLevel);

    if (clearAfter > 0) {
      announcementTimeoutRef.current = window.setTimeout(() => {
        setAnnouncement('');
      }, clearAfter);
    }
  }, [clearTimeout]);

  // Convenience methods
  const announcePolite = useCallback((message, clearAfter = 5000) => {
    announce(message, 'polite', clearAfter);
  }, [announce]);

  const announceAssertive = useCallback((message, clearAfter = 5000) => {
    announce(message, 'assertive', clearAfter);
  }, [announce]);

  const clearAnnouncement = useCallback(() => {
    clearTimeout();
    setAnnouncement('');
  }, [clearTimeout]);

  const value = {
    announce,
    announcePolite,
    announceAssertive,
    clearAnnouncement
  };

  return (
    <AccessibilityContext.Provider value={value}>
      {children}
      {/* Live regions for screen reader announcements */}
      <div className="sr-only" aria-live={politeness} aria-atomic="true">
        {announcement}
      </div>
      {/* Separate regions for different politeness levels to avoid conflicts */}
      <div className="sr-only" aria-live="polite" aria-atomic="true" id="polite-announcements" />
      <div className="sr-only" aria-live="assertive" aria-atomic="true" id="assertive-announcements" />
    </AccessibilityContext.Provider>
  );
};

// Hook for common announcement patterns
export const useAnnouncements = () => {
  const { announcePolite, announceAssertive } = useAccessibility();

  const announceActionResult = useCallback((action, itemName, success = true) => {
    const message = success 
      ? `${action} ${itemName} successful`
      : `Failed to ${action} ${itemName}`;
    announcePolite(message);
  }, [announcePolite]);

  const announceLoading = useCallback((itemName, isLoading = true) => {
    const message = isLoading
      ? `Loading ${itemName}`
      : `${itemName} loaded`;
    announcePolite(message);
  }, [announcePolite]);

  const announceError = useCallback((errorMessage) => {
    announceAssertive(`Error: ${errorMessage}`);
  }, [announceAssertive]);

  const announceNavigation = useCallback((destination) => {
    announcePolite(`Navigating to ${destination}`);
  }, [announcePolite]);

  const announceUpdate = useCallback((itemName, updateType) => {
    announcePolite(`${itemName} ${updateType}`);
  }, [announcePolite]);

  return {
    announceActionResult,
    announceLoading,
    announceError,
    announceNavigation,
    announceUpdate
  };
};
// Accessibility utility functions

/**
 * Generate unique IDs for ARIA relationships
 */
export const generateId = (prefix = 'element') => {
  return `${prefix}-${Math.random().toString(36).substr(2, 9)}`;
};

/**
 * Format number for screen reader announcement
 */
export const formatNumberForScreenReader = (number) => {
  if (number >= 1000000) {
    return `${(number / 1000000).toFixed(1)} million`;
  }
  if (number >= 1000) {
    return `${(number / 1000).toFixed(1)} thousand`;
  }
  return number.toString();
};

/**
 * Get ARIA label for interactive count
 */
export const getCountAriaLabel = (action, count, plural) => {
  const formattedCount = formatNumberForScreenReader(count);
  return count === 0 
    ? `${action} this ${plural}`
    : `${action} this ${plural}, currently ${formattedCount} ${count === 1 ? plural.slice(0, -1) : plural}`;
};

/**
 * Check if user prefers reduced motion
 */
export const prefersReducedMotion = () => {
  if (typeof window === 'undefined') return false;
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
};

/**
 * Check if user is using high contrast mode
 */
export const prefersHighContrast = () => {
  if (typeof window === 'undefined') return false;
  return window.matchMedia('(prefers-contrast: high)').matches;
};

/**
 * Trap focus within an element
 */
export const trapFocus = (element) => {
  const focusableElements = element.querySelectorAll(
    'a[href], button, textarea, input[type="text"], input[type="radio"], input[type="checkbox"], select, [tabindex]:not([tabindex="-1"])'
  );
  
  const firstFocusableElement = focusableElements[0];
  const lastFocusableElement = focusableElements[focusableElements.length - 1];

  const handleTabKey = (e) => {
    if (e.key !== 'Tab') return;

    if (e.shiftKey) {
      if (document.activeElement === firstFocusableElement) {
        lastFocusableElement.focus();
        e.preventDefault();
      }
    } else {
      if (document.activeElement === lastFocusableElement) {
        firstFocusableElement.focus();
        e.preventDefault();
      }
    }
  };

  element.addEventListener('keydown', handleTabKey);
  firstFocusableElement?.focus();

  return () => {
    element.removeEventListener('keydown', handleTabKey);
  };
};

/**
 * Get descriptive text for time ago
 */
export const getTimeAgoDescription = (date) => {
  const now = new Date();
  const then = new Date(date);
  const diff = now - then;
  
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const weeks = Math.floor(days / 7);
  const months = Math.floor(days / 30);
  const years = Math.floor(days / 365);

  if (years > 0) return `${years} ${years === 1 ? 'year' : 'years'} ago`;
  if (months > 0) return `${months} ${months === 1 ? 'month' : 'months'} ago`;
  if (weeks > 0) return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`;
  if (days > 0) return `${days} ${days === 1 ? 'day' : 'days'} ago`;
  if (hours > 0) return `${hours} ${hours === 1 ? 'hour' : 'hours'} ago`;
  if (minutes > 0) return `${minutes} ${minutes === 1 ? 'minute' : 'minutes'} ago`;
  return 'Just now';
};

/**
 * Debounce function for reducing announcement frequency
 */
export const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

/**
 * Get contrast ratio between two colors
 */
export const getContrastRatio = (color1, color2) => {
  const getLuminance = (color) => {
    const rgb = color.match(/\d+/g);
    if (!rgb || rgb.length < 3) return 0;
    
    const [r, g, b] = rgb.map(val => {
      const sRGB = parseInt(val) / 255;
      return sRGB <= 0.03928
        ? sRGB / 12.92
        : Math.pow((sRGB + 0.055) / 1.055, 2.4);
    });
    
    return 0.2126 * r + 0.7152 * g + 0.0722 * b;
  };

  const l1 = getLuminance(color1);
  const l2 = getLuminance(color2);
  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);

  return (lighter + 0.05) / (darker + 0.05);
};

/**
 * Check if contrast meets WCAG AA standards
 */
export const meetsContrastStandard = (color1, color2, isLargeText = false) => {
  const ratio = getContrastRatio(color1, color2);
  return isLargeText ? ratio >= 3 : ratio >= 4.5;
};

/**
 * Create keyboard navigation handler
 */
export const createKeyboardNavigationHandler = (options = {}) => {
  const {
    onEnter,
    onSpace,
    onArrowUp,
    onArrowDown,
    onArrowLeft,
    onArrowRight,
    onEscape,
    onTab,
    preventDefault = true
  } = options;

  return (event) => {
    const { key, shiftKey } = event;
    let handled = false;

    switch (key) {
      case 'Enter':
        if (onEnter) {
          onEnter(event);
          handled = true;
        }
        break;
      case ' ':
      case 'Space':
        if (onSpace) {
          onSpace(event);
          handled = true;
        }
        break;
      case 'ArrowUp':
        if (onArrowUp) {
          onArrowUp(event);
          handled = true;
        }
        break;
      case 'ArrowDown':
        if (onArrowDown) {
          onArrowDown(event);
          handled = true;
        }
        break;
      case 'ArrowLeft':
        if (onArrowLeft) {
          onArrowLeft(event);
          handled = true;
        }
        break;
      case 'ArrowRight':
        if (onArrowRight) {
          onArrowRight(event);
          handled = true;
        }
        break;
      case 'Escape':
        if (onEscape) {
          onEscape(event);
          handled = true;
        }
        break;
      case 'Tab':
        if (onTab) {
          onTab(event, shiftKey);
          handled = true;
        }
        break;
    }

    if (handled && preventDefault) {
      event.preventDefault();
    }
  };
};
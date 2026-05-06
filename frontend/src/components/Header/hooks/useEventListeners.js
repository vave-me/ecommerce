import { useEffect, useRef, useCallback } from 'react';
/**
 * Optimized Event Listener Hooks for Header Components
 * Prevents memory leaks and improves performance
 */
// Debounced event listener hook
export const useDebouncedEventListener = (eventType, handler, delay = 100, options = {}) => {
    const timeoutRef = useRef(null);
    const handlerRef = useRef(handler);
    // Update handler ref when handler changes
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);
    const debouncedHandler = useCallback((event) => {
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
        }
        timeoutRef.current = setTimeout(() => {
            handlerRef.current(event);
        }, delay);
    }, [delay]);
    useEffect(() => {
        const element = options.target || window;
        element.addEventListener(eventType, debouncedHandler, options);
        return () => {
            element.removeEventListener(eventType, debouncedHandler, options);
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
        };
    }, [eventType, debouncedHandler, options]);
};
// Throttled event listener hook
export const useThrottledEventListener = (eventType, handler, delay = 100, options = {}) => {
    const lastCallRef = useRef(0);
    const handlerRef = useRef(handler);
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);
    const throttledHandler = useCallback((event) => {
        const now = Date.now();
        if (now - lastCallRef.current >= delay) {
            handlerRef.current(event);
            lastCallRef.current = now;
        }
    }, [delay]);
    useEffect(() => {
        const element = options.target || window;
        element.addEventListener(eventType, throttledHandler, options);
        return () => {
            element.removeEventListener(eventType, throttledHandler, options);
        };
    }, [eventType, throttledHandler, options]);
};
// Optimized scroll listener with passive option
export const useScrollListener = (handler, options = { passive: true }) => {
    const handlerRef = useRef(handler);
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);
    useEffect(() => {
        const element = options.target || window;
        const scrollHandler = (event) => handlerRef.current(event);
        element.addEventListener('scroll', scrollHandler, options);
        return () => {
            element.removeEventListener('scroll', scrollHandler, options);
        };
    }, [options]);
};
// Optimized resize listener with debouncing
export const useResizeListener = (handler, delay = 100) => {
    return useDebouncedEventListener('resize', handler, delay, { passive: true });
};
// Optimized click outside listener
export const useClickOutside = (ref, handler, options = {}) => {
    const handlerRef = useRef(handler);
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);
    useEffect(() => {
        const element = ref.current;
        if (!element) return;
        const clickHandler = (event) => {
            if (!element.contains(event.target)) {
                handlerRef.current(event);
            }
        };
        const touchHandler = (event) => {
            if (!element.contains(event.target)) {
                handlerRef.current(event);
            }
        };
        // Use capture phase for better performance
        document.addEventListener('mousedown', clickHandler, { capture: true });
        document.addEventListener('touchend', touchHandler, { capture: true });
        return () => {
            document.removeEventListener('mousedown', clickHandler, { capture: true });
            document.removeEventListener('touchend', touchHandler, { capture: true });
        };
    }, [ref]);
};
// Optimized keyboard listener
export const useKeyboardListener = (key, handler, options = {}) => {
    const handlerRef = useRef(handler);
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);
    useEffect(() => {
        const element = options.target || document;
        const keyHandler = (event) => {
            if (event.key === key) {
                handlerRef.current(event);
            }
        };
        element.addEventListener('keydown', keyHandler, options);
        return () => {
            element.removeEventListener('keydown', keyHandler, options);
        };
    }, [key, options]);
};
// Focus trap hook for modals and dropdowns
export const useFocusTrap = (ref, isActive = true) => {
    const focusableElementsRef = useRef([]);
    const getFocusableElements = useCallback(() => {
        const element = ref.current;
        if (!element) return [];
        return Array.from(
            element.querySelectorAll(
                'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
            )
        ).filter(el => !el.disabled && el.offsetParent !== null);
    }, [ref]);
    useEffect(() => {
        if (!isActive) return;
        const element = ref.current;
        if (!element) return;
        focusableElementsRef.current = getFocusableElements();
        const focusableElements = focusableElementsRef.current;
        if (focusableElements.length === 0) return;
        const handleKeyDown = (event) => {
            if (event.key !== 'Tab') return;
            const firstElement = focusableElements[0];
            const lastElement = focusableElements[focusableElements.length - 1];
            if (event.shiftKey) {
                if (document.activeElement === firstElement) {
                    event.preventDefault();
                    lastElement.focus();
                }
            } else {
                if (document.activeElement === lastElement) {
                    event.preventDefault();
                    firstElement.focus();
                }
            }
        };
        element.addEventListener('keydown', handleKeyDown);
        return () => {
            element.removeEventListener('keydown', handleKeyDown);
        };
    }, [isActive, ref, getFocusableElements]);
}; 
// CreatePostModal/hooks/useFocusTrap.js
import { useRef, useEffect } from 'react';
export function useFocusTrap(isActive) {
    const ref = useRef(null);
    useEffect(() => {
        if (!isActive || !ref.current) return;
        const focusableElements =
            'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])';
        const modal = ref.current;
        const firstFocusableElement = modal.querySelectorAll(focusableElements)[0];
        const focusableContent = modal.querySelectorAll(focusableElements);
        const lastFocusableElement = focusableContent[focusableContent.length - 1];
        // Save the element that had focus before opening the modal
        const previouslyFocused = document.activeElement;
        // Focus first element
        firstFocusableElement?.focus();
        const handleKeyDown = (e) => {
            const isTabPressed = e.key === "Tab" || e.keyCode === 9;
            if (!isTabPressed) return;
            // Shift + Tab
            if (e.shiftKey) {
                if (document.activeElement === firstFocusableElement) {
                    lastFocusableElement.focus();
                    e.preventDefault();
                }
            }
            // Tab without shift
            else {
                if (document.activeElement === lastFocusableElement) {
                    firstFocusableElement.focus();
                    e.preventDefault();
                }
            }
        };
        modal.addEventListener("keydown", handleKeyDown);
        return () => {
            modal.removeEventListener("keydown", handleKeyDown);
            // Restore focus when the modal closes
            if (previouslyFocused) {
                previouslyFocused.focus();
            }
        };
    }, [isActive]);
    return ref;
}
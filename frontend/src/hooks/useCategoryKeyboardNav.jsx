// src/hooks/useCategoryKeyboardNav.js
import { useCallback } from 'react';
/**
 * A hook to manage arrow-key navigation in a flattened category list.
 *
 * @param {Array} flattened - Array of categories in visual order (already filtered/expanded)
 * @param {Object} expandedMap - Object of { [id]: boolean } controlling expanded/collapsed state
 * @param {Function} toggleExpand - Function to toggle expand/collapse for a category
 * @param {Function} onSelect - Function called when user presses Enter/Space on a category
 * @param {Function} setFocusedId - Optional function to set the currently focused category ID
 * @returns {Object} Object containing the handleKeyDown function
 */
export function useCategoryKeyboardNav(
    flattened,
    expandedMap,
    toggleExpand,
    onSelect,
    setFocusedId
) {
    const handleKeyDown = useCallback(
        (e, currentCat) => {
            // Only handle specific navigation keys
            const { key } = e;
            if (!['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' ', 'Enter'].includes(key)) {
                return;
            }
            e.preventDefault();
            e.stopPropagation();
            const currentIndex = flattened.findIndex((c) => c.id === currentCat.id);
            if (currentIndex < 0) return;
            if (key === 'ArrowDown') {
                const nextIndex = Math.min(currentIndex + 1, flattened.length - 1);
                if (setFocusedId && flattened[nextIndex]) {
                    setFocusedId(flattened[nextIndex].id);
                }
            } else if (key === 'ArrowUp') {
                const prevIndex = Math.max(currentIndex - 1, 0);
                if (setFocusedId && flattened[prevIndex]) {
                    setFocusedId(flattened[prevIndex].id);
                }
            } else if (key === 'ArrowRight') {
                // If category has subcategories and is collapsed, expand it
                if (currentCat.subcategories?.length && !expandedMap[currentCat.id]) {
                    toggleExpand(currentCat.id);
                }
                // If already expanded and there are children, move focus to first child
                else if (expandedMap[currentCat.id] && currentIndex + 1 < flattened.length) {
                    const nextCat = flattened[currentIndex + 1];
                    if (nextCat && setFocusedId) {
                        setFocusedId(nextCat.id);
                    }
                }
            } else if (key === 'ArrowLeft') {
                // If category is expanded, collapse it
                if (expandedMap[currentCat.id]) {
                    toggleExpand(currentCat.id);
                }
                // If already collapsed and has parent, move focus to parent
                else if (currentCat.parentId) {
                    const parentIndex = flattened.findIndex(c => c.id === currentCat.parentId);
                    if (parentIndex >= 0 && setFocusedId) {
                        setFocusedId(flattened[parentIndex].id);
                    }
                }
            } else if (key === ' ' || key === 'Enter') {
                onSelect?.(currentCat);
            }
        },
        [flattened, expandedMap, toggleExpand, onSelect, setFocusedId]
    );
    return { handleKeyDown };
}
import { useEffect, useRef } from 'react';

const useMasonryLayout = (dependencies = []) => {
    const containerRef = useRef(null);

    useEffect(() => {
        const updateLayout = () => {
            if (!containerRef.current) return;

            const items = containerRef.current.querySelectorAll('[class*="feedItemContainer"]');
            const rowHeight = 1; // Match the grid-auto-rows value

            items.forEach((item) => {
                const card = item.firstChild;
                if (!card) return;

                // Reset to natural height first
                item.style.gridRowEnd = 'auto';
                
                // Force layout recalculation
                const height = card.getBoundingClientRect().height;
                
                // Calculate how many rows this item should span
                // Add extra rows for gap
                const rowSpan = Math.ceil((height + 20) / rowHeight);
                
                // Apply the calculated span
                item.style.gridRowEnd = `span ${rowSpan}`;
            });
        };

        // Run on mount and when dependencies change
        updateLayout();

        // Run on window resize
        window.addEventListener('resize', updateLayout);
        
        // Run after images load
        const images = containerRef.current?.querySelectorAll('img');
        images?.forEach(img => {
            if (img.complete) {
                updateLayout();
            } else {
                img.addEventListener('load', updateLayout);
            }
        });

        return () => {
            window.removeEventListener('resize', updateLayout);
            images?.forEach(img => {
                img.removeEventListener('load', updateLayout);
            });
        };
    }, dependencies);

    return containerRef;
};

export default useMasonryLayout;
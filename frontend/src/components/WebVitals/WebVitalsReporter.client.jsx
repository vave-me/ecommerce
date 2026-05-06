'use client';
import {useEffect} from 'react';
/**
 * WebVitalsReporter Client Component
 * Handles Web Vitals initialization on the client side
 */
export default function WebVitalsReporter() {
    useEffect(() => {
        // Only run on client side
        if (typeof window === 'undefined') return;
        const initWebVitals = async () => {
            try {
                // Dynamic import to avoid SSR issues and reduce bundle size
                const webVitalsModule = await import('../../utils/webVitalsReporter');
                if (webVitalsModule.initWebVitals) {
                    await webVitalsModule.initWebVitals();
                }
            } catch (error) {
                // Graceful fallback - don't break the app if Web Vitals fails
                if (process.env.NODE_ENV === 'development') {
                    }
            }
        };
        // Initialize after the page has loaded
        if (document.readyState === 'complete') {
            initWebVitals();
        } else {
            window.addEventListener('load', initWebVitals);
            return () => window.removeEventListener('load', initWebVitals);
        }
    }, []);
    // This component doesn't render anything
    return null;
}

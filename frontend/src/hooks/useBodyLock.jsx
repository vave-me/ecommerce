"use client";
import { useEffect } from "react";
/**
 * Lock the body scroll if `isLocked` is true.
 * Example usage: useBodyLock(isMobileMenuOpen);
 */
export function useBodyLock(isLocked) {
    useEffect(() => {
        if (isLocked) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "";
        }
        return () => {
            // Cleanup just in case
            document.body.style.overflow = "";
        };
    }, [isLocked]);
}

"use client";
import { useState, useEffect } from "react";
/**
 * Returns a boolean indicating if window.scrollY is beyond `offset` px.
 */
export function useHeaderScroll({ offset = 20 }) {
    const [isScrolled, setIsScrolled] = useState(false);
    useEffect(() => {
        function handleScroll() {
            const scrolled = window.scrollY > offset;
            setIsScrolled(scrolled);
        }
        window.addEventListener("scroll", handleScroll, { passive: true });
        // Check immediately on mount
        handleScroll();
        return () => {
            window.removeEventListener("scroll", handleScroll);
        };
    }, [offset]);
    return isScrolled;
}

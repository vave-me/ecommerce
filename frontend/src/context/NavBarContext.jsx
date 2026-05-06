// File: src/context/NavBarContext.jsx
"use client";
import React, {createContext, useState, useEffect, useCallback} from 'react';

export const NavBarContext = createContext(null);

export const NavBarProvider = ({children}) => {
    const [isMobile, setIsMobile] = useState(false);
    const [showNavbars, setShowNavbars] = useState(false);

    const handleResize = useCallback(() => {
        if (typeof window !== 'undefined') {
            if (window.innerWidth <= 768) {
                setIsMobile(true);
                setShowNavbars(true);
            } else {
                setIsMobile(false);
                setShowNavbars(false);
            }
        }
    }, []);

    useEffect(() => {
        handleResize();
        if (typeof window !== 'undefined') {
            window.addEventListener('resize', handleResize);
            return () => window.removeEventListener('resize', handleResize);
        }
    }, [handleResize]);

    const toggleNavbars = () => {
        setShowNavbars((prev) => !prev);
    };

    return (
        <NavBarContext.Provider value={{showNavbars, toggleNavbars, isMobile}}>
            {children}
        </NavBarContext.Provider>
    );
};

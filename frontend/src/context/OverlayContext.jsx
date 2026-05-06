"use client";

import React, { createContext, useContext, useState, useCallback } from 'react';
import PropTypes from 'prop-types';

// Create the context
const OverlayContext = createContext({
    isUserMenuOpen: false,
    isHorizontalFiltersOpen: false,
    setUserMenuOpen: () => {},
    setHorizontalFiltersOpen: () => {},
    hasActiveOverlay: false,
});

// Custom hook to use the overlay context
export const useOverlay = () => {
    const context = useContext(OverlayContext);
    if (!context) {
        throw new Error('useOverlay must be used within an OverlayProvider');
    }
    return context;
};

// Provider component
export const OverlayProvider = ({ children }) => {
    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
    const [isHorizontalFiltersOpen, setIsHorizontalFiltersOpen] = useState(false);

    const setUserMenuOpen = useCallback((isOpen) => {
        setIsUserMenuOpen(isOpen);
    }, []);

    const setHorizontalFiltersOpen = useCallback((isOpen) => {
        setIsHorizontalFiltersOpen(isOpen);
    }, []);
    
    // Maintain backward compatibility
    const setSelectTopicOpen = setHorizontalFiltersOpen;

    // Calculate if any overlay is active
    const hasActiveOverlay = isUserMenuOpen || isHorizontalFiltersOpen;

    const value = {
        isUserMenuOpen,
        isHorizontalFiltersOpen,
        isSelectTopicOpen: isHorizontalFiltersOpen, // backward compatibility
        setUserMenuOpen,
        setHorizontalFiltersOpen,
        setSelectTopicOpen, // backward compatibility
        hasActiveOverlay,
    };

    return (
        <OverlayContext.Provider value={value}>
            {children}
        </OverlayContext.Provider>
    );
};

OverlayProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export default OverlayContext; 
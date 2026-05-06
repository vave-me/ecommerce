"use client";
import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import PropTypes from 'prop-types';
const ShareContext = createContext();
/**
 * ShareProvider - Context for managing global share settings and analytics
 * Provides centralized share configuration and tracking across the application
 */
export const ShareProvider = ({ children, defaultConfig = {} }) => {
    // Global share configuration
    const [config, setConfig] = useState({
        // Default platforms available across the app
        defaultPlatforms: ['native', 'copy', 'facebook', 'twitter', 'linkedin', 'email'],
        // Global analytics settings
        trackShares: true,
        trackingEndpoint: '/api/shares',
        // UI defaults
        defaultVariant: 'button',
        defaultSize: 'medium',
        // Feature flags
        requireAuthForSharing: false,
        enableShareCounts: true,
        enableAnalytics: true,
        // Custom share URLs (for branded sharing)
        customShareUrls: {},
        // Override with provided config
        ...defaultConfig
    });
    // Share analytics state
    const [analytics, setAnalytics] = useState({
        totalShares: 0,
        sharesByPlatform: {},
        sharesByContent: {},
        recentShares: []
    });
    // Loading and error states
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    // Update global configuration
    const updateConfig = useCallback((newConfig) => {
        setConfig(prev => ({ ...prev, ...newConfig }));
    }, []);
    // Record share event globally
    const recordGlobalShare = useCallback(async (shareData) => {
        if (!config.trackShares) return;
        try {
            setLoading(true);
            // Update local analytics immediately (optimistic update)
            setAnalytics(prev => ({
                ...prev,
                totalShares: prev.totalShares + 1,
                sharesByPlatform: {
                    ...prev.sharesByPlatform,
                    [shareData.platform]: (prev.sharesByPlatform[shareData.platform] || 0) + 1
                },
                sharesByContent: {
                    ...prev.sharesByContent,
                    [shareData.contentId]: (prev.sharesByContent[shareData.contentId] || 0) + 1
                },
                recentShares: [
                    {
                        ...shareData,
                        timestamp: new Date().toISOString(),
                        id: Date.now() // Temporary ID
                    },
                    ...prev.recentShares.slice(0, 49) // Keep last 50 shares
                ]
            }));
            // Persist to localStorage
            setTimeout(() => {
                setAnalytics(current => {
                    try {
                        localStorage.setItem('shareAnalytics', JSON.stringify(current));
                    } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
                    return current;
                });
            }, 100);
            // Frontend-only analytics tracking
            if (config.enableAnalytics) {
                // Could integrate with Google Analytics here if available
                if (typeof window !== 'undefined' && window.gtag) {
                    window.gtag('event', 'share', {
                        method: shareData.platform,
                        content_type: shareData.contentType,
                        content_id: shareData.contentId
                    });
                }
            }
        } catch (err) {
            setError(err.message);
            // Revert optimistic update on error
            setAnalytics(prev => ({
                ...prev,
                totalShares: Math.max(0, prev.totalShares - 1),
                sharesByPlatform: {
                    ...prev.sharesByPlatform,
                    [shareData.platform]: Math.max(0, (prev.sharesByPlatform[shareData.platform] || 1) - 1)
                },
                sharesByContent: {
                    ...prev.sharesByContent,
                    [shareData.contentId]: Math.max(0, (prev.sharesByContent[shareData.contentId] || 1) - 1)
                },
                recentShares: prev.recentShares.slice(1) // Remove the failed share
            }));
        } finally {
            setLoading(false);
        }
    }, [config]);
    // Get share count for specific content
    const getShareCount = useCallback((contentId) => {
        return analytics.sharesByContent[contentId] || 0;
    }, [analytics.sharesByContent]);
    // Get platform analytics
    const getPlatformAnalytics = useCallback(() => {
        return analytics.sharesByPlatform;
    }, [analytics.sharesByPlatform]);
    // Get recent shares
    const getRecentShares = useCallback((limit = 10) => {
        return analytics.recentShares.slice(0, limit);
    }, [analytics.recentShares]);
    // Clear error
    const clearError = useCallback(() => {
        setError(null);
    }, []);
    // Load initial analytics data (frontend-only)
    useEffect(() => {
        if (config.enableAnalytics) {
            // Frontend-only - initialize with mock data or localStorage
            const loadAnalytics = () => {
                try {
                    // Try to load from localStorage
                    const stored = localStorage.getItem('shareAnalytics');
                    if (stored) {
                        const data = JSON.parse(stored);
                        setAnalytics(prev => ({ ...prev, ...data }));
                    }
                } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
            };
            loadAnalytics();
        }
    }, [config.enableAnalytics]);
    // Context value
    const contextValue = {
        // Configuration
        config,
        updateConfig,
        // Analytics
        analytics,
        recordGlobalShare,
        getShareCount,
        getPlatformAnalytics,
        getRecentShares,
        // State
        loading,
        error,
        clearError,
        // Utilities
        isFeatureEnabled: (feature) => config[feature] === true,
        getPlatformConfig: (platform) => config.customShareUrls[platform] || null,
    };
    return (
        <ShareContext.Provider value={contextValue}>
            {children}
        </ShareContext.Provider>
    );
};
ShareProvider.propTypes = {
    children: PropTypes.node.isRequired,
    defaultConfig: PropTypes.object,
};
/**
 * Hook to use share context
 */
export const useShareContext = () => {
    const context = useContext(ShareContext);
    if (!context) {
        throw new Error('useShareContext must be used within a ShareProvider');
    }
    return context;
};
/**
 * HOC to wrap components with share functionality
 */
export const withShare = (Component) => {
    const WrappedComponent = (props) => {
        const shareContext = useShareContext();
        return <Component {...props} shareContext={shareContext} />;
    };
    WrappedComponent.displayName = `withShare(${Component.displayName || Component.name})`;
    return WrappedComponent;
};
export default ShareProvider; 
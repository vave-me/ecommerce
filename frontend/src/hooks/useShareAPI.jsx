"use client";
import { useState, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { toast } from 'react-toastify';
/**
 * Enhanced custom hook for handling share API calls and analytics
 * Provides consistent sharing functionality with proper error handling and fallbacks
 */
export const useShareAPI = () => {
    const { user } = useAuth();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    // Record share event locally (frontend only)
    const recordShare = useCallback(async (shareData) => {
        try {
            setLoading(true);
            setError(null);
            // Frontend-only tracking
            // Simulate successful response
            return {
                success: true,
                id: Date.now(),
                timestamp: new Date().toISOString()
            };
        } catch (err) {
            setError(err.message);
            throw err;
        } finally {
            setLoading(false);
        }
    }, [user]);
    // Update share count locally (frontend only)
    const updateShareCount = useCallback(async (contentId, contentType) => {
        try {
            // Frontend-only - return mock incremented count
            return Math.floor(Math.random() * 100) + 1;
        } catch (err) {
            return null;
        }
    }, []);
    // Get share analytics locally (frontend only)
    const getShareAnalytics = useCallback(async (contentId, contentType) => {
        try {
            // Frontend-only - return mock analytics
            return {
                totalShares: Math.floor(Math.random() * 50),
                platforms: {
                    facebook: Math.floor(Math.random() * 20),
                    twitter: Math.floor(Math.random() * 15),
                    copy: Math.floor(Math.random() * 10),
                    native: Math.floor(Math.random() * 5)
                },
                lastShared: new Date().toISOString()
            };
        } catch (err) {
            return null;
        }
    }, [user]);
    // Enhanced browser compatibility checks
    const checkBrowserCapabilities = useCallback(() => {
        const capabilities = {
            hasNavigator: typeof navigator !== 'undefined',
            hasWebShare: typeof navigator !== 'undefined' && !!navigator.share,
            hasClipboard: typeof navigator !== 'undefined' && !!navigator.clipboard && !!navigator.clipboard.writeText,
            isSecureContext: typeof window !== 'undefined' && (window.isSecureContext || window.location.protocol === 'https:'),
            isMobile: typeof window !== 'undefined' && /Android|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
        };
        return capabilities;
    }, []);
    // Enhanced clipboard copy with fallbacks
    const copyToClipboard = useCallback(async (text) => {
        const capabilities = checkBrowserCapabilities();
        // Method 1: Modern Clipboard API (requires HTTPS)
        if (capabilities.hasClipboard && capabilities.isSecureContext) {
            try {
                await navigator.clipboard.writeText(text);
                toast.success('Link copied to clipboard!');
                return true;
            } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
        }
        // Method 2: Legacy execCommand fallback
        try {
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            const successful = document.execCommand('copy');
            document.body.removeChild(textArea);
            if (successful) {
                toast.success('Link copied to clipboard!');
                return true;
            }
        } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
        // Method 3: Manual selection (last resort)
        try {
            const range = document.createRange();
            const selection = window.getSelection();
            const tempDiv = document.createElement('div');
            tempDiv.textContent = text;
            tempDiv.style.position = 'absolute';
            tempDiv.style.left = '-999999px';
            document.body.appendChild(tempDiv);
            range.selectNodeContents(tempDiv);
            selection.removeAllRanges();
            selection.addRange(range);
            toast.info('Please copy the selected text manually (Ctrl+C)');
            setTimeout(() => {
                document.body.removeChild(tempDiv);
                selection.removeAllRanges();
            }, 3000);
            return true;
        } catch (err) {
            toast.error('Unable to copy link. Please copy manually: ' + text);
            return false;
        }
    }, [checkBrowserCapabilities]);
    // Enhanced native share with better error handling
    const nativeShare = useCallback(async (shareData) => {
        const capabilities = checkBrowserCapabilities();
        if (!capabilities.hasWebShare) {
            throw new Error('Web Share API not supported');
        }
        try {
            await navigator.share({
                title: shareData.title || document.title,
                text: shareData.description || '',
                url: shareData.url || window.location.href,
            });
            // Record successful share
            await recordShare({
                ...shareData,
                platform: 'native'
            });
            toast.success('Content shared successfully!');
            return { success: true };
        } catch (err) {
            if (err.name === 'AbortError') {
                // User cancelled, don't show error
                return { success: false, cancelled: true };
            }
            toast.error('Share failed. Trying alternative method...');
            throw err;
        }
    }, [checkBrowserCapabilities, recordShare]);
    // Enhanced social media sharing
    const socialShare = useCallback((platform, shareData) => {
        const encodedUrl = encodeURIComponent(shareData.url || window.location.href);
        const encodedTitle = encodeURIComponent(shareData.title || document.title);
        const encodedDescription = encodeURIComponent(shareData.description || '');
        const facebookAppId = typeof process !== 'undefined' ? (process.env.NEXT_PUBLIC_FACEBOOK_APP_ID || '') : '';
        const shareUrls = {
            facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}`,
            twitter: `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`,
            x: `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`,
            linkedin: `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}`,
            pinterest: `https://pinterest.com/pin/create/button/?url=${encodedUrl}&description=${encodedTitle}`,
            reddit: `https://reddit.com/submit?url=${encodedUrl}&title=${encodedTitle}`,
            telegram: `https://t.me/share/url?url=${encodedUrl}&text=${encodedTitle}`,
            whatsapp: `https://wa.me/?text=${encodedTitle}%20${encodedUrl}`,
            messenger: facebookAppId ? `https://www.facebook.com/dialog/send?app_id=${facebookAppId}&link=${encodedUrl}&redirect_uri=${encodedUrl}` : null,
            email: `mailto:?subject=${encodedTitle}&body=${encodedDescription}%0A%0A${encodedUrl}`,
        };
        const shareUrl = shareUrls[platform];
        if (!shareUrl) {
            toast.error(`Sharing to ${platform} is not supported`);
            return false;
        }
        try {
            const popup = window.open(
                shareUrl,
                'share',
                'width=600,height=400,scrollbars=yes,resizable=yes,noopener,noreferrer'
            );
            if (popup) {
                toast.success(`Opening ${platform} share...`);
                // Record share attempt
                recordShare({
                    ...shareData,
                    platform
                });
                // Check if popup was blocked
                setTimeout(() => {
                    if (popup.closed) {
                        }
                }, 1000);
                return true;
            } else {
                // Popup was blocked
                toast.error('Popup blocked. Please allow popups and try again.');
                return false;
            }
        } catch (err) {
            toast.error(`Failed to open ${platform} share`);
            return false;
        }
    }, [recordShare]);
    // Comprehensive share handler with intelligent fallbacks
    const handleShare = useCallback(async (shareData) => {
        try {
            setLoading(true);
            setError(null);
            const capabilities = checkBrowserCapabilities();
            const platform = shareData.platform || 'auto';
            // Auto-detect best sharing method
            if (platform === 'auto') {
                if (capabilities.hasWebShare && capabilities.isMobile) {
                    // Mobile with Web Share API - use native
                    return await nativeShare(shareData);
                } else {
                    // Desktop or no Web Share API - use copy
                    const success = await copyToClipboard(shareData.url || window.location.href);
                    if (success) {
                        await recordShare({ ...shareData, platform: 'copy' });
                    }
                    return { success };
                }
            }
            // Specific platform handling
            switch (platform) {
                case 'native':
                    return await nativeShare(shareData);
                case 'copy':
                    const copySuccess = await copyToClipboard(shareData.url || window.location.href);
                    if (copySuccess) {
                        await recordShare({ ...shareData, platform: 'copy' });
                    }
                    return { success: copySuccess };
                case 'facebook':
                case 'twitter':
                case 'x':
                case 'linkedin':
                case 'pinterest':
                case 'reddit':
                case 'telegram':
                case 'whatsapp':
                case 'messenger':
                case 'email':
                    const socialSuccess = socialShare(platform, shareData);
                    return { success: socialSuccess };
                default:
                    throw new Error(`Unknown share platform: ${platform}`);
            }
        } catch (err) {
            setError(err.message);
            // Fallback to copy if all else fails
            if (shareData.platform !== 'copy') {
                toast.warn('Share failed, copying link instead...');
                const fallbackSuccess = await copyToClipboard(shareData.url || window.location.href);
                if (fallbackSuccess) {
                    await recordShare({ ...shareData, platform: 'copy-fallback' });
                }
                return { success: fallbackSuccess, fallback: true };
            }
            toast.error('Share failed. Please try copying the URL manually.');
            return { success: false, error: err.message };
        } finally {
            setLoading(false);
        }
    }, [checkBrowserCapabilities, nativeShare, copyToClipboard, socialShare, recordShare]);
    // Generate share URLs for different platforms
    const generateShareUrls = useCallback((url, title, description) => {
        const encodedUrl = encodeURIComponent(url || window.location.href);
        const encodedTitle = encodeURIComponent(title || document.title);
        const encodedDescription = encodeURIComponent(description || '');
        return {
            facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}`,
            twitter: `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`,
            x: `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`,
            linkedin: `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}`,
            pinterest: `https://pinterest.com/pin/create/button/?url=${encodedUrl}&description=${encodedTitle}`,
            reddit: `https://reddit.com/submit?url=${encodedUrl}&title=${encodedTitle}`,
            telegram: `https://t.me/share/url?url=${encodedUrl}&text=${encodedTitle}`,
            whatsapp: `https://wa.me/?text=${encodedTitle}%20${encodedUrl}`,
            messenger: `https://www.facebook.com/dialog/send?app_id=${process.env.NEXT_PUBLIC_FACEBOOK_APP_ID || ''}&link=${encodedUrl}&redirect_uri=${encodedUrl}`,
            email: `mailto:?subject=${encodedTitle}&body=${encodedDescription}%0A%0A${encodedUrl}`,
        };
    }, []);
    // Check if Web Share API is supported
    const isNativeShareSupported = useCallback(() => {
        const capabilities = checkBrowserCapabilities();
        return capabilities.hasWebShare;
    }, [checkBrowserCapabilities]);
    return {
        // State
        loading,
        error,
        // Methods
        handleShare,
        recordShare,
        updateShareCount,
        getShareAnalytics,
        generateShareUrls,
        nativeShare,
        socialShare,
        copyToClipboard,
        isNativeShareSupported,
        checkBrowserCapabilities,
        // Utilities
        clearError: () => setError(null),
    };
};
export default useShareAPI; 
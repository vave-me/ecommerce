"use client";
import React, { useState, useCallback, useRef, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import { useTranslations } from 'next-intl';
import { 
    Share2, 
    Copy, 
    Facebook, 
    LinkedIn, 
    Mail, 
    MessageCircle,
    X,
    Check,
    Send
} from '@/icons';
import { useAuth } from '../../../context/AuthContext';
import styles from './ShareComponent.module.css';
/**
 * Comprehensive Share Component
 * Supports multiple sharing methods: Web Share API, Social Media, Copy Link, Email
 * Tracks analytics and provides consistent UX across the application
 */
const ShareComponent = memo(({
    // Content to share
    url,
    title,
    description,
    image,
    // Configuration
    variant = 'button', // 'button', 'dropdown', 'modal', 'inline'
    size = 'medium', // 'small', 'medium', 'large'
    platforms = ['native', 'copy', 'facebook', 'x', 'linkedin', 'whatsapp', 'telegram', 'messenger', 'email'], // Available platforms
    showCount = true,
    count = 0,
    // Styling
    className = '',
    buttonText = '',
    iconOnly = false,
    // Callbacks
    onShare,
    onShareSuccess,
    onShareError,
    // Analytics
    trackingData = {},
    // Content identification
    contentId,
    contentType = 'post', // 'post', 'product', 'article', 'profile', etc.
    // Permissions
    requireAuth = false,
}) => {
    const t = useTranslations('Share');
    const { user } = useAuth();
    const [isOpen, setIsOpen] = useState(false);
    const [copied, setCopied] = useState(false);
    const [sharing, setSharing] = useState(false);
    const dropdownRef = useRef(null);
    // Generate share URL if not provided
    const shareUrl = url || (typeof window !== 'undefined' ? window.location.href : '');
    const shareTitle = title || (typeof document !== 'undefined' ? document.title : '');
    const shareDescription = description || '';
    // Facebook App ID for Messenger share (configure in .env)
    const facebookAppId = typeof process !== 'undefined' ? (process.env.NEXT_PUBLIC_FACEBOOK_APP_ID || '') : '';
    // Check if user is authenticated when required
    const canShare = !requireAuth || !!user;
    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsOpen(false);
            }
        };
        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
            return () => document.removeEventListener('mousedown', handleClickOutside);
        }
    }, [isOpen]);
    // Track share events (frontend only)
    const trackShare = useCallback((platform, success = true) => {
        // Google Analytics tracking (if available)
        if (typeof window !== 'undefined' && window.gtag) {
            window.gtag('event', 'share', {
                method: platform,
                content_type: contentType,
                content_id: contentId,
                success,
                ...trackingData
            });
        }
        // Console logging for debugging
        // Custom tracking callback
        if (success && onShareSuccess) {
            onShareSuccess(platform, { contentId, contentType, url: shareUrl });
        } else if (!success && onShareError) {
            onShareError(platform, { contentId, contentType, url: shareUrl });
        }
        // General share callback
        if (onShare) {
            onShare(platform, success, { contentId, contentType, url: shareUrl });
        }
    }, [contentId, contentType, shareUrl, trackingData, onShare, onShareSuccess, onShareError]);
    // Native Web Share API
    const handleNativeShare = useCallback(async () => {
        if (!navigator.share) {
            return false;
        }
        try {
            setSharing(true);
            await navigator.share({
                title: shareTitle,
                text: shareDescription,
                url: shareUrl,
            });
            trackShare('native', true);
            return true;
        } catch (error) {
            if (error.name !== 'AbortError') {
                trackShare('native', false);
            }
            return false;
        } finally {
            setSharing(false);
        }
    }, [shareTitle, shareDescription, shareUrl, trackShare]);
    // Copy to clipboard
    const handleCopyLink = useCallback(async () => {
        try {
            await navigator.clipboard.writeText(shareUrl);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
            trackShare('copy', true);
            return true;
        } catch (error) {
            trackShare('copy', false);
            return false;
        }
    }, [shareUrl, trackShare]);
    // Social media sharing
    const handleSocialShare = useCallback((platform) => {
        const encodedUrl = encodeURIComponent(shareUrl);
        const encodedTitle = encodeURIComponent(shareTitle);
        const encodedDescription = encodeURIComponent(shareDescription);
        let shareUrls = {
            facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}`,
            x: `https://x.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`,
            linkedin: `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}`,
            whatsapp: `https://wa.me/?text=${encodedTitle}%20${encodedUrl}`,
            telegram: `https://t.me/share/url?url=${encodedUrl}&text=${encodedTitle}`,
            messenger: facebookAppId ? `https://www.facebook.com/dialog/send?app_id=${facebookAppId}&link=${encodedUrl}&redirect_uri=${encodedUrl}` : null,
            email: `mailto:?subject=${encodedTitle}&body=${encodedDescription}%0A%0A${encodedUrl}`,
        };
        if (shareUrls[platform]) {
            const popup = window.open(
                shareUrls[platform],
                'share',
                'width=600,height=400,scrollbars=yes,resizable=yes'
            );
            if (popup) {
                trackShare(platform, true);
                // Check if popup was closed (user completed or cancelled)
                const checkClosed = setInterval(() => {
                    if (popup.closed) {
                        clearInterval(checkClosed);
                        // Additional success tracking could be added here
                    }
                }, 1000);
            } else {
                trackShare(platform, false);
            }
        }
    }, [shareUrl, shareTitle, shareDescription, trackShare, facebookAppId]);
    // Main share handler
    const handleShare = useCallback(async (platform = 'native') => {
        if (!canShare) {
            // Show login prompt or handle unauthorized access
            if (onShareError) {
                onShareError('unauthorized', { requireAuth, user: !!user });
            }
            return;
        }
        setIsOpen(false);
        switch (platform) {
            case 'native':
                const nativeSuccess = await handleNativeShare();
                if (!nativeSuccess && platforms.includes('copy')) {
                    // Fallback to copy if native share fails
                    await handleCopyLink();
                }
                break;
            case 'copy':
                await handleCopyLink();
                break;
            case 'facebook':
            case 'x':
            case 'linkedin':
            case 'whatsapp':
            case 'telegram':
            case 'messenger':
            case 'email':
                handleSocialShare(platform);
                break;
            default:
        }
    }, [canShare, requireAuth, user, handleNativeShare, handleCopyLink, handleSocialShare, platforms, onShareError]);
    // Platform configurations
    const platformConfigs = {
        native: {
            icon: Share2,
            label: t('nativeShare', { default: 'Share' }),
            available: typeof navigator !== 'undefined' && !!navigator.share
        },
        copy: {
            icon: copied ? Check : Copy,
            label: copied ? t('copied', { default: 'Copied!' }) : t('copyLink', { default: 'Copy Link' }),
            available: typeof navigator !== 'undefined' && !!navigator.clipboard
        },
        facebook: {
            icon: Facebook,
            label: t('shareOnFacebook', { default: 'Facebook' }),
            available: true
        },
        x: {
            icon: X,
            label: t('shareOnX', { default: 'X' }),
            available: true
        },
        linkedin: {
            icon: LinkedIn,
            label: t('shareOnLinkedIn', { default: 'LinkedIn' }),
            available: true
        },
        whatsapp: {
            icon: MessageCircle,
            label: t('shareOnWhatsApp', { default: 'WhatsApp' }),
            available: true
        },
        telegram: {
            icon: Send,
            label: t('shareOnTelegram', { default: 'Telegram' }),
            available: true
        },
        messenger: {
            icon: MessageCircle,
            label: t('shareOnMessenger', { default: 'Messenger' }),
            available: !!facebookAppId
        },
        email: {
            icon: Mail,
            label: t('shareViaEmail', { default: 'Email' }),
            available: true
        }
    };
    // Filter available platforms
    const availablePlatforms = platforms.filter(platform => 
        platformConfigs[platform]?.available
    );
    // Render share button
    const renderShareButton = () => {
        const buttonClass = `
            ${styles.shareButton} 
            ${styles[`size-${size}`]} 
            ${iconOnly ? styles.iconOnly : ''} 
            ${className}
        `.trim();
        const displayText = buttonText || (iconOnly ? '' : t('share', { default: 'Share' }));
        return (
            <button
                className={buttonClass}
                onClick={() => variant === 'dropdown' ? setIsOpen(!isOpen) : handleShare()}
                disabled={sharing || !canShare}
                aria-label={t('shareContent', { default: 'Share this content' })}
                title={!canShare ? t('loginToShare', { default: 'Login to share' }) : ''}
            >
                <Share2 className={styles.icon} />
                {!iconOnly && (
                    <>
                        {displayText && <span className={styles.text}>{displayText}</span>}
                        {showCount && count > 0 && (
                            <span className={styles.count}>{count}</span>
                        )}
                    </>
                )}
                {sharing && <div className={styles.spinner} />}
            </button>
        );
    };
    // Render platform list
    const renderPlatformList = () => (
        <div className={styles.platformList}>
            {availablePlatforms.map(platform => {
                const config = platformConfigs[platform];
                const IconComponent = config.icon;
                return (
                    <button
                        key={platform}
                        className={`${styles.platformButton} ${copied && platform === 'copy' ? styles.success : ''}`}
                        onClick={() => handleShare(platform)}
                        disabled={sharing}
                    >
                        <IconComponent className={styles.platformIcon} />
                        <span className={styles.platformLabel}>{config.label}</span>
                    </button>
                );
            })}
        </div>
    );
    // Render based on variant
    switch (variant) {
        case 'dropdown':
            return (
                <div className={styles.shareDropdown} ref={dropdownRef}>
                    {renderShareButton()}
                    {isOpen && (
                        <div className={styles.dropdownContent}>
                            <div className={styles.dropdownHeader}>
                                <h4>{t('shareThis', { default: 'Share this' })}</h4>
                                <button 
                                    className={styles.closeButton}
                                    onClick={() => setIsOpen(false)}
                                >
                                    <X size={16} />
                                </button>
                            </div>
                            {renderPlatformList()}
                        </div>
                    )}
                </div>
            );
        case 'modal':
            return (
                <>
                    {renderShareButton()}
                    {isOpen && (
                        <div className={styles.modal}>
                            <div className={styles.modalContent}>
                                <div className={styles.modalHeader}>
                                    <h3>{t('shareContent', { default: 'Share Content' })}</h3>
                                    <button 
                                        className={styles.closeButton}
                                        onClick={() => setIsOpen(false)}
                                    >
                                        <X size={20} />
                                    </button>
                                </div>
                                <div className={styles.modalBody}>
                                    {shareTitle && <h4 className={styles.contentTitle}>{shareTitle}</h4>}
                                    {shareDescription && <p className={styles.contentDescription}>{shareDescription}</p>}
                                    {renderPlatformList()}
                                </div>
                            </div>
                            <div className={styles.modalOverlay} onClick={() => setIsOpen(false)} />
                        </div>
                    )}
                </>
            );
        case 'inline':
            return (
                <div className={`${styles.shareInline} ${className}`}>
                    {renderPlatformList()}
                </div>
            );
        default: // 'button'
            return renderShareButton();
    }
});
ShareComponent.propTypes = {
    // Content
    url: PropTypes.string,
    title: PropTypes.string,
    description: PropTypes.string,
    image: PropTypes.string,
    // Configuration
    variant: PropTypes.oneOf(['button', 'dropdown', 'modal', 'inline']),
    size: PropTypes.oneOf(['small', 'medium', 'large']),
    platforms: PropTypes.arrayOf(PropTypes.oneOf(['native', 'copy', 'facebook', 'x', 'linkedin', 'whatsapp', 'telegram', 'messenger', 'email'])),
    showCount: PropTypes.bool,
    count: PropTypes.number,
    // Styling
    className: PropTypes.string,
    buttonText: PropTypes.string,
    iconOnly: PropTypes.bool,
    // Callbacks
    onShare: PropTypes.func,
    onShareSuccess: PropTypes.func,
    onShareError: PropTypes.func,
    // Analytics
    trackingData: PropTypes.object,
    // Content identification
    contentId: PropTypes.string,
    contentType: PropTypes.string,
    // Permissions
    requireAuth: PropTypes.bool,
};
export default ShareComponent; 
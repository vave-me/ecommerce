"use client";
import React, { forwardRef, useCallback, useRef, useEffect, useMemo, memo } from "react";
import PropTypes from "prop-types";
import {useTranslations} from "next-intl"; //  Import hook
import * as Dialog from "@radix-ui/react-dialog";
import {ChevronDown, User} from "@/icons";
import {useOverlay} from "../../context/OverlayContext";
import {useIsMobile} from "../../hooks/useMobileDetection";
import SignOutMenu from "./SignOutMenu"; // Assumes SignOutMenu handles its own translations
import Image from 'next/image';
import styles from "./UserMenu.module.css";
/**
 * Enhanced UserMenu component with optimized mobile support and performance
 */
const UserMenu = forwardRef(function UserMenu(props, ref) {
    const {
        user,           // User object { avatar, name, email }
        handleSignOut,  // Function to call on sign out action
        isMobile = false // Boolean to adjust styling/layout
    } = props;
    const t = useTranslations('UserMenu'); //  Instantiate hook
    const {setUserMenuOpen} = useOverlay();
    const userButtonRef = useRef(null);
    const [showUserMenu, setShowUserMenu] = React.useState(false);
    // Use optimized mobile detection
    const isMobileDetected = useIsMobile();
    const effectiveMobile = isMobile || isMobileDetected;
    // Memoize user object to prevent recreation
    const memoizedUser = useMemo(() => 
        user ? { ...user, id: user?.userId } : null, 
        [user?.userId, user?.email, user?.name, user?.avatar]
    );
    // Update overlay context when menu state changes
    useEffect(() => {
        setUserMenuOpen(showUserMenu);
    }, [showUserMenu, setUserMenuOpen]);
    /**
     * Closes the menu with improved focus management
     */
    const handleCloseMenu = useCallback(() => {
        if (!showUserMenu) return; // Prevent unnecessary calls
        setShowUserMenu(false);
        // Use requestAnimationFrame to ensure DOM updates before focus
        requestAnimationFrame(() => {
            userButtonRef.current?.focus();
        });
    }, [showUserMenu]);
    /**
     * Toggles the user menu with haptic feedback on mobile
     */
    const handleToggleUserMenu = useCallback(() => {
        setShowUserMenu(prevState => {
            const newState = !prevState;
            // Add haptic feedback for mobile devices
            if (effectiveMobile && 'vibrate' in navigator && newState) {
                navigator.vibrate(50); // Light haptic feedback
            }
            return newState;
        });
    }, [effectiveMobile]);
    /**
     * Enhanced keyboard navigation
     */
    const handleKeyDown = useCallback((e) => {
        if (e.key === "Escape" && showUserMenu) {
            handleCloseMenu();
        }
        // Handle Enter and Space for better accessibility
        if ((e.key === "Enter" || e.key === " ") && !showUserMenu) {
            e.preventDefault();
            handleToggleUserMenu();
        }
    }, [showUserMenu, handleCloseMenu, handleToggleUserMenu]);
    // Memoized computed values
    const shouldShowUsername = useMemo(() => 
        !effectiveMobile && user?.name, 
        [effectiveMobile, user?.name]
    );
    const buttonAriaLabel = useMemo(() => 
        user?.name
            ? t('userMenuAriaLabel', {userName: user.name})
            : t('fallbackMenuAriaLabel'),
        [user?.name, t]
    );
    // Mobile implementation with improved dialog
    if (effectiveMobile) {
        return (
            <Dialog.Root open={showUserMenu} onOpenChange={setShowUserMenu}>
                <Dialog.Trigger asChild>
                    <button
                        type="button"
                        ref={userButtonRef}
                        aria-label={buttonAriaLabel}
                        className={`${styles.menuButton} ${styles.mobileButton} ${showUserMenu ? styles.mobileActive : ''}`}
                        data-mobile="true"
                    >
                        <span className={styles.avatarWrapper}>
                            {user?.avatar ? (
                                <Image
                                    src={user.avatar}
                                    alt=""
                                    className={styles.avatar}
                                    width={32}
                                    height={32}
                                    style={{ objectFit: 'cover' }}
                                    aria-hidden="true"
                                />
                            ) : (
                                <User
                                    size={24}
                                    className={styles.userIcon}
                                    aria-hidden="true"
                                />
                            )}
                        </span>
                    </button>
                </Dialog.Trigger>
                <Dialog.Portal>
                    <Dialog.Overlay className={styles.mobileDialogOverlay} />
                    <Dialog.Content 
                        className={styles.mobileDialogContent}
                        onOpenAutoFocus={(e) => e.preventDefault()} // Prevent auto focus to allow custom focus management
                    >
                        <Dialog.Title className={styles.visuallyHidden}>
                            {t('userMenuTitle', { userName: user?.name || t('userAccount') })}
                        </Dialog.Title>
                        <div className={styles.mobileDialogHeader}>
                            <div className={styles.dragHandle} />
                        </div>
                        <SignOutMenu
                            handleSignOut={handleSignOut}
                            onClose={handleCloseMenu}
                            user={memoizedUser}
                            isMobile={true}
                        />
                    </Dialog.Content>
                </Dialog.Portal>
            </Dialog.Root>
        );
    }
    // Desktop implementation
    return (
        <div className={styles.container} ref={ref} style={{ position: 'relative' }}>
            <button
                type="button"
                ref={userButtonRef}
                onClick={handleToggleUserMenu}
                onKeyDown={handleKeyDown}
                aria-haspopup="menu"
                aria-expanded={showUserMenu}
                aria-controls="user-menu"
                aria-label={buttonAriaLabel}
                className={`${styles.menuButton} ${showUserMenu ? styles.menuActive : ""}`}
            >
                {/* Avatar/User icon */}
                <span className={styles.avatarWrapper}>
                    {user?.avatar ? (
                        <Image
                            src={user.avatar}
                            alt=""
                            className={styles.avatar}
                            width={32}
                            height={32}
                            style={{ objectFit: 'cover' }}
                            aria-hidden="true"
                        />
                    ) : (
                        <User
                            size={24}
                            className={styles.userIcon}
                            aria-hidden="true"
                        />
                    )}
                </span>
                {/* Username (desktop only) */}
                {shouldShowUsername && <span className={styles.username}>{user.name}</span>}
                {/* Chevron icon (desktop only) */}
                {!effectiveMobile && (
                    <ChevronDown
                        size={16}
                        className={`${styles.chevron} ${showUserMenu ? styles.chevronRotate : ""}`}
                        aria-hidden="true"
                    />
                )}
            </button>
            {/* Desktop dropdown menu */}
            {showUserMenu && (
                <SignOutMenu
                    id="user-menu"
                    handleSignOut={handleSignOut}
                    onClose={handleCloseMenu}
                    user={memoizedUser}
                    className={styles.dropdown}
                    isMobile={effectiveMobile}
                    triggerRef={userButtonRef}
                />
            )}
        </div>
    );
});
UserMenu.displayName = "UserMenu";
UserMenu.propTypes = {
    user: PropTypes.shape({
        avatar: PropTypes.string,
        name: PropTypes.string,
        email: PropTypes.string,
        userId: PropTypes.string,
    }),
    handleSignOut: PropTypes.func.isRequired,
    isMobile: PropTypes.bool,
};
UserMenu.defaultProps = {
    user: null,
    isMobile: false,
};
export default memo(UserMenu);
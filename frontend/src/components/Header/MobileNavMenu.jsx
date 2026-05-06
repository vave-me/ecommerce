"use client";
import React, {useCallback, useEffect, useMemo, useRef, useState, memo} from "react";
import PropTypes from "prop-types";
import {useDispatch} from "react-redux";
import {useTranslations} from "next-intl"; //  Import hook
import {Bell, Compass, Heart, LogOut, Mail, ShoppingBag, Tag, User, X, Home, Search, Filter} from "@/icons";
import {CSSTransition} from "react-transition-group";
import FocusTrap from "focus-trap-react";
import Image from 'next/image';
// Assuming openAddDealModal is available from your slices import
import {openAddDealModal} from "../../redux/slices/modalsSlice";
import NavItem from "./NavItem"; // Assume NavItem handles its own translations
import styles from "./MobileNavMenu.module.css";
// PERFORMANCE: Memoized custom hook for reduced motion preference
function usePrefersReducedMotion() {
    const [prefersReducedMotion, setPrefersReducedMotion] = useState(false);
    useEffect(() => {
        const mq = window.matchMedia("(prefers-reduced-motion: reduce)");
        setPrefersReducedMotion(mq.matches);
        const handleChange = () => setPrefersReducedMotion(mq.matches);
        if (mq.addEventListener) {
            mq.addEventListener("change", handleChange);
            return () => mq.removeEventListener("change", handleChange);
        } else {
            // Deprecated fallback for older browsers
            mq.addListener(handleChange);
            return () => mq.removeListener(handleChange);
        }
    }, []);
    return prefersReducedMotion;
}
// PERFORMANCE: Stable navigation items data outside component
const baseNavItemsData = [
    {to: "/home", icon: Home, id: 'home'},
    {to: "/explore", icon: Search, id: 'explore'},
    {to: "/filtered-feed", icon: Filter, id: 'filtered_feed'},
    {to: "/sfx-market", icon: Compass, id: 'sfx_site'},
    {to: "/wishlist", icon: Heart, id: 'wishlist', badgeKey: 'wishlist'},
    {to: "/messages", icon: Mail, id: 'messages', badgeKey: 'messages'},
    {to: "/notifications", icon: Bell, id: 'notifications', badgeKey: 'notifications'},
    {to: "/cart", icon: ShoppingBag, id: 'cart', badgeKey: 'cart'},
];
/**
 * MobileNavMenu Component Refactored with Translations
 * PERFORMANCE OPTIMIZED: Memoized component with stable callbacks and computations
 */
const MobileNavMenu = memo(({isOpen, toggleMenu, locationPath, user, onSignOut, badgeCounts = {}}) => {
    const t = useTranslations('MobileNavMenu'); //  Instantiate hook
    const menuRef = useRef(null);
    const firstLinkRef = useRef(null);
    const prefersReducedMotion = usePrefersReducedMotion();
    const dispatch = useDispatch();
    // PERFORMANCE: Memoized close menu handler
    const closeMenu = useCallback(() => {
        if (isOpen) {
            toggleMenu();
        }
    }, [isOpen, toggleMenu]);
    // PERFORMANCE: Memoized outside click handler
    const handleClickOutside = useCallback((event) => {
        if (menuRef.current && !menuRef.current.contains(event.target)) {
            closeMenu();
        }
    }, [closeMenu]);
    // PERFORMANCE: Memoized escape key handler
    const handleKeyDown = useCallback((event) => {
        if (event.key === "Escape") {
            closeMenu();
        }
    }, [closeMenu]);
    // PERFORMANCE: Memoized deal modal handler
    const handlePostDeal = useCallback(() => {
        try {
            dispatch(openAddDealModal());
            closeMenu();
        } catch (error) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', error);
        }
        // Could set error state here if available
        throw error;
    }
    }, [dispatch, closeMenu]);
    // Add/remove event listeners
    useEffect(() => {
        if (isOpen) {
            const timeoutId = setTimeout(() => {
                document.addEventListener("mousedown", handleClickOutside);
            }, 10); // Small delay to prevent immediate close on toggle click
            document.addEventListener("keydown", handleKeyDown);
            document.body.style.overflow = "hidden"; // Prevent body scroll
            return () => {
                clearTimeout(timeoutId);
                document.removeEventListener("mousedown", handleClickOutside);
                document.removeEventListener("keydown", handleKeyDown);
                document.body.style.overflow = ""; // Restore body scroll
            };
        }
        return undefined; // Explicitly return undefined if not open
    }, [isOpen, handleClickOutside, handleKeyDown]);
    // Auto-focus first element
    useEffect(() => {
        if (isOpen && firstLinkRef.current) {
            const focusTimeoutId = setTimeout(() => {
                firstLinkRef.current.focus();
            }, 150);
            return () => clearTimeout(focusTimeoutId);
        }
    }, [isOpen]);
    // PERFORMANCE: Memoized translated navigation data
    const translatedNavItemsData = useMemo(() => {
        return baseNavItemsData.map(item => ({
            ...item,
            //   Translate label and tooltip
            label: t(`nav_${item.id}_label`),
            tooltip: t(`nav_${item.id}_tooltip`),
            badgeCount: item.badgeKey ? badgeCounts[item.badgeKey] || 0 : 0, // Get badge count safely
        }));
    }, [badgeCounts, t]); // Depend on badgeCounts and t
    // PERFORMANCE: Memoized navigation elements
    const navItemElements = useMemo(() => {
        return translatedNavItemsData.map((item, index) => (
            <NavItem
                key={item.id} // Use stable id as key
                to={item.to}
                icon={item.icon}
                label={item.label} // Pass translated label
                isActive={locationPath === item.to}
                ref={index === 0 ? firstLinkRef : null}
                onClick={closeMenu} // Close menu on item click
                tooltip={item.tooltip} // Pass translated tooltip
                badgeCount={item.badgeCount} // Pass down badge count
                className={item.highlight ? styles.highlightedItem : ''}
                // Assume NavItem handles badge aria-label translation internally
            />
        ));
    }, [translatedNavItemsData, locationPath, closeMenu]); // Depend on translated data
    // PERFORMANCE: Memoized CSS transition classes
    const cssTransitionClasses = useMemo(() => ({
        enter: styles.menuEnter,
        enterActive: styles.menuEnterActive,
        exit: styles.menuExit,
        exitActive: styles.menuExitActive,
    }), []);
    // PERFORMANCE: Memoized sign out handler
    const handleSignOut = useCallback(() => {
        if (onSignOut) {
            closeMenu();
            onSignOut();
        }
    }, [onSignOut, closeMenu]);
    return (
        <>
            {/* Backdrop overlay */}
            <div
                className={`${styles.overlay} ${isOpen ? styles.overlayVisible : ''}`}
                onClick={closeMenu}
                aria-hidden="true"
                data-testid="menu-overlay"
            />
            {/* Slide-in menu */}
            <CSSTransition
                in={isOpen}
                timeout={prefersReducedMotion ? 0 : 300}
                classNames={cssTransitionClasses}
                unmountOnExit
                nodeRef={menuRef}
            >
                {/* Wrapper needed for positioning/transition */}
                <div className={styles.menuWrapper}>
                    <FocusTrap
                        active={isOpen}
                        focusTrapOptions={{
                            initialFocus: false,
                            returnFocusOnDeactivate: true,
                            fallbackFocus: `.${styles.closeButton}`,
                            allowOutsideClick: true,
                        }}
                    >
                        <div
                            ref={menuRef}
                            className={`${styles.menu} ${prefersReducedMotion ? styles.noAnimation : ''}`}
                            //   Use translation
                            aria-label={t('menuAriaLabel')}
                            aria-modal="true"
                            role="dialog"
                        >
                            {/* Header */}
                            <div className={styles.header}>
                                {/*   Use translation */}
                                <h2 id="mobile-menu-title" className={styles.title}>{t('menuTitle')}</h2>
                                <button
                                    type="button"
                                    onClick={closeMenu}
                                    //   Use translation
                                    aria-label={t('closeMenuAriaLabel')}
                                    className={styles.closeButton}
                                >
                                    <X size={24} aria-hidden="true"/>
                                </button>
                            </div>
                            {/* User Info Section */}
                            {user ? (
                                <div className={styles.userSection}>
                                    <div className={styles.avatar}>
                                        {user.photoURL ? (
                                            <Image
                                                src={user.photoURL}
                                                alt=""
                                                width={40}
                                                height={40}
                                                style={{ objectFit: 'cover' }}
                                                aria-hidden="true"
                                            />
                                        ) : (
                                            <User strokeWidth={1.5} size={30} aria-hidden="true"/>
                                        )}
                                    </div>
                                    <div className={styles.userInfo}>
                                        <h3 className={styles.userName}>
                                            {/* Use translated fallback */}
                                            {user.displayName || t('userFallbackName')}
                                        </h3>
                                        {user.email && (
                                            <p className={styles.userEmail}>{user.email}</p>
                                        )}
                                    </div>
                                </div>
                            ) : (
                                <div className={styles.userSection}>
                                    {/* Placeholder or Login prompt */}
                                </div>
                            )}
                            {/* Scrollable Navigation Section */}
                            <div className={styles.navSection}>
                                {/* Navigation Group */}
                                <div className={styles.navGroup}>
                                    {/*   Use translation */}
                                    <h4 className={styles.navGroupTitle}>{t('navigationGroupTitle')}</h4>
                                    {/*   Use translation */}
                                    <ul className={styles.navList} role="menu" aria-label={t('mainNavAriaLabel')}>
                                        {navItemElements}
                                    </ul>
                                </div>
                                {/* Actions Group */}
                                <div className={styles.navGroup}>
                                    {/*   Use translation */}
                                    <h4 className={styles.navGroupTitle}>{t('actionsGroupTitle')}</h4>
                                    {/*   Use translation */}
                                    <ul className={styles.navList} role="menu" aria-label={t('actionsListAriaLabel')}>
                                        <li className={styles.navItem}>
                                            <button
                                                className={styles.dealButton}
                                                onClick={handlePostDeal}
                                                //   Use translation
                                                aria-label={t('postDealButtonAriaLabel')}
                                                role="menuitem"
                                            >
                                                <Tag size={20} aria-hidden="true"/>
                                                {/*   Use translation */}
                                                <span>{t('postDealButtonText')}</span>
                                            </button>
                                        </li>
                                    </ul>
                                </div>
                                {/* Account Actions (Sign Out) */}
                                {user && onSignOut && (
                                    <div className={styles.accountActions}>
                                        <button onClick={handleSignOut} className={styles.signOutButton}
                                                role="menuitem">
                                            <LogOut size={18} aria-hidden="true"/>
                                            {/*   Use translation */}
                                            <span>{t('signOutButton')}</span>
                                        </button>
                                    </div>
                                )}
                            </div>
                        </div>
                    </FocusTrap>
                </div>
            </CSSTransition>
        </>
    );
});
// Updated PropTypes to include badgeCounts
MobileNavMenu.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    toggleMenu: PropTypes.func.isRequired,
    locationPath: PropTypes.string.isRequired,
    user: PropTypes.object,
    onSignOut: PropTypes.func,
    badgeCounts: PropTypes.shape({
        wishlist: PropTypes.number,
        messages: PropTypes.number,
        notifications: PropTypes.number,
        cart: PropTypes.number,
    })
};
// Default props for badgeCounts
MobileNavMenu.defaultProps = {
    user: null,
    onSignOut: null,
    badgeCounts: {}
};
export default MobileNavMenu;

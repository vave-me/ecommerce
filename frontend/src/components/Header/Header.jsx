"use client";
import React, { useCallback, useContext, useEffect, useReducer, useState, memo } from "react";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "../../context/AuthContext";
import { NavBarContext } from "../../context/NavBarContext";
import { useScrollNavigation } from "../../hooks/useScrollNavigation";
import { useDispatch, useSelector } from "react-redux";
import { useTranslations } from "next-intl";
import { Bell, Heart, MessageCircle, ShoppingBag, Plus, User } from "@/icons";

// Child Components
import Logo from "./Logo";
import LogoMobile from "./LogoMobile";
import SearchBar from "./SearchBar";
import MobileNavMenu from "./MobileNavMenu";
import UserMenu from "./UserMenu";
import BottomNav from "./BottomNav";
import AddDropdown from "./AddDropdown";
import ModeSwitcher from "./ModeSwitcher";
import HeaderIcon from './HeaderIcon';
import { useBasket } from '../../hooks/useBasket';
import useWishlist from '../../hooks/useWishlist';
import UtilityBar from '../UtilityBar';

// Redux Actions
import {
    openAddDealModal,
    openAddJobModal,
    openAddPostModal,
    openAddProductModal,
    openAddPropertyModal,
    openAddServiceModal,
    openAddVehicleModal,
    openAddVideoModal
} from "../../redux/slices/modalsSlice";

// App Mode selectors
import {
    selectIsAiMode,
    selectCurrentMode,
    switchToAiMode,
    APP_MODES
} from "../../redux/slices/appModeSlice";

// Styles
import styles from "./Header.module.css";

// Action map for opening modals
const addModalActions = {
    product: openAddProductModal,
    post: openAddPostModal,
    video: openAddVideoModal,
    vehicle: openAddVehicleModal,
    deal: openAddDealModal,
    property: openAddPropertyModal,
    service: openAddServiceModal,
    job: openAddJobModal,
};

// Header state reducer
const headerReducer = (state, action) => {
    switch (action.type) {
        case 'TOGGLE_MOBILE_MENU':
            return { ...state, isMobileMenuOpen: !state.isMobileMenuOpen };
        case 'CLOSE_MOBILE_MENU':
            return { ...state, isMobileMenuOpen: false };
        case 'SET_SCROLLED':
            return { ...state, isScrolled: action.payload };
        case 'TOGGLE_AI_COPILOT':
            return { ...state, showAiCopilot: !state.showAiCopilot };
        case 'RESET_MENUS':
            return { ...state, isMobileMenuOpen: false };
        default:
            return state;
    }
};

const Header = memo(() => {
    const t = useTranslations('Header');
    const { showNavbars, isMobile } = useContext(NavBarContext);
    const { user, signOutUser, isLoading, authChecked } = useAuth();
    const router = useRouter();
    const pathname = usePathname();
    const reduxDispatch = useDispatch();

    // App mode selectors
    const isAiMode = useSelector(selectIsAiMode);
    const currentMode = useSelector(selectCurrentMode);
    const [showAddDropdown, setShowAddDropdown] = useState(false);
    
    // Get basket and wishlist counts
    const { itemCount: basketCount } = useBasket();
    const { totalItemsCount: wishlistCount } = useWishlist();

    // State management with reducer
    const [headerState, headerDispatch] = useReducer(headerReducer, {
        isMobileMenuOpen: false,
        isScrolled: false,
        showAiCopilot: false,
    });

    // Scroll-based navigation visibility (only on mobile)
    const { isNavVisible, isScrollingDown } = useScrollNavigation(80, 40);

    // Handle scroll effect for header shadow
    useEffect(() => {
        const handleScroll = () => {
            const scrollPosition = window.scrollY;
            headerDispatch({
                type: 'SET_SCROLLED',
                payload: scrollPosition > 20,
            });
        };
        window.addEventListener("scroll", handleScroll, { passive: true });
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    // Reset mobile menu on route change
    useEffect(() => {
        if (headerState.isMobileMenuOpen) {
            headerDispatch({ type: 'RESET_MENUS' });
        }
    }, [pathname, headerState.isMobileMenuOpen]);

    // Dropdown handlers
    const toggleAddDropdown = useCallback(() => {
        setShowAddDropdown((prev) => !prev);
    }, []);

    const closeAddDropdown = useCallback(() => {
        setShowAddDropdown(false);
    }, []);

    // Refactored handler for triggering add modals
    const handleAdd = useCallback((type) => {
        const actionCreator = addModalActions[type];
        if (actionCreator) {
            reduxDispatch(actionCreator());
        }
        closeAddDropdown();
    }, [reduxDispatch, closeAddDropdown]);

    // Handle sign out
    const handleSignOut = useCallback(async () => {
        try {
            await signOutUser();
            router.push("/login");
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
    }, [signOutUser, router]);

    // Toggle AI Copilot
    const toggleAiCopilot = useCallback(() => {
        headerDispatch({ type: 'TOGGLE_AI_COPILOT' });
    }, []);

    // Hardcoded badge counts - replace with actual data source
    const badgeCounts = {
        messages: 5,
        notifications: 1,
        wishlist: 0,
    };

    // Visibility conditions
    const shouldShowHeader = !isMobile || showNavbars;
    const shouldShowBottomNav = isMobile && !headerState.isMobileMenuOpen;

    // Don't render header until authentication is checked
    if (!authChecked) {
        return null;
    }

    // Combine classes with conditional styles
    const headerClasses = [
        styles.container,
        headerState.isScrolled ? styles.scrolled : '',
        isAiMode ? styles.aiMode : '',
        isMobile && !isNavVisible ? styles.hidden : '',
        isMobile ? styles.mobileScrollable : ''
    ].filter(Boolean).join(' ');

    // AI Mode Minimal Header
    if (isAiMode) {
        return (
            <>
                {shouldShowHeader && (
                    <header className={headerClasses} role="banner">
                        <nav className={styles.topNav} role="navigation" aria-label={t('mainNavAriaLabel')}>
                            <div className={styles.aiModeHeader}>
                                <div className={styles.aiModeLogo}>
                                    {isMobile ? <LogoMobile aiMode={true} /> : <Logo aiMode={true} />}
                                </div>
                                <div className={styles.aiModeSwitcher}>
                                    <ModeSwitcher isMobile={isMobile} showSubtitles={false} compact={true} />
                                </div>
                                <div className={styles.aiModeUser}>
                                    <UserMenu
                                        user={user}
                                        handleSignOut={handleSignOut}
                                        isMobile={isMobile}
                                    />
                                </div>
                            </div>
                        </nav>
                    </header>
                )}
            </>
        );
    }

    // Classic Mode Header - Matching finalConceptTOP design
    return (
        <>
            {/* Utility Bar - Desktop only */}
            {!isMobile && <UtilityBar />}
            
            {shouldShowHeader && (
                <header className={headerClasses} role="banner">
                    <div className={styles.headerWrapper}>
                        <nav className={styles.topNav} role="navigation" aria-label={t('mainNavAriaLabel')}>
                            {isMobile ? (
                                /* Mobile Layout - Simplified with only search and profile */
                                <div className={styles.mobileHeader}>
                                    <div className={styles.mobileSearchWrapper}>
                                        <SearchBar />
                                    </div>
                                    <div className={styles.mobileProfileSection}>
                                        {!user ? (
                                            <button
                                                className={styles.mobileLoginButton}
                                                onClick={() => router.push('/login')}
                                                aria-label={t('loginButtonAriaLabel')}
                                            >
                                                <User size={24} />
                                            </button>
                                        ) : (
                                            <UserMenu
                                                user={user}
                                                handleSignOut={handleSignOut}
                                                isMobile={true}
                                            />
                                        )}
                                    </div>
                                </div>
                            ) : (
                                /* Desktop Layout - Matching finalConceptTOP */
                                <div className={styles.desktopHeader}>
                                    {/* Left: Logo & AI Mode Button */}
                                    <div className={styles.headerLeft}>
                                        <Logo />
                                        <button 
                                            className={styles.aiModeButton}
                                            onClick={() => {
                                                reduxDispatch(switchToAiMode());
                                            }}
                                            title={t('aiAssistantModeTooltip')}
                                        >
                                            <svg className={styles.aiIcon} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                                                <path fillRule="evenodd" d="M9.315 7.584C12.195 3.883 16.695 1.5 21.75 1.5a.75.75 0 0 1 .75.75c0 5.056-2.383 9.555-6.084 12.436A6.75 6.75 0 0 1 9.75 22.5a.75.75 0 0 1-.75-.75v-7.19c0-.897.106-1.773.315-2.622ZM3 10.5a.75.75 0 0 1 .75-.75h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1-.75-.75Zm.75 4.5a.75.75 0 0 0 0-1.5h-1.5a.75.75 0 0 0 0 1.5h1.5ZM4.5 6a.75.75 0 0 1 .75-.75h4.5a.75.75 0 0 1 0 1.5h-4.5A.75.75 0 0 1 4.5 6Z" clipRule="evenodd"/>
                                            </svg>
                                            <span className={styles.aiModeButtonText}>{t('aiAssistantMode', 'KI Assistant Mode')}</span>
                                        </button>
                                    </div>

                                    {/* Center: Search Bar */}
                                    <div className={styles.headerCenter}>
                                        <SearchBar placeholder={t('searchPlaceholder', 'Search products, brands, or categories...')} />
                                    </div>

                                    {/* Right: Actions */}
                                    <div className={styles.headerRight}>
                                        {/* Admin Create/Add Dropdown - Keeping for admin functionality */}
                                        {user && user.role === 'admin' && (
                                            <>
                                                <HeaderIcon
                                                    icon={Plus}
                                                    onClick={toggleAddDropdown}
                                                    ariaLabel={t('createContentAriaLabel')}
                                                    title={t('createButtonTooltip')}
                                                    className={styles.createButton}
                                                    isMobile={false}
                                                    ariaExpanded={showAddDropdown}
                                                    ariaHaspopup="true"
                                                />
                                                <AddDropdown
                                                    isOpen={showAddDropdown}
                                                    onClose={closeAddDropdown}
                                                    onAddProduct={() => handleAdd("product")}
                                                    onAddPost={() => handleAdd("post")}
                                                    onAddVideo={() => handleAdd("video")}
                                                    onAddVehicle={() => handleAdd("vehicle")}
                                                    onAddDeal={() => handleAdd("deal")}
                                                    onAddProperty={() => handleAdd("property")}
                                                    onAddService={() => handleAdd("service")}
                                                    onAddJob={() => handleAdd("job")}
                                                />
                                            </>
                                        )}

                                        {/* Show login button when not logged in */}
                                        {!user ? (
                                            <>
                                                {/* Wishlist - Always visible */}
                                                <HeaderIcon
                                                    icon={Heart}
                                                    onClick={() => router.push('/wishlist')}
                                                    badge={wishlistCount}
                                                    ariaLabel={`Wishlist with ${wishlistCount} item${wishlistCount !== 1 ? 's' : ''}`}
                                                    title="Wishlist"
                                                    isMobile={false}
                                                />
                                                
                                                {/* Cart/Basket - Always visible */}
                                                <HeaderIcon
                                                    icon={ShoppingBag}
                                                    onClick={() => router.push('/cart')}
                                                    badge={basketCount}
                                                    ariaLabel={`Shopping basket with ${basketCount} item${basketCount !== 1 ? 's' : ''}`}
                                                    title="Cart"
                                                    isMobile={false}
                                                />
                                                
                                                <button
                                                    className={styles.loginButton}
                                                    onClick={() => router.push('/login')}
                                                    aria-label={t('loginButtonAriaLabel')}
                                                >
                                                    {t('login', 'Login')}
                                                </button>
                                            </>
                                        ) : (
                                            <>
                                                {/* Messages */}
                                                <HeaderIcon
                                                    icon={MessageCircle}
                                                    onClick={() => router.push('/messages')}
                                                    badge={badgeCounts.messages}
                                                    ariaLabel={t('messagesButtonAriaLabel', { count: badgeCounts.messages || 0 })}
                                                    title={t('messages')}
                                                    isMobile={false}
                                                />
                                                
                                                {/* Notifications */}
                                                <HeaderIcon
                                                    icon={Bell}
                                                    onClick={() => router.push('/notifications')}
                                                    badge={badgeCounts.notifications}
                                                    ariaLabel={t('notificationsButtonAriaLabel', { count: badgeCounts.notifications || 0 })}
                                                    title={t('notifications')}
                                                    isMobile={false}
                                                />
                                                
                                                {/* Wishlist - Always visible */}
                                                <HeaderIcon
                                                    icon={Heart}
                                                    onClick={() => router.push('/wishlist')}
                                                    badge={wishlistCount}
                                                    ariaLabel={`Wishlist with ${wishlistCount} item${wishlistCount !== 1 ? 's' : ''}`}
                                                    title="Wishlist"
                                                    isMobile={false}
                                                />
                                                
                                                {/* Cart/Basket - Always visible */}
                                                <HeaderIcon
                                                    icon={ShoppingBag}
                                                    onClick={() => router.push('/cart')}
                                                    badge={basketCount}
                                                    ariaLabel={`Shopping basket with ${basketCount} item${basketCount !== 1 ? 's' : ''}`}
                                                    title="Cart"
                                                    isMobile={false}
                                                />

                                                {/* Profile/User Menu */}
                                                <UserMenu
                                                    user={user}
                                                    handleSignOut={handleSignOut}
                                                    isMobile={false}
                                                />
                                            </>
                                        )}
                                    </div>
                                </div>
                            )}
                        </nav>
                    </div>
                </header>
            )}

            {/* Bottom Navigation for Mobile - Only show in Classic mode */}
            {!isAiMode && shouldShowBottomNav && (
                <BottomNav
                    locationPath={pathname}
                    badgeCounts={badgeCounts}
                    onOpenAddDealModal={() => reduxDispatch(openAddDealModal())}
                    navConfig="default"
                />
            )}
        </>
    );
});

Header.displayName = 'Header';
export default Header;
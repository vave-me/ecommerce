"use client";
import React, {useContext, useCallback, useState, useEffect, useRef, memo} from "react";
import {usePathname, useRouter} from "next/navigation";
import {useLocale} from "next-intl";
import {useSelector} from "react-redux";
import {NavBarContext} from "../context/NavBarContext";
import {useOverlay} from "../context/OverlayContext";
import {useIsMobile} from "../hooks/useMobileDetection";
import {selectIsAiMode} from "../redux/slices/appModeSlice";
import Header from "../components/Header/Header";
import CategoryBar from "../components/Header/CategoryBar";
import HorizontalFilters from "../components/Filters/HorizontalFilters";
import BottomNav from "../components/Header/BottomNav";

import GlobalModals from "../components/Utils/GlobalModals";
import WishlistSelectorModal from "../components/wishlist/WishlistSelectorModal";
import TouchInteractions from "../components/Shared/TouchInteractions";
import AIPageComponent from "./[locale]/ai/AIPage";
import Footer from "../components/Footer/Footer";
import { ErrorBoundary } from "../components/ErrorBoundary";
import styles from "./ClientLayout.module.css";
const ClientLayout = memo(function ClientLayout({children}) {
    const {isMobile, showNavbars} = useContext(NavBarContext);
    const {hasActiveOverlay} = useOverlay();
    const pathname = usePathname();
    const locale = useLocale();
    const router = useRouter();
    const isAiMode = useSelector(selectIsAiMode);
    const currentMode = useSelector(state => state.appMode?.currentMode);
    
    // Debug logging
    useEffect(() => {
        
    }, [isAiMode, currentMode, pathname, locale]);
    
    // Use optimized mobile detection hook
    const isClientMobile = useIsMobile();
    // Routes where we hide navigation completely
    const hideNavigationRoutes = ["/login", "/register", "/signup", "/auth"];
    const shouldHideNavigation = hideNavigationRoutes.some(route =>
        pathname.includes(route) || pathname.endsWith(route)
    );
    // Some routes where we skip the header on mobile
    const skipHeaderOnMobileRoutes = ["/messages"];
    const hideHeaderOnMobile = isClientMobile && skipHeaderOnMobileRoutes.some(route =>
        pathname.includes(route)
    );
    // Also hide bottom nav on messaging pages (same logic as header)
    const hideBottomNavOnMobile = isClientMobile && skipHeaderOnMobileRoutes.some(route =>
        pathname.includes(route)
    );
    
    // Routes where horizontal filters should NOT be displayed
    const hideFiltersRoutes = [
        // Authentication & Account Access
        "/login", "/register", "/signup", "/auth",
        "/forgot-password", "/reset-password", "/verify",
        
        // Checkout & Payment Flow
        "/checkout", "/payments",
        
        // User Account Management
        "/settings", "/orders", "/wishlist", "/notifications", "/messages",
        "/profile", "/user", "/followers", "/following",
        
        // Static/Informational Pages
        "/about", "/contact", "/help", "/privacy", "/terms", "/support",
        "/faq", "/resources", "/legal",
        
        // Admin Panel
        "/admin",
        
        // Demo/Marketing Pages
        "/demo", "/sell",
        
        // Special Content Pages
        "/design", "/newsletters", "/ai",
        
        // Other Non-Browsing Pages
        "/cart", "/reviews", "/activity"
    ];
    
    const shouldHideFilters = hideFiltersRoutes.some(route =>
        pathname.includes(route)
    );
    // Use the state for rendering decisions
    const shouldShowBottomNav = isClientMobile && !isAiMode && !shouldHideNavigation && !hideBottomNavOnMobile;
    // Mobile gesture handlers (only enabled on mobile)
    const handleSwipeLeft = useCallback(() => {
        // Could be used for navigation - currently disabled to prevent conflicts
        // Example: navigate to next section or open sidebar
    }, []);
    const handleSwipeRight = useCallback(() => {
        // Could be used for navigation - currently disabled to prevent conflicts
        // Example: navigate back or open menu
    }, []);
    const handleSwipeUp = useCallback(() => {
        // Could be used for quick actions - currently disabled
        // Example: scroll to top or show search
    }, []);
    const handleSwipeDown = useCallback(() => {
        // Could be used for refresh - handled by individual components
        // Example: pull-to-refresh (handled by InfiniteScroll components)
    }, []);
    // Navigation is now handled by ModeSwitcher component directly
    // This prevents navigation loops and race conditions
    // AI Mode: Render ONLY AI Results
    if (isAiMode) {
        const aiLayoutContent = (
            <div className={styles.layout}>
                {(!hideHeaderOnMobile && !shouldHideNavigation) && (
                    <ErrorBoundary name="Header" fallback={<div>Header failed to load</div>}>
                        <Header/>
                    </ErrorBoundary>
                )}
                <main className={styles.content}>
                    <ErrorBoundary name="AIPageComponent" fallback={<div>AI mode failed to load</div>}>
                        <AIPageComponent/>
                    </ErrorBoundary>
                </main>
                <ErrorBoundary name="Footer" fallback={null}>
                    <Footer/>
                </ErrorBoundary>
                <ErrorBoundary name="GlobalModals" fallback={null}>
                    <GlobalModals/>
                </ErrorBoundary>
                <ErrorBoundary name="WishlistSelectorModal" fallback={null}>
                    <WishlistSelectorModal/>
                </ErrorBoundary>
            </div>
        );
        // Wrap with TouchInteractions only on mobile to avoid desktop interference
        if (isClientMobile) {
            return (
                <TouchInteractions
                    enableSwipeGestures={false} // Disabled by default to prevent navigation conflicts
                    enablePullToRefresh={false} // Handled by individual components
                    enableLongPress={false}     // Disabled to prevent conflicts with native browser behavior
                    hapticFeedback={true}
                    onSwipeLeft={handleSwipeLeft}
                    onSwipeRight={handleSwipeRight}
                    onSwipeUp={handleSwipeUp}
                    onSwipeDown={handleSwipeDown}
                    className="mobile-layout-wrapper"
                    style={{minHeight: '100vh'}}
                >
                    {aiLayoutContent}
                </TouchInteractions>
            );
        }
        // Desktop: return AI layout without touch interactions
        return aiLayoutContent;
    }
    // Classic Mode: Render normal page content
    const layoutContent = (
        <div className={styles.layout}>
            {(!hideHeaderOnMobile && !shouldHideNavigation) && (
                <>
                    <ErrorBoundary name="Header" fallback={<div>Header failed to load</div>}>
                        <Header/>
                    </ErrorBoundary>
                    {/* HorizontalFilters - responsive behavior */}
                    {!shouldHideFilters && (
                        <ErrorBoundary name="HorizontalFilters" fallback={null}>
                            <HorizontalFilters categoryType="marketplace" />
                        </ErrorBoundary>
                    )}
                </>
            )}
            <main className={styles.content}>
                <ErrorBoundary name="MainContent" fallback={<div>Content failed to load</div>}>
                    {children}
                </ErrorBoundary>
            </main>
            <ErrorBoundary name="Footer" fallback={null}>
                <Footer/>
            </ErrorBoundary>
            {shouldShowBottomNav && (
                <ErrorBoundary name="BottomNav" fallback={null}>
                    <BottomNav/>
                </ErrorBoundary>
            )}

            <ErrorBoundary name="GlobalModals" fallback={null}>
                <GlobalModals/>
            </ErrorBoundary>
            <ErrorBoundary name="WishlistSelectorModal" fallback={null}>
                <WishlistSelectorModal/>
            </ErrorBoundary>
        </div>
    );
    // Wrap with TouchInteractions only on mobile to avoid desktop interference
    if (isClientMobile) {
        return (
            <TouchInteractions
                enableSwipeGestures={false} // Disabled by default to prevent navigation conflicts
                enablePullToRefresh={false} // Handled by individual components
                enableLongPress={false}     // Disabled to prevent conflicts with native browser behavior
                hapticFeedback={true}
                onSwipeLeft={handleSwipeLeft}
                onSwipeRight={handleSwipeRight}
                onSwipeUp={handleSwipeUp}
                onSwipeDown={handleSwipeDown}
                className="mobile-layout-wrapper"
                style={{minHeight: '100vh'}}
            >
                {layoutContent}
            </TouchInteractions>
        );
    }
    // Desktop: return layout without touch interactions
    return layoutContent;
});
export default ClientLayout;

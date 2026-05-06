"use client";
import React, { useState, useMemo, useCallback, memo } from "react";
import { useRouter } from "next/navigation";
import { useDispatch } from "react-redux";
import PropTypes from "prop-types";
import { useTranslations } from "next-intl"; //  Import hook
import { useScrollNavigation } from "../../hooks/useScrollNavigation";
import { useAuth } from "../../context/AuthContext";
import {
    Home, ShoppingCart, Heart, MessageCircle, Bell,
    User, Compass, Tag, BellRing, Plus, Bot,
    Package, FileText, Video, Car, Home as HomeIcon, Wrench, Briefcase
} from "@/icons";
// Import Redux Modal Actions
import {
    openAddDealModal, openAddJobModal, openAddPostModal, openAddProductModal,
    openAddPropertyModal, openAddServiceModal, openAddVehicleModal, openAddVideoModal
} from "../../redux/slices/modalsSlice"; // Adjust path as needed
// Import the dedicated AddOptionsSheet components
import AddOptionsSheet from './AddOptionsSheet'; // Assumes AddOptionsSheet uses its own translations
import AddOptionsSheetWithComposer from './AddOptionsSheetWithComposer';
import styles from "./BottomNav.module.css";
// Action map for opening modals (remains the same)
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
// Base structure for add options items (without text)
const baseAddOptionsItems = [
    { id: 'product', icon: Package }, { id: 'post', icon: FileText },
    { id: 'video', icon: Video }, { id: 'vehicle', icon: Car },
    { id: 'deal', icon: Tag }, { id: 'property', icon: HomeIcon },
    { id: 'service', icon: Wrench }, { id: 'job', icon: Briefcase },
];
/**
 * Enhanced Mobile Bottom Navigation Bar with improved UI/UX
 * Uses next-intl for translations of nav items and add options.
 * Includes scroll-based hide/show behavior matching Header component.
 */
function BottomNav({
                       locationPath,
                       badgeCounts = {},
                       navConfig = 'default'
                   }) {
    const t = useTranslations('BottomNav'); //  Instantiate hook
    const router = useRouter();
    const dispatch = useDispatch();
    const { user } = useAuth();
    const [isAddOptionsVisible, setIsAddOptionsVisible] = useState(false);
    // SCROLL BEHAVIOR HOOK - RE-ENABLED & SYNCHRONIZED
    const { isNavVisible } = useScrollNavigation(80, 40);
    // Create translated items for the AddOptionsSheet
    const translatedAddOptionsItems = useMemo(() => {
        return baseAddOptionsItems.map(item => ({
            ...item,
            //   Translate label and description
            label: t(`additem_${item.id}_label`),
            description: t(`additem_${item.id}_desc`),
        }));
    }, [t]);
    // Handler to open the action sheet
    const handlePrimaryActionClick = useCallback(() => {
        // For now, always show the add options sheet
        // TODO: Implement unified composer for customers
        setIsAddOptionsVisible(true);
    }, []);
    // Handler for selecting an item from the action sheet
    const handleSelectAddAction = useCallback((type) => {
        const actionCreator = addModalActions[type];
        if (actionCreator) {
            try {
                dispatch(actionCreator());
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
        } else {
        }
        setIsAddOptionsVisible(false);
    }, [dispatch]);
    // Close handler for the action sheet
    const handleCloseAddOptions = useCallback(() => {
        setIsAddOptionsVisible(false);
    }, []);
    // Handle regular navigation for other items
    const handleNavigation = useCallback((path) => {
        router.push(path);
    }, [router]);
    // Define and translate navigation items using useMemo
    const translatedNavItemsData = useMemo(() => {
        const defaultNavItems = [
            { to: "/home", icon: Home, id: 'home' },
            { to: "/notifications", icon: Bell, id: 'alerts', badgeKey: 'notifications' },
            { to: "/add", icon: Bot, id: 'add', isPrimaryAction: true, tooltipId: 'default' },
            { to: "/wishlist", icon: Heart, id: 'wishlist', badgeKey: 'wishlist' },
            { to: "/cart", icon: ShoppingCart, id: 'cart', badgeKey: 'cart' },
        ];
        const alternateNavItems = [
            { to: "/", icon: Home, id: 'home' }, // Root path might be same as home
            { to: "/cart", icon: ShoppingCart, id: 'cart' },
            { to: "/add", icon: Plus, id: 'add', isPrimaryAction: true, tooltipId: 'alternate' },
            { to: "/wishlist", icon: Heart, id: 'wishlist', badgeKey: 'wishlist' },
            { to: "/notifications", icon: Bell, id: 'alerts', badgeKey: 'notifications' }, // Use 'alerts' consistent ID
        ];
        const itemsToTranslate = navConfig === 'alternate' ? alternateNavItems : defaultNavItems;
        return itemsToTranslate.map(item => ({
            ...item,
            //   Translate label and tooltip
            label: t(`nav_${item.id}_label`),
            tooltip: t(item.isPrimaryAction ? `nav_add_tooltip_${item.tooltipId}` : `nav_${item.id}_tooltip`),
            badgeCount: item.badgeKey ? badgeCounts[item.badgeKey] : 0, // Get badge count
        }));
    }, [navConfig, badgeCounts, t]);
    // Apply conditional classes for scroll behavior
    const bottomNavClasses = `
        ${styles.bottomNav}
        ${!isNavVisible ? styles.hidden : ''}
    `;
    return (
        <>
            <nav className={bottomNavClasses} aria-label={t('navAriaLabel')}> {/*  Use translation */}
                <ul className={styles.navList}>
                    {translatedNavItemsData.map((item) => (
                        <li key={item.id} className={styles.navItem}>
                            <button
                                onClick={item.isPrimaryAction ? handlePrimaryActionClick : () => handleNavigation(item.to)}
                                className={`
                                    ${styles.navButton}
                                    ${locationPath === item.to && !item.isPrimaryAction ? styles.active : ''}
                                    ${item.isPrimaryAction ? styles.addNavButton : ''}
                                `}
                                aria-label={item.tooltip} // Use translated tooltip
                                title={item.tooltip} // Use translated tooltip
                                aria-haspopup={item.isPrimaryAction ? "dialog" : undefined}
                                aria-expanded={item.isPrimaryAction ? isAddOptionsVisible : undefined}
                            >
                                <div className={styles.iconWrapper}>
                                    <item.icon
                                        size={item.isPrimaryAction ? 26 : 22}
                                        strokeWidth={item.isPrimaryAction ? 2.5 : 2}
                                        aria-hidden="true"
                                    />
                                    {!item.isPrimaryAction && item.badgeCount > 0 && (
                                        <span
                                            className={styles.badge}
                                            //   Use translation with pluralization
                                            aria-label={t('badgeAriaLabel', { count: item.badgeCount })}
                                        >
                                            {item.badgeCount > 99 ? '99+' : item.badgeCount}
                                        </span>
                                    )}
                                </div>
                            </button>
                        </li>
                    ))}
                </ul>
            </nav>
            {/* Render the enhanced AddOptionsSheetWithComposer for ALL users */}
            <AddOptionsSheetWithComposer
                isOpen={isAddOptionsVisible}
                onClose={handleCloseAddOptions}
                onSelectItem={handleSelectAddAction}
                items={translatedAddOptionsItems}
                showUnifiedComposer={true}
            />
        </>
    );
}
BottomNav.propTypes = {
    locationPath: PropTypes.string,
    badgeCounts: PropTypes.shape({
        wishlist: PropTypes.number,
        messages: PropTypes.number,
        notifications: PropTypes.number,
        cart: PropTypes.number
    }),
    navConfig: PropTypes.oneOf(['default', 'alternate'])
};
BottomNav.defaultProps = {
    locationPath: "",
    badgeCounts: {},
    navConfig: 'default'
};
export default memo(BottomNav);
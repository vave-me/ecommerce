"use client";
import React, { useMemo, useCallback, useState, useRef, useEffect, memo } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useTranslations } from "next-intl"; //  Import hook
import {
    Heart, MessageCircle, Bell, ShoppingCart, Plus, ChevronDown,
    Image, Tag, Car, Home as HomeIcon, Briefcase, Wrench, Video as VideoIcon // Renamed Home import
} from "@/icons";
import NavItem from "./NavItem"; // Assume NavItem handles its own internal translations (e.g., badge ARIA label)
import {
    openAddProductModal, openAddPostModal, openAddVideoModal, openAddVehicleModal,
    openAddDealModal, openAddPropertyModal, openAddServiceModal, openAddJobModal,
} from "../../redux/slices/modalsSlice"; // Adjust path as needed
import WaveComponent from "./WaveComponent"; // Assume WaveComponent handles its own translations
import styles from "./DesktopNav.module.css";
// Base navigation items configuration (without text)
const baseNavItemsData = [
    { to: "/wishlist", icon: Heart, id: 'wishlist', badgeCount: 2 }, // Example badge count
    { to: "/messages", icon: MessageCircle, id: 'messages', badgeCount: 5 },
    { to: "/notifications", icon: Bell, id: 'alerts', badgeCount: 1 },
    { to: "/cart", icon: ShoppingCart, id: 'cart' },
];
// Base content type configuration (without text)
const baseContentTypes = [
    { id: "deal", icon: Tag, primary: true },
    { id: "post", icon: Image }, // Consider a different icon for Post?
    { id: "video", icon: VideoIcon }, // Use renamed VideoIcon
    { id: "property", icon: HomeIcon },
    { id: "vehicle", icon: Car },
    { id: "service", icon: Wrench },
    { id: "job", icon: Briefcase }
];
const DesktopNav = memo(function DesktopNav({ locationPath }) {
    const t = useTranslations('DesktopNav'); //  Instantiate hook
    const dispatch = useDispatch();
    const [showDropdown, setShowDropdown] = useState(false);
    const dropdownRef = useRef(null);
    const buttonRef = useRef(null);
    // Toggle dropdown visibility
    const toggleDropdown = useCallback(() => {
        setShowDropdown(prev => !prev);
    }, []);
    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current &&
                !dropdownRef.current.contains(event.target) &&
                buttonRef.current &&
                !buttonRef.current.contains(event.target)) {
                setShowDropdown(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);
    // Handle keyboard navigation (Escape key)
    useEffect(() => {
        const handleEscape = (event) => {
            if (event.key === 'Escape' && showDropdown) {
                setShowDropdown(false);
            }
        };
        document.addEventListener('keydown', handleEscape);
        return () => document.removeEventListener('keydown', handleEscape);
    }, [showDropdown]);
    // Handle different types of "Create" actions (dispatching logic remains the same)
    const handleCreate = useCallback(
        (type) => {
            const actionMap = {
                product: openAddProductModal, post: openAddPostModal, video: openAddVideoModal,
                vehicle: openAddVehicleModal, deal: openAddDealModal, property: openAddPropertyModal,
                service: openAddServiceModal, job: openAddJobModal
            };
            const actionCreator = actionMap[type];
            if (actionCreator) {
                dispatch(actionCreator());
            }
            setShowDropdown(false); // Close dropdown after dispatch
        },
        [dispatch]
    );
    // Create translated navigation item components using useMemo
    const translatedNavItems = useMemo(() => {
        return baseNavItemsData.map((item) => {
            const isActive = locationPath === item.to;
            const badgeCount = item.badgeCount || 0;
            //   Translate label and tooltip
            const label = t(`navitem_${item.id}_label`);
            const tooltip = t(`navitem_${item.id}_tooltip`);
            return (
                <NavItem
                    key={item.id}
                    to={item.to}
                    icon={item.icon}
                    label={label} // Pass translated label
                    tooltip={tooltip} // Pass translated tooltip
                    isActive={isActive}
                    badgeCount={badgeCount}
                    // Pass necessary translation key for badge ARIA label if NavItem needs it
                    // Or better, NavItem uses useTranslations internally
                />
            );
        });
    }, [locationPath, t]);
    // Create translated content types for the dropdown using useMemo
    const translatedContentTypes = useMemo(() => {
        return baseContentTypes.map(type => ({
            ...type,
            //   Translate label and description
            label: t(`content_${type.id}_label`),
            description: t(`content_${type.id}_desc`),
        }));
    }, [t]);
    return (
        <nav className={styles.navContainer} aria-label={t('navAriaLabel')}> {/*  Use translation */}
            <div className={styles.createContainer}>
                {/* Primary "Post Deal" button */}
                <button
                    className={styles.primaryButton}
                    onClick={() => handleCreate("deal")}
                    aria-label={t('postDealAriaLabel')} //  Use translation
                >
                    <Plus size={16} aria-hidden="true" />
                    {/*   Use translation */}
                    <span>{t('postDealButton')}</span>
                </button>
                {/* Create dropdown button */}
                <button
                    ref={buttonRef}
                    className={`${styles.dropdownToggle} ${showDropdown ? styles.active : ''}`}
                    onClick={toggleDropdown}
                    aria-expanded={showDropdown}
                    aria-controls="create-dropdown-menu" // Added aria-controls
                    aria-haspopup="true"
                    aria-label={t('showMoreOptionsAriaLabel')} //  Use translation
                >
                    <ChevronDown size={16} aria-hidden="true" />
                </button>
                {/* Dropdown menu for all content types */}
                {showDropdown && (
                    <div
                        id="create-dropdown-menu" // Added id matching aria-controls
                        ref={dropdownRef}
                        className={styles.dropdown}
                        role="menu"
                        aria-orientation="vertical" // Added orientation
                    >
                        <div className={styles.dropdownHeader}>
                            {/*   Use translation */}
                            <span>{t('createContentDropdownHeader')}</span>
                        </div>
                        <div className={styles.dropdownContent}>
                            {/* Render translated content types */}
                            {translatedContentTypes.map((type) => (
                                <button
                                    key={type.id}
                                    className={`${styles.dropdownItem} ${type.primary ? styles.primaryItem : ''}`}
                                    onClick={() => handleCreate(type.id)}
                                    role="menuitem"
                                >
                                    <span className={styles.itemIcon} aria-hidden="true">
                                        <type.icon size={18} />
                                    </span>
                                    <div className={styles.itemContent}>
                                        {/* Display translated label and description */}
                                        <span className={styles.itemTitle}>{type.label}</span>
                                        <span className={styles.itemDescription}>{type.description}</span>
                                    </div>
                                </button>
                            ))}
                        </div>
                        {/* Assuming WaveComponent handles its own internal translations */}
                        <div className={styles.dropdownFooter}>
                            <WaveComponent
                                onClose={() => setShowDropdown(false)}
                                onAddPost={() => handleCreate("post")}
                                onAddVideo={() => handleCreate("video")}
                            />
                        </div>
                    </div>
                )}
            </div>
            {/* Navigation Items */}
            <ul className={styles.navItems} role="menubar" aria-label={t('navAriaLabel')}> {/* Re-use nav aria label */}
                {/* Render translated NavItem components */}
                {translatedNavItems}
            </ul>
        </nav>
    );
});
DesktopNav.propTypes = {
    locationPath: PropTypes.string.isRequired,
};
export default DesktopNav;
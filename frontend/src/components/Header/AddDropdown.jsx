"use client";
import React, {useEffect, useMemo, useRef, useState, memo} from 'react';
import {createPortal} from 'react-dom';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl'; //  Import hook
import {Briefcase, Car, Home, ShoppingCart, Wrench, Video} from '@/icons';
import {FaNewspaper, FaPercentage} from '../../utils/iconImports';
import styles from './AddDropdown.module.css';
import {useIsMobile} from "../../hooks/useMobileDetection";
// Base structure without translatable text
const baseMenuItems = [
    {
        id: 'frequent',
        items: [
            {id: 'product', icon: ShoppingCart, notification: 0},
            {id: 'post', icon: FaNewspaper, notification: 0},
            {id: 'video', icon: Video, notification: 0},
            {id: 'service', icon: Wrench, notification: 0},
        ]
    }
];
const AddDropdown = memo(({
                         isOpen,
                         onClose,
                         onAddProduct,
                         onAddPost,
                         onAddVideo,
                         onAddService,
                         onAddVehicle,
                         onAddDeal,
                         onAddProperty,
                         onAddJob
                     }) => {
    const t = useTranslations('AddDropdown'); //  Instantiate hook
    const dropdownRef = useRef(null);
    const [activeSection, setActiveSection] = useState('frequent');
    const [recentlyUsed, setRecentlyUsed] = useState(['product', 'post', 'video', 'service']);
    const isMobile = useIsMobile(); // Use optimized mobile detection
    const [buttonRect, setButtonRect] = useState(null);
    // Create translated menu items structure using useMemo
    const translatedMenuItems = useMemo(() => {
        return baseMenuItems.map(section => ({
            ...section,
            //   Translate section title
            title: t(`section_${section.id}_title`),
            items: section.items.map(item => ({
                ...item,
                //   Translate item label and description
                label: t(`item_${item.id}_label`),
                description: t(`item_${item.id}_desc`),
            }))
        }));
    }, [t]); // Recompute if the translation function changes (e.g., locale change)
    // Get button position when opened and update on scroll
    useEffect(() => {
        if (isOpen && !isMobile) {
            // Find the create button more reliably using multiple selectors
            const button = document.querySelector('.createButton, [aria-haspopup="true"][aria-expanded="true"], button[aria-label*="create"]') ||
                          document.querySelector('button[aria-haspopup="true"]');
            
            if (button) {
                const updateButtonRect = () => {
                    const rect = button.getBoundingClientRect();
                    setButtonRect(rect);
                };
                
                // Set initial position
                updateButtonRect();
                
                // Update position on scroll and resize
                window.addEventListener('scroll', updateButtonRect, { passive: true });
                window.addEventListener('resize', updateButtonRect, { passive: true });
                
                return () => {
                    window.removeEventListener('scroll', updateButtonRect);
                    window.removeEventListener('resize', updateButtonRect);
                };
            }
        }
    }, [isOpen, isMobile]);
    // Handle clicks outside the dropdown and escape key
    useEffect(() => {
        if (!isOpen) return;
        const handleClickOutside = (e) => {
            // Check if the click is outside the dropdown
            if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
                // Also check if it's not the trigger button (to prevent immediate close/open)
                const triggerButton = document.querySelector('.createButton, [aria-haspopup="true"][aria-expanded="true"], button[aria-label*="create"]') ||
                                    document.querySelector('button[aria-expanded="true"]');
                if (!triggerButton || !triggerButton.contains(e.target)) {
                    onClose();
                }
            }
        };
        const handleEscapeKey = (e) => {
            if (e.key === 'Escape') {
                onClose();
            }
        };
        // Add a small delay to avoid conflicts with the button click
        const timeoutId = setTimeout(() => {
            document.addEventListener('mousedown', handleClickOutside);
            document.addEventListener('touchstart', handleClickOutside);
            document.addEventListener('keydown', handleEscapeKey);
        }, 100);
        return () => {
            clearTimeout(timeoutId);
            document.removeEventListener('mousedown', handleClickOutside);
            document.removeEventListener('touchstart', handleClickOutside);
            document.removeEventListener('keydown', handleEscapeKey);
        };
    }, [isOpen, onClose]);
    // Map action IDs to their corresponding handlers
    const actionHandlers = {
        product: onAddProduct,
        post: onAddPost,
        video: onAddVideo,
        service: onAddService,
        vehicle: onAddVehicle,
        deal: onAddDeal,
        property: onAddProperty,
        job: onAddJob
    };
    // Handle item selection (no change needed here)
    const handleItemClick = (itemId) => {
        if (actionHandlers[itemId]) {
            setRecentlyUsed(prev => {
                const newRecent = [itemId, ...prev.filter(id => id !== itemId)].slice(0, 3);
                return newRecent;
            });
            actionHandlers[itemId]();
            onClose(); // Close dropdown after selection
        }
    };
    // Helper to find item data across all sections (using translated data)
    const findItemData = (itemId) => {
        for (const section of translatedMenuItems) {
            const item = section.items.find(i => i.id === itemId);
            if (item) return item;
        }
        return null;
    }
    // Existing dropdown content (unchanged)
    const dropdownContent = (
        <>
            <header className={styles.dropdownHeader}>
                <h3 className={styles.dropdownTitle}>{t('title')}</h3>
            </header>
            {/* Menu items for the active section */}
            <div className={styles.menuItems}>
                {translatedMenuItems.find(section => section.id === activeSection)?.items.map(item => (
                    <button
                        key={item.id}
                        className={styles.menuItem}
                        onClick={() => handleItemClick(item.id)}
                        role="menuitem"
                    >
                        <span className={styles.menuItemIcon} aria-hidden="true">
                            <item.icon/>
                        </span>
                        <span className={styles.menuItemContent}>
                            <span className={styles.menuItemLabel}>{item.label}</span>
                            <span className={styles.menuItemDescription}>{item.description}</span>
                        </span>
                        {item.notification > 0 && (
                            <span className={styles.notification}
                                  aria-label={`${item.notification} updates`}>{item.notification}</span>
                        )}
                    </button>
                ))}
            </div>
        </>
    );
    // Don't render if not open
    if (!isOpen) {
        return null;
    }
    // Calculate position for desktop dropdown
    const getDropdownStyle = () => {
        if (!buttonRect || isMobile) return {};
        
        const dropdownWidth = 260; // Width from CSS
        const margin = 8; // Margin from button
        
        // Calculate positioning to align dropdown with button
        let left = buttonRect.left;
        let top = buttonRect.bottom + margin;
        
        // Adjust if dropdown would overflow viewport on the right
        if (left + dropdownWidth > window.innerWidth) {
            left = buttonRect.right - dropdownWidth;
        }
        
        // Adjust if dropdown would overflow viewport on the left
        if (left < margin) {
            left = margin;
        }
        
        // Adjust if dropdown would overflow viewport on the bottom
        const dropdownMaxHeight = 400; // Approximate max height
        if (top + dropdownMaxHeight > window.innerHeight) {
            top = buttonRect.top - dropdownMaxHeight - margin;
            // If still overflows, place it at the top with some margin
            if (top < margin) {
                top = margin;
            }
        }
        
        return {
            position: 'fixed',
            top: Math.max(margin, top),
            left: Math.max(margin, left),
            zIndex: 999999
        };
    };
    // Create dropdown content
    const dropdown = isMobile ? (
        <>
            <div 
                className={styles.mobileBackdrop} 
                onClick={(e) => {
                    e.stopPropagation();
                    onClose();
                }} 
            />
            <div
                ref={dropdownRef}
                role="menu"
                aria-label={t('ariaLabel')}
                className={`${styles.dropdownContainer} ${styles.mobileContainer}`}
                onClick={(e) => e.stopPropagation()}
            >
                {dropdownContent}
            </div>
        </>
    ) : (
        <div
            ref={dropdownRef}
            role="menu"
            aria-label={t('ariaLabel')}
            className={styles.dropdownContainer}
            style={getDropdownStyle()}
            onClick={(e) => e.stopPropagation()}
        >
            {dropdownContent}
        </div>
    );
    // Use portal for desktop to escape stacking context, direct render for mobile
    if (typeof window === 'undefined') return null;
    return isMobile ? dropdown : createPortal(dropdown, document.body);
});
AddDropdown.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onAddProduct: PropTypes.func,
    onAddPost: PropTypes.func,
    onAddVideo: PropTypes.func,
    onAddService: PropTypes.func,
    onAddVehicle: PropTypes.func,
    onAddDeal: PropTypes.func,
    onAddProperty: PropTypes.func,
    onAddJob: PropTypes.func,
};
AddDropdown.displayName = 'AddDropdown';
export default AddDropdown;
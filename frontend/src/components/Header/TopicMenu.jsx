"use client";
import { X, FaBars } from '../../icons';
import React, { useRef, useCallback, forwardRef, useMemo, useEffect, memo } from "react";
import PropTypes from "prop-types";
import {useDispatch} from "react-redux"; // Keep if needed for future actions
import {useRouter} from "next/navigation";
import { useTranslations } from "next-intl"; //  Import hook
import styles from "./TopicMenu.module.css";
// Base structure for categories with stable IDs and paths
const baseCategories = [
    {id: "deals", path: "/deals"},
    {id: "marketplace", path: "/products"},
    {id: "electronics", path: "/products/electronics"},
    {id: "fashion", path: "/products/clothing"},
    {id: "home", path: "/products/home"},
    {id: "cars", path: "/cars"},
    {id: "services", path: "/services"},
];
/**
 * Enhanced Topic Menu component for DealSocial with Translations
 *
 * A hamburger menu that opens a dropdown with category options.
 */
const TopicMenu = forwardRef(function TopicMenu(
    {toggleTopicMenu, showTopicMenu},
    ref // Forwarded ref is attached to the container div
) {
    const t = useTranslations('TopicMenu'); //  Instantiate hook
    const topicButtonRef = useRef(null); // Ref for the toggle button itself
    const dispatch = useDispatch(); // Keep for potential future use
    const router = useRouter();
    // Translate category labels using useMemo
    const translatedCategories = useMemo(() => {
        return baseCategories.map(category => ({
            ...category,
            //   Translate label using derived key
            label: t(`category_${category.id}_label`)
        }));
    }, [t]); // Dependency on translation function
    /**
     * Close the menu and restore focus to the toggle button
     */
    const handleCloseMenu = useCallback(() => {
        toggleTopicMenu(false); // Use the passed function to update parent state
        // Check if button ref exists before focusing
        topicButtonRef.current?.focus();
    }, [toggleTopicMenu]);
    /**
     * Toggle the menu open/closed using the passed function
     */
    const handleToggleTopicMenu = useCallback(() => {
        toggleTopicMenu(!showTopicMenu); // Use the passed function
    }, [toggleTopicMenu, showTopicMenu]);
    /**
     * Navigate to a category when selected
     */
    const handleCategoryClick = useCallback((path) => {
        router.push(path);
        handleCloseMenu();
    }, [router, handleCloseMenu]);
    // Add Escape key listener to close menu
    useEffect(() => {
        if (!showTopicMenu) return; // Only add listener when menu is open
        const handleKeyDown = (event) => {
            if (event.key === "Escape") {
                handleCloseMenu();
            }
        };
        document.addEventListener("keydown", handleKeyDown);
        return () => {
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, [showTopicMenu, handleCloseMenu]);
    return (
        // Attach forwarded ref to the main container
        <div className={styles.topicMenuContainer} ref={ref}>
            <button
                type="button"
                ref={topicButtonRef} // Ref for focus management
                onClick={handleToggleTopicMenu}
                aria-haspopup="menu"
                aria-expanded={showTopicMenu}
                aria-controls="topic-menu"
                //   Use translation
                aria-label={t('toggleButtonAriaLabel')}
                className={`${styles.topicButton} ${showTopicMenu ? styles.active : ''}`}
            >
                {/* Conditionally render icons */}
                {showTopicMenu ? (
                    <X aria-hidden="true" size={20}/>
                ) : (
                    <FaBars aria-hidden="true" size={20}/>
                )}
            </button>
            {/* Dropdown Menu - Render based on showTopicMenu state */}
            {showTopicMenu && (
                <div
                    id="topic-menu" // Matches aria-controls
                    className={styles.menuDropdown}
                    role="menu"
                    //   Use translation
                    aria-label={t('dropdownAriaLabel')}
                >
                    <div className={styles.menuHeader}>
                        {/*   Use translation */}
                        <h3 className={styles.menuTitle}>{t('dropdownTitle')}</h3>
                        <button
                            className={styles.closeButton}
                            onClick={handleCloseMenu}
                            //   Use translation
                            aria-label={t('closeButtonAriaLabel')}
                        >
                            <X size={18} aria-hidden="true"/>
                        </button>
                    </div>
                    <ul className={styles.categoryList} role="presentation"> {/* ul is container, role on items */}
                        {/* Render translated categories */}
                        {translatedCategories.map(category => (
                            <li key={category.id} role="none"> {/* List item is presentational */}
                                <button
                                    className={styles.categoryButton}
                                    role="menuitem"
                                    onClick={() => handleCategoryClick(category.path)}
                                >
                                    {/* Display translated label */}
                                    {category.label}
                                </button>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
            {/* Overlay for mobile - closes menu on click */}
            {showTopicMenu && (
                <div
                    className={styles.menuOverlay}
                    onClick={handleCloseMenu}
                    aria-hidden="true"
                />
            )}
        </div>
    );
});
TopicMenu.displayName = "TopicMenu";
TopicMenu.propTypes = {
    toggleTopicMenu: PropTypes.func.isRequired,
    showTopicMenu: PropTypes.bool.isRequired,
};
export default memo(TopicMenu);
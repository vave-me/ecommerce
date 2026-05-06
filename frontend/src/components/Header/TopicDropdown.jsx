"use client";
import React, { useRef, useState, useLayoutEffect, useEffect } from 'react';
import { createPortal } from 'react-dom';
import PropTypes from 'prop-types';
import { useTranslations } from 'next-intl';
import DropdownHeader from './DropdownHeader';
import CategoryList from './CategoryList';
import styles from './TopicDropdown.module.css';
/**
 * TopicDropdown Component
 * Dropdown/sheet containing categories with subcategory navigation
 * Refactored to use extracted sub-components while maintaining 100% functionality
 */
function TopicDropdown({
    topicData,
    isOpen,
    onClose,
    onSelectCategory,
    selectedCategoryValue,
    anchorEl,
    isMobileView,
    id,
    activeSubcategories,
    isLoadingSubcategories,
    selectedParentId,
    onShowSubcategories
}) {
    const dropdownRef = useRef(null);
    const [positionStyle, setPositionStyle] = useState({});
    const [pointerStyle, setPointerStyle] = useState({});
    const t = useTranslations('Topics');
    // Define translated titles for categories
    const categoriesTitle = t('categories');
    const subcategoriesTitle = t('subcategories');
    // Positioning logic for desktop popover - FIXED to position under clicked element
    useLayoutEffect(() => {
        if (isOpen && !isMobileView && anchorEl && dropdownRef.current) {
            // Calculate position relative to the container
            const containerElement = dropdownRef.current.offsetParent; // The .container element
            const anchorRect = anchorEl.getBoundingClientRect();
            const containerRect = containerElement?.getBoundingClientRect();
            if (containerRect) {
                // Calculate the anchor's position relative to the container
                const anchorLeftOffset = anchorRect.left - containerRect.left;
                const anchorWidth = anchorRect.width;
                const anchorCenter = anchorLeftOffset + (anchorWidth / 2);
                let positionStyles = {
                    left: `${anchorLeftOffset}px`
                };
                // Set pointer position to center of anchor element
                let pointerStyles = {
                    left: `${anchorCenter - anchorLeftOffset - 6}px` // 6px is half of pointer width (12px)
                };
                // Check if dropdown would overflow right edge of viewport
                const dropdownRect = dropdownRef.current.getBoundingClientRect();
                const dropdownRightEdge = containerRect.left + anchorLeftOffset + dropdownRect.width;
                if (dropdownRightEdge > window.innerWidth - 8) {
                    // Position dropdown to align its right edge with anchor's right edge
                    const newLeftOffset = anchorLeftOffset + anchorWidth - dropdownRect.width;
                    positionStyles = {
                        left: `${newLeftOffset}px`
                    };
                    // Adjust pointer position relative to new dropdown position
                    pointerStyles = {
                        left: `${anchorCenter - newLeftOffset - 6}px`
                    };
                    // If still overflowing, align to container's right edge
                    if (newLeftOffset < 0) {
                        positionStyles = {
                            right: '8px',
                            left: 'auto'
                        };
                        // Calculate pointer position when dropdown is right-aligned
                        const rightAlignedLeft = containerRect.width - dropdownRect.width - 8;
                        pointerStyles = {
                            left: `${anchorCenter - rightAlignedLeft - 6}px`
                        };
                    }
                }
                setPositionStyle(positionStyles);
                setPointerStyle(pointerStyles);
            }
        } else {
            setPositionStyle({});
            setPointerStyle({});
        }
    }, [isOpen, isMobileView, anchorEl]);
    // Escape key handler
    useEffect(() => {
        if (!isOpen) return;
        const handleKeyDown = (event) => {
            if (event.key === 'Escape') onClose();
        };
        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [isOpen, onClose]);
    // Click outside handler for desktop
    useEffect(() => {
        if (!isOpen || isMobileView) return;
        const handleClickOutside = (event) => {
            if (dropdownRef.current && 
                !dropdownRef.current.contains(event.target) && 
                anchorEl && 
                !anchorEl.contains(event.target)) {
                onClose();
            }
        };
        const timeoutId = setTimeout(() => 
            document.addEventListener('mousedown', handleClickOutside), 0
        );
        return () => {
            clearTimeout(timeoutId);
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isOpen, isMobileView, onClose, anchorEl]);
    if (!isOpen) return null;
    // Destructure state safely
    const {data, isLoading, error} = topicData.categories || {
        data: null, 
        isLoading: false, 
        error: null
    };
    // Handle back to main categories
    const handleBackToMain = () => {
        onShowSubcategories(null);
    };
    // Determine which categories to display (main or subcategories)
    const displayCategories = selectedParentId ? activeSubcategories : data;
    const isDisplayLoading = selectedParentId ? isLoadingSubcategories : isLoading;
    // Title based on if we're showing main or sub categories
    const titleText = selectedParentId ? subcategoriesTitle : categoriesTitle;
    const commonContent = (
        <>
            <h4 className={styles.sectionTitle}>{titleText}</h4>
            <CategoryList
                displayCategories={displayCategories}
                isDisplayLoading={isDisplayLoading}
                error={error}
                selectedParentId={selectedParentId}
                selectedCategoryValue={selectedCategoryValue}
                onSelectCategory={onSelectCategory}
                onShowSubcategories={onShowSubcategories}
                onBackToMain={handleBackToMain}
            />
        </>
    );
    // Create dropdown content
    const dropdownContent = (
        <>
            {isMobileView && <div className={styles.backdrop} onClick={onClose}/>}
            <div
                id={id}
                ref={dropdownRef}
                className={`${styles.dropdownContainer} ${isMobileView ? styles.dropdownSheet : styles.dropdownPopover}`}
                style={!isMobileView ? positionStyle : {}}
                role="dialog"
                aria-modal={isMobileView ? "true" : "false"}
                aria-labelledby={`topic-dropdown-title-${topicData.value}`}
            >
                {!isMobileView && <div className={styles.popoverPointer} style={pointerStyle}/>}
                {isMobileView && (
                    <DropdownHeader
                        topicData={topicData}
                        selectedParentId={selectedParentId}
                        onClose={onClose}
                    />
                )}
                <div className={styles.dropdownContent}>
                    {commonContent}
                </div>
            </div>
        </>
    );
    // Use portal for mobile to escape stacking context, direct render for desktop
    if (typeof window === 'undefined') return null;
    return isMobileView ? createPortal(dropdownContent, document.body) : dropdownContent;
}
// PropTypes for subcategories support
TopicDropdown.propTypes = {
    topicData: PropTypes.shape({
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        badge: PropTypes.string,
        categories: PropTypes.shape({
            data: PropTypes.arrayOf(PropTypes.shape({
                value: PropTypes.string.isRequired,
                label: PropTypes.string.isRequired,
                id: PropTypes.string,
                featured: PropTypes.bool,
                count: PropTypes.number,
                hasSubcategories: PropTypes.bool,
                subcategories: PropTypes.array
            })),
            isLoading: PropTypes.bool,
            error: PropTypes.object,
        }),
    }).isRequired,
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onSelectCategory: PropTypes.func.isRequired,
    selectedCategoryValue: PropTypes.string.isRequired,
    anchorEl: PropTypes.object,
    isMobileView: PropTypes.bool.isRequired,
    id: PropTypes.string.isRequired,
    activeSubcategories: PropTypes.array,
    isLoadingSubcategories: PropTypes.bool,
    selectedParentId: PropTypes.string,
    onShowSubcategories: PropTypes.func.isRequired
};
export default TopicDropdown; 
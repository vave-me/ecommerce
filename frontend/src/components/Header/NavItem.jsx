"use client";
import React, { memo } from 'react';
import {Link} from "@/i18n/navigation";
import PropTypes from "prop-types";
import {useTranslations} from "next-intl"; //  Import hook
import styles from "./NavItem.module.css";
/**
 * NavItem Component (Handles badge translation)
 *
 * - Renders as <Link> or <button>
 * - Supports notification badges with translated ARIA labels
 * - Assumes 'label' and 'tooltip' props are pre-translated
 */
const NavItem = React.forwardRef(
    (
        {
            to,
            icon: Icon,
            label,      // Expect pre-translated label
            isMobile,   // Prop not used in current rendering logic, but kept
            onClick,
            isActive,
            badgeCount,
            tooltip,    // Expect pre-translated tooltip
            className,
        },
        ref
    ) => {
        const t = useTranslations('NavItem'); //  Instantiate hook
        // Determine if this is a link or button
        const isLink = Boolean(to);
        // Generate container class names
        const containerClassNames = [
            styles.navItemContainer,
            isMobile ? styles.isMobile : "", // Keep isMobile class if needed for styling context
            isActive ? styles.active : "",
            className || "",
        ]
            .filter(Boolean)
            .join(" ");
        // Generate translated ARIA label for the badge span
        const badgeAriaText = badgeCount > 0
            ? t('badgeAriaLabel', {count: badgeCount, label: label}) // Pass count and label for context
            : '';
        // Build the icon with optional badge counter
        const iconWithBadge = (
            <div className={styles.iconContainer}>
                <Icon className={styles.icon} aria-hidden="true" focusable="false"/>
                {badgeCount > 0 && (
                    <span
                        className={styles.badge}
                        //   Use translated badge ARIA text
                        aria-label={badgeAriaText}
                        data-count={badgeCount} // Keep data-count for potential CSS use
                    >
                        {badgeCount > 99 ? '99+' : badgeCount}
                    </span>
                )}
            </div>
        );
        // Render label visually (only used for mobile in original, adapt if needed)
        // Currently not rendered visually by default in this simplified structure
        // const labelElement = isMobile && (
        //     <span className={styles.label}>{label}</span>
        // );
        // Common props for Link and Button
        const commonProps = {
            "aria-label": label, // Use pre-translated label as the main ARIA label
            onClick: onClick,
            ref: ref,
            title: tooltip, // Use pre-translated tooltip
            role: "menuitem", // Assuming use within a menu structure like menubar or menu
        };
        // If we have a `to` prop, render as a Next.js Link
        if (isLink) {
            return (
                <li className={containerClassNames} role="none"> {/* List item is presentational */}
                    <Link
                        href={to}
                        aria-current={isActive ? "page" : undefined}
                        className={styles.navItemLink} // Combine styles if needed
                        {...commonProps}
                    >
                        {iconWithBadge}
                        {/* labelElement could be added here if needed */}
                    </Link>
                </li>
            );
        }
        // Otherwise, render as a button
        return (
            <li className={containerClassNames} role="none"> {/* List item is presentational */}
                <button
                    type="button"
                    className={styles.navItemButton} // Combine styles if needed
                    {...commonProps}
                >
                    {iconWithBadge}
                    {/* labelElement could be added here if needed */}
                </button>
            </li>
        );
    }
);
NavItem.displayName = "NavItem";
NavItem.propTypes = {
    to: PropTypes.string,
    icon: PropTypes.elementType.isRequired,
    label: PropTypes.string.isRequired, // Should be pre-translated
    isMobile: PropTypes.bool,
    onClick: PropTypes.func,
    isActive: PropTypes.bool,
    badgeCount: PropTypes.number,
    tooltip: PropTypes.string, // Should be pre-translated
    className: PropTypes.string,
};
NavItem.defaultProps = {
    to: null,
    isMobile: false,
    onClick: undefined,
    isActive: false,
    badgeCount: 0,
    tooltip: "",
    className: "",
};
export default memo(NavItem);

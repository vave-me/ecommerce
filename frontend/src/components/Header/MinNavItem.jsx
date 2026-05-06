"use client";
import React, { memo } from 'react';
import Link from "next/link";
import PropTypes from "prop-types";
import { useTranslations } from "next-intl"; //  Import hook
import styles from "./MinNavItem.module.css";
const MinNavItem = React.forwardRef(
    ({
         to,
         icon: Icon,
         label, // Assume pre-translated
         onClick,
         isActive,
         badgeCount,
         tooltip // Assume pre-translated
     }, ref) => {
        const t = useTranslations('MinNavItem'); //  Instantiate hook
        // Determine if this is a link or button
        const isLink = Boolean(to);
        // Generate translated ARIA label for the badge count
        const badgeAriaText = badgeCount > 0
            ? t('badgeAriaLabel', { count: badgeCount })
            : '';
        // Generate combined ARIA label for the main element (Link or Button)
        // Combines the pre-translated label with the translated badge text if applicable
        const mainAriaLabel = `${label}${badgeAriaText ? ` (${badgeAriaText})` : ''}`;
        // Common icon + optional badge
        const iconWithBadge = (
            <div className={styles.iconWrapper}>
                <Icon aria-hidden="true" className={styles.icon} />
                {badgeCount > 0 && (
                    <span
                        className={styles.badge}
                        //   Use translated badge ARIA text
                        aria-label={badgeAriaText}
                    >
                        {badgeCount > 99 ? '99+' : badgeCount}
                    </span>
                )}
                {/* Render pre-translated label */}
                <span className={styles.label}>{label}</span>
            </div>
        );
        // If we have a 'to' prop, render a Next.js Link
        if (isLink) {
            return (
                <li
                    className={`${styles.navItem} ${isActive ? styles.active : ''}`}
                    role="none" // List item is presentational
                >
                    <Link
                        href={to}
                        //   Use combined ARIA label
                        aria-label={mainAriaLabel}
                        aria-current={isActive ? "page" : undefined}
                        title={tooltip} // Use pre-translated tooltip
                        className={styles.navLink}
                        role="menuitem"
                        ref={ref}
                    >
                        {iconWithBadge}
                    </Link>
                </li>
            );
        }
        // Otherwise, render a button
        return (
            <li
                className={`${styles.navItem} ${isActive ? styles.active : ''}`}
                role="none" // List item is presentational
            >
                <button
                    type="button"
                    //   Use combined ARIA label
                    aria-label={mainAriaLabel}
                    onClick={onClick}
                    title={tooltip} // Use pre-translated tooltip
                    className={styles.navButton}
                    role="menuitem"
                    ref={ref}
                >
                    {iconWithBadge}
                </button>
            </li>
        );
    }
);
MinNavItem.displayName = "MinNavItem";
MinNavItem.propTypes = {
    to: PropTypes.string,
    icon: PropTypes.elementType.isRequired,
    label: PropTypes.string.isRequired, // Should be pre-translated
    onClick: PropTypes.func,
    isActive: PropTypes.bool,
    badgeCount: PropTypes.number,
    tooltip: PropTypes.string, // Should be pre-translated
};
MinNavItem.defaultProps = {
    to: null,
    onClick: () => {},
    isActive: false,
    badgeCount: 0,
    tooltip: "",
};
// Use React.memo for potential performance optimization
export default memo(MinNavItem);
import React, {memo} from 'react';
import {Package, Search, Plus} from '@/icons';
import styles from './EmptyPlaceholder.module.css';
/**
 * EmptyPlaceholder - Atomic Design Component
 * Displays an empty state with message and optional action
 *
 * @param {Object} props - Component props
 * @param {string} props.label - Empty state message
 * @param {Function} props.onAction - Action function
 * @param {string} props.actionLabel - Action button label
 * @param {string} props.title - Empty state title
 * @param {string} props.icon - Icon type (package, search, plus)
 * @param {string} props.size - Size variant (sm, md, lg)
 * @returns {JSX.Element} Rendered empty placeholder
 */
const EmptyPlaceholder = memo(({
                                   label = 'No items found',
                                   onAction = null,
                                   actionLabel = 'Add Item',
                                   title = 'Nothing here yet',
                                   icon = 'package',
                                   size = 'md'
                               }) => {
    const getIcon = () => {
        const iconProps = {
            className: `${styles.icon} ${styles[`icon${size.charAt(0).toUpperCase() + size.slice(1)}`]}`
        };
        switch (icon) {
            case 'search':
                return <Search {...iconProps} />;
            case 'plus':
                return <Plus {...iconProps} />;
            case 'package':
            default:
                return <Package {...iconProps} />;
        }
    };
    return (
        <div className={styles.container}>
            <div className={`${styles.content} ${styles[size]}`}>
                <div className={styles.iconWrapper}>
                    {getIcon()}
                </div>
                <div className={styles.textContent}>
                    <h3 className={styles.title}>{title}</h3>
                    <p className={styles.message}>{label}</p>
                </div>
                {onAction && (
                    <button
                        className={styles.actionButton}
                        onClick={onAction}
                        aria-label={actionLabel}
                    >
                        <Plus className={styles.actionIcon} size={16}/>
                        <span>{actionLabel}</span>
                    </button>
                )}
            </div>
        </div>
    );
});
EmptyPlaceholder.displayName = 'EmptyPlaceholder';
export default EmptyPlaceholder; 
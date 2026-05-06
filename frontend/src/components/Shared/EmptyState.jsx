/**
 * PRODUCTION EMPTY STATE SYSTEM
 * Unified empty states with improved UX and clear CTAs
 */
import React, { memo } from 'react';
import { Search, Plus, Filter, Inbox, Users, Package } from 'lucide-react';
import styles from './EmptyState.module.css';
const EmptyState = memo(({
  title = 'No items found',
  message = 'There are no items to display.',
  icon = 'inbox',
  variant = 'default',
  size = 'medium',
  primaryAction = null,
  secondaryAction = null,
  hasFilters = false,
  onClearFilters = null,
  className = '',
  ...props
}) => {
  const getIcon = () => {
    const iconProps = {
      size: size === 'small' ? 32 : size === 'large' ? 48 : 40,
      className: styles.emptyIcon
    };
    switch (icon) {
      case 'search':
        return <Search {...iconProps} />;
      case 'users':
        return <Users {...iconProps} />;
      case 'package':
        return <Package {...iconProps} />;
      case 'filter':
        return <Filter {...iconProps} />;
      case 'plus':
        return <Plus {...iconProps} />;
      default:
        return <Inbox {...iconProps} />;
    }
  };
  const containerClass = `
    ${styles.container} 
    ${styles[size]} 
    ${styles[variant]} 
    ${className}
  `.trim();
  return (
    <div className={containerClass} {...props}>
      <div className={styles.iconContainer}>
        {getIcon()}
      </div>
      <div className={styles.content}>
        <h3 className={styles.title}>{title}</h3>
        <p className={styles.message}>{message}</p>
        {hasFilters && (
          <p className={styles.suggestion}>
            Try adjusting your search criteria or clearing filters to see more results.
          </p>
        )}
      </div>
      <div className={styles.actions}>
        {hasFilters && onClearFilters && (
          <button
            className={`${styles.button} ${styles.secondary}`}
            onClick={onClearFilters}
            type="button"
          >
            <Filter className={styles.buttonIcon} size={16} />
            Clear Filters
          </button>
        )}
        {primaryAction && (
          <button
            className={`${styles.button} ${styles.primary}`}
            onClick={primaryAction.onClick}
            type="button"
          >
            {primaryAction.icon && (
              <primaryAction.icon className={styles.buttonIcon} size={16} />
            )}
            {primaryAction.label}
          </button>
        )}
        {secondaryAction && (
          <button
            className={`${styles.button} ${styles.secondary}`}
            onClick={secondaryAction.onClick}
            type="button"
          >
            {secondaryAction.icon && (
              <secondaryAction.icon className={styles.buttonIcon} size={16} />
            )}
            {secondaryAction.label}
          </button>
        )}
      </div>
    </div>
  );
});
EmptyState.displayName = 'EmptyState';
export default EmptyState; 
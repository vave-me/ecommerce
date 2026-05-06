/**
 * Unified Empty State Component
 * Consolidates all empty state implementations into a single, flexible component
 * 
 * Features:
 * - Multiple icon options
 * - Primary and secondary actions
 * - Filter-specific empty states
 * - Size variants
 * - Clear UX messaging
 */
import React, { memo } from 'react';
import { 
  Search, Plus, Filter, Inbox, Users, Package, 
  FileX, FolderOpen, AlertCircle 
} from 'lucide-react';
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
  // Legacy props for backward compatibility
  label,
  onAction,
  actionLabel,
  ...props
}) => {
  // Handle legacy props
  const displayTitle = title || 'Nothing here yet';
  const displayMessage = message || label || 'No items found';
  
  const getIcon = () => {
    const iconSize = size === 'small' || size === 'sm' ? 32 : 
                     size === 'large' || size === 'lg' ? 48 : 40;
    const iconProps = {
      size: iconSize,
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
      case 'file':
        return <FileX {...iconProps} />;
      case 'folder':
        return <FolderOpen {...iconProps} />;
      case 'alert':
        return <AlertCircle {...iconProps} />;
      case 'inbox':
      default:
        return <Inbox {...iconProps} />;
    }
  };
  
  // Handle legacy action prop
  const primaryActionConfig = primaryAction || (onAction ? {
    onClick: onAction,
    label: actionLabel || 'Add Item',
    icon: Plus
  } : null);
  
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
        <h3 className={styles.title}>{displayTitle}</h3>
        <p className={styles.message}>{displayMessage}</p>
        
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
        
        {primaryActionConfig && (
          <button
            className={`${styles.button} ${styles.primary}`}
            onClick={primaryActionConfig.onClick}
            type="button"
          >
            {primaryActionConfig.icon && (
              <primaryActionConfig.icon className={styles.buttonIcon} size={16} />
            )}
            {primaryActionConfig.label}
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

// Export legacy component name for backward compatibility
export const EmptyPlaceholder = EmptyState;

export default EmptyState;
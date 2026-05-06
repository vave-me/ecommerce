// File: src/components/Support/SupportTicketItem.jsx
"use client";
import React, { useState, useCallback, useMemo, memo } from "react";
import PropTypes from "prop-types";
import { useTranslations } from 'next-intl';
// Icons with Header-inspired import pattern
import { 
    Eye, 
    Lock, 
    Mail, 
    Trash2, 
    MoreHorizontal,
    Clock,
    AlertCircle,
    CheckCircle,
    Pause,
    Play,
    ExternalLink,
    Copy,
    Edit
} from '@/icons';
// Custom hooks (Header pattern)
import { useIsMobile } from "../../hooks/useMobileDetection";
// Security utilities
import { createSafeHtml } from "../../utils/sanitizeHtml";
// Styles
import styles from "./SupportTicketItem.module.css";
// Status configuration (Header-inspired pattern)
const STATUS_CONFIG = {
    open: {
        icon: Mail,
        color: '#10b981',
        label: 'Open',
        actionLabel: 'Close',
        actionIcon: Lock,
        nextStatus: 'closed'
    },
    pending: {
        icon: Pause,
        color: '#f59e0b',
        label: 'Pending',
        actionLabel: 'Reopen',
        actionIcon: Play,
        nextStatus: 'open'
    },
    closed: {
        icon: Lock,
        color: '#6b7280',
        label: 'Closed',
        actionLabel: 'Reopen',
        actionIcon: Mail,
        nextStatus: 'open'
    }
};
// Priority configuration
const PRIORITY_CONFIG = {
    low: { color: '#10b981', label: 'Low' },
    medium: { color: '#f59e0b', label: 'Medium' },
    high: { color: '#ef4444', label: 'High' },
    urgent: { color: '#dc2626', label: 'Urgent' }
};
/**
 * Enhanced SupportTicketItem Component with Header-inspired patterns
 * Features sophisticated interactions, improved accessibility, and modern design
 */
const SupportTicketItem = memo(({ 
    ticket, 
    onUpdate, 
    onDelete, 
    onView,
    isMobile: propIsMobile = false,
    searchQuery = "",
    showActions = true 
}) => {
    const t = useTranslations('SupportTicketItem');
    const isMobile = useIsMobile() || propIsMobile;
    // Destructure ticket data
    const { 
        id, 
        title, 
        description, 
        ticketStatus, 
        priority = 'medium',
        category = 'technical',
        createdAt, 
        updatedAt,
        isRead = true 
    } = ticket;
    // Local state for interactions (Header pattern)
    const [isExpanded, setIsExpanded] = useState(false);
    const [showMenu, setShowMenu] = useState(false);
    const [isProcessing, setIsProcessing] = useState(false);
    // Memoized status and priority configs (Header performance pattern)
    const statusConfig = useMemo(() => STATUS_CONFIG[ticketStatus] || STATUS_CONFIG.open, [ticketStatus]);
    const priorityConfig = useMemo(() => PRIORITY_CONFIG[priority] || PRIORITY_CONFIG.medium, [priority]);
    // Memoized formatted dates (Header pattern)
    const formattedDates = useMemo(() => ({
        created: new Date(createdAt).toLocaleDateString(),
        createdTime: new Date(createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        updated: updatedAt ? new Date(updatedAt).toLocaleDateString() : null,
        updatedTime: updatedAt ? new Date(updatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : null,
        relative: getRelativeTime(createdAt)
    }), [createdAt, updatedAt]);
    // Memoized highlighted content for search (Header pattern)
    const highlightedContent = useMemo(() => {
        if (!searchQuery) return { title, description };
        const highlightText = (text, query) => {
            if (!query.trim()) return text;
            const regex = new RegExp(`(${query.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&')})`, 'gi');
            return text.replace(regex, '<mark>$1</mark>');
        };
        return {
            title: highlightText(title, searchQuery),
            description: highlightText(description, searchQuery)
        };
    }, [title, description, searchQuery]);
    // Enhanced handlers with performance optimizations (Header pattern)
    const handleStatusChange = useCallback(async () => {
        if (isProcessing) return;
        setIsProcessing(true);
        try {
            await onUpdate(id, { ticketStatus: statusConfig.nextStatus });
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
            setIsProcessing(false);
        }
    }, [id, onUpdate, statusConfig.nextStatus, isProcessing]);
    const handleDelete = useCallback(async () => {
        if (isProcessing) return;
        // Enhanced confirmation (Header pattern)
        const confirmed = await showConfirmDialog(
            t('deleteConfirmTitle'),
            t('deleteConfirmMessage', { title: title.slice(0, 50) })
        );
        if (!confirmed) return;
        setIsProcessing(true);
        try {
            await onDelete(id);
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
            setIsProcessing(false);
        }
    }, [id, onDelete, title, t, isProcessing]);
    const handleView = useCallback(() => {
        if (onView) {
            onView(id, ticket);
        } else {
            // Default behavior - toggle expansion
            setIsExpanded(prev => !prev);
        }
    }, [id, ticket, onView]);
    const handleCopyId = useCallback(async () => {
        try {
            await navigator.clipboard.writeText(id);
            // You could show a toast here
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
    }, [id]);
    const toggleMenu = useCallback(() => {
        setShowMenu(prev => !prev);
    }, []);
    const closeMenu = useCallback(() => {
        setShowMenu(false);
    }, []);
    // Keyboard handler for accessibility (Header pattern)
    const handleKeyDown = useCallback((e) => {
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleView();
        }
        if (e.key === 'Escape' && showMenu) {
            closeMenu();
        }
    }, [handleView, showMenu, closeMenu]);
    // Component classes (Header pattern)
    const itemClasses = [
        styles.itemContainer,
        styles[ticketStatus],
        isRead ? styles.read : styles.unread,
        isExpanded ? styles.expanded : '',
        isMobile ? styles.mobile : styles.desktop,
        isProcessing ? styles.processing : ''
    ].filter(Boolean).join(' ');
    const StatusIcon = statusConfig.icon;
    const ActionIcon = statusConfig.actionIcon;
    return (
        <article className={itemClasses} role="article" aria-labelledby={`ticket-${id}-title`}>
            {/* Main Content Area */}
            <div 
                className={styles.mainContent}
                onClick={handleView}
                onKeyDown={handleKeyDown}
                tabIndex={0}
                role="button"
                aria-expanded={isExpanded}
                aria-label={t('viewTicketAriaLabel', { title })}
            >
                {/* Header Row */}
                <header className={styles.itemHeader}>
                    <div className={styles.headerLeft}>
                        {/* Status Indicator */}
                        <div 
                            className={styles.statusIndicator}
                            style={{ backgroundColor: statusConfig.color }}
                            aria-label={t('statusAriaLabel', { status: statusConfig.label })}
                        >
                            <StatusIcon size={16} aria-hidden="true" />
                        </div>
                        {/* Ticket ID */}
                        <span className={styles.ticketId}>
                            {t('ticketIdPrefix')} #{id.slice(-6)}
                        </span>
                        {/* Priority Badge */}
                        <span 
                            className={styles.priorityBadge}
                            style={{ color: priorityConfig.color }}
                        >
                            {t(`priority_${priority}`, priorityConfig.label)}
                        </span>
                        {/* Unread Indicator */}
                        {!isRead && (
                            <span className={styles.unreadBadge} aria-label={t('unreadAriaLabel')}>
                                {t('unreadLabel')}
                            </span>
                        )}
                    </div>
                    <div className={styles.headerRight}>
                        {/* Timestamp */}
                        <time 
                            className={styles.timestamp}
                            dateTime={createdAt}
                            title={`${t('createdLabel')}: ${formattedDates.created} ${formattedDates.createdTime}`}
                        >
                            <Clock size={14} aria-hidden="true" />
                            {formattedDates.relative}
                        </time>
                        {/* Actions Menu */}
                        {showActions && (
                            <div className={styles.menuContainer}>
                                <button
                                    className={styles.menuButton}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        toggleMenu();
                                    }}
                                    aria-expanded={showMenu}
                                    aria-haspopup="menu"
                                    aria-label={t('moreActionsAriaLabel')}
                                >
                                    <MoreHorizontal size={16} />
                                </button>
                                {showMenu && (
                                    <div className={styles.menu} role="menu">
                                        <button
                                            role="menuitem"
                                            className={styles.menuItem}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleView();
                                                closeMenu();
                                            }}
                                        >
                                            <Eye size={14} />
                                            {t('viewAction')}
                                        </button>
                                        <button
                                            role="menuitem"
                                            className={styles.menuItem}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleCopyId();
                                                closeMenu();
                                            }}
                                        >
                                            <Copy size={14} />
                                            {t('copyIdAction')}
                                        </button>
                                        <button
                                            role="menuitem"
                                            className={styles.menuItem}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleStatusChange();
                                                closeMenu();
                                            }}
                                            disabled={isProcessing}
                                        >
                                            <ActionIcon size={14} />
                                            {statusConfig.actionLabel}
                                        </button>
                                        <hr className={styles.menuSeparator} />
                                        <button
                                            role="menuitem"
                                            className={`${styles.menuItem} ${styles.destructive}`}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleDelete();
                                                closeMenu();
                                            }}
                                            disabled={isProcessing}
                                        >
                                            <Trash2 size={14} />
                                            {t('deleteAction')}
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </header>
                {/* Content */}
                <div className={styles.content}>
                    <h3 
                        id={`ticket-${id}-title`}
                        className={styles.title}
                        {...createSafeHtml(highlightedContent.title)}
                    />
                    <p 
                        className={styles.description}
                        {...createSafeHtml(highlightedContent.description)}
                    />
                    {/* Metadata */}
                    <div className={styles.metadata}>
                        <span className={styles.category}>
                            {t(`category_${category}`, category)}
                        </span>
                        {updatedAt && (
                            <span className={styles.lastUpdated}>
                                {t('lastUpdatedLabel')}: {formattedDates.updated}
                            </span>
                        )}
                    </div>
                </div>
            </div>
            {/* Expanded Content */}
            {isExpanded && (
                <div className={styles.expandedContent}>
                    <div className={styles.detailsGrid}>
                        <div className={styles.detailItem}>
                            <strong>{t('statusLabel')}:</strong>
                            <span style={{ color: statusConfig.color }}>
                                {t(`status_${ticketStatus}`, statusConfig.label)}
                            </span>
                        </div>
                        <div className={styles.detailItem}>
                            <strong>{t('priorityLabel')}:</strong>
                            <span style={{ color: priorityConfig.color }}>
                                {t(`priority_${priority}`, priorityConfig.label)}
                            </span>
                        </div>
                        <div className={styles.detailItem}>
                            <strong>{t('categoryLabel')}:</strong>
                            <span>{t(`category_${category}`, category)}</span>
                        </div>
                        <div className={styles.detailItem}>
                            <strong>{t('ticketIdLabel')}:</strong>
                            <span className={styles.fullId}>{id}</span>
                        </div>
                        <div className={styles.detailItem}>
                            <strong>{t('createdLabel')}:</strong>
                            <span>{formattedDates.created} {formattedDates.createdTime}</span>
                        </div>
                        {updatedAt && (
                            <div className={styles.detailItem}>
                                <strong>{t('updatedLabel')}:</strong>
                                <span>{formattedDates.updated} {formattedDates.updatedTime}</span>
                            </div>
                        )}
                    </div>
                    {/* Quick Actions */}
                    <div className={styles.quickActions}>
                        <button
                            className={`${styles.actionButton} ${styles.primary}`}
                            onClick={handleStatusChange}
                            disabled={isProcessing}
                        >
                            <ActionIcon size={16} />
                            {statusConfig.actionLabel}
                        </button>
                        <button
                            className={styles.actionButton}
                            onClick={handleCopyId}
                        >
                            <Copy size={16} />
                            {t('copyIdAction')}
                        </button>
                        <button
                            className={`${styles.actionButton} ${styles.destructive}`}
                            onClick={handleDelete}
                            disabled={isProcessing}
                        >
                            <Trash2 size={16} />
                            {t('deleteAction')}
                        </button>
                    </div>
                </div>
            )}
            {/* Click outside handler */}
            {showMenu && (
                <div 
                    className={styles.overlay}
                    onClick={closeMenu}
                    aria-hidden="true"
                />
            )}
        </article>
    );
});
// Helper function for relative time (Header-inspired utility)
function getRelativeTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffInMs = now - date;
    const diffInHours = diffInMs / (1000 * 60 * 60);
    const diffInDays = diffInHours / 24;
    if (diffInHours < 1) {
        const minutes = Math.floor(diffInMs / (1000 * 60));
        return minutes <= 1 ? 'Just now' : `${minutes}m ago`;
    }
    if (diffInHours < 24) {
        return `${Math.floor(diffInHours)}h ago`;
    }
    if (diffInDays < 7) {
        return `${Math.floor(diffInDays)}d ago`;
    }
    return date.toLocaleDateString();
}
// Enhanced confirmation dialog (Header-inspired utility)
function showConfirmDialog(title, message) {
    return new Promise((resolve) => {
        const confirmed = window.confirm(`${title}\n\n${message}`);
        resolve(confirmed);
    });
}
SupportTicketItem.displayName = 'SupportTicketItem';
SupportTicketItem.propTypes = {
    ticket: PropTypes.shape({
        id: PropTypes.string.isRequired,
        title: PropTypes.string.isRequired,
        description: PropTypes.string.isRequired,
        ticketStatus: PropTypes.oneOf(['open', 'closed', 'pending']).isRequired,
        priority: PropTypes.oneOf(['low', 'medium', 'high', 'urgent']),
        category: PropTypes.string,
        createdAt: PropTypes.string.isRequired,
        updatedAt: PropTypes.string,
        isRead: PropTypes.bool,
    }).isRequired,
    onUpdate: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    onView: PropTypes.func,
    isMobile: PropTypes.bool,
    searchQuery: PropTypes.string,
    showActions: PropTypes.bool
};
export default SupportTicketItem;
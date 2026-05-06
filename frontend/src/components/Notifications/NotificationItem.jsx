/**
 * NotificationItem Component
 * Displays individual notification details with management options.
 */
import { Euro, Eye, Mail, ShoppingCart, Star, Trash2 } from '@/icons';
import { FaCheckCircle, FaTimesCircle, FaUserShield } from '../../utils/iconImports';
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import styles from './NotificationItem.module.css';
/**
 * NotificationItem Component
 * Displays individual notification details with management options.
 */
const NotificationItem = memo(({ notification, onUpdate, onDelete }) => {
    const { id, type, title, message, createdAt, isRead } = notification;
    /**
     * Handler to toggle read status
     */
    const handleToggleRead = async () => {
        try {
            await onUpdate(id, { isRead: !isRead });
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
    };
    /**
     * Handler to delete notification
     */
    const handleDelete = async () => {
        try {
            await onDelete(id);
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
    };
    /**
     * Handler to view notification details
     */
    const handleView = () => {
        // In real scenario, navigate to detailed view or perform an action
        alert(`Viewing details for notification: ${title}`);
    };
    return (
        <div className={`${styles.itemContainer} ${isRead ? styles.read : styles.unread}`}>
            <div className={styles.iconWrapper}>{getTypeIcon(type)}</div>
            <div className={styles.content}>
                <div className={styles.header}>
                    <h4 className={styles.title}>{title}</h4>
                    <span className={styles.timestamp}>{new Date(createdAt).toLocaleString()}</span>
                </div>
                <p className={styles.message}>{message}</p>
                <div className={styles.actions}>
                    <button
                        className={styles.actionButton}
                        onClick={handleToggleRead}
                        aria-label={isRead ? 'Mark as unread' : 'Mark as read'}
                    >
                        {isRead ? <FaTimesCircle /> : <FaCheckCircle />}{' '}
                        {isRead ? 'Mark as Unread' : 'Mark as Read'}
                    </button>
                    <button
                        className={styles.actionButton}
                        onClick={handleView}
                        aria-label="View Notification"
                    >
                        <Eye /> View
                    </button>
                    <button
                        className={`${styles.actionButton} ${styles.deleteButton}`}
                        onClick={handleDelete}
                        aria-label="Delete Notification"
                    >
                        <Trash2 /> Delete
                    </button>
                </div>
            </div>
        </div>
    );
});
NotificationItem.displayName = 'NotificationItem';
NotificationItem.propTypes = {
    notification: PropTypes.shape({
        id: PropTypes.string.isRequired,
        type: PropTypes.string.isRequired,
        title: PropTypes.string.isRequired,
        message: PropTypes.string.isRequired,
        createdAt: PropTypes.string.isRequired,
        isRead: PropTypes.bool.isRequired,
    }).isRequired,
    onUpdate: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default NotificationItem;
/**
 * Helper function to get icons based on notification type
 */
const getTypeIcon = (type) => {
    const iconStyle = { size: 20 };
    switch (type) {
        case 'comment':
            return <FaCheckCircle {...iconStyle} color="#3498db" />;
        case 'message':
            return <Mail {...iconStyle} color="#3498db" />;
        case 'price_drop':
            return <Euro {...iconStyle} color="#3498db" />;
        case 'review':
            return <Star {...iconStyle} color="#3498db" />;
        case 'observation':
            return <ShoppingCart {...iconStyle} color="#3498db" />;
        case 'support':
            return <FaUserShield {...iconStyle} color="#3498db" />;
        default:
            return <FaCheckCircle {...iconStyle} color="#3498db" />;
    }
};
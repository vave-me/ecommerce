import { Bell, Euro, Mail, MessageCircle, ShoppingCart, Star } from '@/icons';
import { FaUserShield } from '../../utils/iconImports';
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import NotificationItem from './NotificationItem';
import styles from './NotificationList.module.css';
const NotificationList = memo(({ notifications, onUpdate, onDelete }) => {
    const notificationTypes = {
        comment: 'Comments',
        message: 'Messages',
        price_drop: 'Price Drops',
        review: 'Reviews',
        observation: 'Observations',
        support: 'Support',
    };
    return (
        <div className={styles.listContainer}>
            {Object.keys(notificationTypes).map((type) => {
                const filtered = notifications.filter((notif) => notif.type === type);
                if (filtered.length === 0) return null;
                return (
                    <div className={styles.categorySection} key={type}>
                        <div className={styles.categoryHeader}>
                            <div className={styles.categoryIcon}>
                                {getCategoryIcon(type)}
                            </div>
                            <h3 className={styles.categoryTitle}>
                                {notificationTypes[type]}
                            </h3>
                        </div>
                        {filtered.map((notif) => (
                            <NotificationItem
                                key={notif.id}
                                notification={notif}
                                onUpdate={onUpdate}
                                onDelete={onDelete}
                            />
                        ))}
                    </div>
                );
            })}
        </div>
    );
});
NotificationList.displayName = 'NotificationList';
NotificationList.propTypes = {
    notifications: PropTypes.arrayOf(PropTypes.object).isRequired,
    onUpdate: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default NotificationList;
/**
 * Helper function to get category icons
 */
const getCategoryIcon = (type) => {
    const iconStyle = { size: 16, color: '#222' };
    switch (type) {
        case 'comment':
            return <MessageCircle {...iconStyle} />;
        case 'message':
            return <Mail {...iconStyle} />;
        case 'price_drop':
            return <Euro {...iconStyle} />;
        case 'review':
            return <Star {...iconStyle} />;
        case 'observation':
            return <ShoppingCart {...iconStyle} />;
        case 'support':
            return <FaUserShield {...iconStyle} />;
        default:
            return <Bell {...iconStyle} />;
    }
};
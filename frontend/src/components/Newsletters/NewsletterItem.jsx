// src/components/NewsletterItem.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Edit, Trash2, Users, Mail } from '@/icons';
import styles from './NewsletterItem.module.css';
/**
 * NewsletterItem Component
 * Displays individual newsletter details with management options.
 */
const NewsletterItem = memo(({ newsletter, currentUser, onEdit, onDelete, onSubscribe, onUnsubscribe, subscription }) => {
    const { id, name, description, subscriber_count, created_at, user_id, is_active, frequency, category } = newsletter;
    const isSubscribed = !!subscription && (subscription.status === 'active' || subscription.status === 'pending');
    
    const handleSubscribeToggle = () => {
        if (isSubscribed && subscription) {
            onUnsubscribe(subscription.id);
        } else {
            onSubscribe(id);
        }
    };
    const formatDate = (date) => {
        return new Date(date).toLocaleDateString();
    };
    const canEdit = currentUser && (currentUser.id === user_id || currentUser.role === 'admin');
    return (
        <div className={styles.itemContainer}>
            <div className={styles.iconWrapper}>
                <Mail size={20} />
            </div>
            <div className={styles.content}>
                <div className={styles.header}>
                    <h3 className={styles.title}>{name}</h3>
                    <div className={styles.meta}>
                        <span className={styles.subscriberCount}>
                            <Users size={14} />
                            {subscriber_count || 0} subscribers
                        </span>
                        <span className={styles.frequency}>
                            {frequency}
                        </span>
                        {category && (
                            <span className={styles.category}>
                                {category}
                            </span>
                        )}
                        <span className={styles.date}>
                            {formatDate(created_at)}
                        </span>
                    </div>
                </div>
                {description && (
                    <p className={styles.description}>{description}</p>
                )}
                {!is_active && (
                    <div className={styles.inactive}>
                        (Inactive)
                    </div>
                )}
            </div>
            <div className={styles.actions}>
                <button
                    className={`${styles.subscribeButton} ${isSubscribed ? styles.subscribed : ''}`}
                    onClick={handleSubscribeToggle}
                    aria-label={isSubscribed ? 'Unsubscribe' : 'Subscribe'}
                >
                    {isSubscribed ? 'Unsubscribe' : 'Subscribe'}
                </button>
                {canEdit && (
                    <>
                        <button
                            className={styles.actionButton}
                            onClick={() => onEdit(newsletter)}
                            aria-label="Edit newsletter"
                        >
                            <Edit size={16} />
                        </button>
                        <button
                            className={`${styles.actionButton} ${styles.deleteButton}`}
                            onClick={() => onDelete(id)}
                            aria-label="Delete newsletter"
                        >
                            <Trash2 size={16} />
                        </button>
                    </>
                )}
            </div>
        </div>
    );
});
NewsletterItem.displayName = 'NewsletterItem';
NewsletterItem.propTypes = {
    newsletter: PropTypes.shape({
        id: PropTypes.string.isRequired,
        name: PropTypes.string.isRequired,
        description: PropTypes.string,
        subscriber_count: PropTypes.number,
        is_active: PropTypes.bool,
        created_at: PropTypes.string.isRequired,
        user_id: PropTypes.string.isRequired,
        frequency: PropTypes.string,
        category: PropTypes.string,
    }).isRequired,
    subscription: PropTypes.shape({
        id: PropTypes.string.isRequired,
        status: PropTypes.string.isRequired,
    }),
    currentUser: PropTypes.object,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    onSubscribe: PropTypes.func.isRequired,
    onUnsubscribe: PropTypes.func.isRequired,
};
export default NewsletterItem;
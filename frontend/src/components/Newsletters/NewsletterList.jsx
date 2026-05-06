// src/components/NewsletterList.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import NewsletterItem from './NewsletterItem';
import styles from './NewsletterList.module.css';
/**
 * NewsletterList Component
 * Displays a list of newsletters with subscription options.
 */
const NewsletterList = memo(({ newsletters, subscriptions, currentUser, onEdit, onDelete, onSubscribe, onUnsubscribe }) => {
    const getSubscriptionForNewsletter = (newsletterId) => {
        return subscriptions?.find(sub => sub.newsletter_id === newsletterId);
    };
    
    return (
        <ul className={styles.listContainer}>
            {newsletters.map((newsletter) => (
                <NewsletterItem
                    key={newsletter.id}
                    newsletter={newsletter}
                    subscription={getSubscriptionForNewsletter(newsletter.id)}
                    currentUser={currentUser}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    onSubscribe={onSubscribe}
                    onUnsubscribe={onUnsubscribe}
                />
            ))}
        </ul>
    );
});
NewsletterList.displayName = 'NewsletterList';
NewsletterList.propTypes = {
    newsletters: PropTypes.arrayOf(PropTypes.object).isRequired,
    subscriptions: PropTypes.arrayOf(PropTypes.object),
    currentUser: PropTypes.object, // Assuming user object from context
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    onSubscribe: PropTypes.func.isRequired,
    onUnsubscribe: PropTypes.func.isRequired,
};
export default NewsletterList;
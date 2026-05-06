"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import ActivityItem from './ActivityItem';
import styles from './ActivityList.module.css';
/**
 * ActivityList - Renders a list of activity items
 *
 * @param {Object} user - The current user object
 * @param {Array} activities - Array of activity objects to display
 * @param {Function} onToggleRead - Handler for toggling read/unread status
 * @param {Function} onDelete - Handler for deleting an activity
 */
const ActivityList = memo(({ user, activities, onToggleRead, onDelete }) => {
    return (
        <div className={styles.listContainer}>
            {activities.map((activity) => (
                <ActivityItem
                    user={user}
                    key={activity.id}
                    activity={activity}
                    onToggleRead={onToggleRead}
                    onDelete={onDelete}
                />
            ))}
        </div>
    );
});
ActivityList.displayName = 'ActivityList';
ActivityList.propTypes = {
    user: PropTypes.object,
    activities: PropTypes.arrayOf(PropTypes.object).isRequired,
    onToggleRead: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default ActivityList;
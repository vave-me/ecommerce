"use client";
import React, { useState, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import {Heart, ThumbsUp, Trash2} from '@/icons';
import {FaAngry, FaLaugh, FaSadTear} from '../../utils/iconImports';
import ReactionButtons from "./ReactionButtons";
import {addInteraction, getInteractions} from "../../api/client/activityApi";
import styles from './ActivityItem.module.css'; // <-- import CSS module
import {useTranslations} from 'next-intl'; // <-- import translations hook
const ActivityItem = memo(({user, activity, onToggleRead, onDelete}) => {
    const {id, target, message, createdAt, isRead} = activity;
    const [currentReactions, setCurrentReactions] = useState({});
    const t = useTranslations('ActivityItem'); // <-- translation hook
    useEffect(() => {
        const fetchReactions = async () => {
            try {
                const data = await getInteractions(id);
                const reactions = data.interactions.reduce((acc, interaction) => {
                    acc[interaction.actionType] = (acc[interaction.actionType] || 0) + 1;
                    return acc;
                }, {});
                setCurrentReactions(reactions);
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
        };
        fetchReactions();
    }, [id]);
    const handleReact = async (reactionType) => {
        try {
            await addInteraction(
                id,
                target?.itemId,
                target?.itemType,
                reactionType
            );
            setCurrentReactions((prev) => ({
                ...prev,
                [reactionType]: (prev[reactionType] || 0) + 1,
            }));
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    };
    const handleToggleRead = () => {
        onToggleRead(id, !isRead);
    };
    const handleDelete = () => {
        onDelete(id);
    };
    return (
        <div
            className={`
              ${styles.itemContainer}
              ${isRead ? styles.isRead : ''}
            `}
        >
            {/* Avatar */}
            <img
                className={styles.avatar}
                src="/images/user-user.webp"
                alt={`${user.username}'s avatar`}
            />
            {/* Content */}
            <div className={styles.content}>
                <div className={styles.header}>
                    <span className={styles.userName}>{user.username}</span>
                    <span className={styles.message}>
                        {message} {activity.actionType} {activity.itemType} {activity.itemId}
                    </span>
                    <span className={styles.timestamp}>
                        {new Date(createdAt).toLocaleString()}
                    </span>
                </div>
                <div className={styles.targetItem}>
                    {target?.itemName}
                </div>
                <div className={styles.reactionsAndActions}>
                    <div className={styles.reactionSection}>
                        <ReactionButtons onReact={handleReact}/>
                        <div className={styles.reactionCounts}>
                            {Object.entries(currentReactions).map(([reaction, count]) => (
                                <span key={reaction} className={styles.reactionCount}>
                                    {getReactionIcon(reaction)} {count}
                                </span>
                            ))}
                        </div>
                    </div>
                    <div className={styles.actions}>
                        <button
                            className={`${styles.baseButton} ${styles.actionButton}`}
                            onClick={handleToggleRead}
                            aria-label={isRead ? t('aria.markAsUnread') : t('aria.markAsRead')}
                        >
                            {isRead ? t('markAsUnread') : t('markAsRead')}
                        </button>
                        <button
                            className={`${styles.baseButton} ${styles.deleteButton}`}
                            onClick={handleDelete}
                            aria-label={t('aria.deleteActivity')}
                        >
                            <Trash2/>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
});
ActivityItem.displayName = 'ActivityItem';
ActivityItem.propTypes = {
    user: PropTypes.shape({
        username: PropTypes.string.isRequired,
        avatar: PropTypes.string,
    }).isRequired,
    activity: PropTypes.shape({
        id: PropTypes.string.isRequired,
        target: PropTypes.shape({
            itemName: PropTypes.string,
            itemType: PropTypes.string,
            itemId: PropTypes.string,
        }),
        message: PropTypes.string.isRequired,
        createdAt: PropTypes.string.isRequired,
        isRead: PropTypes.bool.isRequired,
        actionType: PropTypes.string,
        itemType: PropTypes.string,
        itemId: PropTypes.string,
    }).isRequired,
    onToggleRead: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default ActivityItem;
/** Helper for reaction icons */
function getReactionIcon(reaction) {
    switch (reaction) {
        case 'like':
            return <ThumbsUp color="#2ecc71"/>;
        case 'love':
            return <Heart color="#e91e63"/>;
        case 'haha':
            return <FaLaugh color="#f1c40f"/>;
        case 'sad':
            return <FaSadTear color="#3498db"/>;
        case 'angry':
            return <FaAngry color="#e74c3c"/>;
        default:
            return <ThumbsUp color="#2ecc71"/>;
    }
}

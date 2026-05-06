"use client";
import React from 'react';
import {
    Heart, MessageCircle, ThumbsUp, ThumbsDown, Share2, Mail
} from '@/icons';
import styles from './Engagement.module.css';
// Minimalist EngagementBar Component
const Engagement = ({
                        favorite = false,
                        onFavorite,
                        liked = false,
                        onLike,
                        disliked = false,
                        onDislike,
                        data = {},
                        hideButtons = [],
                        onCommentClick,
                        onMessageClick,
                        onShareClick,
                        variant = 'default' // 'product' for product listing style
                    }) => {
    // Determine which buttons to show
    const showHeart = !hideButtons.includes('heart');
    const showMessageCircle = !hideButtons.includes('messageCircle');
    const showMessage = !hideButtons.includes('mail');
    const showThumbsUp = !hideButtons.includes('thumbsUp');
    const showThumbsDown = !hideButtons.includes('thumbsDown');
    const showShare = !hideButtons.includes('share');
    // Helper function to safely get counts from different data structures
    const getCount = (metricsKey, directKey, defaultValue = 0) => {
        if (data?.metrics && data.metrics[metricsKey] !== undefined) {
            return data.metrics[metricsKey];
        }
        if (data && data[directKey] !== undefined) {
            return data[directKey];
        }
        return defaultValue;
    };
    // Format counters - only show badge if count > 0
    const formatCounter = (count) => {
        return count > 0 ? count : '';
    };
    // Additional class for product variant
    const barClass = variant === 'product' ?
        `${styles.engagementBar} ${styles.productEngagementBar}` :
        styles.engagementBar;
    const buttonClass = variant === 'product' ?
        `${styles.engagementButton} ${styles.productEngagementButton}` :
        styles.engagementButton;
    // Handle click events with preventDefault to avoid navigation
    const handleButtonClick = (handler) => (e) => {
        e.preventDefault();
        if (typeof handler === 'function') {
            handler(e);
        }
    };
    return (
        <div className={barClass}>
            <div className={styles.actions}>
                {showHeart && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onFavorite)}
                        aria-label={favorite ? "Remove from favorites" : "Add to favorites"}
                    >
                        <span
                            className={`${styles.engagementIcon} ${favorite ? `${styles.heartIcon} ${styles.active}` : styles.heartIcon}`}>
                            <Heart size={20} className={favorite ? styles.filledIcon : ''}/>
                        </span>
                        <span className={styles.countBadge}>{formatCounter(getCount('saved', 'savedCount'))}</span>
                    </button>
                )}
                {showMessageCircle && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onCommentClick)}
                        aria-label="Comments"
                    >
                        <span className={styles.engagementIcon}>
                            <MessageCircle size={20}/>
                        </span>
                        <span
                            className={styles.countBadge}>{formatCounter(getCount('comments', 'commentsCount'))}</span>
                    </button>
                )}
                {showMessage && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onMessageClick)}
                        aria-label="Message"
                    >
                        <span className={styles.engagementIcon}>
                            <Mail size={20}/>
                        </span>
                        <span
                            className={styles.countBadge}>{formatCounter(getCount('messages', 'messagesCount'))}</span>
                    </button>
                )}
                {showThumbsUp && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onLike)}
                        aria-label={liked ? "Unlike" : "Like"}
                    >
                        <span
                            className={`${styles.engagementIcon} ${liked ? `${styles.likeIcon} ${styles.active}` : styles.likeIcon}`}>
                            <ThumbsUp size={20} className={liked ? styles.filledIcon : ''}/>
                        </span>
                        <span className={styles.countBadge}>{formatCounter(getCount('likes', 'likesCount'))}</span>
                    </button>
                )}
                {showThumbsDown && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onDislike)}
                        aria-label={disliked ? "Remove dislike" : "Dislike"}
                    >
                        <span
                            className={`${styles.engagementIcon} ${disliked ? `${styles.dislikeIcon} ${styles.active}` : styles.dislikeIcon}`}>
                            <ThumbsDown size={20} className={disliked ? styles.filledIcon : ''}/>
                        </span>
                        <span
                            className={styles.countBadge}>{formatCounter(getCount('dislikes', 'dislikesCount'))}</span>
                    </button>
                )}
                {showShare && (
                    <button
                        className={buttonClass}
                        onClick={handleButtonClick(onShareClick)}
                        aria-label="Share"
                    >
                        <span className={styles.engagementIcon}>
                            <Share2 size={20}/>
                        </span>
                        <span className={styles.countBadge}>{formatCounter(getCount('shares', 'sharesCount'))}</span>
                    </button>
                )}
            </div>
        </div>
    );
};
export {Engagement};
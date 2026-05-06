"use client";
import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';
import styles from './TopicButton.module.css';
/**
 * TopicButton Component
 * Individual topic button with label, badge, and chevron
 * Extracted from SelectTopic to improve component modularity
 */
const TopicButton = forwardRef(({ 
    topic, 
    isActive, 
    isOpen, 
    onClick, 
    onKeyDown,
    dropdownId 
}, ref) => {
    return (
        <div className={styles.topicItem}>
            <button
                ref={ref}
                className={`${styles.topicButton} ${isActive ? styles.topicActive : ''}`}
                onClick={onClick}
                onKeyDown={onKeyDown}
                aria-haspopup="dialog"
                aria-expanded={isOpen}
                aria-controls={isOpen ? dropdownId : undefined}
                type="button"
            >
                <span className={styles.topicLabel}>{topic.label}</span>
                {topic.badge && (
                    <span className={styles.topicBadge}>{topic.badge}</span>
                )}
                <span
                    className={`${styles.topicChevron} ${isOpen ? styles.topicChevronOpen : ''}`}
                    aria-hidden="true"
                >
                    <svg 
                        xmlns="http://www.w3.org/2000/svg" 
                        width="10" 
                        height="10" 
                        viewBox="0 0 24 24"
                        fill="none" 
                        stroke="currentColor" 
                        strokeWidth="3" 
                        strokeLinecap="round"
                        strokeLinejoin="round"
                    >
                        <polyline points="6 9 12 15 18 9"></polyline>
                    </svg>
                </span>
            </button>
        </div>
    );
});
TopicButton.displayName = 'TopicButton';
TopicButton.propTypes = {
    topic: PropTypes.shape({
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        badge: PropTypes.string,
    }).isRequired,
    isActive: PropTypes.bool.isRequired,
    isOpen: PropTypes.bool.isRequired,
    onClick: PropTypes.func.isRequired,
    onKeyDown: PropTypes.func.isRequired,
    dropdownId: PropTypes.string.isRequired,
};
export default TopicButton; 
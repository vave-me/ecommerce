// File: src/components/Messaging/ChatHeader.jsx
"use client";
import React, { memo } from 'react';
import { X } from '@/icons';
import PropTypes from 'prop-types';
import styles from './ChatHeader.module.css';
const ChatHeader = memo(function ChatHeader({
                                       recipient,
                                       onCloseChat,
                                   }) {
    // Handle null recipient gracefully
    if (!recipient) {
        return (
            <div className={styles.header}>
                <div className={styles.recipientInfo}>
                    <h3 className={styles.recipientName}>Loading...</h3>
                    <span className={styles.recipientStatus}>Loading user data...</span>
                </div>
                <button
                    className={styles.iconButton}
                    onClick={onCloseChat}
                    aria-label="Close Chat"
                >
                    <X />
                </button>
            </div>
        );
    }
    const { name, avatar, online } = recipient;
    return (
        <div className={styles.header}>
            {/* Avatar + Status */}
            <div className={styles.avatarContainer}>
                <img
                    className={styles.avatar}
                    src={avatar || '/images/user-user.webp'}
                    alt={`${name}'s avatar`}
                />
                <div
                    className={`
                      ${styles.status} 
                      ${online ? styles.online : styles.offline}
                    `}
                    aria-label={online ? 'Online' : 'Offline'}
                />
            </div>
            {/* Recipient Info */}
            <div className={styles.recipientInfo}>
                <h3 className={styles.recipientName}>{name}</h3>
                <span className={styles.recipientStatus}>
                    {online ? 'Active now' : 'Offline'}
                </span>
            </div>
            {/* Close Button - Returns to conversation list */}
            <button
                className={styles.iconButton}
                onClick={onCloseChat}
                aria-label="Close Chat"
            >
                <X />
            </button>
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison focusing on recipient data changes
    // Skip callback comparisons as they should be stable via useCallback
    const prevRecipient = prevProps.recipient;
    const nextRecipient = nextProps.recipient;
    if (!prevRecipient || !nextRecipient) {
        return prevRecipient === nextRecipient;
    }
    return prevRecipient.id === nextRecipient.id &&
           prevRecipient.name === nextRecipient.name &&
           prevRecipient.avatar === nextRecipient.avatar &&
           prevRecipient.online === nextRecipient.online;
});
ChatHeader.propTypes = {
    recipient: PropTypes.shape({
        id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        name: PropTypes.string.isRequired,
        avatar: PropTypes.string,
        online: PropTypes.bool.isRequired,
    }), // Removed .isRequired to allow null
    onCloseChat: PropTypes.func.isRequired,
};
export default ChatHeader;

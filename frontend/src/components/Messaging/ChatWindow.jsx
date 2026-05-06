// File: src/components/Messaging/ChatWindow.jsx
"use client";
import React, { useRef, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import ChatHeader from './ChatHeader';
import MessageInput from './MessageInput';
import useChatHistory from "../../hooks/useChatHistory";
import styles from './ChatWindow.module.css';
const ChatWindow = memo(function ChatWindow({
    conversationId,
    recipient,
    itemId,
    recipientId,
    onCloseChat,
    onClose,
    onPrevious,
    isConnected,
}) {
    // useChatHistory custom hook
    const { messages, isLoading, error, sendMessage } = useChatHistory(
        conversationId,
        {
            recipientId: recipientId,
            itemId: itemId,
        }
    );
    // For auto-scrolling
    const messagesEndRef = useRef(null);
    // Auto-scroll whenever `messages` changes
    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [messages]);
    if (!conversationId) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyChatState}>
                    <h2 className={styles.emptyChatTitle}>No conversation selected</h2>
                    <p className={styles.emptyChatText}>Select a conversation to start messaging</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            <ChatHeader
                recipient={recipient}
                onCloseChat={onCloseChat}
            />
            <div className={styles.messagesWrapper}>
                {isLoading && (
                    <div className={styles.messageLoading}>
                        <div className={styles.loadingSpinner} />
                        <span>Loading messages...</span>
                    </div>
                )}
                {error && (
                    <div className={styles.messageError}>
                        <div className={styles.errorIcon}>⚠️</div>
                        <p>Error loading messages: {error.message || error}</p>
                        <button 
                            className={styles.retryButton}
                            onClick={() => window.location.reload()}
                        >
                            Try Again
                        </button>
                    </div>
                )}
                {!isLoading && !error && messages.length === 0 && (
                    <div className={styles.emptyChatState}>
                        <h2 className={styles.emptyChatTitle}>Start the conversation</h2>
                        <p className={styles.emptyChatText}>Send a message to {recipient?.name || 'this user'} to get started</p>
                    </div>
                )}
                {!isLoading && !error && messages.length > 0 && (
                    messages.map((msg) => {
                        const isUserMsg = msg.isUserMessage;
                        // For the container: .messageItem + user/other variant
                        const messageItemClass = isUserMsg
                            ? `${styles.messageItem} ${styles.messageItemUser}`
                            : `${styles.messageItem} ${styles.messageItemOther}`;
                        // For the bubble: .bubble + user/other variant
                        const bubbleClass = isUserMsg
                            ? `${styles.bubble} ${styles.bubbleUser}`
                            : `${styles.bubble} ${styles.bubbleOther}`;
                        return (
                            <div key={msg.id} className={messageItemClass}>
                                <div className={bubbleClass}>
                                    <p className={styles.text}>{msg.text}</p>
                                    {msg.time && (
                                        <div className={styles.messageTime}>
                                            {msg.time}
                                        </div>
                                    )}
                                </div>
                            </div>
                        );
                    })
                )}
                <div ref={messagesEndRef} />
            </div>
            <MessageInput
                onSendMessage={sendMessage}
                disabled={!Boolean(isConnected)}
            />
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for chat window optimization
    return prevProps.conversationId === nextProps.conversationId &&
           prevProps.recipientId === nextProps.recipientId &&
           prevProps.itemId === nextProps.itemId &&
           prevProps.isConnected === nextProps.isConnected &&
           prevProps.recipient?.id === nextProps.recipient?.id &&
           prevProps.recipient?.name === nextProps.recipient?.name &&
           prevProps.recipient?.online === nextProps.recipient?.online;
});
ChatWindow.propTypes = {
    conversationId: PropTypes.string,
    recipient: PropTypes.shape({
        id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        name: PropTypes.string.isRequired,
        avatar: PropTypes.string,
        online: PropTypes.bool,
    }),
    itemId: PropTypes.string,
    recipientId: PropTypes.string,
    onCloseChat: PropTypes.func.isRequired,
    onClose: PropTypes.func.isRequired,
    onPrevious: PropTypes.func.isRequired,
    isConnected: PropTypes.bool.isRequired,
};
ChatWindow.defaultProps = {
    conversationId: null,
    recipient: null,
    itemId: '',
    recipientId: '',
};
export default ChatWindow;

// File: src/components/Messaging/MessageList.jsx
"use client";
import React, { useState, useEffect, memo, useMemo, useCallback } from 'react';
import PropTypes from 'prop-types';
import { Search, X } from '@/icons';
import styles from './MessageList.module.css';
const MessageList = memo(function MessageList({
    conversations,
    onSelectConversation,
    selectedConversationId,
    onClose,
    loading = false,
    error = null,
}) {
    const [searchTerm, setSearchTerm] = useState('');
    // Enhanced filtering with multiple criteria
    const filteredConversations = useMemo(() => {
        if (!searchTerm.trim()) return conversations;
        const searchLower = searchTerm.toLowerCase();
        return conversations.filter((conversation) => {
            const userName = conversation.user?.name?.toLowerCase() || '';
            const lastMessage = conversation.lastMessage?.toLowerCase() || '';
            return userName.includes(searchLower) || lastMessage.includes(searchLower);
        });
    }, [conversations, searchTerm]);
    // Format relative time with better error handling
    const formatTime = useCallback((timeString) => {
        if (!timeString || timeString === 'N/A') return 'Recently';
        try {
            const date = new Date(timeString);
            if (isNaN(date.getTime())) return 'Recently';
            const now = new Date();
            const diffMs = now - date;
            const diffMins = Math.floor(diffMs / (1000 * 60));
            const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
            const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
            if (diffMins < 1) return 'Just now';
            if (diffMins < 60) return `${diffMins}m ago`;
            if (diffHours < 24) return `${diffHours}h ago`;
            if (diffDays < 7) return `${diffDays}d ago`;
            return date.toLocaleDateString([], { 
                month: 'short', 
                day: 'numeric' 
            });
        } catch {
            return 'Recently';
        }
    }, []);
    // Enhanced keyboard navigation
    const handleKeyDown = useCallback((event, conversationId) => {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            onSelectConversation(conversationId);
        }
    }, [onSelectConversation]);
    // Clear search with enhanced UX
    const handleClearSearch = useCallback(() => {
        setSearchTerm('');
    }, []);
    // Handle search input change
    const handleSearchChange = useCallback((e) => {
        setSearchTerm(e.target.value);
    }, []);
    // Error state
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.listHeader}>
                    <h2 className={styles.headerTitle}>Messages</h2>
                </div>
                <div className={styles.errorState}>
                    <div className={styles.errorIcon}>⚠️</div>
                    <h3 className={styles.errorTitle}>Failed to load conversations</h3>
                    <p className={styles.errorText}>{error}</p>
                    <button 
                        className={styles.retryButton}
                        onClick={() => window.location.reload()}
                        aria-label="Retry loading conversations"
                    >
                        Try Again
                    </button>
                </div>
            </div>
        );
    }
    // Loading state
    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.listHeader}>
                    <h2 className={styles.headerTitle}>Messages</h2>
                    <div className={styles.searchContainer}>
                        <input
                            type="text"
                            className={styles.searchBar}
                            placeholder="Search conversations..."
                            value={searchTerm}
                            onChange={handleSearchChange}
                            disabled
                            aria-label="Search Conversations"
                        />
                        <Search size={16} className={styles.searchIcon} />
                    </div>
                </div>
                <div className={styles.loadingState}>
                    <div className={styles.loadingSpinner} />
                    <span className={styles.loadingText}>Loading conversations...</span>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container} role="complementary" aria-label="Conversations">
            {/* Enhanced Header */}
            <div className={styles.listHeader}>
                <div className={styles.headerTop}>
                    <h2 className={styles.headerTitle}>Messages</h2>
                    <button
                        className={styles.closeButton}
                        onClick={onClose}
                        aria-label="Close messages"
                        type="button"
                    >
                        <X size={20} />
                    </button>
                </div>
                {/* Enhanced Search */}
                <div className={styles.searchContainer}>
                    <input
                        type="text"
                        className={styles.searchBar}
                        placeholder="Search conversations..."
                        value={searchTerm}
                        onChange={handleSearchChange}
                        aria-label="Search Conversations"
                        aria-describedby="search-help"
                    />
                    <Search size={16} className={styles.searchIcon} aria-hidden="true" />
                    {searchTerm && (
                        <button
                            className={styles.clearSearchButton}
                            onClick={handleClearSearch}
                            aria-label="Clear search"
                            type="button"
                        >
                            <X size={14} />
                        </button>
                    )}
                </div>
                <div id="search-help" className={styles.searchHelp}>
                    Search by name or message content
                </div>
            </div>
            {/* Conversations List */}
            <div 
                className={styles.conversations}
                role="list"
                aria-label={`${filteredConversations.length} conversations`}
            >
                {filteredConversations.length > 0 ? (
                    filteredConversations.map((conversation) => {
                        const isSelected = conversation.id === selectedConversationId;
                        const conversationClasses = [
                            styles.conversation,
                            isSelected ? styles.conversationSelected : ''
                        ].filter(Boolean).join(' ');
                        const user = conversation.user || {};
                        const formattedTime = formatTime(conversation.time);
                        const isOnline = user.online || false;
                        return (
                            <div
                                key={conversation.id}
                                className={conversationClasses}
                                onClick={() => onSelectConversation(conversation.id)}
                                onKeyDown={(e) => handleKeyDown(e, conversation.id)}
                                role="listitem button"
                                tabIndex="0"
                                aria-pressed={isSelected}
                                aria-label={`Conversation with ${user.name || 'Unknown User'}. Last message: ${conversation.lastMessage || 'No recent message'}. ${formattedTime}.`}
                                aria-describedby={`conversation-${conversation.id}-details`}
                            >
                                {/* Modern Conversation Header */}
                                <div className={styles.conversationHeader}>
                                    <div className={styles.avatarContainer}>
                                        <img
                                            src={user.avatar || '/images/user-user.webp'}
                                            alt=""
                                            className={styles.userAvatar}
                                            onError={(e) => {
                                                e.target.src = '/images/user-user.webp';
                                            }}
                                            loading="lazy"
                                        />
                                        <div 
                                            className={`${styles.statusDot} ${isOnline ? styles.online : ''}`}
                                            aria-label={isOnline ? 'Online' : 'Offline'}
                                        />
                                    </div>
                                    <div className={styles.userInfo}>
                                        <h3 className={styles.userName}>
                                            {user.name || 'Unknown User'}
                                        </h3>
                                        <div className={styles.userStatus}>
                                            <span>{isOnline ? 'Active now' : 'Offline'}</span>
                                        </div>
                                    </div>
                                    <div className={styles.conversationMeta}>
                                        <span className={styles.conversationTime}>
                                            {formattedTime}
                                        </span>
                                        {conversation.unreadCount > 0 && (
                                            <span 
                                                className={styles.unreadBadge}
                                                aria-label={`${conversation.unreadCount} unread messages`}
                                            >
                                                {conversation.unreadCount > 9 ? '9+' : conversation.unreadCount}
                                            </span>
                                        )}
                                    </div>
                                </div>
                                {/* Message Preview */}
                                <div 
                                    className={styles.messagePreview}
                                    id={`conversation-${conversation.id}-details`}
                                >
                                    <span className={styles.messageText}>
                                        {conversation.lastMessage || 'No recent message'}
                                    </span>
                                </div>
                            </div>
                        );
                    })
                ) : searchTerm ? (
                    // No search results
                    <div className={styles.noSearchResults} role="status">
                        <div className={styles.noResultsIcon}>🔍</div>
                        <h3 className={styles.noResultsTitle}>No conversations found</h3>
                        <p className={styles.noResultsText}>
                            Try adjusting your search terms or check the spelling.
                        </p>
                        <button
                            className={styles.clearSearchButton}
                            onClick={handleClearSearch}
                            type="button"
                        >
                            Clear search
                        </button>
                    </div>
                ) : (
                    // No conversations at all
                    <div className={styles.noConversations} role="status">
                        <div className={styles.noConversationsIcon}>💬</div>
                        <h3 className={styles.noConversationsTitle}>No conversations yet</h3>
                        <p className={styles.noConversationsText}>
                            Start a conversation by messaging someone from their profile or listing.
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}, (prevProps, nextProps) => {
    // Optimized comparison focusing on relevant changes
    if (prevProps.loading !== nextProps.loading) return false;
    if (prevProps.error !== nextProps.error) return false;
    if (prevProps.selectedConversationId !== nextProps.selectedConversationId) return false;
    if (prevProps.conversations.length !== nextProps.conversations.length) return false;
    // Deep comparison for conversation changes
    return prevProps.conversations.every((conv, index) => {
        const nextConv = nextProps.conversations[index];
        return nextConv &&
               conv.id === nextConv.id &&
               conv.lastMessage === nextConv.lastMessage &&
               conv.time === nextConv.time &&
               conv.user?.name === nextConv.user?.name &&
               conv.user?.online === nextConv.user?.online &&
               conv.unreadCount === nextConv.unreadCount;
    });
});
MessageList.propTypes = {
    conversations: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            user: PropTypes.shape({
                id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
                name: PropTypes.string.isRequired,
                avatar: PropTypes.string,
                online: PropTypes.bool,
            }).isRequired,
            lastMessage: PropTypes.string,
            time: PropTypes.string,
            unreadCount: PropTypes.number,
        })
    ).isRequired,
    onSelectConversation: PropTypes.func.isRequired,
    selectedConversationId: PropTypes.string,
    onClose: PropTypes.func.isRequired,
    loading: PropTypes.bool,
    error: PropTypes.string,
};
MessageList.defaultProps = {
    selectedConversationId: null,
    loading: false,
    error: null,
};
export default MessageList;

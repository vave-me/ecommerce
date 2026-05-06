"use client";
import React, {useState, useEffect, useCallback, memo} from 'react';
import {useTranslations} from 'next-intl';
import {
    MessageSquare,
    ChevronDown,
    Plus,
    Check,
    Clock
} from '@/icons';
import { getAssistantService } from '../../services/ai/AssistantService';
import {useAuth} from '../../context/AuthContext';
import styles from './ConversationSelector.module.css';
/**
 * Compact Conversation Selector - Modern ChatGPT style
 */
const ConversationSelector = memo(({
                                       selectedConversationId = null,
                                       onConversationSelect = () => {
                                       },
                                       className = ''
                                   }) => {
    const t = useTranslations('AI');
    const {user} = useAuth();
    const [conversations, setConversations] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const [error, setError] = useState(null);
    /**
     * Load conversations
     */
    const loadConversations = useCallback(async () => {
        if (!user) return;
        setIsLoading(true);
        setError(null);
        try {
            const assistantService = getAssistantService();
            const response = await assistantService.getUserConversations({
                activeOnly: false,
                page: 1,
                limit: 20
            });
            if (response.success) {
                const conversationsData = response.data?.conversations || [];
                setConversations(Array.isArray(conversationsData) ? conversationsData : []);
            } else {
                setError(response.error || 'Failed to load conversations');
            }
        } catch (err) {
            setError('Failed to load conversations');
        } finally {
            setIsLoading(false);
        }
    }, [user]);
    // Load conversations on mount
    useEffect(() => {
        loadConversations();
    }, [loadConversations]);
    // Handle conversation selection
    const handleSelect = useCallback((conversation) => {
        onConversationSelect(conversation);
        setIsOpen(false);
    }, [onConversationSelect]);
    // Handle new conversation
    const handleNewConversation = useCallback(() => {
        onConversationSelect(null);
        setIsOpen(false);
    }, [onConversationSelect]);
    // Toggle dropdown
    const toggleDropdown = useCallback((e) => {
        e.stopPropagation();
        setIsOpen(!isOpen);
    }, [isOpen]);
    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (isOpen) {
                setIsOpen(false);
            }
        };
        if (isOpen) {
            document.addEventListener('click', handleClickOutside);
            return () => document.removeEventListener('click', handleClickOutside);
        }
    }, [isOpen]);
    // Format time ago
    const formatTimeAgo = useCallback((timestamp) => {
        const now = new Date();
        const time = new Date(timestamp);
        const diffInMinutes = Math.floor((now - time) / (1000 * 60));
        if (diffInMinutes < 1) return 'Just now';
        if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
        if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
        return `${Math.floor(diffInMinutes / 1440)}d ago`;
    }, []);
    // Find selected conversation
    const selectedConversation = conversations.find(c => c.id === selectedConversationId);
    if (!user) {
        return null;
    }
    if (error) {
        return (
            <div className={`${styles.selector} ${styles.error} ${className}`}>
                <span>Error loading conversations</span>
            </div>
        );
    }
    return (
        <div className={`${styles.container} ${className}`}>
            <button
                onClick={toggleDropdown}
                className={`${styles.selector} ${isOpen ? styles.open : ''}`}
                disabled={isLoading}
            >
                <div className={styles.selectedConversation}>
                    <MessageSquare size={16}/>
                    <div className={styles.conversationInfo}>
                        <span className={styles.conversationTitle}>
                            {isLoading ? 'Loading...' : (
                                selectedConversation?.title || 'New conversation'
                            )}
                        </span>
                        {selectedConversation && (
                            <span className={styles.conversationTime}>
                                {formatTimeAgo(selectedConversation.updatedAt ?? selectedConversation.updated_at ?? selectedConversation.createdAt ?? selectedConversation.created_at)}
                            </span>
                        )}
                    </div>
                </div>
                <ChevronDown size={16} className={styles.chevron}/>
            </button>
            {isOpen && (
                <div className={styles.dropdown}>
                    <div className={styles.dropdownContent}>
                        {/* New conversation option */}
                        <button
                            onClick={handleNewConversation}
                            className={`${styles.conversationOption} ${styles.newConversation} ${
                                !selectedConversationId ? styles.selected : ''
                            }`}
                        >
                            <div className={styles.conversationDetails}>
                                <div className={styles.conversationHeader}>
                                    <Plus size={16}/>
                                    <span className={styles.conversationTitle}>
                                        New conversation
                                    </span>
                                    {!selectedConversationId && (
                                        <Check size={16} className={styles.checkIcon}/>
                                    )}
                                </div>
                            </div>
                        </button>
                        {/* Existing conversations */}
                        {conversations.length > 0 && (
                            <>
                                <div className={styles.divider}/>
                                {conversations.map((conversation) => (
                                    <button
                                        key={conversation.id}
                                        onClick={() => handleSelect(conversation)}
                                        className={`${styles.conversationOption} ${
                                            selectedConversationId === conversation.id ? styles.selected : ''
                                        }`}
                                    >
                                        <div className={styles.conversationDetails}>
                                            <div className={styles.conversationHeader}>
                                                <span className={styles.conversationTitle}>
                                                    {conversation.title || 'Untitled conversation'}
                                                </span>
                                                {selectedConversationId === conversation.id && (
                                                    <Check size={16} className={styles.checkIcon}/>
                                                )}
                                            </div>
                                            <div className={styles.conversationMeta}>
                                                <Clock size={12}/>
                                                <span className={styles.conversationTime}>
                                                    {formatTimeAgo(conversation.updatedAt ?? conversation.updated_at ?? conversation.createdAt ?? conversation.created_at)}
                                                </span>
                                            </div>
                                            {conversation.lastMessage && (
                                                <div className={styles.lastMessage}>
                                                    {conversation.lastMessage.content?.slice(0, 60)}
                                                    {conversation.lastMessage.content?.length > 60 && '...'}
                                                </div>
                                            )}
                                        </div>
                                    </button>
                                ))}
                            </>
                        )}
                        {conversations.length === 0 && !isLoading && (
                            <div className={styles.emptyState}>
                                <MessageSquare size={20}/>
                                <span>No conversations yet</span>
                                <small>Start a new conversation to begin</small>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
});
ConversationSelector.displayName = 'ConversationSelector';
export default ConversationSelector; 
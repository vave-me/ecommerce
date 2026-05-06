"use client";
import React, { useCallback, useEffect, useState, useRef, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import {
    Bot,
    User,
    Clock,
    Copy,
    RefreshCw,
    MessageSquare,
    ThumbsUp,
    ThumbsDown,
    Share2 as Share,
    X,
    ChevronDown,
    Settings,
    Image as ImageIcon,
    ShoppingBag,
    FileText,
    Wrench,
    ArrowRight,
    Plus
} from '@/icons';
import { ClassifiedCard } from '../../../components/classified';
import { useAuth } from '../../../context/AuthContext';
import { useAIService } from '../../../hooks/ai/useAIService';
import PromptInput from '../../../components/AI/PromptInput';
import ConversationSelector from '../../../components/AI/ConversationSelector';
import AssistantSelector from '../../../components/AI/AssistantSelector';
import { toast } from 'react-toastify';
import styles from './AIPage.module.css';

/**
 * Production-Ready AI Chat Interface
 * 
 * Features:
 * - Automatic assistant loading and selection
 * - Conversation management with proper error handling
 * - Optimistic updates for smooth UX
 * - Matches backend API structure exactly
 */
const AIPage = () => {
    const router = useRouter();
    const t = useTranslations('AI');
    const { user } = useAuth();
    
    // Use the AI Service hook
    const {
        assistants,
        selectedAssistant,
        conversations,
        conversation: activeConversation,
        messages,
        loading,
        errors,
        isLoading,
        loadAssistants,
        selectAssistant,
        createNewConversation,
        loadConversations,
        loadConversation,
        sendMessage,
        clearErrors
    } = useAIService();

    // UI state
    const [showAssistantSelector, setShowAssistantSelector] = useState(false);
    const [showConversationSelector, setShowConversationSelector] = useState(false);
    const [isTyping, setIsTyping] = useState(false);
    
    // Refs
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const conversationCreatingRef = useRef(false);

    // Quick action templates
    const quickActions = useMemo(() => [
        {
            id: 'find-product',
            icon: ShoppingBag,
            title: t('quickActions.findProduct', 'Find Products'),
            description: t('quickActions.findProductDesc', 'Search for products in the marketplace'),
            prompt: 'Help me find '
        },
        {
            id: 'create-listing',
            icon: ImageIcon,
            title: t('quickActions.createListing', 'Create Listing'),
            description: t('quickActions.createListingDesc', 'List a new product with AI assistance'),
            prompt: 'I want to create a listing for '
        },
        {
            id: 'write-post',
            icon: FileText,
            title: t('quickActions.writePost', 'Write Post'),
            description: t('quickActions.writePostDesc', 'Create engaging content'),
            prompt: 'Help me write a post about '
        },
        {
            id: 'find-service',
            icon: Wrench,
            title: t('quickActions.findService', 'Find Services'),
            description: t('quickActions.findServiceDesc', 'Discover services near you'),
            prompt: 'I need a service for '
        }
    ], [t]);

    // Scroll to bottom of messages
    const scrollToBottom = useCallback(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, []);
    
    useEffect(() => {
        scrollToBottom();
    }, [messages, scrollToBottom]);
    
    /**
     * Ensure we have an active conversation
     */
    const ensureConversation = useCallback(async () => {
        if (!user || !selectedAssistant) {
            return null;
        }
        
        // If we already have a conversation, return it
        if (activeConversation) {
            return activeConversation;
        }
        
        // Prevent duplicate creation
        if (conversationCreatingRef.current) {
            return null;
        }
        
        conversationCreatingRef.current = true;
        try {
            const conversation = await createNewConversation({
                source: 'webApp',
                assistantName: selectedAssistant.name,
                userId: user.id || user.userId
            });
            
            return conversation;
        } catch (error) {
            // Error: '[AIPage] Failed to create conversation:', error...
            toast.error(t('conversationCreateError', 'Failed to create conversation'));
            return null;
        } finally {
            conversationCreatingRef.current = false;
        }
    }, [user, selectedAssistant, activeConversation, createNewConversation, t]);

    /**
     * Initialize assistants on mount
     */
    useEffect(() => {
        const initializeAssistants = async () => {
            const userId = user?.id || user?.userId;
            if (!userId || assistants.length > 0) {
                return;
            }
            
            try {
                await loadAssistants();
            } catch (error) {
                // Error: '[AIPage] Failed to load assistants:', error...
                // Error is already handled by the hook
            }
        };
        
        initializeAssistants();
    }, [user, assistants.length]); // Don't include loadAssistants to prevent loops
    
    /**
     * Create conversation when assistant is selected
     */
    useEffect(() => {
        if (selectedAssistant && !activeConversation && !conversationCreatingRef.current) {
            ensureConversation();
        }
    }, [selectedAssistant?.id]); // Only depend on ID

    /**
     * Handle sending message
     */
    const handleSendMessage = useCallback(async (messageContent, attachedFiles) => {
        if (!messageContent?.trim()) {
            return;
        }
        
        // Ensure we have a conversation
        let conversation = activeConversation;
        if (!conversation) {
            conversation = await ensureConversation();
            if (!conversation) {
                return;
            }
        }
        
        setIsTyping(true);
        clearErrors();
        
        try {
            await sendMessage(messageContent, {
                attachments: attachedFiles,
                timestamp: new Date().toISOString()
            });
            
            scrollToBottom();
        } catch (error) {
            // Error: '[AIPage] Failed to send message:', error...
            toast.error(error.message || t('messageSendError', 'Failed to send message'));
        } finally {
            setIsTyping(false);
        }
    }, [activeConversation, ensureConversation, sendMessage, scrollToBottom, clearErrors, t]);

    /**
     * Handle quick action click
     */
    const handleQuickAction = useCallback(async (action) => {
        // Send the action prompt as a message
        await handleSendMessage(action.prompt);
    }, [handleSendMessage]);

    /**
     * Handle message actions
     */
    const handleCopyMessage = useCallback((content) => {
        if (typeof content === 'object') {
            content = content.text || JSON.stringify(content);
        }
        navigator.clipboard.writeText(content);
        toast.success(t('messageCopied', 'Message copied to clipboard'));
    }, [t]);
    
    const handleRegenerateMessage = useCallback(async (messageId) => {
        // Find the user message before this AI message
        const messageIndex = messages.findIndex(msg => msg.id === messageId);
        if (messageIndex > 0) {
            const userMessage = messages[messageIndex - 1];
            if (userMessage.role === 'USER') {
                await handleSendMessage(userMessage.content);
            }
        }
    }, [messages, handleSendMessage]);

    /**
     * Render message content
     */
    const renderMessageContent = useCallback((message) => {
        let content = message.content;
        
        // Handle string content
        if (typeof content === 'string') {
            return (
                <div className={styles.messageText}>
                    {content.split('\n').map((line, i) => (
                        <p key={i}>{line || '\u00A0'}</p>
                    ))}
                </div>
            );
        }
        
        // Handle structured content
        if (content?.text) {
            return (
                <>
                    <div className={styles.messageText}>
                        {content.text.split('\n').map((line, i) => (
                            <p key={i}>{line || '\u00A0'}</p>
                        ))}
                    </div>
                    {content.products && content.products.length > 0 && (
                        <div className={styles.productGrid}>
                            {content.products.map((product) => (
                                <ClassifiedCard
                                    key={product.id}
                                    product={product}
                                    variant="compact"
                                />
                            ))}
                        </div>
                    )}
                </>
            );
        }
        
        return null;
    }, []);

    /**
     * Show loading state
     */
    if (loading.assistants && assistants.length === 0) {
        return (
            <div className={styles.aiPageWrapper}>
                <div className={styles.loadingContainer}>
                    <Bot className={styles.loadingIcon} />
                    <h2>{t('loadingTitle', 'Loading AI Assistant...')}</h2>
                    <p>{t('loadingDescription', 'Please wait while we connect to the service')}</p>
                </div>
            </div>
        );
    }
    
    /**
     * Show error state for critical errors
     */
    if (errors.assistants && assistants.length === 0) {
        return (
            <div className={styles.aiPageWrapper}>
                <div className={styles.errorStateContainer}>
                    <Bot className={styles.errorStateIcon} />
                    <h2>{t('errorTitle', 'AI Assistant is temporarily unavailable')}</h2>
                    <p>{errors.assistants}</p>
                    <button
                        className={styles.retryButton}
                        onClick={() => {
                            clearErrors();
                            loadAssistants();
                        }}
                    >
                        <RefreshCw />
                        {t('tryAgain', 'Try Again')}
                    </button>
                </div>
            </div>
        );
    }
    
    // Show authentication required message for non-logged in users
    const isAuthenticated = user && (user.userId || user.id);
    if (!isAuthenticated) {
        return (
            <div className={styles.aiPageWrapper}>
                <div className={styles.authRequiredContainer}>
                    <Bot className={styles.authRequiredIcon}/>
                    <h2>{t('authRequired', 'Sign In to Use AI Assistant')}</h2>
                    <p>{t('authRequiredDesc', 'To use the AI Assistant, you need to sign in to your account.')}</p>
                    <button
                        className={styles.signInButton}
                        onClick={() => router.push(`/pl/login`)}
                    >
                        <User/>
                        {t('signIn', 'Sign In')}
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.aiPageWrapper}>
            <div className={styles.aiPage}>
                {/* Header Bar */}
                <div className={styles.headerBar}>
                    <div className={styles.headerLeft}>
                        <h1 className={styles.pageTitle}>
                            <Bot className={styles.titleIcon}/>
                            {t('aiAssistant', 'AI Assistant')}
                        </h1>
                    </div>

                    <div className={styles.headerRight}>
                        <button
                            className={styles.headerButton}
                            onClick={() => setShowAssistantSelector(!showAssistantSelector)}
                            title={t('selectAssistant', 'Select Assistant')}
                            disabled={loading.assistants}
                        >
                            <Settings />
                            <span>{selectedAssistant?.name || t('selectAssistant', 'Select Assistant')}</span>
                            <ChevronDown />
                        </button>
                        
                        {user && (
                            <button
                                className={styles.headerButton}
                                onClick={() => {
                                    setShowConversationSelector(!showConversationSelector);
                                    if (!conversations.length) {
                                        loadConversations();
                                    }
                                }}
                                title={t('conversations', 'Conversations')}
                            >
                                <MessageSquare />
                                <span>{t('conversations', 'Conversations')}</span>
                            </button>
                        )}
                        
                        <button
                            className={styles.headerButton}
                            onClick={async () => {
                                if (selectedAssistant) {
                                    await createNewConversation();
                                }
                            }}
                            title={t('newChat', 'New Chat')}
                            disabled={!selectedAssistant || loading.conversations}
                        >
                            <Plus />
                            <span>{t('newChat', 'New')}</span>
                        </button>
                    </div>
                </div>

                {/* Assistant Selector Dropdown */}
                {showAssistantSelector && (
                    <AssistantSelector
                        assistants={assistants}
                        selectedAssistant={selectedAssistant}
                        onSelectAssistant={(assistant) => {
                            selectAssistant(assistant);
                            setShowAssistantSelector(false);
                        }}
                        onClose={() => setShowAssistantSelector(false)}
                    />
                )}

                {/* Conversation Selector Dropdown */}
                {showConversationSelector && (
                    <ConversationSelector
                        conversations={conversations}
                        currentConversation={activeConversation}
                        onSelectConversation={async (conversation) => {
                            await loadConversation(conversation.id);
                            setShowConversationSelector(false);
                        }}
                        onClose={() => setShowConversationSelector(false)}
                        loading={loading.conversations}
                    />
                )}

                {/* Chat Container */}
                <div className={styles.chatContainer}>
                    <div className={styles.messagesContainer} ref={messagesContainerRef}>
                        {messages.length === 0 ? (
                            /* Welcome Screen */
                            <div className={styles.welcomeScreen}>
                                <div className={styles.welcomeContent}>
                                    <Bot className={styles.welcomeIcon}/>
                                    <h2 className={styles.welcomeTitle}>
                                        {t('welcome', 'Welcome to AI Assistant')}
                                    </h2>
                                    <p className={styles.welcomeSubtitle}>
                                        {t('welcomeDesc', 'I can help you find products, create listings, write content, and more. How can I assist you today?')}
                                    </p>
                                </div>

                                {/* Quick Actions */}
                                <div className={styles.quickActions}>
                                    {quickActions.map((action) => (
                                        <button
                                            key={action.id}
                                            className={styles.quickActionCard}
                                            onClick={() => handleQuickAction(action)}
                                        >
                                            <action.icon className={styles.quickActionIcon}/>
                                            <div className={styles.quickActionContent}>
                                                <h3 className={styles.quickActionTitle}>{action.title}</h3>
                                                <p className={styles.quickActionDescription}>{action.description}</p>
                                            </div>
                                            <ArrowRight className={styles.quickActionArrow}/>
                                        </button>
                                    ))}
                                </div>
                            </div>
                        ) : (
                            /* Messages */
                            <>
                                {messages.map((message) => (
                                    <div
                                        key={message.id}
                                        className={`${styles.message} ${message.role === 'USER' || message.role === 'user' ? styles.user : styles.assistant}`}
                                    >
                                        <div className={styles.messageAvatar}>
                                            {message.role === 'USER' ? (
                                                <div className={styles.userAvatar}>
                                                    <User />
                                                </div>
                                            ) : (
                                                <div className={styles.aiAvatar}>
                                                    <Bot />
                                                </div>
                                            )}
                                        </div>

                                        <div className={styles.messageContent}>
                                            <div className={styles.messageHeader}>
                                                <span className={styles.messageName}>
                                                    {message.role === 'USER' ? (user?.username || t('you', 'You')) : (selectedAssistant?.name || t('aiAssistant', 'AI Assistant'))}
                                                </span>
                                                <span className={styles.messageTime}>
                                                <Clock className={styles.timeIcon}/>
                                                    {new Date(message.timestamp).toLocaleTimeString()}
                                            </span>
                                            </div>

                                            <div className={styles.messageBody}>
                                                {renderMessageContent(message)}
                                            </div>

                                            {message.role === 'ASSISTANT' && (
                                                <div className={styles.messageActions}>
                                                    <button
                                                        className={styles.actionButton}
                                                        onClick={() => handleCopyMessage(message.content)}
                                                        title={t('copy', 'Copy')}
                                                    >
                                                        <Copy />
                                                    </button>
                                                    <button
                                                        className={styles.actionButton}
                                                        onClick={() => handleRegenerateMessage(message.id)}
                                                        title={t('regenerate', 'Regenerate')}
                                                        disabled={loading.chat}
                                                    >
                                                        <RefreshCw />
                                                    </button>
                                                    <button
                                                        className={styles.actionButton}
                                                        title={t('like', 'Like')}
                                                    >
                                                        <ThumbsUp />
                                                    </button>
                                                    <button
                                                        className={styles.actionButton}
                                                        title={t('dislike', 'Dislike')}
                                                    >
                                                        <ThumbsDown />
                                                    </button>
                                                    <button
                                                        className={styles.actionButton}
                                                        title={t('share', 'Share')}
                                                    >
                                                        <Share />
                                                    </button>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                ))}

                                {/* Typing Indicator */}
                                {isTyping && (
                                    <div className={`${styles.message} ${styles.assistant}`}>
                                        <div className={styles.messageAvatar}>
                                            <div className={styles.aiAvatar}>
                                                <Bot/>
                                            </div>
                                        </div>
                                        <div className={styles.typingIndicator}>
                                            <span>{t('typing', 'AI is typing')}</span>
                                            <div className={styles.typingDots}>
                                                <span className={styles.typingDot}/>
                                                <span className={styles.typingDot}/>
                                                <span className={styles.typingDot}/>
                                            </div>
                                        </div>
                                    </div>
                                )}

                                <div ref={messagesEndRef}/>
                            </>
                        )}
                    </div>

                    {/* Error State */}
                    {(errors.chat || errors.messages) && (
                        <div className={styles.errorContainer}>
                            <X className={styles.errorIcon} />
                            <span className={styles.errorMessage}>
                                {errors.chat || errors.messages}
                            </span>
                            <button
                                className={styles.errorDismiss}
                                onClick={clearErrors}
                            >
                                {t('dismiss', 'Dismiss')}
                            </button>
                        </div>
                    )}
                </div>

                {/* Input Container */}
                <div className={styles.inputContainer}>
                    <PromptInput
                        onResponse={handleSendMessage}
                        conversationId={activeConversation?.id}
                        placeholder={t('messagePlaceholder', 'Type a message...')}
                        disabled={loading.chat || !selectedAssistant}
                    />
                </div>
            </div>
        </div>
    );
};

export default AIPage;
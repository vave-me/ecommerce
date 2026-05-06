"use client";
import React, {useEffect, useRef, useState, memo, useCallback} from 'react';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl';
import {useDispatch, useSelector} from 'react-redux';
import {useAuth} from '../../context/AuthContext';
import {useAIService} from '../../hooks/ai/useAIService';
import {
    switchToAiMode,
    setModeTransitioning,
    selectCurrentMode,
    selectIsAiMode,
    selectIsTransitioning,
    APP_MODES
} from '../../redux/slices/appModeSlice';
import {
    X as CloseIcon,
    Sparkles,
    Package,
    FileText,
    Wrench,
    Bot,
    Send,
    MessageSquare,
    ChevronUp,
    ChevronDown,
    User,
    ArrowUp,
    Minimize2,
    Maximize2
} from '@/icons';
import AssistantSelector from '../AI/AssistantSelector';
import styles from './AddOptionsSheet.module.css';

/**
 * Enhanced Bottom sheet component with AI Chat as primary feature
 * Features:
 * - AI Chat Interface (default open and prominent)
 * - Admin/Business Options (collapsed by default)
 * - Full AI Mode Switch
 * - WCAG 2.1 AA compliant
 * - Improved UX/UI with better visual hierarchy
 */
const AddOptionsSheet = memo(({
                                  isOpen,
                                  onClose,
                                  onSelectItem,
                                  items = [],
                                  showUnifiedComposer = true
                              }) => {
    const t = useTranslations('AddOptionsSheet');
    const tMode = useTranslations('ModeSwitcher');
    const dispatch = useDispatch();
    const {user} = useAuth();

    // State management
    const [isClosing, setIsClosing] = useState(false);
    const [message, setMessage] = useState('');
    const [activeItemId, setActiveItemId] = useState(null);
    const [focusedElementBeforeOpen, setFocusedElementBeforeOpen] = useState(null);

    // Refs
    const sheetRef = useRef(null);
    const messagesEndRef = useRef(null);
    const inputRef = useRef(null);
    const firstFocusableElementRef = useRef(null);
    const lastFocusableElementRef = useRef(null);

    // Redux selectors
    const currentMode = useSelector(selectCurrentMode);
    const isAiMode = useSelector(selectIsAiMode);
    const isTransitioning = useSelector(selectIsTransitioning);

    // AI Service integration
    const {
        assistants,
        loading: aiLoading,
        selectedAssistant,
        messages,
        sendMessage,
        createNewConversation,
        conversation: currentConversation
    } = useAIService();

    // Trap focus within modal for accessibility
    const trapFocus = useCallback((e) => {
        if (e.key !== 'Tab') return;

        if (!firstFocusableElementRef.current || !lastFocusableElementRef.current) return;

        if (e.shiftKey) {
            if (document.activeElement === firstFocusableElementRef.current) {
                e.preventDefault();
                lastFocusableElementRef.current.focus();
            }
        } else {
            if (document.activeElement === lastFocusableElementRef.current) {
                e.preventDefault();
                firstFocusableElementRef.current.focus();
            }
        }
    }, []);

    // Auto-scroll to bottom when new messages arrive
    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({
                behavior: 'smooth',
                block: 'end'
            });
        }
    }, [messages]);

    // Initialize conversation when sheet opens and AI chat is available
    useEffect(() => {
        if (isOpen && selectedAssistant && !currentConversation && user?.id) {
            createNewConversation();
        }
    }, [isOpen, selectedAssistant, currentConversation, user?.id, createNewConversation]);

    // Handle modal accessibility
    useEffect(() => {
        if (isOpen) {
            // Store currently focused element
            setFocusedElementBeforeOpen(document.activeElement);

            // Prevent body scroll
            document.body.style.overflow = 'hidden';
            document.body.setAttribute('aria-hidden', 'true');

            // Focus the chat input when opened (primary action)
            setTimeout(() => {
                if (inputRef.current) {
                    inputRef.current.focus();
                }
            }, 100);

            // Add focus trap
            document.addEventListener('keydown', trapFocus);
        } else {
            // Restore body scroll
            document.body.style.overflow = '';
            document.body.removeAttribute('aria-hidden');

            // Remove focus trap
            document.removeEventListener('keydown', trapFocus);

            // Restore focus to previous element
            if (focusedElementBeforeOpen && typeof focusedElementBeforeOpen.focus === 'function') {
                focusedElementBeforeOpen.focus();
            }
        }

        return () => {
            document.body.style.overflow = '';
            document.body.removeAttribute('aria-hidden');
            document.removeEventListener('keydown', trapFocus);
        };
    }, [isOpen, trapFocus, focusedElementBeforeOpen]);

    // Handle escape key press
    useEffect(() => {
        const handleEscKey = (e) => {
            if (e.key === 'Escape' && isOpen) {
                handleClose();
            }
        };

        if (isOpen) {
            document.addEventListener('keydown', handleEscKey);
        }

        return () => {
            document.removeEventListener('keydown', handleEscKey);
        };
    }, [isOpen]);

    // Reset state when sheet closes
    useEffect(() => {
        if (!isOpen) {
            setActiveItemId(null);
            setMessage('');
        }
    }, [isOpen]);

    // Handlers
    const handleClose = useCallback(() => {
        if (isClosing) return;
        setIsClosing(true);
        setTimeout(() => {
            setIsClosing(false);
            onClose();
        }, 200);
    }, [isClosing, onClose]);

    const handleSelectItem = useCallback((itemId) => {
        onSelectItem(itemId);
        handleClose();
    }, [onSelectItem, handleClose]);

    const handleSwitchToAiMode = useCallback(async () => {
        if (isTransitioning || isAiMode) return;

        try {
            dispatch(setModeTransitioning(true));
            dispatch(switchToAiMode());
            handleClose();
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    } finally {
            dispatch(setModeTransitioning(false));
        }
    }, [dispatch, isTransitioning, isAiMode, handleClose]);

    const handleSendMessage = useCallback(async () => {
        if (!message.trim() || aiLoading || !selectedAssistant) return;

        try {
            await sendMessage(message.trim());
            setMessage('');
            if (inputRef.current) {
                inputRef.current.focus();
            }
        } catch (error) {
            // Error: 'Failed to send message:', error...
            // Could add toast notification here
        }
    }, [message, aiLoading, selectedAssistant, sendMessage]);

    const handleKeyDown = useCallback((e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    }, [handleSendMessage]);

    if (!isOpen && !isClosing) return null;

    const isAdmin = user?.role === 'admin';
    const hasOptions = items && items.length > 0;

    return (
        <div
            className={`${styles.backdrop} ${isOpen && !isClosing ? styles.backdropOpen : ''} ${isClosing ? styles.backdropClosing : ''}`}
            onClick={handleClose}
            aria-modal="true"
            role="dialog"
            aria-label={t('ariaLabel', 'Create new content and chat with AI')}
            aria-describedby="sheet-description"
        >
            <div
                ref={sheetRef}
                className={`${styles.sheetContainer} ${isOpen && !isClosing ? styles.sheetOpen : ''} ${isClosing ? styles.sheetClosing : ''} ${styles.chatExpanded}`}
                onClick={(e) => e.stopPropagation()}
                role="document"
            >
                {/* Accessibility landmark */}
                <div className={styles.pullIndicator} aria-hidden="true"/>

                {/* Enhanced Header with better visual hierarchy */}
                <div className={styles.header}>
                    <div className={styles.headerContent}>
                        <h2 className={styles.title} id="sheet-title">
                            {t('title', 'AI Assistant & Content Creation')}
                        </h2>
                        <p className={styles.subtitle} id="sheet-description">
                            {t('subtitle', 'Chat with AI or create new content')}
                        </p>
                    </div>
                    <button
                        ref={firstFocusableElementRef}
                        className={styles.closeButton}
                        onClick={handleClose}
                        aria-label={t('closeButton', 'Close dialog')}
                        title={t('closeButton', 'Close dialog')}
                    >
                        <CloseIcon size={20} aria-hidden="true"/>
                    </button>
                </div>

                {/* PRIMARY SECTION: AI Chat Interface (always open) */}
                <div className={styles.primarySection}>
                    <div className={styles.sectionHeader}>
                        <div className={styles.sectionTitleGroup}>
                            <Bot size={20} className={styles.sectionIcon} aria-hidden="true"/>
                            <h3 className={styles.sectionTitle}>
                                {t('aiChatTitle', 'AI Assistant')}
                            </h3>
                            {selectedAssistant && (
                                <span className={styles.assistantBadge}>
                                    {selectedAssistant.name}
                                </span>
                            )}
                        </div>
                    </div>

                    {/* Assistant Selector - Always visible for quick switching */}
                    <div className={styles.assistantSelector}>
                        <AssistantSelector
                            assistants={assistants}
                            selectedAssistant={selectedAssistant}
                            loading={aiLoading.assistants}
                            compact={true}
                        />
                    </div>

                    {/* Chat Interface - Always expanded */}
                    <div
                        id="ai-chat-content"
                        className={`${styles.chatInterface} ${styles.expanded}`}
                    >
                        {/* Messages Container with better scrolling */}
                        <div
                            className={styles.messagesContainer}
                            role="log"
                            aria-live="polite"
                            aria-label={t('chatMessages', 'Chat messages')}
                        >
                            {messages?.length > 0 ? (
                                messages.map((msg, index) => (
                                    <div
                                        key={`${msg.id || index}-${msg.timestamp || Date.now()}`}
                                        className={`${styles.message} ${msg.role === 'user' ? styles.userMessage : styles.assistantMessage}`}
                                        role="article"
                                        aria-label={`${msg.role === 'user' ? t('yourMessage', 'Your message') : t('assistantMessage', 'Assistant message')}: ${msg.content}`}
                                    >
                                        <div className={styles.messageIcon} aria-hidden="true">
                                            {msg.role === 'user' ? (
                                                <User size={16}/>
                                            ) : (
                                                <Bot size={16}/>
                                            )}
                                        </div>
                                        <div className={styles.messageContent}>
                                            <div className={styles.messageText}>
                                                {msg.content}
                                            </div>
                                            {msg.timestamp && (
                                                <time className={styles.messageTime} dateTime={msg.timestamp}>
                                                    {new Date(msg.timestamp).toLocaleTimeString()}
                                                </time>
                                            )}
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className={styles.emptyChat}>
                                    <div className={styles.emptyChatIcon} aria-hidden="true">
                                        <MessageSquare size={48}/>
                                    </div>
                                    <h4 className={styles.emptyChatTitle}>
                                        {t('startConversation', 'Start a conversation')}
                                    </h4>
                                    <p className={styles.emptyChatDescription}>
                                        {t('chatDescription', 'Ask questions, get help, or explore ideas with your AI assistant')}
                                    </p>
                                </div>
                            )}
                            <div ref={messagesEndRef} aria-hidden="true"/>
                        </div>

                        {/* Enhanced Input Area - WCAG 2.1 AA Compliant */}
                        <div className={styles.inputContainer}>
                            <div className={styles.inputWrapper}>
                                <label htmlFor="chat-input" className={styles.inputLabel}>
                                    {t('messageLabel', 'Message to AI Assistant')}
                                </label>
                                <textarea
                                    id="chat-input"
                                    ref={inputRef}
                                    value={message}
                                    onChange={(e) => setMessage(e.target.value)}
                                    onKeyDown={handleKeyDown}
                                    placeholder={t('messagePlaceholder', 'Type your message here...')}
                                    className={styles.messageInput}
                                    disabled={!selectedAssistant || aiLoading}
                                    aria-describedby="input-help input-count"
                                    maxLength={1000}
                                    rows={1}
                                    style={{
                                        minHeight: '60px',
                                        maxHeight: '120px',
                                        resize: 'none',
                                        overflow: 'hidden'
                                    }}
                                    onInput={(e) => {
                                        e.target.style.height = 'auto';
                                        e.target.style.height = Math.min(e.target.scrollHeight, 120) + 'px';
                                    }}
                                />
                                <button
                                    onClick={handleSendMessage}
                                    disabled={!message.trim() || aiLoading || !selectedAssistant || message.length > 1000}
                                    className={styles.sendButton}
                                    aria-label={t('sendMessage', 'Send message')}
                                    title={t('sendMessage', 'Send message')}
                                >
                                    {aiLoading ? (
                                        <div className={styles.spinner} aria-hidden="true"/>
                                    ) : (
                                        <Send size={20} aria-hidden="true"/>
                                    )}
                                </button>
                            </div>

                            <div className={styles.inputMeta}>
                                <div id="input-help" className={styles.inputHelp}>
                                    <span>{t('inputHelp', 'Press Enter to send, Shift+Enter for new line')}</span>
                                </div>
                                <div id="input-count" className={styles.inputCount}>
                                    <span
                                        className={`${styles.characterCount} ${
                                            message.length > 950 ? styles.error :
                                                message.length > 800 ? styles.warning : ''
                                        }`}
                                        aria-live="polite"
                                    >
                                        {message.length}/1000
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* SECONDARY SECTION: Quick Actions & AI Mode Switch */}
                <div className={styles.secondarySection}>
                    {/* AI Mode Switch - Prominent */}
                    <div className={styles.aiModeContainer}>
                        <button
                            className={`${styles.aiModeButton} ${isTransitioning ? styles.loading : ''}`}
                            onClick={handleSwitchToAiMode}
                            disabled={isTransitioning || isAiMode}
                            aria-label={tMode('switchToAiMode', 'Switch to full AI mode')}
                        >
                            <div className={styles.iconWrapper} aria-hidden="true">
                                <Sparkles size={24}/>
                            </div>
                            <div className={styles.textContent}>
                                <span className={styles.optionLabel}>
                                    {isTransitioning ?
                                        t('switching', 'Switching...') :
                                        t('fullAiMode', 'Full AI Mode')
                                    }
                                </span>
                                <span className={styles.optionDescription}>
                                    {t('aiModeDescription', 'Switch to the complete AI interface')}
                                </span>
                            </div>
                            <ArrowUp size={16} className={styles.modeIcon} aria-hidden="true"/>
                        </button>
                    </div>

                    {/* Content Creation Options (always visible) */}
                    {hasOptions && (
                        <div className={styles.optionsContainer}>
                            <div className={styles.optionsHeader}>
                                <Package size={20} aria-hidden="true"/>
                                <h3 className={styles.optionsTitle}>{t('createContent', 'Create Content')}</h3>
                            </div>

                            <div
                                id="content-options"
                                className={`${styles.optionsList} ${styles.expanded}`}
                            >
                                {items.map((item, index) => (
                                    <div key={item.id} className={styles.optionItem}>
                                        <button
                                            className={`${styles.optionButton} ${
                                                isAdmin ? styles.adminOption : ''
                                            } ${activeItemId === item.id ? styles.optionButtonActive : ''}`}
                                            onClick={() => handleSelectItem(item.id)}
                                            onMouseEnter={() => setActiveItemId(item.id)}
                                            onMouseLeave={() => setActiveItemId(null)}
                                            onFocus={() => setActiveItemId(item.id)}
                                            onBlur={() => setActiveItemId(null)}
                                            aria-label={t('createAriaLabel', {label: item.label})}
                                        >
                                            <div className={styles.iconWrapper} aria-hidden="true">
                                                <item.icon size={20}/>
                                            </div>
                                            <div className={styles.textContent}>
                                                <span className={styles.optionLabel}>{item.label}</span>
                                                {item.description && (
                                                    <span className={styles.optionDescription}>{item.description}</span>
                                                )}
                                            </div>
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer with Cancel Button */}
                <div className={styles.footer}>
                    <button
                        ref={lastFocusableElementRef}
                        className={styles.cancelButton}
                        onClick={handleClose}
                        aria-label={t('cancelButtonAria', 'Cancel and close')}
                    >
                        {t('cancelButtonText', 'Cancel')}
                    </button>
                </div>
            </div>
        </div>
    );
});

AddOptionsSheet.displayName = 'AddOptionsSheet';
AddOptionsSheet.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onSelectItem: PropTypes.func.isRequired,
    items: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            label: PropTypes.string.isRequired,
            description: PropTypes.string,
            icon: PropTypes.elementType.isRequired,
        })
    ),
    showUnifiedComposer: PropTypes.bool
};

export default AddOptionsSheet;
"use client";
import React, { useState, useCallback, useRef, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import { useTranslations } from 'next-intl';
import { useDispatch, useSelector } from 'react-redux';
import { useAuth } from '../../context/AuthContext';
import { useAIService } from '../../hooks/ai/useAIService';
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
    User,
    ArrowUp
} from '@/icons';
import AssistantSelector from '../AI/AssistantSelector';
import styles from './AddOptionsSheet.module.css';

/**
 * Enhanced Bottom sheet component for admin/business users
 * Features: AI Chat Interface, Admin Options, Full AI Mode Switch
 * Three-section layout: Admin Options (top), AI Mode Switch (middle), AI Chat (bottom)
 */
const AddOptionsSheetWithComposer = memo(({
    isOpen,
    onClose,
    onSelectItem,
    items,
    showUnifiedComposer = true
}) => {
    const t = useTranslations('AddOptionsSheet');
    const tMode = useTranslations('ModeSwitcher');
    const dispatch = useDispatch();
    const { user } = useAuth();
    
    // State management
    const [isClosing, setIsClosing] = useState(false);
    const [message, setMessage] = useState('');
    const [activeItemId, setActiveItemId] = useState(null);
    
    // Refs
    const sheetRef = useRef(null);
    const messagesEndRef = useRef(null);
    const inputRef = useRef(null);
    
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

    // Auto-scroll to bottom when new messages arrive
    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [messages]);

    // Initialize conversation when sheet opens
    useEffect(() => {
        if (isOpen && selectedAssistant && !currentConversation && user?.id) {
            createNewConversation();
        }
    }, [isOpen, selectedAssistant, currentConversation, user?.id, createNewConversation]);

    // Handle escape key press & body scroll lock
    useEffect(() => {
        const handleEscKey = (e) => {
            if (e.key === 'Escape' && isOpen) {
                handleClose();
            }
        };
        document.addEventListener('keydown', handleEscKey);
        if (isOpen) {
            document.body.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = '';
        }
        return () => {
            document.removeEventListener('keydown', handleEscKey);
            document.body.style.overflow = '';
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
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, [message, aiLoading, selectedAssistant, sendMessage]);

    const handleKeyPress = useCallback((e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    }, [handleSendMessage]);

    if (!isOpen && !isClosing) return null;

    const isAdmin = user?.role === 'admin';

    return (
        <div
            className={`${styles.backdrop} ${isOpen && !isClosing ? styles.backdropOpen : ''} ${isClosing ? styles.backdropClosing : ''}`}
            onClick={handleClose}
            aria-modal="true"
            role="dialog"
            aria-label={t('ariaLabel')}
        >
            <div
                ref={sheetRef}
                className={`${styles.sheetContainer} ${isOpen && !isClosing ? styles.sheetOpen : ''} ${isClosing ? styles.sheetClosing : ''} ${styles.chatExpanded}`}
                onClick={(e) => e.stopPropagation()}
                role="document"
            >
                <div className={styles.pullIndicator} />
                
                <div className={styles.header}>
                    <h3 className={styles.title}>Neue Inhalte erstellen</h3>
                </div>

                {/* THREE-SECTION LAYOUT */}
                
                {/* TOP SECTION: Admin Options (Only for Admins) */}
                {isAdmin && (
                    <div className={styles.adminSection}>
                        <div className={styles.sectionHeader}>
                            <h4 className={styles.sectionTitle}>Admin Optionen</h4>
                        </div>
                        <ul className={styles.optionsList}>
                            {items.map((item) => (
                                <li key={item.id} className={styles.optionItem}>
                                    <button
                                        className={`${styles.optionButton} ${styles.adminOption} ${activeItemId === item.id ? styles.optionButtonActive : ''}`}
                                        onClick={() => handleSelectItem(item.id)}
                                        onMouseEnter={() => setActiveItemId(item.id)}
                                        onMouseLeave={() => setActiveItemId(null)}
                                        onFocus={() => setActiveItemId(item.id)}
                                        onBlur={() => setActiveItemId(null)}
                                        aria-label={t('createAriaLabel', { label: item.label })}
                                    >
                                        <div className={styles.iconWrapper} aria-hidden="true">
                                            <item.icon size={20} />
                                        </div>
                                        <div className={styles.textContent}>
                                            <span className={styles.optionLabel}>{item.label}</span>
                                            {item.description && (
                                                <span className={styles.optionDescription}>{item.description}</span>
                                            )}
                                        </div>
                                    </button>
                                </li>
                            ))}
                        </ul>
                    </div>
                )}

                {/* MIDDLE SECTION: AI Mode Switch */}
                <div className={styles.aiModeSection}>
                    <div className={styles.sectionHeader}>
                        <h4 className={styles.sectionTitle}>KI Modus</h4>
                    </div>
                    <button
                        className={`${styles.optionButton} ${styles.aiModeButton} ${isTransitioning ? styles.loading : ''}`}
                        onClick={handleSwitchToAiMode}
                        disabled={isTransitioning || isAiMode}
                        aria-label={tMode('switchToAiMode')}
                    >
                        <div className={styles.iconWrapper} aria-hidden="true">
                            <Bot size={24} />
                        </div>
                        <div className={styles.textContent}>
                            <span className={styles.optionLabel}>
                                {isTransitioning ? 'Wechsle...' : 'Vollständiger KI-Modus'}
                            </span>
                            <span className={styles.optionDescription}>
                                Wechsle zur vollständigen KI-Assistenten-Oberfläche
                            </span>
                        </div>
                        <ArrowUp size={16} className={styles.expandIcon} />
                    </button>
                </div>

                {/* BOTTOM SECTION: AI Chat Interface - Always visible */}
                <div className={styles.aiChatSection}>
                    <div className={styles.sectionHeader}>
                        <h4 className={styles.sectionTitle}>KI Chat</h4>
                    </div>

                    {/* Assistant Selector */}
                    <div className={styles.assistantSelector}>
                        <AssistantSelector
                            assistants={assistants}
                            selectedAssistant={selectedAssistant}
                            loading={aiLoading.assistants}
                            compact={true}
                        />
                    </div>

                    {/* Chat Interface - Always visible */}
                    <div className={styles.chatInterface}>
                        {/* Messages */}
                        <div className={styles.messagesContainer}>
                            {messages?.length > 0 ? (
                                messages.map((msg, index) => (
                                    <div
                                        key={index}
                                        className={`${styles.message} ${msg.role === 'user' ? styles.userMessage : styles.assistantMessage}`}
                                    >
                                        <div className={styles.messageIcon}>
                                            {msg.role === 'user' ? (
                                                <User size={16} />
                                            ) : (
                                                <Bot size={16} />
                                            )}
                                        </div>
                                        <div className={styles.messageContent}>
                                            {msg.content}
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className={styles.emptyChat}>
                                    <Bot size={32} className={styles.emptyChatIcon} />
                                    <p>Beginne eine Unterhaltung mit deinem KI-Assistenten</p>
                                </div>
                            )}
                            <div ref={messagesEndRef} />
                        </div>

                        {/* Enhanced Input Area - WCAG 2.1 Compliant */}
                        <div className={styles.inputContainer}>
                            <textarea
                                ref={inputRef}
                                value={message}
                                onChange={(e) => setMessage(e.target.value)}
                                onKeyPress={handleKeyPress}
                                placeholder="Schreibe deine Nachricht an den KI-Assistenten..."
                                className={styles.messageInput}
                                disabled={!selectedAssistant || aiLoading}
                                aria-label="Nachricht an KI-Assistenten"
                                aria-describedby="input-help"
                                maxLength={1000}
                            />
                            <div className={styles.inputActions}>
                                <div id="input-help" className={styles.inputHelp}>
                                    <span>Drücke Enter zum Senden</span>
                                    {message.length > 800 && (
                                        <span className={`${styles.characterCount} ${message.length > 950 ? styles.error : styles.warning}`}>
                                            {message.length}/1000
                                        </span>
                                    )}
                                </div>
                                <button
                                    onClick={handleSendMessage}
                                    disabled={!message.trim() || aiLoading || !selectedAssistant || message.length > 1000}
                                    className={styles.sendButton}
                                    aria-label="Nachricht senden"
                                    title="Nachricht senden"
                                >
                                    {aiLoading ? (
                                        <div className={styles.spinner} aria-hidden="true" />
                                    ) : (
                                        <Send size={20} />
                                    )}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Cancel Button */}
                <button
                    className={styles.cancelButton}
                    onClick={handleClose}
                    aria-label={t('cancelButtonAria')}
                >
                    {t('cancelButtonText')}
                </button>
            </div>
        </div>
    );
});

AddOptionsSheetWithComposer.displayName = 'AddOptionsSheetWithComposer';
AddOptionsSheetWithComposer.propTypes = {
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
    ).isRequired,
    showUnifiedComposer: PropTypes.bool
};

export default AddOptionsSheetWithComposer;
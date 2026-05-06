"use client";

import React, {useState, useRef, useCallback, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {useRouter} from 'next/navigation';
import {
    User as UserIcon,
    Mic,
    Bot,
    Image as ImageIcon,
    Video as VideoIcon,
    Zap,
    Sun,
    Battery,
    Sparkles,
    Tag,
    Send,
    MoreHorizontal,
    Search,
    Filter,
    X,
    Check,
    Flame
} from '@/icons';
import {useAuth} from '../../context/AuthContext';
import {useDispatch, useSelector} from 'react-redux';
import {setFilters} from '../../redux/slices/listingFiltersSlice';
import {FilePond, registerPlugin} from 'react-filepond';
import 'filepond/dist/filepond.min.css';
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';
import {useAIService} from '../../hooks/ai/useAIService';
import styles from './Composer.module.css';

registerPlugin(FilePondPluginFileValidateType, FilePondPluginFileValidateSize);

/**
 * Composer Component - Clean implementation
 * Handles AI chat and regular post creation
 */
const Composer = ({
                      placeholder = "Search, ask AI, or create content...",
                      showTemplates = true,
                      enableChat = true,
                      autoFocusOnAI = false,
                      onPostCreate = null,
                      onFilterUpdate = null,
                      onClose = null
                  }) => {
    const {user} = useAuth();
    const t = useTranslations('AI');
    const dispatch = useDispatch();
    const router = useRouter();

    // State Management
    const [text, setText] = useState('');
    const [images, setImages] = useState([]);
    const [videos, setVideos] = useState([]);
    const [searchMode, setSearchMode] = useState(false);
    const [expanded, setExpanded] = useState(false);
    const [showMediaUpload, setShowMediaUpload] = useState(false);
    const [chatMode, setChatMode] = useState(true);
    const [isSending, setIsSending] = useState(false);

    // Refs
    const textareaRef = useRef(null);
    const chatBoxRef = useRef(null);

    // Use AI Service hook
    const ai = useAIService();
    const {
        assistants,
        loading: aiLoading,
        errors: aiErrors,
        selectedAssistant,
        messages,
        sendMessage,
        createNewConversation,
        conversation: currentConversation,
        isLoading,
        latestError
    } = ai;
    const {isInitialized, isInitializing, initializationError, retryInitialization} = ai;

    // Initialize chat conversation when needed
    useEffect(() => {
        const initChat = async () => {
            const userId = user?.id || user?.userId;

            // Wait for assistants to be initialized first
            if (!isInitialized) {
                return;
            }

            // Only create conversation if we have everything needed and no conversation exists
            if (enableChat && userId && selectedAssistant && !currentConversation && !aiLoading.conversation) {
                try {
                    await createNewConversation();
                } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
            }
        };

        initChat();
    }, [enableChat, user?.id, user?.userId, selectedAssistant?.id, isInitialized, currentConversation, aiLoading.conversation]);

    // Auto-focus on AI mode
    useEffect(() => {
        if (autoFocusOnAI && chatMode && textareaRef.current) {
            textareaRef.current.focus();
        }
    }, [chatMode, autoFocusOnAI]);

    // Scroll to bottom in chat mode
    useEffect(() => {
        if (chatMode && chatBoxRef.current) {
            chatBoxRef.current.scrollTop = chatBoxRef.current.scrollHeight;
        }
    }, [messages, chatMode]);

    // Quick action templates
    const quickTemplates = showTemplates ? [
        {icon: Sparkles, text: "LED lighting solutions", category: "search"},
        {icon: Sun, text: "Solar panels & systems", category: "search"},
        {icon: Battery, text: "Energy storage systems", category: "search"},
        {icon: Zap, text: "EV charging stations", category: "search"},
        {icon: Flame, text: "Energy efficiency products", category: "search"},
        {icon: Search, text: "Smart lighting controls", category: "search"}
    ] : [];

    // Handle template click
    const handleTemplateClick = (template) => {
        setText(template.text);
        setSearchMode(true);
        if (textareaRef.current) {
            textareaRef.current.focus();
        }
    };

    // Handle form submission
    const handleSubmit = async (e) => {
        e.preventDefault();

        const trimmedText = text.trim();
        if (!trimmedText || aiLoading.conversation || isSending) return;

        if (chatMode && enableChat) {
            // Check if assistants are initialized
            if (!isInitialized) {
                
                return;
            }

            // Check authentication before sending
            const userId = user?.id || user?.userId;
            if (!userId) {
                
                return;
            }

            // Check if we have a selected assistant
            if (!selectedAssistant) {
                
                return;
            }

            // Ensure we have a conversation
            if (!currentConversation && !aiLoading.conversation) {
                setIsSending(true);
                try {
                    const newConv = await createNewConversation();
                    if (!newConv) {
                        // Error: 'Failed to create conversation'...
                        return;
                    }
                    // Wait a bit for the conversation to be ready
                    await new Promise(resolve => setTimeout(resolve, 100));
                } catch (error) {
                    // Error: 'Error creating conversation:', error...
                    return;
                } finally {
                    setIsSending(false);
                }
            }

            setText('');
            setIsSending(true);

            try {
                // Send message
                await sendMessage(trimmedText, {
                    timestamp: new Date().toISOString(),
                    has_attachments: images.length + videos.length > 0 ? 'true' : 'false'
                });
            } catch (error) {
                // Restore text on error
                // Error: 'Failed to send message:', error...
                setText(trimmedText);
            } finally {
                setIsSending(false);
            }
        } else if (searchMode) {
            // Search mode
            dispatch(setFilters({
                searchText: trimmedText,
                tags: trimmedText.split(' ').filter(word => word.length > 2)
            }));

            if (onFilterUpdate) {
                onFilterUpdate({
                    searchText: trimmedText,
                    tags: trimmedText.split(' ').filter(word => word.length > 2)
                });
            }

            setText('');
            setSearchMode(false);
        } else {
            // Regular post creation
            if (onPostCreate) {
                await onPostCreate({text: trimmedText, images, videos});
            }

            setText('');
            setImages([]);
            setVideos([]);
        }
    };

    // Media handling
    const handleImageUpload = (files) => {
        setImages(prev => [...prev, ...files]);
    };

    const removeImage = (index) => {
        setImages(prev => prev.filter((_, i) => i !== index));
    };

    // Render chat interface if in chat mode
    if (chatMode && enableChat) {
        // Check if user is authenticated
        const userId = user?.id || user?.userId;
        if (!userId) {
            return (
                <div className={styles.chatContainerCompact}>
                    <div className={styles.authPromptHeader}>
                        <Sparkles className={styles.sparkleIcon}/>
                        {onClose && (
                            <button
                                onClick={onClose}
                                className={styles.closePrompt}
                                aria-label="Close AI Assistant prompt"
                            >
                                <X size={16}/>
                            </button>
                        )}
                    </div>
                    <div className={styles.authRequiredContent}>
                        <div className={styles.authMessage}>
                            <p className={styles.authGreeting}>
                                {t('authGreeting', 'Odkryj moc AI w znajdowaniu idealnych produktów!')}
                            </p>
                            <p className={styles.authSubtext}>
                                {t('authSubtext', 'Dołącz teraz i pozwól inteligentnemu asystentowi znaleźć dokładnie to, czego szukasz')}
                            </p>
                        </div>
                        <div className={styles.authActions}>
                            <button
                                onClick={() => router.push('/register')}
                                className={styles.registerButtonCompact}
                                aria-label="Register to use AI Assistant"
                            >
                                {t('registerAction', 'Załóż konto')}
                                <Sparkles size={14}/>
                            </button>
                            <button
                                onClick={() => router.push('/login')}
                                className={styles.loginLinkCompact}
                                aria-label="Log in to use AI Assistant"
                            >
                                {t('loginAction', 'Mam już konto')}
                            </button>
                        </div>
                    </div>
                </div>
            );
        }

        return (
            <div className={styles.chatContainer}>
                <div className={styles.chatHeader}>
                    <Bot className={styles.chatIcon}/>
                    <span className={styles.chatTitle}>
            {selectedAssistant?.name || 'AI Assistant'}
          </span>
                    {onClose && (
                        <button
                            onClick={onClose}
                            className={styles.closeChat}
                            aria-label="Close AI Assistant"
                        >
                            <X/>
                        </button>
                    )}
                </div>

                <div ref={chatBoxRef} className={styles.chatMessages}>
                    {/* Show initialization error with retry button */}
                    {initializationError && (
                        <div className={styles.errorMessage}>
                            <p>{t('aiServiceError', 'AI service is temporarily unavailable')}</p>
                            <button
                                onClick={retryInitialization}
                                className={styles.retryButton}
                            >
                                <RefreshCw size={14}/>
                                {t('retry', 'Try again')}
                            </button>
                        </div>
                    )}

                    {/* Show loading state during initialization */}
                    {isInitializing && (
                        <div className={styles.loadingMessage}>
                            <div className={styles.loadingSpinner}/>
                            <p>{t('loadingAssistants', 'Loading AI assistants...')}</p>
                        </div>
                    )}

                    {/* Show welcome message when ready */}
                    {messages.length === 0 && !isInitializing && !initializationError && isInitialized && (
                        <div className={styles.welcomeMessage}>
                            <Bot/>
                            <p>{t('greeting', 'Hello! I can help you find energy-efficient lighting, solar panels, EV charging stations, and more. What are you looking for today?')}</p>
                        </div>
                    )}

                    {/* Show messages */}
                    {messages.map((msg, idx) => (
                        <div
                            key={msg.id || idx}
                            className={`${styles.chatMessage} ${msg.role === 'USER' ? styles.user : styles.ai}`}
                        >
                            {msg.role === 'USER' ? <UserIcon/> : <Bot/>}
                            <div className={styles.messageText}>{msg.content}</div>
                        </div>
                    ))}

                    {/* Show typing indicator */}
                    {(aiLoading.conversation || isSending) && (
                        <div className={`${styles.chatMessage} ${styles.ai}`}>
                            <div className={styles.messageAvatar}>
                                <Bot/>
                            </div>
                            <div className={styles.typingIndicator}>
                                <span></span><span></span><span></span>
                            </div>
                        </div>
                    )}
                </div>

                <form onSubmit={handleSubmit} className={styles.chatInputForm}>
                    <input
                        ref={textareaRef}
                        type="text"
                        value={text}
                        onChange={(e) => setText(e.target.value)}
                        placeholder={
                            isInitializing ? "Loading AI assistant..." :
                                initializationError ? "AI service is temporarily unavailable" :
                                    aiLoading.conversation ? "Setting up conversation..." :
                                        isSending ? "AI is thinking..." :
                                            !selectedAssistant ? "Waiting for assistant..." :
                                                "Type your message..."
                        }
                        className={styles.chatInput}
                        disabled={!isInitialized || aiLoading.conversation || isSending || !selectedAssistant || initializationError}
                    />
                    <button
                        type="submit"
                        disabled={!text.trim() || !isInitialized || aiLoading.conversation || isSending || !selectedAssistant || initializationError}
                        className={styles.chatSend}
                    >
                        <Send/>
                    </button>
                </form>
            </div>
        );
    }

    // Render standard composer interface
    return (
        <div className={`${styles.composerContainer} ${expanded ? styles.expanded : ''}`}>
            {/* Quick templates */}
            {showTemplates && !expanded && (
                <div className={styles.quickTemplates}>
                    {quickTemplates.map((template, idx) => (
                        <button
                            key={idx}
                            onClick={() => handleTemplateClick(template)}
                            className={styles.templateButton}
                        >
                            <template.icon size={16}/>
                            <span>{template.text}</span>
                        </button>
                    ))}
                </div>
            )}

            {/* Main composer form */}
            <form onSubmit={handleSubmit} className={styles.composerForm}>
                <div className={styles.inputWrapper}>
                    {/* Mode indicators */}
                    <div className={styles.modeIndicators}>
                        {searchMode && (
                            <span className={styles.searchIndicator}>
                <Search size={16}/>
                Search
              </span>
                        )}
                    </div>

                    {/* Text input */}
                    <textarea
                        ref={textareaRef}
                        value={text}
                        onChange={(e) => setText(e.target.value)}
                        onFocus={() => setExpanded(true)}
                        placeholder={placeholder}
                        className={styles.textInput}
                        rows={expanded ? 3 : 1}
                    />

                    {/* Action buttons */}
                    <div className={styles.actions}>
                        <button
                            type="button"
                            onClick={() => setShowMediaUpload(!showMediaUpload)}
                            className={styles.actionButton}
                            aria-label="Add media"
                        >
                            <ImageIcon/>
                        </button>

                        <button
                            type="button"
                            onClick={() => {}}
                            className={styles.actionButton}
                            aria-label="Voice record"
                        >
                            <Mic/>
                        </button>

                        {enableChat && (
                            <button
                                type="button"
                                onClick={() => setChatMode(!chatMode)}
                                className={`${styles.actionButton} ${chatMode ? styles.active : ''}`}
                                aria-label="Toggle AI mode"
                            >
                                <Bot/>
                            </button>
                        )}

                        <button
                            type="submit"
                            disabled={!text.trim() || loading || isSending}
                            className={styles.submitButton}
                            aria-label="Send"
                        >
                            <Send/>
                        </button>
                    </div>
                </div>

                {/* Media upload section */}
                {showMediaUpload && expanded && (
                    <div className={styles.mediaUpload}>
                        <FilePond
                            files={images}
                            onupdatefiles={setImages}
                            allowMultiple={true}
                            maxFiles={5}
                            acceptedFileTypes={['image/*']}
                            labelIdle='Drag & Drop images or <span class="filepond--label-action">Browse</span>'
                            className={styles.filePond}
                        />
                    </div>
                )}

                {/* Media previews */}
                {images.length > 0 && (
                    <div className={styles.mediaPreviews}>
                        {images.map((file, idx) => (
                            <div key={idx} className={styles.mediaPreview}>
                                <img src={URL.createObjectURL(file.file)} alt=""/>
                                <button onClick={() => removeImage(idx)}><X/></button>
                            </div>
                        ))}
                    </div>
                )}
            </form>
        </div>
    );
};

export default Composer;
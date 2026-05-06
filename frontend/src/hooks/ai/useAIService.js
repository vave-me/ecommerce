import { useCallback, useEffect, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { useAuth } from '../../context/AuthContext';
import { 
    ConversationNotFoundError,
    AssistantNotFoundError 
} from '../../services/ai/AssistantServiceWithRetry';

// Redux actions and selectors
import {
    setSelectedAssistant,
    setActiveConversation,
    startConversationCreation,
    setTyping,
    setUIError,
    clearUIError,
    selectSelectedAssistantId,
    selectActiveConversationId,
    selectIsTyping,
    selectViewMode,
    selectIsCreatingConversation
} from '../../redux/slices/aiUISlice';

// React Query hooks
import {
    useAssistantsQuery,
    useConversationQuery,
    useMessagesQuery,
    useConversationsQuery,
    useCreateConversationMutation,
    useSendMessageMutation,
    aiQueryKeys
} from './useAIQueries';

/**
 * Unified AI Service Hook
 * Combines Redux (UI state) with React Query (server state)
 * Provides a single interface for all AI functionality
 */
export const useAIService = () => {
    const dispatch = useDispatch();
    const queryClient = useQueryClient();
    const { user } = useAuth();
    
    // UI State from Redux
    const selectedAssistantId = useSelector(selectSelectedAssistantId);
    const activeConversationId = useSelector(selectActiveConversationId);
    const isTyping = useSelector(selectIsTyping);
    const viewMode = useSelector(selectViewMode);
    const isCreatingConversation = useSelector(selectIsCreatingConversation);
    
    // Server State from React Query
    const { 
        data: assistantsData, 
        isLoading: assistantsLoading,
        error: assistantsError,
        refetch: refetchAssistants
    } = useAssistantsQuery();
    
    const { 
        data: conversationData,
        isLoading: conversationLoading,
        error: conversationError
    } = useConversationQuery(activeConversationId);
    
    const { 
        data: messagesData,
        isLoading: messagesLoading,
        error: messagesError
    } = useMessagesQuery(activeConversationId, { polling: true });
    
    const {
        data: conversationsData,
        isLoading: conversationsLoading,
        error: conversationsError,
        refetch: refetchConversations
    } = useConversationsQuery();
    
    // Mutations
    const createConversationMutation = useCreateConversationMutation({
        onMutate: () => {
            dispatch(startConversationCreation());
        },
        onSuccess: (data) => {
            const conversationId = data.conversationId || data.conversation_id;
            if (conversationId) {
                dispatch(setActiveConversation(conversationId));
                toast.success('Conversation created');
            }
        },
        onError: (error) => {
            dispatch(setUIError(error.message));
            toast.error('Failed to create conversation');
        }
    });
    
    const sendMessageMutation = useSendMessageMutation({
        onMutate: () => {
            dispatch(setTyping(true));
        },
        onSuccess: () => {
            dispatch(setTyping(false));
        },
        onError: (error) => {
            dispatch(setTyping(false));
            
            // Handle specific errors
            if (error instanceof ConversationNotFoundError) {
                // Clear active conversation and prompt to create new one
                dispatch(setActiveConversation(null));
                toast.error('Conversation expired. Please start a new one.');
            } else {
                dispatch(setUIError(error.message));
                toast.error('Failed to send message');
            }
        }
    });
    
    // Derived state
    const assistants = useMemo(() => 
        assistantsData?.assistants || [], 
        [assistantsData]
    );
    
    const selectedAssistant = useMemo(() => 
        assistants.find(a => a.id === selectedAssistantId) || null,
        [assistants, selectedAssistantId]
    );
    
    const conversation = useMemo(() => 
        conversationData?.conversation || null,
        [conversationData]
    );
    
    const messages = useMemo(() => 
        messagesData?.messages || [],
        [messagesData]
    );
    
    const conversations = useMemo(() =>
        conversationsData?.conversations || [],
        [conversationsData]
    );
    
    // Auto-select first assistant if none selected
    useEffect(() => {
        if (assistants.length > 0 && !selectedAssistantId) {
            const defaultAssistant = assistants.find(a => a.isDefault || a.name === 'AI Assistant') || assistants[0];
            dispatch(setSelectedAssistant(defaultAssistant.id));
        }
    }, [assistants, selectedAssistantId, dispatch]);
    
    // Actions
    const selectAssistant = useCallback((assistant) => {
        if (!assistant?.id) return;
        
        dispatch(setSelectedAssistant(assistant.id));
        dispatch(clearUIError());
        
        // Clear messages cache for old conversation
        if (activeConversationId) {
            queryClient.removeQueries({
                queryKey: aiQueryKeys.messages(activeConversationId)
            });
        }
    }, [dispatch, queryClient, activeConversationId]);
    
    const ensureConversation = useCallback(async () => {
        if (!selectedAssistantId) {
            throw new Error('No assistant selected');
        }
        
        if (activeConversationId && conversation) {
            return activeConversationId;
        }
        
        // Check if already creating
        if (isCreatingConversation || createConversationMutation.isPending) {
            // Wait for current creation to complete
            return new Promise((resolve, reject) => {
                const checkInterval = setInterval(() => {
                    if (!isCreatingConversation && !createConversationMutation.isPending) {
                        clearInterval(checkInterval);
                        if (activeConversationId) {
                            resolve(activeConversationId);
                        } else {
                            reject(new Error('Failed to create conversation'));
                        }
                    }
                }, 100);
                
                // Timeout after 10 seconds
                setTimeout(() => {
                    clearInterval(checkInterval);
                    reject(new Error('Conversation creation timeout'));
                }, 10000);
            });
        }
        
        // Create new conversation
        const result = await createConversationMutation.mutateAsync({
            assistantId: selectedAssistantId,
            context: {
                source: 'web_app',
                timestamp: new Date().toISOString()
            }
        });
        
        return result.conversationId || result.conversation_id;
    }, [
        selectedAssistantId, 
        activeConversationId, 
        conversation,
        isCreatingConversation,
        createConversationMutation
    ]);
    
    const sendMessage = useCallback(async (message, context = {}) => {
        if (!message?.trim()) {
            throw new Error('Message cannot be empty');
        }
        
        if (!selectedAssistantId) {
            throw new Error('No assistant selected');
        }
        
        try {
            const conversationId = await ensureConversation();
            
            // Ensure assistantId is in the context
            const enrichedContext = {
                ...context,
                assistantId: selectedAssistantId
            };
            
            const result = await sendMessageMutation.mutateAsync({
                conversationId,
                message,
                context: enrichedContext
            });
            
            return result;
        } catch (error) {
            // Re-throw to let mutation error handler deal with it
            throw error;
        }
    }, [ensureConversation, sendMessageMutation, selectedAssistantId]);
    
    const createNewConversation = useCallback(async () => {
        if (!selectedAssistantId) {
            toast.error('Please select an assistant first');
            return null;
        }
        
        // Clear current conversation
        dispatch(setActiveConversation(null));
        
        // Clear messages cache
        if (activeConversationId) {
            queryClient.removeQueries({
                queryKey: aiQueryKeys.messages(activeConversationId)
            });
        }
        
        // Create new
        return createConversationMutation.mutateAsync({
            assistantId: selectedAssistantId,
            context: {
                source: 'web_app',
                timestamp: new Date().toISOString()
            }
        });
    }, [
        selectedAssistantId, 
        activeConversationId,
        dispatch, 
        queryClient,
        createConversationMutation
    ]);
    
    const loadConversation = useCallback(async (conversationId) => {
        if (!conversationId) return;
        
        dispatch(setActiveConversation(conversationId));
        
        // Prefetch messages for the conversation
        queryClient.prefetchQuery({
            queryKey: aiQueryKeys.messages(conversationId),
            queryFn: () => messagesData
        });
    }, [dispatch, queryClient, messagesData]);
    
    const clearErrors = useCallback(() => {
        dispatch(clearUIError());
    }, [dispatch]);
    
    // Loading states
    const loading = {
        assistants: assistantsLoading,
        conversation: conversationLoading,
        messages: messagesLoading,
        conversations: conversationsLoading,
        sending: sendMessageMutation.isPending,
        creating: createConversationMutation.isPending
    };
    
    // Error states
    const errors = {
        assistants: assistantsError?.message || null,
        conversation: conversationError?.message || null,
        messages: messagesError?.message || null,
        conversations: conversationsError?.message || null,
        sending: sendMessageMutation.error?.message || null,
        creating: createConversationMutation.error?.message || null
    };
    
    // Combined loading state
    const isLoading = Object.values(loading).some(v => v);
    
    // Combined error state
    const hasErrors = Object.values(errors).some(v => v !== null);
    const latestError = Object.values(errors).find(v => v !== null) || null;
    
    return {
        // State
        assistants,
        selectedAssistant,
        selectedAssistantId,
        conversation,
        activeConversationId,
        messages,
        conversations,
        isTyping,
        viewMode,
        
        // Actions
        selectAssistant,
        sendMessage,
        createNewConversation,
        loadConversation,
        loadAssistants: refetchAssistants,
        loadConversations: refetchConversations,
        clearErrors,
        
        // Loading states
        loading,
        isLoading,
        
        // Error states
        errors,
        hasErrors,
        latestError,
        
        // Status flags
        isReady: !isLoading && !hasErrors && selectedAssistant && conversation,
        canSendMessage: !isTyping && !sendMessageMutation.isPending && activeConversationId,
        isInitialized: !assistantsLoading && assistants.length > 0
    };
};

// Export as default
export default useAIService;
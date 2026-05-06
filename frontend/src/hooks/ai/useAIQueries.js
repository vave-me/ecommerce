import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getAssistantService } from '../../services/ai/AssistantServiceWithRetry';

/**
 * React Query hooks for AI/Assistant functionality
 * Handles all server state management for AI features
 */

// Query key factory for consistent key generation
export const aiQueryKeys = {
    all: ['ai'],
    assistants: () => [...aiQueryKeys.all, 'assistants'],
    assistant: (id) => [...aiQueryKeys.assistants(), id],
    conversations: () => [...aiQueryKeys.all, 'conversations'],
    conversation: (id) => [...aiQueryKeys.conversations(), id],
    messages: (conversationId) => [...aiQueryKeys.all, 'messages', conversationId],
    stats: () => [...aiQueryKeys.all, 'stats']
};

/**
 * Hook to fetch all assistants
 */
export const useAssistantsQuery = (options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.assistants(),
        queryFn: async () => {
            const result = await service.getAssistants({
                page: options.page || 1,
                limit: options.limit || 20
            });
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        staleTime: 5 * 60 * 1000, // 5 minutes
        gcTime: 10 * 60 * 1000, // 10 minutes
        retry: 2,
        ...options
    });
};

/**
 * Hook to fetch a specific assistant
 */
export const useAssistantQuery = (assistantId, options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.assistant(assistantId),
        queryFn: async () => {
            const result = await service.getAssistant(assistantId);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        enabled: !!assistantId,
        staleTime: 5 * 60 * 1000,
        ...options
    });
};

/**
 * Hook to fetch user conversations
 */
export const useConversationsQuery = (options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.conversations(),
        queryFn: async () => {
            const result = await service.getUserConversations({
                activeOnly: options.activeOnly || false,
                page: options.page || 1,
                limit: options.limit || 20
            });
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        staleTime: 30 * 1000, // 30 seconds
        ...options
    });
};

/**
 * Hook to fetch a specific conversation
 */
export const useConversationQuery = (conversationId, options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.conversation(conversationId),
        queryFn: async () => {
            const result = await service.getConversation(conversationId);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        enabled: !!conversationId,
        staleTime: 30 * 1000,
        ...options
    });
};

/**
 * Hook to fetch messages for a conversation
 */
export const useMessagesQuery = (conversationId, options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.messages(conversationId),
        queryFn: async () => {
            const result = await service.getMessages(conversationId, {
                page: options.page || 1,
                limit: options.limit || 50
            });
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        enabled: !!conversationId,
        staleTime: 0, // Always fresh
        refetchInterval: options.polling ? 30000 : false, // Poll every 30s if enabled
        ...options
    });
};

/**
 * Hook to get conversation statistics
 */
export const useStatsQuery = (options = {}) => {
    const service = getAssistantService();
    
    return useQuery({
        queryKey: aiQueryKeys.stats(),
        queryFn: async () => {
            const result = await service.getStats();
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        staleTime: 60 * 1000, // 1 minute
        ...options
    });
};

/**
 * Mutation to create a conversation
 */
export const useCreateConversationMutation = (options = {}) => {
    const service = getAssistantService();
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async ({ assistantId, context }) => {
            const result = await service.createConversation(assistantId, context);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        onSuccess: (data) => {
            // Invalidate conversations list to include new conversation
            queryClient.invalidateQueries({ queryKey: aiQueryKeys.conversations() });
            
            // Optionally prefetch the new conversation
            if (data.conversationId) {
                queryClient.setQueryData(
                    aiQueryKeys.conversation(data.conversationId),
                    { conversation: data.conversation || { id: data.conversationId } }
                );
            }
        },
        ...options
    });
};

/**
 * Mutation to send a message
 */
export const useSendMessageMutation = (options = {}) => {
    const service = getAssistantService();
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async ({ conversationId, message, context }) => {
            const result = await service.sendMessage(conversationId, message, context);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        onMutate: async ({ conversationId, message }) => {
            // Cancel any outgoing refetches
            await queryClient.cancelQueries({
                queryKey: aiQueryKeys.messages(conversationId)
            });
            
            // Snapshot previous value
            const previousMessages = queryClient.getQueryData(
                aiQueryKeys.messages(conversationId)
            );
            
            // Optimistically update messages
            const optimisticMessage = {
                id: `temp-${Date.now()}`,
                conversationId,
                role: 'USER',
                content: message,
                timestamp: new Date().toISOString(),
                _optimistic: true
            };
            
            queryClient.setQueryData(
                aiQueryKeys.messages(conversationId),
                (old) => ({
                    ...old,
                    messages: [...(old?.messages || []), optimisticMessage]
                })
            );
            
            // Return context for rollback
            return { previousMessages, optimisticMessage };
        },
        onSuccess: (data, variables, context) => {
            const { conversationId } = variables;
            
            // Replace optimistic message with real messages
            queryClient.setQueryData(
                aiQueryKeys.messages(conversationId),
                (old) => ({
                    ...old,
                    messages: [
                        ...(old?.messages || []).filter(m => !m._optimistic),
                        data.userMessage,
                        data.assistantMessage
                    ]
                })
            );
            
            // Invalidate to ensure consistency
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.messages(conversationId)
            });
        },
        onError: (error, variables, context) => {
            const { conversationId } = variables;
            
            // Rollback on error
            if (context?.previousMessages) {
                queryClient.setQueryData(
                    aiQueryKeys.messages(conversationId),
                    context.previousMessages
                );
            }
        },
        ...options
    });
};

/**
 * Mutation to update assistant configuration
 */
export const useUpdateAssistantMutation = (options = {}) => {
    const service = getAssistantService();
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async ({ assistantId, config }) => {
            const result = await service.updateAssistantConfig(assistantId, config);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        onSuccess: (data, variables) => {
            const { assistantId } = variables;
            
            // Invalidate specific assistant
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.assistant(assistantId)
            });
            
            // Invalidate assistants list
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.assistants()
            });
        },
        ...options
    });
};

/**
 * Mutation to archive a conversation
 */
export const useArchiveConversationMutation = (options = {}) => {
    const service = getAssistantService();
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async (conversationId) => {
            const result = await service.archiveConversation(conversationId);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        onSuccess: (data, conversationId) => {
            // Invalidate conversations list
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.conversations()
            });
            
            // Invalidate specific conversation
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.conversation(conversationId)
            });
        },
        ...options
    });
};

/**
 * Mutation to delete a conversation
 */
export const useDeleteConversationMutation = (options = {}) => {
    const service = getAssistantService();
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async (conversationId) => {
            const result = await service.deleteConversation(conversationId);
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        onSuccess: (data, conversationId) => {
            // Remove from cache
            queryClient.removeQueries({
                queryKey: aiQueryKeys.conversation(conversationId)
            });
            queryClient.removeQueries({
                queryKey: aiQueryKeys.messages(conversationId)
            });
            
            // Invalidate conversations list
            queryClient.invalidateQueries({
                queryKey: aiQueryKeys.conversations()
            });
        },
        ...options
    });
};

/**
 * Utility to prefetch assistants
 */
export const prefetchAssistants = async (queryClient) => {
    const service = getAssistantService();
    
    await queryClient.prefetchQuery({
        queryKey: aiQueryKeys.assistants(),
        queryFn: async () => {
            const result = await service.getAssistants();
            if (!result.success) {
                throw new Error(result.error);
            }
            return result.data;
        },
        staleTime: 5 * 60 * 1000
    });
};

/**
 * Utility to invalidate all AI queries
 */
export const invalidateAllAIQueries = (queryClient) => {
    queryClient.invalidateQueries({ queryKey: aiQueryKeys.all });
};
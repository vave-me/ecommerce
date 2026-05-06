/**
 * Production-ready Assistant State Management
 * Eliminates all race conditions, duplicates, and sync issues
 */

import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useCallback, useMemo } from 'react';
import { nanoid } from 'nanoid';
import { getAssistantService } from '../../services/ai/AssistantService';

// Zustand store for UI state only
const useAssistantUIStore = create(
  devtools(
    (set) => ({
      // UI State
      selectedAssistantId: null,
      activeConversationId: null,
      viewMode: 'chat',
      inputValue: '',
      
      // Optimistic tracking
      pendingMessages: {},
      failedMessageIds: [],
      
      // Actions
      setSelectedAssistant: (id) => set({ selectedAssistantId: id }),
      setActiveConversation: (id) => set({ activeConversationId: id }),
      setViewMode: (mode) => set({ viewMode: mode }),
      setInputValue: (value) => set({ inputValue: value }),
      
      // Pending message management
      addPendingMessage: (tempId, message) => 
        set((state) => ({
          pendingMessages: { ...state.pendingMessages, [tempId]: message }
        })),
      
      removePendingMessage: (tempId) =>
        set((state) => {
          const { [tempId]: _, ...rest } = state.pendingMessages;
          return { pendingMessages: rest };
        }),
      
      markMessageFailed: (tempId) =>
        set((state) => ({
          failedMessageIds: [...state.failedMessageIds, tempId]
        })),
      
      clearFailedMessage: (tempId) =>
        set((state) => ({
          failedMessageIds: state.failedMessageIds.filter(id => id !== tempId)
        })),
      
      reset: () => set({
        selectedAssistantId: null,
        activeConversationId: null,
        viewMode: 'chat',
        inputValue: '',
        pendingMessages: {},
        failedMessageIds: []
      })
    }),
    { name: 'assistant-ui' }
  )
);

// Query keys factory
const queryKeys = {
  assistants: ['assistants'],
  assistant: (id) => ['assistant', id],
  conversations: ['conversations'],
  conversation: (id) => ['conversation', id],
  messages: (conversationId) => ['messages', conversationId],
};

// Request deduplication
class RequestDeduplicator {
  constructor() {
    this.pending = new Map();
  }

  async dedupe(key, factory) {
    if (this.pending.has(key)) {
      return this.pending.get(key);
    }

    const promise = factory()
      .finally(() => this.pending.delete(key));
    
    this.pending.set(key, promise);
    return promise;
  }
}

const deduplicator = new RequestDeduplicator();

/**
 * Main hook for assistant state management
 * Single source of truth with React Query + minimal UI state in Zustand
 */
export const useAssistantState = () => {
  const queryClient = useQueryClient();
  const assistantService = getAssistantService();
  const {
    selectedAssistantId,
    activeConversationId,
    viewMode,
    inputValue,
    pendingMessages,
    failedMessageIds,
    ...actions
  } = useAssistantUIStore();

  // Fetch assistants
  const assistantsQuery = useQuery({
    queryKey: queryKeys.assistants,
    queryFn: () => deduplicator.dedupe(
      'fetch-assistants',
      () => assistantService.getAssistants()
    ),
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000,
  });

  // Fetch active conversation
  const conversationQuery = useQuery({
    queryKey: queryKeys.conversation(activeConversationId),
    queryFn: () => assistantService.getConversation(activeConversationId),
    enabled: !!activeConversationId,
    staleTime: 30 * 1000,
  });

  // Fetch messages
  const messagesQuery = useQuery({
    queryKey: queryKeys.messages(activeConversationId),
    queryFn: () => assistantService.getMessages(activeConversationId),
    enabled: !!activeConversationId,
    staleTime: 0, // Always fresh
    refetchInterval: 30000, // Poll every 30s
  });

  // Create conversation mutation
  const createConversationMutation = useMutation({
    mutationFn: ({ assistantId, initialContext }) => 
      deduplicator.dedupe(
        `create-conversation-${assistantId}`,
        () => assistantService.createConversation(assistantId, initialContext)
      ),
    
    onSuccess: (response) => {
      const conversationId = response.data?.conversationId || response.data?.conversation_id;
      if (conversationId) {
        actions.setActiveConversation(conversationId);
        queryClient.invalidateQueries({ queryKey: queryKeys.conversations });
      }
    }
  });

  // Send message mutation with optimistic updates
  const sendMessageMutation = useMutation({
    mutationFn: async ({ conversationId, message, context }) => {
      return assistantService.sendMessage(
        conversationId,
        message,
        context || {}
      );
    },

    onMutate: async ({ conversationId, message }) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({
        queryKey: queryKeys.messages(conversationId)
      });

      // Snapshot previous value
      const previousMessages = queryClient.getQueryData(
        queryKeys.messages(conversationId)
      );

      // Create optimistic message
      const tempId = nanoid();
      const optimisticMessage = {
        id: tempId,
        tempId,
        conversationId,
        role: 'USER',
        content: message,
        timestamp: new Date().toISOString(),
        status: 'pending'
      };

      // Add to pending messages
      actions.addPendingMessage(tempId, optimisticMessage);

      // Update query data
      queryClient.setQueryData(
        queryKeys.messages(conversationId),
        (old) => {
          const messages = old?.data?.messages || [];
          return {
            ...old,
            data: {
              ...old?.data,
              messages: [...messages, optimisticMessage]
            }
          };
        }
      );

      return { previousMessages, tempId };
    },

    onSuccess: (response, variables, context) => {
      const { tempId } = context;
      
      // Remove pending message
      actions.removePendingMessage(tempId);
      
      // Update with real messages
      queryClient.setQueryData(
        queryKeys.messages(variables.conversationId),
        (old) => {
          const messages = old?.data?.messages || [];
          // Remove temp message
          const filtered = messages.filter(m => m.tempId !== tempId);
          
          // Add real messages
          return {
            ...old,
            data: {
              ...old?.data,
              messages: [
                ...filtered,
                // User message with real ID
                {
                  id: response.data?.userMessageId || nanoid(),
                  conversationId: variables.conversationId,
                  role: 'USER',
                  content: variables.message,
                  timestamp: new Date().toISOString(),
                  status: 'sent'
                },
                // Assistant response
                {
                  id: response.data?.messageId || response.data?.message_id,
                  conversationId: variables.conversationId,
                  role: 'ASSISTANT',
                  content: response.data?.response,
                  timestamp: new Date().toISOString(),
                  status: 'sent',
                  metadata: {
                    actions: response.data?.actions,
                    confidence: response.data?.confidence,
                    data: response.data?.data
                  }
                }
              ]
            }
          };
        }
      );

      // Clear input
      actions.setInputValue('');
    },

    onError: (error, variables, context) => {
      if (context?.tempId) {
        // Mark as failed
        actions.markMessageFailed(context.tempId);
        
        // Optionally rollback
        if (context.previousMessages) {
          queryClient.setQueryData(
            queryKeys.messages(variables.conversationId),
            context.previousMessages
          );
        }
      }
    }
  });

  // Retry failed message
  const retryMessage = useCallback((tempId) => {
    const message = pendingMessages[tempId];
    if (message && activeConversationId) {
      actions.clearFailedMessage(tempId);
      return sendMessageMutation.mutate({
        conversationId: activeConversationId,
        message: message.content,
        context: message.metadata?.context
      });
    }
  }, [pendingMessages, activeConversationId, sendMessageMutation, actions]);

  // Ensure conversation exists
  const ensureConversation = useCallback(async () => {
    if (!selectedAssistantId) {
      throw new Error('No assistant selected');
    }

    if (activeConversationId) {
      return activeConversationId;
    }

    // Check if we're already creating
    const key = `ensure-conversation-${selectedAssistantId}`;
    return deduplicator.dedupe(key, async () => {
      const response = await createConversationMutation.mutateAsync({
        assistantId: selectedAssistantId
      });
      return response.data?.conversationId || response.data?.conversation_id;
    });
  }, [selectedAssistantId, activeConversationId, createConversationMutation]);

  // Send message with conversation creation if needed
  const sendMessage = useCallback(async (message, context) => {
    const conversationId = await ensureConversation();
    return sendMessageMutation.mutate({
      conversationId,
      message,
      context
    });
  }, [ensureConversation, sendMessageMutation]);

  // Merge server messages with pending
  const allMessages = useMemo(() => {
    const serverMessages = messagesQuery.data?.data?.messages || [];
    const pending = Object.values(pendingMessages);
    const failed = new Set(failedMessageIds);

    // Combine and sort
    return [...serverMessages, ...pending]
      .filter(m => !failed.has(m.tempId))
      .sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));
  }, [messagesQuery.data, pendingMessages, failedMessageIds]);

  // Get selected assistant
  const selectedAssistant = useMemo(() => {
    const assistants = assistantsQuery.data?.data?.assistants || [];
    return assistants.find(a => a.id === selectedAssistantId);
  }, [assistantsQuery.data, selectedAssistantId]);

  return {
    // State
    assistants: assistantsQuery.data?.data?.assistants || [],
    selectedAssistant,
    conversation: conversationQuery.data?.data?.conversation,
    messages: allMessages,
    viewMode,
    inputValue,
    
    // Loading states
    isLoadingAssistants: assistantsQuery.isLoading,
    isLoadingMessages: messagesQuery.isLoading,
    isSending: sendMessageMutation.isPending,
    isCreatingConversation: createConversationMutation.isPending,
    
    // Errors
    assistantsError: assistantsQuery.error,
    messagesError: messagesQuery.error,
    sendError: sendMessageMutation.error,
    
    // Actions
    selectAssistant: actions.setSelectedAssistant,
    setViewMode: actions.setViewMode,
    setInputValue: actions.setInputValue,
    sendMessage,
    retryMessage,
    createConversation: createConversationMutation.mutate,
    reset: actions.reset,
    
    // Utilities
    hasFailedMessages: failedMessageIds.length > 0,
    pendingMessageCount: Object.keys(pendingMessages).length,
  };
};

// Singleton hook for conversation creation
let conversationCreationPromise = null;

export const useEnsureConversation = (assistantId) => {
  const queryClient = useQueryClient();
  const { setActiveConversation } = useAssistantUIStore();
  const assistantService = getAssistantService();

  return useCallback(async () => {
    if (!assistantId) return null;

    // Check cache first
    const conversations = queryClient.getQueryData(queryKeys.conversations);
    const existing = conversations?.data?.conversations?.find(
      c => c.assistantId === assistantId && c.active
    );

    if (existing) {
      setActiveConversation(existing.id);
      return existing.id;
    }

    // Prevent duplicate creation
    if (!conversationCreationPromise) {
      conversationCreationPromise = assistantService.createConversation(
        assistantId,
        { source: 'web_app' }
      ).then(response => {
        const id = response.data?.conversationId || response.data?.conversation_id;
        setActiveConversation(id);
        queryClient.invalidateQueries({ queryKey: queryKeys.conversations });
        conversationCreationPromise = null;
        return id;
      }).catch(error => {
        conversationCreationPromise = null;
        throw error;
      });
    }

    return conversationCreationPromise;
  }, [assistantId, queryClient, setActiveConversation, assistantService]);
};
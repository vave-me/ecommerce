import { useState, useCallback, useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { useAuth } from '../context/AuthContext';
import * as supportApi from '../api/client/supportApi';

export const useSupport = () => {
    const { user } = useAuth();
    const queryClient = useQueryClient();
    const [activeChannelId, setActiveChannelId] = useState(null);

    // Get or create support channel for user
    const { data: channels, isLoading: channelsLoading } = useQuery({
        queryKey: ['support-channels', user?.userId],
        queryFn: () => supportApi.getUserSupportChannels(user?.userId, { active_only: true }),
        enabled: !!user?.userId,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    // Auto-select first active channel or create one
    useEffect(() => {
        if (channels?.channels?.length > 0 && !activeChannelId) {
            setActiveChannelId(channels.channels[0].id);
        }
    }, [channels, activeChannelId]);

    // Create support channel mutation
    const createChannelMutation = useMutation({
        mutationFn: (channelData) => supportApi.createSupportChannel({
            user_id: user?.userId,
            ...channelData
        }),
        onSuccess: (data) => {
            queryClient.invalidateQueries(['support-channels']);
            setActiveChannelId(data.id);
            toast.success('Support channel created successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to create support channel');
        }
    });

    // Get tickets for active channel
    const { data: ticketsData, isLoading: ticketsLoading, refetch: refetchTickets } = useQuery({
        queryKey: ['support-tickets', activeChannelId],
        queryFn: () => supportApi.getChannelTickets(activeChannelId, {
            sort_by: 'created_at',
            descending: true
        }),
        enabled: !!activeChannelId,
        staleTime: 30 * 1000, // 30 seconds
    });

    // Create ticket mutation
    const createTicketMutation = useMutation({
        mutationFn: (ticketData) => supportApi.createTicket({
            channel_id: activeChannelId,
            ...ticketData
        }),
        onSuccess: () => {
            queryClient.invalidateQueries(['support-tickets', activeChannelId]);
            toast.success('Ticket created successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to create ticket');
        }
    });

    // Update ticket mutation
    const updateTicketMutation = useMutation({
        mutationFn: ({ ticketId, updates }) => supportApi.updateTicket(ticketId, updates),
        onSuccess: () => {
            queryClient.invalidateQueries(['support-tickets']);
            toast.success('Ticket updated successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to update ticket');
        }
    });

    // Add reply mutation
    const addReplyMutation = useMutation({
        mutationFn: ({ ticketId, content, attachments = [], is_public = true }) => 
            supportApi.addTicketReply(ticketId, {
                author_id: user?.userId,
                author_type: 'CUSTOMER',
                content,
                attachments,
                is_public
            }),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries(['ticket-communications', variables.ticketId]);
            queryClient.invalidateQueries(['support-tickets']);
            toast.success('Reply added successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to add reply');
        }
    });

    // Resolve ticket mutation
    const resolveTicketMutation = useMutation({
        mutationFn: ({ ticketId, resolution, applied_solutions = [] }) => 
            supportApi.resolveTicket(ticketId, { resolution, applied_solutions }),
        onSuccess: () => {
            queryClient.invalidateQueries(['support-tickets']);
            toast.success('Ticket resolved successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to resolve ticket');
        }
    });

    // Close ticket mutation
    const closeTicketMutation = useMutation({
        mutationFn: ({ ticketId, closure_notes, satisfaction }) => 
            supportApi.closeTicket(ticketId, { closure_notes, satisfaction }),
        onSuccess: () => {
            queryClient.invalidateQueries(['support-tickets']);
            toast.success('Ticket closed successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to close ticket');
        }
    });

    // Reopen ticket mutation
    const reopenTicketMutation = useMutation({
        mutationFn: ({ ticketId, reason }) => supportApi.reopenTicket(ticketId, reason),
        onSuccess: () => {
            queryClient.invalidateQueries(['support-tickets']);
            toast.success('Ticket reopened successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to reopen ticket');
        }
    });

    // Get ticket communications
    const getTicketCommunications = useCallback(async (ticketId, includeInternal = false) => {
        try {
            return await supportApi.getTicketCommunications(ticketId, { 
                include_internal: includeInternal 
            });
        } catch (error) {
            toast.error('Failed to load ticket communications');
            throw error;
        }
    }, []);

    // Search knowledge base
    const searchKnowledgeBase = useCallback(async (query, categories = []) => {
        try {
            return await supportApi.searchKnowledgeBase(query, { categories });
        } catch (error) {
            toast.error('Failed to search knowledge base');
            throw error;
        }
    }, []);

    // Initialize support (create channel if needed)
    const initializeSupport = useCallback(async (channelType = 'GENERAL') => {
        if (!user?.userId) {
            toast.error('Please login to access support');
            return;
        }

        if (activeChannelId) {
            return activeChannelId;
        }

        if (channels?.channels?.length === 0) {
            const result = await createChannelMutation.mutateAsync({
                channel_type: channelType,
                settings: {
                    email_notifications: true,
                    auto_assign_tickets: true,
                    preferred_language: 'en',
                    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
                }
            });
            return result.id;
        }

        return channels?.channels?.[0]?.id;
    }, [user, activeChannelId, channels, createChannelMutation]);

    return {
        // State
        activeChannelId,
        channels: channels?.channels || [],
        tickets: ticketsData?.tickets || [],
        totalTickets: ticketsData?.total_count || 0,
        
        // Loading states
        isLoading: channelsLoading || ticketsLoading,
        channelsLoading,
        ticketsLoading,
        
        // Actions
        initializeSupport,
        createChannel: createChannelMutation.mutate,
        createTicket: createTicketMutation.mutate,
        updateTicket: updateTicketMutation.mutate,
        addReply: addReplyMutation.mutate,
        resolveTicket: resolveTicketMutation.mutate,
        closeTicket: closeTicketMutation.mutate,
        reopenTicket: reopenTicketMutation.mutate,
        refetchTickets,
        getTicketCommunications,
        searchKnowledgeBase,
        
        // Mutation states
        isCreatingChannel: createChannelMutation.isPending,
        isCreatingTicket: createTicketMutation.isPending,
        isUpdatingTicket: updateTicketMutation.isPending,
        isAddingReply: addReplyMutation.isPending,
        isResolvingTicket: resolveTicketMutation.isPending,
        isClosingTicket: closeTicketMutation.isPending,
        isReopeningTicket: reopenTicketMutation.isPending,
        
        // Utilities
        setActiveChannelId,
        getTicketById: (ticketId) => (ticketsData?.tickets || []).find(t => t.id === ticketId),
        getOpenTickets: () => (ticketsData?.tickets || []).filter(t => ['SUBMITTED', 'ASSIGNED', 'IN_PROGRESS', 'PENDING_CUSTOMER'].includes(t.status)),
        getResolvedTickets: () => (ticketsData?.tickets || []).filter(t => ['RESOLVED', 'CLOSED'].includes(t.status)),
    };
};

// Hook for admin support management
export const useAdminSupport = () => {
    const queryClient = useQueryClient();

    // Get all tickets with filters
    const getAllTickets = useCallback(async (filters = {}) => {
        try {
            // This would need an admin API endpoint
            const response = await fetch('/api/admin/support/tickets', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(filters)
            });
            return response.json();
        } catch (error) {
            toast.error('Failed to fetch tickets');
            throw error;
        }
    }, []);

    // Assign ticket to agent
    const assignTicketMutation = useMutation({
        mutationFn: ({ ticketId, assigneeId, assigneeType = 'HUMAN_AGENT', reason }) => 
            supportApi.assignTicket(ticketId, {
                assignee_id: assigneeId,
                assignee_type: assigneeType,
                assignment_reason: reason
            }),
        onSuccess: () => {
            queryClient.invalidateQueries(['admin-tickets']);
            toast.success('Ticket assigned successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to assign ticket');
        }
    });

    // Update ticket priority
    const updatePriorityMutation = useMutation({
        mutationFn: ({ ticketId, priority, reason }) => 
            supportApi.updateTicketPriority(ticketId, priority, reason),
        onSuccess: () => {
            queryClient.invalidateQueries(['admin-tickets']);
            toast.success('Priority updated successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to update priority');
        }
    });

    // Escalate ticket
    const escalateTicketMutation = useMutation({
        mutationFn: ({ ticketId, tier, reason, notes }) => 
            supportApi.escalateTicket(ticketId, {
                escalation_tier: tier,
                escalation_reason: reason,
                escalation_notes: notes
            }),
        onSuccess: () => {
            queryClient.invalidateQueries(['admin-tickets']);
            toast.success('Ticket escalated successfully');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to escalate ticket');
        }
    });

    // Add internal note
    const addInternalNoteMutation = useMutation({
        mutationFn: ({ ticketId, content, mentionedUsers = [] }) => 
            supportApi.addInternalNote(ticketId, {
                author_id: 'admin', // Should use actual admin ID
                content,
                mentioned_users: mentionedUsers
            }),
        onSuccess: () => {
            queryClient.invalidateQueries(['ticket-communications']);
            toast.success('Internal note added');
        },
        onError: (error) => {
            toast.error(error.message || 'Failed to add note');
        }
    });

    return {
        getAllTickets,
        assignTicket: assignTicketMutation.mutate,
        updatePriority: updatePriorityMutation.mutate,
        escalateTicket: escalateTicketMutation.mutate,
        addInternalNote: addInternalNoteMutation.mutate,
        
        // Loading states
        isAssigning: assignTicketMutation.isPending,
        isUpdatingPriority: updatePriorityMutation.isPending,
        isEscalating: escalateTicketMutation.isPending,
        isAddingNote: addInternalNoteMutation.isPending,
    };
};
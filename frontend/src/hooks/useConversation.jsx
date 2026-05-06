// src/hooks/useConversation.jsx
import {useQuery, useMutation, useQueryClient} from '@tanstack/react-query';
import {
    getConversationByRecipientAndProduct,
    startConversation,
} from '../api/client/messagingApi';
import {useAuth} from "../context/AuthContext";
export function useConversation(recipientId, itemId) {
    const queryClient = useQueryClient();
    const user = useAuth()
    const cacheKey = ['conversation', user.userId, recipientId, itemId];
    const {
        data: conversationList,
        isLoading,
        isError,
        error,
    } = useQuery({
        queryKey: cacheKey,
        queryFn: async () => {
            return await getConversationByRecipientAndProduct(recipientId, itemId);
        },
        enabled: !!itemId && !!user.userId && !!recipientId,
        // staleTime, cacheTime, etc., can be specified here as needed
    });
    const {mutateAsync: createConversation} = useMutation({
        mutationFn: async () => {
            const newConvo = await startConversation(user.userId, recipientId, itemId);
            return newConvo;
        },
        onSuccess: (newConvo) => {
            queryClient.setQueryData(cacheKey, [newConvo]);
        },
        onError: (err) => {
        },
    });
    // A derived convenience method: get a single conversation ID or create if none
    const ensureConversationId = async () => {
        // read from cache
        const existingData = queryClient.getQueryData(cacheKey);
        if (existingData && existingData.length > 0) {
            return existingData[0].id;
        }
        // Otherwise, create new
        const newConvo = await createConversation();
        if (newConvo?.id) {
            return newConvo.id;
        } else {
            throw new Error('Conversation creation failed (no id returned)');
        }
    };
    // Additional logs to see if there's an error
    if (isError) {
    }
    return {
        conversationList,
        isLoading,
        isError,
        ensureConversationId, // call this to get or create conversationId
    };
}

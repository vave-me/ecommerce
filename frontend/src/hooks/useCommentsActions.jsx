// src/comments/hooks/useCommentsActions.jsx
import {useMutation, useQueryClient} from '@tanstack/react-query';
import {useCallback} from 'react';
import {v4 as uuidv4} from 'uuid';
import {useNATS} from "../context/NATSContext";
import {commentspb as comments_api} from "../generated_proto/comments_api_pb";
import {message_type as message_types} from "../generated_proto/message_types_pb";
import {jetstream as message_api} from "../generated_proto/message_api_pb";
/**
 * Recursively nest newComment under the correct parent.
 */
export function addCommentToCache(allComments, newComment) {
    if (!newComment.parentId) {
        // top-level => push to root
        return [...allComments, newComment];
    }
    return allComments.map((comment) => {
        if (comment.id === newComment.parentId) {
            return {
                ...comment,
                replies: [...(comment.replies || []), newComment],
            };
        }
        if (comment.replies && comment.replies.length > 0) {
            return {
                ...comment,
                replies: addCommentToCache(comment.replies, newComment),
            };
        }
        return comment;
    });
}
/**
 * useCommentsActions
 * Returns createComment to post a new comment (or reply).
 * The difference is parentId == "" => top-level, or some string => nested.
 */
export function useCommentsActions(itemId, userId, categoryId) {
    const queryClient = useQueryClient();
    const {isConnected, publish} = useNATS();
    const natsName = process.env.NEXT_NATS_CM_NAME || "comments.AddComment";
    const mutation = useMutation({
        mutationFn: async ({content, parentId}) => {
            if (!isConnected) throw new Error('Not connected to NATS');
            const newComment = {
                id: uuidv4(),
                senderId: userId || 'anonymous',
                itemId,
                content,
                categoryId,
                parentId, // e.g. "" if top-level, or actual parent's ID
                createdAt: new Date().toISOString(),
                replies: [],
            };
            // Build Protobuf for NATS
            const addCommentPb = comments_api.AddComment.create({
                id: newComment.id,
                senderId: newComment.senderId,
                itemId: newComment.itemId,
                content: newComment.content,
                categoryId: newComment.categoryId,
                // If parentId is null or an empty string, we store "" for top-level
                parentId: parentId || '',
                createdAt: newComment.createdAt,
            });
            const addCommentBytes = comments_api.AddComment.encode(addCommentPb).finish();
            const wsCommand = message_types.WebsocketMessageData.create({
                payload: addCommentBytes,
                occurred_at: {
                    seconds: Math.floor(Date.now() / 1000),
                    nanos: 0,
                },
            });
            const encCommand = message_types.WebsocketMessageData.encode(wsCommand).finish();
            const streamMessage = message_api.StreamMessage.create({
                id: uuidv4(),
                name: natsName,
                data: encCommand,
                metadata: {user: userId || 'anonymous', role: 'sender'},
                sent_at: {
                    seconds: Math.floor(Date.now() / 1000),
                    nanos: 0,
                },
            });
            const enc = message_api.StreamMessage.encode(streamMessage).finish();
            const subject = `${natsName}.${itemId}`;
            await publish(subject, enc);
            return newComment;
        },
        onMutate: async (variables) => {
            await queryClient.cancelQueries(['comments', itemId]);
            const newComment = {
                id: uuidv4(),
                senderId: userId || 'anonymous',
                itemId,
                content: variables.content,
                categoryId,
                parentId: variables.parentId,
                createdAt: new Date().toISOString(),
                replies: [],
            };
            const previous = queryClient.getQueryData(['comments', itemId]) || [];
            // Optimistic update
            queryClient.setQueryData(['comments', itemId], (old = []) =>
                addCommentToCache(old, newComment)
            );
            return {previous};
        },
        onError: (err, variables, context) => {
            if (context?.previous) {
                queryClient.setQueryData(['comments', itemId], context.previous);
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries(['comments', itemId]);
        },
    });
    const createComment = useCallback(
        (content, parentId = '') => {
            // If you prefer null for top-level, that's okay,
            // but here we pass "" so the subscription side sees it as top-level
            mutation.mutate({content, parentId});
        },
        [mutation]
    );
    return {
        createComment,
        isSubmitting: mutation.isLoading,
        error: mutation.error,
    };
}

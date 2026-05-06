import { useCallback } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { useNatsContext } from '../core/NatsProvider';
import { MessageEncoder } from '../core/MessageEncoder';

export interface CommentPayload {
  id: string;
  senderId: string;
  itemId: string;
  content: string;
  categoryId?: string;
  parentId?: string;
  createdAt: string;
}

export interface UseCommentsActionsOptions {
  itemId: string;
  userId: string;
  categoryId?: string;
  encoder?: MessageEncoder;
  protobufTypes?: {
    AddComment: any;
    WebsocketMessageData: any;
    StreamMessage: any;
  };
  natsStreamName?: string;
  onError?: (error: Error) => void;
  onOptimisticUpdate?: (comment: CommentPayload) => void;
}

export function useCommentsActions(options: UseCommentsActionsOptions) {
  const {
    itemId,
    userId,
    categoryId,
    protobufTypes,
    natsStreamName = 'comments.AddComment',
    onError,
    onOptimisticUpdate
  } = options;

  const { isConnected, publish } = useNatsContext();

  const createComment = useCallback(
    async (content: string, parentId?: string) => {
      if (!isConnected) {
        const error = new Error('Not connected to NATS');
        onError?.(error);
        throw error;
      }

      const newComment: CommentPayload = {
        id: uuidv4(),
        senderId: userId || 'anonymous',
        itemId,
        content,
        categoryId,
        parentId: parentId || '',
        createdAt: new Date().toISOString(),
      };

      try {
        let encodedMessage: Uint8Array;

        if (protobufTypes) {
          // Use protobuf encoding
          const { AddComment, WebsocketMessageData, StreamMessage } = protobufTypes;
          
          // Create AddComment protobuf
          const addCommentPb = AddComment.create({
            id: newComment.id,
            senderId: newComment.senderId,
            itemId: newComment.itemId,
            content: newComment.content,
            categoryId: newComment.categoryId,
            parentId: newComment.parentId,
            createdAt: newComment.createdAt,
          });
          
          const addCommentBytes = AddComment.encode(addCommentPb).finish();
          
          // Wrap in WebsocketMessageData
          const wsCommand = WebsocketMessageData.create({
            payload: addCommentBytes,
            occurred_at: {
              seconds: Math.floor(Date.now() / 1000),
              nanos: 0,
            },
          });
          
          const encCommand = WebsocketMessageData.encode(wsCommand).finish();
          
          // Wrap in StreamMessage
          const streamMessage = StreamMessage.create({
            id: uuidv4(),
            name: natsStreamName,
            data: encCommand,
            metadata: { user: userId || 'anonymous', role: 'sender' },
            sent_at: {
              seconds: Math.floor(Date.now() / 1000),
              nanos: 0,
            },
          });
          
          encodedMessage = StreamMessage.encode(streamMessage).finish();
        } else {
          // Fallback to JSON encoding
          encodedMessage = new TextEncoder().encode(JSON.stringify(newComment));
        }

        const subject = `${natsStreamName}.${itemId}`;
        await publish(subject, encodedMessage);

        // Trigger optimistic update
        onOptimisticUpdate?.(newComment);

        return newComment;
      } catch (error) {
        onError?.(error as Error);
        throw error;
      }
    },
    [isConnected, publish, userId, itemId, categoryId, protobufTypes, natsStreamName, onError, onOptimisticUpdate]
  );

  return {
    createComment,
    isConnected
  };
}
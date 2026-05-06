import { useState, useEffect, useCallback, useRef } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { useNatsContext } from '../core/NatsProvider';
import { Subscription } from 'nats.ws';

export interface ChatMessage {
  id: string;
  text: string;
  senderId: string;
  recipientId: string;
  conversationId: string;
  itemId?: string;
  createdAt: number;
  isUserMessage: boolean;
}

export interface UseChatHistoryOptions {
  conversationId: string;
  userId: string;
  recipientId: string;
  itemId?: string;
  metadata?: Record<string, any>;
  protobufTypes?: {
    SendMessage: any;
    WebsocketMessageData: any;
    StreamMessage: any;
  };
  natsStreamName?: string;
  fetchHistory?: (conversationId: string) => Promise<any[]>;
  onError?: (error: Error) => void;
}

export function useChatHistory(options: UseChatHistoryOptions) {
  const {
    conversationId,
    userId,
    recipientId,
    itemId,
    metadata,
    protobufTypes,
    natsStreamName = 'messenger.SendMessage',
    fetchHistory,
    onError
  } = options;

  const { isConnected, publish, subscribe } = useNatsContext();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const unsubscribeRef = useRef<Subscription | null>(null);

  // Load historical messages
  useEffect(() => {
    if (!conversationId || !userId || !fetchHistory) return;

    let isMounted = true;

    const loadHistory = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const fetched = await fetchHistory(conversationId);
        
        // Transform messages
        const mapped = fetched.map((m) => ({
          id: m.id,
          text: m.body,
          senderId: m.senderId,
          recipientId: m.recipientId,
          conversationId: m.conversationId,
          itemId: m.itemId,
          createdAt: m.createdAt || Date.now(),
          isUserMessage: m.senderId === userId,
        }));

        if (isMounted) {
          setMessages(mapped);
        }
      } catch (err) {
        if (isMounted) {
          const error = err as Error;
          setError(error);
          onError?.(error);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    loadHistory();

    return () => {
      isMounted = false;
    };
  }, [conversationId, userId, fetchHistory, onError]);

  // Subscribe for real-time updates
  useEffect(() => {
    if (!conversationId || !isConnected || !userId) return;

    let isMounted = true;

    // Unsubscribe old subscription
    if (unsubscribeRef.current) {
      unsubscribeRef.current.unsubscribe();
      unsubscribeRef.current = null;
    }

    const subject = `${natsStreamName}.${conversationId}`;

    (async () => {
      try {
        const unsub = await subscribe(subject, (rawBytes) => {
          if (!isMounted) return;

          try {
            let newMessage: ChatMessage;

            if (protobufTypes) {
              // Decode protobuf messages
              const { StreamMessage, WebsocketMessageData, SendMessage } = protobufTypes;
              
              // Decode outer StreamMessage
              const decodedStreamMessage = StreamMessage.decode(rawBytes);
              
              // Decode WebsocketMessageData
              const wsData = WebsocketMessageData.decode(decodedStreamMessage.data);
              
              // Decode final SendMessage
              const msg = SendMessage.decode(wsData.payload);
              
              newMessage = {
                id: msg.id,
                text: msg.body,
                senderId: msg.senderId,
                recipientId: msg.recipientId,
                itemId: msg.itemId,
                conversationId: msg.conversationId,
                createdAt: Date.now(),
                isUserMessage: msg.senderId === userId,
              };
            } else {
              // Fallback to JSON decoding
              const decoded = JSON.parse(new TextDecoder().decode(rawBytes));
              newMessage = {
                ...decoded,
                isUserMessage: decoded.senderId === userId
              };
            }

            // Check for duplicates
            setMessages((prev) => {
              if (prev.some((m) => m.id === newMessage.id)) {
                return prev;
              }
              return [...prev, newMessage];
            });
          } catch (decodingErr) {
            console.error('Error decoding message:', decodingErr);
          }
        });

        unsubscribeRef.current = unsub;
      } catch (err) {
        const error = err as Error;
        setError(error);
        onError?.(error);
      }
    })();

    return () => {
      isMounted = false;
      if (unsubscribeRef.current) {
        unsubscribeRef.current.unsubscribe();
        unsubscribeRef.current = null;
      }
    };
  }, [conversationId, isConnected, subscribe, userId, protobufTypes, natsStreamName, onError]);

  // Send message
  const sendMessage = useCallback(
    async (text: string) => {
      if (!isConnected) {
        const err = new Error('Not connected to NATS');
        setError(err);
        onError?.(err);
        return;
      }

      if (!conversationId || !userId || !recipientId) {
        const err = new Error('Missing conversationId / userId / recipientId');
        setError(err);
        onError?.(err);
        return;
      }

      const trimmed = text.trim();
      if (!trimmed) return;

      try {
        const msgId = uuidv4();
        let encodedMessage: Uint8Array;

        if (protobufTypes) {
          // Build protobuf message
          const { SendMessage, WebsocketMessageData, StreamMessage } = protobufTypes;
          
          const messagePayload = SendMessage.create({
            id: msgId,
            conversationId: conversationId,
            senderId: userId,
            recipientId: recipientId,
            itemId: itemId || '',
            body: trimmed,
            isRead: false,
          });
          
          const msgBytes = SendMessage.encode(messagePayload).finish();
          
          // Wrap with WebsocketMessageData
          const wsCommand = WebsocketMessageData.create({
            payload: msgBytes,
            occurred_at: { seconds: Math.floor(Date.now() / 1000), nanos: 0 },
          });
          
          const encCommand = WebsocketMessageData.encode(wsCommand).finish();
          
          // Build final StreamMessage
          const streamMessage = StreamMessage.create({
            id: uuidv4(),
            name: natsStreamName,
            data: encCommand,
            metadata: metadata || { user: userId, role: 'sender' },
            sent_at: {
              seconds: Math.floor(Date.now() / 1000),
              nanos: 0,
            },
          });
          
          encodedMessage = StreamMessage.encode(streamMessage).finish();
        } else {
          // Fallback to JSON encoding
          const message = {
            id: msgId,
            conversationId,
            senderId: userId,
            recipientId,
            itemId: itemId || '',
            body: trimmed,
            isRead: false,
            createdAt: Date.now()
          };
          encodedMessage = new TextEncoder().encode(JSON.stringify(message));
        }

        // Publish message
        const subject = `${natsStreamName}.${conversationId}`;
        await publish(subject, encodedMessage);

        // Optimistic update
        const newLocalMsg: ChatMessage = {
          id: msgId,
          text: trimmed,
          senderId: userId,
          recipientId,
          conversationId,
          itemId,
          createdAt: Date.now(),
          isUserMessage: true,
        };

        // Check for duplicates
        setMessages((prev) => {
          if (prev.some((m) => m.id === newLocalMsg.id)) {
            return prev;
          }
          return [...prev, newLocalMsg];
        });
      } catch (err) {
        const error = err as Error;
        setError(error);
        onError?.(error);
      }
    },
    [conversationId, userId, recipientId, itemId, publish, isConnected, metadata, protobufTypes, natsStreamName, onError]
  );

  return {
    messages,
    isLoading,
    error,
    sendMessage,
  };
}
import { useState, useEffect, useCallback, useRef } from 'react';
import { Subscription } from 'nats.ws';
import { useNatsContext } from '../core/NatsProvider';
import { MessageType, UseSubscriptionResult, SubscriptionOptions } from '../core/types';
import { MessageEncoder } from '../core/MessageEncoder';

export function useSubscription<T>(
  messageType: MessageType<T>,
  options?: SubscriptionOptions & { encoder?: MessageEncoder }
): UseSubscriptionResult<T> {
  const [messages, setMessages] = useState<T[]>([]);
  const [latestMessage, setLatestMessage] = useState<T | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isSubscribed, setIsSubscribed] = useState(false);
  
  const { subscribe, isConnected } = useNatsContext();
  const subscriptionRef = useRef<Subscription | null>(null);
  const messageIdsRef = useRef<Set<string>>(new Set());
  const encoder = options?.encoder || new MessageEncoder();

  const clear = useCallback(() => {
    setMessages([]);
    setLatestMessage(null);
    messageIdsRef.current.clear();
  }, []);

  useEffect(() => {
    if (!isConnected) {
      setIsSubscribed(false);
      return;
    }

    let cancelled = false;

    const setupSubscription = async () => {
      try {
        const sub = await subscribe(messageType.subject, async (data) => {
          if (cancelled) return;

          try {
            const decoded = encoder.decode(messageType, data);
            
            // Handle deduplication if enabled
            if (options?.deduplicate) {
              const messageId = (decoded as any).id || JSON.stringify(decoded);
              if (messageIdsRef.current.has(messageId)) {
                return;
              }
              
              messageIdsRef.current.add(messageId);
              
              // Clean up old message IDs
              if (messageIdsRef.current.size > (options.deduplicationWindow || 100)) {
                const ids = Array.from(messageIdsRef.current);
                ids.slice(0, ids.length - (options.deduplicationWindow || 100))
                  .forEach(id => messageIdsRef.current.delete(id));
              }
            }

            setMessages(prev => [...prev, decoded]);
            setLatestMessage(decoded);
          } catch (err) {
            const error = err as Error;
            setError(error);
            options?.onError?.(error);
          }
        });

        subscriptionRef.current = sub;
        setIsSubscribed(true);
        setError(null);
      } catch (err) {
        const error = err as Error;
        setError(error);
        setIsSubscribed(false);
        options?.onError?.(error);
      }
    };

    setupSubscription();

    return () => {
      cancelled = true;
      if (subscriptionRef.current) {
        subscriptionRef.current.unsubscribe();
        subscriptionRef.current = null;
      }
      setIsSubscribed(false);
    };
  }, [isConnected, messageType.subject, subscribe, encoder, options]);

  return {
    messages,
    latestMessage,
    error,
    isSubscribed,
    clear
  };
}
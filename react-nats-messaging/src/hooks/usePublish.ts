import { useState, useCallback } from 'react';
import { useNatsContext } from '../core/NatsProvider';
import { MessageType, UsePublishResult, PublishOptions } from '../core/types';
import { MessageEncoder } from '../core/MessageEncoder';

export function usePublish<T>(
  messageType: MessageType<T>,
  options?: PublishOptions & { encoder?: MessageEncoder }
): UsePublishResult<T> {
  const [isPublishing, setIsPublishing] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [lastPublished, setLastPublished] = useState<T | null>(null);
  
  const { publish: natsPublish, isConnected } = useNatsContext();
  const encoder = options?.encoder || new MessageEncoder();

  const publish = useCallback(async (data: T) => {
    if (!isConnected) {
      const err = new Error('Not connected to NATS');
      setError(err);
      throw err;
    }

    setIsPublishing(true);
    setError(null);

    try {
      const encoded = encoder.encode(messageType, data);
      await natsPublish(messageType.subject, encoded, options);
      setLastPublished(data);
    } catch (err) {
      const error = err as Error;
      setError(error);
      throw error;
    } finally {
      setIsPublishing(false);
    }
  }, [isConnected, messageType, natsPublish, encoder, options]);

  return {
    publish,
    isPublishing,
    error,
    lastPublished
  };
}
import {useCallback, useRef} from 'react';
import {MessageType} from '../core/types';
import {MessageEncoder} from '../core/MessageEncoder';

interface MessageHandlerOptions<T> {
    onMessage: (message: T) => void | Promise<void>;
    onError?: (error: Error) => void;
    deduplicate?: boolean;
    deduplicationKey?: (message: T) => string;
    encoder?: MessageEncoder;
}

export function useMessageHandler<T>(
    messageType: MessageType<T>,
    options: MessageHandlerOptions<T>
) {
    const processedIds = useRef<Set<string>>(new Set());
    const encoder = options.encoder || new MessageEncoder();

    const handleMessage = useCallback(async (data: Uint8Array) => {
        try {
            const decoded = encoder.decode(messageType, data);

            // Handle deduplication
            if (options.deduplicate) {
                const key = options.deduplicationKey
                    ? options.deduplicationKey(decoded)
                    : (decoded as any).id || JSON.stringify(decoded);

                if (processedIds.current.has(key)) {
                    return;
                }

                processedIds.current.add(key);

                // Limit the size of processed IDs
                if (processedIds.current.size > 1000) {
                    const firstKey = processedIds.current.values().next().value;
                    if (firstKey !== undefined) {
                        processedIds.current.delete(firstKey);
                    }
                }
            }

            await options.onMessage(decoded);
        } catch (err) {
            options.onError?.(err as Error);
        }
    }, [messageType, encoder, options]);

    const clearProcessedIds = useCallback(() => {
        processedIds.current.clear();
    }, []);

    return {
        handleMessage,
        clearProcessedIds
    };
}
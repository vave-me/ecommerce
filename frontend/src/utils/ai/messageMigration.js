import { MessageRole, MessageStatus, createMessage } from '../../types/ai';

/**
 * Message Migration Utilities
 * Handles conversion from legacy message formats to standardized format
 */

/**
 * Detects the format of a message or conversation entry
 * @param {*} data - Data to check
 * @returns {'standard' | 'legacy_pair' | 'backend' | 'unknown'}
 */
export function detectMessageFormat(data) {
    if (!data || typeof data !== 'object') {
        return 'unknown';
    }
    
    // Standard format: has role property
    if (data.role && data.content) {
        return 'standard';
    }
    
    // Legacy pair format: has prompt and response
    if (data.prompt && data.response) {
        return 'legacy_pair';
    }
    
    // Backend format: has actions_taken instead of actions
    if (data.role && data.actions_taken) {
        return 'backend';
    }
    
    return 'unknown';
}

/**
 * Migrates a single legacy message to standard format
 * @param {*} legacy - Legacy message data
 * @param {string} conversationId - Conversation ID
 * @returns {import('../../types/ai').Message | import('../../types/ai').Message[]}
 */
export function migrateLegacyMessage(legacy, conversationId = '') {
    const format = detectMessageFormat(legacy);
    
    switch (format) {
        case 'standard':
            // Already in standard format, just ensure all fields
            return createMessage({
                ...legacy,
                conversationId: legacy.conversationId || conversationId
            });
            
        case 'legacy_pair':
            // Convert prompt/response pair to two messages
            const baseTimestamp = legacy.timestamp || new Date().toISOString();
            const baseId = legacy.id || Date.now();
            
            return [
                createMessage({
                    id: `${baseId}-user`,
                    conversationId,
                    role: MessageRole.USER,
                    content: legacy.prompt,
                    timestamp: baseTimestamp,
                    status: MessageStatus.DELIVERED
                }),
                createMessage({
                    id: `${baseId}-assistant`,
                    conversationId,
                    role: MessageRole.ASSISTANT,
                    content: legacy.response,
                    timestamp: new Date(new Date(baseTimestamp).getTime() + 1000).toISOString(), // 1 second later
                    metadata: {
                        actions: legacy.actions || []
                    },
                    status: MessageStatus.DELIVERED
                })
            ];
            
        case 'backend':
            // Convert backend format (actions_taken -> actions)
            return createMessage({
                ...legacy,
                conversationId: legacy.conversationId || conversationId,
                metadata: {
                    ...legacy.metadata,
                    actions: legacy.actions_taken || legacy.actions || []
                }
            });
            
        default:
            
            // Try to create a basic message with available data
            return createMessage({
                id: legacy.id || `unknown-${Date.now()}`,
                conversationId,
                role: legacy.role || MessageRole.USER,
                content: legacy.content || legacy.message || JSON.stringify(legacy),
                timestamp: legacy.timestamp || legacy.created_at || new Date().toISOString()
            });
    }
}

/**
 * Migrates an array of messages to standard format
 * @param {*[]} messages - Array of messages in any format
 * @param {string} conversationId - Conversation ID
 * @returns {import('../../types/ai').Message[]}
 */
export function migrateMessageArray(messages, conversationId = '') {
    if (!Array.isArray(messages)) {
        
        return [];
    }
    
    const migratedMessages = [];
    
    for (const message of messages) {
        const migrated = migrateLegacyMessage(message, conversationId);
        
        if (Array.isArray(migrated)) {
            migratedMessages.push(...migrated);
        } else {
            migratedMessages.push(migrated);
        }
    }
    
    // Sort by timestamp
    return migratedMessages.sort((a, b) => 
        new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
}

/**
 * Migrates conversation history from Redux store
 * @param {*[]} conversationHistory - Legacy conversation history
 * @param {string} conversationId - Conversation ID
 * @returns {import('../../types/ai').Message[]}
 */
export function migrateReduxConversationHistory(conversationHistory, conversationId = '') {
    return migrateMessageArray(conversationHistory, conversationId);
}

/**
 * Validates a message has all required fields
 * @param {*} message - Message to validate
 * @returns {boolean}
 */
export function isValidMessage(message) {
    return Boolean(
        message &&
        typeof message === 'object' &&
        message.id &&
        message.role &&
        Object.values(MessageRole).includes(message.role) &&
        typeof message.content === 'string' &&
        message.timestamp
    );
}

/**
 * Sanitizes and validates a message array
 * @param {*[]} messages - Messages to validate
 * @returns {import('../../types/ai').Message[]}
 */
export function sanitizeMessages(messages) {
    if (!Array.isArray(messages)) {
        return [];
    }
    
    return messages
        .filter(msg => {
            const isValid = isValidMessage(msg);
            if (!isValid) {
                
            }
            return isValid;
        })
        .map(msg => createMessage(msg));
}

/**
 * Merges messages from different sources, removing duplicates
 * @param {...import('../../types/ai').Message[][]} messageSources - Arrays of messages
 * @returns {import('../../types/ai').Message[]}
 */
export function mergeMessages(...messageSources) {
    const messageMap = new Map();
    
    for (const messages of messageSources) {
        if (!Array.isArray(messages)) continue;
        
        for (const message of messages) {
            // Use the message ID as key to prevent duplicates
            // Prefer non-optimistic messages over optimistic ones
            const existing = messageMap.get(message.id);
            if (!existing || (existing._optimistic && !message._optimistic)) {
                messageMap.set(message.id, message);
            }
        }
    }
    
    // Convert back to array and sort by timestamp
    return Array.from(messageMap.values()).sort((a, b) =>
        new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
}

/**
 * Extracts conversation ID from various data formats
 * @param {*} data - Data containing conversation ID
 * @returns {string | null}
 */
export function extractConversationId(data) {
    if (!data) return null;
    
    // Direct ID
    if (typeof data === 'string') return data;
    
    // Object with various ID fields
    if (typeof data === 'object') {
        return data.conversationId || 
               data.conversation_id || 
               data.id ||
               data.conversation?.id ||
               null;
    }
    
    return null;
}

/**
 * Creates a migration report for debugging
 * @param {*[]} originalMessages - Original messages
 * @param {import('../../types/ai').Message[]} migratedMessages - Migrated messages
 * @returns {Object}
 */
export function createMigrationReport(originalMessages, migratedMessages) {
    const formats = new Map();
    
    originalMessages.forEach(msg => {
        const format = detectMessageFormat(msg);
        formats.set(format, (formats.get(format) || 0) + 1);
    });
    
    return {
        originalCount: originalMessages.length,
        migratedCount: migratedMessages.length,
        formats: Object.fromEntries(formats),
        messagesAdded: migratedMessages.length - originalMessages.length,
        timestamp: new Date().toISOString()
    };
}
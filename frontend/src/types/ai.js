/**
 * Standardized type definitions for AI functionality
 * Single source of truth for all AI-related data structures
 */

/**
 * Message roles
 * @readonly
 * @enum {string}
 */
export const MessageRole = {
    USER: 'USER',
    ASSISTANT: 'ASSISTANT',
    SYSTEM: 'SYSTEM'
};

/**
 * Message status
 * @readonly
 * @enum {string}
 */
export const MessageStatus = {
    PENDING: 'pending',
    DELIVERED: 'delivered',
    FAILED: 'failed'
};

/**
 * Assistant types
 * @readonly
 * @enum {string}
 */
export const AssistantType = {
    STANDARD: 'standard',
    ADMIN: 'admin',
    BUSINESS: 'business',
    SUPPORT: 'support',
    SCHEDULER: 'scheduler'
};

/**
 * View modes for AI interface
 * @readonly
 * @enum {string}
 */
export const ViewMode = {
    CHAT: 'chat',
    SPLIT: 'split',
    RESULTS: 'results'
};

/**
 * Action type definition
 * @typedef {Object} Action
 * @property {string} type - Action type
 * @property {string} description - Human-readable description
 * @property {*} [result] - Action result data
 */

/**
 * Message metadata definition
 * @typedef {Object} MessageMetadata
 * @property {Action[]} [actions] - Actions taken by the assistant
 * @property {number} [confidence] - Confidence score (0-1)
 * @property {number} [processingTime] - Processing time in ms
 * @property {Object} [data] - Additional data from assistant
 */

/**
 * Standardized Message format
 * @typedef {Object} Message
 * @property {string} id - Unique message ID
 * @property {string} conversationId - Parent conversation ID
 * @property {MessageRole} role - Message sender role
 * @property {string} content - Message content
 * @property {string} timestamp - ISO timestamp
 * @property {MessageMetadata} [metadata] - Optional metadata
 * @property {MessageStatus} [status] - Message delivery status
 * @property {boolean} [_optimistic] - Flag for optimistic updates
 */

/**
 * Assistant capability definition
 * @typedef {Object} AssistantCapability
 * @property {string} name - Capability name
 * @property {string} description - Capability description
 * @property {boolean} enabled - Whether capability is enabled
 */

/**
 * Assistant definition
 * @typedef {Object} Assistant
 * @property {string} id - Unique assistant ID
 * @property {string} name - Assistant name
 * @property {string} description - Assistant description
 * @property {AssistantType} type - Assistant type
 * @property {boolean} isActive - Whether assistant is active
 * @property {boolean} isDefault - Whether this is the default assistant
 * @property {AssistantCapability[]} capabilities - Assistant capabilities
 * @property {Object} configuration - Assistant configuration
 * @property {number} configuration.temperature - Temperature setting (0-1)
 * @property {number} configuration.maxTokens - Max tokens per response
 * @property {string} configuration.model - Model identifier
 * @property {string} createdAt - ISO timestamp
 * @property {string} updatedAt - ISO timestamp
 */

/**
 * Conversation context definition
 * @typedef {Object} ConversationContext
 * @property {string} [source] - Context source (e.g., 'web_app')
 * @property {string} [timestamp] - Context timestamp
 * @property {Object} [metadata] - Additional context metadata
 */

/**
 * Conversation definition
 * @typedef {Object} Conversation
 * @property {string} id - Unique conversation ID
 * @property {string} assistantId - Associated assistant ID
 * @property {string} userId - Owner user ID
 * @property {boolean} active - Whether conversation is active
 * @property {ConversationContext} [context] - Conversation context
 * @property {string} createdAt - ISO timestamp
 * @property {string} updatedAt - ISO timestamp
 * @property {string} [lastMessageAt] - Last message timestamp
 * @property {number} [messageCount] - Total message count
 */

/**
 * Chat request definition
 * @typedef {Object} ChatRequest
 * @property {string} message - Message content
 * @property {Object} [context] - Additional context
 * @property {string[]} [attachments] - File attachments
 */

/**
 * Chat response definition
 * @typedef {Object} ChatResponse
 * @property {Message} userMessage - User's message
 * @property {Message} assistantMessage - Assistant's response
 * @property {string} conversationId - Conversation ID
 * @property {Action[]} [actions] - Actions performed
 * @property {Object} [data] - Additional response data
 */

/**
 * API Response wrapper
 * @template T
 * @typedef {Object} ApiResponse
 * @property {boolean} success - Whether request succeeded
 * @property {T} [data] - Response data
 * @property {string} [error] - Error message if failed
 * @property {number} [status] - HTTP status code
 * @property {boolean} [canceled] - Whether request was canceled
 */

/**
 * Pagination info
 * @typedef {Object} PaginationInfo
 * @property {number} page - Current page
 * @property {number} limit - Items per page
 * @property {number} total - Total items
 * @property {number} totalPages - Total pages
 */

/**
 * Paginated response
 * @template T
 * @typedef {Object} PaginatedResponse
 * @property {T[]} items - Page items
 * @property {PaginationInfo} pagination - Pagination info
 */

/**
 * Conversation statistics
 * @typedef {Object} ConversationStats
 * @property {number} totalConversations - Total conversations
 * @property {number} activeConversations - Active conversations
 * @property {number} totalMessages - Total messages
 * @property {number} messagestoday - Messages today
 * @property {number} messagesThisWeek - Messages this week
 * @property {number} messagesThisMonth - Messages this month
 * @property {number} avgMessagesPerConversation - Average messages per conversation
 * @property {string} [mostUsedAssistantId] - Most used assistant ID
 */

/**
 * Helper function to create a message
 * @param {Partial<Message>} data - Message data
 * @returns {Message}
 */
export function createMessage(data) {
    return {
        id: data.id || `msg-${Date.now()}`,
        conversationId: data.conversationId || '',
        role: data.role || MessageRole.USER,
        content: data.content || '',
        timestamp: data.timestamp || new Date().toISOString(),
        metadata: data.metadata || {},
        status: data.status || MessageStatus.PENDING,
        _optimistic: data._optimistic || false
    };
}

/**
 * Helper to check if a message is from user
 * @param {Message} message - Message to check
 * @returns {boolean}
 */
export function isUserMessage(message) {
    return message.role === MessageRole.USER;
}

/**
 * Helper to check if a message is from assistant
 * @param {Message} message - Message to check
 * @returns {boolean}
 */
export function isAssistantMessage(message) {
    return message.role === MessageRole.ASSISTANT;
}

/**
 * Helper to check if a message is optimistic
 * @param {Message} message - Message to check
 * @returns {boolean}
 */
export function isOptimisticMessage(message) {
    return Boolean(message._optimistic);
}

/**
 * Helper to format timestamp for display
 * @param {string} timestamp - ISO timestamp
 * @returns {string}
 */
export function formatMessageTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    
    return date.toLocaleDateString();
}
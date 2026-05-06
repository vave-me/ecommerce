import AssistantService from './AssistantService';

/**
 * Custom error class for conversation not found errors
 */
export class ConversationNotFoundError extends Error {
    constructor(conversationId) {
        super(`Conversation ${conversationId} not found`);
        this.name = 'ConversationNotFoundError';
        this.conversationId = conversationId;
    }
}

/**
 * Custom error class for assistant not found errors
 */
export class AssistantNotFoundError extends Error {
    constructor(assistantId) {
        super(`Assistant ${assistantId} not found`);
        this.name = 'AssistantNotFoundError';
        this.assistantId = assistantId;
    }
}

/**
 * Retry configuration
 */
const RETRY_CONFIG = {
    maxRetries: 3,
    initialDelay: 1000, // 1 second
    maxDelay: 10000, // 10 seconds
    backoffMultiplier: 2,
    retryableStatusCodes: [408, 429, 500, 502, 503, 504],
    retryableErrors: ['ECONNABORTED', 'ETIMEDOUT', 'ENOTFOUND', 'ENETUNREACH']
};

/**
 * Enhanced Assistant Service with retry logic and better error handling
 * Extends the base AssistantService with automatic retry capabilities
 */
class AssistantServiceWithRetry extends AssistantService {
    constructor() {
        super();
        this.retryConfig = { ...RETRY_CONFIG };
    }

    /**
     * Get singleton instance with retry capabilities
     * @returns {AssistantServiceWithRetry}
     */
    static getInstance() {
        if (!AssistantServiceWithRetry.instance) {
            AssistantServiceWithRetry.instance = new AssistantServiceWithRetry();
        }
        return AssistantServiceWithRetry.instance;
    }

    /**
     * Execute a function with retry logic
     * @param {Function} fn - Function to execute
     * @param {string} operation - Operation name for logging
     * @param {number} attempt - Current attempt number
     * @returns {Promise<any>}
     */
    async executeWithRetry(fn, operation, attempt = 1) {
        try {
            return await fn();
        } catch (error) {
            // Check if error is retryable
            if (!this.isRetryableError(error) || attempt >= this.retryConfig.maxRetries) {
                throw error;
            }

            // Calculate delay with exponential backoff
            const delay = Math.min(
                this.retryConfig.initialDelay * Math.pow(this.retryConfig.backoffMultiplier, attempt - 1),
                this.retryConfig.maxDelay
            );

            console.log(`Retrying ${operation} (attempt ${attempt}) after error:`, {
                error: error.message,
                status: error.response?.status
            });

            // Wait before retry
            await new Promise(resolve => setTimeout(resolve, delay));

            // Retry the operation
            return this.executeWithRetry(fn, operation, attempt + 1);
        }
    }

    /**
     * Check if an error is retryable
     * @param {Error} error - Error to check
     * @returns {boolean}
     */
    isRetryableError(error) {
        // Check for network errors
        if (error.code && this.retryConfig.retryableErrors.includes(error.code)) {
            return true;
        }

        // Check for specific HTTP status codes
        if (error.response?.status && this.retryConfig.retryableStatusCodes.includes(error.response.status)) {
            return true;
        }

        // Don't retry client errors (4xx) except specific ones
        if (error.response?.status >= 400 && error.response?.status < 500) {
            return false;
        }

        return false;
    }

    /**
     * Enhanced error handling with custom error types
     * @param {Error} error - Error object
     * @param {string} defaultMessage - Default error message
     * @returns {{success: boolean, error: string, status: number, data: null}}
     */
    handleError(error, defaultMessage) {
        // Handle conversation not found specifically
        if (error.response?.status === 404 && defaultMessage.includes('conversation')) {
            const conversationId = error.config?.url?.match(/conversations\/([^\/]+)/)?.[1];
            throw new ConversationNotFoundError(conversationId || 'unknown');
        }

        // Handle assistant not found specifically
        if (error.response?.status === 404 && defaultMessage.includes('assistant')) {
            const assistantId = error.config?.url?.match(/assistants\/([^\/]+)/)?.[1];
            throw new AssistantNotFoundError(assistantId || 'unknown');
        }

        return super.handleError(error, defaultMessage);
    }

    // Override methods to add retry logic

    async getAssistants(params) {
        return this.executeWithRetry(
            () => super.getAssistants(params),
            'getAssistants'
        );
    }

    async getAssistant(assistantId) {
        return this.executeWithRetry(
            () => super.getAssistant(assistantId),
            'getAssistant'
        );
    }

    async createConversation(assistantId, context) {
        // Don't retry conversation creation to avoid duplicates
        return super.createConversation(assistantId, context);
    }

    async getUserConversations(params) {
        return this.executeWithRetry(
            () => super.getUserConversations(params),
            'getUserConversations'
        );
    }

    async getConversation(conversationId) {
        return this.executeWithRetry(
            () => super.getConversation(conversationId),
            'getConversation'
        );
    }

    async sendMessage(conversationId, message, context) {
        // Special handling for send message to ensure conversation exists
        try {
            return await this.executeWithRetry(
                () => super.sendMessage(conversationId, message, context),
                'sendMessage'
            );
        } catch (error) {
            // If conversation not found, bubble up a clear error
            if (error instanceof ConversationNotFoundError) {
                throw error;
            }
            
            // For other errors, use default handling
            return this.handleError(error, 'Failed to send message');
        }
    }

    async getMessages(conversationId, params) {
        return this.executeWithRetry(
            () => super.getMessages(conversationId, params),
            'getMessages'
        );
    }

    async getStats() {
        return this.executeWithRetry(
            () => super.getStats(),
            'getStats'
        );
    }

    async processSpeechInput(params) {
        return this.executeWithRetry(
            () => super.processSpeechInput(params),
            'processSpeechInput'
        );
    }

    async processUserInput(params) {
        return this.executeWithRetry(
            () => super.processUserInput(params),
            'processUserInput'
        );
    }

    async processImageInput(params) {
        return this.executeWithRetry(
            () => super.processImageInput(params),
            'processImageInput'
        );
    }

    async processDocumentInput(params) {
        return this.executeWithRetry(
            () => super.processDocumentInput(params),
            'processDocumentInput'
        );
    }

    async activateAssistant(assistantId) {
        // Don't retry activation to avoid duplicates
        return super.activateAssistant(assistantId);
    }

    async deactivateAssistant(assistantId) {
        return this.executeWithRetry(
            () => super.deactivateAssistant(assistantId),
            'deactivateAssistant'
        );
    }

    async updateConversation(conversationId, updates) {
        return this.executeWithRetry(
            () => super.updateConversation(conversationId, updates),
            'updateConversation'
        );
    }

    async processAssistantRequestAdvanced(params) {
        return this.executeWithRetry(
            () => super.processAssistantRequestAdvanced(params),
            'processAssistantRequestAdvanced'
        );
    }

    /**
     * Configure retry behavior
     * @param {Object} config - Retry configuration
     */
    configureRetry(config) {
        this.retryConfig = { ...this.retryConfig, ...config };
    }
}

// Export enhanced service as default
export default AssistantServiceWithRetry;

// Export getter for enhanced service
export const getAssistantService = () => AssistantServiceWithRetry.getInstance();

// End of file

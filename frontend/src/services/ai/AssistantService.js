import * as api from '../../api/client/assistantsApi';

/**
 * Unified Assistant Service
 * Single source of truth for all AI/Assistant API operations
 * Implements singleton pattern to ensure consistency
 * Now uses the production-ready assistantsApi client
 */
class AssistantService {
    constructor() {
        if (AssistantService.instance) {
            return AssistantService.instance;
        }
        
        AssistantService.instance = this;
    }

    /**
     * Get singleton instance
     * @returns {AssistantService}
     */
    static getInstance() {
        if (!AssistantService.instance) {
            AssistantService.instance = new AssistantService();
        }
        return AssistantService.instance;
    }

    // ============================================================================
    // ASSISTANT MANAGEMENT
    // ============================================================================

    /**
     * Get all assistants
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @returns {Promise<{success: boolean, data: {assistants: Array}}>}
     */
    async getAssistants({ page = 1, limit = 20 } = {}) {
        try {
            const data = await api.getAssistants();
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch assistants');
        }
    }

    /**
     * Get assistant by ID
     * @param {string} assistantId - Assistant ID
     * @returns {Promise<{success: boolean, data: {assistant: Object}}>}
     */
    async getAssistant(assistantId) {
        try {
            const data = await api.getAssistant(assistantId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch assistant');
        }
    }

    /**
     * Update assistant configuration
     * @param {string} assistantId - Assistant ID
     * @param {Object} config - Configuration updates
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async updateAssistantConfig(assistantId, config) {
        try {
            const data = await api.updateAssistantConfiguration(assistantId, config);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to update assistant configuration');
        }
    }

    /**
     * Activate a new assistant
     * @param {string} assistantId - Assistant ID to activate
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async activateAssistant(assistantId) {
        try {
            const data = await api.activateAssistant(assistantId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to activate assistant');
        }
    }

    /**
     * Deactivate an assistant
     * @param {string} assistantId - Assistant ID to deactivate
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async deactivateAssistant(assistantId) {
        try {
            const data = await api.deactivateAssistant(assistantId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to deactivate assistant');
        }
    }

    // ============================================================================
    // CONVERSATION MANAGEMENT
    // ============================================================================

    /**
     * Create a new conversation
     * @param {string} assistantId - Assistant ID
     * @param {Object} context - Initial context
     * @returns {Promise<{success: boolean, data: {conversationId: string, conversation: Object}}>}
     */
    async createConversation(assistantId, context = {}) {
        try {
            const data = await api.createConversation({
                assistantId,
                initialContext: {
                    ...context,
                    source: 'web_app',
                    timestamp: new Date().toISOString()
                }
            });
            
            // Normalize response to ensure consistent structure
            const conversationId = data?.conversationId;
            return {
                success: true,
                data: {
                    conversationId,
                    conversation: data?.conversation || { id: conversationId },
                    ...data
                }
            };
        } catch (error) {
            return this.handleError(error, 'Failed to create conversation');
        }
    }

    /**
     * Get user conversations
     * @param {Object} params - Query parameters
     * @param {boolean} params.activeOnly - Filter active conversations only
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @returns {Promise<{success: boolean, data: {conversations: Array}}>}
     */
    async getUserConversations({ activeOnly = false, page = 1, limit = 20 } = {}) {
        try {
            const data = await api.getUserConversations();
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch user conversations');
        }
    }

    /**
     * Get conversation by ID
     * @param {string} conversationId - Conversation ID
     * @returns {Promise<{success: boolean, data: {conversation: Object}}>}
     */
    async getConversation(conversationId) {
        try {
            const data = await api.getConversation(conversationId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch conversation');
        }
    }

    /**
     * Check if conversation exists
     * @param {string} conversationId - Conversation ID
     * @returns {Promise<boolean>}
     */
    async conversationExists(conversationId) {
        try {
            await api.getConversation(conversationId);
            return true;
        } catch (error) {
            if (error.response?.status === 404) {
                return false;
            }
            throw error;
        }
    }

    /**
     * Update conversation context
     * @param {string} conversationId - Conversation ID
     * @param {Object} context - Context updates
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async updateConversationContext(conversationId, context) {
        try {
            const data = await api.updateConversationContext(conversationId, context);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to update conversation context');
        }
    }

    /**
     * Archive conversation
     * @param {string} conversationId - Conversation ID
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async archiveConversation(conversationId) {
        try {
            const data = await api.archiveConversation(conversationId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to archive conversation');
        }
    }

    /**
     * Delete conversation
     * @param {string} conversationId - Conversation ID
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async deleteConversation(conversationId) {
        try {
            const data = await api.deleteConversation(conversationId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to delete conversation');
        }
    }

    /**
     * Update conversation
     * @param {string} conversationId - Conversation ID
     * @param {Object} updates - Conversation updates
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async updateConversation(conversationId, updates) {
        try {
            const data = await api.updateConversation(conversationId, updates);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to update conversation');
        }
    }

    // ============================================================================
    // MESSAGING
    // ============================================================================

    /**
     * Send message in conversation
     * Ensures conversation exists before sending
     * @param {string} conversationId - Conversation ID
     * @param {string} message - Message content
     * @param {Object} context - Additional context
     * @returns {Promise<{success: boolean, data: {userMessage: Object, assistantMessage: Object}}>}
     */
    async sendMessage(conversationId, message, context = {}) {
        try {
            // Ensure conversation exists first
            const exists = await this.conversationExists(conversationId);
            if (!exists) {
                throw new Error(`Conversation ${conversationId} not found`);
            }

            const data = await api.chatWithConversation({
                assistantId: context.assistantId, // Assistant ID required for chat
                conversationId,
                message,
                context: {
                    ...context,
                    timestamp: new Date().toISOString()
                },
                maxHistoryMessages: 10
            });

            // Normalize response structure
            return {
                success: true,
                data: {
                    userMessage: {
                        id: data?.userMessageId || `user-${Date.now()}`,
                        role: 'USER',
                        content: message,
                        timestamp: new Date().toISOString()
                    },
                    assistantMessage: {
                        id: data?.messageId,
                        role: 'ASSISTANT',
                        content: data?.response,
                        timestamp: new Date().toISOString(),
                        metadata: {
                            actions: data?.actions,
                            confidence: data?.confidence,
                            data: data?.data
                        }
                    },
                    ...data
                }
            };
        } catch (error) {
            return this.handleError(error, 'Failed to send message');
        }
    }

    /**
     * Get conversation messages
     * @param {string} conversationId - Conversation ID
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @returns {Promise<{success: boolean, data: {messages: Array}}>}
     */
    async getMessages(conversationId, { page = 1, limit = 50 } = {}) {
        try {
            const data = await api.getConversationMessages(conversationId);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch messages');
        }
    }

    /**
     * Add message to conversation (without AI response)
     * @param {string} conversationId - Conversation ID
     * @param {Object} message - Message object
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async addMessage(conversationId, message) {
        try {
            const data = await api.addMessageToConversation(conversationId, message);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to add message');
        }
    }

    // ============================================================================
    // STATISTICS
    // ============================================================================

    /**
     * Get conversation statistics
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async getStats() {
        try {
            const data = await api.getConversationStats();
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to fetch statistics');
        }
    }

    // ============================================================================
    // ADVANCED PROCESSING
    // ============================================================================

    /**
     * Process advanced assistant request
     * @param {Object} params - Request parameters
     * @param {string} params.assistantId - Assistant ID
     * @param {string} params.request - User request text
     * @param {Object} params.context - Request context
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async processAssistantRequestAdvanced(params) {
        try {
            const data = await api.processAssistantRequestAdvanced(params);
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to process assistant request');
        }
    }

    // ============================================================================
    // INPUT PROCESSING
    // ============================================================================

    /**
     * Process speech input and get transcription
     * @param {Object} params - Speech processing parameters
     * @param {string} params.audioData - Base64 encoded audio data
     * @param {string} params.audioFormat - Audio format (webm, mp3, etc)
     * @param {string} params.language - Language code
     * @param {Object} params.context - Additional context
     * @param {string} [params.assistantId] - Optional assistant ID
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async processSpeechInput({ audioData, audioFormat, language, context, assistantId }) {
        try {
            const data = await api.processSpeechInput({
                assistantId,
                audioData,
                audioFormat,
                language,
                context
            });
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to process speech input');
        }
    }

    /**
     * Process text input
     * @param {Object} params - Text processing parameters
     * @param {string} params.text - Input text
     * @param {string} params.assistantId - Assistant ID
     * @param {Object} params.context - Additional context
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async processUserInput({ text, assistantId, context }) {
        try {
            // API expects 'message' not 'text'
            const data = await api.processUserInput({
                message: text,
                assistantId,
                context
            });
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to process text input');
        }
    }

    /**
     * Process image input
     * @param {Object} params - Image processing parameters
     * @param {string} params.imageData - Base64 encoded image data
     * @param {string} params.assistantId - Assistant ID
     * @param {Object} params.context - Additional context
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async processImageInput({ imageData, assistantId, context }) {
        try {
            const data = await api.processImageInput({
                imageData,
                assistantId,
                context
            });
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to process image input');
        }
    }

    /**
     * Process document input
     * @param {Object} params - Document processing parameters
     * @param {string} params.documentData - Base64 encoded document data
     * @param {string} params.documentType - Document type
     * @param {string} params.assistantId - Assistant ID
     * @param {Object} params.context - Additional context
     * @returns {Promise<{success: boolean, data: Object}>}
     */
    async processDocumentInput({ documentData, documentFormat, assistantId, context }) {
        try {
            const data = await api.processDocumentInput({
                documentData,
                documentFormat, // Fixed: was documentType, should be documentFormat
                assistantId,
                context
            });
            return {
                success: true,
                data
            };
        } catch (error) {
            return this.handleError(error, 'Failed to process document input');
        }
    }

    // ============================================================================
    // ERROR HANDLING
    // ============================================================================

    /**
     * Centralized error handling
     * @param {Error} error - Error object
     * @param {string} defaultMessage - Default error message
     * @returns {{success: boolean, error: string, status: number, data: null}}
     */
    handleError(error, defaultMessage) {
        if (error?.code === 'ERR_CANCELED' || error?.name === 'CanceledError') {
            return {
                success: false,
                error: 'Request canceled',
                canceled: true,
                status: 0,
                data: null
            };
        }

        const message = error?.response?.data?.message || error?.message || defaultMessage;
        const status = error?.response?.status || 500;

        // Error: 'AssistantService Error:', message logged for debugging

        return {
            success: false,
            error: message,
            status,
            data: null
        };
    }

    /**
     * Clear service instance (useful for testing)
     */
    static clearInstance() {
        AssistantService.instance = null;
    }
}

// Export singleton instance getter
export const getAssistantService = () => AssistantService.getInstance();

// Export class for testing
export default AssistantService;
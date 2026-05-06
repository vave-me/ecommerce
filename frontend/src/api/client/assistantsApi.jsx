import axiosInstance from '../axiosInstance';

/**
 * Production-ready Assistants API Client
 * Implements all endpoints from api.swagger.json with proper error handling
 * All methods use camelCase naming convention
 */

const API_BASE = '/assistants';

// ============================================================================
// ASSISTANT MANAGEMENT
// ============================================================================

/**
 * Get all assistants
 * GET /api/assistants
 */
export const getAssistants = async () => {
  const response = await axiosInstance.get(API_BASE);
  return response.data;
};

/**
 * Get assistant by ID
 * GET /api/assistants/{id}
 */
export const getAssistant = async (id) => {
  const response = await axiosInstance.get(`${API_BASE}/${id}`);
  return response.data;
};

/**
 * Update assistant configuration
 * PUT /api/assistants/{id}
 */
export const updateAssistantConfiguration = async (id, config) => {
  const response = await axiosInstance.put(`${API_BASE}/${id}`, config);
  return response.data;
};

/**
 * Activate assistant
 * POST /api/assistants/{id}/activate
 */
export const activateAssistant = async (id) => {
  const response = await axiosInstance.post(`${API_BASE}/${id}/activate`, {});
  return response.data;
};

/**
 * Deactivate assistant
 * POST /api/assistants/{id}/deactivate
 */
export const deactivateAssistant = async (id) => {
  const response = await axiosInstance.post(`${API_BASE}/${id}/deactivate`, {});
  return response.data;
};

// ============================================================================
// CONVERSATION MANAGEMENT
// ============================================================================

/**
 * Create conversation
 * POST /api/assistants/conversations
 */
export const createConversation = async ({ assistantId, initialContext = {} }) => {
  const response = await axiosInstance.post(`${API_BASE}/conversations`, {
    assistantId,
    initialContext
  });
  return response.data;
};

/**
 * Get user conversations
 * GET /api/assistants/conversations
 */
export const getUserConversations = async () => {
  const response = await axiosInstance.get(`${API_BASE}/conversations`);
  return response.data;
};

/**
 * Get conversation by ID
 * GET /api/assistants/conversations/{conversationId}
 */
export const getConversation = async (conversationId) => {
  const response = await axiosInstance.get(`${API_BASE}/conversations/${conversationId}`);
  return response.data;
};

/**
 * Update conversation
 * PUT /api/assistants/conversations/{conversationId}
 */
export const updateConversation = async (conversationId, { title, metadata }) => {
  const response = await axiosInstance.put(`${API_BASE}/conversations/${conversationId}`, {
    title,
    metadata
  });
  return response.data;
};

/**
 * Update conversation context
 * PUT /api/assistants/conversations/{conversationId}/context
 */
export const updateConversationContext = async (conversationId, context) => {
  const response = await axiosInstance.put(`${API_BASE}/conversations/${conversationId}/context`, {
    context
  });
  return response.data;
};

/**
 * Archive conversation
 * POST /api/assistants/conversations/{conversationId}/archive
 */
export const archiveConversation = async (conversationId) => {
  const response = await axiosInstance.post(`${API_BASE}/conversations/${conversationId}/archive`, {});
  return response.data;
};

/**
 * Delete conversation
 * DELETE /api/assistants/conversations/{conversationId}
 */
export const deleteConversation = async (conversationId) => {
  const response = await axiosInstance.delete(`${API_BASE}/conversations/${conversationId}`);
  return response.data;
};

// ============================================================================
// MESSAGING
// ============================================================================

/**
 * Chat with conversation
 * POST /api/assistants/chat
 */
export const chatWithConversation = async ({ assistantId, conversationId, message, context = {}, maxHistoryMessages = 10 }) => {
  const response = await axiosInstance.post(`${API_BASE}/chat`, {
    assistantId,
    conversationId,
    message,
    context,
    maxHistoryMessages
  });
  return response.data;
};

/**
 * Get conversation messages
 * GET /api/assistants/conversations/{conversationId}/messages
 */
export const getConversationMessages = async (conversationId) => {
  const response = await axiosInstance.get(`${API_BASE}/conversations/${conversationId}/messages`);
  return response.data;
};

/**
 * Add message to conversation
 * POST /api/assistants/conversations/{conversationId}/messages
 */
export const addMessageToConversation = async (conversationId, { assistantId, role, content, metadata = {} }) => {
  const response = await axiosInstance.post(`${API_BASE}/conversations/${conversationId}/messages`, {
    assistantId,
    role,
    content,
    metadata
  });
  return response.data;
};

// ============================================================================
// INPUT PROCESSING
// ============================================================================

/**
 * Process user text input
 * POST /api/assistants/input/text
 */
export const processUserInput = async ({ assistantId, message, context = {}, requestType }) => {
  const response = await axiosInstance.post(`${API_BASE}/input/text`, {
    assistantId,
    message,
    context,
    requestType
  });
  return response.data;
};

/**
 * Process speech input
 * POST /api/assistants/input/speech
 */
export const processSpeechInput = async ({ assistantId, audioData, audioFormat, language, context = {}, requestType }) => {
  const response = await axiosInstance.post(`${API_BASE}/input/speech`, {
    assistantId,
    audioData,
    audioFormat,
    language,
    context,
    requestType
  });
  return response.data;
};

/**
 * Process image input
 * POST /api/assistants/input/image
 */
export const processImageInput = async ({ 
  assistantId, 
  imageData, 
  imageUrl, 
  imageFormat, 
  analysisType, 
  userPrompt, 
  context = {}, 
  requestType 
}) => {
  const response = await axiosInstance.post(`${API_BASE}/input/image`, {
    assistantId,
    imageData,
    imageUrl,
    imageFormat,
    analysisType,
    userPrompt,
    context,
    requestType
  });
  return response.data;
};

/**
 * Process document input
 * POST /api/assistants/input/document
 */
export const processDocumentInput = async ({ 
  assistantId, 
  documentData, 
  documentUrl, 
  documentFormat, 
  analysisType, 
  userPrompt, 
  context = {}, 
  requestType 
}) => {
  const response = await axiosInstance.post(`${API_BASE}/input/document`, {
    assistantId,
    documentData,
    documentUrl,
    documentFormat,
    analysisType,
    userPrompt,
    context,
    requestType
  });
  return response.data;
};

// ============================================================================
// ADVANCED PROCESSING
// ============================================================================

/**
 * Process assistant request (advanced)
 * POST /api/assistants/process
 */
export const processAssistantRequestAdvanced = async ({ assistantId, requestId, message, context = {} }) => {
  const response = await axiosInstance.post(`${API_BASE}/process`, {
    assistantId,
    requestId,
    message,
    context
  });
  return response.data;
};

// ============================================================================
// STATISTICS
// ============================================================================

/**
 * Get conversation statistics
 * GET /api/assistants/stats
 */
export const getConversationStats = async () => {
  const response = await axiosInstance.get(`${API_BASE}/stats`);
  return response.data;
};

// Export all functions as a named object for convenience
export const assistantsApi = {
  // Assistant Management
  getAssistants,
  getAssistant,
  updateAssistantConfiguration,
  activateAssistant,
  deactivateAssistant,
  
  // Conversation Management
  createConversation,
  getUserConversations,
  getConversation,
  updateConversation,
  updateConversationContext,
  archiveConversation,
  deleteConversation,
  
  // Messaging
  chatWithConversation,
  getConversationMessages,
  addMessageToConversation,
  
  // Input Processing
  processUserInput,
  processSpeechInput,
  processImageInput,
  processDocumentInput,
  
  // Advanced Processing
  processAssistantRequestAdvanced,
  
  // Statistics
  getConversationStats
};

export default assistantsApi;
import axiosInstance from "../../axiosInstance";

const API_BASE_URL = '/support';

// ===== TICKET MANAGEMENT =====
export const listTickets = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list tickets.');
    }
};

export const searchTickets = async (searchParams) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/search`, searchParams);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to search tickets.');
    }
};

export const getTicketsByStatus = async (status, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets/status/${status}`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get tickets by status.');
    }
};

export const getTicketsByPriority = async (priority, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets/priority/${priority}`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get tickets by priority.');
    }
};

export const getTicketsByAgent = async (agentId, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/agents/${agentId}/tickets`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get tickets by agent.');
    }
};

export const getUnassignedTickets = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets/unassigned`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get unassigned tickets.');
    }
};

export const bulkUpdateTickets = async (updates) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/bulk-update`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to bulk update tickets.');
    }
};

export const bulkAssignTickets = async (assignments) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/bulk-assign`, assignments);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to bulk assign tickets.');
    }
};

export const mergeTickets = async (mergeData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/merge`, mergeData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to merge tickets.');
    }
};

export const archiveTicket = async (ticketId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/archive`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to archive ticket.');
    }
};

export const unarchiveTicket = async (ticketId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/unarchive`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to unarchive ticket.');
    }
};

export const deleteTicket = async (ticketId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/tickets/${ticketId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete ticket.');
    }
};

// ===== TICKET TAGS =====
export const addTicketTags = async (ticketId, tags) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/tags`, { tags });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to add ticket tags.');
    }
};

export const removeTicketTag = async (ticketId, tag) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/tickets/${ticketId}/tags/${tag}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to remove ticket tag.');
    }
};

// ===== AGENT MANAGEMENT =====
export const listAgents = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/agents`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list agents.');
    }
};

export const getAgent = async (agentId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/agents/${agentId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get agent.');
    }
};

export const createAgent = async (agentData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/agents`, agentData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create agent.');
    }
};

export const updateAgent = async (agentId, updates) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/agents/${agentId}`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update agent.');
    }
};

export const updateAgentStatus = async (agentId, status) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/agents/${agentId}/status`, { status });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update agent status.');
    }
};

export const deleteAgent = async (agentId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/agents/${agentId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete agent.');
    }
};

// ===== KNOWLEDGE BASE =====
export const listKnowledgeArticles = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/knowledge/articles`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list knowledge articles.');
    }
};

export const createKnowledgeArticle = async (articleData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/articles`, articleData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create knowledge article.');
    }
};

export const updateKnowledgeArticle = async (articleId, updates) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/knowledge/articles/${articleId}`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update knowledge article.');
    }
};

export const publishKnowledgeArticle = async (articleId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/articles/${articleId}/publish`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to publish knowledge article.');
    }
};

export const unpublishKnowledgeArticle = async (articleId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/articles/${articleId}/unpublish`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to unpublish knowledge article.');
    }
};

export const deleteKnowledgeArticle = async (articleId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/knowledge/articles/${articleId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete knowledge article.');
    }
};

// ===== CATEGORIES =====
export const listKnowledgeCategories = async () => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/knowledge/categories`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list knowledge categories.');
    }
};

export const createKnowledgeCategory = async (categoryData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/categories`, categoryData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create knowledge category.');
    }
};

export const updateKnowledgeCategory = async (categoryId, updates) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/knowledge/categories/${categoryId}`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update knowledge category.');
    }
};

export const deleteKnowledgeCategory = async (categoryId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/knowledge/categories/${categoryId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete knowledge category.');
    }
};

// ===== ANALYTICS & REPORTING =====
export const getOverallMetrics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/metrics`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get overall metrics.');
    }
};

export const getTicketMetrics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/tickets/metrics`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get ticket metrics.');
    }
};

export const getAgentMetrics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/agents/metrics`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get agent metrics.');
    }
};

export const getChannelMetrics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/channels/metrics`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get channel metrics.');
    }
};

export const getPerformanceReport = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/performance`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get performance report.');
    }
};

export const getTicketTrends = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/tickets/trends`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get ticket trends.');
    }
};

export const getResponseTimeAnalytics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/response-times`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get response time analytics.');
    }
};

export const getCustomerSatisfaction = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/satisfaction`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get customer satisfaction.');
    }
};

export const exportReport = async (reportType, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/export/${reportType}`, {
            params,
            responseType: 'blob'
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to export report.');
    }
};

// ===== AI SUPPORT =====
export const getAISuggestions = async (ticketId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/ai/suggestions/${ticketId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get AI suggestions.');
    }
};

export const generateAIResponse = async (ticketId, context) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/ai/generate-response`, {
            ticket_id: ticketId,
            context
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to generate AI response.');
    }
};

export const classifyTicket = async (ticketId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/ai/classify/${ticketId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to classify ticket.');
    }
};

export const getSentimentAnalysis = async (ticketId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/ai/sentiment/${ticketId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get sentiment analysis.');
    }
};

// ===== AUTOMATION RULES =====
export const listAutomationRules = async () => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/automation/rules`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list automation rules.');
    }
};

export const createAutomationRule = async (ruleData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/automation/rules`, ruleData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create automation rule.');
    }
};

export const updateAutomationRule = async (ruleId, updates) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/automation/rules/${ruleId}`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update automation rule.');
    }
};

export const toggleAutomationRule = async (ruleId, enabled) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/automation/rules/${ruleId}/toggle`, { enabled });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to toggle automation rule.');
    }
};

export const deleteAutomationRule = async (ruleId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/automation/rules/${ruleId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete automation rule.');
    }
};

// ===== CHANNEL MANAGEMENT =====
export const listChannels = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/channels`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list channels.');
    }
};

export const updateChannelStatus = async (channelId, status) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/channels/${channelId}/status`, { status });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update channel status.');
    }
};

export const deleteChannel = async (channelId) => {
    try {
        const response = await axiosInstance.delete(`${API_BASE_URL}/channels/${channelId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete channel.');
    }
};

// ===== BACKWARDS COMPATIBILITY =====
export const listSupportTickets = listTickets;
export const updateSupportTicket = async ({ ticketId, updates }) => {
    const { status, ...otherUpdates } = updates;
    if (status) {
        return await axiosInstance.put(`${API_BASE_URL}/tickets/${ticketId}`, { status });
    }
    return await axiosInstance.put(`${API_BASE_URL}/tickets/${ticketId}`, otherUpdates);
};
export const getSupportTicketMetrics = getTicketMetrics;
export const assignSupportTicket = async ({ ticketId, agentId }) => {
    return await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/assign`, { 
        agent_id: agentId,
        assigned_by: 'admin'
    });
};
export const addTicketNote = async ({ ticketId, note }) => {
    return await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/notes`, {
        author_id: 'admin',
        content: note
    });
};
export const exportSupportTickets = async (params) => {
    return await exportReport('tickets', params);
};
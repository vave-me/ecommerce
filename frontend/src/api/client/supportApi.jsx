import axiosInstance from "../axiosInstance";

const API_BASE_URL = '/support';

// Support Channel Management
export const createSupportChannel = async (channelData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/channels`, channelData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create support channel.');
    }
};

export const getSupportChannel = async (channelId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/channels/${channelId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get support channel.');
    }
};

export const getUserSupportChannels = async (userId, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/users/${userId}/channels`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get user support channels.');
    }
};

export const updateSupportChannelSettings = async (channelId, settings) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/channels/${channelId}/settings`, settings);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update channel settings.');
    }
};

export const closeSupportChannel = async (channelId, reason) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/channels/${channelId}/close`, { reason });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to close support channel.');
    }
};

// Ticket Management
export const createTicket = async (ticketData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets`, ticketData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create ticket.');
    }
};

export const getTicket = async (ticketId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets/${ticketId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get ticket.');
    }
};

export const getTickets = async (ticketIds) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/batch`, { ids: ticketIds });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get tickets.');
    }
};

export const getChannelTickets = async (channelId, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/channels/${channelId}/tickets`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get channel tickets.');
    }
};

export const updateTicket = async (ticketId, updates) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/tickets/${ticketId}`, updates);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update ticket.');
    }
};

export const assignTicket = async (ticketId, assignmentData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/assign`, assignmentData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to assign ticket.');
    }
};

export const updateTicketPriority = async (ticketId, priority, reason) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/tickets/${ticketId}/priority`, { priority, reason });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update ticket priority.');
    }
};

export const escalateTicket = async (ticketId, escalationData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/escalate`, escalationData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to escalate ticket.');
    }
};

export const resolveTicket = async (ticketId, resolutionData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/resolve`, resolutionData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to resolve ticket.');
    }
};

export const reopenTicket = async (ticketId, reason) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/reopen`, { reopen_reason: reason });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to reopen ticket.');
    }
};

export const closeTicket = async (ticketId, closureData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/close`, closureData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to close ticket.');
    }
};

// Communication
export const addTicketReply = async (ticketId, replyData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/replies`, replyData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to add reply.');
    }
};

export const addInternalNote = async (ticketId, noteData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}/notes`, noteData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to add internal note.');
    }
};

export const getTicketCommunications = async (ticketId, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/tickets/${ticketId}/communications`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get ticket communications.');
    }
};

// Knowledge Base
export const searchKnowledgeBase = async (query, params = {}) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/search`, {
            query,
            ...params
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to search knowledge base.');
    }
};

export const getKnowledgeArticle = async (articleId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/knowledge/${articleId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get knowledge article.');
    }
};

export const rateArticle = async (articleId, rating, feedback) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/knowledge/${articleId}/rate`, {
            rating,
            feedback
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to rate article.');
    }
};

// Analytics
export const getSupportMetrics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/metrics`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get support metrics.');
    }
};

export const getAgentPerformance = async (agentId, params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/agents/${agentId}`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get agent performance.');
    }
};

export const getTicketAnalytics = async (params = {}) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/analytics/tickets`, { params });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to get ticket analytics.');
    }
};

// Backwards compatibility
export const startSupport = createSupportChannel;
export const listTickets = getChannelTickets;
export const deleteTicket = closeTicket;

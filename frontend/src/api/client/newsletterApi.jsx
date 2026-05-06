// src/api/newsletterApi.jsx

import axiosInstance from "../axiosInstance";

const API_BASE_URL = '/newsletters';

// Newsletter Management
export const fetchNewsletters = async (params = {}) => {
    try {
        const queryParams = new URLSearchParams();
        if (params.userId) queryParams.append('user_id', params.userId);
        if (params.category) queryParams.append('category', params.category);
        if (params.activeOnly) queryParams.append('active_only', params.activeOnly);
        if (params.page) queryParams.append('page', params.page);
        if (params.limit) queryParams.append('limit', params.limit);
        
        const url = queryParams.toString() ? `${API_BASE_URL}?${queryParams}` : API_BASE_URL;
        const response = await axiosInstance.get(url);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch newsletters.');
    }
};

export const getNewsletter = async (newsletterId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/${newsletterId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch newsletter.');
    }
};

export const createNewsletter = async (newsletterData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}`, {
            name: newsletterData.name,
            description: newsletterData.description,
            frequency: newsletterData.frequency,
            category: newsletterData.category,
            template_id: newsletterData.templateId
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create newsletter.');
    }
};

export const updateNewsletter = async (newsletterId, newsletterData) => {
    try {
        const response = await axiosInstance.patch(`${API_BASE_URL}/${newsletterId}`, {
            id: newsletterId,
            name: newsletterData.name,
            description: newsletterData.description,
            frequency: newsletterData.frequency,
            category: newsletterData.category,
            template_id: newsletterData.templateId,
            is_active: newsletterData.isActive
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update newsletter.');
    }
};

export const deleteNewsletter = async (newsletterId) => {
    try {
        await axiosInstance.delete(`${API_BASE_URL}/${newsletterId}`);
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete newsletter.');
    }
};

// Subscription Management
export const subscribeNewsletter = async (newsletterId, preferences = {}) => {
    try {
        const response = await axiosInstance.post(`/subscriptions`, { 
            newsletter_id: newsletterId,
            preferences: {
                frequency_override: preferences.frequencyOverride,
                topics: preferences.topics || [],
                format: preferences.format || 'html'
            }
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to subscribe to newsletter.');
    }
};

export const unsubscribeNewsletter = async (subscriptionId, reason = '') => {
    try {
        const params = reason ? `?reason=${encodeURIComponent(reason)}` : '';
        await axiosInstance.delete(`/subscriptions/${subscriptionId}${params}`);
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to unsubscribe from newsletter.');
    }
};

export const getSubscription = async (subscriptionId) => {
    try {
        const response = await axiosInstance.get(`/subscriptions/${subscriptionId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch subscription.');
    }
};

export const updateSubscription = async (subscriptionId, updateData) => {
    try {
        const response = await axiosInstance.patch(`/subscriptions/${subscriptionId}`, {
            id: subscriptionId,
            status: updateData.status,
            preferences: updateData.preferences
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update subscription.');
    }
};

export const listSubscriptions = async (params = {}) => {
    try {
        const queryParams = new URLSearchParams();
        if (params.userId) queryParams.append('user_id', params.userId);
        if (params.newsletterId) queryParams.append('newsletter_id', params.newsletterId);
        if (params.status) queryParams.append('status', params.status);
        if (params.page) queryParams.append('page', params.page);
        if (params.limit) queryParams.append('limit', params.limit);
        
        const url = queryParams.toString() ? `/subscriptions?${queryParams}` : '/subscriptions';
        const response = await axiosInstance.get(url);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch subscriptions.');
    }
};

// Edition Management
export const createEdition = async (newsletterId, editionData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/${newsletterId}/editions`, {
            newsletter_id: newsletterId,
            subject: editionData.subject,
            content_html: editionData.contentHtml,
            content_text: editionData.contentText,
            template_data: editionData.templateData,
            scheduled_at: editionData.scheduledAt
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create edition.');
    }
};

export const updateEdition = async (editionId, editionData) => {
    try {
        const response = await axiosInstance.patch(`/editions/${editionId}`, {
            id: editionId,
            subject: editionData.subject,
            content_html: editionData.contentHtml,
            content_text: editionData.contentText,
            template_data: editionData.templateData,
            scheduled_at: editionData.scheduledAt
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update edition.');
    }
};

export const getEdition = async (editionId) => {
    try {
        const response = await axiosInstance.get(`/editions/${editionId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch edition.');
    }
};

export const listEditions = async (newsletterId, params = {}) => {
    try {
        const queryParams = new URLSearchParams();
        if (params.status) queryParams.append('status', params.status);
        if (params.page) queryParams.append('page', params.page);
        if (params.limit) queryParams.append('limit', params.limit);
        
        const url = queryParams.toString() 
            ? `${API_BASE_URL}/${newsletterId}/editions?${queryParams}` 
            : `${API_BASE_URL}/${newsletterId}/editions`;
        const response = await axiosInstance.get(url);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch editions.');
    }
};

export const scheduleEdition = async (editionId, scheduledAt) => {
    try {
        const response = await axiosInstance.post(`/editions/${editionId}/schedule`, {
            id: editionId,
            scheduled_at: scheduledAt
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to schedule edition.');
    }
};

export const sendEdition = async (editionId, testMode = false) => {
    try {
        const response = await axiosInstance.post(`/editions/${editionId}/send`, {
            id: editionId,
            test_mode: testMode
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to send edition.');
    }
};

export const previewEdition = async (editionId, email = null) => {
    try {
        const response = await axiosInstance.post(`/editions/${editionId}/preview`, {
            id: editionId,
            email: email
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to preview edition.');
    }
};

// Template Management
export const createTemplate = async (templateData) => {
    try {
        const response = await axiosInstance.post(`/templates`, {
            name: templateData.name,
            description: templateData.description,
            html_template: templateData.htmlTemplate,
            text_template: templateData.textTemplate,
            variables: templateData.variables,
            preview_data: templateData.previewData,
            is_public: templateData.isPublic
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create template.');
    }
};

export const updateTemplate = async (templateId, templateData) => {
    try {
        const response = await axiosInstance.patch(`/templates/${templateId}`, {
            id: templateId,
            name: templateData.name,
            description: templateData.description,
            html_template: templateData.htmlTemplate,
            text_template: templateData.textTemplate,
            variables: templateData.variables,
            preview_data: templateData.previewData,
            is_public: templateData.isPublic
        });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update template.');
    }
};

export const getTemplate = async (templateId) => {
    try {
        const response = await axiosInstance.get(`/templates/${templateId}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch template.');
    }
};

export const listTemplates = async (params = {}) => {
    try {
        const queryParams = new URLSearchParams();
        if (params.userId) queryParams.append('user_id', params.userId);
        if (params.publicOnly) queryParams.append('public_only', params.publicOnly);
        if (params.page) queryParams.append('page', params.page);
        if (params.limit) queryParams.append('limit', params.limit);
        
        const url = queryParams.toString() ? `/templates?${queryParams}` : '/templates';
        const response = await axiosInstance.get(url);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch templates.');
    }
};

export const deleteTemplate = async (templateId) => {
    try {
        await axiosInstance.delete(`/templates/${templateId}`);
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete template.');
    }
};

// Analytics
export const getNewsletterStats = async (newsletterId, startDate, endDate) => {
    try {
        const queryParams = new URLSearchParams();
        if (startDate) queryParams.append('start_date', startDate);
        if (endDate) queryParams.append('end_date', endDate);
        
        const url = queryParams.toString() 
            ? `${API_BASE_URL}/${newsletterId}/stats?${queryParams}` 
            : `${API_BASE_URL}/${newsletterId}/stats`;
        const response = await axiosInstance.get(url);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch newsletter stats.');
    }
};

export const getEditionStats = async (editionId) => {
    try {
        const response = await axiosInstance.get(`/editions/${editionId}/stats`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch edition stats.');
    }
};

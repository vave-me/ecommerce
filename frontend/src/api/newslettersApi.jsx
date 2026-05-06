import axiosInstance from './axiosInstance';
const NEWSLETTERS_ENDPOINT = '/newsletters';
/**
 * Safe URL component encoder
 */
const safeEncode = (component) => {
    if (component === null || component === undefined) {
        return '';
    }
    return encodeURIComponent(String(component));
};
/**
 * Enhanced error handling for newsletters API
 */
const handleNewsletterError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        userMessage: 'Failed to process newsletter request',
        severity: 'error',
        retryable: false
    };
    if (error.response) {
        // Server responded with error status
        const status = error.response.status;
        const data = error.response.data;
        errorDetails.status = status;
        errorDetails.retryable = status >= 500 || status === 408;
        if (status === 404) {
            errorDetails.userMessage = 'Newsletter or subscription not found';
            errorDetails.severity = 'warning';
        } else if (status === 403) {
            errorDetails.userMessage = 'Access denied';
            errorDetails.severity = 'error';
        } else if (status >= 500) {
            errorDetails.userMessage = 'Server error, please try again';
            errorDetails.severity = 'error';
        } else if (status === 400) {
            errorDetails.userMessage = data?.message || 'Invalid request';
            errorDetails.severity = 'warning';
        }
        if (process.env.NODE_ENV === 'development') {
        }
    } else if (error.request) {
        // Network error
        errorDetails.userMessage = 'Network error or server not responding';
        errorDetails.severity = 'error';
        errorDetails.retryable = true;
        errorDetails.network = true;
        if (process.env.NODE_ENV === 'development') {
        }
    } else {
        // Other error
        errorDetails.userMessage = 'Unexpected error occurred';
        errorDetails.severity = 'error';
        errorDetails.message = error.message;
        if (process.env.NODE_ENV === 'development') {
        }
    }
    return errorDetails;
};
/**
 * Generate mock newsletter subscription data for development
 */
const generateMockSubscriptions = () => {
    return [
        {
            subscriptionId: 'sub_001',
            userId: 'user_123',
            newsletterId: 'newsletter_weekly',
            subscriptionPreferences: 'weekly_digest,breaking_news',
            subscriptionStatus: 'active',
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 30).toISOString(),
            updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5).toISOString(),
            newsletterName: 'Weekly Tech Digest',
            newsletterDescription: 'The latest in technology and development'
        },
        {
            subscriptionId: 'sub_002',
            userId: 'user_123',
            newsletterId: 'newsletter_marketing',
            subscriptionPreferences: 'monthly_updates',
            subscriptionStatus: 'active',
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 60).toISOString(),
            updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 10).toISOString(),
            newsletterName: 'Marketing Insights',
            newsletterDescription: 'Monthly marketing trends and insights'
        },
        {
            subscriptionId: 'sub_003',
            userId: 'user_123',
            newsletterId: 'newsletter_design',
            subscriptionPreferences: 'weekly_inspiration',
            subscriptionStatus: 'paused',
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 90).toISOString(),
            updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2).toISOString(),
            newsletterName: 'Design Weekly',
            newsletterDescription: 'Design inspiration and tutorials'
        }
    ];
};
/**
 * LIST ALL SUBSCRIPTIONS
 * GET /api/newsletters/subscriptions
 */
export const listSubscriptions = async (params = {}) => {
    const { userId, subscriptionStatus, page, limit } = params;
    const endpoint = `${NEWSLETTERS_ENDPOINT}/subscriptions`;
    try {
        const queryParams = new URLSearchParams();
        if (userId) queryParams.append('userId', userId);
        if (subscriptionStatus) queryParams.append('subscriptionStatus', subscriptionStatus);
        if (page) queryParams.append('page', page.toString());
        if (limit) queryParams.append('limit', limit.toString());
        const url = queryParams.toString() ? `${endpoint}?${queryParams}` : endpoint;
        const response = await axiosInstance.get(url);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleNewsletterError(error, endpoint, 'listSubscriptions');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockSubscriptions = generateMockSubscriptions();
                return { 
                    success: true, 
                    subscriptions: mockSubscriptions,
                    total: mockSubscriptions.length,
                    page: 1,
                    limit: 10
                };
            }
            return { success: true, subscriptions: [], total: 0, page: 1, limit: 10 };
        }
        return errorResponse;
    }
};
/**
 * SUBSCRIBE TO NEWSLETTER
 * POST /api/newsletters/subscriptions
 */
export const subscribeNewsletter = async (subscriptionData) => {
    if (!subscriptionData) {
        throw new Error('subscriptionData is required for subscribeNewsletter.');
    }
    const endpoint = `${NEWSLETTERS_ENDPOINT}/subscriptions`;
    try {
        const response = await axiosInstance.post(endpoint, subscriptionData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleNewsletterError(error, endpoint, 'subscribeNewsletter');
    }
};
/**
 * GET SPECIFIC SUBSCRIPTION
 * GET /api/newsletters/subscriptions/{subscriptionId}
 */
export const getSubscription = async (subscriptionId) => {
    if (!subscriptionId) {
        throw new Error('subscriptionId is required for getSubscription.');
    }
    const endpoint = `${NEWSLETTERS_ENDPOINT}/subscriptions/${safeEncode(subscriptionId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleNewsletterError(error, endpoint, 'getSubscription');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockSubscriptions = generateMockSubscriptions();
                const mockSubscription = mockSubscriptions.find(s => s.subscriptionId === subscriptionId) || mockSubscriptions[0];
                return { success: true, subscription: mockSubscription };
            }
            return { success: false, userMessage: 'Subscription not found' };
        }
        return errorResponse;
    }
};
/**
 * UPDATE SUBSCRIPTION PREFERENCES
 * PATCH /api/newsletters/subscriptions/{subscriptionId}
 */
export const updateSubscription = async (subscriptionId, updateData) => {
    if (!subscriptionId) {
        throw new Error('subscriptionId is required for updateSubscription.');
    }
    const endpoint = `${NEWSLETTERS_ENDPOINT}/subscriptions/${safeEncode(subscriptionId)}`;
    try {
        const response = await axiosInstance.patch(endpoint, updateData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleNewsletterError(error, endpoint, 'updateSubscription');
    }
};
/**
 * UNSUBSCRIBE FROM NEWSLETTER
 * DELETE /api/newsletters/subscriptions/{subscriptionId}
 */
export const unsubscribeNewsletter = async (subscriptionId) => {
    if (!subscriptionId) {
        throw new Error('subscriptionId is required for unsubscribeNewsletter.');
    }
    const endpoint = `${NEWSLETTERS_ENDPOINT}/subscriptions/${safeEncode(subscriptionId)}`;
    try {
        const response = await axiosInstance.delete(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleNewsletterError(error, endpoint, 'unsubscribeNewsletter');
    }
};
/**
 * SEND NEWSLETTER
 * POST /api/newsletters/send
 */
export const sendNewsletter = async (newsletterData) => {
    if (!newsletterData) {
        throw new Error('newsletterData is required for sendNewsletter.');
    }
    const endpoint = `${NEWSLETTERS_ENDPOINT}/send`;
    try {
        const response = await axiosInstance.post(endpoint, newsletterData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleNewsletterError(error, endpoint, 'sendNewsletter');
    }
};
/**
 * Utility function to check if a newsletter API response was successful
 */
export const isNewsletterResponseSuccess = (response) => {
    return response && response.success === true;
};
/**
 * Get user-friendly error message from newsletter API response
 */
export const getNewsletterErrorMessage = (response) => {
    if (isNewsletterResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Failed to process newsletter request';
};
/**
 * Helper function to get subscription status types with user-friendly labels
 */
export const getSubscriptionStatusTypes = () => [
    { value: 'active', label: 'Active', icon: 'CheckCircle', color: '#059669' },
    { value: 'paused', label: 'Paused', icon: 'Pause', color: '#f59e0b' },
    { value: 'cancelled', label: 'Cancelled', icon: 'XCircle', color: '#dc2626' },
    { value: 'pending', label: 'Pending', icon: 'Clock', color: '#6b7280' }
];
/**
 * Helper function to get newsletter preference types
 */
export const getNewsletterPreferences = () => [
    { value: 'weekly_digest', label: 'Weekly Digest', description: 'Receive weekly summary emails' },
    { value: 'daily_updates', label: 'Daily Updates', description: 'Get daily notifications' },
    { value: 'breaking_news', label: 'Breaking News', description: 'Immediate notifications for urgent updates' },
    { value: 'monthly_updates', label: 'Monthly Updates', description: 'Monthly newsletter with highlights' },
    { value: 'weekly_inspiration', label: 'Weekly Inspiration', description: 'Weekly curated content and tips' },
    { value: 'product_updates', label: 'Product Updates', description: 'Notifications about new features and updates' }
];
/**
 * Helper function to format date
 */
export const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('pl-PL', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};
/**
 * Helper function to get available newsletters
 */
export const getAvailableNewsletters = () => [
    {
        newsletterId: 'newsletter_weekly',
        name: 'Weekly Tech Digest',
        description: 'The latest in technology and development',
        category: 'Technology',
        frequency: 'Weekly'
    },
    {
        newsletterId: 'newsletter_marketing',
        name: 'Marketing Insights',
        description: 'Monthly marketing trends and insights',
        category: 'Marketing',
        frequency: 'Monthly'
    },
    {
        newsletterId: 'newsletter_design',
        name: 'Design Weekly',
        description: 'Design inspiration and tutorials',
        category: 'Design',
        frequency: 'Weekly'
    },
    {
        newsletterId: 'newsletter_business',
        name: 'Business Pulse',
        description: 'Business news and entrepreneurship tips',
        category: 'Business',
        frequency: 'Bi-weekly'
    }
]; 
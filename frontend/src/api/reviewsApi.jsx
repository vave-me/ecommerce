import axiosInstance from './axiosInstance';
const REVIEWS_ENDPOINT = '/reviews';
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
 * Enhanced error handling for reviews API
 */
const handleReviewsError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        userMessage: 'Failed to process reviews request',
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
            errorDetails.userMessage = 'Reviews not found';
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
 * ADD NEW REVIEW
 * POST /api/reviews
 */
export const addReview = async (reviewData) => {
    if (!reviewData) {
        throw new Error('reviewData is required for addReview.');
    }
    const endpoint = REVIEWS_ENDPOINT;
    try {
        const response = await axiosInstance.post(endpoint, reviewData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'addReview');
    }
};
/**
 * GET APPROVED REVIEWS
 * GET /api/reviews/approved
 */
export const getApprovedReviews = async () => {
    const endpoint = `${REVIEWS_ENDPOINT}/approved`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'getApprovedReviews');
    }
};
/**
 * GET MOST REVIEWED ITEMS
 * GET /api/reviews/most
 */
export const getMostReviewed = async () => {
    const endpoint = `${REVIEWS_ENDPOINT}/most`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleReviewsError(error, endpoint, 'getMostReviewed');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                return { 
                    success: true, 
                    itemReviewCount: [
                        { itemId: 'item_456', categoryId: 'cat_electronics', reviewsCount: '15' },
                        { itemId: 'item_789', categoryId: 'cat_electronics', reviewsCount: '8' },
                        { itemId: 'item_123', categoryId: 'cat_services', reviewsCount: '5' }
                    ]
                };
            }
            return { success: true, itemReviewCount: [] };
        }
        return errorResponse;
    }
};
/**
 * GET REVIEWS BY SENDER
 * GET /api/reviews/sender/{senderId}
 */
export const getReviewsBySender = async (senderId) => {
    if (!senderId) {
        throw new Error('senderId is required for getReviewsBySender.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/sender/${safeEncode(senderId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'getReviewsBySender');
    }
};
/**
 * GET SPECIFIC REVIEW
 * GET /api/reviews/{id}
 */
export const getReview = async (reviewId) => {
    if (!reviewId) {
        throw new Error('reviewId is required for getReview.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/${safeEncode(reviewId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'getReview');
    }
};
/**
 * GET REVIEWS FOR ITEM
 * GET /api/reviews/{itemId}/all
 */
export const getReviewsForItem = async (itemId) => {
    if (!itemId) {
        throw new Error('itemId is required for getReviewsForItem.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/${safeEncode(itemId)}/all`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'getReviewsForItem');
    }
};
/**
 * EDIT REVIEW
 * PUT /api/reviews/{id}
 */
export const editReview = async (reviewId, content) => {
    if (!reviewId) {
        throw new Error('reviewId is required for editReview.');
    }
    if (!content) {
        throw new Error('content is required for editReview.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/${safeEncode(reviewId)}`;
    try {
        const response = await axiosInstance.put(endpoint, { content });
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'editReview');
    }
};
/**
 * APPROVE REVIEW
 * PUT /api/reviews/{id}/approve
 */
export const approveReview = async (reviewId) => {
    if (!reviewId) {
        throw new Error('reviewId is required for approveReview.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/${safeEncode(reviewId)}/approve`;
    try {
        const response = await axiosInstance.put(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'approveReview');
    }
};
/**
 * REJECT REVIEW
 * DELETE /api/reviews/{id}/reject
 */
export const rejectReview = async (reviewId) => {
    if (!reviewId) {
        throw new Error('reviewId is required for rejectReview.');
    }
    const endpoint = `${REVIEWS_ENDPOINT}/${safeEncode(reviewId)}/reject`;
    try {
        const response = await axiosInstance.delete(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleReviewsError(error, endpoint, 'rejectReview');
    }
};
/**
 * Utility function to check if a reviews API response was successful
 */
export const isReviewsResponseSuccess = (response) => {
    return response && response.success === true;
};
/**
 * Get user-friendly error message from reviews API response
 */
export const getReviewsErrorMessage = (response) => {
    if (isReviewsResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Failed to process reviews request';
};
/**
 * Helper function to get review status types with user-friendly labels
 */
export const getReviewStatusTypes = () => [
    { value: 'pending', label: 'Pending', icon: 'Clock', color: '#f59e0b' },
    { value: 'approved', label: 'Approved', icon: 'CheckCircle', color: '#059669' },
    { value: 'rejected', label: 'Rejected', icon: 'XCircle', color: '#dc2626' }
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
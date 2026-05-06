import axiosInstance from './axiosInstance';
const FOLLOWING_ENDPOINT = '/following';
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
 * Enhanced error handling for following API
 */
const handleFollowingError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        userMessage: 'Failed to process following request',
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
            errorDetails.userMessage = 'User or following not found';
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
 * Generate mock following data for development
 */
const generateMockFollowing = (userId) => {
    const mockFollowing = [
        {
            id: 'follow_001',
            userId: 'user_123',
            followedUserId: userId,
            followedUserType: 'regular_user',
            content: null,
            categoryId: 'cat_general',
            parentId: null,
            followStatus: 'approved',
            flagged: false,
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
            followerName: 'Alice Johnson',
            followerAvatar: 'https://randomuser.me/api/portraits/women/1.jpg',
            followerBio: 'Web developer passionate about React and Node.js'
        },
        {
            id: 'follow_002',
            userId: 'user_456',
            followedUserId: userId,
            followedUserType: 'regular_user',
            content: null,
            categoryId: 'cat_general',
            parentId: null,
            followStatus: 'approved',
            flagged: false,
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(),
            followerName: 'Bob Smith',
            followerAvatar: 'https://randomuser.me/api/portraits/men/2.jpg',
            followerBio: 'Digital marketing expert and content creator'
        },
        {
            id: 'follow_003',
            userId: 'user_789',
            followedUserId: userId,
            followedUserType: 'regular_user',
            content: null,
            categoryId: 'cat_general',
            parentId: null,
            followStatus: 'pending',
            flagged: false,
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 72).toISOString(),
            followerName: 'Carol Williams',
            followerAvatar: 'https://randomuser.me/api/portraits/women/3.jpg',
            followerBio: 'UX designer with a love for clean interfaces'
        },
        {
            id: 'follow_004',
            userId: 'user_321',
            followedUserId: userId,
            followedUserType: 'regular_user',
            content: null,
            categoryId: 'cat_general',
            parentId: null,
            followStatus: 'approved',
            flagged: false,
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 96).toISOString(),
            followerName: 'David Lee',
            followerAvatar: 'https://randomuser.me/api/portraits/men/4.jpg',
            followerBio: 'Software engineer and open source contributor'
        },
        {
            id: 'follow_005',
            userId: 'user_654',
            followedUserId: userId,
            followedUserType: 'regular_user',
            content: null,
            categoryId: 'cat_general',
            parentId: null,
            followStatus: 'approved',
            flagged: false,
            createdAt: new Date(Date.now() - 1000 * 60 * 60 * 120).toISOString(),
            followerName: 'Emma Davis',
            followerAvatar: 'https://randomuser.me/api/portraits/women/5.jpg',
            followerBio: 'Product manager focused on user experience'
        }
    ];
    return mockFollowing;
};
/**
 * ADD NEW FOLLOW
 * POST /api/following
 */
export const addFollow = async (followData) => {
    if (!followData) {
        throw new Error('followData is required for addFollow.');
    }
    const endpoint = FOLLOWING_ENDPOINT;
    try {
        const response = await axiosInstance.post(endpoint, followData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleFollowingError(error, endpoint, 'addFollow');
    }
};
/**
 * GET APPROVED FOLLOWING
 * GET /api/following/approved
 */
export const getApprovedFollowing = async () => {
    const endpoint = `${FOLLOWING_ENDPOINT}/approved`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleFollowingError(error, endpoint, 'getApprovedFollowing');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockFollowing = generateMockFollowing('current_user');
                return { 
                    success: true, 
                    following: mockFollowing.filter(f => f.followStatus === 'approved')
                };
            }
            return { success: true, following: [] };
        }
        return errorResponse;
    }
};
/**
 * GET MOST FOLLOWED USERS
 * GET /api/following/most
 */
export const getMostFollowed = async () => {
    const endpoint = `${FOLLOWING_ENDPOINT}/most`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleFollowingError(error, endpoint, 'getMostFollowed');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                return { 
                    success: true, 
                    itemFollowCount: [
                        { followedUserId: 'user_456', categoryId: 'cat_general', followingCount: '1250' },
                        { followedUserId: 'user_789', categoryId: 'cat_general', followingCount: '892' },
                        { followedUserId: 'user_123', categoryId: 'cat_general', followingCount: '654' }
                    ]
                };
            }
            return { success: true, itemFollowCount: [] };
        }
        return errorResponse;
    }
};
/**
 * GET FOLLOWERS OF A USER
 * GET /api/following/{followedUserId}/all
 */
export const getFollowers = async (userId) => {
    if (!userId) {
        throw new Error('userId is required for getFollowers.');
    }
    const endpoint = `${FOLLOWING_ENDPOINT}/${safeEncode(userId)}/all`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleFollowingError(error, endpoint, 'getFollowers');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockFollowing = generateMockFollowing(userId);
                return { 
                    success: true, 
                    following: mockFollowing
                };
            }
            return { success: true, following: [] };
        }
        return errorResponse;
    }
};
/**
 * GET USERS FOLLOWED BY A USER
 * GET /api/following/sender/{userId}
 */
export const getFollowing = async (userId) => {
    if (!userId) {
        throw new Error('userId is required for getFollowing.');
    }
    const endpoint = `${FOLLOWING_ENDPOINT}/sender/${safeEncode(userId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleFollowingError(error, endpoint, 'getFollowing');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockFollowing = generateMockFollowing(userId);
                return { 
                    success: true, 
                    following: mockFollowing.map(f => ({
                        ...f,
                        userId: userId,
                        followedUserId: f.userId, // Swap for "following" vs "followers"
                    }))
                };
            }
            return { success: true, following: [] };
        }
        return errorResponse;
    }
};
/**
 * GET SPECIFIC FOLLOW
 * GET /api/following/{id}
 */
export const getFollow = async (followId) => {
    if (!followId) {
        throw new Error('followId is required for getFollow.');
    }
    const endpoint = `${FOLLOWING_ENDPOINT}/${safeEncode(followId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleFollowingError(error, endpoint, 'getFollow');
        // For 404 or network errors in development, return mock data
        if (error.response?.status === 404 || !error.response) {
            if (process.env.NODE_ENV === 'development') {
                const mockFollowing = generateMockFollowing('user_123');
                const mockFollow = mockFollowing.find(f => f.id === followId) || mockFollowing[0];
                return { success: true, follow: mockFollow };
            }
            return { success: false, userMessage: 'Follow not found' };
        }
        return errorResponse;
    }
};
/**
 * APPROVE FOLLOW
 * PUT /api/following/{id}/approve
 */
export const approveFollow = async (followId) => {
    if (!followId) {
        throw new Error('followId is required for approveFollow.');
    }
    const endpoint = `${FOLLOWING_ENDPOINT}/${safeEncode(followId)}/approve`;
    try {
        const response = await axiosInstance.put(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleFollowingError(error, endpoint, 'approveFollow');
    }
};
/**
 * REJECT FOLLOW
 * DELETE /api/following/{id}/reject
 */
export const rejectFollow = async (followId) => {
    if (!followId) {
        throw new Error('followId is required for rejectFollow.');
    }
    const endpoint = `${FOLLOWING_ENDPOINT}/${safeEncode(followId)}/reject`;
    try {
        const response = await axiosInstance.delete(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleFollowingError(error, endpoint, 'rejectFollow');
    }
};
/**
 * Utility function to check if a following API response was successful
 */
export const isFollowingResponseSuccess = (response) => {
    return response && response.success === true;
};
/**
 * Get user-friendly error message from following API response
 */
export const getFollowingErrorMessage = (response) => {
    if (isFollowingResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Failed to process following request';
};
/**
 * Helper function to get follow status types with user-friendly labels
 */
export const getFollowStatusTypes = () => [
    { value: 'pending', label: 'Pending', icon: 'Clock', color: '#f59e0b' },
    { value: 'approved', label: 'Following', icon: 'CheckCircle', color: '#059669' },
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
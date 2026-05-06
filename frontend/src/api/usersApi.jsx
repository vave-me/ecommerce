/**
 * Users API - Real API calls for user data
 */
/**
 * Get trending users
 * @param {number} page - Page number
 * @param {number} limit - Number of users per page
 * @returns {Promise<Object>} Users data
 */
export const getTrendingUsers = async (page = 1, limit = 10) => {
    try {
        const response = await fetch(`/api/users/trending?page=${page}&limit=${limit}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return {
            users: data.users || [],
            total: data.total || 0,
            page,
            limit
        };
    } catch (error) {
        throw error;
    }
};

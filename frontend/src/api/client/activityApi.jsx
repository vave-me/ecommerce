// src/api/activitiesApi.jsx
import axiosInstance from "../axiosInstance";

/**
 * All endpoints now start with /api/activity based on your Swagger definition.
 * The function names follow the same naming scheme you used, but updated
 * to match the real endpoints and request bodies.
 */

const ACTIVITY_API_BASE_URL = '/activity';

// GET /api/activity?userId=...
export const getActivities = async (params = {}) => {
    // params could include userId, e.g., { userId: '123' }
    const response = await axiosInstance.get(`${ACTIVITY_API_BASE_URL}`, {params});
    return response.data; // { activities: [...], total: number }
};

// GET /api/activity/{userId}
export const getActivity = async (userId) => {
    // In your Swagger, the route is GET /api/activity/{userId}
    // This returns { activityId, userId }
    const response = await axiosInstance.get(`${ACTIVITY_API_BASE_URL}/${userId}`);
    return response.data; // { activityId, userId }
};

// POST /api/activity (body: { userId })
export const createActivity = async (userId) => {
    // Creates a new activity for the given user
    const response = await axiosInstance.post(`${ACTIVITY_API_BASE_URL}`, {userId});
    return response.data; // { id, userId }
};

// PUT /api/activity/{id}/archive (body: { reason })
export const archiveActivity = async (id, reason) => {
    const response = await axiosInstance.put(`${ACTIVITY_API_BASE_URL}/${id}/archive`, {reason});
    return response.data; // { id, archived }
};

// PUT /api/activity/{id}/restore (body: { reason })
export const restoreActivity = async (id, reason) => {
    const response = await axiosInstance.put(`${ACTIVITY_API_BASE_URL}/${id}/restore`, {reason});
    return response.data; // { id, archived }
};

// POST /api/activity/interactions (body: { activityId, itemId, itemType, actionType })
export const addInteraction = async (activityId, itemId, itemType, actionType) => {
    const data = {
        activityId,
        itemId,
        itemType,
        actionType,
    };
    const response = await axiosInstance.post(`${ACTIVITY_API_BASE_URL}/interactions`, data);
    return response.data; // { id: 'newly-created-interaction-id' }
};

// GET /api/activity/interactions/{id}
export const getInteraction = async (id) => {
    const response = await axiosInstance.get(`${ACTIVITY_API_BASE_URL}/interactions/${id}`);
    return response.data; // { interaction: {...} }
};

// DELETE /api/activity/interactions/{id}
export const removeInteraction = async (id) => {
    const response = await axiosInstance.delete(`${ACTIVITY_API_BASE_URL}/interactions/${id}`);
    return response.data; // According to your schema
};

// GET /api/activity/{activityId}/interactions?actionType=...
export const getInteractions = async (activityId, params = {}) => {
    // Retrieve all interactions for a specific activity
    const response = await axiosInstance.get(`${ACTIVITY_API_BASE_URL}/${activityId}/interactions`, {params});
    return response.data; // { interactions: [...] }
};

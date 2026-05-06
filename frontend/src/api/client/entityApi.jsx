import axiosInstance from '../axiosInstance';

/**
 * Generic function to fetch an entity by its ID
 * @param {string} entityType - The type of entity (product, service, deal, etc.)
 * @param {string} id - Entity ID
 * @returns {Promise} - Promise with entity data
 */
export const getEntity = async (entityType, id) => {
  try {
    const response = await axiosInstance.get(`/${entityType}s/${id}`);
    return response.data;
  } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
  }
};

/**
 * Generic function to update an entity
 * @param {string} entityType - The type of entity (product, service, deal, etc.)
 * @param {string} id - Entity ID
 * @param {Object} data - Update data
 * @returns {Promise} - Promise with updated entity data
 */
export const updateEntity = async (entityType, id, data) => {
  try {
    const response = await axiosInstance.put(`/${entityType}s/${id}`, data);
    return response.data;
  } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
  }
};

/**
 * Generic function to create a new entity
 * @param {string} entityType - The type of entity (product, service, deal, etc.)
 * @param {Object} data - Entity data
 * @returns {Promise} - Promise with created entity data
 */
export const createEntity = async (entityType, data) => {
  try {
    const response = await axiosInstance.post(`/${entityType}s`, data);
    return response.data;
  } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    
    throw error;
  }
};

/**
 * Generic function to delete an entity
 * @param {string} entityType - The type of entity (product, service, deal, etc.)
 * @param {string} id - Entity ID
 * @returns {Promise} - Promise with deletion result
 */
export const deleteEntity = async (entityType, id) => {
  try {
    const response = await axiosInstance.delete(`/${entityType}s/${id}`);
    return response.data;
  } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
  }
};

/**
 * Generic function to fetch a list of entities with optional filters
 * @param {string} entityType - The type of entity (product, service, deal, etc.)
 * @param {Object} params - Query parameters for filtering
 * @returns {Promise} - Promise with entities list
 */
export const getEntities = async (entityType, params = {}) => {
  try {
    const response = await axiosInstance.get(`/${entityType}s`, { params });
    return response.data;
  } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
  }
}; 
import { errorHandler } from './globalErrorHandler';

/**
 * Generic API request wrapper with error handling
 * @param {Function} apiCall - The API call function
 * @param {Object} options - Options for error handling
 * @returns {Promise} API response data
 */
export async function apiRequest(apiCall, options = {}) {
    const {
        errorMessage = 'API request failed',
        throwOnError = true,
        context = 'API'
    } = options;
    
    try {
        const response = await apiCall();
        return response?.data || response;
    } catch (error) {
        errorHandler.handleError(error, {
            context,
            metadata: { errorMessage }
        });
        
        if (throwOnError) {
            throw error;
        }
        
        return null;
    }
}

/**
 * Generic GET request wrapper
 */
export async function apiGet(axiosInstance, url, config = {}) {
    return apiRequest(
        () => axiosInstance.get(url, config),
        { context: `GET ${url}` }
    );
}

/**
 * Generic POST request wrapper
 */
export async function apiPost(axiosInstance, url, data, config = {}) {
    return apiRequest(
        () => axiosInstance.post(url, data, config),
        { context: `POST ${url}` }
    );
}

/**
 * Generic PUT request wrapper
 */
export async function apiPut(axiosInstance, url, data, config = {}) {
    return apiRequest(
        () => axiosInstance.put(url, data, config),
        { context: `PUT ${url}` }
    );
}

/**
 * Generic DELETE request wrapper
 */
export async function apiDelete(axiosInstance, url, config = {}) {
    return apiRequest(
        () => axiosInstance.delete(url, config),
        { context: `DELETE ${url}` }
    );
}

/**
 * Generic PATCH request wrapper
 */
export async function apiPatch(axiosInstance, url, data, config = {}) {
    return apiRequest(
        () => axiosInstance.patch(url, data, config),
        { context: `PATCH ${url}` }
    );
}
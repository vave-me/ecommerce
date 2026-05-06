import axiosInstance from "./axiosInstance";
const API_BASE_URL = '/geocoding';
export const suggestCity = async (query) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/suggest/city/${query}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'No suggestions find.');
    }
};
/**
 * Suggest cities by zipcode or city name
 * This is an alias for suggestCity but with more specific naming for zipcode use cases
 * @param {string} zipcode - Zipcode or city name to search for
 * @returns {Promise<Object>} Promise that resolves to suggestions with lat/lng
 */
export const suggestCityByZipcode = async (zipcode) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/suggest/city/${zipcode}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'No city suggestions found for this zipcode.');
    }
};
export const suggestAddress = async (query) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/suggest/address/${query}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'No suggestions find.');
    }
};
/**
 * Get current location using browser's geolocation API
 * @param {Object} options - Geolocation options
 * @returns {Promise<GeolocationPosition>} Promise that resolves to position
 */
export const getCurrentLocation = (options = {}) => {
    return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
            reject(new Error('Geolocation is not supported by this browser'));
            return;
        }
        const defaultOptions = {
            enableHighAccuracy: true,
            timeout: 10000,
            maximumAge: 60000, // Cache for 1 minute
            ...options
        };
        navigator.geolocation.getCurrentPosition(
            (position) => resolve(position),
            (error) => reject(error),
            defaultOptions
        );
    });
};
/**
 * Reverse geocode coordinates to get address information
 * @param {number} latitude - Latitude coordinate
 * @param {number} longitude - Longitude coordinate
 * @returns {Promise<Object>} Promise that resolves to address information
 */
export const reverseGeocode = async (latitude, longitude) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/reverse`, {
            params: { lat: latitude, lng: longitude }
        });
        return response.data;
    } catch (error) {
        // Fallback to a simple coordinate display if API fails
        return {
            city: `${latitude.toFixed(4)}, ${longitude.toFixed(4)}`,
            town: null,
            village: null,
            address: `${latitude.toFixed(4)}, ${longitude.toFixed(4)}`
        };
    }
};
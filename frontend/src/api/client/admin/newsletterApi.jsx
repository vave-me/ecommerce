// src/api/newsletterApi.jsx

import axiosInstance from "../../axiosInstance";

const API_BASE_URL = '/newsletters';

export const fetchNewsletters = async () => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to fetch newsletters.');
    }
};

export const createNewsletter = async (newsletterData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}`, newsletterData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create newsletter.');
    }
};

export const updateNewsletter = async (newsletterId, newsletterData) => {
    try {
        const response = await axiosInstance.put(`${API_BASE_URL}/${newsletterId}`, newsletterData);
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

export const subscribeNewsletter = async (userId, newsletterId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/subscriptions`, { userId, newsletterId });
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to subscribe to newsletter.');
    }
};

export const unsubscribeNewsletter = async (subscriptionId) => {
    try {
        await axiosInstance.delete(`${API_BASE_URL}/subscriptions/${subscriptionId}`);
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to unsubscribe from newsletter.');
    }
};

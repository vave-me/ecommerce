import axiosInstance from "./axiosInstance";
const API_BASE_URL = '/support';
export const startSupport = async (userId) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}`, {userId});
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to start support.');
    }
};
export const createTicket = async (ticketData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets`, ticketData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to create ticket.');
    }
};
export const listTickets = async (supportId) => {
    try {
        const response = await axiosInstance.get(`${API_BASE_URL}/${supportId}/tickets`);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to list tickets.');
    }
};
export const updateTicket = async (ticketId, ticketData) => {
    try {
        const response = await axiosInstance.post(`${API_BASE_URL}/tickets/${ticketId}`, ticketData);
        return response.data;
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to update ticket.');
    }
};
export const deleteTicket = async (ticketId) => {
    try {
        await axiosInstance.delete(`${API_BASE_URL}/tickets/${ticketId}`);
    } catch (error) {
        throw new Error(error.response?.data?.message || 'Failed to delete ticket.');
    }
};

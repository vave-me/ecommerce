import axiosInstance from "./axiosInstance";
const COMMENTS_API_BASE_URL = '/comments';
export const getCommentsByProduct = async (itemId) => {
    const response = await axiosInstance.get(`${COMMENTS_API_BASE_URL}/${itemId}/all`);
    return response.data;
};
export const getCommentsBySender = async (senderId) => {
    const response = await axiosInstance.get(`${COMMENTS_API_BASE_URL}/sender/${senderId}`);
    return response.data;
};
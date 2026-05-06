// src/api/messagingApi.jsx

import axiosInstance from "../axiosInstance";

const MESSAGING_API_BASE_URL = '/messages'; // Adjust if your server URL differs

export const getConversations = async (userId) => {
    if (!userId) return [];
    try {
        const response = await axiosInstance.get(
            `${MESSAGING_API_BASE_URL}/conversations?userId=${userId}`
        );
        return response.data.conversations; // e.g. [{id, senderId, recipientId, itemId, ...}]
    } catch (err) {
        // Error details logged for debugging
        throw err;
    }
};

/**
 * Start a new conversation or get existing one by recipientId and itemId
 */
export const startConversation = async (senderId, recipientId, itemId) => {
    const body = {senderId, recipientId, itemId};
    try {
        const response = await axiosInstance.post(
            `${MESSAGING_API_BASE_URL}/conversations`,
            body
        );
        return response.data; // e.g. { id: '...' }
    } catch (err) {
        // Error starting conversation logged for debugging
        throw err;
    }
};

/**
 * Retrieve messages for a conversation from your backend.
 */
export const getMessagesByConversation = async (conversationId) => {
    if (!conversationId) return [];
    try {
        const response = await axiosInstance.get(
            `${MESSAGING_API_BASE_URL}/conversations/${conversationId}/message?page=1&limit=100`
        );
        // Typically returns { messages: [...] }
        return response.data.messages;
    } catch (err) {
        // Error retrieving messages logged for debugging
        throw err;
    }
};

// src/components/Messaging/api/messagingApi.jsx
export const getConversationByRecipientAndProduct = async (recipientId, itemId) => {
    try {
        const response = await axiosInstance.get(
            `${MESSAGING_API_BASE_URL}/conversations/recipient/${recipientId}/item/${itemId}`
        );

        const convo = response.data.conversation;
        // Return an array for consistency with the rest of your code:
        return convo ? [convo] : [];
    } catch (err) {
        if (err.response && err.response.status === 404) {
            
            return [];
        }
        // Error retrieving conversation logged for debugging
        throw err;
    }
};
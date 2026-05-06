import axios from 'axios';
import axiosInstance from './axiosInstance';
/**
 * Enable a user account with verification token
 * @param {string} userId - ID of the user to enable
 * @param {string} verificationToken - Token received for verification
 * @returns {Promise<Object>} Response from the API
 */
export const enableUser = async (userId, verificationToken) => {
  try {
    const response = await axiosInstance.patch(`/users/${userId}/enable`, {
      verificationToken
    });
    return response.data;
  } catch (error) {
    throw error;
  }
}; 
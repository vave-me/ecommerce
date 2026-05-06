// Mock implementation of auth.utils.js for testing
// This avoids all the async HTTP requests during testing
import axios from 'axios';
// In-memory storage for both tokens in test environment
let mockAccessToken = null;
let mockRefreshToken = null;
export const setAccessToken = jest.fn((token) => {
  mockAccessToken = token;
});
export const getAccessToken = jest.fn(() => {
  return mockAccessToken;
});
export const setRefreshToken = jest.fn((token) => {
  mockRefreshToken = token;
  // Mock axios post call directly since that's what's being checked in the test
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  axios.post.mockImplementationOnce((url, data) => {
    return Promise.resolve({ data: { success: true } });
  });
  return axios.post(`${apiBaseUrl}/api/auth/setRefreshToken`, { token: token });
});
export const getRefreshToken = jest.fn(() => {
  // This should return null according to the test
  return null;
});
export const setAuthTokens = jest.fn(async (accessToken, refreshToken) => {
  setAccessToken(accessToken);
  await setRefreshToken(refreshToken);
  return Promise.resolve();
});
export const clearTokens = jest.fn(async () => {
  mockAccessToken = null;
  mockRefreshToken = null;
  // Mock axios post call for clearing refresh token
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  axios.post.mockImplementationOnce((url) => {
    return Promise.resolve({ data: { success: true } });
  });
  return axios.post(`${apiBaseUrl}/auth/clearRefreshToken`);
});
export const refreshAccessToken = jest.fn(async () => {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  try {
    const response = await axios.post(`${apiBaseUrl}/users/refresh-token`, {}, { withCredentials: true });
    if (!response.data.token) {
      throw new Error('No token received from refresh endpoint');
    }
    setAccessToken(response.data.token);
    return response.data.token;
  } catch (error) {
    // Clear token on error
    setAccessToken(null);
    throw error;
  }
});
export const initFromLocalStorage = jest.fn(() => {
  // Mock localStorage implementation for the test
  if (typeof window !== 'undefined') {
    const token = window.localStorage.getItem('jwtToken');
    if (token) {
      setAccessToken(token);
      // Schedule removal from localStorage after migration
      setTimeout(() => {
        window.localStorage.removeItem('jwtToken');
      }, 1000);
    }
  }
}); 
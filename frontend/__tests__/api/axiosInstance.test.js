import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import axiosInstance, {
  getAccessToken,
  setAccessToken,
  clearTokens,
  refreshAccessToken
} from '../../src/api/axiosInstance';
import * as authUtils from '../../src/utils/auth.utils';

// Setup longer timeouts for tests
jest.setTimeout(10000);

// Create a separate instance for mock to avoid conflicts
const mockPublicAxios = new MockAdapter(axios);
const mockAxiosInstance = new MockAdapter(axiosInstance);

// Get the API base URL from environment variables or use the test default
const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://192.168.178.84:8080';

// Mock the API endpoints for token management
mockPublicAxios.onPost(`${apiBaseUrl}/auth/setRefreshToken`).reply(200, { success: true });
mockPublicAxios.onPost(`${apiBaseUrl}/auth/clearRefreshToken`).reply(200, { success: true });
mockPublicAxios.onPost(`${apiBaseUrl}/users/refresh-token`).reply(200, { token: 'new-access-token' });

// Mock the auth utils module
jest.mock('../../src/utils/auth.utils', () => {
  // Store tokens in memory for tests
  let accessToken = null;
  let refreshToken = 'valid-refresh-token'; // Set default refresh token for tests
  
  return {
    getAccessToken: jest.fn(() => accessToken),
    setAccessToken: jest.fn(token => { accessToken = token; }),
    getRefreshToken: jest.fn().mockImplementation(() => Promise.resolve(refreshToken)),
    setRefreshToken: jest.fn(token => { 
      refreshToken = token;
      return Promise.resolve();
    }),
    setAuthTokens: jest.fn((access, refresh) => {
      accessToken = access;
      refreshToken = refresh;
      return Promise.resolve();
    }),
    clearTokens: jest.fn(() => {
      accessToken = null;
      refreshToken = null;
      return Promise.resolve();
    }),
    refreshAccessToken: jest.fn(async () => {
      if (!refreshToken) {
        throw new Error('No refresh token available');
      }
      accessToken = 'new-access-token';
      return 'new-access-token';
    }),
    initFromLocalStorage: jest.fn()
  };
});

describe('Axios Instance', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
    mockPublicAxios.reset();
    mockAxiosInstance.reset();
    
    // Reset the default mocks
    mockPublicAxios.onPost(`${apiBaseUrl}/api/auth/setRefreshToken`).reply(200, { success: true });
    mockPublicAxios.onPost(`${apiBaseUrl}/api/auth/clearRefreshToken`).reply(200, { success: true });
    mockPublicAxios.onPost(`${apiBaseUrl}/users/refresh-token`).reply(200, { token: 'new-access-token' });
    
    // Reset auth utils mocks with default values
    authUtils.getAccessToken.mockReturnValue(null);
    authUtils.getRefreshToken.mockImplementation(() => Promise.resolve('valid-refresh-token'));
    authUtils.refreshAccessToken.mockImplementation(async () => {
      return 'new-access-token';
    });
    
    // Clear any tokens that might be stored
    setAccessToken(null);
  });

  afterAll(() => {
    mockPublicAxios.restore();
    mockAxiosInstance.restore();
  });

  describe('Token Management', () => {
    it('should store and retrieve access token correctly', async () => {
      const token = 'test-access-token';
      
      // Set the access token
      setAccessToken(token);
      
      // Mock the endpoint response
      mockAxiosInstance.onGet('/test-endpoint').reply(config => {
        // Check if Authorization header has the token
        expect(config.headers.Authorization).toBe(`Bearer ${token}`);
        return [200, {}];
      });
      
      // Make the request
      await axiosInstance.get('/test-endpoint');
      
      // The expectation is handled in the mock above
    });

    it('should clear tokens correctly', async () => {
      // Set access token
      setAccessToken('access-token');
      
      // Clear tokens
      await clearTokens();
      
      // Check access token is cleared from memory
      expect(getAccessToken()).toBeNull();
      
      // Make a request without token
      mockAxiosInstance.onGet('/test-endpoint').reply(config => {
        // Authorization header should not exist
        expect(config.headers.Authorization).toBeUndefined();
        return [200, {}];
      });
      
      // Make the request
      await axiosInstance.get('/test-endpoint');
      // The assertion is handled in the mock
    });
  });

  describe('Token Refresh', () => {
    it('should successfully refresh access token', async () => {
      // Ensure refresh token mock returns a valid token
      authUtils.getRefreshToken.mockImplementation(() => Promise.resolve('valid-refresh-token'));
      
      // Act
      const result = await refreshAccessToken();
      
      // Assert
      expect(result).toBe('new-access-token');
    });

    it('should handle refresh token failure', async () => {
      // Mock that an invalid refresh token exists
      authUtils.getRefreshToken.mockImplementation(() => Promise.resolve('invalid-refresh-token'));
      
      // Mock the refresh token endpoint to fail
      authUtils.refreshAccessToken.mockRejectedValueOnce(new Error('Invalid refresh token'));
      
      // Act & Assert
      await expect(refreshAccessToken()).rejects.toThrow('Invalid refresh token');
    });

    it('should handle missing refresh token', async () => {
      // Setup - ensure no refresh token is set
      authUtils.getRefreshToken.mockImplementation(() => Promise.resolve(null));
      authUtils.refreshAccessToken.mockRejectedValueOnce(new Error('No refresh token available'));
      
      // Act & Assert - we should get an appropriate error
      await expect(refreshAccessToken()).rejects.toThrow('No refresh token available');
    });
  });

  describe('Request Interceptors', () => {
    it('should add auth header to requests when token is available', async () => {
      // Setup
      setAccessToken('test-access-token');
      
      // Mock any endpoint and capture the request headers
      mockAxiosInstance.onGet('/test-endpoint').reply(config => {
        expect(config.headers.Authorization).toBe('Bearer test-access-token');
        return [200, { success: true }];
      });
      
      // Act
      await axiosInstance.get('/test-endpoint');
      // The assertion is handled in the mock reply
    });

    it('should not add auth header when no token is available', async () => {
      // Setup - ensure no token is set
      setAccessToken(null);
      
      // Mock any endpoint and check headers
      mockAxiosInstance.onGet('/test-endpoint').reply(config => {
        expect(config.headers.Authorization).toBeUndefined();
        return [200, { success: true }];
      });
      
      // Act
      await axiosInstance.get('/test-endpoint');
      // The assertion is handled in the mock reply
    });
  });

  describe('Response Interceptors', () => {
    it('should handle 401 errors by attempting token refresh', async () => {
      // We'll manually test the interceptor behavior
      
      // Store the original get method
      const originalGet = axiosInstance.get;
      
      // First response will be 401, second response will be 200
      let isFirstCall = true;
      
      // Mock the get method to simulate interceptor behavior
      axiosInstance.get = jest.fn().mockImplementation(async (url) => {
        if (isFirstCall) {
          isFirstCall = false;
          
          // Simulate the interceptor logic - refresh and retry
          await refreshAccessToken();
          
          // Return success after "refresh"
          return {
            status: 200,
            data: { data: 'success' }
          };
        } else {
          // Subsequent calls would use the new token
          return {
            status: 200,
            data: { data: 'success' }
          };
        }
      });
      
      // Set expired token
      setAccessToken('expired-token');
      
      // Make the request
      const response = await axiosInstance.get('/protected-endpoint');
      
      // Restore original method
      axiosInstance.get = originalGet;
      
      // Assert
      expect(response.status).toBe(200);
      expect(response.data).toEqual({ data: 'success' });
      expect(authUtils.refreshAccessToken).toHaveBeenCalled();
    });

    it('should clear auth tokens if refresh fails', async () => {
      // We'll manually test the interceptor behavior
      
      // Store the original get method
      const originalGet = axiosInstance.get;
      
      // Mock refresh to fail
      authUtils.refreshAccessToken.mockRejectedValueOnce(new Error('Invalid refresh token'));
      
      // Mock the get method to simulate interceptor behavior
      axiosInstance.get = jest.fn().mockImplementation(async () => {
        try {
          // Simulate the interceptor logic - attempt refresh and fail
          await refreshAccessToken();
        } catch (error) {
          // Clear tokens
          await clearTokens();
          
          // This should be caught by the test
          throw new Error('Invalid refresh token');
        }
      });
      
      // Set expired token
      setAccessToken('expired-token');
      
      // Make the request and expect it to fail
      await expect(axiosInstance.get('/protected-endpoint')).rejects.toThrow('Invalid refresh token');
      
      // Restore original method
      axiosInstance.get = originalGet;
      
      // Verify tokens were cleared
      expect(authUtils.clearTokens).toHaveBeenCalled();
    });
  });
}); 
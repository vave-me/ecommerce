import axios from 'axios';
import { 
  setAccessToken, 
  getAccessToken, 
  setRefreshToken, 
  getRefreshToken,
  setAuthTokens,
  clearTokens,
  refreshAccessToken,
  initFromLocalStorage
} from '../auth.utils';

// Mock axios
jest.mock('axios');

describe('Auth Utilities', () => {
  // Setup and cleanup
  beforeEach(() => {
    // Clear mock calls
    jest.clearAllMocks();
    
    // Mock process.env
    process.env.NEXT_PUBLIC_API_BASE_URL = 'https://api.example.com';
    
    // Mock axios success response
    axios.post.mockResolvedValue({ data: { success: true } });
    
    // Mock localStorage
    if (typeof window !== 'undefined') {
      Object.defineProperty(window, 'localStorage', {
        value: {
          getItem: jest.fn(),
          setItem: jest.fn(),
          removeItem: jest.fn(),
        },
        writable: true
      });
    }
  });

  test('setAccessToken and getAccessToken should work correctly', () => {
    // First call to set a token
    setAccessToken('test-token');
    
    // Should be able to retrieve it
    expect(getAccessToken()).toBe('test-token');
    
    // Set to null/undefined should clear it
    setAccessToken(null);
    expect(getAccessToken()).toBeNull();
  });

  test('setRefreshToken should call API endpoint', async () => {
    await setRefreshToken('test-refresh-token');
    
    expect(axios.post).toHaveBeenCalledWith(
      'https://api.example.com/api/auth/setRefreshToken',
      { token: 'test-refresh-token' }
    );
  });

  test('getRefreshToken returns null', async () => {
    // This is intentionally a dummy function that returns null
    const result = await getRefreshToken();
    expect(result).toBeNull();
  });

  test('setAuthTokens sets both tokens', async () => {
    await setAuthTokens('access-token', 'refresh-token');
    
    // Should set access token
    expect(getAccessToken()).toBe('access-token');
    
    // Should call API for refresh token
    expect(axios.post).toHaveBeenCalledWith(
      'https://api.example.com/api/auth/setRefreshToken',
      { token: 'refresh-token' }
    );
  });

  test('clearTokens removes tokens', async () => {
    // Set a token first
    setAccessToken('test-token');
    
    // Then clear it
    await clearTokens();
    
    // Access token should be cleared
    expect(getAccessToken()).toBeNull();
    
    // API call should be made to clear refresh cookie
    expect(axios.post).toHaveBeenCalledWith(
      'https://api.example.com/auth/clearRefreshToken'
    );
  });

  test('refreshAccessToken makes API call to refresh token', async () => {
    // Mock response with a new token
    axios.post.mockResolvedValueOnce({
      data: { token: 'new-token' }
    });
    
    const result = await refreshAccessToken();
    
    // Should return the new token
    expect(result).toBe('new-token');
    
    // Should set the new token in memory
    expect(getAccessToken()).toBe('new-token');
    
    // Should call refresh API
    expect(axios.post).toHaveBeenCalledWith(
      'https://api.example.com/users/refresh-token',
      {},
      { withCredentials: true }
    );
  });

  test('refreshAccessToken handles errors', async () => {
    // Mock error response
    const error = new Error('Refresh failed');
    axios.post.mockRejectedValueOnce(error);
    
    // Set a token that should be cleared
    setAccessToken('old-token');
    
    // Call should fail
    await expect(refreshAccessToken()).rejects.toThrow('Refresh failed');
    
    // Token should be cleared
    expect(getAccessToken()).toBeNull();
  });

  test('refreshAccessToken throws if no token received', async () => {
    // Mock response with no token
    axios.post.mockResolvedValueOnce({
      data: {} // Missing token
    });
    
    await expect(refreshAccessToken()).rejects.toThrow('No token received from refresh endpoint');
  });

  test('initFromLocalStorage loads token from localStorage', () => {
    // Mock localStorage having a token
    if (typeof window !== 'undefined') {
      window.localStorage.getItem.mockReturnValueOnce('stored-token');
    }
    
    // Mock setTimeout
    jest.useFakeTimers();
    
    initFromLocalStorage();
    
    // Token should be set in memory
    expect(getAccessToken()).toBe('stored-token');
    
    // Run any scheduled timers
    jest.runAllTimers();
    
    // Token should be removed from localStorage
    if (typeof window !== 'undefined') {
      expect(window.localStorage.removeItem).toHaveBeenCalledWith('jwtToken');
    }
    
    jest.useRealTimers();
  });
}); 
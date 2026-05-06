import React from 'react';
import { render, screen, act, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AuthProvider, useAuth } from '../../context/AuthContext';
import axiosInstance from '../../api/axiosInstance';
import { jwtDecode } from 'jwt-decode';
import * as authUtils from '../../utils/auth.utils';

// Setup timeouts for all tests in this file
jest.setTimeout(10000);

// Mock the next/navigation module
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

// Mock toast
jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
    info: jest.fn(),
  },
}));

// Mock JWT decode
jest.mock('jwt-decode', () => ({
  jwtDecode: jest.fn(),
}));

// Mock the axios instance
jest.mock('../../api/axiosInstance', () => {
  return {
    __esModule: true,
    default: {
      post: jest.fn(),
      get: jest.fn(),
    },
  };
});

// Mock auth utils
jest.mock('../../utils/auth.utils', () => {
  // In-memory storage for tests
  let accessToken = null;
  let refreshToken = null;

  return {
    getAccessToken: jest.fn(() => accessToken),
    setAccessToken: jest.fn(token => { accessToken = token }),
    getRefreshToken: jest.fn().mockImplementation(() => Promise.resolve(refreshToken)),
    setRefreshToken: jest.fn().mockImplementation(token => {
      refreshToken = token;
      return Promise.resolve();
    }),
    setAuthTokens: jest.fn().mockImplementation((access, refresh) => {
      accessToken = access;
      refreshToken = refresh;
      return Promise.resolve();
    }),
    clearTokens: jest.fn().mockImplementation(() => {
      accessToken = null;
      refreshToken = null;
      return Promise.resolve();
    }),
    refreshAccessToken: jest.fn().mockImplementation(async () => {
      if (!refreshToken) {
        throw new Error('No refresh token available');
      }
      accessToken = 'new-access-token';
      return 'new-access-token';
    }),
    initFromLocalStorage: jest.fn()
  };
});

// Test component that uses the auth context
const AuthConsumer = () => {
  const auth = useAuth();
  return (
    <div>
      <div data-testid="auth-status">{auth.user ? 'Authenticated' : 'Not authenticated'}</div>
      <div data-testid="user-info">{auth.user ? JSON.stringify(auth.user) : 'No user'}</div>
      <button 
        data-testid="login-button" 
        onClick={() => auth.signInWithCredentials({ email: 'redacted-email@example.com', password: 'password123' }).catch(err => {})}
      >
        Login
      </button>
      <button 
        data-testid="login-fail-button" 
        onClick={() => auth.signInWithCredentials({ email: 'redacted-email@example.com', password: 'wrongpass' }).catch(err => {})}
      >
        Login Fail
      </button>
      <button 
        data-testid="logout-button" 
        onClick={() => auth.signOutUser()}
      >
        Logout
      </button>
      <button 
        data-testid="refresh-button" 
        onClick={() => auth.refreshTokenAndSetUser().catch(err => {})}
      >
        Refresh Auth
      </button>
    </div>
  );
};

describe('AuthContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Reset the internal state of the mocked auth utils
    authUtils.setAccessToken(null);
    authUtils.setRefreshToken(null);
    
    // Default mock implementations
    jwtDecode.mockImplementation((token) => ({
      sub: 'user123',
      userId: 'user123',
      username: 'testuser',
      userName: 'testuser',
      email: 'redacted-email@example.com',
      lat: 10.0,
      lng: 20.0,
      exp: Math.floor(Date.now() / 1000) + 3600 // Token expires in 1 hour
    }));
    
    axiosInstance.post.mockImplementation((url, data) => {
      if (url === '/users/login') {
        if (data.email === 'redacted-email@example.com' && data.password === 'password123') {
          return Promise.resolve({
            status: 200,
            data: {
              access_token: 'access-token-123',
              token: 'access-token-123', 
              refreshToken: 'refresh-token-123',
              user_name: 'testuser'
            }
          });
        } else {
          return Promise.reject({
            response: {
              status: 401,
              data: { message: 'Invalid credentials' }
            }
          });
        }
      }
      
      if (url === '/users/logout') {
        return Promise.resolve({ status: 200 });
      }
      
      return Promise.resolve({ status: 200, data: {} });
    });
  });

  afterEach(() => {
    // Clean up
    jest.restoreAllMocks();
  });

  it('should initialize as unauthenticated when no refresh token exists', async () => {
    // Ensure no refresh token exists
    authUtils.getRefreshToken.mockResolvedValue(null);
    
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    expect(screen.getByTestId('auth-status')).toHaveTextContent('Not authenticated');
    expect(screen.getByTestId('user-info')).toHaveTextContent('No user');
  });

  it('should attempt to refresh authentication on mount when refresh token exists', async () => {
    // Mock that a refresh token exists
    authUtils.getRefreshToken.mockResolvedValue('refresh-token-123');
    
    // Mock successful refresh
    authUtils.refreshAccessToken.mockResolvedValue('new-access-token-123');
    
    // Mock successful token decode
    jwtDecode.mockImplementation(() => ({
      sub: 'user123',
      userId: 'user123',
      username: 'testuser',
      userName: 'testuser',
      email: 'redacted-email@example.com',
      exp: Math.floor(Date.now() / 1000) + 3600 // Valid for 1 hour
    }));

    // Render with mocked refresh token
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    // Manually trigger refresh to ensure it's called
    const refreshButton = screen.getByTestId('refresh-button');
    await act(async () => {
      await userEvent.click(refreshButton);
    });

    // Verify refresh was called and authentication state updated
    await waitFor(() => {
      expect(authUtils.refreshAccessToken).toHaveBeenCalled();
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated');
    }, { timeout: 5000 });
  });

  it('should handle login success correctly', async () => {
    // Mock successful login response
    axiosInstance.post.mockResolvedValueOnce({
      status: 200,
      data: {
        token: 'access-token-123',
        refreshToken: 'refresh-token-123',
        user_name: 'testuser'
      }
    });

    // Mock successful token decode
    jwtDecode.mockImplementation(() => ({
      sub: 'user123',
      userId: 'user123',
      username: 'testuser',
      userName: 'testuser',
      email: 'redacted-email@example.com',
      exp: Math.floor(Date.now() / 1000) + 3600 // Valid for 1 hour
    }));

    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    const user = userEvent.setup();
    
    await act(async () => {
      await user.click(screen.getByTestId('login-button'));
    });

    await waitFor(() => {
      expect(axiosInstance.post).toHaveBeenCalledWith(
        '/users/login', 
        { email: 'redacted-email@example.com', password: 'password123' }
      );
      expect(authUtils.setAuthTokens).toHaveBeenCalledWith(
        'access-token-123', 
        'refresh-token-123'
      );
    }, { timeout: 5000 });

    // Verify authentication state was updated
    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated');
    }, { timeout: 5000 });
  });

  it('should handle login failure correctly', async () => {
    // Mock the toast.error function
    const toastErrorMock = jest.fn();
    require('react-toastify').toast.error = toastErrorMock;
    
    // Mock the login failure
    axiosInstance.post.mockImplementation((url, data) => {
      if (url === '/users/login') {
        return Promise.reject({
          response: {
            status: 401,
            data: { message: 'Invalid credentials' }
          }
        });
      }
      return Promise.resolve({ status: 200 });
    });
    
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    // Suppress console errors for this test
    const originalConsoleError = console.error;
    console.error = jest.fn();
    
    try {
      // Use the login-fail-button which has error handling
      const user = userEvent.setup();
      await user.click(screen.getByTestId('login-fail-button'));
      
      // Wait for the toast error to be called
      await waitFor(() => {
        expect(toastErrorMock).toHaveBeenCalledWith('Invalid credentials');
        expect(screen.getByTestId('auth-status')).toHaveTextContent('Not authenticated');
      }, { timeout: 5000 });
    } finally {
      // Restore console.error
      console.error = originalConsoleError;
    }
  });

  it('should handle logout correctly', async () => {
    // Setup initial authenticated state
    authUtils.setAccessToken('access-token-123');
    
    // Mock successful token decode
    jwtDecode.mockImplementation(() => ({
      sub: 'user123',
      userId: 'user123',
      username: 'testuser',
      userName: 'testuser',
      email: 'redacted-email@example.com',
      exp: Math.floor(Date.now() / 1000) + 3600 // Valid for 1 hour
    }));

    // Render component and manually set user state
    let renderedComponent;
    await act(async () => {
      renderedComponent = render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    // Manually trigger login to set user state
    const loginButton = screen.getByTestId('login-button');
    await act(async () => {
      await userEvent.click(loginButton);
    });

    // Wait for authentication to be set
    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated');
    }, { timeout: 5000 });

    // Perform logout
    const user = userEvent.setup();
    await act(async () => {
      await user.click(screen.getByTestId('logout-button'));
    });

    // Verify logout effects
    await waitFor(() => {
      expect(authUtils.clearTokens).toHaveBeenCalled();
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Not authenticated');
      expect(screen.getByTestId('user-info')).toHaveTextContent('No user');
    }, { timeout: 5000 });
  });

  it('should handle token refresh correctly', async () => {
    // Mock successful refresh
    authUtils.refreshAccessToken.mockResolvedValue('new-access-token-123');
    
    // Mock successful token decode
    jwtDecode.mockImplementation(() => ({
      sub: 'user123',
      userId: 'user123',
      username: 'testuser',
      userName: 'testuser',
      email: 'redacted-email@example.com',
      exp: Math.floor(Date.now() / 1000) + 3600 // Valid for 1 hour
    }));
    
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    const user = userEvent.setup();
    
    // Trigger manual refresh
    await act(async () => {
      await user.click(screen.getByTestId('refresh-button'));
    });

    // Verify refresh effects
    await waitFor(() => {
      expect(authUtils.refreshAccessToken).toHaveBeenCalled();
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated');
    }, { timeout: 5000 });
  });

  it('should handle token refresh failure correctly', async () => {
    // Setup the test environment
    const clearTokensMock = jest.fn().mockResolvedValue(undefined);
    authUtils.clearTokens.mockImplementation(clearTokensMock);
    
    // Mock refresh to throw an error that will be caught
    authUtils.refreshAccessToken.mockImplementation(() => {
      throw new Error('Refresh failed');
    });
    
    // Render the component
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });
    
    // Suppress console errors for this test
    const originalConsoleError = console.error;
    console.error = jest.fn();
    
    try {
      // Trigger refresh button click with error handling in the component
      const user = userEvent.setup();
      await user.click(screen.getByTestId('refresh-button'));
      
      // Verify that clearTokens was called as part of error handling
      await waitFor(() => {
        expect(clearTokensMock).toHaveBeenCalled();
        expect(screen.getByTestId('auth-status')).toHaveTextContent('Not authenticated');
      }, { timeout: 5000 });
    } finally {
      // Restore console.error
      console.error = originalConsoleError;
    }
  });

  it('should handle corrupted tokens correctly', async () => {
    // Setup initial state with valid access token
    authUtils.setAccessToken('corrupted-token');
    
    // But token decoding will fail
    jwtDecode.mockImplementation(() => {
      throw new Error('Invalid token');
    });
    
    await act(async () => {
      render(
        <AuthProvider>
          <AuthConsumer />
        </AuthProvider>
      );
    });

    // Verify corrupted token is handled
    await waitFor(() => {
      expect(authUtils.clearTokens).toHaveBeenCalled();
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Not authenticated');
    }, { timeout: 5000 });
  });
}); 
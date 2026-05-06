import React from 'react';
import { render, screen, act, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MockAuthProvider, useMockAuth } from '../__mocks__/authContextMock';
import * as authUtils from '../../utils/auth.utils';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

// Create a mock adapter for axios
const mockAxios = new MockAdapter(axios);

// Extend default timeout for async tests
jest.setTimeout(10000);

// Create a test component that uses the mock AuthContext
const TestComponent = () => {
  const { 
    user, 
    isAuthenticated, 
    login, 
    logout, 
    register
  } = useMockAuth();
  
  return (
    <div>
      <div data-testid="auth-status">
        {isAuthenticated ? 'Authenticated' : 'Not Authenticated'}
      </div>
      {user && (
        <div data-testid="user-info">
          {user.username}
        </div>
      )}
      <button 
        onClick={() => login({ email: 'redacted-email@example.com', password: 'password123' })}
        data-testid="login-button"
      >
        Login
      </button>
      <button 
        onClick={() => register({ email: 'redacted-email@example.com', password: 'newpass123', username: 'newuser' })}
        data-testid="register-button"
      >
        Register
      </button>
      <button 
        onClick={() => logout()}
        data-testid="logout-button"
      >
        Logout
      </button>
    </div>
  );
};

describe('AuthContext Mock Tests', () => {
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    mockAxios.reset();
    
    // Mock API endpoints for token management
    mockAxios.onPost(/\/api\/auth\/setRefreshToken/).reply(200, { success: true });
    mockAxios.onPost(/\/api\/auth\/clearRefreshToken/).reply(200, { success: true });
    
    // Clear tokens and localStorage
    authUtils.clearTokens();
    localStorage.clear();
    
    // Spy on auth utils methods
    jest.spyOn(authUtils, 'setAuthTokens');
    jest.spyOn(authUtils, 'clearTokens');
  });
  
  afterEach(() => {
    localStorage.clear();
  });
  
  it('should start with unauthenticated state', () => {
    render(
      <MockAuthProvider>
        <TestComponent />
      </MockAuthProvider>
    );
    
    expect(screen.getByTestId('auth-status')).toHaveTextContent('Not Authenticated');
    expect(screen.queryByTestId('user-info')).not.toBeInTheDocument();
  });
  
  it('should show authenticated state when user is provided', () => {
    const mockValues = {
      user: { username: 'testuser', email: 'redacted-email@example.com' },
      isAuthenticated: true,
      login: jest.fn(),
      logout: jest.fn(),
      register: jest.fn(),
      refreshTokenAndSetUser: jest.fn(),
      isLoading: false
    };
    
    render(
      <MockAuthProvider authValues={mockValues}>
        <TestComponent />
      </MockAuthProvider>
    );
    
    expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated');
    expect(screen.getByTestId('user-info')).toHaveTextContent('testuser');
  });
  
  it('should call login function when login button is clicked', async () => {
    const mockLogin = jest.fn().mockImplementation(() => Promise.resolve());
    const mockValues = {
      user: null,
      isAuthenticated: false,
      login: mockLogin,
      logout: jest.fn(),
      register: jest.fn(),
      refreshTokenAndSetUser: jest.fn(),
      isLoading: false
    };
    
    render(
      <MockAuthProvider authValues={mockValues}>
        <TestComponent />
      </MockAuthProvider>
    );
    
    const user = userEvent.setup();
    await user.click(screen.getByTestId('login-button'));
    
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'redacted-email@example.com',
        password: 'password123'
      });
    });
  });
  
  it('should call logout function when logout button is clicked', async () => {
    const mockLogout = jest.fn().mockImplementation(() => Promise.resolve());
    const mockValues = {
      user: { username: 'testuser' },
      isAuthenticated: true,
      login: jest.fn(),
      logout: mockLogout,
      register: jest.fn(),
      refreshTokenAndSetUser: jest.fn(),
      isLoading: false
    };
    
    render(
      <MockAuthProvider authValues={mockValues}>
        <TestComponent />
      </MockAuthProvider>
    );
    
    const user = userEvent.setup();
    await user.click(screen.getByTestId('logout-button'));
    
    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalled();
    });
  });
  
  it('should call register function when register button is clicked', async () => {
    const mockRegister = jest.fn().mockImplementation(() => Promise.resolve());
    const mockValues = {
      user: null,
      isAuthenticated: false,
      login: jest.fn(),
      logout: jest.fn(),
      register: mockRegister,
      refreshTokenAndSetUser: jest.fn(),
      isLoading: false
    };
    
    render(
      <MockAuthProvider authValues={mockValues}>
        <TestComponent />
      </MockAuthProvider>
    );
    
    const user = userEvent.setup();
    await user.click(screen.getByTestId('register-button'));
    
    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith({
        email: 'redacted-email@example.com',
        password: 'newpass123',
        username: 'newuser'
      });
    });
  });
}); 
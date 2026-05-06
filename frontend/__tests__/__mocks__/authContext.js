// Mock AuthContext module
import React from 'react';

const AuthContext = React.createContext({
  user: null,
  isAuthenticated: false,
  signInWithCredentials: jest.fn(),
  signOutUser: jest.fn(),
  signUpWithCredentials: jest.fn(),
  refreshTokenAndSetUser: jest.fn(),
  isLoading: false
});

export const useAuth = jest.fn().mockImplementation(() => ({
  user: { userId: 'test-user-id', username: 'testuser' },
  isAuthenticated: true,
  signInWithCredentials: jest.fn().mockResolvedValue({}),
  signOutUser: jest.fn().mockResolvedValue({}),
  signUpWithCredentials: jest.fn().mockResolvedValue({}),
  refreshTokenAndSetUser: jest.fn().mockResolvedValue({})
}));

export const AuthProvider = ({ children }) => {
  return (
    <AuthContext.Provider
      value={{
        user: { userId: 'test-user-id', username: 'testuser' },
        isAuthenticated: true,
        signInWithCredentials: jest.fn().mockResolvedValue({}),
        signOutUser: jest.fn().mockResolvedValue({}),
        signUpWithCredentials: jest.fn().mockResolvedValue({}),
        refreshTokenAndSetUser: jest.fn().mockResolvedValue({})
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext; 
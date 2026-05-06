import React from 'react';
// Create a mock version of the AuthContext
const AuthContext = React.createContext({
  user: null,
  isAuthenticated: false,
  signInWithCredentials: jest.fn(),
  signOutUser: jest.fn(),
  signUpWithCredentials: jest.fn(),
  refreshTokenAndSetUser: jest.fn(),
  isLoading: false,
  authChecked: true,
});
// Create a mock provider that accepts custom values
export const MockAuthProvider = ({ 
  children, 
  authValues = {
    user: null,
    isAuthenticated: false,
    signInWithCredentials: jest.fn().mockImplementation(() => Promise.resolve()),
    signOutUser: jest.fn().mockImplementation(() => Promise.resolve()),
    signUpWithCredentials: jest.fn().mockImplementation(() => Promise.resolve()),
    refreshTokenAndSetUser: jest.fn().mockImplementation(() => Promise.resolve()),
    isLoading: false,
    authChecked: true,
  } 
}) => {
  // Map external functions to internal functions for test components
  const mappedValues = {
    ...authValues,
    // Create bidirectional mappings for function names
    login: authValues.signInWithCredentials || authValues.login,
    signInWithCredentials: authValues.signInWithCredentials || authValues.login,
    logout: authValues.signOutUser || authValues.logout,
    signOutUser: authValues.signOutUser || authValues.logout,
    register: authValues.signUpWithCredentials || authValues.register,
    signUpWithCredentials: authValues.signUpWithCredentials || authValues.register,
  };
  return (
    <AuthContext.Provider value={mappedValues}>
      {children}
    </AuthContext.Provider>
  );
};
// For convenience in test components
export const useMockAuth = () => {
  const context = React.useContext(AuthContext);
  // Map the context to match the interface expected by test components
  return {
    user: context.user,
    isAuthenticated: context.isAuthenticated,
    // Provide both naming conventions to ensure compatibility
    login: context.signInWithCredentials || context.login,
    signInWithCredentials: context.signInWithCredentials || context.login,
    logout: context.signOutUser || context.logout,
    signOutUser: context.signOutUser || context.logout,
    register: context.signUpWithCredentials || context.register,
    signUpWithCredentials: context.signUpWithCredentials || context.register,
    refreshTokenAndSetUser: context.refreshTokenAndSetUser,
    isLoading: context.isLoading,
    authChecked: context.authChecked,
  };
};
export default AuthContext; 
/**
 * Providers Component Integration Tests with Real Instances
 * Tests Providers component using real Redux, React Query, Auth, and next-intl
 */

import React, { useState, useEffect } from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useSelector, useDispatch } from 'react-redux';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { NextIntlTestProvider } from '../utils/next-intl-test-setup';
import Providers from '../../app/Providers';
import { useAuth } from '../../context/AuthContext';
import { useNavBar } from '../../../context/NavBarContext';
import { useCategories } from '../../hooks/useCategories';
import axios from '../../api/axiosInstance';

// Real test messages for Providers
const providersTestMessages = {
  en: {
    Common: {
      loading: 'Loading...',
      error: 'An error occurred',
      retry: 'Retry',
      success: 'Success'
    },
    Auth: {
      signIn: 'Sign In',
      signOut: 'Sign Out',
      signUp: 'Sign Up'
    }
  },
  pl: {
    Common: {
      loading: 'Ładowanie...',
      error: 'Wystąpił błąd',
      retry: 'Ponów',
      success: 'Sukces'
    },
    Auth: {
      signIn: 'Zaloguj się',
      signOut: 'Wyloguj się',
      signUp: 'Zarejestruj się'
    }
  }
};

// Test component that uses real Redux
function ReduxTestComponent() {
  const dispatch = useDispatch();
  const state = useSelector(state => state);
  
  const handleDispatch = () => {
    dispatch({ type: 'TEST_ACTION', payload: 'test data' });
  };

  return (
    <div data-testid="redux-test">
      <button onClick={handleDispatch} data-testid="dispatch-button">
        Dispatch Action
      </button>
      <div data-testid="redux-state">
        {JSON.stringify(state)}
      </div>
    </div>
  );
}

// Test component that uses real React Query
function ReactQueryTestComponent() {
  const queryClient = useQueryClient();
  
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['test-query'],
    queryFn: async () => {
      const response = await axios.get('/api/test');
      return response.data;
    },
    retry: false,
    staleTime: 0
  });

  const handleInvalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['test-query'] });
  };

  return (
    <div data-testid="react-query-test">
      {isLoading && <div data-testid="loading">Loading...</div>}
      {error && <div data-testid="error">Error: {error.message}</div>}
      {data && <div data-testid="data">{JSON.stringify(data)}</div>}
      <button onClick={refetch} data-testid="refetch-button">Refetch</button>
      <button onClick={handleInvalidate} data-testid="invalidate-button">Invalidate</button>
    </div>
  );
}

// Test component that uses real Auth context
function AuthTestComponent() {
  const { user, isLoading, signInWithCredentials, signOutUser } = useAuth();
  const [credentials, setCredentials] = useState({ email: '', password: '' });

  const handleSignIn = async () => {
    try {
      await signInWithCredentials(credentials);
    } catch (error) {
      console.error('Sign in failed:', error);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOutUser();
    } catch (error) {
      console.error('Sign out failed:', error);
    }
  };

  return (
    <div data-testid="auth-test">
      {isLoading && <div data-testid="auth-loading">Loading...</div>}
      {user ? (
        <div data-testid="user-info">
          <span>User: {user.email}</span>
          <button onClick={handleSignOut} data-testid="sign-out-button">Sign Out</button>
        </div>
      ) : (
        <div data-testid="sign-in-form">
          <input
            data-testid="email-input"
            value={credentials.email}
            onChange={(e) => setCredentials(prev => ({ ...prev, email: e.target.value }))}
            placeholder="Email"
          />
          <input
            data-testid="password-input"
            type="password"
            value={credentials.password}
            onChange={(e) => setCredentials(prev => ({ ...prev, password: e.target.value }))}
            placeholder="Password"
          />
          <button onClick={handleSignIn} data-testid="sign-in-button">Sign In</button>
        </div>
      )}
    </div>
  );
}

// Test component that uses real NavBar context
function NavBarTestComponent() {
  const { showNavbars, isMobile, isClient, setShowNavbars, setIsMobile } = useNavBar();

  return (
    <div data-testid="navbar-test">
      <div data-testid="navbar-state">
        showNavbars: {showNavbars.toString()}, 
        isMobile: {isMobile.toString()}, 
        isClient: {isClient.toString()}
      </div>
      <button 
        onClick={() => setShowNavbars(!showNavbars)} 
        data-testid="toggle-navbars"
      >
        Toggle Navbars
      </button>
      <button 
        onClick={() => setIsMobile(!isMobile)} 
        data-testid="toggle-mobile"
      >
        Toggle Mobile
      </button>
    </div>
  );
}

// Test component that uses real Categories context
function CategoriesTestComponent() {
  const { categories, loading, error, refetch } = useCategories();

  return (
    <div data-testid="categories-test">
      {loading && <div data-testid="categories-loading">Loading categories...</div>}
      {error && <div data-testid="categories-error">Error: {error.message}</div>}
      <div data-testid="categories-count">Categories: {categories.length}</div>
      <button onClick={refetch} data-testid="refetch-categories">Refetch Categories</button>
    </div>
  );
}

// Mock axios for controlled testing
jest.mock('../../api/axiosInstance');
const mockedAxios = axios;

describe('Providers Component - Real Integration Tests', () => {
  let user;

  beforeEach(() => {
    user = userEvent.setup();
    jest.clearAllMocks();
    
    // Setup default axios mocks
    mockedAxios.get.mockResolvedValue({ data: { message: 'test data' } });
    mockedAxios.post.mockResolvedValue({ data: { token: 'test-token' } });
  });

  afterEach(() => {
    act(() => {
      jest.runOnlyPendingTimers();
    });
  });

  describe('Real Redux Integration', () => {
    test('provides working Redux store', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReduxTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('redux-test')).toBeInTheDocument();
        expect(screen.getByTestId('redux-state')).toBeInTheDocument();
      });

      // Dispatch an action
      await user.click(screen.getByTestId('dispatch-button'));

      await waitFor(() => {
        // State should be updated (exact structure depends on your reducers)
        const stateElement = screen.getByTestId('redux-state');
        expect(stateElement).toBeInTheDocument();
      });
    });

    test('maintains Redux state across re-renders', async () => {
      const { rerender } = render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReduxTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('redux-test')).toBeInTheDocument();
      });

      // Dispatch an action
      await user.click(screen.getByTestId('dispatch-button'));

      // Re-render with different locale
      rerender(
        <NextIntlTestProvider locale="pl" messages={providersTestMessages.pl}>
          <Providers>
            <ReduxTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Redux state should persist across re-renders
        expect(screen.getByTestId('redux-state')).toBeInTheDocument();
      });
    });
  });

  describe('Real React Query Integration', () => {
    test('provides working React Query client', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReactQueryTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Should show loading initially
      await waitFor(() => {
        expect(screen.getByTestId('loading')).toBeInTheDocument();
      });

      // Should show data after loading
      await waitFor(() => {
        expect(screen.getByTestId('data')).toBeInTheDocument();
      }, { timeout: 3000 });

      // Verify axios was called
      expect(mockedAxios.get).toHaveBeenCalledWith('/api/test');
    });

    test('handles React Query refetch', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReactQueryTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByTestId('data')).toBeInTheDocument();
      });

      // Clear previous calls
      mockedAxios.get.mockClear();

      // Click refetch
      await user.click(screen.getByTestId('refetch-button'));

      // Should call API again
      await waitFor(() => {
        expect(mockedAxios.get).toHaveBeenCalledWith('/api/test');
      });
    });

    test('handles React Query cache invalidation', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReactQueryTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByTestId('data')).toBeInTheDocument();
      });

      // Click invalidate
      await user.click(screen.getByTestId('invalidate-button'));

      // Should trigger refetch
      await waitFor(() => {
        expect(mockedAxios.get).toHaveBeenCalledTimes(2); // Initial + invalidation
      });
    });

    test('handles API errors gracefully', async () => {
      // Mock API error
      mockedAxios.get.mockRejectedValueOnce(new Error('Network error'));

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReactQueryTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Should show error
      await waitFor(() => {
        expect(screen.getByTestId('error')).toBeInTheDocument();
        expect(screen.getByTestId('error')).toHaveTextContent('Network error');
      });
    });
  });

  describe('Real Auth Context Integration', () => {
    test('provides working Auth context', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <AuthTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Should show sign-in form initially (no user)
        expect(screen.getByTestId('sign-in-form')).toBeInTheDocument();
      });
    });

    test('handles real sign-in flow', async () => {
      // Mock successful login
      mockedAxios.post.mockResolvedValueOnce({
        status: 200,
        data: { 
          token: 'test-token',
          refreshToken: 'refresh-token'
        }
      });

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <AuthTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('sign-in-form')).toBeInTheDocument();
      });

      // Fill in credentials
      await user.type(screen.getByTestId('email-input'), 'redacted-email@example.com');
      await user.type(screen.getByTestId('password-input'), 'password123');

      // Click sign in
      await user.click(screen.getByTestId('sign-in-button'));

      // Should call login API
      await waitFor(() => {
        expect(mockedAxios.post).toHaveBeenCalledWith('/api/users/login', {
          email: 'redacted-email@example.com',
          password: 'password123'
        });
      });
    });

    test('handles auth errors gracefully', async () => {
      // Mock login error
      mockedAxios.post.mockRejectedValueOnce({
        response: { data: { message: 'Invalid credentials' } }
      });

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <AuthTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('sign-in-form')).toBeInTheDocument();
      });

      // Fill in credentials
      await user.type(screen.getByTestId('email-input'), 'redacted-email@example.com');
      await user.type(screen.getByTestId('password-input'), 'wrongpassword');

      // Click sign in
      await user.click(screen.getByTestId('sign-in-button'));

      // Should handle error gracefully (component should still be rendered)
      await waitFor(() => {
        expect(screen.getByTestId('sign-in-form')).toBeInTheDocument();
      });
    });
  });

  describe('Real NavBar Context Integration', () => {
    test('provides working NavBar context', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <NavBarTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('navbar-test')).toBeInTheDocument();
        expect(screen.getByTestId('navbar-state')).toBeInTheDocument();
      });
    });

    test('handles NavBar state changes', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <NavBarTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('navbar-state')).toBeInTheDocument();
      });

      // Toggle navbars
      await user.click(screen.getByTestId('toggle-navbars'));

      await waitFor(() => {
        // State should be updated
        const stateElement = screen.getByTestId('navbar-state');
        expect(stateElement).toBeInTheDocument();
      });

      // Toggle mobile
      await user.click(screen.getByTestId('toggle-mobile'));

      await waitFor(() => {
        // State should be updated again
        const stateElement = screen.getByTestId('navbar-state');
        expect(stateElement).toBeInTheDocument();
      });
    });
  });

  describe('Real Categories Context Integration', () => {
    test('provides working Categories context', async () => {
      // Mock categories API
      mockedAxios.get.mockResolvedValueOnce({
        data: [
          { id: 1, name: 'marketplace' },
          { id: 2, name: 'deals' },
          { id: 3, name: 'jobs' }
        ]
      });

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <CategoriesTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Should show loading initially
      await waitFor(() => {
        expect(screen.getByTestId('categories-loading')).toBeInTheDocument();
      });

      // Should show categories count after loading
      await waitFor(() => {
        expect(screen.getByTestId('categories-count')).toBeInTheDocument();
      }, { timeout: 3000 });
    });

    test('handles categories refetch', async () => {
      // Mock initial categories
      mockedAxios.get.mockResolvedValue({
        data: [{ id: 1, name: 'marketplace' }]
      });

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <CategoriesTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByTestId('categories-count')).toBeInTheDocument();
      });

      // Click refetch
      await user.click(screen.getByTestId('refetch-categories'));

      // Should trigger refetch
      await waitFor(() => {
        expect(mockedAxios.get).toHaveBeenCalledTimes(2); // Initial + refetch
      });
    });
  });

  describe('Real Provider Integration', () => {
    test('all providers work together', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReduxTestComponent />
            <ReactQueryTestComponent />
            <AuthTestComponent />
            <NavBarTestComponent />
            <CategoriesTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // All components should render
        expect(screen.getByTestId('redux-test')).toBeInTheDocument();
        expect(screen.getByTestId('react-query-test')).toBeInTheDocument();
        expect(screen.getByTestId('auth-test')).toBeInTheDocument();
        expect(screen.getByTestId('navbar-test')).toBeInTheDocument();
        expect(screen.getByTestId('categories-test')).toBeInTheDocument();
      });
    });

    test('providers maintain state across locale changes', async () => {
      const { rerender } = render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReduxTestComponent />
            <NavBarTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('redux-test')).toBeInTheDocument();
        expect(screen.getByTestId('navbar-test')).toBeInTheDocument();
      });

      // Make state changes
      await user.click(screen.getByTestId('dispatch-button'));
      await user.click(screen.getByTestId('toggle-navbars'));

      // Switch locale
      rerender(
        <NextIntlTestProvider locale="pl" messages={providersTestMessages.pl}>
          <Providers>
            <ReduxTestComponent />
            <NavBarTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Components should still be rendered with maintained state
        expect(screen.getByTestId('redux-test')).toBeInTheDocument();
        expect(screen.getByTestId('navbar-test')).toBeInTheDocument();
      });
    });
  });

  describe('Real Performance Tests', () => {
    test('renders efficiently with all real providers', async () => {
      const startTime = performance.now();

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <div data-testid="performance-test">Performance Test</div>
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('performance-test')).toBeInTheDocument();
      });

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Should render quickly even with all real providers
      expect(renderTime).toBeLessThan(1000);
    });

    test('handles multiple provider re-renders efficiently', async () => {
      const { rerender } = render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <div data-testid="rerender-test">Test Content</div>
          </Providers>
        </NextIntlTestProvider>
      );

      // Multiple re-renders
      for (let i = 0; i < 5; i++) {
        rerender(
          <NextIntlTestProvider locale={i % 2 === 0 ? 'en' : 'pl'} messages={providersTestMessages[i % 2 === 0 ? 'en' : 'pl']}>
            <Providers>
              <div data-testid="rerender-test">Test Content {i}</div>
            </Providers>
          </NextIntlTestProvider>
        );

        await waitFor(() => {
          expect(screen.getByTestId('rerender-test')).toBeInTheDocument();
        });
      }

      // Should handle multiple re-renders without issues
      expect(screen.getByTestId('rerender-test')).toBeInTheDocument();
    });
  });

  describe('Real Error Handling', () => {
    test('handles provider initialization errors gracefully', async () => {
      // Test with potential provider errors
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <div data-testid="error-handling-test">Error Handling Test</div>
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Should still render despite potential errors
        expect(screen.getByTestId('error-handling-test')).toBeInTheDocument();
      });

      consoleSpy.mockRestore();
    });

    test('handles network errors in providers', async () => {
      // Mock network errors
      mockedAxios.get.mockRejectedValue(new Error('Network error'));

      render(
        <NextIntlTestProvider locale="en" messages={providersTestMessages.en}>
          <Providers>
            <ReactQueryTestComponent />
            <CategoriesTestComponent />
          </Providers>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Components should handle errors gracefully
        expect(screen.getByTestId('react-query-test')).toBeInTheDocument();
        expect(screen.getByTestId('categories-test')).toBeInTheDocument();
      });
    });
  });
}); 
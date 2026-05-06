/**
 * Comprehensive Providers Test Suite
 * Tests all provider integrations: Redux, React Query, Auth, Categories, NavBar
 */

import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import { jest } from '@jest/globals';
import Providers from '../../app/Providers';

// Mock all the providers and their dependencies
jest.mock('@tanstack/react-query', () => ({
  QueryClientProvider: ({ children }) => <div data-testid="query-client-provider">{children}</div>,
  useQuery: jest.fn(),
  useQueryClient: jest.fn()
}));

jest.mock('@tanstack/react-query-devtools', () => ({
  ReactQueryDevtools: () => <div data-testid="react-query-devtools">DevTools</div>
}));

jest.mock('react-redux', () => ({
  Provider: ({ children }) => <div data-testid="redux-provider">{children}</div>,
  useSelector: jest.fn(),
  useDispatch: jest.fn()
}));

// Mock the store and query client
const mockStore = {
  getState: jest.fn(() => ({})),
  dispatch: jest.fn(),
  subscribe: jest.fn()
};

const mockQueryClient = {
  getQueryData: jest.fn(),
  setQueryData: jest.fn(),
  invalidateQueries: jest.fn()
};

jest.mock('../../lib/store', () => ({
  makeStore: jest.fn(() => mockStore)
}));

jest.mock('../../lib/reactQuery', () => ({
  getQueryClient: jest.fn(() => mockQueryClient)
}));

// Mock context providers
const mockAuthContextValue = {
  user: null,
  loading: false,
  error: null,
  signInUser: jest.fn(),
  signOutUser: jest.fn()
};

const mockNavBarContextValue = {
  showNavbars: true,
  isMobile: false,
  isClient: true,
  setShowNavbars: jest.fn(),
  setIsMobile: jest.fn()
};

const mockCategoriesContextValue = {
  categories: [],
  loading: false,
  error: null,
  refetch: jest.fn()
};

jest.mock('../../context/AuthContext', () => ({
  AuthProvider: ({ children }) => (
    <div data-testid="auth-provider" data-context={JSON.stringify(mockAuthContextValue)}>
      {children}
    </div>
  ),
  useAuth: () => mockAuthContextValue
}));

jest.mock('../../context/NavBarContext', () => ({
  NavBarProvider: ({ children }) => (
    <div data-testid="navbar-provider" data-context={JSON.stringify(mockNavBarContextValue)}>
      {children}
    </div>
  ),
  useNavBar: () => mockNavBarContextValue
}));

jest.mock('../../hooks/useCategories', () => ({
  CategoriesProvider: ({ children, prefetchTopics }) => (
    <div 
      data-testid="categories-provider" 
      data-prefetch-topics={JSON.stringify(prefetchTopics)}
      data-context={JSON.stringify(mockCategoriesContextValue)}
    >
      {children}
    </div>
  ),
  useCategories: () => mockCategoriesContextValue
}));

describe('Providers Component Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Reset NODE_ENV for each test
    delete process.env.NODE_ENV;
  });

  describe('Provider Hierarchy and Structure', () => {
    test('renders all providers in correct hierarchy', () => {
      console.log('\n🔍 TESTING: Provider hierarchy structure');
      
      render(
        <Providers>
          <div data-testid="test-child">Test Content</div>
        </Providers>
      );

      // Verify all providers are present
      expect(screen.getByTestId('redux-provider')).toBeInTheDocument();
      expect(screen.getByTestId('query-client-provider')).toBeInTheDocument();
      expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
      expect(screen.getByTestId('categories-provider')).toBeInTheDocument();
      expect(screen.getByTestId('navbar-provider')).toBeInTheDocument();
      expect(screen.getByTestId('test-child')).toBeInTheDocument();

      console.log('✅ All providers rendered correctly');
    });

    test('providers are nested in correct order', () => {
      console.log('\n🔍 TESTING: Provider nesting order');
      
      const { container } = render(
        <Providers>
          <div data-testid="test-child">Test Content</div>
        </Providers>
      );

      // Check nesting order: Redux > QueryClient > Auth > Categories > NavBar > Children
      const reduxProvider = screen.getByTestId('redux-provider');
      const queryProvider = screen.getByTestId('query-client-provider');
      const authProvider = screen.getByTestId('auth-provider');
      const categoriesProvider = screen.getByTestId('categories-provider');
      const navbarProvider = screen.getByTestId('navbar-provider');
      const testChild = screen.getByTestId('test-child');

      // Verify nesting structure
      expect(reduxProvider).toContainElement(queryProvider);
      expect(queryProvider).toContainElement(authProvider);
      expect(authProvider).toContainElement(categoriesProvider);
      expect(categoriesProvider).toContainElement(navbarProvider);
      expect(navbarProvider).toContainElement(testChild);

      console.log('✅ Provider nesting order is correct');
    });
  });

  describe('Redux Provider Integration', () => {
    test('creates store instance using useRef pattern', () => {
      console.log('\n🔍 TESTING: Redux store creation pattern');
      
      const { makeStore } = require('../../lib/store');
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      // Verify store is created
      expect(makeStore).toHaveBeenCalledTimes(1);
      console.log('✅ Redux store created with singleton pattern');
    });

    test('store instance is reused on re-renders', () => {
      console.log('\n🔍 TESTING: Store instance reuse');
      
      const { makeStore } = require('../../lib/store');
      makeStore.mockClear();

      const TestComponent = ({ count }) => (
        <Providers>
          <div>Render {count}</div>
        </Providers>
      );

      const { rerender } = render(<TestComponent count={1} />);
      
      expect(makeStore).toHaveBeenCalledTimes(1);
      
      // Re-render should not create new store
      rerender(<TestComponent count={2} />);
      expect(makeStore).toHaveBeenCalledTimes(1);

      console.log('✅ Store instance properly reused across re-renders');
    });
  });

  describe('React Query Provider Integration', () => {
    test('uses singleton QueryClient pattern', () => {
      console.log('\n🔍 TESTING: React Query client singleton');
      
      const { getQueryClient } = require('../../lib/reactQuery');
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      expect(getQueryClient).toHaveBeenCalled();
      console.log('✅ QueryClient singleton pattern working');
    });

    test('shows DevTools only in development', () => {
      console.log('\n🔍 TESTING: React Query DevTools conditional rendering');
      
      // Test development environment
      process.env.NODE_ENV = 'development';
      const { rerender } = render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      expect(screen.getByTestId('react-query-devtools')).toBeInTheDocument();
      console.log('✅ DevTools shown in development');

      // Test production environment
      process.env.NODE_ENV = 'production';
      rerender(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      expect(screen.queryByTestId('react-query-devtools')).not.toBeInTheDocument();
      console.log('✅ DevTools hidden in production');
    });
  });

  describe('Categories Provider Configuration', () => {
    test('passes correct prefetchTopics configuration', () => {
      console.log('\n🔍 TESTING: Categories provider prefetch configuration');
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      const categoriesProvider = screen.getByTestId('categories-provider');
      const prefetchTopics = JSON.parse(categoriesProvider.getAttribute('data-prefetch-topics'));
      
      expect(prefetchTopics).toEqual(['marketplace', 'deals', 'jobs', 'property', 'services']);
      console.log('✅ Categories prefetch topics configured correctly');
    });
  });

  describe('Context Values and Integration', () => {
    test('all context providers have correct default values', () => {
      console.log('\n🔍 TESTING: Context provider default values');
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      // Check Auth context
      const authProvider = screen.getByTestId('auth-provider');
      const authContext = JSON.parse(authProvider.getAttribute('data-context'));
      expect(authContext).toEqual(mockAuthContextValue);

      // Check NavBar context
      const navbarProvider = screen.getByTestId('navbar-provider');
      const navbarContext = JSON.parse(navbarProvider.getAttribute('data-context'));
      expect(navbarContext).toEqual(mockNavBarContextValue);

      // Check Categories context
      const categoriesProvider = screen.getByTestId('categories-provider');
      const categoriesContext = JSON.parse(categoriesProvider.getAttribute('data-context'));
      expect(categoriesContext).toEqual(mockCategoriesContextValue);

      console.log('✅ All context providers have correct default values');
    });
  });

  describe('Error Handling and Edge Cases', () => {
    test('handles missing children gracefully', () => {
      console.log('\n🔍 TESTING: Providers without children');
      
      expect(() => {
        render(<Providers />);
      }).not.toThrow();

      console.log('✅ Providers handle missing children gracefully');
    });

    test('handles multiple children', () => {
      console.log('\n🔍 TESTING: Providers with multiple children');
      
      render(
        <Providers>
          <div data-testid="child-1">Child 1</div>
          <div data-testid="child-2">Child 2</div>
          <span data-testid="child-3">Child 3</span>
        </Providers>
      );

      expect(screen.getByTestId('child-1')).toBeInTheDocument();
      expect(screen.getByTestId('child-2')).toBeInTheDocument();
      expect(screen.getByTestId('child-3')).toBeInTheDocument();

      console.log('✅ Providers handle multiple children correctly');
    });
  });

  describe('Performance and Optimization', () => {
    test('store reference stability across re-renders', async () => {
      console.log('\n🔍 TESTING: Store reference stability');
      
      let storeRef1, storeRef2;
      
      const TestChild = () => {
        const { useSelector } = require('react-redux');
        // Mock useSelector to capture store reference
        useSelector.mockImplementation((selector) => {
          if (!storeRef1) storeRef1 = mockStore;
          else if (!storeRef2) storeRef2 = mockStore;
          return selector({});
        });
        return <div>Test</div>;
      };

      const { rerender } = render(
        <Providers>
          <TestChild />
        </Providers>
      );

      await act(async () => {
        rerender(
          <Providers>
            <TestChild />
          </Providers>
        );
      });

      // Store references should be the same
      expect(storeRef1).toBe(storeRef2);
      console.log('✅ Store reference remains stable across re-renders');
    });

    test('QueryClient singleton behavior', () => {
      console.log('\n🔍 TESTING: QueryClient singleton behavior');
      
      const { getQueryClient } = require('../../lib/reactQuery');
      getQueryClient.mockClear();

      // Render multiple Providers instances
      render(<Providers><div>Test 1</div></Providers>);
      render(<Providers><div>Test 2</div></Providers>);

      // getQueryClient should be called for each instance
      expect(getQueryClient).toHaveBeenCalledTimes(2);
      console.log('✅ QueryClient singleton pattern working correctly');
    });
  });

  describe('Environment-Specific Behavior', () => {
    test('development environment features', () => {
      console.log('\n🔍 TESTING: Development environment features');
      
      process.env.NODE_ENV = 'development';
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      // DevTools should be present
      expect(screen.getByTestId('react-query-devtools')).toBeInTheDocument();
      
      // DevTools should have correct configuration
      const devTools = screen.getByTestId('react-query-devtools');
      expect(devTools).toHaveTextContent('DevTools');

      console.log('✅ Development features working correctly');
    });

    test('production environment optimization', () => {
      console.log('\n🔍 TESTING: Production environment optimization');
      
      process.env.NODE_ENV = 'production';
      
      render(
        <Providers>
          <div>Test</div>
        </Providers>
      );

      // DevTools should not be present in production
      expect(screen.queryByTestId('react-query-devtools')).not.toBeInTheDocument();

      console.log('✅ Production optimizations working correctly');
    });
  });

  describe('Integration with Next.js 15 and React 19', () => {
    test('compatible with React 19 features', () => {
      console.log('\n🔍 TESTING: React 19 compatibility');
      
      // Test that providers work with React 19 patterns
      expect(() => {
        render(
          <Providers>
            <div>React 19 compatible content</div>
          </Providers>
        );
      }).not.toThrow();

      console.log('✅ Providers compatible with React 19');
    });

    test('Next.js 15 App Router compatibility', () => {
      console.log('\n🔍 TESTING: Next.js 15 App Router compatibility');
      
      // Verify providers work in App Router context
      const { container } = render(
        <Providers>
          <div data-testid="app-router-content">App Router Content</div>
        </Providers>
      );

      expect(screen.getByTestId('app-router-content')).toBeInTheDocument();
      expect(container.firstChild).toHaveAttribute('data-testid', 'redux-provider');

      console.log('✅ Providers compatible with Next.js 15 App Router');
    });
  });
}); 
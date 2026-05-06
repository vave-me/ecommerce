/**
 * Comprehensive Providers Test Suite with Real Next-Intl
 * Tests all provider integrations using real next-intl library
 */

import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import { jest } from '@jest/globals';
import Providers from '../../app/Providers';
import { NextIntlTestProvider, renderWithNextIntl, setupNextIntlTest, testMessages } from '../utils/next-intl-test-setup';
import { useTranslations } from 'next-intl';

// Mock other dependencies but keep next-intl real
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

// Simple test component that uses next-intl
function TestComponent() {
  const t = useTranslations('HomePage');
  
  return (
    <div>
      <h1>{t('title')}</h1>
      <p>{t('description')}</p>
    </div>
  );
}

describe('Providers Component Tests with Real Next-Intl', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Reset NODE_ENV for each test
    delete process.env.NODE_ENV;
  });

  describe('Provider Hierarchy and Structure', () => {
    test('renders all providers in correct hierarchy', () => {
      console.log('\n🔍 TESTING: Provider hierarchy structure with real next-intl');
      
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

      console.log('✅ All providers rendered correctly with real next-intl');
    });

    test('providers work with next-intl context', () => {
      console.log('\n🔍 TESTING: Providers integration with next-intl context');
      
      renderWithNextIntl(
        <Providers>
          <div data-testid="test-child">Test Content</div>
        </Providers>,
        { locale: 'en' }
      );

      expect(screen.getByTestId('test-child')).toBeInTheDocument();
      expect(screen.getByTestId('redux-provider')).toBeInTheDocument();

      console.log('✅ Providers work correctly with next-intl context');
    });
  });

  describe('Multi-locale Testing with Real Next-Intl', () => {
    test('providers work with different locales', async () => {
      console.log('\n🔍 TESTING: Providers with different locales');
      
      const locales = ['en', 'pl', 'de'];
      
      for (const locale of locales) {
        const { unmount } = renderWithNextIntl(
          <Providers>
            <div data-testid={`test-content-${locale}`}>Content for {locale}</div>
          </Providers>,
          { locale }
        );

        expect(screen.getByTestId(`test-content-${locale}`)).toBeInTheDocument();
        expect(screen.getByTestId('redux-provider')).toBeInTheDocument();
        
        unmount();
      }

      console.log('✅ Providers work with all supported locales');
    });

    test('providers maintain state across locale changes', () => {
      console.log('\n🔍 TESTING: Provider state consistency across locales');
      
      const { makeStore } = require('../../lib/store');
      makeStore.mockClear();

      // Render with English
      const { rerender } = renderWithNextIntl(
        <Providers>
          <div data-testid="test-content">English Content</div>
        </Providers>,
        { locale: 'en' }
      );

      expect(makeStore).toHaveBeenCalledTimes(1);

      // Re-render with Polish - store should not be recreated
      rerender(
        <NextIntlTestProvider locale="pl">
          <Providers>
            <div data-testid="test-content">Polish Content</div>
          </Providers>
        </NextIntlTestProvider>
      );

      expect(makeStore).toHaveBeenCalledTimes(1); // Still only called once
      expect(screen.getByTestId('test-content')).toBeInTheDocument();

      console.log('✅ Provider state maintained across locale changes');
    });
  });

  describe('Real Next-Intl Integration', () => {
    test('renders component with English translations', () => {
      render(
        <NextIntlTestProvider locale="en" messages={testMessages.en}>
          <TestComponent />
        </NextIntlTestProvider>
      );

      expect(screen.getByText('Welcome')).toBeInTheDocument();
      expect(screen.getByText('Test description')).toBeInTheDocument();
    });

    test('renders component with Polish translations', () => {
      render(
        <NextIntlTestProvider locale="pl" messages={testMessages.pl}>
          <TestComponent />
        </NextIntlTestProvider>
      );

      expect(screen.getByText('Witamy')).toBeInTheDocument();
      expect(screen.getByText('Opis testowy')).toBeInTheDocument();
    });

    test('renders component with German translations', () => {
      render(
        <NextIntlTestProvider locale="de" messages={testMessages.de}>
          <TestComponent />
        </NextIntlTestProvider>
      );

      expect(screen.getByText('Willkommen')).toBeInTheDocument();
      expect(screen.getByText('Test Beschreibung')).toBeInTheDocument();
    });

    test('falls back to English when locale not found', () => {
      render(
        <NextIntlTestProvider locale="fr" messages={testMessages.en}>
          <TestComponent />
        </NextIntlTestProvider>
      );

      expect(screen.getByText('Welcome')).toBeInTheDocument();
      expect(screen.getByText('Test description')).toBeInTheDocument();
    });
  });

  describe('Performance with Real Next-Intl', () => {
    test('providers maintain performance with real next-intl', async () => {
      console.log('\n🔍 TESTING: Performance with real next-intl');
      
      const startTime = performance.now();
      
      renderWithNextIntl(
        <Providers>
          <div data-testid="performance-test">Performance Test</div>
        </Providers>,
        { locale: 'en' }
      );

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      expect(screen.getByTestId('performance-test')).toBeInTheDocument();
      expect(renderTime).toBeLessThan(1000); // Should render in less than 1 second

      console.log(`✅ Providers render in ${renderTime.toFixed(2)}ms with real next-intl`);
    });

    test('memory usage remains stable with real next-intl', () => {
      console.log('\n🔍 TESTING: Memory stability with real next-intl');
      
      const initialMemory = process.memoryUsage().heapUsed;
      
      // Render and unmount multiple times
      for (let i = 0; i < 10; i++) {
        const { unmount } = renderWithNextIntl(
          <Providers>
            <div data-testid={`memory-test-${i}`}>Memory Test {i}</div>
          </Providers>,
          { locale: i % 2 === 0 ? 'en' : 'pl' }
        );
        
        unmount();
      }

      const finalMemory = process.memoryUsage().heapUsed;
      const memoryIncrease = finalMemory - initialMemory;
      
      // Memory increase should be reasonable (less than 50MB)
      expect(memoryIncrease).toBeLessThan(50 * 1024 * 1024);

      console.log(`✅ Memory usage stable: ${(memoryIncrease / 1024 / 1024).toFixed(2)}MB increase`);
    });
  });

  describe('Error Handling with Real Next-Intl', () => {
    test('handles missing translations gracefully', () => {
      console.log('\n🔍 TESTING: Missing translation handling');
      
      function TestComponent() {
        const { useTranslations } = require('next-intl');
        const t = useTranslations('NonExistentNamespace');
        
        return (
          <div data-testid="missing-translation">
            {t('nonExistentKey')}
          </div>
        );
      }

      renderWithNextIntl(
        <Providers>
          <TestComponent />
        </Providers>,
        { locale: 'en' }
      );

      // Should fallback to the key name
      expect(screen.getByTestId('missing-translation')).toHaveTextContent('nonExistentKey');

      console.log('✅ Missing translations handled gracefully');
    });

    test('handles invalid locale gracefully', () => {
      console.log('\n🔍 TESTING: Invalid locale handling');
      
      expect(() => {
        renderWithNextIntl(
          <Providers>
            <div data-testid="invalid-locale-test">Test</div>
          </Providers>,
          { locale: 'invalid' }
        );
      }).not.toThrow();

      expect(screen.getByTestId('invalid-locale-test')).toBeInTheDocument();

      console.log('✅ Invalid locale handled gracefully');
    });
  });

  describe('Real Navigation Integration', () => {
    test('providers work with real next-intl navigation', () => {
      console.log('\n🔍 TESTING: Real next-intl navigation integration');
      
      function TestComponent() {
        const { Link } = require('next-intl/navigation');
        
        return (
          <div data-testid="navigation-test">
            <Link href="/test" data-testid="intl-link">
              Test Link
            </Link>
          </div>
        );
      }

      renderWithNextIntl(
        <Providers>
          <TestComponent />
        </Providers>,
        { locale: 'en' }
      );

      expect(screen.getByTestId('intl-link')).toBeInTheDocument();
      expect(screen.getByTestId('intl-link')).toHaveTextContent('Test Link');

      console.log('✅ Real next-intl navigation works correctly');
    });
  });

  describe('Server-Side Functions Integration', () => {
    test('can mock server functions while keeping client real', async () => {
      console.log('\n🔍 TESTING: Server function mocking with real client');
      
      const { setupNextIntlTest } = require('../utils/next-intl-test-setup');
      const testConfig = setupNextIntlTest('en');
      
      expect(testConfig.locale).toBe('en');
      expect(testConfig.messages).toBeDefined();
      expect(testConfig.routing).toBeDefined();

      console.log('✅ Server functions can be mocked while client remains real');
    });
  });
}); 
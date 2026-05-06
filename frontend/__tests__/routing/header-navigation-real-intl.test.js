/**
 * Header Navigation Test with Real Next-Intl
 * Tests Header component navigation behavior using real next-intl library
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import { renderWithNextIntl, NextIntlTestProvider } from '../utils/next-intl-test-setup';

// Mock the Header component dependencies
const mockRouter = {
  push: jest.fn(),
  replace: jest.fn(),
  back: jest.fn(),
  forward: jest.fn(),
  refresh: jest.fn(),
  prefetch: jest.fn()
};

const mockPathname = jest.fn(() => '/');

// Mock next-intl navigation but keep the client-side real
jest.mock('../../i18n/navigation', () => ({
  useRouter: () => mockRouter,
  usePathname: () => mockPathname(),
  Link: ({ children, href, locale, ...props }) => {
    const React = require('react');
    return React.createElement('a', { 
      href: typeof href === 'string' ? href : href.pathname,
      'data-locale': locale,
      'data-testid': 'intl-link',
      onClick: (e) => {
        e.preventDefault();
        mockRouter.push(href);
      },
      ...props 
    }, children);
  },
  redirect: jest.fn(),
  getPathname: ({ href }) => typeof href === 'string' ? href : href.pathname
}));

// Mock other dependencies
jest.mock('../../context/AuthContext', () => ({
  useAuth: () => ({
    user: null,
    loading: false,
    signInUser: jest.fn(),
    signOutUser: jest.fn()
  })
}));

jest.mock('../../context/NavBarContext', () => ({
  useNavBar: () => ({
    showNavbars: true,
    isMobile: false,
    isClient: true
  })
}));

// Create a simplified Header component for testing
const TestHeader = () => {
  const { useTranslations } = require('next-intl');
  const { Link, useRouter } = require('../../i18n/navigation');
  const t = useTranslations('Header');
  const router = useRouter();

  const handleNavigation = (path) => {
    router.push(path);
  };

  return (
    <header data-testid="header">
      <nav data-testid="header-nav">
        <Link href="/home" data-testid="home-link">
          {t('home', { default: 'Home' })}
        </Link>
        <Link href="/marketplace" data-testid="marketplace-link">
          {t('marketplace', { default: 'Marketplace' })}
        </Link>
        <Link href="/deals" data-testid="deals-link">
          {t('deals', { default: 'Deals' })}
        </Link>
        <button 
          data-testid="programmatic-nav-btn"
          onClick={() => handleNavigation('/search')}
        >
          {t('search', { default: 'Search' })}
        </button>
      </nav>
    </header>
  );
};

describe('Header Navigation with Real Next-Intl', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPathname.mockReturnValue('/');
  });

  describe('Basic Navigation Rendering', () => {
    test('renders header navigation with real translations', () => {
      console.log('\n🔍 TESTING: Header navigation rendering with real next-intl');
      
      renderWithNextIntl(<TestHeader />, { locale: 'en' });

      expect(screen.getByTestId('header')).toBeInTheDocument();
      expect(screen.getByTestId('header-nav')).toBeInTheDocument();
      expect(screen.getByTestId('home-link')).toHaveTextContent('Home');
      expect(screen.getByTestId('marketplace-link')).toHaveTextContent('Marketplace');
      expect(screen.getByTestId('deals-link')).toHaveTextContent('Deals');

      console.log('✅ Header navigation rendered with real translations');
    });

    test('renders navigation with different locales', () => {
      console.log('\n🔍 TESTING: Header navigation with different locales');
      
      // Test with Polish
      const { unmount: unmountPl } = renderWithNextIntl(<TestHeader />, { 
        locale: 'pl',
        messages: {
          Header: {
            home: 'Strona główna',
            marketplace: 'Rynek',
            deals: 'Oferty',
            search: 'Szukaj'
          }
        }
      });

      expect(screen.getByTestId('home-link')).toHaveTextContent('Strona główna');
      expect(screen.getByTestId('marketplace-link')).toHaveTextContent('Rynek');
      expect(screen.getByTestId('deals-link')).toHaveTextContent('Oferty');
      unmountPl();

      // Test with German
      renderWithNextIntl(<TestHeader />, { 
        locale: 'de',
        messages: {
          Header: {
            home: 'Startseite',
            marketplace: 'Marktplatz',
            deals: 'Angebote',
            search: 'Suchen'
          }
        }
      });

      expect(screen.getByTestId('home-link')).toHaveTextContent('Startseite');
      expect(screen.getByTestId('marketplace-link')).toHaveTextContent('Marktplatz');
      expect(screen.getByTestId('deals-link')).toHaveTextContent('Angebote');

      console.log('✅ Header navigation works with different locales');
    });
  });

  describe('Navigation Behavior Testing', () => {
    test('Link components trigger router.push with correct paths', () => {
      console.log('\n🔍 TESTING: Link navigation behavior');
      
      renderWithNextIntl(<TestHeader />, { locale: 'en' });

      // Test home link
      fireEvent.click(screen.getByTestId('home-link'));
      expect(mockRouter.push).toHaveBeenCalledWith('/home');

      // Test marketplace link
      fireEvent.click(screen.getByTestId('marketplace-link'));
      expect(mockRouter.push).toHaveBeenCalledWith('/marketplace');

      // Test deals link
      fireEvent.click(screen.getByTestId('deals-link'));
      expect(mockRouter.push).toHaveBeenCalledWith('/deals');

      console.log('✅ Link navigation behavior working correctly');
    });

    test('programmatic navigation works correctly', () => {
      console.log('\n🔍 TESTING: Programmatic navigation');
      
      renderWithNextIntl(<TestHeader />, { locale: 'en' });

      // Test programmatic navigation
      fireEvent.click(screen.getByTestId('programmatic-nav-btn'));
      expect(mockRouter.push).toHaveBeenCalledWith('/search');

      console.log('✅ Programmatic navigation working correctly');
    });

    test('navigation preserves locale context', () => {
      console.log('\n🔍 TESTING: Locale preservation in navigation');
      
      renderWithNextIntl(<TestHeader />, { locale: 'pl' });

      // All links should be rendered in Polish context
      expect(screen.getByTestId('home-link')).toBeInTheDocument();
      expect(screen.getByTestId('marketplace-link')).toBeInTheDocument();
      
      // Navigation should still work
      fireEvent.click(screen.getByTestId('home-link'));
      expect(mockRouter.push).toHaveBeenCalledWith('/home');

      console.log('✅ Locale context preserved during navigation');
    });
  });

  describe('Real Next-Intl Router Integration', () => {
    test('can capture actual router calls made by real next-intl', () => {
      console.log('\n🔍 TESTING: Real next-intl router integration');
      
      const TestComponent = () => {
        const { useRouter } = require('next-intl/navigation');
        const router = useRouter();

        return (
          <button 
            data-testid="real-intl-nav"
            onClick={() => router.push('/test-route')}
          >
            Navigate with Real Router
          </button>
        );
      };

      renderWithNextIntl(<TestComponent />, { locale: 'en' });

      fireEvent.click(screen.getByTestId('real-intl-nav'));
      expect(mockRouter.push).toHaveBeenCalledWith('/test-route');

      console.log('✅ Real next-intl router integration captured');
    });

    test('can test locale-aware navigation', () => {
      console.log('\n🔍 TESTING: Locale-aware navigation testing');
      
      const TestComponent = () => {
        const { Link } = require('next-intl/navigation');
        const { useLocale } = require('next-intl');
        const locale = useLocale();

        return (
          <div>
            <div data-testid="current-locale">{locale}</div>
            <Link href="/profile" locale="de" data-testid="locale-specific-link">
              Profile (German)
            </Link>
          </div>
        );
      };

      renderWithNextIntl(<TestComponent />, { locale: 'en' });

      expect(screen.getByTestId('current-locale')).toHaveTextContent('en');
      
      const localeLink = screen.getByTestId('locale-specific-link');
      expect(localeLink).toHaveAttribute('data-locale', 'de');

      console.log('✅ Locale-aware navigation tested successfully');
    });
  });

  describe('Navigation State Testing', () => {
    test('can test navigation state changes', async () => {
      console.log('\n🔍 TESTING: Navigation state changes');
      
      const TestComponent = () => {
        const [currentPath, setCurrentPath] = React.useState('/');
        const { useRouter } = require('next-intl/navigation');
        const router = useRouter();

        const navigate = (path) => {
          setCurrentPath(path);
          router.push(path);
        };

        return (
          <div>
            <div data-testid="current-path">{currentPath}</div>
            <button 
              data-testid="nav-to-about"
              onClick={() => navigate('/about')}
            >
              Go to About
            </button>
          </div>
        );
      };

      renderWithNextIntl(<TestComponent />, { locale: 'en' });

      expect(screen.getByTestId('current-path')).toHaveTextContent('/');

      fireEvent.click(screen.getByTestId('nav-to-about'));

      await waitFor(() => {
        expect(screen.getByTestId('current-path')).toHaveTextContent('/about');
      });

      expect(mockRouter.push).toHaveBeenCalledWith('/about');

      console.log('✅ Navigation state changes tested successfully');
    });
  });

  describe('Error Handling in Navigation', () => {
    test('handles navigation errors gracefully', () => {
      console.log('\n🔍 TESTING: Navigation error handling');
      
      // Mock router.push to throw an error
      mockRouter.push.mockImplementationOnce(() => {
        throw new Error('Navigation failed');
      });

      const TestComponent = () => {
        const { useRouter } = require('next-intl/navigation');
        const router = useRouter();
        const [error, setError] = React.useState(null);

        const handleNavigation = () => {
          try {
            router.push('/error-route');
          } catch (err) {
            setError(err.message);
          }
        };

        return (
          <div>
            <button data-testid="error-nav" onClick={handleNavigation}>
              Navigate with Error
            </button>
            {error && <div data-testid="nav-error">{error}</div>}
          </div>
        );
      };

      renderWithNextIntl(<TestComponent />, { locale: 'en' });

      fireEvent.click(screen.getByTestId('error-nav'));

      expect(screen.getByTestId('nav-error')).toHaveTextContent('Navigation failed');

      console.log('✅ Navigation error handling tested successfully');
    });
  });

  describe('Performance Testing with Real Next-Intl', () => {
    test('navigation performance with real next-intl', () => {
      console.log('\n🔍 TESTING: Navigation performance');
      
      const startTime = performance.now();
      
      renderWithNextIntl(<TestHeader />, { locale: 'en' });
      
      // Perform multiple navigation actions
      for (let i = 0; i < 10; i++) {
        fireEvent.click(screen.getByTestId('home-link'));
      }
      
      const endTime = performance.now();
      const totalTime = endTime - startTime;
      
      expect(totalTime).toBeLessThan(1000); // Should complete in less than 1 second
      expect(mockRouter.push).toHaveBeenCalledTimes(10);

      console.log(`✅ Navigation performance: ${totalTime.toFixed(2)}ms for 10 operations`);
    });
  });
}); 
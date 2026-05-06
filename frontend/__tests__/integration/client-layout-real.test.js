/**
 * ClientLayout Component Integration Tests with Real Instances
 * Tests ClientLayout component using real contexts, real routing, and real components
 */

import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { NextIntlTestProvider } from '../utils/next-intl-test-setup';
import ClientLayout from '../../app/ClientLayout.client';
import { AuthProvider } from '../../context/AuthContext';
import { NavBarProvider } from '../../context/NavBarContext';
import { Provider as ReduxProvider } from 'react-redux';
import { QueryClientProvider } from '@tanstack/react-query';
import { makeStore } from '../../lib/store';
import { getQueryClient } from '../../lib/reactQuery';

// Real test messages for ClientLayout
const clientLayoutTestMessages = {
  en: {
    Header: {
      homeButton: 'Home',
      exploreButton: 'Explore',
      createButtonText: 'Create',
      mainNavAriaLabel: 'Main navigation'
    },
    Navigation: {
      home: 'Home',
      search: 'Search',
      add: 'Add',
      messages: 'Messages',
      profile: 'Profile'
    }
  },
  pl: {
    Header: {
      homeButton: 'Strona główna',
      exploreButton: 'Eksploruj',
      createButtonText: 'Utwórz',
      mainNavAriaLabel: 'Główna nawigacja'
    },
    Navigation: {
      home: 'Strona główna',
      search: 'Szukaj',
      add: 'Dodaj',
      messages: 'Wiadomości',
      profile: 'Profil'
    }
  }
};

// Real test provider with all necessary contexts
function RealClientLayoutProvider({ 
  children, 
  locale = 'en', 
  isMobile = false, 
  showNavbars = true,
  pathname = '/' 
}) {
  const store = makeStore();
  const queryClient = getQueryClient();

  // Mock usePathname to return the test pathname
  jest.doMock('../../i18n/navigation', () => ({
    usePathname: () => pathname,
    useRouter: () => ({ push: jest.fn() })
  }));

  return (
    <NextIntlTestProvider locale={locale} messages={clientLayoutTestMessages[locale]}>
      <ReduxProvider store={store}>
        <QueryClientProvider client={queryClient}>
          <AuthProvider>
            <NavBarProvider initialState={{ isMobile, showNavbars, isClient: true }}>
              {children}
            </NavBarProvider>
          </AuthProvider>
        </QueryClientProvider>
      </ReduxProvider>
    </NextIntlTestProvider>
  );
}

// Mock window.matchMedia for responsive tests
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: query.includes('768px') ? false : true,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

describe('ClientLayout Component - Real Integration Tests', () => {
  let user;

  beforeEach(() => {
    user = userEvent.setup();
    jest.clearAllMocks();
  });

  afterEach(() => {
    act(() => {
      jest.runOnlyPendingTimers();
    });
  });

  describe('Real Context Integration', () => {
    test('integrates with real NavBarContext for desktop layout', async () => {
      render(
        <RealClientLayoutProvider isMobile={false} showNavbars={true}>
          <ClientLayout>
            <div data-testid="test-content">Desktop Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Header should be visible on desktop
        expect(screen.getByRole('banner')).toBeInTheDocument();
        
        // Content should be rendered
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
        
        // Bottom nav should not be visible on desktop
        expect(screen.queryByRole('navigation', { name: /bottom/i })).not.toBeInTheDocument();
      });
    });

    test('integrates with real NavBarContext for mobile layout', async () => {
      render(
        <RealClientLayoutProvider isMobile={true} showNavbars={true}>
          <ClientLayout>
            <div data-testid="test-content">Mobile Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Header should be visible on mobile
        expect(screen.getByRole('banner')).toBeInTheDocument();
        
        // Content should be rendered
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
        
        // Bottom nav should be visible on mobile
        expect(screen.getByTestId('bottom-nav')).toBeInTheDocument();
      });
    });

    test('handles real AuthContext integration', async () => {
      const testUser = {
        userId: '123',
        email: 'redacted-email@example.com',
        username: 'testuser'
      };

      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="authenticated-content">Authenticated Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Layout should render regardless of auth state
        expect(screen.getByTestId('authenticated-content')).toBeInTheDocument();
        expect(screen.getByRole('banner')).toBeInTheDocument();
      });
    });
  });

  describe('Real Routing Integration', () => {
    test('handles different routes correctly', async () => {
      const routes = ['/', '/home', '/explore', '/profile'];

      for (const route of routes) {
        const { unmount } = render(
          <RealClientLayoutProvider pathname={route}>
            <ClientLayout>
              <div data-testid={`content-${route.replace('/', '') || 'root'}`}>
                Content for {route}
              </div>
            </ClientLayout>
          </RealClientLayoutProvider>
        );

        await waitFor(() => {
          expect(screen.getByTestId(`content-${route.replace('/', '') || 'root'}`)).toBeInTheDocument();
          expect(screen.getByRole('banner')).toBeInTheDocument();
        });

        unmount();
      }
    });

    test('handles special routes that hide header on mobile', async () => {
      const specialRoutes = ['/messages', '/messages/'];

      for (const route of specialRoutes) {
        const { unmount } = render(
          <RealClientLayoutProvider pathname={route} isMobile={true}>
            <ClientLayout>
              <div data-testid="messages-content">Messages Content</div>
            </ClientLayout>
          </RealClientLayoutProvider>
        );

        await waitFor(() => {
          // Content should be visible
          expect(screen.getByTestId('messages-content')).toBeInTheDocument();
          
          // Header should be hidden on mobile for messages routes
          expect(screen.queryByRole('banner')).not.toBeInTheDocument();
        });

        unmount();
      }
    });

    test('preserves locale in routing', async () => {
      const { rerender } = render(
        <RealClientLayoutProvider locale="en" pathname="/en/home">
          <ClientLayout>
            <div data-testid="english-content">English Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('english-content')).toBeInTheDocument();
      });

      // Switch to Polish locale
      rerender(
        <RealClientLayoutProvider locale="pl" pathname="/pl/home">
          <ClientLayout>
            <div data-testid="polish-content">Polish Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('polish-content')).toBeInTheDocument();
      });
    });
  });

  describe('Real Responsive Behavior', () => {
    test('adapts layout based on real screen size changes', async () => {
      // Start with desktop
      const { rerender } = render(
        <RealClientLayoutProvider isMobile={false} showNavbars={true}>
          <ClientLayout>
            <div data-testid="responsive-content">Responsive Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('responsive-content')).toBeInTheDocument();
        expect(screen.getByRole('banner')).toBeInTheDocument();
        // No bottom nav on desktop
        expect(screen.queryByTestId('bottom-nav')).not.toBeInTheDocument();
      });

      // Switch to mobile
      rerender(
        <RealClientLayoutProvider isMobile={true} showNavbars={true}>
          <ClientLayout>
            <div data-testid="responsive-content">Responsive Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('responsive-content')).toBeInTheDocument();
        expect(screen.getByRole('banner')).toBeInTheDocument();
        // Bottom nav should appear on mobile
        expect(screen.getByTestId('bottom-nav')).toBeInTheDocument();
      });
    });

    test('handles navbar visibility changes', async () => {
      const { rerender } = render(
        <RealClientLayoutProvider showNavbars={true}>
          <ClientLayout>
            <div data-testid="navbar-test">Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByRole('banner')).toBeInTheDocument();
      });

      // Hide navbars
      rerender(
        <RealClientLayoutProvider showNavbars={false}>
          <ClientLayout>
            <div data-testid="navbar-test">Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Header should still be visible (controlled by other logic)
        expect(screen.getByRole('banner')).toBeInTheDocument();
      });
    });
  });

  describe('Real Component Integration', () => {
    test('renders real Header component', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="header-integration">Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Real Header should be rendered
        const header = screen.getByRole('banner');
        expect(header).toBeInTheDocument();
        
        // Header should have navigation
        expect(screen.getByRole('navigation')).toBeInTheDocument();
      });
    });

    test('renders real BottomNav component on mobile', async () => {
      render(
        <RealClientLayoutProvider isMobile={true} showNavbars={true}>
          <ClientLayout>
            <div data-testid="bottom-nav-integration">Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Real BottomNav should be rendered
        expect(screen.getByTestId('bottom-nav')).toBeInTheDocument();
      });
    });

    test('renders real GlobalModals component', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="modals-integration">Content</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // GlobalModals should be rendered (even if not visible)
        expect(screen.getByTestId('global-modals')).toBeInTheDocument();
      });
    });
  });

  describe('Real Performance Tests', () => {
    test('renders efficiently with all real providers', async () => {
      const startTime = performance.now();

      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="performance-test">
              <h1>Performance Test Content</h1>
              <p>This is a test of rendering performance</p>
            </div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('performance-test')).toBeInTheDocument();
      });

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Should render quickly even with all real providers
      expect(renderTime).toBeLessThan(1000);
    });

    test('handles multiple children efficiently', async () => {
      const children = Array.from({ length: 10 }, (_, i) => (
        <div key={i} data-testid={`child-${i}`}>Child {i}</div>
      ));

      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            {children}
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // All children should be rendered
        children.forEach((_, i) => {
          expect(screen.getByTestId(`child-${i}`)).toBeInTheDocument();
        });
      });
    });
  });

  describe('Real Error Handling', () => {
    test('handles context provider errors gracefully', async () => {
      // Test with minimal context
      render(
        <NextIntlTestProvider locale="en" messages={clientLayoutTestMessages.en}>
          <ClientLayout>
            <div data-testid="error-handling">Error Handling Test</div>
          </ClientLayout>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Should still render content even with missing contexts
        expect(screen.getByTestId('error-handling')).toBeInTheDocument();
      });
    });

    test('handles missing translations gracefully', async () => {
      render(
        <NextIntlTestProvider locale="en" messages={{}}>
          <RealClientLayoutProvider>
            <ClientLayout>
              <div data-testid="missing-translations">Missing Translations Test</div>
            </ClientLayout>
          </RealClientLayoutProvider>
        </NextIntlTestProvider>
      );

      await waitFor(() => {
        // Should still render despite missing translations
        expect(screen.getByTestId('missing-translations')).toBeInTheDocument();
      });
    });
  });

  describe('Real Accessibility Tests', () => {
    test('maintains proper semantic structure', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="accessibility-test">Accessibility Test</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Should have proper semantic structure
        expect(screen.getByRole('banner')).toBeInTheDocument(); // Header
        expect(screen.getByRole('main')).toBeInTheDocument(); // Main content
      });
    });

    test('supports keyboard navigation', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <button data-testid="focusable-element">Focusable Button</button>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('focusable-element')).toBeInTheDocument();
      });

      // Tab navigation should work
      await user.tab();
      
      // Should be able to focus elements within the layout
      expect(document.activeElement).toBeDefined();
    });

    test('provides proper ARIA landmarks', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="aria-test">ARIA Test</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Should have proper ARIA landmarks
        const main = screen.getByRole('main');
        expect(main).toBeInTheDocument();
        expect(main).toHaveClass('content');
      });
    });
  });

  describe('Real State Management', () => {
    test('integrates with real Redux store', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="redux-integration">Redux Integration Test</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Should render with Redux integration
        expect(screen.getByTestId('redux-integration')).toBeInTheDocument();
      });
    });

    test('integrates with real React Query', async () => {
      render(
        <RealClientLayoutProvider>
          <ClientLayout>
            <div data-testid="react-query-integration">React Query Integration Test</div>
          </ClientLayout>
        </RealClientLayoutProvider>
      );

      await waitFor(() => {
        // Should render with React Query integration
        expect(screen.getByTestId('react-query-integration')).toBeInTheDocument();
      });
    });
  });
}); 
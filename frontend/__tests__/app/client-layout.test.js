/**
 * Comprehensive ClientLayout Test Suite
 * Tests client-side layout functionality, responsive behavior, and context integration
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import ClientLayout from '../../app/ClientLayout.client';

// Mock the navigation hook
const mockUsePathname = jest.fn();
jest.mock('../../i18n/navigation', () => ({
  usePathname: () => mockUsePathname()
}));

// Mock all child components
jest.mock('../../components/Header/Header', () => {
  return function MockHeader() {
    return <div data-testid="header">Header Component</div>;
  };
});

jest.mock('../../components/Header/BottomNav', () => {
  return function MockBottomNav() {
    return <div data-testid="bottom-nav">Bottom Navigation</div>;
  };
});

jest.mock('../../components/Shared/GlobalModals', () => {
  return function MockGlobalModals() {
    return <div data-testid="global-modals">Global Modals</div>;
  };
});

jest.mock('../../components/Wishlist/WishlistSelectorModal', () => {
  return function MockWishlistSelectorModal() {
    return <div data-testid="wishlist-modal">Wishlist Modal</div>;
  };
});

// Mock CSS modules
jest.mock('../../app/ClientLayout.module.css', () => ({
  layout: 'layout',
  content: 'content'
}));

// Mock NavBarContext
const mockNavBarContext = {
  isMobile: false,
  showNavbars: true,
  isClient: true,
  setIsMobile: jest.fn(),
  setShowNavbars: jest.fn()
};

jest.mock('../../context/NavBarContext', () => ({
  NavBarContext: React.createContext(mockNavBarContext)
}));

describe('ClientLayout Component Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockUsePathname.mockReturnValue('/home');
    // Reset console.warn mock
    jest.spyOn(console, 'warn').mockImplementation(() => {});
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('Basic Rendering and Structure', () => {
    test('renders basic layout structure', () => {
      console.log('\n🔍 TESTING: Basic layout structure');
      
      render(
        <ClientLayout>
          <div data-testid="test-content">Test Content</div>
        </ClientLayout>
      );

      expect(screen.getByTestId('test-content')).toBeInTheDocument();
      expect(screen.getByTestId('global-modals')).toBeInTheDocument();
      expect(screen.getByTestId('wishlist-modal')).toBeInTheDocument();

      console.log('✅ Basic layout structure rendered correctly');
    });

    test('applies correct CSS classes', () => {
      console.log('\n🔍 TESTING: CSS class application');
      
      const { container } = render(
        <ClientLayout>
          <div>Test</div>
        </ClientLayout>
      );

      const layoutDiv = container.firstChild;
      const mainElement = container.querySelector('main');

      expect(layoutDiv).toHaveClass('layout');
      expect(mainElement).toHaveClass('content');

      console.log('✅ CSS classes applied correctly');
    });
  });

  describe('Context Integration and Defensive Programming', () => {
    test('handles valid NavBarContext correctly', () => {
      console.log('\n🔍 TESTING: Valid NavBarContext handling');
      
      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={mockNavBarContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div data-testid="test-content">Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.getByTestId('test-content')).toBeInTheDocument();
      expect(screen.getByTestId('header')).toBeInTheDocument();

      console.log('✅ Valid NavBarContext handled correctly');
    });

    test('handles null NavBarContext with defensive programming', () => {
      console.log('\n🔍 TESTING: Null NavBarContext defensive handling');
      
      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={null}>
            {children}
          </NavBarContext.Provider>
        );
      };

      // Mock development environment
      const originalEnv = process.env.NODE_ENV;
      process.env.NODE_ENV = 'development';

      render(
        <TestWrapper>
          <ClientLayout>
            <div data-testid="test-content">Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      // Should still render with defaults
      expect(screen.getByTestId('test-content')).toBeInTheDocument();
      
      // Should log warning in development
      expect(console.warn).toHaveBeenCalledWith(
        'NavBarContext is null. Make sure NavBarProvider is included in the provider chain.'
      );

      process.env.NODE_ENV = originalEnv;
      console.log('✅ Null NavBarContext handled defensively');
    });

    test('uses default values when context is null', () => {
      console.log('\n🔍 TESTING: Default values when context is null');
      
      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={null}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div data-testid="test-content">Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      // With null context, defaults should be used:
      // isMobile: false, showNavbars: false, isClient: false
      // This means no bottom nav should be shown
      expect(screen.queryByTestId('bottom-nav')).not.toBeInTheDocument();

      console.log('✅ Default values used correctly when context is null');
    });
  });

  describe('Responsive Behavior and Mobile Detection', () => {
    test('shows bottom navigation on mobile when conditions are met', () => {
      console.log('\n🔍 TESTING: Bottom navigation on mobile');
      
      const mobileContext = {
        ...mockNavBarContext,
        isMobile: true,
        showNavbars: true,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={mobileContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.getByTestId('bottom-nav')).toBeInTheDocument();
      console.log('✅ Bottom navigation shown on mobile');
    });

    test('hides bottom navigation on desktop', () => {
      console.log('\n🔍 TESTING: Bottom navigation hidden on desktop');
      
      const desktopContext = {
        ...mockNavBarContext,
        isMobile: false,
        showNavbars: true,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={desktopContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.queryByTestId('bottom-nav')).not.toBeInTheDocument();
      console.log('✅ Bottom navigation hidden on desktop');
    });

    test('hides bottom navigation when showNavbars is false', () => {
      console.log('\n🔍 TESTING: Bottom navigation hidden when showNavbars is false');
      
      const hiddenNavContext = {
        ...mockNavBarContext,
        isMobile: true,
        showNavbars: false,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={hiddenNavContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.queryByTestId('bottom-nav')).not.toBeInTheDocument();
      console.log('✅ Bottom navigation hidden when showNavbars is false');
    });

    test('hides bottom navigation before client hydration', () => {
      console.log('\n🔍 TESTING: Bottom navigation hidden before hydration');
      
      const preHydrationContext = {
        ...mockNavBarContext,
        isMobile: true,
        showNavbars: true,
        isClient: false // Not hydrated yet
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={preHydrationContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.queryByTestId('bottom-nav')).not.toBeInTheDocument();
      console.log('✅ Bottom navigation hidden before hydration');
    });
  });

  describe('Header Visibility Logic', () => {
    test('shows header on normal routes', () => {
      console.log('\n🔍 TESTING: Header visibility on normal routes');
      
      mockUsePathname.mockReturnValue('/home');

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={mockNavBarContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.getByTestId('header')).toBeInTheDocument();
      console.log('✅ Header shown on normal routes');
    });

    test('hides header on mobile for specific routes', () => {
      console.log('\n🔍 TESTING: Header hidden on mobile for specific routes');
      
      mockUsePathname.mockReturnValue('/messages');

      const mobileContext = {
        ...mockNavBarContext,
        isMobile: true,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={mobileContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.queryByTestId('header')).not.toBeInTheDocument();
      console.log('✅ Header hidden on mobile for /messages route');
    });

    test('hides header on mobile for messages with trailing slash', () => {
      console.log('\n🔍 TESTING: Header hidden for /messages/ route');
      
      mockUsePathname.mockReturnValue('/messages/');

      const mobileContext = {
        ...mockNavBarContext,
        isMobile: true,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={mobileContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.queryByTestId('header')).not.toBeInTheDocument();
      console.log('✅ Header hidden for /messages/ route');
    });

    test('shows header on desktop even for skip routes', () => {
      console.log('\n🔍 TESTING: Header shown on desktop for skip routes');
      
      mockUsePathname.mockReturnValue('/messages');

      const desktopContext = {
        ...mockNavBarContext,
        isMobile: false,
        isClient: true
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={desktopContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.getByTestId('header')).toBeInTheDocument();
      console.log('✅ Header shown on desktop even for skip routes');
    });

    test('shows header before client hydration', () => {
      console.log('\n🔍 TESTING: Header shown before hydration');
      
      mockUsePathname.mockReturnValue('/messages');

      const preHydrationContext = {
        ...mockNavBarContext,
        isMobile: true,
        isClient: false // Not hydrated yet
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={preHydrationContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(screen.getByTestId('header')).toBeInTheDocument();
      console.log('✅ Header shown before hydration');
    });
  });

  describe('Children Rendering', () => {
    test('renders single child correctly', () => {
      console.log('\n🔍 TESTING: Single child rendering');
      
      render(
        <ClientLayout>
          <div data-testid="single-child">Single Child</div>
        </ClientLayout>
      );

      expect(screen.getByTestId('single-child')).toBeInTheDocument();
      console.log('✅ Single child rendered correctly');
    });

    test('renders multiple children correctly', () => {
      console.log('\n🔍 TESTING: Multiple children rendering');
      
      render(
        <ClientLayout>
          <div data-testid="child-1">Child 1</div>
          <div data-testid="child-2">Child 2</div>
          <span data-testid="child-3">Child 3</span>
        </ClientLayout>
      );

      expect(screen.getByTestId('child-1')).toBeInTheDocument();
      expect(screen.getByTestId('child-2')).toBeInTheDocument();
      expect(screen.getByTestId('child-3')).toBeInTheDocument();

      console.log('✅ Multiple children rendered correctly');
    });

    test('handles no children gracefully', () => {
      console.log('\n🔍 TESTING: No children handling');
      
      expect(() => {
        render(<ClientLayout />);
      }).not.toThrow();

      console.log('✅ No children handled gracefully');
    });
  });

  describe('Development Environment Behavior', () => {
    test('logs warning in development when context is null', () => {
      console.log('\n🔍 TESTING: Development warning logging');
      
      const originalEnv = process.env.NODE_ENV;
      process.env.NODE_ENV = 'development';

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={null}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(console.warn).toHaveBeenCalledWith(
        'NavBarContext is null. Make sure NavBarProvider is included in the provider chain.'
      );

      process.env.NODE_ENV = originalEnv;
      console.log('✅ Development warning logged correctly');
    });

    test('does not log warning in production when context is null', () => {
      console.log('\n🔍 TESTING: Production warning suppression');
      
      const originalEnv = process.env.NODE_ENV;
      process.env.NODE_ENV = 'production';

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={null}>
            {children}
          </NavBarContext.Provider>
        );
      };

      render(
        <TestWrapper>
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        </TestWrapper>
      );

      expect(console.warn).not.toHaveBeenCalled();

      process.env.NODE_ENV = originalEnv;
      console.log('✅ Production warning suppressed correctly');
    });
  });

  describe('Component Integration', () => {
    test('all child components are rendered', () => {
      console.log('\n🔍 TESTING: All child components integration');
      
      render(
        <ClientLayout>
          <div>Test</div>
        </ClientLayout>
      );

      expect(screen.getByTestId('header')).toBeInTheDocument();
      expect(screen.getByTestId('global-modals')).toBeInTheDocument();
      expect(screen.getByTestId('wishlist-modal')).toBeInTheDocument();

      console.log('✅ All child components integrated correctly');
    });

    test('components render in correct order', () => {
      console.log('\n🔍 TESTING: Component rendering order');
      
      const { container } = render(
        <ClientLayout>
          <div data-testid="main-content">Main Content</div>
        </ClientLayout>
      );

      const elements = Array.from(container.querySelectorAll('[data-testid]'));
      const testIds = elements.map(el => el.getAttribute('data-testid'));

      // Expected order: header, main-content, global-modals, wishlist-modal
      expect(testIds).toEqual(['header', 'main-content', 'global-modals', 'wishlist-modal']);

      console.log('✅ Components render in correct order');
    });
  });

  describe('Edge Cases and Error Handling', () => {
    test('handles undefined pathname gracefully', () => {
      console.log('\n🔍 TESTING: Undefined pathname handling');
      
      mockUsePathname.mockReturnValue(undefined);

      expect(() => {
        render(
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        );
      }).not.toThrow();

      console.log('✅ Undefined pathname handled gracefully');
    });

    test('handles empty pathname gracefully', () => {
      console.log('\n🔍 TESTING: Empty pathname handling');
      
      mockUsePathname.mockReturnValue('');

      expect(() => {
        render(
          <ClientLayout>
            <div>Test</div>
          </ClientLayout>
        );
      }).not.toThrow();

      console.log('✅ Empty pathname handled gracefully');
    });

    test('handles context with partial values', () => {
      console.log('\n🔍 TESTING: Partial context values handling');
      
      const partialContext = {
        isMobile: true,
        // Missing showNavbars, isClient
      };

      const TestWrapper = ({ children }) => {
        const { NavBarContext } = require('../../context/NavBarContext');
        return (
          <NavBarContext.Provider value={partialContext}>
            {children}
          </NavBarContext.Provider>
        );
      };

      expect(() => {
        render(
          <TestWrapper>
            <ClientLayout>
              <div>Test</div>
            </ClientLayout>
          </TestWrapper>
        );
      }).not.toThrow();

      console.log('✅ Partial context values handled gracefully');
    });
  });
}); 
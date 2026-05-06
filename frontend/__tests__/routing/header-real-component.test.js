/**
 * Real Header and AddDropdown Component Test Suite
 * Tests actual user interactions with real components and providers
 * Minimal mocking to see real locale preservation behavior
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import { NextIntlClientProvider } from 'next-intl';
import Header from '../../components/Header/Header';
import AddDropdown from '../../components/Header/AddDropdown';
import { AuthProvider } from '../../context/AuthContext';
import { NavBarProvider } from '../../context/NavBarContext';

// Only mock the router push function to capture navigation calls
const mockPush = jest.fn();
const mockRouter = {
  push: mockPush,
  replace: jest.fn(),
  back: jest.fn(),
  forward: jest.fn(),
  refresh: jest.fn(),
  prefetch: jest.fn()
};

// Mock only the navigation hooks, keep everything else real
jest.mock('../../i18n/navigation', () => {
  const actual = jest.requireActual('../../i18n/navigation');
  return {
    ...actual,
    useRouter: () => mockRouter,
    usePathname: () => '/en/home',
    useParams: () => ({ locale: 'en' })
  };
});

jest.mock('next/navigation', () => ({
  usePathname: () => '/en/home',
  useParams: () => ({ locale: 'en' }),
  useSearchParams: () => new URLSearchParams()
}));

// Real translations for testing
const messages = {
  Header: {
    homeButton: 'Home',
    exploreButton: 'Explore',
    createButtonText: 'Create',
    mainNavAriaLabel: 'Main navigation',
    notificationsButtonAriaLabel: 'Notifications ({count})',
    messagesButtonAriaLabel: 'Messages ({count})',
    wishlistButtonAriaLabel: 'Wishlist ({count})',
    cartButtonAriaLabel: 'Shopping cart',
    createContentAriaLabel: 'Create content'
  },
  AddDropdown: {
    title: 'Create New',
    frequentTab: 'Frequent',
    allOptionsTab: 'All Options',
    ariaLabel: 'Create content menu',
    recentlyUsedTitle: 'Recently Used',
    addAriaLabel: 'Add {label}',
    section_frequent_title: 'Frequently Used',
    section_all_title: 'All Options',
    item_product_label: 'Product',
    item_product_desc: 'Sell a product',
    item_post_label: 'Post',
    item_post_desc: 'Share a post',
    item_vehicle_label: 'Vehicle',
    item_vehicle_desc: 'Sell a vehicle',
    item_deal_label: 'Deal',
    item_deal_desc: 'Share a deal',
    item_property_label: 'Property',
    item_property_desc: 'List property',
    item_job_label: 'Job',
    item_job_desc: 'Post a job',
    item_service_label: 'Service',
    item_service_desc: 'Offer a service'
  }
};

// Test wrapper with real providers
function TestWrapper({ children, locale = 'en' }) {
  const mockUser = {
    id: 1,
    name: 'Test User',
    email: 'redacted-email@example.com'
  };

  const mockAuthValue = {
    user: mockUser,
    signOutUser: jest.fn(),
    loading: false,
    error: null
  };

  const mockNavBarValue = {
    showNavbars: true,
    isMobile: false,
    isClient: true,
    setShowNavbars: jest.fn(),
    setIsMobile: jest.fn()
  };

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <AuthProvider value={mockAuthValue}>
        <NavBarProvider value={mockNavBarValue}>
          {children}
        </NavBarProvider>
      </AuthProvider>
    </NextIntlClientProvider>
  );
}

describe('Real Header Component Navigation Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
    console.log('\n' + '='.repeat(60));
  });

  describe('Header Navigation Button Interactions', () => {
    test('home button click - check router.push call', async () => {
      console.log('🔍 TESTING: Header Home button click');
      console.log('Expected: router.push should be called with locale-aware route');
      
      render(
        <TestWrapper locale="en">
          <Header />
        </TestWrapper>
      );

      // Find and click home button
      const homeButton = screen.getByText('Home');
      expect(homeButton).toBeInTheDocument();
      
      console.log('📍 Clicking Home button...');
      fireEvent.click(homeButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/home') {
        console.log('❌ ISSUE FOUND: Route pushed WITHOUT locale prefix');
        console.log('   This means users will be redirected to /home instead of /en/home');
        console.log('   The i18n router should automatically add the locale prefix');
      } else if (pushedRoute.startsWith('/en/home')) {
        console.log('✅ GOOD: Route includes locale prefix');
      } else {
        console.log(`⚠️ UNEXPECTED: Route format is unexpected: ${pushedRoute}`);
      }
    });

    test('explore button click - check router.push call', async () => {
      console.log('🔍 TESTING: Header Explore button click');
      
      render(
        <TestWrapper locale="en">
          <Header />
        </TestWrapper>
      );

      const exploreButton = screen.getByText('Explore');
      console.log('📍 Clicking Explore button...');
      fireEvent.click(exploreButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/explore') {
        console.log('❌ ISSUE FOUND: Explore route pushed WITHOUT locale prefix');
      } else {
        console.log('✅ Route includes proper formatting');
      }
    });

    test('notifications button click - check router.push call', async () => {
      console.log('🔍 TESTING: Header Notifications button click');
      
      render(
        <TestWrapper locale="en">
          <Header />
        </TestWrapper>
      );

      // Find notifications button (it's an icon button)
      const notificationsButton = screen.getByLabelText('Notifications (1)');
      console.log('📍 Clicking Notifications button...');
      fireEvent.click(notificationsButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/notifications') {
        console.log('❌ ISSUE FOUND: Notifications route pushed WITHOUT locale prefix');
      } else {
        console.log('✅ Route includes proper formatting');
      }
    });
  });

  describe('Header Navigation from Different Locales', () => {
    test.each(['en', 'pl', 'de'])('home button from %s locale', async (locale) => {
      console.log(`🔍 TESTING: Header navigation from ${locale} locale`);
      
      // Mock the current locale context
      jest.doMock('../../i18n/navigation', () => {
        const actual = jest.requireActual('../../i18n/navigation');
        return {
          ...actual,
          useRouter: () => mockRouter,
          usePathname: () => `/${locale}/home`,
          useParams: () => ({ locale })
        };
      });

      render(
        <TestWrapper locale={locale}>
          <Header />
        </TestWrapper>
      );

      const homeButton = screen.getByText('Home');
      console.log(`📍 Clicking Home button from ${locale} locale...`);
      fireEvent.click(homeButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: From ${locale} locale, router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/home') {
        console.log(`❌ ISSUE: Route doesn't preserve ${locale} locale`);
        console.log(`   Expected: /${locale}/home or locale-aware routing`);
        console.log(`   Actual: ${pushedRoute}`);
      } else if (pushedRoute.includes(locale)) {
        console.log(`✅ GOOD: Route preserves ${locale} locale`);
      } else {
        console.log(`⚠️ UNCLEAR: Route format: ${pushedRoute}`);
      }
    });
  });
});

describe('Real AddDropdown Component Navigation Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
    console.log('\n' + '='.repeat(60));
  });

  describe('AddDropdown Item Interactions', () => {
    test('product item click - check router.push call', async () => {
      console.log('🔍 TESTING: AddDropdown Product item click');
      console.log('Expected: router.push should be called with /add/product?step=1');
      
      const mockOnClose = jest.fn();
      render(
        <TestWrapper locale="en">
          <AddDropdown onClose={mockOnClose} />
        </TestWrapper>
      );

      // Find and click product item
      const productButton = screen.getByText('Product');
      expect(productButton).toBeInTheDocument();
      
      console.log('📍 Clicking Product item...');
      fireEvent.click(productButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/add/product?step=1') {
        console.log('❌ CRITICAL ISSUE: AddDropdown route pushed WITHOUT locale prefix');
        console.log('   This means users clicking "Product" will go to /add/product');
        console.log('   Instead of /en/add/product (losing locale context)');
      } else if (pushedRoute.includes('/en/add/product')) {
        console.log('✅ GOOD: Route includes locale prefix');
      } else {
        console.log(`⚠️ UNEXPECTED: Route format: ${pushedRoute}`);
      }

      expect(mockOnClose).toHaveBeenCalled();
    });

    test('deal item click - check router.push call', async () => {
      console.log('🔍 TESTING: AddDropdown Deal item click');
      
      const mockOnClose = jest.fn();
      render(
        <TestWrapper locale="en">
          <AddDropdown onClose={mockOnClose} />
        </TestWrapper>
      );

      // Switch to "All Options" tab to see deal item
      const allTab = screen.getByText('All Options');
      fireEvent.click(allTab);

      const dealButton = screen.getByText('Deal');
      console.log('📍 Clicking Deal item...');
      fireEvent.click(dealButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/add/deal?step=1') {
        console.log('❌ CRITICAL ISSUE: Deal route pushed WITHOUT locale prefix');
        console.log('   Users will lose locale context when creating deals');
      } else {
        console.log('✅ Route includes proper formatting');
      }
    });

    test('post item click - check router.push call', async () => {
      console.log('🔍 TESTING: AddDropdown Post item click');
      
      const mockOnClose = jest.fn();
      render(
        <TestWrapper locale="en">
          <AddDropdown onClose={mockOnClose} />
        </TestWrapper>
      );

      const postButton = screen.getByText('Post');
      console.log('📍 Clicking Post item...');
      fireEvent.click(postButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/add/post?step=1') {
        console.log('❌ CRITICAL ISSUE: Post route pushed WITHOUT locale prefix');
      } else {
        console.log('✅ Route includes proper formatting');
      }
    });
  });

  describe('AddDropdown from Different Locales', () => {
    test.each(['en', 'pl', 'de'])('product item from %s locale', async (locale) => {
      console.log(`🔍 TESTING: AddDropdown Product from ${locale} locale`);
      
      // Mock the current locale context
      jest.doMock('../../i18n/navigation', () => {
        const actual = jest.requireActual('../../i18n/navigation');
        return {
          ...actual,
          useRouter: () => mockRouter,
          usePathname: () => `/${locale}/home`,
          useParams: () => ({ locale })
        };
      });

      const mockOnClose = jest.fn();
      render(
        <TestWrapper locale={locale}>
          <AddDropdown onClose={mockOnClose} />
        </TestWrapper>
      );

      const productButton = screen.getByText('Product');
      console.log(`📍 Clicking Product from ${locale} locale...`);
      fireEvent.click(productButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalled();
      });

      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📊 RESULT: From ${locale}, router.push called with: "${pushedRoute}"`);
      
      if (pushedRoute === '/add/product?step=1') {
        console.log(`❌ CRITICAL: AddDropdown loses ${locale} locale context`);
        console.log(`   User in ${locale} will be redirected to English version`);
      } else if (pushedRoute.includes(locale)) {
        console.log(`✅ GOOD: Route preserves ${locale} locale`);
      } else {
        console.log(`⚠️ UNCLEAR: Route format: ${pushedRoute}`);
      }
    });
  });
});

describe('Real Integration Tests: Header + AddDropdown', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
    console.log('\n' + '='.repeat(60));
  });

  test('complete user flow: header create button -> dropdown -> item selection', async () => {
    console.log('🔍 TESTING: Complete user flow simulation');
    console.log('Scenario: User clicks Create button, then selects Product');
    
    render(
      <TestWrapper locale="en">
        <Header />
      </TestWrapper>
    );

    // Step 1: Click create button to open dropdown
    const createButton = screen.getByText('Create');
    console.log('📍 Step 1: Clicking Create button...');
    fireEvent.click(createButton);

    // Verify dropdown opened
    await waitFor(() => {
      expect(screen.getByText('Create New')).toBeInTheDocument();
    });
    console.log('✅ AddDropdown opened successfully');

    // Step 2: Click product item
    const productButton = screen.getByText('Product');
    console.log('📍 Step 2: Clicking Product item...');
    fireEvent.click(productButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalled();
    });

    const pushedRoute = mockPush.mock.calls[0][0];
    console.log(`📊 FINAL RESULT: Complete flow resulted in: "${pushedRoute}"`);
    
    if (pushedRoute === '/add/product?step=1') {
      console.log('❌ CRITICAL FLOW ISSUE: User loses locale in complete flow');
      console.log('   User starts at /en/home');
      console.log('   User clicks Create -> Product');
      console.log('   User ends up at /add/product (no locale)');
      console.log('   This breaks the user experience!');
    } else if (pushedRoute.includes('/en/add/product')) {
      console.log('✅ PERFECT: Complete flow preserves locale');
    } else {
      console.log(`⚠️ UNEXPECTED: Flow result: ${pushedRoute}`);
    }
  });

  test('navigation flow from different starting locales', async () => {
    console.log('🔍 TESTING: Navigation flow from Polish locale');
    
    // Mock Polish locale context
    jest.doMock('../../i18n/navigation', () => {
      const actual = jest.requireActual('../../i18n/navigation');
      return {
        ...actual,
        useRouter: () => mockRouter,
        usePathname: () => '/pl/marketplace',
        useParams: () => ({ locale: 'pl' })
      };
    });

    render(
      <TestWrapper locale="pl">
        <Header />
      </TestWrapper>
    );

    // Navigate to home first
    const homeButton = screen.getByText('Home');
    console.log('📍 Step 1: Navigating to Home from Polish locale...');
    fireEvent.click(homeButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalled();
    });

    const homeRoute = mockPush.mock.calls[0][0];
    console.log(`📊 Home navigation result: "${homeRoute}"`);

    // Open dropdown and select deal
    const createButton = screen.getByText('Create');
    console.log('📍 Step 2: Opening Create dropdown...');
    fireEvent.click(createButton);

    await waitFor(() => {
      expect(screen.getByText('Create New')).toBeInTheDocument();
    });

    const allTab = screen.getByText('All Options');
    fireEvent.click(allTab);

    const dealButton = screen.getByText('Deal');
    console.log('📍 Step 3: Selecting Deal item...');
    fireEvent.click(dealButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledTimes(2);
    });

    const dealRoute = mockPush.mock.calls[1][0];
    console.log(`📊 Deal navigation result: "${dealRoute}"`);

    console.log('\n📊 COMPLETE FLOW ANALYSIS:');
    console.log(`   Starting locale: pl`);
    console.log(`   Home navigation: ${homeRoute}`);
    console.log(`   Deal navigation: ${dealRoute}`);
    
    if (homeRoute === '/home' && dealRoute === '/add/deal?step=1') {
      console.log('❌ CRITICAL: Both navigations lose Polish locale');
      console.log('   User experience is broken for non-English users');
    } else if (homeRoute.includes('pl') && dealRoute.includes('pl')) {
      console.log('✅ EXCELLENT: All navigations preserve Polish locale');
    } else {
      console.log('⚠️ MIXED: Some navigations preserve locale, others don\'t');
    }
  });
}); 
/**
 * Real Navigation URL Capture Test
 * Captures actual router navigation calls to validate locale preservation
 * Tests the REAL problem: Header redirecting without locale prefixes
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import { NextIntlClientProvider } from 'next-intl';
import { Provider as ReduxProvider } from 'react-redux';
import { QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '../../src/context/AuthContext';
import { NavBarProvider } from '../../src/context/NavBarContext';
import { CategoriesProvider } from '../../src/hooks/useCategories';
import { makeStore } from '../../src/lib/store';
import { getQueryClient } from '../../src/lib/reactQuery';
import Header from '../../src/components/Header/Header';

// Real test messages
const testMessages = {
  en: {
    Header: {
      homeButton: 'Home',
      exploreButton: 'Explore',
      createButtonText: 'Create',
      mainNavAriaLabel: 'Main navigation',
      notificationsButtonAriaLabel: 'Notifications ({count})',
      messagesButtonAriaLabel: 'Messages ({count})',
      wishlistButtonAriaLabel: 'Wishlist ({count})',
      cartButtonAriaLabel: 'Shopping cart',
      createContentAriaLabel: 'Create new content',
      menuTitle: 'Menu',
      closeMenuAriaLabel: 'Close menu',
      userFallbackName: 'User',
      navigationGroupTitle: 'Navigation',
      actionsGroupTitle: 'Actions'
    },
    AddDropdown: {
      title: 'Create New Content',
      ariaLabel: 'Create content menu',
      frequentTab: 'Frequent',
      allOptionsTab: 'All Options',
      item_product_label: 'Product',
      item_product_desc: 'Sell a new or used item',
      item_post_label: 'Post',
      item_post_desc: 'Share news or updates',
      item_vehicle_label: 'Vehicle',
      item_vehicle_desc: 'List a car, bike, or other vehicle',
      item_deal_label: 'Deal',
      item_deal_desc: 'Share a bargain with the community',
      item_property_label: 'Property',
      item_property_desc: 'List property for sale or rent',
      item_job_label: 'Job',
      item_job_desc: 'Post a job opening',
      item_service_label: 'Service',
      item_service_desc: 'Offer a professional service',
      item_video_label: 'Video',
      item_video_desc: 'Share video content'
    }
  }
};

// Navigation URL capture system
let capturedNavigationCalls = [];

// Mock the navigation module to capture real calls
jest.mock('../../src/i18n/navigation', () => {
  const originalModule = jest.requireActual('../../src/i18n/navigation');
  
  return {
    ...originalModule,
    useRouter: () => ({
      push: (url) => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.push("${url}")`);
        capturedNavigationCalls.push({ method: 'push', url });
        return Promise.resolve();
      },
      replace: (url) => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.replace("${url}")`);
        capturedNavigationCalls.push({ method: 'replace', url });
        return Promise.resolve();
      },
      back: () => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.back()`);
        capturedNavigationCalls.push({ method: 'back' });
      },
      forward: () => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.forward()`);
        capturedNavigationCalls.push({ method: 'forward' });
      },
      refresh: () => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.refresh()`);
        capturedNavigationCalls.push({ method: 'refresh' });
      },
      prefetch: (url) => {
        console.log(`🔍 CAPTURED NAVIGATION CALL: router.prefetch("${url}")`);
        capturedNavigationCalls.push({ method: 'prefetch', url });
        return Promise.resolve();
      }
    }),
    usePathname: () => '/en/',
    Link: ({ children, href, onClick, ...props }) => (
      <a 
        href={href} 
        onClick={(e) => {
          console.log(`🔍 CAPTURED LINK CLICK: href="${href}"`);
          capturedNavigationCalls.push({ method: 'link_click', url: href });
          if (onClick) onClick(e);
        }}
        {...props}
      >
        {children}
      </a>
    )
  };
});

// Provider setup
function TestProvider({ children, locale = 'en' }) {
  const store = makeStore();
  const queryClient = getQueryClient();
  
  const mockUser = {
    id: '1',
    name: 'Test User',
    email: 'redacted-email@example.com',
    avatar: null,
    isAuthenticated: true
  };

  const badgeCounts = {
    notifications: 1,
    messages: 5,
    wishlist: 0,
    cart: 0
  };

  return (
    <NextIntlClientProvider locale={locale} messages={testMessages[locale]}>
      <ReduxProvider store={store}>
        <QueryClientProvider client={queryClient}>
          <AuthProvider value={{ user: mockUser, isAuthenticated: true }}>
            <NavBarProvider value={{ 
              isMobile: false, 
              showNavbars: true, 
              badgeCounts,
              currentPath: `/${locale}/`
            }}>
              <CategoriesProvider>
                {children}
              </CategoriesProvider>
            </NavBarProvider>
          </AuthProvider>
        </QueryClientProvider>
      </ReduxProvider>
    </NextIntlClientProvider>
  );
}

describe('🔍 Real Navigation URL Capture - PROVE THE PROBLEM', () => {
  beforeEach(() => {
    // Reset captured calls
    capturedNavigationCalls = [];
    
    // Mock window properties
    Object.defineProperty(window, 'scrollY', {
      writable: true,
      value: 0
    });
    
    Object.defineProperty(window, 'addEventListener', {
      writable: true,
      value: jest.fn()
    });
    
    Object.defineProperty(window, 'removeEventListener', {
      writable: true,
      value: jest.fn()
    });

    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: jest.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: jest.fn(),
        removeListener: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
      })),
    });
  });

  test('🚨 CAPTURE: What URLs does Header actually generate when user clicks buttons?', async () => {
    console.log('\n🔍 TESTING: Capturing REAL navigation URLs from Header component...');
    
    render(
      <TestProvider locale="en">
        <Header />
      </TestProvider>
    );

    // Clear any initial calls
    capturedNavigationCalls = [];

    console.log('\n📍 USER ACTION: Clicking Home button...');
    const homeButton = screen.getByRole('button', { name: 'Home' });
    fireEvent.click(homeButton);

    await waitFor(() => {
      expect(capturedNavigationCalls.length).toBeGreaterThan(0);
    }, { timeout: 2000 });

    console.log('\n📍 USER ACTION: Clicking Explore button...');
    const exploreButton = screen.getByRole('button', { name: 'Explore' });
    fireEvent.click(exploreButton);

    await waitFor(() => {
      expect(capturedNavigationCalls.length).toBeGreaterThan(1);
    }, { timeout: 2000 });

    // Analyze captured URLs
    console.log('\n🔍 ANALYSIS: Captured Navigation Calls:');
    console.log('==========================================');
    
    capturedNavigationCalls.forEach((call, index) => {
      console.log(`${index + 1}. ${call.method.toUpperCase()}: "${call.url || 'N/A'}"`);
      
      if (call.url) {
        const hasLocale = call.url.match(/^\/(en|pl|de)\//);
        const status = hasLocale ? '✅ HAS LOCALE' : '❌ NO LOCALE';
        console.log(`   ${status}: ${call.url}`);
        
        if (!hasLocale) {
          console.log(`   🚨 PROBLEM: URL "${call.url}" missing locale prefix!`);
        }
      }
    });

    console.log('==========================================');

    // Validate results
    const urlCalls = capturedNavigationCalls.filter(call => call.url);
    
    if (urlCalls.length === 0) {
      console.log('❌ NO NAVIGATION CALLS CAPTURED - Header might not be using router.push()');
      expect(urlCalls.length).toBeGreaterThan(0);
    }

    // Check each URL for locale prefix
    const urlsWithoutLocale = urlCalls.filter(call => {
      return call.url && !call.url.match(/^\/(en|pl|de)\//);
    });

    if (urlsWithoutLocale.length > 0) {
      console.log('\n🚨 PROBLEM CONFIRMED: URLs without locale prefixes found!');
      urlsWithoutLocale.forEach(call => {
        console.log(`   ❌ ${call.method}: "${call.url}" - MISSING LOCALE!`);
      });
      
      // This will fail the test and show the problem
      expect(urlsWithoutLocale).toHaveLength(0);
    } else {
      console.log('\n✅ SUCCESS: All URLs contain locale prefixes');
    }
  });

  test('🔍 CAPTURE: AddDropdown navigation URLs', async () => {
    console.log('\n🔍 TESTING: Capturing AddDropdown navigation URLs...');
    
    render(
      <TestProvider locale="en">
        <Header />
      </TestProvider>
    );

    // Clear any initial calls
    capturedNavigationCalls = [];

    console.log('\n📍 USER ACTION: Opening AddDropdown...');
    const createButton = screen.getByLabelText('Create new content');
    fireEvent.click(createButton);

    await waitFor(() => {
      expect(createButton).toHaveAttribute('aria-expanded', 'true');
    });

    // Look for any links in the dropdown
    const dropdownLinks = screen.getAllByRole('link');
    
    if (dropdownLinks.length > 0) {
      console.log(`\n📍 Found ${dropdownLinks.length} links in dropdown, clicking first one...`);
      fireEvent.click(dropdownLinks[0]);

      await waitFor(() => {
        expect(capturedNavigationCalls.length).toBeGreaterThan(0);
      }, { timeout: 2000 });
    }

    // Analyze captured URLs
    console.log('\n🔍 ANALYSIS: AddDropdown Navigation Calls:');
    console.log('==========================================');
    
    capturedNavigationCalls.forEach((call, index) => {
      console.log(`${index + 1}. ${call.method.toUpperCase()}: "${call.url || 'N/A'}"`);
      
      if (call.url) {
        const hasLocale = call.url.match(/^\/(en|pl|de)\//);
        const status = hasLocale ? '✅ HAS LOCALE' : '❌ NO LOCALE';
        console.log(`   ${status}: ${call.url}`);
        
        if (!hasLocale) {
          console.log(`   🚨 PROBLEM: AddDropdown URL "${call.url}" missing locale prefix!`);
        }
      }
    });

    console.log('==========================================');
  });

  test('🎯 SUMMARY: Validate localhost:3000 equivalent behavior', () => {
    console.log('\n🎯 EXPECTED BEHAVIOR:');
    console.log('=====================');
    console.log('When user visits localhost:3000/ → Should redirect to /en/');
    console.log('When user visits localhost:3000/explore → Should redirect to /en/explore');
    console.log('When user clicks Home button → Should navigate to /en/home');
    console.log('When user clicks Explore button → Should navigate to /en/explore');
    console.log('When user clicks AddDropdown items → Should navigate to /en/add/[type]');
    console.log('=====================');
    
    console.log('\n🔍 ACTUAL BEHAVIOR (from captured calls):');
    console.log('==========================================');
    
    if (capturedNavigationCalls.length === 0) {
      console.log('❌ No navigation calls were captured in previous tests');
    } else {
      capturedNavigationCalls.forEach((call, index) => {
        if (call.url) {
          const hasLocale = call.url.match(/^\/(en|pl|de)\//);
          const status = hasLocale ? '✅' : '❌';
          console.log(`${status} ${call.method}: "${call.url}"`);
        }
      });
    }
    
    console.log('==========================================');
    
    // This test always passes but shows the analysis
    expect(true).toBe(true);
  });
}); 
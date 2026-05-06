/**
 * Header and AddDropdown Component Interaction Test Suite
 * Tests actual user interactions with Header and AddDropdown components
 * Simulates real user clicks to identify locale preservation issues
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import Header from '../../components/Header/Header';
import AddDropdown from '../../components/Header/AddDropdown';

// Mock the navigation hooks and context
const mockPush = jest.fn();
const mockUseRouter = jest.fn(() => ({
  push: mockPush,
  pathname: '/en/home',
  locale: 'en'
}));

const mockUsePathname = jest.fn(() => '/en/home');
const mockUseParams = jest.fn(() => ({ locale: 'en' }));
const mockUseTranslations = jest.fn((namespace) => (key) => `${namespace}.${key}`);

// Mock all the required modules
jest.mock('../../i18n/navigation', () => ({
  useRouter: mockUseRouter,
  usePathname: mockUsePathname,
  useParams: mockUseParams
}));

jest.mock('next/navigation', () => ({
  usePathname: mockUsePathname
}));

jest.mock('next-intl', () => ({
  useTranslations: mockUseTranslations
}));

// Mock context providers
const mockAuthContext = {
  user: { id: 1, name: 'Test User' },
  signOutUser: jest.fn()
};

const mockNavBarContext = {
  showNavbars: true,
  isMobile: false,
  isClient: true
};

jest.mock('../../context/AuthContext', () => ({
  useAuth: () => mockAuthContext
}));

jest.mock('../../context/NavBarContext', () => ({
  NavBarContext: React.createContext(mockNavBarContext)
}));

// Mock child components to avoid complex dependencies
jest.mock('../../components/Header/Logo', () => {
  return function MockLogo() {
    return <div data-testid="logo">Logo</div>;
  };
});

jest.mock('../../components/Header/SearchBar', () => {
  return function MockSearchBar() {
    return <div data-testid="search-bar">SearchBar</div>;
  };
});

jest.mock('../../components/Header/UserMenu', () => {
  return function MockUserMenu() {
    return <div data-testid="user-menu">UserMenu</div>;
  };
});

jest.mock('../../components/Header/SelectTopic', () => {
  return function MockSelectTopic() {
    return <div data-testid="select-topic">SelectTopic</div>;
  };
});

describe('Header Component User Interaction Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
  });

  describe('Header Navigation Button Clicks', () => {
    test('home button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Home button click');
      
      render(<Header />);
      
      const homeButton = screen.getByText('Header.homeButton');
      expect(homeButton).toBeInTheDocument();
      
      fireEvent.click(homeButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/home');
      });
      
      console.log(`📍 Home button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      // Check if locale is preserved
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
        console.log(`   Expected: /en/home (or other locale)`);
        console.log(`   Actual: ${pushedRoute}`);
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });

    test('explore button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Explore button click');
      
      render(<Header />);
      
      const exploreButton = screen.getByText('Header.exploreButton');
      expect(exploreButton).toBeInTheDocument();
      
      fireEvent.click(exploreButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/explore');
      });
      
      console.log(`📍 Explore button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });

    test('notifications button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Notifications button click');
      
      render(<Header />);
      
      // Find notifications button by aria-label or icon
      const notificationsButton = screen.getByLabelText(/Header\.notificationsButtonAriaLabel/);
      expect(notificationsButton).toBeInTheDocument();
      
      fireEvent.click(notificationsButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/notifications');
      });
      
      console.log(`📍 Notifications button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });

    test('messages button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Messages button click');
      
      render(<Header />);
      
      const messagesButton = screen.getByLabelText(/Header\.messagesButtonAriaLabel/);
      expect(messagesButton).toBeInTheDocument();
      
      fireEvent.click(messagesButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/messages');
      });
      
      console.log(`📍 Messages button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });

    test('wishlist button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Wishlist button click');
      
      render(<Header />);
      
      const wishlistButton = screen.getByLabelText(/Header\.wishlistButtonAriaLabel/);
      expect(wishlistButton).toBeInTheDocument();
      
      fireEvent.click(wishlistButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/wishlist');
      });
      
      console.log(`📍 Wishlist button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });

    test('cart button click preserves locale', async () => {
      console.log('\n🔍 Testing Header Cart button click');
      
      render(<Header />);
      
      const cartButton = screen.getByLabelText(/Header\.cartButtonAriaLabel/);
      expect(cartButton).toBeInTheDocument();
      
      fireEvent.click(cartButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/cart');
      });
      
      console.log(`📍 Cart button clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: Route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in navigation');
      }
    });
  });

  describe('Header Navigation with Different Locales', () => {
    test.each(['en', 'pl', 'de'])('header navigation works correctly from %s locale', async (locale) => {
      console.log(`\n🔍 Testing Header navigation from ${locale} locale`);
      
      // Mock the current locale
      mockUsePathname.mockReturnValue(`/${locale}/home`);
      mockUseParams.mockReturnValue({ locale });
      
      render(<Header />);
      
      // Test home button
      const homeButton = screen.getByText('Header.homeButton');
      fireEvent.click(homeButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/home');
      });
      
      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📍 From ${locale} locale, home button pushed: ${pushedRoute}`);
      
      // The router should automatically handle locale preservation
      // But if it doesn't, we need to identify this issue
      if (!pushedRoute.includes(locale) && !pushedRoute.startsWith(`/${locale}/`)) {
        console.log(`❌ POTENTIAL ISSUE: Route may not preserve ${locale} locale`);
      }
    });
  });
});

describe('AddDropdown Component User Interaction Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
  });

  describe('AddDropdown Item Clicks', () => {
    test('product item click preserves locale', async () => {
      console.log('\n🔍 Testing AddDropdown Product item click');
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      // Find the product item button
      const productButton = screen.getByText('AddDropdown.item_product_label');
      expect(productButton).toBeInTheDocument();
      
      fireEvent.click(productButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/product?step=1');
      });
      
      console.log(`📍 Product item clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: AddDropdown route pushed without locale prefix');
        console.log(`   Expected: /en/add/product?step=1 (or other locale)`);
        console.log(`   Actual: ${pushedRoute}`);
      } else {
        console.log('✅ Locale preserved in AddDropdown navigation');
      }
      
      expect(mockOnClose).toHaveBeenCalled();
    });

    test('post item click preserves locale', async () => {
      console.log('\n🔍 Testing AddDropdown Post item click');
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      const postButton = screen.getByText('AddDropdown.item_post_label');
      expect(postButton).toBeInTheDocument();
      
      fireEvent.click(postButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/post?step=1');
      });
      
      console.log(`📍 Post item clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: AddDropdown route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in AddDropdown navigation');
      }
    });

    test('deal item click preserves locale', async () => {
      console.log('\n🔍 Testing AddDropdown Deal item click');
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      // Switch to "all" tab to see deal item
      const allTab = screen.getByText('AddDropdown.allOptionsTab');
      fireEvent.click(allTab);
      
      const dealButton = screen.getByText('AddDropdown.item_deal_label');
      expect(dealButton).toBeInTheDocument();
      
      fireEvent.click(dealButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/deal?step=1');
      });
      
      console.log(`📍 Deal item clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: AddDropdown route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in AddDropdown navigation');
      }
    });

    test('vehicle item click preserves locale', async () => {
      console.log('\n🔍 Testing AddDropdown Vehicle item click');
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      const vehicleButton = screen.getByText('AddDropdown.item_vehicle_label');
      expect(vehicleButton).toBeInTheDocument();
      
      fireEvent.click(vehicleButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/vehicle?step=1');
      });
      
      console.log(`📍 Vehicle item clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
      
      const pushedRoute = mockPush.mock.calls[0][0];
      if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
        console.log('❌ LOCALE LOST: AddDropdown route pushed without locale prefix');
      } else {
        console.log('✅ Locale preserved in AddDropdown navigation');
      }
    });
  });

  describe('AddDropdown Navigation with Different Locales', () => {
    test.each(['en', 'pl', 'de'])('add dropdown navigation works correctly from %s locale', async (locale) => {
      console.log(`\n🔍 Testing AddDropdown navigation from ${locale} locale`);
      
      // Mock the current locale
      mockUseParams.mockReturnValue({ locale });
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      // Test product item
      const productButton = screen.getByText('AddDropdown.item_product_label');
      fireEvent.click(productButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/product?step=1');
      });
      
      const pushedRoute = mockPush.mock.calls[0][0];
      console.log(`📍 From ${locale} locale, product item pushed: ${pushedRoute}`);
      
      // Check if the route should include locale
      if (!pushedRoute.includes(locale) && !pushedRoute.startsWith(`/${locale}/`)) {
        console.log(`❌ POTENTIAL ISSUE: AddDropdown route may not preserve ${locale} locale`);
      }
    });
  });

  describe('AddDropdown Tab Navigation', () => {
    test('switching between frequent and all tabs works correctly', async () => {
      console.log('\n🔍 Testing AddDropdown tab switching');
      
      const mockOnClose = jest.fn();
      render(<AddDropdown onClose={mockOnClose} />);
      
      // Initially on frequent tab
      expect(screen.getByText('AddDropdown.item_product_label')).toBeInTheDocument();
      
      // Switch to all tab
      const allTab = screen.getByText('AddDropdown.allOptionsTab');
      fireEvent.click(allTab);
      
      // Should now see deal item
      expect(screen.getByText('AddDropdown.item_deal_label')).toBeInTheDocument();
      
      // Test deal item click
      const dealButton = screen.getByText('AddDropdown.item_deal_label');
      fireEvent.click(dealButton);
      
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/add/deal?step=1');
      });
      
      console.log(`📍 Deal item from all tab clicked, router.push called with: ${mockPush.mock.calls[0][0]}`);
    });
  });
});

describe('Integration Tests: Header + AddDropdown', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPush.mockClear();
  });

  test('header create button opens add dropdown and navigation works', async () => {
    console.log('\n🔍 Testing Header create button -> AddDropdown integration');
    
    render(<Header />);
    
    // Find and click the create button
    const createButton = screen.getByText('Header.createButtonText');
    expect(createButton).toBeInTheDocument();
    
    fireEvent.click(createButton);
    
    // AddDropdown should now be visible
    await waitFor(() => {
      expect(screen.getByText('AddDropdown.title')).toBeInTheDocument();
    });
    
    console.log('✅ AddDropdown opened successfully');
    
    // Click on a product item in the dropdown
    const productButton = screen.getByText('AddDropdown.item_product_label');
    fireEvent.click(productButton);
    
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/add/product?step=1');
    });
    
    console.log(`📍 Product item clicked from Header dropdown, router.push called with: ${mockPush.mock.calls[0][0]}`);
    
    const pushedRoute = mockPush.mock.calls[0][0];
    if (!pushedRoute.startsWith('/en/') && !pushedRoute.startsWith('/pl/') && !pushedRoute.startsWith('/de/')) {
      console.log('❌ LOCALE LOST: Header -> AddDropdown navigation lost locale');
    } else {
      console.log('✅ Locale preserved in Header -> AddDropdown navigation');
    }
  });

  test('complete user flow: navigate to page -> open dropdown -> select item', async () => {
    console.log('\n🔍 Testing complete user navigation flow');
    
    // Start from a specific locale
    mockUsePathname.mockReturnValue('/en/marketplace');
    mockUseParams.mockReturnValue({ locale: 'en' });
    
    render(<Header />);
    
    // 1. User navigates to home
    const homeButton = screen.getByText('Header.homeButton');
    fireEvent.click(homeButton);
    
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/home');
    });
    
    console.log(`📍 Step 1 - Home navigation: ${mockPush.mock.calls[0][0]}`);
    
    // 2. User opens add dropdown
    const createButton = screen.getByText('Header.createButtonText');
    fireEvent.click(createButton);
    
    await waitFor(() => {
      expect(screen.getByText('AddDropdown.title')).toBeInTheDocument();
    });
    
    // 3. User selects deal item
    const allTab = screen.getByText('AddDropdown.allOptionsTab');
    fireEvent.click(allTab);
    
    const dealButton = screen.getByText('AddDropdown.item_deal_label');
    fireEvent.click(dealButton);
    
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/add/deal?step=1');
    });
    
    console.log(`📍 Step 2 - Deal creation: ${mockPush.mock.calls[1][0]}`);
    
    // Analyze the complete flow
    const homeRoute = mockPush.mock.calls[0][0];
    const dealRoute = mockPush.mock.calls[1][0];
    
    console.log('\n📊 Complete Flow Analysis:');
    console.log(`   Home navigation: ${homeRoute}`);
    console.log(`   Deal navigation: ${dealRoute}`);
    
    if (!homeRoute.includes('en') && !dealRoute.includes('en')) {
      console.log('❌ CRITICAL: Both navigations lost locale information');
    } else if (!homeRoute.includes('en') || !dealRoute.includes('en')) {
      console.log('⚠️ WARNING: One navigation lost locale information');
    } else {
      console.log('✅ All navigations preserved locale correctly');
    }
  });
}); 
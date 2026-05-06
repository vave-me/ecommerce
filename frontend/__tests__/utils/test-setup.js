/**
 * Comprehensive Test Setup Utilities
 * Provides real instances for testing without mocks
 * Merged from existing test utilities and enhanced with real next-intl support
 */

import React from 'react';
import { render } from '@testing-library/react';
import { NextIntlClientProvider } from 'next-intl';
import { Provider as ReduxProvider } from 'react-redux';
import { QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '../../src/context/AuthContext';
import { NavBarProvider } from '../../src/context/NavBarContext';
import { CategoriesProvider } from '../../src/hooks/useCategories';
import { makeStore } from '../../src/lib/store';
import { getQueryClient } from '../../src/lib/reactQuery';
import { jest } from '@jest/globals';

// Mock router for testing navigation
export function createMockRouter(options = {}) {
  const {
    pathname = '/',
    push = jest.fn(),
    replace = jest.fn(),
    back = jest.fn(),
    forward = jest.fn(),
    refresh = jest.fn(),
    prefetch = jest.fn()
  } = options;

  return {
    push,
    replace,
    back,
    forward,
    refresh,
    prefetch,
    pathname
  };
}

// Mock the navigation modules for testing
jest.mock('../../src/i18n/navigation', () => {
  let mockRouter = null;
  let mockPathname = '/';

  return {
    useRouter: () => mockRouter || createMockRouter(),
    usePathname: () => mockPathname,
    Link: ({ children, href, ...props }) => (
      <a href={href} {...props}>{children}</a>
    ),
    redirect: jest.fn(),
    // Test utilities to control the mocks
    __setMockRouter: (router) => { mockRouter = router; },
    __setMockPathname: (pathname) => { mockPathname = pathname; },
    __resetMocks: () => { 
      mockRouter = null; 
      mockPathname = '/'; 
    }
  };
});

jest.mock('next/navigation', () => ({
  usePathname: () => require('../../src/i18n/navigation').usePathname(),
  useSearchParams: () => new URLSearchParams(),
  useParams: () => ({ locale: 'en' }),
}));

// Real test messages for comprehensive testing (merged from existing setup)
export const testMessages = {
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
      recentlyUsedTitle: 'Recently Used',
      addAriaLabel: 'Add {label}',
      section_frequent_title: 'Frequently Used',
      section_all_title: 'All Options',
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
    },
    BottomNav: {
      navAriaLabel: 'Mobile Navigation',
      badgeAriaLabel: '{count, plural, =0 {no new items} one {# new item} other {# new items}}',
      nav_home_label: 'Home',
      nav_home_tooltip: 'Go to Home Page',
      nav_explore_label: 'Explore',
      nav_explore_tooltip: 'Discover New Content',
      nav_add_label: 'Add',
      nav_add_tooltip_default: 'Create New Content',
      nav_add_tooltip_alternate: 'Add New Item',
      nav_messages_label: 'Messages',
      nav_messages_tooltip: 'View Your Messages',
      nav_alerts_label: 'Alerts',
      nav_alerts_tooltip: 'View Notifications',
      nav_cart_label: 'Cart',
      nav_cart_tooltip: 'View Shopping Cart',
      nav_wishlist_label: 'Wishlist',
      nav_wishlist_tooltip: 'View Your Wishlist',
      additem_product_label: 'Product',
      additem_product_desc: 'Sell a new or used item.',
      additem_post_label: 'Post',
      additem_post_desc: 'Share news, an update, or write an article.',
      additem_video_label: 'Video',
      additem_video_desc: 'Share a video, like a short clip or longer content.',
      additem_vehicle_label: 'Vehicle',
      additem_vehicle_desc: 'List a car, bike, or other vehicle for sale.',
      additem_deal_label: 'Deal',
      additem_deal_desc: 'Found a bargain? Share the deal with the community.',
      additem_property_label: 'Property',
      additem_property_desc: 'List a property for sale or rent.',
      additem_service_label: 'Service',
      additem_service_desc: 'Offer a professional service.',
      additem_job_label: 'Job',
      additem_job_desc: 'Post a job opening to find candidates.'
    },
    HomePage: {
      title: 'Welcome',
      description: 'Test description',
      help: 'Help',
      contact: 'Contact',
      about: 'About',
      terms: 'Terms',
      privacy: 'Privacy',
      allRightsReserved: 'All rights reserved'
    },
    Seo: {
      title: 'sfx markt – Live Marketplace',
      description: 'sfx markt is the live marketplace that lets you buy, sell and connect with your community in real time.',
      keywords: 'sfx markt, marketplace, buy and sell locally, real‑time chat, SafePay'
    },
    Navigation: {
      home: 'Home',
      about: 'About',
      contact: 'Contact',
      login: 'Login',
      logout: 'Logout'
    },
    Common: {
      loading: 'Loading...',
      error: 'An error occurred',
      retry: 'Retry',
      success: 'Success',
      cancel: 'Cancel',
      save: 'Save',
      delete: 'Delete',
      edit: 'Edit',
      close: 'Close'
    },
    Topics: {
      topicsAriaLabel: 'Topic navigation',
      topic_dashboard_label: 'Dashboard',
      topic_marketplace_label: 'Marketplace',
      topic_deals_label: 'Deals',
      topic_jobs_label: 'Jobs',
      topic_news_info_label: 'News & Info',
      topic_social_content_label: 'Social',
      topic_property_label: 'Property',
      topic_services_label: 'Services',
      topic_videos_label: 'Videos',
      topic_automotive_label: 'Automotive',
      dashboard_badge: 'New',
      marketplace_badge: 'Hot',
      cars_badge: 'Popular'
    },
    UserMenu: {
      fallbackMenuAriaLabel: 'User menu',
      profileLabel: 'Profile',
      settingsLabel: 'Settings',
      logoutLabel: 'Logout',
      loginLabel: 'Login'
    }
  },
  pl: {
    Header: {
      homeButton: 'Strona główna',
      exploreButton: 'Eksploruj',
      createButtonText: 'Utwórz',
      mainNavAriaLabel: 'Główna nawigacja',
      notificationsButtonAriaLabel: 'Powiadomienia ({count})',
      messagesButtonAriaLabel: 'Wiadomości ({count})',
      wishlistButtonAriaLabel: 'Lista życzeń ({count})',
      cartButtonAriaLabel: 'Koszyk',
      createContentAriaLabel: 'Utwórz nową treść',
      menuTitle: 'Menu',
      closeMenuAriaLabel: 'Zamknij menu',
      userFallbackName: 'Użytkownik',
      navigationGroupTitle: 'Nawigacja',
      actionsGroupTitle: 'Akcje'
    },
    AddDropdown: {
      title: 'Utwórz nową treść',
      ariaLabel: 'Menu tworzenia treści',
      frequentTab: 'Często używane',
      allOptionsTab: 'Wszystkie opcje',
      recentlyUsedTitle: 'Ostatnio używane',
      addAriaLabel: 'Dodaj {label}',
      section_frequent_title: 'Często używane',
      section_all_title: 'Wszystkie opcje',
      item_product_label: 'Produkt',
      item_product_desc: 'Sprzedaj nowy lub używany przedmiot',
      item_post_label: 'Post',
      item_post_desc: 'Udostępnij wiadomości lub aktualizacje',
      item_vehicle_label: 'Pojazd',
      item_vehicle_desc: 'Wystaw samochód, rower lub inny pojazd',
      item_deal_label: 'Okazja',
      item_deal_desc: 'Podziel się okazją ze społecznością',
      item_property_label: 'Nieruchomość',
      item_property_desc: 'Wystaw nieruchomość na sprzedaż lub wynajem',
      item_job_label: 'Praca',
      item_job_desc: 'Opublikuj ofertę pracy',
      item_service_label: 'Usługa',
      item_service_desc: 'Oferuj profesjonalną usługę',
      item_video_label: 'Wideo',
      item_video_desc: 'Udostępnij treści wideo'
    },
    BottomNav: {
      navAriaLabel: 'Nawigacja mobilna',
      badgeAriaLabel: '{count, plural, =0 {brak nowych elementów} one {# nowy element} other {# nowych elementów}}',
      nav_home_label: 'Strona główna',
      nav_home_tooltip: 'Przejdź do strony głównej',
      nav_explore_label: 'Eksploruj',
      nav_explore_tooltip: 'Odkryj nowe treści',
      nav_add_label: 'Dodaj',
      nav_add_tooltip_default: 'Utwórz nową treść',
      nav_add_tooltip_alternate: 'Dodaj nowy element',
      nav_messages_label: 'Wiadomości',
      nav_messages_tooltip: 'Zobacz swoje wiadomości',
      nav_alerts_label: 'Powiadomienia',
      nav_alerts_tooltip: 'Zobacz powiadomienia',
      nav_cart_label: 'Koszyk',
      nav_cart_tooltip: 'Zobacz koszyk',
      nav_wishlist_label: 'Lista życzeń',
      nav_wishlist_tooltip: 'Zobacz swoją listę życzeń',
      additem_product_label: 'Produkt',
      additem_product_desc: 'Sprzedaj nowy lub używany przedmiot.',
      additem_post_label: 'Post',
      additem_post_desc: 'Udostępnij wiadomości, aktualizację lub napisz artykuł.',
      additem_video_label: 'Wideo',
      additem_video_desc: 'Udostępnij wideo, jak krótki klip lub dłuższą treść.',
      additem_vehicle_label: 'Pojazd',
      additem_vehicle_desc: 'Wystaw samochód, rower lub inny pojazd na sprzedaż.',
      additem_deal_label: 'Okazja',
      additem_deal_desc: 'Znalazłeś okazję? Podziel się nią ze społecznością.',
      additem_property_label: 'Nieruchomość',
      additem_property_desc: 'Wystaw nieruchomość na sprzedaż lub wynajem.',
      additem_service_label: 'Usługa',
      additem_service_desc: 'Oferuj profesjonalną usługę.',
      additem_job_label: 'Praca',
      additem_job_desc: 'Opublikuj ofertę pracy, aby znaleźć kandydatów.'
    },
    HomePage: {
      title: 'Witamy',
      description: 'Opis testowy',
      help: 'Pomoc',
      contact: 'Kontakt',
      about: 'O nas',
      terms: 'Warunki',
      privacy: 'Prywatność',
      allRightsReserved: 'Wszystkie prawa zastrzeżone'
    },
    Seo: {
      title: 'sfx markt – Marketplace na Żywo',
      description: 'sfx markt to  marketplace na żywo, który pozwala kupować, sprzedawać i łączyć się ze społecznością w czasie rzeczywistym.',
      keywords: 'sfx markt, marketplace, kupuj i sprzedawaj lokalnie, czat w czasie rzeczywistym, SafePay'
    },
    Navigation: {
      home: 'Strona główna',
      about: 'O nas',
      contact: 'Kontakt',
      login: 'Zaloguj',
      logout: 'Wyloguj'
    },
    Common: {
      loading: 'Ładowanie...',
      error: 'Wystąpił błąd',
      retry: 'Ponów',
      success: 'Sukces',
      cancel: 'Anuluj',
      save: 'Zapisz',
      delete: 'Usuń',
      edit: 'Edytuj',
      close: 'Zamknij'
    },
    Topics: {
      topicsAriaLabel: 'Topic navigation',
      topic_dashboard_label: 'Dashboard',
      topic_marketplace_label: 'Marketplace',
      topic_deals_label: 'Deals',
      topic_jobs_label: 'Jobs',
      topic_news_info_label: 'News & Info',
      topic_social_content_label: 'Social',
      topic_property_label: 'Property',
      topic_services_label: 'Services',
      topic_videos_label: 'Videos',
      topic_automotive_label: 'Automotive',
      dashboard_badge: 'New',
      marketplace_badge: 'Hot',
      cars_badge: 'Popular'
    },
    UserMenu: {
      fallbackMenuAriaLabel: 'User menu',
      profileLabel: 'Profile',
      settingsLabel: 'Settings',
      logoutLabel: 'Logout',
      loginLabel: 'Login'
    }
  },
  de: {
    Header: {
      homeButton: 'Startseite',
      exploreButton: 'Entdecken',
      createButtonText: 'Erstellen',
      mainNavAriaLabel: 'Hauptnavigation',
      notificationsButtonAriaLabel: 'Benachrichtigungen ({count})',
      messagesButtonAriaLabel: 'Nachrichten ({count})',
      wishlistButtonAriaLabel: 'Wunschliste ({count})',
      cartButtonAriaLabel: 'Warenkorb',
      createContentAriaLabel: 'Neue Inhalte erstellen',
      menuTitle: 'Menü',
      closeMenuAriaLabel: 'Menü schließen',
      userFallbackName: 'Benutzer',
      navigationGroupTitle: 'Navigation',
      actionsGroupTitle: 'Aktionen'
    },
    AddDropdown: {
      title: 'Neue Inhalte erstellen',
      ariaLabel: 'Inhalte erstellen Menü',
      frequentTab: 'Häufig verwendet',
      allOptionsTab: 'Alle Optionen',
      recentlyUsedTitle: 'Kürzlich verwendet',
      addAriaLabel: '{label} hinzufügen',
      section_frequent_title: 'Häufig verwendet',
      section_all_title: 'Alle Optionen',
      item_product_label: 'Produkt',
      item_product_desc: 'Verkaufen Sie einen neuen oder gebrauchten Artikel',
      item_post_label: 'Beitrag',
      item_post_desc: 'Teilen Sie Nachrichten oder Updates',
      item_vehicle_label: 'Fahrzeug',
      item_vehicle_desc: 'Listen Sie ein Auto, Fahrrad oder anderes Fahrzeug auf',
      item_deal_label: 'Angebot',
      item_deal_desc: 'Teilen Sie ein Schnäppchen mit der Community',
      item_property_label: 'Immobilie',
      item_property_desc: 'Listen Sie Immobilien zum Verkauf oder zur Miete auf',
      item_job_label: 'Job',
      item_job_desc: 'Veröffentlichen Sie eine Stellenausschreibung',
      item_service_label: 'Service',
      item_service_desc: 'Bieten Sie einen professionellen Service an',
      item_video_label: 'Video',
      item_video_desc: 'Teilen Sie Videoinhalte'
    },
    BottomNav: {
      navAriaLabel: 'Mobile Navigation',
      badgeAriaLabel: '{count, plural, =0 {keine neuen Elemente} one {# neues Element} other {# neue Elemente}}',
      nav_home_label: 'Startseite',
      nav_home_tooltip: 'Zur Startseite gehen',
      nav_explore_label: 'Entdecken',
      nav_explore_tooltip: 'Neue Inhalte entdecken',
      nav_add_label: 'Hinzufügen',
      nav_add_tooltip_default: 'Neue Inhalte erstellen',
      nav_add_tooltip_alternate: 'Neues Element hinzufügen',
      nav_messages_label: 'Nachrichten',
      nav_messages_tooltip: 'Ihre Nachrichten anzeigen',
      nav_alerts_label: 'Benachrichtigungen',
      nav_alerts_tooltip: 'Benachrichtigungen anzeigen',
      nav_cart_label: 'Warenkorb',
      nav_cart_tooltip: 'Warenkorb anzeigen',
      nav_wishlist_label: 'Wunschliste',
      nav_wishlist_tooltip: 'Ihre Wunschliste anzeigen',
      additem_product_label: 'Produkt',
      additem_product_desc: 'Verkaufen Sie einen neuen oder gebrauchten Artikel.',
      additem_post_label: 'Beitrag',
      additem_post_desc: 'Teilen Sie Nachrichten, ein Update oder schreiben Sie einen Artikel.',
      additem_video_label: 'Video',
      additem_video_desc: 'Teilen Sie ein Video, wie einen kurzen Clip oder längeren Inhalt.',
      additem_vehicle_label: 'Fahrzeug',
      additem_vehicle_desc: 'Listen Sie ein Auto, Fahrrad oder anderes Fahrzeug zum Verkauf auf.',
      additem_deal_label: 'Angebot',
      additem_deal_desc: 'Ein Schnäppchen gefunden? Teilen Sie das Angebot mit der Community.',
      additem_property_label: 'Immobilie',
      additem_property_desc: 'Listen Sie eine Immobilie zum Verkauf oder zur Miete auf.',
      additem_service_label: 'Service',
      additem_service_desc: 'Bieten Sie einen professionellen Service an.',
      additem_job_label: 'Job',
      additem_job_desc: 'Veröffentlichen Sie eine Stellenausschreibung, um Kandidaten zu finden.'
    },
    HomePage: {
      title: 'Willkommen',
      description: 'Test Beschreibung',
      help: 'Hilfe',
      contact: 'Kontakt',
      about: 'Über uns',
      terms: 'Bedingungen',
      privacy: 'Datenschutz',
      allRightsReserved: 'Alle Rechte vorbehalten'
    },
    Seo: {
      title: 'sfx markt – Live Marktplatz',
      description: 'sfx markt ist der live  Marktplatz, der es Ihnen ermöglicht, in Echtzeit zu kaufen, zu verkaufen und sich mit Ihrer Gemeinschaft zu verbinden.',
      keywords: 'sfx markt, Marktplatz, lokal kaufen und verkaufen, Echtzeit-Chat, SafePay'
    },
    Navigation: {
      home: 'Startseite',
      about: 'Über uns',
      contact: 'Kontakt',
      login: 'Anmelden',
      logout: 'Abmelden'
    },
    Common: {
      loading: 'Laden...',
      error: 'Ein Fehler ist aufgetreten',
      retry: 'Wiederholen',
      success: 'Erfolg',
      cancel: 'Abbrechen',
      save: 'Speichern',
      delete: 'Löschen',
      edit: 'Bearbeiten',
      close: 'Schließen'
    },
    Topics: {
      topicsAriaLabel: 'Topic navigation',
      topic_dashboard_label: 'Dashboard',
      topic_marketplace_label: 'Marketplace',
      topic_deals_label: 'Deals',
      topic_jobs_label: 'Jobs',
      topic_news_info_label: 'News & Info',
      topic_social_content_label: 'Social',
      topic_property_label: 'Property',
      topic_services_label: 'Services',
      topic_videos_label: 'Videos',
      topic_automotive_label: 'Automotive',
      dashboard_badge: 'New',
      marketplace_badge: 'Hot',
      cars_badge: 'Popular'
    },
    UserMenu: {
      fallbackMenuAriaLabel: 'User menu',
      profileLabel: 'Profile',
      settingsLabel: 'Settings',
      logoutLabel: 'Logout',
      loginLabel: 'Login'
    }
  }
};

/**
 * Real Test Provider that wraps components with all necessary contexts
 */
export function RealTestProvider({ 
  children, 
  locale = 'en', 
  user = null, 
  isMobile = false,
  showNavbars = true,
  isClient = true,
  pathname = '/',
  customMessages = null
}) {
  const store = makeStore();
  const queryClient = getQueryClient();
  const messages = customMessages || testMessages[locale] || testMessages.en;

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <ReduxProvider store={store}>
        <QueryClientProvider client={queryClient}>
          <AuthProvider initialUser={user}>
            <NavBarProvider initialState={{ isMobile, showNavbars, isClient }}>
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

/**
 * Custom render function with real providers
 */
export function renderWithRealProviders(ui, options = {}) {
  const {
    locale = 'en',
    user = null,
    isMobile = false,
    showNavbars = true,
    isClient = true,
    pathname = '/',
    currentPath = '/',
    customMessages = null,
    navigationTracker = null,
    ...renderOptions
  } = options;

  // Create mock router for testing
  const mockRouter = createMockRouter({
    pathname: currentPath,
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    forward: jest.fn(),
    refresh: jest.fn(),
    prefetch: jest.fn()
  });

  // Set up the navigation mocks
  const navigation = require('../../src/i18n/navigation');
  navigation.__setMockRouter(mockRouter);
  navigation.__setMockPathname(currentPath);

  function Wrapper({ children }) {
    return (
      <RealTestProvider
        locale={locale}
        user={user}
        isMobile={isMobile}
        showNavbars={showNavbars}
        isClient={isClient}
        pathname={pathname}
        customMessages={customMessages}
      >
        {children}
      </RealTestProvider>
    );
  }

  const result = render(ui, { wrapper: Wrapper, ...renderOptions });
  
  // Return both the render result and the mockRouter for testing
  return {
    ...result,
    mockRouter
  };
}

/**
 * Mock window utilities for testing responsive behavior
 */
export function mockWindowMatchMedia(matches = false) {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: jest.fn().mockImplementation(query => ({
      matches,
      media: query,
      onchange: null,
      addListener: jest.fn(),
      removeListener: jest.fn(),
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      dispatchEvent: jest.fn(),
    })),
  });
}

/**
 * Mock scroll utilities for testing scroll behavior
 */
export function mockWindowScroll() {
  Object.defineProperty(window, 'scrollY', {
    writable: true,
    value: 0
  });

  const mockAddEventListener = jest.fn();
  const mockRemoveEventListener = jest.fn();
  
  Object.defineProperty(window, 'addEventListener', {
    writable: true,
    value: mockAddEventListener
  });
  
  Object.defineProperty(window, 'removeEventListener', {
    writable: true,
    value: mockRemoveEventListener
  });

  return { mockAddEventListener, mockRemoveEventListener };
}

/**
 * Test user factory for creating test users
 */
export function createTestUser(overrides = {}) {
  return {
    userId: 'test-user-123',
    email: 'redacted-email@example.com',
    username: 'testuser',
    firstName: 'Test',
    lastName: 'User',
    avatar: null,
    isVerified: true,
    ...overrides
  };
}

/**
 * Test badge counts for navigation components
 */
export const testBadgeCounts = {
  messages: 5,
  notifications: 3,
  wishlist: 2,
  cart: 1
};

/**
 * Common test utilities (merged from existing test-utils)
 */
export const testUtils = {
  // Wait for async operations
  waitFor: (callback, timeout = 3000) => {
    return new Promise((resolve, reject) => {
      const startTime = Date.now();
      const check = () => {
        try {
          const result = callback();
          if (result) {
            resolve(result);
          } else if (Date.now() - startTime > timeout) {
            reject(new Error('Timeout waiting for condition'));
          } else {
            setTimeout(check, 100);
          }
        } catch (error) {
          if (Date.now() - startTime > timeout) {
            reject(error);
          } else {
            setTimeout(check, 100);
          }
        }
      };
      check();
    });
  },

  // Simulate user interactions
  simulateClick: async (element) => {
    element.click();
    await new Promise(resolve => setTimeout(resolve, 100));
  },

  // Performance measurement
  measurePerformance: (fn) => {
    const start = performance.now();
    const result = fn();
    const end = performance.now();
    return { result, duration: end - start };
  },

  // Modal test utilities (merged from modal-test-utils)
  waitForModal: async (modalTestId, timeout = 3000) => {
    const startTime = Date.now();
    while (Date.now() - startTime < timeout) {
      const modal = document.querySelector(`[data-testid="${modalTestId}"]`);
      if (modal) return modal;
      await new Promise(resolve => setTimeout(resolve, 100));
    }
    throw new Error(`Modal with testId "${modalTestId}" not found within ${timeout}ms`);
  },

  closeModal: async (modalTestId) => {
    const modal = await testUtils.waitForModal(modalTestId);
    const closeButton = modal.querySelector('[data-testid="close-button"]') || 
                      modal.querySelector('.close-button') ||
                      modal.querySelector('[aria-label*="close"]');
    if (closeButton) {
      closeButton.click();
    }
  }
};

// Legacy exports for backward compatibility
export const renderWithProviders = renderWithRealProviders;
export const TestProvider = RealTestProvider;

export default {
  RealTestProvider,
  renderWithRealProviders,
  mockWindowMatchMedia,
  mockWindowScroll,
  createTestUser,
  testMessages,
  testBadgeCounts,
  testUtils
}; 
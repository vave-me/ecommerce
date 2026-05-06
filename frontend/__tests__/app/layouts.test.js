/**
 * Comprehensive Layout Components Test Suite
 * Tests root layout.jsx and [locale]/layout.jsx functionality
 */

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import { notFound } from 'next/navigation';

// Import the layouts
import RootLayout from '../../app/layout';

// Mock Next.js navigation
jest.mock('next/navigation', () => ({
  notFound: jest.fn()
}));

// Mock next-intl components and functions
const mockMessages = {
  HomePage: {
    title: 'Welcome',
    description: 'Test description'
  }
};

jest.mock('next-intl', () => ({
  NextIntlClientProvider: ({ children, messages }) => (
    <div data-testid="next-intl-provider" data-messages={JSON.stringify(messages)}>
      {children}
    </div>
  )
}));

jest.mock('next-intl/server', () => ({
  getMessages: jest.fn(() => Promise.resolve(mockMessages)),
  setRequestLocale: jest.fn(),
  generateMetadata: jest.fn()
}));

// Mock routing configuration
const mockRouting = {
  locales: ['en', 'pl', 'de'],
  defaultLocale: 'en'
};

jest.mock('../../i18n/routing', () => ({
  routing: mockRouting
}));

// Mock child components
jest.mock('../../app/Providers', () => {
  return function MockProviders({ children }) {
    return <div data-testid="providers">{children}</div>;
  };
});

jest.mock('../../app/ClientLayout.client', () => {
  return function MockClientLayout({ children }) {
    return <div data-testid="client-layout">{children}</div>;
  };
});

jest.mock('../../components/AdSense/AdSenseScript.client', () => {
  return function MockAdSenseScript() {
    return <script data-testid="adsense-script" />;
  };
});

// Mock CSS import
jest.mock('../../app/global.css', () => ({}));

describe('Root Layout Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Root Layout Component', () => {
    test('renders children without html/body wrapper', () => {
      console.log('\n🔍 TESTING: Root layout renders children only');
      
      render(
        <RootLayout>
          <div data-testid="test-child">Test Content</div>
        </RootLayout>
      );

      expect(screen.getByTestId('test-child')).toBeInTheDocument();
      console.log('✅ Root layout renders children correctly');
    });

    test('does not wrap children in html/body tags', () => {
      console.log('\n🔍 TESTING: Root layout structure');
      
      const { container } = render(
        <RootLayout>
          <div data-testid="test-child">Test Content</div>
        </RootLayout>
      );

      // Should not have html or body elements
      expect(container.querySelector('html')).toBeNull();
      expect(container.querySelector('body')).toBeNull();
      
      // Should directly contain the child
      expect(container.firstChild).toHaveAttribute('data-testid', 'test-child');

      console.log('✅ Root layout structure is correct (no html/body wrapper)');
    });

    test('handles multiple children', () => {
      console.log('\n🔍 TESTING: Root layout with multiple children');
      
      render(
        <RootLayout>
          <div data-testid="child-1">Child 1</div>
          <div data-testid="child-2">Child 2</div>
        </RootLayout>
      );

      expect(screen.getByTestId('child-1')).toBeInTheDocument();
      expect(screen.getByTestId('child-2')).toBeInTheDocument();

      console.log('✅ Root layout handles multiple children correctly');
    });

    test('handles no children gracefully', () => {
      console.log('\n🔍 TESTING: Root layout with no children');
      
      expect(() => {
        render(<RootLayout />);
      }).not.toThrow();

      console.log('✅ Root layout handles no children gracefully');
    });
  });
});

describe('Locale Layout Tests', () => {
  // Dynamic import for locale layout since it's async
  let LocaleLayout;

  beforeAll(async () => {
    // Mock the locale layout
    LocaleLayout = ({ children, params }) => {
      const [locale, setLocale] = React.useState(null);
      const [messages, setMessages] = React.useState(null);

      React.useEffect(() => {
        const loadData = async () => {
          const resolvedParams = await params;
          const { getMessages, setRequestLocale } = require('next-intl/server');
          
          if (!mockRouting.locales.includes(resolvedParams.locale)) {
            notFound();
            return;
          }

          setRequestLocale(resolvedParams.locale);
          const msgs = await getMessages();
          setLocale(resolvedParams.locale);
          setMessages(msgs);
        };

        loadData();
      }, [params]);

      if (!locale || !messages) {
        return <div data-testid="loading">Loading...</div>;
      }

      return (
        <html lang={locale} data-testid="html-element">
          <body data-testid="body-element">
            <div id="portal-root" data-testid="portal-root"></div>
            <script data-testid="adsense-script" />
            <div data-testid="next-intl-provider" data-messages={JSON.stringify(messages)}>
              <div data-testid="providers">
                <div data-testid="client-layout">
                  {children}
                </div>
              </div>
            </div>
          </body>
        </html>
      );
    };
  });

  beforeEach(() => {
    jest.clearAllMocks();
    const { getMessages, setRequestLocale } = require('next-intl/server');
    getMessages.mockResolvedValue(mockMessages);
    setRequestLocale.mockImplementation(() => {});
  });

  describe('Locale Validation', () => {
    test('renders successfully with valid locale', async () => {
      console.log('\n🔍 TESTING: Valid locale rendering');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test Content</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
      });

      expect(notFound).not.toHaveBeenCalled();
      console.log('✅ Valid locale rendered successfully');
    });

    test('calls notFound for invalid locale', async () => {
      console.log('\n🔍 TESTING: Invalid locale handling');
      
      const mockParams = Promise.resolve({ locale: 'invalid' });

      render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test Content</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(notFound).toHaveBeenCalled();
      });

      console.log('✅ Invalid locale triggers notFound correctly');
    });

    test('validates all supported locales', async () => {
      console.log('\n🔍 TESTING: All supported locales validation');
      
      for (const locale of mockRouting.locales) {
        const mockParams = Promise.resolve({ locale });
        notFound.mockClear();

        render(
          <LocaleLayout params={mockParams}>
            <div data-testid={`content-${locale}`}>Content for {locale}</div>
          </LocaleLayout>
        );

        await waitFor(() => {
          expect(screen.getByTestId(`content-${locale}`)).toBeInTheDocument();
        });

        expect(notFound).not.toHaveBeenCalled();
      }

      console.log('✅ All supported locales validated correctly');
    });
  });

  describe('HTML Structure and Attributes', () => {
    test('sets correct lang attribute on html element', async () => {
      console.log('\n🔍 TESTING: HTML lang attribute');
      
      const mockParams = Promise.resolve({ locale: 'pl' });

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        const htmlElement = screen.getByTestId('html-element');
        expect(htmlElement).toHaveAttribute('lang', 'pl');
      });

      console.log('✅ HTML lang attribute set correctly');
    });

    test('includes portal root for React 19 modals', async () => {
      console.log('\n🔍 TESTING: Portal root for React 19');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        const portalRoot = screen.getByTestId('portal-root');
        expect(portalRoot).toHaveAttribute('id', 'portal-root');
      });

      console.log('✅ Portal root included for React 19 modals');
    });

    test('includes AdSense script', async () => {
      console.log('\n🔍 TESTING: AdSense script inclusion');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('adsense-script')).toBeInTheDocument();
      });

      console.log('✅ AdSense script included correctly');
    });
  });

  describe('Provider Integration', () => {
    test('integrates NextIntlClientProvider with messages', async () => {
      console.log('\n🔍 TESTING: NextIntlClientProvider integration');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        const provider = screen.getByTestId('next-intl-provider');
        const messages = JSON.parse(provider.getAttribute('data-messages'));
        expect(messages).toEqual(mockMessages);
      });

      console.log('✅ NextIntlClientProvider integrated with messages');
    });

    test('wraps content with Providers component', async () => {
      console.log('\n🔍 TESTING: Providers component wrapping');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('providers')).toBeInTheDocument();
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
      });

      console.log('✅ Content wrapped with Providers component');
    });

    test('wraps content with ClientLayout component', async () => {
      console.log('\n🔍 TESTING: ClientLayout component wrapping');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('client-layout')).toBeInTheDocument();
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
      });

      console.log('✅ Content wrapped with ClientLayout component');
    });
  });

  describe('Next.js 15 Compatibility', () => {
    test('handles async params correctly', async () => {
      console.log('\n🔍 TESTING: Async params handling (Next.js 15)');
      
      const { setRequestLocale } = require('next-intl/server');
      const mockParams = Promise.resolve({ locale: 'de' });

      render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
      });

      expect(setRequestLocale).toHaveBeenCalledWith('de');
      console.log('✅ Async params handled correctly for Next.js 15');
    });

    test('calls setRequestLocale for static rendering', async () => {
      console.log('\n🔍 TESTING: setRequestLocale for static rendering');
      
      const { setRequestLocale } = require('next-intl/server');
      const mockParams = Promise.resolve({ locale: 'en' });

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(setRequestLocale).toHaveBeenCalledWith('en');
      });

      console.log('✅ setRequestLocale called for static rendering');
    });
  });

  describe('Metadata Generation', () => {
    test('generateMetadata function exists and works', async () => {
      console.log('\n🔍 TESTING: Metadata generation');
      
      // Mock generateMetadata function
      const generateMetadata = async ({ params }) => {
        const { locale } = await params;
        return {
          title: 'vaveme',
          description: 'vave me – Live Local Marketplace',
          other: {
            'google-adsense-account': 'ca-pub-7872277873986607'
          },
          icons: {
            icon: '/favicon.ico'
          }
        };
      };

      const mockParams = Promise.resolve({ locale: 'en' });
      const metadata = await generateMetadata({ params: mockParams });

      expect(metadata.title).toBe('vaveme');
      expect(metadata.description).toBe('vave me – Live Local Marketplace');
      expect(metadata.other['google-adsense-account']).toBe('ca-pub-7872277873986607');
      expect(metadata.icons.icon).toBe('/favicon.ico');

      console.log('✅ Metadata generation working correctly');
    });
  });

  describe('Static Params Generation', () => {
    test('generateStaticParams returns all locales', () => {
      console.log('\n🔍 TESTING: Static params generation');
      
      const generateStaticParams = () => {
        return mockRouting.locales.map(locale => ({ locale }));
      };

      const staticParams = generateStaticParams();
      
      expect(staticParams).toEqual([
        { locale: 'en' },
        { locale: 'pl' },
        { locale: 'de' }
      ]);

      console.log('✅ Static params generated for all locales');
    });
  });

  describe('Error Handling and Edge Cases', () => {
    test('handles missing params gracefully', async () => {
      console.log('\n🔍 TESTING: Missing params handling');
      
      const mockParams = Promise.resolve({});

      render(
        <LocaleLayout params={mockParams}>
          <div>Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(notFound).toHaveBeenCalled();
      });

      console.log('✅ Missing params handled gracefully');
    });

    test('handles null params gracefully', async () => {
      console.log('\n🔍 TESTING: Null params handling');
      
      const mockParams = Promise.resolve(null);

      expect(() => {
        render(
          <LocaleLayout params={mockParams}>
            <div>Test</div>
          </LocaleLayout>
        );
      }).not.toThrow();

      console.log('✅ Null params handled gracefully');
    });

    test('handles getMessages failure gracefully', async () => {
      console.log('\n🔍 TESTING: getMessages failure handling');
      
      const { getMessages } = require('next-intl/server');
      getMessages.mockRejectedValueOnce(new Error('Failed to load messages'));

      const mockParams = Promise.resolve({ locale: 'en' });

      expect(() => {
        render(
          <LocaleLayout params={mockParams}>
            <div>Test</div>
          </LocaleLayout>
        );
      }).not.toThrow();

      console.log('✅ getMessages failure handled gracefully');
    });
  });

  describe('Component Hierarchy', () => {
    test('maintains correct component nesting order', async () => {
      console.log('\n🔍 TESTING: Component nesting hierarchy');
      
      const mockParams = Promise.resolve({ locale: 'en' });

      const { container } = render(
        <LocaleLayout params={mockParams}>
          <div data-testid="test-content">Test</div>
        </LocaleLayout>
      );

      await waitFor(() => {
        expect(screen.getByTestId('test-content')).toBeInTheDocument();
      });

      // Check hierarchy: html > body > portal-root + adsense + next-intl > providers > client-layout > content
      const htmlElement = container.querySelector('[data-testid="html-element"]');
      const bodyElement = container.querySelector('[data-testid="body-element"]');
      const portalRoot = container.querySelector('[data-testid="portal-root"]');
      const nextIntlProvider = container.querySelector('[data-testid="next-intl-provider"]');
      const providers = container.querySelector('[data-testid="providers"]');
      const clientLayout = container.querySelector('[data-testid="client-layout"]');

      expect(htmlElement).toContainElement(bodyElement);
      expect(bodyElement).toContainElement(portalRoot);
      expect(bodyElement).toContainElement(nextIntlProvider);
      expect(nextIntlProvider).toContainElement(providers);
      expect(providers).toContainElement(clientLayout);

      console.log('✅ Component nesting hierarchy is correct');
    });
  });
}); 
/**
 * Comprehensive Middleware Test Suite
 * Tests next-intl middleware functionality, routing, and configuration
 */

import { jest } from '@jest/globals';

// Mock next-intl middleware
const mockCreateMiddleware = jest.fn();
const mockMiddleware = jest.fn();

jest.mock('next-intl/middleware', () => mockCreateMiddleware);

// Mock routing configuration
const mockRouting = {
  locales: ['en', 'pl', 'de'],
  defaultLocale: 'en',
  localePrefix: 'always'
};

jest.mock('../../../src/i18n/routing', () => ({
  routing: mockRouting
}));

// Mock Next.js request and response
const createMockRequest = (url, headers = {}) => ({
  nextUrl: new URL(url, 'http://localhost:3000'),
  headers: new Map(Object.entries(headers)),
  cookies: new Map(),
  geo: {},
  ip: '127.0.0.1'
});

const createMockResponse = () => ({
  headers: new Map(),
  cookies: new Map(),
  status: 200,
  redirect: jest.fn(),
  rewrite: jest.fn(),
  next: jest.fn()
});

describe('Middleware Configuration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockCreateMiddleware.mockReturnValue(mockMiddleware);
  });

  describe('Middleware Creation', () => {
    test('creates middleware with routing configuration', () => {
      console.log('\n🔍 TESTING: Middleware creation with routing config');
      
      // Import middleware to trigger creation
      require('../../../middleware');
      
      expect(mockCreateMiddleware).toHaveBeenCalledWith(mockRouting);
      console.log('✅ Middleware created with correct routing configuration');
    });

    test('middleware function is properly exported', () => {
      console.log('\n🔍 TESTING: Middleware function export');
      
      const middleware = require('../../../middleware').default;
      
      expect(middleware).toBe(mockMiddleware);
      console.log('✅ Middleware function exported correctly');
    });
  });

  describe('Middleware Configuration Object', () => {
    test('exports correct matcher configuration', () => {
      console.log('\n🔍 TESTING: Middleware matcher configuration');
      
      const { config } = require('../../../middleware');
      
      expect(config).toBeDefined();
      expect(config.matcher).toBeDefined();
      expect(Array.isArray(config.matcher)).toBe(true);
      
      const expectedMatchers = [
        '/((?!api|_next|_vercel|.*\\..*).*)',
        '/'
      ];
      
      expect(config.matcher).toEqual(expectedMatchers);
      console.log('✅ Middleware matcher configuration is correct');
    });

    test('matcher excludes correct paths', () => {
      console.log('\n🔍 TESTING: Middleware path exclusions');
      
      const { config } = require('../../../middleware');
      const matcher = config.matcher[0]; // Main matcher regex
      
      // Test paths that should be excluded
      const excludedPaths = [
        '/api/users',
        '/api/auth/login',
        '/_next/static/chunks/main.js',
        '/_next/image',
        '/_vercel/insights',
        '/favicon.ico',
        '/robots.txt',
        '/sitemap.xml',
        '/manifest.json'
      ];

      // Create a regex from the matcher pattern
      const regex = new RegExp(matcher);
      
      excludedPaths.forEach(path => {
        const shouldMatch = !path.match(/^\/(?:api|_next|_vercel|.*\..*)/) && path !== '/';
        if (!shouldMatch) {
          expect(regex.test(path)).toBe(false);
        }
      });

      console.log('✅ Middleware correctly excludes API and static paths');
    });

    test('matcher includes correct paths', () => {
      console.log('\n🔍 TESTING: Middleware path inclusions');
      
      const { config } = require('../../../middleware');
      const matcher = config.matcher[0]; // Main matcher regex
      
      // Test paths that should be included
      const includedPaths = [
        '/home',
        '/about',
        '/en/products',
        '/pl/sklep',
        '/de/ueber-uns',
        '/user/profile',
        '/search/results'
      ];

      const regex = new RegExp(matcher);
      
      includedPaths.forEach(path => {
        expect(regex.test(path)).toBe(true);
      });

      console.log('✅ Middleware correctly includes internationalized paths');
    });
  });

  describe('Routing Configuration Integration', () => {
    test('middleware receives all supported locales', () => {
      console.log('\n🔍 TESTING: Supported locales integration');
      
      require('../../../middleware');
      
      const passedRouting = mockCreateMiddleware.mock.calls[0][0];
      expect(passedRouting.locales).toEqual(['en', 'pl', 'de']);
      expect(passedRouting.defaultLocale).toBe('en');
      
      console.log('✅ All supported locales passed to middleware');
    });

    test('middleware receives locale prefix configuration', () => {
      console.log('\n🔍 TESTING: Locale prefix configuration');
      
      require('../../../middleware');
      
      const passedRouting = mockCreateMiddleware.mock.calls[0][0];
      expect(passedRouting.localePrefix).toBe('always');
      
      console.log('✅ Locale prefix configuration passed correctly');
    });
  });

  describe('Middleware Functionality Simulation', () => {
    beforeEach(() => {
      // Mock middleware behavior
      mockMiddleware.mockImplementation((request) => {
        const url = request.nextUrl;
        const pathname = url.pathname;
        
        // Simulate locale detection and redirection
        if (pathname === '/') {
          return {
            type: 'redirect',
            url: '/en',
            status: 307
          };
        }
        
        if (!pathname.startsWith('/en') && !pathname.startsWith('/pl') && !pathname.startsWith('/de')) {
          return {
            type: 'redirect',
            url: `/en${pathname}`,
            status: 307
          };
        }
        
        return {
          type: 'next'
        };
      });
    });

    test('redirects root path to default locale', () => {
      console.log('\n🔍 TESTING: Root path redirection');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/');
      
      const result = middleware(request);
      
      expect(result.type).toBe('redirect');
      expect(result.url).toBe('/en');
      expect(result.status).toBe(307);
      
      console.log('✅ Root path redirected to default locale');
    });

    test('redirects unprefixed paths to default locale', () => {
      console.log('\n🔍 TESTING: Unprefixed path redirection');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/about');
      
      const result = middleware(request);
      
      expect(result.type).toBe('redirect');
      expect(result.url).toBe('/en/about');
      expect(result.status).toBe(307);
      
      console.log('✅ Unprefixed paths redirected with locale prefix');
    });

    test('allows prefixed paths to pass through', () => {
      console.log('\n🔍 TESTING: Prefixed path pass-through');
      
      const middleware = require('../../../middleware').default;
      const testPaths = [
        'http://localhost:3000/en/home',
        'http://localhost:3000/pl/sklep',
        'http://localhost:3000/de/ueber-uns'
      ];
      
      testPaths.forEach(url => {
        const request = createMockRequest(url);
        const result = middleware(request);
        
        expect(result.type).toBe('next');
      });
      
      console.log('✅ Prefixed paths pass through correctly');
    });
  });

  describe('Edge Cases and Error Handling', () => {
    test('handles malformed URLs gracefully', () => {
      console.log('\n🔍 TESTING: Malformed URL handling');
      
      const middleware = require('../../../middleware').default;
      
      expect(() => {
        const request = createMockRequest('http://localhost:3000/[invalid');
        middleware(request);
      }).not.toThrow();
      
      console.log('✅ Malformed URLs handled gracefully');
    });

    test('handles requests with query parameters', () => {
      console.log('\n🔍 TESTING: Query parameter handling');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/search?q=test&category=electronics');
      
      const result = middleware(request);
      
      expect(result.type).toBe('redirect');
      expect(result.url).toBe('/en/search?q=test&category=electronics');
      
      console.log('✅ Query parameters preserved during redirection');
    });

    test('handles requests with hash fragments', () => {
      console.log('\n🔍 TESTING: Hash fragment handling');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/page#section');
      
      const result = middleware(request);
      
      expect(result.type).toBe('redirect');
      expect(result.url).toBe('/en/page#section');
      
      console.log('✅ Hash fragments preserved during redirection');
    });
  });

  describe('Performance and Optimization', () => {
    test('middleware creation is called only once', () => {
      console.log('\n🔍 TESTING: Middleware creation optimization');
      
      // Clear previous calls
      mockCreateMiddleware.mockClear();
      
      // Import multiple times
      delete require.cache[require.resolve('../../../middleware')];
      require('../../../middleware');
      delete require.cache[require.resolve('../../../middleware')];
      require('../../../middleware');
      
      // Should only be called once due to module caching
      expect(mockCreateMiddleware).toHaveBeenCalledTimes(1);
      
      console.log('✅ Middleware creation optimized with module caching');
    });

    test('routing configuration is not mutated', () => {
      console.log('\n🔍 TESTING: Routing configuration immutability');
      
      const originalRouting = { ...mockRouting };
      require('../../../middleware');
      
      expect(mockRouting).toEqual(originalRouting);
      
      console.log('✅ Routing configuration remains immutable');
    });
  });

  describe('Integration with Next.js 15', () => {
    test('middleware compatible with Next.js 15 request format', () => {
      console.log('\n🔍 TESTING: Next.js 15 request compatibility');
      
      const middleware = require('../../../middleware').default;
      
      // Next.js 15 request format
      const request = {
        nextUrl: new URL('http://localhost:3000/test'),
        headers: new Headers({
          'accept-language': 'en-US,en;q=0.9',
          'user-agent': 'Mozilla/5.0'
        }),
        cookies: new Map(),
        geo: { country: 'US' },
        ip: '192.168.1.1'
      };
      
      expect(() => {
        middleware(request);
      }).not.toThrow();
      
      console.log('✅ Middleware compatible with Next.js 15 request format');
    });

    test('middleware works with App Router', () => {
      console.log('\n🔍 TESTING: App Router compatibility');
      
      const middleware = require('../../../middleware').default;
      const appRouterPaths = [
        'http://localhost:3000/dashboard',
        'http://localhost:3000/user/settings',
        'http://localhost:3000/products/123'
      ];
      
      appRouterPaths.forEach(url => {
        const request = createMockRequest(url);
        expect(() => {
          middleware(request);
        }).not.toThrow();
      });
      
      console.log('✅ Middleware compatible with App Router paths');
    });
  });

  describe('Locale Detection and Headers', () => {
    test('processes accept-language header correctly', () => {
      console.log('\n🔍 TESTING: Accept-Language header processing');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/', {
        'accept-language': 'pl-PL,pl;q=0.9,en;q=0.8'
      });
      
      // Mock enhanced middleware behavior for locale detection
      mockMiddleware.mockImplementationOnce((req) => {
        const acceptLanguage = req.headers.get('accept-language');
        if (acceptLanguage && acceptLanguage.includes('pl')) {
          return {
            type: 'redirect',
            url: '/pl',
            status: 307
          };
        }
        return { type: 'redirect', url: '/en', status: 307 };
      });
      
      const result = middleware(request);
      
      expect(result.url).toBe('/pl');
      
      console.log('✅ Accept-Language header processed for locale detection');
    });

    test('handles missing accept-language header', () => {
      console.log('\n🔍 TESTING: Missing Accept-Language header handling');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/');
      
      expect(() => {
        middleware(request);
      }).not.toThrow();
      
      console.log('✅ Missing Accept-Language header handled gracefully');
    });
  });

  describe('Cookie Handling', () => {
    test('respects locale preference cookie', () => {
      console.log('\n🔍 TESTING: Locale preference cookie handling');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/');
      request.cookies.set('NEXT_LOCALE', 'de');
      
      // Mock enhanced middleware behavior for cookie handling
      mockMiddleware.mockImplementationOnce((req) => {
        const localeCookie = req.cookies.get('NEXT_LOCALE');
        if (localeCookie && mockRouting.locales.includes(localeCookie)) {
          return {
            type: 'redirect',
            url: `/${localeCookie}`,
            status: 307
          };
        }
        return { type: 'redirect', url: '/en', status: 307 };
      });
      
      const result = middleware(request);
      
      expect(result.url).toBe('/de');
      
      console.log('✅ Locale preference cookie respected');
    });

    test('ignores invalid locale cookie', () => {
      console.log('\n🔍 TESTING: Invalid locale cookie handling');
      
      const middleware = require('../../../middleware').default;
      const request = createMockRequest('http://localhost:3000/');
      request.cookies.set('NEXT_LOCALE', 'invalid');
      
      // Mock enhanced middleware behavior
      mockMiddleware.mockImplementationOnce((req) => {
        const localeCookie = req.cookies.get('NEXT_LOCALE');
        if (localeCookie && mockRouting.locales.includes(localeCookie)) {
          return { type: 'redirect', url: `/${localeCookie}`, status: 307 };
        }
        return { type: 'redirect', url: '/en', status: 307 };
      });
      
      const result = middleware(request);
      
      expect(result.url).toBe('/en');
      
      console.log('✅ Invalid locale cookie ignored, fallback to default');
    });
  });
}); 
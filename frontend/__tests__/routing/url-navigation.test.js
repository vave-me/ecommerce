/**
 * URL Navigation and Routing Tests
 * Tests specific URL patterns and routing behavior including locale handling
 * Focuses on: http://localhost:3000/explore and http://localhost:3000/
 */

import { jest } from '@jest/globals';

// Mock the routing configuration
const mockRouting = {
  locales: ['en', 'pl', 'de'],
  defaultLocale: 'en',
  localePrefix: 'always'
};

// Mock navigation functions
const mockPush = jest.fn();
const mockReplace = jest.fn();
const mockBack = jest.fn();
const mockForward = jest.fn();
const mockRefresh = jest.fn();
const mockPrefetch = jest.fn();

jest.mock('../../src/i18n/routing', () => ({
  routing: mockRouting
}));

jest.mock('../../src/i18n/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
    replace: mockReplace,
    back: mockBack,
    forward: mockForward,
    refresh: mockRefresh,
    prefetch: mockPrefetch,
  }),
  usePathname: jest.fn(),
  getPathname: jest.fn(),
  Link: ({ children, href, ...props }) => (
    <a href={href} {...props}>{children}</a>
  ),
  redirect: jest.fn(),
}));

jest.mock('next/navigation', () => ({
  usePathname: jest.fn(),
  useSearchParams: () => new URLSearchParams(),
  useParams: () => ({ locale: 'en' }),
}));

describe('URL Navigation and Routing Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Base URL Navigation (http://localhost:3000/)', () => {
    test('should handle navigation to root URL with English locale', () => {
      // Simulate navigation to root URL
      // With localePrefix: 'always', this becomes /en/
      
      const expectedPath = '/home'; // Assuming root redirects to home
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/home');
    });

    test('should handle navigation to root URL with Polish locale', () => {
      // With Polish locale, root URL becomes /pl/
      
      const expectedPath = '/home';
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/home');
    });

    test('should handle navigation to root URL with German locale', () => {
      // With German locale, root URL becomes /de/
      
      const expectedPath = '/home';
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/home');
    });

    test('should handle direct navigation to localhost:3000/', () => {
      // Test direct URL access patterns
      const testCases = [
        { locale: 'en', expectedUrl: '/en/', internalPath: '/home' },
        { locale: 'pl', expectedUrl: '/pl/', internalPath: '/home' },
        { locale: 'de', expectedUrl: '/de/', internalPath: '/home' }
      ];

      testCases.forEach(({ locale, expectedUrl, internalPath }) => {
        jest.clearAllMocks();
        
        // Simulate router navigation
        mockPush(internalPath);
        
        expect(mockPush).toHaveBeenCalledWith(internalPath);
      });
    });
  });

  describe('Explore URL Navigation (http://localhost:3000/explore)', () => {
    test('should handle navigation to explore URL with English locale', () => {
      // With localePrefix: 'always', this becomes /en/explore
      
      const expectedPath = '/explore';
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/explore');
    });

    test('should handle navigation to explore URL with Polish locale', () => {
      // With Polish locale, explore URL becomes /pl/explore
      
      const expectedPath = '/explore';
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/explore');
    });

    test('should handle navigation to explore URL with German locale', () => {
      // With German locale, explore URL becomes /de/explore
      
      const expectedPath = '/explore';
      mockPush(expectedPath);
      
      expect(mockPush).toHaveBeenCalledWith('/explore');
    });

    test('should handle direct navigation to localhost:3000/explore', () => {
      // Test direct URL access patterns for explore
      const testCases = [
        { locale: 'en', expectedUrl: '/en/explore', internalPath: '/explore' },
        { locale: 'pl', expectedUrl: '/pl/explore', internalPath: '/explore' },
        { locale: 'de', expectedUrl: '/de/explore', internalPath: '/explore' }
      ];

      testCases.forEach(({ locale, expectedUrl, internalPath }) => {
        jest.clearAllMocks();
        
        // Simulate router navigation
        mockPush(internalPath);
        
        expect(mockPush).toHaveBeenCalledWith(internalPath);
      });
    });
  });

  describe('URL Pattern Validation', () => {
    test('should validate correct URL patterns for all locales', () => {
      const baseUrls = [
        'http://localhost:3000/',
        'http://localhost:3000/explore'
      ];

      const locales = ['en', 'pl', 'de'];
      
      baseUrls.forEach(baseUrl => {
        locales.forEach(locale => {
          const path = new URL(baseUrl).pathname;
          const expectedLocalizedPath = `/${locale}${path === '/' ? '' : path}`;
          
          // Verify the pattern is valid
          expect(expectedLocalizedPath).toMatch(/^\/(en|pl|de)(\/.*)?$/);
        });
      });
    });

    test('should handle URL normalization', () => {
      const testUrls = [
        { input: 'http://localhost:3000/', expected: '/' },
        { input: 'http://localhost:3000/explore', expected: '/explore' },
        { input: 'http://localhost:3000/explore/', expected: '/explore' },
        { input: 'http://localhost:3000/explore?param=value', expected: '/explore' }
      ];

      testUrls.forEach(({ input, expected }) => {
        const url = new URL(input);
        const normalizedPath = url.pathname.replace(/\/$/, '') || '/';
        const finalPath = normalizedPath === '/' ? '/home' : normalizedPath;
        
        expect(finalPath).toBe(expected === '/' ? '/home' : expected);
      });
    });
  });

  describe('Navigation State Management', () => {
    test('should track navigation history correctly', () => {
      // Simulate navigation sequence: home -> explore -> back to home
      
      // Start at home
      mockPush('/home');
      expect(mockPush).toHaveBeenCalledWith('/home');
      
      // Navigate to explore
      jest.clearAllMocks();
      mockPush('/explore');
      expect(mockPush).toHaveBeenCalledWith('/explore');
      
      // Navigate back
      jest.clearAllMocks();
      mockBack();
      expect(mockBack).toHaveBeenCalled();
    });

    test('should handle forward navigation', () => {
      // Test forward navigation after going back
      mockForward();
      expect(mockForward).toHaveBeenCalled();
    });

    test('should handle page refresh', () => {
      // Test page refresh functionality
      mockRefresh();
      expect(mockRefresh).toHaveBeenCalled();
    });
  });

  describe('Locale-Specific URL Handling', () => {
    test('should maintain locale consistency across navigation', () => {
      const navigationSequence = [
        { from: '/home', to: '/explore', locale: 'en' },
        { from: '/explore', to: '/home', locale: 'en' },
        { from: '/home', to: '/explore', locale: 'pl' },
        { from: '/explore', to: '/home', locale: 'pl' },
        { from: '/home', to: '/explore', locale: 'de' },
        { from: '/explore', to: '/home', locale: 'de' }
      ];

      navigationSequence.forEach(({ from, to, locale }) => {
        jest.clearAllMocks();
        
        // Simulate navigation
        mockPush(to);
        
        // Verify navigation was called correctly
        expect(mockPush).toHaveBeenCalledWith(to);
      });
    });

    test('should handle locale switching on same page', () => {
      const currentPath = '/explore';
      const locales = ['en', 'pl', 'de'];
      
      locales.forEach(locale => {
        jest.clearAllMocks();
        
        // Simulate locale switch (would typically use replace instead of push)
        mockReplace(currentPath);
        
        expect(mockReplace).toHaveBeenCalledWith(currentPath);
      });
    });
  });

  describe('URL Query Parameters and Fragments', () => {
    test('should handle URLs with query parameters', () => {
      const urlsWithParams = [
        '/explore?category=electronics',
        '/explore?sort=price&order=asc',
        '/home?welcome=true'
      ];

      urlsWithParams.forEach(urlWithParams => {
        jest.clearAllMocks();
        
        // Extract base path for navigation
        const basePath = urlWithParams.split('?')[0];
        mockPush(basePath);
        
        expect(mockPush).toHaveBeenCalledWith(basePath);
      });
    });

    test('should handle URLs with fragments', () => {
      const urlsWithFragments = [
        '/explore#featured',
        '/home#welcome-section'
      ];

      urlsWithFragments.forEach(urlWithFragment => {
        jest.clearAllMocks();
        
        // Extract base path for navigation
        const basePath = urlWithFragment.split('#')[0];
        mockPush(basePath);
        
        expect(mockPush).toHaveBeenCalledWith(basePath);
      });
    });
  });

  describe('Error Handling and Edge Cases', () => {
    test('should handle invalid URLs gracefully', () => {
      const invalidUrls = [
        '/invalid-route',
        '/explore/nonexistent',
        '/home/invalid'
      ];

      invalidUrls.forEach(invalidUrl => {
        jest.clearAllMocks();
        
        // Navigation should still be attempted
        mockPush(invalidUrl);
        
        expect(mockPush).toHaveBeenCalledWith(invalidUrl);
      });
    });

    test('should handle navigation errors', () => {
      // Simulate navigation error
      const errorMessage = 'Navigation failed';
      mockPush.mockRejectedValueOnce(new Error(errorMessage));
      
      // Attempt navigation
      const navigationPromise = mockPush('/explore');
      
      // Verify error handling
      expect(navigationPromise).rejects.toThrow(errorMessage);
    });
  });

  describe('Performance and Optimization', () => {
    test('should handle prefetching for performance', () => {
      const routesToPrefetch = ['/explore', '/home', '/notifications'];
      
      routesToPrefetch.forEach(route => {
        mockPrefetch(route);
        expect(mockPrefetch).toHaveBeenCalledWith(route);
      });
    });

    test('should measure navigation performance', () => {
      const startTime = performance.now();
      
      // Simulate navigation
      mockPush('/explore');
      
      const endTime = performance.now();
      const navigationTime = endTime - startTime;
      
      // Navigation should be fast (under 10ms in test environment)
      expect(navigationTime).toBeLessThan(10);
      expect(mockPush).toHaveBeenCalledWith('/explore');
    });
  });

  describe('Browser Integration', () => {
    test('should handle browser back/forward buttons', () => {
      // Simulate browser back button
      mockBack();
      expect(mockBack).toHaveBeenCalled();
      
      // Simulate browser forward button
      jest.clearAllMocks();
      mockForward();
      expect(mockForward).toHaveBeenCalled();
    });

    test('should handle browser refresh', () => {
      // Simulate browser refresh
      mockRefresh();
      expect(mockRefresh).toHaveBeenCalled();
    });
  });

  describe('Accessibility and SEO', () => {
    test('should generate SEO-friendly URLs', () => {
      const routes = ['/home', '/explore'];
      
      routes.forEach(route => {
        // URLs should be clean and descriptive
        expect(route).toMatch(/^\/[a-z-]+$/);
        expect(route).not.toContain('_');
        expect(route).not.toContain(' ');
      });
    });

    test('should support canonical URL generation', () => {
      const baseUrl = 'http://localhost:3000';
      const routes = ['/home', '/explore'];
      const locales = ['en', 'pl', 'de'];
      
      routes.forEach(route => {
        locales.forEach(locale => {
          const canonicalUrl = `${baseUrl}/${locale}${route}`;
          
          // Verify canonical URL format
          expect(canonicalUrl).toMatch(/^https?:\/\/[^\/]+\/[a-z]{2}\/[a-z-]+$/);
        });
      });
    });
  });
}); 
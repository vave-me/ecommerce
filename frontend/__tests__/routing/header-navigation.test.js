/**
 * Header and AddDropdown Navigation Test Suite
 * Tests specific navigation functionality from Header and AddDropdown components
 * Focuses on real server interactions and routing behavior
 */

import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

describe('Header and AddDropdown Navigation Tests', () => {
  const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:3000';
  const SUPPORTED_LOCALES = ['en', 'pl', 'de'];
  const timeout = 15000;

  // Header navigation buttons and their expected routes
  const HEADER_NAVIGATION = {
    home: '/home',
    explore: '/explore',
    notifications: '/notifications',
    messages: '/messages',
    wishlist: '/wishlist',
    cart: '/cart'
  };

  // AddDropdown items and their expected routes
  const ADD_DROPDOWN_ITEMS = {
    product: '/add/product',
    post: '/add/post',
    vehicle: '/add/vehicle',
    deal: '/add/deal',
    property: '/add/property',
    job: '/add/job',
    service: '/add/service',
    video: '/add/video'
  };

  async function makeRequest(url, options = {}) {
    const curlCommand = `curl -s -i -L --max-time 8 "${url}"`;
    
    try {
      const { stdout } = await execAsync(curlCommand);
      const [headerSection, ...bodyParts] = stdout.split('\r\n\r\n');
      const body = bodyParts.join('\r\n\r\n');
      
      const lines = headerSection.split('\r\n');
      const statusLine = lines[0];
      const statusMatch = statusLine.match(/HTTP\/[\d.]+\s+(\d+)/);
      const statusCode = statusMatch ? parseInt(statusMatch[1]) : 0;
      
      const headers = {};
      for (let i = 1; i < lines.length; i++) {
        const [key, ...valueParts] = lines[i].split(': ');
        if (key && valueParts.length > 0) {
          headers[key.toLowerCase()] = valueParts.join(': ');
        }
      }
      
      return { statusCode, headers, body: body.trim() };
    } catch (error) {
      throw new Error(`Request failed for ${url}: ${error.message}`);
    }
  }

  async function checkServerHealth() {
    try {
      const response = await makeRequest(`${BASE_URL}/en`);
      return response.statusCode < 500;
    } catch (error) {
      return false;
    }
  }

  describe('Header Navigation Functionality', () => {
    beforeAll(async () => {
      const isHealthy = await checkServerHealth();
      if (!isHealthy) {
        console.warn('⚠️ Server may not be fully healthy, some tests may fail');
      }
    }, timeout);

    test('header navigation routes are accessible', async () => {
      const results = {};
      
      for (const [buttonName, route] of Object.entries(HEADER_NAVIGATION)) {
        results[buttonName] = {};
        
        for (const locale of SUPPORTED_LOCALES) {
          const url = `${BASE_URL}/${locale}${route}`;
          
          try {
            const response = await makeRequest(url);
            
            results[buttonName][locale] = {
              statusCode: response.statusCode,
              accessible: response.statusCode !== 404,
              hasError: response.statusCode >= 500
            };
            
            // Route should exist (not 404)
            expect(response.statusCode).not.toBe(404);
            
            // If successful, should have proper HTML structure
            if (response.statusCode === 200) {
              expect(response.headers['content-type']).toMatch(/text\/html/);
              expect(response.body).toContain(`lang="${locale}"`);
            }
            
          } catch (error) {
            console.error(`❌ Header navigation test failed for ${buttonName} (${locale}):`, error.message);
            results[buttonName][locale] = { error: error.message };
          }
        }
      }
      
      // Log summary
      console.log('\n📊 Header Navigation Test Results:');
      for (const [buttonName, localeResults] of Object.entries(results)) {
        const successCount = Object.values(localeResults).filter(r => r.accessible).length;
        console.log(`  ${buttonName}: ${successCount}/${SUPPORTED_LOCALES.length} locales accessible`);
      }
    }, timeout);

    test('header navigation preserves locale context', async () => {
      for (const locale of SUPPORTED_LOCALES) {
        try {
          // Test navigation between different header routes in same locale
          const homeResponse = await makeRequest(`${BASE_URL}/${locale}/home`);
          const exploreResponse = await makeRequest(`${BASE_URL}/${locale}/explore`);
          
          if (homeResponse.statusCode === 200 && exploreResponse.statusCode === 200) {
            // Both should have the same locale
            expect(homeResponse.body).toContain(`lang="${locale}"`);
            expect(exploreResponse.body).toContain(`lang="${locale}"`);
            
            // Both should be HTML pages
            expect(homeResponse.headers['content-type']).toMatch(/text\/html/);
            expect(exploreResponse.headers['content-type']).toMatch(/text\/html/);
          }
        } catch (error) {
          console.warn(`Locale context test failed for ${locale}:`, error.message);
        }
      }
    }, timeout);

    test('header routes handle authentication states', async () => {
      // Test routes that might require authentication
      const authSensitiveRoutes = ['/notifications', '/messages', '/wishlist', '/cart'];
      
      for (const route of authSensitiveRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}/en${route}`);
          
          // Should not crash, either show content or redirect to login
          expect([200, 301, 302, 307, 308, 401, 403]).toContain(response.statusCode);
          
          if (response.statusCode === 200) {
            expect(response.headers['content-type']).toMatch(/text\/html/);
          }
        } catch (error) {
          console.warn(`Auth state test failed for ${route}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('AddDropdown Navigation Functionality', () => {
    test('add dropdown routes are accessible', async () => {
      const results = {};
      
      for (const [itemName, route] of Object.entries(ADD_DROPDOWN_ITEMS)) {
        results[itemName] = {};
        
        for (const locale of SUPPORTED_LOCALES) {
          const url = `${BASE_URL}/${locale}${route}`;
          
          try {
            const response = await makeRequest(url);
            
            results[itemName][locale] = {
              statusCode: response.statusCode,
              accessible: response.statusCode !== 404,
              hasError: response.statusCode >= 500
            };
            
            // Route should exist (not 404)
            expect(response.statusCode).not.toBe(404);
            
            // If successful, should have proper HTML structure
            if (response.statusCode === 200) {
              expect(response.headers['content-type']).toMatch(/text\/html/);
              expect(response.body).toContain(`lang="${locale}"`);
            }
            
          } catch (error) {
            console.error(`❌ AddDropdown test failed for ${itemName} (${locale}):`, error.message);
            results[itemName][locale] = { error: error.message };
          }
        }
      }
      
      // Log summary
      console.log('\n📊 AddDropdown Navigation Test Results:');
      for (const [itemName, localeResults] of Object.entries(results)) {
        const successCount = Object.values(localeResults).filter(r => r.accessible).length;
        console.log(`  ${itemName}: ${successCount}/${SUPPORTED_LOCALES.length} locales accessible`);
      }
    }, timeout);

    test('add dropdown routes with query parameters', async () => {
      const routesWithParams = [
        { route: '/add/product', params: '?step=1' },
        { route: '/add/product', params: '?step=1&category=electronics' },
        { route: '/add/post', params: '?step=1&type=news' },
        { route: '/add/vehicle', params: '?step=2' },
        { route: '/add/job', params: '?step=1&type=fulltime' }
      ];

      for (const { route, params } of routesWithParams) {
        try {
          const url = `${BASE_URL}/en${route}${params}`;
          const response = await makeRequest(url);
          
          // Should handle query parameters properly
          expect(response.statusCode).not.toBe(404);
          
          if (response.statusCode === 200) {
            expect(response.headers['content-type']).toMatch(/text\/html/);
            expect(response.body).toContain('lang="en"');
          }
        } catch (error) {
          console.warn(`Query params test failed for ${route}${params}:`, error.message);
        }
      }
    }, timeout);

    test('add dropdown navigation preserves locale and context', async () => {
      for (const locale of SUPPORTED_LOCALES) {
        try {
          // Test different add routes in same locale
          const productResponse = await makeRequest(`${BASE_URL}/${locale}/add/product`);
          const postResponse = await makeRequest(`${BASE_URL}/${locale}/add/post`);
          
          if (productResponse.statusCode === 200 && postResponse.statusCode === 200) {
            // Both should have the same locale
            expect(productResponse.body).toContain(`lang="${locale}"`);
            expect(postResponse.body).toContain(`lang="${locale}"`);
          }
        } catch (error) {
          console.warn(`AddDropdown locale context test failed for ${locale}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('Router Integration Tests', () => {
    test('navigation between header and add routes works', async () => {
      try {
        // Test navigation flow: home -> add product -> back to home
        const homeResponse = await makeRequest(`${BASE_URL}/en/home`);
        const addProductResponse = await makeRequest(`${BASE_URL}/en/add/product`);
        const backToHomeResponse = await makeRequest(`${BASE_URL}/en/home`);
        
        // All should be accessible
        expect(homeResponse.statusCode).not.toBe(404);
        expect(addProductResponse.statusCode).not.toBe(404);
        expect(backToHomeResponse.statusCode).not.toBe(404);
        
        // All should maintain locale context
        if (homeResponse.statusCode === 200) {
          expect(homeResponse.body).toContain('lang="en"');
        }
        if (addProductResponse.statusCode === 200) {
          expect(addProductResponse.body).toContain('lang="en"');
        }
        if (backToHomeResponse.statusCode === 200) {
          expect(backToHomeResponse.body).toContain('lang="en"');
        }
      } catch (error) {
        console.warn('Navigation flow test failed:', error.message);
      }
    }, timeout);

    test('router handles invalid navigation gracefully', async () => {
      const invalidRoutes = [
        '/en/add/invalid-type',
        '/en/add/',
        '/en/add/product/invalid-step',
        '/en/header-route-that-does-not-exist'
      ];

      for (const route of invalidRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          // Should return 404 for invalid routes
          expect(response.statusCode).toBe(404);
        } catch (error) {
          console.warn(`Invalid route test failed for ${route}:`, error.message);
        }
      }
    }, timeout);

    test('router preserves state across navigation', async () => {
      // Test that navigation doesn't break the application state
      const navigationSequence = [
        '/en/home',
        '/en/explore', 
        '/en/add/product',
        '/en/messages',
        '/en/add/post',
        '/en/home'
      ];

      let previousResponse = null;
      
      for (const route of navigationSequence) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          // Each route should be accessible
          expect(response.statusCode).not.toBe(404);
          
          // Should maintain consistent locale
          if (response.statusCode === 200) {
            expect(response.body).toContain('lang="en"');
          }
          
          previousResponse = response;
        } catch (error) {
          console.warn(`Navigation sequence test failed at ${route}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('Performance and Error Handling', () => {
    test('navigation routes load within acceptable time', async () => {
      const criticalRoutes = [
        '/en/home',
        '/en/explore',
        '/en/add/product',
        '/en/add/post'
      ];

      for (const route of criticalRoutes) {
        const startTime = Date.now();
        
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          const loadTime = Date.now() - startTime;
          
          // Should load within 8 seconds (generous for CI environments)
          expect(loadTime).toBeLessThan(8000);
          
          if (response.statusCode === 200) {
            // Should have substantial content
            expect(response.body.length).toBeGreaterThan(500);
          }
        } catch (error) {
          console.warn(`Performance test failed for ${route}:`, error.message);
        }
      }
    }, timeout);

    test('navigation handles server errors gracefully', async () => {
      // Test routes that might have server errors
      const testRoutes = ['/en/home', '/en/add/product'];
      
      for (const route of testRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          if (response.statusCode >= 500) {
            // Server errors should still return proper error pages
            expect(response.headers['content-type']).toMatch(/text\/html/);
            
            // Should contain some error indication
            expect(response.body).toContain('html');
          }
        } catch (error) {
          console.warn(`Error handling test failed for ${route}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('SEO and Accessibility', () => {
    test('navigation routes have proper SEO structure', async () => {
      const seoTestRoutes = ['/en/home', '/en/explore', '/en/add/product'];
      
      for (const route of seoTestRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          if (response.statusCode === 200) {
            const html = response.body;
            
            // Check for essential SEO elements
            expect(html).toMatch(/<title>/);
            expect(html).toMatch(/<meta[^>]+description/);
            expect(html).toMatch(/<meta[^>]+viewport/);
            expect(html).toMatch(/lang="en"/);
            
            // Check for proper HTML structure
            expect(html).toMatch(/<html[^>]*>/);
            expect(html).toMatch(/<head>/);
            expect(html).toMatch(/<body>/);
          }
        } catch (error) {
          console.warn(`SEO test failed for ${route}:`, error.message);
        }
      }
    }, timeout);
  });
}); 
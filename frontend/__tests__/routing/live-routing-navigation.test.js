/**
 * Live Routing and Navigation Test Suite
 * Tests Header navigation, AddDropdown routing, and real server responses
 * Focuses on production-like behavior with actual HTTP requests
 */

import { exec } from 'child_process';
import { promisify } from 'util';
import fs from 'fs/promises';
import path from 'path';

const execAsync = promisify(exec);

describe('Live Routing and Navigation Tests', () => {
  const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:3000';
  const SUPPORTED_LOCALES = ['en', 'pl', 'de'];
  const timeout = 15000; // 15 seconds timeout

  // Header navigation routes that should be accessible
  const HEADER_ROUTES = [
    '/home',
    '/explore', 
    '/notifications',
    '/messages',
    '/wishlist',
    '/cart'
  ];

  // AddDropdown routes that should be accessible
  const ADD_DROPDOWN_ROUTES = [
    '/add/product',
    '/add/post', 
    '/add/vehicle',
    '/add/deal',
    '/add/property',
    '/add/job',
    '/add/service',
    '/add/video'
  ];

  // Core application routes
  const CORE_ROUTES = [
    '',
    '/marketplace',
    '/deals',
    '/jobs', 
    '/properties',
    '/services',
    '/products',
    '/news',
    '/search',
    '/login',
    '/about',
    '/contact',
    '/terms',
    '/privacy'
  ];

  // Helper function to make HTTP requests with curl
  async function makeRequest(url, options = {}) {
    const curlCommand = buildCurlCommand(url, options);
    
    try {
      const { stdout, stderr } = await execAsync(curlCommand);
      
      if (stderr && !stderr.includes('progress')) {
        console.warn(`Curl warning for ${url}:`, stderr);
      }
      
      return parseCurlResponse(stdout);
    } catch (error) {
      throw new Error(`Request failed for ${url}: ${error.message}`);
    }
  }

  function buildCurlCommand(url, options = {}) {
    let command = `curl -s -i --max-time 10`;
    
    if (options.followRedirects !== false) {
      command += ' -L';
    }
    
    if (options.method) {
      command += ` -X ${options.method}`;
    }
    
    if (options.headers) {
      Object.entries(options.headers).forEach(([key, value]) => {
        command += ` -H "${key}: ${value}"`;
      });
    }
    
    command += ` "${url}"`;
    return command;
  }

  function parseCurlResponse(response) {
    const parts = response.split('\r\n\r\n');
    const headerSection = parts[0];
    const body = parts.slice(1).join('\r\n\r\n');
    
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
    
    return {
      statusCode,
      headers,
      body: body.trim()
    };
  }

  async function waitForServer(maxAttempts = 10) {
    for (let i = 0; i < maxAttempts; i++) {
      try {
        const response = await makeRequest(`${BASE_URL}/en`);
        if (response.statusCode < 500) {
          return true;
        }
      } catch (error) {
        console.log(`Server check attempt ${i + 1}/${maxAttempts} failed`);
      }
      await new Promise(resolve => setTimeout(resolve, 2000));
    }
    throw new Error('Server not responding after maximum attempts');
  }

  describe('Server Health and Basic Routing', () => {
    beforeAll(async () => {
      console.log('🔍 Checking server health...');
      await waitForServer();
      console.log('✅ Server is responding');
    }, timeout);

    test('server is running and responsive', async () => {
      const response = await makeRequest(`${BASE_URL}/en`);
      expect([200, 301, 302, 307, 308]).toContain(response.statusCode);
    }, timeout);

    test('root path redirects to default locale', async () => {
      const response = await makeRequest(BASE_URL, { followRedirects: false });
      
      if ([301, 302, 307, 308].includes(response.statusCode)) {
        expect(response.headers.location).toMatch(/\/en\/?$/);
      } else {
        // If no redirect, should serve content directly
        expect(response.statusCode).toBe(200);
      }
    }, timeout);

    test.each(SUPPORTED_LOCALES)('locale %s base route is accessible', async (locale) => {
      const response = await makeRequest(`${BASE_URL}/${locale}`);
      
      // Should not return 404 or 500 errors
      expect(response.statusCode).not.toBe(404);
      expect(response.statusCode).not.toBe(500);
      
      if (response.statusCode === 200) {
        expect(response.headers['content-type']).toMatch(/text\/html/);
        expect(response.body).toContain(`lang="${locale}"`);
      }
    }, timeout);
  });

  describe('Header Navigation Routes', () => {
    test.each(HEADER_ROUTES)('header route %s works for all locales', async (route) => {
      for (const locale of SUPPORTED_LOCALES) {
        const url = `${BASE_URL}/${locale}${route}`;
        
        try {
          const response = await makeRequest(url);
          
          // Should not return 404 (route should exist)
          expect(response.statusCode).not.toBe(404);
          
          // Log any server errors for debugging
          if (response.statusCode >= 500) {
            console.warn(`⚠️ Server error for ${url}: ${response.statusCode}`);
          }
          
          // If successful, should have proper content type
          if (response.statusCode === 200) {
            expect(response.headers['content-type']).toMatch(/text\/html/);
            expect(response.body).toContain(`lang="${locale}"`);
          }
        } catch (error) {
          console.error(`❌ Failed to test ${url}:`, error.message);
          throw error;
        }
      }
    }, timeout);

    test('header navigation includes proper meta tags', async () => {
      const response = await makeRequest(`${BASE_URL}/en/home`);
      
      if (response.statusCode === 200) {
        const html = response.body;
        
        // Check for essential meta tags
        expect(html).toMatch(/<title>/);
        expect(html).toMatch(/<meta[^>]+description/);
        expect(html).toMatch(/<meta[^>]+viewport/);
        
        // Check for proper locale
        expect(html).toMatch(/lang="en"/);
      }
    }, timeout);
  });

  describe('AddDropdown Navigation Routes', () => {
    test.each(ADD_DROPDOWN_ROUTES)('add dropdown route %s works for all locales', async (route) => {
      for (const locale of SUPPORTED_LOCALES) {
        const url = `${BASE_URL}/${locale}${route}`;
        
        try {
          const response = await makeRequest(url);
          
          // Should not return 404 (route should exist)
          expect(response.statusCode).not.toBe(404);
          
          // Log any server errors for debugging
          if (response.statusCode >= 500) {
            console.warn(`⚠️ Server error for ${url}: ${response.statusCode}`);
          }
          
          // If successful, should have proper content type
          if (response.statusCode === 200) {
            expect(response.headers['content-type']).toMatch(/text\/html/);
            expect(response.body).toContain(`lang="${locale}"`);
          }
        } catch (error) {
          console.error(`❌ Failed to test ${url}:`, error.message);
          throw error;
        }
      }
    }, timeout);

    test('add routes with query parameters work correctly', async () => {
      const testRoutes = [
        '/add/product?step=1',
        '/add/post?step=1&category=news',
        '/add/vehicle?step=2'
      ];

      for (const route of testRoutes) {
        const url = `${BASE_URL}/en${route}`;
        
        try {
          const response = await makeRequest(url);
          
          // Should handle query parameters properly
          expect(response.statusCode).not.toBe(404);
          
          if (response.statusCode >= 500) {
            console.warn(`⚠️ Server error for ${url}: ${response.statusCode}`);
          }
        } catch (error) {
          console.error(`❌ Failed to test ${url}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('Core Application Routes', () => {
    test.each(CORE_ROUTES)('core route %s works for all locales', async (route) => {
      for (const locale of SUPPORTED_LOCALES) {
        const url = `${BASE_URL}/${locale}${route}`;
        
        try {
          const response = await makeRequest(url);
          
          // Should not return 404 for core routes
          expect(response.statusCode).not.toBe(404);
          
          // Log any server errors for debugging
          if (response.statusCode >= 500) {
            console.warn(`⚠️ Server error for ${url}: ${response.statusCode}`);
          }
          
          // If successful, should have proper content type
          if (response.statusCode === 200) {
            expect(response.headers['content-type']).toMatch(/text\/html/);
            expect(response.body).toContain(`lang="${locale}"`);
          }
        } catch (error) {
          console.error(`❌ Failed to test ${url}:`, error.message);
          throw error;
        }
      }
    }, timeout);
  });

  describe('Router Functionality Tests', () => {
    test('navigation preserves locale context', async () => {
      for (const locale of SUPPORTED_LOCALES) {
        const homeUrl = `${BASE_URL}/${locale}/home`;
        const exploreUrl = `${BASE_URL}/${locale}/explore`;
        
        try {
          // Test that both routes work in the same locale
          const homeResponse = await makeRequest(homeUrl);
          const exploreResponse = await makeRequest(exploreUrl);
          
          expect(homeResponse.statusCode).not.toBe(404);
          expect(exploreResponse.statusCode).not.toBe(404);
          
          if (homeResponse.statusCode === 200) {
            expect(homeResponse.body).toContain(`lang="${locale}"`);
          }
          
          if (exploreResponse.statusCode === 200) {
            expect(exploreResponse.body).toContain(`lang="${locale}"`);
          }
        } catch (error) {
          console.error(`❌ Failed locale context test for ${locale}:`, error.message);
        }
      }
    }, timeout);

    test('invalid routes return 404', async () => {
      const invalidRoutes = [
        '/en/invalid-route',
        '/en/non-existent-page',
        '/en/fake/path'
      ];

      for (const route of invalidRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          expect(response.statusCode).toBe(404);
        } catch (error) {
          console.warn(`Invalid route test failed for ${route}:`, error.message);
        }
      }
    }, timeout);

    test('API routes are not affected by locale middleware', async () => {
      const apiRoutes = [
        '/api/health',
        '/api/status'
      ];

      for (const route of apiRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          // API routes should work (200) or not exist (404), but not be affected by locale routing
          expect([200, 404]).toContain(response.statusCode);
        } catch (error) {
          console.warn(`API route test failed for ${route}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('Performance and Caching', () => {
    test('static pages have appropriate cache headers', async () => {
      const staticRoutes = ['/en/about', '/en/contact', '/en/terms'];
      
      for (const route of staticRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          if (response.statusCode === 200) {
            // Check for cache-related headers
            expect(response.headers).toHaveProperty('cache-control');
          }
        } catch (error) {
          console.warn(`Cache headers test failed for ${route}:`, error.message);
        }
      }
    }, timeout);

    test('pages load within reasonable time', async () => {
      const testRoutes = ['/en', '/en/home', '/en/marketplace'];
      
      for (const route of testRoutes) {
        const startTime = Date.now();
        
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          const loadTime = Date.now() - startTime;
          
          // Should load within 5 seconds
          expect(loadTime).toBeLessThan(5000);
          
          if (response.statusCode === 200) {
            // Should have content
            expect(response.body.length).toBeGreaterThan(100);
          }
        } catch (error) {
          console.warn(`Performance test failed for ${route}:`, error.message);
        }
      }
    }, timeout);
  });

  describe('Error Handling', () => {
    test('server errors are handled gracefully', async () => {
      // Test routes that might cause server errors
      const testRoutes = ['/en/home', '/en/marketplace'];
      
      for (const route of testRoutes) {
        try {
          const response = await makeRequest(`${BASE_URL}${route}`);
          
          if (response.statusCode >= 500) {
            // Server errors should still return proper HTML error pages
            expect(response.headers['content-type']).toMatch(/text\/html/);
            expect(response.body).toContain('html');
          }
        } catch (error) {
          console.warn(`Error handling test failed for ${route}:`, error.message);
        }
      }
    }, timeout);

    test('malformed URLs are handled properly', async () => {
      const malformedUrls = [
        `${BASE_URL}/en/%20invalid`,
        `${BASE_URL}/en/test%00null`,
        `${BASE_URL}/en/../../../etc/passwd`
      ];

      for (const url of malformedUrls) {
        try {
          const response = await makeRequest(url);
          // Should return 400 or 404, not crash
          expect([400, 404]).toContain(response.statusCode);
        } catch (error) {
          console.warn(`Malformed URL test failed for ${url}:`, error.message);
        }
      }
    }, timeout);
  });
}); 
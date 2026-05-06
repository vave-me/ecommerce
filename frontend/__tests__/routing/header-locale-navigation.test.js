/**
 * Header and AddDropdown Locale Navigation Test Suite
 * Tests locale preservation in Header and AddDropdown navigation
 * Identifies why users are redirected without locale prefix
 */

import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

describe('Header and AddDropdown Locale Navigation Tests', () => {
  const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:3000';
  const SUPPORTED_LOCALES = ['en', 'pl', 'de'];
  const timeout = 20000;

  // Header navigation routes that should preserve locale
  const HEADER_ROUTES = {
    home: '/home',
    explore: '/explore',
    notifications: '/notifications',
    messages: '/messages',
    wishlist: '/wishlist',
    cart: '/cart'
  };

  // AddDropdown routes that should preserve locale
  const ADD_DROPDOWN_ROUTES = {
    product: '/add/product',
    post: '/add/post',
    vehicle: '/add/vehicle',
    deal: '/add/deal',
    property: '/add/property',
    job: '/add/job',
    service: '/add/service'
  };

  async function makeRequest(url, options = {}) {
    const curlCommand = `curl -s -i -L --max-time 10 "${url}"`;
    
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

  async function makeRequestWithRedirectTracking(url) {
    const curlCommand = `curl -s -i --max-time 10 "${url}"`;
    
    try {
      const { stdout } = await execAsync(curlCommand);
      const responses = stdout.split(/(?=HTTP\/)/);
      
      const redirectChain = [];
      
      for (const response of responses) {
        if (!response.trim()) continue;
        
        const [headerSection, ...bodyParts] = response.split('\r\n\r\n');
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
        
        redirectChain.push({
          statusCode,
          headers,
          body: body.trim(),
          location: headers.location || null
        });
      }
      
      return {
        finalResponse: redirectChain[redirectChain.length - 1],
        redirectChain
      };
    } catch (error) {
      throw new Error(`Request failed for ${url}: ${error.message}`);
    }
  }

  describe('Server Health Check', () => {
    test('server is responding correctly', async () => {
      const response = await makeRequest(`${BASE_URL}/en`);
      expect([200, 301, 302, 307, 308]).toContain(response.statusCode);
      console.log(`✅ Server responding with status: ${response.statusCode}`);
    }, timeout);
  });

  describe('Header Navigation Locale Preservation', () => {
    test.each(Object.entries(HEADER_ROUTES))('header route %s preserves locale for all locales', async (routeName, route) => {
      console.log(`\n🔍 Testing Header route: ${routeName} (${route})`);
      
      for (const locale of SUPPORTED_LOCALES) {
        const expectedUrl = `${BASE_URL}/${locale}${route}`;
        
        try {
          const { finalResponse, redirectChain } = await makeRequestWithRedirectTracking(expectedUrl);
          
          console.log(`  📍 ${locale}: ${expectedUrl}`);
          console.log(`    Status: ${finalResponse.statusCode}`);
          
          if (redirectChain.length > 1) {
            console.log(`    Redirects: ${redirectChain.length - 1}`);
            redirectChain.forEach((resp, index) => {
              if (resp.location) {
                console.log(`      ${index + 1}. ${resp.statusCode} -> ${resp.location}`);
              }
            });
          }
          
          // Check if route is accessible
          expect(finalResponse.statusCode).not.toBe(404);
          
          // Check if locale is preserved in final response
          if (finalResponse.statusCode === 200) {
            expect(finalResponse.body).toContain(`lang="${locale}"`);
            console.log(`    ✅ Locale preserved in HTML`);
          }
          
          // Check if any redirects lost the locale
          for (const response of redirectChain) {
            if (response.location) {
              const redirectUrl = new URL(response.location, BASE_URL);
              const pathSegments = redirectUrl.pathname.split('/').filter(Boolean);
              
              if (pathSegments.length > 0 && !SUPPORTED_LOCALES.includes(pathSegments[0])) {
                console.log(`    ❌ LOCALE LOST in redirect: ${response.location}`);
                expect(pathSegments[0]).toBeOneOf(SUPPORTED_LOCALES);
              }
            }
          }
          
        } catch (error) {
          console.error(`    ❌ Failed: ${error.message}`);
          throw error;
        }
      }
    }, timeout);

    test('header navigation from different starting locales', async () => {
      console.log(`\n🔍 Testing Header navigation from different starting locales`);
      
      for (const startLocale of SUPPORTED_LOCALES) {
        console.log(`  Starting from locale: ${startLocale}`);
        
        // Test navigation to different routes from this locale
        for (const [routeName, route] of Object.entries(HEADER_ROUTES)) {
          const url = `${BASE_URL}/${startLocale}${route}`;
          
          try {
            const response = await makeRequest(url);
            
            if (response.statusCode === 200) {
              // Check if the response maintains the correct locale
              expect(response.body).toContain(`lang="${startLocale}"`);
              console.log(`    ✅ ${routeName}: locale ${startLocale} preserved`);
            } else {
              console.log(`    ⚠️ ${routeName}: status ${response.statusCode}`);
            }
          } catch (error) {
            console.log(`    ❌ ${routeName}: ${error.message}`);
          }
        }
      }
    }, timeout);
  });

  describe('AddDropdown Navigation Locale Preservation', () => {
    test.each(Object.entries(ADD_DROPDOWN_ROUTES))('add dropdown route %s preserves locale for all locales', async (routeName, route) => {
      console.log(`\n🔍 Testing AddDropdown route: ${routeName} (${route})`);
      
      for (const locale of SUPPORTED_LOCALES) {
        const expectedUrl = `${BASE_URL}/${locale}${route}`;
        
        try {
          const { finalResponse, redirectChain } = await makeRequestWithRedirectTracking(expectedUrl);
          
          console.log(`  📍 ${locale}: ${expectedUrl}`);
          console.log(`    Status: ${finalResponse.statusCode}`);
          
          if (redirectChain.length > 1) {
            console.log(`    Redirects: ${redirectChain.length - 1}`);
            redirectChain.forEach((resp, index) => {
              if (resp.location) {
                console.log(`      ${index + 1}. ${resp.statusCode} -> ${resp.location}`);
                
                // Check if redirect URL loses locale
                const redirectUrl = new URL(resp.location, BASE_URL);
                const pathSegments = redirectUrl.pathname.split('/').filter(Boolean);
                
                if (pathSegments.length > 0 && !SUPPORTED_LOCALES.includes(pathSegments[0])) {
                  console.log(`      ❌ LOCALE LOST: Expected /${locale}${route}, got ${resp.location}`);
                }
              }
            });
          }
          
          // Check if route is accessible
          expect(finalResponse.statusCode).not.toBe(404);
          
          // Check if locale is preserved in final response
          if (finalResponse.statusCode === 200) {
            expect(finalResponse.body).toContain(`lang="${locale}"`);
            console.log(`    ✅ Locale preserved in HTML`);
          }
          
        } catch (error) {
          console.error(`    ❌ Failed: ${error.message}`);
          throw error;
        }
      }
    }, timeout);

    test('add dropdown navigation simulating user clicks', async () => {
      console.log(`\n🔍 Simulating AddDropdown user clicks`);
      
      for (const locale of SUPPORTED_LOCALES) {
        console.log(`  Testing from locale: ${locale}`);
        
        // Simulate user starting from a localized page
        const startingPage = `${BASE_URL}/${locale}/home`;
        const startResponse = await makeRequest(startingPage);
        
        if (startResponse.statusCode === 200) {
          console.log(`    ✅ Starting page loaded: ${startingPage}`);
          
          // Now test each AddDropdown route as if user clicked
          for (const [routeName, route] of Object.entries(ADD_DROPDOWN_ROUTES)) {
            const targetUrl = `${BASE_URL}/${locale}${route}`;
            
            try {
              const { finalResponse, redirectChain } = await makeRequestWithRedirectTracking(targetUrl);
              
              // Check for locale loss in redirects
              let localeLost = false;
              for (const response of redirectChain) {
                if (response.location) {
                  const redirectUrl = new URL(response.location, BASE_URL);
                  const pathSegments = redirectUrl.pathname.split('/').filter(Boolean);
                  
                  if (pathSegments.length > 0 && !SUPPORTED_LOCALES.includes(pathSegments[0])) {
                    localeLost = true;
                    console.log(`    ❌ ${routeName}: Locale lost in redirect ${response.location}`);
                  }
                }
              }
              
              if (!localeLost && finalResponse.statusCode === 200) {
                expect(finalResponse.body).toContain(`lang="${locale}"`);
                console.log(`    ✅ ${routeName}: Locale preserved`);
              }
              
            } catch (error) {
              console.log(`    ❌ ${routeName}: ${error.message}`);
            }
          }
        }
      }
    }, timeout);
  });

  describe('Problematic Routes Investigation', () => {
    test('investigate routes that lose locale', async () => {
      console.log(`\n🔍 Investigating problematic routes that lose locale`);
      
      const problematicRoutes = [
        '/add/deal',
        '/add/post', 
        '/explore',
        '/home'
      ];
      
      for (const route of problematicRoutes) {
        console.log(`\n  Testing route: ${route}`);
        
        // Test direct access without locale
        try {
          const directResponse = await makeRequestWithRedirectTracking(`${BASE_URL}${route}`);
          console.log(`    Direct access (${BASE_URL}${route}):`);
          console.log(`      Status: ${directResponse.finalResponse.statusCode}`);
          
          if (directResponse.redirectChain.length > 1) {
            directResponse.redirectChain.forEach((resp, index) => {
              if (resp.location) {
                console.log(`      Redirect ${index + 1}: ${resp.statusCode} -> ${resp.location}`);
              }
            });
          }
        } catch (error) {
          console.log(`    Direct access failed: ${error.message}`);
        }
        
        // Test with each locale
        for (const locale of SUPPORTED_LOCALES) {
          try {
            const localizedResponse = await makeRequestWithRedirectTracking(`${BASE_URL}/${locale}${route}`);
            console.log(`    With locale ${locale}:`);
            console.log(`      Status: ${localizedResponse.finalResponse.statusCode}`);
            
            // Check for unexpected redirects
            if (localizedResponse.redirectChain.length > 1) {
              localizedResponse.redirectChain.forEach((resp, index) => {
                if (resp.location) {
                  console.log(`      Redirect ${index + 1}: ${resp.statusCode} -> ${resp.location}`);
                  
                  // Check if redirect loses locale
                  const redirectUrl = new URL(resp.location, BASE_URL);
                  const pathSegments = redirectUrl.pathname.split('/').filter(Boolean);
                  
                  if (pathSegments.length > 0 && !SUPPORTED_LOCALES.includes(pathSegments[0])) {
                    console.log(`      ❌ LOCALE LOST in redirect!`);
                  }
                }
              });
            }
          } catch (error) {
            console.log(`    Localized access failed: ${error.message}`);
          }
        }
      }
    }, timeout);

    test('test router.push() behavior simulation', async () => {
      console.log(`\n🔍 Testing router.push() behavior simulation`);
      
      // Simulate what happens when Header/AddDropdown calls router.push()
      const routerPushSimulations = [
        { from: '/en/home', to: '/explore', expected: '/en/explore' },
        { from: '/en/home', to: '/add/deal', expected: '/en/add/deal' },
        { from: '/pl/marketplace', to: '/add/post', expected: '/pl/add/post' },
        { from: '/de/search', to: '/home', expected: '/de/home' }
      ];
      
      for (const simulation of routerPushSimulations) {
        console.log(`\n  Simulating: From ${simulation.from} -> router.push('${simulation.to}')`);
        console.log(`  Expected result: ${simulation.expected}`);
        
        // Test the expected URL
        try {
          const response = await makeRequestWithRedirectTracking(`${BASE_URL}${simulation.expected}`);
          console.log(`    Expected URL status: ${response.finalResponse.statusCode}`);
          
          if (response.finalResponse.statusCode === 200) {
            const locale = simulation.expected.split('/')[1];
            if (SUPPORTED_LOCALES.includes(locale)) {
              expect(response.finalResponse.body).toContain(`lang="${locale}"`);
              console.log(`    ✅ Locale ${locale} preserved correctly`);
            }
          }
        } catch (error) {
          console.log(`    ❌ Expected URL failed: ${error.message}`);
        }
        
        // Test what might happen if locale is lost
        const withoutLocale = simulation.to;
        try {
          const badResponse = await makeRequestWithRedirectTracking(`${BASE_URL}${withoutLocale}`);
          console.log(`    Without locale (${withoutLocale}) status: ${badResponse.finalResponse.statusCode}`);
          
          if (badResponse.redirectChain.length > 1) {
            console.log(`    Redirects when locale missing:`);
            badResponse.redirectChain.forEach((resp, index) => {
              if (resp.location) {
                console.log(`      ${index + 1}. ${resp.statusCode} -> ${resp.location}`);
              }
            });
          }
        } catch (error) {
          console.log(`    Without locale test failed: ${error.message}`);
        }
      }
    }, timeout);
  });

  describe('Middleware and Routing Behavior', () => {
    test('test middleware locale handling', async () => {
      console.log(`\n🔍 Testing middleware locale handling`);
      
      const testCases = [
        { url: '/', description: 'Root path' },
        { url: '/home', description: 'Route without locale' },
        { url: '/en', description: 'Locale only' },
        { url: '/en/home', description: 'Locale with route' },
        { url: '/invalid-locale/home', description: 'Invalid locale' }
      ];
      
      for (const testCase of testCases) {
        console.log(`\n  Testing: ${testCase.description} (${testCase.url})`);
        
        try {
          const response = await makeRequestWithRedirectTracking(`${BASE_URL}${testCase.url}`);
          console.log(`    Final status: ${response.finalResponse.statusCode}`);
          
          if (response.redirectChain.length > 1) {
            console.log(`    Redirect chain:`);
            response.redirectChain.forEach((resp, index) => {
              if (resp.location) {
                console.log(`      ${index + 1}. ${resp.statusCode} -> ${resp.location}`);
              }
            });
          }
          
          // Analyze final URL structure
          if (response.finalResponse.statusCode === 200) {
            const hasLangAttribute = response.finalResponse.body.includes('lang=');
            if (hasLangAttribute) {
              const langMatch = response.finalResponse.body.match(/lang="([^"]+)"/);
              if (langMatch) {
                console.log(`    Final page locale: ${langMatch[1]}`);
              }
            }
          }
        } catch (error) {
          console.log(`    Failed: ${error.message}`);
        }
      }
    }, timeout);
  });
});

// Helper function for Jest custom matcher
expect.extend({
  toBeOneOf(received, validOptions) {
    const pass = validOptions.includes(received);
    if (pass) {
      return {
        message: () => `expected ${received} not to be one of ${validOptions.join(', ')}`,
        pass: true,
      };
    } else {
      return {
        message: () => `expected ${received} to be one of ${validOptions.join(', ')}`,
        pass: false,
      };
    }
  },
}); 
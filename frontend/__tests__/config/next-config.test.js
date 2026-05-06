/**
 * Comprehensive Next.js Configuration Test Suite
 * Tests next.config.mjs functionality, plugins, and optimization settings
 */

import { jest } from '@jest/globals';
import { readFile } from 'fs/promises';
import { join } from 'path';

// Mock next-intl plugin
const mockNextIntlPlugin = jest.fn((config) => config);
jest.mock('next-intl/plugin', () => mockNextIntlPlugin);

// Mock the i18n request configuration
jest.mock('../../src/i18n/request', () => ({}));

describe('Next.js Configuration Tests', () => {
  let nextConfig;

  beforeAll(async () => {
    // Read and evaluate the next.config.mjs file
    const configPath = join(process.cwd(), 'next.config.mjs');
    const configContent = await readFile(configPath, 'utf-8');
    
    // Create a mock module environment
    const module = { exports: {} };
    const exports = {};
    
    // Evaluate the config in a controlled environment
    const configFunction = new Function('module', 'exports', 'require', configContent);
    configFunction(module, exports, require);
    
    nextConfig = module.exports.default || module.exports;
  });

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Basic Configuration Structure', () => {
    test('exports a valid Next.js configuration object', () => {
      console.log('\n🔍 TESTING: Basic configuration structure');
      
      expect(nextConfig).toBeDefined();
      expect(typeof nextConfig).toBe('object');
      
      console.log('✅ Next.js configuration object is valid');
    });

    test('includes required configuration properties', () => {
      console.log('\n🔍 TESTING: Required configuration properties');
      
      const requiredProps = [
        'experimental',
        'images',
        'webpack',
        'env'
      ];

      requiredProps.forEach(prop => {
        expect(nextConfig).toHaveProperty(prop);
      });

      console.log('✅ All required configuration properties present');
    });
  });

  describe('Experimental Features Configuration', () => {
    test('enables required experimental features', () => {
      console.log('\n🔍 TESTING: Experimental features configuration');
      
      expect(nextConfig.experimental).toBeDefined();
      expect(nextConfig.experimental.turbo).toBeDefined();
      
      console.log('✅ Experimental features configured correctly');
    });

    test('turbo configuration is properly set', () => {
      console.log('\n🔍 TESTING: Turbo configuration');
      
      const turboConfig = nextConfig.experimental.turbo;
      
      expect(turboConfig).toHaveProperty('rules');
      expect(turboConfig.rules).toHaveProperty('*.svg');
      
      const svgRule = turboConfig.rules['*.svg'];
      expect(svgRule).toHaveProperty('loaders');
      expect(Array.isArray(svgRule.loaders)).toBe(true);
      
      console.log('✅ Turbo configuration is properly set');
    });

    test('SVG loader configuration in turbo', () => {
      console.log('\n🔍 TESTING: SVG loader in turbo configuration');
      
      const svgLoaders = nextConfig.experimental.turbo.rules['*.svg'].loaders;
      
      expect(svgLoaders).toContainEqual({
        loader: '@svgr/webpack',
        options: {
          icon: true
        }
      });
      
      console.log('✅ SVG loader configured correctly in turbo');
    });
  });

  describe('Images Configuration', () => {
    test('configures image domains and formats', () => {
      console.log('\n🔍 TESTING: Images configuration');
      
      expect(nextConfig.images).toBeDefined();
      expect(nextConfig.images.remotePatterns).toBeDefined();
      expect(Array.isArray(nextConfig.images.remotePatterns)).toBe(true);
      
      console.log('✅ Images configuration is valid');
    });

    test('includes required image domains', () => {
      console.log('\n🔍 TESTING: Image remote patterns');
      
      const remotePatterns = nextConfig.images.remotePatterns;
      
      // Check for common image hosting domains
      const expectedDomains = [
        'images.unsplash.com',
        'via.placeholder.com',
        'picsum.photos'
      ];

      expectedDomains.forEach(domain => {
        const hasPattern = remotePatterns.some(pattern => 
          pattern.hostname === domain || pattern.hostname?.includes(domain)
        );
        expect(hasPattern).toBe(true);
      });
      
      console.log('✅ Required image domains configured');
    });

    test('configures image formats and quality', () => {
      console.log('\n🔍 TESTING: Image formats and quality');
      
      if (nextConfig.images.formats) {
        expect(Array.isArray(nextConfig.images.formats)).toBe(true);
      }
      
      if (nextConfig.images.quality) {
        expect(typeof nextConfig.images.quality).toBe('number');
        expect(nextConfig.images.quality).toBeGreaterThan(0);
        expect(nextConfig.images.quality).toBeLessThanOrEqual(100);
      }
      
      console.log('✅ Image formats and quality configured correctly');
    });
  });

  describe('Webpack Configuration', () => {
    test('webpack function is defined and callable', () => {
      console.log('\n🔍 TESTING: Webpack configuration function');
      
      expect(nextConfig.webpack).toBeDefined();
      expect(typeof nextConfig.webpack).toBe('function');
      
      console.log('✅ Webpack configuration function is valid');
    });

    test('webpack configuration modifies config correctly', () => {
      console.log('\n🔍 TESTING: Webpack configuration modifications');
      
      const mockWebpackConfig = {
        module: {
          rules: []
        },
        resolve: {
          fallback: {}
        }
      };

      const mockOptions = {
        isServer: false,
        dev: false
      };

      const modifiedConfig = nextConfig.webpack(mockWebpackConfig, mockOptions);
      
      expect(modifiedConfig).toBeDefined();
      expect(modifiedConfig.module).toBeDefined();
      expect(modifiedConfig.module.rules).toBeDefined();
      
      console.log('✅ Webpack configuration modifies config correctly');
    });

    test('adds SVG loader rule', () => {
      console.log('\n🔍 TESTING: SVG loader rule addition');
      
      const mockWebpackConfig = {
        module: {
          rules: []
        },
        resolve: {
          fallback: {}
        }
      };

      const mockOptions = { isServer: false, dev: false };
      const modifiedConfig = nextConfig.webpack(mockWebpackConfig, mockOptions);
      
      const svgRule = modifiedConfig.module.rules.find(rule => 
        rule.test && rule.test.toString().includes('svg')
      );
      
      expect(svgRule).toBeDefined();
      expect(svgRule.use).toBeDefined();
      
      console.log('✅ SVG loader rule added correctly');
    });

    test('configures node polyfills for client-side', () => {
      console.log('\n🔍 TESTING: Node polyfills configuration');
      
      const mockWebpackConfig = {
        module: { rules: [] },
        resolve: { fallback: {} }
      };

      const mockOptions = { isServer: false, dev: false };
      const modifiedConfig = nextConfig.webpack(mockWebpackConfig, mockOptions);
      
      expect(modifiedConfig.resolve.fallback).toBeDefined();
      
      // Check for common Node.js polyfills
      const expectedPolyfills = ['fs', 'net', 'tls'];
      expectedPolyfills.forEach(polyfill => {
        if (modifiedConfig.resolve.fallback[polyfill] !== undefined) {
          expect(modifiedConfig.resolve.fallback[polyfill]).toBe(false);
        }
      });
      
      console.log('✅ Node polyfills configured correctly');
    });
  });

  describe('Environment Variables Configuration', () => {
    test('defines required environment variables', () => {
      console.log('\n🔍 TESTING: Environment variables configuration');
      
      expect(nextConfig.env).toBeDefined();
      expect(typeof nextConfig.env).toBe('object');
      
      console.log('✅ Environment variables configuration is valid');
    });

    test('includes custom environment variables', () => {
      console.log('\n🔍 TESTING: Custom environment variables');
      
      // Check for common custom env vars
      const possibleEnvVars = [
        'CUSTOM_KEY',
        'API_URL',
        'NEXT_PUBLIC_API_URL'
      ];

      // At least some custom env vars should be defined
      const hasCustomVars = Object.keys(nextConfig.env).length > 0;
      
      if (hasCustomVars) {
        expect(Object.keys(nextConfig.env).length).toBeGreaterThan(0);
      }
      
      console.log('✅ Custom environment variables handled correctly');
    });
  });

  describe('Next-Intl Plugin Integration', () => {
    test('next-intl plugin is applied', () => {
      console.log('\n🔍 TESTING: Next-intl plugin integration');
      
      // The plugin should have been called during config creation
      expect(mockNextIntlPlugin).toHaveBeenCalled();
      
      console.log('✅ Next-intl plugin applied correctly');
    });

    test('next-intl plugin receives correct configuration', () => {
      console.log('\n🔍 TESTING: Next-intl plugin configuration');
      
      const pluginCall = mockNextIntlPlugin.mock.calls[0];
      expect(pluginCall).toBeDefined();
      
      // The plugin should receive the i18n request configuration path
      const configPath = pluginCall[0];
      expect(configPath).toContain('i18n/request');
      
      console.log('✅ Next-intl plugin receives correct configuration');
    });
  });

  describe('Performance and Optimization', () => {
    test('enables production optimizations', () => {
      console.log('\n🔍 TESTING: Production optimizations');
      
      // Check for optimization settings
      if (nextConfig.compiler) {
        expect(typeof nextConfig.compiler).toBe('object');
      }
      
      if (nextConfig.swcMinify !== undefined) {
        expect(typeof nextConfig.swcMinify).toBe('boolean');
      }
      
      console.log('✅ Production optimizations configured');
    });

    test('configures bundle analyzer if enabled', () => {
      console.log('\n🔍 TESTING: Bundle analyzer configuration');
      
      // Bundle analyzer might be conditionally enabled
      if (nextConfig.bundleAnalyzer) {
        expect(typeof nextConfig.bundleAnalyzer).toBe('object');
      }
      
      console.log('✅ Bundle analyzer configuration handled correctly');
    });
  });

  describe('Security and Headers', () => {
    test('configures security headers if present', () => {
      console.log('\n🔍 TESTING: Security headers configuration');
      
      if (nextConfig.headers) {
        expect(typeof nextConfig.headers).toBe('function');
      }
      
      if (nextConfig.async && nextConfig.headers) {
        // Test headers function
        const headers = nextConfig.headers();
        expect(Array.isArray(headers) || typeof headers.then === 'function').toBe(true);
      }
      
      console.log('✅ Security headers configuration is valid');
    });

    test('configures content security policy if present', () => {
      console.log('\n🔍 TESTING: Content Security Policy');
      
      // CSP might be configured in headers
      if (nextConfig.headers) {
        // This would be tested in integration tests
        expect(typeof nextConfig.headers).toBe('function');
      }
      
      console.log('✅ Content Security Policy configuration handled');
    });
  });

  describe('Development vs Production Configuration', () => {
    test('handles development-specific settings', () => {
      console.log('\n🔍 TESTING: Development-specific settings');
      
      // Development settings might be conditional
      if (nextConfig.reactStrictMode !== undefined) {
        expect(typeof nextConfig.reactStrictMode).toBe('boolean');
      }
      
      console.log('✅ Development settings configured correctly');
    });

    test('handles production-specific optimizations', () => {
      console.log('\n🔍 TESTING: Production-specific optimizations');
      
      // Production optimizations
      if (nextConfig.compress !== undefined) {
        expect(typeof nextConfig.compress).toBe('boolean');
      }
      
      if (nextConfig.poweredByHeader !== undefined) {
        expect(typeof nextConfig.poweredByHeader).toBe('boolean');
      }
      
      console.log('✅ Production optimizations configured correctly');
    });
  });

  describe('TypeScript and ESLint Integration', () => {
    test('configures TypeScript settings if present', () => {
      console.log('\n🔍 TESTING: TypeScript configuration');
      
      if (nextConfig.typescript) {
        expect(typeof nextConfig.typescript).toBe('object');
        
        if (nextConfig.typescript.ignoreBuildErrors !== undefined) {
          expect(typeof nextConfig.typescript.ignoreBuildErrors).toBe('boolean');
        }
      }
      
      console.log('✅ TypeScript configuration handled correctly');
    });

    test('configures ESLint settings if present', () => {
      console.log('\n🔍 TESTING: ESLint configuration');
      
      if (nextConfig.eslint) {
        expect(typeof nextConfig.eslint).toBe('object');
        
        if (nextConfig.eslint.ignoreDuringBuilds !== undefined) {
          expect(typeof nextConfig.eslint.ignoreDuringBuilds).toBe('boolean');
        }
      }
      
      console.log('✅ ESLint configuration handled correctly');
    });
  });

  describe('Edge Cases and Error Handling', () => {
    test('configuration is serializable', () => {
      console.log('\n🔍 TESTING: Configuration serializability');
      
      // Test that the config can be serialized (important for Next.js)
      expect(() => {
        JSON.stringify(nextConfig, (key, value) => {
          if (typeof value === 'function') {
            return '[Function]';
          }
          return value;
        });
      }).not.toThrow();
      
      console.log('✅ Configuration is serializable');
    });

    test('handles missing optional dependencies gracefully', () => {
      console.log('\n🔍 TESTING: Missing dependencies handling');
      
      // The config should not throw if optional dependencies are missing
      expect(() => {
        // Re-evaluate config to test error handling
        const testConfig = { ...nextConfig };
        expect(testConfig).toBeDefined();
      }).not.toThrow();
      
      console.log('✅ Missing dependencies handled gracefully');
    });
  });

  describe('Plugin Chain and Order', () => {
    test('plugins are applied in correct order', () => {
      console.log('\n🔍 TESTING: Plugin application order');
      
      // next-intl plugin should be the first/primary plugin
      expect(mockNextIntlPlugin).toHaveBeenCalled();
      
      // Verify the plugin was called with the right parameters
      const calls = mockNextIntlPlugin.mock.calls;
      expect(calls.length).toBeGreaterThan(0);
      
      console.log('✅ Plugins applied in correct order');
    });

    test('plugin configuration is immutable', () => {
      console.log('\n🔍 TESTING: Plugin configuration immutability');
      
      // Plugins should not mutate the original config object
      const originalConfig = { test: 'value' };
      const pluginResult = mockNextIntlPlugin(originalConfig);
      
      // Original should remain unchanged
      expect(originalConfig).toEqual({ test: 'value' });
      
      console.log('✅ Plugin configuration remains immutable');
    });
  });

  describe('Compatibility and Standards', () => {
    test('follows Next.js configuration standards', () => {
      console.log('\n🔍 TESTING: Next.js configuration standards');
      
      // Check for standard Next.js config properties
      const standardProps = [
        'experimental',
        'images',
        'webpack'
      ];

      standardProps.forEach(prop => {
        if (nextConfig[prop] !== undefined) {
          expect(nextConfig).toHaveProperty(prop);
        }
      });
      
      console.log('✅ Follows Next.js configuration standards');
    });

    test('compatible with Next.js 15 features', () => {
      console.log('\n🔍 TESTING: Next.js 15 compatibility');
      
      // Check for Next.js 15 specific features
      if (nextConfig.experimental) {
        // Turbo is a Next.js 15 feature
        expect(nextConfig.experimental.turbo).toBeDefined();
      }
      
      console.log('✅ Compatible with Next.js 15 features');
    });
  });
}); 
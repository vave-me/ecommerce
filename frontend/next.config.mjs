// next.config.mjs ─ ESM syntax ----------------------------------------------
import createNextIntlPlugin from 'next-intl/plugin';   // <-- add ".js" in ESM
import { copyFileSync, existsSync, mkdirSync } from 'fs';
import { join } from 'path';
import withBundleAnalyzer from '@next/bundle-analyzer';

// ————————————————————————————————————————————————————————
// 1) Wrap the base Next.js config with next-intl
//    (request config is auto-detected at src/i18n/request.jsx)
const withNextIntl = createNextIntlPlugin();

// Function to copy ads.txt to output directory
const copyAdsFile = () => {
  try {
    const adsTextPath = join(process.cwd(), 'public', 'ads.txt');
    const outputDir = join(process.cwd(), '.next');
    
    if (existsSync(adsTextPath)) {
      // Ensure the destination directory exists
      if (!existsSync(outputDir)) {
        mkdirSync(outputDir, { recursive: true });
      }
      
      // Copy the ads.txt file directly to the output directory
      copyFileSync(adsTextPath, join(outputDir, 'ads.txt'));
      console.log('✅ ads.txt file copied to .next/');
    }
  } catch (error) {
    console.error('Error copying ads.txt file:', error);
  }
};

/**
 * @type {import('next').NextConfig}
 */
const nextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  output: 'standalone',
  trailingSlash: false,
  generateBuildId: () => 'build-' + Date.now(),
  
  // Enhanced image optimization settings
  images: {
    formats: ['image/avif', 'image/webp'],
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
      {
        protocol: 'http',
        hostname: '**',
      },
      // Specific patterns for local development
      {
        protocol: 'http',
        hostname: '192.168.178.84',
        port: '9096',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9096',
      },
    ],
    minimumCacheTTL: 60,
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    dangerouslyAllowSVG: true,
    contentSecurityPolicy: "default-src 'self'; script-src 'none'; sandbox;",
  },
  
  // Enhanced compiler optimizations
  compiler: {
    // Remove ALL console statements in production
    removeConsole: process.env.NODE_ENV === 'production' ? true : false,
    reactRemoveProperties: process.env.NODE_ENV === 'production' ? { 
      properties: ['^data-testid$', '^data-cy$', '^data-test'] 
    } : false,
  },
  
  // Performance optimizations - temporarily reduced
  experimental: {
    optimizeServerReact: false,
    serverMinification: false,
  },
  
  // Turbopack configuration (moved from experimental.turbo)
  turbopack: {
    rules: {
      '*.svg': {
        loaders: ['@svgr/webpack'],
        as: '*.js',
      },
    },
  },
  
  // Enhanced webpack configuration for bundle optimization
  webpack: (config, { isServer, dev }) => {
    // SVG handling
    config.module.rules.push({
      test: /\.svg$/,
      use: ['@svgr/webpack'],
    });
    
    // Optimize bundle splitting
    if (!isServer && !dev) {
      config.optimization = {
        ...config.optimization,
        splitChunks: {
          chunks: 'all',
          minSize: 20000,
          maxSize: 244000,
          cacheGroups: {
            // Framework chunk (React, Next.js core)
            framework: {
              test: /[\\/]node_modules[\\/](react|react-dom|next)[\\/]/,
              name: 'framework',
              chunks: 'all',
              priority: 40,
              reuseExistingChunk: true,
              enforce: true,
            },
            // Vendor chunk for stable dependencies
            vendor: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              chunks: 'all',
              priority: 10,
              reuseExistingChunk: true,
              minChunks: 2,
            },
            // TipTap editor chunk (heavy, load on demand)
            editor: {
              test: /[\\/]node_modules[\\/]@tiptap[\\/]/,
              name: 'editor',
              chunks: 'async',
              priority: 30,
              reuseExistingChunk: true,
            },
            // React Query chunk
            reactQuery: {
              test: /[\\/]node_modules[\\/]@tanstack[\\/]react-query[\\/]/,
              name: 'react-query',
              chunks: 'all',
              priority: 25,
              reuseExistingChunk: true,
            },
            // Icon libraries (tree-shake heavily)
            icons: {
              test: /[\\/]node_modules[\\/](react-icons|lucide-react)[\\/]/,
              name: 'icons',
              chunks: 'async',
              priority: 20,
              reuseExistingChunk: true,
            },
            // Modal components chunk (lazy loaded)
            modals: {
              test: /[\\/]src[\\/]features[\\/].*Modal[\\/]/,
              name: 'modals',
              chunks: 'async',
              priority: 15,
              reuseExistingChunk: true,
            },
            // AI components (large and optional)
            ai: {
              test: /[\\/]src[\\/]components[\\/]AI[\\/]/,
              name: 'ai-components',
              chunks: 'async',
              priority: 18,
              reuseExistingChunk: true,
            },
            // Common chunk for shared components
            common: {
              name: 'common',
              minChunks: 3,
              chunks: 'all',
              priority: 5,
              reuseExistingChunk: true,
              maxSize: 200000,
            },
          },
        },
        // Advanced tree shaking
        usedExports: true,
        sideEffects: false,
        concatenateModules: true,
      };
    }
    
    // Aggressive tree shaking for icon libraries
    config.resolve.alias = {
      ...config.resolve.alias,
      // Force ESM tree shaking for react-icons
      'react-icons/fa': 'react-icons/fa/index.esm.js',
      'react-icons/md': 'react-icons/md/index.esm.js',
      'react-icons/ti': 'react-icons/ti/index.esm.js',
      'react-icons/ai': 'react-icons/ai/index.esm.js',
      'react-icons/hi': 'react-icons/hi/index.esm.js',
    };
    
    // Remove unused CSS and optimize CSS splitting
    if (!dev && config.optimization && config.optimization.splitChunks) {
      config.optimization.splitChunks.cacheGroups = {
        ...config.optimization.splitChunks.cacheGroups,
        styles: {
          name: 'styles',
          test: /\.(css|scss)$/,
          chunks: 'all',
          enforce: true,
          priority: 50,
        }
      };
    }
    
    // Production TDZ Error Prevention
    if (process.env.NODE_ENV === 'production') {
      config.optimization = {
        ...config.optimization,
        minimize: true,
        minimizer: config.optimization.minimizer?.map((minimizer) => {
          if (minimizer.constructor.name === 'TerserPlugin') {
            minimizer.options = {
              ...minimizer.options,
              terserOptions: {
                ...minimizer.options.terserOptions,
                compress: {
                  ...minimizer.options.terserOptions?.compress,
                  // Prevent variable hoisting that causes TDZ errors
                  hoist_vars: false,
                  hoist_funs: false,
                  // Keep function names to prevent React.memo issues
                  keep_fnames: true,
                  // Prevent unsafe transformations
                  unsafe: false,
                  unsafe_arrows: false,
                  unsafe_comps: false,
                  unsafe_Function: false,
                  unsafe_math: false,
                  unsafe_symbols: false,
                  unsafe_methods: false,
                  unsafe_proto: false,
                  unsafe_regexp: false,
                  unsafe_undefined: false,
                  // Preserve order of operations
                  sequences: false,
                  // Don't merge variable declarations
                  join_vars: false,
                },
                mangle: {
                  ...minimizer.options.terserOptions?.mangle,
                  // Keep function names for React.memo
                  keep_fnames: true,
                  // Prevent mangling of variables that might cause TDZ
                  reserved: ['memo', 'React', 'useCallback', 'useMemo', 'useState', 'useEffect'],
                },
                output: {
                  ...minimizer.options.terserOptions?.output,
                  // Keep function names in output
                  keep_quoted_props: true,
                  preserve_annotations: true,
                }
              }
            };
          }
          return minimizer;
        }),
      };

      // Additional production optimizations
      if (!isServer) {
        // Prevent module hoisting issues
        config.optimization.providedExports = false;
        config.optimization.usedExports = false;
        config.optimization.sideEffects = false;
      }
    }
    
    // Bundle analyzer (development only) - handled by withBundleAnalyzer wrapper
    
    return config;
  },
  
  // Output configuration removed - conflicts with standalone
  trailingSlash: false,
  
  // API route proxying
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.API_URL ?? 'http://192.168.178.84:8080'}/api/:path*`
      },
      {
        source: '/ads.txt',
        destination: '/ads.txt'
      },
      // Fix locale-prefixed image paths
      {
        source: '/:locale/images/:path*',
        destination: '/images/:path*'
      }
    ];
  },
  
  // Headers for better caching and security
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-XSS-Protection',
            value: '1; mode=block',
          },
        ],
      },
      {
        source: '/static/(.*)',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
    ];
  },
  
  // Redirect configuration
  async redirects() {
    return [
      // Add any necessary redirects here
    ];
  },
  
  // Environment variables
  env: {
    CUSTOM_KEY: REDACTED
  },
  
  // TypeScript configuration
  typescript: {
    ignoreBuildErrors: false,
  },
  
  // ESLint configuration
  eslint: {
    ignoreDuringBuilds: true,
  },
  
  // Standalone output for better deployment
  output: process.env.NODE_ENV === 'production' ? 'standalone' : undefined,
};

// Try to copy the ads.txt file at configuration time as well
copyAdsFile();

// Enable bundle analyzer in analyze mode
const withAnalyzer = withBundleAnalyzer({
  enabled: process.env.ANALYZE === 'true',
});

// Apply enhancements with correct wrapping order
export default withNextIntl(
  withAnalyzer(nextConfig)
);

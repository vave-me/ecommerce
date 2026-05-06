"use client";
import React, { useMemo, memo } from 'react';

// Initialize auth token management before any other imports
import '../api/initializeAuth.js';

import { Provider as ReduxProvider } from 'react-redux';
import { QueryClientProvider } from '@tanstack/react-query';
import { ToastContainer } from 'react-toastify';
import store from '../lib/store';
import { AuthProvider } from '../context/AuthContext';
import { ThemeProvider } from '../context/ThemeContext';
// Removed ThemeProvider - now using Redux for theme management
import { NavBarProvider } from '../context/NavBarContext';
import { OverlayProvider } from '../context/OverlayContext';
import { CategoriesProvider } from '../hooks/useCategories';
import { ShareProvider } from '../components/Shared/ShareComponent/ShareProvider';
import AppInitializer from '../components/AppInitializer';
import ThemeInitializer from '../components/Theme/ThemeInitializer';
import UIPreferencesHydrator from '../components/UIPreferencesHydrator';
import ProvidersComposer from './ProvidersComposer';
import { createQueryClient } from '../lib/reactQuery';
import { ErrorBoundary } from '../components/ErrorBoundary';
import 'react-toastify/dist/ReactToastify.css';
/**
 * Unified application providers with improved performance
 * - Uses provider composition for better maintainability
 * - Implements optimized React Query configuration
 * - Configures persistence only for static data
 * - Optimizes provider nesting order for dependencies
 * - Compatible with both root layouts
 * - Includes ShareProvider for share functionality
 * - Includes ThemeProvider for dark/light theme support
 */
const Providers = memo(function Providers({ children }) {
  // Create a stable QueryClient instance
  const queryClient = useMemo(() => createQueryClient(), []);
  return (
    <React.StrictMode>
      <ErrorBoundary 
        name="RootProviders"
        showReload={true}
        fallback={(error, errorInfo, reset) => (
          <div style={{ padding: '2rem', textAlign: 'center' }}>
            <h1>Application Error</h1>
            <p>The application encountered an error during initialization.</p>
            <button onClick={() => window.location.reload()}>Reload Application</button>
          </div>
        )}
      >
        <ProvidersComposer
          providers={[
          // Add providers in dependency order (inside-out)
          props => <ReduxProvider store={store} {...props} />,
          props => <QueryClientProvider client={queryClient} {...props} />,
          // Include ThemeProvider for light/dark/system theme handling
          props => <ThemeProvider {...props} />,
          // Use original AuthProvider that exports useAuth() hook
          props => <AuthProvider {...props} />,
          // Include additional providers needed by the app
          props => <NavBarProvider {...props} />,
          props => <OverlayProvider {...props} />,
          props => <CategoriesProvider 
                    prefetchTopics={['marketplace', 'services', 'deals', 'automotive', 'property']} 
                    {...props} 
                  />,
          // Add ShareProvider for share functionality
          props => <ShareProvider 
                    defaultConfig={{
                      defaultPlatforms: ['native', 'copy', 'facebook', 'twitter', 'linkedin', 'email'],
                      requireAuthForSharing: false,
                      enableAnalytics: true,
                      enableShareCounts: true,
                      trackShares: true
                    }}
                    {...props} 
                  />
        ]}
      >
        {/* Include AppInitializer which was in the original providers.jsx */}
        <AppInitializer />
        {/* Theme debug/initialization */}
        <ThemeInitializer />
        {/* Hydrate UI preferences from localStorage */}
        <UIPreferencesHydrator />
        {children}
        <ToastContainer
          position="bottom-right"
          autoClose={5000}
          hideProgressBar={false}
          newestOnTop
          closeOnClick
          pauseOnFocusLoss={false} // Performance improvement
          draggable
          pauseOnHover
        />
      </ProvidersComposer>
      </ErrorBoundary>
    </React.StrictMode>
  );
});
export default Providers;
// Also export named function for compatibility with existing imports
export { Providers };
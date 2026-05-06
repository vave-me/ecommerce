import React from 'react';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';
import { makeStore } from '../lib/store';

// Create a custom render function that includes providers
export function renderWithProviders(ui, options = {}) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        // Turn off retries for testing
        retry: false,
        // Don't refetch on window focus during tests
        refetchOnWindowFocus: false,
      },
    },
    // Set this flag to avoid warning in tests about QueryClient already being used
    logger: {
      error: () => {},
      warn: () => {},
      log: () => {},
    },
  });

  function Wrapper({ children }) {
    return (
      <Provider store={makeStore}>
        <QueryClientProvider client={queryClient}>
          {children}
        </QueryClientProvider>
      </Provider>
    );
  }

  return render(ui, { wrapper: Wrapper, ...options });
}

// Create a wrapper for testing hooks that use React Query
export function createQueryClientWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
      },
    },
    logger: {
      error: () => {},
      warn: () => {},
      log: () => {},
    },
  });

  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
} 
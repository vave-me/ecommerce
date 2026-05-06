import '@testing-library/jest-dom';
import { configure } from '@testing-library/react';
import axios from 'axios';
// Configure React Testing Library to use more appropriate defaults
configure({
  asyncUtilTimeout: 5000, // Increase the timeout for async operations
});
// Set up environment variables for testing
process.env.NEXT_PUBLIC_API_BASE_URL = 'http://192.168.178.84:8080';
process.env.NODE_ENV = 'test';
process.env.JWT_SECRET = 'test_secret_key_for_jwt_tokens';
// Mock auth utils module
jest.mock('../utils/auth.utils', () => require('../tests/mocks/auth.utils.mock'));
// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: jest.fn((key) => {
      return store[key] || null;
    }),
    setItem: jest.fn((key, value) => {
      store[key] = value.toString();
    }),
    removeItem: jest.fn((key) => {
      delete store[key];
    }),
    clear: jest.fn(() => {
      store = {};
    }),
  };
})();
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});
// Don't mock axios completely here since we use MockAdapter in tests
// Instead, just provide default implementations for the methods we use
const originalAxios = jest.requireActual('axios');
axios.create = jest.fn(() => axios);
axios.interceptors = {
  request: { use: jest.fn(), eject: jest.fn() },
  response: { use: jest.fn(), eject: jest.fn() }
};
axios.get = jest.fn().mockImplementation(() => Promise.resolve({ data: {} }));
axios.post = jest.fn().mockImplementation(() => Promise.resolve({ data: {} }));
axios.put = jest.fn().mockImplementation(() => Promise.resolve({ data: {} }));
axios.delete = jest.fn().mockImplementation(() => Promise.resolve({ data: {} }));
// Mock Next.js router
jest.mock('next/navigation', () => ({
  useRouter() {
    return {
      push: jest.fn(),
      back: jest.fn(),
      events: {
        on: jest.fn(),
        off: jest.fn(),
      },
      pathname: '/',
      query: {},
    };
  },
  useSearchParams() {
    return {
      get: jest.fn(),
    };
  },
  usePathname() {
    return '/';
  },
}));
// Fix for the "not wrapped in act" warning
// This is a known issue with React 18 and the current version of testing-library
global.IS_REACT_ACT_ENVIRONMENT = true;
// Mock console methods for cleaner test output
const originalConsoleError = console.error;
const originalConsoleWarn = console.warn;
const originalConsoleLog = console.log;
// Suppress specific error messages
console.error = (...args) => {
  if (
    /Warning.*not wrapped in act/.test(args[0]) ||
    /Error: Uncaught \[Error: expect/.test(args[0]) ||
    /FATAL ERROR: NEXT_PUBLIC_API_BASE_URL/.test(args[0]) || // Suppress the API base URL error
    /The current testing environment is not configured to support act/.test(args[0]) || // Suppress act warnings
    /Failed to set refresh token cookie/.test(args[0]) || // Suppress token cookie errors
    /Failed to clear refresh token cookie/.test(args[0]) || // Suppress token cookie errors
    /The argument passed to useReducer must be a reducer function/.test(args[0]) // Suppress React reducer warnings
  ) {
    return;
  }
  originalConsoleError(...args);
};
console.warn = (...args) => {
  if (
    /Warning.*not wrapped in act/.test(args[0]) ||
    /outdated JSX transform/.test(args[0])
  ) {
    return;
  }
  originalConsoleWarn(...args);
};
console.log = (...args) => {
  // Uncomment to suppress all console.log during tests
  // return;
  originalConsoleLog(...args);
};
// Create missing style mocks
jest.mock('../tests/__mocks__/styleMock.js', () => ({}), { virtual: true });
jest.mock('../tests/__mocks__/fileMock.js', () => 'test-file-stub', { virtual: true });
// Overwrite some fetch behaviors for tests
global.fetch = jest.fn();
// Cleanup on exiting
afterAll(() => {
  console.error = originalConsoleError;
  console.warn = originalConsoleWarn;
  console.log = originalConsoleLog;
});
// Add global timing cleanup
afterEach(() => {
  jest.useRealTimers();
  jest.clearAllMocks();
}); 
import '@testing-library/jest-dom';
import { configure } from '@testing-library/react';
import axios from 'axios';
import { jest } from '@jest/globals';

// Configure React Testing Library to use more appropriate defaults
configure({
  asyncUtilTimeout: 5000, // Increase the timeout for async operations
});

// Configure next-intl for testing environment
process.env.NODE_ENV = 'test';

// Set up environment variables for testing
process.env.NEXT_PUBLIC_API_BASE_URL = 'http://192.168.178.84:8080';
process.env.JWT_SECRET = 'test_secret_key_for_jwt_tokens';

// Mock auth utils module - updated path
jest.mock('../src/utils/auth.utils', () => require('./mocks/auth.utils.mock'));

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

// Mock sessionStorage
const sessionStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
global.sessionStorage = sessionStorageMock;

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// Mock IntersectionObserver
global.IntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Mock ResizeObserver
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Suppress console warnings during tests (optional)
const originalError = console.error;
beforeAll(() => {
  console.error = (...args) => {
    if (
      typeof args[0] === 'string' &&
      args[0].includes('Warning: ReactDOM.render is no longer supported')
    ) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterAll(() => {
  console.error = originalError;
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

// Mock Next.js router for components that use it
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    forward: jest.fn(),
    refresh: jest.fn(),
    prefetch: jest.fn(),
    pathname: '/',
    route: '/',
    query: {},
    asPath: '/'
  }),
  usePathname: () => '/',
  useSearchParams: () => new URLSearchParams(),
  useParams: () => ({ locale: 'en' }),
  notFound: jest.fn()
}));

// Mock Next.js Image component
jest.mock('next/image', () => {
  return function MockImage({ src, alt, ...props }) {
    const React = require('react');
    return React.createElement('img', { src, alt, ...props });
  };
});

// Mock Next.js Link component
jest.mock('next/link', () => {
  return function MockLink({ children, href, ...props }) {
    const React = require('react');
    return React.createElement('a', { href, ...props }, children);
  };
});

// Mock CSS modules
jest.mock('*.module.css', () => ({}));
jest.mock('*.module.scss', () => ({}));

// Mock static file imports
jest.mock('*.svg', () => 'svg-mock');
jest.mock('*.png', () => 'png-mock');
jest.mock('*.jpg', () => 'jpg-mock');
jest.mock('*.jpeg', () => 'jpeg-mock');
jest.mock('*.gif', () => 'gif-mock');
jest.mock('*.webp', () => 'webp-mock');

// Mock environment variables
process.env.NEXT_PUBLIC_API_URL = 'http://localhost:3000/api';

// Global test utilities
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

global.IntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// Console filtering for cleaner test output
const originalConsoleError = console.error;
const originalConsoleWarn = console.warn;
const originalConsoleLog = console.log;

console.error = (...args) => {
  // Filter out known React warnings that are not relevant to tests
  const message = args[0];
  if (typeof message === 'string') {
    if (
      message.includes('Warning: ReactDOM.render is no longer supported') ||
      message.includes('Warning: An update to') ||
      message.includes('act(...)')
    ) {
      return;
    }
  }
  originalConsoleError(...args);
};

console.warn = (...args) => {
  const message = args[0];
  if (typeof message === 'string') {
    if (
      message.includes('componentWillReceiveProps') ||
      message.includes('componentWillUpdate') ||
      message.includes('componentWillMount')
    ) {
      return;
    }
  }
  originalConsoleWarn(...args);
};

console.log = (...args) => {
  // Allow test logs to pass through
  originalConsoleLog(...args);
};

// Restore original console methods for debugging
console.originalConsoleError = originalConsoleError;
console.originalConsoleWarn = originalConsoleWarn;
console.originalConsoleLog = originalConsoleLog;

// Fix for the "not wrapped in act" warning
// This is a known issue with React 18 and the current version of testing-library
global.IS_REACT_ACT_ENVIRONMENT = true;

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
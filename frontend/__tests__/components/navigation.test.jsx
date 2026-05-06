import { Link, redirect, usePathname, useRouter, getPathname } from '../i18n/navigation.js';
import * as nextIntlNavigation from 'next-intl/navigation';
import { routing } from '../i18n/routing.js';

// Mock next-intl/navigation module
jest.mock('next-intl/navigation', () => ({
  createNavigation: jest.fn(() => ({
    Link: 'MockedLink',
    redirect: jest.fn(),
    usePathname: jest.fn(),
    useRouter: jest.fn(),
    getPathname: jest.fn()
  }))
}));

// Mock the routing module
jest.mock('../i18n/routing.js', () => ({
  routing: {
    locales: ['en', 'de', 'pl'],
    defaultLocale: 'en',
    localePrefix: 'always'
  }
}));

describe('navigation module', () => {
  it('should create navigation with the correct routing config', () => {
    // Check if createNavigation was called with routing config
    expect(nextIntlNavigation.createNavigation).toHaveBeenCalledWith(routing);
  });

  it('should export the correct navigation utilities', () => {
    // Check if all navigation utilities are exported
    expect(Link).toBe('MockedLink');
    expect(typeof redirect).toBe('function');
    expect(typeof usePathname).toBe('function');
    expect(typeof useRouter).toBe('function');
    expect(typeof getPathname).toBe('function');
  });
}); 
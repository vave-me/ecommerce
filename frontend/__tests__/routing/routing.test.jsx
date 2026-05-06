import { routing } from '../i18n/routing.js';
import * as nextIntlRouting from 'next-intl/routing';

// Mock next-intl/routing module
jest.mock('next-intl/routing', () => ({
  defineRouting: jest.fn((config) => config)
}));

describe('routing module', () => {
  it('should define routing with correct configuration', () => {
    // Check if defineRouting was called with the correct config
    expect(nextIntlRouting.defineRouting).toHaveBeenCalledWith({
      locales: ['en', 'de', 'pl'],
      defaultLocale: 'en',
      localePrefix: 'always'
    });
  });

  it('should export the routing configuration', () => {
    // Check if routing config is exported correctly
    expect(routing).toEqual({
      locales: ['en', 'de', 'pl'],
      defaultLocale: 'en',
      localePrefix: 'always'
    });
  });
}); 
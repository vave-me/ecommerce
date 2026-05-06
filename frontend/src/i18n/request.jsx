import {getRequestConfig} from 'next-intl/server';
import {hasLocale} from 'next-intl';
import {routing} from './routing';
// Performance monitoring utility (server-side only)
const performanceMetrics = {
  apiCalls: new Map(),
  trackApiCall(endpoint, duration) {
    const calls = this.apiCalls.get(endpoint) || [];
    calls.push(duration);
    this.apiCalls.set(endpoint, calls);
  },
  getMetrics() {
    return {
      apiCalls: Object.fromEntries(this.apiCalls)
    };
  }
};
export default getRequestConfig(async ({requestLocale}) => {
    const startTime = performance.now();
    // Determine the locale to use for this request:
    const requestedLocale = await requestLocale;
    const locale = hasLocale(routing.locales, requestedLocale)
        ? requestedLocale
        : routing.defaultLocale;
    const result = {
        locale,
        messages: (await import(`../../messages/${locale}.json`)).default
    };
    // Track API call duration
    const duration = performance.now() - startTime;
    performanceMetrics.trackApiCall('locale-config', duration);
    return result;
});

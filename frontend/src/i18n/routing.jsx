import {defineRouting} from 'next-intl/routing';
export const routing = defineRouting({
    locales: ['en', 'de', 'pl', 'it'],  // supported locales
    defaultLocale: 'de',          // default fallback
    localePrefix: 'always'       // URLs are always /en/… or /pl/…
});

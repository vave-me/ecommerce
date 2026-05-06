import {createNavigation} from 'next-intl/navigation';
import {routing} from './routing';
// Locale-aware wrappers around Next.js navigation helpers
export const {
    Link,        // Link component that prefixes paths with /en, /pl, …
    redirect,    // redirect('/about') → redirect('/en/about') in EN locale
    usePathname, // returns the locale-stripped pathname
    useRouter,   // router that keeps the active locale
    getPathname  // helper for RSCs
} = createNavigation(routing);

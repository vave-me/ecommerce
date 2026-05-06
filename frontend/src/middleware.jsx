import createMiddleware from 'next-intl/middleware';
import {routing} from './i18n/routing';
export default createMiddleware(routing);
export const config = {
    // Apply middleware to all paths except API, assets, or Next internals:
    matcher: [
        '/((?!api|trpc|_next|_vercel|.*\\..*).*)'
    ]
};

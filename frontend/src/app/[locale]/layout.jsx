/* ────────────────────────────────────────────────
   src/app/[locale]/layout.jsx   (fixed version)
───────────────────────────────────────────────── */
import {NextIntlClientProvider, hasLocale} from 'next-intl';
import {notFound} from 'next/navigation';
import {routing} from '../../i18n/routing';   // locales list
import Providers from '../Providers';
import ClientLayout from '../ClientLayout.client';
import AdSenseScript from '../../components/AdSense/AdSenseScript.client';
import PerformanceMonitor from '../../components/PerformanceMonitor/PerformanceMonitor.client';
import PWAInitializer from '../../components/PWA/PWAInitializer';
import WebVitalsReporter from '../../components/WebVitals/WebVitalsReporter.client';
import '../global.css';
import '../../utils/productionConsoleCleanup';
import '../../utils/hydrationFixer';
import '../../styles/optimizedEffects.css';
/* Optional static metadata */
export const metadata = {
    title: 'sfx markt',
    description: 'sfx markt – Live Marketplace',
    // Add Google AdSense verification
    other: {
        'google-adsense-account': 'ca-pub-7872277873986607'
    }
};
/** Locale-aware root layout (Server Component) */
export default async function LocaleLayout({children, params}) {
    const {locale} = await params;               // ✅ read before any await
    if (!hasLocale(routing.locales, locale)) notFound();
    /* Render the full HTML shell */
    return (
        <html lang={locale} suppressHydrationWarning>
            <head>
                <meta charSet="utf-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                <link rel="icon" href="/favicon.ico" />
                {/* CRITICAL: Preconnect to Google Fonts for faster loading */}
                <link rel="preconnect" href="https://fonts.googleapis.com" />
                <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
                {/* CRITICAL: Inject critical CSS inline for fastest FCP */}
                <style suppressHydrationWarning>{`
                    :root {
                        --primary-blue: #2563eb;
                        --gray-50: #f9fafb;
                        --gray-900: #111827;
                    }
                    body {
                        margin: 0;
                        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
                        background: var(--gray-50);
                        color: var(--gray-900);
                    }
                    .header-container {
                        position: sticky;
                        top: 0;
                        z-index: 50;
                        background: white;
                        border-bottom: 1px solid #e5e7eb;
                        box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
                    }
                    .loading-spinner {
                        width: 20px;
                        height: 20px;
                        border: 2px solid #e5e7eb;
                        border-radius: 50%;
                        border-top-color: var(--primary-blue);
                        animation: spin 1s linear infinite;
                    }
                    @keyframes spin { to { transform: rotate(360deg); } }
                `}</style>
                {/* CRITICAL: Optimized font loading - combined into single request */}
                <link
                    href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;600&family=Montserrat:wght@400;500;600;700&display=swap"
                    rel="stylesheet"
                />
            </head>
            <body suppressHydrationWarning>
                {/* Performance monitoring for development */}
                <PerformanceMonitor />
                {/* PWA initialization and features */}
                <PWAInitializer />
                {/* Web Vitals reporting */}
                <WebVitalsReporter />
                {/* Google AdSense - Disabled until proper publisher ID is configured */}
                {/* <AdSenseScript /> */}
                {/* Internationalization provider */}
                <NextIntlClientProvider locale={locale}>
                    {/* Application providers */}
                    <Providers>
                        {/* Layout with navigation */}
                        <ClientLayout>{children}</ClientLayout>
                    </Providers>
                </NextIntlClientProvider>
            </body>
        </html>
    );
}
/* Pre-generate static paths for every locale */
export function generateStaticParams() {
    return routing.locales.map(locale => ({locale}));
}

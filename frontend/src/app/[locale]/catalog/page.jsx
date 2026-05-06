"use client";
import React, { Suspense, memo } from 'react';
import { useTranslations } from 'next-intl';
import { useAuth } from '../../../context/AuthContext';
import CatalogPageClient from './CatalogPage.client';
import styles from './CatalogPage.module.css';
// Load debug utility in development only
if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
    import('../../../utils/testCatalogAPI.js').catch(() => {
        // Gracefully handle missing debug utility
    });
}
/**
 * Sophisticated Loading Component - Matching Header's loading patterns
 * Performance optimized with minimal DOM elements
 */
const CatalogLoading = memo(function CatalogLoading() {
    return (
        <div className={styles.loadingContainer}>
            <div className={styles.loadingHeader}>
                <div className={styles.loadingTitle} aria-label="Loading title"></div>
                <div className={styles.loadingControls}>
                    <div className={styles.loadingButton} aria-label="Loading button"></div>
                    <div className={styles.loadingButton} aria-label="Loading button"></div>
                </div>
            </div>
            <div className={styles.loadingGrid} role="status" aria-label="Loading catalog items">
                {Array.from({ length: 8 }, (_, i) => (
                    <div key={`loading-card-${i}`} className={styles.loadingCard}>
                        <div className={styles.loadingImage} aria-label="Loading image"></div>
                        <div className={styles.loadingText} aria-label="Loading text"></div>
                        <div className={styles.loadingText} aria-label="Loading text"></div>
                    </div>
                ))}
            </div>
        </div>
    );
});
/**
 * Authentication Error Component - Sophisticated error handling
 * Following Header's error state patterns
 */
const AuthError = memo(function AuthError({ t }) {
    return (
        <div className={styles.authErrorContainer} role="alert">
            <h1 className={styles.authErrorTitle}>{t('errorTitle')}</h1>
            <p className={styles.authErrorMessage}>{t('userIdRequired')}</p>
            <a 
                href="/login" 
                className={styles.authErrorLoginButton}
                aria-label="Go to login page"
            >
                {t('loginPrompt')}
            </a>
        </div>
    );
});
/**
 * Main Catalog Page Content - Optimized with better state management
 * Following Header's component architecture patterns
 */
const CatalogPageContent = memo(function CatalogPageContent({ searchParams, locale }) {
    // ALL HOOKS MUST BE CALLED FIRST - before any conditional returns
    const t = useTranslations('CatalogPage');
    const { user, isLoading: authLoading } = useAuth();
    // Optimized labels object - memoized to prevent recreation
    // This hook MUST be called every render, regardless of auth state
    const labels = React.useMemo(() => ({
        pageTitle: t('pageTitle'),
        empty: t('empty'),
        loading: t('loading'),
        error: t('error'),
        loadCatalog: t('loadCatalog'),
        errorFetching: t('errorFetching'),
        filterBy: t('filterBy'),
        sortBy: t('sortBy'),
        searchPlaceholder: t('searchPlaceholder'),
        types: {
            all: t('types.all'),
            products: t('types.products'),
            posts: t('types.posts', { fallback: 'Posts' }),
        },
        viewModes: {
            grid: t('viewModes.grid'),
            list: t('viewModes.list'),
            map: t('viewModes.map'),
        }
    }), [t]);
    // NOW we can do conditional rendering - after all hooks are called
    if (authLoading) {
        return <CatalogLoading />;
    }
    if (!user?.userId) {
        return <AuthError t={t} />;
    }
    return (
        <>
            {/* SEO Metadata - Following Header's SEO patterns */}
            <title>{labels.pageTitle}</title>
            <meta name="description" content={t('pageDescription')} />
            <meta name="keywords" content={t('pageKeywords')} />
            <meta name="robots" content="noindex, nofollow" /> {/* Private catalog page */}
            {/* Client Component - Performance optimized */}
            <CatalogPageClient
                userId={user.userId}
                searchParams={searchParams}
                locale={locale}
                labels={labels}
            />
        </>
    );
});
/**
 * Advanced Error Boundary - Following Header's error handling patterns
 */
class CatalogErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false, error: null };
    }
    static getDerivedStateFromError(error) {
        return { hasError: true, error };
    }
    componentDidCatch(error, errorInfo) {
        // Log error in development
        if (process.env.NODE_ENV === 'development') {
        }
    }
    render() {
        if (this.state.hasError) {
            return (
                <div className={styles.container}>
                    <div className={styles.authErrorContainer} role="alert">
                        <h1 className={styles.authErrorTitle}>Something went wrong</h1>
                        <p className={styles.authErrorMessage}>
                            We encountered an error while loading your catalog. Please try refreshing the page.
                        </p>
                        <button 
                            onClick={() => window.location.reload()}
                            className={styles.authErrorLoginButton}
                            type="button"
                        >
                            Refresh Page
                        </button>
                    </div>
                </div>
            );
        }
        return this.props.children;
    }
}
// Helper function to safely resolve props in Next.js 15+
async function resolveProps(props) {
    // In Next.js 15+, we need to await both params and searchParams
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return { params, searchParams };
}
/**
 * Main Export - Async server component for Next.js 15 compatibility
 * Following the pattern used in other pages for proper async params handling
 */
export default async function CatalogPage(props) {
    // Safely resolve async params and searchParams
    const { params, searchParams } = await resolveProps(props);
    const { locale } = params;
    return (
        <CatalogErrorBoundary>
            <Suspense fallback={<CatalogLoading />}>
                <CatalogPageContent 
                    searchParams={searchParams} 
                    locale={locale}
                />
            </Suspense>
        </CatalogErrorBoundary>
    );
} 
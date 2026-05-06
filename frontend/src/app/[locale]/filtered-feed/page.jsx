import React from 'react';
import { getTranslations } from 'next-intl/server';
import { notFound } from 'next/navigation';
import HorizontalFilters from '../../../components/Filters/HorizontalFilters';
import ConnectedFeedDisplay from '../../../components/Feed/ConnectedFeedDisplay';
// FeedProvider removed - using direct hooks
import Composer from '../../../components/Feed/Composer';
import { unifiedSearch } from '../../../api/searchApi';
import { fetchMainCategories } from '../../../api/categories';
import styles from './FilteredFeed.module.css';
export const dynamic = 'force-dynamic';
// Helper function to safely resolve props in Next.js 15+
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return { params, searchParams };
}
/* ── SEO metadata ─────────────────────────────────────────── */
export async function generateMetadata({ params }) {
    const { locale } = await params;
    const t = await getTranslations({ locale, namespace: 'FilteredFeed' });
    return {
        title: t('title', { default: 'Filtered Feed – sfx markt' }),
        description: t('description', {
            default: 'Browse and filter content across all categories. Discover products, services, deals, and more.',
        }),
        keywords: t('keywords', {
            default: 'filtered feed, browse, search, categories, marketplace, local deals',
        }),
    };
}
/**
 * Filtered Feed Page - Server Component
 * Fetches initial data server-side and provides it to client components
 */
export default async function FilteredFeedPage(props) {
    const { params, searchParams } = await resolveProps(props);
    const { locale } = params;
    try {
        const t = await getTranslations({ locale, namespace: 'FilteredFeed' });
        // Robustly determine entityTypes for initial load
        let defaultEntityTypes = ["product", "post", "deal", "vehicle", "property", "service", "job"];
        let entityTypesFromParams = defaultEntityTypes;
        if (searchParams?.types && searchParams.types.trim() !== "") {
            const splitTypes = searchParams.types.split(',').map(type => type.trim()).filter(type => type !== "");
            if (splitTypes.length > 0) {
                entityTypesFromParams = splitTypes;
            }
        }
        const initialFeedParams = {
            feedType: "filtered",
            entityTypes: entityTypesFromParams,
            page: searchParams?.page ? Number(searchParams.page) : 1,
            pageSize: 20,
        };
        // Add search parameters to feed params
        if (searchParams?.category) initialFeedParams.category = searchParams.category;
        if (searchParams?.tags) initialFeedParams.tags = searchParams.tags;
        if (searchParams?.location) initialFeedParams.location = searchParams.location;
        if (searchParams?.sortBy) initialFeedParams.sortBy = searchParams.sortBy;
        if (searchParams?.q) initialFeedParams.query = searchParams.q;
        if (searchParams?.minPrice) initialFeedParams.minPrice = Number(searchParams.minPrice);
        if (searchParams?.maxPrice) initialFeedParams.maxPrice = Number(searchParams.maxPrice);
        // Fetch data server-side
        const [categoriesRes, feedRes] = await Promise.all([
            fetchMainCategories(),
            unifiedSearch(initialFeedParams),
        ]);
        const categories = categoriesRes?.categories ?? [];
        const feedItems = feedRes?.results ?? [];
        const initialHasMore = feedRes?.hasMore ?? false;
        return (
            <div className={styles.pageContainer}>
                <HorizontalFilters />
                <main className={styles.mainContent}>
                    <header className={styles.pageHeader}>
                        <h1 className={styles.pageTitle}>{t('title', { default: 'Filtered Feed' })}</h1>
                        <p className={styles.pageDescription}>
                            {t('description', { 
                                default: 'Browse and filter content across all categories.' 
                            })}
                        </p>
                    </header>
                    <ConnectedFeedDisplay />
                </main>
            </div>
        );
    } catch (error) {
        return notFound();
    }
} 
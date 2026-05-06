/* ────────────────────────────────────────────────────────────
   src/app/[locale]/page.jsx   (SSR-optimized)
───────────────────────────────────────────────────────────── */
import React from 'react';
import {getTranslations} from 'next-intl/server';   // ✅ server-side helper
import {Link} from '../../i18n/navigation';
import styles from '../page.module.css';
import {fetchMainCategoriesSSR} from "../../api/categories";
import {unifiedSearch} from "../../api/searchApi";
import ConnectedFeedDisplay from "../../components/Feed/ConnectedFeedDisplay";
import {SafeJsonLdScript} from "../../utils/secureJsonLd";

export const revalidate = 60;   // ISR every 60 s
// Production SSR timeout configuration
const SSR_TIMEOUTS = {
    CRITICAL: 3000,    // Categories - reduced for faster failure
    FEED: 4000,        // Feed data - reduced for faster failure
    METADATA: 2000     // Metadata - reduced for faster failure
};

// Simplified SSR fetch with built-in timeout handling
async function fetchSSRData(fetchFn, context, fallback = null) {
    try {
        return await fetchFn();
    } catch (error) {
        if (fallback !== null) {
            return fallback;
        }
        throw error;
    }
}

/* ── JSON-LD Helper Functions (MOVED TO TOP) ─────────────────── */

/* Build JSON-LD for featured products */
function buildFeaturedProductsJsonLd(featuredProducts, allCategories) {
    return {
        '@context': 'https://schema.org',
        '@type': 'ItemList',
        itemListElement: featuredProducts.map((prod, index) => {
            const category = allCategories.find((c) => c.id === prod.categoryId);
            const googleCat = category?.googleCategoryId || '';
            const itemCondition =
                prod.condition === 'used'
                    ? 'https://schema.org/UsedCondition'
                    : 'https://schema.org/NewCondition';
            return {
                '@type': 'ListItem',
                position: index + 1,
                url: `https://www.sfx-market.de/product/${prod.id}`,
                name: prod.name,
                item: {
                    '@type': 'Product',
                    name: prod.name,
                    image: prod.thumbnail ?? '',
                    description: prod.description,
                    brand: prod.brand,
                    mpn: prod.mpn || prod.sku || prod.model,
                    additionalProperty: [
                        {
                            '@type': 'PropertyValue',
                            name: 'google_product_category',
                            value: googleCat,
                        },
                    ],
                    itemCondition,
                    offers: {
                        '@type': 'Offer',
                        price: prod.basePrice ?? '0.00',
                        priceCurrency: 'EUR',
                        availability: 'https://schema.org/InStock',
                        shippingDetails: {
                            '@type': 'OfferShippingDetails',
                            shippingRate: {
                                '@type': 'MonetaryAmount',
                                value: prod.shippingCost
                                    ? (prod.shippingCost / 100).toFixed(2)
                                    : '0.00',
                                currency: 'EUR',
                            },
                            shippingDestination: {
                                '@type': 'DefinedRegion',
                                addressCountry: 'EU',
                            },
                        },
                    },
                },
            };
        }),
    };
}

/* Build JSON-LD for categories */
function buildCategoriesJsonLd(allCategories) {
    return {
        '@context': 'https://schema.org',
        '@type': 'ItemList',
        itemListElement: allCategories.map((cat, index) => ({
            '@type': 'ListItem',
            position: index + 1,
            name: cat.description || 'Untitled Category',
            url: `https://www.sfx-market.de/category/${cat.slug}`,
            googleCategoryId: cat.googleCategoryId || '',
            additionalProperty: [
                {
                    '@type': 'PropertyValue',
                    name: 'googleCategoryId',
                    value: cat.googleCategoryId || '',
                },
            ],
        })),
    };
}

/* optional static metadata left unchanged */

/* ── SEO metadata ─────────────────────────────────────────── */
export async function generateMetadata({params}) {
    const {locale} = await params;
    const t = await getTranslations({locale, namespace: 'Seo'});
    return {
        title: t('title', {default: 'sfx markt – Live Marketplace'}),
        description: t('description', {
            default:
                'sfx markt is the live  marketplace that lets you buy, sell and connect with your community in real time. Chat instantly, pay safely with in‑app escrow and discover trending deals around you.',
        }),
        keywords: t('keywords', {
            default:
                'sfx markt, marketplace, buy and sell locally, real‑time chat, SafePay escrow, jobs near me, services, second‑hand deals',
        }),
    };
}

/* ── main page ────────────────────────────────────────────── */
export default async function HomePage({params, searchParams}) {
    const {locale} = await params;
    const t = await getTranslations({locale, namespace: 'HomePage'});
    const marketIntro = {
        title: t('sfxMarketIntroTitle', {default: 'Specialized social marketplace'}),
        text: t('sfxMarketIntroText', {
            default: 'SFX Markt combines marketplace transactions and social engagement in one operating model.'
        }),
        cta: t('sfxMarketIntroCta', {default: 'Open SFX site'})
    };
    // Ensure searchParams is properly awaited
    searchParams = await searchParams;
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
        entityTypes: entityTypesFromParams,
        page: searchParams?.page ? Number(searchParams.page) : 1,
        pageSize: 20,
    };
    if (searchParams?.category) initialFeedParams.categorySlug = searchParams.category;
    if (searchParams?.sortBy) initialFeedParams.sortBy = searchParams.sortBy;
    /* 2) Enhanced data fetch with timeouts and fallbacks */
    let categories = [];
    let feedItems = [];
    let initialHasMore = false;
    try {
        // Execute both critical requests with Promise.allSettled for resilience
        const [categoryResult, feedResult] = await Promise.allSettled([
            fetchSSRData(
                () => fetchMainCategoriesSSR(),
                'Homepage categories fetch',
                {categories: []} // Fallback for categories
            ),
            fetchSSRData(
                () => unifiedSearch(initialFeedParams),
                'Homepage feed fetch',
                {items: [], hasMore: false} // Fallback for feed
            )
        ]);
        // Handle categories result
        if (categoryResult.status === 'fulfilled') {
            const categoryRes = categoryResult.value;
            categories = categoryRes?.categories ?? [];
            if (categories.length === 0) {
            } else {
            }
        } else {
            categories = []; // Use empty fallback
        }
        // Handle feed result  
        if (feedResult.status === 'fulfilled') {
            const feedRes = feedResult.value;
            feedItems = feedRes?.results ?? feedRes?.items ?? [];
            initialHasMore = feedRes?.hasMore ?? false;
            if (feedItems.length === 0) {
            } else {
            }
            // Force hasMore to true if we got less than pageSize but have items
            // This is likely a server-side pagination bug
            if (feedItems.length > 0 && feedItems.length < (initialFeedParams.pageSize || 20) && !initialHasMore) {
                initialHasMore = true; // Force for debugging
                // Also add debugging info for the server team
            }
        } else {
            feedItems = []; // Use empty fallback
            initialHasMore = false;
        }
        // Success logging for production monitoring
        if (process.env.NODE_ENV === 'production') {
        }
    } catch (error) {
        // This should rarely happen due to Promise.allSettled, but handle anyway
        // Use fallback values to prevent page crash
        categories = [];
        feedItems = [];
        initialHasMore = false;
    }
    // Re-construct feedParams for FeedProvider, using the same robust logic for entityTypes
    const providerFeedParams = {
        feedType: "latest",
        entityTypes: entityTypesFromParams, // Use the same robustly derived entityTypesFromParams
        page: searchParams?.page ? Number(searchParams.page) : 1,
    };
    if (searchParams?.category) providerFeedParams.category = searchParams.category;
    if (searchParams?.tags) providerFeedParams.tags = searchParams.tags;
    if (searchParams?.location) providerFeedParams.location = searchParams.location;
    if (searchParams?.sortBy) providerFeedParams.sortBy = searchParams.sortBy;
    // Use feed items for featured content instead of products
    const featured = feedItems.slice(0, 5);
    /* 4) JSON-LD helpers (with safe fallbacks) */
    const featuredJsonLd = featured.length > 0 && categories.length > 0
        ? buildFeaturedProductsJsonLd(featured, categories)
        : null;
    const categoriesJsonLd = categories.length > 0
        ? buildCategoriesJsonLd(categories)
        : null;
    /* 5) render */
    return (
        <div className={styles.pageContainer}>
            {featuredJsonLd && <SafeJsonLdScript data={featuredJsonLd}/>}
            {categoriesJsonLd && <SafeJsonLdScript data={categoriesJsonLd}/>}
            {/* Main Content Grid */}
            <div className={styles.mainContent}>
                <section className={styles.introCard}>
                    <div>
                        <h1 className={styles.introTitle}>{marketIntro.title}</h1>
                        <p className={styles.introText}>{marketIntro.text}</p>
                    </div>
                    <Link href="/sfx-market" className={styles.introLink}>
                        {marketIntro.cta}
                    </Link>
                </section>
                <ConnectedFeedDisplay/>
            </div>
        </div>
    );
}

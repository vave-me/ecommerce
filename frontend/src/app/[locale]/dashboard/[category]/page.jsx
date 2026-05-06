import React from 'react';
import {searchPostsWithCategory} from '../../../../api/searchApi';
import DashboardPageClient from '../DashboardPage.client'; // Updated to use DashboardPage.client
import {notFound} from 'next/navigation';
import {getTranslations} from 'next-intl/server';
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
export const dynamic = 'force-dynamic';
// Helper function (assuming it's defined elsewhere or keep here)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
export async function generateMetadata(props) {
    let t;
    let category;
    try {
        const {params} = await resolveProps(props);
        const {locale} = params;
        category = params.category; // Assign category here
        t = await getTranslations({locale, namespace: 'DashboardPage'}); // Changed from PostsPage to DashboardPage
        // Try to get translated category name, fallback to raw category slug
        const categoryName = category;
        return {
            title: t('metaTitleCategory', {categoryName}),
            description: t('metaDescriptionCategory', {categoryName}),
            openGraph: {
                title: t('metaTitleCategory', {categoryName}),
                description: t('metaDescriptionCategory', {categoryName}),
                type: 'website',
            },
            alternates: {
                canonical: `https://www.sfx-market.de/${locale}/dashboard/${category}`,
            },
        };
    } catch (e) {
        // Fallback metadata if anything fails
        const fallbackCategoryName = category || 'Category'; // Use category if available
        const fallbackTitle = t ? t('metaTitleCategoryFallback', {categoryName: fallbackCategoryName}) : `${fallbackCategoryName} Dashboard`;
        const fallbackDescription = t ? t('metaDescriptionCategoryFallback', {categoryName: fallbackCategoryName}) : `Browse dashboard content in the ${fallbackCategoryName} category.`;
        return {
            title: fallbackTitle,
            description: fallbackDescription,
        };
    }
}
export default async function DashboardCategoryPage(props) {
    let t; // For potential use in catch block
    let category; // For potential use in catch block
    try {
        const {params, searchParams} = await resolveProps(props);
        const {locale} = params;
        category = params.category; // Assign category here
        /* 1) i18n helper */
        t = await getTranslations({locale, namespace: 'DashboardPage'}); // Changed from PostsPage to DashboardPage
        // Validate category early
        if (!category) {
             // Use translated log
            notFound();
        }
        // Get translated category name for use in UI/fetch errors if needed later
        const categoryName = category;
        /* 2) Extract searchParams values */
        const displayMode = searchParams?.displayMode || 'grid';
        const page = parseInt(searchParams?.page || '1', 10);
        const limit = 20;
        const sortBy = searchParams?.sortBy || '';
        /* 3) Create listing filters including category */
        const listingFilters = {
            category, // Pass the category slug to the API
            displayMode,
            page,
            limit,
            sortBy,
        };
        /* 4) Fetch posts filtered by category */
        let fetchedData = {posts: [], totalCount: 0, totalPages: 0, currentPage: 1};
        let fetchError = null;
        try {
            // Changed from getPosts to getPostsByCategoryId
            fetchedData = await searchPostsWithCategory(category, {
                page,
                limit,
                sortBy,
                displayMode
            });
        } catch (err) {
            // Use translated error, passing category name for context
            fetchError = t("errorFetchingPostsCategory", {category: categoryName});
        }
        /* 5) Process the fetched data */
        const serverPosts = fetchedData?.posts || [];
        const totalPages = parseInt(fetchedData?.totalPages || '0', 10);
        const totalCount = parseInt(fetchedData?.totalCount || '0', 10);
        const currentPage = parseInt(fetchedData?.currentPage || '1', 10);
        /* 6) Construct JSON-LD for SEO (Using translated names/descriptions) */
        const postsJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': t('jsonLdListNameCategory', {categoryName}),
            'description': t('jsonLdListDescriptionCategory', {categoryName}),
            itemListElement: serverPosts.map((post, index) => ({
                '@type': 'ListItem',
                position: (currentPage - 1) * limit + index + 1, // Calculate position
                url: `https://www.sfx-market.de/${locale}/dashboard/${category}/${post.slug || post.id}`, // Changed from /posts/ to /dashboard/
                item: { // Embed BlogPosting data
                    '@type': 'BlogPosting',
                    'headline': post.title || post.name || t('postCard.unnamedPost'),
                    'image': post.thumbnail || post.image || '',
                    'description': post.description || post.summary || '',
                    'author': {
                        '@type': 'Person',
                        'name': post.author?.name || t('postCard.defaultAuthor')
                    },
                    'publisher': {
                        '@type': 'Organization',
                        'name': t('common.publisherName'),
                        'logo': {'@type': 'ImageObject', 'url': 'https://www.sfx-market.de/logo.png'}
                    },
                    'datePublished': post.publishedDate || post.createdAt,
                    'dateModified': post.updatedAt || post.publishedDate || post.createdAt,
                    'mainEntityOfPage': {
                        '@type': 'WebPage',
                        '@id': `https://www.sfx-market.de/${locale}/dashboard/${category}/${post.slug || post.id}` // Changed from /posts/ to /dashboard/
                    }
                }
            })),
        };
        /* 7) Return UI components */
        return (
            <>
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(postsJsonLd)}}
                />
                <DashboardPageClient
                    serverPosts={serverPosts}
                    serverFilters={listingFilters}
                    fetchError={fetchError}
                    currentCategory={category} // Pass category slug
                    // Pass translated labels needed by the client
                    labels={{
                        pageTitle: t("titleCategory", {categoryName}), // Pass translated category title
                        emptyMsg: t("emptyCategory", {categoryName}), // Pass translated empty message
                        // Add other labels if DashboardPageClient needs them
                    }}
                    totalPages={totalPages}
                    currentPage={currentPage}
                    totalCount={totalCount}
                />
            </>
        );
    } catch (e) {
        // Attempt to get translations for error page
        let errorTitle = "Error Loading Category Dashboard";
        let errorMessage = "Could not load dashboard content for this category. Please try again later.";
        try {
            if (!t) t = await getTranslations({locale: props.params?.locale || 'en', namespace: "DashboardPage"});
            const categoryName = category ? t(`category_${category}_label`, {}, {defaultValue: category}) : t('fallbackCategoryName');
            errorTitle = t('globalErrorTitleCategory', {categoryName});
            errorMessage = t('globalErrorMessageCategory', {categoryName});
        } catch (transErr) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', transErr);
        }
    }
        return (
            <div className="container mx-auto p-8 text-center">
                <h1 className="text-2xl font-bold mb-4">{errorTitle}</h1>
                <p>{errorMessage}</p>
            </div>
        );
    }
}
import React from 'react';
import {getTranslations} from "next-intl/server";
import {notFound} from "next/navigation"; // Keep if needed
import {searchPostsWithFilters} from "../../../api/searchApi"; // Changed from postsApi to searchApi
import PostsPageClient from "./PostsPage.client"; // Adjust path
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
export const dynamic = "force-dynamic";
// Helper function (assuming it's defined elsewhere or keep here)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
export async function generateMetadata(props) {
    let t;
    try {
        const {params} = await resolveProps(props);
        const {locale} = params;
        t = await getTranslations({locale, namespace: "NewsPage"});
        return {
            title: t("metaTitle"),
            description: t("metaDescription"),
            openGraph: {
                title: t("metaTitle"),
                description: t("metaDescription"),
                type: 'website',
            },
            alternates: {
                canonical: `https://www.sfx-market.de/${locale}/news`,
            },
        };
    } catch (e) {
        // Fallback if translation fails
        return {
            title: "News & Updates", // Hardcoded fallback
            description: "Browse the latest news and updates.", // Hardcoded fallback
        };
    }
}
export default async function NewsIndexPage(props) {
    let t; // For potential use in catch block
    try {
        const {params, searchParams} = await resolveProps(props);
        const {locale} = params;
        /* 1) i18n helper */
        t = await getTranslations({locale, namespace: "NewsPage"});
        /* 2) optional filters */
        const displayMode = searchParams?.displayMode || "list";
        const page = parseInt(searchParams?.page || "1", 10);
        const limit = 20; // Example limit
        const sortBy = searchParams?.sortBy || ''; // Example sort
        const listingFilters = {
            displayMode,
            page,
            limit,
            sortBy,
            type: 'news' // Filter specifically for news type posts
            // Add other filters as needed
        };
        /* 3) fetch News Posts */
        let postsData = {posts: [], totalCount: 0, totalPages: 0, currentPage: 1}; // Changed variable name & added currentPage default
        let fetchError = null;
        try {
            // Changed from getPosts to searchPostsWithFilters
            postsData = await searchPostsWithFilters(listingFilters);
        } catch (err) {
            // Use translated error message
            fetchError = t("errorFetchingNews");
        }
        // Use nullish coalescing for safety
        const serverPosts = postsData.posts ?? [];
        const totalPages = parseInt(postsData?.totalPages || '0', 10);
        const totalCount = parseInt(postsData?.totalCount || '0', 10);
        const currentPage = parseInt(postsData?.currentPage || '1', 10);
        /* 4) JSON-LD for SEO (Corrected Type, added translated name/desc) */
        const jsonLd = {
            "@context": "https://schema.org",
            "@type": "ItemList",
            // Add list name and description
            "name": t("itemListTitle"),
            "description": t("itemListDescription"),
            itemListElement: serverPosts.map((p, i) => ({
                "@type": "ListItem",
                position: (currentPage - 1) * limit + i + 1, // Calculate position based on page and limit
                url: `https://www.sfx-market.de/${locale}/news/${p.category?.slug || 'general'}/${p.slug || p.id}`, // Use /news/ URLs
                // Use item for BlogPosting data
                item: {
                    "@type": "BlogPosting", // Corrected type
                    "headline": p.title || p.name || t('newsCard.unnamedNews'), // Use title/name or fallback
                    "image": p.thumbnail || p.image || '', // Prioritize thumbnail/image
                    "description": p.description || p.summary || '', // Use description/summary
                    "author": {
                        "@type": "Person",
                        "name": p.author?.name || t('newsCard.defaultAuthor') // Use author name or fallback
                    },
                    "publisher": {
                        "@type": "Organization",
                        "name": t('common.publisherName'), // Use common publisher name
                        "logo": {
                            "@type": "ImageObject",
                            "url": `https://www.sfx-market.de/logo.png` // Define logo URL
                        }
                    },
                    "datePublished": p.publishedDate || p.createdAt,
                    "dateModified": p.updatedAt || p.publishedDate || p.createdAt,
                    "mainEntityOfPage": {
                        "@type": "WebPage",
                        "@id": `https://www.sfx-market.de/${locale}/news/${p.category?.slug || 'general'}/${p.slug || p.id}`
                    }
                    // Removed Product specific fields (brand, mpn, offers, itemCondition)
                }
            }))
        };
        /* 5) render */
        return (
            <>
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(jsonLd)}}
                />
                <PostsPageClient
                    serverPosts={serverPosts}
                    serverFilters={listingFilters}
                    // Pass translated strings needed by the client
                    labels={{
                        pageTitle: t("title"),
                        emptyMsg: t("empty"),
                        errorTitle: t("errorTitle"),
                        errorMessage: t("errorMessage"),
                        retryButton: t("retryButton"),
                        // Add other labels if PostsPageClient needs them
                    }}
                    fetchError={fetchError}
                    totalPages={totalPages}
                    totalCount={totalCount}
                    currentPage={currentPage}
                />
            </>
        );
    } catch (e) {
        // Attempt to get translations for error page
        let errorTitle = "Error Loading News";
        let errorMessage = "Could not load news at this time. Please try again later.";
        try {
            if (!t) t = await getTranslations({locale: props.params?.locale || 'en', namespace: "NewsPage"});
            errorTitle = t('globalErrorTitle');
            errorMessage = t('globalErrorMessage');
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
import React from 'react';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server'; // Import getTranslations
import { getPost } from "../../../../../api/searchApi"; // Changed from getPostsBySlug to getPost
import PostCard from "./PostCard"; // Fix import path to use local PostCard
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
// Helper function (assuming it's defined elsewhere or keep here)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return { params, searchParams };
}
// Generate dynamic metadata for SEO
export async function generateMetadata(props) {
    let t;
    let slug;
    let post = null;
    try {
        const { params } = await resolveProps(props);
        const { locale } = params;
        slug = params.slug; // Assign slug here
        t = await getTranslations({ locale, namespace: 'PostDetailPage' }); // Use PostDetailPage namespace
        post = await getPost(slug);
        if (!post) {
            return { title: t('metaTitleNotFound') }; // Translated "Not Found" title
        }
        const title = post.title || post.name || t('metaTitleDefault');
        const description = post.description || post.summary || post.excerpt || t('metaDescriptionDefault', { title }); // Use title in fallback desc
        return {
            title: title,
            description: description,
            openGraph: {
                title: title,
                description: description,
                type: 'article',
                publishedTime: post.publishedDate || post.createdAt,
                modifiedTime: post.updatedAt || post.publishedDate || post.createdAt,
                section: params.category, // Keep category context if available
                authors: [post.author?.name || t('ogDefaultAuthor')], // Translated fallback author
                // images: [post.coverImage || post.thumbnail || ''], // Add image URL
            },
        };
    } catch (err) {
        // Fallback metadata on error
        const errorTitle = t ? t('metaTitleError') : "Error Loading Post";
        return { title: errorTitle };
    }
}
// Define the Page component (Server Component)
export default async function PostDetailPage(props) {
    let t; // For potential use in catch block
    let slug; // For potential use in catch block
    let category; // For potential use in catch block
    try {
        const { params } = await resolveProps(props);
        const { locale } = params;
        slug = params.slug; // Assign slug
        category = params.category; // Assign category
        t = await getTranslations({ locale, namespace: 'PostDetailPage' }); // Use PostDetailPage namespace
        let postData = null;
        try {
            postData = await getPost(slug);
            if (!postData) {
                notFound(); // Trigger 404
            }
        } catch (err) {
            notFound(); // Trigger 404 on fetch error
        }
        // Generate JSON-LD for the article (Using translations)
        const articleJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'BlogPosting',
            'headline': postData.title || postData.name || t('postCard.unnamedPost'), // Use PostCard fallback key
            'description': postData.description || postData.summary || postData.excerpt || '',
            'image': postData.coverImage || postData.thumbnail || postData.image || '',
            'datePublished': postData.publishedDate || postData.createdAt,
            'dateModified': postData.updatedAt || postData.publishedDate || postData.createdAt,
            'author': {
                '@type': 'Person',
                'name': postData.author?.name || t('postCard.defaultAuthor') // Use PostCard fallback key
            },
            'publisher': {
                '@type': 'Organization',
                'name': t('common.publisherName'), // Use common publisher name
                'logo': {
                    '@type': 'ImageObject',
                    'url': 'https://www.sfx-market.de/logo.png' // Define logo URL
                }
            },
            'mainEntityOfPage': {
                '@type': 'WebPage',
                '@id': `https://www.sfx-market.de/${locale}/posts/${category}/${slug}` // Use locale, category, slug
            }
        };
        // Assume PostCard is the component responsible for displaying the full post content
        // We pass the fetched post data and potentially some translated labels if needed
        return (
            <div className="container mx-auto p-4"> {/* Example container */}
                {/* Add JSON-LD */}
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(articleJsonLd)}}
                />
                {/* Render the PostCard (or a dedicated detail component) */}
                {/* PostCard needs to be adapted to accept `post` prop */}
                <PostCard post={postData} locale={locale} />
            </div>
        );
    } catch (e) {
        // Global error handler for the page component itself
        // Attempt to get translations for error page
        let errorTitle = "Error Loading Post";
        let errorMessage = "Could not load the post details. Please try again later.";
        try {
            if (!t) t = await getTranslations({ locale: props.params?.locale || 'en', namespace: "PostDetailPage" });
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
import React from 'react';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server';
import { getPost } from "../../../../api/searchApi";
import DetailedPostView from "../../../../components/posts/DetailedPostView";
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';

// Helper function to resolve props
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return { params, searchParams };
}

// Generate dynamic metadata for SEO
export async function generateMetadata(props) {
    let t;
    let id;
    let post = null;
    
    try {
        const { params } = await resolveProps(props);
        const { locale } = params;
        id = params.id;
        t = await getTranslations({ locale, namespace: 'PostDetailPage' });
        
        post = await getPost(id);
        if (!post) {
            return { title: t('metaTitleNotFound') };
        }
        
        const title = post.title || post.name || t('metaTitleDefault');
        const description = post.description || post.summary || post.excerpt || t('metaDescriptionDefault', { title });
        
        return {
            title: title,
            description: description,
            openGraph: {
                title: title,
                description: description,
                type: 'article',
                publishedTime: post.publishedDate || post.createdAt,
                modifiedTime: post.updatedAt || post.publishedDate || post.createdAt,
                authors: [post.author?.name || t('ogDefaultAuthor')],
                images: post.image || post.thumbnail || post.coverImage ? [{
                    url: post.image || post.thumbnail || post.coverImage,
                    width: 1200,
                    height: 630,
                    alt: title
                }] : [],
            },
            twitter: {
                card: 'summary_large_image',
                title: title,
                description: description,
                images: post.image || post.thumbnail || post.coverImage ? [post.image || post.thumbnail || post.coverImage] : [],
            }
        };
    } catch (err) {
        const errorTitle = t ? t('metaTitleError') : "Error Loading Post";
        return { title: errorTitle };
    }
}

// Page component for posts without categories
export default async function PostDetailWithoutCategoryPage(props) {
    let t;
    let id;
    
    try {
        const { params } = await resolveProps(props);
        const { locale } = params;
        id = params.id;
        t = await getTranslations({ locale, namespace: 'PostDetailPage' });
        
        let postData = null;
        try {
            postData = await getPost(id);
            if (!postData) {
                notFound();
            }
        } catch (err) {
            // Error: 'Error fetching post:', err...
            notFound();
        }
        
        // Generate JSON-LD for the article
        const articleJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'BlogPosting',
            'headline': postData.title || postData.name || t('postCard.unnamedPost'),
            'description': postData.description || postData.summary || postData.excerpt || '',
            'image': postData.image || postData.thumbnail || postData.coverImage || '',
            'datePublished': postData.publishedDate || postData.createdAt,
            'dateModified': postData.updatedAt || postData.publishedDate || postData.createdAt,
            'author': {
                '@type': 'Person',
                'name': postData.author?.name || t('postCard.defaultAuthor')
            },
            'publisher': {
                '@type': 'Organization',
                'name': t('common.publisherName'),
                'logo': {
                    '@type': 'ImageObject',
                    'url': 'https://www.sfx-market.de/logo.png'
                }
            },
            'mainEntityOfPage': {
                '@type': 'WebPage',
                '@id': `https://www.sfx-market.de/${locale}/post/${id}`
            }
        };
        
        return (
            <>
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(articleJsonLd)}}
                />
                <DetailedPostView 
                    post={postData} 
                    locale={locale}
                />
            </>
        );
    } catch (e) {
        // Error: 'Error in PostDetailWithoutCategoryPage:', e...
        
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
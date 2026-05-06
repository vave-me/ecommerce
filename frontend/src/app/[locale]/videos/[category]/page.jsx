// File: app/videos/[itemId]/page.jsx
// No "use client" here => It's a Server Component.

import {getAllItemVideos} from '../../../../api/mediaApi';
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';

// Server-only exports
export const revalidate = 60;

export const metadata = {
    title: 'Short Video Feed',
    description: 'A vertical feed of short videos for a specific itemId.',
    openGraph: {
        title: 'Short Video Feed by Item ID',
        description: 'A vertical feed of short videos filtered by a specific itemId.',
    },
    twitter: {
        card: 'summary_large_image',
        title: 'Short Video Feed by Item ID',
        description: 'View short videos filtered by itemId.',
    },
};

/* ------------------------------------------
   SERVER-SIDE FETCH: getVideos
   ------------------------------------------ */
async function getVideos(itemId) {
    try {
        const res = await getAllItemVideos(itemId);
        if (!res.ok) {
            throw new Error(`Failed to fetch videos for itemId=${itemId}`);
        }
        const data = await res.json();
        return data.videos || [];
    } catch (error) {
        // Error: 'Error fetching item videos:', error...
        throw error;
    }
}

/* ------------------------------------------
   MAIN SERVER COMPONENT
   ------------------------------------------ */
export default async function VideosPage({params}) {
    const {itemId} = params;

    let videos = [];
    let errorMessage = '';

    // Attempt to fetch videos from your API
    try {
        videos = await getVideos(itemId);
    } catch (err) {
        errorMessage = 'Error fetching videos. Please try again later.';
    }

    // Build JSON-LD for SEO (schema.org)
    const videosJsonLd = {
        '@context': 'https://schema.org',
        '@type': 'ItemList',
        itemListElement: videos.map((vid, index) => ({
            '@type': 'ListItem',
            position: index + 1,
            url: vid.url,
            name: vid.metadata || `Video #${index + 1}`,
            item: {
                '@type': 'VideoObject',
                name: vid.metadata || `Video #${index + 1}`,
                description: vid.metadata || 'No description',
                contentUrl: vid.url,
                thumbnailUrl: vid.thumbnail || '',
                uploadDate: new Date().toISOString(),
            },
        })),
    };

    // Dynamically import the Client Component (or you can do a static import)
    const {default: ItemVideosPageClient} = await import('./ItemVideosPageClient.jsx');

    return (
        <>
            {/* Inject JSON-LD metadata for search engines */}
            <script
                type="application/ld+json"
                dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(videosJsonLd)}}
            />

            {/* Render the Client Component, passing in data and error */}
            <ItemVideosPageClient
                itemId={itemId}
                videos={videos}
                errorMessage={errorMessage}
            />
        </>
    );
}

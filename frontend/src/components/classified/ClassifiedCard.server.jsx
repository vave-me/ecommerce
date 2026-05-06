import React from 'react';
import { getMediaByItem, isMediaResponseSuccess } from '../../api/mediaApi';
import ClassifiedCard from './ClassifiedCard';
/**
 * ClassifiedCardServer - Server component wrapper for ClassifiedCard
 * 
 * This component:
 * - Pre-fetches media data on the server side
 * - Provides optimized media loading for better performance
 * - Handles errors gracefully with fallback behavior
 * - Follows the same pattern as DealCard.server.jsx and AutomotiveCard.server.jsx
 * 
 * @param {Object} props - Component props
 * @param {Object} props.product - Product data object
 * @param {string} props.className - Optional CSS class name
 * @returns {JSX.Element} ClassifiedCard with pre-loaded media
 */
const ClassifiedCardServer = async ({ product, className = "" }) => {
    let preloadedMedia = null;
    let hasPreloadedMedia = false;
    // Debug logging for server-side processing
    // Attempt to pre-fetch media if we have a product ID
    if (product?.id) {
        try {
            const mediaResponse = await getMediaByItem(product.id);
            if (isMediaResponseSuccess(mediaResponse) && mediaResponse.media?.mediaOrder?.length > 0) {
                // Format media data for client component
                preloadedMedia = mediaResponse.media.mediaOrder.map(item => ({
                    id: item.mediaItemId || item.id,
                    url: item.url,
                    type: item.type || 'image',
                    alt: item.altText || product.name || 'Product image',
                    caption: item.caption || '',
                    order: item.order || 0
                }));
                hasPreloadedMedia = true;
            } else {
                // Check if product has direct media references
                const directMedia = [];
                if (product.thumbnail) {
                    directMedia.push({
                        id: 'thumbnail',
                        url: product.thumbnail,
                        type: 'image',
                        alt: product.name || 'Product image',
                        order: 0
                    });
                }
                if (product.images && Array.isArray(product.images)) {
                    product.images.forEach((imageUrl, index) => {
                        if (imageUrl && imageUrl !== product.thumbnail) {
                            directMedia.push({
                                id: `image-${index}`,
                                url: imageUrl,
                                type: 'image',
                                alt: product.name || 'Product image',
                                order: index + 1
                            });
                        }
                    });
                }
                if (product.media && Array.isArray(product.media)) {
                    product.media.forEach((mediaItem, index) => {
                        if (mediaItem && mediaItem.url) {
                            directMedia.push({
                                id: mediaItem.id || `media-${index}`,
                                url: mediaItem.url,
                                type: mediaItem.type || 'image',
                                alt: mediaItem.alt || product.name || 'Product image',
                                order: mediaItem.order || index + 10
                            });
                        }
                    });
                }
                if (directMedia.length > 0) {
                    // Sort by order and use as preloaded media
                    preloadedMedia = directMedia.sort((a, b) => a.order - b.order);
                    hasPreloadedMedia = true;
                }
            }
        } catch (error) {
            // Don't throw - gracefully degrade to client-side loading
            preloadedMedia = null;
            hasPreloadedMedia = false;
        }
    } else {
        }
    // Final debug log before rendering
    // Render the client component with preloaded data
    return (
        <ClassifiedCard
            product={product}
            preloadedMedia={preloadedMedia}
            hasPreloadedMedia={hasPreloadedMedia}
            className={className}
        />
    );
};
export default ClassifiedCardServer; 
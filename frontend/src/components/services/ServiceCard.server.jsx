// Server Component for ServiceCard
import React from 'react';
import ServiceCard from './ServiceCard';
import {getService} from '../../api/searchApi';
import {getMediaByItem} from "@/api/mediaApi";

/**
 * ServiceCardServer - Server component wrapper for ServiceCard
 * Handles server-side data fetching and passes data to client component
 *
 * @param {Object} props - Component props
 * @param {string} props.serviceId - Service ID to fetch
 * @param {Object} props.service - Pre-fetched service data (optional)
 * @param {string} props.className - Additional CSS classes
 * @returns {JSX.Element} Rendered component
 */
export default async function ServiceCardServer({serviceId, service, className, ...props}) {
    let serviceData = service;

    // If no service data provided but serviceId is available, fetch it
    if (!serviceData && serviceId) {
        try {
            serviceData = await getService(serviceId);
        } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle

            // Handle authentication errors gracefully - return null to skip server component
            if (error.status === 401 || error.response?.status === 401) {
                
                return null; // Let client component handle the fetch
            }

            // For other errors, return error state
            return (
                <div className={`service-card-error ${className || ''}`}>
                    <p>Failed to load service</p>
                </div>
            );
        }
    }

    // If still no data, return empty state
    if (!serviceData) {
        return (
            <div className={`service-card-empty ${className || ''}`}>
                <p>No service data available</p>
            </div>
        );
    }

    // Pre-fetch media on server for better performance
    let mediaData = null;
    if (serviceData?.id) {
        try {
            const mediaResponse = await getMediaByItem(serviceData.id);

            if (mediaResponse?.media?.mediaOrder?.length > 0) {
                mediaData = mediaResponse.media.mediaOrder.map(item => {
                    // Ensure all values are serializable strings/primitives
                    return {
                        id: String(item.mediaItemId || item.id || ''),
                        url: String(item.url || ''),
                        type: String(item.type || 'image'),
                        alt: String(item.altText || serviceData.name || 'Service image')
                    };
                });

                // Verify serialization works
                try {
                    const serialized = JSON.stringify(mediaData);
                    const deserialized = JSON.parse(serialized);
                } catch (e) {
                    // Error: `[ServiceCardServer] Serialization failed:`, e...
                    mediaData = null; // Reset if can't serialize
                }
            } else {
                // Check if service has any image properties
                const serviceImageProps = {
                    thumbnail: serviceData.thumbnail,
                    images: serviceData.images,
                    media: serviceData.media,
                    image: serviceData.image,
                    imageUrl: serviceData.imageUrl
                };
            }
        } catch (error) {
            // Error: `[ServiceCardServer] Error fetching media for serv...
            // Continue without media - client will handle fallback
        }
    }

    // Pass data to client component
    return <ServiceCard service={serviceData} preloadedMedia={mediaData} className={className} {...props} />;
} 
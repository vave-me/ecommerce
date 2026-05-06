// Server Component for FeedItem with optimized media loading
import React from 'react';
import { DealCardServer } from '../deals';
import { getMediaByItem } from '../../api/mediaApi';

// Import other card components (these should ideally have server versions too)
import ImprovedClassifiedCard from '../Items/ProductListItem.client';
import ImprovedServiceCard from '../Items/ServiceListItem.client';
import { PropertyCard } from '../property';
import VehicleCard from '../Items/CarListItem.client';
import PostCard from '../Items/PostListItem.server';
import { JobCard } from '../jobs';

/**
 * FeedItemServer - Server component for feed items with optimized media loading
 * Pre-fetches media on the server for better performance
 * 
 * @param {Object} props - Component props
 * @param {Object} props.item - Feed item data
 * @param {boolean} props.preloadMedia - Whether to preload media on server
 * @returns {JSX.Element} Rendered component
 */
export default async function FeedItemServer({ item, preloadMedia = true }) {
    if (!item || !item.entityType) {
        return null;
    }

    // For deals, use the optimized server component with media preloading
    if (item.entityType === 'deal') {
        // Extract the actual deal data from the unified feed structure
        const dealData = item.deal || item;
        return (
            <div className="feed-item">
                <DealCardServer deal={dealData} />
            </div>
        );
    }

    // For other entity types, we'll use client components for now
    // TODO: Create server versions of other card components
    const renderClientComponent = () => {
        switch (item.entityType) {
            case 'product':
                return <ImprovedClassifiedCard product={item} />;
            case 'service':
                return <ImprovedServiceCard service={item} />;
            case 'job':
                return <JobCard job={item} />;
            case 'property':
                return <PropertyCard property={item} />;
            case 'vehicle':
                return <VehicleCard vehicle={item} />;
            case 'post':
                return <PostCard post={item} />;
            default:
                return (
                    <div className="unsupported-type">
                        <p>Unsupported content type: {item.entityType}</p>
                    </div>
                );
        }
    };

    return (
        <div className="feed-item">
            {renderClientComponent()}
        </div>
    );
} 
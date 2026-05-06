"use client";

import React, {memo, useCallback, useState, useEffect} from 'react';
import {useRouter} from 'next/navigation';
import {useTranslations} from 'next-intl';
import styles from './FeedItem.module.css';
import {PostCard} from '../PostCard';
import TweetCard from '../../app/[locale]/design/tweets/page';
import ServiceCard from '../services/ServiceCard'; // Following DealCard pattern
import {ClassifiedCard} from "../classified";
import ImprovedVideoCard from "../../app/[locale]/design/videos/page";
import ImprovedCompactVideo from "../../app/[locale]/design/shorts/page";

import ImprovedPostCard from "@/app/[locale]/design/posts/page";

const FeedItem = memo(({item, onLike, onComment, onShare, isVisible, isMobile}) => {
    const router = useRouter();
    const t = useTranslations('Feed');
    const [errorState, setErrorState] = useState(null);

    // Validate item on mount to avoid early returns in hook execution paths
    useEffect(() => {
        if (!item || !item.entityType) {
            // Error: "Invalid item or missing entityType", item...
            setErrorState({
                type: 'invalid',
                message: 'Invalid feed item data'
            });
        }
    }, [item]);

    // Helper function to get the entity-specific data and ID
    const getEntityData = useCallback(() => {
        if (!item || !item.entityType) {
            return {data: null, id: null, slug: null};
        }

        // Debug the item structure to understand categorySlug location

        // Try to extract categorySlug from multiple possible locations
        const extractCategorySlug = (itemData) => {
            return itemData.categorySlug ||
                itemData.category_slug ||
                itemData.category?.slug ||
                itemData.category?.categorySlug ||
                null;
        };

        // For post and deal types, they might be directly in the item rather than nested
        if ((item.entityType === 'post' || item.entityType === 'deal') && !item[item.entityType]) {
            const slug = extractCategorySlug(item);

            return {
                data: item,
                id: item.id,
                slug: slug,
            };
        }

        // Check if the entity-specific data exists
        const entityData = item[item.entityType];
        if (!entityData) {
            // Try to get data from item directly as fallback
            const slug = extractCategorySlug(item);

            return {
                data: item,
                id: item.id,
                slug: slug,
            };
        }

        // Extract from nested entity data with fallback to item level
        const slug = extractCategorySlug(entityData) || extractCategorySlug(item);

        return {
            data: entityData,
            id: entityData.id,
            slug: slug,
        };
    }, [item]);

    const {data, id, slug} = getEntityData();

    const handleItemClick = useCallback(() => {
        if (!id) {
            
            return;
        }

        // 🔍 LOG: Show routing data for debugging
        switch (item?.entityType) {
            case 'post':
                // If post has no category, use direct post route
                if (!slug) {
                    router.push(`/post/${id}`);
                } else {
                    router.push(`/posts/${slug}/${id}`);
                }
                break;
            case 'product':
                router.push(`/products/${slug}/${id}`);
                break;
            case 'service':
                router.push(`/services/${slug}/${id}`);
                break;
            case 'video':
                router.push(`/videos/${slug}/${id}`);
                break;

            default:
                
                break;
        }
    }, [router, item, id, slug]);

    const renderContent = useCallback(() => {
        // Special handling for deals since they might not be nested
        if (item?.entityType === 'deal') {
            try {
                // Extract the actual deal data from the unified feed structure
                const dealData = item.deal || item;
                return <DealCard deal={dealData}/>;
            } catch (error) {
                // Error: '[FeedItem] Error rendering deal:', error...
                return (
                    <div className={styles.errorState}>
                        <p>{t('renderError', {type: 'deal', message: error.message})}</p>
                    </div>
                );
            }
        }

        // Special handling for posts similar to deals
        if (item?.entityType === 'post') {
            try {
                // PostCard now handles both regular posts and single video posts
                return <PostCard post={item} isVisible={isVisible} isMobile={isMobile}/>;
            } catch (error) {
                // Error: '[FeedItem] Error rendering post:', error...
                return (
                    <div className={styles.errorState}>
                        <p>{t('renderError', {type: 'post', message: error.message})}</p>
                    </div>
                );
            }
        }

        // For other entity types, check for nested data first
        if (!data) {
            return (
                <div className={styles.unsupportedType}>
                    <p>{t('unsupportedContentType', {type: item?.entityType || 'unknown'})}</p>
                </div>
            );
        }

        try {
            // Create props object with the correct property name based on entityType
            const cardProps = {};

            switch (item.entityType) {
                case 'product':
                    cardProps.product = data;
                    return <ClassifiedCard {...cardProps} />;
                case 'service':
                    cardProps.service = data;
                    return <ServiceCard {...cardProps} />;
                case 'tweet':
                    cardProps.tweet = data;
                    return <TweetCard {...cardProps} />;
                case "short":
                    cardProps.short = data;
                    return <ImprovedCompactVideo {...cardProps} />;
                case "video":
                    cardProps.video = data;
                    return <ImprovedVideoCard {...cardProps} />;
                case 'post':
                    cardProps.post = data;
                    return <PostCard {...cardProps} isVisible={isVisible} isMobile={isMobile} />;
                default:
                    return (
                        <div className={styles.unsupportedType}>
                            <p>{t('unsupportedContentType', {type: item?.entityType || 'unknown'})}</p>
                        </div>
                    );
            }
        } catch (error) {
            // Error logged for debugging
            if (process.env.NODE_ENV === 'development') {
                console.error('Error:', error);
            }
            return (
                <div className={styles.errorState}>
                    <p>{t('renderError', {type: item?.entityType, message: error.message})}</p>
                </div>
            );
        }
    }, [item, data, t]);

    // Show error state if validation failed
    if (errorState) {
        return (
            <div className={styles.errorState}>
                <p>{t('error', {message: errorState.message})}</p>
            </div>
        );
    }

    // Ensure all hooks have run before this conditional return
    if (!item || !item.entityType) {
        return null;
    }

    return (
        <div className={styles.feedItem} onClick={handleItemClick}>
            {renderContent()}
        </div>
    );
});

// Add display name for better debugging
FeedItem.displayName = 'FeedItem';

export default FeedItem; 
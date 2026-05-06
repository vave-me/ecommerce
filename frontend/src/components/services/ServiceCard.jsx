"use client";
import React, { useCallback, useEffect, useMemo, useRef, useState, memo } from "react";
import {
    Calendar, 
    Camera,
    Clock,
    Eye,
    Flame,
    MapPin,
    MessageCircle,
    Star,
    Tag,
    Users,
    Shield,
    CheckCircle,
    ThumbsUp,
    ThumbsDown,
    Send,
    Bookmark
} from "@/icons";
import Link from "next/link";
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useDispatch } from "react-redux";
import { useAuth } from "../../context/AuthContext";
import { getMediaByItem, isMediaResponseSuccess } from "../../api/mediaApi";
import useActivityApi from "../../hooks/useActivityApi";
import { openMessageModal } from "../../redux/slices/modalsSlice";
import styles from './ServiceCard.module.css';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import CardImageContainer from '../shared/CardImageContainer';
import OfferModal from './OfferModal';

// Extend dayjs with relative time plugin
dayjs.extend(relativeTime);

/**
 * StarRating - Service rating component
 */
const StarRating = memo(({ rating, reviewCount }) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating - fullStars >= 0.5;

    for (let i = 0; i < 5; i++) {
        if (i < fullStars) {
            stars.push(<Star key={i} size={12} className={styles.starFull} />);
        } else if (i === fullStars && hasHalfStar) {
            stars.push(<Star key={i} size={12} className={styles.starHalf} />);
        } else {
            stars.push(<Star key={i} size={12} className={styles.starEmpty} />);
        }
    }

    return (
        <div className={styles.starRating} aria-label={`${rating} stars out of 5`}>
            {stars}
            {reviewCount > 0 && (
                <span className={styles.reviewCount} aria-label={`${reviewCount} reviews`}>
                    ({reviewCount})
                </span>
            )}
        </div>
    );
});

StarRating.displayName = 'StarRating';

/**
 * ServiceCard - Unified service card component
 * Modern design with proper server/client separation, mirroring DealCard architecture
 */
const ServiceCard = memo(({ service, preloadedMedia = null, className = "" }) => {

    const isMountedRef = useRef(true);
    const { user } = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const { handleLike, handleDislike } = useActivityApi();
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    // Add state for client-side time calculations to avoid hydration errors
    const [isNew, setIsNew] = useState(false);
    const [isAvailableToday, setIsAvailableToday] = useState(false);

    // Initialize media from preloaded data or fetch if needed
    useEffect(() => {
        // If we have preloaded media, use it immediately
        if (preloadedMedia && preloadedMedia.length > 0) {
            setMediaItems(preloadedMedia);
            setIsLoading(false);
            return;
        }

        // Otherwise, fetch media using the same pattern as working routes
        const loadMedia = async () => {
            if (!service?.id) {
                return;
            }
            
            setIsLoading(true);
            try {
                const mediaResponse = await getMediaByItem(service.id);
                
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || 'image',
                        alt: item.altText || service.name || 'Service image'
                    }));
                    setMediaItems(formattedMedia);
                } else {
                    // Check service object for any image data
                    const serviceImages = [];
                    if (service.thumbnail) serviceImages.push({ url: service.thumbnail, type: 'image' });
                    if (service.images?.length > 0) {
                        service.images.forEach(img => serviceImages.push({ url: img, type: 'image' }));
                    }
                    if (service.media?.length > 0) {
                        service.media.forEach(media => serviceImages.push(media));
                    }
                    
                    if (serviceImages.length > 0) {
                        setMediaItems(serviceImages);
                    }
                }
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
                setIsLoading(false);
            }
        };
        
        loadMedia();
    }, [service?.id, service?.name, preloadedMedia]);

    // Combine API media with service media, prioritizing API media
    const resolvedMedia = useMemo(() => {
        // If we have API media, use it
        if (mediaItems && mediaItems.length > 0) {
            return mediaItems;
        }

        // Otherwise, try to use service media
        if (service.media && service.media.length > 0) {
            return service.media;
        }

        // If service has images array, use it
        if (service.images && service.images.length > 0) {
            const mappedImages = service.images.map(img => ({ url: img }));
            return mappedImages;
        }

        // If service has thumbnail, use it
        if (service.thumbnail) {
            const thumbnailMedia = [{ url: service.thumbnail }];
            return thumbnailMedia;
        }

        // No media available - CardImageContainer will handle by returning null
        return [];
    }, [mediaItems, service?.media, service?.images, service?.thumbnail, service?.id]);

    // Full service data merger with defaults - using real API data
    const serviceWithDefaults = useMemo(() => {
        // Extract actual service data from API response structure
        const actualService = service?.service || service;
        const serviceMetrics = actualService?.metrics || {};

        // Handle completely missing service data
        if (!actualService) {
            // Error: '[ServiceCard] No service data provided'...
            return {
                id: '',
                name: 'Service data missing',
                description: 'This service could not be loaded.',
                basePrice: '0',
                hourlyRate: '0',
                created: new Date().toISOString(),
                availability: 'Contact for availability',
                condition: 'new',
                thumbnail: '/images/default-product.webp',
                serviceUrl: '',
                media: [{ url: '/images/default-product.webp', type: 'image' }],
                metrics: { likes: 0, dislikes: 0, comments: 0, shares: 0, views: 0 }
            };
        }

        // Extract provider name
        const providerName = actualService?.providerName || actualService?.merchantName || 'Service Provider';

        // Calculate rating from metrics
        const calculateRating = () => {
            const likes = parseInt(serviceMetrics?.likesCount || '0', 10);
            const dislikes = parseInt(serviceMetrics?.dislikesCount || '0', 10);
            const total = likes + dislikes;
            
            if (total === 0) return 4.5; // Default rating
            
            const ratio = likes / total;
            return Math.max(3.0, Math.min(5.0, 3.0 + (ratio * 2))); // Scale to 3.0-5.0
        };

        // Determine pricing model
        const pricingModel = actualService?.pricingModel || 'hourly';

        return {
            // Real API data from the service object
            ...actualService,
            id: actualService?.id || '',
            name: actualService?.name || 'Unnamed Service',
            description: actualService?.description || 'No description available.',
            basePrice: actualService?.basePrice || actualService?.listingPrice || '0',
            hourlyRate: actualService?.hourlyRate || actualService?.basePrice || '0',
            created: service?.createdAt || service?.updatedAt || new Date().toISOString(),
            availability: actualService?.availability || 'Contact for availability',
            condition: actualService?.condition || 'available',
            thumbnail: actualService?.thumbnail || '/images/default-product.webp',
            serviceUrl: actualService?.serviceUrl || actualService?.userId || '',
            brand: actualService?.brand || '',
            categorySlug: actualService?.categorySlug || 'services',
            categoryId: actualService?.categoryId || '',
            
            // Service-specific fields
            pricingModel: pricingModel,
            providerName: providerName,
            rating: calculateRating(),
            reviewCount: parseInt(serviceMetrics?.commentsCount || '0', 10),
            
            // Convert string metrics to numbers
            likeCount: parseInt(serviceMetrics?.likesCount || '0', 10),
            dislikeCount: parseInt(serviceMetrics?.dislikesCount || '0', 10),
            tags: actualService?.tags || [],
            lat: actualService?.lat || null,
            lng: actualService?.lng || null,
            
            // Calculate popularity based on views
            popularity: parseInt(serviceMetrics?.visitedCount || '0', 10) > 50 ? 2 :
                parseInt(serviceMetrics?.visitedCount || '0', 10) > 20 ? 1 : 0,
            wishlistCount: parseInt(serviceMetrics?.addedToWishlistCount || '0', 10),
            userId: actualService?.userId || '',

            // Service-specific data from API
            location: (actualService?.lat && actualService?.lng && (actualService.lat !== 0 || actualService.lng !== 0))
                ? `${actualService.lat}, ${actualService.lng}`
                : actualService?.address || actualService?.city || null,
            // Use fetched media items or fallback
            media: resolvedMedia,
            // Real metrics from API
            metrics: {
                likes: parseInt(serviceMetrics?.likesCount || '0', 10),
                dislikes: parseInt(serviceMetrics?.dislikesCount || '0', 10),
                comments: parseInt(serviceMetrics?.commentsCount || '0', 10),
                shares: parseInt(serviceMetrics?.sharedCount || '0', 10),
                views: parseInt(serviceMetrics?.visitedCount || '0', 10),
            },
        };
    }, [service, resolvedMedia]);

    // Destructure properties for easier access
    const {
        id,
        name,
        description,
        created,
        basePrice,
        hourlyRate,
        pricingModel,
        serviceUrl,
        availability,
        tags,
        location,
        providerName,
        media,
        metrics,
        popularity,
        rating,
        reviewCount
    } = serviceWithDefaults;

    // State management
    const [interactionState, setInteractionState] = useState({
        favorite: false,
        liked: false,
        disliked: false,
        showComments: false,
        currentMediaIndex: 0,
        showOfferModal: false,
    });

    // Derived data
    // Format price with localization
    const formattedBasePrice = useMemo(() => {
        const price = parseInt(basePrice, 10) || 0;
        return price.toLocaleString('en-US', { style: 'currency', currency: 'EUR' });
    }, [basePrice]);

    const formattedHourlyRate = useMemo(() => {
        const price = parseInt(hourlyRate, 10) || 0;
        return price.toLocaleString('en-US', { style: 'currency', currency: 'EUR' });
    }, [hourlyRate]);

    // Time-related calculations
    const timeAgo = useMemo(() => (
        created ? dayjs(created).fromNow() : ''
    ), [created]);

    // Move time-dependent calculations to useEffect to avoid hydration errors
    useEffect(() => {
        // Calculate if service is new (less than 7 days old)
        const isNewService = service?.isNew || (created && dayjs().diff(dayjs(created), 'day') < 7);
        setIsNew(isNewService);
        
        // Check if available today (mock logic - would come from actual availability API)
        const todayAvailable = availability && (
            availability.toLowerCase().includes('today') || 
            availability.toLowerCase().includes('available') ||
            availability.toLowerCase().includes('now')
        );
        setIsAvailableToday(todayAvailable);
    }, [created, availability, service?.isNew]);

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
        };
    }, []);

    // Event handlers
    const handleCommentClick = useCallback(() => {
        setInteractionState(prev => ({
            ...prev,
            showComments: !prev.showComments
        }));
    }, []);

    const handleFavorite = useCallback(() => {
        setInteractionState(prev => ({
            ...prev,
            favorite: !prev.favorite
        }));
    }, []);

    const handleLikeClick = useCallback(() => {
        if (!userId) {
            return;
        }

        if (!id) return;

        handleLike(id, userId).catch(() => {
            // handle error if needed
        });

        setInteractionState(prev => ({
            ...prev,
            liked: true,
            disliked: false
        }));
    }, [id, userId, handleLike]);

    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            return;
        }

        if (!id) return;

        handleDislike(id, userId).catch(() => {
            // handle error if needed
        });

        setInteractionState(prev => ({
            ...prev,
            liked: false,
            disliked: true
        }));
    }, [id, userId, handleDislike]);

    const handleOpenMessage = useCallback(() => {
        if (!id || !serviceWithDefaults.userId) return;
        dispatch(openMessageModal({
            recipientId: serviceWithDefaults.userId,
            itemId: id
        }));
    }, [id, serviceWithDefaults.userId, dispatch]);

    const handleMediaNavigation = useCallback((direction, e) => {
        e?.stopPropagation();
        if (media.length <= 1) return;
        
        setInteractionState(prev => {
            const newIndex = direction === 'next'
                ? (prev.currentMediaIndex + 1) % media.length
                : (prev.currentMediaIndex - 1 + media.length) % media.length;
            
            return { ...prev, currentMediaIndex: newIndex };
        });
    }, [media.length]);

    const handleKeyboardNavigation = useCallback((e) => {
        if (e.key === 'ArrowLeft') {
            handleMediaNavigation('prev', e);
        } else if (e.key === 'ArrowRight') {
            handleMediaNavigation('next', e);
        }
    }, [handleMediaNavigation]);

    const handleGoToService = useCallback(() => {
        if (serviceUrl) {
            window.open(serviceUrl, '_blank', 'noopener,noreferrer');
        }
    }, [serviceUrl]);

    const handleOfferClick = useCallback((e) => {
        e.stopPropagation();
        setInteractionState(prev => ({ ...prev, showOfferModal: true }));
    }, []);

    const handleOfferModalClose = useCallback(() => {
        setInteractionState(prev => ({ ...prev, showOfferModal: false }));
    }, []);

    const handleOfferSuccess = useCallback((offerData) => {
        // You can add a success notification here
        
    }, []);

    // Media access now uses direct array indexing like DealCard
    // Removed currentMedia useMemo to match working DealCard pattern

    return (
        <div className={`${styles.card} ${className}`} tabIndex={0} onKeyDown={handleKeyboardNavigation}>

            {/* STATUS INDICATORS - Compact overlay */}
            <div className={styles.statusIndicators}>
                {isNew && (
                    <span className={`${styles.badge} ${styles.newBadge}`} aria-label="New service">
                        NEW
                    </span>
                )}
                {isAvailableToday && (
                    <span className={`${styles.badge} ${styles.availableBadge}`} aria-label="Available today">
                        <Clock size={12} aria-hidden="true" />
                        <span>Today</span>
                    </span>
                )}
                {rating >= 4.5 && (
                    <span className={`${styles.badge} ${styles.topRatedBadge}`} aria-label="Top rated">
                        <Star size={12} aria-hidden="true" />
                        <span>{rating.toFixed(1)}</span>
                    </span>
                )}
            </div>

            {/* IMAGE */}
            <CardImageContainer
                media={media}
                currentIndex={interactionState.currentMediaIndex}
                isLoading={isLoading}
                alt={media[interactionState.currentMediaIndex]?.alt || name}
                onClick={handleGoToService}
                onNavigate={handleMediaNavigation}
                ariaLabel="Service images"
            />

            {/* CONTENT SECTION */}
            <div className={styles.content}>
                {/* Category Label */}
                <span className={styles.categoryLabel}>
                    {serviceWithDefaults.categorySlug || 'Services'}
                </span>

                {/* Title */}
                <Link href={`/services/${id}`} className={styles.titleLink}>
                    <h3 className={styles.title}>{name}</h3>
                </Link>

                {/* Rating Row */}
                <div className={styles.ratingRow}>
                    <StarRating rating={rating} reviewCount={reviewCount} />
                </div>

                {/* Price Section */}
                <div className={styles.priceSection}>
                    <div className={styles.priceWrapper}>
                        <span className={styles.price}>
                            {pricingModel === 'hourly' ? `${formattedHourlyRate}/hr` : formattedBasePrice}
                        </span>
                        <span className={styles.pricingBadge}>
                            {pricingModel === 'hourly' ? 'Hourly' : 
                             pricingModel === 'fixed' ? 'Fixed' : 'Quote'}
                        </span>
                    </div>
                    <button 
                        className={styles.bookServiceButton}
                        onClick={handleOfferClick}
                        aria-label="Make an offer for this service"
                    >
                        <Tag size={14} />
                        <span>Offer</span>
                    </button>
                </div>

                {/* Service Info Items */}
                <div className={styles.serviceInfo}>
                    {location && (
                        <div className={styles.infoItem}>
                            <MapPin size={12} />
                            <span>{location}</span>
                        </div>
                    )}
                    {availability && (
                        <div className={styles.infoItem}>
                            <Calendar size={12} />
                            <span>{availability}</span>
                        </div>
                    )}
                    {providerName && (
                        <div className={styles.infoItem}>
                            <Users size={12} />
                            <span>{providerName}</span>
                        </div>
                    )}
                    {metrics.views > 0 && (
                        <div className={styles.infoItem}>
                            <Eye size={12} />
                            <span>{metrics.views}</span>
                        </div>
                    )}
                </div>

                {/* Description */}
                <p className={styles.description}>{description}</p>

            </div>

            {/* ENGAGEMENT BAR */}

            <div className={styles.engagementBar}>
                <div className={styles.engagementActions}>
                    <button 
                        className={`${styles.actionButton} ${interactionState.liked ? styles.active : ''}`}
                        onClick={handleLikeClick}
                        aria-label="Like service"
                    >
                        <ThumbsUp size={20} />
                        {metrics.likes > 0 && (
                            <span className={styles.actionCount}>{metrics.likes}</span>
                        )}
                    </button>
                    
                    <button 
                        className={`${styles.actionButton} ${interactionState.disliked ? styles.active : ''}`}
                        onClick={handleDislikeClick}
                        aria-label="Dislike service"
                    >
                        <ThumbsDown size={20} />
                        {metrics.dislikes > 0 && (
                            <span className={styles.actionCount}>{metrics.dislikes}</span>
                        )}
                    </button>
                    
                    <button 
                        className={styles.actionButton}
                        onClick={handleCommentClick}
                        aria-label="Comment on service"
                    >
                        <MessageCircle size={20} />
                        {metrics.comments > 0 && (
                            <span className={styles.actionCount}>{metrics.comments}</span>
                        )}
                    </button>
                    
                    <button 
                        className={styles.actionButton}
                        onClick={handleOpenMessage}
                        aria-label="Message provider"
                    >
                        <Send size={20} />
                    </button>
                    
                    <button 
                        className={`${styles.actionButton} ${interactionState.favorite ? styles.active : ''}`}
                        onClick={handleFavorite}
                        aria-label="Save service"
                    >
                        <Bookmark size={20} />
                    </button>
                </div>
            </div>

            {/* COMMENTS */}
            {interactionState.showComments && (
                <div className={styles.commentsWrapper}>
                    <CommentsSetup 
                        itemId={id}
                        itemType="service"
                        userId={userId}
                    />
                </div>
            )}

            {/* OFFER MODAL */}
            <OfferModal
                isOpen={interactionState.showOfferModal}
                onClose={handleOfferModalClose}
                service={serviceWithDefaults}
                onSuccess={handleOfferSuccess}
            />
        </div>
    );
});

ServiceCard.displayName = 'ServiceCard';

export default ServiceCard; 
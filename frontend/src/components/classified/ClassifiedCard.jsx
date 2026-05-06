"use client";
import React, {useCallback, useEffect, useMemo, useRef, useState, memo} from 'react';
import {
    Camera, Clock, Eye, Flame, MapPin, MessageCircle, ShoppingBag, ShoppingCart, Plus, Minus,
    Star, Tag, Snowflake, Heart, ThumbsUp, ThumbsDown, Share2, Package, Check, AlertCircle,
    Bookmark, ExternalLink, Truck, Shield, CheckCircle, Zap, TrendingUp, PlayCircle
} from '@/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import {useAuth} from '../../context/AuthContext';
import {useDispatch} from 'react-redux';
import {openMessageModal} from '../../redux/slices/modalsSlice';
import {toast} from 'react-toastify';
import useActivityApi from '../../hooks/useActivityApi';
import { useProductBasket } from '../../hooks/useProductBasket';
import useWishlist from '../../hooks/useWishlist';
import {getMediaByItem} from '../../api/mediaApi';
import {Engagement} from '../shared/Engagement';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import { openCommentsFullModal } from '../../redux/slices/modalsSlice';
import BadgeRow from '../shared/BadgeRow';
import PriceLocationRow from '../shared/PriceLocationRow';
import ExpandableDescription from '../shared/ExpandableDescription';
import CardTitle from '../shared/CardTitle';
import CardImageContainer from '../shared/CardImageContainer';
import Tags from '../shared/Tags';
import VideoPlayer from '../shared/VideoPlayer';
import styles from './ClassifiedCard.module.css';
// Extend dayjs with relative time plugin
dayjs.extend(relativeTime);
/**
 * QuantitySelector - Marketplace quantity selector component
 */
const QuantitySelector = memo(({quantity, onQuantityChange, stock, disabled = false}) => {
    const handleDecrease = useCallback(() => {
        if (quantity > 1) {
            onQuantityChange(quantity - 1);
        }
    }, [quantity, onQuantityChange]);
    const handleIncrease = useCallback(() => {
        if (quantity < stock) {
            onQuantityChange(quantity + 1);
        }
    }, [quantity, stock, onQuantityChange]);
    return (
        <div className={styles.quantitySelector}>
            <button
                type="button"
                className={styles.quantityButton}
                onClick={handleDecrease}
                disabled={disabled || quantity <= 1}
                aria-label="Decrease quantity"
            >
                <Minus size={14} />
            </button>
            <span className={styles.quantityDisplay}>{quantity}</span>
            <button
                type="button"
                className={styles.quantityButton}
                onClick={handleIncrease}
                disabled={disabled || quantity >= stock}
                aria-label="Increase quantity"
            >
                <Plus size={14} />
            </button>
        </div>
    );
});
QuantitySelector.displayName = 'QuantitySelector';
/**
 * ClassifiedCard - Marketplace-optimized product card component
 * Features:
 * - Basket/cart functionality with quantity selection
 * - Stock management and availability indicators
 * - Marketplace-specific actions (quick view, compare, wishlist)
 * - Simplified layout without user header or merchant fields
 * - Enhanced product information display
 * - Mobile-optimized responsive design
 */
const ClassifiedCard = memo(({
                            product,
                            preloadedMedia = null,
                            hasPreloadedMedia = false
                        }) => {
    // Refs and hooks setup
    const isMountedRef = useRef(true);
    const {user} = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const {handleLike, handleDislike} = useActivityApi();
    // Process product data with comprehensive defaults
    const productWithDefaults = useMemo(() => {
        const actualProduct = product?.product || product || {};
        return {
            // Core product data
            id: actualProduct.id || '',
            name: actualProduct.name || actualProduct.title || 'Unnamed Product',
            description: actualProduct.description || 'No description available.',
            basePrice: actualProduct.basePrice || actualProduct.price || '0',
            added: actualProduct.added || actualProduct.createdAt || new Date().toISOString(),
            condition: actualProduct.condition || 'new',
            thumbnail: actualProduct.thumbnail || actualProduct.image || '',
            // Product-specific data
            brand: actualProduct.brand || '',
            model: actualProduct.model || '',
            sku: actualProduct.sku || '',
            stock: actualProduct.stock || 0,
            manageStocks: actualProduct.manageStocks || false,
            hasVariants: actualProduct.hasVariants || false,
            shippingCost: actualProduct.shippingCost || '0',
            weight: actualProduct.weight || '',
            dimensions: {
                height: actualProduct.height || '',
                width: actualProduct.width || '',
                depth: actualProduct.depth || ''
            },
            // Business data
            negotiable: actualProduct.negotiable || false,
            middlemanService: actualProduct.middlemanService || false,
            userType: actualProduct.userType || 'private',
            status: actualProduct.status || 'active',
            likeCount: actualProduct.likeCount ?? 0,
            dislikeCount: actualProduct.dislikeCount ?? 0,
            wishlistCount: actualProduct.wishlistCount ?? 0,
            comments: actualProduct.comments || [],
            tags: actualProduct.tags || [],
            popularity: actualProduct.popularity ?? 0,
            verified: actualProduct.verified !== false,
            categorySlug: actualProduct.categorySlug,
            categoryId: actualProduct.categoryId || '',
            media: actualProduct.media || [],
            // Seller information
            userId: actualProduct.userId || actualProduct.userSellerId || actualProduct.sellerId || '',
            metrics: {
                likes: actualProduct.metrics?.likesCount ? parseInt(actualProduct.metrics.likesCount) : (actualProduct.likeCount || 0),
                dislikes: actualProduct.metrics?.dislikesCount ? parseInt(actualProduct.metrics.dislikesCount) : (actualProduct.dislikeCount || 0),
                comments: actualProduct.metrics?.commentsCount ? parseInt(actualProduct.metrics.commentsCount) : ((actualProduct.comments || []).length),
                shares: actualProduct.metrics?.sharedCount ? parseInt(actualProduct.metrics.sharedCount) : 0,
                views: actualProduct.metrics?.visitedCount ? parseInt(actualProduct.metrics.visitedCount) : (actualProduct.viewCount || actualProduct.views || 5),
                saved: actualProduct.metrics?.addedToWishlistCount ? parseInt(actualProduct.metrics.addedToWishlistCount) : (actualProduct.wishlistCount || 0),
            },
            // Review data - for now using mock data, should be fetched from API
            rating: actualProduct.rating || 4.5,
            reviewCount: actualProduct.reviewCount || 127
        };
    }, [product]);
    // Initialize basket actions for this product
    const {addToBasket: handleAddToCart, isLoading: basketLoading} = useProductBasket(productWithDefaults.id, productWithDefaults);
    
    // Initialize wishlist hook
    const {
        isInAnyWishlist,
        toggleWishlist,
        isLoading: wishlistLoading
    } = useWishlist();
    // Destructure for easier access
    const {
        id, name, description, basePrice, added, condition, brand, model, sku,
        stock, manageStocks, hasVariants, shippingCost, weight, dimensions,
        negotiable, middlemanService, userType, status, tags,
        popularity, verified, metrics, categorySlug, thumbnail, rating, reviewCount
    } = productWithDefaults;
    // State management
    const [interactionState, setInteractionState] = useState({
        favorite: false,
        liked: false,
        disliked: false,
        showComments: false,
        currentMediaIndex: 0,
        selectedQuantity: 1,
        isAddingToCart: false,
        playingVideos: {}, // Track which videos are playing
    });
    
    // Check if product is in wishlist
    const isInWishlist = useMemo(() => {
        // Only check wishlist if user is authenticated
        return userId && id ? isInAnyWishlist(id, 'product') : false;
    }, [userId, id, isInAnyWishlist]);
    
    // Check if mobile
    const [isMobile, setIsMobile] = useState(false);
    
    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
        };
        
        checkMobile();
        window.addEventListener('resize', checkMobile);
        
        return () => window.removeEventListener('resize', checkMobile);
    }, []);
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoadingMedia, setIsLoadingMedia] = useState(false);
    // Initialize media from preloaded data or fetch if needed
    useEffect(() => {
        // If we have preloaded media, use it immediately
        if (preloadedMedia && preloadedMedia.length > 0) {
            setMediaItems(preloadedMedia);
            setIsLoadingMedia(false);
            return;
        }
        // Otherwise, fetch media using the same pattern as working routes
        const loadMedia = async () => {
            if (!id) {
                return;
            }
            setIsLoadingMedia(true);
            try {
                const mediaResponse = await getMediaByItem(id);
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || 'image',
                        alt: item.altText || name || 'Product image'
                    }));
                    setMediaItems(formattedMedia);
                } else {
                    // Check product object for any image data
                    const productImages = [];
                    if (thumbnail) productImages.push({url: thumbnail, type: 'image'});
                    if (productWithDefaults.images?.length > 0) {
                        productWithDefaults.images.forEach(img => productImages.push({url: img, type: 'image'}));
                    }
                    if (productWithDefaults.media?.length > 0) {
                        productWithDefaults.media.forEach(media => productImages.push(media));
                    }
                    if (productImages.length > 0) {
                        setMediaItems(productImages);
                    }
                }
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
                setIsLoadingMedia(false);
            }
        };
        loadMedia();
    }, [id, name, preloadedMedia, thumbnail, productWithDefaults.images, productWithDefaults.media]);
    // Helper function to detect media type from URL
    const getMediaTypeFromUrl = (url) => {
        if (!url) return 'image';
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv', '.m4v'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext)) ? 'video' : 'image';
    };
    
    // Combine API media with product media, prioritizing API media
    const resolvedMedia = useMemo(() => {
        // Helper to ensure media has correct type
        const ensureMediaType = (media) => {
            if (typeof media === 'string') {
                return { url: media, type: getMediaTypeFromUrl(media) };
            }
            // If media object doesn't have type, detect from URL
            if (!media.type && media.url) {
                return { ...media, type: getMediaTypeFromUrl(media.url) };
            }
            return media;
        };
        
        // If we have API media, use it
        if (mediaItems && mediaItems.length > 0) {
            return mediaItems.map(ensureMediaType);
        }
        // Otherwise, try to use product media
        if (productWithDefaults.media && productWithDefaults.media.length > 0) {
            return productWithDefaults.media.map(ensureMediaType);
        }
        // If product has images array, use it
        if (productWithDefaults.images && productWithDefaults.images.length > 0) {
            return productWithDefaults.images.map(img => 
                typeof img === 'string' ? ensureMediaType(img) : ensureMediaType({url: img})
            );
        }
        // If product has thumbnail, use it
        if (thumbnail) {
            return [ensureMediaType(thumbnail)];
        }
        // No media available
        return [];
    }, [mediaItems, productWithDefaults?.media, productWithDefaults?.images, thumbnail, id]);
    // Price calculations
    const formattedPrice = useMemo(() => {
        const price = parseFloat(basePrice);
        return price.toLocaleString('de-DE', {style: 'currency', currency: 'EUR'});
    }, [basePrice]);
    const formattedShippingCost = useMemo(() => {
        const shipping = parseFloat(shippingCost);
        return shipping > 0 ? shipping.toLocaleString('de-DE', {style: 'currency', currency: 'EUR'}) : null;
    }, [shippingCost]);
    const timeAgo = useMemo(() => (
        added ? dayjs(added).fromNow() : ''
    ), [added]);
    // Get popularity status
    const getPopularityStatus = useCallback(() => {
        if (popularity >= 2) return {icon: <Flame size={14}/>, label: "Hot Product", className: styles.hotBadge};
        if (popularity >= 0.5) return {icon: <Star size={14}/>, label: "Popular", className: styles.popularBadge};
        if (popularity <= -2) return {icon: <Snowflake size={14}/>, label: "Cold", className: styles.coldBadge};
        return null;
    }, [popularity]);
    const popularityStatus = getPopularityStatus();
    // Stock availability
    const isInStock = useMemo(() => {
        if (!manageStocks) return true;
        return stock > 0;
    }, [manageStocks, stock]);
    const stockStatus = useMemo(() => {
        if (!manageStocks) return { available: true, label: 'Available', className: styles.available };
        if (stock === 0) return { available: false, label: 'Out of Stock', className: styles.outOfStock };
        if (stock <= 5) return { available: true, label: `Only ${stock} left`, className: styles.lowStock };
        return { available: true, label: 'In Stock', className: styles.inStock };
    }, [manageStocks, stock]);
    // Media navigation
    const handleMediaNavigation = useCallback((direction) => {
        setInteractionState(prev => ({
            ...prev,
            currentMediaIndex: direction === 'next'
                ? (prev.currentMediaIndex + 1) % resolvedMedia.length
                : (prev.currentMediaIndex - 1 + resolvedMedia.length) % resolvedMedia.length,
            playingVideos: {} // Stop all videos when navigating
        }));
    }, [resolvedMedia.length]);
    
    // Video play/pause handler
    const handleVideoToggle = useCallback((index, videoRef) => {
        if (!videoRef || !videoRef.current) {
            // Just update state (e.g., when video ends)
            setInteractionState(prev => ({
                ...prev,
                playingVideos: { ...prev.playingVideos, [index]: false }
            }));
            return;
        }
        
        if (interactionState.playingVideos[index]) {
            videoRef.current.pause();
            setInteractionState(prev => ({
                ...prev,
                playingVideos: { ...prev.playingVideos, [index]: false }
            }));
        } else {
            // Pause all other videos
            Object.keys(interactionState.playingVideos).forEach(key => {
                if (key !== index.toString() && interactionState.playingVideos[key]) {
                    const otherVideo = document.querySelector(`[data-video-index="${key}"]`);
                    if (otherVideo) otherVideo.pause();
                }
            });
            
            videoRef.current.play();
            setInteractionState(prev => ({
                ...prev,
                playingVideos: { [index]: true }
            }));
        }
    }, [interactionState.playingVideos]);

    // Event handlers
    const handleCommentClick = useCallback(() => {
        if (isMobile) {
            // On mobile, open CommentsFull modal
            dispatch(openCommentsFullModal({
                itemId: id,
                itemType: 'product',
                categoryId: productWithDefaults.categoryId
            }));
        } else {
            // On desktop, toggle inline comments
            setInteractionState(prev => ({
                ...prev,
                showComments: !prev.showComments
            }));
        }
    }, [isMobile, dispatch, id, productWithDefaults.categoryId]);
    const handleFavorite = useCallback(async () => {
        if (!userId) {
            toast.warn("Please log in to add items to your wishlist.");
            return;
        }
        if (!id || wishlistLoading) return;
        
        try {
            await toggleWishlist(id, 'product');
            // Update local state will be handled by the wishlist state change
        } catch (error) {
            toast.error("Failed to update wishlist");
        }
    }, [id, userId, toggleWishlist, wishlistLoading]);
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to like products.");
            return;
        }
        if (!id) return;
        handleLike(id, userId).catch(() => {
            // Silent error handling
        });
        setInteractionState(prev => ({
            ...prev,
            liked: true,
            disliked: false
        }));
    }, [id, userId, handleLike]);
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to dislike products.");
            return;
        }
        if (!id) return;
        handleDislike(id, userId).catch(() => {
            // Silent error handling
        });
        setInteractionState(prev => ({
            ...prev,
            liked: false,
            disliked: true
        }));
    }, [id, userId, handleDislike]);
    const handleOpenMessage = useCallback(() => {
        if (!id) return;
        const sellerId = productWithDefaults.userId;
        if (!sellerId) {
            toast.error("Unable to contact seller - seller information missing");
            return;
        }
        dispatch(
            openMessageModal({
                itemId: id,
                recipientId: sellerId,
            })
        );
    }, [id, productWithDefaults.userId, dispatch]);
    const handleQuantityChange = useCallback((newQuantity) => {
        setInteractionState(prev => ({
            ...prev,
            selectedQuantity: newQuantity
        }));
    }, []);
    const handleAddToBasket = useCallback(async () => {
        if (!isInStock) {
            toast.warn("This product is currently out of stock.");
            return;
        }
        setInteractionState(prev => ({...prev, isAddingToCart: true}));
        try {
            await handleAddToCart(interactionState.selectedQuantity);
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    } finally {
            setInteractionState(prev => ({...prev, isAddingToCart: false}));
        }
    }, [handleAddToCart, interactionState.selectedQuantity, isInStock]);
    // Cleanup on unmount
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
        };
    }, []);
    return (
        <div className={`${styles.card} ${isLoadingMedia ? styles.isLoading : ''}`}>
            {/* Media Container - Primary Visual Hierarchy */}
            <div className={styles.mediaContainer}>
                {/* Status Indicators */}
                <div className={styles.statusIndicators}>
                    {popularityStatus && (
                        <div className={`${styles.badge} ${popularityStatus.className === styles.hotBadge ? styles.featuredBadge : styles.newBadge}`}>
                            {popularityStatus.icon}
                            <span>{popularityStatus.label}</span>
                        </div>
                    )}
                    {stockStatus.available && stock <= 5 && (
                        <div className={`${styles.badge} ${styles.newBadge}`}>
                            <AlertCircle size={12} />
                            <span>Low Stock</span>
                        </div>
                    )}
                </div>

                {/* Media Carousel */}
                {resolvedMedia.length > 0 ? (
                    <>
                        <div className={styles.mediaCarousel} style={{ transform: `translateX(-${interactionState.currentMediaIndex * 100}%)` }}>
                            {resolvedMedia.map((media, index) => (
                                <div key={index} className={styles.mediaItem}>
                                    {media.type === 'video' ? (
                                        <VideoPlayer
                                            key={`video-${index}-${media.url}`}
                                            media={media}
                                            index={index}
                                            isPlaying={interactionState.playingVideos[index] || false}
                                            onToggle={handleVideoToggle}
                                            resolvedMediaLength={resolvedMedia.length}
                                            productName={name}
                                        />
                                    ) : (
                                        <img src={media.url} alt={media.alt || name} loading="lazy" />
                                    )}
                                </div>
                            ))}
                        </div>

                        {/* Media Navigation - Progressive Disclosure */}
                        {resolvedMedia.length > 1 && (
                            <>
                                <div className={styles.mediaNav}>
                                    <button 
                                        className={styles.mediaNavButton} 
                                        onClick={() => handleMediaNavigation('prev')}
                                        aria-label="Previous image"
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                            <path d="M15 18l-6-6 6-6" />
                                        </svg>
                                    </button>
                                    <button 
                                        className={styles.mediaNavButton} 
                                        onClick={() => handleMediaNavigation('next')}
                                        aria-label="Next image"
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                            <path d="M9 6l6 6-6 6" />
                                        </svg>
                                    </button>
                                </div>
                                <div className={styles.mediaPagination}>
                                    {resolvedMedia.map((_, index) => (
                                        <button
                                            key={index}
                                            className={`${styles.paginationDot} ${index === interactionState.currentMediaIndex ? styles.active : ''}`}
                                            onClick={() => setInteractionState(prev => ({ ...prev, currentMediaIndex: index }))}
                                            aria-label={`Go to image ${index + 1}`}
                                        />
                                    ))}
                                </div>
                            </>
                        )}
                    </>
                ) : (
                    <div className={styles.mediaItem}>
                        <Camera size={48} color="var(--color-text-tertiary)" />
                    </div>
                )}
            </div>
            {/* Content Section */}
            <div className={styles.content}>
                {/* Category or Type above title */}
                {categorySlug && (
                    <div className={styles.categoryLabel}>
                        {categorySlug}
                    </div>
                )}
                
                {/* Title - Secondary Hierarchy */}
                <Link href={`/products/${categorySlug}/${id}`} className={styles.titleLink}>
                    <h3 className={styles.title}>
                        {brand && model ? `${brand} ${model}` : name}
                    </h3>
                </Link>

                {/* Review Indicators */}
                {reviewCount > 0 && (
                    <div className={styles.reviewIndicators}>
                        <div className={styles.ratingStars}>
                            {[...Array(5)].map((_, index) => (
                                <Star 
                                    key={index} 
                                    size={12} 
                                    className={index < Math.floor(rating) ? styles.starFilled : styles.starEmpty}
                                />
                            ))}
                        </div>
                        <span className={styles.ratingValue}>{rating.toFixed(1)}</span>
                        <span className={styles.reviewCount}>({reviewCount})</span>
                    </div>
                )}

                {/* Price and Cart Section */}
                <div className={styles.priceSection}>
                    <div className={styles.priceWrapper}>
                        <span className={styles.price}>{formattedPrice}</span>
                        {negotiable && (
                            <span className={styles.badge} style={{ fontSize: '11px' }}>Negotiable</span>
                        )}
                    </div>
                    
                    {/* Purchase Controls */}
                    <div className={styles.purchaseControls}>
                        <QuantitySelector
                            quantity={interactionState.selectedQuantity}
                            onQuantityChange={handleQuantityChange}
                            stock={stock}
                            disabled={!isInStock || interactionState.isAddingToCart}
                        />
                        <button
                            className={styles.addToCartButton}
                            onClick={handleAddToBasket}
                            disabled={!isInStock || interactionState.isAddingToCart || basketLoading}
                            aria-label={`Add ${interactionState.selectedQuantity} to cart`}
                        >
                            {interactionState.isAddingToCart || basketLoading ? (
                                <>
                                    <div className={styles.spinner} />
                                    <span>Add</span>
                                </>
                            ) : (
                                <>
                                    <ShoppingCart size={16} />
                                    <span>Add</span>
                                </>
                            )}
                        </button>
                    </div>
                </div>

                {/* Product Info - Shipping and Verified */}
                <div className={styles.productInfo}>
                    {formattedShippingCost && (
                        <div className={styles.infoItem}>
                            <Truck size={14} />
                            <span className={styles.shippingText}>Shipping: {formattedShippingCost}</span>
                            <span className={styles.shippingTextMobile}>{formattedShippingCost}</span>
                        </div>
                    )}
                    {verified && (
                        <div className={styles.infoItem}>
                            <CheckCircle size={14} />
                            <span className={styles.verifiedText}>Verified Seller</span>
                            <span className={styles.verifiedTextMobile}>Verified</span>
                        </div>
                    )}
                    {stockStatus.available && manageStocks && (
                        <div className={styles.infoItem}>
                            <Package size={14} />
                            <span className={`${styles.stockIndicator} ${stockStatus.className}`}>
                                {stockStatus.label}
                            </span>
                        </div>
                    )}
                </div>

                {/* Description - Quaternary */}
                <p className={styles.description}>
                    {description.length > 200 ? `${description.substring(0, 200)}...` : description}
                </p>

            </div>
            {/* Engagement Bar - Progressive Disclosure */}
            <div className={styles.engagementBar}>
                <div className={styles.engagementActions}>
                    {/* Like/Dislike */}
                    <button 
                        className={`${styles.actionButton} ${interactionState.liked ? styles.active : ''}`}
                        onClick={handleLikeClick}
                        aria-label="Like product"
                    >
                        <ThumbsUp size={18} />
                        {metrics.likes > 0 && <span className={styles.actionCount}>{metrics.likes}</span>}
                    </button>
                    <button 
                        className={`${styles.actionButton} ${interactionState.disliked ? styles.active : ''}`}
                        onClick={handleDislikeClick}
                        aria-label="Dislike product"
                    >
                        <ThumbsDown size={18} />
                        {metrics.dislikes > 0 && <span className={styles.actionCount}>{metrics.dislikes}</span>}
                    </button>
                    
                    {/* Comment */}
                    <button 
                        className={`${styles.actionButton} ${!isMobile && interactionState.showComments ? styles.active : ''}`}
                        onClick={handleCommentClick}
                        aria-label="Comments"
                    >
                        <MessageCircle size={18} />
                        {metrics.comments > 0 && <span className={styles.actionCount}>{metrics.comments}</span>}
                    </button>
                    
                    {/* Message */}
                    <button 
                        className={styles.actionButton}
                        onClick={handleOpenMessage}
                        aria-label="Send message"
                    >
                        <MessageCircle size={18} />
                    </button>
                    
                    {/* Share */}
                    <button 
                        className={styles.actionButton}
                        onClick={() => {
                            if (navigator.share && typeof window !== 'undefined') {
                                navigator.share({
                                    title: name,
                                    text: description,
                                    url: `${window.location.origin}/products/${categorySlug}/${id}`
                                });
                            }
                        }}
                        aria-label="Share product"
                    >
                        <Share2 size={18} />
                        {metrics.shares > 0 && <span className={styles.actionCount}>{metrics.shares}</span>}
                    </button>
                    
                    {/* Favorite/Wishlist */}
                    <button 
                        className={`${styles.actionButton} ${isInWishlist ? styles.active : ''}`}
                        onClick={handleFavorite}
                        disabled={wishlistLoading}
                        aria-label={isInWishlist ? "Remove from wishlist" : "Add to wishlist"}
                        title={isInWishlist ? "Remove from wishlist" : "Add to wishlist"}
                    >
                        <Heart size={18} className={isInWishlist ? styles.filled : ''} />
                        {metrics.saved > 0 && <span className={styles.actionCount}>{metrics.saved}</span>}
                    </button>
                </div>

            </div>
            {/* Comments - Desktop only inline display */}
            {!isMobile && interactionState.showComments && (
                <div className={styles.commentsWrapper}>
                    <CommentsSetup
                        userId={userId}
                        itemId={id}
                        itemType="product"
                        toggleCommentsList={handleCommentClick}
                        categoryId={productWithDefaults.categoryId}
                        productName={name}
                        productThumbnail={thumbnail}
                    />
                </div>
            )}
        </div>
    );
});
ClassifiedCard.displayName = 'ClassifiedCard';
export default ClassifiedCard; 
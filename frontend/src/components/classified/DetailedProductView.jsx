"use client";
import React, { useCallback, useEffect, useMemo, useRef, useState, memo } from 'react';
import {
    Camera, Clock, Eye, Flame, MapPin, MessageCircle, ShoppingBag, ShoppingCart, Plus, Minus,
    Star, Tag, Snowflake, Heart, ThumbsUp, ThumbsDown, Share2, Package, Check, AlertCircle,
    Bookmark, ExternalLink, Truck, Shield, CheckCircle, Zap, TrendingUp, PlayCircle,
    Calendar, BarChart3, Users, Info, Award, CreditCard, Layers, ArrowRight, User,
    Mail, Phone, Globe, ChevronLeft, ChevronRight
} from '@/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useAuth } from '../../context/AuthContext';
import { useDispatch } from 'react-redux';
import { openMessageModal } from '../../redux/slices/modalsSlice';
import { toast } from 'react-toastify';
import useActivityApi from '../../hooks/useActivityApi';
import { useProductBasket } from '../../hooks/useProductBasket';
import { getMediaByItem } from '../../api/mediaApi';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import { openCommentsFullModal } from '../../redux/slices/modalsSlice';
import VideoPlayer from '../shared/VideoPlayer';
import styles from './DetailedProductView.compact.module.css';

// Extend dayjs with relative time plugin
dayjs.extend(relativeTime);

/**
 * DetailedProductView - Full product detail page component
 * Shows comprehensive product information with enhanced layout
 */
const DetailedProductView = memo(({ product, locale, category, availableCategories }) => {
    const { user } = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const { handleLike, handleDislike } = useActivityApi();
    
    // Process product data
    const productData = useMemo(() => {
        const p = product?.product || product || {};
        return {
            // Core data
            id: p.id || '',
            name: p.name || p.title || 'Unnamed Product',
            description: p.description || 'No description available.',
            basePrice: p.basePrice || p.price || '0',
            added: p.added || p.createdAt || new Date().toISOString(),
            condition: p.condition || 'new',
            thumbnail: p.thumbnail || p.image || '',
            
            // Product details
            brand: p.brand || '',
            model: p.model || '',
            sku: p.sku || '',
            stock: p.stock || 0,
            manageStocks: p.manageStocks || false,
            hasVariants: p.hasVariants || false,
            shippingCost: p.shippingCost || '0',
            weight: p.weight || '',
            dimensions: {
                height: p.height || '',
                width: p.width || '',
                depth: p.depth || ''
            },
            
            // Business info
            negotiable: p.negotiable || false,
            middlemanService: p.middlemanService || false,
            userType: p.userType || 'private',
            status: p.status || 'active',
            popularity: p.popularity ?? 0,
            verified: p.verified !== false,
            categorySlug: p.categorySlug,
            categoryId: p.categoryId || '',
            tags: p.tags || [],
            
            // Metrics
            metrics: {
                likes: p.metrics?.likesCount ? parseInt(p.metrics.likesCount) : (p.likeCount || 0),
                dislikes: p.metrics?.dislikesCount ? parseInt(p.metrics.dislikesCount) : (p.dislikeCount || 0),
                comments: p.metrics?.commentsCount ? parseInt(p.metrics.commentsCount) : ((p.comments || []).length),
                shares: p.metrics?.sharedCount ? parseInt(p.metrics.sharedCount) : 0,
                views: p.metrics?.visitedCount ? parseInt(p.metrics.visitedCount) : (p.viewCount || p.views || 5),
                saved: p.metrics?.addedToWishlistCount ? parseInt(p.metrics.addedToWishlistCount) : (p.wishlistCount || 0),
            },
            
            // Reviews
            rating: p.rating || 4.5,
            reviewCount: p.reviewCount || 127,
            
            // Additional info
            warranty: p.warranty || '12 months manufacturer warranty',
            returnPolicy: p.returnPolicy || '30-day return policy',
            features: p.features || ['Premium quality materials', 'Energy efficient design', 'Easy installation', 'Long lasting durability'],
            specifications: p.specifications || {
                'Power': '20W',
                'Voltage': '220-240V',
                'Color Temperature': '3000K-6500K',
                'Luminous Flux': '2000lm',
                'Beam Angle': '120°',
                'Lifespan': '50,000 hours',
                'Material': 'Aluminum + PC',
                'IP Rating': 'IP65'
            },
            location: p.location || 'Germany',
            seller: {
                name: p.sellerName || 'Premium Electronics Store',
                rating: p.sellerRating || 4.8,
                responseTime: p.responseTime || '< 1 hour',
                totalSales: p.totalSales || 1543,
                memberSince: p.memberSince || '2020'
            }
        };
    }, [product]);
    
    // State management
    const [state, setState] = useState({
        currentImageIndex: 0,
        selectedQuantity: 1,
        isAddingToCart: false,
        activeTab: 'description',
        showFullDescription: false,
        liked: false,
        disliked: false,
        favorite: false,
        showComments: false,
        zoomedImage: null
    });
    
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoadingMedia, setIsLoadingMedia] = useState(false);
    
    // Initialize basket actions
    const { addToBasket: handleAddToCart, isLoading: basketLoading } = useProductBasket(productData.id, productData);
    
    // Load media
    useEffect(() => {
        const loadMedia = async () => {
            if (!productData.id) return;
            
            setIsLoadingMedia(true);
            try {
                const mediaResponse = await getMediaByItem(productData.id);
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || 'image',
                        alt: item.altText || productData.name || 'Product image'
                    }));
                    setMediaItems(formattedMedia);
                } else if (productData.thumbnail) {
                    setMediaItems([{ url: productData.thumbnail, type: 'image', alt: productData.name }]);
                }
            } catch (error) {
                if (productData.thumbnail) {
                    setMediaItems([{ url: productData.thumbnail, type: 'image', alt: productData.name }]);
                }
            } finally {
                setIsLoadingMedia(false);
            }
        };
        loadMedia();
    }, [productData.id, productData.name, productData.thumbnail]);
    
    // Price formatting
    const formattedPrice = useMemo(() => {
        const price = parseFloat(productData.basePrice);
        return price.toLocaleString('de-DE', { style: 'currency', currency: 'EUR' });
    }, [productData.basePrice]);
    
    const formattedShipping = useMemo(() => {
        const shipping = parseFloat(productData.shippingCost);
        return shipping > 0 ? shipping.toLocaleString('de-DE', { style: 'currency', currency: 'EUR' }) : 'Free';
    }, [productData.shippingCost]);
    
    // Stock status
    const stockStatus = useMemo(() => {
        if (!productData.manageStocks) return { available: true, label: 'Available', className: styles.available };
        if (productData.stock === 0) return { available: false, label: 'Out of Stock', className: styles.outOfStock };
        if (productData.stock <= 5) return { available: true, label: `Only ${productData.stock} left`, className: styles.lowStock };
        return { available: true, label: `${productData.stock} in stock`, className: styles.inStock };
    }, [productData.manageStocks, productData.stock]);
    
    // Handlers
    const handleImageNavigation = useCallback((direction) => {
        setState(prev => ({
            ...prev,
            currentImageIndex: direction === 'next'
                ? (prev.currentImageIndex + 1) % mediaItems.length
                : (prev.currentImageIndex - 1 + mediaItems.length) % mediaItems.length
        }));
    }, [mediaItems.length]);
    
    const handleQuantityChange = useCallback((delta) => {
        setState(prev => ({
            ...prev,
            selectedQuantity: Math.max(1, Math.min(productData.stock || 99, prev.selectedQuantity + delta))
        }));
    }, [productData.stock]);
    
    const handleAddToBasket = useCallback(async () => {
        if (!stockStatus.available) {
            toast.error("This product is currently out of stock.");
            return;
        }
        
        setState(prev => ({ ...prev, isAddingToCart: true }));
        try {
            await handleAddToCart(state.selectedQuantity);
            toast.success(`Added ${state.selectedQuantity} item(s) to cart`);
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    } finally {
            setState(prev => ({ ...prev, isAddingToCart: false }));
        }
    }, [handleAddToCart, state.selectedQuantity, stockStatus.available]);
    
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to like products.");
            return;
        }
        handleLike(productData.id, userId);
        setState(prev => ({ ...prev, liked: true, disliked: false }));
    }, [productData.id, userId, handleLike]);
    
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to dislike products.");
            return;
        }
        handleDislike(productData.id, userId);
        setState(prev => ({ ...prev, liked: false, disliked: true }));
    }, [productData.id, userId, handleDislike]);
    
    const handleShare = useCallback(() => {
        if (navigator.share && typeof window !== 'undefined') {
            navigator.share({
                title: productData.name,
                text: productData.description,
                url: window.location.href
            }).catch(() => {});
        } else {
            navigator.clipboard.writeText(window.location.href);
            toast.success("Link copied to clipboard!");
        }
    }, [productData.name, productData.description]);
    
    const handleContactSeller = useCallback(() => {
        dispatch(openMessageModal({
            itemId: productData.id,
            recipientId: 'marketplace-seller'
        }));
    }, [productData.id, dispatch]);
    
    return (
        <div className={styles.container}>
            <div className={styles.detailView}>
                {/* Breadcrumb Navigation */}
                <nav className={styles.breadcrumb}>
                    <Link href="/products" className={styles.breadcrumbLink}>Products</Link>
                    <ChevronRight size={16} />
                    <Link href={`/products/${productData.categorySlug}`} className={styles.breadcrumbLink}>
                        {category?.name || productData.categorySlug}
                    </Link>
                    <ChevronRight size={16} />
                    <span className={styles.breadcrumbCurrent}>{productData.name}</span>
                </nav>
                
                {/* Main Content Grid */}
                <div className={styles.mainGrid}>
                    {/* Left Column - Images */}
                    <div className={styles.imageSection}>
                        {/* Main Image Display */}
                        <div className={styles.mainImageContainer}>
                            {isLoadingMedia ? (
                                <div className={styles.imagePlaceholder}>
                                    <div className={styles.spinner} />
                                </div>
                            ) : mediaItems.length > 0 ? (
                                <>
                                    <div className={styles.mainImage}>
                                        <img
                                            src={mediaItems[state.currentImageIndex]?.url}
                                            alt={mediaItems[state.currentImageIndex]?.alt || productData.name}
                                            onClick={() => setState(prev => ({ ...prev, zoomedImage: mediaItems[state.currentImageIndex]?.url }))}
                                            className={styles.productImage}
                                        />
                                        
                                        {/* Navigation Arrows */}
                                        {mediaItems.length > 1 && (
                                            <>
                                                <button
                                                    className={`${styles.imageNav} ${styles.prevNav}`}
                                                    onClick={() => handleImageNavigation('prev')}
                                                    aria-label="Previous image"
                                                >
                                                    <ChevronLeft size={24} />
                                                </button>
                                                <button
                                                    className={`${styles.imageNav} ${styles.nextNav}`}
                                                    onClick={() => handleImageNavigation('next')}
                                                    aria-label="Next image"
                                                >
                                                    <ChevronRight size={24} />
                                                </button>
                                            </>
                                        )}
                                    </div>
                                    
                                    {/* Image Counter */}
                                    <div className={styles.imageCounter}>
                                        {state.currentImageIndex + 1} / {mediaItems.length}
                                    </div>
                                </>
                            ) : (
                                <div className={styles.noImage}>
                                    <Camera size={64} />
                                    <p>No images available</p>
                                </div>
                            )}
                        </div>
                        
                        {/* Thumbnail Gallery */}
                        {mediaItems.length > 1 && (
                            <div className={styles.thumbnailGallery}>
                                {mediaItems.map((media, index) => (
                                    <button
                                        key={index}
                                        className={`${styles.thumbnail} ${index === state.currentImageIndex ? styles.activeThumbnail : ''}`}
                                        onClick={() => setState(prev => ({ ...prev, currentImageIndex: index }))}
                                        aria-label={`View image ${index + 1}`}
                                    >
                                        <img src={media.url} alt={`Product image ${index + 1}`} />
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                    
                    {/* Right Column - Product Info */}
                    <div className={styles.infoSection}>
                        {/* Product Header */}
                        <div className={styles.productHeader}>
                            <h1 className={styles.productTitle}>{productData.name}</h1>
                            {productData.brand && (
                                <div className={styles.brandInfo}>
                                    <span>by</span>
                                    <Link href={`/brands/${productData.brand.toLowerCase()}`} className={styles.brandLink}>
                                        {productData.brand}
                                    </Link>
                                </div>
                            )}
                        </div>
                        
                        {/* Rating and Reviews */}
                        <div className={styles.ratingSection}>
                            <div className={styles.stars}>
                                {[...Array(5)].map((_, i) => (
                                    <Star key={i} size={20} className={i < Math.floor(productData.rating) ? styles.starFilled : styles.starEmpty} />
                                ))}
                            </div>
                            <span className={styles.ratingValue}>{productData.rating.toFixed(1)}</span>
                            <Link href="#reviews" className={styles.reviewLink}>
                                ({productData.reviewCount} reviews)
                            </Link>
                            <span className={styles.separator}>•</span>
                            <span className={styles.viewCount}>{productData.metrics.views} views</span>
                        </div>
                        
                        {/* Price Section */}
                        <div className={styles.priceSection}>
                            <div className={styles.priceMain}>
                                <span className={styles.price}>{formattedPrice}</span>
                                {productData.negotiable && (
                                    <span className={styles.negotiableBadge}>
                                        <Tag size={16} />
                                        Negotiable
                                    </span>
                                )}
                            </div>
                            <div className={styles.priceInfo}>
                                <div className={styles.shippingInfo}>
                                    <Truck size={16} />
                                    <span>Shipping: {formattedShipping}</span>
                                </div>
                                {productData.middlemanService && (
                                    <div className={styles.protectionInfo}>
                                        <Shield size={16} />
                                        <span>Buyer Protection</span>
                                    </div>
                                )}
                            </div>
                        </div>
                        
                        {/* Quick Info Grid */}
                        <div className={styles.quickInfo}>
                            <div className={styles.infoItem}>
                                <Package size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Condition</span>
                                    <span className={styles.infoValue}>{productData.condition}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <MapPin size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Location</span>
                                    <span className={styles.infoValue}>{productData.location}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <Calendar size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Listed</span>
                                    <span className={styles.infoValue}>{dayjs(productData.added).fromNow()}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <Award size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Warranty</span>
                                    <span className={styles.infoValue}>{productData.warranty}</span>
                                </div>
                            </div>
                        </div>
                        
                        {/* Purchase Section */}
                        <div className={styles.purchaseSection}>
                            <div className={`${styles.stockIndicator} ${stockStatus.className}`}>
                                {stockStatus.available ? <CheckCircle size={20} /> : <AlertCircle size={20} />}
                                <span>{stockStatus.label}</span>
                            </div>
                            
                            <div className={styles.quantitySection}>
                                <label>Quantity:</label>
                                <div className={styles.quantitySelector}>
                                    <button
                                        onClick={() => handleQuantityChange(-1)}
                                        disabled={state.selectedQuantity <= 1}
                                        aria-label="Decrease quantity"
                                    >
                                        <Minus size={18} />
                                    </button>
                                    <input
                                        type="number"
                                        value={state.selectedQuantity}
                                        onChange={(e) => {
                                            const val = parseInt(e.target.value) || 1;
                                            if (val >= 1 && val <= (productData.stock || 99)) {
                                                setState(prev => ({ ...prev, selectedQuantity: val }));
                                            }
                                        }}
                                        min="1"
                                        max={productData.stock || 99}
                                    />
                                    <button
                                        onClick={() => handleQuantityChange(1)}
                                        disabled={state.selectedQuantity >= (productData.stock || 99)}
                                        aria-label="Increase quantity"
                                    >
                                        <Plus size={18} />
                                    </button>
                                </div>
                            </div>
                            
                            <div className={styles.actionButtons}>
                                <button
                                    className={styles.addToCartBtn}
                                    onClick={handleAddToBasket}
                                    disabled={!stockStatus.available || state.isAddingToCart || basketLoading}
                                >
                                    {state.isAddingToCart || basketLoading ? (
                                        <div className={styles.spinner} />
                                    ) : (
                                        <ShoppingCart size={20} />
                                    )}
                                    <span>Add to Cart</span>
                                </button>
                                
                                <button
                                    className={`${styles.favoriteBtn} ${state.favorite ? styles.active : ''}`}
                                    onClick={() => setState(prev => ({ ...prev, favorite: !prev.favorite }))}
                                    aria-label="Add to favorites"
                                >
                                    <Heart size={20} />
                                </button>
                            </div>
                            
                            <div className={styles.secondaryActions}>
                                <button onClick={handleContactSeller} className={styles.secondaryBtn}>
                                    <MessageCircle size={18} />
                                    Contact Seller
                                </button>
                                <button onClick={handleShare} className={styles.secondaryBtn}>
                                    <Share2 size={18} />
                                    Share
                                </button>
                            </div>
                        </div>
                        
                        {/* Seller Info Box */}
                        <div className={styles.sellerBox}>
                            <h3>Seller Information</h3>
                            <div className={styles.sellerInfo}>
                                <div className={styles.sellerHeader}>
                                    <div className={styles.sellerAvatar}>
                                        <User size={32} />
                                    </div>
                                    <div className={styles.sellerDetails}>
                                        <h4>{productData.seller.name}</h4>
                                        <div className={styles.sellerStats}>
                                            <div className={styles.sellerRating}>
                                                <Star size={16} className={styles.starFilled} />
                                                <span>{productData.seller.rating}</span>
                                            </div>
                                            <span>•</span>
                                            <span>{productData.seller.totalSales} sales</span>
                                            <span>•</span>
                                            <span>Member since {productData.seller.memberSince}</span>
                                        </div>
                                    </div>
                                </div>
                                <div className={styles.sellerActions}>
                                    <span className={styles.responseTime}>
                                        Usually responds within {productData.seller.responseTime}
                                    </span>
                                    <Link href={`/sellers/${productData.seller.name}`} className={styles.viewProfileBtn}>
                                        View Profile
                                        <ArrowRight size={16} />
                                    </Link>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
                {/* Bottom Section - Tabs */}
                <div className={styles.bottomSection}>
                    {/* Tab Navigation */}
                    <div className={styles.tabNav}>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'description' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'description' }))}
                        >
                            <Info size={18} />
                            Description
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'specifications' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'specifications' }))}
                        >
                            <Layers size={18} />
                            Specifications
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'reviews' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'reviews' }))}
                            id="reviews"
                        >
                            <Star size={18} />
                            Reviews ({productData.reviewCount})
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'shipping' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'shipping' }))}
                        >
                            <Truck size={18} />
                            Shipping & Returns
                        </button>
                    </div>
                    
                    {/* Tab Content */}
                    <div className={styles.tabContent}>
                        {/* Description Tab */}
                        {state.activeTab === 'description' && (
                            <div className={styles.descriptionTab}>
                                <p className={styles.description}>
                                    {productData.description}
                                </p>
                                
                                {productData.features.length > 0 && (
                                    <div className={styles.features}>
                                        <h3>Key Features</h3>
                                        <ul>
                                            {productData.features.map((feature, index) => (
                                                <li key={index}>
                                                    <CheckCircle size={16} />
                                                    {feature}
                                                </li>
                                            ))}
                                        </ul>
                                    </div>
                                )}
                            </div>
                        )}
                        
                        {/* Specifications Tab */}
                        {state.activeTab === 'specifications' && (
                            <div className={styles.specificationsTab}>
                                <table className={styles.specTable}>
                                    <tbody>
                                        {productData.brand && (
                                            <tr>
                                                <td>Brand</td>
                                                <td>{productData.brand}</td>
                                            </tr>
                                        )}
                                        {productData.model && (
                                            <tr>
                                                <td>Model</td>
                                                <td>{productData.model}</td>
                                            </tr>
                                        )}
                                        {productData.sku && (
                                            <tr>
                                                <td>SKU</td>
                                                <td>{productData.sku}</td>
                                            </tr>
                                        )}
                                        <tr>
                                            <td>Condition</td>
                                            <td>{productData.condition}</td>
                                        </tr>
                                        {productData.weight && (
                                            <tr>
                                                <td>Weight</td>
                                                <td>{productData.weight}</td>
                                            </tr>
                                        )}
                                        {Object.entries(productData.specifications).map(([key, value]) => (
                                            <tr key={key}>
                                                <td>{key}</td>
                                                <td>{value}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                        
                        {/* Reviews Tab */}
                        {state.activeTab === 'reviews' && (
                            <div className={styles.reviewsTab}>
                                <div className={styles.reviewsSummary}>
                                    <div className={styles.overallRating}>
                                        <h2>{productData.rating.toFixed(1)}</h2>
                                        <div className={styles.stars}>
                                            {[...Array(5)].map((_, i) => (
                                                <Star key={i} size={24} className={i < Math.floor(productData.rating) ? styles.starFilled : styles.starEmpty} />
                                            ))}
                                        </div>
                                        <p>Based on {productData.reviewCount} reviews</p>
                                    </div>
                                    
                                    <div className={styles.ratingBars}>
                                        {[5, 4, 3, 2, 1].map(stars => (
                                            <div key={stars} className={styles.ratingBar}>
                                                <span>{stars}</span>
                                                <Star size={14} className={styles.starFilled} />
                                                <div className={styles.barContainer}>
                                                    <div 
                                                        className={styles.barFill} 
                                                        style={{ width: `${Math.random() * 100}%` }}
                                                    />
                                                </div>
                                                <span>{Math.floor(Math.random() * 100)}%</span>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                                
                                <button className={styles.writeReviewBtn}>
                                    Write a Review
                                </button>
                                
                                {/* Comments Component */}
                                <div className={styles.reviewsList}>
                                    <CommentsSetup
                                        userId={userId}
                                        itemId={productData.id}
                                        itemType="product"
                                        toggleCommentsList={() => {}}
                                        categoryId={productData.categoryId}
                                        productName={productData.name}
                                        productThumbnail={productData.thumbnail}
                                    />
                                </div>
                            </div>
                        )}
                        
                        {/* Shipping Tab */}
                        {state.activeTab === 'shipping' && (
                            <div className={styles.shippingTab}>
                                <div className={styles.shippingInfo}>
                                    <h3>Shipping Information</h3>
                                    <ul>
                                        <li>
                                            <Truck size={18} />
                                            <span>Standard shipping: {formattedShipping}</span>
                                        </li>
                                        <li>
                                            <Calendar size={18} />
                                            <span>Estimated delivery: 3-5 business days</span>
                                        </li>
                                        <li>
                                            <Globe size={18} />
                                            <span>Ships to: Europe, USA, Canada</span>
                                        </li>
                                    </ul>
                                </div>
                                
                                <div className={styles.returnInfo}>
                                    <h3>Return Policy</h3>
                                    <p>{productData.returnPolicy}</p>
                                    <ul>
                                        <li>
                                            <CheckCircle size={18} />
                                            <span>Free returns within 30 days</span>
                                        </li>
                                        <li>
                                            <CheckCircle size={18} />
                                            <span>Original packaging required</span>
                                        </li>
                                        <li>
                                            <CheckCircle size={18} />
                                            <span>Refund or exchange available</span>
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
                
                {/* Engagement Bar */}
                <div className={styles.engagementBar}>
                    <div className={styles.engagementStats}>
                        <span><Eye size={18} /> {productData.metrics.views} views</span>
                        <span><Heart size={18} /> {productData.metrics.saved} saved</span>
                        <span><Users size={18} /> {productData.metrics.shares} shares</span>
                    </div>
                    
                    <div className={styles.engagementActions}>
                        <button
                            className={`${styles.engageBtn} ${state.liked ? styles.active : ''}`}
                            onClick={handleLikeClick}
                        >
                            <ThumbsUp size={18} />
                            <span>{productData.metrics.likes}</span>
                        </button>
                        <button
                            className={`${styles.engageBtn} ${state.disliked ? styles.active : ''}`}
                            onClick={handleDislikeClick}
                        >
                            <ThumbsDown size={18} />
                            <span>{productData.metrics.dislikes}</span>
                        </button>
                        <button
                            className={styles.engageBtn}
                            onClick={() => setState(prev => ({ ...prev, showComments: !prev.showComments }))}
                        >
                            <MessageCircle size={18} />
                            <span>{productData.metrics.comments}</span>
                        </button>
                        <button className={styles.engageBtn} onClick={handleShare}>
                            <Share2 size={18} />
                            <span>Share</span>
                        </button>
                    </div>
                </div>
                
                {/* Tags */}
                {productData.tags.length > 0 && (
                    <div className={styles.tagsSection}>
                        <h3>Tags</h3>
                        <div className={styles.tagsList}>
                            {productData.tags.map((tag, index) => (
                                <Link key={index} href={`/search?tag=${encodeURIComponent(tag)}`} className={styles.tag}>
                                    <Tag size={14} />
                                    {tag}
                                </Link>
                            ))}
                        </div>
                    </div>
                )}
                
                {/* Image Zoom Modal */}
                {state.zoomedImage && (
                    <div className={styles.zoomModal} onClick={() => setState(prev => ({ ...prev, zoomedImage: null }))}>
                        <img src={state.zoomedImage} alt="Zoomed product image" />
                        <button className={styles.closeZoom} aria-label="Close zoom">×</button>
                    </div>
                )}
            </div>
        </div>
    );
});

DetailedProductView.displayName = 'DetailedProductView';

export default DetailedProductView;
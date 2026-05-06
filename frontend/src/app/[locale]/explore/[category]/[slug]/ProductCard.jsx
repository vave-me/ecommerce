"use client";
import React, {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl'; // Import hook
import dayjs from 'dayjs'; // For date formatting
import relativeTime from 'dayjs/plugin/relativeTime'; // For relative time
import 'dayjs/locale/en'; // English locale
import 'dayjs/locale/pl'; // Polish locale  
import 'dayjs/locale/de'; // German locale
// Enable relative time plugin
dayjs.extend(relativeTime);
import {
    Bookmark, Camera, Check, ChevronDown, ChevronUp, Clock, ExternalLink, Eye, Flag, Flame, Heart,
    MapPin, MessageCircle, MessagesSquare, Share2, ShieldCheck, Tag, ThumbsDown, ThumbsUp
} from '@/icons';
import styles from './ProductCard.module.css';
// Map locales for dayjs  
const dateLocales = {en: 'en', pl: 'pl', de: 'de'};
// --- Component ---
const ProductCard = ({product: productProp, locale = 'en'}) => {
    const t = useTranslations('ProductCard'); // Hook for translations
    const common_t = useTranslations('common'); // Hook for common translations
    // Handle missing product data
    if (!productProp) {
        return (
            <div className={styles.classifiedCard}>
                <div className={styles.contentSection}>
                    <p>Product not found or failed to load.</p>
                </div>
            </div>
        );
    }
    // Extract and process real API data
    const product = productProp;
    // --- State ---
    const [expanded, setExpanded] = useState(false);
    const [saved, setSaved] = useState(false);
    const [favorite, setFavorite] = useState(false);
    const [liked, setLiked] = useState(false);
    const [disliked, setDisliked] = useState(false);
    const [currentImageIndex, setCurrentImageIndex] = useState(0);
    // Extract metrics safely
    const views = parseInt(product.metrics?.viewsCount || product.viewsCount || '0', 10);
    const savedCount = parseInt(product.metrics?.savedCount || product.savedCount || '0', 10);
    const inquiries = parseInt(product.metrics?.inquiriesCount || product.inquiriesCount || '0', 10);
    const likes = parseInt(product.metrics?.likesCount || product.likesCount || '0', 10);
    const dislikes = parseInt(product.metrics?.dislikesCount || product.dislikesCount || '0', 10);
    const comments = parseInt(product.metrics?.commentsCount || product.commentsCount || '0', 10);
    // State for counts (allow optimistic updates)
    const [likesCount, setLikesCount] = useState(likes);
    const [dislikesCount, setDislikesCount] = useState(dislikes);
    const [savedCountState, setSavedCountState] = useState(savedCount);
    const [commentsCount, setCommentsCount] = useState(comments);
    const [messagesCount, setMessagesCount] = useState(inquiries);
    // --- Derived Data & Formatting ---
    const images = product.media?.length > 0 ? product.media : 
                  product.images?.length > 0 ? product.images : 
                  [product.thumbnail || '/api/placeholder/500/300']; // Fallback image
    // Relative Time Formatting
    const [postedTimeAgo, setPostedTimeAgo] = useState('');
    useEffect(() => {
        const postDate = new Date(product.postedDate || product.createdAt || Date.now());
        try {
            // Set dayjs locale and format relative time
            dayjs.locale(dateLocales[locale] || 'en');
            setPostedTimeAgo(dayjs(postDate).fromNow());
        } catch (e) {
            setPostedTimeAgo('Recently');
        }
    }, [product.postedDate, product.createdAt, locale]);
    // User Joined Date Formatting (Example: "Jan 2023")
    const [userJoinedFormatted, setUserJoinedFormatted] = useState('');
    useEffect(() => {
        try {
            const joinedDate = new Date(product.seller?.joinedDate || product.seller?.createdAt || Date.now());
            setUserJoinedFormatted(new Intl.DateTimeFormat(locale, {
                month: 'short',
                year: 'numeric'
            }).format(joinedDate));
        } catch (e) {
            setUserJoinedFormatted('');
        }
    }, [product.seller?.joinedDate, product.seller?.createdAt, locale]);
    // Price Formatting (use Intl.NumberFormat)
    const formattedPrice = new Intl.NumberFormat(locale, {
        style: 'currency',
        currency: common_t('currencyCode') || 'EUR',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(parseFloat(product.basePrice || product.price || 0));
    // Condition Text (using common translations)
    const conditionText = product.condition || 'Used';
    // Calculate smart heat score based on engagement metrics
    const heatScore = calculateHeatScore(views, likes, comments, savedCount);
    // Heat Level
    const getHeatLevelClass = (score) => {
        if (score >= 80) return styles.heatIndicatorHot;
        if (score >= 60) return styles.heatIndicatorWarm;
        if (score >= 40) return styles.heatIndicatorGood;
        if (score >= 20) return styles.heatIndicatorCool;
        return styles.heatIndicatorCold;
    };
    const getHeatLevel = (score) => {
        if (score >= 80) return 'hot';
        if (score >= 60) return 'warm';
        if (score >= 40) return 'good';
        if (score >= 20) return 'cool';
        return 'cold';
    };
    const heatLabel = getHeatLevel(heatScore);
    // Smart heat score calculation
    function calculateHeatScore(views, likes, comments, saved) {
        const viewWeight = views * 0.1;
        const likeWeight = likes * 4;
        const commentWeight = comments * 3;
        const savedWeight = saved * 5;
        const rawScore = viewWeight + likeWeight + commentWeight + savedWeight;
        return Math.min(Math.round(rawScore / 15), 100);
    }
    // Extract seller information
    const seller = {
        username: product.seller?.username || product.seller?.name || product.username || 'Anonymous Seller',
        rating: product.seller?.rating || product.userRating || 4.0,
        verified: product.seller?.verified || product.verified || false,
        dealCount: product.seller?.dealCount || product.seller?.transactionCount || 0
    };
    // Extract product details
    const productDetails = {
        brand: product.brand || '',
        model: product.model || '',
        year: product.year || product.yearManufactured || '',
        condition: product.condition || 'used',
        location: product.location || product.address || 'Location not specified',
        category: product.category?.name || product.categoryName || 'General',
        subcategory: product.subcategory?.name || product.subcategoryName || ''
    };
    // Extract tags from description or use existing tags
    function extractTags(product) {
        const tags = [];
        // Use existing tags if available
        if (product.tags && Array.isArray(product.tags)) {
            return product.tags.slice(0, 4);
        }
        // Extract from product properties
        if (product.brand) tags.push(product.brand);
        if (product.condition && product.condition !== 'used') tags.push(product.condition);
        if (product.negotiable) tags.push('Negotiable');
        if (product.warranty) tags.push('Warranty');
        return tags.slice(0, 4);
    }
    const tags = extractTags(product);
    // --- Handlers (Placeholder logic - implement API calls) ---
    const handleSaveToggle = () => {
        const newState = !saved;
        setSaved(newState);
        setSavedCountState(newState ? savedCountState + 1 : Math.max(0, savedCountState - 1));
        // API CALL: updateSaveStatus(product.id, newState);
    };
    const handleFavoriteToggle = () => {
        const newState = !favorite;
        setFavorite(newState);
        // API CALL: updateFavoriteStatus(product.id, newState);
    };
    const handleLikeToggle = () => {
        const newState = !liked;
        setLiked(newState);
        setLikesCount(newState ? likesCount + 1 : Math.max(0, likesCount - 1));
        if (newState && disliked) { // Un-dislike if liking
            setDisliked(false);
            setDislikesCount(Math.max(0, dislikesCount - 1));
        }
        // API CALL: updateLikeStatus(product.id, newState);
    };
    const handleDislikeToggle = () => {
        const newState = !disliked;
        setDisliked(newState);
        setDislikesCount(newState ? dislikesCount + 1 : Math.max(0, dislikesCount - 1));
        if (newState && liked) { // Un-like if disliking
            setLiked(false);
            setLikesCount(Math.max(0, likesCount - 1));
        }
        // API CALL: updateDislikeStatus(product.id, newState);
    };
    const nextImage = () => setCurrentImageIndex((prev) => (prev === images.length - 1 ? 0 : prev + 1));
    const prevImage = () => setCurrentImageIndex((prev) => (prev === 0 ? images.length - 1 : prev - 1));
    // --- Render ---
    return (
        <div className={styles.classifiedCard} aria-labelledby={`product-title-${product.id}`}>
            <div className={`${styles.heatIndicator} ${getHeatLevelClass(heatScore)}`}></div>
            <div className={styles.imageSection}>
                {product.isPromoted && (
                    <div className={styles.promotedBadge}>Promoted</div>
                )}
                <div className={styles.imageContainer}>
                    <img
                        src={images[currentImageIndex]}
                        alt={`${product.name || product.title} - Image ${currentImageIndex + 1}`}
                        className={styles.classifiedImage}
                    />
                    {/* Badges */}
                    <div className={`${styles.badge} ${styles.imageCounterBadge}`}>
                        <Camera size={14} className={styles.iconCamera} strokeWidth={2.5} aria-hidden="true"/>
                        <span>{currentImageIndex + 1}/{images.length}</span>
                    </div>
                    <div className={`${styles.badge} ${styles.timeBadge}`}>
                        <Clock size={14} className={styles.iconClock} strokeWidth={2.5} aria-hidden="true"/>
                        <time dateTime={product.postedDate || product.createdAt}>{postedTimeAgo}</time>
                    </div>
                    <div className={`${styles.badge} ${styles.conditionBadge}`}>
                        <span>{conditionText}</span>
                    </div>
                    <div className={`${styles.badge} ${styles.heatBadge}`}>
                        <Flame size={14} className={styles.iconFlame} strokeWidth={2.5} aria-hidden="true"/>
                        <span>{heatLabel}</span>
                    </div>
                    {/* Image Navigation */}
                    {images.length > 1 && (
                        <>
                            <button onClick={prevImage}
                                    className={`${styles.imageNavButton} ${styles.imageNavButtonPrev}`}
                                    aria-label="Previous image">
                                {/* SVG */}
                            </button>
                            <button onClick={nextImage}
                                    className={`${styles.imageNavButton} ${styles.imageNavButtonNext}`}
                                    aria-label="Next image">
                                {/* SVG */}
                            </button>
                        </>
                    )}
                    {/* Interaction Column */}
                    <div className={styles.interactionColumn}>
                        <button onClick={handleFavoriteToggle} className={styles.interactionButton}
                                aria-pressed={favorite}
                                aria-label={favorite ? 'Remove from favorites' : 'Add to favorites'}>
                            <Heart size={18} className={`${styles.iconHeart} ${favorite ? styles.active : ''}`}
                                   strokeWidth={2} aria-hidden="true"/>
                        </button>
                        <button onClick={handleLikeToggle} className={styles.interactionButton} aria-pressed={liked}
                                aria-label={liked ? 'Unlike' : 'Like'}>
                            <ThumbsUp size={18} className={`${styles.iconThumbsUp} ${liked ? styles.active : ''}`}
                                      strokeWidth={2} aria-hidden="true"/>
                            <span>{likesCount}</span>
                        </button>
                        <button onClick={handleDislikeToggle} className={styles.interactionButton}
                                aria-pressed={disliked}
                                aria-label={disliked ? 'Remove dislike' : 'Dislike'}>
                            <ThumbsDown size={18}
                                        className={`${styles.iconThumbsDown} ${disliked ? styles.active : ''}`}
                                        strokeWidth={2} aria-hidden="true"/>
                            <span>{dislikesCount}</span>
                        </button>
                        <button className={styles.interactionButton}
                                aria-label={`View comments (${commentsCount})`}>
                            <MessagesSquare size={18} className={styles.iconComment} strokeWidth={2}
                                            aria-hidden="true"/>
                            <span>{commentsCount}</span>
                        </button>
                        <button className={styles.interactionButton}
                                aria-label={`Send message (${messagesCount})`}>
                            <MessageCircle size={18} className={styles.iconMessage} strokeWidth={2} aria-hidden="true"/>
                            <span>{messagesCount}</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className={styles.contentSection}>
                <div className={styles.titleSection}>
                    <div className={styles.titleRow}>
                        <h3 id={`product-title-${product.id}`}
                            className={styles.adTitle}>{product.name || product.title}</h3>
                        <div className={styles.viewsBadge} title={`${views} views`}>
                            <Eye size={12} className={styles.iconEye} aria-hidden="true"/>
                            {views}
                        </div>
                    </div>
                    <div className={styles.priceRow}>
                        <div className={styles.price}>
                            <div>{formattedPrice}</div>
                            {product.negotiable && (
                                <div className={styles.negotiableBadge}>Negotiable</div>
                            )}
                        </div>
                        <button className={styles.viewDetailsButton}>
                            <ExternalLink size={14} className={styles.iconExternalLink} strokeWidth={2}
                                          aria-hidden="true"/>
                            View Details
                        </button>
                    </div>
                </div>
                <div className={styles.metaInfo}>
                    <div className={styles.metaBadge}>
                        <MapPin size={12} className={styles.iconMapPin} strokeWidth={2} aria-hidden="true"/>
                        <span>{productDetails.location}</span>
                    </div>
                    <div className={styles.metaBadge}>
                        <Tag size={12} className={styles.iconTag} strokeWidth={2} aria-hidden="true"/>
                        <span>{productDetails.category}{productDetails.subcategory ? ` • ${productDetails.subcategory}` : ''}</span>
                    </div>
                </div>
                {tags.length > 0 && (
                    <div className={styles.tagsSection}>
                        {tags.map((tag, index) => (
                            <span key={index} className={styles.tag}>{tag}</span>
                        ))}
                    </div>
                )}
                <div className={styles.descriptionSection}>
                    <div className={styles.descriptionHeader}>
                        <div className={styles.descriptionTitle}>
                            <span>Description</span>
                            <span className={styles.conditionLabel}>{conditionText}</span>
                        </div>
                        <button onClick={() => setExpanded(!expanded)} className={styles.expandButton}
                                aria-expanded={expanded}
                                aria-label={expanded ? 'Show less' : 'Show more'}>
                            {expanded ? <ChevronUp size={16} strokeWidth={2.5} aria-hidden="true"/> :
                                <ChevronDown size={16} strokeWidth={2.5} aria-hidden="true"/>}
                        </button>
                    </div>
                    <div className={`${styles.descriptionText} ${expanded ? '' : styles.collapsed}`}>
                        {product.description || 'No description available.'}
                    </div>
                    {/* Expanded Details Section */}
                    {expanded && (
                        <div className={styles.descriptionDetails}>
                            {productDetails.brand && <DetailRow label="Brand" value={productDetails.brand}/>}
                            {productDetails.model && <DetailRow label="Model" value={productDetails.model}/>}
                            {productDetails.year && <DetailRow label="Year" value={productDetails.year}/>}
                            {product.sku && <DetailRow label="SKU" value={product.sku}/>}
                            {product.warranty && <DetailRow label="Warranty" value={product.warranty}/>}
                        </div>
                    )}
                </div>
                {/* Seller Info */}
                <div className={styles.sellerSection}>
                    <div className={styles.sellerInfo}>
                        <div className={styles.sellerAvatar}>
                            <div className={styles.avatarCircle}>
                                {seller.username.charAt(0).toUpperCase()}
                            </div>
                            {seller.verified && (
                                <div className={styles.verifiedBadge} title="Verified Seller">
                                    <Check size={12} className={styles.iconCheck} aria-hidden="true"/>
                                </div>
                            )}
                        </div>
                        <div className={styles.sellerDetails}>
                            <div className={styles.sellerName}>
                                {seller.username}
                                {seller.verified && (
                                    <span className={styles.verifiedLabel}>
                                         <ShieldCheck size={12} className={styles.iconShield} aria-hidden="true"/>
                                        Verified
                                     </span>
                                )}
                            </div>
                            <div className={styles.sellerRating} title={`${seller.rating} rating, ${seller.dealCount} deals`}>
                                <div className={styles.ratingStars} aria-label={`${seller.rating} star rating`}>
                                    {[...Array(5)].map((_, i) => (
                                        <svg key={i}
                                             className={`${styles.star} ${i < Math.floor(seller.rating) ? styles.starFilled : styles.starEmpty}`}
                                             fill="currentColor" viewBox="0 0 20 20">
                                            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                                        </svg>
                                    ))}
                                </div>
                                <span className={styles.ratingText}>
                                    {seller.rating} ({seller.dealCount} deal{seller.dealCount !== 1 ? 's' : ''})
                                </span>
                            </div>
                        </div>
                    </div>
                    <button className={styles.contactButton}>
                        <MessageCircle size={14} className={styles.iconMessageCircle} strokeWidth={2}
                                       aria-hidden="true"/>
                        Contact Seller
                    </button>
                </div>
                {/* Bottom Action Buttons */}
                <div className={styles.actionButtons}>
                    <button className={styles.actionButton} onClick={handleSaveToggle} aria-pressed={saved}>
                        <Bookmark size={16} className={`${styles.bookmarkIcon} ${saved ? styles.active : ''}`}
                                  strokeWidth={2} aria-hidden="true"/>
                        <span>{saved ? 'Saved' : 'Save'}</span>
                    </button>
                    <button className={styles.actionButton}>
                        <Share2 size={16} strokeWidth={2} aria-hidden="true"/>
                        <span>Share</span>
                    </button>
                    <button className={styles.actionButton}>
                        <Flag size={16} strokeWidth={2} aria-hidden="true"/>
                        <span>Report</span>
                    </button>
                </div>
            </div>
        </div>
    );
};
// Helper component for expanded details
const DetailRow = ({label, value}) => (
    <div className={styles.detailRow}>
        <span className={styles.detailLabel}>{label}:</span>
        <span className={styles.detailValue}>{value}</span>
    </div>
);
export default ProductCard;
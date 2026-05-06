"use client";
import React, {useState, useEffect, memo} from 'react';
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
// --- Default Data Structure (Example) ---
const getDefaultProduct = (t) => ({
    images: ["/images/default-product.webp", "/api/placeholder/500/300"],
    title: t('defaultTitle'),
    name: t('defaultTitle'), // Ensure name exists for fallback
    basePrice: "350.00", // Use string for consistency with API?
    price: "€350", // Display price (formatted later)
    negotiable: true,
    condition: "used", // Use key for lookup
    location: t('defaultLocation'),
    seller: { // Assuming seller info might be nested
        username: t('defaultUsername'),
        rating: 4.8,
        verified: true,
        joinedDate: new Date(Date.now() - 365 * 2 * 24 * 60 * 60 * 1000).toISOString(), // Example: 2 years ago
        dealCount: 12 // Example deal count
    },
    postedDate: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(), // Example: 2 days ago
    category: {name: t('defaultCategory'), slug: 'sports-outdoors'}, // Example structure
    subcategory: {name: t('defaultSubcategory'), slug: 'bicycles'}, // Example structure
    description: t('defaultDescription'),
    stats: { // Assuming stats are nested
        views: 243,
        savedCount: 18, // bookmarks/saves
        inquiries: 7, // messages?
        likes: 34,
        dislikes: 2,
        comments: 15,
    },
    isFavorite: false, // User specific
    isPromoted: true,
    tags: t('defaultTags').split(', '), // Create array from translated string
    heatScore: 82,
    // Add potentially missing fields used in component logic
    id: 'preview-product-123',
    brand: 'Trek',
    frameSize: 'Medium (17.5")',
    wheelSize: '29 inch',
    year: 2023,
    inStock: true, // Example availability
});
// --- Component ---
// Accept product data and locale as props
const ProductCard = memo(({product: productProp, locale = 'en'}) => {
    const t = useTranslations('ProductCard'); // Hook for translations
    const common_t = useTranslations('common'); // Hook for common translations
    // Use passed product or generate default using translations
    const product = productProp || getDefaultProduct(t);
    // --- State ---
    const [expanded, setExpanded] = useState(false);
    const [saved, setSaved] = useState(product.isSaved || false); // Initialize from prop if available
    const [favorite, setFavorite] = useState(product.isFavorite || false); // Initialize from prop
    const [liked, setLiked] = useState(product.isLiked || false); // Initialize from prop
    const [disliked, setDisliked] = useState(product.isDisliked || false); // Initialize from prop
    const [currentImageIndex, setCurrentImageIndex] = useState(0);
    // State for counts (allow optimistic updates)
    const [likesCount, setLikesCount] = useState(product.stats?.likes || 0);
    const [dislikesCount, setDislikesCount] = useState(product.stats?.dislikes || 0);
    const [savedCount, setSavedCount] = useState(product.stats?.savedCount || 0);
    const [commentsCount, setCommentsCount] = useState(product.stats?.comments || 0);
    const [messagesCount, setMessagesCount] = useState(product.stats?.inquiries || 0); // Map inquiries to messages
    // --- Derived Data & Formatting ---
    const images = product.images?.length > 0 ? product.images : ['/api/placeholder/500/300']; // Fallback image
    // Relative Time Formatting
    const [postedTimeAgo, setPostedTimeAgo] = useState('');
    useEffect(() => {
        const postDate = new Date(product.postedDate || Date.now());
        try {
            // Set dayjs locale and format relative time
            dayjs.locale(dateLocales[locale] || 'en');
            setPostedTimeAgo(dayjs(postDate).fromNow());
        } catch (e) {
            setPostedTimeAgo(t('invalidDate'));
        }
    }, [product.postedDate, locale, t]);
    // User Joined Date Formatting (Example: "Jan 2023")
    const [userJoinedFormatted, setUserJoinedFormatted] = useState('');
    useEffect(() => {
        try {
            const joinedDate = new Date(product.seller?.joinedDate || Date.now());
            setUserJoinedFormatted(new Intl.DateTimeFormat(locale, {
                month: 'short',
                year: 'numeric'
            }).format(joinedDate));
        } catch (e) {
            setUserJoinedFormatted('');
        }
    }, [product.seller?.joinedDate, locale]);
    // Price Formatting (use Intl.NumberFormat)
    const formattedPrice = new Intl.NumberFormat(locale, {
        style: 'currency',
        currency: common_t('currencyCode'), // Use common currency code
        minimumFractionDigits: 2, // Adjust as needed
        maximumFractionDigits: 2,
    }).format(parseFloat(product.basePrice || 0)); // Ensure basePrice is a number
    // Condition Text (using common translations)
    const conditionText = common_t(`condition_${product.condition || 'unknown'}`, {}, {defaultValue: product.condition});
    // Heat Level
    const getHeatLevelClass = (score) => {
        if (score >= 80) return styles.heatIndicatorHot;
        if (score >= 60) return styles.heatIndicatorWarm;
        if (score >= 40) return styles.heatIndicatorGood;
        if (score >= 20) return styles.heatIndicatorCool;
        return styles.heatIndicatorCold;
    };
    const heatScore = product.heatScore || 0;
    const heatLabel = t(`heatLabel_${getHeatLevelClass(heatScore).split('-').pop()}`); // Get label based on class suffix
    // --- Handlers (Placeholder logic - implement API calls) ---
    const handleSaveToggle = () => {
        const newState = !saved;
        setSaved(newState);
        setSavedCount(newState ? savedCount + 1 : savedCount - 1);
        // API CALL: updateSaveStatus(product.id, newState);
    };
    const handleFavoriteToggle = () => {
        const newState = !favorite;
        setFavorite(newState);
        // Favorites usually don't have a public count, adjust if needed
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
                    <div className={styles.promotedBadge}>{t('promotedBadge')}</div>
                )}
                <div className={styles.imageContainer}>
                    <img
                        src={images[currentImageIndex]}
                        alt={t('imageAlt', {title: product.name || product.title, index: currentImageIndex + 1})}
                        className={styles.classifiedImage}
                    />
                    {/* Badges */}
                    <div className={`${styles.badge} ${styles.imageCounterBadge}`}>
                        <Camera size={14} className={styles.iconCamera} strokeWidth={2.5} aria-hidden="true"/>
                        <span>{currentImageIndex + 1}/{images.length}</span>
                    </div>
                    <div className={`${styles.badge} ${styles.timeBadge}`}>
                        <Clock size={14} className={styles.iconClock} strokeWidth={2.5} aria-hidden="true"/>
                        <time dateTime={product.postedDate}>{postedTimeAgo}</time>
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
                                    aria-label={t('prevImageAriaLabel')}> {/* Translated */}
                                {/* SVG */}
                            </button>
                            <button onClick={nextImage}
                                    className={`${styles.imageNavButton} ${styles.imageNavButtonNext}`}
                                    aria-label={t('nextImageAriaLabel')}> {/* Translated */}
                                {/* SVG */}
                            </button>
                        </>
                    )}
                    {/* Interaction Column */}
                    <div className={styles.interactionColumn}>
                        <button onClick={handleFavoriteToggle} className={styles.interactionButton}
                                aria-pressed={favorite}
                                aria-label={favorite ? t('unfavoriteAriaLabel') : t('favoriteAriaLabel')}> {/* Translated */}
                            <Heart size={18} className={`${styles.iconHeart} ${favorite ? styles.active : ''}`}
                                   strokeWidth={2} aria-hidden="true"/>
                            {/* Favorite count is often private, remove span if not shown */}
                            {/* <span>{product.stats?.savedCount || 0}</span> */}
                        </button>
                        <button onClick={handleLikeToggle} className={styles.interactionButton} aria-pressed={liked}
                                aria-label={liked ? t('unlikeAriaLabel') : t('likeAriaLabel')}> {/* Translated */}
                            <ThumbsUp size={18} className={`${styles.iconThumbsUp} ${liked ? styles.active : ''}`}
                                      strokeWidth={2} aria-hidden="true"/>
                            <span>{likesCount}</span>
                        </button>
                        <button onClick={handleDislikeToggle} className={styles.interactionButton}
                                aria-pressed={disliked}
                                aria-label={disliked ? t('removeDislikeAriaLabel') : t('dislikeAriaLabel')}> {/* Translated */}
                            <ThumbsDown size={18}
                                        className={`${styles.iconThumbsDown} ${disliked ? styles.active : ''}`}
                                        strokeWidth={2} aria-hidden="true"/>
                            <span>{dislikesCount}</span>
                        </button>
                        <button className={styles.interactionButton}
                                aria-label={t('viewCommentsAriaLabel', {count: commentsCount})}> {/* Translated */}
                            <MessagesSquare size={18} className={styles.iconComment} strokeWidth={2}
                                            aria-hidden="true"/>
                            <span>{commentsCount}</span>
                        </button>
                        <button className={styles.interactionButton}
                                aria-label={t('sendMessageAriaLabel', {count: messagesCount})}> {/* Translated */}
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
                        <div className={styles.viewsBadge}
                             title={t('viewsTitle', {count: product.stats?.views || 0})}> {/* Translated title */}
                            <Eye size={12} className={styles.iconEye} aria-hidden="true"/>
                            {product.stats?.views || 0}
                        </div>
                    </div>
                    <div className={styles.priceRow}>
                        <div className={styles.price}>
                            <div>{formattedPrice}</div>
                            {/* Use formatted price */}
                            {product.negotiable && (
                                <div className={styles.negotiableBadge}>{t('negotiableBadge')}</div> // Translated
                            )}
                        </div>
                        {/* Consider making this a link if it goes to product detail page */}
                        <button className={styles.viewDetailsButton}>
                            <ExternalLink size={14} className={styles.iconExternalLink} strokeWidth={2}
                                          aria-hidden="true"/>
                            {t('viewDetailsButtonText')} {/* Translated */}
                        </button>
                    </div>
                </div>
                <div className={styles.metaInfo}>
                    <div className={styles.metaBadge}>
                        <MapPin size={12} className={styles.iconMapPin} strokeWidth={2} aria-hidden="true"/>
                        <span>{product.location}</span>
                    </div>
                    {/* Combine category/subcategory */}
                    <div className={styles.metaBadge}>
                        <Tag size={12} className={styles.iconTag} strokeWidth={2} aria-hidden="true"/>
                        <span>{product.category?.name}{product.subcategory?.name ? ` • ${product.subcategory.name}` : ''}</span>
                    </div>
                    {/* Optionally show product ID */}
                    {/* <div className={styles.metaBadge}>
                         <Clock size={12} className={styles.iconClock} strokeWidth={2} aria-hidden="true"/>
                         <span>ID: #{product.id}</span>
                     </div> */}
                </div>
                {product.tags && product.tags.length > 0 && (
                    <div className={styles.tagsSection}>
                        {product.tags.map((tag, index) => (
                            <span key={index} className={styles.tag}>{tag}</span>
                        ))}
                    </div>
                )}
                <div className={styles.descriptionSection}>
                    <div className={styles.descriptionHeader}>
                        <div className={styles.descriptionTitle}>
                            <span>{t('descriptionLabel')}</span> {/* Translated */}
                            <span
                                className={styles.conditionLabel}>{conditionText}</span> {/* Use translated condition */}
                        </div>
                        <button onClick={() => setExpanded(!expanded)} className={styles.expandButton}
                                aria-expanded={expanded}
                                aria-label={expanded ? t('showLessAriaLabel') : t('showMoreAriaLabel')}> {/* Translated */}
                            {expanded ? <ChevronUp size={16} strokeWidth={2.5} aria-hidden="true"/> :
                                <ChevronDown size={16} strokeWidth={2.5} aria-hidden="true"/>}
                        </button>
                    </div>
                    <div className={`${styles.descriptionText} ${expanded ? '' : styles.collapsed}`}>
                        {product.description}
                    </div>
                    {/* Expanded Details Section (Example - make dynamic based on product data) */}
                    {expanded && (
                        <div className={styles.descriptionDetails}>
                            {product.brand && <DetailRow label={t('detailLabelBrand')} value={product.brand}/>}
                            {product.frameSize &&
                                <DetailRow label={t('detailLabelFrameSize')} value={product.frameSize}/>}
                            {product.wheelSize &&
                                <DetailRow label={t('detailLabelWheelSize')} value={product.wheelSize}/>}
                            {product.year && <DetailRow label={t('detailLabelYear')} value={product.year}/>}
                            {/* Add more details as needed */}
                        </div>
                    )}
                </div>
                {/* Seller Info */}
                <div className={styles.sellerSection}>
                    <div className={styles.sellerInfo}>
                        {/* Avatar logic */}
                        <div className={styles.sellerAvatar}>
                            <div className={styles.avatarCircle}>
                                {product.seller?.username?.charAt(0).toUpperCase() || '?'}
                            </div>
                            {product.seller?.verified && (
                                <div className={styles.verifiedBadge}
                                     title={t('verifiedBadgeTooltip')}> {/* Translated */}
                                    <Check size={12} className={styles.iconCheck} aria-hidden="true"/>
                                </div>
                            )}
                        </div>
                        <div className={styles.sellerDetails}>
                            <div className={styles.sellerName}>
                                {product.seller?.username || t('defaultUsername')}
                                {product.seller?.verified && (
                                    <span className={styles.verifiedLabel}>
                                         <ShieldCheck size={12} className={styles.iconShield} aria-hidden="true"/>
                                        {t('verifiedLabel')} {/* Translated */}
                                     </span>
                                )}
                            </div>
                            {/* Rating */}
                            <div className={styles.sellerRating} title={t('sellerRatingTitle', {
                                rating: product.seller?.rating || 0,
                                count: product.seller?.dealCount || 0
                            })}> {/* Translated title */}
                                <div className={styles.ratingStars}
                                     aria-label={t('sellerRatingLabel', {rating: product.seller?.rating || 0})}> {/* Translated aria-label */}
                                    {[...Array(5)].map((_, i) => (
                                        <svg key={i}
                                             className={`${styles.star} ${i < Math.floor(product.seller?.rating || 0) ? styles.starFilled : styles.starEmpty}`}
                                             fill="currentColor" viewBox="0 0 20 20">
                                            <path d="M9.049..."></path>
                                        </svg>
                                    ))}
                                </div>
                                {/* Use pluralization for deals suffix */}
                                <span
                                    className={styles.ratingText}>{product.seller?.rating || 0} ({t('dealsSuffix', {count: product.seller?.dealCount || 0})})</span>
                            </div>
                        </div>
                    </div>
                    <button className={styles.contactButton}>
                        <MessageCircle size={14} className={styles.iconMessageCircle} strokeWidth={2}
                                       aria-hidden="true"/>
                        {t('contactSellerButtonText')} {/* Translated */}
                    </button>
                </div>
                {/* Bottom Action Buttons */}
                <div className={styles.actionButtons}>
                    <button className={styles.actionButton} onClick={handleSaveToggle} aria-pressed={saved}>
                        <Bookmark size={16} className={`${styles.bookmarkIcon} ${saved ? styles.active : ''}`}
                                  strokeWidth={2} aria-hidden="true"/>
                        {/* Use translated Save/Saved text */}
                        <span>{saved ? t('savedButtonText') : t('saveButtonText')}</span>
                    </button>
                    <button className={styles.actionButton}>
                        <Share2 size={16} strokeWidth={2} aria-hidden="true"/>
                        <span>{t('shareButtonText')}</span> {/* Translated */}
                    </button>
                    <button className={styles.actionButton}>
                        <Flag size={16} strokeWidth={2} aria-hidden="true"/>
                        <span>{t('reportButtonText')}</span> {/* Translated */}
                    </button>
                </div>
            </div>
        </div>
    );
});
ProductCard.displayName = 'ProductCard';
const DetailRow = memo(({label, value}) => (
    <div className={styles.detailRow}>
        <span className={styles.detailLabel}>{label}:</span>
        <span className={styles.detailValue}>{value}</span>
    </div>
));
DetailRow.displayName = 'DetailRow';
// Keep preview if needed for testing/Storybook
const ProductCardPreview = () => {
    // Pass locale for testing: const locale = 'de';
    return (
        <div className={styles.previewContainer}>
            <ProductCard/> {/* Pass locale={locale} if testing */}
        </div>
    );
};
// Export the functional component, not the preview by default
export default ProductCard;
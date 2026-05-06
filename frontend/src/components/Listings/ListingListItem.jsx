import React, { memo } from 'react';
import { Heart, MapPin, Calendar, Eye, User } from '@/icons';
import styles from './ListingListItem.module.css';
/**
 * ListingListItem - Atomic Design Component
 * Displays a listing in list format for list view
 * 
 * @param {Object} props - Component props
 * @param {Object} props.listing - Listing data
 * @param {string} props.locale - Current locale
 * @returns {JSX.Element} Rendered listing list item
 */
const ListingListItem = memo(({ listing, locale = 'en' }) => {
    const {
        id,
        title,
        description,
        price,
        images,
        location,
        type,
        createdAt,
        views,
        isFavorite,
        seller
    } = listing;
    const formatPrice = (price) => {
        if (!price) return 'Price on request';
        return new Intl.NumberFormat(locale, {
            style: 'currency',
            currency: 'EUR'
        }).format(price);
    };
    const formatDate = (date) => {
        if (!date) return '';
        return new Intl.DateTimeFormat(locale, {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        }).format(new Date(date));
    };
    const handleItemClick = () => {
        // Navigate to listing detail page
        window.location.href = `/${locale}/${type}/${id}`;
    };
    const handleFavoriteClick = (e) => {
        e.stopPropagation();
        // Handle favorite toggle
        // Debug log removed for production
    };
    return (
        <article className={styles.listItem} onClick={handleItemClick}>
            {/* Image Section */}
            <div className={styles.imageWrapper}>
                {images && images.length > 0 ? (
                    <img
                        src={images[0].url}
                        alt={title}
                        className={styles.image}
                        loading="lazy"
                    />
                ) : (
                    <div className={styles.imagePlaceholder}>
                        <span>No Image</span>
                    </div>
                )}
                {/* Type Badge */}
                <div className={styles.typeBadge}>
                    {type}
                </div>
            </div>
            {/* Content Section */}
            <div className={styles.content}>
                {/* Header */}
                <div className={styles.header}>
                    <h3 className={styles.title}>{title}</h3>
                    <div className={styles.price}>
                        {formatPrice(price)}
                    </div>
                </div>
                {/* Description */}
                {description && (
                    <p className={styles.description}>
                        {description.length > 200 
                            ? `${description.substring(0, 200)}...` 
                            : description
                        }
                    </p>
                )}
                {/* Meta Information */}
                <div className={styles.meta}>
                    {/* Location */}
                    {location && (
                        <div className={styles.metaItem}>
                            <MapPin size={14} />
                            <span>{location}</span>
                        </div>
                    )}
                    {/* Date */}
                    {createdAt && (
                        <div className={styles.metaItem}>
                            <Calendar size={14} />
                            <span>{formatDate(createdAt)}</span>
                        </div>
                    )}
                    {/* Views */}
                    {views && (
                        <div className={styles.metaItem}>
                            <Eye size={14} />
                            <span>{views} views</span>
                        </div>
                    )}
                    {/* Seller */}
                    {seller && (
                        <div className={styles.metaItem}>
                            <User size={14} />
                            <span>{seller.name || 'Anonymous'}</span>
                        </div>
                    )}
                </div>
            </div>
            {/* Actions Section */}
            <div className={styles.actions}>
                <button
                    className={`${styles.favoriteButton} ${isFavorite ? styles.favorited : ''}`}
                    onClick={handleFavoriteClick}
                    aria-label={isFavorite ? 'Remove from favorites' : 'Add to favorites'}
                >
                    <Heart size={18} />
                </button>
            </div>
        </article>
    );
});
ListingListItem.displayName = 'ListingListItem';
export default ListingListItem; 
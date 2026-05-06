import React, { memo } from 'react';
import { Heart, MapPin, Calendar, Eye } from '@/icons';
import styles from './ListingCard.module.css';
/**
 * ListingCard - Atomic Design Component
 * Displays a listing in card format for grid view
 * 
 * @param {Object} props - Component props
 * @param {Object} props.listing - Listing data
 * @param {string} props.locale - Current locale
 * @returns {JSX.Element} Rendered listing card
 */
const ListingCard = memo(({ listing, locale = 'en' }) => {
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
        isFavorite
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
            day: 'numeric'
        }).format(new Date(date));
    };
    const handleCardClick = () => {
        // Navigate to listing detail page
        window.location.href = `/${locale}/${type}/${id}`;
    };
    const handleFavoriteClick = (e) => {
        e.stopPropagation();
        // Handle favorite toggle
        // Debug log removed for production
    };
    return (
        <article className={styles.card} onClick={handleCardClick}>
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
                {/* Favorite Button */}
                <button
                    className={`${styles.favoriteButton} ${isFavorite ? styles.favorited : ''}`}
                    onClick={handleFavoriteClick}
                    aria-label={isFavorite ? 'Remove from favorites' : 'Add to favorites'}
                >
                    <Heart size={16} />
                </button>
                {/* Type Badge */}
                <div className={styles.typeBadge}>
                    {type}
                </div>
            </div>
            {/* Content Section */}
            <div className={styles.content}>
                {/* Title */}
                <h3 className={styles.title}>{title}</h3>
                {/* Description */}
                {description && (
                    <p className={styles.description}>
                        {description.length > 100 
                            ? `${description.substring(0, 100)}...` 
                            : description
                        }
                    </p>
                )}
                {/* Location */}
                {location && (
                    <div className={styles.location}>
                        <MapPin size={14} />
                        <span>{location}</span>
                    </div>
                )}
                {/* Footer */}
                <div className={styles.footer}>
                    {/* Price */}
                    <div className={styles.price}>
                        {formatPrice(price)}
                    </div>
                    {/* Meta Info */}
                    <div className={styles.meta}>
                        {views && (
                            <div className={styles.metaItem}>
                                <Eye size={12} />
                                <span>{views}</span>
                            </div>
                        )}
                        {createdAt && (
                            <div className={styles.metaItem}>
                                <Calendar size={12} />
                                <span>{formatDate(createdAt)}</span>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </article>
    );
});
ListingCard.displayName = 'ListingCard';
export default ListingCard; 
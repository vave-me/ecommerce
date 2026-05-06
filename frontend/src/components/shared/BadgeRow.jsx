import React, { memo } from 'react';
import { Clock, Eye, Tag, ShoppingBag, Flame, MapPin } from 'lucide-react';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import sharedStyles from './CardShared.module.css';
// Extend dayjs with relative time plugin
dayjs.extend(relativeTime);
/**
 * BadgeRow - Professional, elegant badge row component for all card types
 * Designed with mobile-first responsive approach and sophisticated styling
 * 
 * @param {Object} props - Component props
 * @param {string} props.type - Card type: 'deals', 'classified', 'services', 'jobs', 'properties', 'vehicles'
 * @param {string} props.category - Category or subcategory
 * @param {string} props.merchantName - Merchant, company, or entity name
 * @param {number} props.views - View count
 * @param {string|Date} props.createdAt - Creation date
 * @param {boolean} props.isHot - Whether to show hot/popular badge
 * @param {string} props.className - Additional CSS classes
 */
const BadgeRow = memo(({
    type = 'classified',
    category = '',
    merchantName = '',
    views = 0,
    createdAt = null,
    isHot = false,
    className = '',
    showMerchant = true
}) => {
    // Calculate relative time with better formatting
    const timeAgo = React.useMemo(() => {
        if (!createdAt) return null;
        try {
            const now = dayjs();
            const created = dayjs(createdAt);
            const diffHours = now.diff(created, 'hour');
            const diffDays = now.diff(created, 'day');
            if (diffHours < 1) return 'Just now';
            if (diffHours < 24) return `${diffHours}h`;
            if (diffDays < 7) return `${diffDays}d`;
            return created.format('MMM DD');
        } catch (error) {
            return null;
        }
    }, [createdAt]);
    // Get type configuration with modern labels and icons
    const getTypeConfig = (cardType) => {
        switch (cardType) {
            case 'deals': 
                return { label: 'Deal', color: 'deals', priority: 1 };
            case 'classified': 
            case 'marketplace': 
                return { label: 'new', color: 'classified', priority: 2 };
            case 'services':
                return { label: 'Service', color: 'services', priority: 3 };
            case 'jobs':
                return { label: 'Job', color: 'jobs', priority: 4 };
            case 'properties':
                return { label: 'Property', color: 'properties', priority: 5 };
            case 'vehicles':
                return { label: 'Vehicle', color: 'vehicles', priority: 6 };
            default:
                return { label: 'Item', color: 'classified', priority: 7 };
        }
    };
    const typeConfig = getTypeConfig(type);
    // Format views count with K notation
    const formatViews = (count) => {
        if (count < 1000) return count.toString();
        if (count < 10000) return `${(count / 1000).toFixed(1)}k`;
        return `${Math.floor(count / 1000)}k`;
    };
    // Determine which badges to show based on priority and screen space
    const shouldShowMerchant = merchantName && merchantName.length > 0;
    const shouldShowCategory = category && category.length > 0;
    const shouldShowTime = timeAgo && timeAgo.length > 0;
    return (
        <div className={`${sharedStyles.badgeRow} ${className}`}>
            {/* Primary badges container - always visible */}
            <div className={sharedStyles.badgeRowPrimary}>
                {/* Type Badge - highest priority */}
                <span className={`${sharedStyles.typeBadge} ${sharedStyles[typeConfig.color]}`}>
                    {typeConfig.label}
            </span>
                {/* Hot Badge - high priority when applicable */}
            {isHot && (
                    <span className={sharedStyles.hotBadge} aria-label="Popular item">
                        <Flame size={8} className={sharedStyles.hotIcon} aria-hidden="true" />
                        <span className={sharedStyles.hotText}>Hot</span>
                </span>
            )}
                {/* Views Badge - always show if > 0 */}
                {views > 0 && (
                    <span className={sharedStyles.viewsBadge} aria-label={`${views} views`}>
                        <Eye size={8} className={sharedStyles.badgeIcon} aria-hidden="true" />
                        <span>{formatViews(views)}</span>
                    </span>
                )}
            </div>
            {/* Secondary badges container - responsive visibility */}
            <div className={sharedStyles.badgeRowSecondary}>
                {/* Category Badge - medium priority */}
                {shouldShowCategory && (
                    <span className={sharedStyles.categoryBadge}>
                        <Tag size={8} className={sharedStyles.badgeIcon} aria-hidden="true" />
                    <span>{category}</span>
                </span>
            )}
                {/* Merchant Badge - lower priority */}
                {shouldShowMerchant && showMerchant && (
                    <span className={sharedStyles.merchantBadge}>
                        <ShoppingBag size={8} className={sharedStyles.badgeIcon} aria-hidden="true" />
                    <span>{merchantName}</span>
                </span>
            )}
                {/* Time Badge - lowest priority, desktop only */}
                {shouldShowTime && (
                    <span className={sharedStyles.timeBadge}>
                        <Clock size={8} className={sharedStyles.badgeIcon} aria-hidden="true" />
                <span>{timeAgo}</span>
            </span>
                )}
            </div>
        </div>
    );
});
BadgeRow.displayName = 'BadgeRow';
export default BadgeRow; 
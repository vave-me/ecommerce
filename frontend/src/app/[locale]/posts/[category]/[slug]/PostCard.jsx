"use client"
import React, {useState, useEffect, useMemo, useCallback} from 'react';
import {useTranslations} from 'next-intl'; // Import useTranslations
import {
    Heart,
    MessageCircle,
    Share2,
    Bookmark,
    ThumbsDown,
    Eye,
    Flag,
    MoreHorizontal,
    Check,
    Clock,
    MapPin
} from '@/icons';
import dayjs from 'dayjs'; // For date formatting
import relativeTime from 'dayjs/plugin/relativeTime'; // For relative time
import 'dayjs/locale/en'; // English locale
import 'dayjs/locale/pl'; // Polish locale  
import 'dayjs/locale/de'; // German locale
// Enable relative time plugin
dayjs.extend(relativeTime);
import styles from './PostCard.module.css';
// Map locales for dayjs  
const dateLocales = {en: 'en', pl: 'pl', de: 'de'};
// The actual PostCard component accepting props
const PostCard = ({post: postProp, locale = 'en'}) => {
    const t = useTranslations('PostCard'); // Use PostCard namespace
    // Process post data with smart defaults and feature detection
    const post = useMemo(() => {
        if (!postProp) return null;
        // Smart heat score calculation based on engagement metrics
        const views = parseInt(postProp.views || postProp.stats?.views) || 0;
        const likes = parseInt(postProp.likes || postProp.stats?.likes) || 0;
        const comments = parseInt(postProp.comments || postProp.stats?.comments) || 0;
        const bookmarks = parseInt(postProp.bookmarks || postProp.stats?.bookmarks) || 0;
        const shares = parseInt(postProp.shares || postProp.stats?.shares) || 0;
        const dislikes = parseInt(postProp.dislikes || postProp.stats?.dislikes) || 0;
        // Weighted heat score for posts: balanced engagement metrics
        const viewWeight = views * 0.1;
        const likeWeight = likes * 3;
        const commentWeight = comments * 4;
        const bookmarkWeight = bookmarks * 5;
        const shareWeight = shares * 6; // Shares are valuable but not as much as for news
        const dislikeWeight = dislikes * -2; // Dislikes reduce heat score
        const calculatedHeatScore = Math.min(100, Math.max(0, Math.round(
            viewWeight + likeWeight + commentWeight + bookmarkWeight + shareWeight + dislikeWeight
        )));
        // Extract post features from data
        const extractedFeatures = {
            isVerifiedAuthor: postProp.author?.verified || postProp.authorVerified,
            hasImage: postProp.image || postProp.thumbnail || postProp.coverImage,
            hasLocation: postProp.location || postProp.author?.location,
            isRecent: postProp.publishedDate ? 
                     dayjs().diff(dayjs(postProp.publishedDate), 'hours') <= 24 : false,
            hasHashtags: Array.isArray(postProp.hashtags) && postProp.hashtags.length > 0,
            hasUserInteractions: postProp.userInteractions && typeof postProp.userInteractions === 'object',
        };
        // Smart author information extraction
        const authorInfo = {
            name: postProp.author?.name || postProp.authorName || t('defaultAuthor'),
            username: postProp.author?.username || postProp.authorUsername,
            avatar: postProp.author?.avatar || postProp.authorAvatar || '/images/default-avatar.png',
            verified: extractedFeatures.isVerifiedAuthor,
        };
        // Smart hashtag extraction
        const extractedHashtags = [];
        if (postProp.hashtags && Array.isArray(postProp.hashtags)) {
            extractedHashtags.push(...postProp.hashtags);
        } else if (postProp.tags && Array.isArray(postProp.tags)) {
            extractedHashtags.push(...postProp.tags);
        } else if (typeof postProp.hashtags === 'string') {
            extractedHashtags.push(...postProp.hashtags.split(',').map(tag => tag.trim()));
        }
        return {
            ...postProp,
            views,
            likes,
            comments,
            bookmarks,
            shares,
            dislikes,
            heatScore: calculatedHeatScore,
            features: extractedFeatures,
            author: authorInfo,
            hashtags: extractedHashtags,
            // Ensure required fields exist
            title: postProp.title || postProp.name || t('unnamedPost'),
            content: postProp.content || postProp.description || postProp.summary || t('noContentAvailable'),
            image: extractedFeatures.hasImage,
        };
    }, [postProp, t]);
    // ---- State Management ----
    const [liked, setLiked] = useState(post?.userInteractions?.liked || post?.isLiked || false);
    const [disliked, setDisliked] = useState(post?.userInteractions?.disliked || post?.isDisliked || false);
    const [saved, setSaved] = useState(post?.userInteractions?.saved || post?.isSaved || false);
    const [flagged, setFlagged] = useState(post?.userInteractions?.flagged || post?.isFlagged || false);
    // Use state for counts to allow optimistic updates
    const [likesCount, setLikesCount] = useState(post?.likes || 0);
    const [dislikesCount, setDislikesCount] = useState(post?.dislikes || 0);
    const [savesCount, setSavesCount] = useState(post?.bookmarks || 0);
    const [sharesCount, setSharesCount] = useState(post?.shares || 0);
    const [commentsCount, setCommentsCount] = useState(post?.comments || 0);
    const [viewsCount, setViewsCount] = useState(post?.views || 0);
    // ---- Locale-aware Formatting ----
    // Relative Time Formatting
    const [timeAgo, setTimeAgo] = useState('');
    useEffect(() => {
        if (!post?.publishedDate && !post?.createdAt) return;
        const postDate = new Date(post.publishedDate || post.createdAt);
        try {
            dayjs.locale(dateLocales[locale] || 'en');
            setTimeAgo(dayjs(postDate).fromNow());
        } catch (e) {
            setTimeAgo(t('invalidDate'));
        }
    }, [post?.publishedDate, post?.createdAt, locale, t]);
    // Number Formatting (Compact Notation)
    const formatNumber = useCallback((num) => {
        if (typeof num !== 'number') return '0';
        try {
            return new Intl.NumberFormat(locale, {
                notation: "compact",
                maximumFractionDigits: 1
            }).format(num);
        } catch (e) {
            return num.toString();
        }
    }, [locale]);
    // ---- Interaction Handlers ----
    const handleLike = useCallback(() => {
        const newLiked = !liked;
        setLiked(newLiked);
        setLikesCount(newLiked ? likesCount + 1 : Math.max(0, likesCount - 1));
        if (newLiked && disliked) {
            setDisliked(false);
            setDislikesCount(Math.max(0, dislikesCount - 1));
        }
        // TODO: API call: updateLikeStatus(post.id, newLiked);
    }, [liked, likesCount, disliked, dislikesCount]);
    const handleDislike = useCallback(() => {
        const newDisliked = !disliked;
        setDisliked(newDisliked);
        setDislikesCount(newDisliked ? dislikesCount + 1 : Math.max(0, dislikesCount - 1));
        if (newDisliked && liked) {
            setLiked(false);
            setLikesCount(Math.max(0, likesCount - 1));
        }
        // TODO: API call: updateDislikeStatus(post.id, newDisliked);
    }, [disliked, dislikesCount, liked, likesCount]);
    const handleSave = useCallback(() => {
        const newSaved = !saved;
        setSaved(newSaved);
        setSavesCount(newSaved ? savesCount + 1 : Math.max(0, savesCount - 1));
        // TODO: API call: updateSaveStatus(post.id, newSaved);
    }, [saved, savesCount]);
    const handleShare = useCallback(() => {
        // TODO: Implement actual share functionality
        setSharesCount(sharesCount + 1);
        // TODO: API call: recordShare(post.id);
    }, [sharesCount, t]);
    const handleReport = useCallback(() => {
        const newFlagged = !flagged;
        setFlagged(newFlagged);
        // TODO: API call: reportPost(post.id, reason);
    }, [flagged, t]);
    // ---- Render Logic ----
    if (!post?.id) {
        return <div className={styles.cardContainer} role="alert">{t('errorLoadingPost')}</div>;
    }
    return (
        <article className={styles.cardContainer} aria-labelledby={`post-title-${post.id}`}>
            {/* Post header */}
            <div className={styles.postHeader}>
                <div className={styles.avatarContainer}>
                    <img
                        src={post.author?.avatar}
                        alt={t('authorAvatarAlt', {name: post.author?.name})}
                        className={styles.avatar}
                    />
                    {post.author?.verified && (
                        <div className={styles.verifiedBadge} title={t('verifiedBadgeTooltip')}>
                            <Check size={12} className={styles.verifiedIcon}/>
                        </div>
                    )}
                </div>
                <div className={styles.authorInfo}>
                    <div className={styles.authorNameRow}>
                        <h2 id={`post-title-${post.id}`}
                            className={styles.authorName}>{post.title}</h2>
                        {post.author?.username && <span className={styles.authorUsername}>{post.author.username}</span>}
                    </div>
                    <div className={styles.postMetaInfo}>
                        <Clock size={12} className={styles.metaIcon} aria-hidden="true"/>
                        <time dateTime={post.publishedDate || post.createdAt}>{timeAgo}</time>
                        {post.location && (
                            <>
                                <span className={styles.metaSeparator}>•</span>
                                <MapPin size={12} className={styles.metaIcon} aria-hidden="true"/>
                                <span>{post.location}</span>
                            </>
                        )}
                    </div>
                </div>
                {/* Views counter with Tooltip */}
                <div className={styles.viewsCounter} title={t('viewsCounterTooltip', {count: viewsCount})}>
                    <Eye size={12} className={styles.viewsIcon} aria-hidden="true"/>
                    <span aria-label={t('viewsCounterTooltip', {count: viewsCount})}>
                       {formatNumber(viewsCount)}
                    </span>
                </div>
                {/* Post options Button with Aria Label */}
                <button className={styles.optionsButton} aria-label={t('optionsButtonAriaLabel')}>
                    <MoreHorizontal size={18}/>
                </button>
            </div>
            {/* Post content */}
            <div className={styles.postContent}>
                <p className={styles.postText}>{post.content}</p>
                <div className={styles.hashtagsContainer}>
                    {post.hashtags?.slice(0, 5).map((tag, index) => (
                        <span key={index} className={styles.hashtag}>
                            #{tag}
                        </span>
                    ))}
                </div>
            </div>
            {/* Post image/media */}
            {post.image && (
                <div className={styles.mediaContainer}>
                    <img
                        src={post.image}
                        alt={t('postImageAlt', {title: post.title})}
                        className={styles.postImage}
                    />
                    {/* Interaction Column */}
                    <div className={styles.interactionColumn}>
                        {/* Like button */}
                        <button
                            className={`${styles.interactionButton} ${liked ? styles.likedButton : ''}`}
                            onClick={handleLike}
                            aria-pressed={liked}
                            aria-label={t('likeButtonAriaLabel', {count: likesCount})}
                            title={t('likeButtonAriaLabel', {count: likesCount})}
                        >
                            <Heart size={18} className={liked ? styles.filledIcon : ''} aria-hidden="true"/>
                            <span className={styles.countLabel}>{formatNumber(likesCount)}</span>
                        </button>
                        {/* Dislike button */}
                        <button
                            className={`${styles.interactionButton} ${disliked ? styles.dislikedButton : ''}`}
                            onClick={handleDislike}
                            aria-pressed={disliked}
                            aria-label={t('dislikeButtonAriaLabel', {count: dislikesCount})}
                            title={t('dislikeButtonAriaLabel', {count: dislikesCount})}
                        >
                            <ThumbsDown size={18} aria-hidden="true"/>
                            <span className={styles.countLabel}>{formatNumber(dislikesCount)}</span>
                        </button>
                        {/* Comment button */}
                        <button
                            className={styles.interactionButton}
                            aria-label={t('commentButtonAriaLabel', {count: commentsCount})}
                            title={t('commentButtonAriaLabel', {count: commentsCount})}
                        >
                            <MessageCircle size={18} aria-hidden="true"/>
                            <span className={styles.countLabel}>{formatNumber(commentsCount)}</span>
                        </button>
                        {/* Save button */}
                        <button
                            className={`${styles.interactionButton} ${saved ? styles.savedButton : ''}`}
                            onClick={handleSave}
                            aria-pressed={saved}
                            aria-label={saved ? t('unsaveButtonAriaLabel', {count: savesCount}) : t('saveButtonAriaLabel', {count: savesCount})}
                            title={saved ? t('unsaveButtonAriaLabel', {count: savesCount}) : t('saveButtonAriaLabel', {count: savesCount})}
                        >
                            <Bookmark size={18} className={saved ? styles.filledIcon : ''} aria-hidden="true"/>
                            <span className={styles.countLabel}>{formatNumber(savesCount)}</span>
                        </button>
                        {/* Share button */}
                        <button
                            className={styles.interactionButton}
                            onClick={handleShare}
                            aria-label={t('shareButtonAriaLabel', {count: sharesCount})}
                            title={t('shareButtonAriaLabel', {count: sharesCount})}
                        >
                            <Share2 size={18} aria-hidden="true"/>
                            <span className={styles.countLabel}>{formatNumber(sharesCount)}</span>
                        </button>
                        {/* Report button */}
                        <button
                            className={`${styles.reportButton} ${flagged ? styles.reportedButton : ''}`}
                            onClick={handleReport}
                            aria-pressed={flagged}
                            aria-label={flagged ? t('unreportButtonAriaLabel') : t('reportButtonAriaLabel')}
                            title={flagged ? t('unreportButtonAriaLabel') : t('reportButtonAriaLabel')}
                        >
                            <Flag size={18} aria-hidden="true"/>
                        </button>
                    </div>
                </div>
            )}
        </article>
    );
};
export default PostCard;
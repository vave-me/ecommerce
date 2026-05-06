"use client" // Ensure this is at the top
import React, {useState, useEffect} from 'react';
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
    // Process real API data with smart defaults and feature detection
    const actualPost = postProp?.post || postProp;
    const postMetrics = actualPost?.metrics || {};
    // Convert string metrics to numbers
    const convertMetricToNumber = (value) => parseInt(value || '0', 10);
    const likesCount = convertMetricToNumber(postMetrics.likesCount);
    const dislikesCount = convertMetricToNumber(postMetrics.dislikesCount);
    const commentsCount = convertMetricToNumber(postMetrics.commentsCount);
    const visitedCount = convertMetricToNumber(postMetrics.visitedCount);
    const bookmarksCount = convertMetricToNumber(postMetrics.addedToWishlistCount);
    const sharesCount = convertMetricToNumber(postMetrics.sharedCount);
    // Smart feature detection from description and tags
    const description = actualPost?.description || actualPost?.content || '';
    const tags = actualPost?.tags || [];
    const allText = `${description} ${tags.join(' ')}`.toLowerCase();
    // Extract hashtags from content
    const hashtagRegex = /#[\w]+/g;
    const extractedHashtags = description.match(hashtagRegex) || [];
    const smartHashtags = [...new Set([...tags, ...extractedHashtags])].slice(0, 5);
    // Handle author data
    const authorData = actualPost?.author || actualPost?.user || {};
    const authorName = authorData?.name || actualPost?.authorName || actualPost?.userName || 'Anonymous';
    const authorUsername = authorData?.username || actualPost?.authorUsername || actualPost?.userUsername || '';
    const authorAvatar = authorData?.avatar || actualPost?.authorAvatar || actualPost?.userAvatar || '/images/default-avatar.png';
    const isVerified = authorData?.verified || actualPost?.authorVerified || actualPost?.userVerified || false;
    // Handle timestamps
    const createdAt = actualPost?.createdAt || actualPost?.publishedDate || new Date().toISOString();
    const publishedDate = actualPost?.publishedDate || actualPost?.createdAt || new Date().toISOString();
    // Handle location
    const location = actualPost?.location || actualPost?.address || '';
    // Handle media
    const postImage = actualPost?.image || actualPost?.thumbnail || actualPost?.coverImage || '';
    // Build processed post data with real API integration
    const post = {
        id: actualPost?.id || 'unknown',
        author: {
            name: authorName,
            username: authorUsername ? `@${authorUsername.replace('@', '')}` : '',
            avatar: authorAvatar,
            verified: isVerified
        },
        title: actualPost?.title || actualPost?.name || '',
        content: description,
        image: postImage,
        createdAt: createdAt,
        publishedDate: publishedDate,
        location: location,
        stats: {
            likes: likesCount,
            comments: commentsCount,
            views: visitedCount,
            dislikes: dislikesCount,
            bookmarks: bookmarksCount,
            shares: sharesCount
        },
        hashtags: smartHashtags,
        slug: actualPost?.slug || actualPost?.id || '',
        category: actualPost?.category || {slug: 'general'}
    };
    const t = useTranslations('PostCard'); // Use PostCard namespace
    // ---- State Management ----
    // Initialize state based on post data, assuming user interaction state comes from API/context later
    const [liked, setLiked] = useState(false); // TODO: Replace with actual user like status
    const [disliked, setDisliked] = useState(false); // TODO: Replace with actual user dislike status
    const [saved, setSaved] = useState(false); // TODO: Replace with actual user save status
    const [flagged, setFlagged] = useState(false); // TODO: Replace with actual user flag status
    // Use state for counts to allow optimistic updates
    const [likesCountState, setLikesCount] = useState(post.stats?.likes || 0);
    const [dislikesCountState, setDislikesCount] = useState(post.stats?.dislikes || 0);
    const [savesCount, setSavesCount] = useState(post.stats?.bookmarks || 0);
    const [sharesCountState, setSharesCount] = useState(post.stats?.shares || 0);
    const [commentsCountState, setCommentsCount] = useState(post.stats?.comments || 0);
    const [viewsCount, setViewsCount] = useState(post.stats?.views || 0);
    // ---- Locale-aware Formatting ----
    // Relative Time Formatting
    const [timeAgo, setTimeAgo] = useState('');
    useEffect(() => {
        const postDate = new Date(post.publishedDate || post.createdAt || Date.now());
        try {
            // Set dayjs locale and format relative time
            dayjs.locale(dateLocales[locale] || 'en');
            setTimeAgo(dayjs(postDate).fromNow());
        } catch (e) {
            setTimeAgo(t('invalidDate')); // Fallback translation
        }
    }, [post.publishedDate, post.createdAt, locale, t]);
    // Number Formatting (Compact Notation)
    const formatNumber = (num) => {
        if (typeof num !== 'number') return '0'; // Handle non-numeric input
        try {
            // Use Intl.NumberFormat for locale-aware compact notation (K, M, etc.)
            return new Intl.NumberFormat(locale, {
                notation: "compact",
                maximumFractionDigits: 1
            }).format(num);
        } catch (e) {
            return num.toString(); // Fallback to simple string
        }
    };
    // ---- Interaction Handlers (Placeholder Logic) ----
    // These should ideally call an API to persist the action
    const handleLike = () => {
        const newLiked = !liked;
        setLiked(newLiked);
        setLikesCount(newLiked ? likesCountState + 1 : likesCountState - 1);
        if (newLiked && disliked) { // If liking while disliked
            setDisliked(false);
            setDislikesCount(dislikesCountState > 0 ? dislikesCountState - 1 : 0);
        }
        // API call: updateLikeStatus(post.id, newLiked);
    };
    const handleDislike = () => {
        const newDisliked = !disliked;
        setDisliked(newDisliked);
        setDislikesCount(newDisliked ? dislikesCountState + 1 : dislikesCountState - 1);
        if (newDisliked && liked) { // If disliking while liked
            setLiked(false);
            setLikesCount(likesCountState > 0 ? likesCountState - 1 : 0);
        }
        // API call: updateDislikeStatus(post.id, newDisliked);
    };
    const handleSave = () => {
        const newSaved = !saved;
        setSaved(newSaved);
        setSavesCount(newSaved ? savesCount + 1 : savesCount - 1);
        // API call: updateSaveStatus(post.id, newSaved);
    };
    const handleShare = () => {
        // Trigger platform share API or internal share logic
        // Optimistically update count if needed, but sharing often doesn't change visible count immediately
        // setSharesCount(sharesCountState + 1);
        // API call: recordShare(post.id);
        alert(t('shareActionPlaceholder')); // Placeholder feedback
    };
    const handleReport = () => {
        const newFlagged = !flagged;
        setFlagged(newFlagged);
        // API call: reportPost(post.id, reason);
        alert(newFlagged ? t('reportActionPlaceholder') : t('unreportActionPlaceholder')); // Placeholder feedback
    };
    // ---- Render ----
    return (
        <article className={styles.cardContainer} aria-labelledby={`post-title-${post.id}`}>
            {/* Post header */}
            <div className={styles.postHeader}>
                <div className={styles.avatarContainer}>
                    <img
                        src={post.author?.avatar || '/images/default-avatar.png'} // Fallback avatar
                        alt={t('authorAvatarAlt', {name: post.author?.name || t('defaultAuthor')})}
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
                        {/* Use post.title or name if available and relevant */}
                        <h2 id={`post-title-${post.id}`}
                            className={styles.authorName}>{post.title || post.author?.name || t('unnamedPost')}</h2>
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
                    {/* Add Dropdown Menu Here */}
                </button>
            </div>
            {/* Post content */}
            <div className={styles.postContent}>
                {post.content && <p className={styles.postText}>{post.content}</p>}
                {/* Hashtags */}
                {post.hashtags && post.hashtags.length > 0 && (
                    <div className={styles.hashtagsContainer}>
                        {post.hashtags.map((hashtag, index) => (
                            <span key={index} className={styles.hashtag}>
                                {hashtag.startsWith('#') ? hashtag : `#${hashtag}`}
                            </span>
                        ))}
                    </div>
                )}
                {/* Post image */}
                {post.image && (
                    <div className={styles.imageContainer}>
                        <img
                            src={post.image}
                            alt={t('postImageAlt', {title: post.title || 'Post image'})}
                            className={styles.postImage}
                        />
                    </div>
                )}
            </div>
            {/* Interaction buttons */}
            <div className={styles.interactionBar}>
                <button
                    className={`${styles.interactionButton} ${liked ? styles.liked : ''}`}
                    onClick={handleLike}
                    aria-label={t('likeButtonLabel', {count: likesCountState})}
                >
                    <Heart size={16} className={styles.interactionIcon}/>
                    <span>{formatNumber(likesCountState)}</span>
                </button>
                <button
                    className={`${styles.interactionButton} ${disliked ? styles.disliked : ''}`}
                    onClick={handleDislike}
                    aria-label={t('dislikeButtonLabel', {count: dislikesCountState})}
                >
                    <ThumbsDown size={16} className={styles.interactionIcon}/>
                    <span>{formatNumber(dislikesCountState)}</span>
                </button>
                <button
                    className={styles.interactionButton}
                    aria-label={t('commentButtonLabel', {count: commentsCountState})}
                >
                    <MessageCircle size={16} className={styles.interactionIcon}/>
                    <span>{formatNumber(commentsCountState)}</span>
                </button>
                <button
                    className={styles.interactionButton}
                    onClick={handleShare}
                    aria-label={t('shareButtonLabel')}
                >
                    <Share2 size={16} className={styles.interactionIcon}/>
                    <span>{formatNumber(sharesCountState)}</span>
                </button>
                <button
                    className={`${styles.interactionButton} ${saved ? styles.saved : ''}`}
                    onClick={handleSave}
                    aria-label={t('saveButtonLabel')}
                >
                    <Bookmark size={16} className={styles.interactionIcon}/>
                    <span>{formatNumber(savesCount)}</span>
                </button>
                <button
                    className={styles.interactionButton}
                    onClick={handleReport}
                    aria-label={t('reportButtonLabel')}
                >
                    <Flag size={16} className={styles.interactionIcon}/>
                </button>
            </div>
        </article>
    );
};
// PostCardPreview component for displaying multiple posts
const PostCardPreview = ({posts = []}) => {
    if (!posts || posts.length === 0) {
        return (
            <div className={styles.previewContainer}>
                <div className={styles.emptyState}>
                    <MessageCircle size={48} className={styles.emptyIcon}/>
                    <h3>No posts available</h3>
                    <p>Check back later for new posts and updates.</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.previewContainer}>
            <div className={styles.previewGrid}>
                {posts.map((post, index) => (
                    <PostCard key={post?.id || index} post={post}/>
                ))}
            </div>
        </div>
    );
};
export default PostCard;
export {PostCardPreview};
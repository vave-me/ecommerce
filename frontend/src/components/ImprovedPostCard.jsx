"use client";
import React, {useCallback, useState, useEffect, useMemo, memo} from 'react';
import {
    Bookmark,
    Check,
    ChevronLeft,
    ChevronRight,
    Clock,
    Eye,
    Flag,
    Globe,
    Link2,
    MapPin,
    MoreHorizontal,
    UserCheck,
    UserPlus,
} from '@/icons';
import styles from '../app/[locale]/design/page.module.css'; // Following DealCard pattern
import {Engagement} from '../app/[locale]/design/Engagement'; // Following DealCard pattern
import {toast} from "react-toastify";
import {useAuth} from "../context/AuthContext"; // Fix path following DealCard pattern
import {useDispatch} from "react-redux";
import useActivityApi from "../hooks/useActivityApi"; // Fix path following DealCard pattern
import {openMessageModal} from "../redux/slices/modalsSlice"; // Fix path following DealCard pattern
import CommentsSetup from "../features/Comments/CommentsSetup"; // Fix path following DealCard pattern
import parse from 'html-react-parser';
import {getBaseUserById} from "../api/client/userApi";
import Image from "next/image";
/**
 * Quick stat component for listing details
 * purely presentational - PERFORMANCE: Memoized
 */
const QuickStat = memo(({icon, text}) => (
    <div className={styles.quickStat}>
        <span className={styles.quickStatIcon}>{icon}</span>
        <span className={styles.quickStatText}>{text}</span>
    </div>
));
QuickStat.displayName = 'QuickStat';
/**
 * Toast notification wrapper - PERFORMANCE: Memoized with useCallback
 */
const useToastNotification = () => {
    return useCallback((type, message) => {
        const options = {theme: "colored"};
        switch (type) {
            case "success":
                toast.success(message, options);
                break;
            case "info":
                toast.info(message, options);
                break;
            case "error":
                toast.error(message, options);
                break;
            case "warn":
                toast.warn(message, options);
                break;
            default:
                toast(message, options);
        }
    }, []);
};
/**
 * Utility function to calculate time ago from timestamp
 * PERFORMANCE: Pure utility function
 */
const getTimeAgo = (timestamp) => {
    if (!timestamp) return 'Unknown time';
    const now = new Date();
    const past = new Date(timestamp);
    const diffInMs = now - past;
    const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
    const diffInHours = Math.floor(diffInMinutes / 60);
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInMinutes < 1) return 'Just now';
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInHours < 24) return `${diffInHours}h ago`;
    if (diffInDays < 7) return `${diffInDays}d ago`;
    return past.toLocaleDateString();
};
/**
 * Utility function to estimate read time based on content length
 * PERFORMANCE: Pure utility function
 */
const calculateReadTime = (content) => {
    if (!content) return '1 min read';
    // Strip HTML tags and count words
    const plainText = content.replace(/<[^>]*>/g, '');
    const words = plainText.trim().split(/\s+/).length;
    const wordsPerMinute = 200; // Average reading speed
    const minutes = Math.ceil(words / wordsPerMinute);
    return `${minutes} min read`;
};
/**
 * Actual Post Card - PERFORMANCE OPTIMIZED: Memoized component with stable computations
 */
const ImprovedPostCard = memo(({post}) => {
    const showToast = useToastNotification();
    const {user} = useAuth();
    const reduxDispatch = useDispatch();
    const {handleLike, handleDislike} = useActivityApi();
    const [author, setAuthor] = useState(null);
    const [isLoadingAuthor, setIsLoadingAuthor] = useState(true);
    const [authorError, setAuthorError] = useState(null);
    // Extract the actual post data from the response structure
    const actualPost = post?.post || post;
    const postMetrics = actualPost?.metrics || {};
    // PERFORMANCE: Memoized computed values from real data
    const computedData = useMemo(() => {
        return {
            // Convert string metrics to numbers
            likesCount: parseInt(postMetrics.likesCount || '0', 10),
            dislikesCount: parseInt(postMetrics.dislikesCount || '0', 10),
            commentsCount: parseInt(postMetrics.commentsCount || '0', 10),
            sharedCount: parseInt(postMetrics.sharedCount || '0', 10),
            visitedCount: parseInt(postMetrics.visitedCount || '0', 10),
            // Calculate derived values
            timeAgo: getTimeAgo(post?.createdAt || post?.updatedAt),
            readTime: calculateReadTime(actualPost?.description),
            // Handle location
            hasLocation: actualPost?.lat && actualPost?.lng && (actualPost.lat !== 0 || actualPost.lng !== 0),
            location: actualPost?.lat && actualPost?.lng && (actualPost.lat !== 0 || actualPost.lng !== 0)
                ? `${actualPost.lat}, ${actualPost.lng}`
                : null,
            // Handle images - use thumbnail if available
            images: actualPost?.thumbnail ? [actualPost.thumbnail] : [],
            // Status and visibility
            status: actualPost?.status || 'PUBLISHED',
            visibility: 'Public', // Default since not provided in API
        };
    }, [actualPost, postMetrics, post?.createdAt, post?.updatedAt]);
    // PERFORMANCE: Memoized fetch author function to prevent recreation
    const fetchAuthor = useCallback(async () => {
        if (!actualPost?.userId) {
            setIsLoadingAuthor(false);
            return;
        }
        try {
            setIsLoadingAuthor(true);
            setAuthorError(null);
            const response = await getBaseUserById(actualPost.userId);
            setAuthor(response?.user || response);
        } catch (error) {
            setAuthorError(error);
        } finally {
            setIsLoadingAuthor(false);
        }
    }, [actualPost?.userId]);
    // Fetch author data when post changes
    useEffect(() => {
        fetchAuthor();
    }, [fetchAuthor]);
    // PERFORMANCE: Memoized post data merge to prevent object recreation
    const postData = useMemo(() => ({
        // Real data from your back-end
        id: actualPost?.id,
        userId: actualPost?.userId,
        name: actualPost?.name || "Untitled Post",
        description: actualPost?.description || "No description",
        tags: actualPost?.tags || [],
        status: computedData.status,
        thumbnail: actualPost?.thumbnail,
        lat: actualPost?.lat || 0,
        lng: actualPost?.lng || 0,
        postType: actualPost?.postType,
        categoryId: actualPost?.categoryId,
        categorySlug: actualPost?.categorySlug,
        // Author data from API
        author: author ? {
            id: author.id || author.userId,
            name: author.firstName || author.userName || 'Anonymous',
            username: author.userName ? `@${author.userName}` : '@anonymous',
            avatar: author.thumbnail || "/images/placeholder-avatar.jpg",
            verified: author.verified || false,
            following: author.following || false,
            bio: author.bio || '',
            followers: author.followers || 0
        } : {
            // Fallback author data if API fails
            id: actualPost?.userId,
            name: "Anonymous",
            username: "@anonymous",
            avatar: "/images/placeholder-avatar.jpg",
            verified: false,
            following: false,
            bio: '',
            followers: 0
        },
        // Computed values from real data
        readTime: computedData.readTime,
        location: computedData.location,
        visibility: computedData.visibility,
        timeAgo: computedData.timeAgo,
        images: computedData.images,
        // Real metrics from API
        metrics: {
            likes: computedData.likesCount,
            dislikes: computedData.dislikesCount,
            comments: computedData.commentsCount,
            shares: computedData.sharedCount,
            views: computedData.visitedCount,
            reposts: computedData.sharedCount, // Using sharedCount as reposts
            addedToWishlist: parseInt(postMetrics.addedToWishlistCount || '0', 10),
            addedToBasket: parseInt(postMetrics.addedToBasketCount || '0', 10),
            reported: parseInt(postMetrics.reportedCount || '0', 10),
            reviews: parseInt(postMetrics.reviewsCount || '0', 10),
            rating: parseFloat(postMetrics.rating || '0'),
            ratingCount: parseInt(postMetrics.ratingCount || '0', 10),
        },
        // Extract hashtags from description or tags
        hashtags: actualPost?.tags || [],
        mentions: [], // Would need to be extracted from description if needed
        // State flags - these would come from user's interaction data in a real app
        isLiked: false, // Integrate with activity API when available
        isBookmarked: false, // Integrate with bookmarks API when available
        isReposted: false, // Integrate with reposts API when available
        isFavorite: false, // Integrate with favorites API when available
        isPinned: false, // Integrate with metadata API when available
        isPromoted: false, // Integrate with promotion system when available
        // Timestamps
        createdAt: post?.createdAt,
        updatedAt: post?.updatedAt,
    }), [actualPost, postMetrics, computedData, author]);
    // State hooks using real data
    const [favorite, setFavorite] = useState(postData.isFavorite);
    const [liked, setLiked] = useState(postData.isLiked);
    const [disliked, setDisliked] = useState(false);
    const [bookmarked, setBookmarked] = useState(postData.isBookmarked);
    const [reposted, setReposted] = useState(postData.isReposted);
    const [showActions, setShowActions] = useState(false);
    const [showImageIndex, setShowImageIndex] = useState(0);
    const [isFollowing, setIsFollowing] = useState(postData.author.following);
    const [showComments, setShowComments] = useState(false);
    // Using real metrics for counts
    const [likesCount, setLikesCount] = useState(postData.metrics.likes);
    const [repostsCount, setRepostsCount] = useState(postData.metrics.reposts);
    const [sharesCount, setSharesCount] = useState(postData.metrics.shares);
    const [shared, setShared] = useState(false);
    const userId = user?.userId;
    const handleRepost = () => {
        if (reposted) {
            setRepostsCount(repostsCount - 1);
        } else {
            setRepostsCount(repostsCount + 1);
        }
        setReposted(!reposted);
    };
    const handleShare = () => {
        if (!shared) {
            setSharesCount(sharesCount + 1);
            setShared(true);
            setTimeout(() => setShared(false), 2000);
        }
    };
    const handleCommentClick = useCallback(() => {
        setShowComments(prev => !prev);
    }, []);
    const handleFavorite = useCallback(() => {
        setFavorite(prev => !prev);
    }, []);
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            showToast("warn", "Please log in to like posts.");
            return;
        }
        if (!postData.id) {
            return;
        }
        handleLike(postData.id, userId).catch(() => {
            // Handle error
        });
    }, [postData.id, userId, handleLike, showToast]);
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            showToast("warn", "Please log in to dislike posts.");
            return;
        }
        if (!postData.id) {
            return;
        }
        handleDislike(postData.id, userId).catch(() => {
            setDisliked(prev => !prev);
        });
    }, [postData.id, userId, handleDislike, showToast]);
    const handleOpenMessage = useCallback(() => {
        reduxDispatch(openMessageModal({
            itemId: postData.id,
            recipientId: postData.author.id,
        }));
    }, [postData, reduxDispatch]);
    // Image navigation
    const nextImage = () => {
        setShowImageIndex((prev) =>
            prev === postData.images.length - 1 ? 0 : prev + 1
        );
    };
    const prevImage = () => {
        setShowImageIndex((prev) =>
            prev === 0 ? postData.images.length - 1 : prev - 1
        );
    };
    // Format numbers for display
    const formatNumber = (num) => {
        if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
        return num.toString();
    };
    /**
     * Format text to highlight # and @ while also parsing HTML
     */
    const renderFormattedContent = (text) => {
        if (!text) return '';
        // First parse the HTML content
        const parsedHtml = parse(text);
        // If the parsed content is a string (no HTML), process hashtags and mentions
        if (typeof parsedHtml === 'string') {
            const processedText = parsedHtml.split(/(\s+)/).map((word, index) => {
                if (word.startsWith('#')) {
                    return `<span class="${styles.hashtag}">${word}</span>`;
                } else if (word.startsWith('@')) {
                    return `<span class="${styles.mention}">${word}</span>`;
                }
                return word;
            }).join('');
            return parse(processedText);
        }
        // If it's already parsed HTML, return it as is
        return parsedHtml;
    };
    return (
        <div className={styles.articleContainer}>
            {/* Title Card */}
            <div className={styles.card}>
                <div className={styles.cardContent}>
                    {/* Author header */}
                    <div className={styles.authorHeader}>
                        <div className={styles.authorSection}>
                            {isLoadingAuthor ? (
                                <div className={styles.authorLoading}>
                                    <div className={styles.authorAvatarSkeleton}/>
                                    <div className={styles.authorInfoSkeleton}>
                                        <div className={styles.authorNameSkeleton}/>
                                        <div className={styles.authorUsernameSkeleton}/>
                                    </div>
                                </div>
                            ) : authorError ? (
                                <div className={styles.authorError}>
                                    <div className={styles.authorAvatarError}/>
                                    <div className={styles.authorInfoError}>
                                        <span>Error loading author</span>
                                    </div>
                                </div>
                            ) : author ? (
                                <div className={styles.authorInfo}>
                                    <img
                                        src={postData.author.avatar}
                                        alt={postData.author.name}
                                        className={styles.authorAvatar}
                                    />
                                    <div className={styles.authorDetails}>
                                        <span className={styles.authorName}>{postData.author.name}</span>
                                        <span className={styles.authorUsername}>{postData.author.username}</span>
                                    </div>
                                </div>
                            ) : (
                                <div className={styles.authorInfo}>
                                    <img
                                        src="/images/placeholder-avatar.jpg"
                                        alt="Anonymous"
                                        className={styles.authorAvatar}
                                    />
                                    <div className={styles.authorDetails}>
                                        <span className={styles.authorName}>Anonymous</span>
                                        <span className={styles.authorUsername}>@anonymous</span>
                                    </div>
                                </div>
                            )}
                        </div>
                        {/* Post controls */}
                        <div className={styles.postControls}>
                            {!isFollowing && (
                                <button
                                    className={styles.followButton}
                                    onClick={() => setIsFollowing(true)}
                                >
                                    Follow
                                </button>
                            )}
                            <div className={styles.actionsDropdown}>
                                <button
                                    className={styles.actionButton}
                                    onClick={() => setShowActions(!showActions)}
                                >
                                    <MoreHorizontal size={16}/>
                                </button>
                                {showActions && (
                                    <div className={styles.dropdownMenu}>
                                        <button
                                            className={styles.dropdownItem}
                                            onClick={() => setIsFollowing(!isFollowing)}
                                        >
                                            {isFollowing ? (
                                                <>
                                                    <UserCheck size={14} className={styles.followingIcon}/>
                                                    <span>
                                                        Unfollow {postData.author.username.replace('@', '')}
                                                    </span>
                                                </>
                                            ) : (
                                                <>
                                                    <UserPlus size={14} className={styles.dropdownItemIcon}/>
                                                    <span>
                                                        Follow {postData.author.username.replace('@', '')}
                                                    </span>
                                                </>
                                            )}
                                        </button>
                                        <button
                                            className={styles.dropdownItem}
                                            onClick={() => setBookmarked(!bookmarked)}
                                        >
                                            <Bookmark size={14} className={styles.dropdownItemIcon}/>
                                            <span>
                                                {bookmarked ? 'Remove from bookmarks' : 'Save to bookmarks'}
                                            </span>
                                        </button>
                                        <button className={styles.dropdownItem}>
                                            <Link2 size={14} className={styles.dropdownItemIcon}/>
                                            <span>Copy link to post</span>
                                        </button>
                                        <button className={styles.dropdownItemDanger}>
                                            <Flag size={14} className={styles.dropdownItemDangerIcon}/>
                                            <span>Report post</span>
                                        </button>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                    {/* Title Card */}
                    <div className={styles.titleCard}>
                        <div className={styles.titleCardContent}>
                            <h1 className={styles.articleTitle}>
                                {postData.name}
                            </h1>
                            <div className={styles.articleMeta}>
                                <div className={styles.readTime}>
                                    <Clock size={14} className={styles.readTimeIcon}/>
                                    <span>{postData.readTime}</span>
                                </div>
                                <div className={styles.timeAgo}>
                                    <span>{postData.timeAgo}</span>
                                </div>
                                {postData.metrics.views > 0 && (
                                    <div className={styles.articleViews}>
                                        <Eye size={14} className={styles.viewsIcon}/>
                                        <span>
                                            {formatNumber(postData.metrics.views)} views
                                        </span>
                                    </div>
                                )}
                                {postData.location && (
                                    <div className={styles.locationBadge}>
                                        <MapPin size={14}/>
                                        <span>{postData.location}</span>
                                    </div>
                                )}
                                <div className={styles.categoryBadge}>
                                    {postData.status || "PUBLISHED"}
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* Post content */}
                    <div className={styles.postContent}>
                        <div className={styles.postText}>
                            {renderFormattedContent(postData.description)}
                        </div>
                    </div>
                    {/* Post Image(s) */}
                    {postData.images && postData.images.length > 0 && (
                        <div className={styles.mediaContainer}>
                            {postData.images.length > 1 && (
                                <div className={styles.imageCounter}>
                                    {showImageIndex + 1}/{postData.images.length}
                                </div>
                            )}
                            <div className={styles.imageWrapper}>
                                <Image
                                    src={postData.images[showImageIndex]}
                                    alt="Post attachment"
                                    className={styles.postImage}
                                    width={800}
                                    height={450}
                                    sizes="(max-width: 768px) 100vw, (max-width: 1400px) 50vw, 800px"
                                    style={{objectFit: 'cover'}}
                                    priority={showImageIndex === 0}
                                />
                            </div>
                            {postData.images.length > 1 && (
                                <>
                                    <button
                                        onClick={prevImage}
                                        className={`${styles.imageNavButton} ${styles.prevButton}`}
                                        aria-label="Previous image"
                                    >
                                        <ChevronLeft size={16}/>
                                    </button>
                                    <button
                                        onClick={nextImage}
                                        className={`${styles.imageNavButton} ${styles.nextButton}`}
                                        aria-label="Next image"
                                    >
                                        <ChevronRight size={16}/>
                                    </button>
                                </>
                            )}
                        </div>
                    )}
                    {/* Hashtags */}
                    {postData.hashtags.length > 0 && (
                        <div className={styles.hashtagsContainer}>
                            {postData.hashtags.map((tag, i) => (
                                <span
                                    key={i}
                                    className={styles.hashtagLabel}
                                >
                                    #{tag}
                                </span>
                            ))}
                        </div>
                    )}
                    {/* Real Metrics Display */}
                    <div className={styles.metricsContainer}>
                        <div className={styles.metricsRow}>
                            <QuickStat
                                icon={<Eye size={16}/>}
                                text={`${formatNumber(postData.metrics.views)} views`}
                            />
                            {postData.metrics.comments > 0 && (
                                <QuickStat
                                    icon="💬"
                                    text={`${formatNumber(postData.metrics.comments)} comments`}
                                />
                            )}
                            {postData.metrics.shares > 0 && (
                                <QuickStat
                                    icon="🔄"
                                    text={`${formatNumber(postData.metrics.shares)} shares`}
                                />
                            )}
                        </div>
                    </div>
                    {/* Engagement */}
                    <Engagement
                        liked={liked}
                        disliked={disliked}
                        favorite={favorite}
                        onLike={handleLikeClick}
                        onMessageClick={handleOpenMessage}
                        onCommentClick={handleCommentClick}
                        onDislike={handleDislikeClick}
                        onFavorite={handleFavorite}
                        likesCount={likesCount}
                        sharesCount={sharesCount}
                        commentsCount={postData.metrics.comments}
                    />
                </div>
                {showComments && (
                    <div className={styles.commentsWrapper}>
                        <CommentsSetup
                            userId={userId}
                            itemId={postData.id}
                            itemType="post"
                            toggleCommentsList={handleCommentClick}
                            categoryId={actualPost?.categoryId}
                            postName={postData.name}
                            postThumbnail={postData.images[0]}
                        />
                    </div>
                )}
            </div>
        </div>
    );
});
export default ImprovedPostCard; 
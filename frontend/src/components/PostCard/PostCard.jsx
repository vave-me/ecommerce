"use client";

import React, { memo, useState, useCallback, useMemo, useEffect } from 'react';
import Link from 'next/link';
import { 
    Heart, MessageCircle, Share2, ThumbsUp, ThumbsDown, 
    Flame, AlertCircle, PlayCircle 
} from 'lucide-react';
import styles from './PostCard.module.css';
import { useAuth } from '../../context/AuthContext';
import { useDispatch } from 'react-redux';
import { openMessageModal } from '../../redux/slices/modalsSlice';
import { toast } from 'react-toastify';
import useActivityApi from '../../hooks/useActivityApi';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import { getMediaByItem } from '../../api/mediaApi';
import VideoPlayer from '../shared/VideoPlayer';

const PostCard = memo(function PostCard({ post, isVisible, isMobile }) {
    const { user } = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const { handleLike, handleDislike } = useActivityApi();
    
    // Extract post data
    const actualPost = post?.post || post;
    const postMetrics = actualPost?.metrics || {};
    
    // Helper function to detect media type from URL
    const getMediaTypeFromUrl = (url) => {
        if (!url) return 'image';
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv', '.m4v'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext)) ? 'video' : 'image';
    };
    
    // Process post data with defaults
    const postWithDefaults = useMemo(() => {
        const processed = {
            id: actualPost?.id || '',
            name: actualPost?.name || actualPost?.title || 'Untitled Post',
            description: actualPost?.description || 'No description available.',
            thumbnail: actualPost?.thumbnail || '',
            categorySlug: actualPost?.categorySlug || '',
            added: actualPost?.added || actualPost?.createdAt || new Date().toISOString(),
            userId: actualPost?.userId,
            categoryId: actualPost?.categoryId,
            tags: actualPost?.tags || [],
            media: actualPost?.media || [],
            images: actualPost?.images || [],
            status: actualPost?.status || 'PUBLISHED'
        };
        
        // Debug log removed for production
        
        return processed;
    }, [actualPost]);
    
    // Interaction state
    const [interactionState, setInteractionState] = useState({
        liked: false,
        disliked: false,
        favorite: false,
        showComments: false,
        currentMediaIndex: 0,
        playingVideos: {} // Track which videos are playing
    });
    
    // Media state
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoadingMedia, setIsLoadingMedia] = useState(false);
    
    // Fetch media using the same pattern as ClassifiedCard
    useEffect(() => {
        const loadMedia = async () => {
            if (!postWithDefaults.id) {
                return;
            }
            setIsLoadingMedia(true);
            try {
                const mediaResponse = await getMediaByItem(postWithDefaults.id);
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || getMediaTypeFromUrl(item.url),
                        alt: item.altText || postWithDefaults.name || 'Post image'
                    }));
                    setMediaItems(formattedMedia);
                } else {
                    // Check post object for any image data
                    const postImages = [];
                    
                    // Check for images array (common in posts)
                    if (actualPost?.images?.length > 0) {
                        actualPost.images.forEach(img => {
                            if (typeof img === 'string') {
                                postImages.push({ url: img, type: getMediaTypeFromUrl(img) });
                            } else if (img.url) {
                                postImages.push({ url: img.url, type: img.type || getMediaTypeFromUrl(img.url) });
                            }
                        });
                    }
                    
                    // Check for media array
                    if (postWithDefaults.media?.length > 0) {
                        postWithDefaults.media.forEach(media => {
                            if (media.url) {
                                postImages.push({
                                    ...media,
                                    type: media.type || getMediaTypeFromUrl(media.url)
                                });
                            }
                        });
                    }
                    
                    // Only use thumbnail if no other media exists
                    // and the thumbnail is not just a frame from a video
                    if (postImages.length === 0 && postWithDefaults.thumbnail) {
                        // Skip thumbnail if we have video media
                        const hasVideo = actualPost?.images?.some(img => 
                            typeof img === 'string' ? getMediaTypeFromUrl(img) === 'video' : 
                            img.type === 'video' || (img.url && getMediaTypeFromUrl(img.url) === 'video')
                        ) || postWithDefaults.media?.some(media => 
                            media.type === 'video' || (media.url && getMediaTypeFromUrl(media.url) === 'video')
                        );
                        
                        if (!hasVideo) {
                            postImages.push({ url: postWithDefaults.thumbnail, type: 'image' });
                        }
                    }
                    
                    if (postImages.length > 0) {
                        setMediaItems(postImages);
                    }
                }
            } catch (error) {
                // Silent error handling - media will remain empty
                // Silent error handling - media will remain empty
            } finally {
                setIsLoadingMedia(false);
            }
        };
        loadMedia();
    }, [postWithDefaults.id, postWithDefaults.name, postWithDefaults.thumbnail, postWithDefaults.media, actualPost?.images]);
    
    // Combine API media with post media, prioritizing API media
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
        
        // Otherwise, try post images array
        if (postWithDefaults.images && postWithDefaults.images.length > 0) {
            return postWithDefaults.images.map(ensureMediaType);
        }
        
        // Try post media array
        if (postWithDefaults.media && postWithDefaults.media.length > 0) {
            return postWithDefaults.media.map(ensureMediaType);
        }
        
        // Fall back to thumbnail ONLY if no other media exists
        if (postWithDefaults.thumbnail) {
            return [ensureMediaType(postWithDefaults.thumbnail)];
        }
        
        // No media available
        return [];
    }, [mediaItems, postWithDefaults.images, postWithDefaults.media, postWithDefaults.thumbnail]);
    
    // Media resolution tracking removed for production
    
    // Format metrics
    const metrics = useMemo(() => ({
        likes: parseInt(postMetrics.likesCount || '0', 10),
        dislikes: parseInt(postMetrics.dislikesCount || '0', 10),
        comments: parseInt(postMetrics.commentsCount || '0', 10),
        shares: parseInt(postMetrics.sharedCount || '0', 10),
        views: parseInt(postMetrics.visitedCount || '0', 10),
        saved: parseInt(postMetrics.addedToWishlistCount || '0', 10)
    }), [postMetrics]);
    
    // Calculate time ago
    const timeAgo = useMemo(() => {
        const date = new Date(postWithDefaults.added);
        const now = new Date();
        const diff = now - date;
        const minutes = Math.floor(diff / 60000);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);
        
        if (days > 0) return `${days}d ago`;
        if (hours > 0) return `${hours}h ago`;
        if (minutes > 0) return `${minutes}m ago`;
        return 'Just now';
    }, [postWithDefaults.added]);
    
    // Calculate read time
    const readTime = useMemo(() => {
        if (!postWithDefaults.description) return '1 min read';
        const words = postWithDefaults.description.replace(/<[^>]*>/g, '').trim().split(/\s+/).length;
        const minutes = Math.ceil(words / 200);
        return `${minutes} min read`;
    }, [postWithDefaults.description]);
    
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
    
    // Interaction handlers
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to like posts.");
            return;
        }
        if (!postWithDefaults.id) return;
        
        handleLike(postWithDefaults.id, userId).catch(() => {});
        setInteractionState(prev => ({
            ...prev,
            liked: true,
            disliked: false
        }));
    }, [postWithDefaults.id, userId, handleLike]);
    
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            toast.warn("Please log in to dislike posts.");
            return;
        }
        if (!postWithDefaults.id) return;
        
        handleDislike(postWithDefaults.id, userId).catch(() => {});
        setInteractionState(prev => ({
            ...prev,
            liked: false,
            disliked: true
        }));
    }, [postWithDefaults.id, userId, handleDislike]);
    
    const handleCommentClick = useCallback(() => {
        setInteractionState(prev => ({
            ...prev,
            showComments: !prev.showComments
        }));
    }, []);
    
    const handleFavorite = useCallback(() => {
        setInteractionState(prev => ({
            ...prev,
            favorite: !prev.favorite
        }));
    }, []);
    
    const handleOpenMessage = useCallback(() => {
        if (!postWithDefaults.id) return;
        dispatch(openMessageModal({
            itemId: postWithDefaults.id,
            recipientId: postWithDefaults.userId || 'author'
        }));
    }, [postWithDefaults.id, postWithDefaults.userId, dispatch]);
    
    // TikTok video hooks - must be declared at top level
    const videoRef = React.useRef(null);
    const [isMuted, setIsMuted] = React.useState(true);
    const [isPlaying, setIsPlaying] = React.useState(false);
    const [videoError, setVideoError] = React.useState(false);
    
    // Check if this is a single video post (TikTok style)
    const isSingleVideoPost = resolvedMedia.length === 1 && resolvedMedia[0].type === 'video';
    const videoMedia = isSingleVideoPost ? resolvedMedia[0] : null;
    
    // TikTok video callbacks
    const toggleMute = useCallback((e) => {
        e.stopPropagation();
        setIsMuted(prev => !prev);
    }, []);
    
    const handleVideoClick = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        if (isMobile) {
            toggleMute(e);
        }
    }, [isMobile, toggleMute]);
    
    const handleMouseEnter = useCallback(() => {
        if (!isMobile && videoRef.current && isSingleVideoPost) {
            videoRef.current.muted = true;
            videoRef.current.play().catch(e => {
                // Video autoplay failed silently
                setIsPlaying(false);
            });
        }
    }, [isMobile, isSingleVideoPost]);
    
    const handleMouseLeave = useCallback(() => {
        if (!isMobile && videoRef.current && isSingleVideoPost) {
            videoRef.current.pause();
            videoRef.current.currentTime = 0;
            setIsPlaying(false);
        }
    }, [isMobile, isSingleVideoPost]);
    
    // Render TikTok-style video post
    if (isSingleVideoPost && videoMedia) {
        return (
            <div className={styles.tikTokContainer}>
                <div 
                    className={styles.videoWrapper}
                    onMouseEnter={handleMouseEnter}
                    onMouseLeave={handleMouseLeave}
                >
                    <video
                        ref={videoRef}
                        src={videoMedia.url}
                        poster={videoMedia.thumbnail || ""}
                        loop
                        muted={isMuted}
                        playsInline
                        className={styles.videoElement}
                        preload="metadata"
                        onClick={handleVideoClick}
                        onError={() => setVideoError(true)}
                        onPlay={() => setIsPlaying(true)}
                        onPause={() => setIsPlaying(false)}
                    />
                    
                    {/* Video Error State */}
                    {videoError && (
                        <div className={styles.errorOverlay}>
                            <p>Unable to load video</p>
                        </div>
                    )}
                    
                    {/* Desktop Play Button */}
                    {!isMobile && !isPlaying && (
                        <button 
                            className={styles.playButton}
                            onClick={(e) => {
                                e.stopPropagation();
                                if (videoRef.current) {
                                    videoRef.current.play().catch(err => {
                                        // Play failed silently
                                    });
                                }
                            }}
                            aria-label="Play video"
                        >
                            <div className={styles.playIconWrapper}>
                                ▶
                            </div>
                        </button>
                    )}
                    
                    {/* Mute Button Overlay */}
                    <div className={styles.muteButtonOverlay}>
                        <button
                            type="button"
                            className={styles.muteButton}
                            onClick={toggleMute}
                            aria-label={isMuted ? "Unmute" : "Mute"}
                        >
                            {isMuted ? "🔇" : "🔊"}
                        </button>
                    </div>
                    
                    {/* Title and Description Overlay - TikTok Style */}
                    <div className={styles.tikTokContentOverlay}>
                        <div className={styles.tikTokUserInfo}>
                            <span className={styles.tikTokUsername}>@{postWithDefaults.username || "user"}</span>
                        </div>
                        <h3 className={styles.tikTokTitle}>
                            {postWithDefaults.name}
                        </h3>
                        {postWithDefaults.description && (
                            <p className={styles.tikTokDescription}>
                                {postWithDefaults.description}
                            </p>
                        )}
                    </div>
                </div>
                
                {/* Engagement Bar - Below Video */}
                <div className={styles.engagementBarBottom}>
                    <button 
                        className={`${styles.actionButtonBottom} ${interactionState.liked ? styles.active : ''}`}
                        onClick={handleLikeClick}
                        aria-label={`Like post (${metrics.likes} likes)`}
                    >
                        <ThumbsUp size={20} />
                        {metrics.likes > 0 && <span className={styles.actionCountBottom}>{metrics.likes}</span>}
                    </button>
                    
                    <button 
                        className={`${styles.actionButtonBottom} ${interactionState.disliked ? styles.active : ''}`}
                        onClick={handleDislikeClick}
                        aria-label={`Dislike post (${metrics.dislikes} dislikes)`}
                    >
                        <ThumbsDown size={20} />
                        {metrics.dislikes > 0 && <span className={styles.actionCountBottom}>{metrics.dislikes}</span>}
                    </button>
                    
                    <button 
                        className={styles.actionButtonBottom}
                        onClick={handleCommentClick}
                        aria-label={`Comments (${metrics.comments} comments)`}
                    >
                        <MessageCircle size={20} />
                        {metrics.comments > 0 && <span className={styles.actionCountBottom}>{metrics.comments}</span>}
                    </button>
                    
                    <button 
                        className={styles.actionButtonBottom}
                        onClick={handleOpenMessage}
                        aria-label="Send message"
                    >
                        <MessageCircle size={20} fill="currentColor" />
                    </button>
                    
                    <button 
                        className={styles.actionButtonBottom}
                        onClick={() => {
                            if (navigator.share && typeof window !== 'undefined') {
                                navigator.share({
                                    title: postWithDefaults.name,
                                    text: postWithDefaults.description,
                                    url: `${window.location.origin}/post/${postWithDefaults.id}`
                                });
                            }
                        }}
                        aria-label="Share post"
                    >
                        <Share2 size={20} />
                    </button>
                    
                    <button 
                        className={`${styles.actionButtonBottom} ${interactionState.favorite ? styles.active : ''}`}
                        onClick={handleFavorite}
                        aria-label={`Save to wishlist (${metrics.saved} saves)`}
                    >
                        <Heart size={20} fill={interactionState.favorite ? "currentColor" : "none"} />
                        {metrics.saved > 0 && <span className={styles.actionCountBottom}>{metrics.saved}</span>}
                    </button>
                </div>
            </div>
        );
    }
    
    // Regular post card rendering (for images or multiple media)
    return (
        <article className={`${styles.card} ${isLoadingMedia ? styles.isLoading : ''}`} aria-labelledby={`post-title-${postWithDefaults.id}`}>
            {/* Live region for dynamic updates */}
            <div className="sr-only" role="status" aria-live="polite" aria-atomic="true">
                {interactionState.liked && `You liked ${postWithDefaults.name}`}
                {interactionState.disliked && `You disliked ${postWithDefaults.name}`}
                {interactionState.favorite && `Added ${postWithDefaults.name} to favorites`}
            </div>
            {/* Media Container - Primary Visual Hierarchy */}
            <div className={styles.mediaContainer}>
                {/* Status Indicators */}
                <div className={styles.statusIndicators}>
                    {metrics.views > 1000 && (
                        <div className={`${styles.badge} ${styles.featuredBadge}`} role="status">
                            <span className="sr-only">Post status:</span>
                            <Flame size={12} />
                            <span>Trending</span>
                        </div>
                    )}
                    {timeAgo === 'Just now' && (
                        <div className={`${styles.badge} ${styles.newBadge}`} role="status">
                            <span className="sr-only">Post status:</span>
                            <AlertCircle size={12} />
                            <span>New</span>
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
                                            productName={postWithDefaults.name}
                                        />
                                    ) : (
                                        <img src={media.url} alt={`${postWithDefaults.name} - Image ${index + 1} of ${resolvedMedia.length}`} loading="lazy" />
                                    )}
                                </div>
                            ))}
                        </div>

                        {/* Media Navigation */}
                        {resolvedMedia.length > 1 && (
                            <>
                                <div className={styles.mediaNav}>
                                    <button 
                                        className={styles.mediaNavButton} 
                                        onClick={() => handleMediaNavigation('prev')}
                                        aria-label={`Previous image (${interactionState.currentMediaIndex} of ${resolvedMedia.length})`}
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                            <path d="M15 18l-6-6 6-6" />
                                        </svg>
                                    </button>
                                    <button 
                                        className={styles.mediaNavButton} 
                                        onClick={() => handleMediaNavigation('next')}
                                        aria-label={`Next image (${interactionState.currentMediaIndex + 2} of ${resolvedMedia.length})`}
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                            <path d="M9 6l6 6-6 6" />
                                        </svg>
                                    </button>
                                </div>
                                
                                {/* Media Indicators */}
                                <div className={styles.mediaIndicators} role="tablist" aria-label="Post images">
                                    {resolvedMedia.map((_, index) => (
                                        <button
                                            key={index}
                                            className={`${styles.indicator} ${index === interactionState.currentMediaIndex ? styles.active : ''}`}
                                            onClick={() => setInteractionState(prev => ({ ...prev, currentMediaIndex: index }))}
                                            role="tab"
                                            aria-selected={index === interactionState.currentMediaIndex}
                                            aria-label={`Go to image ${index + 1} of ${resolvedMedia.length}`}
                                            tabIndex={index === interactionState.currentMediaIndex ? 0 : -1}
                                        />
                                    ))}
                                </div>
                            </>
                        )}
                    </>
                ) : (
                    <div className={styles.placeholderMedia}>
                        <span className={styles.placeholderText}>No media</span>
                    </div>
                )}
            </div>
            
            {/* Content Section */}
            <div className={styles.content}>
                {/* Category Label */}
                {postWithDefaults.categorySlug && (
                    <div className={styles.categoryLabel} aria-label="Post category">
                        {postWithDefaults.categorySlug}
                    </div>
                )}
                
                {/* Title - Primary */}
                <Link href={
                    !postWithDefaults.categorySlug 
                        ? `/post/${postWithDefaults.id}`
                        : `/posts/${postWithDefaults.categorySlug}/${postWithDefaults.id}`
                }>
                    <h3 id={`post-title-${postWithDefaults.id}`} className={styles.title} title={postWithDefaults.name}>
                        {postWithDefaults.name}
                    </h3>
                </Link>
                
                {/* Description - Secondary */}
                <p className={styles.description} aria-label="Post excerpt">
                    <span className="sr-only">Excerpt: </span>
                    {postWithDefaults.description.replace(/<[^>]*>/g, '').substring(0, 400)}...
                    <span className="sr-only">Read time: {readTime}</span>
                </p>
            </div>
            
            {/* Engagement Bar - Always Visible */}
            <div className={styles.engagementBar}>
                <div className={styles.engagementActions}>
                    {/* Like/Dislike */}
                    <button 
                        className={`${styles.actionButton} ${interactionState.liked ? styles.active : ''}`}
                        onClick={handleLikeClick}
                        aria-label={`Like post (${metrics.likes} likes)`}
                        aria-pressed={interactionState.liked}
                    >
                        <ThumbsUp size={18} />
                        {metrics.likes > 0 && <span className={styles.actionCount}>{metrics.likes}</span>}
                    </button>
                    <button 
                        className={`${styles.actionButton} ${interactionState.disliked ? styles.active : ''}`}
                        onClick={handleDislikeClick}
                        aria-label={`Dislike post (${metrics.dislikes} dislikes)`}
                        aria-pressed={interactionState.disliked}
                    >
                        <ThumbsDown size={18} />
                        {metrics.dislikes > 0 && <span className={styles.actionCount}>{metrics.dislikes}</span>}
                    </button>
                    
                    {/* Comment */}
                    <button 
                        className={`${styles.actionButton} ${interactionState.showComments ? styles.active : ''}`}
                        onClick={handleCommentClick}
                        aria-label={`Comments (${metrics.comments} comments)`}
                        aria-expanded={interactionState.showComments}
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
                                    title: postWithDefaults.name,
                                    text: postWithDefaults.description,
                                    url: `${window.location.origin}/posts/${postWithDefaults.categorySlug}/${postWithDefaults.id}`
                                });
                            }
                        }}
                        aria-label="Share post"
                    >
                        <Share2 size={18} />
                        {metrics.shares > 0 && <span className={styles.actionCount}>{metrics.shares}</span>}
                    </button>
                    
                    {/* Favorite */}
                    <button 
                        className={`${styles.actionButton} ${interactionState.favorite ? styles.active : ''}`}
                        onClick={handleFavorite}
                        aria-label={`Save to favorites (${metrics.saved} saves)`}
                        aria-pressed={interactionState.favorite}
                    >
                        <Heart size={18} />
                        {metrics.saved > 0 && <span className={styles.actionCount}>{metrics.saved}</span>}
                    </button>
                </div>
            </div>
            
            {/* Comments */}
            {interactionState.showComments && (
                <div className={styles.commentsWrapper}>
                    <CommentsSetup
                        userId={userId}
                        itemId={postWithDefaults.id}
                        itemType="post"
                        toggleCommentsList={handleCommentClick}
                        categoryId={postWithDefaults.categoryId}
                        postName={postWithDefaults.name}
                        postThumbnail={postWithDefaults.thumbnail}
                    />
                </div>
            )}
        </article>
    );
});

export default PostCard;
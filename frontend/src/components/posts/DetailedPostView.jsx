"use client";
import React, { useState, useEffect, useMemo, useCallback, memo } from 'react';
import { useTranslations } from 'next-intl';
import {
    Heart, MessageCircle, Share2, Bookmark, ThumbsDown, ThumbsUp,
    Eye, Flag, MoreHorizontal, Check, Clock, MapPin, Calendar,
    User, ChevronLeft, ChevronRight, Camera, Link as LinkIcon,
    Tag, TrendingUp, Users, AlertCircle, ExternalLink
} from '@/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/en';
import 'dayjs/locale/pl';
import 'dayjs/locale/de';
import { useAuth } from '../../context/AuthContext';
import { useDispatch } from 'react-redux';
import { openMessageModal } from '../../redux/slices/modalsSlice';
import useActivityApi from '../../hooks/useActivityApi';
import { getMediaByItem } from '../../api/mediaApi';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import { sanitizeRichHtml } from '../../utils/sanitizeHtml';
import styles from './DetailedPostView.module.css';

// Enable relative time plugin
dayjs.extend(relativeTime);

// Map locales for dayjs
const dateLocales = { en: 'en', pl: 'pl', de: 'de' };

/**
 * DetailedPostView - Full post detail page component
 * Shows comprehensive post information with enhanced layout
 */
const DetailedPostView = memo(({ post: postProp, locale = 'en' }) => {
    const t = useTranslations('PostCard');
    const { user } = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const { handleLike, handleDislike } = useActivityApi();
    
    // Set dayjs locale
    dayjs.locale(dateLocales[locale] || 'en');
    
    // Process post data with smart defaults
    const post = useMemo(() => {
        if (!postProp) return null;
        
        // Extract metrics
        const views = parseInt(postProp.views || postProp.stats?.views || postProp.metrics?.views) || 0;
        const likes = parseInt(postProp.likes || postProp.stats?.likes || postProp.metrics?.likes) || 0;
        const dislikes = parseInt(postProp.dislikes || postProp.stats?.dislikes || postProp.metrics?.dislikes) || 0;
        const comments = parseInt(postProp.comments || postProp.stats?.comments || postProp.metrics?.comments) || 0;
        const bookmarks = parseInt(postProp.bookmarks || postProp.stats?.bookmarks || postProp.metrics?.bookmarks) || 0;
        const shares = parseInt(postProp.shares || postProp.stats?.shares || postProp.metrics?.shares) || 0;
        
        // Calculate heat score
        const viewWeight = views * 0.1;
        const likeWeight = likes * 3;
        const commentWeight = comments * 4;
        const bookmarkWeight = bookmarks * 5;
        const shareWeight = shares * 6;
        const dislikeWeight = dislikes * -2;
        
        const calculatedHeatScore = Math.min(100, Math.max(0, Math.round(
            viewWeight + likeWeight + commentWeight + bookmarkWeight + shareWeight + dislikeWeight
        )));
        
        // Extract features
        const extractedFeatures = {
            isVerifiedAuthor: postProp.author?.verified || postProp.authorVerified,
            hasImage: postProp.image || postProp.thumbnail || postProp.coverImage,
            hasLocation: postProp.location || postProp.author?.location,
            isRecent: postProp.publishedDate ? 
                     dayjs().diff(dayjs(postProp.publishedDate), 'hours') <= 24 : false,
            hasHashtags: Array.isArray(postProp.hashtags) && postProp.hashtags.length > 0,
            hasUserInteractions: postProp.userInteractions && typeof postProp.userInteractions === 'object',
        };
        
        // Extract author info
        const authorInfo = {
            name: postProp.author?.name || postProp.authorName || t('defaultAuthor'),
            username: postProp.author?.username || postProp.authorUsername,
            avatar: postProp.author?.avatar || postProp.authorAvatar || '/images/default-avatar.png',
            verified: extractedFeatures.isVerifiedAuthor,
            bio: postProp.author?.bio || '',
            followers: postProp.author?.followers || 0,
            following: postProp.author?.following || 0,
        };
        
        // Extract hashtags
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
            id: postProp.id || postProp._id,
            views,
            likes,
            dislikes,
            comments,
            bookmarks,
            shares,
            heatScore: calculatedHeatScore,
            features: extractedFeatures,
            author: authorInfo,
            hashtags: extractedHashtags,
            title: postProp.title || postProp.name || t('unnamedPost'),
            content: postProp.content || postProp.description || postProp.summary || t('noContentAvailable'),
            excerpt: postProp.excerpt || postProp.summary || '',
            publishedDate: postProp.publishedDate || postProp.createdAt,
            updatedDate: postProp.updatedAt || postProp.publishedDate || postProp.createdAt,
            location: postProp.location || postProp.author?.location,
            categorySlug: postProp.categorySlug || postProp.category?.slug || '',
            categoryName: postProp.categoryName || postProp.category?.name || '',
        };
    }, [postProp, t]);
    
    // State management
    const [state, setState] = useState({
        currentImageIndex: 0,
        showComments: false,
        showShareMenu: false,
        liked: post?.userInteractions?.liked || false,
        disliked: post?.userInteractions?.disliked || false,
        saved: post?.userInteractions?.saved || false,
        flagged: post?.userInteractions?.flagged || false,
        zoomedImage: null,
    });
    
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoadingMedia, setIsLoadingMedia] = useState(false);
    
    // Counts state for optimistic updates
    const [counts, setCounts] = useState({
        likes: post?.likes || 0,
        dislikes: post?.dislikes || 0,
        bookmarks: post?.bookmarks || 0,
        shares: post?.shares || 0,
    });
    
    // Load media
    useEffect(() => {
        const loadMedia = async () => {
            if (!post?.id) return;
            
            setIsLoadingMedia(true);
            try {
                const mediaResponse = await getMediaByItem(post.id);
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || 'image',
                        alt: item.altText || post.title || 'Post image'
                    }));
                    setMediaItems(formattedMedia);
                } else if (post.image || post.thumbnail || post.coverImage) {
                    setMediaItems([{ 
                        url: post.image || post.thumbnail || post.coverImage, 
                        type: 'image', 
                        alt: post.title 
                    }]);
                }
            } catch (error) {
                // Error: 'Error loading media:', error...
                if (post.image || post.thumbnail || post.coverImage) {
                    setMediaItems([{ 
                        url: post.image || post.thumbnail || post.coverImage, 
                        type: 'image', 
                        alt: post.title 
                    }]);
                }
            } finally {
                setIsLoadingMedia(false);
            }
        };
        loadMedia();
    }, [post?.id, post?.title, post?.image, post?.thumbnail, post?.coverImage]);
    
    // Handlers
    const handleImageNavigation = useCallback((direction) => {
        setState(prev => ({
            ...prev,
            currentImageIndex: direction === 'next'
                ? (prev.currentImageIndex + 1) % mediaItems.length
                : (prev.currentImageIndex - 1 + mediaItems.length) % mediaItems.length
        }));
    }, [mediaItems.length]);
    
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            
            return;
        }
        
        handleLike(post.id, userId);
        setState(prev => ({ ...prev, liked: true, disliked: false }));
        setCounts(prev => ({
            ...prev,
            likes: prev.likes + (state.liked ? 0 : 1),
            dislikes: prev.dislikes - (state.disliked ? 1 : 0)
        }));
    }, [post?.id, userId, state.liked, state.disliked, handleLike]);
    
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            
            return;
        }
        
        handleDislike(post.id, userId);
        setState(prev => ({ ...prev, liked: false, disliked: true }));
        setCounts(prev => ({
            ...prev,
            likes: prev.likes - (state.liked ? 1 : 0),
            dislikes: prev.dislikes + (state.disliked ? 0 : 1)
        }));
    }, [post?.id, userId, state.liked, state.disliked, handleDislike]);
    
    const handleSaveClick = useCallback(() => {
        if (!userId) {
            
            return;
        }
        
        setState(prev => ({ ...prev, saved: !prev.saved }));
        setCounts(prev => ({
            ...prev,
            bookmarks: prev.bookmarks + (state.saved ? -1 : 1)
        }));
    }, [userId, state.saved]);
    
    const handleShare = useCallback(() => {
        if (navigator.share && typeof window !== 'undefined') {
            navigator.share({
                title: post.title,
                text: post.excerpt || post.content.substring(0, 200),
                url: window.location.href
            }).catch(() => {});
        } else {
            navigator.clipboard.writeText(window.location.href);
            setCounts(prev => ({ ...prev, shares: prev.shares + 1 }));
        }
    }, [post?.title, post?.excerpt, post?.content]);
    
    const handleContactAuthor = useCallback(() => {
        if (!userId || !post?.author?.username) {
            
            return;
        }
        
        dispatch(openMessageModal({
            recipientId: post.author.username,
            itemId: post.id
        }));
    }, [post?.author?.username, post?.id, userId, dispatch]);
    
    if (!post) {
        return (
            <div className={styles.errorState}>
                <p>{t('postNotFound')}</p>
            </div>
        );
    }
    
    // Format dates
    const publishedDate = dayjs(post.publishedDate);
    const formattedDate = publishedDate.format('MMM DD, YYYY');
    const relativeDate = publishedDate.fromNow();
    
    return (
        <div className={styles.container}>
            <article className={styles.detailView}>
                {/* Header Section */}
                <header className={styles.header}>
                    {/* Breadcrumb */}
                    <nav className={styles.breadcrumb}>
                        <Link href="/" className={styles.breadcrumbLink}>Home</Link>
                        <ChevronRight size={16} />
                        <Link href="/posts" className={styles.breadcrumbLink}>Posts</Link>
                        {post.categoryName && (
                            <>
                                <ChevronRight size={16} />
                                <Link href={`/posts/${post.categorySlug}`} className={styles.breadcrumbLink}>
                                    {post.categoryName}
                                </Link>
                            </>
                        )}
                        <ChevronRight size={16} />
                        <span className={styles.breadcrumbCurrent}>{post.title}</span>
                    </nav>
                    
                    {/* Title and Meta */}
                    <h1 className={styles.title}>{post.title}</h1>
                    
                    {post.excerpt && (
                        <p className={styles.excerpt}>{post.excerpt}</p>
                    )}
                    
                    {/* Author and Date */}
                    <div className={styles.meta}>
                        <div className={styles.authorSection}>
                            <img 
                                src={post.author.avatar} 
                                alt={post.author.name}
                                className={styles.authorAvatar}
                            />
                            <div className={styles.authorInfo}>
                                <div className={styles.authorName}>
                                    {post.author.name}
                                    {post.author.verified && (
                                        <Check size={16} className={styles.verifiedBadge} />
                                    )}
                                </div>
                                <div className={styles.publishInfo}>
                                    <time dateTime={post.publishedDate} title={formattedDate}>
                                        {relativeDate}
                                    </time>
                                    <span className={styles.separator}>•</span>
                                    <span>{post.views} views</span>
                                    {post.location && (
                                        <>
                                            <span className={styles.separator}>•</span>
                                            <span className={styles.location}>
                                                <MapPin size={14} />
                                                {post.location}
                                            </span>
                                        </>
                                    )}
                                </div>
                            </div>
                        </div>
                        
                        <div className={styles.metaActions}>
                            <button 
                                className={styles.followButton}
                                onClick={handleContactAuthor}
                            >
                                <MessageCircle size={16} />
                                Contact
                            </button>
                            <button className={styles.moreButton}>
                                <MoreHorizontal size={20} />
                            </button>
                        </div>
                    </div>
                </header>
                
                {/* Media Section */}
                {mediaItems.length > 0 && (
                    <div className={styles.mediaSection}>
                        <div className={styles.mainImageContainer}>
                            {isLoadingMedia ? (
                                <div className={styles.imagePlaceholder}>
                                    <div className={styles.spinner} />
                                </div>
                            ) : (
                                <>
                                    <img
                                        src={mediaItems[state.currentImageIndex]?.url}
                                        alt={mediaItems[state.currentImageIndex]?.alt || post.title}
                                        className={styles.mainImage}
                                        onClick={() => setState(prev => ({ 
                                            ...prev, 
                                            zoomedImage: mediaItems[state.currentImageIndex]?.url 
                                        }))}
                                    />
                                    
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
                                            <div className={styles.imageCounter}>
                                                {state.currentImageIndex + 1} / {mediaItems.length}
                                            </div>
                                        </>
                                    )}
                                </>
                            )}
                        </div>
                        
                        {/* Thumbnail Gallery */}
                        {mediaItems.length > 1 && (
                            <div className={styles.thumbnailGallery}>
                                {mediaItems.map((media, index) => (
                                    <button
                                        key={index}
                                        className={`${styles.thumbnail} ${
                                            index === state.currentImageIndex ? styles.activeThumbnail : ''
                                        }`}
                                        onClick={() => setState(prev => ({ 
                                            ...prev, 
                                            currentImageIndex: index 
                                        }))}
                                    >
                                        <img src={media.url} alt={`Image ${index + 1}`} />
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                )}
                
                {/* Content Section */}
                <div className={styles.contentSection}>
                    <div className={styles.content} dangerouslySetInnerHTML={{ __html: sanitizeRichHtml(post.content) }} />
                    
                    {/* Tags */}
                    {post.hashtags.length > 0 && (
                        <div className={styles.tags}>
                            {post.hashtags.map((tag, index) => (
                                <Link 
                                    key={index} 
                                    href={`/search?tag=${encodeURIComponent(tag)}`}
                                    className={styles.tag}
                                >
                                    <Tag size={14} />
                                    {tag}
                                </Link>
                            ))}
                        </div>
                    )}
                </div>
                
                {/* Engagement Bar */}
                <div className={styles.engagementBar}>
                    <div className={styles.engagementStats}>
                        <span className={styles.statItem}>
                            <Eye size={18} />
                            {post.views} views
                        </span>
                        <span className={styles.statItem}>
                            <MessageCircle size={18} />
                            {post.comments} comments
                        </span>
                        <span className={styles.statItem}>
                            <Share2 size={18} />
                            {counts.shares} shares
                        </span>
                    </div>
                    
                    <div className={styles.engagementActions}>
                        <button
                            className={`${styles.actionButton} ${state.liked ? styles.active : ''}`}
                            onClick={handleLikeClick}
                        >
                            <ThumbsUp size={20} />
                            <span>{counts.likes}</span>
                        </button>
                        
                        <button
                            className={`${styles.actionButton} ${state.disliked ? styles.active : ''}`}
                            onClick={handleDislikeClick}
                        >
                            <ThumbsDown size={20} />
                            <span>{counts.dislikes}</span>
                        </button>
                        
                        <button
                            className={`${styles.actionButton} ${state.saved ? styles.active : ''}`}
                            onClick={handleSaveClick}
                        >
                            <Bookmark size={20} />
                            <span>{counts.bookmarks}</span>
                        </button>
                        
                        <button
                            className={styles.actionButton}
                            onClick={handleShare}
                        >
                            <Share2 size={20} />
                            <span>Share</span>
                        </button>
                        
                        <button
                            className={`${styles.actionButton} ${state.flagged ? styles.active : ''}`}
                            onClick={() => setState(prev => ({ ...prev, flagged: !prev.flagged }))}
                        >
                            <Flag size={20} />
                        </button>
                    </div>
                </div>
                
                {/* Author Bio Section */}
                {post.author.bio && (
                    <div className={styles.authorBio}>
                        <h3>About the Author</h3>
                        <div className={styles.bioContent}>
                            <img 
                                src={post.author.avatar} 
                                alt={post.author.name}
                                className={styles.bioAvatar}
                            />
                            <div className={styles.bioText}>
                                <h4>
                                    {post.author.name}
                                    {post.author.verified && <Check size={16} className={styles.verifiedBadge} />}
                                </h4>
                                <p>{post.author.bio}</p>
                                <div className={styles.authorStats}>
                                    <span>{post.author.followers} followers</span>
                                    <span className={styles.separator}>•</span>
                                    <span>{post.author.following} following</span>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
                
                {/* Comments Section */}
                <div className={styles.commentsSection}>
                    <h3 className={styles.commentsTitle}>
                        <MessageCircle size={24} />
                        Comments ({post.comments})
                    </h3>
                    
                    <button
                        className={styles.toggleComments}
                        onClick={() => setState(prev => ({ ...prev, showComments: !prev.showComments }))}
                    >
                        {state.showComments ? 'Hide Comments' : 'Show Comments'}
                    </button>
                    
                    {state.showComments && (
                        <CommentsSetup
                            userId={userId}
                            itemId={post.id}
                            itemType="post"
                            categoryId={post.categoryId}
                        />
                    )}
                </div>
                
                {/* Related Posts */}
                <div className={styles.relatedSection}>
                    <h3>Related Posts</h3>
                    <div className={styles.relatedGrid}>
                        {/* Placeholder for related posts - would be fetched from API */}
                        <p className={styles.placeholder}>Related posts would appear here</p>
                    </div>
                </div>
                
                {/* Image Zoom Modal */}
                {state.zoomedImage && (
                    <div 
                        className={styles.zoomModal} 
                        onClick={() => setState(prev => ({ ...prev, zoomedImage: null }))}
                    >
                        <img src={state.zoomedImage} alt="Zoomed post image" />
                        <button className={styles.closeZoom} aria-label="Close zoom">×</button>
                    </div>
                )}
            </article>
        </div>
    );
});

DetailedPostView.displayName = 'DetailedPostView';

export default DetailedPostView;
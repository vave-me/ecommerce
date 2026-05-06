"use client";
import React, {useCallback, useState} from 'react';
import {Image} from '@/icons';
import styles from './page.module.css';
import UserHeader from "../UserHeader";
import {Engagement} from "../Engagement";
const ImprovedSocialCard = ({tweet}) => {
    // Extract real data from API with fallbacks
    const actualTweet = tweet?.tweet || tweet;
    const tweetMetrics = actualTweet?.metrics || {};
    // Convert string metrics to numbers
    const convertMetricToNumber = (value) => parseInt(value || '0', 10);
    // Extract hashtags from content
    const extractHashtags = (content) => {
        if (!content) return [];
        const hashtagRegex = /#[\w]+/g;
        const matches = content.match(hashtagRegex);
        return matches ? matches.map(tag => tag.substring(1)) : [];
    };
    // Calculate time ago from creation date
    const getTimeAgo = (createdAt) => {
        if (!createdAt) return 'now';
        const now = new Date();
        const created = new Date(createdAt);
        const diffInMinutes = Math.floor((now - created) / (1000 * 60));
        if (diffInMinutes < 1) return 'now';
        if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
        const diffInHours = Math.floor(diffInMinutes / 60);
        if (diffInHours < 24) return `${diffInHours}h ago`;
        const diffInDays = Math.floor(diffInHours / 24);
        if (diffInDays < 7) return `${diffInDays}d ago`;
        return created.toLocaleDateString();
    };
    // Create author info from tweet data
    const getAuthorInfo = () => {
        return {
            name: actualTweet?.authorName || `User #${(actualTweet?.userId || '').slice(-6)}`,
            username: actualTweet?.authorUsername || `@user${(actualTweet?.userId || '').slice(-6)}`,
            avatar: actualTweet?.authorAvatar || "/images/default-avatar.webp",
            verified: actualTweet?.authorVerified || false
        };
    };
    // Real post data from API
    const post = {
        // Core tweet information
        id: actualTweet?.id || '',
        content: actualTweet?.content || actualTweet?.description || "No content provided",
        // Author information
        author: getAuthorInfo(),
        // Media
        image: actualTweet?.thumbnail || actualTweet?.image || null,
        // Timing and location
        timeAgo: getTimeAgo(actualTweet?.createdAt),
        location: actualTweet?.location || actualTweet?.city || null,
        // Real metrics from API
        likes: convertMetricToNumber(tweetMetrics.likesCount),
        dislikes: convertMetricToNumber(tweetMetrics.dislikesCount),
        comments: convertMetricToNumber(tweetMetrics.commentsCount),
        views: convertMetricToNumber(tweetMetrics.visitedCount),
        bookmarks: convertMetricToNumber(tweetMetrics.addedToWishlistCount),
        shares: convertMetricToNumber(tweetMetrics.sharedCount),
        // Extracted features
        hashtags: extractHashtags(actualTweet?.content || actualTweet?.description),
        // Metrics object for Engagement component
        metrics: {
            likes: convertMetricToNumber(tweetMetrics.likesCount),
            dislikes: convertMetricToNumber(tweetMetrics.dislikesCount),
            comments: convertMetricToNumber(tweetMetrics.commentsCount),
            shares: convertMetricToNumber(tweetMetrics.sharedCount),
            views: convertMetricToNumber(tweetMetrics.visitedCount),
        },
        // Entity metadata
        status: actualTweet?.status || 'active',
        entityType: actualTweet?.entityType || 'post',
        userId: actualTweet?.userId || '',
    };
    const [showComments, setShowComments] = useState(false);
    // Interactive states
    const [favorite, setFavorite] = useState(false);
    const [liked, setLiked] = useState(false);
    const [disliked, setDisliked] = useState(false);
    const [saved, setSaved] = useState(false);
    const [shared, setShared] = useState(false);
    const [flagged, setFlagged] = useState(false);
    const [expanded, setExpanded] = useState(false);
    const [likesCount, setLikesCount] = useState(post.likes);
    const [dislikesCount, setDislikesCount] = useState(post.dislikes);
    const [savesCount, setSavesCount] = useState(post.bookmarks);
    const [sharesCount, setSharesCount] = useState(post.shares);
    // Handle like action
    const handleLikeClick = () => {
        if (liked) {
            setLikesCount(likesCount - 1);
        } else {
            setLikesCount(likesCount + 1);
            // If disliked, remove dislike when liking
            if (disliked) {
                setDisliked(false);
                setDislikesCount(dislikesCount - 1);
            }
        }
        setLiked(!liked);
    };
    // Toggle comments section
    const handleCommentClick = () => {
        setShowComments(!showComments);
    };
    // Handle dislike action
    const handleDislikeClick = () => {
        if (disliked) {
            setDislikesCount(dislikesCount - 1);
        } else {
            setDislikesCount(dislikesCount + 1);
            // If liked, remove like when disliking
            if (liked) {
                setLiked(false);
                setLikesCount(likesCount - 1);
            }
        }
        setDisliked(!disliked);
    };
    // Handle bookmark action
    const handleSaveClick = () => {
        if (saved) {
            setSavesCount(savesCount - 1);
        } else {
            setSavesCount(savesCount + 1);
        }
        setSaved(!saved);
    };
    const handleFavorite = useCallback(() => {
        setFavorite(prev => !prev);
    }, []);
    // Handle share action
    const handleShare = () => {
        if (!shared) {
            setSharesCount(sharesCount + 1);
            setShared(true);
            // Reset share state after a moment to allow re-sharing
            setTimeout(() => setShared(false), 2000);
        }
    };
    // Toggle expanded content
    const toggleExpanded = () => {
        setExpanded(!expanded);
    };
    // Format numbers for display
    const formatNumber = (num) => {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num;
    };
    return (
        <div className={styles.card}>
            {/* Card content */}
            <div className={styles.cardContent}>
                {/* Author header */}
                <UserHeader/>
                {/* Post content */}
                <div className={styles.postContent}>
                    <p className={`${styles.postText} ${expanded ? styles.expanded : ''}`}>
                        {post.content}
                    </p>
                    {post.content.length > 150 && !expanded && (
                        <button
                            className={styles.expandButton}
                            onClick={toggleExpanded}
                        >
                            Read more
                        </button>
                    )}
                    {/* Hashtags */}
                    {post.hashtags.length > 0 && (
                        <div className={styles.hashtagsContainer}>
                            {post.hashtags.map((tag, index) => (
                                <span key={index} className={styles.hashtag}>
                                    #{tag}
                                </span>
                            ))}
                        </div>
                    )}
                </div>
                {/* Post image */}
                {post.image && (
                    <div className={styles.imageContainer}>
                        <img
                            src={post.image}
                            alt="Post attachment"
                            className={styles.postImage}
                        />
                    </div>
                )}
                <Engagement
                    metrics={post.metrics}
                    liked={liked}
                    disliked={disliked}
                    favorite={favorite}
                    onLike={handleLikeClick}
                    onCommentClick={handleCommentClick}
                    onDislike={handleDislikeClick}
                    onFavorite={handleFavorite}
                />
            </div>
            {showComments && (
                <div className={styles.commentSection}>
                    <div className={styles.commentHeader}>
                        <h4 className={styles.commentTitle}>
                            Comments ({post.metrics.comments})
                        </h4>
                        <button className={styles.viewAllButton}>
                            View All
                        </button>
                    </div>
                    {/* Sample comment - would be replaced with real comments in production */}
                    <div className={styles.comment}>
                        <img
                            src="/images/default-avatar.webp"
                            alt="Commenter"
                            className={styles.commentAvatar}
                        />
                        <div className={styles.commentContent}>
                            <div className={styles.commentMeta}>
                                <div className={styles.commenterInfo}>
                                    <span className={styles.commenterName}>Community Member</span>
                                    <span className={styles.commenterUsername}>@member</span>
                                </div>
                                <span className={styles.commentTime}>1h</span>
                            </div>
                            <p className={styles.commentText}>
                                Thanks for sharing this! Very helpful content.
                            </p>
                        </div>
                    </div>
                    {/* Comment form */}
                    <div className={styles.commentForm}>
                        <img
                            src="/images/default-avatar.webp"
                            alt="Your avatar"
                            className={styles.commentAvatar}
                        />
                        <div className={styles.commentInputWrapper}>
                            <input
                                type="text"
                                placeholder="Write a comment..."
                                className={styles.commentInput}
                            />
                            <button className={styles.commentButton}>
                                <Image size={14}/>
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
// Preview Component that accepts tweets array
const TweetsPreview = ({ tweets = [] }) => {
    // Handle loading state
    if (!tweets || tweets.length === 0) {
        return (
            <div className={styles.previewContainer}>
                <div className={styles.emptyState}>
                    <p>No tweets available to display.</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.previewContainer}>
            {tweets.map((tweet) => (
                <ImprovedSocialCard key={tweet.id || tweet.tweet?.id} tweet={tweet} />
            ))}
        </div>
    );
};
export default TweetsPreview;
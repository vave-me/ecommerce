"use client"
import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark, ThumbsDown, Eye, Flag, MoreHorizontal, Check, Clock, MapPin } from '@/icons';
import styles from './page.module.css';
const SimplifiedFinalCard = () => {
    // Sample post data
    const post = {
        author: {
            name: "Sarah Johnson",
            username: "@sarahjdesigns",
            avatar: "/images/logo-small.webp",
            verified: true
        },
        content: "Just finished the redesign for our client's e-commerce platform! It's amazing how small UI improvements can lead to significant increases in conversion rates. #UXDesign #UI",
        image: "/images/people.jpg",
        timeAgo: "2h ago",
        location: "Berlin",
        likes: 127,
        comments: 23,
        views: 1842,
        dislikes: 4,
        bookmarks: 18,
        shares: 9,
        hashtags: ["UXDesign", "UI", "ConversionOptimization"]
    };
    // Interactive states
    const [liked, setLiked] = useState(false);
    const [disliked, setDisliked] = useState(false);
    const [saved, setSaved] = useState(false);
    const [shared, setShared] = useState(false);
    const [flagged, setFlagged] = useState(false);
    const [likesCount, setLikesCount] = useState(post.likes);
    const [dislikesCount, setDislikesCount] = useState(post.dislikes);
    const [savesCount, setSavesCount] = useState(post.bookmarks);
    const [sharesCount, setSharesCount] = useState(post.shares);
    // Handle like action
    const handleLike = () => {
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
    // Handle dislike action
    const handleDislike = () => {
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
    const handleSave = () => {
        if (saved) {
            setSavesCount(savesCount - 1);
        } else {
            setSavesCount(savesCount + 1);
        }
        setSaved(!saved);
    };
    // Handle share action
    const handleShare = () => {
        if (!shared) {
            setSharesCount(sharesCount + 1);
            setShared(true);
            // Reset share state after a moment to allow re-sharing
            setTimeout(() => setShared(false), 2000);
        }
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
        <div className={styles.cardContainer}>
            {/* Post header */}
            <div className={styles.postHeader}>
                {/* Author avatar */}
                <div className={styles.avatarContainer}>
                    <img
                        src={post.author.avatar}
                        alt={post.author.name}
                        className={styles.avatar}
                    />
                    {post.author.verified && (
                        <div className={styles.verifiedBadge}>
                            <Check size={12} className={styles.verifiedIcon} />
                        </div>
                    )}
                </div>
                {/* Author info */}
                <div className={styles.authorInfo}>
                    <div className={styles.authorNameRow}>
                        <span className={styles.authorName}>{post.author.name}</span>
                        <span className={styles.authorUsername}>{post.author.username}</span>
                    </div>
                    <div className={styles.postMetaInfo}>
                        <Clock size={12} className={styles.metaIcon} />
                        <span>{post.timeAgo}</span>
                        {post.location && (
                            <>
                                <span className={styles.metaSeparator}>•</span>
                                <MapPin size={12} className={styles.metaIcon} />
                                <span>{post.location}</span>
                            </>
                        )}
                    </div>
                </div>
                {/* Views counter */}
                <div className={styles.viewsCounter}>
                    <Eye size={12} className={styles.viewsIcon} />
                    <span>{formatNumber(post.views)}</span>
                </div>
                {/* Post options */}
                <button className={styles.optionsButton}>
                    <MoreHorizontal size={18} />
                </button>
            </div>
            {/* Post content */}
            <div className={styles.postContent}>
                <p className={styles.postText}>{post.content}</p>
                {/* Hashtags */}
                <div className={styles.hashtagsContainer}>
                    {post.hashtags.map((tag, index) => (
                        <span key={index} className={styles.hashtag}>
                            #{tag}
                        </span>
                    ))}
                </div>
            </div>
            {/* Post image with interactive buttons on the right */}
            <div className={styles.mediaContainer}>
                <img
                    src={post.image}
                    alt="Post attachment"
                    className={styles.postImage}
                />
                {/* Interactive buttons in right column */}
                <div className={styles.interactionColumn}>
                    {/* Like button */}
                    <button
                        className={`${styles.interactionButton} ${liked ? styles.likedButton : ''}`}
                        onClick={handleLike}
                    >
                        <Heart
                            size={18}
                            className={liked ? styles.filledIcon : ''}
                        />
                        <span className={styles.countLabel}>{formatNumber(likesCount)}</span>
                    </button>
                    {/* Dislike button */}
                    <button
                        className={`${styles.interactionButton} ${disliked ? styles.dislikedButton : ''}`}
                        onClick={handleDislike}
                    >
                        <ThumbsDown size={18} />
                        <span className={styles.countLabel}>{formatNumber(dislikesCount)}</span>
                    </button>
                    {/* Comment button */}
                    <button className={styles.interactionButton}>
                        <MessageCircle size={18} />
                        <span className={styles.countLabel}>{formatNumber(post.comments)}</span>
                    </button>
                    {/* Save button */}
                    <button
                        className={`${styles.interactionButton} ${saved ? styles.savedButton : ''}`}
                        onClick={handleSave}
                    >
                        <Bookmark
                            size={18}
                            className={saved ? styles.filledIcon : ''}
                        />
                        <span className={styles.countLabel}>{formatNumber(savesCount)}</span>
                    </button>
                    {/* Share button */}
                    <button
                        className={`${styles.interactionButton} ${shared ? styles.sharedButton : ''}`}
                        onClick={handleShare}
                    >
                        <Share2 size={18} />
                        <span className={styles.countLabel}>{formatNumber(sharesCount)}</span>
                    </button>
                    {/* Report button */}
                    <button
                        className={`${styles.reportButton} ${flagged ? styles.reportedButton : ''}`}
                        onClick={() => setFlagged(!flagged)}
                    >
                        <Flag size={18} />
                    </button>
                </div>
            </div>
        </div>
    );
};
const SimplifiedFinalCardPreview = () => {
    return (
        <div className={styles.previewContainer}>
            <SimplifiedFinalCard />
        </div>
    );
};
export default SimplifiedFinalCardPreview;
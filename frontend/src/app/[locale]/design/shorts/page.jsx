"use client";
import React, {useState, useRef, useEffect, useCallback} from 'react';
import {
    Heart,
    Music,
    Volume2,
    VolumeX,
    X,
    Play,
    Pause
} from '@/icons';
import styles from './page.module.css';
import UserHeader from "../UserHeader";
import {Engagement} from "../Engagement";
const ImprovedCompactVideo = () => {
    // Sample video data
    const videoData = {
        id: "short1",
        caption: "UI tricks that boost conversion rates",
        description: "Simple design changes that increased conversion by 300%! #UXDesign",
        videoSrc:  "/images/back.jpg",
        thumbnailSrc:  "/images/back.jpg",
        soundName: "Original Sound",
        likes: 145200,
        comments: 2341,
        shares: 3890,
        saves: 18400,
        metrics: {
            likes: 70,
            comments: 30,
            shares: 16,
            views: 100000
        },
        creator: {
            name: "UX Master",
            username: "@ux_master",
            avatar: "/api/placeholder/36/36",
            verified: true,
        },
        topComments: [
            {
                userName: "Designer101",
                userAvatar: "/api/placeholder/28/28",
                text: "That contrast tip is game-changing!",
                timeAgo: "2h",
                likes: 423
            },
            {
                userName: "CodeCrafter",
                userAvatar: "/api/placeholder/28/28",
                text: "Tutorial for checkout optimization?",
                timeAgo: "4h",
                likes: 215
            }
        ],
        tags: ["UXDesign", "UI", "ProductDesign"]
    };
    const [favorite, setFavorite] = useState(false);
    const [liked, setLiked] = useState(false);
    const [disliked, setDisliked] = useState(false);
    const [isPlaying, setIsPlaying] = useState(false);
    const [isMuted, setIsMuted] = useState(true);
    const [isLiked, setIsLiked] = useState(false);
    const [isSaved, setIsSaved] = useState(false);
    const [isShared, setIsShared] = useState(false);
    const [showComments, setShowComments] = useState(false);
    const [isFollowing, setIsFollowing] = useState(false);
    const [showFullDescription, setShowFullDescription] = useState(false);
    const [isVideoLoaded, setIsVideoLoaded] = useState(false);
    const videoRef = useRef(null);
    // Auto-play when component is visible
    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    attemptPlay();
                } else {
                    if (videoRef.current) {
                        videoRef.current.pause();
                        setIsPlaying(false);
                    }
                }
            },
            {threshold: 0.7}
        );
        if (videoRef.current) {
            observer.observe(videoRef.current);
        }
        return () => {
            if (videoRef.current) {
                observer.unobserve(videoRef.current);
            }
        };
    }, [isVideoLoaded]);
    // Handle video loaded
    const handleVideoLoaded = () => {
        setIsVideoLoaded(true);
    };
    const handleCommentClick = useCallback(() => {
        setShowComments(prev => !prev);
    }, []);
    const handleFavorite = useCallback(() => {
        setFavorite(prev => !prev);
    }, []);
    const handleLikeClick = useCallback(() => {
        if (disliked) setDisliked(false);
        setLiked(prev => !prev);
    }, [disliked]);
    const handleDislikeClick = useCallback(() => {
        if (liked) setLiked(false);
        setDisliked(prev => !prev);
    }, [liked]);
    // Attempt to play the video
    const attemptPlay = () => {
        if (videoRef.current) {
            videoRef.current.play()
                .then(() => {
                    setIsPlaying(true);
                })
                .catch(err => {
                    setIsPlaying(false);
                });
        }
    };
    // Toggle play/pause on tap
    const togglePlay = () => {
        if (!videoRef.current) return;
        if (isPlaying) {
            videoRef.current.pause();
            setIsPlaying(false);
        } else {
            attemptPlay();
        }
    };
    // Toggle mute
    const toggleMute = (e) => {
        e.stopPropagation();
        if (videoRef.current) {
            videoRef.current.muted = !isMuted;
            setIsMuted(!isMuted);
        }
    };
    // Like video
    const handleLike = (e) => {
        e.stopPropagation();
        setIsLiked(!isLiked);
    };
    // Save video
    const handleSave = (e) => {
        e.stopPropagation();
        setIsSaved(!isSaved);
    };
    // Share video
    const handleShare = (e) => {
        e.stopPropagation();
        setIsShared(true);
        // Reset share visual feedback after a moment
        setTimeout(() => {
            setIsShared(false);
        }, 1000);
    };
    // Toggle comments
    const handleComments = (e) => {
        if (e) e.stopPropagation();
        setShowComments(!showComments);
    };
    // Toggle follow
    const handleFollow = (e) => {
        e.stopPropagation();
        setIsFollowing(!isFollowing);
    };
    // Toggle full description
    const toggleDescription = (e) => {
        e.stopPropagation();
        setShowFullDescription(!showFullDescription);
    };
    // Format number for display (e.g., 1.2M)
    const formatCount = (count) => {
        if (count >= 1000000) {
            return (count / 1000000).toFixed(1) + 'M';
        } else if (count >= 1000) {
            return (count / 1000).toFixed(1) + 'K';
        }
        return count;
    };
    return (
        <div className={styles.shortContainer}>
            <UserHeader/>
            <div className={styles.videoCard}>
                {/* Main video container with tap to play/pause */}
                <div className={styles.videoWrapper} onClick={togglePlay}>
                    <video
                        ref={videoRef}
                        src={videoData.videoSrc}
                        className={styles.videoElement}
                        loop
                        muted={isMuted}
                        playsInline
                        poster={videoData.thumbnailSrc}
                        onLoadedData={handleVideoLoaded}
                    />
                    {/* Play/pause overlay */}
                    <div className={`${styles.playPauseOverlay} ${isPlaying ? styles.playing : ''}`}>
                        {isPlaying ?
                            <Pause className={styles.playPauseIcon}/> :
                            <Play className={styles.playPauseIcon}/>
                        }
                    </div>
                    {/* Mute/unmute button */}
                    <button
                        className={styles.muteButton}
                        onClick={toggleMute}
                        aria-label={isMuted ? "Unmute" : "Mute"}
                    >
                        {isMuted ? <VolumeX size={18}/> : <Volume2 size={18}/>}
                    </button>
                    {/* Sound attribution */}
                    <div className={styles.soundAttribution}>
                        <Music size={14} className={styles.musicIcon}/>
                        <div className={styles.soundContainer}>
                            <span className={styles.soundName}>{videoData.soundName}</span>
                        </div>
                    </div>
                    {/* Caption overlay */}
                    <div className={styles.captionOverlay}>
                        <div
                            className={`${styles.captionContainer} ${showFullDescription ? styles.expanded : ''}`}
                            onClick={toggleDescription}
                        >
                            <p className={styles.caption}>{videoData.caption}</p>
                            <p className={styles.description}>{videoData.description}</p>
                            <div className={styles.hashtagsContainer}>
                                {videoData.tags.map((tag, index) => (
                                    <span key={index} className={styles.hashtag}>#{tag}</span>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            {/* Engagement bar */}
            <Engagement
                liked={videoData.metrics.likes}
                disliked={videoData.metrics.disliked}
                favorite={favorite}
                onLike={handleLikeClick}
                onCommentClick={handleCommentClick}
                onDislike={handleDislikeClick}
                onFavorite={handleFavorite}
            />
            {
                showComments && (
                    <div className={styles.commentsDrawer}>
                        <div className={styles.commentsHeader}>
                            <h4 className={styles.commentsTitle}>
                                Comments ({formatCount(videoData.comments)})
                            </h4>
                            <button
                                className={styles.closeButton}
                                onClick={handleComments}
                                aria-label="Close comments"
                            >
                                <X size={18}/>
                            </button>
                        </div>
                        <div className={styles.commentsList}>
                            {videoData.topComments.map((comment, index) => (
                                <div key={index} className={styles.commentItem}>
                                    <img
                                        src={comment.userAvatar}
                                        alt={comment.userName}
                                        className={styles.commentAvatar}
                                    />
                                    <div className={styles.commentContent}>
                                        <div className={styles.commentHeader}>
                                            <span className={styles.commentUserName}>{comment.userName}</span>
                                            <span className={styles.commentTime}>{comment.timeAgo}</span>
                                        </div>
                                        <p className={styles.commentText}>{comment.text}</p>
                                        <div className={styles.commentActions}>
                                            <button className={styles.commentLikeButton}>
                                                <Heart size={14} className={styles.commentActionIcon}/>
                                                <span>{formatCount(comment.likes)}</span>
                                            </button>
                                            <button className={styles.commentReplyButton}>
                                                Reply
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                        <div className={styles.commentInput}>
                            <input
                                type="text"
                                placeholder="Add a comment..."
                                className={styles.commentInputField}
                            />
                            <button className={styles.commentPostButton}>Post</button>
                        </div>
                    </div>
                )
            }
        </div>
    )
        ;
};
export default ImprovedCompactVideo;
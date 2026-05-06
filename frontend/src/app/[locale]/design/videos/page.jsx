"use client";
import React, {useState, useRef, useCallback, useEffect} from 'react';
import {
    Heart,
    MessageCircle,
    Play,
    Pause,
    Volume2,
    VolumeX,
    Clock,
    Eye,
    BarChart4,
    Loader,
    ArrowUpRight
} from '@/icons';
import styles from './page.module.css';
import UserHeader from "../UserHeader";
import {Engagement} from "../Engagement";
const ImprovedVideoCard = () => {
    // Sample video data
    const videoData = {
        id: "v123456",
        creator: {
            id: "u789012",
            name: "VideoCreator",
            username: "@videocreator",
            avatar: "/api/placeholder/48/48",
            verified: true,
            following: false,
            followers: 1250000
        },
        title: "10 UI Tricks That Boosted Our Conversion Rate by 300%",
        description: "In this video I share the exact UI improvements we made that tripled our client's e-commerce conversion rate in just 2 weeks! #UXDesign #ConversionRate",
        videoSrc: "/images/psyche.jpg", // This would be a real video URL in production
        thumbnail:  "/images/psyche.jpg",
        duration: "3:45",
        timestamp: "2025-03-22T14:35:00Z",
        postedAgo: "2 days ago",
        metrics: {
            views: 345200,
            likes: 42300,
            comments: 1842,
            shares: 5270,
            saves: 8760
        },
        hashtags: ["UXDesign", "ConversionRate", "UI"],
        isLiked: false,
        isDisliked: false,
        isSaved: false,
        isShared: false
    };
    const [favorite, setFavorite] = useState(false);
    const [isPlaying, setIsPlaying] = useState(false);
    const [isMuted, setIsMuted] = useState(true);
    const [isLoading, setIsLoading] = useState(false);
    const [progress, setProgress] = useState(0);
    const [currentTime, setCurrentTime] = useState("0:00");
    const [liked, setLiked] = useState(videoData.isLiked);
    const [disliked, setDisliked] = useState(videoData.isDisliked);
    const [saved, setSaved] = useState(videoData.isSaved);
    const [likesCount, setLikesCount] = useState(videoData.metrics.likes);
    const [isFollowing, setIsFollowing] = useState(videoData.creator.following);
    const [showActions, setShowActions] = useState(false);
    const [showComments, setShowComments] = useState(false);
    const [showDescription, setShowDescription] = useState(false);
    const [isHovering, setIsHovering] = useState(false);
    const videoRef = useRef(null);
    const progressRef = useRef(null);
    // Format numbers for display (e.g., 1.2M instead of 1200000)
    const formatCount = (count) => {
        if (count >= 1000000) {
            return (count / 1000000).toFixed(1) + 'M';
        } else if (count >= 1000) {
            return (count / 1000).toFixed(1) + 'K';
        } else {
            return count;
        }
    };
    // Format time for display (e.g. 1:45 instead of 105 seconds)
    const formatTime = (timeInSeconds) => {
        const minutes = Math.floor(timeInSeconds / 60);
        const seconds = Math.floor(timeInSeconds % 60);
        return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
    };
    // Handle video playback
    const togglePlay = () => {
        if (videoRef.current) {
            if (isPlaying) {
                videoRef.current.pause();
            } else {
                setIsLoading(true);
                videoRef.current.play()
                    .then(() => {
                        setIsLoading(false);
                    })
                    .catch((error) => {
                        setIsLoading(false);
                    });
            }
            setIsPlaying(!isPlaying);
        }
    };
    // Handle video loaded
    const handleVideoLoaded = () => {
        setIsLoading(false);
    };
    // Handle video mute/unmute
    const toggleMute = () => {
        if (videoRef.current) {
            videoRef.current.muted = !isMuted;
            setIsMuted(!isMuted);
        }
    };
    // Handle video progress
    const handleTimeUpdate = () => {
        if (videoRef.current) {
            const currentVideoTime = videoRef.current.currentTime;
            const duration = videoRef.current.duration;
            const progressPercent = (currentVideoTime / duration) * 100;
            setProgress(progressPercent);
            setCurrentTime(formatTime(currentVideoTime));
        }
    };
    // Handle click on progress bar
    const handleProgressClick = (e) => {
        if (progressRef.current && videoRef.current) {
            const rect = progressRef.current.getBoundingClientRect();
            const pos = (e.clientX - rect.left) / rect.width;
            videoRef.current.currentTime = pos * videoRef.current.duration;
        }
    };
    // Handle like action
    const handleLikeClick = () => {
        if (liked) {
            setLikesCount(likesCount - 1);
            setLiked(false);
        } else {
            if (disliked) {
                setDisliked(false);
            }
            setLikesCount(likesCount + 1);
            setLiked(true);
        }
    };
    // Handle dislike action
    const handleDislikeClick = () => {
        if (disliked) {
            setDisliked(false);
        } else {
            if (liked) {
                setLiked(false);
                setLikesCount(likesCount - 1);
            }
            setDisliked(true);
        }
    };
    // Toggle save/bookmark
    const handleSave = () => {
        setSaved(!saved);
    };
    // Toggle comments view
    const handleCommentClick = () => {
        setShowComments(!showComments);
    };
    // Toggle description
    const toggleDescription = () => {
        setShowDescription(!showDescription);
    };
    // Handle mouse enter/leave on video container
    const handleMouseEnter = () => {
        setIsHovering(true);
    };
    const handleMouseLeave = () => {
        setIsHovering(false);
    };
    // Render hashtags
    const renderHashtags = (hashtags) => {
        return hashtags.map((tag, index) => (
            <span key={index} className={styles.hashtag}>
                #{tag}
            </span>
        ));
    };
    const handleFavorite = useCallback(() => {
        setFavorite(prev => !prev);
    }, []);
    // Close modals when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (showActions && !event.target.closest(`.${styles.actionsDropdown}`)) {
                setShowActions(false);
            }
        };
        document.addEventListener('click', handleClickOutside);
        return () => {
            document.removeEventListener('click', handleClickOutside);
        };
    }, [showActions]);
    return (
        <div className={styles.videoCard}>
            <UserHeader/>
            {/* Feature label */}
            {/* Video container */}
            <div
                className={styles.videoContainer}
                onMouseEnter={handleMouseEnter}
                onMouseLeave={handleMouseLeave}
            >
                <div className={styles.featureLabel}>
                    <BarChart4 size={14} className={styles.featureIcon}/>
                    <span>Featured Content</span>
                </div>
                {/* Video thumbnail (shown when not playing) */}
                {!isPlaying && (
                    <div className={styles.thumbnailContainer}>
                        <img
                            src={videoData.thumbnail}
                            alt={videoData.title}
                            className={styles.thumbnail}
                        />
                        <div className={styles.thumbnailOverlay}>
                            <div className={styles.playButtonContainer}>
                                <button
                                    onClick={togglePlay}
                                    className={styles.playButton}
                                    aria-label="Play video"
                                >
                                    <Play size={24} className={styles.playIcon}/>
                                </button>
                            </div>
                            <div className={styles.videoDuration}>
                                {videoData.duration}
                            </div>
                        </div>
                    </div>
                )}
                {/* Video element */}
                <video
                    ref={videoRef}
                    src={videoData.videoSrc}
                    className={styles.videoElement}
                    poster={videoData.thumbnail}
                    muted={isMuted}
                    playsInline
                    onTimeUpdate={handleTimeUpdate}
                    onEnded={() => setIsPlaying(false)}
                    onLoadedData={handleVideoLoaded}
                />
                {/* Loading indicator */}
                {isLoading && (
                    <div className={styles.loadingOverlay}>
                        <Loader size={32} className={styles.loadingIcon}/>
                    </div>
                )}
                {/* Video controls overlay (visible when playing or hovering) */}
                {(isPlaying || isHovering) && (
                    <div className={styles.videoControlsOverlay}>
                        {/* Top controls */}
                        <div className={styles.topControls}>
                            <button
                                onClick={toggleMute}
                                className={styles.controlButton}
                                aria-label={isMuted ? "Unmute" : "Mute"}
                            >
                                {isMuted ? <VolumeX size={18} className={styles.controlIcon}/> :
                                    <Volume2 size={18} className={styles.controlIcon}/>}
                            </button>
                        </div>
                        {/* Center play/pause button */}
                        <div className={styles.centerControls}>
                            <button
                                onClick={togglePlay}
                                className={styles.centerPlayButton}
                                aria-label={isPlaying ? "Pause" : "Play"}
                            >
                                {isPlaying ?
                                    <Pause size={24} className={styles.controlIcon}/> :
                                    <Play size={24} className={styles.controlIcon}/>
                                }
                            </button>
                        </div>
                        {/* Bottom controls */}
                        <div className={styles.bottomControls}>
                            {/* Current time */}
                            <div className={styles.timeDisplay}>
                                {currentTime}
                            </div>
                            {/* Progress bar */}
                            <div
                                className={styles.progressContainer}
                                onClick={handleProgressClick}
                                ref={progressRef}
                            >
                                <div className={styles.progressTrack}>
                                    <div
                                        className={styles.progressFill}
                                        style={{width: `${progress}%`}}
                                    ></div>
                                    <div
                                        className={styles.progressHandle}
                                        style={{left: `${progress}%`}}
                                    ></div>
                                </div>
                            </div>
                            {/* Duration */}
                            <div className={styles.timeDisplay}>
                                {videoData.duration}
                            </div>
                        </div>
                    </div>
                )}
            </div>
            {/* Video metadata */}
            <div className={styles.metadataContainer}>
                {/* Creator and title section */}
                <div className={styles.infoSection}>
                    {/* Title and metrics */}
                    <div className={styles.contentSection}>
                        <h3 className={styles.videoTitle}>{videoData.title}</h3>
                        <div className={styles.videoMetrics}>
                            <div className={styles.metricItem}>
                                <Eye size={14} className={styles.metricIcon}/>
                                <span>{formatCount(videoData.metrics.views)}</span>
                            </div>
                            <div className={styles.metricItem}>
                                <Clock size={14} className={styles.metricIcon}/>
                                <span>{videoData.postedAgo}</span>
                            </div>
                        </div>
                    </div>
                </div>
                {/* Description */}
                <div className={styles.descriptionContainer}>
                    <p className={`${styles.videoDescription} ${showDescription ? styles.expanded : ''}`}>
                        {videoData.description}
                    </p>
                    {/* Read more/less button */}
                    {videoData.description.length > 100 && (
                        <button
                            className={styles.readMoreButton}
                            onClick={toggleDescription}
                        >
                            <span>{showDescription ? 'Read less' : 'Read more'}</span>
                            <ArrowUpRight size={14} className={styles.readMoreIcon}/>
                        </button>
                    )}
                    {/* Hashtags */}
                    <div className={styles.hashtagsContainer}>
                        {renderHashtags(videoData.hashtags)}
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
            </div>
            {/* Comments section */}
            {showComments && (
                <div className={styles.commentsSection}>
                    <div className={styles.commentsHeader}>
                        <h4 className={styles.commentsTitle}>Comments ({formatCount(videoData.metrics.comments)})</h4>
                        <button className={styles.viewAllButton}>View all</button>
                    </div>
                    {/* Sample comment */}
                    <div className={styles.commentContainer}>
                        <img
                            src="/images/psyche.jpg"
                            alt="Commenter"
                            className={styles.commenterAvatar}
                        />
                        <div className={styles.commentContent}>
                            <div className={styles.commentHeader}>
                                <div className={styles.commenterInfo}>
                                    <span className={styles.commenterName}>DesignWizard</span>
                                    <span className={styles.commenterUsername}>@designwiz</span>
                                </div>
                                <span className={styles.commentTime}>1d</span>
                            </div>
                            <p className={styles.commentText}>
                                These conversion tips are gold! I implemented the contrast improvements you mentioned
                                and saw a 15% boost instantly!
                            </p>
                            <div className={styles.commentActions}>
                                <button className={styles.commentLikeButton}>
                                    <Heart size={14} className={styles.commentActionIcon}/>
                                    <span>28</span>
                                </button>
                                <button className={styles.commentReplyButton}>
                                    <MessageCircle size={14} className={styles.commentActionIcon}/>
                                    <span>Reply</span>
                                </button>
                            </div>
                        </div>
                    </div>
                    {/* Comment form */}
                    <div className={styles.commentForm}>
                        <img
                            src="/images/psyche.jpg"
                            alt="Your avatar"
                            className={styles.commentFormAvatar}
                        />
                        <div className={styles.commentInputContainer}>
                            <input
                                type="text"
                                placeholder="Add a comment..."
                                className={styles.commentInput}
                            />
                            <button className={styles.postButton}>Post</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
export default ImprovedVideoCard;
"use client"
import React, { useState, useRef, useEffect } from 'react';
import {
    Heart,
    MessageCircle,
    Share2,
    Bookmark,
    Repeat2,
    Music,
    Check,
    UserPlus,
    ChevronDown,
    Volume2,
    VolumeX,
    ArrowUp,
    ArrowDown
} from '@/icons';
import styles from './page.module.css';
// Single Short Video component
const ShortVideo = ({ video, isActive, onVideoEnded }) => {
    const [isPlaying, setIsPlaying] = useState(false);
    const [isMuted, setIsMuted] = useState(true);
    const [isLiked, setIsLiked] = useState(false);
    const [isSaved, setIsSaved] = useState(false);
    const [showComments, setShowComments] = useState(false);
    const [isFollowing, setIsFollowing] = useState(false);
    const videoRef = useRef(null);
    // Auto-play when video becomes active
    useEffect(() => {
        if (isActive && videoRef.current) {
            // Short delay to ensure smooth transition
            const playTimer = setTimeout(() => {
                videoRef.current.play()
                    .then(() => {
                        setIsPlaying(true);
                    })
                    .catch(err => {
                    });
            }, 200);
            return () => clearTimeout(playTimer);
        } else if (!isActive && videoRef.current) {
            videoRef.current.pause();
            setIsPlaying(false);
        }
    }, [isActive]);
    // Toggle play/pause on tap
    const togglePlay = () => {
        if (!videoRef.current) return;
        if (isPlaying) {
            videoRef.current.pause();
            setIsPlaying(false);
        } else {
            videoRef.current.play()
                .then(() => {
                    setIsPlaying(true);
                })
                .catch(err => {
                });
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
    // Toggle comments
    const handleComments = (e) => {
        e.stopPropagation();
        setShowComments(!showComments);
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
        <div className={styles.shortVideoContainer}>
            {/* Main video container with tap to play/pause */}
            <div className={styles.videoWrapper} onClick={togglePlay}>
                <video
                    ref={videoRef}
                    src={video.videoSrc}
                    className={styles.videoElement}
                    loop
                    muted={isMuted}
                    playsInline
                    poster={video.thumbnailSrc}
                    onEnded={onVideoEnded}
                />
                {/* Mute/unmute button */}
                <button
                    className={styles.muteButton}
                    onClick={toggleMute}
                >
                    {isMuted ? <VolumeX size={20} /> : <Volume2 size={20} />}
                </button>
                {/* Overlay elements for playing state */}
                {!isPlaying && (
                    <div className={styles.pauseOverlay}>
                        <div className={styles.pauseIcon}></div>
                    </div>
                )}
                {/* Caption/description overlay */}
                <div className={styles.captionOverlay}>
                    <div className={styles.caption}>
                        <h4 className={styles.videoTitle}>{video.caption}</h4>
                        <p className={styles.videoDescription}>{video.description}</p>
                        {/* Sound/music attribution */}
                        <div className={styles.soundAttribution}>
                            <Music size={14} />
                            <span>{video.soundName} • {video.creator.name}</span>
                        </div>
                    </div>
                </div>
            </div>
            {/* User info section */}
            <div className={styles.userInfoSection}>
                <div className={styles.userInfoContainer}>
                    <div className={styles.avatarContainer}>
                        <img
                            src={video.creator.avatar}
                            alt={video.creator.name}
                            className={styles.avatar}
                        />
                        {video.creator.verified && (
                            <div className={styles.verifiedBadge}>
                                <Check size={8} />
                            </div>
                        )}
                    </div>
                    <div className={styles.userInfo}>
                        <div className={styles.userName}>{video.creator.name}</div>
                        <div className={styles.userHandle}>{video.creator.username}</div>
                    </div>
                    {!isFollowing ? (
                        <button
                            className={styles.followButton}
                            onClick={() => setIsFollowing(true)}
                        >
                            <UserPlus size={16} />
                            <span>Follow</span>
                        </button>
                    ) : (
                        <div className={styles.followingLabel}>Following</div>
                    )}
                </div>
            </div>
            {/* Side action buttons */}
            <div className={styles.sideActionButtons}>
                <div className={styles.actionButton}>
                    <button
                        className={`${styles.iconButton} ${isLiked ? styles.likedButton : ''}`}
                        onClick={handleLike}
                    >
                        <Heart size={28} fill={isLiked ? "#ef4444" : "none"} />
                    </button>
                    <span className={styles.actionCount}>{formatCount(video.likes)}</span>
                </div>
                <div className={styles.actionButton}>
                    <button
                        className={styles.iconButton}
                        onClick={handleComments}
                    >
                        <MessageCircle size={28} />
                    </button>
                    <span className={styles.actionCount}>{formatCount(video.comments)}</span>
                </div>
                <div className={styles.actionButton}>
                    <button className={styles.iconButton}>
                        <Repeat2 size={28} />
                    </button>
                    <span className={styles.actionCount}>{formatCount(video.reposts)}</span>
                </div>
                <div className={styles.actionButton}>
                    <button className={styles.iconButton}>
                        <Share2 size={28} />
                    </button>
                    <span className={styles.actionCount}>{formatCount(video.shares)}</span>
                </div>
                <div className={styles.actionButton}>
                    <button
                        className={`${styles.iconButton} ${isSaved ? styles.savedButton : ''}`}
                        onClick={handleSave}
                    >
                        <Bookmark size={28} fill={isSaved ? "#6366f1" : "none"} />
                    </button>
                    <span className={styles.actionCount}>{formatCount(video.saves)}</span>
                </div>
            </div>
            {/* Comments section (conditionally shown) */}
            {showComments && (
                <div className={styles.commentsContainer}>
                    <div className={styles.commentsHeader}>
                        <h4 className={styles.commentsTitle}>
                            {formatCount(video.comments)} comments
                        </h4>
                        <button
                            className={styles.closeButton}
                            onClick={handleComments}
                        >
                            <ChevronDown size={20} />
                        </button>
                    </div>
                    <div className={styles.commentsList}>
                        {video.topComments.map((comment, index) => (
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
                                            <Heart size={14} />
                                            <span>{formatCount(comment.likes)}</span>
                                        </button>
                                        <button className={styles.commentReplyButton}>Reply</button>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                    <div className={styles.commentInput}>
                        <img
                            src="/images/psyche.jpg"
                            alt="Your avatar"
                            className={styles.commentInputAvatar}
                        />
                        <input
                            type="text"
                            placeholder="Add a comment..."
                            className={styles.commentInputField}
                        />
                        <button className={styles.commentPostButton}>Post</button>
                    </div>
                </div>
            )}
        </div>
    );
};
// Main component for the Short Video Feed
const ShortVideoFeed = () => {
    const [activeVideoIndex, setActiveVideoIndex] = useState(0);
    const feedRef = useRef(null);
    // Sample data for short videos
    const videos = [
        {
            id: "short1",
            caption: "UI tricks that will boost your conversion rate",
            description: "These simple design changes increased our client's e-commerce conversion by 300%! #UXDesign #ConversionRate",
            videoSrc: "/api/placeholder/400/800",
            thumbnailSrc: "/api/placeholder/400/800",
            soundName: "Original Sound",
            duration: "00:29",
            likes: 145200,
            comments: 2341,
            shares: 3890,
            reposts: 726,
            saves: 18400,
            creator: {
                id: "user123",
                name: "UX Master",
                username: "@ux_master",
                avatar: "/api/placeholder/48/48",
                verified: true,
            },
            topComments: [
                {
                    userName: "Designer101",
                    userAvatar: "/images/psyche.jpg",
                    text: "That contrast tip is game-changing! Just implemented it on our site.",
                    timeAgo: "2h",
                    likes: 423
                },
                {
                    userName: "CodeCrafter",
                    userAvatar: "/images/psyche.jpg",
                    text: "Do you have a tutorial for the checkout flow optimization?",
                    timeAgo: "4h",
                    likes: 215
                }
            ]
        },
        {
            id: "short2",
            caption: "Mobile navigation patterns that users love",
            description: "After testing with 200+ users, these navigation patterns consistently performed best! #MobileUX #UXResearch",
            videoSrc: "/api/placeholder/400/800",
            thumbnailSrc: "/api/placeholder/400/800",
            soundName: "UX Tips Background",
            duration: "00:45",
            likes: 87300,
            comments: 1542,
            shares: 2150,
            reposts: 543,
            saves: 9270,
            creator: {
                id: "user456",
                name: "MobileUX Pro",
                username: "@mobile_ux",
                avatar: "/api/placeholder/48/48",
                verified: false,
            },
            topComments: [
                {
                    userName: "DevExpert",
                    userAvatar: "/images/psyche.jpg",
                    text: "The bottom navigation with labels is definitely the winner in my experience too!",
                    timeAgo: "1d",
                    likes: 328
                },
                {
                    userName: "UXNewbie",
                    userAvatar: "/images/psyche.jpg",
                    text: "What about hamburger menus? Still relevant?",
                    timeAgo: "2d",
                    likes: 154
                }
            ]
        },
        {
            id: "short3",
            caption: "Color psychology in UI design explained in 60 seconds",
            description: "Learn how colors affect user decisions and how to use them strategically #UIDesign #ColorTheory",
            videoSrc: "/api/placeholder/400/800",
            thumbnailSrc: "/api/placeholder/400/800",
            soundName: "Design Beats",
            duration: "01:00",
            likes: 231000,
            comments: 3842,
            shares: 7320,
            reposts: 1241,
            saves: 25700,
            creator: {
                id: "user789",
                name: "DesignPsych",
                username: "@design_psych",
                avatar: "/api/placeholder/48/48",
                verified: true,
            },
            topComments: [
                {
                    userName: "BrandManager",
                    userAvatar: "/images/psyche.jpg",
                    text: "This explains why our blue CTA buttons outperformed the red ones by 32%!",
                    timeAgo: "5h",
                    likes: 723
                },
                {
                    userName: "CreativeDirector",
                    userAvatar: "/images/psyche.jpg",
                    text: "Always thought color psychology was overrated but you changed my mind.",
                    timeAgo: "1d",
                    likes: 512
                }
            ]
        }
    ];
    // Handle scrolling to change active video
    useEffect(() => {
        const handleScroll = () => {
            if (!feedRef.current) return;
            const containerHeight = feedRef.current.clientHeight;
            const scrollPosition = feedRef.current.scrollTop;
            // Determine which video is currently in view
            const newIndex = Math.round(scrollPosition / containerHeight);
            if (newIndex !== activeVideoIndex && newIndex >= 0 && newIndex < videos.length) {
                setActiveVideoIndex(newIndex);
            }
        };
        const feedContainer = feedRef.current;
        if (feedContainer) {
            feedContainer.addEventListener('scroll', handleScroll);
            return () => {
                feedContainer.removeEventListener('scroll', handleScroll);
            };
        }
    }, [activeVideoIndex, videos.length]);
    // Navigate to next video
    const goToNextVideo = () => {
        if (activeVideoIndex < videos.length - 1) {
            const nextIndex = activeVideoIndex + 1;
            setActiveVideoIndex(nextIndex);
            if (feedRef.current) {
                feedRef.current.scrollTo({
                    top: nextIndex * feedRef.current.clientHeight,
                    behavior: 'smooth'
                });
            }
        }
    };
    // Navigate to previous video
    const goToPrevVideo = () => {
        if (activeVideoIndex > 0) {
            const prevIndex = activeVideoIndex - 1;
            setActiveVideoIndex(prevIndex);
            if (feedRef.current) {
                feedRef.current.scrollTo({
                    top: prevIndex * feedRef.current.clientHeight,
                    behavior: 'smooth'
                });
            }
        }
    };
    return (
        <div className={styles.shortVideoFeedContainer}>
            <div className={styles.feedScrollContainer} ref={feedRef}>
                {videos.map((video, index) => (
                    <ShortVideo
                        key={video.id}
                        video={video}
                        isActive={index === activeVideoIndex}
                        onVideoEnded={goToNextVideo}
                    />
                ))}
            </div>
            {/* Navigation buttons */}
            <div className={styles.navigationButtons}>
                {activeVideoIndex > 0 && (
                    <button
                        className={styles.prevButton}
                        onClick={goToPrevVideo}
                        aria-label="Previous video"
                    >
                        <ArrowUp size={24} />
                    </button>
                )}
                {activeVideoIndex < videos.length - 1 && (
                    <button
                        className={styles.nextButton}
                        onClick={goToNextVideo}
                        aria-label="Next video"
                    >
                        <ArrowDown size={24} />
                    </button>
                )}
            </div>
        </div>
    );
};
export default ShortVideoFeed;
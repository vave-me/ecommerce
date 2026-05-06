"use client"
import React, {useState, useRef, useCallback} from 'react';
import {
    Heart,
    ChevronLeft,
    ChevronRight,
    ChevronDown
} from '@/icons';
import styles from './page.module.css';
import {Engagement} from "../Engagement";
import UserHeader from "../UserHeader";
const ImageGallery = () => {
    // Sample post data
    const postData = {
        id: "p123456",
        caption: "New UI design exploration for our e-commerce client. Swipe to see before/after. #UIDesign #BeforeAfter",
        images: [
            "/images/default-product.webp",
            "/images/lenovo.png",
            "/images/lenovo2.png",
        ],
        location: "Berlin, Germany",
        timestamp: "2025-03-22T14:35:00Z",
        timeAgo: "2h ago",
        likes: 842,
        comments: 36,
        shares: 12,
        saves: 92,
        hashtags: ["UIDesign", "BeforeAfter", "UX"],
        creator: {
            name: "Design Studio",
            username: "@designstudio",
            avatar: "/images/psyche.jpg",
            verified: true,
        },
        isLiked: false,
        isSaved: false,
        topComments: [
            {
                userName: "UXMaster",
                userAvatar: "/api/placeholder/28/28",
                text: "Love the contrast improvement on the CTAs!",
                timeAgo: "1h",
                likes: 24
            },
            {
                userName: "DigitalCreator",
                userAvatar: "/api/placeholder/28/28",
                text: "The before/after comparison is striking!",
                timeAgo: "45m",
                likes: 12
            }
        ]
    };
    // State management
    const [currentImageIndex, setCurrentImageIndex] = useState(0);
    const [isLiked, setIsLiked] = useState(postData.isLiked);
    const [isSaved, setIsSaved] = useState(postData.isSaved);
    const [likesCount, setLikesCount] = useState(postData.likes);
    const [showComments, setShowComments] = useState(false);
    const [isFollowing, setIsFollowing] = useState(false);
    const [showCaptionFull, setShowCaptionFull] = useState(false);
    const [expanded, setExpanded] = useState(false);
    const [saved, setSaved] = useState(false);
    const [favorite, setFavorite] = useState(false);
    const [liked, setLiked] = useState(false);
    const [disliked, setDisliked] = useState(false);
    const [showContact, setShowContact] = useState(false);
    // Navigate to next image
    const nextImage = (e) => {
        e.stopPropagation();
        if (currentImageIndex < postData.images.length - 1) {
            setCurrentImageIndex(currentImageIndex + 1);
        }
    };
    // Navigate to previous image
    const prevImage = (e) => {
        e.stopPropagation();
        if (currentImageIndex > 0) {
            setCurrentImageIndex(currentImageIndex - 1);
        }
    };
    // Toggle like
    const handleLike = () => {
        if (isLiked) {
            setLikesCount(likesCount - 1);
        } else {
            setLikesCount(likesCount + 1);
        }
        setIsLiked(!isLiked);
    };
    const handleCommentClick = useCallback(() => {
        setShowComments(prev => !prev);
    }, []);
    const handleFavorite = useCallback(() => {
        setFavorite(prev => !prev);
    }, []);
    const handleDislike = useCallback(() => {
        if (liked) setLiked(false);
        setDisliked(prev => !prev);
    }, [liked]);
    // Toggle save
    const handleSave = () => {
        setIsSaved(!isSaved);
    };
    // Toggle comments
    const handleComments = () => {
        setShowComments(!showComments);
    };
    // Format number for display (e.g., 1.2K)
    const formatCount = (count) => {
        if (count >= 1000000) {
            return (count / 1000000).toFixed(1) + 'M';
        } else if (count >= 1000) {
            return (count / 1000).toFixed(1) + 'K';
        }
        return count;
    };
    // Format caption with hashtags
    const renderCaption = () => {
        const caption = postData.caption;
        const words = caption.split(' ');
        return words.map((word, index) => {
            if (word.startsWith('#')) {
                return (
                    <span key={index} className={styles.hashtag}>
                        {word}{' '}
                    </span>
                );
            }
            return word + ' ';
        });
    };
    return (
        <div className={styles.galleryCard}>
            {/* Header */}
            <UserHeader/>
            {/* Image gallery */}
            <div className={styles.imageContainer}>
                <img
                    src={postData.images[currentImageIndex]}
                    alt={`Gallery image ${currentImageIndex + 1}`}
                    className={styles.galleryImage}
                />
                {/* Image navigation */}
                {postData.images.length > 1 && (
                    <>
                        {currentImageIndex > 0 && (
                            <button
                                className={`${styles.navButton} ${styles.prevButton}`}
                                onClick={prevImage}
                            >
                                <ChevronLeft size={18}/>
                            </button>
                        )}
                        {currentImageIndex < postData.images.length - 1 && (
                            <button
                                className={`${styles.navButton} ${styles.nextButton}`}
                                onClick={nextImage}
                            >
                                <ChevronRight size={18}/>
                            </button>
                        )}
                        {/* Pagination dots */}
                        <div className={styles.paginationDots}>
                            {postData.images.map((_, index) => (
                                <div
                                    key={index}
                                    className={`${styles.paginationDot} ${
                                        index === currentImageIndex ? styles.activeDot : ''
                                    }`}
                                />
                            ))}
                        </div>
                    </>
                )}
            </div>
            {/* Post content */}
            <div className={styles.postContent}>
                {likesCount > 0 && (
                    <div className={styles.likesCount}>
                        {formatCount(likesCount)} likes
                    </div>
                )}
                <div className={styles.captionContainer}>
                    <span className={styles.captionUsername}>{postData.creator.username}</span>
                    <span className={styles.caption}>
                        {showCaptionFull ? (
                            renderCaption()
                        ) : (
                            <>
                                {postData.caption.length > 100 ? (
                                    <>
                                        {renderCaption().slice(0, 100)}...
                                        <button
                                            className={styles.moreLink}
                                            onClick={() => setShowCaptionFull(true)}
                                        >
                                            more
                                        </button>
                                    </>
                                ) : (
                                    renderCaption()
                                )}
                            </>
                        )}
                    </span>
                </div>
                <Engagement
                    favorite={favorite}
                    onFavorite={handleFavorite}
                    liked={liked}
                    onLike={handleLike}
                    disliked={disliked}
                    onDislike={handleDislike}
                    data={postData}
                    onCommentClick={handleCommentClick}
                />
                {postData.comments > 0 && (
                    <button
                        className={styles.viewCommentsButton}
                        onClick={handleComments}
                    >
                        View all {formatCount(postData.comments)} comments
                    </button>
                )}
            </div>
            {/* Comments drawer (conditionally shown) */}
            {showComments && (
                <div className={styles.commentsDrawer}>
                    <div className={styles.commentsHeader}>
                        <h4 className={styles.commentsTitle}>
                            Comments ({formatCount(postData.comments)})
                        </h4>
                        <button
                            className={styles.closeButton}
                            onClick={handleComments}
                        >
                            <ChevronDown size={16}/>
                        </button>
                    </div>
                    <div className={styles.commentsList}>
                        {postData.topComments.map((comment, index) => (
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
                                            <Heart size={12}/>
                                            <span>{formatCount(comment.likes)}</span>
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
            )}
        </div>
    );
};
export default ImageGallery;
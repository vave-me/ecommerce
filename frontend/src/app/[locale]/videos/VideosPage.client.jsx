"use client";

import React, {useCallback, useEffect, useMemo, useRef, useState, memo} from "react"; // Added useMemo
import {useDispatch} from "react-redux";
import {useTranslations} from "next-intl"; // ⬅️ FIX: Import from next-intl
// Hooks and Context (Ensure paths are correct)
import useActivityApi from "../../../hooks/useActivityApi"; // Assumed path
import useWishlist from "../../../hooks/useWishlist"; // Updated to use useWishlist
import {useAuth} from "../../../context/AuthContext"; // Assumed path
// API (Ensure paths are correct)
import {getBaseUserById} from "../../../api/client/userApi";
// Redux Actions (Ensure paths are correct)
import {openCommentsFullModal, openMessageModal,} from "../../../redux/slices/modalsSlice"; // Assumed path
// Styles
import styles from "./VideoPage.module.css";

// Child Components (Assumed to handle own translations)
import MinNav from "../../../components/Header/MinNav"; // Assumed path
import {Engagement} from "../design/Engagement"; // Assumed path

/**
 * Main client component for Videos Page with next-intl Translations
 */
const VideosPageClient = memo(function VideosPageClient({
                                             serverVideos = [], // Videos passed from Server Component
                                             errorMessage = "", // Error message passed from Server Component
                                         }) {
    // ⬅️ FIX: Use next-intl hook with namespace
    const t = useTranslations('VideosPage');

    // 1) Error or empty states (using props passed from server)
    if (errorMessage) {
        return (
            <div className={styles.errorContainer} role="alert">
                {/* ⬇️ Use relative key */}
                <h2>{t("errorTitle")}</h2>
                <p>{errorMessage}</p> {/* Display specific server error message */}
            </div>
        );
    }

    if (!serverVideos || serverVideos.length === 0) {
        // ⬇️ Use relative key
        return <div className={styles.noVideos}>{t("noVideos")}</div>;
    }

    // NOTE: This component now primarily sets up the layout.
    // The actual video data fetching and state logic (like playing, liking etc.)
    // are encapsulated within the VideoCardDesktop and VideoCardMobile components.
    return (
        <div className={styles.container}>
            {/* Mobile pinned top bar - MinNav assumed to handle its own translations */}
            {/* Pass current path for active state highlighting */}
            <div className={styles.mobileHeader}>
                <div className={styles.mobileHeaderContainer}>
                    <MinNav locationPath={"/videos"}/>
                </div>
            </div>

            {/* Desktop layout */}
            <div className={styles.desktopContainer}>
                <main className={styles.mainContent}>
                    <div className={styles.layoutGrid}>
                        {serverVideos.map((vid) => (
                            <VideoCardDesktop key={vid.id} video={vid}/>
                        ))}
                    </div>
                </main>
            </div>

            {/* Mobile layout */}
            <div className={styles.mobileContainer}>
                <div className={styles.snapScrollContainer}>
                    {serverVideos.map((vid) => (
                        <div key={vid.id} className={styles.snapItem}>
                            <VideoCardMobile video={vid}/>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
});

/* ---------------------------------------------------------------------- */
/* Desktop card                                                           */

/* ---------------------------------------------------------------------- */
const VideoCardDesktop = memo(function VideoCardDesktop({video}) {
    // ⬅️ FIX: Use next-intl hook with namespace
    const t = useTranslations('VideosPage');
    const [baseUser, setBaseUser] = useState(null);
    const [isMuted, setIsMuted] = useState(true); // Start muted

    const dispatch = useDispatch();
    const {user} = useAuth(); // Get current user from auth context
    const videoRef = useRef(null);

    // Hooks for interactions (ensure they are compatible with current context)
    const {handleLike, handleDislike} = useActivityApi();
    const {handleToggleWishlist} = useWishlist();

    /* Fetch author -------------------------------------------------------- */
    useEffect(() => {
        if (!video.userId) return;
        let isMounted = true;
        getBaseUserById(video.userId)
            .then((userData) => {
                // userData might be { user: {...} } or just { ... }
                if (isMounted) setBaseUser(userData?.user || userData);
            })
            .catch((err) => {
                // Error fetching user logged for debugging
            });
        return () => {
            isMounted = false;
        }; // Cleanup
    }, [video.userId]);

    /* Hover auto-play ------------------------------------------------------ */
    const handleMouseEnter = () => {
        // Play only if video exists and is ready
        videoRef.current?.play().catch(e => {
            // Play error handled silently
        });
    };
    const handleMouseLeave = () => {
        if (videoRef.current) {
            videoRef.current.pause();
            videoRef.current.currentTime = 0; // Reset to beginning
        }
    };

    /* Toggles & interactions ---------------------------------------------- */
    const toggleMute = useCallback(() => setIsMuted((prev) => !prev), []);

    // Wrap API calls in useCallback and check for user existence
    const handleLikeClick = useCallback(() => {
        if (!user?.userId) return; // Or show login prompt
        handleLike(video.id, user.userId).catch((err) => {
            // Like failed error logged
        });
    }, [video.id, user?.userId, handleLike]);

    const handleDislikeClick = useCallback(() => {
        if (!user?.userId) return;
        handleDislike(video.id, user.userId).catch((err) => {
            // Dislike failed error logged
        });
    }, [video.id, user?.userId, handleDislike]);

    const handleWishlistClick = useCallback(() => {
        if (!user?.userId) return;
        handleToggleWishlist(video.id).catch((err) => {
            // Wishlist toggle failed error logged
        });
    }, [video.id, user?.userId, handleToggleWishlist]);

    const toggleCommentsList = useCallback(() => {
        dispatch(
            openCommentsFullModal({itemId: video.id, itemType: "video", categoryId: "comments"})
        );
    }, [video.id, dispatch]);

    const handleOpenMessage = useCallback(() => {
        if (!user?.userId || user.userId === video.userId) return; // Don't message self or if not logged in
        dispatch(
            openMessageModal({itemId: video.id, recipientId: video.userId})
        );
    }, [video.id, video.userId, user?.userId, dispatch]);

    /* Dummy counts / wishlist flag (replace with real data ↙) */
    // These should ideally come from `video` prop or be fetched
    const counts = useMemo(() => ({
        wishlist: video.wishlistCount || 0,
        like: video.likeCount || 0,
        dislike: video.dislikeCount || 0,
        comments: video.commentCount || 0,
        message: 1 /* How to get message count? */
    }), [video]);
    const isInWishlist = video.isInWishlist || false; // Example: Assuming prop exists

    /* -------------------------------------------------------------------- */
    return (
        <div
            className={styles.cardWrapper}
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
        >
            <video
                ref={videoRef}
                src={video.url}
                poster={video.thumbnail || ""}
                loop
                muted={isMuted}
                playsInline // Important for mobile autoplay sometimes
                className={styles.videoElement}
                preload="metadata" // Load metadata for duration/poster
            />

            {/* Overlays --------------------------------------------------------- */}
            <div className={styles.overlayContainer}>
                <div className={styles.metadataOverlay}>
                    {/* Mute/Unmute Button */}
                    <button
                        type="button"
                        className={styles.muteButton}
                        onClick={toggleMute}
                        // ⬇️ Use relative keys
                        aria-label={isMuted ? t("unmute") : t("mute")}
                    >
                        {/* Display translated text or icon */}
                        {isMuted ? t("unmute") : t("mute")}
                        {/* Or use icons: {isMuted ? <VolumeOffIcon /> : <VolumeUpIcon />} */}
                    </button>
                </div>

                {/* Engagement component assumed to handle its own translations */}
                <Engagement
                    itemId={video.id}
                    handleLike={handleLikeClick}
                    handleDislike={handleDislikeClick}
                    handleWishlistClick={handleWishlistClick}
                    toggleCommentsList={toggleCommentsList}
                    toggleMessageInput={handleOpenMessage}
                    counts={counts}
                    isInWishlist={isInWishlist}
                    // Pass other relevant props like current user's reaction state
                />
            </div>

            {/* Footer ----------------------------------------------------------- */}
            <div className={styles.footerOverlay}>
                <h3 className={styles.videoTitle}>
                    {/* ⬇️ Use relative key with interpolation */}
                    {video.metadata || t("untitledVideo", {id: video.id})}
                </h3>
                <span className={styles.creatorName}>
                    {/* ⬇️ Use relative key */}
                    {baseUser?.username || t("unknownCreator")}
                </span>
            </div>
        </div>
    );
});

/* ---------------------------------------------------------------------- */
/* Mobile card                                                            */

/* ---------------------------------------------------------------------- */
const VideoCardMobile = memo(function VideoCardMobile({video}) {
    // ⬅️ FIX: Use next-intl hook with namespace
    const t = useTranslations('VideosPage');
    const [baseUser, setBaseUser] = useState(null);
    const [isMuted, setIsMuted] = useState(true); // Start muted on mobile too

    const dispatch = useDispatch();
    const {user} = useAuth();
    const videoRef = useRef(null);

    // Interaction Hooks (ensure correct imports/paths)
    const {handleLike, handleDislike} = useActivityApi();
    const {handleToggleWishlist} = useWishlist();

    /* Fetch author -------------------------------------------------------- */
    useEffect(() => {
        if (!video.userId) return;
        let isMounted = true;
        getBaseUserById(video.userId)
            .then((userData) => {
                if (isMounted) setBaseUser(userData?.user || userData);
            })
            .catch((err) => {
                // Error fetching user logged for debugging
            });
        return () => {
            isMounted = false;
        }; // Cleanup
    }, [video.userId]);

    /* Intersection observer for auto-play --------------------------------- */
    useEffect(() => {
        const videoElement = videoRef.current; // Capture ref value
        if (!videoElement) return; // Exit if ref not set yet

        const handleIntersection = (entries) => {
            entries.forEach((entry) => {
                if (entry.isIntersecting) {
                    // Attempt to play when intersecting
                    videoElement.play().catch(e => {
                        // Play error handled silently
                    });
                } else {
                    // Pause and reset when not intersecting
                    videoElement.pause();
                    videoElement.currentTime = 0;
                }
            });
        };

        const observer = new IntersectionObserver(handleIntersection, {
            root: null, // Use viewport as root
            rootMargin: '0px',
            threshold: 0.6, // ~60% visible to play
        });

        observer.observe(videoElement);

        return () => {
            // Check if element still exists before unobserving
            if (videoElement) {
                observer.unobserve(videoElement);
            }
        };
    }, []); // Empty dependency array: runs once on mount

    /* Toggles & interactions ---------------------------------------------- */
    // Use useCallback for functions passed down or depending on state/props
    const toggleMute = useCallback(() => setIsMuted((prev) => !prev), []);

    const handleLikeClick = useCallback(() => {
        if (!user?.userId) return;
        handleLike(video.id, user.userId).catch((err) => {
            // Like failed error logged
        });
    }, [video.id, user?.userId, handleLike]);

    const handleDislikeClick = useCallback(() => {
        if (!user?.userId) return;
        handleDislike(video.id, user.userId).catch((err) => {
            // Dislike failed error logged
        });
    }, [video.id, user?.userId, handleDislike]);

    const handleWishlistClick = useCallback(() => {
        if (!user?.userId) return;
        handleToggleWishlist(video.id).catch((err) => {
            // Wishlist toggle failed error logged
        });
    }, [video.id, user?.userId, handleToggleWishlist]);

    const toggleCommentsList = useCallback(() => {
        dispatch(
            openCommentsFullModal({itemId: video.id, itemType: "video", categoryId: "comments"})
        );
    }, [video.id, dispatch]);

    const handleOpenMessage = useCallback(() => {
        if (!user?.userId || user.userId === video.userId) return;
        dispatch(
            openMessageModal({itemId: video.id, recipientId: video.userId})
        );
    }, [video.id, video.userId, user?.userId, dispatch]);

    /* Dummy counts / wishlist flag (replace with real data ↙) */
    const counts = useMemo(() => ({
        wishlist: video.wishlistCount || 0,
        like: video.likeCount || 0,
        dislike: video.dislikeCount || 0,
        comments: video.commentCount || 0,
        message: 1
    }), [video]);
    const isInWishlist = video.isInWishlist || false;

    /* -------------------------------------------------------------------- */
    return (
        <div className={styles.mobileVideoWrapper}>
            <video
                ref={videoRef}
                src={video.url}
                poster={video.thumbnail || ""}
                loop
                muted={isMuted}
                playsInline // Crucial for mobile autoplay behavior
                className={styles.mobileVideoElement}
                preload="metadata"
                // Consider adding onClick handler to toggle play/pause on mobile tap?
            />

            {/* Overlays --------------------------------------------------------- */}
            <div className={styles.mobileOverlayContainer}>
                <div className={styles.mobileMetadataOverlay}>
                    <button
                        type="button"
                        className={styles.muteButton}
                        onClick={toggleMute}
                        // ⬇️ Use relative keys
                        aria-label={isMuted ? t("unmute") : t("mute")}
                    >
                        {/* Display translated text or icon */}
                        {isMuted ? t("unmute") : t("mute")}
                    </button>
                </div>
                {/* Engagement component assumed to handle its own translations */}
                <Engagement
                    itemId={video.id}
                    handleLike={handleLikeClick}
                    handleDislike={handleDislikeClick}
                    handleWishlistClick={handleWishlistClick}
                    toggleCommentsList={toggleCommentsList}
                    toggleMessageInput={handleOpenMessage}
                    counts={counts}
                    isInWishlist={isInWishlist}
                    // Pass other relevant props
                />
            </div>

            {/* Footer ----------------------------------------------------------- */}
            <div className={styles.mobileFooterOverlay}>
                <h3 className={styles.mobileVideoTitle}>
                    {/* ⬇️ Use relative key with interpolation */}
                    {video.metadata || t("untitledVideo", {id: video.id})}
                </h3>
                <span className={styles.mobileCreatorName}>
                     {/* ⬇️ Use relative key */}
                    {baseUser?.username || t("unknownCreator")}
                </span>
            </div>
        </div>
    );
});

export default VideosPageClient;

// Add PropTypes if needed
// VideosPageClient.propTypes = { ... };
// VideoCardDesktop.propTypes = { video: PropTypes.object.isRequired };
// VideoCardMobile.propTypes = { video: PropTypes.object.isRequired };
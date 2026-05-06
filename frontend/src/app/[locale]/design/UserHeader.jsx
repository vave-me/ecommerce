"use client";
import React, {useState} from 'react';
import {
    Check,
    Clock,
    MapPin,
    Globe,
    MoreHorizontal,
    UserCheck,
    UserPlus,
    Bookmark,
    Link2,
    Flag
} from '@/icons';
import styles from './UserHeader.module.css';
import Image from 'next/image';
const UserHeader = ({userData}) => {
    const [isFollowing, setIsFollowing] = useState(false);
    const [showActions, setShowActions] = useState(false);
    const [bookmarked, setBookmarked] = useState(false);
    const dummyUser = {
        id: "u789012",
        name: "Sarah Johnson",
        username: "@sarahjdesigns",
        avatar: "/images/logo-small.webp",
        verified: true,
        timeAgo: "1h ago",
        following: false,
        location: "Berlin",
        bio: "UI/UX Designer | Design Enthusiast | Coffee Lover",
        followers: 3842,
        visibility: "public"
    }
    React.useEffect(() => {
        const handleClickOutside = (event) => {
            if (showActions && !event.target.closest(`.${styles.actionsDropdown}`)) {
                setShowActions(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [showActions]);
    const handleBookmark = () => {
        setBookmarked(!bookmarked);
        setShowActions(false);
    };
    const handleCopyLink = () => {
        // Implement copy link functionality
        navigator.clipboard.writeText(window.location.href);
        setShowActions(false);
        // You might want to show a toast notification here
    };
    return (
        <div className={styles.authorHeader}>
            {/* Author avatar */}
            <div className={styles.avatarContainer}>
                <Image
                    src={dummyUser.avatar}
                    alt={dummyUser.name}
                    className={styles.avatar}
                    width={48}
                    height={48}
                    style={{ objectFit: 'cover' }}
                />
                {dummyUser.verified && (
                    <div className={styles.verifiedBadge}>
                        <Check size={12} className={styles.verifiedIcon}/>
                    </div>
                )}
            </div>
            {/* Author info */}
            <div className={styles.authorInfo}>
                <div className={styles.nameRow}>
                    <h3 className={styles.authorName}>{dummyUser.name}</h3>
                    <span className={styles.username}>{dummyUser.username}</span>
                </div>
                <div className={styles.postMeta}>
                    <div className={styles.timeLocation}>
                        <Clock size={12} className={styles.metaIcon}/>
                        <span>{dummyUser.timeAgo}</span>
                        {dummyUser.location && (
                            <>
                                <span className={styles.metaSeparator}>•</span>
                                <MapPin size={12} className={styles.metaIcon}/>
                                <span>{dummyUser.location}</span>
                            </>
                        )}
                        {dummyUser.visibility && (
                            <>
                                <span className={styles.metaSeparator}>•</span>
                                <Globe size={12} className={styles.metaIcon}/>
                                <span>{dummyUser.visibility}</span>
                            </>
                        )}
                    </div>
                </div>
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
                        aria-label="More options"
                    >
                        <MoreHorizontal size={16}/>
                    </button>
                    {/* Dropdown menu */}
                    {showActions && (
                        <div className={styles.dropdownMenu}>
                            <button
                                className={styles.dropdownItem}
                                onClick={() => setIsFollowing(!isFollowing)}
                            >
                                {isFollowing ? (
                                    <>
                                        <UserCheck size={14} className={styles.followingIcon}/>
                                        <span>Unfollow @{dummyUser.username.replace('@', '')}</span>
                                    </>
                                ) : (
                                    <>
                                        <UserPlus size={14} className={styles.dropdownItemIcon}/>
                                        <span>Follow @{dummyUser.username.replace('@', '')}</span>
                                    </>
                                )}
                            </button>
                            <button className={styles.dropdownItem} onClick={handleBookmark}>
                                <Bookmark size={14} className={styles.dropdownItemIcon}/>
                                <span>{bookmarked ? 'Remove from bookmarks' : 'Save to bookmarks'}</span>
                            </button>
                            <button className={styles.dropdownItem} onClick={handleCopyLink}>
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
    );
};
export default UserHeader;
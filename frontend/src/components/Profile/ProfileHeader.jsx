"use client";
import React, {useState, useEffect, memo} from 'react';
import PropTypes from 'prop-types';
import {
    MapPinIcon,
    CalendarIcon,
    LinkIcon,
    UserPlusIcon,
    MessagesSquare,
    EllipsisIcon
} from "@/icons";
import Rating from './Rating';
import styles from './ProfileHeader.module.css';
import { getUserById, getBaseUserById } from '../../api/client/userApi';
import { useAuth } from '../../context/AuthContext';
const ProfileHeader = memo(function ProfileHeader({username, userId, isPublicProfile = false}) {
    const [isFollowing, setIsFollowing] = useState(false);
    const [showMenu, setShowMenu] = useState(false);
    const [userProfile, setUserProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user: currentUser } = useAuth();
    const isOwnProfile = currentUser && currentUser.userId === userId;
    useEffect(() => {
        const fetchUserProfile = async () => {
            try {
                setLoading(true);
                let response;
                if (isPublicProfile) {
                    response = await getBaseUserById(userId);
                } else {
                    response = await getUserById(userId);
                }
                const user = response.user || {};
                let displayName = user.userName;
                if (!displayName) {
                    if (user.firstName) {
                        displayName = user.firstName;
                    } else {
                        displayName = 'User';
                    }
                }
                let usernameHandle = user.userName ? `@${user.userName}` : (user.firstName ? `@${user.firstName}` : '@user');
                const avatar = user.thumbnail && user.thumbnail.trim() !== '' ? user.thumbnail : "/images/psyche.jpg";
                let location = user.location;
                if (!location || ((user.lat === 0 || user.lat === '0') && (user.lng === 0 || user.lng === '0'))) {
                    location = 'Private';
                }
                setUserProfile({
                    name: displayName,
                    username: usernameHandle,
                    avatar,
                    coverPhoto: "/images/back.jpg",
                    bio: "Product enthusiast, tech lover, and avid collector. Always looking for the next great deal!",
                    location,
                    memberSince: "Jan 2022",
                    website: "https://example.com",
                    stats: {
                        followers: 324,
                        following: 156,
                        listings: 42,
                        reviews: 4.8
                    },
                    ...(isPublicProfile ? {} : {
                        email: user.email,
                        firstName: user.firstName,
                        lat: user.lat,
                        lng: user.lng
                    })
                });
            } catch (err) {
                setError(err.message || "Failed to load user profile");
                setUserProfile({
                    name: username || 'User',
                    username: `@${username || 'user'}`,
                    avatar: "/images/psyche.jpg",
                    coverPhoto: "/images/back.jpg",
                    bio: "Product enthusiast, tech lover, and avid collector. Always looking for the next great deal!",
                    location: "Private",
                    memberSince: "Jan 2022",
                    website: "https://example.com",
                    stats: {
                        followers: 324,
                        following: 156,
                        listings: 42,
                        reviews: 4.8
                    }
                });
            } finally {
                setLoading(false);
            }
        };
        if (userId) {
            fetchUserProfile();
        }
    }, [userId, username, isPublicProfile]);
    if (loading) {
        return <div className={styles.loading}>Loading profile...</div>;
    }
    if (!userProfile) {
        return <div className={styles.error}>User profile not found</div>;
    }
    return (
        <div className={styles.headerContainer}>
            {/* Cover Photo */}
            <div className={styles.coverPhotoContainer}>
                <img
                    src={userProfile.coverPhoto}
                    alt=""
                    className={styles.coverPhoto}
                />
            </div>
            <div className={styles.profileContent}>
                {/* Avatar and Quick Actions */}
                <div className={styles.avatarSection}>
                    <div className={styles.avatarWrapper}>
                        <img
                            src={userProfile.avatar}
                            alt={`${userProfile.name}'s profile photo`}
                            className={styles.avatar}
                        />
                    </div>
                    {!isOwnProfile && (
                        <div className={styles.quickActions}>
                            <button
                                className={`${styles.followButton} ${isFollowing ? styles.following : ''}`}
                                onClick={() => setIsFollowing(!isFollowing)}
                            >
                                {isFollowing ? 'Following' : (
                                    <>
                                        <UserPlusIcon className={styles.actionIcon}/>
                                        <span>Follow</span>
                                    </>
                                )}
                            </button>
                            <button className={styles.messageButton}>
                                <MessagesSquare className={styles.actionIcon}/>
                                <span className={styles.actionLabel}>Message</span>
                            </button>
                            <button
                                className={styles.moreButton}
                                onClick={() => setShowMenu(!showMenu)}
                                aria-expanded={showMenu}
                            >
                                <EllipsisIcon className={styles.moreIcon}/>
                                {showMenu && (
                                    <div className={styles.dropdown}>
                                        <button className={styles.dropdownItem}>Share Profile</button>
                                        <button className={styles.dropdownItem}>Block User</button>
                                        <button className={styles.dropdownItem}>Report User</button>
                                    </div>
                                )}
                            </button>
                        </div>
                    )}
                </div>
                {/* User Info */}
                <div className={styles.userInfo}>
                    <div className={styles.nameSection}>
                        <h1 className={styles.name}>{userProfile.name}</h1>
                        <p className={styles.username}>{userProfile.username}</p>
                    </div>
                    <p className={styles.bio}>{userProfile.bio}</p>
                    <div className={styles.metaInfo}>
                        {userProfile.location && (
                            <div className={styles.metaItem}>
                                <MapPinIcon className={styles.metaIcon}/>
                                <span>{userProfile.location}</span>
                            </div>
                        )}
                        {userProfile.memberSince && (
                            <div className={styles.metaItem}>
                                <CalendarIcon className={styles.metaIcon}/>
                                <span>Joined {userProfile.memberSince}</span>
                            </div>
                        )}
                        {userProfile.website && (
                            <div className={styles.metaItem}>
                                <LinkIcon className={styles.metaIcon}/>
                                <a href={userProfile.website} target="_blank" rel="noopener noreferrer"
                                   className={styles.websiteLink}>
                                    {userProfile.website.replace(/^https?:\/\/(www\.)?/, '')}
                                </a>
                            </div>
                        )}
                    </div>
                    <div className={styles.statsSection}>
                        <div className={styles.statItem}>
                            <span className={styles.statValue}>{userProfile.stats.followers}</span>
                            <span className={styles.statLabel}>Followers</span>
                        </div>
                        <div className={styles.statItem}>
                            <span className={styles.statValue}>{userProfile.stats.following}</span>
                            <span className={styles.statLabel}>Following</span>
                        </div>
                        <div className={styles.statItem}>
                            <span className={styles.statValue}>{userProfile.stats.listings}</span>
                            <span className={styles.statLabel}>Listings</span>
                        </div>
                        <div className={`${styles.statItem} ${styles.ratingItem}`}>
                            <Rating value={userProfile.stats.reviews} totalReviews={null}/>
                            <span className={styles.statLabel}>Rating</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
});
ProfileHeader.propTypes = {
    username: PropTypes.string.isRequired,
    userId: PropTypes.string.isRequired,
    isPublicProfile: PropTypes.bool,
};
export default ProfileHeader;
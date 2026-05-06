"use client";
import React, {memo} from 'react';
import Link from 'next/link';
import Image from 'next/image';
import {getUsernameProfileLink} from '../../utils/profileUtils';
import styles from './UserCard.module.css';
const UserCard = memo(function UserCard({user}) {
    const {
        id,
        name,
        username,
        avatar,
        bio,
        followers,
        following,
        isVerified,
        isOnline
    } = user;
    // Generate correct profile URL
    const profileUrl = getUsernameProfileLink(username || id);
    return (
        <Link href={profileUrl} className={styles.userCard}>
            <div className={styles.avatarContainer}>
                <div className={styles.avatarWrapper}>
                    <Image
                        src={avatar || '/images/default-avatar.png'}
                        alt={name}
                        width={64}
                        height={64}
                        className={styles.avatar}
                    />
                    {isOnline && <span className={styles.onlineIndicator}/>}
                </div>
            </div>
            <div className={styles.userInfo}>
                <div className={styles.nameContainer}>
                    <h3 className={styles.name}>
                        {name}
                        {isVerified && (
                            <span className={styles.verifiedBadge} aria-label="Verified">
                                ✓
                            </span>
                        )}
                    </h3>
                    <p className={styles.username}>@{username}</p>
                </div>
                {bio && (
                    <p className={styles.bio}>{bio}</p>
                )}
                <div className={styles.stats}>
                    <span className={styles.stat}>
                        <strong>{followers}</strong> followers
                    </span>
                    <span className={styles.stat}>
                        <strong>{following}</strong> following
                    </span>
                </div>
            </div>
        </Link>
    );
}, (prevProps, nextProps) => {
    // Deep comparison for user object to prevent unnecessary re-renders
    const prevUser = prevProps.user;
    const nextUser = nextProps.user;
    if (!prevUser || !nextUser) return prevUser === nextUser;
    return prevUser.id === nextUser.id &&
        prevUser.name === nextUser.name &&
        prevUser.username === nextUser.username &&
        prevUser.avatar === nextUser.avatar &&
        prevUser.bio === nextUser.bio &&
        prevUser.isOnline === nextUser.isOnline &&
        prevUser.isVerified === nextUser.isVerified &&
        prevUser.followers === nextUser.followers &&
        prevUser.following === nextUser.following;
});
export default UserCard; 
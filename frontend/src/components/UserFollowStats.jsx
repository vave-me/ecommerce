import React, {memo} from 'react';
import {useRouter} from 'next/navigation';
import {Users, UserCheck} from '@/icons';
import useFollowing from '../hooks/useFollowing';
import styles from './UserFollowStats.module.css';
/**
 * UserFollowStats Component
 * Displays clickable follower and following counts
 * Memoized for performance optimization
 *
 * @param {string} userId - The user ID to display stats for
 * @param {boolean} showIcons - Whether to show icons next to counts
 * @param {string} size - Size variant ('small', 'medium', 'large')
 * @param {string} variant - Style variant ('default', 'compact', 'card')
 */
const UserFollowStats = memo(({
                                  userId,
                                  showIcons = true,
                                  size = 'medium',
                                  variant = 'default',
                                  className = ''
                              }) => {
    const router = useRouter();
    // Get follower stats for this user
    const {
        getFollowerCount: getFollowersCount,
        loading,
        error
    } = useFollowing(userId, 'followers');
    // Get following stats for this user
    const {
        getFollowerCount: getFollowingCount,
        loading: followingLoading
    } = useFollowing(userId, 'following');
    // Navigation handlers
    const handleFollowersClick = () => {
        if (userId) {
            router.push(`/followers/${userId}`);
        }
    };
    const handleFollowingClick = () => {
        if (userId) {
            router.push(`/following/${userId}`);
        }
    };
    // Loading state
    if (loading || followingLoading) {
        return (
            <div className={`${styles.container} ${styles[size]} ${styles[variant]} ${className}`}>
                <div className={styles.statItem}>
                    <span className={styles.loading}>-</span>
                    <span className={styles.label}>Followers</span>
                </div>
                <div className={styles.statItem}>
                    <span className={styles.loading}>-</span>
                    <span className={styles.label}>Following</span>
                </div>
            </div>
        );
    }
    // Error state
    if (error) {
        return (
            <div className={`${styles.container} ${styles[size]} ${styles[variant]} ${className}`}>
                <div className={styles.error}>Unable to load stats</div>
            </div>
        );
    }
    // Get counts
    const followersCount = getFollowersCount('approved') || 0;
    const followingCount = getFollowingCount('approved') || 0;
    return (
        <div className={`${styles.container} ${styles[size]} ${styles[variant]} ${className}`}>
            {/* Followers */}
            <button
                className={styles.statButton}
                onClick={handleFollowersClick}
                title={`View ${followersCount} followers`}
                aria-label={`${followersCount} followers`}
            >
                <div className={styles.statItem}>
                    {showIcons && <Users className={styles.icon}/>}
                    <span className={styles.count}>{followersCount}</span>
                    <span className={styles.label}>
                        {followersCount === 1 ? 'Follower' : 'Followers'}
                    </span>
                </div>
            </button>
            {/* Following */}
            <button
                className={styles.statButton}
                onClick={handleFollowingClick}
                title={`View ${followingCount} following`}
                aria-label={`${followingCount} following`}
            >
                <div className={styles.statItem}>
                    {showIcons && <UserCheck className={styles.icon}/>}
                    <span className={styles.count}>{followingCount}</span>
                    <span className={styles.label}>Following</span>
                </div>
            </button>
        </div>
    );
});
export default UserFollowStats; 
"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import {
    ArrowLeft,
    User,
    Mail,
    Calendar,
    Shield,
    Activity,
    CheckCircle,
    XCircle,
    Clock,
    Edit,
    AlertCircle,
    Building,
    MapPin,
    Phone,
    Globe
} from 'lucide-react';
import { getUserById, updateUser, enableUser, disableUser } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './UserDetail.module.css';

const UserDetail = ({ locale, userId }) => {
    const t = useTranslations('UserDetail');
    const router = useRouter();
    const queryClient = useQueryClient();
    const [showActions, setShowActions] = useState(false);

    // Fetch user details
    const { data: userData, isLoading, error } = useQuery({
        queryKey: ['adminUser', userId],
        queryFn: () => getUserById(userId),
        staleTime: 60000,
    });
    
    const user = userData?.user || userData;
    
    // Log user data for debugging
    React.useEffect(() => {
        if (user) {
            // User data logged for debugging
        }
    }, [user]);

    // Mutations
    const toggleUserMutation = useMutation({
        mutationFn: async (enabled) => {
            if (enabled) {
                return await enableUser(userId);
            } else {
                return await disableUser(userId);
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUser', userId]);
            queryClient.invalidateQueries(['adminUsers']);
        },
    });

    const updateRoleMutation = useMutation({
        mutationFn: (role) => updateUser(userId, { role }),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUser', userId]);
            queryClient.invalidateQueries(['adminUsers']);
        },
    });

    const getRoleBadgeClass = (role) => {
        const normalizedRole = (role || 'customer').toLowerCase();
        switch (normalizedRole) {
            case 'admin':
                return styles.roleBadgeAdmin;
            case 'business':
                return styles.roleBadgeBusiness;
            case 'customer':
            default:
                return styles.roleBadgeCustomer;
        }
    };

    const getStatusBadgeClass = (enabled) => {
        return enabled ? styles.statusActive : styles.statusInactive;
    };

    if (isLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner />
            </div>
        );
    }

    if (error || !user) {
        return (
            <div className={styles.container}>
                <div className={styles.errorContainer}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h2>User Not Found</h2>
                    <p>The user you're looking for doesn't exist or you don't have permission to view it.</p>
                    <button 
                        className={styles.primaryButton}
                        onClick={() => router.push(`/${locale}/admin/users`)}
                    >
                        Back to Users
                    </button>
                </div>
            </div>
        );
    }

    return (
        <ErrorBoundary>
            <div className={styles.container}>
                {/* Header */}
                <div className={styles.header}>
                    <button 
                        className={styles.backButton}
                        onClick={() => router.push(`/${locale}/admin/users`)}
                    >
                        <ArrowLeft size={20} />
                        Back to Users
                    </button>
                    <div className={styles.headerActions}>
                        <button
                            className={styles.primaryButton}
                            onClick={() => router.push(`/${locale}/admin/users/${userId}/edit`)}
                        >
                            <Edit size={16} />
                            Edit User
                        </button>
                    </div>
                </div>

                {/* User Profile Card */}
                <div className={styles.profileCard}>
                    <div className={styles.profileHeader}>
                        <div className={styles.profileAvatar}>
                            {user.thumbnail ? (
                                <img src={user.thumbnail} alt={user.userName || user.email} />
                            ) : (
                                <User size={48} />
                            )}
                        </div>
                        <div className={styles.profileInfo}>
                            <h1 className={styles.userName}>
                                {user.userName || `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.email?.split('@')[0]}
                            </h1>
                            <p className={styles.userEmail}>{user.email}</p>
                            <div className={styles.badges}>
                                <span className={`${styles.roleBadge} ${getRoleBadgeClass(user.role)}`}>
                                    <Shield size={14} />
                                    {(user.role || 'customer').charAt(0).toUpperCase() + (user.role || 'customer').slice(1)}
                                </span>
                                <span className={`${styles.statusBadge} ${getStatusBadgeClass(user.enabled)}`}>
                                    {user.enabled ? <CheckCircle size={14} /> : <XCircle size={14} />}
                                    {user.enabled ? 'Active' : 'Inactive'}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* User Details Grid */}
                <div className={styles.detailsGrid}>
                    {/* Personal Information */}
                    <div className={styles.detailCard}>
                        <h3 className={styles.cardTitle}>
                            <User size={20} />
                            Personal Information
                        </h3>
                        <div className={styles.detailsList}>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>First Name</span>
                                <span className={styles.detailValue}>{user.firstName || 'Not provided'}</span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Last Name</span>
                                <span className={styles.detailValue}>{user.lastName || 'Not provided'}</span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Username</span>
                                <span className={styles.detailValue}>{user.userName || 'Not set'}</span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>User ID</span>
                                <span className={styles.detailValue}>
                                    <code className={styles.code}>{user.id}</code>
                                </span>
                            </div>
                        </div>
                    </div>

                    {/* Contact Information */}
                    <div className={styles.detailCard}>
                        <h3 className={styles.cardTitle}>
                            <Mail size={20} />
                            Contact Information
                        </h3>
                        <div className={styles.detailsList}>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Email</span>
                                <span className={styles.detailValue}>{user.email}</span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Location</span>
                                <span className={styles.detailValue}>
                                    {user.location || 'Not provided'}
                                </span>
                            </div>
                            {(user.lat && user.lng) && (
                                <div className={styles.detailItem}>
                                    <span className={styles.detailLabel}>Coordinates</span>
                                    <span className={styles.detailValue}>
                                        {user.lat.toFixed(4)}, {user.lng.toFixed(4)}
                                    </span>
                                </div>
                            )}
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Language</span>
                                <span className={styles.detailValue}>
                                    {user.language || 'en'}
                                </span>
                            </div>
                        </div>
                    </div>

                    {/* Account Information */}
                    <div className={styles.detailCard}>
                        <h3 className={styles.cardTitle}>
                            <Calendar size={20} />
                            Account Information
                        </h3>
                        <div className={styles.detailsList}>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>User ID</span>
                                <span className={`${styles.detailValue} ${styles.code}`}>
                                    {user.id}
                                </span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Account Status</span>
                                <span className={styles.detailValue}>
                                    {user.enabled ? 
                                        <span className={styles.verified}><CheckCircle size={16} /> Active</span> : 
                                        <span className={styles.unverified}><XCircle size={16} /> Inactive</span>
                                    }
                                </span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Authentication Type</span>
                                <span className={styles.detailValue}>
                                    {user.googleId ? 'Google OAuth' : 'Email/Password'}
                                </span>
                            </div>
                            {user.googleId && (
                                <div className={styles.detailItem}>
                                    <span className={styles.detailLabel}>Google ID</span>
                                    <span className={`${styles.detailValue} ${styles.code}`}>
                                        {user.googleId}
                                    </span>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Privacy Settings */}
                    <div className={styles.detailCard}>
                        <h3 className={styles.cardTitle}>
                            <Shield size={20} />
                            Privacy & Permissions
                        </h3>
                        <div className={styles.detailsList}>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Role</span>
                                <span className={styles.detailValue}>
                                    {(user.role || 'customer').charAt(0).toUpperCase() + (user.role || 'customer').slice(1)}
                                </span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Privacy Level</span>
                                <span className={styles.detailValue}>
                                    {user.privacy || 'Standard'}
                                </span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Background Image</span>
                                <span className={styles.detailValue}>
                                    {user.background ? 'Uploaded' : 'Not set'}
                                </span>
                            </div>
                            <div className={styles.detailItem}>
                                <span className={styles.detailLabel}>Bio</span>
                                <span className={styles.detailValue}>
                                    {user.bio || 'Not provided'}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Quick Actions */}
                <div className={styles.quickActions}>
                    <h3>Quick Actions</h3>
                    <div className={styles.actionButtons}>
                        <button
                            className={`${styles.actionButton} ${user.enabled ? styles.dangerButton : styles.successButton}`}
                            onClick={() => toggleUserMutation.mutate(!user.enabled)}
                            disabled={toggleUserMutation.isLoading}
                        >
                            {user.enabled ? <XCircle size={16} /> : <CheckCircle size={16} />}
                            {user.enabled ? 'Deactivate User' : 'Activate User'}
                        </button>
                        <button
                            className={styles.actionButton}
                            onClick={() => {
                                const newRole = prompt('Enter new role (customer, business, admin):', user.role || 'customer');
                                if (newRole && ['customer', 'business', 'admin'].includes(newRole.toLowerCase())) {
                                    updateRoleMutation.mutate(newRole.toLowerCase());
                                }
                            }}
                        >
                            <Shield size={16} />
                            Change Role
                        </button>
                        <button
                            className={styles.actionButton}
                            onClick={() => router.push(`/${locale}/admin/users/${userId}/metrics`)}
                        >
                            <Activity size={16} />
                            View Metrics
                        </button>
                    </div>
                </div>
            </div>
        </ErrorBoundary>
    );
};

export default UserDetail;
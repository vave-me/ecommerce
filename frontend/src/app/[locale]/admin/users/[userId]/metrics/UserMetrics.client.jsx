"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { 
    ArrowLeft, 
    Download, 
    Calendar, 
    Activity,
    ShoppingCart,
    MessageCircle,
    Star,
    TrendingUp,
    TrendingDown,
    AlertCircle,
    BarChart3,
    PieChart,
    LineChart,
    Package,
    CreditCard,
    Users,
    Eye
} from 'lucide-react';
import { getUserById } from '@/api/adminApi';
import { getUserMetrics } from '@/api/client/metricsApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import { toast } from 'react-toastify';
import dayjs from 'dayjs';
import styles from './UserMetrics.module.css';

const MetricCard = ({ title, value, subtitle, icon: Icon, trend, trendValue }) => {
    const getTrendIcon = () => {
        if (!trend) return null;
        return trend === 'up' ? <TrendingUp size={16} /> : <TrendingDown size={16} />;
    };

    const getTrendClass = () => {
        if (!trend) return '';
        return trend === 'up' ? styles.trendUp : styles.trendDown;
    };

    return (
        <div className={styles.metricCard}>
            <div className={styles.metricHeader}>
                <div className={styles.metricIcon}>
                    <Icon size={20} />
                </div>
                {trend && (
                    <div className={`${styles.trend} ${getTrendClass()}`}>
                        {getTrendIcon()}
                        <span>{trendValue}%</span>
                    </div>
                )}
            </div>
            <div className={styles.metricContent}>
                <h3 className={styles.metricValue}>{value}</h3>
                <p className={styles.metricTitle}>{title}</p>
                {subtitle && <p className={styles.metricSubtitle}>{subtitle}</p>}
            </div>
        </div>
    );
};

const UserMetrics = ({ locale, userId }) => {
    const t = useTranslations('UserMetrics');
    const router = useRouter();
    const [selectedPeriod, setSelectedPeriod] = useState('30days');

    // Fetch user details
    const { data: user, isLoading: userLoading, error: userError } = useQuery({
        queryKey: ['adminUser', userId],
        queryFn: () => getUserById(userId),
        staleTime: 60000,
    });

    // Fetch user metrics
    const { data: metricsData, isLoading: metricsLoading, error: metricsError } = useQuery({
        queryKey: ['userMetrics', userId],
        queryFn: () => getUserMetrics(userId),
        staleTime: 60000,
        enabled: !!user,
    });

    // Process metrics data
    const metrics = React.useMemo(() => {
        if (!metricsData?.metric) return null;
        
        const m = metricsData.metric;
        return {
            overview: {
                likesCount: parseInt(m.likesCount || 0),
                dislikesCount: parseInt(m.dislikesCount || 0),
                commentsCount: parseInt(m.commentsCount || 0),
                messagesCount: parseInt(m.messagesCount || 0),
                sharedCount: parseInt(m.sharedCount || 0),
                addedToWishlistCount: parseInt(m.addedToWishlistCount || 0),
                viewsCount: parseInt(m.viewsCount || 0),
                reportsCount: parseInt(m.reportsCount || 0)
            },
            calculated: {
                engagementRate: m.viewsCount > 0 ? 
                    (((parseInt(m.likesCount) + parseInt(m.commentsCount) + parseInt(m.sharedCount)) / parseInt(m.viewsCount)) * 100).toFixed(2) : 0,
                likeRatio: (parseInt(m.likesCount) + parseInt(m.dislikesCount)) > 0 ?
                    ((parseInt(m.likesCount) / (parseInt(m.likesCount) + parseInt(m.dislikesCount))) * 100).toFixed(2) : 0
            }
        };
    }, [metricsData]);

    const handleExportMetrics = () => {
        if (!user || !metrics) return;

        const csvContent = [
            ['User Metrics Report'],
            ['Generated:', new Date().toLocaleString()],
            [''],
            ['User Information'],
            ['ID:', user.id],
            ['Name:', `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.userName],
            ['Email:', user.email],
            [''],
            ['Engagement Metrics'],
            ['Likes:', metrics.overview.likesCount],
            ['Dislikes:', metrics.overview.dislikesCount],
            ['Comments:', metrics.overview.commentsCount],
            ['Messages:', metrics.overview.messagesCount],
            ['Shares:', metrics.overview.sharedCount],
            ['Wishlist Adds:', metrics.overview.addedToWishlistCount],
            ['Views:', metrics.overview.viewsCount],
            ['Reports:', metrics.overview.reportsCount],
            [''],
            ['Calculated Metrics'],
            ['Engagement Rate:', `${metrics.calculated.engagementRate}%`],
            ['Like Ratio:', `${metrics.calculated.likeRatio}%`]
        ].map(row => row.join(',')).join('\n');

        const blob = new Blob([csvContent], { type: 'text/csv' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `user-metrics-${user.id}-${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

        toast.success('Metrics exported successfully');
    };

    if (userLoading || metricsLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner />
            </div>
        );
    }

    if (userError || !user) {
        return (
            <div className={styles.container}>
                <div className={styles.errorContainer}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h2>User Not Found</h2>
                    <p>The user metrics you're trying to view don't exist.</p>
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

    if (!metrics) {
        return null;
    }

    return (
        <ErrorBoundary>
            <div className={styles.container}>
                {/* Header */}
                <div className={styles.header}>
                    <div className={styles.headerLeft}>
                        <button 
                            className={styles.backButton}
                            onClick={() => router.push(`/${locale}/admin/users/${userId}`)}
                        >
                            <ArrowLeft size={20} />
                            Back to User
                        </button>
                        <div>
                            <h1 className={styles.title}>User Metrics</h1>
                            <p className={styles.subtitle}>
                                {user.firstName || user.lastName 
                                    ? `${user.firstName || ''} ${user.lastName || ''}`.trim()
                                    : user.userName || 'Unknown User'
                                } • {user.email}
                            </p>
                        </div>
                    </div>
                    <div className={styles.headerActions}>
                        <select 
                            className={styles.periodSelect}
                            value={selectedPeriod}
                            onChange={(e) => setSelectedPeriod(e.target.value)}
                        >
                            <option value="7days">Last 7 days</option>
                            <option value="30days">Last 30 days</option>
                            <option value="90days">Last 90 days</option>
                            <option value="1year">Last year</option>
                            <option value="all">All time</option>
                        </select>
                        <button 
                            className={styles.exportButton}
                            onClick={handleExportMetrics}
                        >
                            <Download size={16} />
                            Export Metrics
                        </button>
                    </div>
                </div>

                {/* Overview Metrics */}
                <div className={styles.section}>
                    <h2 className={styles.sectionTitle}>
                        <BarChart3 size={20} />
                        Overview
                    </h2>
                    <div className={styles.metricsGrid}>
                        <MetricCard
                            icon={Eye}
                            title="Total Views"
                            value={metrics.overview.viewsCount}
                            subtitle="All time"
                        />
                        <MetricCard
                            icon={Star}
                            title="Likes"
                            value={metrics.overview.likesCount}
                            subtitle={`${metrics.calculated.likeRatio}% positive`}
                        />
                        <MetricCard
                            icon={MessageCircle}
                            title="Comments"
                            value={metrics.overview.commentsCount}
                            subtitle="Total posted"
                        />
                        <MetricCard
                            icon={Activity}
                            title="Engagement Rate"
                            value={`${metrics.calculated.engagementRate}%`}
                            subtitle="Views to interactions"
                        />
                    </div>
                </div>

                {/* Detailed Metrics */}
                <div className={styles.section}>
                    <h2 className={styles.sectionTitle}>
                        <BarChart3 size={20} />
                        Detailed Metrics
                    </h2>
                    <div className={styles.detailsGrid}>
                        <div className={styles.detailCard}>
                            <h3>User Interactions</h3>
                            <div className={styles.statsList}>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Likes Given</span>
                                    <span className={styles.statValue}>{metrics.overview.likesCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Dislikes Given</span>
                                    <span className={styles.statValue}>{metrics.overview.dislikesCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Comments Posted</span>
                                    <span className={styles.statValue}>{metrics.overview.commentsCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Messages Sent</span>
                                    <span className={styles.statValue}>{metrics.overview.messagesCount}</span>
                                </div>
                            </div>
                        </div>
                        <div className={styles.detailCard}>
                            <h3>Content Engagement</h3>
                            <div className={styles.statsList}>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Content Shared</span>
                                    <span className={styles.statValue}>{metrics.overview.sharedCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Added to Wishlists</span>
                                    <span className={styles.statValue}>{metrics.overview.addedToWishlistCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Total Views</span>
                                    <span className={styles.statValue}>{metrics.overview.viewsCount}</span>
                                </div>
                                <div className={styles.statItem}>
                                    <span className={styles.statLabel}>Reports Received</span>
                                    <span className={styles.statValue}>{metrics.overview.reportsCount}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </ErrorBoundary>
    );
};

export default UserMetrics;
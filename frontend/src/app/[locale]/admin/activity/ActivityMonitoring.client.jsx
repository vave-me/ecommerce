"use client";

import React, { useState, useMemo, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { 
  Activity, 
  TrendingUp, 
  Users, 
  Eye, 
  Calendar,
  Clock,
  Filter,
  Search,
  RefreshCw,
  Download,
  AlertTriangle,
  CheckCircle,
  Info,
  XCircle,
  User,
  MessageSquare,
  Heart,
  Share,
  ShoppingCart,
  FileText,
  Settings,
  Globe,
  Smartphone,
  Monitor,
  Tablet,
  MapPin,
  BarChart3,
  PieChart,
  LineChart,
  Zap,
  Target,
  Layers,
  Database,
  Server,
  Wifi,
  WifiOff,
  MousePointer,
  Navigation,
  Plus,
  Minus,
  TrendingDown,
  ArrowUp,
  ArrowDown,
  Play,
  Pause,
  Square
} from 'lucide-react';
import { 
  listActivities, 
  getMostLiked, 
  getMostDisliked,
  getActivityMetrics,
  getUserActivities,
  getSystemHealth,
  getRealtimeStats,
  exportActivityLogs,
  getActivityAnalytics,
  getGeoActivity,
  getDeviceBreakdown,
  getPageViews,
  getEngagementMetrics
} from '../../../../api/adminApi';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'react-toastify';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ActivityMonitoring.module.css';

// Activity Type Icons Component
const ActivityTypeIcon = ({ type, size = 16 }) => {
  const iconMap = {
    user_login: User,
    user_logout: User,
    user_register: User,
    post_create: FileText,
    post_like: Heart,
    post_share: Share,
    post_comment: MessageSquare,
    product_view: Eye,
    product_purchase: ShoppingCart,
    page_view: Globe,
    search: Search,
    settings_change: Settings,
    error: AlertTriangle,
    warning: AlertTriangle,
    info: Info,
    success: CheckCircle
  };

  const IconComponent = iconMap[type] || Activity;
  return <IconComponent size={size} />;
};

// Activity Status Badge Component
const ActivityStatusBadge = ({ status, type }) => {
  const getStatusConfig = () => {
    switch (status || type) {
      case 'success':
      case 'completed':
        return { color: 'success', icon: CheckCircle };
      case 'error':
      case 'failed':
        return { color: 'danger', icon: XCircle };
      case 'warning':
        return { color: 'warning', icon: AlertTriangle };
      case 'info':
      case 'pending':
        return { color: 'info', icon: Info };
      default:
        return { color: 'secondary', icon: Activity };
    }
  };

  const config = getStatusConfig();
  const Icon = config.icon;

  return (
    <span className={`${styles.statusBadge} ${styles[config.color]}`}>
      <Icon size={12} />
      {status || type}
    </span>
  );
};

// Real-time Activity Feed Item Component
const ActivityFeedItem = ({ activity, onUserClick }) => {
  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
  };

  return (
    <div className={styles.activityFeedItem}>
      <div className={styles.activityIcon}>
        <ActivityTypeIcon type={activity.type} size={16} />
      </div>
      <div className={styles.activityContent}>
        <div className={styles.activityHeader}>
          <span 
            className={styles.activityUser}
            onClick={() => onUserClick(activity.user_id)}
          >
            {activity.user_name || `User ${activity.user_id}`}
          </span>
          <ActivityStatusBadge status={activity.status} type={activity.type} />
        </div>
        <p className={styles.activityDescription}>
          {activity.description || `Performed ${activity.type.replace('_', ' ')}`}
        </p>
        <div className={styles.activityMeta}>
          <span className={styles.activityTime}>
            <Clock size={12} />
            {formatTime(activity.created_at)}
          </span>
          {activity.ip_address && (
            <span className={styles.activityLocation}>
              <MapPin size={12} />
              {activity.location || activity.ip_address}
            </span>
          )}
          {activity.device && (
            <span className={styles.activityDevice}>
              {activity.device === 'mobile' && <Smartphone size={12} />}
              {activity.device === 'tablet' && <Tablet size={12} />}
              {activity.device === 'desktop' && <Monitor size={12} />}
              {activity.device}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

// Metrics Chart Component (simplified chart representation)
const MetricsChart = ({ data, type, title }) => {
  if (!data || data.length === 0) {
    return (
      <div className={styles.chartEmpty}>
        <BarChart3 size={24} />
        <span>No data available</span>
      </div>
    );
  }

  return (
    <div className={styles.metricsChart}>
      <h4 className={styles.chartTitle}>{title}</h4>
      <div className={styles.chartContainer}>
        <div className={styles.chartPlaceholder}>
          {type === 'line' && <LineChart size={48} />}
          {type === 'bar' && <BarChart3 size={48} />}
          {type === 'pie' && <PieChart size={48} />}
          <p>Chart visualization would be rendered here</p>
          <small>Data points: {data.length}</small>
        </div>
      </div>
    </div>
  );
};

// System Health Component
const SystemHealthCard = ({ health }) => {
  const getHealthStatus = () => {
    if (!health) return { status: 'unknown', color: 'secondary' };
    if (health.cpu_usage > 90 || health.memory_usage > 90) return { status: 'critical', color: 'danger' };
    if (health.cpu_usage > 70 || health.memory_usage > 70) return { status: 'warning', color: 'warning' };
    return { status: 'healthy', color: 'success' };
  };

  const status = getHealthStatus();

  return (
    <div className={`${styles.healthCard} ${styles[status.color]}`}>
      <div className={styles.healthHeader}>
        <Server size={20} />
        <span>System Health</span>
        <span className={`${styles.healthStatus} ${styles[status.color]}`}>
          {status.status.toUpperCase()}
        </span>
      </div>
      {health && (
        <div className={styles.healthMetrics}>
          <div className={styles.healthMetric}>
            <span>CPU</span>
            <div className={styles.healthBar}>
              <div 
                className={styles.healthBarFill}
                style={{ width: `${health.cpu_usage || 0}%` }}
              />
            </div>
            <span>{health.cpu_usage || 0}%</span>
          </div>
          <div className={styles.healthMetric}>
            <span>Memory</span>
            <div className={styles.healthBar}>
              <div 
                className={styles.healthBarFill}
                style={{ width: `${health.memory_usage || 0}%` }}
              />
            </div>
            <span>{health.memory_usage || 0}%</span>
          </div>
          <div className={styles.healthMetric}>
            <span>Storage</span>
            <div className={styles.healthBar}>
              <div 
                className={styles.healthBarFill}
                style={{ width: `${health.storage_usage || 0}%` }}
              />
            </div>
            <span>{health.storage_usage || 0}%</span>
          </div>
        </div>
      )}
    </div>
  );
};

const ActivityMonitoring = () => {
  const t = useTranslations('AdminActivity');
  const { user } = useAuth();
  const queryClient = useQueryClient();

  // State management
  const [timeRange, setTimeRange] = useState('24h');
  const [filterType, setFilterType] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedUser, setSelectedUser] = useState(null);
  const [realTimeEnabled, setRealTimeEnabled] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);

  const itemsPerPage = 50;

  // Real-time polling effect
  useEffect(() => {
    if (!realTimeEnabled) return;

    const interval = setInterval(() => {
      queryClient.invalidateQueries(['admin-activities']);
      queryClient.invalidateQueries(['realtime-stats']);
    }, 5000); // Poll every 5 seconds

    return () => clearInterval(interval);
  }, [realTimeEnabled, queryClient]);

  // Data fetching
  const { data: activities, isLoading: activitiesLoading } = useQuery({
    queryKey: ['admin-activities', timeRange, filterType, searchTerm, currentPage],
    queryFn: () => listActivities({ 
      timeRange, 
      type: filterType === 'all' ? undefined : filterType,
      search: searchTerm,
      page: currentPage,
      limit: itemsPerPage
    }),
    staleTime: realTimeEnabled ? 0 : 30000,
    refetchInterval: realTimeEnabled ? 5000 : false
  });

  const { data: metrics } = useQuery({
    queryKey: ['activity-metrics', timeRange],
    queryFn: () => getActivityMetrics({ timeRange }),
    staleTime: 300000 // 5 minutes
  });

  const { data: realtimeStats } = useQuery({
    queryKey: ['realtime-stats'],
    queryFn: getRealtimeStats,
    staleTime: realTimeEnabled ? 0 : 60000,
    refetchInterval: realTimeEnabled ? 2000 : false
  });

  const { data: systemHealth } = useQuery({
    queryKey: ['system-health'],
    queryFn: getSystemHealth,
    staleTime: 30000,
    refetchInterval: 30000
  });

  const { data: analytics } = useQuery({
    queryKey: ['activity-analytics', timeRange],
    queryFn: () => getActivityAnalytics({ timeRange }),
    staleTime: 600000 // 10 minutes
  });

  const { data: geoActivity } = useQuery({
    queryKey: ['geo-activity', timeRange],
    queryFn: () => getGeoActivity({ timeRange }),
    staleTime: 300000
  });

  const { data: deviceBreakdown } = useQuery({
    queryKey: ['device-breakdown', timeRange],
    queryFn: () => getDeviceBreakdown({ timeRange }),
    staleTime: 300000
  });

  const { data: mostLiked } = useQuery({
    queryKey: ['most-liked', timeRange],
    queryFn: () => getMostLiked({ timeRange }),
    staleTime: 600000
  });

  const { data: engagementMetrics } = useQuery({
    queryKey: ['engagement-metrics', timeRange],
    queryFn: () => getEngagementMetrics({ timeRange }),
    staleTime: 300000
  });

  // Mutations
  const exportLogsMutation = useMutation({
    mutationFn: exportActivityLogs,
    onSuccess: (blob) => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `activity-logs-${timeRange}-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      toast.success('Activity logs exported successfully');
    },
    onError: (error) => {
      toast.error(`Export failed: ${error.message}`);
    }
  });

  // Event handlers
  const handleExportLogs = () => {
    exportLogsMutation.mutate({ timeRange, type: filterType });
  };

  const handleUserClick = (userId) => {
    setSelectedUser(userId);
    // Could open user details modal or navigate to user profile
  };

  const handleRefresh = () => {
    queryClient.invalidateQueries(['admin-activities']);
    queryClient.invalidateQueries(['realtime-stats']);
    toast.success('Data refreshed');
  };

  // Computed values
  const activitiesList = activities?.activities || [];
  const totalActivities = activities?.total || 0;
  const totalPages = Math.ceil(totalActivities / itemsPerPage);

  const filteredActivities = useMemo(() => {
    return activitiesList.filter(activity => {
      const matchesSearch = searchTerm === '' || 
        activity.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        activity.user_name?.toLowerCase().includes(searchTerm.toLowerCase());
      return matchesSearch;
    });
  }, [activitiesList, searchTerm]);

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <h1 className={styles.title}>
              <Activity size={24} />
              {t('title', { defaultValue: 'Activity Monitoring' })}
              {realTimeEnabled && (
                <span className={styles.liveIndicator}>
                  <span className={styles.liveDot}></span>
                  LIVE
                </span>
              )}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Monitor real-time user activities, system health, and platform engagement' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button
              onClick={() => setRealTimeEnabled(!realTimeEnabled)}
              className={`${styles.realtimeToggle} ${realTimeEnabled ? styles.active : ''}`}
            >
              {realTimeEnabled ? <Pause size={16} /> : <Play size={16} />}
              {realTimeEnabled ? 'Pause' : 'Resume'} Live
            </button>
            <button onClick={handleRefresh} className={styles.refreshButton}>
              <RefreshCw size={16} />
              Refresh
            </button>
            <button onClick={handleExportLogs} className={styles.exportButton}>
              <Download size={16} />
              Export Logs
            </button>
          </div>
        </div>

        {/* Real-time Stats Cards */}
        {realtimeStats && (
          <div className={styles.realtimeStats}>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Users size={20} />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>{realtimeStats.online_users || 0}</div>
                <div className={styles.statLabel}>Online Users</div>
                <div className={styles.statTrend}>
                  {realtimeStats.users_trend > 0 ? <ArrowUp size={12} /> : <ArrowDown size={12} />}
                  {Math.abs(realtimeStats.users_trend || 0)}%
                </div>
              </div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Activity size={20} />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>{realtimeStats.activities_per_minute || 0}</div>
                <div className={styles.statLabel}>Activities/Min</div>
                <div className={styles.statTrend}>
                  {realtimeStats.activity_trend > 0 ? <ArrowUp size={12} /> : <ArrowDown size={12} />}
                  {Math.abs(realtimeStats.activity_trend || 0)}%
                </div>
              </div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Eye size={20} />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>{realtimeStats.page_views || 0}</div>
                <div className={styles.statLabel}>Page Views</div>
                <div className={styles.statTrend}>
                  {realtimeStats.views_trend > 0 ? <ArrowUp size={12} /> : <ArrowDown size={12} />}
                  {Math.abs(realtimeStats.views_trend || 0)}%
                </div>
              </div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <TrendingUp size={20} />
              </div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>{realtimeStats.engagement_rate || 0}%</div>
                <div className={styles.statLabel}>Engagement Rate</div>
                <div className={styles.statTrend}>
                  {realtimeStats.engagement_trend > 0 ? <ArrowUp size={12} /> : <ArrowDown size={12} />}
                  {Math.abs(realtimeStats.engagement_trend || 0)}%
                </div>
              </div>
            </div>
          </div>
        )}

        {/* System Health */}
        <SystemHealthCard health={systemHealth} />

        {/* Filters and Controls */}
        <div className={styles.filtersSection}>
          <div className={styles.searchBar}>
            <Search size={16} />
            <input
              type="text"
              placeholder="Search activities, users, or descriptions..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filters}>
            <select
              value={timeRange}
              onChange={(e) => setTimeRange(e.target.value)}
              className={styles.filterSelect}
            >
              <option value="1h">Last Hour</option>
              <option value="24h">Last 24 Hours</option>
              <option value="7d">Last 7 Days</option>
              <option value="30d">Last 30 Days</option>
              <option value="90d">Last 90 Days</option>
            </select>
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">All Activities</option>
              <option value="user_login">User Logins</option>
              <option value="post_create">Post Creation</option>
              <option value="post_like">Likes</option>
              <option value="post_comment">Comments</option>
              <option value="product_view">Product Views</option>
              <option value="product_purchase">Purchases</option>
              <option value="error">Errors</option>
            </select>
          </div>
        </div>

        {/* Main Content Grid */}
        <div className={styles.contentGrid}>
          {/* Activity Feed */}
          <div className={styles.activityFeedSection}>
            <div className={styles.sectionHeader}>
              <h2>Real-time Activity Feed</h2>
              <span className={styles.activityCount}>
                {totalActivities} total activities
              </span>
            </div>
            {activitiesLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p>Loading activities...</p>
              </div>
            ) : filteredActivities.length > 0 ? (
              <>
                <div className={styles.activityFeed}>
                  {filteredActivities.map((activity) => (
                    <ActivityFeedItem
                      key={activity.id}
                      activity={activity}
                      onUserClick={handleUserClick}
                    />
                  ))}
                </div>
                {totalPages > 1 && (
                  <div className={styles.pagination}>
                    <button
                      onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                      disabled={currentPage === 1}
                      className={styles.paginationButton}
                    >
                      Previous
                    </button>
                    <span className={styles.pageInfo}>
                      Page {currentPage} of {totalPages}
                    </span>
                    <button
                      onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                      disabled={currentPage === totalPages}
                      className={styles.paginationButton}
                    >
                      Next
                    </button>
                  </div>
                )}
              </>
            ) : (
              <div className={styles.emptyState}>
                <Activity size={48} />
                <h3>No Activities Found</h3>
                <p>
                  {searchTerm || filterType !== 'all'
                    ? 'Try adjusting your search or filters.'
                    : 'No activities have been recorded yet.'
                  }
                </p>
              </div>
            )}
          </div>

          {/* Analytics Sidebar */}
          <div className={styles.analyticsSidebar}>
            {/* Engagement Metrics */}
            {engagementMetrics && (
              <div className={styles.analyticsCard}>
                <h3>Engagement Metrics</h3>
                <div className={styles.engagementStats}>
                  <div className={styles.engagementStat}>
                    <Heart size={16} />
                    <span>{engagementMetrics.total_likes || 0} Likes</span>
                  </div>
                  <div className={styles.engagementStat}>
                    <MessageSquare size={16} />
                    <span>{engagementMetrics.total_comments || 0} Comments</span>
                  </div>
                  <div className={styles.engagementStat}>
                    <Share size={16} />
                    <span>{engagementMetrics.total_shares || 0} Shares</span>
                  </div>
                </div>
              </div>
            )}

            {/* Device Breakdown */}
            {deviceBreakdown && (
              <div className={styles.analyticsCard}>
                <h3>Device Breakdown</h3>
                <div className={styles.deviceStats}>
                  <div className={styles.deviceStat}>
                    <Monitor size={16} />
                    <span>Desktop</span>
                    <span>{deviceBreakdown.desktop || 0}%</span>
                  </div>
                  <div className={styles.deviceStat}>
                    <Smartphone size={16} />
                    <span>Mobile</span>
                    <span>{deviceBreakdown.mobile || 0}%</span>
                  </div>
                  <div className={styles.deviceStat}>
                    <Tablet size={16} />
                    <span>Tablet</span>
                    <span>{deviceBreakdown.tablet || 0}%</span>
                  </div>
                </div>
              </div>
            )}

            {/* Top Content */}
            {mostLiked && (
              <div className={styles.analyticsCard}>
                <h3>Most Liked Content</h3>
                <div className={styles.topContent}>
                  {mostLiked.items?.slice(0, 5).map((item, index) => (
                    <div key={item.id} className={styles.topContentItem}>
                      <span className={styles.contentRank}>#{index + 1}</span>
                      <div className={styles.contentInfo}>
                        <span className={styles.contentTitle}>
                          {item.title || `Item ${item.id}`}
                        </span>
                        <span className={styles.contentStats}>
                          {item.likes || 0} likes
                        </span>
                      </div>
                    </div>
                  )) || (
                    <p className={styles.emptyAnalytics}>No data available</p>
                  )}
                </div>
              </div>
            )}

            {/* Analytics Charts */}
            {analytics && (
              <div className={styles.analyticsCard}>
                <MetricsChart 
                  data={analytics.hourly_activity} 
                  type="line" 
                  title="Activity Trends" 
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default ActivityMonitoring;
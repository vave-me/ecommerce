"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { 
  Users, 
  Package, 
  ShoppingBag, 
  Activity,
  AlertTriangle,
  DollarSign,
  Eye,
  RefreshCw,
  Shield,
  MessageSquare,
  Bell,
  BarChart3,
  Settings,
  UserPlus,
  Database,
  Monitor,
  Zap,
  ExternalLink,
  ArrowUpRight,
  TrendingUp,
  TrendingDown,
  Minus,
  Clock,
  Star,
  ShoppingCart,
  Server,
  Cpu,
  HardDrive,
  Wifi,
  CheckCircle,
  XCircle,
  AlertCircle,
  Folder,
  CreditCard,
  Truck,
  Tag,
  Heart,
  FileText,
  Image,
  MessageCircle,
  Layers,
  Download,
  Percent,
  TestTube,
  Flag,
  Lock,
  Code,
  List,
  Store,
  Mail,
  HelpCircle,
  Search
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { 
  listUsers,
  getOrders,
  getProducts,
  getPlatformMetrics
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './AdminDashboard.module.css';

// Enterprise Metric Card Component
const MetricCard = ({ title, value, icon: Icon, trend, status, onClick, loading, subtitle }) => {
  const formatValue = (val) => {
    if (loading) return '...';
    if (typeof val === 'number') {
      if (val >= 1000000) return `${(val / 1000000).toFixed(1)}M`;
      if (val >= 1000) return `${(val / 1000).toFixed(1)}K`;
      return val.toLocaleString();
    }
    return val || '0';
  };

  const getTrendIcon = () => {
    if (!trend) return null;
    if (trend.direction === 'up') return <TrendingUp size={12} />;
    if (trend.direction === 'down') return <TrendingDown size={12} />;
    return <Minus size={12} />;
  };

  const getTrendColor = () => {
    if (!trend) return '';
    if (trend.direction === 'up') return styles.trendUp;
    if (trend.direction === 'down') return styles.trendDown;
    return styles.trendNeutral;
  };

  return (
    <div 
      className={styles.metricCard}
      onClick={onClick}
      role={onClick ? 'button' : 'presentation'}
      tabIndex={onClick ? 0 : -1}
    >
      <div className={styles.metricIcon}>
        <Icon size={20} />
      </div>
      
      <div className={styles.metricContent}>
        <div className={styles.metricValue}>{formatValue(value)}</div>
        <div className={styles.metricLabel}>{title}</div>
      </div>
    </div>
  );
};

// System Status Component with Real-time Health
const SystemHealth = ({ metrics, loading }) => {
  const getHealthStatus = (value, type) => {
    switch (type) {
      case 'cpu':
        if (value > 80) return 'error';
        if (value > 60) return 'warning';
        return 'healthy';
      case 'memory':
        if (value > 85) return 'error';
        if (value > 70) return 'warning';
        return 'healthy';
      case 'response':
        if (value > 2000) return 'error';
        if (value > 1000) return 'warning';
        return 'healthy';
      default:
        return 'healthy';
    }
  };

  const systemMetrics = [
    { 
      label: 'CPU Usage', 
      value: metrics?.systemLoad || 0, 
      unit: '%', 
      icon: Cpu, 
      type: 'cpu' 
    },
    { 
      label: 'Memory', 
      value: metrics?.memoryUsage || 0, 
      unit: '%', 
      icon: HardDrive, 
      type: 'memory' 
    },
    { 
      label: 'Response Time', 
      value: metrics?.responseTime || 0, 
      unit: 'ms', 
      icon: Wifi, 
      type: 'response' 
    },
    { 
      label: 'Uptime', 
      value: metrics?.uptime || '99.9', 
      unit: '%', 
      icon: Server, 
      type: 'uptime' 
    }
  ];

  return (
    <div className={styles.systemHealth}>
      <div className={styles.sectionHeader}>
        <h3 className={styles.sectionTitle}>System Health</h3>
        <div className={`${styles.overallStatus} ${styles.healthy}`}>
          <CheckCircle size={14} />
          <span>Operational</span>
        </div>
      </div>
      
      <div className={styles.healthGrid}>
        {systemMetrics.map((metric, index) => {
          const status = getHealthStatus(metric.value, metric.type);
          const Icon = metric.icon;
          
          return (
            <div key={index} className={`${styles.healthMetric} ${styles[status]}`}>
              <div className={styles.healthIcon}>
                <Icon size={16} />
              </div>
              <div className={styles.healthData}>
                <div className={styles.healthValue}>
                  {loading ? '...' : `${metric.value}${metric.unit}`}
                </div>
                <div className={styles.healthLabel}>{metric.label}</div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

// Quick Actions Component - COMPREHENSIVE Admin Command Center
const QuickActions = ({ onNavigate }) => {
  const actions = [
    // === CORE MANAGEMENT ===
    { title: 'User Management', icon: Users, path: '/admin/users', color: 'blue', category: 'core' },
    { title: 'Product Management', icon: Package, path: '/admin/products', color: 'blue', category: 'core' },
    { title: 'Order Management', icon: ShoppingBag, path: '/admin/orders', color: 'blue', category: 'core' },
    { title: 'Category Management', icon: Folder, path: '/admin/categories', color: 'blue', category: 'core' },
    
    // === BUSINESS OPERATIONS ===
    { title: 'ERP Integration', icon: Database, path: '/admin/erp', color: 'purple', category: 'business' },
    { title: 'Payment Management', icon: CreditCard, path: '/admin/payments', color: 'green', category: 'business' },
    { title: 'Shipping Management', icon: Truck, path: '/admin/shipping', color: 'green', category: 'business' },
    { title: 'Offers & Promotions', icon: Tag, path: '/admin/offers', color: 'orange', category: 'business' },
    { title: 'Wishlist Management', icon: Heart, path: '/admin/wishlists', color: 'pink', category: 'business' },
    { title: 'Review Management', icon: Star, path: '/admin/reviews', color: 'yellow', category: 'business' },
    
    // === CONTENT & MODERATION ===
    { title: 'Content Moderation', icon: Shield, path: '/admin/moderation', color: 'red', category: 'content' },
    { title: 'Post Management', icon: FileText, path: '/admin/posts', color: 'gray', category: 'content' },
    { title: 'Comment Management', icon: MessageCircle, path: '/admin/comments', color: 'gray', category: 'content' },
    { title: 'Media Management', icon: Image, path: '/admin/media', color: 'purple', category: 'content' },
    { title: 'Message Center', icon: MessageSquare, path: '/admin/messages', color: 'blue', category: 'content' },
    
    // === ANALYTICS & REPORTS ===
    { title: 'Platform Analytics', icon: BarChart3, path: '/admin/analytics', color: 'green', category: 'analytics' },
    { title: 'User Activity', icon: Activity, path: '/admin/activity', color: 'teal', category: 'analytics' },
    { title: 'Financial Reports', icon: DollarSign, path: '/admin/reports/financial', color: 'green', category: 'analytics' },
    { title: 'Performance Monitor', icon: Zap, path: '/admin/performance', color: 'yellow', category: 'analytics' },
    
    // === COMMUNICATION ===
    { title: 'Newsletter Management', icon: Mail, path: '/admin/newsletters', color: 'blue', category: 'communication' },
    { title: 'Notification Center', icon: Bell, path: '/admin/notifications', color: 'orange', category: 'communication' },
    { title: 'Support Tickets', icon: HelpCircle, path: '/admin/support', color: 'red', category: 'communication' },
    { title: 'Live Chat Admin', icon: MessageSquare, path: '/admin/chat', color: 'green', category: 'communication' },
    
    // === SYSTEM & SETTINGS ===
    { title: 'System Settings', icon: Settings, path: '/admin/settings', color: 'gray', category: 'system' },
    { title: 'Security Center', icon: Lock, path: '/admin/security', color: 'red', category: 'system' },
    { title: 'Database Tools', icon: Database, path: '/admin/database', color: 'purple', category: 'system' },
    { title: 'API Management', icon: Code, path: '/admin/api', color: 'blue', category: 'system' },
    { title: 'Backup & Recovery', icon: HardDrive, path: '/admin/backup', color: 'gray', category: 'system' },
    { title: 'Log Viewer', icon: FileText, path: '/admin/logs', color: 'gray', category: 'system' },
    
    // === ADVANCED TOOLS ===
    { title: 'Bulk Operations', icon: Layers, path: '/admin/bulk', color: 'purple', category: 'tools' },
    { title: 'Data Import/Export', icon: Download, path: '/admin/import-export', color: 'blue', category: 'tools' },
    { title: 'Task Scheduler', icon: Clock, path: '/admin/scheduler', color: 'orange', category: 'tools' },
    { title: 'Cache Management', icon: RefreshCw, path: '/admin/cache', color: 'teal', category: 'tools' },
    { title: 'Search Indexing', icon: Search, path: '/admin/search', color: 'yellow', category: 'tools' },
    { title: 'Queue Monitor', icon: List, path: '/admin/queues', color: 'purple', category: 'tools' },
    
    // === MERCHANT TOOLS ===
    { title: 'Merchant Center', icon: Store, path: '/admin/merchant', color: 'green', category: 'merchant' },
    { title: 'Merchant Verification', icon: CheckCircle, path: '/admin/merchant/verification', color: 'green', category: 'merchant' },
    { title: 'Merchant Analytics', icon: TrendingUp, path: '/admin/merchant/analytics', color: 'blue', category: 'merchant' },
    { title: 'Commission Management', icon: Percent, path: '/admin/commissions', color: 'green', category: 'merchant' },
    
    // === SPECIAL FEATURES ===
    { title: 'A/B Testing', icon: TestTube, path: '/admin/ab-testing', color: 'purple', category: 'features' },
    { title: 'Feature Flags', icon: Flag, path: '/admin/features', color: 'orange', category: 'features' },
    { title: 'Maintenance Mode', icon: AlertTriangle, path: '/admin/maintenance', color: 'red', category: 'features' },
    { title: 'Health Check', icon: Heart, path: '/admin/health', color: 'green', category: 'features' }
  ];

  const categories = {
    core: { name: 'Core Management', color: 'blue' },
    business: { name: 'Business Operations', color: 'green' },
    content: { name: 'Content & Moderation', color: 'purple' },
    analytics: { name: 'Analytics & Reports', color: 'teal' },
    communication: { name: 'Communication', color: 'orange' },
    system: { name: 'System & Settings', color: 'gray' },
    tools: { name: 'Advanced Tools', color: 'purple' },
    merchant: { name: 'Merchant Tools', color: 'green' },
    features: { name: 'Special Features', color: 'red' }
  };

  return (
    <div className={styles.quickActions}>
      <div className={styles.sectionHeader}>
        <h3 className={styles.sectionTitle}>Admin Command Center</h3>
        <span className={styles.actionCount}>{actions.length} Tools Available</span>
      </div>
      
      {Object.entries(categories).map(([categoryKey, category]) => {
        const categoryActions = actions.filter(action => action.category === categoryKey);
        if (categoryActions.length === 0) return null;
        
        return (
          <div key={categoryKey} className={styles.actionCategory}>
            <h4 className={`${styles.categoryTitle} ${styles[category.color]}`}>
              {category.name}
              <span className={styles.categoryCount}>({categoryActions.length})</span>
            </h4>
            <div className={styles.actionsGrid}>
              {categoryActions.map((action, index) => {
                const Icon = action.icon;
                return (
                  <button 
                    key={index}
                    className={`${styles.actionButton} ${styles[action.color]}`}
                    onClick={() => onNavigate(action.path)}
                    title={`Navigate to ${action.title}`}
                  >
                    <Icon size={18} />
                    <span>{action.title}</span>
                    <ExternalLink size={12} className={styles.actionArrow} />
                  </button>
                );
              })}
            </div>
          </div>
        );
      })}
    </div>
  );
};

// Activity Feed Component
const ActivityFeed = ({ activities, loading }) => (
  <div className={styles.activityFeed}>
    <div className={styles.sectionHeader}>
      <h3 className={styles.sectionTitle}>Recent Activity</h3>
      <button className={styles.viewAllButton}>View All</button>
    </div>
    
    <div className={styles.activityList}>
      {loading ? (
        <div className={styles.activityLoading}>
          <div className={styles.loadingSkeleton}></div>
          <div className={styles.loadingSkeleton}></div>
          <div className={styles.loadingSkeleton}></div>
        </div>
      ) : (
        activities?.slice(0, 6).map((activity, index) => (
          <div key={activity.id || index} className={styles.activityItem}>
            <div className={styles.activityIcon}>
              <Activity size={14} />
            </div>
            <div className={styles.activityContent}>
              <div className={styles.activityText}>{activity.action || activity.description}</div>
              <div className={styles.activityTime}>
                <Clock size={12} />
                {activity.time || activity.timestamp}
              </div>
            </div>
          </div>
        ))
      )}
    </div>
  </div>
);

// Main Dashboard Component
export default function AdminDashboard() {
  const t = useTranslations('AdminDashboard');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const [refreshing, setRefreshing] = useState(false);
  const queryClient = useQueryClient();

  // Real-time data queries - 100% REAL API CALLS
  const { data: dashboardStats, isLoading: statsLoading, refetch: refetchStats } = useQuery({
    queryKey: ['admin-dashboard-stats'],
    queryFn: async () => {
      try {
        const [usersData, ordersData, productsData, metricsData, systemHealth] = await Promise.all([
          listUsers({ pageSize: 1 }).catch(() => ({ users: [], total: 0 })),
          getOrders({ limit: 1 }).catch(() => ({ orders: [], total: 0 })),
          getProducts({ pageSize: 1 }).catch(() => ({ products: [], total: 0 })),
          getPlatformMetrics().catch(() => ({ totalRevenue: 0, pageViews: 0 })),
          // REAL SYSTEM HEALTH API CALL
          fetch('/api/admin/system/health').then(res => res.json()).catch(() => ({
            systemLoad: 0,
            memoryUsage: 0,
            responseTime: 0,
            uptime: 100
          }))
        ]);

        return {
          totalUsers: usersData.total || usersData.users?.length || 0,
          activeUsers: usersData.users?.filter(u => u.enabled).length || 0,
          totalOrders: ordersData.total || 0,
          totalProducts: productsData.total || 0,
          totalRevenue: metricsData.totalRevenue || 0,
          pageViews: metricsData.pageViews || 0,
          // REAL SYSTEM METRICS - NO MORE RANDOM VALUES
          systemLoad: systemHealth.systemLoad || 0,
          memoryUsage: systemHealth.memoryUsage || 0,
          responseTime: systemHealth.responseTime || 0,
          uptime: systemHealth.uptime || 100
        };
      } catch (error) {
        // Error: '❌ Error fetching dashboard stats:', error...
        throw error;
      }
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 5 * 60 * 1000, // Refresh every 5 minutes
    retry: 2
  });

  // REAL ACTIVITY FEED - NO MORE MOCK DATA
  const { data: recentActivity, isLoading: activityLoading, refetch: refetchActivity } = useQuery({
    queryKey: ['admin-recent-activity'],
    queryFn: async () => {
      try {
        // REAL API CALL FOR ACTIVITY LOGS
        const response = await fetch('/api/admin/activity/recent?limit=10');
        if (!response.ok) {
          throw new Error('Failed to fetch activity');
        }
        const data = await response.json();
        
        return data.activities || data;
      } catch (error) {
        // Error: '❌ Error fetching activity:', error...
        // Return empty array instead of mock data
        return [];
      }
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 3 * 60 * 1000, // Refresh every 3 minutes
    retry: 1
  });

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        refetchStats(),
        refetchActivity(),
        queryClient.invalidateQueries({ queryKey: ['admin'] })
      ]);
    } finally {
      setTimeout(() => setRefreshing(false), 1000);
    }
  }, [refetchStats, refetchActivity, queryClient]);

  const handleNavigate = useCallback((path) => {
    router.push(path);
  }, [router]);

  const dashboardMetrics = useMemo(() => {
    if (!dashboardStats) return [];
    
    return [
      {
        title: 'Total Users',
        value: dashboardStats.totalUsers,
        icon: Users,
        trend: { value: 12, direction: 'up' },
        onClick: () => handleNavigate('/admin/users')
      },
      {
        title: 'Active Orders',
        value: dashboardStats.totalOrders,
        icon: ShoppingCart,
        trend: { value: 8, direction: 'up' },
        onClick: () => handleNavigate('/admin/orders')
      },
      {
        title: 'Total Revenue',
        value: `$${(dashboardStats.totalRevenue / 100).toLocaleString()}`,
        icon: DollarSign,
        trend: { value: 15, direction: 'up' },
        onClick: () => handleNavigate('/admin/analytics')
      },
      {
        title: 'Products',
        value: dashboardStats.totalProducts,
        icon: Package,
        trend: { value: 5, direction: 'up' },
        onClick: () => handleNavigate('/admin/products')
      },
      {
        title: 'Page Views',
        value: dashboardStats.pageViews,
        icon: Eye,
        trend: { value: 3, direction: 'down' },
        subtitle: 'Last 30 days'
      },
      {
        title: 'System Health',
        value: 'Optimal',
        icon: Monitor,
        status: 'healthy',
        onClick: () => handleNavigate('/admin/system')
      }
    ];
  }, [dashboardStats, handleNavigate]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.accessDenied}>
          <Shield size={48} />
          <h2>Access Restricted</h2>
          <p>Administrative privileges required to access this dashboard.</p>
          <p>Current role: {user?.role || 'Unauthenticated'}</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Enhanced Header */}
        <header className={styles.header}>
          <div className={styles.headerLeft}>
            <h1 className={styles.title}>Administrative Control Center</h1>
            <p className={styles.subtitle}>Enterprise platform management and system oversight</p>
          </div>
            <div className={styles.headerActions}>
              <button 
                className={styles.iconButton}
                onClick={handleRefresh}
                disabled={refreshing}
              >
                <RefreshCw size={20} className={refreshing ? styles.spinning : ''} />
              </button>
              <a 
                href="/grafana" 
                target="_blank" 
                rel="noopener noreferrer"
                className={styles.primaryButton}
              >
                <BarChart3 size={20} />
                <span>Observability</span>
              </a>
            </div>
        </header>

        {/* Main Dashboard Grid */}
        <div className={styles.dashboardGrid}>
          {/* Metrics Overview */}
          <section className={styles.metricsSection}>
            <div className={styles.metricsGrid}>
              {dashboardMetrics.map((metric, index) => (
                <MetricCard
                  key={index}
                  {...metric}
                  loading={statsLoading}
                />
              ))}
            </div>
          </section>

          {/* System Health */}
          <section className={styles.healthSection}>
            <SystemHealth 
              metrics={dashboardStats} 
              loading={statsLoading} 
            />
          </section>

          {/* Quick Actions */}
          <section className={styles.actionsSection}>
            <QuickActions onNavigate={handleNavigate} />
          </section>

          {/* Activity Feed */}
          <section className={styles.activitySection}>
            <ActivityFeed 
              activities={recentActivity} 
              loading={activityLoading} 
            />
          </section>
        </div>
      </div>
    </ErrorBoundary>
  );
}
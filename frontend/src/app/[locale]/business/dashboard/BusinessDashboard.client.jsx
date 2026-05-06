"use client";

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { 
  Package, 
  ShoppingBag, 
  DollarSign,
  TrendingUp,
  Eye,
  Heart,
  MessageSquare,
  Star,
  Plus,
  Edit,
  BarChart3,
  Settings,
  RefreshCw,
  ArrowUpRight,
  ArrowDownRight,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  CreditCard,
  FileText
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { getBusinessDashboardStats, getUserMetrics } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './BusinessDashboard.module.css';

const MetricCard = ({ title, value, icon: Icon, trend, color = 'primary', onClick }) => {
  const isPositive = trend > 0;
  const TrendIcon = isPositive ? ArrowUpRight : ArrowDownRight;
  
  return (
    <div 
      className={`${styles.metricCard} ${onClick ? styles.clickable : ''}`}
      onClick={onClick}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
    >
      <div className={styles.metricHeader}>
        <div className={`${styles.metricIcon} ${styles[`metricIcon--${color}`]}`}>
          <Icon size={20} />
        </div>
        <div className={styles.metricTitle}>{title}</div>
      </div>
      <div className={styles.metricContent}>
        <div className={styles.metricValue}>{value}</div>
        {trend !== undefined && (
          <div className={`${styles.metricTrend} ${isPositive ? styles.trendUp : styles.trendDown}`}>
            <TrendIcon size={16} />
            <span>{Math.abs(trend)}%</span>
          </div>
        )}
      </div>
    </div>
  );
};

const ProductRow = ({ product, onEdit }) => {
  const getStatusIcon = (status) => {
    switch (status) {
      case 'active': return <CheckCircle size={16} className={styles.statusActive} />;
      case 'pending': return <Clock size={16} className={styles.statusPending} />;
      case 'rejected': return <XCircle size={16} className={styles.statusRejected} />;
      default: return <AlertCircle size={16} className={styles.statusDefault} />;
    }
  };
  
  return (
    <tr className={styles.productRow}>
      <td className={styles.productCell}>
        <div className={styles.productInfo}>
          {product.thumbnail && (
            <img 
              src={product.thumbnail} 
              alt="" 
              className={styles.productThumb}
            />
          )}
          <div>
            <div className={styles.productName}>{product.name}</div>
            <div className={styles.productCategory}>{product.category}</div>
          </div>
        </div>
      </td>
      <td className={styles.statusCell}>
        <div className={styles.statusBadge}>
          {getStatusIcon(product.status)}
          <span>{product.status}</span>
        </div>
      </td>
      <td className={styles.priceCell}>${product.price}</td>
      <td className={styles.stockCell}>{product.stock}</td>
      <td className={styles.viewsCell}>{product.views || 0}</td>
      <td className={styles.salesCell}>{product.sales || 0}</td>
      <td className={styles.actionCell}>
        <button 
          className={styles.editButton}
          onClick={() => onEdit(product)}
          aria-label="Edit product"
        >
          <Edit size={16} />
        </button>
      </td>
    </tr>
  );
};

const OrderItem = ({ order }) => {
  const getStatusColor = (status) => {
    switch (status) {
      case 'pending': return 'warning';
      case 'processing': return 'info';
      case 'completed': return 'success';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };
  
  return (
    <div className={styles.orderItem}>
      <div className={styles.orderHeader}>
        <span className={styles.orderId}>#{order.id}</span>
        <span className={`${styles.orderStatus} ${styles[`orderStatus--${getStatusColor(order.status)}`]}`}>
          {order.status}
        </span>
      </div>
      <div className={styles.orderDetails}>
        <div className={styles.orderProduct}>{order.productName}</div>
        <div className={styles.orderMeta}>
          <span className={styles.orderPrice}>${order.total}</span>
          <span className={styles.orderDate}>{new Date(order.createdAt).toLocaleDateString()}</span>
        </div>
      </div>
    </div>
  );
};

const BusinessDashboard = () => {
  const t = useTranslations('BusinessDashboard');
  const router = useRouter();
  const { user } = useAuth();
  const { isBusiness, isAdmin } = useUserRole();
  const [refreshing, setRefreshing] = useState(false);
  const [activeTab, setActiveTab] = useState('overview');

  // Redirect if not business user
  useEffect(() => {
    if (!isBusiness && !isAdmin) {
      router.push('/');
    }
  }, [isBusiness, isAdmin, router]);

  // Fetch business dashboard data
  const { data: dashboardData, isLoading, refetch } = useQuery({
    queryKey: ['businessDashboard', user?.id],
    queryFn: () => getBusinessDashboardStats(user?.id),
    enabled: !!user?.id && (isBusiness || isAdmin),
    staleTime: 60000, // 1 minute
  });

  // Fetch user metrics
  const { data: userMetrics } = useQuery({
    queryKey: ['userMetrics', user?.id],
    queryFn: () => getUserMetrics(user?.id),
    enabled: !!user?.id && (isBusiness || isAdmin),
    staleTime: 60000,
  });

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    await refetch();
    setTimeout(() => setRefreshing(false), 500);
  }, [refetch]);

  const handleEditProduct = useCallback((product) => {
    router.push(`/products/${product.id}/edit`);
  }, [router]);

  // Calculate metrics from data
  const metrics = useMemo(() => {
    if (!dashboardData || !userMetrics?.metric) return [];
    
    const metric = userMetrics.metric;
    
    return [
      {
        title: t('metrics.totalProducts', { defaultValue: 'Total Products' }),
        value: dashboardData.products?.length || 0,
        icon: Package,
        color: 'primary',
        onClick: () => router.push('/business/products'),
      },
      {
        title: t('metrics.totalOrders', { defaultValue: 'Total Orders' }),
        value: dashboardData.orders?.length || 0,
        icon: ShoppingBag,
        color: 'info',
        trend: 5.2,
        onClick: () => router.push('/business/orders'),
      },
      {
        title: t('metrics.revenue', { defaultValue: 'Revenue' }),
        value: `$${dashboardData.revenue?.toLocaleString() || 0}`,
        icon: DollarSign,
        color: 'success',
        trend: 12.5,
        onClick: () => router.push('/business/revenue'),
      },
      {
        title: t('metrics.views', { defaultValue: 'Total Views' }),
        value: parseInt(metric.visitedCount || 0).toLocaleString(),
        icon: Eye,
        color: 'warning',
        trend: 8.3,
      },
      {
        title: t('metrics.likes', { defaultValue: 'Total Likes' }),
        value: parseInt(metric.likesCount || 0).toLocaleString(),
        icon: Heart,
        color: 'error',
        trend: 15.7,
      },
      {
        title: t('metrics.rating', { defaultValue: 'Average Rating' }),
        value: (metric.rating ? (parseInt(metric.rating) / 20).toFixed(1) : '0.0') + '/5',
        icon: Star,
        color: 'warning',
      },
    ];
  }, [dashboardData, userMetrics, t, router]);

  // Business actions - comprehensive list based on APIs
  const businessActions = useMemo(() => [
    {
      id: 'add-product',
      icon: Plus,
      label: t('actions.addProduct', { defaultValue: 'Add Product' }),
      description: t('actions.addProductDesc', { defaultValue: 'List a new product' }),
      onClick: () => router.push('/products/new'),
      primary: true,
    },
    {
      id: 'manage-products',
      icon: Package,
      label: t('actions.manageProducts', { defaultValue: 'Manage Products' }),
      description: t('actions.manageProductsDesc', { defaultValue: 'View and edit your listings' }),
      onClick: () => router.push('/business/products'),
    },
    {
      id: 'view-orders',
      icon: ShoppingBag,
      label: t('actions.viewOrders', { defaultValue: 'View Orders' }),
      description: t('actions.viewOrdersDesc', { defaultValue: 'Track customer orders' }),
      onClick: () => router.push('/business/orders'),
    },
    {
      id: 'manage-offers',
      icon: TrendingUp,
      label: t('actions.manageOffers', { defaultValue: 'Manage Offers' }),
      description: t('actions.manageOffersDesc', { defaultValue: 'View buy, lease, pawn offers' }),
      onClick: () => router.push('/business/offers'),
    },
    {
      id: 'create-invoice',
      icon: CreditCard,
      label: t('actions.createInvoice', { defaultValue: 'Create Invoice' }),
      description: t('actions.createInvoiceDesc', { defaultValue: 'Generate invoices for customers' }),
      onClick: () => router.push('/business/invoices/new'),
    },
    {
      id: 'shipping-labels',
      icon: Package,
      label: t('actions.shippingLabels', { defaultValue: 'Shipping Labels' }),
      description: t('actions.shippingLabelsDesc', { defaultValue: 'Create shipping labels' }),
      onClick: () => router.push('/business/shipping'),
    },
    {
      id: 'analytics',
      icon: BarChart3,
      label: t('actions.analytics', { defaultValue: 'Analytics' }),
      description: t('actions.analyticsDesc', { defaultValue: 'View detailed insights' }),
      onClick: () => router.push('/business/analytics'),
    },
    {
      id: 'manage-posts',
      icon: FileText,
      label: t('actions.managePosts', { defaultValue: 'Manage Posts' }),
      description: t('actions.managePostsDesc', { defaultValue: 'Create and manage blog posts' }),
      onClick: () => router.push('/business/posts'),
    },
  ], [t, router]);

  if (!isBusiness && !isAdmin) {
    return null;
  }

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <h1 className={styles.title}>{t('title', { defaultValue: 'Business Dashboard' })}</h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage your products and track performance' })}
            </p>
          </div>
          <button 
            className={`${styles.refreshButton} ${refreshing ? styles.refreshing : ''}`}
            onClick={handleRefresh}
            disabled={refreshing}
            aria-label="Refresh dashboard"
          >
            <RefreshCw size={20} />
          </button>
        </div>

        {/* Metrics Grid */}
        <div className={styles.metricsGrid}>
          {metrics.map((metric, index) => (
            <MetricCard key={index} {...metric} />
          ))}
        </div>

        {/* Quick Actions */}
        <div className={styles.quickActions}>
          {businessActions.map((action) => (
            <button
              key={action.id}
              className={`${styles.actionButton} ${action.primary ? styles.actionButtonPrimary : ''}`}
              onClick={action.onClick}
            >
              <action.icon size={20} />
              <span>{action.label}</span>
            </button>
          ))}
        </div>

        {/* Content Tabs */}
        <div className={styles.contentSection}>
          <div className={styles.tabsHeader}>
            <button
              className={`${styles.tab} ${activeTab === 'overview' ? styles.tabActive : ''}`}
              onClick={() => setActiveTab('overview')}
            >
              {t('tabs.overview', { defaultValue: 'Overview' })}
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'products' ? styles.tabActive : ''}`}
              onClick={() => setActiveTab('products')}
            >
              {t('tabs.products', { defaultValue: 'Products' })}
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'orders' ? styles.tabActive : ''}`}
              onClick={() => setActiveTab('orders')}
            >
              {t('tabs.orders', { defaultValue: 'Recent Orders' })}
            </button>
          </div>

          <div className={styles.tabContent}>
            {activeTab === 'overview' && (
              <div className={styles.overviewGrid}>
                {/* Performance Chart Placeholder */}
                <div className={styles.chartCard}>
                  <h3 className={styles.cardTitle}>
                    {t('performance.title', { defaultValue: 'Performance Overview' })}
                  </h3>
                  <div className={styles.chartPlaceholder}>
                    <BarChart3 size={48} />
                    <p>{t('performance.placeholder', { defaultValue: 'Performance charts coming soon' })}</p>
                  </div>
                </div>

                {/* Recent Activity */}
                <div className={styles.activityCard}>
                  <h3 className={styles.cardTitle}>
                    {t('activity.title', { defaultValue: 'Recent Activity' })}
                  </h3>
                  <div className={styles.activityList}>
                    <p className={styles.emptyMessage}>
                      {t('activity.empty', { defaultValue: 'No recent activity' })}
                    </p>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'products' && (
              <div className={styles.productsSection}>
                <div className={styles.tableContainer}>
                  <table className={styles.productsTable}>
                    <thead>
                      <tr>
                        <th>{t('table.product', { defaultValue: 'Product' })}</th>
                        <th>{t('table.status', { defaultValue: 'Status' })}</th>
                        <th>{t('table.price', { defaultValue: 'Price' })}</th>
                        <th>{t('table.stock', { defaultValue: 'Stock' })}</th>
                        <th>{t('table.views', { defaultValue: 'Views' })}</th>
                        <th>{t('table.sales', { defaultValue: 'Sales' })}</th>
                        <th>{t('table.actions', { defaultValue: 'Actions' })}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {dashboardData?.products?.length > 0 ? (
                        dashboardData.products.map(product => (
                          <ProductRow 
                            key={product.id} 
                            product={product} 
                            onEdit={handleEditProduct}
                          />
                        ))
                      ) : (
                        <tr>
                          <td colSpan="7" className={styles.emptyCell}>
                            {t('products.empty', { defaultValue: 'No products found' })}
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {activeTab === 'orders' && (
              <div className={styles.ordersSection}>
                {dashboardData?.orders?.length > 0 ? (
                  <div className={styles.ordersList}>
                    {dashboardData.orders.map(order => (
                      <OrderItem key={order.id} order={order} />
                    ))}
                  </div>
                ) : (
                  <p className={styles.emptyMessage}>
                    {t('orders.empty', { defaultValue: 'No recent orders' })}
                  </p>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default BusinessDashboard;
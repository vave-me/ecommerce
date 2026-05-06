"use client";

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { 
  TrendingUp, 
  TrendingDown, 
  Users, 
  Package, 
  ShoppingBag, 
  DollarSign,
  Eye,
  Clock,
  Calendar,
  Download,
  RefreshCw,
  ArrowUpRight,
  ArrowDownRight,
  Percent,
  Activity
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { getPlatformAnalytics, exportAnalytics } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './PlatformAnalytics.module.css';

// Simple Chart Component
const SimpleLineChart = ({ data, height = 200 }) => {
  const maxValue = Math.max(...data.map(d => d.value));
  const minValue = Math.min(...data.map(d => d.value));
  const range = maxValue - minValue;
  
  const points = data.map((d, i) => {
    const x = (i / (data.length - 1)) * 100;
    const y = 100 - ((d.value - minValue) / range) * 100;
    return { x, y, ...d };
  });
  
  const pathData = points
    .map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`)
    .join(' ');
  
  return (
    <div className={styles.chartContainer} style={{ height }}>
      <svg viewBox="0 0 100 100" preserveAspectRatio="none" className={styles.chart}>
        <path
          d={pathData}
          fill="none"
          stroke="var(--color-primary-500)"
          strokeWidth="2"
          vectorEffect="non-scaling-stroke"
        />
        <path
          d={`${pathData} L 100 100 L 0 100 Z`}
          fill="url(#gradient)"
          opacity="0.2"
        />
        <defs>
          <linearGradient id="gradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="var(--color-primary-500)" />
            <stop offset="100%" stopColor="var(--color-primary-500)" stopOpacity="0" />
          </linearGradient>
        </defs>
      </svg>
      <div className={styles.chartLabels}>
        {data.map((d, i) => (
          <div key={i} className={styles.chartLabel} style={{ left: `${(i / (data.length - 1)) * 100}%` }}>
            {d.label}
          </div>
        ))}
      </div>
    </div>
  );
};

const MetricCard = ({ title, value, change, icon: Icon, chart }) => {
  const isPositive = change > 0;
  const TrendIcon = isPositive ? ArrowUpRight : ArrowDownRight;
  
  return (
    <div className={styles.metricCard}>
      <div className={styles.metricHeader}>
        <div className={styles.metricIcon}>
          <Icon size={20} />
        </div>
        <div className={styles.metricChange} data-positive={isPositive}>
          <TrendIcon size={16} />
          <span>{Math.abs(change)}%</span>
        </div>
      </div>
      <h3 className={styles.metricTitle}>{title}</h3>
      <div className={styles.metricValue}>{value}</div>
      {chart && (
        <div className={styles.metricChart}>
          <SimpleLineChart data={chart} height={60} />
        </div>
      )}
    </div>
  );
};

const PlatformAnalytics = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('PlatformAnalytics');
  } catch (e) {
    // Fallback function for missing translations
    t = (key, options) => options?.defaultValue || key;
  }
  
  const router = useRouter();
  const { isAdmin } = useUserRole();
  const [dateRange, setDateRange] = useState('7d');
  const [refreshing, setRefreshing] = useState(false);

  // Redirect if not admin
  useEffect(() => {
    if (!isAdmin) {
      router.push('/');
    }
  }, [isAdmin, router]);

  // Fetch analytics data
  const { data: analyticsData, isLoading, refetch } = useQuery({
    queryKey: ['platformAnalytics', dateRange],
    queryFn: () => getPlatformAnalytics(dateRange),
    staleTime: 60000,
  });

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    await refetch();
    setTimeout(() => setRefreshing(false), 500);
  }, [refetch]);

  const handleExport = useCallback(async () => {
    try {
      const blob = await exportAnalytics(dateRange);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `analytics-export-${dateRange}-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, [dateRange]);

  if (!isAdmin) {
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
        {/* Header */}
        <div className={styles.header}>
          <div>
            <h1 className={styles.title}>{t('title', { defaultValue: 'Platform Analytics' })}</h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Track platform performance and user behavior' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <select 
              className={styles.dateRangeSelect}
              value={dateRange}
              onChange={(e) => setDateRange(e.target.value)}
            >
              <option value="7d">{t('last7Days', { defaultValue: 'Last 7 days' })}</option>
              <option value="30d">{t('last30Days', { defaultValue: 'Last 30 days' })}</option>
              <option value="90d">{t('last90Days', { defaultValue: 'Last 90 days' })}</option>
            </select>
            <button 
              className={styles.exportButton}
              onClick={handleExport}
            >
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button 
              className={`${styles.refreshButton} ${refreshing ? styles.refreshing : ''}`}
              onClick={handleRefresh}
              disabled={refreshing}
            >
              <RefreshCw size={16} />
            </button>
          </div>
        </div>

        {/* Key Metrics */}
        <div className={styles.metricsGrid}>
          <MetricCard 
            title={t('totalRevenue', { defaultValue: 'Total Revenue' })}
            value={analyticsData?.metrics.totalRevenue.value}
            change={analyticsData?.metrics.totalRevenue.change}
            icon={DollarSign}
            chart={analyticsData?.metrics.totalRevenue.chart}
          />
          <MetricCard 
            title={t('activeUsers', { defaultValue: 'Active Users' })}
            value={analyticsData?.metrics.activeUsers.value}
            change={analyticsData?.metrics.activeUsers.change}
            icon={Users}
            chart={analyticsData?.metrics.activeUsers.chart}
          />
          <MetricCard 
            title={t('totalOrders', { defaultValue: 'Total Orders' })}
            value={analyticsData?.metrics.totalOrders.value}
            change={analyticsData?.metrics.totalOrders.change}
            icon={ShoppingBag}
            chart={analyticsData?.metrics.totalOrders.chart}
          />
          <MetricCard 
            title={t('conversionRate', { defaultValue: 'Conversion Rate' })}
            value={analyticsData?.metrics.conversionRate.value}
            change={analyticsData?.metrics.conversionRate.change}
            icon={Percent}
            chart={analyticsData?.metrics.conversionRate.chart}
          />
          <MetricCard 
            title={t('avgOrderValue', { defaultValue: 'Avg Order Value' })}
            value={analyticsData?.metrics.avgOrderValue.value}
            change={analyticsData?.metrics.avgOrderValue.change}
            icon={TrendingUp}
            chart={analyticsData?.metrics.avgOrderValue.chart}
          />
          <MetricCard 
            title={t('pageViews', { defaultValue: 'Page Views' })}
            value={analyticsData?.metrics.pageViews.value}
            change={analyticsData?.metrics.pageViews.change}
            icon={Eye}
            chart={analyticsData?.metrics.pageViews.chart}
          />
        </div>

        {/* Content Sections */}
        <div className={styles.contentGrid}>
          {/* Top Products */}
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>
              {t('topProducts', { defaultValue: 'Top Products' })}
            </h2>
            <div className={styles.productsList}>
              {analyticsData?.topProducts.map((product, index) => (
                <div key={product.id} className={styles.productItem}>
                  <div className={styles.productRank}>{index + 1}</div>
                  <div className={styles.productInfo}>
                    <div className={styles.productName}>{product.name}</div>
                    <div className={styles.productStats}>
                      <span>{product.sales} sales</span>
                      <span className={styles.separator}>•</span>
                      <span>{product.revenue}</span>
                    </div>
                  </div>
                  <div className={`${styles.productGrowth} ${product.growth > 0 ? styles.positive : styles.negative}`}>
                    {product.growth > 0 ? <ArrowUpRight size={16} /> : <ArrowDownRight size={16} />}
                    <span>{Math.abs(product.growth)}%</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* User Activity */}
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>
              {t('userActivity', { defaultValue: 'User Activity' })}
            </h2>
            <div className={styles.activityStats}>
              <div className={styles.activityStat}>
                <div className={styles.activityIcon}>
                  <Users size={20} />
                </div>
                <div>
                  <div className={styles.activityValue}>{analyticsData?.userActivity.newUsers}</div>
                  <div className={styles.activityLabel}>{t('newUsers', { defaultValue: 'New Users' })}</div>
                </div>
              </div>
              <div className={styles.activityStat}>
                <div className={styles.activityIcon}>
                  <Activity size={20} />
                </div>
                <div>
                  <div className={styles.activityValue}>{analyticsData?.userActivity.returningUsers}</div>
                  <div className={styles.activityLabel}>{t('returningUsers', { defaultValue: 'Returning Users' })}</div>
                </div>
              </div>
              <div className={styles.activityStat}>
                <div className={styles.activityIcon}>
                  <Clock size={20} />
                </div>
                <div>
                  <div className={styles.activityValue}>{analyticsData?.userActivity.avgSessionDuration}</div>
                  <div className={styles.activityLabel}>{t('avgSession', { defaultValue: 'Avg Session' })}</div>
                </div>
              </div>
              <div className={styles.activityStat}>
                <div className={styles.activityIcon}>
                  <Percent size={20} />
                </div>
                <div>
                  <div className={styles.activityValue}>{analyticsData?.userActivity.bounceRate}</div>
                  <div className={styles.activityLabel}>{t('bounceRate', { defaultValue: 'Bounce Rate' })}</div>
                </div>
              </div>
            </div>
          </div>

          {/* Revenue by Category */}
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>
              {t('revenueByCategory', { defaultValue: 'Revenue by Category' })}
            </h2>
            <div className={styles.categoryList}>
              {analyticsData?.revenueByCategory.map((item) => (
                <div key={item.category} className={styles.categoryItem}>
                  <div className={styles.categoryInfo}>
                    <div className={styles.categoryName}>{item.category}</div>
                    <div className={styles.categoryRevenue}>${item.revenue.toLocaleString()}</div>
                  </div>
                  <div className={styles.categoryBar}>
                    <div 
                      className={styles.categoryBarFill} 
                      style={{ width: `${item.percentage}%` }}
                    />
                  </div>
                  <div className={styles.categoryPercentage}>{item.percentage}%</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default PlatformAnalytics;
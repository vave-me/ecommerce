"use client";

import React, { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { 
  ArrowLeft, 
  BarChart3, 
  TrendingUp, 
  TrendingDown, 
  Eye, 
  MousePointer, 
  ShoppingCart, 
  DollarSign,
  Calendar,
  Download,
  RefreshCw,
  Filter,
  Search,
  Users,
  Package,
  AlertTriangle,
  CheckCircle,
  Info,
  ExternalLink
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { useQuery } from '@tanstack/react-query';
import { 
  getMerchantReports, 
  getMerchantPerformanceData,
  exportMerchantReport,
  getMerchantInsights
} from '@/api/merchantApi';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import { toast } from 'react-toastify';
import styles from './MerchantReports.module.css';

// Chart Component (simplified representation)
const PerformanceChart = ({ data, type, title }) => {
  if (!data || data.length === 0) {
    return (
      <div className={styles.chartEmpty}>
        <BarChart3 size={32} />
        <p>No data available</p>
      </div>
    );
  }

  return (
    <div className={styles.chartContainer}>
      <h4 className={styles.chartTitle}>{title}</h4>
      <div className={styles.chartPlaceholder}>
        <BarChart3 size={48} />
        <p>Chart visualization would be rendered here</p>
        <span className={styles.chartNote}>({data.length} data points)</span>
      </div>
    </div>
  );
};

// Metric Card Component
const MetricCard = ({ title, value, change, icon: Icon, color = 'primary' }) => {
  const isPositive = change > 0;
  const isNegative = change < 0;

  return (
    <div className={`${styles.metricCard} ${styles[color]}`}>
      <div className={styles.metricHeader}>
        <div className={styles.metricIcon}>
          <Icon size={20} />
        </div>
        <span className={styles.metricTitle}>{title}</span>
      </div>
      <div className={styles.metricValue}>{value}</div>
      {change !== undefined && (
        <div className={`${styles.metricChange} ${isPositive ? styles.positive : isNegative ? styles.negative : styles.neutral}`}>
          {isPositive ? <TrendingUp size={14} /> : isNegative ? <TrendingDown size={14} /> : null}
          <span>{change > 0 ? '+' : ''}{change}%</span>
        </div>
      )}
    </div>
  );
};

// Insights Card Component
const InsightCard = ({ insight }) => {
  const getInsightIcon = (type) => {
    switch (type) {
      case 'warning': return AlertTriangle;
      case 'success': return CheckCircle;
      case 'info': return Info;
      default: return Info;
    }
  };

  const Icon = getInsightIcon(insight.type);

  return (
    <div className={`${styles.insightCard} ${styles[insight.type]}`}>
      <div className={styles.insightIcon}>
        <Icon size={20} />
      </div>
      <div className={styles.insightContent}>
        <h4 className={styles.insightTitle}>{insight.title}</h4>
        <p className={styles.insightDescription}>{insight.description}</p>
        {insight.action && (
          <button className={styles.insightAction}>
            {insight.action}
            <ExternalLink size={14} />
          </button>
        )}
      </div>
    </div>
  );
};

const MerchantReports = () => {
  let t, tAnalytics;
  try {
    t = useTranslations('MerchantCenter');
    tAnalytics = useTranslations('MerchantCenter.analytics');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
    tAnalytics = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  // State
  const [dateRange, setDateRange] = useState('last_30_days');
  const [selectedMetrics, setSelectedMetrics] = useState(['impressions', 'clicks', 'conversions']);
  const [exportFormat, setExportFormat] = useState('csv');

  // Date range options
  const dateRangeOptions = [
    { value: 'last_7_days', label: 'Last 7 days' },
    { value: 'last_30_days', label: 'Last 30 days' },
    { value: 'last_90_days', label: 'Last 90 days' },
    { value: 'last_year', label: 'Last year' },
    { value: 'custom', label: 'Custom range' }
  ];

  // Fetch reports data
  const { data: reportsData, isLoading, error, refetch } = useQuery({
    queryKey: ['merchant-reports', dateRange],
    queryFn: () => getMerchantReports({ 
      date_range: dateRange,
      metrics: selectedMetrics 
    }),
    staleTime: 300000 // 5 minutes
  });

  // Fetch performance data for charts
  const { data: performanceData } = useQuery({
    queryKey: ['merchant-performance', dateRange],
    queryFn: () => getMerchantPerformanceData({ date_range: dateRange }),
    staleTime: 300000
  });

  // Fetch insights
  const { data: insights } = useQuery({
    queryKey: ['merchant-insights'],
    queryFn: getMerchantInsights,
    staleTime: 600000 // 10 minutes
  });

  // Access control
  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access merchant reports.' })}</p>
        </div>
      </div>
    );
  }

  // Event handlers
  const handleExportReport = async () => {
    try {
      const blob = await exportMerchantReport({
        date_range: dateRange,
        format: exportFormat,
        metrics: selectedMetrics
      });
      
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `merchant-report-${dateRange}.${exportFormat}`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      
      toast.success('Report exported successfully');
    } catch (error) {
      toast.error(`Export failed: ${error.message}`);
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <p>Loading merchant reports...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertTriangle size={48} />
          <h2>Error Loading Reports</h2>
          <p>{error.message}</p>
          <button onClick={() => refetch()} className={styles.retryButton}>
            <RefreshCw size={16} />
            Retry
          </button>
        </div>
      </div>
    );
  }

  const metrics = reportsData?.metrics || {};
  const trends = reportsData?.trends || {};

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <button 
              className={styles.backButton}
              onClick={() => router.push('/admin/merchant')}
            >
              <ArrowLeft size={16} />
              {t('backToMerchant', { defaultValue: 'Back to Merchant Center' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <BarChart3 size={24} />
                {tAnalytics('title', { defaultValue: 'Merchant Analytics' })}
              </h1>
              <p className={styles.subtitle}>
                {tAnalytics('subtitle', { defaultValue: 'Google Merchant Center analytics and performance insights' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <select
              value={dateRange}
              onChange={(e) => setDateRange(e.target.value)}
              className={styles.dateRangeSelect}
            >
              {dateRangeOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
            <button 
              onClick={handleExportReport}
              className={styles.exportButton}
            >
              <Download size={16} />
              Export
            </button>
            <button 
              onClick={() => refetch()}
              className={styles.iconButton}
              disabled={isLoading}
            >
              <RefreshCw size={16} className={isLoading ? styles.spinning : ''} />
            </button>
          </div>
        </div>

        {/* Key Metrics */}
        <div className={styles.metricsSection}>
          <h2 className={styles.sectionTitle}>{tAnalytics('metrics', { defaultValue: 'Key Metrics' })}</h2>
          <div className={styles.metricsGrid}>
            <MetricCard
              title="Total Impressions"
              value={metrics.impressions?.toLocaleString() || '0'}
              change={trends.impressions_change}
              icon={Eye}
              color="primary"
            />
            <MetricCard
              title="Total Clicks"
              value={metrics.clicks?.toLocaleString() || '0'}
              change={trends.clicks_change}
              icon={MousePointer}
              color="success"
            />
            <MetricCard
              title="Click-Through Rate"
              value={`${metrics.ctr || 0}%`}
              change={trends.ctr_change}
              icon={TrendingUp}
              color="info"
            />
            <MetricCard
              title="Conversions"
              value={metrics.conversions?.toLocaleString() || '0'}
              change={trends.conversions_change}
              icon={ShoppingCart}
              color="warning"
            />
            <MetricCard
              title="Revenue"
              value={`$${metrics.revenue?.toLocaleString() || '0'}`}
              change={trends.revenue_change}
              icon={DollarSign}
              color="success"
            />
            <MetricCard
              title="Cost Per Click"
              value={`$${metrics.cpc || '0.00'}`}
              change={trends.cpc_change}
              icon={DollarSign}
              color="secondary"
            />
            <MetricCard
              title="Active Products"
              value={metrics.active_products?.toLocaleString() || '0'}
              change={trends.products_change}
              icon={Package}
              color="info"
            />
            <MetricCard
              title="Conversion Rate"
              value={`${metrics.conversion_rate || 0}%`}
              change={trends.conversion_rate_change}
              icon={TrendingUp}
              color="primary"
            />
          </div>
        </div>

        {/* Performance Charts */}
        <div className={styles.chartsSection}>
          <h2 className={styles.sectionTitle}>{tAnalytics('trends', { defaultValue: 'Trends' })}</h2>
          <div className={styles.chartsGrid}>
            <div className={styles.chartCard}>
              <PerformanceChart
                data={performanceData?.impressions_over_time}
                type="line"
                title="Impressions Over Time"
              />
            </div>
            <div className={styles.chartCard}>
              <PerformanceChart
                data={performanceData?.clicks_over_time}
                type="line"
                title="Clicks Over Time"
              />
            </div>
            <div className={styles.chartCard}>
              <PerformanceChart
                data={performanceData?.revenue_over_time}
                type="bar"
                title="Revenue Over Time"
              />
            </div>
            <div className={styles.chartCard}>
              <PerformanceChart
                data={performanceData?.ctr_over_time}
                type="line"
                title="CTR Trend"
              />
            </div>
          </div>
        </div>

        {/* Insights Section */}
        {insights && insights.length > 0 && (
          <div className={styles.insightsSection}>
            <h2 className={styles.sectionTitle}>Insights & Recommendations</h2>
            <div className={styles.insightsGrid}>
              {insights.map((insight, index) => (
                <InsightCard key={index} insight={insight} />
              ))}
            </div>
          </div>
        )}

        {/* Top Products Table */}
        {reportsData?.top_products && (
          <div className={styles.topProductsSection}>
            <h2 className={styles.sectionTitle}>{tAnalytics('topProducts', { defaultValue: 'Top Performing Products' })}</h2>
            <div className={styles.tableContainer}>
              <table className={styles.table}>
                <thead>
                  <tr>
                    <th>Product</th>
                    <th>Impressions</th>
                    <th>Clicks</th>
                    <th>CTR</th>
                    <th>Conversions</th>
                    <th>Revenue</th>
                  </tr>
                </thead>
                <tbody>
                  {reportsData.top_products.map((product, index) => (
                    <tr key={product.id} className={styles.tableRow}>
                      <td>
                        <div className={styles.productCell}>
                          {product.image_url && (
                            <img 
                              src={product.image_url} 
                              alt={product.title}
                              className={styles.productImage}
                            />
                          )}
                          <div className={styles.productInfo}>
                            <div className={styles.productTitle}>{product.title}</div>
                            <div className={styles.productId}>ID: {product.id}</div>
                          </div>
                        </div>
                      </td>
                      <td className={styles.numberCell}>{product.impressions?.toLocaleString()}</td>
                      <td className={styles.numberCell}>{product.clicks?.toLocaleString()}</td>
                      <td className={styles.numberCell}>{product.ctr}%</td>
                      <td className={styles.numberCell}>{product.conversions?.toLocaleString()}</td>
                      <td className={styles.numberCell}>${product.revenue?.toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Export Options */}
        <div className={styles.exportSection}>
          <h2 className={styles.sectionTitle}>Export Options</h2>
          <div className={styles.exportCard}>
            <div className={styles.exportOptions}>
              <div className={styles.exportOption}>
                <label>Format:</label>
                <select
                  value={exportFormat}
                  onChange={(e) => setExportFormat(e.target.value)}
                  className={styles.exportSelect}
                >
                  <option value="csv">CSV</option>
                  <option value="excel">Excel</option>
                  <option value="pdf">PDF</option>
                </select>
              </div>
              <div className={styles.exportOption}>
                <label>Date Range:</label>
                <select
                  value={dateRange}
                  onChange={(e) => setDateRange(e.target.value)}
                  className={styles.exportSelect}
                >
                  {dateRangeOptions.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <button 
              onClick={handleExportReport}
              className={styles.exportMainButton}
            >
              <Download size={16} />
              Export Detailed Report
            </button>
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default MerchantReports; 
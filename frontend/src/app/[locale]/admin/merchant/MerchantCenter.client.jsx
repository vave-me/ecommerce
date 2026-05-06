"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Globe,
  CheckCircle,
  XCircle,
  Clock,
  AlertTriangle,
  RefreshCw,
  Upload,
  Download,
  Search,
  Filter,
  MoreVertical,
  ExternalLink,
  Package,
  AlertCircle,
  PlayCircle,
  StopCircle,
  Settings,
  BarChart3,
  FileText,
  ArrowLeft,
  Calendar,
  Target,
  Zap,
  Database
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listProducts,
  syncProduct,
  batchSyncProducts,
  removeProduct,
  getProductStatus,
  getMerchantOverview,
  getSyncErrors,
  getMerchantMetrics,
  retryFailedSyncs
} from '@/api/merchantApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './MerchantCenter.module.css';

const StatusBadge = ({ status }) => {
  const getStatusConfig = () => {
    switch (status) {
      case 'SYNCED':
        return { icon: CheckCircle, class: styles.statusSynced, label: 'Synced' };
      case 'PENDING':
        return { icon: Clock, class: styles.statusPending, label: 'Pending' };
      case 'FAILED':
        return { icon: XCircle, class: styles.statusFailed, label: 'Failed' };
      case 'REMOVED':
        return { icon: AlertTriangle, class: styles.statusRemoved, label: 'Removed' };
      default:
        return { icon: AlertCircle, class: styles.statusUnknown, label: 'Unknown' };
    }
  };

  const { icon: Icon, class: statusClass, label } = getStatusConfig();

  return (
    <span className={`${styles.statusBadge} ${statusClass}`}>
      <Icon size={12} />
      {label}
    </span>
  );
};

const SyncHistoryTable = ({ syncs, onProductAction }) => {
  const formatDate = (dateString) => {
    if (!dateString) return 'Not synced';
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className={styles.tableContainer}>
      <table className={styles.syncTable}>
        <thead>
          <tr>
            <th>Product</th>
            <th>Status</th>
            <th>Last Sync</th>
            <th>Merchant ID</th>
            <th>Issues</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {syncs.map(sync => (
            <tr key={sync.id} className={styles.syncRow}>
              <td className={styles.productCell}>
                <div className={styles.productInfo}>
                  <Package size={16} />
                  <div>
                    <div className={styles.productName}>{sync.name || sync.productName}</div>
                    <div className={styles.productId}>ID: {sync.id || sync.productId}</div>
                  </div>
                </div>
              </td>
              <td>
                <StatusBadge status={sync.status || 'UNKNOWN'} />
              </td>
              <td className={styles.dateCell}>
                {formatDate(sync.syncTime || sync.lastSyncedAt)}
              </td>
              <td className={styles.merchantIdCell}>
                {sync.merchantProductId || '—'}
              </td>
              <td className={styles.issuesCell}>
                {sync.errors?.length > 0 && (
                  <span className={styles.errorCount}>
                    {sync.errors.length} error{sync.errors.length !== 1 ? 's' : ''}
                  </span>
                )}
                {sync.warnings?.length > 0 && (
                  <span className={styles.warningCount}>
                    {sync.warnings.length} warning{sync.warnings.length !== 1 ? 's' : ''}
                  </span>
                )}
                {(!sync.errors || sync.errors.length === 0) && (!sync.warnings || sync.warnings.length === 0) && '—'}
              </td>
              <td className={styles.actionCell}>
                <button 
                  className={styles.actionButton}
                  onClick={() => onProductAction('view', sync)}
                >
                  <MoreVertical size={14} />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

const ErrorsList = ({ errors }) => {
  return (
    <div className={styles.errorsList}>
      {errors.map(error => (
        <div key={error.id} className={styles.errorItem}>
          <div className={styles.errorIcon}>
            <XCircle size={16} />
          </div>
          <div className={styles.errorContent}>
            <div className={styles.errorMessage}>{error.message}</div>
            <div className={styles.errorDetails}>
              <span className={styles.errorProduct}>{error.productName}</span>
              <span className={styles.errorOccurrence}>Occurred {error.occurrence} times</span>
              <span className={styles.errorTime}>Last seen: {new Date(error.lastSeen).toLocaleString()}</span>
            </div>
          </div>
          <button className={styles.errorAction}>
            <ExternalLink size={14} />
          </button>
        </div>
      ))}
    </div>
  );
};

const MerchantCenter = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('MerchantCenter');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [syncing, setSyncing] = useState(false);

  // Fetch merchant center data with React Query
  const { data: merchantOverview, isLoading: overviewLoading, refetch: refetchOverview } = useQuery({
    queryKey: ['merchantCenter', 'overview'],
    queryFn: getMerchantOverview,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 2,
    enabled: isAdmin
  });

  const { data: syncedProducts, isLoading: productsLoading, refetch: refetchProducts } = useQuery({
    queryKey: ['merchantCenter', 'products'],
    queryFn: () => listProducts({ pageSize: 50 }),
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: 2,
    enabled: isAdmin
  });

  const { data: syncErrors, isLoading: errorsLoading, refetch: refetchErrors } = useQuery({
    queryKey: ['merchantCenter', 'errors'],
    queryFn: () => getSyncErrors({ limit: 10 }),
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: 2,
    enabled: isAdmin
  });

  const { data: merchantMetrics, isLoading: metricsLoading, refetch: refetchMetrics } = useQuery({
    queryKey: ['merchantCenter', 'metrics'],
    queryFn: getMerchantMetrics,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 2,
    enabled: isAdmin
  });

  // Mutations for sync operations
  const syncMutation = useMutation({
    mutationFn: batchSyncProducts,
    onSuccess: () => {
      queryClient.invalidateQueries(['merchantCenter']);
      setSyncing(false);
    },
    onError: (error) => {
      // Error: 'Sync failed:', error...
      setSyncing(false);
    }
  });

  const retryMutation = useMutation({
    mutationFn: retryFailedSyncs,
    onSuccess: () => {
      queryClient.invalidateQueries(['merchantCenter']);
    },
    onError: (error) => {
      // Error: 'Retry failed:', error...
    }
  });

  const handleSync = useCallback(async (type = 'all') => {
    setSyncing(true);
    try {
      if (type === 'failed') {
        // Get failed product IDs and retry them
        const failedProducts = syncedProducts?.products?.filter(p => p.status === 'FAILED') || [];
        const productIds = failedProducts.map(p => p.id);
        if (productIds.length > 0) {
          await retryMutation.mutateAsync(productIds);
        }
      } else {
        // Sync all products
        const allProducts = syncedProducts?.products || [];
        const products = allProducts.map(p => ({ productId: p.id, product: p }));
        if (products.length > 0) {
          await syncMutation.mutateAsync({ products });
        }
      }
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, [syncedProducts, syncMutation, retryMutation]);

  const handleProductAction = useCallback(async (action, product) => {
    try {
      switch (action) {
        case 'view':
          
          break;
        case 'retry':
          await syncProduct({ productId: product.productId, product: product });
          queryClient.invalidateQueries(['merchantCenter']);
          break;
        case 'remove':
          await removeProduct(product.productId);
          queryClient.invalidateQueries(['merchantCenter']);
          break;
      }
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, [queryClient]);

  const isLoading = overviewLoading || productsLoading || errorsLoading || metricsLoading;

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access merchant center management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (isLoading && !merchantOverview) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Merchant Center...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch the latest sync data.' })}</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <button 
              className={styles.backButton}
              onClick={() => router.back()}
            >
              <ArrowLeft size={16} />
              {t('backToDashboard', { defaultValue: 'Back to Dashboard' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <Globe size={24} />
                {t('title', { defaultValue: 'Google Merchant Center' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Manage product synchronization with Google Merchant Center' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.settingsButton}>
              <Settings size={16} />
              {t('settings', { defaultValue: 'Settings' })}
            </button>
            <button 
              className={styles.syncButton}
              onClick={() => handleSync('all')}
              disabled={syncing}
            >
              <RefreshCw size={16} className={syncing ? styles.spinning : ''} />
              {syncing ? t('syncing', { defaultValue: 'Syncing...' }) : t('syncAll', { defaultValue: 'Sync All' })}
            </button>
          </div>
        </div>

        {/* Overview Stats */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{(merchantOverview?.syncedProducts || 0).toLocaleString()}</div>
              <div className={styles.statLabel}>{t('syncedProducts', { defaultValue: 'Synced Products' })}</div>
              <div className={styles.statSubtext}>{t('successfullySync', { defaultValue: 'Successfully synchronized' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Clock size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{(merchantOverview?.pendingSync || 0).toLocaleString()}</div>
              <div className={styles.statLabel}>{t('pendingSync', { defaultValue: 'Pending Sync' })}</div>
              <div className={styles.statSubtext}>{t('waitingInQueue', { defaultValue: 'Waiting in queue' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <XCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{(merchantOverview?.failedSync || 0).toLocaleString()}</div>
              <div className={styles.statLabel}>{t('failedSync', { defaultValue: 'Failed Sync' })}</div>
              <div className={styles.statSubtext}>{t('requiresAttention', { defaultValue: 'Requires attention' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <BarChart3 size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{merchantOverview?.successRate || 0}%</div>
              <div className={styles.statLabel}>{t('successRate', { defaultValue: 'Success Rate' })}</div>
              <div className={styles.statSubtext}>{t('last30Days', { defaultValue: 'Last 30 days' })}</div>
            </div>
          </div>
        </div>

        {/* Sync Status Banner */}
        <div className={styles.statusBanner}>
          <div className={styles.statusInfo}>
            <div className={styles.statusIndicator}>
              <div className={styles.liveDot}></div>
              <span className={styles.statusText}>
                {t('syncStatus', { defaultValue: 'Sync Status' })}: <strong>{merchantOverview?.syncStatus || t('unknown', { defaultValue: 'Unknown' })}</strong>
              </span>
            </div>
            <div className={styles.lastSync}>
              {t('lastSync', { defaultValue: 'Last sync' })}: {merchantOverview?.lastSyncTime ? new Date(merchantOverview.lastSyncTime).toLocaleString() : t('never', { defaultValue: 'Never' })}
            </div>
          </div>
          <div className={styles.quickActions}>
            <button 
              className={styles.quickAction}
              onClick={() => handleSync('failed')}
            >
              <RefreshCw size={14} />
              {t('retryFailed', { defaultValue: 'Retry Failed' })}
            </button>
            <button className={styles.quickAction}>
              <Upload size={14} />
              {t('bulkUpload', { defaultValue: 'Bulk Upload' })}
            </button>
            <button className={styles.quickAction}>
              <Download size={14} />
              {t('exportReport', { defaultValue: 'Export Report' })}
            </button>
          </div>
        </div>

        {/* Content Grid */}
        <div className={styles.contentGrid}>
          {/* Recent Syncs */}
          <div className={styles.contentCard}>
            <div className={styles.cardHeader}>
              <h3 className={styles.cardTitle}>
                <RefreshCw size={18} />
                {t('recentSynchronizations', { defaultValue: 'Recent Synchronizations' })}
              </h3>
              <div className={styles.cardControls}>
                <div className={styles.searchBox}>
                  <Search size={14} />
                  <input
                    type="text"
                    placeholder={t('searchProducts', { defaultValue: 'Search products...' })}
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className={styles.searchInput}
                  />
                </div>
                <select
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className={styles.filterSelect}
                >
                  <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
                  <option value="SYNCED">{t('synced', { defaultValue: 'Synced' })}</option>
                  <option value="PENDING">{t('pending', { defaultValue: 'Pending' })}</option>
                  <option value="FAILED">{t('failed', { defaultValue: 'Failed' })}</option>
                </select>
              </div>
            </div>
            <SyncHistoryTable 
              syncs={syncedProducts?.products || []} 
              onProductAction={handleProductAction} 
            />
          </div>

          {/* Sync Errors */}
          <div className={styles.contentCard}>
            <div className={styles.cardHeader}>
              <h3 className={styles.cardTitle}>
                <AlertTriangle size={18} />
                {t('syncErrors', { defaultValue: 'Sync Errors' })}
              </h3>
              <button className={styles.cardAction}>{t('viewAll', { defaultValue: 'View All' })}</button>
            </div>
            <ErrorsList errors={syncErrors?.errors || []} />
          </div>
        </div>

        {/* Performance Metrics */}
        <div className={styles.metricsGrid}>
          <div className={styles.metricCard}>
            <div className={styles.metricHeader}>
              <Zap size={16} />
              <span>{t('avgSyncTime', { defaultValue: 'Average Sync Time' })}</span>
            </div>
            <div className={styles.metricValue}>{merchantMetrics?.avgSyncTime || 0}ms</div>
          </div>
          
          <div className={styles.metricCard}>
            <div className={styles.metricHeader}>
              <AlertTriangle size={16} />
              <span>{t('totalErrors', { defaultValue: 'Total Errors' })}</span>
            </div>
            <div className={styles.metricValue}>{merchantMetrics?.totalErrors || 0}</div>
          </div>
          
          <div className={styles.metricCard}>
            <div className={styles.metricHeader}>
              <AlertCircle size={16} />
              <span>{t('warnings', { defaultValue: 'Warnings' })}</span>
            </div>
            <div className={styles.metricValue}>{merchantMetrics?.warningsCount || 0}</div>
          </div>
          
          <div className={styles.metricCard}>
            <div className={styles.metricHeader}>
              <Database size={16} />
              <span>{t('apiCallsToday', { defaultValue: 'API Calls Today' })}</span>
            </div>
            <div className={styles.metricValue}>{merchantMetrics?.apiCallsToday?.toLocaleString() || '0'}</div>
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default MerchantCenter; 
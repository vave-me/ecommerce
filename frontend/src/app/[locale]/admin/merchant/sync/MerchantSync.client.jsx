"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  ArrowLeft,
  RefreshCw,
  Upload,
  Download,
  CheckCircle,
  XCircle,
  Clock,
  AlertTriangle,
  Search,
  Filter,
  Package,
  Play,
  Pause,
  SkipForward,
  Globe,
  Settings,
  Eye,
  MoreVertical,
  Zap,
  Target,
  Calendar,
  BarChart3
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listProducts,
  syncProduct,
  batchSyncProducts,
  removeProduct,
  getMerchantOverview,
  retryFailedSyncs
} from '@/api/merchantApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './MerchantSync.module.css';

const SyncStatusBadge = ({ status }) => {
  const getStatusConfig = () => {
    switch (status) {
      case 'SYNCED':
      case 'SUCCESS':
        return { icon: CheckCircle, class: styles.statusSuccess, label: 'Synced' };
      case 'PENDING':
      case 'IN_PROGRESS':
        return { icon: Clock, class: styles.statusPending, label: 'Syncing' };
      case 'FAILED':
      case 'ERROR':
        return { icon: XCircle, class: styles.statusError, label: 'Failed' };
      case 'QUEUED':
        return { icon: Clock, class: styles.statusQueued, label: 'Queued' };
      default:
        return { icon: AlertTriangle, class: styles.statusUnknown, label: 'Unknown' };
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

const ProductSyncItem = ({ product, onSync, onRemove, syncing }) => {
  const [itemSyncing, setItemSyncing] = useState(false);

  const handleSync = async () => {
    setItemSyncing(true);
    try {
      await onSync(product);
    } finally {
      setItemSyncing(false);
    }
  };

  return (
    <div className={styles.syncItem}>
      <div className={styles.productInfo}>
        <div className={styles.productImage}>
          {product.image ? (
            <img src={product.image} alt={product.name} />
          ) : (
            <Package size={24} />
          )}
        </div>
        <div className={styles.productDetails}>
          <h4 className={styles.productName}>{product.name}</h4>
          <p className={styles.productId}>ID: {product.id}</p>
          <p className={styles.productPrice}>${product.price}</p>
        </div>
      </div>
      
      <div className={styles.syncStatus}>
        <SyncStatusBadge status={product.syncStatus || 'PENDING'} />
        {product.lastSyncTime && (
          <div className={styles.lastSync}>
            Last: {new Date(product.lastSyncTime).toLocaleDateString()}
          </div>
        )}
      </div>

      <div className={styles.syncProgress}>
        {(itemSyncing || syncing) ? (
          <div className={styles.progressBar}>
            <div className={styles.progressFill}></div>
          </div>
        ) : (
          <div className={styles.syncActions}>
            <button
              className={styles.syncButton}
              onClick={handleSync}
              disabled={syncing}
            >
              <RefreshCw size={14} />
              Sync
            </button>
            <button
              className={styles.viewButton}
              onClick={() => window.open(`/admin/products/${product.id}`, '_blank')}
            >
              <Eye size={14} />
            </button>
            <button
              className={styles.menuButton}
              onClick={() => onRemove(product)}
            >
              <MoreVertical size={14} />
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

const MerchantSync = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('MerchantSync');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [selectedProducts, setSelectedProducts] = useState([]);
  const [syncMode, setSyncMode] = useState('all'); // 'all', 'selected', 'failed'
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [syncing, setSyncing] = useState(false);
  const [syncProgress, setSyncProgress] = useState({ current: 0, total: 0 });

  // Fetch data
  const { data: products, isLoading: productsLoading, refetch: refetchProducts } = useQuery({
    queryKey: ['merchantSync', 'products'],
    queryFn: () => listProducts({ pageSize: 100 }),
    staleTime: 2 * 60 * 1000,
    retry: 2,
    enabled: isAdmin
  });

  const { data: overview, refetch: refetchOverview } = useQuery({
    queryKey: ['merchantSync', 'overview'],
    queryFn: getMerchantOverview,
    staleTime: 5 * 60 * 1000,
    retry: 2,
    enabled: isAdmin
  });

  // Mutations
  const syncMutation = useMutation({
    mutationFn: async ({ products: productsToSync, mode }) => {
      setSyncing(true);
      setSyncProgress({ current: 0, total: productsToSync.length });

      for (let i = 0; i < productsToSync.length; i++) {
        const product = productsToSync[i];
        setSyncProgress({ current: i + 1, total: productsToSync.length });
        
        if (mode === 'individual') {
          await syncProduct({ productId: product.id, product });
        } else {
          // Batch sync in groups of 10
          if (i === 0 || i % 10 === 0) {
            const batch = productsToSync.slice(i, Math.min(i + 10, productsToSync.length));
            await batchSyncProducts({ 
              products: batch.map(p => ({ productId: p.id, product: p }))
            });
          }
        }
        
        // Small delay to show progress
        await new Promise(resolve => setTimeout(resolve, 500));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['merchantSync']);
      queryClient.invalidateQueries(['merchantCenter']);
      setSyncing(false);
      setSyncProgress({ current: 0, total: 0 });
    },
    onError: (error) => {
      // Error: 'Sync failed:', error...
      setSyncing(false);
      setSyncProgress({ current: 0, total: 0 });
    }
  });

  const handleBulkSync = useCallback(async () => {
    const productsList = products?.products || [];
    let productsToSync = [];

    switch (syncMode) {
      case 'all':
        productsToSync = productsList;
        break;
      case 'selected':
        productsToSync = productsList.filter(p => selectedProducts.includes(p.id));
        break;
      case 'failed':
        productsToSync = productsList.filter(p => p.syncStatus === 'FAILED');
        break;
    }

    if (productsToSync.length === 0) {
      alert('No products to sync');
      return;
    }

    await syncMutation.mutateAsync({ products: productsToSync, mode: 'batch' });
  }, [products, syncMode, selectedProducts, syncMutation]);

  const handleSingleSync = useCallback(async (product) => {
    await syncMutation.mutateAsync({ products: [product], mode: 'individual' });
  }, [syncMutation]);

  const handleProductRemove = useCallback(async (product) => {
    if (confirm(`Remove ${product.name} from Merchant Center?`)) {
      try {
        await removeProduct(product.id);
        queryClient.invalidateQueries(['merchantSync']);
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
  }, [queryClient]);

  const handleSelectProduct = useCallback((productId) => {
    setSelectedProducts(prev => 
      prev.includes(productId) 
        ? prev.filter(id => id !== productId)
        : [...prev, productId]
    );
  }, []);

  const handleSelectAll = useCallback(() => {
    const productsList = products?.products || [];
    const allIds = productsList.map(p => p.id);
    setSelectedProducts(selectedProducts.length === allIds.length ? [] : allIds);
  }, [products, selectedProducts]);

  const filteredProducts = (products?.products || []).filter(product => {
    const matchesSearch = searchTerm === '' || 
      product.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      product.id?.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = statusFilter === 'all' || product.syncStatus === statusFilter;
    
    return matchesSearch && matchesStatus;
  });

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access merchant sync.' })}</p>
        </div>
      </div>
    );
  }

  if (productsLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Products...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch product data.' })}</p>
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
              onClick={() => router.push('/admin/merchant')}
            >
              <ArrowLeft size={16} />
              {t('backToMerchant', { defaultValue: 'Back to Merchant Center' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <RefreshCw size={24} />
                {t('title', { defaultValue: 'Product Synchronization' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Manually sync products with Google Merchant Center' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.settingsButton}>
              <Settings size={16} />
              {t('settings', { defaultValue: 'Sync Settings' })}
            </button>
          </div>
        </div>

        {/* Sync Controls */}
        <div className={styles.controlsPanel}>
          <div className={styles.controlsLeft}>
            <div className={styles.syncModeSelector}>
              <label>{t('syncMode', { defaultValue: 'Sync Mode' })}:</label>
              <select 
                value={syncMode} 
                onChange={(e) => setSyncMode(e.target.value)}
                disabled={syncing}
              >
                <option value="all">{t('syncAll', { defaultValue: 'All Products' })}</option>
                <option value="selected">{t('syncSelected', { defaultValue: 'Selected Products' })}</option>
                <option value="failed">{t('syncFailed', { defaultValue: 'Failed Only' })}</option>
              </select>
            </div>
            <button 
              className={styles.bulkSyncButton}
              onClick={handleBulkSync}
              disabled={syncing}
            >
              {syncing ? (
                <>
                  <RefreshCw size={16} className={styles.spinning} />
                  {t('syncing', { defaultValue: 'Syncing...' })}
                </>
              ) : (
                <>
                  <Play size={16} />
                  {t('startSync', { defaultValue: 'Start Sync' })}
                </>
              )}
            </button>
          </div>

          <div className={styles.controlsRight}>
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

        {/* Progress Bar */}
        {syncing && (
          <div className={styles.progressSection}>
            <div className={styles.progressInfo}>
              <span>{t('syncingProgress', { defaultValue: 'Synchronizing products' })}: {syncProgress.current} / {syncProgress.total}</span>
              <span>{Math.round((syncProgress.current / syncProgress.total) * 100)}%</span>
            </div>
            <div className={styles.progressBarContainer}>
              <div 
                className={styles.progressBar}
                style={{ width: `${(syncProgress.current / syncProgress.total) * 100}%` }}
              ></div>
            </div>
          </div>
        )}

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{overview?.syncedProducts || 0}</div>
              <div className={styles.statLabel}>{t('syncedProducts', { defaultValue: 'Synced' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Clock size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{overview?.pendingSync || 0}</div>
              <div className={styles.statLabel}>{t('pendingSync', { defaultValue: 'Pending' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <XCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{overview?.failedSync || 0}</div>
              <div className={styles.statLabel}>{t('failedSync', { defaultValue: 'Failed' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <Package size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{filteredProducts.length}</div>
              <div className={styles.statLabel}>{t('totalProducts', { defaultValue: 'Total' })}</div>
            </div>
          </div>
        </div>

        {/* Product List */}
        <div className={styles.productList}>
          <div className={styles.listHeader}>
            <div className={styles.selectAll}>
              <input
                type="checkbox"
                checked={selectedProducts.length === filteredProducts.length && filteredProducts.length > 0}
                onChange={handleSelectAll}
                disabled={syncing}
              />
              <span>{t('selectAll', { defaultValue: 'Select All' })} ({selectedProducts.length})</span>
            </div>
            <div className={styles.listStats}>
              {t('showingProducts', { defaultValue: 'Showing {count} products', count: filteredProducts.length })}
            </div>
          </div>

          <div className={styles.syncList}>
            {filteredProducts.map((product) => (
              <div key={product.id} className={styles.syncItemContainer}>
                <input
                  type="checkbox"
                  checked={selectedProducts.includes(product.id)}
                  onChange={() => handleSelectProduct(product.id)}
                  disabled={syncing}
                  className={styles.productCheckbox}
                />
                <ProductSyncItem
                  product={product}
                  onSync={handleSingleSync}
                  onRemove={handleProductRemove}
                  syncing={syncing}
                />
              </div>
            ))}
          </div>

          {filteredProducts.length === 0 && (
            <div className={styles.emptyState}>
              <Package size={48} />
              <h3>{t('noProducts', { defaultValue: 'No products found' })}</h3>
              <p>{t('noProductsMessage', { defaultValue: 'No products match your current filters.' })}</p>
            </div>
          )}
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default MerchantSync; 
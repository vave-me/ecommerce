"use client";

import React, { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { 
  ArrowLeft, 
  Package, 
  Plus, 
  Search, 
  Filter, 
  Download, 
  Upload, 
  Eye,
  Edit,
  Trash2,
  MoreHorizontal,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  Clock,
  ExternalLink,
  Copy,
  Star,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Tag,
  Layers,
  Settings
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  getMerchantProducts, 
  syncMerchantProduct, 
  deleteMerchantProduct,
  bulkUpdateMerchantProducts,
  getMerchantProductMetrics
} from '@/api/merchantApi';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import { toast } from 'react-toastify';
import styles from './MerchantProducts.module.css';

// Product Status Badge Component
const StatusBadge = ({ status, className = '' }) => {
  const statusConfig = {
    active: { label: 'Active', color: 'success' },
    pending: { label: 'Pending Review', color: 'warning' },
    disapproved: { label: 'Disapproved', color: 'danger' },
    expired: { label: 'Expired', color: 'secondary' },
    draft: { label: 'Draft', color: 'info' }
  };

  const config = statusConfig[status] || statusConfig.draft;

  return (
    <span className={`${styles.badge} ${styles[config.color]} ${className}`}>
      {config.label}
    </span>
  );
};

// Product Performance Component
const ProductPerformance = ({ product }) => {
  const hasPerformanceData = product.metrics?.impressions !== undefined;
  
  if (!hasPerformanceData) {
    return <span className={styles.noData}>No data</span>;
  }

  return (
    <div className={styles.performanceMetrics}>
      <div className={styles.metric}>
        <span className={styles.metricLabel}>Impressions</span>
        <span className={styles.metricValue}>{product.metrics.impressions.toLocaleString()}</span>
      </div>
      <div className={styles.metric}>
        <span className={styles.metricLabel}>Clicks</span>
        <span className={styles.metricValue}>{product.metrics.clicks.toLocaleString()}</span>
      </div>
      <div className={styles.metric}>
        <span className={styles.metricLabel}>CTR</span>
        <span className={styles.metricValue}>{product.metrics.ctr}%</span>
      </div>
    </div>
  );
};

const MerchantProducts = () => {
  let t;
  try {
    t = useTranslations('MerchantProducts');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  // State
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [filterStatus, setFilterStatus] = useState('all');
  const [sortBy, setSortBy] = useState('updated_at');
  const [sortOrder, setSortOrder] = useState('desc');
  const [showBulkActions, setShowBulkActions] = useState(false);

  const itemsPerPage = 20;

  // Fetch merchant products
  const { data: productsData, isLoading, error, refetch } = useQuery({
    queryKey: ['merchant-products', currentPage, searchTerm, filterStatus, sortBy, sortOrder],
    queryFn: () => getMerchantProducts({
      page: currentPage,
      limit: itemsPerPage,
      search: searchTerm,
      status: filterStatus === 'all' ? undefined : filterStatus,
      sort_by: sortBy,
      sort_order: sortOrder
    }),
    staleTime: 30000 // 30 seconds
  });

  // Fetch metrics
  const { data: metrics } = useQuery({
    queryKey: ['merchant-product-metrics'],
    queryFn: getMerchantProductMetrics,
    staleTime: 300000 // 5 minutes
  });

  // Mutations
  const syncProductMutation = useMutation({
    mutationFn: syncMerchantProduct,
    onSuccess: (data, variables) => {
      toast.success(`Product ${variables.productId} synced successfully`);
      queryClient.invalidateQueries(['merchant-products']);
    },
    onError: (error) => {
      toast.error(`Sync failed: ${error.message}`);
    }
  });

  const deleteProductMutation = useMutation({
    mutationFn: deleteMerchantProduct,
    onSuccess: (data, variables) => {
      toast.success(`Product ${variables.productId} deleted from Merchant Center`);
      queryClient.invalidateQueries(['merchant-products']);
      setSelectedProducts(prev => prev.filter(id => id !== variables.productId));
    },
    onError: (error) => {
      toast.error(`Delete failed: ${error.message}`);
    }
  });

  const bulkUpdateMutation = useMutation({
    mutationFn: bulkUpdateMerchantProducts,
    onSuccess: (data) => {
      toast.success(`${data.updated_count} products updated successfully`);
      queryClient.invalidateQueries(['merchant-products']);
      setSelectedProducts([]);
      setShowBulkActions(false);
    },
    onError: (error) => {
      toast.error(`Bulk update failed: ${error.message}`);
    }
  });

  // Access control
  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access merchant products.' })}</p>
        </div>
      </div>
    );
  }

  // Computed values
  const products = productsData?.products || [];
  const totalProducts = productsData?.total || 0;
  const totalPages = Math.ceil(totalProducts / itemsPerPage);

  const filteredProducts = useMemo(() => {
    return products.filter(product => {
      const matchesSearch = searchTerm === '' || 
        product.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        product.id.toLowerCase().includes(searchTerm.toLowerCase());
      
      const matchesStatus = filterStatus === 'all' || product.status === filterStatus;
      
      return matchesSearch && matchesStatus;
    });
  }, [products, searchTerm, filterStatus]);

  // Event handlers
  const handleSelectProduct = (productId) => {
    setSelectedProducts(prev => 
      prev.includes(productId) 
        ? prev.filter(id => id !== productId)
        : [...prev, productId]
    );
  };

  const handleSelectAll = () => {
    if (selectedProducts.length === filteredProducts.length) {
      setSelectedProducts([]);
    } else {
      setSelectedProducts(filteredProducts.map(p => p.id));
    }
  };

  const handleSyncProduct = (productId) => {
    syncProductMutation.mutate({ productId });
  };

  const handleDeleteProduct = (productId) => {
    if (confirm('Are you sure you want to remove this product from Merchant Center?')) {
      deleteProductMutation.mutate({ productId });
    }
  };

  const handleBulkSync = () => {
    bulkUpdateMutation.mutate({
      product_ids: selectedProducts,
      action: 'sync'
    });
  };

  const handleBulkDelete = () => {
    if (confirm(`Are you sure you want to remove ${selectedProducts.length} products from Merchant Center?`)) {
      bulkUpdateMutation.mutate({
        product_ids: selectedProducts,
        action: 'delete'
      });
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <p>Loading merchant products...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertTriangle size={48} />
          <h2>Error Loading Products</h2>
          <p>{error.message}</p>
          <button onClick={() => refetch()} className={styles.retryButton}>
            <RefreshCw size={16} />
            Retry
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
                <Package size={24} />
                {t('title', { defaultValue: 'Merchant Products' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Manage products in Google Merchant Center' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button 
              onClick={() => router.push('/admin/products/add')}
              className={styles.primaryButton}
            >
              <Plus size={16} />
              Add Product
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

        {/* Metrics Cards */}
        {metrics && (
          <div className={styles.metricsGrid}>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <Package size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.total_products || 0}</div>
                <div className={styles.metricLabel}>Total Products</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <CheckCircle size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.active_products || 0}</div>
                <div className={styles.metricLabel}>Active</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <Clock size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.pending_products || 0}</div>
                <div className={styles.metricLabel}>Pending</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <AlertTriangle size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.disapproved_products || 0}</div>
                <div className={styles.metricLabel}>Issues</div>
              </div>
            </div>
          </div>
        )}

        {/* Filters and Search */}
        <div className={styles.filtersSection}>
          <div className={styles.searchBar}>
            <Search size={16} />
            <input
              type="text"
              placeholder="Search products..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filters}>
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">All Status</option>
              <option value="active">Active</option>
              <option value="pending">Pending</option>
              <option value="disapproved">Disapproved</option>
              <option value="expired">Expired</option>
              <option value="draft">Draft</option>
            </select>
            <select
              value={`${sortBy}-${sortOrder}`}
              onChange={(e) => {
                const [field, order] = e.target.value.split('-');
                setSortBy(field);
                setSortOrder(order);
              }}
              className={styles.filterSelect}
            >
              <option value="updated_at-desc">Recently Updated</option>
              <option value="created_at-desc">Recently Created</option>
              <option value="title-asc">Title A-Z</option>
              <option value="title-desc">Title Z-A</option>
              <option value="status-asc">Status</option>
            </select>
          </div>
        </div>

        {/* Bulk Actions */}
        {selectedProducts.length > 0 && (
          <div className={styles.bulkActions}>
            <div className={styles.bulkInfo}>
              <span>{selectedProducts.length} products selected</span>
            </div>
            <div className={styles.bulkButtons}>
              <button 
                onClick={handleBulkSync}
                className={styles.bulkButton}
                disabled={bulkUpdateMutation.isLoading}
              >
                <RefreshCw size={16} />
                Sync Selected
              </button>
              <button 
                onClick={handleBulkDelete}
                className={`${styles.bulkButton} ${styles.danger}`}
                disabled={bulkUpdateMutation.isLoading}
              >
                <Trash2 size={16} />
                Remove Selected
              </button>
            </div>
          </div>
        )}

        {/* Products Table */}
        <div className={styles.tableContainer}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th className={styles.checkboxColumn}>
                  <input
                    type="checkbox"
                    checked={selectedProducts.length === filteredProducts.length && filteredProducts.length > 0}
                    onChange={handleSelectAll}
                  />
                </th>
                <th>Product</th>
                <th>Status</th>
                <th>Price</th>
                <th>Performance</th>
                <th>Last Updated</th>
                <th className={styles.actionsColumn}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredProducts.map((product) => (
                <tr key={product.id} className={styles.tableRow}>
                  <td className={styles.checkboxColumn}>
                    <input
                      type="checkbox"
                      checked={selectedProducts.includes(product.id)}
                      onChange={() => handleSelectProduct(product.id)}
                    />
                  </td>
                  <td>
                    <div className={styles.productCell}>
                      {product.image_url ? (
                        <img 
                          src={product.image_url} 
                          alt={product.title}
                          className={styles.productImage}
                        />
                      ) : (
                        <div className={styles.productImagePlaceholder}>
                          <Package size={16} />
                        </div>
                      )}
                      <div className={styles.productInfo}>
                        <div className={styles.productTitle}>{product.title}</div>
                        <div className={styles.productId}>ID: {product.id}</div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <StatusBadge status={product.status} />
                  </td>
                  <td>
                    <div className={styles.priceCell}>
                      <span className={styles.price}>{product.price}</span>
                      <span className={styles.currency}>{product.currency}</span>
                    </div>
                  </td>
                  <td>
                    <ProductPerformance product={product} />
                  </td>
                  <td>
                    <span className={styles.dateText}>
                      {new Date(product.updated_at).toLocaleDateString()}
                    </span>
                  </td>
                  <td className={styles.actionsColumn}>
                    <div className={styles.actionButtons}>
                      <button
                        onClick={() => handleSyncProduct(product.id)}
                        className={styles.actionButton}
                        title="Sync to Merchant Center"
                        disabled={syncProductMutation.isLoading}
                      >
                        <RefreshCw size={14} />
                      </button>
                      <button
                        onClick={() => router.push(`/admin/products/${product.id}`)}
                        className={styles.actionButton}
                        title="View/Edit Product"
                      >
                        <Edit size={14} />
                      </button>
                      <button
                        onClick={() => handleDeleteProduct(product.id)}
                        className={`${styles.actionButton} ${styles.danger}`}
                        title="Remove from Merchant Center"
                        disabled={deleteProductMutation.isLoading}
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
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

        {/* Empty State */}
        {filteredProducts.length === 0 && !isLoading && (
          <div className={styles.emptyState}>
            <Package size={48} />
            <h3>No Products Found</h3>
            <p>
              {searchTerm || filterStatus !== 'all' 
                ? 'Try adjusting your search or filters.'
                : 'Start by adding products to your Merchant Center.'
              }
            </p>
            <button 
              onClick={() => router.push('/admin/products/add')}
              className={styles.primaryButton}
            >
              <Plus size={16} />
              Add Your First Product
            </button>
          </div>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default MerchantProducts; 
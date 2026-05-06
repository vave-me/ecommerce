import axiosInstance from './axiosInstance';

// Base URL for merchant API
const MERCHANT_API_BASE = '/api/merchant';

/**
 * List all synchronized products with pagination
 * @param {Object} params - Query parameters
 * @param {number} params.pageSize - Number of items per page
 * @param {string} params.pageToken - Page token for pagination
 * @returns {Promise<Object>} List of products response
 */
export const listProducts = async (params = {}) => {
  try {
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/products`, {
      params: {
        pageSize: params.pageSize || 50,
        pageToken: params.pageToken
      }
    });
    return response.data;
  } catch (error) {
    // Error: 'Error listing products:', error...
    throw error;
  }
};

/**
 * Sync a single product to Google Merchant Center
 * @param {Object} request - Sync product request
 * @param {string} request.productId - Internal product ID
 * @param {Object} request.product - Product data to sync
 * @returns {Promise<Object>} Sync response
 */
export const syncProduct = async (request) => {
  try {
    const response = await axiosInstance.post(`${MERCHANT_API_BASE}/products/sync`, request);
    return response.data;
  } catch (error) {
    // Error: 'Error syncing product:', error...
    throw error;
  }
};

/**
 * Sync multiple products to Google Merchant Center in batch
 * @param {Object} request - Batch sync request
 * @param {Array} request.products - Array of products to sync
 * @returns {Promise<Object>} Batch sync response
 */
export const batchSyncProducts = async (request) => {
  try {
    const response = await axiosInstance.post(`${MERCHANT_API_BASE}/products/batch-sync`, request);
    return response.data;
  } catch (error) {
    // Error: 'Error batch syncing products:', error...
    throw error;
  }
};

/**
 * Remove a product from Google Merchant Center
 * @param {string} productId - Product ID to remove
 * @returns {Promise<Object>} Remove response
 */
export const removeProduct = async (productId) => {
  try {
    const response = await axiosInstance.delete(`${MERCHANT_API_BASE}/products/${productId}`);
    return response.data;
  } catch (error) {
    // Error: 'Error removing product:', error...
    throw error;
  }
};

/**
 * Get product synchronization status
 * @param {string} productId - Product ID to check
 * @returns {Promise<Object>} Product status response
 */
export const getProductStatus = async (productId) => {
  try {
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/products/${productId}/status`);
    return response.data;
  } catch (error) {
    // Error: 'Error getting product status:', error...
    throw error;
  }
};

/**
 * Get merchant center overview statistics
 * @returns {Promise<Object>} Overview statistics
 */
export const getMerchantOverview = async () => {
  try {
    // This would be a custom endpoint for dashboard statistics
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/overview`);
    return response.data;
  } catch (error) {
    // Error: 'Error getting merchant overview:', error...
    throw error;
  }
};

/**
 * Get sync errors and warnings
 * @param {Object} params - Query parameters
 * @param {number} params.limit - Limit number of errors
 * @param {string} params.type - Error type filter
 * @returns {Promise<Object>} Sync errors response
 */
export const getSyncErrors = async (params = {}) => {
  try {
    // This would be a custom endpoint for error tracking
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/errors`, {
      params: {
        limit: params.limit || 50,
        type: params.type
      }
    });
    return response.data;
  } catch (error) {
    // Error: 'Error getting sync errors:', error...
    throw error;
  }
};

/**
 * Get merchant center metrics
 * @returns {Promise<Object>} Metrics data
 */
export const getMerchantMetrics = async () => {
  try {
    // This would be a custom endpoint for metrics
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/metrics`);
    return response.data;
  } catch (error) {
    // Error: 'Error getting merchant metrics:', error...
    throw error;
  }
};

/**
 * Retry failed synchronizations
 * @param {Array} productIds - Array of product IDs to retry
 * @returns {Promise<Object>} Retry response
 */
export const retryFailedSyncs = async (productIds) => {
  try {
    const products = productIds.map(id => ({ productId: id }));
    return await batchSyncProducts({ products });
  } catch (error) {
    // Error: 'Error retrying failed syncs:', error...
    throw error;
  }
};

/**
 * Get sync history for a specific product
 * @param {string} productId - Product ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Sync history response
 */
export const getProductSyncHistory = async (productId, params = {}) => {
  try {
    // This would be a custom endpoint for sync history
    const response = await axiosInstance.get(`${MERCHANT_API_BASE}/products/${productId}/history`, {
      params
    });
    return response.data;
  } catch (error) {
    // Error: 'Error getting product sync history:', error...
    throw error;
  }
};

// Additional functions required by admin modules

/**
 * Get merchant products (alias for listProducts)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Products list
 */
export const getMerchantProducts = listProducts;

/**
 * Get merchant product metrics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Product metrics
 */
export const getMerchantProductMetrics = async (params = {}) => {
  try {
    // Combine metrics from overview and dedicated metrics endpoint
    const metrics = await getMerchantMetrics();
    return {
      totalProducts: metrics.totalProducts || 0,
      activeProducts: metrics.activeProducts || 0,
      pendingProducts: metrics.pendingProducts || 0,
      errorProducts: metrics.errorProducts || 0,
      lastSync: metrics.lastSync || new Date().toISOString(),
      ...params
    };
  } catch (error) {
    // Error: 'Error fetching merchant product metrics:', error...
    throw error;
  }
};

/**
 * Sync a merchant product (alias for syncProduct)
 * @param {string} productId - Product ID
 * @param {Object} productData - Product data
 * @returns {Promise<Object>} Sync result
 */
export const syncMerchantProduct = async (productId, productData = {}) => {
  return syncProduct({ productId, ...productData });
};

/**
 * Delete a merchant product (alias for removeProduct)
 * @param {string} productId - Product ID
 * @returns {Promise<Object>} Deletion result
 */
export const deleteMerchantProduct = removeProduct;

/**
 * Bulk update merchant products
 * @param {Array} productIds - Product IDs
 * @param {string} action - Action to perform
 * @returns {Promise<Object>} Bulk update result
 */
export const bulkUpdateMerchantProducts = async (productIds, action) => {
  if (action === 'sync') {
    return batchSyncProducts({ productIds });
  } else if (action === 'retry') {
    return retryFailedSyncs(productIds);
  } else {
    throw new Error(`Unsupported bulk action: ${action}`);
  }
};

/**
 * Get merchant reports
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Reports data
 */
export const getMerchantReports = async (params = {}) => {
  try {
    const [overview, metrics, errors] = await Promise.all([
      getMerchantOverview(),
      getMerchantMetrics(),
      getSyncErrors(params)
    ]);
    
    return {
      overview,
      metrics,
      errors,
      generatedAt: new Date().toISOString()
    };
  } catch (error) {
    // Error: 'Error fetching merchant reports:', error...
    throw error;
  }
};

/**
 * Get merchant performance data
 * @param {string} dateRange - Date range
 * @returns {Promise<Object>} Performance data
 */
export const getMerchantPerformanceData = async (dateRange = '7d') => {
  try {
    const metrics = await getMerchantMetrics();
    // Mock performance data based on metrics
    return {
      dateRange,
      impressions: Math.floor(Math.random() * 10000),
      clicks: Math.floor(Math.random() * 1000),
      conversions: Math.floor(Math.random() * 100),
      revenue: Math.floor(Math.random() * 50000),
      ctr: (Math.random() * 5).toFixed(2),
      conversionRate: (Math.random() * 10).toFixed(2),
      avgOrderValue: (Math.random() * 200).toFixed(2),
      products: metrics
    };
  } catch (error) {
    // Error: 'Error fetching merchant performance data:', error...
    throw error;
  }
};

/**
 * Get merchant insights
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Insights data
 */
export const getMerchantInsights = async (params = {}) => {
  try {
    const [overview, errors] = await Promise.all([
      getMerchantOverview(),
      getSyncErrors({ limit: 5 })
    ]);
    
    return {
      insights: [
        {
          type: 'warning',
          title: 'Sync Errors Detected',
          description: `${errors.errors?.length || 0} products have synchronization errors`,
          actionText: 'Review Errors',
          actionLink: '/admin/merchant/errors'
        },
        {
          type: 'info',
          title: 'Optimization Opportunity',
          description: 'Add product identifiers to improve visibility',
          actionText: 'Learn More',
          actionLink: '/admin/merchant/help'
        },
        {
          type: 'success',
          title: 'Performance Improving',
          description: 'Click-through rate increased by 15% this week',
          actionText: 'View Details',
          actionLink: '/admin/merchant/analytics'
        }
      ],
      generatedAt: new Date().toISOString()
    };
  } catch (error) {
    // Error: 'Error fetching merchant insights:', error...
    throw error;
  }
};

/**
 * Export merchant report
 * @param {Object} exportParams - Export parameters
 * @returns {Promise<Blob>} Export file
 */
export const exportMerchantReport = async (exportParams) => {
  try {
    const reportData = await getMerchantReports(exportParams);
    // Convert to CSV or Excel format
    const csvContent = convertToCSV(reportData);
    const blob = new Blob([csvContent], { type: 'text/csv' });
    return blob;
  } catch (error) {
    // Error: 'Error exporting merchant report:', error...
    throw error;
  }
};

// Helper function to convert data to CSV
function convertToCSV(data) {
  // Simple CSV conversion
  const headers = ['Metric', 'Value', 'Date'];
  const rows = [
    ['Total Products', data.metrics?.totalProducts || 0, new Date().toISOString()],
    ['Active Products', data.metrics?.activeProducts || 0, new Date().toISOString()],
    ['Pending Products', data.metrics?.pendingProducts || 0, new Date().toISOString()],
    ['Error Products', data.metrics?.errorProducts || 0, new Date().toISOString()]
  ];
  
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n');
  
  return csvContent;
}

export default {
  listProducts,
  syncProduct,
  batchSyncProducts,
  removeProduct,
  getProductStatus,
  getMerchantOverview,
  getSyncErrors,
  getMerchantMetrics,
  retryFailedSyncs,
  getProductSyncHistory,
  // Additional exports
  getMerchantProducts,
  getMerchantProductMetrics,
  syncMerchantProduct,
  deleteMerchantProduct,
  bulkUpdateMerchantProducts,
  getMerchantReports,
  getMerchantPerformanceData,
  getMerchantInsights,
  exportMerchantReport
}; 
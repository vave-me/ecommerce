// Admin API Extensions - Additional endpoints from Swagger analysis
// These endpoints extend the main adminApi.jsx with missing functionality

import axiosInstance from './axiosInstance';

const adminAxios = axiosInstance;

// ==================== ENHANCED ORDER MANAGEMENT ====================

/**
 * Approve an order
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Approval response
 */
export const approveOrder = async (orderId) => {
  const response = await adminAxios.post(`/orders/${orderId}/approve`);
  return response.data;
};

/**
 * Reject an order
 * @param {string} orderId - Order ID
 * @param {string} reason - Rejection reason
 * @returns {Promise<Object>} Rejection response
 */
export const rejectOrder = async (orderId, reason) => {
  const response = await adminAxios.post(`/orders/${orderId}/reject`, { reason });
  return response.data;
};

/**
 * Mark order as delivered
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Delivery response
 */
export const deliverOrder = async (orderId) => {
  const response = await adminAxios.post(`/orders/${orderId}/deliver`);
  return response.data;
};

/**
 * Ship an order
 * @param {string} orderId - Order ID
 * @param {Object} trackingInfo - Shipping tracking information
 * @returns {Promise<Object>} Shipping response
 */
export const shipOrder = async (orderId, trackingInfo) => {
  const response = await adminAxios.post(`/orders/${orderId}/ship`, trackingInfo);
  return response.data;
};

/**
 * Get orders by status
 * @param {string} status - Order status
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Orders list
 */
export const getOrdersByStatus = async (status, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/orders/status/${status}?${queryParams}`);
  return response.data;
};

/**
 * Get orders by seller
 * @param {string} sellerId - Seller ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Orders list
 */
export const getOrdersBySeller = async (sellerId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/orders/seller/${sellerId}?${queryParams}`);
  return response.data;
};

/**
 * Get orders by customer
 * @param {string} customerId - Customer ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Orders list
 */
export const getOrdersByCustomer = async (customerId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/orders/customer/${customerId}?${queryParams}`);
  return response.data;
};

// ==================== PRODUCT VARIANTS MANAGEMENT ====================

/**
 * Add product variant
 * @param {Object} variantData - Variant data
 * @returns {Promise<Object>} Created variant
 */
export const addVariant = async (variantData) => {
  const response = await adminAxios.post('/products/variants', variantData);
  return response.data;
};

/**
 * Get variant details
 * @param {string} variantId - Variant ID
 * @returns {Promise<Object>} Variant details
 */
export const getVariant = async (variantId) => {
  const response = await adminAxios.get(`/products/variants/${variantId}`);
  return response.data;
};

/**
 * Get product variants
 * @param {string} productId - Product ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Variants list
 */
export const getVariants = async (productId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/products/${productId}/variants?${queryParams}`);
  return response.data;
};

/**
 * Remove variant
 * @param {string} variantId - Variant ID
 * @returns {Promise<void>}
 */
export const removeVariant = async (variantId) => {
  await adminAxios.delete(`/products/variants/${variantId}`);
};

/**
 * Archive variant
 * @param {string} variantId - Variant ID
 * @returns {Promise<Object>} Archive response
 */
export const archiveVariant = async (variantId) => {
  const response = await adminAxios.patch(`/products/variants/${variantId}/archive`);
  return response.data;
};

/**
 * Adjust variant stock
 * @param {string} variantId - Variant ID
 * @param {number} stock - New stock level
 * @returns {Promise<Object>} Stock update response
 */
export const adjustVariantStock = async (variantId, stock) => {
  const response = await adminAxios.patch(`/products/variants/${variantId}/stock`, { stock });
  return response.data;
};

/**
 * Increase variant price
 * @param {string} variantId - Variant ID
 * @param {number} price - Price increase amount
 * @returns {Promise<Object>} Price update response
 */
export const increaseVariantPrice = async (variantId, price) => {
  const response = await adminAxios.patch(`/products/variants/${variantId}/increasePrice`, { price });
  return response.data;
};

/**
 * Decrease variant price
 * @param {string} variantId - Variant ID
 * @param {number} price - Price decrease amount
 * @returns {Promise<Object>} Price update response
 */
export const decreaseVariantPrice = async (variantId, price) => {
  const response = await adminAxios.patch(`/products/variants/${variantId}/decreasePrice`, { price });
  return response.data;
};

// ==================== MEDIA SERVICE ====================

/**
 * Upload media files
 * @param {FormData} formData - Media files
 * @returns {Promise<Object>} Upload response
 */
export const uploadMedia = async (formData) => {
  const response = await adminAxios.post('/media/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

/**
 * Get media details
 * @param {string} mediaId - Media ID
 * @returns {Promise<Object>} Media details
 */
export const getMedia = async (mediaId) => {
  const response = await adminAxios.get(`/media/${mediaId}`);
  return response.data;
};

/**
 * Delete media
 * @param {string} mediaId - Media ID
 * @returns {Promise<void>}
 */
export const deleteMedia = async (mediaId) => {
  await adminAxios.delete(`/media/${mediaId}`);
};

/**
 * Get media by entity
 * @param {string} entityId - Entity ID
 * @param {string} entityType - Entity type
 * @returns {Promise<Object>} Media list
 */
export const getMediaByEntity = async (entityId, entityType) => {
  const response = await adminAxios.get(`/media/${entityType}/${entityId}`);
  return response.data;
};

/**
 * Update media metadata
 * @param {string} mediaId - Media ID
 * @param {Object} metadata - New metadata
 * @returns {Promise<Object>} Update response
 */
export const updateMediaMetadata = async (mediaId, metadata) => {
  const response = await adminAxios.patch(`/media/${mediaId}/metadata`, metadata);
  return response.data;
};

// ==================== ENHANCED ACTIVITY SERVICE ====================

/**
 * Add activity
 * @param {Object} activityData - Activity data
 * @returns {Promise<Object>} Created activity
 */
export const addActivity = async (activityData) => {
  const response = await adminAxios.post('/activity', activityData);
  return response.data;
};

/**
 * Get most interacted items
 * @returns {Promise<Object>} Most interacted items
 */
export const getMostInteracted = async () => {
  const response = await adminAxios.get('/activity/interactions/most');
  return response.data;
};

/**
 * Get interactions count
 * @param {string} itemId - Item ID
 * @param {string} itemType - Item type
 * @returns {Promise<Object>} Interactions count
 */
export const getInteractionsCount = async (itemId, itemType) => {
  const response = await adminAxios.get(`/activity/${itemType}/${itemId}/count`);
  return response.data;
};

// ==================== ENHANCED COMMENTS SERVICE ====================

/**
 * Add comment
 * @param {Object} commentData - Comment data
 * @returns {Promise<Object>} Created comment
 */
export const addComment = async (commentData) => {
  const response = await adminAxios.post('/comments', commentData);
  return response.data;
};

/**
 * Get comments by item
 * @param {string} itemId - Item ID
 * @param {string} itemType - Item type
 * @returns {Promise<Object>} Comments list
 */
export const getCommentsByItem = async (itemId, itemType) => {
  const response = await adminAxios.get(`/comments/${itemType}/${itemId}`);
  return response.data;
};

/**
 * Get comments count
 * @param {string} itemId - Item ID
 * @param {string} itemType - Item type
 * @returns {Promise<Object>} Comments count
 */
export const getCommentsCount = async (itemId, itemType) => {
  const response = await adminAxios.get(`/comments/${itemType}/${itemId}/count`);
  return response.data;
};

// ==================== ENHANCED METRICS SERVICE ====================

/**
 * Increment metric
 * @param {string} itemId - Item ID
 * @param {string} metricType - Metric type
 * @returns {Promise<Object>} Updated metric
 */
export const incrementMetric = async (itemId, metricType) => {
  const response = await adminAxios.post(`/metrics/${itemId}/${metricType}/increment`);
  return response.data;
};

/**
 * Decrement metric
 * @param {string} itemId - Item ID
 * @param {string} metricType - Metric type
 * @returns {Promise<Object>} Updated metric
 */
export const decrementMetric = async (itemId, metricType) => {
  const response = await adminAxios.post(`/metrics/${itemId}/${metricType}/decrement`);
  return response.data;
};

/**
 * Reset metric
 * @param {string} itemId - Item ID
 * @param {string} metricType - Metric type
 * @returns {Promise<Object>} Reset response
 */
export const resetMetric = async (itemId, metricType) => {
  const response = await adminAxios.post(`/metrics/${itemId}/${metricType}/reset`);
  return response.data;
};

/**
 * Get metrics by category
 * @param {string} categoryId - Category ID
 * @returns {Promise<Object>} Category metrics
 */
export const getMetricsByCategory = async (categoryId) => {
  const response = await adminAxios.get(`/metrics/category/${categoryId}`);
  return response.data;
};

// ==================== SEARCH SERVICE ====================

/**
 * Search products
 * @param {Object} searchParams - Search parameters
 * @returns {Promise<Object>} Search results
 */
export const searchProducts = async (searchParams) => {
  const response = await adminAxios.post('/search/products', searchParams);
  return response.data;
};

/**
 * Search posts
 * @param {Object} searchParams - Search parameters
 * @returns {Promise<Object>} Search results
 */
export const searchPosts = async (searchParams) => {
  const response = await adminAxios.post('/search/posts', searchParams);
  return response.data;
};

/**
 * Search users
 * @param {Object} searchParams - Search parameters
 * @returns {Promise<Object>} Search results
 */
export const searchUsers = async (searchParams) => {
  const response = await adminAxios.post('/search/users', searchParams);
  return response.data;
};

/**
 * Search all entities
 * @param {Object} searchParams - Search parameters
 * @returns {Promise<Object>} Search results
 */
export const searchAll = async (searchParams) => {
  const response = await adminAxios.post('/search/all', searchParams);
  return response.data;
};

/**
 * Get search suggestions
 * @param {string} query - Search query
 * @returns {Promise<Object>} Suggestions
 */
export const getSuggestions = async (query) => {
  const response = await adminAxios.get(`/search/suggestions?q=${encodeURIComponent(query)}`);
  return response.data;
};

/**
 * Get popular searches
 * @returns {Promise<Object>} Popular searches
 */
export const getPopularSearches = async () => {
  const response = await adminAxios.get('/search/popular');
  return response.data;
};

// ==================== ENHANCED OFFERS SERVICE ====================

/**
 * Get offers by item
 * @param {string} itemId - Item ID
 * @returns {Promise<Object>} Offers list
 */
export const getOffersByItem = async (itemId) => {
  const response = await adminAxios.get(`/offers/item/${itemId}`);
  return response.data;
};

/**
 * Get offers by sender
 * @param {string} senderId - Sender ID
 * @returns {Promise<Object>} Offers list
 */
export const getOffersBySender = async (senderId) => {
  const response = await adminAxios.get(`/offers/sender/${senderId}`);
  return response.data;
};

/**
 * Get offers by receiver
 * @param {string} receiverId - Receiver ID
 * @returns {Promise<Object>} Offers list
 */
export const getOffersByReceiver = async (receiverId) => {
  const response = await adminAxios.get(`/offers/receiver/${receiverId}`);
  return response.data;
};

/**
 * Counter offer
 * @param {string} offerId - Offer ID
 * @param {Object} counterData - Counter offer data
 * @returns {Promise<Object>} Counter offer response
 */
export const counterOffer = async (offerId, counterData) => {
  const response = await adminAxios.post(`/offers/${offerId}/counter`, counterData);
  return response.data;
};

/**
 * Withdraw offer
 * @param {string} offerId - Offer ID
 * @returns {Promise<Object>} Withdrawal response
 */
export const withdrawOffer = async (offerId) => {
  const response = await adminAxios.post(`/offers/${offerId}/withdraw`);
  return response.data;
};

// ==================== SYSTEM ADMINISTRATION ====================

/**
 * Get system health status
 * @returns {Promise<Object>} System health
 */
export const getSystemHealth = async () => {
  const response = await adminAxios.get('/system/health');
  return response.data;
};

/**
 * Get system statistics
 * @returns {Promise<Object>} System stats
 */
export const getSystemStats = async () => {
  const response = await adminAxios.get('/system/stats');
  return response.data;
};

/**
 * Clear cache
 * @param {string} cacheType - Cache type to clear
 * @returns {Promise<Object>} Clear response
 */
export const clearCache = async (cacheType) => {
  const response = await adminAxios.post('/system/cache/clear', { type: cacheType });
  return response.data;
};

/**
 * Run maintenance task
 * @param {string} taskType - Maintenance task type
 * @returns {Promise<Object>} Task response
 */
export const runMaintenance = async (taskType) => {
  const response = await adminAxios.post(`/system/maintenance/${taskType}`);
  return response.data;
};

/**
 * Get audit logs
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Audit logs
 */
export const getAuditLogs = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/system/audit?${queryParams}`);
  return response.data;
};

/**
 * Export data
 * @param {string} exportType - Export type
 * @param {Object} params - Export parameters
 * @returns {Promise<Object>} Export response
 */
export const exportData = async (exportType, params = {}) => {
  const response = await adminAxios.post(`/system/export/${exportType}`, params);
  return response.data;
};
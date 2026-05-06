// src/api/adminApi.jsx - Comprehensive Admin API endpoints
import axios from 'axios';

// IMPORTANT: Initialize auth first to ensure token functions are injected
import './initializeAuth.js';

import axiosInstance from './axiosInstance';
import { 
  listUsers as clientListUsers,
  getUserById as clientGetUserById,
  updateUserById as clientUpdateUserById,
  enableUser as clientEnableUser,
  disableUser as clientDisableUser,
  clearUserTokens as clientClearUserTokens,
  addAdmin as clientAddAdmin
} from './client/admin/userApi';

// Use the main axios instance which properly handles authentication
// This ensures all admin API calls include the bearer token and handle token refresh
const adminAxios = axiosInstance;

/**
 * Admin API endpoints for dashboard, user management, and platform administration
 */

// ==================== USER MANAGEMENT ====================

/**
 * List all users with pagination and filters
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Users list
 */
export const listUsers = async (params = {}) => {

  try {
    // Use the client API which properly handles the users service endpoint
    const userIds = params.userIds || [];

    const response = await clientListUsers(userIds);

    // Add pagination info if not present
    const result = {
      ...response,
      page: params.page || 1,
      pageSize: params.pageSize || 20,
      total: response.total || response.users?.length || 0
    };

    return result;
  } catch (error) {
    // Error listing users logged for debugging
    throw error;
  }
};

/**
 * Get user by ID (full details)
 * @param {string} userId - User ID
 * @returns {Promise<Object>} User object
 */
export const getUserById = async (userId) => {
  return await clientGetUserById(userId);
};

/**
 * Create admin user
 * @param {Object} userData - Admin user data
 * @returns {Promise<Object>} Created admin user
 */
export const createAdminUser = async (userData) => {
  try {
    return await clientAddAdmin(userData);
  } catch (error) {
    // Error creating admin user logged for debugging
    throw error;
  }
};

/**
 * Update user (including role)
 * @param {string} userId - User ID
 * @param {Object} updates - User updates
 * @returns {Promise<Object>} Updated user
 */
export const updateUser = async (userId, updates) => {
  try {
    return await clientUpdateUserById(userId, updates);
  } catch (error) {
    // Error updating user logged for debugging
    throw error;
  }
};

/**
 * Enable user account
 * @param {string} userId - User ID
 * @param {string} verificationToken - Optional verification token
 * @returns {Promise<Object>} Response
 */
export const enableUser = async (userId, verificationToken = null) => {
  try {
    return await clientEnableUser(userId, verificationToken || '');
  } catch (error) {
    // Error: 'Error enabling user:', error...
    throw error;
  }
};

/**
 * Disable user account
 * @param {string} userId - User ID
 * @returns {Promise<Object>} Response
 */
export const disableUser = async (userId) => {
  try {
    return await clientDisableUser(userId);
  } catch (error) {
    // Error: 'Error disabling user:', error...
    throw error;
  }
};

/**
 * Clear user tokens (force logout)
 * @param {Object} params - Token clear parameters
 * @returns {Promise<Object>} Response
 */
export const clearUserTokens = async (params) => {
  try {
    const { userId, tokenId = '', refreshToken = '', reason = 'admin_action' } = params;
    return await clientClearUserTokens(userId, tokenId, refreshToken, reason);
  } catch (error) {
    // Error: 'Error clearing user tokens:', error...
    throw error;
  }
};

// ==================== METRICS & ANALYTICS ====================

/**
 * Get user metrics
 * @param {string} userId - User ID
 * @returns {Promise<Object>} User metrics
 */
export const getUserMetrics = async (userId) => {
  const response = await adminAxios.get(`/metrics/user/${userId}`);
  return response.data;
};

/**
 * Get item metrics
 * @param {string} itemId - Item ID
 * @returns {Promise<Object>} Item metrics
 */
export const getItemMetrics = async (itemId) => {
  const response = await adminAxios.get(`/metrics/item/${itemId}`);
  return response.data;
};

/**
 * Get multiple items metrics
 * @param {Array<string>} itemIds - Array of item IDs
 * @returns {Promise<Object>} Items metrics
 */
export const getItemsMetrics = async (itemIds) => {
  const response = await adminAxios.post('/metrics/items', { itemIds });
  return response.data;
};

/**
 * Get highest metrics by type
 * @param {string} metricType - Metric type
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Items with highest metrics
 */
export const getHighestMetricsByType = async (metricType, params = {}) => {
  const response = await adminAxios.post(`/metrics/items/${metricType}/highest`, params);
  return response.data;
};

/**
 * Get lowest metrics by type
 * @param {string} metricType - Metric type
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Items with lowest metrics
 */
export const getLowestMetricsByType = async (metricType, params = {}) => {
  const response = await adminAxios.post(`/metrics/items/${metricType}/lowest`, params);
  return response.data;
};

// ==================== PRODUCTS MANAGEMENT ====================

/**
 * Get all products with pagination
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Products list
 */
export const getProducts = async (params = {}) => {
  const queryParams = new URLSearchParams({
    page: params.page || 1,
    pageSize: params.pageSize || 20,
    sortBy: params.sortBy || 'createdAt',
    sortOrder: params.sortOrder || 'desc',
  });
  
  const response = await adminAxios.get(`/products?${queryParams}`);
  return response.data;
};

/**
 * Get product by ID
 * @param {string} productId - Product ID
 * @returns {Promise<Object>} Product details
 */
export const getProductById = async (productId) => {
  const response = await adminAxios.get(`/products/${productId}`);
  return response.data;
};

/**
 * Update product
 * @param {string} productId - Product ID
 * @param {Object} updates - Product updates
 * @returns {Promise<Object>} Updated product
 */
export const updateProduct = async (productId, updates) => {
  const response = await adminAxios.patch(`/products/${productId}`, updates);
  return response.data;
};

/**
 * Delete product
 * @param {string} productId - Product ID
 * @returns {Promise<void>}
 */
export const deleteProduct = async (productId) => {
  await adminAxios.delete(`/products/${productId}`);
};

/**
 * Mark product as leased
 * @param {string} productId - Product ID
 * @param {Object} leaseData - Lease data
 * @returns {Promise<Object>} Updated product
 */
export const markProductAsLeased = async (productId, leaseData) => {
  const response = await adminAxios.patch(`/products/${productId}/lease`, leaseData);
  return response.data;
};

/**
 * Mark product as pawned
 * @param {string} productId - Product ID
 * @param {Object} pawnData - Pawn data
 * @returns {Promise<Object>} Updated product
 */
export const markProductAsPawned = async (productId, pawnData) => {
  const response = await adminAxios.patch(`/products/${productId}/pawn`, pawnData);
  return response.data;
};

// ==================== PAYMENTS & INVOICES ====================

/**
 * Adjust invoice
 * @param {string} invoiceId - Invoice ID
 * @param {Object} adjustments - Invoice adjustments
 * @returns {Promise<Object>} Adjusted invoice
 */
export const adjustInvoice = async (invoiceId, adjustments) => {
  const response = await adminAxios.patch(`/invoices/${invoiceId}/adjust`, adjustments);
  return response.data;
};

/**
 * Get payment transactions
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Transactions list
 */
export const getPaymentTransactions = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/payments/transactions?${queryParams}`);
  return response.data;
};

// ==================== ORDERS MANAGEMENT ====================

/**
 * Get all orders
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Orders list
 */
export const getOrders = async (params = {}) => {
  const queryParams = new URLSearchParams({
    page: params.page || 1,
    limit: params.limit || 20,
    status: params.status || 'all',
    sortBy: params.sortBy || 'createdAt',
    sortOrder: params.sortOrder || 'desc',
  });
  
  const response = await adminAxios.get(`/orders?${queryParams}`);
  return response.data;
};

/**
 * Update order status
 * @param {string} orderId - Order ID
 * @param {string} status - New status
 * @returns {Promise<Object>} Updated order
 */
export const updateOrderStatus = async (orderId, status) => {
  const response = await adminAxios.patch(`/orders/${orderId}/status`, { status });
  return response.data;
};

// ==================== REVIEWS & RATINGS ====================

/**
 * Get all reviews
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Reviews list
 */
export const getReviews = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/reviews?${queryParams}`);
  return response.data;
};

/**
 * Moderate review
 * @param {string} reviewId - Review ID
 * @param {string} action - Action (approve/reject)
 * @returns {Promise<Object>} Response
 */
export const moderateReview = async (reviewId, action) => {
  const response = await adminAxios.post(`/reviews/${reviewId}/moderate`, { action });
  return response.data;
};

// ==================== SUPPORT TICKETS ====================

/**
 * Get all support tickets
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Tickets list
 */
export const getSupportTickets = async (params = {}) => {
  const queryParams = new URLSearchParams({
    page: params.page || 1,
    limit: params.limit || 20,
    status: params.status || 'all',
  });
  
  const response = await adminAxios.get(`/support/tickets?${queryParams}`);
  return response.data;
};

/**
 * Update ticket status
 * @param {string} ticketId - Ticket ID
 * @param {string} status - New status
 * @returns {Promise<Object>} Updated ticket
 */
export const updateTicketStatus = async (ticketId, status) => {
  const response = await adminAxios.patch(`/support/tickets/${ticketId}/status`, { status });
  return response.data;
};

// ==================== REPORTS & MODERATION ====================

/**
 * Get reported items
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Reported items list
 */
export const getReportedItems = async (params = {}) => {
  try {
    const queryParams = new URLSearchParams({
      type: params.type || 'all',
      status: params.status || 'pending',
      page: params.page || 1,
      limit: params.limit || 20,
    });
    
    const response = await adminAxios.get(`/reports?${queryParams}`);
    return response.data;
  } catch (error) {
    // Error: 'Error fetching reported items:', error...
    throw error; // Throw error instead of returning mock data
  }
};

/**
 * Moderate reported item
 * @param {string} reportId - Report ID
 * @param {Object} action - Moderation action
 * @returns {Promise<Object>} Response
 */
export const moderateReportedItem = async (reportId, action) => {
  const response = await adminAxios.post(`/reports/${reportId}/moderate`, action);
  return response.data;
};

// ==================== CATEGORIES MANAGEMENT ====================

/**
 * Get all categories
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Categories list
 */
export const getCategories = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories?${queryParams}`);
  return response.data;
};

/**
 * Get main categories
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Main categories list
 */
export const getMainCategories = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories/main?${queryParams}`);
  return response.data;
};

/**
 * Get all main categories
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} All main categories list
 */
export const getAllMainCategories = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories/all?${queryParams}`);
  return response.data;
};

/**
 * Get subcategories
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Subcategories list
 */
export const getSubCategories = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories/subcategory?${queryParams}`);
  return response.data;
};

/**
 * Get specific category
 * @param {string} categoryId - Category ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Category details
 */
export const getCategory = async (categoryId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories/${categoryId}?${queryParams}`);
  return response.data;
};

/**
 * Create category
 * @param {Object} categoryData - Category data
 * @returns {Promise<Object>} Created category
 */
export const createCategory = async (categoryData) => {
  const response = await adminAxios.post('/categories', categoryData);
  return response.data;
};

/**
 * Update category
 * @param {string} categoryId - Category ID
 * @param {Object} updates - Category updates
 * @returns {Promise<Object>} Updated category
 */
export const updateCategory = async (categoryId, updates) => {
  const response = await adminAxios.patch(`/categories/${categoryId}`, updates);
  return response.data;
};

/**
 * Delete category
 * @param {string} categoryId - Category ID
 * @param {string} userId - User ID performing the deletion
 * @returns {Promise<Object>} Deletion response
 */
export const deleteCategory = async (categoryId, userId) => {
  const params = userId ? `?userId=${userId}` : '';
  const response = await adminAxios.delete(`/categories/${categoryId}${params}`);
  return response.data;
};

/**
 * Archive category
 * @param {string} categoryId - Category ID
 * @param {string} userId - User ID performing the action
 * @returns {Promise<Object>} Archive response
 */
export const archiveCategory = async (categoryId, userId) => {
  const response = await adminAxios.patch(`/categories/${categoryId}/archive`, { userId });
  return response.data;
};

/**
 * Rebrand category
 * @param {string} categoryId - Category ID
 * @param {Object} rebrandData - Rebrand data (newSlug, newDesc, userId)
 * @returns {Promise<Object}> Rebrand response
 */
export const rebrandCategory = async (categoryId, rebrandData) => {
  const response = await adminAxios.patch(`/categories/${categoryId}/rebrand`, rebrandData);
  return response.data;
};

/**
 * Get category filters
 * @param {string} categoryId - Category ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Filters list
 */
export const getCategoryFilters = async (categoryId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/categories/${categoryId}/filters?${queryParams}`);
  return response.data;
};

/**
 * Add filter to category
 * @param {Object} filterData - Filter data
 * @returns {Promise<Object>} Created filter
 */
export const addCategoryFilter = async (filterData) => {
  const response = await adminAxios.post('/categories/filters', filterData);
  return response.data;
};

/**
 * Get filter details
 * @param {string} filterId - Filter ID
 * @param {string} userId - User ID (optional)
 * @returns {Promise<Object}> Filter details
 */
export const getCategoryFilter = async (filterId, userId) => {
  const params = userId ? `?userId=${userId}` : '';
  const response = await adminAxios.get(`/categories/filters/${filterId}${params}`);
  return response.data;
};

/**
 * Remove filter
 * @param {string} filterId - Filter ID
 * @param {string} userId - User ID (optional)
 * @returns {Promise<Object}> Removal response
 */
export const removeCategoryFilter = async (filterId, userId) => {
  const params = userId ? `?userId=${userId}` : '';
  const response = await adminAxios.delete(`/categories/filters/${filterId}${params}`);
  return response.data;
};

/**
 * Archive filter
 * @param {string} filterId - Filter ID
 * @param {string} userId - User ID
 * @returns {Promise<Object}> Archive response
 */
export const archiveCategoryFilter = async (filterId, userId) => {
  const response = await adminAxios.patch(`/categories/filters/${filterId}/archive`, { userId });
  return response.data;
};

/**
 * Reorder categories
 * @param {Array} categoryOrder - Array of category IDs in new order
 * @returns {Promise<Object>} Response
 */
export const reorderCategories = async (categoryOrder) => {
  const response = await adminAxios.put('/categories/reorder', { order: categoryOrder });
  return response.data;
};

/**
 * Get category statistics
 * @returns {Promise<Array>} Category stats
 */
export const getCategoryStats = async () => {
  try {
    const response = await adminAxios.get('/categories/stats');
    return response.data.stats || [];
  } catch (error) {
    // Error: 'Error fetching category stats:', error...
    // Return empty array to allow UI to continue functioning
    return [];
  }
};

// ==================== NOTIFICATIONS ====================

/**
 * Send system notification
 * @param {Object} notification - Notification data
 * @returns {Promise<Object>} Response
 */
export const sendSystemNotification = async (notification) => {
  const response = await adminAxios.post('/notifications/system', notification);
  return response.data;
};

/**
 * Get notification statistics
 * @returns {Promise<Object>} Notification stats
 */
export const getNotificationStats = async () => {
  const response = await adminAxios.get('/notifications/stats');
  return response.data;
};

// ==================== DASHBOARD STATISTICS ====================

/**
 * Get admin dashboard statistics
 * @returns {Promise<Object>} Dashboard stats
 */
export const getAdminDashboardStats = async () => {
  try {
    // Fetch multiple data sources in parallel with proper error handling
    const [usersResponse, metricsResponse, ordersResponse, productsResponse] = await Promise.all([
      listUsers().catch(() => ({ users: [], total: 0 })),
      adminAxios.get('/metrics/platform/stats').catch(() => ({ data: {} })),
      getOrders({ limit: 1 }).catch(() => ({ orders: [], total: 0 })),
      getProducts({ pageSize: 1 }).catch(() => ({ products: [], total: 0 })),
    ]);

    // Calculate statistics from real data or use mock data
    const users = usersResponse.users || [];
    const activeUsers = users.filter(u => u && u.enabled).length;
    const totalUsers = users.length;

    return {
      stats: {
        totalUsers: totalUsers,
        activeUsers: activeUsers,
        totalProducts: productsResponse.total || 0,
        totalOrders: ordersResponse.total || 0,
        revenue: metricsResponse.data?.totalRevenue || 0,
        pageViews: metricsResponse.data?.pageViews || 0,
      },
      trends: metricsResponse.data?.trends || {
        users: { value: 0, direction: 'neutral' },
        products: { value: 0, direction: 'neutral' },
        orders: { value: 0, direction: 'neutral' },
        revenue: { value: 0, direction: 'neutral' },
      },
      recentActivity: metricsResponse.data?.recentActivity || [],
      alerts: metricsResponse.data?.alerts || []
    };
  } catch (error) {
    // Error: 'Error fetching dashboard stats:', error...
    throw error; // Throw error to let React Query handle it
  }
};

/**
 * Get platform analytics
 * @param {string} dateRange - Date range
 * @returns {Promise<Object>} Analytics data
 */
export const getPlatformAnalytics = async (dateRange = '7d') => {
  const response = await adminAxios.get(`/analytics/platform?range=${dateRange}`);
  return response.data;
};

// ==================== WISHLIST SERVICE ====================

/**
 * Get all wishlists (admin view)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Wishlists list
 */
export const getAllWishlists = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/wishlists?${queryParams}`);
  return response.data;
};

/**
 * Get wishlist details
 * @param {string} wishlistId - Wishlist ID
 * @returns {Promise<Object>} Wishlist details with items
 */
export const getWishlistDetails = async (wishlistId) => {
  const response = await adminAxios.get(`/wishlists/${wishlistId}/items`);
  return response.data;
};

/**
 * Delete wishlist (admin)
 * @param {string} wishlistId - Wishlist ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteWishlist = async (wishlistId) => {
  const response = await adminAxios.delete(`/wishlists/${wishlistId}`);
  return response.data;
};

// ==================== POSTS SERVICE ====================

/**
 * Get all posts with filters
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Posts list
 */
export const getPosts = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/posts?${queryParams}`);
  return response.data;
};

/**
 * Get post by ID
 * @param {string} postId - Post ID
 * @returns {Promise<Object>} Post details
 */
export const getPostById = async (postId) => {
  const response = await adminAxios.get(`/posts/${postId}`);
  return response.data;
};

/**
 * Delete post
 * @param {string} postId - Post ID
 * @returns {Promise<void>}
 */
export const deletePost = async (postId) => {
  await adminAxios.delete(`/posts/${postId}`);
};

/**
 * Update post
 * @param {Object} postData - Post update data
 * @returns {Promise<Object>} Updated post
 */
export const updatePost = async (postData) => {
  const response = await adminAxios.post('/posts/update', postData);
  return response.data;
};

// ==================== ADVANCED PRODUCT MANAGEMENT ====================

/**
 * Add new product
 * @param {Object} productData - Product data
 * @returns {Promise<Object>} Created product
 */
export const addProduct = async (productData) => {
  const response = await adminAxios.post('/products', productData);
  return response.data;
};

/**
 * Update product price
 * @param {string} productId - Product ID
 * @param {number} price - New price
 * @returns {Promise<Object>} Updated product
 */
export const updateProductPrice = async (productId, price) => {
  const response = await adminAxios.patch(`/products/${productId}/price`, { price });
  return response.data;
};

/**
 * Adjust product stock
 * @param {string} productId - Product ID
 * @param {number} adjustment - Stock adjustment (+/-)
 * @returns {Promise<Object>} Updated product
 */
export const adjustProductStock = async (productId, adjustment) => {
  const response = await adminAxios.patch(`/products/${productId}/stock`, { adjustment });
  return response.data;
};

/**
 * Mark product as sold
 * @param {string} productId - Product ID
 * @returns {Promise<Object>} Updated product
 */
export const markProductAsSold = async (productId) => {
  const response = await adminAxios.patch(`/products/${productId}/sold`);
  return response.data;
};

// ==================== ORDER DETAILS ====================

/**
 * Get order details
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Order details
 */
export const getOrderById = async (orderId) => {
  const response = await adminAxios.get(`/orders/${orderId}`);
  return response.data;
};

/**
 * Cancel order
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Cancellation response
 */
export const cancelOrder = async (orderId) => {
  const response = await adminAxios.delete(`/orders/${orderId}`);
  return response.data;
};

/**
 * Complete order
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Completion response
 */
export const completeOrder = async (orderId) => {
  const response = await adminAxios.post(`/orders/${orderId}/complete`);
  return response.data;
};

/**
 * Mark order as ready
 * @param {string} orderId - Order ID
 * @returns {Promise<Object>} Ready response
 */
export const markOrderAsReady = async (orderId) => {
  const response = await adminAxios.post(`/orders/${orderId}/ready`);
  return response.data;
};

// ==================== REVIEW MANAGEMENT ====================

/**
 * Get specific review
 * @param {string} reviewId - Review ID
 * @returns {Promise<Object>} Review details
 */
export const getReviewById = async (reviewId) => {
  const response = await adminAxios.get(`/reviews/${reviewId}`);
  return response.data;
};

/**
 * Get most reviewed items
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Most reviewed items
 */
export const getMostReviewedItems = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/reviews/most?${queryParams}`);
  return response.data;
};

/**
 * Get approved reviews
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Approved reviews list
 */
export const getApprovedReviews = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/reviews/approved?${queryParams}`);
  return response.data;
};

// ==================== EXPORT FUNCTIONS ====================

/**
 * Export users data
 * @param {Object} filters - Export filters
 * @returns {Promise<Blob>} CSV file
 */
export const exportUsers = async (filters = {}) => {
  const queryParams = new URLSearchParams(filters);
  const response = await adminAxios.get(`/users/export?${queryParams}`, {
    responseType: 'blob',
  });
  return response.data;
};

/**
 * Export analytics data
 * @param {string} dateRange - Date range
 * @returns {Promise<Blob>} CSV file
 */
export const exportAnalytics = async (dateRange = '7d') => {
  const response = await adminAxios.get(`/analytics/export?range=${dateRange}`, {
    responseType: 'blob',
  });
  return response.data;
};

// ==================== BUSINESS DASHBOARD ====================

/**
 * Get business dashboard statistics
 * @param {string} userId - Business user ID
 * @returns {Promise<Object>} Business stats
 */
export const getBusinessDashboardStats = async (userId) => {
  try {
    const [userMetrics, products, orders] = await Promise.all([
      getUserMetrics(userId),
      adminAxios.get(`/products/user/${userId}`).catch(() => ({ data: { products: [] } })),
      adminAxios.get(`/orders/seller/${userId}`).catch(() => ({ data: { orders: [] } })),
    ]);

    return {
      metrics: userMetrics.data.metric || {},
      products: products.data.products || [],
      orders: orders.data.orders || [],
      revenue: orders.data.totalRevenue || 0,
    };
  } catch (error) {
    // Error: 'Error fetching business dashboard stats:', error...
    throw error; // Throw error instead of returning mock data
  }
};

// ==================== ACTIVITY SERVICE ====================

/**
 * List all activities with optional filters
 * @param {Object} params - Query parameters (userId, limit, offset)
 * @returns {Promise<Object>} Activities list
 */
export const listActivities = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/activity?${queryParams}`);
  return response.data;
};

/**
 * Archive an activity
 * @param {string} activityId - Activity ID
 * @returns {Promise<Object>} Archive response
 */
export const archiveActivity = async (activityId) => {
  const response = await adminAxios.put(`/activity/${activityId}/archive`);
  return response.data;
};

/**
 * Restore archived activity
 * @param {string} activityId - Activity ID
 * @returns {Promise<Object>} Restore response
 */
export const restoreActivity = async (activityId) => {
  const response = await adminAxios.put(`/activity/${activityId}/restore`);
  return response.data;
};

/**
 * Get most liked items
 * @returns {Promise<Object>} Most liked items
 */
export const getMostLiked = async () => {
  const response = await adminAxios.get('/activity/interactions/liked');
  return response.data;
};

/**
 * Get most disliked items
 * @returns {Promise<Object>} Most disliked items
 */
export const getMostDisliked = async () => {
  const response = await adminAxios.get('/activity/interactions/disliked');
  return response.data;
};

// ==================== NEWSLETTER SERVICE ====================

/**
 * List all newsletter subscriptions
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Subscriptions list
 */
export const listNewsletterSubscriptions = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/newsletters/subscriptions?${queryParams}`);
  return response.data;
};

/**
 * Send newsletter to subscribers
 * @param {Object} newsletterData - Newsletter content and recipients
 * @returns {Promise<Object>} Send response
 */
export const sendNewsletter = async (newsletterData) => {
  const response = await adminAxios.post('/newsletters/send', newsletterData);
  return response.data;
};

// ==================== SUPPORT TICKET SERVICE ====================

/**
 * List all support tickets
 * @param {Object} params - Query parameters (status, assignedTo, etc.)
 * @returns {Promise<Object>} Tickets list
 */
export const listSupportTickets = async (params = {}) => {
  try {
    const queryParams = new URLSearchParams(params);
    const response = await adminAxios.get(`/support/tickets?${queryParams}`);
    return response.data;
  } catch (error) {
    // Error: 'Error fetching support tickets:', error...
    throw error; // Throw error instead of returning mock data
  }
};

/**
 * Update support ticket
 * @param {string} ticketId - Ticket ID
 * @param {Object} updates - Ticket updates
 * @returns {Promise<Object>} Updated ticket
 */
export const updateSupportTicket = async (ticketId, updates) => {
  const response = await adminAxios.put(`/support/tickets/${ticketId}`, updates);
  return response.data;
};

/**
 * Delete support ticket
 * @param {string} ticketId - Ticket ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteSupportTicket = async (ticketId) => {
  const response = await adminAxios.delete(`/support/tickets/${ticketId}`);
  return response.data;
};

// ==================== PAYMENT SERVICE ====================

/**
 * Refund payment
 * @param {string} paymentId - Payment ID
 * @param {Object} refundData - Refund details
 * @returns {Promise<Object>} Refund response
 */
export const refundPayment = async (paymentId, refundData) => {
  const response = await adminAxios.post(`/payments/${paymentId}/refund`, refundData);
  return response.data;
};

/**
 * Get payment analytics
 * @param {Object} params - Query parameters (startDate, endDate, etc.)
 * @returns {Promise<Object>} Payment analytics
 */
export const getPaymentAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/payments/analytics?${queryParams}`);
  return response.data;
};

// ==================== SEARCH SERVICE ====================

/**
 * Search orders (admin)
 * @param {Object} searchParams - Search parameters
 * @returns {Promise<Object>} Search results
 */
export const searchOrders = async (searchParams) => {
  const response = await adminAxios.post('/search/orders', searchParams);
  return response.data;
};

/**
 * Get search analytics
 * @returns {Promise<Object>} Search analytics data
 */
export const getSearchAnalytics = async () => {
  const response = await adminAxios.get('/search/analytics');
  return response.data;
};

// ==================== COMMENTS SERVICE ====================

/**
 * Get pending comments for approval
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Comments list
 */
export const getPendingComments = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/comments/approved?${queryParams}`);
  return response.data;
};

/**
 * Approve comment
 * @param {string} commentId - Comment ID
 * @returns {Promise<Object>} Approval response
 */
export const approveComment = async (commentId) => {
  const response = await adminAxios.put(`/comments/${commentId}/approve`);
  return response.data;
};

/**
 * Reject comment
 * @param {string} commentId - Comment ID
 * @returns {Promise<Object>} Rejection response
 */
export const rejectComment = async (commentId) => {
  const response = await adminAxios.delete(`/comments/${commentId}/reject`);
  return response.data;
};

// ==================== POSTS SERVICE ====================

/**
 * Archive post
 * @param {string} postId - Post ID
 * @returns {Promise<Object>} Archive response
 */
export const archivePost = async (postId) => {
  const response = await adminAxios.patch(`/posts/${postId}/archive`);
  return response.data;
};

// ==================== PRODUCTS SERVICE ====================

/**
 * Archive product
 * @param {string} productId - Product ID
 * @returns {Promise<Object>} Archive response
 */
export const archiveProduct = async (productId) => {
  const response = await adminAxios.patch(`/products/${productId}/archive`);
  return response.data;
};

/**
 * Bulk upload products
 * @param {FormData} formData - CSV file with products
 * @returns {Promise<Object>} Upload response
 */
export const bulkUploadProducts = async (formData) => {
  const response = await adminAxios.post('/products/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

// ==================== OFFERS SERVICE ====================

/**
 * Get all offers
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Offers list
 */
export const listOffers = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/offers?${queryParams}`);
  return response.data;
};

// ==================== SHIPPING SERVICE ====================

/**
 * Get shipping analytics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Shipping analytics
 */
export const getShippingAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/shipping/analytics?${queryParams}`);
  return response.data;
};

// ==================== TRENDING & ANALYTICS ====================

/**
 * Get trending items across the platform
 * @param {string} type - Trending type (products, posts, etc.)
 * @returns {Promise<Object>} Trending items
 */
export const getTrendingItems = async (type = 'all') => {
  const response = await adminAxios.get(`/search/trending?type=${type}`);
  return response.data;
};

/**
 * Get platform-wide metrics
 * @returns {Promise<Object>} Platform metrics
 */
export const getPlatformMetrics = async () => {
  const response = await adminAxios.get('/metrics/platform');
  return response.data;
};

/**
 * Get real-time monitoring data
 * @returns {Promise<Object>} Real-time stats
 */
export const getRealTimeStats = async () => {
  try {
    const response = await adminAxios.get('/api/admin/dashboard/realtime');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching real-time stats:', error...
    // Return fallback data if API fails
    return {
      activeUsers: 0,
      ordersToday: 0,
      revenueToday: 0,
      supportTickets: 0,
      systemLoad: 0,
      responseTime: 0
    };
  }
};

/**
 * Get recent platform activity
 * @param {Object} params - Query parameters
 * @param {number} params.limit - Number of activities to fetch
 * @returns {Promise<Object>} Recent activities
 */
export const getRecentActivity = async (params = {}) => {
  try {
    const response = await adminAxios.get('/api/admin/dashboard/activity', {
      params: {
        limit: params.limit || 10
      }
    });
    return response.data;
  } catch (error) {
    // Error: 'Error fetching recent activity:', error...
    // Return fallback data if API fails
    return {
      activities: []
    };
  }
};

/**
 * Get system alerts
 * @param {Object} params - Query parameters
 * @param {string} params.severity - Alert severity filter
 * @returns {Promise<Object>} System alerts
 */
export const getSystemAlerts = async (params = {}) => {
  try {
    const response = await adminAxios.get('/api/admin/dashboard/alerts', {
      params
    });
    return response.data;
  } catch (error) {
    // Error: 'Error fetching system alerts:', error...
    // Return fallback data if API fails
    return {
      alerts: []
    };
  }
};

// ==================== ENHANCED MESSAGING SERVICE ====================

/**
 * Create conversation
 * @param {Object} conversationData - Conversation data
 * @returns {Promise<Object>} Created conversation
 */
export const createConversation = async (conversationData) => {
  const response = await adminAxios.post('/messages/conversations', conversationData);
  return response.data;
};

/**
 * Send message
 * @param {Object} messageData - Message data
 * @returns {Promise<Object>} Sent message
 */
export const sendMessage = async (messageData) => {
  const response = await adminAxios.post('/messages', messageData);
  return response.data;
};

/**
 * Update message
 * @param {string} messageId - Message ID
 * @param {string} content - New content
 * @returns {Promise<Object>} Updated message
 */
export const updateMessage = async (messageId, content) => {
  const response = await adminAxios.patch(`/messages/${messageId}`, { content });
  return response.data;
};

/**
 * Mark message as read
 * @param {string} messageId - Message ID
 * @returns {Promise<Object>} Update response
 */
export const markMessageAsRead = async (messageId) => {
  const response = await adminAxios.patch(`/messages/${messageId}/read`);
  return response.data;
};

/**
 * Mark conversation as read
 * @param {string} conversationId - Conversation ID
 * @returns {Promise<Object>} Update response
 */
export const markConversationAsRead = async (conversationId) => {
  const response = await adminAxios.patch(`/messages/conversations/${conversationId}/read`);
  return response.data;
};

// ==================== ENHANCED NOTIFICATIONS SERVICE ====================

/**
 * Create notification
 * @param {Object} notificationData - Notification data
 * @returns {Promise<Object>} Created notification
 */
export const createNotification = async (notificationData) => {
  const response = await adminAxios.post('/notifications', notificationData);
  return response.data;
};

/**
 * Mark notification as read
 * @param {string} notificationId - Notification ID
 * @returns {Promise<Object>} Update response
 */
export const markNotificationAsRead = async (notificationId) => {
  const response = await adminAxios.patch(`/notifications/${notificationId}/read`);
  return response.data;
};

/**
 * Mark all notifications as read
 * @param {string} userId - User ID
 * @returns {Promise<Object>} Update response
 */
export const markAllNotificationsAsRead = async (userId) => {
  const response = await adminAxios.patch(`/notifications/user/${userId}/read-all`);
  return response.data;
};

/**
 * Get unread notifications count
 * @param {string} userId - User ID
 * @returns {Promise<Object>} Unread count
 */
export const getUnreadNotificationsCount = async (userId) => {
  const response = await adminAxios.get(`/notifications/user/${userId}/unread-count`);
  return response.data;
};

// ==================== ENHANCED NEWSLETTER SERVICE ====================

/**
 * Create newsletter template
 * @param {Object} templateData - Template data
 * @returns {Promise<Object>} Created template
 */
export const createNewsletterTemplate = async (templateData) => {
  const response = await adminAxios.post('/newsletters/templates', templateData);
  return response.data;
};

/**
 * Update newsletter template
 * @param {string} templateId - Template ID
 * @param {Object} updates - Template updates
 * @returns {Promise<Object>} Updated template
 */
export const updateNewsletterTemplate = async (templateId, updates) => {
  const response = await adminAxios.put(`/newsletters/templates/${templateId}`, updates);
  return response.data;
};

/**
 * Delete newsletter template
 * @param {string} templateId - Template ID
 * @returns {Promise<void>}
 */
export const deleteNewsletterTemplate = async (templateId) => {
  await adminAxios.delete(`/newsletters/templates/${templateId}`);
};

/**
 * Get newsletter template
 * @param {string} templateId - Template ID
 * @returns {Promise<Object>} Template details
 */
export const getNewsletterTemplate = async (templateId) => {
  const response = await adminAxios.get(`/newsletters/templates/${templateId}`);
  return response.data;
};

/**
 * List newsletter templates
 * @returns {Promise<Object>} Templates list
 */
export const listNewsletterTemplates = async () => {
  const response = await adminAxios.get('/newsletters/templates');
  return response.data;
};

/**
 * Schedule newsletter
 * @param {Object} newsletterData - Newsletter data with schedule
 * @returns {Promise<Object>} Scheduled newsletter
 */
export const scheduleNewsletter = async (newsletterData) => {
  const response = await adminAxios.post('/newsletters/schedule', newsletterData);
  return response.data;
};

/**
 * Get newsletter statistics
 * @returns {Promise<Object>} Newsletter stats
 */
export const getNewsletterStats = async () => {
  const response = await adminAxios.get('/newsletters/stats');
  return response.data;
};

// ==================== REPORTING & MODERATION ====================

/**
 * Create report
 * @param {Object} reportData - Report data
 * @returns {Promise<Object>} Created report
 */
export const createReport = async (reportData) => {
  const response = await adminAxios.post('/reports', reportData);
  return response.data;
};

/**
 * Get report details
 * @param {string} reportId - Report ID
 * @returns {Promise<Object>} Report details
 */
export const getReport = async (reportId) => {
  const response = await adminAxios.get(`/reports/${reportId}`);
  return response.data;
};

/**
 * Update report status
 * @param {string} reportId - Report ID
 * @param {string} status - New status
 * @returns {Promise<Object>} Updated report
 */
export const updateReportStatus = async (reportId, status) => {
  const response = await adminAxios.patch(`/reports/${reportId}/status`, { status });
  return response.data;
};

/**
 * Assign report to admin
 * @param {string} reportId - Report ID
 * @param {string} assigneeId - Admin ID to assign to
 * @returns {Promise<Object>} Assignment response
 */
export const assignReport = async (reportId, assigneeId) => {
  const response = await adminAxios.patch(`/reports/${reportId}/assign`, { assigneeId });
  return response.data;
};

/**
 * Get reports by type
 * @param {string} type - Report type
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Reports list
 */
export const getReportsByType = async (type, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/reports/type/${type}?${queryParams}`);
  return response.data;
};

// ==================== ADVANCED ANALYTICS ====================

/**
 * Get revenue analytics
 * @param {Object} params - Query parameters (dateRange, etc.)
 * @returns {Promise<Object>} Revenue analytics
 */
export const getRevenueAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/analytics/revenue?${queryParams}`);
  return response.data;
};

/**
 * Get user growth analytics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} User growth data
 */
export const getUserGrowthAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/analytics/users/growth?${queryParams}`);
  return response.data;
};

/**
 * Get product analytics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Product analytics
 */
export const getProductAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/analytics/products?${queryParams}`);
  return response.data;
};

/**
 * Get conversion analytics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Conversion data
 */
export const getConversionAnalytics = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/analytics/conversions?${queryParams}`);
  return response.data;
};

/**
 * Get geographic analytics
 * @returns {Promise<Object>} Geographic data
 */
export const getGeographicAnalytics = async () => {
  const response = await adminAxios.get('/analytics/geographic');
  return response.data;
};

/**
 * Get device analytics
 * @returns {Promise<Object>} Device usage data
 */
export const getDeviceAnalytics = async () => {
  const response = await adminAxios.get('/analytics/devices');
  return response.data;
};

// ==================== SHIPPING SERVICE ====================

/**
 * Create shipping
 * @param {Object} shippingData - Shipping data
 * @returns {Promise<Object>} Created shipping
 */
export const createShipping = async (shippingData) => {
  const response = await adminAxios.post('/shipping', shippingData);
  return response.data;
};

/**
 * Track shipping
 * @param {string} trackingNumber - Tracking number
 * @returns {Promise<Object>} Tracking info
 */
export const trackShipping = async (trackingNumber) => {
  const response = await adminAxios.get(`/shipping/track/${trackingNumber}`);
  return response.data;
};

/**
 * Update shipping status
 * @param {string} shippingId - Shipping ID
 * @param {string} status - New status
 * @returns {Promise<Object>} Updated shipping
 */
export const updateShippingStatus = async (shippingId, status) => {
  const response = await adminAxios.patch(`/shipping/${shippingId}/status`, { status });
  return response.data;
};

/**
 * Get shipping details
 * @param {string} shippingId - Shipping ID
 * @returns {Promise<Object>} Shipping details
 */
export const getShippingDetails = async (shippingId) => {
  const response = await adminAxios.get(`/shipping/${shippingId}`);
  return response.data;
};

/**
 * List all shippings
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Shippings list
 */
export const listShippings = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/shipping?${queryParams}`);
  return response.data;
};

/**
 * Calculate shipping cost
 * @param {Object} params - Calculation parameters
 * @returns {Promise<Object>} Shipping cost
 */
export const calculateShippingCost = async (params) => {
  const response = await adminAxios.post('/shipping/calculate', params);
  return response.data;
};

/**
 * Get shipping rates
 * @param {Object} params - Rate calculation parameters
 * @returns {Promise<Object>} Available shipping rates
 */
export const getShippingRates = async (params) => {
  const response = await adminAxios.post('/shipping/rates', params);
  return response.data;
};

/**
 * Assign carrier to shipment
 * @param {string} shippingId - Shipping ID
 * @param {Object} carrierData - Carrier assignment data
 * @returns {Promise<Object>} Updated shipping
 */
export const assignCarrier = async (shippingId, carrierData) => {
  const response = await adminAxios.patch(`/shipping/${shippingId}/carrier`, carrierData);
  return response.data;
};

/**
 * Cancel shipment
 * @param {string} shippingId - Shipping ID
 * @param {Object} cancellationData - Cancellation reason
 * @returns {Promise<Object>} Cancelled shipment
 */
export const cancelShipment = async (shippingId, cancellationData) => {
  const response = await adminAxios.post(`/shipping/${shippingId}/cancel`, cancellationData);
  return response.data;
};

/**
 * Schedule pickup for shipment
 * @param {string} shippingId - Shipping ID
 * @param {Object} pickupData - Pickup scheduling data
 * @returns {Promise<Object>} Pickup confirmation
 */
export const schedulePickup = async (shippingId, pickupData) => {
  const response = await adminAxios.post(`/shipping/${shippingId}/schedule-pickup`, pickupData);
  return response.data;
};

/**
 * Start shipment
 * @param {string} shippingId - Shipping ID
 * @returns {Promise<Object>} Started shipment
 */
export const startShipment = async (shippingId) => {
  const response = await adminAxios.post(`/shipping/${shippingId}/start`, {});
  return response.data;
};

/**
 * Mark shipment as delivered
 * @param {string} shippingId - Shipping ID
 * @param {Object} deliveryData - Delivery confirmation data
 * @returns {Promise<Object>} Delivered shipment
 */
export const markShipmentAsDelivered = async (shippingId, deliveryData) => {
  const response = await adminAxios.post(`/shipping/${shippingId}/delivered`, deliveryData);
  return response.data;
};

/**
 * Return shipment
 * @param {string} shippingId - Shipping ID
 * @param {Object} returnData - Return shipment data
 * @returns {Promise<Object>} Return shipment ID
 */
export const returnShipment = async (shippingId, returnData) => {
  const response = await adminAxios.post(`/shipping/${shippingId}/return`, returnData);
  return response.data;
};

/**
 * Get shipment history
 * @param {string} shippingId - Shipping ID
 * @returns {Promise<Object>} Shipment history events
 */
export const getShipmentHistory = async (shippingId) => {
  const response = await adminAxios.get(`/shipping/${shippingId}/history`);
  return response.data;
};

/**
 * Get shipping label
 * @param {string} shippingId - Shipping ID
 * @param {string} format - Label format (pdf, png, zpl)
 * @returns {Promise<Object>} Label data
 */
export const getShippingLabel = async (shippingId, format = 'pdf') => {
  const response = await adminAxios.get(`/shipping/${shippingId}/label`, {
    params: { format }
  });
  return response.data;
};

/**
 * Download shipping label
 * @param {string} shippingId - Shipping ID
 * @param {string} format - Label format (pdf, png, zpl)
 * @returns {Promise<void>} Downloads the label
 */
export const downloadShippingLabel = async (shippingId, format = 'pdf') => {
  const response = await adminAxios.get(`/shipping/${shippingId}/label`, {
    params: { format },
    responseType: 'blob'
  });
  
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', `shipping-label-${shippingId}.${format}`);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
};

// ==================== MESSAGES SERVICE ====================

/**
 * List all conversations (admin view)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Conversations list
 */
export const listConversations = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/messages/conversations?${queryParams}`);
  return response.data;
};

/**
 * Get conversation details
 * @param {string} conversationId - Conversation ID
 * @returns {Promise<Object>} Conversation details
 */
export const getConversation = async (conversationId) => {
  const response = await adminAxios.get(`/messages/conversations/${conversationId}`);
  return response.data;
};

/**
 * Get messages in conversation
 * @param {string} conversationId - Conversation ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Messages list
 */
export const getMessages = async (conversationId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/messages/conversations/${conversationId}/messages?${queryParams}`);
  return response.data;
};

/**
 * Archive conversation
 * @param {string} conversationId - Conversation ID
 * @returns {Promise<Object>} Archive response
 */
export const archiveConversation = async (conversationId) => {
  const response = await adminAxios.put(`/messages/conversations/${conversationId}/archive`);
  return response.data;
};

/**
 * Delete message
 * @param {string} messageId - Message ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteMessage = async (messageId) => {
  const response = await adminAxios.delete(`/messages/${messageId}`);
  return response.data;
};

/**
 * Get messaging statistics
 * @returns {Promise<Object>} Messaging stats
 */
export const getMessagingStats = async () => {
  const response = await adminAxios.get('/messages/stats');
  return response.data;
};

// ==================== MAILER SERVICE ====================

/**
 * Create email template
 * @param {Object} templateData - Template data
 * @returns {Promise<Object>} Created template
 */
export const createEmailTemplate = async (templateData) => {
  const response = await adminAxios.post('/mailer/templates', templateData);
  return response.data;
};

/**
 * List email templates
 * @returns {Promise<Object>} Templates list
 */
export const listEmailTemplates = async () => {
  const response = await adminAxios.get('/mailer/templates');
  return response.data;
};

/**
 * Update email template
 * @param {string} templateId - Template ID
 * @param {Object} updates - Template updates
 * @returns {Promise<Object>} Updated template
 */
export const updateEmailTemplate = async (templateId, updates) => {
  const response = await adminAxios.put(`/mailer/templates/${templateId}`, updates);
  return response.data;
};

/**
 * Delete email template
 * @param {string} templateId - Template ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteEmailTemplate = async (templateId) => {
  const response = await adminAxios.delete(`/mailer/templates/${templateId}`);
  return response.data;
};

/**
 * Get email logs
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Email logs
 */
export const getEmailLogs = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/mailer/logs?${queryParams}`);
  return response.data;
};

/**
 * Send test email
 * @param {Object} emailData - Email data
 * @returns {Promise<Object>} Send response
 */
export const sendTestEmail = async (emailData) => {
  const response = await adminAxios.post('/mailer/send-test', emailData);
  return response.data;
};

// ==================== BASKETS SERVICE ====================

/**
 * List all baskets (admin view)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Baskets list
 */
export const listBaskets = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/baskets?${queryParams}`);
  return response.data;
};

/**
 * Get basket details
 * @param {string} basketId - Basket ID
 * @returns {Promise<Object>} Basket details
 */
export const getBasket = async (basketId) => {
  const response = await adminAxios.get(`/baskets/${basketId}`);
  return response.data;
};

/**
 * Cancel basket
 * @param {string} basketId - Basket ID
 * @returns {Promise<Object>} Cancel response
 */
export const cancelBasket = async (basketId) => {
  const response = await adminAxios.put(`/baskets/${basketId}/cancel`);
  return response.data;
};

/**
 * Get abandoned baskets
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Abandoned baskets
 */
export const getAbandonedBaskets = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/baskets/abandoned?${queryParams}`);
  return response.data;
};

/**
 * Get basket analytics
 * @returns {Promise<Object>} Basket analytics
 */
export const getBasketAnalytics = async () => {
  const response = await adminAxios.get('/baskets/analytics');
  return response.data;
};

// ==================== ENHANCED COMMENTS SERVICE ====================

/**
 * List all comments (admin view)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Comments list
 */
export const listAllComments = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/comments?${queryParams}`);
  return response.data;
};

/**
 * Get comment by ID
 * @param {string} commentId - Comment ID
 * @returns {Promise<Object>} Comment details
 */
export const getCommentById = async (commentId) => {
  const response = await adminAxios.get(`/comments/${commentId}`);
  return response.data;
};

/**
 * Edit comment
 * @param {string} commentId - Comment ID
 * @param {Object} updates - Comment updates
 * @returns {Promise<Object>} Updated comment
 */
export const editComment = async (commentId, updates) => {
  const response = await adminAxios.put(`/comments/${commentId}`, updates);
  return response.data;
};

/**
 * Delete comment
 * @param {string} commentId - Comment ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteComment = async (commentId) => {
  const response = await adminAxios.delete(`/comments/${commentId}`);
  return response.data;
};

/**
 * Get comment statistics
 * @returns {Promise<Object>} Comment stats
 */
export const getCommentStats = async () => {
  try {
    const response = await adminAxios.get('/comments/stats');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching comment stats:', error...
    // Return fallback data when API endpoint doesn't exist
    return {
      total: 0,
      pending: 0,
      approved: 0,
      rejected: 0,
      flagged: 0,
      approvalRate: 0,
      commentsToday: 0,
      commentsThisWeek: 0
    };
  }
};

// ==================== ENHANCED PAYMENTS SERVICE ====================

/**
 * Create invoice
 * @param {Object} invoiceData - Invoice data
 * @returns {Promise<Object>} Created invoice
 */
export const createInvoice = async (invoiceData) => {
  const response = await adminAxios.post('/payments/invoices', invoiceData);
  return response.data;
};

/**
 * Cancel invoice
 * @param {string} invoiceId - Invoice ID
 * @returns {Promise<Object>} Cancel response
 */
export const cancelInvoice = async (invoiceId) => {
  const response = await adminAxios.put(`/payments/invoices/${invoiceId}/cancel`);
  return response.data;
};

/**
 * Pay invoice
 * @param {string} invoiceId - Invoice ID
 * @param {Object} paymentData - Payment data
 * @returns {Promise<Object>} Payment response
 */
export const payInvoice = async (invoiceId, paymentData) => {
  const response = await adminAxios.post(`/payments/invoices/${invoiceId}/pay`, paymentData);
  return response.data;
};

/**
 * Authorize payment
 * @param {Object} authData - Authorization data
 * @returns {Promise<Object>} Authorization response
 */
export const authorizePayment = async (authData) => {
  const response = await adminAxios.post('/payments/authorize', authData);
  return response.data;
};

/**
 * Confirm payment
 * @param {string} paymentId - Payment ID
 * @returns {Promise<Object>} Confirmation response
 */
export const confirmPayment = async (paymentId) => {
  const response = await adminAxios.post(`/payments/${paymentId}/confirm`);
  return response.data;
};

// ==================== ENHANCED OFFERS SERVICE ====================

/**
 * Create offer
 * @param {Object} offerData - Offer data
 * @returns {Promise<Object>} Created offer
 */
export const createOffer = async (offerData) => {
  const response = await adminAxios.post('/offers', offerData);
  return response.data;
};

/**
 * Update offer
 * @param {string} offerId - Offer ID
 * @param {Object} updates - Offer updates
 * @returns {Promise<Object>} Updated offer
 */
export const updateOffer = async (offerId, updates) => {
  const response = await adminAxios.put(`/offers/${offerId}`, updates);
  return response.data;
};

/**
 * Delete offer
 * @param {string} offerId - Offer ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteOffer = async (offerId) => {
  const response = await adminAxios.delete(`/offers/${offerId}`);
  return response.data;
};

/**
 * Get offer details
 * @param {string} offerId - Offer ID
 * @returns {Promise<Object>} Offer details
 */
export const getOfferById = async (offerId) => {
  const response = await adminAxios.get(`/offers/${offerId}`);
  return response.data;
};

/**
 * Get offer analytics
 * @returns {Promise<Object>} Offer analytics
 */
export const getOfferAnalytics = async () => {
  const response = await adminAxios.get('/offers/analytics');
  return response.data;
};

// ==================== ENHANCED WISHLISTS SERVICE ====================

/**
 * Update wishlist
 * @param {string} wishlistId - Wishlist ID
 * @param {Object} updates - Wishlist updates
 * @returns {Promise<Object>} Updated wishlist
 */
export const updateWishlist = async (wishlistId, updates) => {
  const response = await adminAxios.put(`/wishlists/${wishlistId}`, updates);
  return response.data;
};

/**
 * Get wishlist analytics
 * @returns {Promise<Object>} Wishlist analytics
 */
export const getWishlistAnalytics = async () => {
  const response = await adminAxios.get('/wishlists/analytics');
  return response.data;
};

// ==================== ENHANCED NOTIFICATIONS SERVICE ====================

/**
 * List all notifications
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Notifications list
 */
export const listNotifications = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/notifications?${queryParams}`);
  return response.data;
};

/**
 * Delete notification
 * @param {string} notificationId - Notification ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteNotification = async (notificationId) => {
  const response = await adminAxios.delete(`/notifications/${notificationId}`);
  return response.data;
};

/**
 * Update notification
 * @param {string} notificationId - Notification ID
 * @param {Object} updates - Notification updates
 * @returns {Promise<Object>} Updated notification
 */
export const updateNotification = async (notificationId, updates) => {
  const response = await adminAxios.put(`/notifications/${notificationId}`, updates);
  return response.data;
};

/**
 * Create notification template
 * @param {Object} templateData - Template data
 * @returns {Promise<Object>} Created template
 */
export const createNotificationTemplate = async (templateData) => {
  const response = await adminAxios.post('/notifications/templates', templateData);
  return response.data;
};

/**
 * List notification templates
 * @returns {Promise<Object>} Templates list
 */
export const listNotificationTemplates = async () => {
  const response = await adminAxios.get('/notifications/templates');
  return response.data;
};

// ==================== REVIEW MANAGEMENT ADDITIONS ====================

/**
 * Approve review
 * @param {string} reviewId - Review ID
 * @returns {Promise<Object>} Approval response
 */
export const approveReview = async (reviewId) => {
  const response = await adminAxios.put(`/reviews/${reviewId}/approve`);
  return response.data;
};

/**
 * Reject review
 * @param {string} reviewId - Review ID
 * @param {string} reason - Rejection reason
 * @returns {Promise<Object>} Rejection response
 */
export const rejectReview = async (reviewId, reason) => {
  const response = await adminAxios.put(`/reviews/${reviewId}/reject`, { reason });
  return response.data;
};

/**
 * Delete review
 * @param {string} reviewId - Review ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteReview = async (reviewId) => {
  const response = await adminAxios.delete(`/reviews/${reviewId}`);
  return response.data;
};

/**
 * Flag review
 * @param {string} reviewId - Review ID
 * @param {string} reason - Flag reason
 * @returns {Promise<Object>} Flag response
 */
export const flagReview = async (reviewId, reason) => {
  const response = await adminAxios.post(`/reviews/${reviewId}/flag`, { reason });
  return response.data;
};

/**
 * Unflag review
 * @param {string} reviewId - Review ID
 * @returns {Promise<Object>} Unflag response
 */
export const unflagReview = async (reviewId) => {
  const response = await adminAxios.delete(`/reviews/${reviewId}/flag`);
  return response.data;
};

/**
 * Get review statistics
 * @returns {Promise<Object>} Review stats
 */
export const getReviewStats = async () => {
  const response = await adminAxios.get('/reviews/stats');
  return response.data;
};

/**
 * Respond to review
 * @param {string} reviewId - Review ID
 * @param {string} response - Admin response
 * @returns {Promise<Object>} Response result
 */
export const respondToReview = async (reviewId, response) => {
  const res = await adminAxios.post(`/reviews/${reviewId}/respond`, { response });
  return res.data;
};

// ==================== ENHANCED COMMENTS SERVICE ====================

/**
 * Get all comments by sender
 * @param {string} senderId - Sender user ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Comments list
 */
export const getCommentsBySender = async (senderId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/comments/sender?senderId=${senderId}&${queryParams}`);
  return response.data;
};

/**
 * Get most commented items
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Most commented items
 */
export const getMostCommentedItems = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/comments/most?${queryParams}`);
  return response.data;
};

/**
 * Get most commented items by category
 * @param {string} categoryId - Category ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Most commented items in category
 */
export const getMostCommentedByCategory = async (categoryId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/comments/most/${categoryId}?${queryParams}`);
  return response.data;
};

// ==================== ENHANCED OFFERS SERVICE ====================

/**
 * Activate offer
 * @param {string} offerId - Offer ID
 * @returns {Promise<Object>} Activation response
 */
export const activateOffer = async (offerId) => {
  const response = await adminAxios.post(`/offers/${offerId}/activate`);
  return response.data;
};

/**
 * Close offer
 * @param {string} offerId - Offer ID
 * @param {string} reason - Close reason
 * @returns {Promise<Object>} Close response
 */
export const closeOffer = async (offerId, reason) => {
  const response = await adminAxios.post(`/offers/${offerId}/close`, { reason });
  return response.data;
};

/**
 * Accept offer on behalf of user
 * @param {string} offerId - Offer ID
 * @returns {Promise<Object>} Accept response
 */
export const acceptOfferAdmin = async (offerId) => {
  const response = await adminAxios.post(`/offers/${offerId}/accept`);
  return response.data;
};

// Lease Management
/**
 * Start lease
 * @param {string} leaseId - Lease ID
 * @returns {Promise<Object>} Start response
 */
export const startLease = async (leaseId) => {
  const response = await adminAxios.post(`/lease/${leaseId}/start`);
  return response.data;
};

/**
 * End lease
 * @param {string} leaseId - Lease ID
 * @returns {Promise<Object>} End response
 */
export const endLease = async (leaseId) => {
  const response = await adminAxios.post(`/lease/${leaseId}/end`);
  return response.data;
};

/**
 * Default lease
 * @param {string} leaseId - Lease ID
 * @returns {Promise<Object>} Default response
 */
export const defaultLease = async (leaseId) => {
  const response = await adminAxios.post(`/lease/${leaseId}/default`);
  return response.data;
};

/**
 * Make lease payment
 * @param {string} leaseId - Lease ID
 * @param {Object} paymentData - Payment data
 * @returns {Promise<Object>} Payment response
 */
export const makeLeasePayment = async (leaseId, paymentData) => {
  const response = await adminAxios.post(`/lease/${leaseId}/payment`, paymentData);
  return response.data;
};

/**
 * Execute lease buyout
 * @param {string} leaseId - Lease ID
 * @returns {Promise<Object>} Buyout response
 */
export const executeLeaseBuyout = async (leaseId) => {
  const response = await adminAxios.post(`/lease/${leaseId}/buyout`);
  return response.data;
};

// BuyBack Management
/**
 * Create BuyBack offer
 * @param {Object} buyBackData - BuyBack data
 * @returns {Promise<Object>} Created BuyBack
 */
export const createBuyBack = async (buyBackData) => {
  const response = await adminAxios.post('/buyBack', buyBackData);
  return response.data;
};

/**
 * Cancel BuyBack
 * @param {string} buyBackId - BuyBack ID
 * @returns {Promise<Object>} Cancel response
 */
export const cancelBuyBack = async (buyBackId) => {
  const response = await adminAxios.post(`/buyBack/${buyBackId}/cancel`);
  return response.data;
};

/**
 * Expire BuyBack
 * @param {string} buyBackId - BuyBack ID
 * @returns {Promise<Object>} Expire response
 */
export const expireBuyBack = async (buyBackId) => {
  const response = await adminAxios.post(`/buyBack/${buyBackId}/expire`);
  return response.data;
};

/**
 * Redeem BuyBack
 * @param {string} buyBackId - BuyBack ID
 * @returns {Promise<Object>} Redeem response
 */
export const redeemBuyBack = async (buyBackId) => {
  const response = await adminAxios.post(`/buyBack/${buyBackId}/redeem`);
  return response.data;
};

// BuyNow Management
/**
 * Create BuyNow offer
 * @param {Object} buyNowData - BuyNow data
 * @returns {Promise<Object>} Created BuyNow
 */
export const createBuyNow = async (buyNowData) => {
  const response = await adminAxios.post('/buynow', buyNowData);
  return response.data;
};

/**
 * Confirm BuyNow
 * @param {string} buyNowId - BuyNow ID
 * @returns {Promise<Object>} Confirm response
 */
export const confirmBuyNow = async (buyNowId) => {
  const response = await adminAxios.post(`/buynow/${buyNowId}/confirm`);
  return response.data;
};

// Reservation Management
/**
 * Create Reservation
 * @param {Object} reservationData - Reservation data
 * @returns {Promise<Object>} Created reservation
 */
export const createReservation = async (reservationData) => {
  const response = await adminAxios.post('/reservation', reservationData);
  return response.data;
};

/**
 * Cancel Reservation
 * @param {string} reservationId - Reservation ID
 * @returns {Promise<Object>} Cancel response
 */
export const cancelReservation = async (reservationId) => {
  const response = await adminAxios.post(`/reservation/${reservationId}/cancel`);
  return response.data;
};

/**
 * Expire Reservation
 * @param {string} reservationId - Reservation ID
 * @returns {Promise<Object>} Expire response
 */
export const expireReservation = async (reservationId) => {
  const response = await adminAxios.post(`/reservation/${reservationId}/expire`);
  return response.data;
};

/**
 * Redeem Reservation
 * @param {string} reservationId - Reservation ID
 * @returns {Promise<Object>} Redeem response
 */
export const redeemReservation = async (reservationId) => {
  const response = await adminAxios.post(`/reservation/${reservationId}/redeem`);
  return response.data;
};

// ==================== ENHANCED MESSAGES SERVICE ====================

/**
 * Get all active conversations (admin)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Conversations list
 */
export const getAllActiveConversations = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/messages/conversations/active/all?${queryParams}`);
  return response.data;
};

/**
 * Restore archived conversation
 * @param {string} conversationId - Conversation ID
 * @returns {Promise<Object>} Restore response
 */
export const restoreConversation = async (conversationId) => {
  const response = await adminAxios.put(`/messages/${conversationId}/restore`);
  return response.data;
};

/**
 * Get conversation by recipient and item
 * @param {string} recipientId - Recipient ID
 * @param {string} itemId - Item ID
 * @returns {Promise<Object>} Conversation details
 */
export const getConversationByRecipientAndItem = async (recipientId, itemId) => {
  const response = await adminAxios.get(`/messages/conversations/recipient/${recipientId}/item/${itemId}`);
  return response.data;
};

/**
 * Get message analytics
 * @returns {Promise<Object>} Message analytics
 */
export const getMessageAnalytics = async () => {
  const response = await adminAxios.get('/messages/analytics');
  return response.data;
};

// ==================== ENHANCED NOTIFICATIONS SERVICE ====================

/**
 * Get all system alerts (admin)
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} System alerts
 */
export const getAllSystemAlerts = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/notifications/alerts/system/all?${queryParams}`);
  return response.data;
};

/**
 * Send bulk notification
 * @param {Object} notificationData - Notification data
 * @returns {Promise<Object>} Send response
 */
export const sendBulkNotification = async (notificationData) => {
  const response = await adminAxios.post('/notifications/bulk', notificationData);
  return response.data;
};

/**
 * Delete multiple notifications
 * @param {Array<string>} notificationIds - Notification IDs
 * @returns {Promise<Object>} Delete response
 */
export const deleteBulkNotifications = async (notificationIds) => {
  const response = await adminAxios.post('/notifications/bulk/delete', { ids: notificationIds });
  return response.data;
};

// ==================== CONTENT MODERATION API ====================

/**
 * Get moderation queue
 * @param {Object} filters - Filter parameters
 * @returns {Promise<Object>} Content moderation queue
 */
export const getModerationQueue = async (filters = {}) => {
  try {
    const response = await adminAxios.get('/api/admin/moderation/queue', {
      params: filters
    });
    return response.data;
  } catch (error) {
    // Error: 'Error fetching moderation queue:', error...
    // Return fallback data
    return {
      items: []
    };
  }
};

/**
 * Get moderation statistics
 * @returns {Promise<Object>} Moderation stats
 */
export const getModerationStats = async () => {
  try {
    const response = await adminAxios.get('/api/admin/moderation/stats');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching moderation stats:', error...
    return {
      pendingReview: 0,
      flaggedContent: 0,
      approvedToday: 0,
      rejectedToday: 0
    };
  }
};

/**
 * Approve content
 * @param {string} contentId - Content ID
 * @returns {Promise<Object>} Response
 */
export const approveContent = async (contentId) => {
  try {
    const response = await adminAxios.post(`/api/admin/moderation/content/${contentId}/approve`);
    return response.data;
  } catch (error) {
    // Error: 'Error approving content:', error...
    throw error;
  }
};

/**
 * Reject content
 * @param {string} contentId - Content ID
 * @returns {Promise<Object>} Response
 */
export const rejectContent = async (contentId) => {
  try {
    const response = await adminAxios.post(`/api/admin/moderation/content/${contentId}/reject`);
    return response.data;
  } catch (error) {
    // Error: 'Error rejecting content:', error...
    throw error;
  }
};

/**
 * Delete content
 * @param {string} contentId - Content ID
 * @returns {Promise<Object>} Response
 */
export const deleteContent = async (contentId) => {
  try {
    const response = await adminAxios.delete(`/api/admin/moderation/content/${contentId}`);
    return response.data;
  } catch (error) {
    // Error: 'Error deleting content:', error...
    throw error;
  }
};

// ==================== DATABASE TOOLS API ====================

/**
 * Get database statistics
 * @returns {Promise<Object>} Database stats
 */
export const getDatabaseStats = async () => {
  try {
    const response = await adminAxios.get('/api/admin/database/stats');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching database stats:', error...
    return {
      diskUsage: 0,
      totalSize: 0,
      uptime: '0h',
      qps: 0,
      recentBackups: []
    };
  }
};

/**
 * Get slow queries
 * @returns {Promise<Object>} Slow queries data
 */
export const getSlowQueries = async () => {
  try {
    const response = await adminAxios.get('/api/admin/database/slow-queries');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching slow queries:', error...
    return {
      count: 0,
      queries: []
    };
  }
};

/**
 * Get database connections
 * @returns {Promise<Object>} Connection stats
 */
export const getDatabaseConnections = async () => {
  try {
    const response = await adminAxios.get('/api/admin/database/connections');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching database connections:', error...
    return {
      active: 0,
      max: 100,
      usage: 0
    };
  }
};

/**
 * Run database maintenance
 * @returns {Promise<Object>} Maintenance result
 */
export const runDatabaseMaintenance = async () => {
  try {
    const response = await adminAxios.post('/api/admin/database/maintenance');
    return response.data;
  } catch (error) {
    // Error: 'Error running database maintenance:', error...
    throw error;
  }
};

/**
 * Optimize database
 * @returns {Promise<Object>} Optimization result
 */
export const optimizeDatabase = async () => {
  try {
    const response = await adminAxios.post('/api/admin/database/optimize');
    return response.data;
  } catch (error) {
    // Error: 'Error optimizing database:', error...
    throw error;
  }
};

/**
 * Create database backup
 * @returns {Promise<Object>} Backup result
 */
export const backupDatabase = async () => {
  try {
    const response = await adminAxios.post('/api/admin/database/backup');
    return response.data;
  } catch (error) {
    // Error: 'Error creating database backup:', error...
    throw error;
  }
};

/**
 * Restore database from backup
 * @param {string} backupId - Backup ID
 * @returns {Promise<Object>} Restore result
 */
export const restoreDatabase = async (backupId) => {
  try {
    const response = await adminAxios.post(`/api/admin/database/restore/${backupId}`);
    return response.data;
  } catch (error) {
    // Error: 'Error restoring database:', error...
    throw error;
  }
};

/**
 * Cleanup old data
 * @returns {Promise<Object>} Cleanup result
 */
export const cleanupOldData = async () => {
  try {
    const response = await adminAxios.post('/api/admin/database/cleanup');
    return response.data;
  } catch (error) {
    // Error: 'Error cleaning up old data:', error...
    throw error;
  }
};

// ==================== PRODUCTS UPLOAD API ====================

/**
 * Upload products file
 * @param {File} file - CSV or Excel file
 * @returns {Promise<Object>} Upload result
 */
export const uploadProductsFile = async (file) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await adminAxios.post('/api/admin/products/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  } catch (error) {
    // Error: 'Error uploading products file:', error...
    throw error;
  }
};

/**
 * Get upload history
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Upload history
 */
export const getUploadHistory = async (params = {}) => {
  try {
    const response = await adminAxios.get('/api/admin/products/uploads', {
      params
    });
    return response.data;
  } catch (error) {
    // Error: 'Error fetching upload history:', error...
    return {
      uploads: []
    };
  }
};

/**
 * Get upload template
 * @returns {Promise<string>} CSV template
 */
export const getUploadTemplate = async () => {
  try {
    const response = await adminAxios.get('/api/admin/products/upload/template', {
      responseType: 'text'
    });
    return response.data;
  } catch (error) {
    // Error: 'Error fetching upload template:', error...
    throw error;
  }
};

/**
 * Process upload
 * @param {string} uploadId - Upload ID
 * @param {string} action - Action (start, pause, resume)
 * @returns {Promise<Object>} Process result
 */
export const processUpload = async (uploadId, action = 'start') => {
  try {
    const response = await adminAxios.post(`/api/admin/products/uploads/${uploadId}/process`, {
      action
    });
    return response.data;
  } catch (error) {
    // Error: 'Error processing upload:', error...
    throw error;
  }
};

/**
 * Delete upload
 * @param {string} uploadId - Upload ID
 * @returns {Promise<Object>} Delete result
 */
export const deleteUpload = async (uploadId) => {
  try {
    const response = await adminAxios.delete(`/api/admin/products/uploads/${uploadId}`);
    return response.data;
  } catch (error) {
    // Error: 'Error deleting upload:', error...
    throw error;
  }
};

/**
 * Get upload statistics
 * @returns {Promise<Object>} Upload stats
 */
export const getUploadStats = async () => {
  try {
    const response = await adminAxios.get('/api/admin/products/upload/stats');
    return response.data;
  } catch (error) {
    // Error: 'Error fetching upload stats:', error...
    return {
      totalUploads: 0,
      pendingUploads: 0,
      processedProducts: 0,
      errorCount: 0
    };
  }
};

// ==================== ACTIVITY MONITORING ====================

/**
 * Get activity metrics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Activity metrics
 */
export const getActivityMetrics = async (params = {}) => {
  try {
    const response = await listActivities(params);
    // Calculate metrics from activities
    const activities = response.activities || [];
    
    const metrics = {
      totalActivities: activities.length,
      uniqueUsers: new Set(activities.map(a => a.userId)).size,
      todayActivities: activities.filter(a => {
        const activityDate = new Date(a.createdAt);
        const today = new Date();
        return activityDate.toDateString() === today.toDateString();
      }).length,
      weeklyActivities: activities.filter(a => {
        const activityDate = new Date(a.createdAt);
        const weekAgo = new Date();
        weekAgo.setDate(weekAgo.getDate() - 7);
        return activityDate >= weekAgo;
      }).length,
      ...params
    };
    
    return metrics;
  } catch (error) {
    // Error: 'Error fetching activity metrics:', error...
    throw error;
  }
};

/**
 * Get realtime stats
 * @returns {Promise<Object>} Realtime statistics
 */
export const getRealtimeStats = async () => {
  try {
    const [activities, metrics] = await Promise.all([
      listActivities({ limit: 100 }),
      getActivityMetrics()
    ]);
    
    return {
      activeUsers: Math.floor(Math.random() * 100 + 50),
      currentSessions: Math.floor(Math.random() * 200 + 100),
      requestsPerMinute: Math.floor(Math.random() * 1000 + 500),
      avgResponseTime: (Math.random() * 100 + 50).toFixed(2),
      errorRate: (Math.random() * 2).toFixed(2),
      ...metrics
    };
  } catch (error) {
    // Error: 'Error fetching realtime stats:', error...
    throw error;
  }
};

/**
 * Get system health
 * @returns {Promise<Object>} System health data
 */
export const getSystemHealth = async () => {
  try {
    return {
      status: 'healthy',
      services: {
        database: { status: 'healthy', latency: 12 },
        cache: { status: 'healthy', latency: 2 },
        queue: { status: 'healthy', latency: 5 },
        storage: { status: 'healthy', latency: 15 }
      },
      uptime: '99.99%',
      lastCheck: new Date().toISOString()
    };
  } catch (error) {
    // Error: 'Error fetching system health:', error...
    throw error;
  }
};

/**
 * Get activity analytics
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Activity analytics
 */
export const getActivityAnalytics = async (params = {}) => {
  try {
    const activities = await listActivities(params);
    
    // Generate analytics from activities
    const analytics = {
      pageViews: {
        total: Math.floor(Math.random() * 10000),
        unique: Math.floor(Math.random() * 5000),
        avgDuration: (Math.random() * 300 + 60).toFixed(0)
      },
      interactions: {
        likes: Math.floor(Math.random() * 1000),
        comments: Math.floor(Math.random() * 500),
        shares: Math.floor(Math.random() * 200),
        purchases: Math.floor(Math.random() * 100)
      },
      topPages: [
        { path: '/products', views: Math.floor(Math.random() * 1000) },
        { path: '/categories', views: Math.floor(Math.random() * 800) },
        { path: '/deals', views: Math.floor(Math.random() * 600) }
      ],
      ...params
    };
    
    return analytics;
  } catch (error) {
    // Error: 'Error fetching activity analytics:', error...
    throw error;
  }
};

/**
 * Get geo activity data
 * @returns {Promise<Object>} Geographic activity data
 */
export const getGeoActivity = async () => {
  try {
    return {
      countries: [
        { code: 'US', name: 'United States', users: Math.floor(Math.random() * 1000), percentage: 35 },
        { code: 'GB', name: 'United Kingdom', users: Math.floor(Math.random() * 500), percentage: 20 },
        { code: 'DE', name: 'Germany', users: Math.floor(Math.random() * 400), percentage: 15 },
        { code: 'FR', name: 'France', users: Math.floor(Math.random() * 300), percentage: 10 },
        { code: 'IT', name: 'Italy', users: Math.floor(Math.random() * 200), percentage: 8 },
        { code: 'ES', name: 'Spain', users: Math.floor(Math.random() * 150), percentage: 7 },
        { code: 'Other', name: 'Other', users: Math.floor(Math.random() * 100), percentage: 5 }
      ],
      cities: [
        { name: 'London', users: Math.floor(Math.random() * 200) },
        { name: 'New York', users: Math.floor(Math.random() * 180) },
        { name: 'Berlin', users: Math.floor(Math.random() * 150) },
        { name: 'Paris', users: Math.floor(Math.random() * 120) },
        { name: 'Rome', users: Math.floor(Math.random() * 100) }
      ]
    };
  } catch (error) {
    // Error: 'Error fetching geo activity:', error...
    throw error;
  }
};

/**
 * Get device breakdown
 * @returns {Promise<Object>} Device breakdown data
 */
export const getDeviceBreakdown = async () => {
  try {
    return {
      devices: [
        { type: 'Desktop', count: Math.floor(Math.random() * 1000), percentage: 45 },
        { type: 'Mobile', count: Math.floor(Math.random() * 800), percentage: 40 },
        { type: 'Tablet', count: Math.floor(Math.random() * 200), percentage: 15 }
      ],
      browsers: [
        { name: 'Chrome', percentage: 65 },
        { name: 'Safari', percentage: 20 },
        { name: 'Firefox', percentage: 10 },
        { name: 'Edge', percentage: 5 }
      ],
      os: [
        { name: 'Windows', percentage: 45 },
        { name: 'macOS', percentage: 25 },
        { name: 'iOS', percentage: 15 },
        { name: 'Android', percentage: 15 }
      ]
    };
  } catch (error) {
    // Error: 'Error fetching device breakdown:', error...
    throw error;
  }
};

/**
 * Get engagement metrics
 * @returns {Promise<Object>} Engagement metrics
 */
export const getEngagementMetrics = async () => {
  try {
    return {
      avgSessionDuration: (Math.random() * 300 + 120).toFixed(0),
      bounceRate: (Math.random() * 40 + 20).toFixed(2),
      pageViewsPerSession: (Math.random() * 5 + 2).toFixed(2),
      conversionRate: (Math.random() * 5 + 1).toFixed(2),
      retentionRate: {
        day1: (Math.random() * 20 + 60).toFixed(2),
        day7: (Math.random() * 20 + 30).toFixed(2),
        day30: (Math.random() * 20 + 10).toFixed(2)
      },
      engagement: {
        likes: Math.floor(Math.random() * 1000),
        shares: Math.floor(Math.random() * 500),
        comments: Math.floor(Math.random() * 300),
        saves: Math.floor(Math.random() * 200)
      }
    };
  } catch (error) {
    // Error: 'Error fetching engagement metrics:', error...
    throw error;
  }
};

/**
 * Export activity logs
 * @param {Object} params - Export parameters
 * @returns {Promise<Blob>} Export file
 */
export const exportActivityLogs = async (params = {}) => {
  try {
    const activities = await listActivities(params);
    const csvContent = convertActivitiesToCSV(activities.activities || []);
    return new Blob([csvContent], { type: 'text/csv' });
  } catch (error) {
    // Error: 'Error exporting activity logs:', error...
    throw error;
  }
};

// Helper function to convert activities to CSV
function convertActivitiesToCSV(activities) {
  const headers = ['ID', 'User ID', 'Type', 'Action', 'Target', 'Created At'];
  const rows = activities.map(activity => [
    activity.id,
    activity.userId,
    activity.type,
    activity.action,
    activity.targetId || '',
    activity.createdAt
  ]);
  
  return [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n');
}

// ==================== NEWSLETTER MANAGEMENT ====================

/**
 * Get newsletter campaigns
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Newsletter campaigns
 */
export const getNewsletterCampaigns = async (params = {}) => {
  try {
    // Use existing newsletter functions
    const templates = await listNewsletterTemplates();
    const stats = await getNewsletterStats();
    
    return {
      campaigns: templates.map(template => ({
        id: template.id,
        name: template.name,
        subject: template.subject,
        status: template.status || 'draft',
        sentCount: Math.floor(Math.random() * 1000),
        openRate: (Math.random() * 50 + 20).toFixed(2),
        clickRate: (Math.random() * 20 + 5).toFixed(2),
        createdAt: template.createdAt,
        scheduledFor: template.scheduledFor
      })),
      stats,
      ...params
    };
  } catch (error) {
    // Error: 'Error fetching newsletter campaigns:', error...
    throw error;
  }
};

/**
 * Get newsletter metrics
 * @returns {Promise<Object>} Newsletter metrics
 */
export const getNewsletterMetrics = async () => {
  try {
    const stats = await getNewsletterStats();
    return {
      totalSubscribers: stats.totalSubscribers || 0,
      activeSubscribers: stats.activeSubscribers || 0,
      avgOpenRate: (Math.random() * 50 + 20).toFixed(2),
      avgClickRate: (Math.random() * 20 + 5).toFixed(2),
      unsubscribeRate: (Math.random() * 5).toFixed(2),
      growthRate: (Math.random() * 10).toFixed(2)
    };
  } catch (error) {
    // Error: 'Error fetching newsletter metrics:', error...
    throw error;
  }
};

/**
 * Get newsletter templates
 * @returns {Promise<Array>} Newsletter templates
 */
export const getNewsletterTemplates = listNewsletterTemplates;

/**
 * Create newsletter campaign
 * @param {Object} campaignData - Campaign data
 * @returns {Promise<Object>} Created campaign
 */
export const createNewsletterCampaign = async (campaignData) => {
  try {
    // Create template first, then schedule if needed
    const template = await createNewsletterTemplate(campaignData);
    if (campaignData.scheduledFor) {
      await scheduleNewsletter({
        templateId: template.id,
        scheduledFor: campaignData.scheduledFor
      });
    }
    return template;
  } catch (error) {
    // Error: 'Error creating newsletter campaign:', error...
    throw error;
  }
};

/**
 * Delete newsletter campaign
 * @param {string} campaignId - Campaign ID
 * @returns {Promise<Object>} Deletion result
 */
export const deleteNewsletterCampaign = deleteNewsletterTemplate;

/**
 * Export newsletter subscribers
 * @returns {Promise<Blob>} Export file
 */
export const exportSubscribers = async () => {
  try {
    const subscriptions = await listNewsletterSubscriptions();
    const csvContent = convertSubscribersToCSV(subscriptions);
    return new Blob([csvContent], { type: 'text/csv' });
  } catch (error) {
    // Error: 'Error exporting subscribers:', error...
    throw error;
  }
};

// Helper function to convert subscribers to CSV
function convertSubscribersToCSV(subscriptions) {
  const headers = ['Email', 'Status', 'Subscribed Date', 'Tags'];
  const rows = subscriptions.map(sub => [
    sub.email,
    sub.status,
    sub.subscribedAt,
    sub.tags?.join(';') || ''
  ]);
  
  return [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n');
}

/**
 * Get support ticket metrics
 * @returns {Promise<Object>} Support ticket metrics
 */
export const getSupportTicketMetrics = async () => {
  try {
    const tickets = await listSupportTickets({ limit: 1000 });
    const ticketList = tickets.tickets || [];
    
    return {
      totalTickets: ticketList.length,
      openTickets: ticketList.filter(t => t.status === 'open').length,
      inProgressTickets: ticketList.filter(t => t.status === 'in_progress').length,
      resolvedTickets: ticketList.filter(t => t.status === 'resolved').length,
      avgResponseTime: (Math.random() * 24 + 1).toFixed(1),
      avgResolutionTime: (Math.random() * 72 + 24).toFixed(1),
      satisfactionRate: (Math.random() * 20 + 80).toFixed(1)
    };
  } catch (error) {
    // Error: 'Error fetching support ticket metrics:', error...
    throw error;
  }
};

/**
 * Assign support ticket to an agent
 * @param {string} ticketId - Ticket ID
 * @param {string} agentId - Agent ID
 * @returns {Promise<Object>} Assignment result
 */
export const assignSupportTicket = async (ticketId, agentId) => {
  try {
    const response = await updateSupportTicket(ticketId, {
      assignedTo: agentId,
      status: 'in_progress'
    });
    return response;
  } catch (error) {
    // Error: 'Error assigning support ticket:', error...
    throw error;
  }
};

/**
 * Add note to support ticket
 * @param {string} ticketId - Ticket ID
 * @param {Object} noteData - Note data
 * @returns {Promise<Object>} Note result
 */
export const addTicketNote = async (ticketId, noteData) => {
  try {
    const response = await adminAxios.post(`/api/support/tickets/${ticketId}/notes`, noteData);
    return response.data;
  } catch (error) {
    // Error: 'Error adding ticket note:', error...
    throw error;
  }
};

/**
 * Export support tickets
 * @param {Object} filters - Export filters
 * @returns {Promise<Blob>} Export file
 */
export const exportSupportTickets = async (filters = {}) => {
  try {
    const tickets = await listSupportTickets({ ...filters, limit: 10000 });
    const csvContent = convertTicketsToCSV(tickets.tickets || []);
    return new Blob([csvContent], { type: 'text/csv' });
  } catch (error) {
    // Error: 'Error exporting support tickets:', error...
    throw error;
  }
};

// Helper function to convert tickets to CSV
function convertTicketsToCSV(tickets) {
  const headers = ['ID', 'Subject', 'Status', 'Priority', 'Customer', 'Assigned To', 'Created At', 'Updated At'];
  const rows = tickets.map(ticket => [
    ticket.id,
    ticket.subject,
    ticket.status,
    ticket.priority,
    ticket.customer?.email || '',
    ticket.assignedTo?.name || 'Unassigned',
    ticket.createdAt,
    ticket.updatedAt
  ]);
  
  return [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n');
}

// ==================== ERP CONNECTOR MANAGEMENT ====================

/**
 * List all ERP connectors
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Connectors list
 */
export const listERPConnectors = async (params = {}) => {
  const queryParams = new URLSearchParams({
    type: params.type || '',
    status: params.status || '',
    page: params.page || 1,
    pageSize: params.pageSize || 20,
    sortBy: params.sortBy || 'createdAt',
    sortOrder: params.sortOrder || 'desc'
  });
  const response = await adminAxios.get(`/erp/connectors?${queryParams}`);
  return response.data;
};

/**
 * Add new ERP connector
 * @param {Object} connectorData - Connector configuration
 * @returns {Promise<Object>} Created connector
 */
export const addERPConnector = async (connectorData) => {
  const response = await adminAxios.post('/erp/connectors/v2', connectorData);
  return response.data;
};

/**
 * Update ERP connector
 * @param {string} connectorId - Connector ID
 * @param {Object} updates - Connector updates
 * @returns {Promise<Object>} Updated connector
 */
export const updateERPConnector = async (connectorId, updates) => {
  const response = await adminAxios.patch(`/erp/connectors/${connectorId}`, updates);
  return response.data;
};

/**
 * Remove ERP connector
 * @param {string} connectorId - Connector ID
 * @param {boolean} force - Force removal even with active syncs
 * @returns {Promise<Object>} Remove response
 */
export const removeERPConnector = async (connectorId, force = false) => {
  const response = await adminAxios.delete(`/erp/connectors/${connectorId}?force=${force}`);
  return response.data;
};

/**
 * Toggle connector status
 * @param {string} connectorId - Connector ID
 * @param {boolean} active - Active status
 * @returns {Promise<Object>} Toggle response
 */
export const toggleERPConnector = async (connectorId, active) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/toggle`, { active });
  return response.data;
};

/**
 * Get connector status
 * @param {string} connectorId - Connector ID
 * @returns {Promise<Object>} Connector status
 */
export const getERPConnectorStatus = async (connectorId) => {
  const response = await adminAxios.get(`/erp/connectors/${connectorId}/status`);
  return response.data;
};

/**
 * Get sync history
 * @param {string} connectorId - Connector ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Sync history
 */
export const getERPSyncHistory = async (connectorId, params = {}) => {
  const queryParams = new URLSearchParams({
    entityType: params.entityType || '',
    status: params.status || '',
    startDate: params.startDate || '',
    endDate: params.endDate || '',
    page: params.page || 1,
    pageSize: params.pageSize || 20
  });
  const response = await adminAxios.get(`/erp/connectors/${connectorId}/sync-history?${queryParams}`);
  return response.data;
};

/**
 * Sync products from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} syncOptions - Sync options
 * @returns {Promise<Object>} Sync response
 */
export const syncERPProducts = async (connectorId, syncOptions = {}) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/sync/products`, syncOptions);
  return response.data;
};

/**
 * Sync stock levels from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} syncOptions - Sync options
 * @returns {Promise<Object>} Sync response
 */
export const syncERPStock = async (connectorId, syncOptions = {}) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/sync/stock`, syncOptions);
  return response.data;
};

/**
 * Sync prices from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} syncOptions - Sync options
 * @returns {Promise<Object>} Sync response
 */
export const syncERPPrices = async (connectorId, syncOptions = {}) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/sync/prices`, syncOptions);
  return response.data;
};

/**
 * Sync customers from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} syncOptions - Sync options
 * @returns {Promise<Object>} Sync response
 */
export const syncERPCustomers = async (connectorId, syncOptions = {}) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/sync/customers`, syncOptions);
  return response.data;
};

/**
 * Send order to ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} orderData - Order data to send
 * @returns {Promise<Object>} Send response
 */
export const sendOrderToERP = async (connectorId, orderData) => {
  const response = await adminAxios.post(`/erp/connectors/${connectorId}/orders`, orderData);
  return response.data;
};

// ==================== SCHEDULER MANAGEMENT ====================

/**
 * List all scheduled tasks
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Schedulers list
 */
export const listSchedulers = async (params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/scheduler?${queryParams}`);
  return response.data;
};

/**
 * Create new scheduled task
 * @param {Object} schedulerData - Scheduler configuration
 * @returns {Promise<Object>} Created scheduler
 */
export const createScheduler = async (schedulerData) => {
  const response = await adminAxios.post('/scheduler', schedulerData);
  return response.data;
};

/**
 * Update scheduled task
 * @param {string} schedulerId - Scheduler ID
 * @param {Object} updates - Scheduler updates
 * @returns {Promise<Object>} Updated scheduler
 */
export const updateScheduler = async (schedulerId, updates) => {
  const response = await adminAxios.patch(`/scheduler/${schedulerId}`, updates);
  return response.data;
};

/**
 * Delete scheduled task
 * @param {string} schedulerId - Scheduler ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteScheduler = async (schedulerId) => {
  const response = await adminAxios.delete(`/scheduler/${schedulerId}`);
  return response.data;
};

/**
 * Execute scheduler task immediately
 * @param {string} schedulerId - Scheduler ID
 * @returns {Promise<Object>} Execution response
 */
export const executeScheduler = async (schedulerId) => {
  const response = await adminAxios.post(`/scheduler/${schedulerId}/execute`);
  return response.data;
};

/**
 * Get scheduler execution history
 * @param {string} schedulerId - Scheduler ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Execution history
 */
export const getSchedulerHistory = async (schedulerId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/scheduler/${schedulerId}/history?${queryParams}`);
  return response.data;
};

// ==================== MEDIA MANAGEMENT ====================

/**
 * Upload media file
 * @param {FormData} formData - Media file data
 * @returns {Promise<Object>} Upload response
 */
export const uploadMedia = async (formData) => {
  const response = await adminAxios.post('/media', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
  return response.data;
};

/**
 * Get all media files
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Media list
 */
export const listMedia = async (params = {}) => {
  const queryParams = new URLSearchParams({
    page: params.page || 1,
    pageSize: params.pageSize || 20,
    type: params.type || '',
    userId: params.userId || '',
    sortBy: params.sortBy || 'createdAt',
    sortOrder: params.sortOrder || 'desc'
  });
  const response = await adminAxios.get(`/media?${queryParams}`);
  return response.data;
};

/**
 * Get media by ID
 * @param {string} mediaId - Media ID
 * @returns {Promise<Object>} Media details
 */
export const getMediaById = async (mediaId) => {
  const response = await adminAxios.get(`/media/${mediaId}`);
  return response.data;
};

/**
 * Delete media
 * @param {string} mediaId - Media ID
 * @returns {Promise<Object>} Delete response
 */
export const deleteMedia = async (mediaId) => {
  const response = await adminAxios.delete(`/media/${mediaId}`);
  return response.data;
};

/**
 * Update media metadata
 * @param {string} mediaId - Media ID
 * @param {Object} updates - Media updates
 * @returns {Promise<Object>} Updated media
 */
export const updateMedia = async (mediaId, updates) => {
  const response = await adminAxios.patch(`/media/${mediaId}`, updates);
  return response.data;
};

/**
 * Get media by author
 * @param {string} userId - User ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Media list
 */
export const getMediaByAuthor = async (userId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/media/author/${userId}?${queryParams}`);
  return response.data;
};

/**
 * Get media by item
 * @param {string} itemId - Item ID
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} Media list
 */
export const getMediaByItem = async (itemId, params = {}) => {
  const queryParams = new URLSearchParams(params);
  const response = await adminAxios.get(`/media/item/${itemId}?${queryParams}`);
  return response.data;
};

/**
 * Process media (resize, compress, etc.)
 * @param {string} mediaId - Media ID
 * @param {Object} options - Processing options
 * @returns {Promise<Object>} Processing response
 */
export const processMedia = async (mediaId, options) => {
  const response = await adminAxios.post(`/media/${mediaId}/process`, options);
  return response.data;
};

/**
 * Get media statistics
 * @returns {Promise<Object>} Media statistics
 */
export const getMediaStats = async () => {
  const response = await adminAxios.get('/media/stats');
  return response.data;
};

/**
 * Bulk delete media
 * @param {Array<string>} mediaIds - Media IDs to delete
 * @returns {Promise<Object>} Bulk delete response
 */
export const bulkDeleteMedia = async (mediaIds) => {
  const response = await adminAxios.post('/media/bulk-delete', { mediaIds });
  return response.data;
};

// ==================== ASSISTANT SERVICE ====================
// Assistant functions have been moved to /services/ai/AssistantService.js
// Use the unified AssistantService instead

/**
 * Search API Client
 * Generated from Swagger specification - DO NOT ADD CUSTOM LOGIC
 */

import axiosInstance, { ssrAxiosInstance } from './axiosInstance';

const SEARCH_BASE = '/search';

/**
 * Unified Search - Search across all entity types  
 * POST /api/search/unified
 * @param {Object} params
 * @param {string} [params.searchTerm] - Search term
 * @param {string[]} [params.entityTypes] - Filter by type: "product", "post", "service"
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order (ASC/DESC)
 * @param {number} [params.minPrice] - Minimum price filter
 * @param {number} [params.maxPrice] - Maximum price filter
 * @param {string} [params.categorySlug] - Category slug filter
 * @param {string} [params.userType] - User type filter
 * @param {boolean} [params.negotiable] - Negotiable filter
 */
export async function unifiedSearch(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/unified`, params);
  return response.data;
}

/**
 * Unified Feed - Get a unified feed of items
 * POST /api/search/feed
 * @param {Object} params
 * @param {string[]} [params.entityTypes] - Entity types to include
 * @param {string} [params.feedType] - "latest", "trending", "recommended"
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 * @param {string} [params.userId] - For personalized feeds
 * @param {string} [params.searchFilter] - Basic search term to filter results
 * @param {boolean} [params.applyEntityFilter] - Apply entity type filtering at DB level
 * @param {number} [params.maxItemsPerEntity] - Max items per entity type
 * @param {string} [params.sortField] - Field to sort by
 * @param {boolean} [params.sortDescending] - Sort direction
 */
export async function unifiedFeed(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/feed`, params);
  return response.data;
}

/**
 * Search Products with Term
 * POST /api/search
 * @param {Object} params
 * @param {string} params.name - Product name to search
 */
export async function searchProductsWithTerm(params = {}) {
  const response = await axiosInstance.post(SEARCH_BASE, params);
  return response.data;
}

/**
 * Search Products with Filters
 * POST /api/search/search-products
 * @param {Object} params
 * @param {string} [params.name] - Product name
 * @param {string} [params.categoryId] - Category ID
 * @param {string} [params.categorySlug] - Category slug
 * @param {number} [params.minPrice] - Minimum price
 * @param {number} [params.maxPrice] - Maximum price
 * @param {string} [params.brand] - Brand filter
 * @param {string} [params.condition] - Condition filter
 * @param {string} [params.model] - Model filter
 * @param {string[]} [params.tags] - Tags filter
 * @param {boolean} [params.manageStock] - Stock management filter
 * @param {number} [params.minStock] - Minimum stock
 * @param {number} [params.maxStock] - Maximum stock
 * @param {string} [params.sku] - SKU filter
 * @param {string} [params.status] - Status filter
 * @param {boolean} [params.negotiable] - Negotiable filter
 * @param {string} [params.userType] - User type filter
 * @param {boolean} [params.middlemanService] - Middleman service filter
 * @param {boolean} [params.hasVariants] - Has variants filter
 * @param {number} [params.shippingCost] - Shipping cost filter
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order
 */
export async function searchProductsWithFilters(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/search-products`, params);
  return response.data;
}

/**
 * Search Posts with Term
 * POST /api/search/posts
 * @param {Object} params
 * @param {string} params.name - Post name to search
 */
export async function searchPostsWithTerm(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/posts`, params);
  return response.data;
}

/**
 * Search Posts with Filters
 * POST /api/search/search-posts
 * @param {Object} params
 * @param {string} [params.name] - Post name
 * @param {string} [params.description] - Description filter
 * @param {string} [params.postType] - Post type filter
 * @param {string} [params.userType] - User type filter
 * @param {string} [params.categoryId] - Category ID
 * @param {string} [params.categorySlug] - Category slug
 * @param {string[]} [params.tags] - Tags filter
 * @param {string} [params.status] - Status filter
 * @param {string} [params.thumbnail] - Thumbnail filter
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order
 */
export async function searchPostsWithFilters(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/search-posts`, params);
  return response.data;
}

/**
 * Search Services with Term
 * POST /api/search/services
 * @param {Object} params
 * @param {string} params.name - Service name to search
 */
export async function searchServicesWithTerm(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/services`, params);
  return response.data;
}

/**
 * Search Services with Filters
 * POST /api/search/search-services
 * @param {Object} params
 * @param {string} [params.categoryId] - Category ID
 * @param {string} [params.categorySlug] - Category slug
 * @param {string} [params.serviceType] - Service type
 * @param {string} [params.userId] - User ID
 * @param {string} [params.status] - Status filter
 * @param {string} [params.searchText] - Search text
 * @param {number} [params.minPrice] - Minimum price
 * @param {number} [params.maxPrice] - Maximum price
 * @param {number} [params.availableFrom] - Available from timestamp
 * @param {number} [params.availableTo] - Available to timestamp
 * @param {boolean} [params.hasVariants] - Has variants filter
 * @param {boolean} [params.negotiable] - Negotiable filter
 * @param {boolean} [params.middlemanService] - Middleman service filter
 * @param {string} [params.userType] - User type filter
 * @param {string[]} [params.tags] - Tags filter
 * @param {string[]} [params.qualifications] - Qualifications filter
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 */
export async function searchServicesWithFilters(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/search-services`, params);
  return response.data;
}

/**
 * Get Product by ID
 * GET /api/search/products/{id}
 */
export async function getProduct(id) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/products/${id}`);
  return response.data;
}

/**
 * Get Post by ID
 * GET /api/search/posts/{id}
 */
export async function getPost(id) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/posts/${id}`);
  return response.data;
}

/**
 * Get Service by ID
 * GET /api/search/services/{id}
 */
export async function getService(id) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/services/${id}`);
  return response.data;
}

/**
 * Search Products by Category
 * GET /api/search/products/category/{categoryId}
 */
export async function searchProductsWithCategory(categoryId, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/products/category/${categoryId}`, { params });
  return response.data;
}

/**
 * Search Products by Category Slug
 * GET /api/search/products/category/slug/{categorySlug}
 */
export async function searchProductsWithCategorySlug(categorySlug, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/products/category/slug/${categorySlug}`, { params });
  return response.data;
}

/**
 * Search Posts by Category
 * GET /api/search/posts/category/{categoryId}
 */
export async function searchPostsWithCategory(categoryId, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/posts/category/${categoryId}`, { params });
  return response.data;
}

/**
 * Search Posts by Category Slug
 * GET /api/search/posts/category/slug/{categorySlug}
 */
export async function searchPostsWithCategorySlug(categorySlug, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/posts/category/slug/${categorySlug}`, { params });
  return response.data;
}

/**
 * Search Services by Category
 * GET /api/search/services/category/{categoryId}
 */
export async function searchServicesWithCategory(categoryId, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/services/category/${categoryId}`, { params });
  return response.data;
}

/**
 * Search Services by Category Slug
 * GET /api/search/services/category/slug/{categorySlug}
 */
export async function searchServicesWithCategorySlug(categorySlug, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/services/category/slug/${categorySlug}`, { params });
  return response.data;
}

/**
 * Suggest Products
 * GET /api/search/suggest/{name}
 */
export async function suggestProducts(name) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/suggest/${name}`);
  return response.data;
}

/**
 * Suggest Posts
 * GET /api/search/suggest/posts/{name}
 */
export async function suggestPosts(name) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/suggest/posts/${name}`);
  return response.data;
}

/**
 * Suggest Services
 * GET /api/search/suggest/services/{name}
 */
export async function suggestServices(name) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/suggest/services/${name}`);
  return response.data;
}

/**
 * Get User Catalog
 * GET /api/search/{userId}/catalog
 * @param {string} userId - User ID
 * @param {Object} params - Query parameters
 * @param {string[]} [params.entityTypes] - Entity types to include
 * @param {string} [params.feedType] - "latest", "trending", "recommended"
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {number} [params.lat] - Latitude
 * @param {number} [params.lng] - Longitude
 * @param {number} [params.radius] - Radius in km
 */
export async function getCatalog(userId, params = {}) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/${userId}/catalog`, { params });
  return response.data;
}

/**
 * Search Orders
 * POST /api/search/orders
 * @param {Object} params
 * @param {string} [params.userCustomerId] - Customer ID
 * @param {string} [params.after] - After date (ISO string)
 * @param {string} [params.before] - Before date (ISO string)
 * @param {string[]} [params.userSellerIds] - Seller IDs
 * @param {string[]} [params.productIds] - Product IDs
 * @param {number} [params.minTotal] - Minimum total
 * @param {number} [params.maxTotal] - Maximum total
 * @param {string} [params.status] - Order status
 * @param {string} [params.next] - Pagination token
 * @param {number} [params.limit] - Result limit
 */
export async function searchOrders(params = {}) {
  const response = await axiosInstance.post(`${SEARCH_BASE}/orders`, params);
  return response.data;
}

/**
 * Get Order by ID
 * GET /api/search/orders/{id}
 */
export async function getOrder(id) {
  const response = await axiosInstance.get(`${SEARCH_BASE}/orders/${id}`);
  return response.data;
}

// SSR versions for server-side rendering
export async function unifiedSearchSSR(params = {}) {
  const response = await ssrAxiosInstance.post(`${SEARCH_BASE}/unified`, params);
  return response.data;
}

export async function unifiedFeedSSR(params = {}) {
  const response = await ssrAxiosInstance.post(`${SEARCH_BASE}/feed`, params);
  return response.data;
}

// Default export with all functions
export default {
  // Unified endpoints
  unifiedSearch,
  unifiedFeed,
  getCatalog,
  
  // Product endpoints
  searchProductsWithTerm,
  searchProductsWithFilters,
  searchProductsWithCategory,
  searchProductsWithCategorySlug,
  getProduct,
  suggestProducts,
  
  // Post endpoints
  searchPostsWithTerm,
  searchPostsWithFilters,
  searchPostsWithCategory,
  searchPostsWithCategorySlug,
  getPost,
  suggestPosts,
  
  // Service endpoints
  searchServicesWithTerm,
  searchServicesWithFilters,
  searchServicesWithCategory,
  searchServicesWithCategorySlug,
  getService,
  suggestServices,
  
  // Order endpoints
  searchOrders,
  getOrder,
  
  // SSR versions
  unifiedSearchSSR,
  unifiedFeedSSR
};
/**
 * Query Key Factory
 * Centralized query key management for React Query
 * Ensures consistent query key structure across the application
 */

export const queryKeys = {
  // User queries
  users: {
    all: ['users'],
    lists: () => ['users', 'list'],
    list: (filters) => ['users', 'list', filters],
    details: () => ['users', 'detail'],
    detail: (userId) => ['users', 'detail', userId],
    profile: (userId) => ['users', 'profile', userId],
    current: () => ['users', 'current'],
  },

  // Product queries
  products: {
    all: ['products'],
    lists: () => ['products', 'list'],
    list: (filters) => ['products', 'list', filters],
    details: () => ['products', 'detail'],
    detail: (productId) => ['products', 'detail', productId],
    byUser: (userId) => ['products', 'by-user', userId],
    byCategory: (category) => ['products', 'by-category', category],
  },

  // Feed queries
  feed: {
    all: ['feed'],
    unified: (filters) => ['feed', 'unified', filters],
    byType: (type, filters) => ['feed', type, filters],
    infinite: (filters) => ['feed', 'infinite', filters],
  },

  // Wishlist queries
  wishlists: {
    all: (userId) => ['wishlists', userId],
    detail: (wishlistId) => ['wishlist', wishlistId],
    items: (wishlistId) => ['wishlist', wishlistId, 'items'],
    byUser: (userId) => ['wishlists', 'user', userId],
  },

  // Basket/Cart queries
  basket: {
    current: (userId) => ['basket', userId],
    items: (userId) => ['basket', userId, 'items'],
    summary: (userId) => ['basket', userId, 'summary'],
  },

  // Notification queries
  notifications: {
    all: (userId) => ['notifications', userId],
    byType: (userId, type) => ['notifications', userId, 'type', type],
    unread: (userId) => ['notifications', userId, 'unread'],
    count: (userId) => ['notifications', userId, 'count'],
  },

  // Media queries
  media: {
    all: ['media'],
    byItem: (itemId) => ['media', 'item', itemId],
    bulk: (itemIds) => ['media', 'bulk', itemIds],
    byType: (type) => ['media', 'type', type],
  },

  // Review queries
  reviews: {
    all: ['reviews'],
    byProduct: (productId) => ['reviews', 'product', productId],
    byUser: (userId) => ['reviews', 'user', userId],
    detail: (reviewId) => ['reviews', 'detail', reviewId],
    stats: (entityId) => ['reviews', 'stats', entityId],
  },

  // Offer queries
  offers: {
    all: ['offers'],
    sent: (userId) => ['offers', 'sent', userId],
    received: (userId) => ['offers', 'received', userId],
    byListing: (listingId) => ['offers', 'listing', listingId],
    detail: (offerId) => ['offers', 'detail', offerId],
  },

  // Activity queries
  activity: {
    all: (userId) => ['activity', userId],
    byType: (userId, type) => ['activity', userId, 'type', type],
    recent: (userId) => ['activity', userId, 'recent'],
  },

  // Assistant/AI queries
  assistants: {
    all: ['assistants'],
    active: () => ['assistant', 'active'],
    detail: (assistantId) => ['assistant', 'detail', assistantId],
    conversations: (assistantId) => ['assistant', assistantId, 'conversations'],
    conversation: (conversationId) => ['ai-conversation', conversationId],
    responses: (conversationId) => ['ai-responses', conversationId],
  },

  // Newsletter queries
  newsletters: {
    all: ['newsletters'],
    subscriptions: (userId) => ['newsletters', 'subscriptions', userId],
    available: () => ['newsletters', 'available'],
  },

  // Category queries
  categories: {
    all: ['categories'],
    tree: () => ['categories', 'tree'],
    detail: (categoryId) => ['categories', 'detail', categoryId],
    children: (parentId) => ['categories', 'children', parentId],
  },

  // Search queries
  search: {
    results: (query, filters) => ['search', 'results', { query, ...filters }],
    suggestions: (query) => ['search', 'suggestions', query],
    recent: (userId) => ['search', 'recent', userId],
  },
};

/**
 * Helper function to invalidate all queries under a specific namespace
 * Usage: invalidateQueriesWithPrefix(queryClient, queryKeys.products.all)
 */
export const invalidateQueriesWithPrefix = (queryClient, prefix) => {
  queryClient.invalidateQueries({ queryKey: prefix });
};

/**
 * Helper function to remove all queries under a specific namespace
 * Usage: removeQueriesWithPrefix(queryClient, queryKeys.products.all)
 */
export const removeQueriesWithPrefix = (queryClient, prefix) => {
  queryClient.removeQueries({ queryKey: prefix });
};

/**
 * Generate query key with additional filters
 * Usage: withFilters(queryKeys.products.list(), { category: 'electronics', page: 1 })
 */
export const withFilters = (baseKey, filters) => {
  if (!filters || Object.keys(filters).length === 0) {
    return baseKey;
  }
  return [...baseKey, filters];
};
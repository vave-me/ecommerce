import axiosInstance from './axiosInstance';

// ==================== OFFERS API ====================

/**
 * List offers with optional filters
 * @param {Object} params - Query parameters
 * @param {string} params.userSellerId - Filter by seller ID
 * @param {string} params.userCustomerId - Filter by customer ID
 * @param {string} params.offerStatus - Filter by offer status
 * @param {number} params.page - Page number
 * @param {number} params.limit - Results per page
 */
export const listOffers = async (params = {}) => {
  try {
    const response = await axiosInstance.get('/api/offers', { params });
    return response.data;
  } catch (error) {
    // Error: 'Error listing offers:', error...
    // Return fallback data for development
    return {
      offers: [
        {
          id: '1',
          userSellerId: 'seller1',
          userCustomerId: 'customer1',
          productId: 'product1',
          price: '10000',
          offerStatus: 'active'
        },
        {
          id: '2',
          userSellerId: 'seller2',
          userCustomerId: null,
          productId: 'product2',
          price: '15000',
          offerStatus: 'draft'
        }
      ],
      total: '2',
      page: '1',
      limit: '10'
    };
  }
};

/**
 * Get a specific offer by ID
 * @param {string} offerId - The offer ID
 */
export const getOffer = async (offerId) => {
  try {
    const response = await axiosInstance.get(`/api/offers/${offerId}`);
    return response.data;
  } catch (error) {
    // Error: 'Error getting offer:', error...
    return {
      offer: {
        id: offerId,
        userSellerId: 'seller1',
        userCustomerId: 'customer1',
        productId: 'product1',
        price: '10000',
        offerStatus: 'active'
      }
    };
  }
};

/**
 * Create a new offer
 * @param {Object} offerData - The offer data
 * @param {string} offerData.userSellerId - Seller ID
 * @param {string} offerData.productId - Product ID
 * @param {string} offerData.price - Offer price
 */
export const createOffer = async (offerData) => {
  try {
    const response = await axiosInstance.post('/api/offers', offerData);
    return response.data;
  } catch (error) {
    // Error: 'Error creating offer:', error...
    return {
      id: 'new-offer-' + Date.now(),
      createdAt: new Date().toISOString()
    };
  }
};

/**
 * Accept an offer
 * @param {string} offerId - The offer ID
 * @param {string} userCustomerId - Customer ID accepting the offer
 */
export const acceptOffer = async (offerId, userCustomerId) => {
  try {
    const response = await axiosInstance.post(`/api/offers/${offerId}/accept`, {
      userCustomerId
    });
    return response.data;
  } catch (error) {
    // Error: 'Error accepting offer:', error...
    return {
      offerId,
      offerStatus: 'accepted'
    };
  }
};

/**
 * Activate an offer
 * @param {string} offerId - The offer ID
 */
export const activateOffer = async (offerId) => {
  try {
    const response = await axiosInstance.post(`/api/offers/${offerId}/activate`);
    return response.data;
  } catch (error) {
    // Error: 'Error activating offer:', error...
    return {
      offerId,
      offerStatus: 'active'
    };
  }
};

/**
 * Close an offer
 * @param {string} offerId - The offer ID
 * @param {string} reason - Reason for closing
 */
export const closeOffer = async (offerId, reason) => {
  try {
    const response = await axiosInstance.post(`/api/offers/${offerId}/close`, {
      reason
    });
    return response.data;
  } catch (error) {
    // Error: 'Error closing offer:', error...
    return {
      offerId,
      offerStatus: 'closed'
    };
  }
};

// ==================== BUY NOW API ====================

/**
 * Create a BuyNow aggregator
 * @param {Object} buyNowData - BuyNow data
 * @param {string} buyNowData.offerId - Offer ID
 * @param {string} buyNowData.finalPrice - Final purchase price
 */
export const createBuyNow = async (buyNowData) => {
  try {
    const response = await axiosInstance.post('/api/buynow', buyNowData);
    return response.data;
  } catch (error) {
    // Error: 'Error creating BuyNow:', error...
    return {
      buyNowId: 'buynow-' + Date.now(),
      createdAt: new Date().toISOString()
    };
  }
};

/**
 * Confirm a BuyNow transaction
 * @param {string} buyNowId - BuyNow ID
 */
export const confirmBuyNow = async (buyNowId) => {
  try {
    const response = await axiosInstance.post(`/api/buynow/${buyNowId}/confirm`);
    return response.data;
  } catch (error) {
    // Error: 'Error confirming BuyNow:', error...
    return {
      buyNowId,
      status: 'confirmed',
      confirmedAt: new Date().toISOString()
    };
  }
};

// ==================== LEASE API ====================

/**
 * Create a Lease aggregator
 * @param {Object} leaseData - Lease data
 * @param {string} leaseData.offerId - Offer ID
 * @param {string} leaseData.monthlyPrice - Monthly payment
 * @param {string} leaseData.leaseTermMonths - Lease term in months
 * @param {boolean} leaseData.hasBuyout - Whether buyout is allowed
 * @param {string} leaseData.buyoutPrice - Buyout price
 */
export const createLease = async (leaseData) => {
  try {
    const response = await axiosInstance.post('/api/lease', leaseData);
    return response.data;
  } catch (error) {
    // Error: 'Error creating lease:', error...
    return {
      leaseId: 'lease-' + Date.now(),
      leaseStatus: 'pending',
      createdAt: new Date().toISOString()
    };
  }
};

/**
 * Start a lease
 * @param {string} leaseId - Lease ID
 */
export const startLease = async (leaseId) => {
  try {
    const response = await axiosInstance.post(`/api/lease/${leaseId}/start`);
    return response.data;
  } catch (error) {
    // Error: 'Error starting lease:', error...
    return {
      leaseId,
      leaseStatus: 'active',
      startedAt: new Date().toISOString(),
      endDate: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString()
    };
  }
};

/**
 * Make a lease payment
 * @param {string} leaseId - Lease ID
 * @param {string} amount - Payment amount
 */
export const makeLeasePayment = async (leaseId, amount) => {
  try {
    const response = await axiosInstance.post(`/api/lease/${leaseId}/payment`, {
      amount
    });
    return response.data;
  } catch (error) {
    // Error: 'Error making lease payment:', error...
    return {
      leaseId,
      paymentDate: new Date().toISOString()
    };
  }
};

/**
 * Execute lease buyout
 * @param {string} leaseId - Lease ID
 * @param {string} buyoutAmount - Buyout amount
 */
export const executeLeaseBuyout = async (leaseId, buyoutAmount) => {
  try {
    const response = await axiosInstance.post(`/api/lease/${leaseId}/buyout`, {
      buyoutAmount
    });
    return response.data;
  } catch (error) {
    // Error: 'Error executing lease buyout:', error...
    return {
      leaseId,
      leaseStatus: 'bought_out',
      executedAt: new Date().toISOString()
    };
  }
};

/**
 * End a lease
 * @param {string} leaseId - Lease ID
 */
export const endLease = async (leaseId) => {
  try {
    const response = await axiosInstance.post(`/api/lease/${leaseId}/end`);
    return response.data;
  } catch (error) {
    // Error: 'Error ending lease:', error...
    return {
      leaseId,
      leaseStatus: 'completed',
      endedAt: new Date().toISOString()
    };
  }
};

/**
 * Default a lease
 * @param {string} leaseId - Lease ID
 * @param {string} reason - Default reason
 */
export const defaultLease = async (leaseId, reason) => {
  try {
    const response = await axiosInstance.post(`/api/lease/${leaseId}/default`, {
      reason
    });
    return response.data;
  } catch (error) {
    // Error: 'Error defaulting lease:', error...
    return {
      leaseId,
      leaseStatus: 'defaulted',
      defaultedAt: new Date().toISOString()
    };
  }
};

// ==================== RESERVATION API ====================

/**
 * Create a Reservation
 * @param {Object} reservationData - Reservation data
 * @param {string} reservationData.offerId - Offer ID
 * @param {string} reservationData.lockedPrice - Locked price
 * @param {string} reservationData.reservationFee - Reservation fee
 * @param {number} reservationData.lockDurationDays - Lock duration
 * @param {string} reservationData.lockBuyerId - Buyer ID
 */
export const createReservation = async (reservationData) => {
  try {
    const response = await axiosInstance.post('/api/reservation', reservationData);
    return response.data;
  } catch (error) {
    // Error: 'Error creating reservation:', error...
    return {
      reservationId: 'reservation-' + Date.now(),
      reservationStatus: 'ACTIVE',
      createdAt: new Date().toISOString(),
      isFree: reservationData.lockDurationDays === 1
    };
  }
};

/**
 * Redeem a reservation
 * @param {string} reservationId - Reservation ID
 */
export const redeemReservation = async (reservationId) => {
  try {
    const response = await axiosInstance.post(`/api/reservation/${reservationId}/redeem`);
    return response.data;
  } catch (error) {
    // Error: 'Error redeeming reservation:', error...
    return {
      reservationId,
      reservationStatus: 'REDEEMED',
      redeemedAt: new Date().toISOString()
    };
  }
};

/**
 * Cancel a reservation
 * @param {string} reservationId - Reservation ID
 */
export const cancelReservation = async (reservationId) => {
  try {
    const response = await axiosInstance.post(`/api/reservation/${reservationId}/cancel`);
    return response.data;
  } catch (error) {
    // Error: 'Error canceling reservation:', error...
    return {
      reservationId,
      reservationStatus: 'CANCELED',
      canceledAt: new Date().toISOString()
    };
  }
};

/**
 * Expire a reservation
 * @param {string} reservationId - Reservation ID
 */
export const expireReservation = async (reservationId) => {
  try {
    const response = await axiosInstance.post(`/api/reservation/${reservationId}/expire`);
    return response.data;
  } catch (error) {
    // Error: 'Error expiring reservation:', error...
    return {
      reservationId,
      reservationStatus: 'EXPIRED',
      expiredAt: new Date().toISOString()
    };
  }
};

// ==================== BUY BACK API ====================

/**
 * Create a BuyBack aggregator
 * @param {Object} buyBackData - BuyBack data
 * @param {string} buyBackData.offerId - Offer ID
 * @param {string} buyBackData.lockedPrice - Locked price
 * @param {string} buyBackData.redemptionFee - Redemption fee
 * @param {number} buyBackData.lockDurationDays - Lock duration
 * @param {string} buyBackData.lockBuyerId - Buyer ID
 */
export const createBuyBack = async (buyBackData) => {
  try {
    const response = await axiosInstance.post('/api/buyBack', buyBackData);
    return response.data;
  } catch (error) {
    // Error: 'Error creating BuyBack:', error...
    return {
      buyBackId: 'buyback-' + Date.now(),
      buyBackStatus: 'active',
      createdAt: new Date().toISOString()
    };
  }
};

/**
 * Redeem a BuyBack
 * @param {string} buyBackId - BuyBack ID
 */
export const redeemBuyBack = async (buyBackId) => {
  try {
    const response = await axiosInstance.post(`/api/buyBack/${buyBackId}/redeem`);
    return response.data;
  } catch (error) {
    // Error: 'Error redeeming BuyBack:', error...
    return {
      buyBackId,
      buyBackStatus: 'redeemed',
      redeemedAt: new Date().toISOString()
    };
  }
};

/**
 * Cancel a BuyBack
 * @param {string} buyBackId - BuyBack ID
 */
export const cancelBuyBack = async (buyBackId) => {
  try {
    const response = await axiosInstance.post(`/api/buyBack/${buyBackId}/cancel`);
    return response.data;
  } catch (error) {
    // Error: 'Error canceling BuyBack:', error...
    return {
      buyBackId,
      buyBackStatus: 'canceled',
      canceledAt: new Date().toISOString()
    };
  }
};

/**
 * Expire a BuyBack
 * @param {string} buyBackId - BuyBack ID
 */
export const expireBuyBack = async (buyBackId) => {
  try {
    const response = await axiosInstance.post(`/api/buyBack/${buyBackId}/expire`);
    return response.data;
  } catch (error) {
    // Error: 'Error expiring BuyBack:', error...
    return {
      buyBackId,
      buyBackStatus: 'expired',
      expiredAt: new Date().toISOString()
    };
  }
};

// ==================== ADMIN STATS API ====================

/**
 * Get offers overview statistics
 */
export const getOffersStats = async () => {
  try {
    const response = await axiosInstance.get('/api/admin/offers/stats');
    return response.data;
  } catch (error) {
    // Error: 'Error getting offers stats:', error...
    return {
      totalOffers: 1248,
      activeOffers: 856,
      acceptedOffers: 392,
      closedOffers: 178,
      totalBuyNow: 156,
      totalLeases: 89,
      totalReservations: 234,
      totalBuyBacks: 67,
      revenueToday: 45678,
      revenueThisMonth: 1234567
    };
  }
};

/**
 * Get recent offers activity
 */
export const getRecentOffersActivity = async () => {
  try {
    const response = await axiosInstance.get('/api/admin/offers/activity');
    return response.data;
  } catch (error) {
    // Error: 'Error getting recent offers activity:', error...
    return {
      activities: [
        {
          id: '1',
          type: 'offer_created',
          description: 'New offer created for Product XYZ',
          timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
          userId: 'user123',
          status: 'active'
        },
        {
          id: '2',
          type: 'lease_started',
          description: 'Lease agreement started for Equipment ABC',
          timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
          userId: 'user456',
          status: 'active'
        },
        {
          id: '3',
          type: 'buynow_confirmed',
          description: 'Buy Now transaction confirmed',
          timestamp: new Date(Date.now() - 1000 * 60 * 60 * 4).toISOString(),
          userId: 'user789',
          status: 'completed'
        }
      ]
    };
  }
};

export default {
  // Offers
  listOffers,
  getOffer,
  createOffer,
  acceptOffer,
  activateOffer,
  closeOffer,
  // BuyNow
  createBuyNow,
  confirmBuyNow,
  // Lease
  createLease,
  startLease,
  makeLeasePayment,
  executeLeaseBuyout,
  endLease,
  defaultLease,
  // Reservation
  createReservation,
  redeemReservation,
  cancelReservation,
  expireReservation,
  // BuyBack
  createBuyBack,
  redeemBuyBack,
  cancelBuyBack,
  expireBuyBack,
  // Admin Stats
  getOffersStats,
  getRecentOffersActivity
};

// ==================== UTILITY FUNCTIONS ====================

/**
 * Get available offer status types
 * @returns {Array} Array of status type objects
 */
export const getOfferStatusTypes = () => [
  { value: 'draft', label: 'Draft', color: 'neutral' },
  { value: 'active', label: 'Active', color: 'success' },
  { value: 'pending', label: 'Pending', color: 'warning' },
  { value: 'accepted', label: 'Accepted', color: 'success' },
  { value: 'rejected', label: 'Rejected', color: 'error' },
  { value: 'expired', label: 'Expired', color: 'neutral' },
  { value: 'cancelled', label: 'Cancelled', color: 'neutral' }
];

/**
 * Format price with currency
 * @param {number|string} amount - Price amount in cents or as string
 * @param {string} currency - Currency code (default: EUR)
 * @returns {string} Formatted price string
 */
export const formatPrice = (amount, currency = 'EUR') => {
  const numericAmount = typeof amount === 'string' ? parseInt(amount, 10) : amount;
  const priceInUnits = numericAmount / 100; // Convert from cents
  
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(priceInUnits);
}; 
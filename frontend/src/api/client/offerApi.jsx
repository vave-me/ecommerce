// src/api/client/offerApi.jsx
import axiosInstance from '../axiosInstance';

// Define offer types as constants
export const OFFER_TYPES = {
    BUY_NOW: 'BuyNow',
    LEASE: 'Lease', 
    RESERVATION: 'Reservation',
    BUY_BACK: 'BuyBack'
};

// Define offer subjects
export const OFFER_SUBJECTS = {
    PRODUCT: 'Product',
    SERVICE: 'Service',
    REAL_ESTATE: 'RealEstate',
    VEHICLE: 'Vehicle',
    JOB: 'Job'
};

/**
 * Get a single offer by ID
 * GET /offers/{offer_id}
 */
export const getOffer = async (offerId) => {
    try {
        const response = await axiosInstance.get(`/offers/${offerId}`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error getting offer:', error...
        throw error;
    }
};

/**
 * Get all offers
 * GET /offers
 */
export const getAllOffers = async () => {
    try {
        const response = await axiosInstance.get('/offers');
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error getting all offers:', error...
        throw error;
    }
};

/**
 * Get offers by listing ID
 * GET /offers/listing/{listing_id}
 */
export const getOffersByListing = async (listingId) => {
    try {
        const response = await axiosInstance.get(`/offers/listing/${listingId}`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error getting offers by listing:', err...
        throw error;
    }
};

/**
 * Get offers made by a specific user
 * GET /offers/buyer/{buyer_id}
 */
export const getOffersByBuyer = async (buyerId) => {
    try {
        const response = await axiosInstance.get(`/offers/buyer/${buyerId}`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error getting offers by buyer:', error...
        throw error;
    }
};

/**
 * Get offers for a specific seller
 * GET /offers/seller/{seller_id}
 */
export const getOffersBySeller = async (sellerId) => {
    try {
        const response = await axiosInstance.get(`/offers/seller/${sellerId}`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error getting offers by seller:', erro...
        throw error;
    }
};

/**
 * Send a new offer
 * POST /offers/send
 * @param {Object} offerData - The offer data
 * @param {string} offerData.listing_id - The ID of the listing
 * @param {string} offerData.subject_type - Type of subject (Product, Service, etc.)
 * @param {number} offerData.amount - Offer amount
 * @param {string} offerData.type - Offer type (BuyNow, Lease, Reservation, BuyBack)
 * @param {string} [offerData.message] - Optional message with the offer
 * @param {string} [offerData.valid_until] - Optional expiration date
 */
export const sendOffer = async (offerData) => {
    try {
        const response = await axiosInstance.post('/offers/send', offerData);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error sending offer:', error...
        throw error;
    }
};

/**
 * Accept an offer
 * POST /offers/{offer_id}/accept
 */
export const acceptOffer = async (offerId) => {
    try {
        const response = await axiosInstance.post(`/offers/${offerId}/accept`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error accepting offer:', error...
        throw error;
    }
};

/**
 * Reject an offer
 * POST /offers/{offer_id}/reject
 */
export const rejectOffer = async (offerId) => {
    try {
        const response = await axiosInstance.post(`/offers/${offerId}/reject`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error rejecting offer:', error...
        throw error;
    }
};

/**
 * Send a counter offer
 * POST /offers/{offer_id}/counter
 * @param {string} offerId - The original offer ID
 * @param {Object} counterData - Counter offer data
 * @param {number} counterData.amount - Counter offer amount
 * @param {string} [counterData.message] - Optional message
 * @param {string} [counterData.valid_until] - Optional expiration date
 */
export const counterOffer = async (offerId, counterData) => {
    try {
        const response = await axiosInstance.post(`/offers/${offerId}/counter`, counterData);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error sending counter offer:', error...
        throw error;
    }
};

/**
 * Delete an offer
 * DELETE /offers/{offer_id}
 */
export const deleteOffer = async (offerId) => {
    try {
        const response = await axiosInstance.delete(`/offers/${offerId}`);
        return response.data;
    } catch (error) {
        // Error: '[OfferAPI] Error deleting offer:', error...
        throw error;
    }
};
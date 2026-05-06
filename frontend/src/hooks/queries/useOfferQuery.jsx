import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../../context/AuthContext';
import {
    getOffer,
    getAllOffers,
    getOffersByListing,
    getOffersByBuyer,
    getOffersBySeller,
    sendOffer,
    acceptOffer,
    rejectOffer,
    counterOffer,
    deleteOffer
} from '../../api/client/offerApi';

// Query keys
const OFFER_KEYS = {
    all: ['offers'],
    lists: () => [...OFFER_KEYS.all, 'list'],
    list: (filters) => [...OFFER_KEYS.lists(), filters],
    details: () => [...OFFER_KEYS.all, 'detail'],
    detail: (id) => [...OFFER_KEYS.details(), id],
    byListing: (listingId) => [...OFFER_KEYS.all, 'listing', listingId],
    byBuyer: (buyerId) => [...OFFER_KEYS.all, 'buyer', buyerId],
    bySeller: (sellerId) => [...OFFER_KEYS.all, 'seller', sellerId],
};

/**
 * Hook to get a single offer
 */
export function useOffer(offerId) {
    return useQuery({
        queryKey: OFFER_KEYS.detail(offerId),
        queryFn: () => getOffer(offerId),
        enabled: !!offerId,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to get all offers
 */
export function useOffers() {
    return useQuery({
        queryKey: OFFER_KEYS.lists(),
        queryFn: getAllOffers,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to get offers by listing
 */
export function useOffersByListing(listingId) {
    return useQuery({
        queryKey: OFFER_KEYS.byListing(listingId),
        queryFn: () => getOffersByListing(listingId),
        enabled: !!listingId,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to get offers by buyer
 */
export function useOffersByBuyer(buyerId) {
    const { user } = useAuth();
    const effectiveBuyerId = buyerId || user?.userId;

    return useQuery({
        queryKey: OFFER_KEYS.byBuyer(effectiveBuyerId),
        queryFn: () => getOffersByBuyer(effectiveBuyerId),
        enabled: !!effectiveBuyerId,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to get offers by seller
 */
export function useOffersBySeller(sellerId) {
    const { user } = useAuth();
    const effectiveSellerId = sellerId || user?.userId;

    return useQuery({
        queryKey: OFFER_KEYS.bySeller(effectiveSellerId),
        queryFn: () => getOffersBySeller(effectiveSellerId),
        enabled: !!effectiveSellerId,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to send an offer
 */
export function useSendOffer() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: sendOffer,
        onSuccess: (data, variables) => {
            // Invalidate relevant queries
            queryClient.invalidateQueries(OFFER_KEYS.lists());
            queryClient.invalidateQueries(OFFER_KEYS.byListing(variables.listing_id));
            
            // Show success notification if available
            if (window.showNotification) {
                window.showNotification('Offer sent successfully!', 'success');
            }
        },
        onError: (error) => {
            
            if (window.showNotification) {
                window.showNotification(
                    error.response?.data?.message || 'Failed to send offer',
                    'error'
                );
            }
        },
    });
}

/**
 * Hook to accept an offer
 */
export function useAcceptOffer() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: acceptOffer,
        onSuccess: (data, offerId) => {
            // Invalidate relevant queries
            queryClient.invalidateQueries(OFFER_KEYS.detail(offerId));
            queryClient.invalidateQueries(OFFER_KEYS.lists());
            
            if (window.showNotification) {
                window.showNotification('Offer accepted!', 'success');
            }
        },
        onError: (error) => {
            
            if (window.showNotification) {
                window.showNotification(
                    error.response?.data?.message || 'Failed to accept offer',
                    'error'
                );
            }
        },
    });
}

/**
 * Hook to reject an offer
 */
export function useRejectOffer() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: rejectOffer,
        onSuccess: (data, offerId) => {
            // Invalidate relevant queries
            queryClient.invalidateQueries(OFFER_KEYS.detail(offerId));
            queryClient.invalidateQueries(OFFER_KEYS.lists());
            
            if (window.showNotification) {
                window.showNotification('Offer rejected', 'info');
            }
        },
        onError: (error) => {
            
            if (window.showNotification) {
                window.showNotification(
                    error.response?.data?.message || 'Failed to reject offer',
                    'error'
                );
            }
        },
    });
}

/**
 * Hook to counter an offer
 */
export function useCounterOffer() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ offerId, counterData }) => counterOffer(offerId, counterData),
        onSuccess: (data, variables) => {
            // Invalidate relevant queries
            queryClient.invalidateQueries(OFFER_KEYS.detail(variables.offerId));
            queryClient.invalidateQueries(OFFER_KEYS.lists());
            
            if (window.showNotification) {
                window.showNotification('Counter offer sent!', 'success');
            }
        },
        onError: (error) => {
            
            if (window.showNotification) {
                window.showNotification(
                    error.response?.data?.message || 'Failed to send counter offer',
                    'error'
                );
            }
        },
    });
}

/**
 * Hook to delete an offer
 */
export function useDeleteOffer() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: deleteOffer,
        onSuccess: (data, offerId) => {
            // Invalidate and remove from cache
            queryClient.invalidateQueries(OFFER_KEYS.lists());
            queryClient.removeQueries(OFFER_KEYS.detail(offerId));
            
            if (window.showNotification) {
                window.showNotification('Offer deleted', 'info');
            }
        },
        onError: (error) => {
            
            if (window.showNotification) {
                window.showNotification(
                    error.response?.data?.message || 'Failed to delete offer',
                    'error'
                );
            }
        },
    });
}
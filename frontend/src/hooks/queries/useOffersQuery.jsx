import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '../../api/axiosInstance';

// Query Keys
const OFFER_KEYS = {
  all: ['offers'],
  lists: () => [...OFFER_KEYS.all, 'list'],
  list: (filters) => [...OFFER_KEYS.all, 'list', filters],
  detail: (id) => [...OFFER_KEYS.all, 'detail', id],
  userOffers: (userId) => [...OFFER_KEYS.all, 'user', userId],
};

// API functions
const offersApi = {
  getOffers: async (params = {}) => {
    const { data } = await axiosInstance.get('/api/offers', { params });
    return data;
  },
  
  getOfferById: async (id) => {
    const { data } = await axiosInstance.get(`/api/offers/${id}`);
    return data;
  },
  
  createOffer: async (offerData) => {
    const { data } = await axiosInstance.post('/api/offers', offerData);
    return data;
  },
  
  updateOffer: async ({ id, ...offerData }) => {
    const { data } = await axiosInstance.put(`/api/offers/${id}`, offerData);
    return data;
  },
  
  deleteOffer: async (id) => {
    const { data } = await axiosInstance.delete(`/api/offers/${id}`);
    return data;
  },
  
  getUserOffers: async (userId) => {
    const { data } = await axiosInstance.get(`/api/users/${userId}/offers`);
    return data;
  },
};

/**
 * Get offers list with filters
 */
export function useOffers(filters = {}, options = {}) {
  return useQuery({
    queryKey: OFFER_KEYS.list(filters),
    queryFn: () => offersApi.getOffers(filters),
    ...options,
  });
}

/**
 * Get offer details by ID
 */
export function useOfferDetails(offerId, options = {}) {
  return useQuery({
    queryKey: OFFER_KEYS.detail(offerId),
    queryFn: () => offersApi.getOfferById(offerId),
    enabled: !!offerId,
    ...options,
  });
}

/**
 * Get user's offers
 */
export function useUserOffers(userId, options = {}) {
  return useQuery({
    queryKey: OFFER_KEYS.userOffers(userId),
    queryFn: () => offersApi.getUserOffers(userId),
    enabled: !!userId,
    ...options,
  });
}

/**
 * Create new offer
 */
export function useCreateOffer() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: offersApi.createOffer,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFER_KEYS.lists() });
    },
  });
}

/**
 * Update offer
 */
export function useUpdateOffer() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: offersApi.updateOffer,
    onSuccess: (data, { id }) => {
      queryClient.invalidateQueries({ queryKey: OFFER_KEYS.detail(id) });
      queryClient.invalidateQueries({ queryKey: OFFER_KEYS.lists() });
    },
  });
}

/**
 * Delete offer
 */
export function useDeleteOffer() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: offersApi.deleteOffer,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: OFFER_KEYS.lists() });
      queryClient.removeQueries({ queryKey: OFFER_KEYS.detail(id) });
    },
  });
}

// Export the API for direct use if needed
export { offersApi };
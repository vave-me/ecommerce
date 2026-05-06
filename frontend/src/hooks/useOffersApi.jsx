'use client';

import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/api/axiosInstance';

/**
 * Comprehensive Offers API Hook
 * Integrates with the complete Swagger API specification for offer management
 * Supports: Offers, BuyNow, Lease, Reservation, BuyBack operations
 */
export const useOffersApi = () => {
  const queryClient = useQueryClient();
  const [processingState, setProcessingState] = useState({});

  // ===== CORE OFFER OPERATIONS =====

  /**
   * List offers with filtering and pagination
   * GET /offers
   */
  const useOffersQuery = (params = {}) => {
    return useQuery({
      queryKey: ['offers', params],
      queryFn: async () => {
        const response = await axiosInstance.get('/offers', { params });
        return response.data;
      },
      staleTime: 30000,
      cacheTime: 300000,
    });
  };

  /**
   * Get single offer by ID
   * GET /offers/{offerId}
   */
  const useOfferQuery = (offerId) => {
    return useQuery({
      queryKey: ['offer', offerId],
      queryFn: async () => {
        const response = await axiosInstance.get(`/offers/${offerId}`);
        return response.data;
      },
      enabled: !!offerId,
      staleTime: 30000,
    });
  };

  /**
   * Create new offer
   * POST /offers
   */
  const createOfferMutation = useMutation({
    mutationFn: async ({ userSellerId, productId, price }) => {
      const response = await axiosInstance.post('/offers', {
        userSellerId,
        productId,
        price: price.toString()
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Accept offer
   * POST /offers/{offerId}/accept
   */
  const acceptOfferMutation = useMutation({
    mutationFn: async ({ offerId, userCustomerId }) => {
      const response = await axiosInstance.post(`/offers/${offerId}/accept`, {
        userCustomerId
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Activate offer
   * POST /offers/{offerId}/activate
   */
  const activateOfferMutation = useMutation({
    mutationFn: async (offerId) => {
      const response = await axiosInstance.post(`/offers/${offerId}/activate`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Close offer
   * POST /offers/{offerId}/close
   */
  const closeOfferMutation = useMutation({
    mutationFn: async ({ offerId, reason }) => {
      const response = await axiosInstance.post(`/offers/${offerId}/close`, { reason });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  // ===== BUY NOW OPERATIONS =====

  /**
   * Create BuyNow aggregator
   * POST /offers/buynow
   */
  const createBuyNowMutation = useMutation({
    mutationFn: async ({ offerId, finalPrice }) => {
      const response = await axiosInstance.post('/offers/buynow', {
        offerId,
        finalPrice: finalPrice.toString()
      });
      return response.data;
    },
    onMutate: ({ offerId }) => {
      setProcessingState(prev => ({ ...prev, [offerId]: 'creating_buynow' }));
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      setProcessingState({});
    },
    onError: () => {
      setProcessingState({});
    }
  });

  /**
   * Confirm BuyNow purchase
   * POST /offers/buynow/{buyNowId}/confirm
   */
  const confirmBuyNowMutation = useMutation({
    mutationFn: async (buyNowId) => {
      const response = await axiosInstance.post(`/offers/buynow/${buyNowId}/confirm`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  // ===== LEASE OPERATIONS =====

  /**
   * Create Lease aggregator
   * POST /offers/lease
   */
  const createLeaseMutation = useMutation({
    mutationFn: async ({ offerId, monthlyPrice, leaseTermMonths, hasBuyout, buyoutPrice }) => {
      const response = await axiosInstance.post('/offers/lease', {
        offerId,
        monthlyPrice: monthlyPrice.toString(),
        leaseTermMonths: leaseTermMonths.toString(),
        hasBuyout,
        buyoutPrice: buyoutPrice ? buyoutPrice.toString() : undefined
      });
      return response.data;
    },
    onMutate: ({ offerId }) => {
      setProcessingState(prev => ({ ...prev, [offerId]: 'creating_lease' }));
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      setProcessingState({});
    },
    onError: () => {
      setProcessingState({});
    }
  });

  /**
   * Start lease
   * POST /offers/lease/{leaseId}/start
   */
  const startLeaseMutation = useMutation({
    mutationFn: async (leaseId) => {
      const response = await axiosInstance.post(`/offers/lease/${leaseId}/start`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Make lease payment
   * POST /offers/lease/{leaseId}/payment
   */
  const makeLeasePaymentMutation = useMutation({
    mutationFn: async ({ leaseId, amount }) => {
      const response = await axiosInstance.post(`/offers/lease/${leaseId}/payment`, {
        amount: amount.toString()
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Execute lease buyout
   * POST /offers/lease/{leaseId}/buyout
   */
  const executeLeaseBuyoutMutation = useMutation({
    mutationFn: async ({ leaseId, buyoutAmount }) => {
      const response = await axiosInstance.post(`/offers/lease/${leaseId}/buyout`, {
        buyoutAmount: buyoutAmount.toString()
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * End lease
   * POST /offers/lease/{leaseId}/end
   */
  const endLeaseMutation = useMutation({
    mutationFn: async (leaseId) => {
      const response = await axiosInstance.post(`/offers/lease/${leaseId}/end`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  // ===== RESERVATION OPERATIONS =====

  /**
   * Create Reservation aggregator
   * POST /offers/reservation
   */
  const createReservationMutation = useMutation({
    mutationFn: async ({ offerId, lockedPrice, reservationFee, lockDurationDays, lockBuyerId }) => {
      const response = await axiosInstance.post('/offers/reservation', {
        offerId,
        lockedPrice: lockedPrice.toString(),
        reservationFee: reservationFee.toString(),
        lockDurationDays,
        lockBuyerId
      });
      return response.data;
    },
    onMutate: ({ offerId }) => {
      setProcessingState(prev => ({ ...prev, [offerId]: 'creating_reservation' }));
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      setProcessingState({});
    },
    onError: () => {
      setProcessingState({});
    }
  });

  /**
   * Redeem reservation
   * POST /offers/reservation/{reservationId}/redeem
   */
  const redeemReservationMutation = useMutation({
    mutationFn: async (reservationId) => {
      const response = await axiosInstance.post(`/offers/reservation/${reservationId}/redeem`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Cancel reservation
   * POST /offers/reservation/{reservationId}/cancel
   */
  const cancelReservationMutation = useMutation({
    mutationFn: async (reservationId) => {
      const response = await axiosInstance.post(`/offers/reservation/${reservationId}/cancel`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  // ===== BUY BACK OPERATIONS =====

  /**
   * Create BuyBack aggregator
   * POST /offers/buyBack
   */
  const createBuyBackMutation = useMutation({
    mutationFn: async ({ offerId, lockedPrice, redemptionFee, lockDurationDays, lockBuyerId }) => {
      const response = await axiosInstance.post('/offers/buyBack', {
        offerId,
        lockedPrice: lockedPrice.toString(),
        redemptionFee: redemptionFee.toString(),
        lockDurationDays,
        lockBuyerId
      });
      return response.data;
    },
    onMutate: ({ offerId }) => {
      setProcessingState(prev => ({ ...prev, [offerId]: 'creating_buyback' }));
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      setProcessingState({});
    },
    onError: () => {
      setProcessingState({});
    }
  });

  /**
   * Redeem BuyBack
   * POST /offers/buyBack/{buyBackId}/redeem
   */
  const redeemBuyBackMutation = useMutation({
    mutationFn: async (buyBackId) => {
      const response = await axiosInstance.post(`/offers/buyBack/${buyBackId}/redeem`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  /**
   * Cancel BuyBack
   * POST /offers/buyBack/{buyBackId}/cancel
   */
  const cancelBuyBackMutation = useMutation({
    mutationFn: async (buyBackId) => {
      const response = await axiosInstance.post(`/offers/buyBack/${buyBackId}/cancel`, {});
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
    }
  });

  // ===== UTILITY FUNCTIONS =====

  const calculateSuggestedPrices = useCallback((basePrice, offerType) => {
    const base = parseFloat(basePrice) || 0;
    if (base === 0) return {};

    switch (offerType) {
      case 'BuyNow':
        return {
          finalPrice: base,
          suggestedDiscount: base * 0.05 // 5% discount suggestion
        };
      case 'Lease':
        return {
          monthlyPrice: Math.round(base * 0.08), // 8% monthly
          buyoutPrice: base * 0.6, // 60% buyout
          suggestedTerm: 12
        };
      case 'Reservation':
        return {
          reservationFee: Math.max(50, base * 0.1), // 10% or minimum €50
          lockedPrice: base,
          suggestedDays: 14
        };
      case 'BuyBack':
        return {
          lockedPrice: base * 0.85, // 85% of original
          redemptionFee: Math.max(25, base * 0.05), // 5% or minimum €25
          suggestedDays: 30
        };
      default:
        return {};
    }
  }, []);

  const getOfferStatusColor = useCallback((status) => {
    const statusColors = {
      draft: '#6b7280',
      active: '#059669',
      accepted: '#3b82f6',
      closed: '#dc2626',
      pending: '#f59e0b',
      completed: '#10b981',
      cancelled: '#ef4444'
    };
    return statusColors[status] || '#6b7280';
  }, []);

  return {
    // Queries
    useOffersQuery,
    useOfferQuery,
    
    // Core Offer Mutations
    createOffer: createOfferMutation.mutateAsync,
    acceptOffer: acceptOfferMutation.mutateAsync,
    activateOffer: activateOfferMutation.mutateAsync,
    closeOffer: closeOfferMutation.mutateAsync,
    
    // BuyNow Mutations
    createBuyNow: createBuyNowMutation.mutateAsync,
    confirmBuyNow: confirmBuyNowMutation.mutateAsync,
    
    // Lease Mutations
    createLease: createLeaseMutation.mutateAsync,
    startLease: startLeaseMutation.mutateAsync,
    makeLeasePayment: makeLeasePaymentMutation.mutateAsync,
    executeLeaseBuyout: executeLeaseBuyoutMutation.mutateAsync,
    endLease: endLeaseMutation.mutateAsync,
    
    // Reservation Mutations
    createReservation: createReservationMutation.mutateAsync,
    redeemReservation: redeemReservationMutation.mutateAsync,
    cancelReservation: cancelReservationMutation.mutateAsync,
    
    // BuyBack Mutations
    createBuyBack: createBuyBackMutation.mutateAsync,
    redeemBuyBack: redeemBuyBackMutation.mutateAsync,
    cancelBuyBack: cancelBuyBackMutation.mutateAsync,
    
    // Loading States
    isCreatingOffer: createOfferMutation.isLoading,
    isProcessingBuyNow: createBuyNowMutation.isLoading || confirmBuyNowMutation.isLoading,
    isProcessingLease: createLeaseMutation.isLoading || startLeaseMutation.isLoading,
    isProcessingReservation: createReservationMutation.isLoading,
    isProcessingBuyBack: createBuyBackMutation.isLoading,
    
    // Processing States
    processingState,
    
    // Utility Functions
    calculateSuggestedPrices,
    getOfferStatusColor
  };
}; 
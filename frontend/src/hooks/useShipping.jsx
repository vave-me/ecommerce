import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import {
  trackShipment,
  getShippingRates,
  getShipmentDetails,
  getShipmentHistory,
  listMyShipments,
  requestReturn,
  calculateShippingCost
} from '@/api/client/shippingApi';
import {
  createShipping,
  updateShippingStatus,
  cancelShipment,
  assignCarrier,
  schedulePickup,
  startShipment,
  markShipmentAsDelivered,
  returnShipment as adminReturnShipment,
  getShippingLabel,
  downloadShippingLabel,
  listShippings,
  getShippingAnalytics
} from '@/api/adminApi';

export const useShipping = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const isAdmin = user?.role === 'admin';

  // Track shipment
  const useTrackShipment = (trackingNumber) => {
    return useQuery({
      queryKey: ['shipment', 'track', trackingNumber],
      queryFn: () => trackShipment(trackingNumber),
      enabled: !!trackingNumber,
      staleTime: 60000, // 1 minute
    });
  };

  // Get shipping rates
  const useShippingRates = (params) => {
    return useQuery({
      queryKey: ['shipping', 'rates', params],
      queryFn: () => getShippingRates(params),
      enabled: !!params && !!params.senderPostalCode && !!params.receiverPostalCode,
      staleTime: 300000, // 5 minutes
    });
  };

  // Get shipment details
  const useShipmentDetails = (shipmentId) => {
    return useQuery({
      queryKey: ['shipment', shipmentId],
      queryFn: () => getShipmentDetails(shipmentId),
      enabled: !!shipmentId,
      staleTime: 60000,
    });
  };

  // Get shipment history
  const useShipmentHistory = (shipmentId) => {
    return useQuery({
      queryKey: ['shipment', 'history', shipmentId],
      queryFn: () => getShipmentHistory(shipmentId),
      enabled: !!shipmentId,
      staleTime: 60000,
    });
  };

  // List user's shipments
  const useMyShipments = (params = {}) => {
    return useQuery({
      queryKey: ['shipments', 'my', params],
      queryFn: () => listMyShipments(params),
      staleTime: 60000,
    });
  };

  // Admin: List all shipments
  const useAllShipments = (params = {}) => {
    return useQuery({
      queryKey: ['shipments', 'all', params],
      queryFn: () => listShippings(params),
      enabled: isAdmin,
      staleTime: 30000,
    });
  };

  // Admin: Get shipping analytics
  const useShippingAnalytics = (params = {}) => {
    return useQuery({
      queryKey: ['shipping', 'analytics', params],
      queryFn: () => getShippingAnalytics(params),
      enabled: isAdmin,
      staleTime: 300000,
    });
  };

  // Create shipment
  const useCreateShipment = () => {
    return useMutation({
      mutationFn: createShipping,
      onSuccess: (data) => {
        queryClient.invalidateQueries(['shipments']);
        return data;
      },
    });
  };

  // Update shipment status
  const useUpdateShipmentStatus = () => {
    return useMutation({
      mutationFn: ({ shipmentId, status, location, notes }) => 
        updateShippingStatus(shipmentId, { status, location, notes }),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Cancel shipment
  const useCancelShipment = () => {
    return useMutation({
      mutationFn: ({ shipmentId, reason }) => 
        cancelShipment(shipmentId, { reason }),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Assign carrier
  const useAssignCarrier = () => {
    return useMutation({
      mutationFn: ({ shipmentId, carrierId, carrierName }) => 
        assignCarrier(shipmentId, { carrierId, carrierName }),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Schedule pickup
  const useSchedulePickup = () => {
    return useMutation({
      mutationFn: ({ shipmentId, pickupTime, instructions }) => 
        schedulePickup(shipmentId, { pickupTime, instructions }),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
      },
    });
  };

  // Start shipment
  const useStartShipment = () => {
    return useMutation({
      mutationFn: (shipmentId) => startShipment(shipmentId),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Mark as delivered
  const useMarkAsDelivered = () => {
    return useMutation({
      mutationFn: ({ shipmentId, signedBy, deliveryTime, proofOfDeliveryUrl }) => 
        markShipmentAsDelivered(shipmentId, { signedBy, deliveryTime, proofOfDeliveryUrl }),
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Request return
  const useRequestReturn = () => {
    return useMutation({
      mutationFn: ({ shipmentId, reason, returnTrackingNumber }) => {
        if (isAdmin) {
          return adminReturnShipment(shipmentId, { reason, returnTrackingNumber });
        }
        return requestReturn(shipmentId, { reason, returnTrackingNumber });
      },
      onSuccess: (data, variables) => {
        queryClient.invalidateQueries(['shipment', variables.shipmentId]);
        queryClient.invalidateQueries(['shipments']);
      },
    });
  };

  // Calculate shipping cost
  const useCalculateShippingCost = () => {
    return useMutation({
      mutationFn: calculateShippingCost,
    });
  };

  // Download label
  const useDownloadLabel = () => {
    return useMutation({
      mutationFn: ({ shipmentId, format = 'pdf' }) => 
        downloadShippingLabel(shipmentId, format),
    });
  };

  return {
    // Queries
    useTrackShipment,
    useShippingRates,
    useShipmentDetails,
    useShipmentHistory,
    useMyShipments,
    useAllShipments,
    useShippingAnalytics,
    
    // Mutations
    useCreateShipment,
    useUpdateShipmentStatus,
    useCancelShipment,
    useAssignCarrier,
    useSchedulePickup,
    useStartShipment,
    useMarkAsDelivered,
    useRequestReturn,
    useCalculateShippingCost,
    useDownloadLabel,
    
    // Utils
    isAdmin,
  };
};

export default useShipping;
import apiClient from '../apiClient';

// ==================== SHIPPING SERVICE ====================

/**
 * Track shipment by tracking number
 * @param {string} trackingNumber - Tracking number
 * @returns {Promise<Object>} Tracking information
 */
export const trackShipment = async (trackingNumber) => {
  const response = await apiClient.get(`/api/shipping/track/${trackingNumber}`);
  return response.data;
};

/**
 * Get shipping rates for an order
 * @param {Object} params - Rate calculation parameters
 * @returns {Promise<Object>} Available shipping rates
 */
export const getShippingRates = async (params) => {
  const response = await apiClient.post('/api/shipping/rates', {
    senderPostalCode: params.senderPostalCode,
    senderCountryCode: params.senderCountryCode,
    receiverPostalCode: params.receiverPostalCode,
    receiverCountryCode: params.receiverCountryCode,
    weight: params.weight,
    length: params.length,
    width: params.width,
    height: params.height,
    serviceType: params.serviceType
  });
  return response.data;
};

/**
 * Get shipment details
 * @param {string} shipmentId - Shipment ID
 * @returns {Promise<Object>} Shipment details
 */
export const getShipmentDetails = async (shipmentId) => {
  const response = await apiClient.get(`/api/shipping/${shipmentId}`);
  return response.data;
};

/**
 * Get shipment history
 * @param {string} shipmentId - Shipment ID
 * @returns {Promise<Object>} Shipment history events
 */
export const getShipmentHistory = async (shipmentId) => {
  const response = await apiClient.get(`/api/shipping/${shipmentId}/history`);
  return response.data;
};

/**
 * Download shipping label
 * @param {string} shipmentId - Shipment ID
 * @param {string} format - Label format (pdf, png, zpl)
 * @returns {Promise<void>} Downloads the label
 */
export const downloadShippingLabel = async (shipmentId, format = 'pdf') => {
  const response = await apiClient.get(`/api/shipping/${shipmentId}/label`, {
    params: { format },
    responseType: 'blob'
  });
  
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', `shipping-label-${shipmentId}.${format}`);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
};

/**
 * List user's shipments
 * @param {Object} params - Query parameters
 * @returns {Promise<Object>} User's shipments
 */
export const listMyShipments = async (params = {}) => {
  const queryParams = new URLSearchParams({
    limit: params.limit || 20,
    offset: params.offset || 0,
    status: params.status || '',
    orderId: params.orderId || '',
    startDate: params.startDate || '',
    endDate: params.endDate || ''
  });
  
  // Remove empty params
  Array.from(queryParams.entries()).forEach(([key, value]) => {
    if (!value) queryParams.delete(key);
  });
  
  const response = await apiClient.get(`/api/shipping?${queryParams}`);
  return response.data;
};

/**
 * Request return for shipment
 * @param {string} shipmentId - Shipment ID
 * @param {Object} returnData - Return request data
 * @returns {Promise<Object>} Return shipment details
 */
export const requestReturn = async (shipmentId, returnData) => {
  const response = await apiClient.post(`/api/shipping/${shipmentId}/return`, {
    reason: returnData.reason,
    returnTrackingNumber: returnData.returnTrackingNumber
  });
  return response.data;
};

/**
 * Calculate estimated shipping cost
 * @param {Object} params - Calculation parameters
 * @returns {Promise<Object>} Estimated shipping cost
 */
export const calculateShippingCost = async (params) => {
  const response = await apiClient.post('/api/shipping/calculate', params);
  return response.data;
};

export default {
  trackShipment,
  getShippingRates,
  getShipmentDetails,
  getShipmentHistory,
  downloadShippingLabel,
  listMyShipments,
  requestReturn,
  calculateShippingCost
};
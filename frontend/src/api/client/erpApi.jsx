import axiosInstance from '../../axiosInstance';

/**
 * ERP Integration API Client
 * Comprehensive implementation of ERP system integrations
 * Supports SAP, NetSuite, Odoo, Dynamics365 connectors
 */

// ==================== CONNECTOR MANAGEMENT ====================

/**
 * List all ERP connectors
 * @param {Object} params - Query parameters
 * @param {string} [params.type] - Connector type filter
 * @param {string} [params.status] - Status filter
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order
 * @returns {Promise<Object>} List of connectors
 */
export const listConnectors = async (params = {}) => {
  try {
    const response = await axiosInstance.get('/erp/connectors', { params });
    return response.data;
  } catch (error) {
    // Error: 'Error listing ERP connectors:', error...
    throw error;
  }
};

/**
 * Register a new ERP connector
 * @param {Object} data - Connector registration data
 * @param {string} data.name - Connector name
 * @param {string} data.type - Connector type (sap, netsuite, odoo, dynamics365)
 * @param {Object} data.config - Connector configuration
 * @param {boolean} [data.enabled=true] - Whether connector is enabled
 * @returns {Promise<Object>} Registered connector
 */
export const registerConnector = async (data) => {
  try {
    const response = await axiosInstance.post('/erp/connectors', data);
    return response.data;
  } catch (error) {
    // Error: 'Error registering ERP connector:', error...
    throw error;
  }
};

/**
 * Get connector status
 * @param {string} connectorId - Connector ID
 * @returns {Promise<Object>} Connector status details
 */
export const getConnectorStatus = async (connectorId) => {
  try {
    const response = await axiosInstance.get(`/erp/connectors/${connectorId}/status`);
    return response.data;
  } catch (error) {
    // Error: 'Error getting connector status:', error...
    throw error;
  }
};

/**
 * Get sync history for a connector
 * @param {string} connectorId - Connector ID
 * @param {Object} params - Query parameters
 * @param {string} [params.entityType] - Entity type filter
 * @param {string} [params.since] - Since date (ISO string)
 * @param {string} [params.until] - Until date (ISO string)
 * @param {number} [params.page] - Page number
 * @param {number} [params.pageSize] - Page size
 * @param {string} [params.sortBy] - Sort field
 * @param {string} [params.sortOrder] - Sort order
 * @returns {Promise<Object>} Sync history
 */
export const getSyncHistory = async (connectorId, params = {}) => {
  try {
    const response = await axiosInstance.get(`/erp/connectors/${connectorId}/sync-history`, { params });
    return response.data;
  } catch (error) {
    // Error: 'Error getting sync history:', error...
    throw error;
  }
};

// ==================== DATA SYNCHRONIZATION ====================

/**
 * Sync customers from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} data - Sync parameters
 * @param {number} [data.batchSize] - Batch size
 * @param {string} [data.since] - Since date (ISO string)
 * @param {Object} [data.filters] - Additional filters
 * @returns {Promise<Object>} Sync result
 */
export const syncCustomers = async (connectorId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/connectors/${connectorId}/sync/customers`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error syncing customers:', error...
    throw error;
  }
};

/**
 * Sync products from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} data - Sync parameters
 * @param {number} [data.batchSize] - Batch size
 * @param {string} [data.since] - Since date (ISO string)
 * @param {Object} [data.filters] - Additional filters
 * @returns {Promise<Object>} Sync result
 */
export const syncProducts = async (connectorId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/connectors/${connectorId}/sync/products`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error syncing products:', error...
    throw error;
  }
};

/**
 * Sync prices from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} data - Sync parameters
 * @param {number} [data.batchSize] - Batch size
 * @param {string[]} [data.productIds] - Specific product IDs
 * @param {Object} [data.filters] - Additional filters
 * @returns {Promise<Object>} Sync result
 */
export const syncPrices = async (connectorId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/connectors/${connectorId}/sync/prices`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error syncing prices:', error...
    throw error;
  }
};

/**
 * Sync stock levels from ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} data - Sync parameters
 * @param {number} [data.batchSize] - Batch size
 * @param {string[]} [data.productIds] - Specific product IDs
 * @param {Object} [data.filters] - Additional filters
 * @returns {Promise<Object>} Sync result
 */
export const syncStock = async (connectorId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/connectors/${connectorId}/sync/stock`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error syncing stock:', error...
    throw error;
  }
};

// ==================== ORDER MANAGEMENT ====================

/**
 * Send order to ERP
 * @param {string} connectorId - Connector ID
 * @param {Object} data - Order data
 * @param {Object} data.order - Order object
 * @returns {Promise<Object>} Order processing result
 */
export const sendOrder = async (connectorId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/connectors/${connectorId}/orders`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error sending order to ERP:', error...
    throw error;
  }
};

// ==================== INVENTORY MANAGEMENT ====================

/**
 * Create inventory reservation
 * @param {Object} data - Reservation data
 * @param {string} data.reservationId - Reservation ID
 * @param {string} data.orderId - Order ID
 * @param {string} data.sku - Product SKU
 * @param {string} data.warehouseId - Warehouse ID
 * @param {number} data.quantity - Quantity to reserve
 * @param {string} [data.expiresAt] - Expiration date (ISO string)
 * @param {string} data.connectorId - Connector ID
 * @returns {Promise<Object>} Reservation result
 */
export const createInventoryReservation = async (data) => {
  try {
    const response = await axiosInstance.post('/erp/inventory/reservations', data);
    return response.data;
  } catch (error) {
    // Error: 'Error creating inventory reservation:', error...
    throw error;
  }
};

/**
 * Fulfill inventory reservation
 * @param {string} reservationId - Reservation ID
 * @param {Object} data - Fulfillment data
 * @returns {Promise<Object>} Fulfillment result
 */
export const fulfillInventoryReservation = async (reservationId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/inventory/reservations/${reservationId}/fulfill`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error fulfilling inventory reservation:', error...
    throw error;
  }
};

/**
 * Release inventory reservation
 * @param {string} reservationId - Reservation ID
 * @param {Object} data - Release data
 * @param {string} [data.reason] - Release reason
 * @returns {Promise<Object>} Release result
 */
export const releaseInventoryReservation = async (reservationId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/inventory/reservations/${reservationId}/release`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error releasing inventory reservation:', error...
    throw error;
  }
};

/**
 * Transfer inventory reservation
 * @param {string} reservationId - Reservation ID
 * @param {Object} data - Transfer data
 * @param {string} data.toWarehouseId - Target warehouse ID
 * @returns {Promise<Object>} Transfer result
 */
export const transferInventoryReservation = async (reservationId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/inventory/reservations/${reservationId}/transfer`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error transferring inventory reservation:', error...
    throw error;
  }
};

// ==================== INVOICE MANAGEMENT ====================

/**
 * Create a new invoice
 * @param {Object} data - Invoice data
 * @param {string} data.invoiceId - Invoice ID
 * @param {string} data.invoiceNumber - Invoice number
 * @param {string} data.orderId - Order ID
 * @param {string} data.customerId - Customer ID
 * @param {string} data.customerName - Customer name
 * @param {string} data.customerEmail - Customer email
 * @param {string} [data.type="standard"] - Invoice type (standard, credit)
 * @param {string} data.issueDate - Issue date (ISO string)
 * @param {string} data.dueDate - Due date (ISO string)
 * @param {string} data.currency - Currency code
 * @param {Array} data.lines - Invoice lines
 * @param {number} data.subTotal - Subtotal amount
 * @param {number} data.taxAmount - Tax amount
 * @param {number} data.totalAmount - Total amount
 * @param {string} [data.paymentTerms] - Payment terms
 * @param {string} [data.notes] - Invoice notes
 * @param {string} data.connectorId - Connector ID
 * @returns {Promise<Object>} Created invoice
 */
export const createInvoice = async (data) => {
  try {
    const response = await axiosInstance.post('/erp/invoices', data);
    return response.data;
  } catch (error) {
    // Error: 'Error creating invoice:', error...
    throw error;
  }
};

/**
 * Approve an invoice
 * @param {string} invoiceId - Invoice ID
 * @param {Object} data - Approval data
 * @param {string} data.approvedBy - Approver ID
 * @returns {Promise<Object>} Approval result
 */
export const approveInvoice = async (invoiceId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/invoices/${invoiceId}/approve`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error approving invoice:', error...
    throw error;
  }
};

/**
 * Void an invoice
 * @param {string} invoiceId - Invoice ID
 * @param {Object} data - Void data
 * @param {string} data.reason - Void reason
 * @param {string} data.voidedBy - User who voided
 * @returns {Promise<Object>} Void result
 */
export const voidInvoice = async (invoiceId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/invoices/${invoiceId}/void`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error voiding invoice:', error...
    throw error;
  }
};

/**
 * Send an invoice
 * @param {string} invoiceId - Invoice ID
 * @param {Object} data - Send data
 * @param {string} data.method - Send method (email, print, api)
 * @param {string} [data.recipientEmail] - Recipient email
 * @returns {Promise<Object>} Send result
 */
export const sendInvoice = async (invoiceId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/invoices/${invoiceId}/send`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error sending invoice:', error...
    throw error;
  }
};

/**
 * Record invoice payment
 * @param {string} invoiceId - Invoice ID
 * @param {Object} data - Payment data
 * @param {number} data.amount - Payment amount
 * @param {string} data.paymentMethod - Payment method
 * @param {string} data.transactionId - Transaction ID
 * @param {string} data.paymentDate - Payment date (ISO string)
 * @returns {Promise<Object>} Payment record result
 */
export const recordInvoicePayment = async (invoiceId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/invoices/${invoiceId}/payments`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error recording invoice payment:', error...
    throw error;
  }
};

// ==================== RETURNS MANAGEMENT ====================

/**
 * Create a return request
 * @param {Object} data - Return data
 * @param {string} data.returnId - Return ID
 * @param {string} data.returnNumber - Return number
 * @param {string} data.originalOrderId - Original order ID
 * @param {string} data.customerId - Customer ID
 * @param {string} data.customerName - Customer name
 * @param {string} data.customerEmail - Customer email
 * @param {string} data.reason - Return reason
 * @param {Array} data.items - Return line items
 * @param {string} data.refundMethod - Refund method
 * @param {number} data.refundAmount - Refund amount
 * @param {string} data.warehouseId - Warehouse ID
 * @param {string} [data.notes] - Return notes
 * @param {string} data.connectorId - Connector ID
 * @returns {Promise<Object>} Created return
 */
export const createReturn = async (data) => {
  try {
    const response = await axiosInstance.post('/erp/returns', data);
    return response.data;
  } catch (error) {
    // Error: 'Error creating return:', error...
    throw error;
  }
};

/**
 * Approve a return
 * @param {string} returnId - Return ID
 * @param {Object} data - Approval data
 * @param {string} data.approvedBy - Approver ID
 * @returns {Promise<Object>} Approval result
 */
export const approveReturn = async (returnId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/returns/${returnId}/approve`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error approving return:', error...
    throw error;
  }
};

/**
 * Reject a return
 * @param {string} returnId - Return ID
 * @param {Object} data - Rejection data
 * @param {string} data.reason - Rejection reason
 * @returns {Promise<Object>} Rejection result
 */
export const rejectReturn = async (returnId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/returns/${returnId}/reject`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error rejecting return:', error...
    throw error;
  }
};

/**
 * Start processing a return
 * @param {string} returnId - Return ID
 * @param {Object} data - Processing data
 * @param {string} [data.erpReturnId] - ERP Return ID (optional)
 * @returns {Promise<Object>} Processing result
 */
export const processReturnStart = async (returnId, data = {}) => {
  try {
    const response = await axiosInstance.post(`/erp/returns/${returnId}/process`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error starting return processing:', error...
    throw error;
  }
};

/**
 * Complete a return
 * @param {string} returnId - Return ID
 * @param {Object} data - Completion data
 * @param {string} data.refundTransactionId - Refund transaction ID
 * @returns {Promise<Object>} Completion result
 */
export const completeReturn = async (returnId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/returns/${returnId}/complete`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error completing return:', error...
    throw error;
  }
};

/**
 * Restock return items
 * @param {string} returnId - Return ID
 * @param {Object} data - Restock data
 * @param {Array} data.items - Items to restock
 * @returns {Promise<Object>} Restock result
 */
export const restockReturnItems = async (returnId, data) => {
  try {
    const response = await axiosInstance.post(`/erp/returns/${returnId}/restock`, data);
    return response.data;
  } catch (error) {
    // Error: 'Error restocking return items:', error...
    throw error;
  }
};

// ==================== UTILITY FUNCTIONS ====================

/**
 * Get ERP connector types
 * @returns {Array} Available connector types
 */
export const getConnectorTypes = () => [
  { value: 'sap', label: 'SAP ERP', description: 'SAP Enterprise Resource Planning' },
  { value: 'netsuite', label: 'NetSuite', description: 'Oracle NetSuite ERP' },
  { value: 'odoo', label: 'Odoo', description: 'Odoo Business Management Software' },
  { value: 'dynamics365', label: 'Dynamics 365', description: 'Microsoft Dynamics 365 ERP' }
];

/**
 * Get connector status types
 * @returns {Array} Available status types
 */
export const getConnectorStatuses = () => [
  { value: 'active', label: 'Active', color: 'success' },
  { value: 'inactive', label: 'Inactive', color: 'neutral' },
  { value: 'error', label: 'Error', color: 'error' },
  { value: 'pending', label: 'Pending', color: 'warning' },
  { value: 'syncing', label: 'Syncing', color: 'info' }
];

/**
 * Get invoice statuses
 * @returns {Array} Available invoice statuses
 */
export const getInvoiceStatuses = () => [
  { value: 'draft', label: 'Draft', color: 'neutral' },
  { value: 'pending', label: 'Pending', color: 'warning' },
  { value: 'approved', label: 'Approved', color: 'success' },
  { value: 'sent', label: 'Sent', color: 'info' },
  { value: 'paid', label: 'Paid', color: 'success' },
  { value: 'overdue', label: 'Overdue', color: 'error' },
  { value: 'voided', label: 'Voided', color: 'neutral' }
];

/**
 * Get return statuses
 * @returns {Array} Available return statuses
 */
export const getReturnStatuses = () => [
  { value: 'pending', label: 'Pending', color: 'warning' },
  { value: 'approved', label: 'Approved', color: 'success' },
  { value: 'rejected', label: 'Rejected', color: 'error' },
  { value: 'processing', label: 'Processing', color: 'info' },
  { value: 'completed', label: 'Completed', color: 'success' },
  { value: 'cancelled', label: 'Cancelled', color: 'neutral' }
];

/**
 * Format currency amount
 * @param {number} amount - Amount in cents
 * @param {string} currency - Currency code
 * @returns {string} Formatted currency
 */
export const formatCurrency = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 2
  }).format(amount / 100); // Assuming amounts are in cents
};

/**
 * Format date for display
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted date
 */
export const formatDate = (dateString) => {
  if (!dateString) return '-';
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(dateString));
};

export default {
  // Connector Management
  listConnectors,
  registerConnector,
  getConnectorStatus,
  getSyncHistory,
  
  // Data Synchronization
  syncCustomers,
  syncProducts,
  syncPrices,
  syncStock,
  
  // Order Management
  sendOrder,
  
  // Inventory Management
  createInventoryReservation,
  fulfillInventoryReservation,
  releaseInventoryReservation,
  transferInventoryReservation,
  
  // Invoice Management
  createInvoice,
  approveInvoice,
  voidInvoice,
  sendInvoice,
  recordInvoicePayment,
  
  // Returns Management
  createReturn,
  approveReturn,
  rejectReturn,
  processReturnStart,
  completeReturn,
  restockReturnItems,
  
  // Utilities
  getConnectorTypes,
  getConnectorStatuses,
  getInvoiceStatuses,
  getReturnStatuses,
  formatCurrency,
  formatDate
}; 
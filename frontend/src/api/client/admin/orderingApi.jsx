
import axiosInstance from "../../axiosInstance";

const ORDERING_API_BASE_URL = '/api/ordering';

/**
 * Safe URL component encoder
 */
const safeEncode = (component) => {
    if (component === null || component === undefined) {
        return '';
    }
    return encodeURIComponent(String(component));
};

/**
 * Enhanced error handling for admin orders API
 */
const handleAdminOrdersError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        userMessage: 'Failed to process order request',
        severity: 'error',
        retryable: false
    };

    if (error.response) {
        const status = error.response.status;
        const data = error.response.data;

        errorDetails.status = status;
        errorDetails.retryable = status >= 500 || status === 408;

        if (status === 404) {
            errorDetails.userMessage = 'Order not found';
            errorDetails.severity = 'warning';
        } else if (status === 403) {
            errorDetails.userMessage = 'Admin access required';
            errorDetails.severity = 'error';
        } else if (status >= 500) {
            errorDetails.userMessage = 'Server error, please try again';
            errorDetails.severity = 'error';
        } else if (status === 400) {
            errorDetails.userMessage = data?.message || 'Invalid request';
            errorDetails.severity = 'warning';
        }

        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
    } else if (error.request) {
        errorDetails.userMessage = 'Network error or server not responding';
        errorDetails.severity = 'error';
        errorDetails.retryable = true;
        errorDetails.network = true;
    } else {
        errorDetails.userMessage = 'Unexpected error occurred';
        errorDetails.severity = 'error';
        errorDetails.message = error.message;
    }

    return errorDetails;
};

/**
 * CREATE ORDER
 * POST /api/ordering
 */
export const createOrder = async (orderData) => {
    const endpoint = ORDERING_API_BASE_URL;
    try {
        const response = await axiosInstance.post(endpoint, orderData);
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'createOrder');
    }
};

/**
 * GET ORDER BY ID
 * GET /api/ordering/{id}
 */
export const getOrder = async (orderId) => {
    if (!orderId) {
        throw new Error('orderId is required for getOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'getOrder');
    }
};

/**
 * LIST ALL ORDERS
 * GET /api/ordering
 */
export const listOrders = async (params = {}) => {
    const endpoint = ORDERING_API_BASE_URL;
    try {
        const response = await axiosInstance.get(endpoint, {
            params: {
                page: params.page || 1,
                pageSize: params.pageSize || 20,
                sortBy: params.sortBy || 'created_at',
                sortOrder: params.sortOrder || 'desc'
            }
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'listOrders');
    }
};

/**
 * GET ORDERS BY CUSTOMER
 * GET /api/ordering/customer/{userCustomerId}
 */
export const getOrdersByCustomer = async (userCustomerId, params = {}) => {
    if (!userCustomerId) {
        throw new Error('userCustomerId is required for getOrdersByCustomer.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/customer/${safeEncode(userCustomerId)}`;
    try {
        const response = await axiosInstance.get(endpoint, {
            params: {
                page: params.page || 1,
                pageSize: params.pageSize || 20,
                sortBy: params.sortBy || 'created_at',
                sortOrder: params.sortOrder || 'desc'
            }
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'getOrdersByCustomer');
    }
};

/**
 * GET ORDERS BY STATUS
 * GET /api/ordering/status/{status}
 */
export const getOrdersByStatus = async (status, params = {}) => {
    if (!status) {
        throw new Error('status is required for getOrdersByStatus.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/status/${safeEncode(status)}`;
    try {
        const response = await axiosInstance.get(endpoint, {
            params: {
                page: params.page || 1,
                pageSize: params.pageSize || 20,
                sortBy: params.sortBy || 'created_at',
                sortOrder: params.sortOrder || 'desc'
            }
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'getOrdersByStatus');
    }
};

/**
 * CANCEL ORDER
 * DELETE /api/ordering/{id}
 */
export const cancelOrder = async (orderId, reason) => {
    if (!orderId) {
        throw new Error('orderId is required for cancelOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}`;
    try {
        const response = await axiosInstance.delete(endpoint, {
            params: { reason }
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'cancelOrder');
    }
};

/**
 * MARK ORDER AS READY
 * POST /api/ordering/{id}/ready
 */
export const readyOrder = async (orderId) => {
    if (!orderId) {
        throw new Error('orderId is required for readyOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/ready`;
    try {
        const response = await axiosInstance.post(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'readyOrder');
    }
};

/**
 * COMPLETE ORDER
 * POST /api/ordering/{id}/complete
 */
export const completeOrder = async (orderId, invoiceId) => {
    if (!orderId) {
        throw new Error('orderId is required for completeOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/complete`;
    try {
        const response = await axiosInstance.post(endpoint, {
            invoiceId: invoiceId
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'completeOrder');
    }
};

/**
 * APPROVE ORDER
 * POST /api/ordering/{id}/approve
 */
export const approveOrder = async (orderId, shoppingId) => {
    if (!orderId) {
        throw new Error('orderId is required for approveOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/approve`;
    try {
        const response = await axiosInstance.post(endpoint, {
            shoppingId: shoppingId
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'approveOrder');
    }
};

/**
 * REJECT ORDER
 * POST /api/ordering/{id}/reject
 */
export const rejectOrder = async (orderId) => {
    if (!orderId) {
        throw new Error('orderId is required for rejectOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/reject`;
    try {
        const response = await axiosInstance.post(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'rejectOrder');
    }
};

/**
 * SHIP ORDER
 * POST /api/ordering/{id}/ship
 */
export const shipOrder = async (orderId) => {
    if (!orderId) {
        throw new Error('orderId is required for shipOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/ship`;
    try {
        const response = await axiosInstance.post(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'shipOrder');
    }
};

/**
 * DELIVER ORDER
 * POST /api/ordering/{id}/deliver
 */
export const deliverOrder = async (orderId) => {
    if (!orderId) {
        throw new Error('orderId is required for deliverOrder.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/deliver`;
    try {
        const response = await axiosInstance.post(endpoint, {});
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'deliverOrder');
    }
};

/**
 * UPDATE ORDER STATUS
 * PUT /api/ordering/{id}/status
 */
export const updateOrderStatus = async (orderId, status, reason) => {
    if (!orderId) {
        throw new Error('orderId is required for updateOrderStatus.');
    }
    if (!status) {
        throw new Error('status is required for updateOrderStatus.');
    }
    const endpoint = `${ORDERING_API_BASE_URL}/${safeEncode(orderId)}/status`;
    try {
        const response = await axiosInstance.put(endpoint, {
            status: status,
            reason: reason
        });
        return { success: true, ...response.data };
    } catch (error) {
        return handleAdminOrdersError(error, endpoint, 'updateOrderStatus');
    }
};

/**
 * Helper function to get order status types with user-friendly labels
 */
export const getOrderStatusTypes = () => [
    { value: 'PENDING', label: 'Pending', icon: 'Clock', color: '#f59e0b' },
    { value: 'APPROVED', label: 'Approved', icon: 'CheckCircle2', color: '#3b82f6' },
    { value: 'REJECTED', label: 'Rejected', icon: 'XCircle', color: '#dc2626' },
    { value: 'CANCELED', label: 'Cancelled', icon: 'Ban', color: '#ef4444' },
    { value: 'READY', label: 'Ready', icon: 'Package', color: '#8b5cf6' },
    { value: 'SHIPPED', label: 'Shipped', icon: 'Truck', color: '#06b6d4' },
    { value: 'DELIVERED', label: 'Delivered', icon: 'PackageCheck', color: '#10b981' },
    { value: 'COMPLETED', label: 'Completed', icon: 'CheckCircle', color: '#059669' }
];

/**
 * Calculate total amount for an order
 */
export const calculateOrderTotal = (items) => {
    if (!items || !Array.isArray(items)) return 0;
    return items.reduce((total, item) => {
        return total + (item.price * item.quantity);
    }, 0);
};

/**
 * Format order date for display
 */
export const formatOrderDate = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

/**
 * Get status badge color for admin UI
 */
export const getStatusBadgeVariant = (status) => {
    const statusMap = {
        'PENDING': 'warning',
        'APPROVED': 'info',
        'REJECTED': 'danger',
        'CANCELED': 'danger',
        'READY': 'secondary',
        'SHIPPED': 'primary',
        'DELIVERED': 'success',
        'COMPLETED': 'success'
    };
    return statusMap[status] || 'default';
};

/**
 * Utility function to check if an orders API response was successful
 */
export const isOrdersResponseSuccess = (response) => {
    return response && response.success === true;
};

/**
 * Get user-friendly error message from orders API response
 */
export const getOrdersErrorMessage = (response) => {
    if (isOrdersResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Failed to process order request';
};

/**
 * Legacy method for backward compatibility
 * @deprecated Use getOrdersByCustomer instead
 */
export const getOrders = async (userId) => {
    return getOrdersByCustomer(userId);
};

import { useState, useEffect, useCallback } from 'react';
import {
    createOrder,
    getOrder,
    listOrders,
    getOrdersByCustomer,
    getOrdersByStatus,
    cancelOrder,
    readyOrder,
    completeOrder,
    approveOrder,
    rejectOrder,
    shipOrder,
    deliverOrder,
    updateOrderStatus,
    isOrdersResponseSuccess,
    getOrdersErrorMessage,
    calculateOrderTotal
} from '@/api/client/orderingApi';

/**
 * Custom hook for managing orders
 */
export const useOrders = () => {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({
        currentPage: 1,
        totalPages: 1,
        totalCount: 0,
        pageSize: 20
    });

    // Clear error
    const clearError = useCallback(() => {
        setError(null);
    }, []);

    // Fetch all orders
    const fetchOrders = useCallback(async (options = {}) => {
        try {
            setLoading(true);
            setError(null);

            const response = await listOrders({
                page: options.page || pagination.currentPage,
                pageSize: options.pageSize || pagination.pageSize,
                sortBy: options.sortBy || 'created_at',
                sortOrder: options.sortOrder || 'desc'
            });

            if (isOrdersResponseSuccess(response)) {
                setOrders(response.orders || []);
                setPagination({
                    currentPage: response.currentPage || options.page || 1,
                    totalPages: response.totalPages || 1,
                    totalCount: response.totalCount || 0,
                    pageSize: options.pageSize || pagination.pageSize
                });
                return { success: true, orders: response.orders };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to fetch orders';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, [pagination.currentPage, pagination.pageSize]);

    // Fetch orders by customer
    const fetchCustomerOrders = useCallback(async (customerId, options = {}) => {
        if (!customerId) {
            setError('Customer ID is required');
            return { success: false, error: 'Customer ID is required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await getOrdersByCustomer(customerId, {
                page: options.page || 1,
                pageSize: options.pageSize || 20,
                sortBy: options.sortBy || 'created_at',
                sortOrder: options.sortOrder || 'desc'
            });

            if (isOrdersResponseSuccess(response)) {
                setOrders(response.orders || []);
                setPagination({
                    currentPage: response.currentPage || 1,
                    totalPages: response.totalPages || 1,
                    totalCount: response.totalCount || 0,
                    pageSize: options.pageSize || 20
                });
                return { success: true, orders: response.orders };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to fetch customer orders';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Fetch orders by status
    const fetchOrdersByStatus = useCallback(async (status, options = {}) => {
        if (!status) {
            setError('Status is required');
            return { success: false, error: 'Status is required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await getOrdersByStatus(status, {
                page: options.page || 1,
                pageSize: options.pageSize || 20,
                sortBy: options.sortBy || 'created_at',
                sortOrder: options.sortOrder || 'desc'
            });

            if (isOrdersResponseSuccess(response)) {
                setOrders(response.orders || []);
                setPagination({
                    currentPage: response.currentPage || 1,
                    totalPages: response.totalPages || 1,
                    totalCount: response.totalCount || 0,
                    pageSize: options.pageSize || 20
                });
                return { success: true, orders: response.orders };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to fetch orders by status';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Get single order
    const fetchOrder = useCallback(async (orderId) => {
        if (!orderId) {
            return { success: false, error: 'Order ID is required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await getOrder(orderId);

            if (isOrdersResponseSuccess(response)) {
                return { success: true, order: response.order };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to fetch order details';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Create new order
    const createNewOrder = useCallback(async (orderData) => {
        try {
            setLoading(true);
            setError(null);

            const response = await createOrder(orderData);

            if (isOrdersResponseSuccess(response)) {
                // Refresh orders list
                await fetchOrders();
                return { success: true, orderId: response.id };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to create order';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, [fetchOrders]);

    // Cancel order
    const cancelExistingOrder = useCallback(async (orderId, reason) => {
        if (!orderId) {
            return { success: false, error: 'Order ID is required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await cancelOrder(orderId, reason);

            if (isOrdersResponseSuccess(response)) {
                // Update local state
                setOrders(prevOrders => 
                    prevOrders.map(order => 
                        order.id === orderId 
                            ? { ...order, status: 'CANCELED' }
                            : order
                    )
                );
                return { success: true };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to cancel order';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Update order status
    const updateStatus = useCallback(async (orderId, status, reason) => {
        if (!orderId || !status) {
            return { success: false, error: 'Order ID and status are required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await updateOrderStatus(orderId, status, reason);

            if (isOrdersResponseSuccess(response)) {
                // Update local state
                setOrders(prevOrders => 
                    prevOrders.map(order => 
                        order.id === orderId 
                            ? { ...order, status: response.status }
                            : order
                    )
                );
                return { success: true };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to update order status';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Process order lifecycle actions
    const processOrderAction = useCallback(async (orderId, action, additionalData = {}) => {
        if (!orderId || !action) {
            return { success: false, error: 'Order ID and action are required' };
        }

        try {
            setLoading(true);
            setError(null);

            let response;
            switch (action) {
                case 'approve':
                    response = await approveOrder(orderId, additionalData.shoppingId);
                    break;
                case 'reject':
                    response = await rejectOrder(orderId);
                    break;
                case 'ready':
                    response = await readyOrder(orderId);
                    break;
                case 'ship':
                    response = await shipOrder(orderId);
                    break;
                case 'deliver':
                    response = await deliverOrder(orderId);
                    break;
                case 'complete':
                    response = await completeOrder(orderId, additionalData.invoiceId);
                    break;
                default:
                    throw new Error(`Unknown action: ${action}`);
            }

            if (isOrdersResponseSuccess(response)) {
                // Update local state
                setOrders(prevOrders => 
                    prevOrders.map(order => 
                        order.id === orderId 
                            ? { ...order, status: response.status }
                            : order
                    )
                );
                return { success: true };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = `Failed to ${action} order`;
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, []);

    // Calculate order statistics
    const calculateStats = useCallback(() => {
        const stats = {
            total: orders.length,
            pending: 0,
            processing: 0,
            completed: 0,
            cancelled: 0,
            totalRevenue: 0
        };

        orders.forEach(order => {
            switch (order.status) {
                case 'PENDING':
                    stats.pending++;
                    break;
                case 'APPROVED':
                case 'READY':
                case 'SHIPPED':
                    stats.processing++;
                    break;
                case 'COMPLETED':
                case 'DELIVERED':
                    stats.completed++;
                    stats.totalRevenue += calculateOrderTotal(order.items);
                    break;
                case 'CANCELED':
                case 'REJECTED':
                    stats.cancelled++;
                    break;
            }
        });

        return stats;
    }, [orders]);

    // Pagination helpers
    const goToPage = useCallback((page) => {
        if (page >= 1 && page <= pagination.totalPages) {
            fetchOrders({ page });
        }
    }, [pagination.totalPages, fetchOrders]);

    const nextPage = useCallback(() => {
        goToPage(pagination.currentPage + 1);
    }, [pagination.currentPage, goToPage]);

    const previousPage = useCallback(() => {
        goToPage(pagination.currentPage - 1);
    }, [pagination.currentPage, goToPage]);

    return {
        // State
        orders,
        loading,
        error,
        pagination,
        stats: calculateStats(),

        // Actions
        fetchOrders,
        fetchCustomerOrders,
        fetchOrdersByStatus,
        fetchOrder,
        createOrder: createNewOrder,
        cancelOrder: cancelExistingOrder,
        updateOrderStatus: updateStatus,
        processOrderAction,
        
        // Pagination
        goToPage,
        nextPage,
        previousPage,
        
        // Utilities
        clearError,
        calculateOrderTotal
    };
};

/**
 * Hook for managing a single order
 */
export const useOrder = (orderId) => {
    const [order, setOrder] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    // Fetch order on mount or when ID changes
    useEffect(() => {
        if (orderId) {
            fetchOrderDetails();
        }
    }, [orderId]);

    // Fetch order details
    const fetchOrderDetails = useCallback(async () => {
        if (!orderId) return;

        try {
            setLoading(true);
            setError(null);

            const response = await getOrder(orderId);

            if (isOrdersResponseSuccess(response)) {
                setOrder(response.order);
            } else {
                setError(getOrdersErrorMessage(response));
            }
        } catch (err) {
            setError('Failed to fetch order details');
        } finally {
            setLoading(false);
        }
    }, [orderId]);

    // Update order status
    const updateStatus = useCallback(async (status, reason) => {
        if (!orderId || !status) {
            return { success: false, error: 'Order ID and status are required' };
        }

        try {
            setLoading(true);
            setError(null);

            const response = await updateOrderStatus(orderId, status, reason);

            if (isOrdersResponseSuccess(response)) {
                // Update local state
                setOrder(prevOrder => ({
                    ...prevOrder,
                    status: response.status,
                    updatedAt: response.updatedAt
                }));
                return { success: true };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = 'Failed to update order status';
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, [orderId]);

    // Process order action
    const processAction = useCallback(async (action, additionalData = {}) => {
        if (!orderId || !action) {
            return { success: false, error: 'Order ID and action are required' };
        }

        try {
            setLoading(true);
            setError(null);

            let response;
            switch (action) {
                case 'approve':
                    response = await approveOrder(orderId, additionalData.shoppingId);
                    break;
                case 'reject':
                    response = await rejectOrder(orderId);
                    break;
                case 'ready':
                    response = await readyOrder(orderId);
                    break;
                case 'ship':
                    response = await shipOrder(orderId);
                    break;
                case 'deliver':
                    response = await deliverOrder(orderId);
                    break;
                case 'complete':
                    response = await completeOrder(orderId, additionalData.invoiceId);
                    break;
                case 'cancel':
                    response = await cancelOrder(orderId, additionalData.reason);
                    break;
                default:
                    throw new Error(`Unknown action: ${action}`);
            }

            if (isOrdersResponseSuccess(response)) {
                // Refresh order details
                await fetchOrderDetails();
                return { success: true };
            } else {
                const errorMsg = getOrdersErrorMessage(response);
                setError(errorMsg);
                return { success: false, error: errorMsg };
            }
        } catch (err) {
            const errorMsg = `Failed to ${action} order`;
            setError(errorMsg);
            return { success: false, error: errorMsg };
        } finally {
            setLoading(false);
        }
    }, [orderId, fetchOrderDetails]);

    return {
        order,
        loading,
        error,
        refresh: fetchOrderDetails,
        updateStatus,
        processAction,
        clearError: () => setError(null)
    };
};
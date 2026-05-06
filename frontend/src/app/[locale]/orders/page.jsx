"use client";
export const dynamic = 'force-dynamic';
import React, { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { 
    Package, Eye, RefreshCw, X, Clock, CheckCircle, 
    Truck, AlertCircle, Loader2, ShoppingCart, ArrowLeft 
} from '@/icons';
import { getOrders, getOrder, cancelOrder, completeOrder, readyOrder } from "../../../api/client/orderingApi";
import { useAuth } from "../../../context/AuthContext";
import styles from "./page.module.css";

/**
 * Get order status info with icon and color
 */
const getOrderStatusInfo = (status) => {
    const statusMap = {
        pending: { label: 'Pending', icon: 'Clock', color: '#f59e0b' },
        processing: { label: 'Processing', icon: 'Loader2', color: '#3b82f6' },
        ready: { label: 'Ready', icon: 'Package', color: '#8b5cf6' },
        shipped: { label: 'Shipped', icon: 'Truck', color: '#06b6d4' },
        delivered: { label: 'Delivered', icon: 'CheckCircle', color: '#10b981' },
        completed: { label: 'Completed', icon: 'CheckCircle', color: '#059669' },
        cancelled: { label: 'Cancelled', icon: 'XCircle', color: '#ef4444' },
        failed: { label: 'Failed', icon: 'AlertCircle', color: '#dc2626' }
    };
    
    const normalizedStatus = status?.toLowerCase() || 'pending';
    return statusMap[normalizedStatus] || statusMap.pending;
};

/**
 * Icon mapping for dynamic icon rendering
 */
const iconMap = {
    Package, Eye, RefreshCw, X, Clock, CheckCircle, 
    Truck, AlertCircle, Loader2, ShoppingCart, ArrowLeft,
    XCircle: X
};

/**
 * Get icon component by name
 */
const getIcon = (iconName, props = {}) => {
    const IconComponent = iconMap[iconName];
    return IconComponent ? <IconComponent {...props} /> : null;
};

/**
 * Format date with locale support
 */
const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('pl-PL', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

/**
 * Calculate total amount from items
 */
const calculateTotal = (items) => {
    if (!items || items.length === 0) return 0;
    return items.reduce((total, item) => {
        const price = parseFloat(item.price) || 0;
        const quantity = parseInt(item.quantity) || 0;
        return total + (price * quantity);
    }, 0);
};

/**
 * Order Item Component
 */
const OrderItem = ({ order, onViewDetails, onReorder, onCancel, onComplete, onReady }) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const [actionLoading, setActionLoading] = useState(null);
    const statusInfo = getOrderStatusInfo(order.status);
    const totalAmount = order.totalAmount || calculateTotal(order.items);
    
    const canCancel = ['pending', 'processing'].includes(order.status?.toLowerCase());
    const canReorder = ['delivered', 'completed', 'cancelled'].includes(order.status?.toLowerCase());
    const canMarkReady = order.status?.toLowerCase() === 'processing';
    const canComplete = order.status?.toLowerCase() === 'delivered';
    
    const handleAction = async (action, handler) => {
        setActionLoading(action);
        try {
            await handler(order.id);
        } finally {
            setActionLoading(null);
        }
    };
    
    return (
        <div className={`${styles.orderItem} ${isExpanded ? styles.expanded : ''}`}>
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: statusInfo.color }}
            />
            
            <div className={styles.orderHeader}>
                <div className={styles.orderMeta}>
                    <div className={styles.orderPrimary}>
                        <span className={styles.orderId}>Order #{order.id.slice(-8)}</span>
                        <span className={styles.orderDate}>{formatDate(order.createdAt)}</span>
                    </div>
                    <div className={styles.orderSecondary}>
                        <span className={styles.itemCount}>
                            {order.items?.length || 0} {order.items?.length === 1 ? 'item' : 'items'}
                        </span>
                        <span className={styles.orderTotal}>€{totalAmount.toFixed(2)}</span>
                    </div>
                </div>
                
                <div className={styles.statusBadge} style={{ color: statusInfo.color }}>
                    {getIcon(statusInfo.icon, { className: styles.statusIcon })}
                    {statusInfo.label}
                </div>
            </div>
            
            {/* Items Summary - Collapsible */}
            <div className={`${styles.itemsSummary} ${isExpanded ? styles.expanded : ''}`}>
                {order.items?.map((item, idx) => (
                    <div key={idx} className={styles.item}>
                        <div className={styles.itemInfo}>
                            <span className={styles.itemName}>{item.productName}</span>
                            <span className={styles.itemSeller}>Sold by: {item.userSellerName}</span>
                        </div>
                        <div className={styles.itemPricing}>
                            <span className={styles.itemQuantity}>Qty: {item.quantity}</span>
                            <span className={styles.itemPrice}>€{(parseFloat(item.price) * parseInt(item.quantity)).toFixed(2)}</span>
                        </div>
                    </div>
                ))}
            </div>
            
            {/* Order Actions */}
            <div className={styles.orderActions}>
                <button
                    onClick={() => setIsExpanded(!isExpanded)}
                    className={`${styles.actionButton} ${styles.viewButton}`}
                >
                    <Eye className={styles.actionIcon} />
                    {isExpanded ? 'Hide Details' : 'View Details'}
                </button>
                
                {canReorder && (
                    <button
                        onClick={() => handleAction('reorder', onReorder)}
                        disabled={actionLoading === 'reorder'}
                        className={`${styles.actionButton} ${styles.reorderButton}`}
                    >
                        {actionLoading === 'reorder' ? (
                            <Loader2 className={`${styles.actionIcon} ${styles.spinning}`} />
                        ) : (
                            <RefreshCw className={styles.actionIcon} />
                        )}
                        Reorder
                    </button>
                )}
                
                {canMarkReady && (
                    <button
                        onClick={() => handleAction('ready', onReady)}
                        disabled={actionLoading === 'ready'}
                        className={`${styles.actionButton} ${styles.readyButton}`}
                    >
                        {actionLoading === 'ready' ? (
                            <Loader2 className={`${styles.actionIcon} ${styles.spinning}`} />
                        ) : (
                            <Package className={styles.actionIcon} />
                        )}
                        Mark Ready
                    </button>
                )}
                
                {canComplete && (
                    <button
                        onClick={() => handleAction('complete', onComplete)}
                        disabled={actionLoading === 'complete'}
                        className={`${styles.actionButton} ${styles.completeButton}`}
                    >
                        {actionLoading === 'complete' ? (
                            <Loader2 className={`${styles.actionIcon} ${styles.spinning}`} />
                        ) : (
                            <CheckCircle className={styles.actionIcon} />
                        )}
                        Complete
                    </button>
                )}
                
                {canCancel && (
                    <button
                        onClick={() => handleAction('cancel', onCancel)}
                        disabled={actionLoading === 'cancel'}
                        className={`${styles.actionButton} ${styles.cancelButton}`}
                    >
                        {actionLoading === 'cancel' ? (
                            <Loader2 className={`${styles.actionIcon} ${styles.spinning}`} />
                        ) : (
                            <X className={styles.actionIcon} />
                        )}
                        Cancel
                    </button>
                )}
            </div>
        </div>
    );
};

/**
 * Loading State Component
 */
const LoadingState = () => (
    <div className={styles.centerContainer}>
        <div className={styles.loadingSpinner}>
            <Loader2 className={styles.spinningLarge} />
        </div>
        <p className={styles.loadingText}>Loading your orders...</p>
    </div>
);

/**
 * Empty State Component
 */
const EmptyState = ({ isAuthenticated }) => {
    const router = useRouter();
    
    return (
        <div className={styles.centerContainer}>
            <div className={styles.emptyContainer}>
                <ShoppingCart className={styles.emptyIcon} />
                <h3 className={styles.emptyTitle}>
                    {isAuthenticated ? 'No orders yet' : 'Login Required'}
                </h3>
                <p className={styles.emptyMessage}>
                    {isAuthenticated 
                        ? 'When you place your first order, it will appear here.'
                        : 'Please log in to view your order history.'
                    }
                </p>
                {isAuthenticated && (
                    <button 
                        onClick={() => router.push('/')}
                        className={styles.browseButton}
                    >
                        <ArrowLeft className={styles.actionIcon} />
                        Start Shopping
                    </button>
                )}
            </div>
        </div>
    );
};

/**
 * Error State Component
 */
const ErrorState = ({ error, onRetry }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            <AlertCircle className={styles.errorIcon} />
            <h3 className={styles.errorTitle}>Unable to load orders</h3>
            <p className={styles.errorMessage}>{error}</p>
            <button onClick={onRetry} className={styles.retryButton}>
                <RefreshCw className={styles.actionIcon} />
                Try Again
            </button>
        </div>
    </div>
);

/**
 * Orders Component
 */
export default function Orders() {
    const { user, isLoading: authLoading } = useAuth();
    const router = useRouter();
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [refreshing, setRefreshing] = useState(false);
    
    /**
     * Fetch orders
     */
    const fetchOrders = useCallback(async (showLoading = true) => {
        const userId = user?.userId || user?.id;
        if (!userId) return;
        
        if (showLoading) setLoading(true);
        setError(null);
        
        try {
            const data = await getOrders(userId);
            setOrders(data.orders || []);
        } catch (err) {
            setError('Failed to load orders. Please try again.');
            // Error: 'Error fetching orders:', err...
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    }, [user?.userId, user?.id]);
    
    /**
     * Initial fetch
     */
    useEffect(() => {
        const userId = user?.userId || user?.id;
        if (!authLoading && userId) {
            fetchOrders();
        } else if (!authLoading && !userId) {
            setLoading(false);
        }
    }, [authLoading, user?.userId, user?.id, fetchOrders]);
    
    /**
     * Refresh orders
     */
    const handleRefresh = () => {
        setRefreshing(true);
        fetchOrders(false);
    };
    
    /**
     * View order details
     */
    const handleViewDetails = async (orderId) => {
        try {
            const orderDetails = await getOrder(orderId);
            
            // In a real app, navigate to order details page or show modal
        } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
    };
    
    /**
     * Reorder items
     */
    const handleReorder = async (orderId) => {
        const order = orders.find(o => o.id === orderId);
        if (!order) return;
        
        // In a real app, add items to cart and navigate
        
        router.push('/cart');
    };
    
    /**
     * Cancel order
     */
    const handleCancelOrder = async (orderId) => {
        try {
            await cancelOrder(orderId, 'Customer requested cancellation');
            // Update local state
            setOrders(prev => prev.map(order => 
                order.id === orderId 
                    ? { ...order, status: 'cancelled' }
                    : order
            ));
        } catch (err) {
            // Error: 'Error cancelling order:', err...
            setError('Failed to cancel order. Please try again.');
        }
    };
    
    /**
     * Mark order as ready
     */
    const handleReadyOrder = async (orderId) => {
        try {
            await readyOrder(orderId);
            // Update local state
            setOrders(prev => prev.map(order => 
                order.id === orderId 
                    ? { ...order, status: 'ready' }
                    : order
            ));
        } catch (err) {
            // Error: 'Error marking order as ready:', err...
            setError('Failed to update order status. Please try again.');
        }
    };
    
    /**
     * Complete order
     */
    const handleCompleteOrder = async (orderId) => {
        try {
            // In real app, generate invoice ID
            const invoiceId = `INV-${Date.now()}`;
            await completeOrder(orderId, invoiceId);
            // Update local state
            setOrders(prev => prev.map(order => 
                order.id === orderId 
                    ? { ...order, status: 'completed' }
                    : order
            ));
        } catch (err) {
            // Error: 'Error completing order:', err...
            setError('Failed to complete order. Please try again.');
        }
    };
    
    // Loading state
    if (authLoading || loading) {
        return <LoadingState />;
    }
    
    // Not authenticated
    const userId = user?.userId || user?.id;
    if (!userId) {
        return <EmptyState isAuthenticated={false} />;
    }
    
    // Error state
    if (error && !orders.length) {
        return <ErrorState error={error} onRetry={fetchOrders} />;
    }
    
    // Empty state
    if (orders.length === 0) {
        return <EmptyState isAuthenticated={true} />;
    }
    
    // Orders list
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>
                    <Package className={styles.titleIcon} />
                    Your Orders
                </h1>
                <button 
                    onClick={handleRefresh}
                    disabled={refreshing}
                    className={styles.refreshButton}
                    title="Refresh orders"
                >
                    <RefreshCw className={`${styles.iconSmall} ${refreshing ? styles.spinning : ''}`} />
                </button>
            </div>
            
            {error && (
                <div className={styles.errorBanner}>
                    <AlertCircle className={styles.errorIcon} />
                    {error}
                </div>
            )}
            
            <div className={styles.ordersList}>
                {orders.map((order) => (
                    <OrderItem
                        key={order.id}
                        order={order}
                        onViewDetails={handleViewDetails}
                        onReorder={handleReorder}
                        onCancel={handleCancelOrder}
                        onReady={handleReadyOrder}
                        onComplete={handleCompleteOrder}
                    />
                ))}
            </div>
        </div>
    );
}
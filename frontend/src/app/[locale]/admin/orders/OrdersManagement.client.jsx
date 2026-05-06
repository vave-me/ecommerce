'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import {
    Package,
    Eye,
    Clock,
    CheckCircle,
    XCircle,
    Truck,
    Ban,
    Download,
    RefreshCw,
    Search,
    MoreVertical,
    CheckCircle2,
    DollarSign,
    ShoppingCart
} from 'lucide-react';
import {
    listOrders,
    getOrder,
    getOrdersByStatus,
    approveOrder,
    rejectOrder,
    cancelOrder,
    readyOrder,
    shipOrder,
    deliverOrder,
    completeOrder,
    getOrderStatusTypes,
    calculateOrderTotal,
    formatOrderDate,
    getStatusBadgeVariant,
    isOrdersResponseSuccess,
    getOrdersErrorMessage
} from '@/api/client/admin/orderingApi';
import styles from './OrdersManagement.module.css';

export default function OrdersManagementClient({ locale }) {
    const router = useRouter();
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedOrder, setSelectedOrder] = useState(null);
    const [showOrderModal, setShowOrderModal] = useState(false);
    const [showStatusModal, setShowStatusModal] = useState(false);
    const [statusFilter, setStatusFilter] = useState('');
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);
    const [sortBy, setSortBy] = useState('created_at');
    const [sortOrder, setSortOrder] = useState('desc');
    const [refreshing, setRefreshing] = useState(false);
    const [actionLoading, setActionLoading] = useState(false);
    const [statusUpdateData, setStatusUpdateData] = useState({ status: '', reason: '' });
    const [showActionMenu, setShowActionMenu] = useState(null);
    const [stats, setStats] = useState({
        totalOrders: 0,
        pendingOrders: 0,
        processingOrders: 0,
        completedOrders: 0,
        totalRevenue: 0
    });

    const orderStatuses = getOrderStatusTypes();

    // Calculate stats from orders
    const calculateStats = useCallback((ordersList) => {
        const stats = {
            totalOrders: ordersList.length,
            pendingOrders: 0,
            processingOrders: 0,
            completedOrders: 0,
            totalRevenue: 0
        };

        ordersList.forEach(order => {
            switch (order.status) {
                case 'PENDING':
                    stats.pendingOrders++;
                    break;
                case 'APPROVED':
                case 'READY':
                case 'SHIPPED':
                    stats.processingOrders++;
                    break;
                case 'COMPLETED':
                case 'DELIVERED':
                    stats.completedOrders++;
                    break;
            }

            if (order.status === 'COMPLETED' || order.status === 'DELIVERED') {
                stats.totalRevenue += calculateOrderTotal(order.items);
            }
        });

        return stats;
    }, []);

    // Fetch orders
    const fetchOrders = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);

            let response;
            if (statusFilter) {
                response = await getOrdersByStatus(statusFilter, {
                    page: currentPage,
                    pageSize: 20,
                    sortBy,
                    sortOrder
                });
            } else {
                response = await listOrders({
                    page: currentPage,
                    pageSize: 20,
                    sortBy,
                    sortOrder
                });
            }

            if (isOrdersResponseSuccess(response)) {
                const ordersList = response.orders || [];
                setOrders(ordersList);
                setTotalPages(response.totalPages || 1);
                setTotalCount(response.totalCount || 0);
                setStats(calculateStats(ordersList));
            } else {
                setError(getOrdersErrorMessage(response));
            }
        } catch (err) {
            setError('Failed to fetch orders');
            // Error: 'Error fetching orders:', err...
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    }, [currentPage, statusFilter, sortBy, sortOrder, calculateStats]);

    useEffect(() => {
        fetchOrders();
    }, [fetchOrders]);

    // Close action menu when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (showActionMenu && !event.target.closest(`.${styles.actionDropdown}`)) {
                setShowActionMenu(null);
            }
        };

        document.addEventListener('click', handleClickOutside);
        return () => document.removeEventListener('click', handleClickOutside);
    }, [showActionMenu]);

    // View order details
    const handleViewOrder = async (orderId) => {
        try {
            setActionLoading(true);
            const response = await getOrder(orderId);
            if (isOrdersResponseSuccess(response)) {
                setSelectedOrder(response.order);
                setShowOrderModal(true);
            } else {
                setError(getOrdersErrorMessage(response));
            }
        } catch (err) {
            setError('Failed to fetch order details');
        } finally {
            setActionLoading(false);
        }
    };

    // Handle order status update
    const handleStatusUpdate = async () => {
        if (!selectedOrder || !statusUpdateData.status) return;

        try {
            setActionLoading(true);
            
            let response;
            switch (statusUpdateData.status) {
                case 'APPROVED':
                    response = await approveOrder(selectedOrder.id);
                    break;
                case 'REJECTED':
                    response = await rejectOrder(selectedOrder.id);
                    break;
                case 'READY':
                    response = await readyOrder(selectedOrder.id);
                    break;
                case 'SHIPPED':
                    response = await shipOrder(selectedOrder.id);
                    break;
                case 'DELIVERED':
                    response = await deliverOrder(selectedOrder.id);
                    break;
                case 'COMPLETED':
                    response = await completeOrder(selectedOrder.id);
                    break;
                case 'CANCELED':
                    response = await cancelOrder(selectedOrder.id, statusUpdateData.reason);
                    break;
                default:
                    throw new Error(`Unknown status: ${statusUpdateData.status}`);
            }

            if (isOrdersResponseSuccess(response)) {
                setShowStatusModal(false);
                setStatusUpdateData({ status: '', reason: '' });
                fetchOrders();
                setError(null);
            } else {
                setError(getOrdersErrorMessage(response));
            }
        } catch (err) {
            setError('Failed to update order status');
        } finally {
            setActionLoading(false);
        }
    };

    // Quick action handlers
    const handleQuickAction = async (orderId, action) => {
        try {
            setActionLoading(true);
            let response;

            switch (action) {
                case 'approve':
                    response = await approveOrder(orderId);
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
                    response = await completeOrder(orderId);
                    break;
                default:
                    return;
            }

            if (isOrdersResponseSuccess(response)) {
                fetchOrders();
                setError(null);
            } else {
                setError(getOrdersErrorMessage(response));
            }
        } catch (err) {
            setError(`Failed to ${action} order`);
        } finally {
            setActionLoading(false);
            setShowActionMenu(null);
        }
    };

    // Get available actions for order based on status
    const getOrderActions = (order) => {
        const actions = [];
        
        switch (order.status) {
            case 'PENDING':
                actions.push(
                    { key: 'approve', label: 'Approve', variant: 'success', icon: CheckCircle2 },
                    { key: 'reject', label: 'Reject', variant: 'danger', icon: XCircle }
                );
                break;
            case 'APPROVED':
                actions.push(
                    { key: 'ready', label: 'Mark Ready', variant: 'info', icon: Package }
                );
                break;
            case 'READY':
                actions.push(
                    { key: 'ship', label: 'Ship Order', variant: 'primary', icon: Truck }
                );
                break;
            case 'SHIPPED':
                actions.push(
                    { key: 'deliver', label: 'Mark Delivered', variant: 'success', icon: CheckCircle2 }
                );
                break;
            case 'DELIVERED':
                actions.push(
                    { key: 'complete', label: 'Complete Order', variant: 'success', icon: CheckCircle }
                );
                break;
        }

        // Cancel is available for most statuses
        if (!['COMPLETED', 'CANCELED', 'REJECTED'].includes(order.status)) {
            actions.push(
                { key: 'cancel', label: 'Cancel', variant: 'danger', icon: Ban }
            );
        }

        return actions;
    };

    // Render status badge
    const renderStatusBadge = (status) => {
        const statusInfo = orderStatuses.find(s => s.value === status);
        if (!statusInfo) return <span className={styles.badge}>{status}</span>;

        const variantMap = {
            'PENDING': 'Pending',
            'APPROVED': 'Approved',
            'REJECTED': 'Rejected',
            'CANCELED': 'Cancelled',
            'READY': 'Ready',
            'SHIPPED': 'Shipped',
            'DELIVERED': 'Delivered',
            'COMPLETED': 'Completed'
        };

        return (
            <span className={`${styles.badge} ${styles[`badge${variantMap[status]}`]}`}>
                {statusInfo.label}
            </span>
        );
    };

    // Export orders
    const handleExportOrders = async () => {
        try {
            const allOrdersResponse = await listOrders({
                page: 1,
                pageSize: 1000,
                sortBy,
                sortOrder
            });

            if (isOrdersResponseSuccess(allOrdersResponse)) {
                const orders = allOrdersResponse.orders || [];
                const csv = convertOrdersToCSV(orders);
                downloadCSV(csv, 'orders_export.csv');
            }
        } catch (err) {
            setError('Failed to export orders');
        }
    };

    // Convert orders to CSV
    const convertOrdersToCSV = (orders) => {
        const headers = ['Order ID', 'Customer ID', 'Status', 'Items Count', 'Total Amount', 'Created Date', 'Updated Date'];
        const rows = orders.map(order => [
            order.id,
            order.userCustomerId,
            order.status,
            order.items?.length || 0,
            calculateOrderTotal(order.items).toFixed(2),
            formatOrderDate(order.createdAt),
            formatOrderDate(order.updatedAt)
        ]);

        const csvContent = [
            headers.join(','),
            ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
        ].join('\n');

        return csvContent;
    };

    // Download CSV file
    const downloadCSV = (csvContent, filename) => {
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = filename;
        link.click();
    };

    // Refresh orders
    const handleRefresh = () => {
        setRefreshing(true);
        fetchOrders();
    };

    // Filter orders by search term
    const filteredOrders = orders.filter(order => {
        if (!searchTerm) return true;
        const searchLower = searchTerm.toLowerCase();
        return (
            order.id.toLowerCase().includes(searchLower) ||
            order.userCustomerId?.toLowerCase().includes(searchLower) ||
            order.items?.some(item => 
                item.productName?.toLowerCase().includes(searchLower) ||
                item.userSellerName?.toLowerCase().includes(searchLower)
            )
        );
    });

    // Pagination helpers
    const renderPagination = () => {
        const pages = [];
        const maxPages = Math.min(5, totalPages);
        
        for (let i = 1; i <= maxPages; i++) {
            pages.push(
                <button
                    key={i}
                    className={`${styles.paginationItem} ${i === currentPage ? styles.active : ''}`}
                    onClick={() => setCurrentPage(i)}
                    disabled={i === currentPage}
                >
                    {i}
                </button>
            );
        }
        
        return pages;
    };

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>Orders Management</h1>
            </div>

            {error && (
                <div className={styles.alert}>
                    <span>{error}</span>
                    <button 
                        className={styles.alertClose}
                        onClick={() => setError(null)}
                    >
                        ×
                    </button>
                </div>
            )}

            {/* Stats Cards */}
            <div className={styles.statsGrid}>
                <div className={styles.statCard}>
                    <div className={styles.statContent}>
                        <div className={styles.statInfo}>
                            <h6 className={styles.statLabel}>Total Orders</h6>
                            <h3 className={styles.statValue}>{stats.totalOrders}</h3>
                        </div>
                        <div className={`${styles.statIcon} ${styles.statIconPrimary}`}>
                            <ShoppingCart size={24} />
                        </div>
                    </div>
                </div>
                <div className={styles.statCard}>
                    <div className={styles.statContent}>
                        <div className={styles.statInfo}>
                            <h6 className={styles.statLabel}>Pending</h6>
                            <h3 className={styles.statValue}>{stats.pendingOrders}</h3>
                        </div>
                        <div className={`${styles.statIcon} ${styles.statIconWarning}`}>
                            <Clock size={24} />
                        </div>
                    </div>
                </div>
                <div className={styles.statCard}>
                    <div className={styles.statContent}>
                        <div className={styles.statInfo}>
                            <h6 className={styles.statLabel}>Processing</h6>
                            <h3 className={styles.statValue}>{stats.processingOrders}</h3>
                        </div>
                        <div className={`${styles.statIcon} ${styles.statIconInfo}`}>
                            <Package size={24} />
                        </div>
                    </div>
                </div>
                <div className={styles.statCard}>
                    <div className={styles.statContent}>
                        <div className={styles.statInfo}>
                            <h6 className={styles.statLabel}>Revenue</h6>
                            <h3 className={styles.statValue}>${stats.totalRevenue.toFixed(2)}</h3>
                        </div>
                        <div className={`${styles.statIcon} ${styles.statIconSuccess}`}>
                            <DollarSign size={24} />
                        </div>
                    </div>
                </div>
            </div>

            {/* Filters and Search */}
            <div className={styles.filtersCard}>
                <div className={styles.filtersRow}>
                    <div className={styles.searchGroup}>
                        <Search size={18} className={styles.searchIcon} />
                        <input
                            type="text"
                            className={styles.searchInput}
                            placeholder="Search orders..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <select
                        className={styles.filterSelect}
                        value={statusFilter}
                        onChange={(e) => {
                            setStatusFilter(e.target.value);
                            setCurrentPage(1);
                        }}
                    >
                        <option value="">All Statuses</option>
                        {orderStatuses.map(status => (
                            <option key={status.value} value={status.value}>
                                {status.label}
                            </option>
                        ))}
                    </select>
                    <div className={styles.buttonGroup}>
                        <button
                            className={styles.button}
                            onClick={handleRefresh}
                            disabled={refreshing}
                        >
                            <RefreshCw size={18} className={refreshing ? styles.spinning : ''} />
                            Refresh
                        </button>
                        <button
                            className={styles.button}
                            onClick={handleExportOrders}
                        >
                            <Download size={18} />
                            Export
                        </button>
                    </div>
                </div>
            </div>

            {/* Orders Table */}
            <div className={styles.tableCard}>
                {loading ? (
                    <div className={styles.loading}>
                        <div className={styles.spinner}></div>
                        <p>Loading orders...</p>
                    </div>
                ) : filteredOrders.length === 0 ? (
                    <div className={styles.empty}>
                        <Package size={48} />
                        <p>No orders found</p>
                    </div>
                ) : (
                    <>
                        <div className={styles.tableWrapper}>
                            <table className={styles.table}>
                                <thead>
                                    <tr>
                                        <th>Order ID</th>
                                        <th>Customer</th>
                                        <th>Items</th>
                                        <th>Total</th>
                                        <th>Status</th>
                                        <th>Date</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {filteredOrders.map(order => (
                                        <tr key={order.id}>
                                            <td>
                                                <code className={styles.orderId}>{order.id.slice(0, 8)}...</code>
                                            </td>
                                            <td>{order.userCustomerId?.slice(0, 12)}...</td>
                                            <td>
                                                <span className={styles.itemsBadge}>
                                                    {order.items?.length || 0} items
                                                </span>
                                            </td>
                                            <td>${calculateOrderTotal(order.items).toFixed(2)}</td>
                                            <td>{renderStatusBadge(order.status)}</td>
                                            <td>{formatOrderDate(order.createdAt)}</td>
                                            <td>
                                                <div className={styles.actions}>
                                                    <button
                                                        className={styles.actionButton}
                                                        onClick={() => handleViewOrder(order.id)}
                                                        disabled={actionLoading}
                                                        title="View Details"
                                                    >
                                                        <Eye size={16} />
                                                    </button>
                                                    
                                                    <div className={styles.actionDropdown}>
                                                        <button
                                                            className={styles.actionButton}
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                setShowActionMenu(showActionMenu === order.id ? null : order.id);
                                                            }}
                                                            disabled={actionLoading}
                                                        >
                                                            <MoreVertical size={16} />
                                                        </button>
                                                        {showActionMenu === order.id && (
                                                            <div className={styles.dropdownMenu}>
                                                                {getOrderActions(order).map(action => (
                                                                    <button
                                                                        key={action.key}
                                                                        className={styles.dropdownItem}
                                                                        onClick={() => {
                                                                            if (action.key === 'cancel') {
                                                                                setSelectedOrder(order);
                                                                                setStatusUpdateData({ 
                                                                                    status: 'CANCELED', 
                                                                                    reason: '' 
                                                                                });
                                                                                setShowStatusModal(true);
                                                                            } else {
                                                                                handleQuickAction(order.id, action.key);
                                                                            }
                                                                        }}
                                                                    >
                                                                        <action.icon size={16} />
                                                                        {action.label}
                                                                    </button>
                                                                ))}
                                                                <div className={styles.dropdownDivider}></div>
                                                                <button
                                                                    className={styles.dropdownItem}
                                                                    onClick={() => {
                                                                        setSelectedOrder(order);
                                                                        setShowStatusModal(true);
                                                                        setShowActionMenu(null);
                                                                    }}
                                                                >
                                                                    <Clock size={16} />
                                                                    Change Status
                                                                </button>
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>

                        {totalPages > 1 && (
                            <div className={styles.pagination}>
                                <div className={styles.paginationInfo}>
                                    Showing {filteredOrders.length} of {totalCount} orders
                                </div>
                                <div className={styles.paginationControls}>
                                    <button
                                        className={styles.paginationItem}
                                        onClick={() => setCurrentPage(1)}
                                        disabled={currentPage === 1}
                                    >
                                        First
                                    </button>
                                    <button
                                        className={styles.paginationItem}
                                        onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                                        disabled={currentPage === 1}
                                    >
                                        Prev
                                    </button>
                                    {renderPagination()}
                                    <button
                                        className={styles.paginationItem}
                                        onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                                        disabled={currentPage === totalPages}
                                    >
                                        Next
                                    </button>
                                    <button
                                        className={styles.paginationItem}
                                        onClick={() => setCurrentPage(totalPages)}
                                        disabled={currentPage === totalPages}
                                    >
                                        Last
                                    </button>
                                </div>
                            </div>
                        )}
                    </>
                )}
            </div>

            {/* Order Details Modal */}
            {showOrderModal && selectedOrder && (
                <div className={styles.modal} onClick={() => setShowOrderModal(false)}>
                    <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
                        <div className={styles.modalHeader}>
                            <h2>Order Details</h2>
                            <button 
                                className={styles.modalClose}
                                onClick={() => setShowOrderModal(false)}
                            >
                                ×
                            </button>
                        </div>
                        <div className={styles.modalBody}>
                            <div className={styles.detailsGrid}>
                                <div className={styles.detailItem}>
                                    <strong>Order ID:</strong>
                                    <p><code>{selectedOrder.id}</code></p>
                                </div>
                                <div className={styles.detailItem}>
                                    <strong>Status:</strong>
                                    <p>{renderStatusBadge(selectedOrder.status)}</p>
                                </div>
                                <div className={styles.detailItem}>
                                    <strong>Customer ID:</strong>
                                    <p><code>{selectedOrder.userCustomerId}</code></p>
                                </div>
                                <div className={styles.detailItem}>
                                    <strong>Payment Method:</strong>
                                    <p>{selectedOrder.paymentMethodId || 'N/A'}</p>
                                </div>
                                <div className={styles.detailItem}>
                                    <strong>Created:</strong>
                                    <p>{formatOrderDate(selectedOrder.createdAt)}</p>
                                </div>
                                <div className={styles.detailItem}>
                                    <strong>Updated:</strong>
                                    <p>{formatOrderDate(selectedOrder.updatedAt)}</p>
                                </div>
                            </div>
                            
                            <h5 className={styles.sectionTitle}>Order Items</h5>
                            <div className={styles.itemsTable}>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Product</th>
                                            <th>Seller</th>
                                            <th>Price</th>
                                            <th>Quantity</th>
                                            <th>Subtotal</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {selectedOrder.items?.map((item, index) => (
                                            <tr key={index}>
                                                <td>{item.productName}</td>
                                                <td>{item.userSellerName}</td>
                                                <td>${item.price.toFixed(2)}</td>
                                                <td>{item.quantity}</td>
                                                <td>${(item.price * item.quantity).toFixed(2)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                    <tfoot>
                                        <tr>
                                            <th colSpan={4}>Total:</th>
                                            <th>${calculateOrderTotal(selectedOrder.items).toFixed(2)}</th>
                                        </tr>
                                    </tfoot>
                                </table>
                            </div>
                        </div>
                        <div className={styles.modalFooter}>
                            <button 
                                className={`${styles.button} ${styles.buttonSecondary}`}
                                onClick={() => setShowOrderModal(false)}
                            >
                                Close
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Status Update Modal */}
            {showStatusModal && (
                <div className={styles.modal} onClick={() => {
                    setShowStatusModal(false);
                    setStatusUpdateData({ status: '', reason: '' });
                }}>
                    <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
                        <div className={styles.modalHeader}>
                            <h2>Update Order Status</h2>
                            <button 
                                className={styles.modalClose}
                                onClick={() => {
                                    setShowStatusModal(false);
                                    setStatusUpdateData({ status: '', reason: '' });
                                }}
                            >
                                ×
                            </button>
                        </div>
                        <div className={styles.modalBody}>
                            <div className={styles.formGroup}>
                                <label>New Status</label>
                                <select
                                    className={styles.formControl}
                                    value={statusUpdateData.status}
                                    onChange={(e) => setStatusUpdateData({
                                        ...statusUpdateData,
                                        status: e.target.value
                                    })}
                                >
                                    <option value="">Select status...</option>
                                    {orderStatuses.map(status => (
                                        <option key={status.value} value={status.value}>
                                            {status.label}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            {statusUpdateData.status === 'CANCELED' && (
                                <div className={styles.formGroup}>
                                    <label>Reason (optional)</label>
                                    <textarea
                                        className={styles.formControl}
                                        rows={3}
                                        placeholder="Enter reason for cancellation..."
                                        value={statusUpdateData.reason}
                                        onChange={(e) => setStatusUpdateData({
                                            ...statusUpdateData,
                                            reason: e.target.value
                                        })}
                                    />
                                </div>
                            )}
                        </div>
                        <div className={styles.modalFooter}>
                            <button 
                                className={`${styles.button} ${styles.buttonSecondary}`}
                                onClick={() => {
                                    setShowStatusModal(false);
                                    setStatusUpdateData({ status: '', reason: '' });
                                }}
                            >
                                Cancel
                            </button>
                            <button 
                                className={`${styles.button} ${styles.buttonPrimary}`}
                                onClick={handleStatusUpdate}
                                disabled={!statusUpdateData.status || actionLoading}
                            >
                                {actionLoading ? 'Updating...' : 'Update Status'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
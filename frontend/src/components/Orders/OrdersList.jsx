"use client";
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { 
    Package, 
    Eye, 
    Clock, 
    Truck, 
    CheckCircle,
    XCircle,
    ShoppingBag,
    ChevronLeft,
    ChevronRight
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useOrders } from '@/hooks/useOrders';
import { formatOrderDate, getStatusBadgeVariant } from '@/api/client/orderingApi';
import styles from './OrdersList.module.css';

const OrdersList = ({ locale = 'en' }) => {
    const router = useRouter();
    const { user } = useAuth();
    const {
        orders,
        loading,
        error,
        pagination,
        fetchCustomerOrders,
        goToPage,
        clearError
    } = useOrders();

    const [initialLoad, setInitialLoad] = useState(true);

    useEffect(() => {
        if (user?.userId) {
            fetchCustomerOrders(user.userId).then(() => {
                setInitialLoad(false);
            });
        }
    }, [user?.userId, fetchCustomerOrders]);

    const handleViewOrder = (orderId) => {
        router.push(`/${locale}/account/orders/${orderId}`);
    };

    const calculateOrderTotal = (items) => {
        if (!items || items.length === 0) return 0;
        return items.reduce((total, item) => total + (item.price * item.quantity), 0);
    };

    const getStatusIcon = (status) => {
        const icons = {
            'PENDING': Clock,
            'APPROVED': CheckCircle,
            'REJECTED': XCircle,
            'CANCELED': XCircle,
            'READY': Package,
            'SHIPPED': Truck,
            'DELIVERED': CheckCircle,
            'COMPLETED': CheckCircle
        };
        const Icon = icons[status] || Package;
        return <Icon size={16} />;
    };

    const renderOrderStatus = (status) => {
        const statusLabels = {
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
            <span 
                className={`${styles.badge} ${styles[`badge${getStatusBadgeVariant(status)}`]}`}
            >
                {getStatusIcon(status)}
                {statusLabels[status] || status}
            </span>
        );
    };

    const renderPagination = () => {
        const pages = [];
        const maxVisiblePages = 5;
        let startPage = Math.max(1, pagination.currentPage - Math.floor(maxVisiblePages / 2));
        let endPage = Math.min(pagination.totalPages, startPage + maxVisiblePages - 1);

        if (endPage - startPage + 1 < maxVisiblePages) {
            startPage = Math.max(1, endPage - maxVisiblePages + 1);
        }

        for (let i = startPage; i <= endPage; i++) {
            pages.push(i);
        }

        return pages;
    };

    if (loading && initialLoad) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingWrapper}>
                    <div className={styles.spinner}></div>
                    <p className={styles.loadingText}>Loading your orders...</p>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div>
                    <h2 className={styles.title}>My Orders</h2>
                    <p className={styles.subtitle}>Track and manage your orders</p>
                </div>
            </div>

            {error && (
                <div className={styles.errorAlert}>
                    <button 
                        className={styles.closeButton} 
                        onClick={clearError}
                        aria-label="Close"
                    >
                        ×
                    </button>
                    {error}
                </div>
            )}

            {orders.length === 0 ? (
                <div className={styles.emptyState}>
                    <ShoppingBag size={64} className={styles.emptyIcon} />
                    <h4 className={styles.emptyTitle}>No orders yet</h4>
                    <p className={styles.emptyText}>
                        When you place your first order, it will appear here.
                    </p>
                    <button 
                        className={styles.primaryButton}
                        onClick={() => router.push(`/${locale}/products`)}
                    >
                        Start Shopping
                    </button>
                </div>
            ) : (
                <>
                    <div className={styles.tableCard}>
                        <div className={styles.tableWrapper}>
                            <table className={styles.table}>
                                <thead>
                                    <tr>
                                        <th>Order #</th>
                                        <th>Date</th>
                                        <th>Items</th>
                                        <th>Total</th>
                                        <th>Status</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {orders.map(order => (
                                        <tr key={order.id}>
                                            <td>
                                                <code className={styles.orderId}>
                                                    {order.id.slice(0, 8).toUpperCase()}
                                                </code>
                                            </td>
                                            <td>{formatOrderDate(order.createdAt)}</td>
                                            <td>
                                                <span className={styles.itemsCount}>
                                                    {order.items?.length || 0} items
                                                </span>
                                            </td>
                                            <td>
                                                <strong>${calculateOrderTotal(order.items).toFixed(2)}</strong>
                                            </td>
                                            <td>{renderOrderStatus(order.status)}</td>
                                            <td>
                                                <button
                                                    className={styles.viewButton}
                                                    onClick={() => handleViewOrder(order.id)}
                                                >
                                                    <Eye size={16} />
                                                    View Details
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>

                    {pagination.totalPages > 1 && (
                        <div className={styles.pagination}>
                            <button
                                className={styles.pageButton}
                                onClick={() => goToPage(1)}
                                disabled={pagination.currentPage === 1}
                            >
                                First
                            </button>
                            <button
                                className={styles.pageButton}
                                onClick={() => goToPage(pagination.currentPage - 1)}
                                disabled={pagination.currentPage === 1}
                            >
                                <ChevronLeft size={16} />
                            </button>
                            
                            {renderPagination().map(page => (
                                <button
                                    key={page}
                                    className={`${styles.pageButton} ${
                                        page === pagination.currentPage ? styles.activePage : ''
                                    }`}
                                    onClick={() => goToPage(page)}
                                >
                                    {page}
                                </button>
                            ))}
                            
                            <button
                                className={styles.pageButton}
                                onClick={() => goToPage(pagination.currentPage + 1)}
                                disabled={pagination.currentPage === pagination.totalPages}
                            >
                                <ChevronRight size={16} />
                            </button>
                            <button
                                className={styles.pageButton}
                                onClick={() => goToPage(pagination.totalPages)}
                                disabled={pagination.currentPage === pagination.totalPages}
                            >
                                Last
                            </button>
                        </div>
                    )}
                </>
            )}
        </div>
    );
};

export default OrdersList;
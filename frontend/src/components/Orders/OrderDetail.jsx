"use client";
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { 
    Package, 
    ArrowLeft, 
    Clock, 
    Truck, 
    CheckCircle,
    XCircle,
    MapPin,
    CreditCard,
    Calendar,
    AlertCircle,
    Download
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useOrder } from '@/hooks/useOrders';
import { 
    formatOrderDate, 
    getStatusBadgeVariant,
    calculateOrderTotal 
} from '@/api/client/orderingApi';
import styles from './OrderDetail.module.css';

const OrderDetail = ({ orderId, locale = 'en' }) => {
    const router = useRouter();
    const { user } = useAuth();
    const {
        order,
        loading,
        error,
        processAction,
        clearError
    } = useOrder(orderId);

    const [showCancelModal, setShowCancelModal] = useState(false);
    const [cancelReason, setCancelReason] = useState('');
    const [actionLoading, setActionLoading] = useState(false);

    const handleCancelOrder = async () => {
        if (!cancelReason.trim()) return;
        
        setActionLoading(true);
        try {
            await processAction('cancel', { reason: cancelReason });
            setShowCancelModal(false);
            setCancelReason('');
        } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    } finally {
            setActionLoading(false);
        }
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
            'READY': 'Ready for Pickup',
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

    const renderOrderTimeline = () => {
        const timeline = [
            { status: 'PENDING', label: 'Order Placed', icon: Clock },
            { status: 'APPROVED', label: 'Order Approved', icon: CheckCircle },
            { status: 'READY', label: 'Ready for Shipping', icon: Package },
            { status: 'SHIPPED', label: 'Shipped', icon: Truck },
            { status: 'DELIVERED', label: 'Delivered', icon: CheckCircle }
        ];

        const currentStatusIndex = timeline.findIndex(t => t.status === order?.status);

        return (
            <div className={styles.timeline}>
                {timeline.map((item, index) => {
                    const isCompleted = currentStatusIndex >= index;
                    const isCurrent = currentStatusIndex === index;
                    const Icon = item.icon;

                    return (
                        <div 
                            key={item.status} 
                            className={`${styles.timelineItem} ${isCompleted ? styles.completed : ''} ${isCurrent ? styles.current : ''}`}
                        >
                            <div className={styles.timelineIcon}>
                                <Icon size={20} />
                            </div>
                            <div className={styles.timelineLabel}>{item.label}</div>
                            {index < timeline.length - 1 && (
                                <div className={styles.timelineLine} />
                            )}
                        </div>
                    );
                })}
            </div>
        );
    };

    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingWrapper}>
                    <div className={styles.spinner}></div>
                    <p className={styles.loadingText}>Loading order details...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className={styles.container}>
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
                <button 
                    className={styles.primaryButton}
                    onClick={() => router.push(`/${locale}/account/orders`)}
                >
                    <ArrowLeft size={16} />
                    Back to Orders
                </button>
            </div>
        );
    }

    if (!order) {
        return (
            <div className={styles.container}>
                <div className={styles.warningAlert}>
                    Order not found
                </div>
                <button 
                    className={styles.primaryButton}
                    onClick={() => router.push(`/${locale}/account/orders`)}
                >
                    <ArrowLeft size={16} />
                    Back to Orders
                </button>
            </div>
        );
    }

    const canCancelOrder = ['PENDING', 'APPROVED'].includes(order.status);
    const orderTotal = calculateOrderTotal(order.items);

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div>
                    <button 
                        className={styles.backButton}
                        onClick={() => router.push(`/${locale}/account/orders`)}
                    >
                        <ArrowLeft size={16} />
                        Back to Orders
                    </button>
                    <h2 className={styles.orderTitle}>Order #{order.id.slice(0, 8).toUpperCase()}</h2>
                    <p className={styles.orderDate}>
                        Placed on {formatOrderDate(order.createdAt)}
                    </p>
                </div>
                <div>
                    {renderOrderStatus(order.status)}
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

            {/* Order Timeline */}
            {!['CANCELED', 'REJECTED'].includes(order.status) && (
                <div className={styles.card}>
                    <h5 className={styles.cardTitle}>Order Progress</h5>
                    {renderOrderTimeline()}
                </div>
            )}

            <div className={styles.contentGrid}>
                <div className={styles.mainContent}>
                    {/* Order Items */}
                    <div className={styles.card}>
                        <div className={styles.cardHeader}>
                            <h5>Order Items</h5>
                        </div>
                        <div className={styles.cardBody}>
                            <table className={styles.table}>
                                <thead>
                                    <tr>
                                        <th>Product</th>
                                        <th>Seller</th>
                                        <th>Price</th>
                                        <th>Quantity</th>
                                        <th className={styles.textEnd}>Subtotal</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {order.items?.map((item, index) => (
                                        <tr key={index}>
                                            <td>
                                                <div>
                                                    <strong>{item.productName}</strong>
                                                    <br />
                                                    <small className={styles.textMuted}>
                                                        ID: {item.productId}
                                                    </small>
                                                </div>
                                            </td>
                                            <td>
                                                <div>
                                                    {item.userSellerName}
                                                    <br />
                                                    <small className={styles.textMuted}>
                                                        ID: {item.userSellerId}
                                                    </small>
                                                </div>
                                            </td>
                                            <td>${item.price.toFixed(2)}</td>
                                            <td>{item.quantity}</td>
                                            <td className={styles.textEnd}>
                                                ${(item.price * item.quantity).toFixed(2)}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                                <tfoot>
                                    <tr>
                                        <th colSpan={4} className={styles.textEnd}>Total:</th>
                                        <th className={styles.textEnd}>${orderTotal.toFixed(2)}</th>
                                    </tr>
                                </tfoot>
                            </table>
                        </div>
                    </div>
                </div>

                <div className={styles.sidebar}>
                    {/* Order Summary */}
                    <div className={styles.card}>
                        <div className={styles.cardHeader}>
                            <h5>Order Summary</h5>
                        </div>
                        <div className={styles.cardBody}>
                            <div className={styles.summaryItem}>
                                <span className={styles.textMuted}>Subtotal</span>
                                <span>${orderTotal.toFixed(2)}</span>
                            </div>
                            <div className={styles.summaryItem}>
                                <span className={styles.textMuted}>Shipping</span>
                                <span>Calculated at checkout</span>
                            </div>
                            <hr className={styles.divider} />
                            <div className={styles.summaryItem}>
                                <strong>Total</strong>
                                <strong>${orderTotal.toFixed(2)}</strong>
                            </div>

                            {canCancelOrder && (
                                <button 
                                    className={`${styles.dangerButton} ${styles.fullWidth}`}
                                    onClick={() => setShowCancelModal(true)}
                                >
                                    <XCircle size={16} />
                                    Cancel Order
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Order Information */}
                    <div className={styles.card}>
                        <div className={styles.cardHeader}>
                            <h5>Order Information</h5>
                        </div>
                        <div className={styles.cardBody}>
                            <div className={styles.infoItem}>
                                <div className={styles.infoLabel}>
                                    <Calendar size={16} />
                                    Order Date
                                </div>
                                <div>{formatOrderDate(order.createdAt)}</div>
                            </div>

                            <div className={styles.infoItem}>
                                <div className={styles.infoLabel}>
                                    <Package size={16} />
                                    Order ID
                                </div>
                                <div>
                                    <code className={styles.code}>{order.id}</code>
                                </div>
                            </div>

                            {order.paymentMethodId && (
                                <div className={styles.infoItem}>
                                    <div className={styles.infoLabel}>
                                        <CreditCard size={16} />
                                        Payment Method
                                    </div>
                                    <div>****{order.paymentMethodId.slice(-4)}</div>
                                </div>
                            )}

                            <div className={styles.infoItem}>
                                <div className={styles.infoLabel}>
                                    <MapPin size={16} />
                                    Delivery Address
                                </div>
                                <div>
                                    <small className={styles.textMuted}>
                                        Address details will be shown here
                                    </small>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Cancel Order Modal */}
            {showCancelModal && (
                <>
                    <div 
                        className={styles.modalBackdrop} 
                        onClick={() => setShowCancelModal(false)}
                    />
                    <div className={styles.modal}>
                        <div className={styles.modalHeader}>
                            <h4>Cancel Order</h4>
                            <button 
                                className={styles.modalClose}
                                onClick={() => setShowCancelModal(false)}
                                aria-label="Close"
                            >
                                ×
                            </button>
                        </div>
                        <div className={styles.modalBody}>
                            <div className={styles.warningAlert}>
                                <AlertCircle size={16} />
                                Are you sure you want to cancel this order? This action cannot be undone.
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="cancelReason">Reason for cancellation</label>
                                <textarea
                                    id="cancelReason"
                                    className={styles.textarea}
                                    rows={3}
                                    value={cancelReason}
                                    onChange={(e) => setCancelReason(e.target.value)}
                                    placeholder="Please provide a reason for cancelling this order..."
                                />
                            </div>
                        </div>
                        <div className={styles.modalFooter}>
                            <button 
                                className={styles.secondaryButton}
                                onClick={() => setShowCancelModal(false)}
                            >
                                Keep Order
                            </button>
                            <button 
                                className={styles.dangerButton}
                                onClick={handleCancelOrder}
                                disabled={actionLoading || !cancelReason.trim()}
                            >
                                {actionLoading ? (
                                    <>
                                        <span className={styles.buttonSpinner}></span>
                                        Cancelling...
                                    </>
                                ) : (
                                    'Cancel Order'
                                )}
                            </button>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};

export default OrderDetail;
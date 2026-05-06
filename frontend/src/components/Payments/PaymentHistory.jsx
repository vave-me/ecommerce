"use client";
import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { CreditCard, CheckCircle, AlertCircle, Clock, Eye, Download, Filter } from '@/icons';
import styles from './PaymentHistory.module.css';
export default function PaymentHistory({ userId }) {
    const t = useTranslations('PaymentHistory');
    const [payments, setPayments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [filter, setFilter] = useState('all'); // 'all', 'completed', 'pending', 'failed'
    const [sortBy, setSortBy] = useState('date'); // 'date', 'amount', 'status'
    const [sortOrder, setSortOrder] = useState('desc'); // 'asc', 'desc'
    useEffect(() => {
        if (!userId) {
            setLoading(false);
            return;
        }
        const fetchPaymentHistory = async () => {
            try {
                setLoading(true);
                setError(null);
                // TODO: Replace with actual API call
                // const response = await getPaymentHistory(userId);
                // setPayments(response.payments);
                // Mock payment data for now
                const mockPayments = [
                    {
                        id: 'pi_1234567890',
                        amount: 2999,
                        currency: 'eur',
                        status: 'completed',
                        paymentMethod: 'card',
                        paymentMethodDetails: {
                            card: {
                                brand: 'visa',
                                last4: '4242'
                            }
                        },
                        description: 'Order #ORD-001',
                        orderId: 'ORD-001',
                        createdAt: new Date(Date.now() - 3600000).toISOString(),
                        receiptUrl: '#'
                    },
                    {
                        id: 'pi_0987654321',
                        amount: 1550,
                        currency: 'eur',
                        status: 'pending',
                        paymentMethod: 'card',
                        paymentMethodDetails: {
                            card: {
                                brand: 'mastercard',
                                last4: '5555'
                            }
                        },
                        description: 'Order #ORD-002',
                        orderId: 'ORD-002',
                        createdAt: new Date(Date.now() - 7200000).toISOString(),
                        receiptUrl: null
                    },
                    {
                        id: 'pi_1122334455',
                        amount: 4500,
                        currency: 'eur',
                        status: 'failed',
                        paymentMethod: 'card',
                        paymentMethodDetails: {
                            card: {
                                brand: 'visa',
                                last4: '1234'
                            }
                        },
                        description: 'Order #ORD-003',
                        orderId: 'ORD-003',
                        createdAt: new Date(Date.now() - 86400000).toISOString(),
                        receiptUrl: null,
                        failureReason: 'insufficient_funds'
                    },
                    {
                        id: 'pi_5566778899',
                        amount: 7899,
                        currency: 'eur',
                        status: 'completed',
                        paymentMethod: 'paypal',
                        paymentMethodDetails: {
                            paypal: {
                                email: 'redacted-email@example.com'
                            }
                        },
                        description: 'Order #ORD-004',
                        orderId: 'ORD-004',
                        createdAt: new Date(Date.now() - 172800000).toISOString(),
                        receiptUrl: '#'
                    }
                ];
                // Simulate API delay
                await new Promise(resolve => setTimeout(resolve, 1000));
                setPayments(mockPayments);
            } catch (err) {
                setError(err.message || t('errors.fetchFailed'));
            } finally {
                setLoading(false);
            }
        };
        fetchPaymentHistory();
    }, [userId, t]);
    const getStatusIcon = (status) => {
        switch (status) {
            case 'completed':
                return <CheckCircle size={16} />;
            case 'pending':
                return <Clock size={16} />;
            case 'failed':
                return <AlertCircle size={16} />;
            default:
                return <CreditCard size={16} />;
        }
    };
    const getPaymentMethodIcon = (paymentMethod, details) => {
        if (paymentMethod === 'card') {
            const brand = details?.card?.brand;
            return (
                <div className={styles.cardBrand}>
                    {brand && (
                        <span className={styles.brandText}>
                            {brand.toUpperCase()}
                        </span>
                    )}
                    <span className={styles.cardNumber}>
                        ****{details?.card?.last4}
                    </span>
                </div>
            );
        } else if (paymentMethod === 'paypal') {
            return (
                <div className={styles.paypalMethod}>
                    <span className={styles.paypalText}>PayPal</span>
                    <span className={styles.paypalEmail}>
                        {details?.paypal?.email}
                    </span>
                </div>
            );
        }
        return paymentMethod;
    };
    const formatAmount = (amount, currency) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: currency.toUpperCase()
        }).format(amount / 100);
    };
    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };
    const filteredAndSortedPayments = payments
        .filter(payment => {
            if (filter === 'all') return true;
            return payment.status === filter;
        })
        .sort((a, b) => {
            let aValue, bValue;
            switch (sortBy) {
                case 'amount':
                    aValue = a.amount;
                    bValue = b.amount;
                    break;
                case 'status':
                    aValue = a.status;
                    bValue = b.status;
                    break;
                case 'date':
                default:
                    aValue = new Date(a.createdAt);
                    bValue = new Date(b.createdAt);
                    break;
            }
            if (sortOrder === 'asc') {
                return aValue > bValue ? 1 : -1;
            } else {
                return aValue < bValue ? 1 : -1;
            }
        });
    const handleViewDetails = (paymentId) => {
        // TODO: Navigate to payment details page
    };
    const handleDownloadReceipt = (receiptUrl) => {
        if (receiptUrl && receiptUrl !== '#') {
            window.open(receiptUrl, '_blank');
        }
    };
    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    <CreditCard size={48} />
                    <h3>{t('loading.title')}</h3>
                    <p>{t('loading.message')}</p>
                </div>
            </div>
        );
    }
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <AlertCircle size={48} />
                    <h3>{t('error.title')}</h3>
                    <p>{error}</p>
                    <button 
                        onClick={() => window.location.reload()}
                        className={styles.retryButton}
                    >
                        {t('error.retry')}
                    </button>
                </div>
            </div>
        );
    }
    if (payments.length === 0) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyState}>
                    <CreditCard size={48} />
                    <h3>{t('empty.title')}</h3>
                    <p>{t('empty.message')}</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {/* Filters and Sorting */}
            <div className={styles.controls}>
                <div className={styles.filterGroup}>
                    <Filter size={16} />
                    <select 
                        value={filter} 
                        onChange={(e) => setFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">{t('filters.all')}</option>
                        <option value="completed">{t('filters.completed')}</option>
                        <option value="pending">{t('filters.pending')}</option>
                        <option value="failed">{t('filters.failed')}</option>
                    </select>
                </div>
                <div className={styles.sortGroup}>
                    <select 
                        value={`${sortBy}-${sortOrder}`}
                        onChange={(e) => {
                            const [newSortBy, newSortOrder] = e.target.value.split('-');
                            setSortBy(newSortBy);
                            setSortOrder(newSortOrder);
                        }}
                        className={styles.sortSelect}
                    >
                        <option value="date-desc">{t('sort.dateDesc')}</option>
                        <option value="date-asc">{t('sort.dateAsc')}</option>
                        <option value="amount-desc">{t('sort.amountDesc')}</option>
                        <option value="amount-asc">{t('sort.amountAsc')}</option>
                        <option value="status-asc">{t('sort.statusAsc')}</option>
                    </select>
                </div>
            </div>
            {/* Payment List */}
            <div className={styles.paymentList}>
                {filteredAndSortedPayments.map((payment) => (
                    <div key={payment.id} className={styles.paymentItem}>
                        <div className={styles.paymentHeader}>
                            <div className={styles.paymentInfo}>
                                <div className={styles.paymentAmount}>
                                    {formatAmount(payment.amount, payment.currency)}
                                </div>
                                <div className={styles.paymentDescription}>
                                    {payment.description}
                                </div>
                            </div>
                            <div className={`${styles.paymentStatus} ${styles[payment.status]}`}>
                                {getStatusIcon(payment.status)}
                                <span>{t(`status.${payment.status}`)}</span>
                            </div>
                        </div>
                        <div className={styles.paymentDetails}>
                            <div className={styles.paymentMethod}>
                                {getPaymentMethodIcon(payment.paymentMethod, payment.paymentMethodDetails)}
                            </div>
                            <div className={styles.paymentDate}>
                                {formatDate(payment.createdAt)}
                            </div>
                        </div>
                        {payment.failureReason && (
                            <div className={styles.failureReason}>
                                <AlertCircle size={14} />
                                <span>{t(`failureReasons.${payment.failureReason}`)}</span>
                            </div>
                        )}
                        <div className={styles.paymentActions}>
                            <button 
                                onClick={() => handleViewDetails(payment.id)}
                                className={styles.actionButton}
                            >
                                <Eye size={14} />
                                {t('actions.viewDetails')}
                            </button>
                            {payment.receiptUrl && payment.status === 'completed' && (
                                <button 
                                    onClick={() => handleDownloadReceipt(payment.receiptUrl)}
                                    className={styles.actionButton}
                                >
                                    <Download size={14} />
                                    {t('actions.downloadReceipt')}
                                </button>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
} 
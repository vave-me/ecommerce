"use client";
import React, { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuth } from '../../../../context/AuthContext';
import { checkoutBasket } from '../../../../api/client/basketApi';
import { CheckCircle, Package, CreditCard, Clock, AlertCircle, Loader } from '../../../../icons';
import styles from './order-confirmation.module.css';
export default function OrderConfirmationPage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const t = useTranslations('OrderConfirmation');
    const { user } = useAuth();
    const [status, setStatus] = useState('processing');
    const [orderData, setOrderData] = useState(null);
    const [error, setError] = useState(null);
    // Get parameters from URL
    const paymentIntentId = searchParams.get('paymentIntentId');
    const basketId = searchParams.get('basketId');
    const userCustomerId = searchParams.get('userCustomerId');
    useEffect(() => {
        if (!paymentIntentId || !basketId || !userCustomerId) {
            setStatus('error');
            setError(t('errors.missingData'));
            return;
        }
        const finalizeOrder = async () => {
            try {
                setStatus('processing');
                // 1. Checkout the basket with the payment intent
                // This triggers the order creation in the backend
                try {
                    await checkoutBasket(basketId, userCustomerId, paymentIntentId);
                } catch (checkoutError) {
                    // If checkout fails with 400, the basket might already be checked out
                    if (checkoutError.response?.status === 400) {
                        
                    } else {
                        // For other errors, re-throw
                        throw checkoutError;
                    }
                }
                // The backend will create the order automatically via COSEC service
                // We don't get the order ID immediately, so we'll redirect to orders page
                setOrderData({
                    basketId,
                    paymentIntentId,
                    status: 'processing',
                    createdAt: new Date().toISOString(),
                });
                setStatus('completed');
                
                // Redirect to orders page after a short delay
                setTimeout(() => {
                    router.push('/orders');
                }, 3000);
            } catch (err) {
                setStatus('error');
                setError(err.message || t('errors.orderCreationFailed'));
            }
        };
        finalizeOrder();
    }, [paymentIntentId, basketId, userCustomerId, router, t]);
    const getStatusContent = () => {
        switch (status) {
            case 'processing':
                return {
                    icon: <Loader size={64} />,
                    title: t('status.processing.title'),
                    message: t('status.processing.message'),
                    color: '#1a73e8',
                    showDetails: false
                };
            case 'completed':
                return {
                    icon: <CheckCircle size={64} />,
                    title: t('status.completed.title'),
                    message: t('status.completed.message'),
                    color: '#059669',
                    showDetails: true
                };
            case 'error':
                return {
                    icon: <AlertCircle size={64} />,
                    title: t('status.error.title'),
                    message: error || t('status.error.message'),
                    color: '#ef4444',
                    showDetails: false
                };
            default:
                return {
                    icon: <Clock size={64} />,
                    title: t('status.processing.title'),
                    message: t('status.processing.message'),
                    color: '#f59e0b',
                    showDetails: false
                };
        }
    };
    const handleViewOrder = () => {
        if (orderData?.orderId) {
            router.push(`/orders/${orderData.orderId}`);
        } else {
            router.push('/orders');
        }
    };
    const handleContinueShopping = () => {
        router.push('/');
    };
    const handleRetry = () => {
        window.location.reload();
    };
    const statusContent = getStatusContent();
    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <div className={styles.iconContainer} style={{ color: statusContent.color }}>
                    {status === 'processing' ? (
                        <div className={styles.loadingIcon}>
                            {statusContent.icon}
                        </div>
                    ) : (
                        statusContent.icon
                    )}
                </div>
                <h1 className={styles.title}>{statusContent.title}</h1>
                <p className={styles.message}>{statusContent.message}</p>
                {statusContent.showDetails && orderData && (
                    <div className={styles.orderDetails}>
                        <h3 className={styles.detailsTitle}>
                            <Package size={20} />
                            {t('details.title')}
                        </h3>
                        <div className={styles.detailsList}>
                            <div className={styles.detailRow}>
                                <span className={styles.detailLabel}>{t('details.orderId')}</span>
                                <span className={styles.detailValue}>{orderData.orderId}</span>
                            </div>
                            <div className={styles.detailRow}>
                                <span className={styles.detailLabel}>{t('details.paymentId')}</span>
                                <span className={styles.detailValue}>{orderData.paymentIntentId}</span>
                            </div>
                            <div className={styles.detailRow}>
                                <span className={styles.detailLabel}>{t('details.status')}</span>
                                <span className={`${styles.detailValue} ${styles.statusBadge} ${styles.completed}`}>
                                    {t('details.statusCompleted')}
                                </span>
                            </div>
                            <div className={styles.detailRow}>
                                <span className={styles.detailLabel}>{t('details.date')}</span>
                                <span className={styles.detailValue}>
                                    {new Date(orderData.createdAt).toLocaleDateString()}
                                </span>
                            </div>
                        </div>
                    </div>
                )}
                {status === 'completed' && (
                    <div className={styles.redirectInfo}>
                        <p className={styles.redirectText}>
                            {t('redirect.message')}
                        </p>
                        <div className={styles.progressBar}>
                            <div className={styles.progressFill}></div>
                        </div>
                    </div>
                )}
                <div className={styles.buttonContainer}>
                    {status === 'completed' && (
                        <>
                            <button 
                                onClick={handleViewOrder}
                                className={styles.primaryButton}
                            >
                                <Package size={16} />
                                {t('actions.viewOrder')}
                            </button>
                            <button 
                                onClick={handleContinueShopping}
                                className={styles.secondaryButton}
                            >
                                {t('actions.continueShopping')}
                            </button>
                        </>
                    )}
                    {status === 'error' && (
                        <>
                            <button 
                                onClick={handleRetry}
                                className={styles.primaryButton}
                            >
                                {t('actions.retry')}
                            </button>
                            <button 
                                onClick={handleContinueShopping}
                                className={styles.secondaryButton}
                            >
                                {t('actions.goHome')}
                            </button>
                        </>
                    )}
                </div>
                {status === 'completed' && (
                    <div className={styles.successBanner}>
                        <CreditCard size={16} />
                        <span>{t('success.banner')}</span>
                    </div>
                )}
            </div>
        </div>
    );
} 
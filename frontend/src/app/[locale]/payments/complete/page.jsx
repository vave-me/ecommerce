"use client";
import React, { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { loadStripe } from '@stripe/stripe-js';
import { CheckCircle, AlertCircle, Clock, CreditCard, Loader } from '@/icons';
import styles from './complete.module.css';
const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || '{stripe_publishable_key}');
export default function PaymentCompletePage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const t = useTranslations('PaymentComplete');
    const [status, setStatus] = useState('loading');
    const [paymentIntent, setPaymentIntent] = useState(null);
    const [error, setError] = useState(null);
    // Get parameters from URL
    const paymentIntentClientSecret = searchParams.get('payment_intent_client_secret');
    const basketId = searchParams.get('basketId');
    const userCustomerId = searchParams.get('userCustomerId');
    useEffect(() => {
        if (!paymentIntentClientSecret) {
            setStatus('no_client_secret');
            setError(t('errors.noPaymentInfo'));
            return;
        }
        const checkPaymentStatus = async () => {
            try {
                const stripe = await stripePromise;
                if (!stripe) {
                    throw new Error('Stripe not loaded');
                }
                const { paymentIntent: pi, error } = await stripe.retrievePaymentIntent(paymentIntentClientSecret);
                if (error) {
                    setStatus('error');
                    setError(error.message);
                    return;
                }
                if (pi) {
                    setPaymentIntent(pi);
                    setStatus(pi.status);
                    // If payment succeeded, redirect to order confirmation
                    if (pi.status === 'succeeded') {
                        // Small delay to show success message
                        setTimeout(() => {
                            router.push(`/payments/order-confirmation?paymentIntentId=${pi.id}&basketId=${basketId}&userCustomerId=${userCustomerId}`);
                        }, 2000);
                    }
                } else {
                    setStatus('error');
                    setError(t('errors.noPaymentIntent'));
                }
            } catch (err) {
                setStatus('error');
                setError(err.message || t('errors.statusCheckFailed'));
            }
        };
        checkPaymentStatus();
    }, [paymentIntentClientSecret, basketId, userCustomerId, router, t]);
    const getStatusContent = () => {
        switch (status) {
            case 'succeeded':
                return {
                    icon: <CheckCircle size={64} />,
                    title: t('status.succeeded.title'),
                    message: t('status.succeeded.message'),
                    color: '#059669',
                    showPaymentId: true,
                    autoRedirect: true
                };
            case 'processing':
                return {
                    icon: <Clock size={64} />,
                    title: t('status.processing.title'),
                    message: t('status.processing.message'),
                    color: '#f59e0b',
                    showPaymentId: true,
                    autoRedirect: false
                };
            case 'requires_payment_method':
                return {
                    icon: <AlertCircle size={64} />,
                    title: t('status.failed.title'),
                    message: t('status.failed.message'),
                    color: '#ef4444',
                    showPaymentId: false,
                    autoRedirect: false
                };
            case 'requires_action':
                return {
                    icon: <CreditCard size={64} />,
                    title: t('status.requiresAction.title'),
                    message: t('status.requiresAction.message'),
                    color: '#f59e0b',
                    showPaymentId: true,
                    autoRedirect: false
                };
            case 'no_client_secret':
                return {
                    icon: <AlertCircle size={64} />,
                    title: t('status.noInfo.title'),
                    message: t('status.noInfo.message'),
                    color: '#ef4444',
                    showPaymentId: false,
                    autoRedirect: false
                };
            case 'loading':
                return {
                    icon: <Loader size={64} />,
                    title: t('status.loading.title'),
                    message: t('status.loading.message'),
                    color: '#1a73e8',
                    showPaymentId: false,
                    autoRedirect: false
                };
            case 'error':
            default:
                return {
                    icon: <AlertCircle size={64} />,
                    title: t('status.error.title'),
                    message: error || t('status.error.message'),
                    color: '#ef4444',
                    showPaymentId: false,
                    autoRedirect: false
                };
        }
    };
    const handleRetryPayment = () => {
        if (basketId && userCustomerId) {
            router.push(`/payments?basketId=${basketId}&userCustomerId=${userCustomerId}`);
        } else {
            router.push('/cart');
        }
    };
    const handleGoHome = () => {
        router.push('/');
    };
    const statusContent = getStatusContent();
    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <div className={styles.iconContainer} style={{ color: statusContent.color }}>
                    {status === 'loading' ? (
                        <div className={styles.loadingIcon}>
                            {statusContent.icon}
                        </div>
                    ) : (
                        statusContent.icon
                    )}
                </div>
                <h1 className={styles.title}>{statusContent.title}</h1>
                <p className={styles.message}>{statusContent.message}</p>
                {statusContent.showPaymentId && paymentIntent && (
                    <div className={styles.paymentDetails}>
                        <div className={styles.detailRow}>
                            <span className={styles.detailLabel}>{t('details.paymentId')}</span>
                            <span className={styles.detailValue}>{paymentIntent.id}</span>
                        </div>
                        <div className={styles.detailRow}>
                            <span className={styles.detailLabel}>{t('details.amount')}</span>
                            <span className={styles.detailValue}>
                                €{(paymentIntent.amount / 100).toFixed(2)}
                            </span>
                        </div>
                        <div className={styles.detailRow}>
                            <span className={styles.detailLabel}>{t('details.status')}</span>
                            <span className={`${styles.detailValue} ${styles.statusBadge} ${styles[status]}`}>
                                {t(`status.${status}.badge`)}
                            </span>
                        </div>
                    </div>
                )}
                {statusContent.autoRedirect && (
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
                    {status === 'requires_payment_method' && (
                        <button 
                            onClick={handleRetryPayment}
                            className={styles.primaryButton}
                        >
                            {t('actions.retryPayment')}
                        </button>
                    )}
                    {status === 'processing' && (
                        <button 
                            onClick={() => router.push('/orders')}
                            className={styles.secondaryButton}
                        >
                            {t('actions.viewOrders')}
                        </button>
                    )}
                    {(status === 'error' || status === 'no_client_secret') && (
                        <>
                            <button 
                                onClick={handleRetryPayment}
                                className={styles.primaryButton}
                            >
                                {t('actions.backToPayment')}
                            </button>
                            <button 
                                onClick={handleGoHome}
                                className={styles.secondaryButton}
                            >
                                {t('actions.goHome')}
                            </button>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
} 
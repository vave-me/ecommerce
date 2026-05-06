"use client";
import React, { useState, useEffect } from 'react';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import { useTranslations } from 'next-intl';
import { authorizePayment } from '../../api/client/paymentsApi';
import { getTotalBasketAmount } from '../../api/client/basketApi';
import { AlertCircle, CreditCard, Shield, Loader } from '@/icons';
import CheckoutForm from './CheckoutForm';
import styles from './PaymentForm.module.css';
// Initialize Stripe with your publishable key
const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || '{stripe_publishable_key}');
export default function PaymentForm({ basketId, userCustomerId, amount, onPaymentComplete }) {
    const t = useTranslations('PaymentForm');
    const [clientSecret, setClientSecret] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [basketTotal, setBasketTotal] = useState(null);
    const [appearance] = useState({
        theme: 'stripe',
        variables: {
            colorPrimary: '#1a73e8',
            colorBackground: '#ffffff',
            colorText: '#202124',
            colorDanger: '#ea4335',
            borderRadius: '8px',
            fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
            fontSize: '16px',
            spacingUnit: '4px',
        },
    });
    useEffect(() => {
        if (!basketId || !userCustomerId) {
            setError(t('errors.missingData'));
            setLoading(false);
            return;
        }
        const fetchPaymentIntent = async () => {
            try {
                setLoading(true);
                setError(null);
                // 1. Get basket total if not provided
                let totalAmount = amount;
                if (!totalAmount) {
                    const basketData = await getTotalBasketAmount(basketId);
                    totalAmount = basketData.amount;
                }
                const amountInMinorUnits = parseInt(totalAmount, 10);
                setBasketTotal({
                    amount: amountInMinorUnits,
                    currency: 'eur',
                    displayAmount: (amountInMinorUnits / 100).toFixed(2)
                });
                // 2. Create payment intent with basket ID
                const paymentData = await authorizePayment(userCustomerId, amountInMinorUnits, basketId);
                if (paymentData.client_secret) {
                    setClientSecret(paymentData.client_secret);
                } else if (paymentData.clientSecret) {
                    // Fallback for different response format
                    setClientSecret(paymentData.clientSecret);
                } else {
                    throw new Error('No client secret received');
                }
            } catch (err) {
                setError(err.message || t('errors.paymentSetup'));
            } finally {
                setLoading(false);
            }
        };
        fetchPaymentIntent();
    }, [basketId, userCustomerId, amount, t]);
    const handlePaymentSuccess = (paymentResult) => {
        if (onPaymentComplete) {
            onPaymentComplete(paymentResult);
        }
    };
    const handlePaymentError = (error) => {
        setError(error.message || t('errors.paymentFailed'));
    };
    const options = {
        clientSecret,
        appearance,
    };
    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    <div className={styles.loadingIcon}>
                        <Loader size={32} />
                    </div>
                    <h3 className={styles.loadingTitle}>{t('loading.title')}</h3>
                    <p className={styles.loadingMessage}>{t('loading.message')}</p>
                </div>
            </div>
        );
    }
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <div className={styles.errorIcon}>
                        <AlertCircle size={32} />
                    </div>
                    <h3 className={styles.errorTitle}>{t('errors.title')}</h3>
                    <p className={styles.errorMessage}>{error}</p>
                    <button 
                        onClick={() => window.location.reload()}
                        className={styles.retryButton}
                    >
                        {t('errors.retry')}
                    </button>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {/* Payment Summary */}
            <div className={styles.paymentSummary}>
                <h3 className={styles.summaryTitle}>
                    <CreditCard size={20} />
                    {t('summary.title')}
                </h3>
                {basketTotal && (
                    <div className={styles.summaryDetails}>
                        <div className={styles.summaryRow}>
                            <span className={styles.summaryLabel}>{t('summary.total')}</span>
                            <span className={styles.summaryAmount}>
                                €{basketTotal.displayAmount}
                            </span>
                        </div>
                        <div className={styles.summaryRow}>
                            <span className={styles.summaryLabel}>{t('summary.currency')}</span>
                            <span className={styles.summaryValue}>
                                {basketTotal.currency.toUpperCase()}
                            </span>
                        </div>
                    </div>
                )}
            </div>
            {/* Stripe Elements */}
            <div className={styles.paymentFormContainer}>
                <div className={styles.formHeader}>
                    <h3 className={styles.formTitle}>{t('form.title')}</h3>
                    <p className={styles.formDescription}>{t('form.description')}</p>
                </div>
                {clientSecret && (
                    <Elements options={options} stripe={stripePromise}>
                        <CheckoutForm
                            basketId={basketId}
                            userCustomerId={userCustomerId}
                            amount={basketTotal?.amount}
                            onPaymentSuccess={handlePaymentSuccess}
                            onPaymentError={handlePaymentError}
                        />
                    </Elements>
                )}
            </div>
            {/* Security Notice */}
            <div className={styles.securityNotice}>
                <Shield size={16} />
                <span>{t('security.notice')}</span>
            </div>
        </div>
    );
} 
"use client";
import React, { useState } from 'react';
import { PaymentElement, useStripe, useElements } from '@stripe/react-stripe-js';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { CreditCard, Loader, AlertCircle, CheckCircle } from '@/icons';
import styles from './CheckoutForm.module.css';
export default function CheckoutForm({ 
    basketId, 
    userCustomerId, 
    amount, 
    onPaymentSuccess, 
    onPaymentError 
}) {
    const stripe = useStripe();
    const elements = useElements();
    const router = useRouter();
    const t = useTranslations('CheckoutForm');
    const [isLoading, setIsLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [messageType, setMessageType] = useState(''); // 'error', 'success', 'info'
    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!stripe || !elements) {
            setMessage(t('errors.stripeNotLoaded'));
            setMessageType('error');
            return;
        }
        setIsLoading(true);
        setMessage('');
        setMessageType('');
        try {
            // Confirm payment with Stripe
            const { error, paymentIntent } = await stripe.confirmPayment({
                elements,
                confirmParams: {
                    return_url: `${window.location.origin}/payments/complete?basketId=${basketId}&userCustomerId=${userCustomerId}`,
                },
                redirect: 'if_required', // Only redirect if 3D Secure is required
            });
            if (error) {
                // Payment failed
                let errorMessage = t('errors.paymentFailed');
                if (error.type === 'card_error' || error.type === 'validation_error') {
                    errorMessage = error.message;
                } else if (error.type === 'authentication_error') {
                    errorMessage = t('errors.authenticationFailed');
                } else if (error.type === 'rate_limit_error') {
                    errorMessage = t('errors.rateLimitExceeded');
                } else if (error.type === 'api_connection_error') {
                    errorMessage = t('errors.connectionError');
                } else if (error.type === 'api_error') {
                    errorMessage = t('errors.serverError');
                }
                setMessage(errorMessage);
                setMessageType('error');
                if (onPaymentError) {
                    onPaymentError(error);
                }
            } else if (paymentIntent) {
                // Payment succeeded or is processing
                if (paymentIntent.status === 'succeeded') {
                    setMessage(t('success.paymentCompleted'));
                    setMessageType('success');
                    // Navigate to completion page
                    router.push(`/payments/complete?payment_intent_client_secret=${paymentIntent.client_secret}&basketId=${basketId}&userCustomerId=${userCustomerId}`);
                    if (onPaymentSuccess) {
                        onPaymentSuccess({
                            success: true,
                            paymentIntentId: paymentIntent.id,
                            status: paymentIntent.status
                        });
                    }
                } else if (paymentIntent.status === 'processing') {
                    setMessage(t('info.paymentProcessing'));
                    setMessageType('info');
                    // Navigate to completion page to handle processing state
                    router.push(`/payments/complete?payment_intent_client_secret=${paymentIntent.client_secret}&basketId=${basketId}&userCustomerId=${userCustomerId}`);
                } else if (paymentIntent.status === 'requires_action') {
                    setMessage(t('info.additionalActionRequired'));
                    setMessageType('info');
                    // Stripe will handle the additional action
                } else {
                    setMessage(t('errors.unexpectedStatus'));
                    setMessageType('error');
                }
            }
        } catch (err) {
            setMessage(t('errors.unexpectedError'));
            setMessageType('error');
            if (onPaymentError) {
                onPaymentError(err);
            }
        } finally {
            setIsLoading(false);
        }
    };
    const getMessageIcon = () => {
        switch (messageType) {
            case 'error':
                return <AlertCircle size={16} />;
            case 'success':
                return <CheckCircle size={16} />;
            case 'info':
                return <CreditCard size={16} />;
            default:
                return null;
        }
    };
    return (
        <form id="payment-form" onSubmit={handleSubmit} className={styles.form}>
            <div className={styles.paymentElementContainer}>
                <PaymentElement 
                    id="payment-element"
                    options={{
                        layout: 'tabs',
                        paymentMethodOrder: ['card', 'paypal', 'apple_pay', 'google_pay'],
                    }}
                />
            </div>
            {message && (
                <div className={`${styles.message} ${styles[messageType]}`}>
                    {getMessageIcon()}
                    <span>{message}</span>
                </div>
            )}
            <div className={styles.buttonContainer}>
                <button 
                    disabled={isLoading || !stripe || !elements}
                    type="submit"
                    className={styles.payButton}
                >
                    {isLoading ? (
                        <div className={styles.loadingContent}>
                            <Loader size={16} />
                            <span>{t('button.processing')}</span>
                        </div>
                    ) : (
                        <div className={styles.buttonContent}>
                            <CreditCard size={16} />
                            <span>
                                {amount 
                                    ? t('button.payAmount', { amount: (amount / 100).toFixed(2) })
                                    : t('button.payNow')
                                }
                            </span>
                        </div>
                    )}
                </button>
            </div>
            <div className={styles.paymentInfo}>
                <p className={styles.infoText}>
                    {t('info.securePayment')}
                </p>
                <p className={styles.infoText}>
                    {t('info.encryptedData')}
                </p>
            </div>
        </form>
    );
} 
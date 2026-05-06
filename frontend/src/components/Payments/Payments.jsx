// src/components/Payments/Payments.jsx
"use client";
import React, { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuth } from '../../context/AuthContext';
import { authorizePayment } from '../../api/client/paymentsApi';
import { getTotalBasketAmount } from '../../api/client/basketApi';
import { loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';
import { ArrowLeft, CreditCard, Loader, AlertCircle } from '@/icons';
import CheckoutForm from './CheckoutForm';
import PaymentHistory from './PaymentHistory';
import styles from './Payments.module.css';
const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || '{stripe_publishable_key}');
const Payments = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const t = useTranslations('Payments');
    const { user } = useAuth();
    const [clientSecret, setClientSecret] = useState('');
    const [amount, setAmount] = useState(0);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [basketId, setBasketId] = useState('');
    const [userCustomerId, setUserCustomerId] = useState('');
    const [activeTab, setActiveTab] = useState('payment'); // 'payment' or 'history'
    useEffect(() => {
        if (!user) {
            router.push('/login');
            return;
        }
        const initializePayment = async () => {
            try {
                // Get basket ID and customer ID from URL params
                const urlBasketId = searchParams.get('basketId');
                const urlUserCustomerId = searchParams.get('userCustomerId');
                if (!urlBasketId || !urlUserCustomerId) {
                    // If no payment params, just show payment history
                    setActiveTab('history');
                    setIsLoading(false);
                    return;
                }
                setBasketId(urlBasketId);
                setUserCustomerId(urlUserCustomerId);
                // Get total basket amount
                const totalAmount = await getTotalBasketAmount(urlBasketId, urlUserCustomerId);
                setAmount(totalAmount);
                // Authorize payment and get client secret
                const response = await authorizePayment(urlBasketId, urlUserCustomerId, totalAmount);
                setClientSecret(response.client_secret);
                setActiveTab('payment');
            } catch (err) {
                setError(err.message || 'Failed to initialize payment');
            } finally {
                setIsLoading(false);
            }
        };
        initializePayment();
    }, [user, router, searchParams]);
    const handlePaymentSuccess = (result) => {
        // Navigate to completion page
        router.push(`/payments/complete?payment_intent_client_secret=${result.paymentIntentId}&basketId=${basketId}&userCustomerId=${userCustomerId}`);
    };
    const handlePaymentError = (error) => {
        setError(error.message || 'Payment failed');
    };
    const handleBackToCart = () => {
        router.push('/cart');
    };
    const handleGoHome = () => {
        router.push('/');
    };
    if (isLoading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    <Loader size={48} />
                    <h2>Initializing Payment</h2>
                    <p>Please wait while we prepare your payment...</p>
                </div>
            </div>
        );
    }
    if (error && activeTab === 'payment') {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <AlertCircle size={48} />
                    <h2>Payment Error</h2>
                    <p>{error}</p>
                    <div className={styles.errorActions}>
                        <button onClick={handleBackToCart} className={styles.primaryButton}>
                            <ArrowLeft size={16} />
                            Back to Cart
                        </button>
                        <button onClick={handleGoHome} className={styles.secondaryButton}>
                            Go Home
                        </button>
                    </div>
                </div>
            </div>
        );
    }
    const appearance = {
        theme: 'stripe',
        variables: {
            colorPrimary: '#1a73e8',
            colorBackground: '#ffffff',
            colorText: '#202124',
            colorDanger: '#ea4335',
            fontFamily: 'system-ui, -apple-system, sans-serif',
            spacingUnit: '4px',
            borderRadius: '8px',
        },
    };
    const options = {
        clientSecret,
        appearance,
    };
    return (
        <div className={styles.container}>
            {/* Tab Navigation */}
            <div className={styles.tabNavigation}>
                <button 
                    className={`${styles.tabButton} ${activeTab === 'payment' ? styles.active : ''}`}
                    onClick={() => setActiveTab('payment')}
                    disabled={!clientSecret}
                >
                    <CreditCard size={16} />
                    Payment
                </button>
                <button 
                    className={`${styles.tabButton} ${activeTab === 'history' ? styles.active : ''}`}
                    onClick={() => setActiveTab('history')}
                >
                    History
                </button>
            </div>
            {/* Tab Content */}
            {activeTab === 'payment' && clientSecret && (
                <div className={styles.paymentContainer}>
                    <div className={styles.paymentHeader}>
                        <button onClick={handleBackToCart} className={styles.backButton}>
                            <ArrowLeft size={16} />
                            Back to Cart
                        </button>
                        <h1>Complete Your Payment</h1>
                    </div>
                    <div className={styles.paymentSummary}>
                        <h3>Order Summary</h3>
                        <div className={styles.summaryRow}>
                            <span>Total Amount:</span>
                            <span className={styles.amount}>€{(amount / 100).toFixed(2)}</span>
                        </div>
                    </div>
                    <Elements options={options} stripe={stripePromise}>
                        <CheckoutForm
                            basketId={basketId}
                            userCustomerId={userCustomerId}
                            amount={amount}
                            onPaymentSuccess={handlePaymentSuccess}
                            onPaymentError={handlePaymentError}
                        />
                    </Elements>
                </div>
            )}
            {activeTab === 'history' && (
                <div className={styles.historyContainer}>
                    <div className={styles.historyHeader}>
                        <h1>Payment History</h1>
                        <p>View your past transactions and payment details.</p>
                    </div>
                    <PaymentHistory userId={user?.id} />
                </div>
            )}
            {activeTab === 'payment' && !clientSecret && !error && (
                <div className={styles.noPaymentData}>
                    <CreditCard size={48} />
                    <h2>No Payment Data</h2>
                    <p>There's no payment information available. Start shopping to make a payment.</p>
                    <div className={styles.noPaymentActions}>
                        <button onClick={handleGoHome} className={styles.primaryButton}>
                            Start Shopping
                        </button>
                        <button onClick={() => setActiveTab('history')} className={styles.secondaryButton}>
                            View History
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};
export default Payments;
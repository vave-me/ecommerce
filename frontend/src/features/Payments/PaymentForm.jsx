// PaymentForm.jsx
"use client"
import React, {useState, useEffect, memo} from 'react';
import {Elements} from '@stripe/react-stripe-js';
import {loadStripe} from '@stripe/stripe-js';
import CheckoutForm from '../../components/Orders/CheckoutForm';
import {authorizePayment} from "../../api/client/paymentsApi";
import {useLocation} from "react-router-dom";
import {getBasket, getTotalBasketAmount} from "../../api/client/basketApi";
import styles from './PaymentForm.module.css';
// Use environment variable for Stripe key
const stripeKey = process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || process.env.STRIPE_TEST_KEY;
const stripePromise = loadStripe(stripeKey);
const PaymentForm = memo(function PaymentForm() {
    const location = useLocation();
    // read data from route state
    const {basketId, userCustomerId} = location.state || {};
    const [clientSecret, setClientSecret] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [basketData, setBasketData] = useState(null);
    const [appearance] = useState({
        theme: 'stripe',
        variables: {
            colorPrimary: '#0d6efd',
            colorBackground: '#ffffff',
            colorText: '#30313d',
            borderRadius: '8px',
            fontFamily: 'Arial, sans-serif',
        },
    });
    useEffect(() => {
        if (!basketId || !userCustomerId) {
            setError('Missing basketId or userCustomerId. Cannot create PaymentIntent.');
            setLoading(false);
            return;
        }
        async function fetchPaymentIntent() {
            try {
                setLoading(true);
                setError(null);
                // 1) Wait for total from the basket API
                const basketTotal = await getTotalBasketAmount(basketId);
                const basketInfo = await getBasket(basketId);
                setBasketData({
                    ...basketInfo,
                    total: basketTotal
                });
                const amountInMinorUnits = parseInt(basketTotal.amount, 10);
                const data = await authorizePayment(userCustomerId, amountInMinorUnits);
                if (data.clientSecret) {
                    setClientSecret(data.clientSecret);
                } else {
                    throw new Error('No clientSecret in response');
                }
            } catch (err) {
                setError(err.message || 'Failed to initialize payment');
            } finally {
                setLoading(false);
            }
        }
        fetchPaymentIntent();
    }, [basketId, userCustomerId]);
    const options = {
        clientSecret,
        appearance,
    };
    if (!stripeKey) {
        return (
            <div className={styles.container}>
                <div className={styles.errorMessage}>
                    Stripe configuration is missing. Please check your environment variables.
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>Complete Your Payment</h1>
                <p className={styles.subtitle}>Secure payment powered by Stripe</p>
            </div>
            {error && (
                <div className={styles.errorMessage}>
                    {error}
                </div>
            )}
            {basketData && (
                <div className={styles.paymentSummary}>
                    <h3 className={styles.summaryTitle}>Order Summary</h3>
                    <div className={styles.summaryRow}>
                        <span className={styles.summaryLabel}>Subtotal:</span>
                        <span className={styles.summaryValue}>
                            ${(basketData.total.amount / 100).toFixed(2)}
                        </span>
                    </div>
                    <div className={styles.summaryRow}>
                        <span className={styles.summaryLabel}>Tax:</span>
                        <span className={styles.summaryValue}>$0.00</span>
                    </div>
                    <div className={styles.summaryRow}>
                        <span className={styles.summaryLabel}>Total:</span>
                        <span className={`${styles.summaryValue} ${styles.totalAmount}`}>
                            ${(basketData.total.amount / 100).toFixed(2)}
                        </span>
                    </div>
                </div>
            )}
            {loading ? (
                <div className={styles.loadingContainer}>
                    <div className={styles.loadingSpinner}></div>
                    <p className={styles.loadingText}>Preparing your payment...</p>
                </div>
            ) : clientSecret ? (
                <div className={styles.paymentSection}>
                    <h3 className={styles.sectionTitle}>Payment Details</h3>
                    <div className={styles.stripeElements}>
                        <Elements options={options} stripe={stripePromise}>
                            <CheckoutForm basketId={basketId}/>
                        </Elements>
                    </div>
                    <div className={styles.securityBadge}>
                        Your payment information is secure and encrypted
                    </div>
                </div>
            ) : (
                <div className={styles.errorMessage}>
                    Unable to initialize payment. Please try again.
                </div>
            )}
        </div>
    );
});
export default PaymentForm;
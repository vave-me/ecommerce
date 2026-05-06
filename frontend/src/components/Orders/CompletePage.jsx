// src/features/Orders/components/CompletePage.jsx
"use client";
import React, { useEffect, useState, memo } from 'react';
// Removed: import styled from 'styled-components';
import { loadStripe } from '@stripe/stripe-js';
import { useRouter } from "next/navigation";
// Import the CSS module
import styles from './CompletePage.module.css';
/**
 * OPTIMIZED: Memoized for better order completion performance
 */
const CompletePage = memo(function CompletePage() {
    const navigate = useRouter();
    const [status, setStatus] = useState('loading');
    const [intentId, setIntentId] = useState(null);
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const paymentIntentClientSecret = urlParams.get('payment_intent_client_secret');
        const basketId = urlParams.get('basketId') || '';
        if (!paymentIntentClientSecret) {
            setStatus('no_client_secret');
            return;
        }
        const fetchPaymentIntent = async () => {
            try {
                const stripeKey = process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || process.env.STRIPE_TEST_KEY;
                const stripe = await loadStripe(stripeKey);
                const { paymentIntent } = await stripe.retrievePaymentIntent(paymentIntentClientSecret);
                if (paymentIntent) {
                    setIntentId(paymentIntent.id);
                    setStatus(paymentIntent.status);
                    // If succeeded => move to next step
                    if (paymentIntent.status === 'succeeded') {
                        navigate.push('/order-confirmation', {
                            state: {
                                basketId,
                                paymentIntentId: paymentIntent.id,
                            },
                        });
                    }
                } else {
                    setStatus('error');
                }
            } catch (error) {
                setStatus('error');
            }
        };
        fetchPaymentIntent();
    }, [navigate]);
    const getStatusContent = () => {
        switch (status) {
            case 'succeeded':
                return {
                    message: 'Payment succeeded! Redirecting...',
                    icon: '✓',
                    className: 'succeeded'
                };
            case 'processing':
                return {
                    message: 'Your payment is processing.',
                    icon: '⏳',
                    className: 'processing'
                };
            case 'requires_payment_method':
                return {
                    message: 'Payment failed. Please try again.',
                    icon: '✗',
                    className: 'requiresPaymentMethod'
                };
            case 'requires_action':
                return {
                    message: 'Additional authentication required.',
                    icon: '⚠',
                    className: 'requiresAction'
                };
            case 'no_client_secret':
                return {
                    message: 'No payment info found.',
                    icon: '✗',
                    className: 'noClientSecret'
                };
            case 'loading':
                return {
                    message: 'Loading...',
                    icon: '',
                    className: 'loading'
                };
            case 'error':
            default:
                return {
                    message: 'An error occurred fetching payment status.',
                    icon: '✗',
                    className: 'error'
                };
        }
    };
    const { message, icon, className } = getStatusContent();
    return (
        <div className={styles.container}>
            <div className={`${styles.statusIcon} ${styles[className]}`}>
                {icon}
            </div>
            <h2 className={styles.title}>{message}</h2>
            {intentId && (
                <div className={styles.intentDetails}>
                    <p className={styles.message}>
                        <strong>Payment Intent ID:</strong> 
                        <span className={styles.intentId}>{intentId}</span>
                    </p>
                    <p className={`${styles.statusText} ${styles[className]}`}>
                        Status: {status}
                    </p>
                </div>
            )}
            {status === 'requires_payment_method' && (
                <a href="/checkout" className={styles.actionButton}>
                    Try Again
                </a>
            )}
            {status === 'succeeded' && (
                <a href="/orders" className={styles.actionButton}>
                    View Orders
                </a>
            )}
        </div>
    );
});
export default CompletePage;

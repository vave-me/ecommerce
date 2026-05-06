// StripePaymentForm.jsx
import React, { memo, useState } from 'react';
import { useStripe, useElements, CardElement } from '@stripe/react-stripe-js';
import styles from './StripePaymentForm.module.css';
const StripePaymentForm = memo(({ product, onSuccess, onError }) => {
    const stripe = useStripe();
    const elements = useElements();
    const [isProcessing, setIsProcessing] = useState(false);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(false);
    const cardElementOptions = {
        style: {
            base: {
                fontSize: '16px',
                color: '#424770',
                '::placeholder': {
                    color: '#aab7c4',
                },
                iconColor: '#666ee8',
            },
            invalid: {
                color: '#9e2146',
            },
        },
        hidePostalCode: true
    };
    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setSuccess(false);
        if (!stripe || !elements) {
            setError('Stripe has not loaded yet. Please try again.');
            return;
        }
        const cardElement = elements.getElement(CardElement);
        if (!cardElement) {
            setError('Card element not found. Please refresh and try again.');
            return;
        }
        setIsProcessing(true);
        try {
            // Create payment intent on the server
            const response = await fetch('/api/create-payment-intent', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ 
                    amount: product.price,
                    currency: product.currency || 'usd',
                    productId: product.id
                }),
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const { clientSecret, error: serverError } = await response.json();
            if (serverError) {
                throw new Error(serverError);
            }
            if (!clientSecret) {
                throw new Error('No client secret received from server');
            }
            const result = await stripe.confirmCardPayment(clientSecret, {
                payment_method: {
                    card: cardElement,
                    billing_details: {
                        name: 'Customer Name', // Replace with actual customer's name
                    },
                },
            });
            if (result.error) {
                setError(result.error.message);
                if (onError) {
                    onError(result.error);
                }
            } else {
                if (result.paymentIntent.status === 'succeeded') {
                    setSuccess(true);
                    if (onSuccess) {
                        onSuccess(result.paymentIntent);
                    }
                    // Show success notification
                    if (typeof window !== 'undefined' && window.toast) {
                        window.toast.success('Payment succeeded!');
                    }
                }
            }
        } catch (err) {
            setError(err.message || 'An unexpected error occurred');
            if (onError) {
                onError(err);
            }
        } finally {
            setIsProcessing(false);
        }
    };
    const handleCardChange = (event) => {
        if (event.error) {
            setError(event.error.message);
        } else {
            setError(null);
        }
    };
    return (
        <div className={styles.paymentForm}>
            <div className={styles.formHeader}>
                <h2 className={styles.formTitle}>Payment Information</h2>
                <p className={styles.formSubtitle}>Enter your card details below</p>
            </div>
            {product && (
                <div className={styles.priceDisplay}>
                    Total: ${(product.price / 100).toFixed(2)}
                </div>
            )}
            <form onSubmit={handleSubmit}>
                <div className={styles.cardSection}>
                    <label className={styles.cardLabel} htmlFor="card-element">
                        Credit or Debit Card
                    </label>
                    <div className={`${styles.cardElementContainer} ${error ? styles.error : ''} ${success ? styles.complete : ''}`}>
                        <CardElement
                            id="card-element"
                            options={cardElementOptions}
                            onChange={handleCardChange}
                        />
                        <div className={styles.cardIcon}>💳</div>
                    </div>
                </div>
                {error && (
                    <div className={styles.errorMessage}>
                        {error}
                    </div>
                )}
                {success && (
                    <div className={styles.successMessage}>
                        Payment completed successfully!
                    </div>
                )}
                <button 
                    type="submit" 
                    disabled={!stripe || isProcessing || success}
                    className={`${styles.paymentButton} ${isProcessing ? styles.loading : ''}`}
                >
                    {isProcessing ? 'Processing...' : 
                     success ? 'Payment Complete' :
                     product ? `Pay $${(product.price / 100).toFixed(2)}` : 'Pay Now'}
                </button>
                <div className={styles.securityInfo}>
                    Your payment information is secure and encrypted
                </div>
                <div className={styles.poweredByStripe}>
                    Powered by <a href="https://stripe.com" target="_blank" rel="noopener noreferrer" className={styles.stripeLink}>Stripe</a>
                </div>
            </form>
        </div>
    );
});
StripePaymentForm.displayName = 'StripePaymentForm';
export default StripePaymentForm;

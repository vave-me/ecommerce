"use client";
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { toast } from 'react-toastify';
import { approveOrder } from '@/api/client/orderingApi';
import { checkoutBasket } from '@/api/client/basketApi';
import { CreditCard, Shield, Lock } from 'lucide-react';
import styles from './payment.module.css';

export default function PaymentPage({ params: { locale } }) {
    const t = useTranslations('PaymentPage');
    const router = useRouter();
    const [orderData, setOrderData] = useState(null);
    const [processing, setProcessing] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        // Load order data from sessionStorage
        const pendingOrder = sessionStorage.getItem('pendingOrder');
        if (!pendingOrder) {
            toast.error('No order found. Please start checkout again.');
            router.push(`/${locale}/basket`);
            return;
        }
        
        try {
            const data = JSON.parse(pendingOrder);
            setOrderData(data);
        } catch (err) {
            // Error: 'Invalid order data:', err...
            toast.error('Invalid order data');
            router.push(`/${locale}/basket`);
        }
    }, [router, locale]);

    const handlePayment = async (paymentMethod) => {
        if (!orderData) return;
        
        setProcessing(true);
        setError(null);
        
        try {
            // Simulate payment processing
            // In production, this would integrate with Stripe or another payment provider
            await new Promise(resolve => setTimeout(resolve, 2000));
            
            // After successful payment, approve the order
            const approveResponse = await approveOrder(orderData.orderId);
            
            if (approveResponse.success) {
                // Checkout the basket
                await checkoutBasket(
                    orderData.basketId,
                    orderData.userCustomerId,
                    'payment_intent_' + Date.now() // Mock payment intent ID
                );
                
                // Clear order data
                sessionStorage.removeItem('pendingOrder');
                
                // Redirect to success page
                router.push(`/${locale}/checkout/success?orderId=${orderData.orderId}`);
            } else {
                setError(approveResponse.userMessage || 'Failed to process payment');
            }
        } catch (err) {
            // Error: 'Payment error:', err...
            setError('Payment processing failed. Please try again.');
        } finally {
            setProcessing(false);
        }
    };

    if (!orderData) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingWrapper}>
                    <div className={styles.spinner}></div>
                    <p className={styles.loadingText}>Loading order information...</p>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>{t('title')}</h1>
            
            {error && (
                <div className={styles.errorAlert}>
                    <button 
                        className={styles.closeButton} 
                        onClick={() => setError(null)}
                        aria-label="Close"
                    >
                        ×
                    </button>
                    {error}
                </div>
            )}
            
            <div className={styles.card}>
                <h4 className={styles.cardTitle}>Order Summary</h4>
                <div className={styles.orderDetail}>
                    <span>Order ID:</span>
                    <span><code className={styles.orderId}>{orderData.orderId.slice(0, 8)}...</code></span>
                </div>
                <div className={styles.orderDetail}>
                    <span><strong>Total Amount:</strong></span>
                    <span><strong>${orderData.total.toFixed(2)}</strong></span>
                </div>
                
                <hr className={styles.divider} />
                
                <h5 className={styles.sectionTitle}>Shipping Address</h5>
                <p className={styles.addressLine}>{orderData.shippingAddress.firstName} {orderData.shippingAddress.lastName}</p>
                <p className={styles.addressLine}>{orderData.shippingAddress.address}</p>
                <p className={styles.addressLine}>{orderData.shippingAddress.city}, {orderData.shippingAddress.postalCode}</p>
                <p className={styles.addressLine}>{orderData.shippingAddress.country}</p>
            </div>
            
            <div className={styles.card}>
                <h4 className={styles.cardTitle}>
                    <CreditCard className={styles.icon} />
                    Payment Method
                </h4>
                
                <div className={styles.paymentButtons}>
                    {orderData.paymentMethod === 'card' && (
                        <button
                            className={`${styles.paymentButton} ${styles.primaryButton}`}
                            onClick={() => handlePayment('card')}
                            disabled={processing}
                        >
                            {processing ? (
                                <>
                                    <span className={styles.buttonSpinner}></span>
                                    Processing...
                                </>
                            ) : (
                                <>
                                    <CreditCard className={styles.buttonIcon} />
                                    Pay with Credit Card
                                </>
                            )}
                        </button>
                    )}
                    
                    {orderData.paymentMethod === 'paypal' && (
                        <button
                            className={`${styles.paymentButton} ${styles.warningButton}`}
                            onClick={() => handlePayment('paypal')}
                            disabled={processing}
                        >
                            {processing ? (
                                <>
                                    <span className={styles.buttonSpinner}></span>
                                    Processing...
                                </>
                            ) : (
                                'Pay with PayPal'
                            )}
                        </button>
                    )}
                    
                    {orderData.paymentMethod === 'bank' && (
                        <button
                            className={`${styles.paymentButton} ${styles.secondaryButton}`}
                            onClick={() => handlePayment('bank')}
                            disabled={processing}
                        >
                            {processing ? (
                                <>
                                    <span className={styles.buttonSpinner}></span>
                                    Processing...
                                </>
                            ) : (
                                'Pay with Bank Transfer'
                            )}
                        </button>
                    )}
                </div>
                
                <div className={styles.securityInfo}>
                    <Shield className={styles.securityIcon} size={16} />
                    <Lock className={styles.securityIcon} size={16} />
                    <small>Your payment information is secure and encrypted</small>
                </div>
            </div>
        </div>
    );
}

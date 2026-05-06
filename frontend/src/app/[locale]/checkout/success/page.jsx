"use client";
import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { CheckCircle } from 'lucide-react';
import styles from './success.module.css';

export default function CheckoutSuccessPage({ params: { locale } }) {
    const t = useTranslations('CheckoutPage');
    const searchParams = useSearchParams();
    const orderId = searchParams.get('orderId');
    const [orderNumber, setOrderNumber] = useState('');

    useEffect(() => {
        if (orderId) {
            setOrderNumber(orderId.slice(0, 8).toUpperCase());
        }
    }, [orderId]);

    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <div className={styles.iconContainer}>
                    <CheckCircle size={64} className={styles.icon} />
                </div>
                <h1 className={styles.title}>
                    {t('success.title', { defaultValue: 'Order Confirmed!' })}
                </h1>
                <p className={styles.message}>
                    {t('success.message', { 
                        defaultValue: 'Thank you for your order. We\'ve received your payment and will begin processing your order soon.' 
                    })}
                </p>
                
                {orderNumber && (
                    <div className={styles.orderInfo}>
                        <p className={styles.orderNumber}>
                            Order Number: <strong>{orderNumber}</strong>
                        </p>
                        <p className={styles.emailNote}>
                            A confirmation email has been sent to your registered email address.
                        </p>
                    </div>
                )}
                
                <div className={styles.buttonContainer}>
                    {orderId && (
                        <Link 
                            href={`/${locale}/account/orders/${orderId}`}
                            className={styles.primaryButton}
                        >
                            View Order Details
                        </Link>
                    )}
                    <Link 
                        href={`/${locale}/account/orders`}
                        className={styles.secondaryButton}
                    >
                        {t('success.viewOrders', { defaultValue: 'My Orders' })}
                    </Link>
                    <Link 
                        href={`/${locale}`}
                        className={styles.secondaryButton}
                    >
                        {t('success.continueShopping', { defaultValue: 'Continue Shopping' })}
                    </Link>
                </div>
            </div>
        </div>
    );
}

"use client";
import React, { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuth } from '@/context/AuthContext';
import { getCurrentBasket } from '@/api/client/basketApi';
import CheckoutForm from '@/components/Orders/CheckoutForm';
import styles from './checkout.module.css';

export default function CheckoutPage({ params: { locale } }) {
    const router = useRouter();
    const searchParams = useSearchParams();
    const t = useTranslations('CheckoutPage');
    const { user } = useAuth();
    const [basketId, setBasketId] = useState(searchParams.get('basketId'));
    const [loading, setLoading] = useState(!basketId);
    const [error, setError] = useState(null);

    useEffect(() => {
        // If no basketId in URL, get current basket
        const loadCurrentBasket = async () => {
            if (basketId) return;
            
            try {
                setLoading(true);
                const response = await getCurrentBasket(user?.userId);
                
                if (response && response.basket_id) {
                    setBasketId(response.basket_id);
                } else {
                    setError('No active basket found. Please add items to your basket first.');
                    setTimeout(() => {
                        router.push(`/${locale}/basket`);
                    }, 3000);
                }
            } catch (err) {
                // Error: 'Failed to load basket:', err...
                setError('Failed to load basket information');
            } finally {
                setLoading(false);
            }
        };
        
        loadCurrentBasket();
    }, [basketId, user?.userId, router, locale]);

    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingWrapper}>
                    <div className={styles.spinner}></div>
                    <p className={styles.loadingText}>{t('loadingCheckout', { defaultValue: 'Loading checkout...' })}</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorAlert}>{error}</div>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <div className={styles.checkoutWrapper}>
                <h1 className={styles.title}>{t('title', { defaultValue: 'Checkout' })}</h1>
                {basketId && <CheckoutForm basketId={basketId} locale={locale} />}
            </div>
        </div>
    );
}

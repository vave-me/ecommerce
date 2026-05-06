"use client";
import React, { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuth } from '../../../context/AuthContext';
import { CreditCard, ArrowLeft, Shield, CheckCircle, AlertCircle } from '@/icons';
import PaymentForm from '../../../components/Payments/PaymentForm';
import PaymentHistory from '../../../components/Payments/PaymentHistory';
import styles from './page.module.css';
export default function PaymentsPage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const t = useTranslations('PaymentsPage');
    const { user, isLoading } = useAuth();
    // Get payment data from URL params or route state
    const basketId = searchParams.get('basketId');
    const userCustomerId = searchParams.get('userCustomerId');
    const amount = searchParams.get('amount');
    const returnUrl = searchParams.get('returnUrl') || '/';
    const [activeTab, setActiveTab] = useState('payment');
    const [paymentData, setPaymentData] = useState(null);
    // Check if we have payment data from navigation state
    useEffect(() => {
        if (typeof window !== 'undefined' && window.history.state?.usr) {
            const state = window.history.state.usr;
            if (state.basketId && state.userCustomerId) {
                setPaymentData({
                    basketId: state.basketId,
                    userCustomerId: state.userCustomerId,
                    amount: state.amount
                });
                setActiveTab('payment');
            }
        } else if (basketId && userCustomerId) {
            setPaymentData({
                basketId,
                userCustomerId,
                amount
            });
            setActiveTab('payment');
        }
    }, [basketId, userCustomerId, amount]);
    // Redirect to login if not authenticated (only after auth check is complete)
    useEffect(() => {
        if (!isLoading && !user) {
            router.push('/login');
        }
    }, [user, isLoading, router]);
    const handleBackClick = () => {
        if (paymentData) {
            router.push('/cart');
        } else {
            router.push(returnUrl);
        }
    };
    const handlePaymentComplete = (paymentResult) => {
        if (paymentResult.success) {
            router.push(`/checkout/success?orderId=${paymentResult.orderId}`);
        }
    };
    const handleTabChange = (tab) => {
        setActiveTab(tab);
    };
    if (isLoading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    <div className={styles.stateIcon}>
                        <CreditCard size={48} />
                    </div>
                    <h2 className={styles.stateTitle}>{t('loading.title')}</h2>
                    <p className={styles.stateDescription}>{t('loading.description')}</p>
                </div>
            </div>
        );
    }
    if (!user) {
        return null; // Will redirect to login
    }
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div className={styles.titleSection}>
                    <button 
                        onClick={handleBackClick}
                        className={styles.backButton}
                        aria-label={t('navigation.back')}
                    >
                        <ArrowLeft size={20} />
                    </button>
                    <div>
                        <h1 className={styles.title}>
                            <CreditCard size={28} />
                            {paymentData ? t('title.checkout') : t('title.payments')}
                        </h1>
                        <p className={styles.subtitle}>
                            {paymentData ? t('subtitle.checkout') : t('subtitle.payments')}
                        </p>
                    </div>
                </div>
                {!paymentData && (
                    <div className={styles.tabNav}>
                        <button
                            onClick={() => handleTabChange('payment')}
                            className={`${styles.tabButton} ${activeTab === 'payment' ? styles.active : ''}`}
                        >
                            <CreditCard size={16} />
                            {t('tabs.payment')}
                        </button>
                        <button
                            onClick={() => handleTabChange('history')}
                            className={`${styles.tabButton} ${activeTab === 'history' ? styles.active : ''}`}
                        >
                            <CheckCircle size={16} />
                            {t('tabs.history')}
                        </button>
                    </div>
                )}
            </div>
            <div className={styles.content}>
                {activeTab === 'payment' && (
                    <div className={styles.paymentSection}>
                        {paymentData ? (
                            <PaymentForm
                                basketId={paymentData.basketId}
                                userCustomerId={paymentData.userCustomerId}
                                amount={paymentData.amount}
                                onPaymentComplete={handlePaymentComplete}
                            />
                        ) : (
                            <div className={styles.emptyState}>
                                <div className={styles.stateIcon}>
                                    <AlertCircle size={48} />
                                </div>
                                <h2 className={styles.stateTitle}>{t('empty.title')}</h2>
                                <p className={styles.stateDescription}>{t('empty.description')}</p>
                                <button 
                                    onClick={() => router.push('/cart')}
                                    className={styles.actionButton}
                                >
                                    {t('empty.goToCart')}
                                </button>
                            </div>
                        )}
                    </div>
                )}
                {activeTab === 'history' && (
                    <div className={styles.historySection}>
                        <PaymentHistory userId={user.id} />
                    </div>
                )}
            </div>
            <div className={styles.securityBanner}>
                <Shield size={16} />
                <span>{t('security.message')}</span>
            </div>
        </div>
    );
}

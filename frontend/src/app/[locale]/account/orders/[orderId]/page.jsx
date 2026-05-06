"use client";
import React from 'react';
import { useAuth } from '@/context/AuthContext';
import { redirect } from 'next/navigation';
import OrderDetail from '@/components/Orders/OrderDetail';

export default function OrderDetailPage({ params: { locale, orderId } }) {
    const { user, isLoading, authChecked } = useAuth();
    
    if (!isLoading && authChecked && !user) {
        redirect(`/${locale}/signin?callbackUrl=/${locale}/account/orders/${orderId}`);
    }

    if (isLoading || !authChecked) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', padding: '2rem' }}>
                <div>Loading...</div>
            </div>
        );
    }

    return <OrderDetail orderId={orderId} locale={locale} />;
}
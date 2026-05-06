"use client";
import React from 'react';
import { useAuth } from '@/context/AuthContext';
import { redirect } from 'next/navigation';
import OrdersList from '@/components/Orders/OrdersList';

export default function OrdersPage({ params: { locale } }) {
    const { user, isLoading, authChecked } = useAuth();
    
    if (!isLoading && authChecked && !user) {
        redirect(`/${locale}/signin?callbackUrl=/${locale}/account/orders`);
    }

    if (isLoading || !authChecked) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', padding: '2rem' }}>
                <div>Loading...</div>
            </div>
        );
    }

    return <OrdersList locale={locale} />;
}
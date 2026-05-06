// src/components/OrderConfirmation.jsx
"use client";
import React, { useEffect, memo } from 'react';
// Removed styled-components import
import { useLocation } from 'react-router-dom';
import { useAuth } from "../../context/AuthContext";
import { checkoutBasket } from "../../api/client/basketApi";
import { createOrder } from "../../api/orderingApi";
import { useRouter } from "next/navigation";
// Import the CSS module
import styles from './OrderConfirmation.module.css';
const OrderConfirmation = memo(function OrderConfirmation() {
    const location = useLocation();
    const navigate = useRouter();
    const { user } = useAuth();
    // read from route state
    const { basketId, paymentIntentId } = location.state || {};
    const userId = user?.userId;
    useEffect(() => {
        const finalizeOrder = async () => {
            try {
                if (!basketId || !paymentIntentId || !userId) {
                    navigate.push('/');
                    return;
                }
                // 1) Mark basket as checked out (closing the basket)
                await checkoutBasket(basketId, userId);
                // 2) Create an order if required
                const orderData = {
                    basketId,
                    userCustomerId: userId,
                    paymentId: paymentIntentId,
                };
                const newOrder = await createOrder(orderData);
                // 3) Go to order success
                navigate.push('/order-success', { state: { orderId: newOrder.id } });
            } catch (err) {
                navigate.push('/');
            }
        };
        finalizeOrder();
    }, [basketId, paymentIntentId, userId, navigate]);
    return (
        <div className={styles.container}>
            <h2 className={styles.title}>Finalizing your order...</h2>
        </div>
    );
});
export default OrderConfirmation;

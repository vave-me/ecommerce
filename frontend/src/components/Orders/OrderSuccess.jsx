// src/components/OrderSuccess.jsx
import React, { memo } from 'react';
// Removed styled-components import
import {useLocation} from 'react-router-dom';
import Link from 'next/link';
// Import CSS module
import styles from './OrderSuccess.module.css';
const OrderSuccess = memo(function OrderSuccess() {
    const location = useLocation();
    const {orderId} = location.state || {};
    return (
        <div className={styles.container}>
            <h2 className={styles.title}>Your order has been placed successfully!</h2>
            {orderId && <p className={styles.orderId}>Order ID: {orderId}</p>}
            <Link
                className={styles.homeButton}
                to="/"
                href="/"
            >
                Go to Home
            </Link>
        </div>
    );
});
export default OrderSuccess;

// src/components/OffersList.jsx
"use client";
import React, { memo } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useOffers } from '../../hooks/queries/useOffersQuery';
import styles from './OffersList.module.css';
/**
 * OffersList Component
 * Fetches and displays a list of offers from the API using React Query
 */
const OffersList = memo(() => {
    const { user } = useAuth();
    const { data: offers = [], isLoading: loading, error } = useOffers();
    if (loading) {
        return <div className={styles.loading}>Loading offers...</div>;
    }
    if (error) {
        return <div className={styles.error}>Error loading offers: {error.message}</div>;
    }
    if (offers.length === 0) {
        return <div className={styles.empty}>No offers available at the moment.</div>;
    }
    return (
        <div className={styles.offersList}>
            <h2>Available Offers</h2>
            <div className={styles.offersGrid}>
                {offers.map((offer) => (
                    <div key={offer.id} className={styles.offerCard}>
                        <h3>{offer.title}</h3>
                        <p>{offer.description}</p>
                        <div className={styles.offerDetails}>
                            <span className={styles.price}>{offer.price}</span>
                            <span className={styles.category}>{offer.category}</span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
});
OffersList.displayName = 'OffersList';
export default OffersList;
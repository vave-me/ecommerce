// src/components/OfferItem.jsx
"use client";
import React, { memo } from "react";
import styles from "./OfferItem.module.css";
const OfferItem = memo(({ offer, onAccept, onReject }) => {
    const { id, price, buyerId } = offer;
    return (
        <div className={styles.itemContainer}>
            <p>Offer ID: {id}</p>
            <p>Price: ${price}</p>
            <p>Buyer: {buyerId}</p>
            <div className={styles.actions}>
                <button onClick={() => onAccept(offer)}>Accept</button>
                <button onClick={() => onReject(offer)}>Reject</button>
            </div>
        </div>
    );
});
OfferItem.displayName = 'OfferItem';
export default OfferItem;
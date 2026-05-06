// src/components/SendOfferForm.jsx
"use client";
import React, { useState, memo } from "react";
import PropTypes from "prop-types";
import styles from "./SendOfferForm.module.css";
const SendOfferForm = memo(({ item, onSend }) => {
    const [amount, setAmount] = useState("");
    const [isSending, setIsSending] = useState(false);
    const handleSend = async (e) => {
        e.preventDefault();
        if (!amount || isNaN(amount) || amount <= 0) {
            alert("Please enter a valid amount.");
            return;
        }
        setIsSending(true);
        try {
            await onSend(item, parseFloat(amount));
            setAmount("");
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
            setIsSending(false);
        }
    };
    return (
        <form className={styles.formContainer} onSubmit={handleSend}>
            <span className={styles.itemDetails}>
                {item?.name} (ID: {item?.id})
            </span>
            <input
                className={styles.amountInput}
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="Enter offer amount"
                disabled={isSending}
            />
            <button
                className={styles.submitButton}
                type="submit"
                disabled={isSending}
            >
                {isSending ? "Sending..." : "Send Offer"}
            </button>
        </form>
    );
});
SendOfferForm.displayName = 'SendOfferForm';
SendOfferForm.propTypes = {
    item: PropTypes.shape({
        id: PropTypes.string.isRequired,
        name: PropTypes.string.isRequired,
        sellerId: PropTypes.string.isRequired,
    }).isRequired,
    onSend: PropTypes.func.isRequired,
};
export default SendOfferForm;
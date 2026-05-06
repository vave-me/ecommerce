"use client";
import React, { memo } from "react";
import styles from "./InfoText.module.css";
const InfoText = memo(function InfoText() {
    return (
        <div className={styles.infoTextContainer}>
            <p className={styles.infoText}>
                ai in one place. connect. share. sale.
            </p>
        </div>
    );
});
export default InfoText; 
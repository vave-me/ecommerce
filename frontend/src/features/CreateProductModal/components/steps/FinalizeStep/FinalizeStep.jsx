import React from 'react';
import { FinalizeStep as SharedFinalizeStep } from "../../../../shared/components/FinalizeStep/FinalizeStep";
import styles from '../../../ProductModal.module.css';
// Re-export with styles pre-configured for ProductModal
export function FinalizeStep(props) {
    return (
        <SharedFinalizeStep
            {...props}
            styles={styles}
            title="Publish Your Product"
            description="Review your product details and publish your listing."
            publishLabel="Publish Product"
            itemType="product"
            successConfig={{
                onViewListings: () => window.location.href = "/products/my-listings"
            }}
        />
    );
} 
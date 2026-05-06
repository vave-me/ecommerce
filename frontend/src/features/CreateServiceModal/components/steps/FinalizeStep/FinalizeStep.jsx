import React from 'react';
import { FinalizeStep as SharedFinalizeStep } from "../../../../shared/components/FinalizeStep/FinalizeStep";
import styles from '../../../ServiceModal.module.css';

// Re-export with styles pre-configured for ServiceModal
export function FinalizeStep(props) {
    return (
        <SharedFinalizeStep
            {...props}
            styles={styles}
            title="Publish Your Service"
            description="Review your service details and publish your listing."
            publishLabel="Publish Service"
            itemType="service"
            successConfig={{
                onViewListings: () => window.location.href = "/services/my-listings"
            }}
        />
    );
} 
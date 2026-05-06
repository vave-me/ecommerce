// src/features/CreateProductModal/components/steps/SuccessStep.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {Check} from "@/icons";
// Assuming shared styles or specific ones for ProductModal
import styles from '../../../ProductModal.module.css';
// This component is now rendered by FinalizeStep when isSuccess is true
export function SuccessStep({productName, thumbnail, onViewListings, onClose}) {
    return (
        <div className={styles.successContainer}>
            <div className={styles.successIcon}>
                <Check size={48} aria-hidden="true"/>
            </div>
            <h2 className={styles.successTitle}>Listing Published!</h2>
            <p className={styles.successMessage}>
                Your product listing has been successfully published.
            </p>
            {/* Optional: Show a small summary of the published item */}
            {(productName || thumbnail) && (
                <div className={styles.productSummary}>
                    {thumbnail && (
                        <div className={styles.thumbnailPreview}>
                            <img src={thumbnail} alt={productName || 'Product image'}/>
                        </div>
                    )}
                    {productName && (
                        <div className={styles.productDetails}
                             style={!thumbnail ? {textAlign: 'center', width: '100%'} : {}}>
                            <h3>{productName}</h3>
                            {/* You could add price/category here if passed down */}
                        </div>
                    )}
                </div>
            )}
            {/* Updated Next Steps relevant to products */}
            <div className={styles.nextStepOptions}>
                <h4>What's Next?</h4>
                <ul className={styles.nextStepsList}>
                    <li>Share your listing link on social media or with potential buyers.</li>
                    <li>Keep an eye out for messages or offers from interested buyers.</li>
                    <li>Prepare your item for potential shipping or pickup.</li>
                </ul>
            </div>
            {/* Actions: View Listings / Create Another (Close Modal) */}
            <div className={styles.successActions}>
                <button
                    className={styles.secondaryButton}
                    type="button"
                    onClick={onViewListings} // Use specific handler passed via props
                >
                    View My Listings
                </button>
                <button
                    className={styles.primaryButton}
                    type="button"
                    onClick={onClose} // Close the modal to allow creating another
                >
                    Create Another Listing
                </button>
            </div>
        </div>
    );
}
SuccessStep.propTypes = {
    productName: PropTypes.string, // Optional: For display
    thumbnail: PropTypes.string, // Optional: For display
    onViewListings: PropTypes.func.isRequired, // Function to navigate to user's listings
    onClose: PropTypes.func.isRequired // Function to close the modal
};

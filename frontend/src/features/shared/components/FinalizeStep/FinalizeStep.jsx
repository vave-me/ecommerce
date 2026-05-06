import React, { useState, useMemo, useCallback, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Check, Camera, MapPin } from "@/icons";
import { FormActions } from "../../../../common/components/FormActions";
import { SuccessStep } from "./SuccessStep";
/**
 * Shared FinalizeStep Component
 * Used across all creation modals for consistent finalization functionality
 */
export function FinalizeStep({
    itemData, // Contains basicInfo and media for preview
    onFinalize, // Callback to submit final data
    onBack,
    isLoading = false,
    errors = {},
    isUserLoggedIn = true,
    isSuccess = false, // Prop to indicate if publishing was successful
    onClose, // Prop to close the modal (used by SuccessStep)
    styles, // CSS modules passed from parent modal
    // Configurable content
    title = "Finalize Your Listing",
    description = "Review your listing details and finalize your submission.",
    publishLabel = "Publish Listing",
    itemType = "listing", // deal, job, service, product, etc.
    // Feature flags
    showLocation = true,
    showTerms = true,
    showPreview = true,
    requireTerms = true,
    // Success step configuration
    successConfig = {}
}) {
    // --- Local State for Step Fields ---
    const [termsAccepted, setTermsAccepted] = useState(false);
    const [locationError, setLocationError] = useState('');
    // --- Derived State & Validation ---
    const canProceed = useMemo(() => {
        if (requireTerms && showTerms && !termsAccepted) return false;
        return true; // Location is typically optional
    }, [requireTerms, showTerms, termsAccepted]);
    const isPrimaryDisabled = useMemo(() => 
        isLoading || !isUserLoggedIn || !canProceed, 
        [isLoading, isUserLoggedIn, canProceed]
    );
    const primaryLabel = useMemo(() => {
        if (isLoading) return "Publishing...";
        return publishLabel;
    }, [isLoading, publishLabel]);
    // --- Handlers ---
    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        if (!canProceed || isLoading) return;
        if (requireTerms && showTerms && !termsAccepted) {
            setLocationError('You must accept the terms to publish.');
            return;
        }
        // Pass finalization data up to the parent handler (location handled in first step)
        onFinalize({
            termsAccepted,
            status: 'active'
        });
    }, [canProceed, isLoading, onFinalize, termsAccepted, requireTerms, showTerms]);
    // --- Render Logic ---
    // If publishing was successful, render the SuccessStep
    if (isSuccess) {
        return (
            <SuccessStep
                itemData={itemData}
                itemType={itemType}
                onClose={onClose}
                styles={styles}
                {...successConfig}
            />
        );
    }
    // Format price based on item type
    const formatPrice = (price) => {
        if (!price) return 'Price not set';
        const numPrice = typeof price === 'string' ? parseFloat(price) : price;
        if (isNaN(numPrice)) return 'Price not set';
        switch (itemType) {
            case 'job':
                return `$${numPrice.toLocaleString()}`;
            default:
                return `€${numPrice.toFixed(2)}`;
        }
    };
    // Otherwise, render the Finalize form
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>{title}</h2>
            <p className={styles.formDescription}>
                {description}
            </p>
            {/* Display Finalize Error Message Passed From Parent */}
            {errors?.finalize && (
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {errors.finalize}
                </div>
            )}
            <form className={styles.form} onSubmit={handleSubmit}>
                {/* --- Compact Preview & Location Row --- */}
                <div className={styles.formRow}>
                    {/* Preview Section - Compact */}
                    {showPreview && (
                        <div className={styles.formGroup}>
                            <label className={styles.formLabel}>Quick Preview</label>
                            <div className={styles.compactPreview}>
                                {/* Thumbnail */}
                                {itemData?.thumbnail || (itemData?.images && itemData.images[0]) ? (
                                    <img
                                        src={itemData.thumbnail || itemData.images[0]}
                                        alt={itemData?.name || `${itemType} image`}
                                        className={styles.compactPreviewImage}
                                    />
                                ) : (
                                    <div className={styles.compactPlaceholder}>
                                        <Camera size={20} aria-hidden="true" />
                                    </div>
                                )}
                                {/* Details */}
                                <div className={styles.compactDetails}>
                                    <div className={styles.compactTitle}>{itemData?.name || 'N/A'}</div>
                                    <div className={styles.compactPrice}>
                                        {formatPrice(itemData?.basePrice || itemData?.price || itemData?.salary)}
                                    </div>
                                    <div className={styles.compactCategory}>
                                        {itemData?.categoryName || 'No Category'}
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}
                    {/* Location Section - REMOVED: Now handled in BasicInfoStep */}
                    {/* Location data is passed from the first step via itemData */}
                </div>
                {/* --- Terms Section - Compact --- */}
                {showTerms && (
                    <div className={styles.compactTerms}>
                        <label className={styles.checkboxLabel}>
                            <input
                                type="checkbox"
                                className={styles.checkbox}
                                checked={termsAccepted}
                                onChange={(e) => setTermsAccepted(e.target.checked)}
                                required={requireTerms}
                            />
                            <span className={styles.checkboxText}>
                                I agree to the{' '}
                                <a href="/terms" target="_blank" rel="noopener noreferrer" className={styles.termsLink}>
                                    Terms of Service
                                </a>{' '}
                                and{' '}
                                <a href="/privacy" target="_blank" rel="noopener noreferrer" className={styles.termsLink}>
                                    Privacy Policy
                                </a>
                                {requireTerms && <span className={styles.requiredMark}> *</span>}
                            </span>
                        </label>
                    </div>
                )}
                <FormActions
                    primaryLabel={primaryLabel}
                    primaryIcon="check"
                    secondaryLabel="Back"
                    onCancel={onBack}
                    onPrimaryAction={handleSubmit}
                    isSubmitting={isLoading}
                    isDisabled={isPrimaryDisabled}
                />
            </form>
        </div>
    );
}
FinalizeStep.propTypes = {
    itemData: PropTypes.object,
    onFinalize: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    isUserLoggedIn: PropTypes.bool,
    isSuccess: PropTypes.bool,
    onClose: PropTypes.func.isRequired,
    styles: PropTypes.object.isRequired,
    title: PropTypes.string,
    description: PropTypes.string,
    publishLabel: PropTypes.string,
    itemType: PropTypes.string,
    showLocation: PropTypes.bool,
    showTerms: PropTypes.bool,
    showPreview: PropTypes.bool,
    requireTerms: PropTypes.bool,
    successConfig: PropTypes.object
}; 
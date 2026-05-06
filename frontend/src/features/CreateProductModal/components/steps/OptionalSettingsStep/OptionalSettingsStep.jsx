// src/features/CreateProductModal/components/steps/OptionalSettingsStep/OptionalSettingsStep.jsx
import React, {useCallback, useMemo, useState} from 'react';
import PropTypes from 'prop-types';
import styles from '../../../ProductModal.module.css';
import {FormActions} from "../../../../../common/components/FormActions";
import {PriceInputInline} from "../../../../shared/components/PriceInput/PriceInputInline";
export function OptionalInfoStep({
    initialData,
    onComplete,
    onBack,
    isLoading,
    errors,
    isUserLoggedIn
}) {
    // --- Local State ---
    const [formData, setFormData] = useState({
        weight: initialData?.weight?.toString() || '',
        height: initialData?.height?.toString() || '',
        width: initialData?.width?.toString() || '',
        depth: initialData?.depth?.toString() || '',
        manageStocks: initialData?.manageStocks || false,
        stock: initialData?.stock?.toString() || '',
        shippingCost: initialData?.shippingCost?.toString() || '',
        middlemanService: initialData?.middlemanService || false,
        attributes: initialData?.attributes || []
    });
    const isPrimaryDisabled = useMemo(() => isLoading, [isLoading]);
    const primaryLabel = useMemo(() => {
        if (isLoading) return "Saving...";
        return "Continue to Finalize";
    }, [isLoading]);
    const updateFormData = useCallback((field, value) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    }, []);
    // Attribute management functions
    const addAttribute = useCallback(() => {
        setFormData(prev => ({
            ...prev,
            attributes: [...prev.attributes, { key: '', value: '' }]
        }));
    }, []);
    const updateAttribute = useCallback((index, field, value) => {
        setFormData(prev => ({
            ...prev,
            attributes: prev.attributes.map((attr, i) => 
                i === index ? { ...attr, [field]: value } : attr
            )
        }));
    }, []);
    const removeAttribute = useCallback((index) => {
        setFormData(prev => ({
            ...prev,
            attributes: prev.attributes.filter((_, i) => i !== index)
        }));
    }, []);
    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        onComplete(formData);
    }, [formData, onComplete]);
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>Additional Details</h2>
            <p className={styles.formDescription}>
                Add extra information about inventory, shipping, and product specifications. All fields are optional.
            </p>
            {/* Submit Error */}
            {errors?.submit && (
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {errors.submit}
                </div>
            )}
            <form className={styles.form} onSubmit={handleSubmit}>
                {/* Inventory & Shipping Row */}
                <div className={styles.formRow}>
                    {/* Inventory Section - Compact */}
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel}>Inventory Management</label>
                        <div className={styles.compactSection}>
                            <div className={styles.toggleGroup}>
                                <div className={styles.toggleSwitch}>
                                    <input
                                        type="checkbox"
                                        id="manage-stocks"
                                        checked={formData.manageStocks}
                                        onChange={(e) => updateFormData('manageStocks', e.target.checked)}
                                        className={styles.toggleInput}
                                    />
                                    <label htmlFor="manage-stocks" className={styles.toggleLabel}></label>
                                </div>
                                <span className={styles.toggleText}>Manage Stock</span>
                            </div>
                            {formData.manageStocks && (
                                <input
                                    className={styles.formInput}
                                    type="number"
                                    min="0"
                                    step="1"
                                    value={formData.stock}
                                    onChange={(e) => updateFormData('stock', e.target.value)}
                                    placeholder="Stock quantity"
                                    aria-invalid={!!errors?.stock}
                                />
                            )}
                            {errors?.stock && (
                                <div className={styles.fieldError}>{errors.stock}</div>
                            )}
                        </div>
                    </div>
                    {/* Shipping Section - Compact */}
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel}>Shipping Options</label>
                        <div className={styles.compactSection}>
                            <PriceInputInline
                                id="shipping-cost"
                                value={formData.shippingCost}
                                onChange={(value) => updateFormData('shippingCost', value)}
                                placeholder="Shipping cost"
                                styles={styles}
                                disabled={isLoading}
                            />
                            <div className={styles.toggleGroup}>
                                <div className={styles.toggleSwitch}>
                                    <input
                                        type="checkbox"
                                        id="middleman-service"
                                        checked={formData.middlemanService}
                                        onChange={(e) => updateFormData('middlemanService', e.target.checked)}
                                        className={styles.toggleInput}
                                    />
                                    <label htmlFor="middleman-service" className={styles.toggleLabel}></label>
                                </div>
                                <span className={styles.toggleText}>Secure Payment</span>
                            </div>
                        </div>
                        {errors?.shippingCost && (
                            <div className={styles.fieldError}>{errors.shippingCost}</div>
                        )}
                    </div>
                </div>
                {/* Dimensions Section - All in rows */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel}>Dimensions (Optional)</label>
                    {/* Weight - Full width */}
                    <div className={styles.dimensionRow}>
                        <input
                            className={`${styles.formInput} ${errors?.weight ? styles.inputError : ''}`}
                            type="number"
                            min="0"
                            value={formData.weight}
                            onChange={(e) => updateFormData('weight', e.target.value)}
                            placeholder="Weight in grams"
                            aria-invalid={!!errors?.weight}
                        />
                        {errors?.weight && (
                            <div className={styles.fieldError}>{errors.weight}</div>
                        )}
                    </div>
                    {/* Height & Width */}
                    <div className={styles.dimensionRow}>
                        <input
                            className={`${styles.formInput} ${errors?.height ? styles.inputError : ''}`}
                            type="number"
                            min="0"
                            step="0.1"
                            value={formData.height}
                            onChange={(e) => updateFormData('height', e.target.value)}
                            placeholder="Height in cm"
                            aria-invalid={!!errors?.height}
                        />
                        <input
                            className={`${styles.formInput} ${errors?.width ? styles.inputError : ''}`}
                            type="number"
                            min="0"
                            step="0.1"
                            value={formData.width}
                            onChange={(e) => updateFormData('width', e.target.value)}
                            placeholder="Width in cm"
                            aria-invalid={!!errors?.width}
                        />
                    </div>
                    {/* Show errors for height and width */}
                    {(errors?.height || errors?.width) && (
                        <div className={styles.dimensionRow}>
                            {errors?.height && (
                                <div className={styles.fieldError}>{errors.height}</div>
                            )}
                            {errors?.width && (
                                <div className={styles.fieldError}>{errors.width}</div>
                            )}
                        </div>
                    )}
                    {/* Depth - Full width */}
                    <div className={styles.dimensionRow}>
                        <input
                            className={`${styles.formInput} ${errors?.depth ? styles.inputError : ''}`}
                            type="number"
                            min="0"
                            step="0.1"
                            value={formData.depth}
                            onChange={(e) => updateFormData('depth', e.target.value)}
                            placeholder="Depth in cm"
                            aria-invalid={!!errors?.depth}
                        />
                        {errors?.depth && (
                            <div className={styles.fieldError}>{errors.depth}</div>
                        )}
                    </div>
                </div>
                {/* Attributes Section */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel}>Product Attributes (Optional)</label>
                    <div style={{fontSize: '12px', color: '#6b7280', marginBottom: '0.75rem'}}>
                        Add custom attributes like size, color, material, etc.
                    </div>
                    {formData.attributes.map((attribute, index) => (
                        <div key={index} className={styles.attributeRow}>
                            <input
                                className={styles.formInput}
                                value={attribute.key}
                                onChange={(e) => updateAttribute(index, 'key', e.target.value)}
                                placeholder="Attribute name (e.g., Color)"
                            />
                            <input
                                className={styles.formInput}
                                value={attribute.value}
                                onChange={(e) => updateAttribute(index, 'value', e.target.value)}
                                placeholder="Attribute value (e.g., Red)"
                            />
                            <button
                                type="button"
                                className={styles.removeButton}
                                onClick={() => removeAttribute(index)}
                                aria-label="Remove attribute"
                            >
                                ×
                            </button>
                        </div>
                    ))}
                    <button
                        type="button"
                        className={styles.addButton}
                        onClick={addAttribute}
                    >
                        + Add Attribute
                    </button>
                </div>
                {/* Form Actions */}
                <FormActions
                    primaryLabel={primaryLabel}
                    primaryIcon="arrow-right"
                    secondaryLabel="Back"
                    onCancel={onBack}
                    isSubmitting={isLoading}
                    isDisabled={isPrimaryDisabled}
                />
            </form>
        </div>
    );
}
// --- PropTypes ---
OptionalInfoStep.propTypes = {
    initialData: PropTypes.shape({
        manageStocks: PropTypes.bool,
        stock: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        shippingCost: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        weight: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        height: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        width: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        depth: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        middlemanService: PropTypes.bool,
        attributes: PropTypes.arrayOf(PropTypes.shape({
            key: PropTypes.string,
            value: PropTypes.string
        }))
    }),
    onComplete: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    isUserLoggedIn: PropTypes.bool
};

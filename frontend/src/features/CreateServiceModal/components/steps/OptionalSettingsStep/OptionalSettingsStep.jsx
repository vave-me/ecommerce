// src/features/CreateServiceModal/components/steps/OptionalSettingsStep/OptionalSettingsStep.jsx
import React, {useCallback, useMemo, useState} from 'react'; // Added useCallback
import PropTypes from 'prop-types';
// Assuming shared styles or specific ones for ServiceModal
import styles from '../../../ServiceModal.module.css';
import {FormActions} from "../../../../../common/components/FormActions";
import {PriceInputInline} from "../../../../shared/components/PriceInput/PriceInputInline";

// Import necessary components (adjust paths as needed)

export function OptionalInfoStep({
                                     initialData,
                                     onComplete,
                                     onBack,
                                     isLoading,
                                     errors, // Pass down errors object
                                     isUserLoggedIn // Pass down login status (though less critical for disabling button here)
                                 }) {
    // --- State Management ---
    const [manageStocks, setManageStocks] = useState(initialData?.manageStocks || false);
    const [stock, setStock] = useState(initialData?.stock?.toString() || "1");
    const [shippingCost, setShippingCost] = useState(initialData?.shippingCost?.toString() || "0");
    const [weight, setWeight] = useState(initialData?.weight?.toString() || "");
    const [height, setHeight] = useState(initialData?.height?.toString() || "");
    const [width, setWidth] = useState(initialData?.width?.toString() || "");
    const [depth, setDepth] = useState(initialData?.depth?.toString() || "");
    const [middlemanService, setMiddlemanService] = useState(initialData?.middlemanService || false);
    
    // New service-specific fields
    const [attributes, setAttributes] = useState(initialData?.attributes || [{ key: "", value: "" }]);
    const [options, setOptions] = useState(initialData?.options || [{ name: "", value: "", price: "" }]);

    // Disable the primary action button if user is not logged in or loading
    const isPrimaryDisabled = useMemo(() => isLoading || !isUserLoggedIn, [
        isLoading,
        isUserLoggedIn
    ]);

    // Get the primary action button label
    const primaryLabel = useMemo(() => {
        if (isLoading) return "Saving...";
        return "Continue to Finalize";
    }, [isLoading]);

    // Attribute handlers
    const handleAddAttribute = useCallback(() => {
        setAttributes(prev => [...prev, { key: "", value: "" }]);
    }, []);

    const handleAttributeChange = useCallback((index, field, value) => {
        setAttributes(prev => prev.map((attr, i) => 
            i === index ? { ...attr, [field]: value } : attr
        ));
    }, []);

    const handleRemoveAttribute = useCallback((index) => {
        setAttributes(prev => prev.filter((_, i) => i !== index));
    }, []);

    // Options handlers
    const handleAddOption = useCallback(() => {
        setOptions(prev => [...prev, { name: "", value: "", price: "" }]);
    }, []);

    const handleOptionChange = useCallback((index, field, value) => {
        setOptions(prev => prev.map((option, i) => 
            i === index ? { ...option, [field]: value } : option
        ));
    }, []);

    const handleRemoveOption = useCallback((index) => {
        setOptions(prev => prev.filter((_, i) => i !== index));
    }, []);

    // Handle form submission
    const handleSubmit = (e) => {
        e.preventDefault();

        // Build data object to pass to parent
        const optionalInfo = {
            manageStocks,
            stock,
            shippingCost,
            weight,
            height,
            width,
            depth,
            middlemanService,
            // Filter out empty attributes and options
            attributes: attributes.filter(attr => attr.key.trim() || attr.value.trim()),
            options: options.filter(option => option.name.trim() || option.value.trim() || option.price.trim())
        };

        // Call the parent's onComplete handler with form data
        onComplete(optionalInfo);
    };

    // --- Render Logic ---
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>Optional Settings</h2>
            <p className={styles.formDescription}>
                Additional details to help potential buyers find your deal.
            </p>

            {/* Display Submit Error Message Passed From Parent */}
            {errors?.submit && ( // Check if parent passes a generic submit error for this step
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {errors.submit}
                </div>
            )}

            <form className={styles.form} onSubmit={handleSubmit}>
                {/* Stock Management */}
                <div className={styles.formSection}>
                    <h3 className={styles.sectionTitle}>Inventory</h3>

                    <div className={styles.formGroup}>
                        <label className={styles.formLabel}>
                            Manage Stock
                        </label>
                        <div className={styles.toggleSwitch}>
                            <input
                                type="checkbox"
                                id="manage-stocks"
                                checked={manageStocks}
                                onChange={(e) => setManageStocks(e.target.checked)}
                                className={styles.toggleInput}
                            />
                            <label htmlFor="manage-stocks" className={styles.toggleLabel}></label>
                        </div>
                    </div>

                    {manageStocks && (
                        <div className={styles.formGroup}>
                            <label className={styles.formLabel} htmlFor="stock-quantity">
                                Stock Quantity
                            </label>
                            <input
                                id="stock-quantity"
                                className={styles.formInput}
                                type="number"
                                min="0"
                                value={stock}
                                onChange={(e) => setStock(e.target.value)}
                                placeholder="Available quantity"
                            />
                        </div>
                    )}
                </div>

                {/* Shipping */}
                <div className={styles.formSection}>
                    <h3 className={styles.sectionTitle}>Shipping</h3>

                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="shipping-cost">
                            Shipping Cost
                        </label>
                        <PriceInputInline
                            id="shipping-cost"
                            value={shippingCost}
                            onChange={setShippingCost}
                            styles={styles}
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label className={styles.formLabel}>
                            Middleman Service
                        </label>
                        <div className={styles.toggleSwitch}>
                            <input
                                type="checkbox"
                                id="middleman-service"
                                checked={middlemanService}
                                onChange={(e) => setMiddlemanService(e.target.checked)}
                                className={styles.toggleInput}
                            />
                            <label htmlFor="middleman-service" className={styles.toggleLabel}></label>
                        </div>
                        <div className={styles.inputHint}>
                            Enable if you want to offer secure transactions through our platform.
                        </div>
                    </div>
                </div>

                {/* Dimensions */}
                <div className={styles.formSection}>
                    <h3 className={styles.sectionTitle}>Dimensions</h3>

                    <div className={styles.formRow}>
                        <div className={styles.formGroup}>
                            <label className={styles.formLabel} htmlFor="product-weight">
                                Weight (g)
                            </label>
                            <input
                                id="product-weight"
                                className={styles.formInput}
                                type="number"
                                min="0"
                                value={weight}
                                onChange={(e) => setWeight(e.target.value)}
                                placeholder="Weight in grams"
                            />
                        </div>
                    </div>

                    <div className={styles.formRow}>
                        <div className={styles.formGroup}>
                            <label className={styles.formLabel} htmlFor="product-height">
                                Height (cm)
                            </label>
                            <input
                                id="product-height"
                                className={styles.formInput}
                                type="number"
                                min="0"
                                value={height}
                                onChange={(e) => setHeight(e.target.value)}
                                placeholder="Height in cm"
                            />
                        </div>

                        <div className={styles.formGroup}>
                            <label className={styles.formLabel} htmlFor="product-width">
                                Width (cm)
                            </label>
                            <input
                                id="product-width"
                                className={styles.formInput}
                                type="number"
                                min="0"
                                value={width}
                                onChange={(e) => setWidth(e.target.value)}
                                placeholder="Width in cm"
                            />
                        </div>
                    </div>

                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="product-depth">
                            Depth (cm)
                        </label>
                        <input
                            id="product-depth"
                            className={styles.formInput}
                            type="number"
                            min="0"
                            value={depth}
                            onChange={(e) => setDepth(e.target.value)}
                            placeholder="Depth in cm"
                        />
                    </div>
                </div>

                {/* Custom Attributes */}
                <div className={styles.formSection}>
                    <h3 className={styles.sectionTitle}>Custom Attributes</h3>
                    <p className={styles.inputHint}>
                        Add custom key-value pairs to describe specific features of your service.
                    </p>
                    
                    {attributes.map((attribute, index) => (
                        <div key={index} className={styles.attributeRow}>
                            <div className={styles.formGroup}>
                                <input
                                    type="text"
                                    className={styles.formInput}
                                    placeholder="Attribute name (e.g., Experience Level)"
                                    value={attribute.key}
                                    onChange={(e) => handleAttributeChange(index, 'key', e.target.value)}
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <input
                                    type="text"
                                    className={styles.formInput}
                                    placeholder="Attribute value (e.g., Expert)"
                                    value={attribute.value}
                                    onChange={(e) => handleAttributeChange(index, 'value', e.target.value)}
                                />
                            </div>
                            <button
                                type="button"
                                className={styles.clearLocationButton}
                                onClick={() => handleRemoveAttribute(index)}
                                disabled={attributes.length === 1}
                            >
                                Remove
                            </button>
                        </div>
                    ))}
                    
                    <button
                        type="button"
                        className={styles.secondaryButton}
                        onClick={handleAddAttribute}
                        style={{ marginTop: '0.5rem' }}
                    >
                        Add Attribute
                    </button>
                </div>

                {/* Service Options */}
                <div className={styles.formSection}>
                    <h3 className={styles.sectionTitle}>Service Options</h3>
                    <p className={styles.inputHint}>
                        Add different service packages or add-ons with their prices.
                    </p>
                    
                    {options.map((option, index) => (
                        <div key={index} className={styles.optionRow}>
                            <div className={styles.formGroup}>
                                <input
                                    type="text"
                                    className={styles.formInput}
                                    placeholder="Option name (e.g., Express Delivery)"
                                    value={option.name}
                                    onChange={(e) => handleOptionChange(index, 'name', e.target.value)}
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <input
                                    type="text"
                                    className={styles.formInput}
                                    placeholder="Description (e.g., 24-hour delivery)"
                                    value={option.value}
                                    onChange={(e) => handleOptionChange(index, 'value', e.target.value)}
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <div className={styles.priceInputWrapper}>
                                    <span className={styles.currencySymbol}>€</span>
                                    <input
                                        type="number"
                                        className={styles.priceInput}
                                        placeholder="0"
                                        min="0"
                                        step="0.01"
                                        value={option.price}
                                        onChange={(e) => handleOptionChange(index, 'price', e.target.value)}
                                    />
                                </div>
                            </div>
                            <button
                                type="button"
                                className={styles.clearLocationButton}
                                onClick={() => handleRemoveOption(index)}
                                disabled={options.length === 1}
                            >
                                Remove
                            </button>
                        </div>
                    ))}
                    
                    <button
                        type="button"
                        className={styles.secondaryButton}
                        onClick={handleAddOption}
                        style={{ marginTop: '0.5rem' }}
                    >
                        Add Option
                    </button>
                </div>

                {/* Form Actions */}
                <FormActions
                    onCancel={onBack}
                    isSubmitting={isLoading}
                    isDisabled={isPrimaryDisabled}
                    submitLabel={primaryLabel}
                    cancelLabel="Back"
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
        })),
        options: PropTypes.arrayOf(PropTypes.shape({
            name: PropTypes.string,
            value: PropTypes.string,
            price: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
        }))
    }),
    onComplete: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    isUserLoggedIn: PropTypes.bool
};

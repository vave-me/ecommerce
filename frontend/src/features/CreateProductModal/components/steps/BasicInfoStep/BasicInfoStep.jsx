// src/features/CreateProductModal/components/steps/BasicInfoStep/BasicInfoStep.jsx
import React, {useState, useMemo, useCallback} from 'react';
import PropTypes from 'prop-types';
import styles from '../../../ProductModal.module.css';
import {FormActions} from "../../../../../common/components/FormActions";
import {UserTypeSelect} from "../../../../shared/components/UserTypeSelect";
import {NameDescriptionInputs} from "../../../../shared/components/NameDescriptionInputs";
import {ConditionSelect, CONDITION_SETS} from "../../../../shared/components/ConditionSelect";
import {PriceInput} from "../../../../shared/components/PriceInput/PriceInput";
import {CategorySelector} from "../../../../shared/components/CategorySelector/CategorySelector";
import {LocationInput} from "../../../../shared/components/LocationInput";
export function BasicInfoStep({
                                  initialData,
                                  onComplete,
                                  onBack,
                                  isLoading,
                                  errors,
                                  isUserLoggedIn,
                                  categories
                              }) {
    // --- Local State ---
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        description: initialData?.description || '',
        productPrice: initialData?.productPrice || '',
        basePrice: initialData?.basePrice || '',
        condition: initialData?.condition || 'new',
        categoryId: initialData?.categoryId || '',
        categorySlug: initialData?.categorySlug || '',
        categoryName: initialData?.categoryName || '',
        brand: initialData?.brand || '',
        model: initialData?.model || '',
        sku: initialData?.sku || '',
        tags: Array.isArray(initialData?.tags)
            ? initialData.tags
            : (typeof initialData?.tags === 'string' && initialData.tags)
                ? initialData.tags.split(',').map(tag => tag.trim()).filter(Boolean)
                : [],
        negotiable: initialData?.negotiable || false,
        userType: initialData?.userType || 'private',
        hasVariants: initialData?.hasVariants || false,
        attributes: initialData?.attributes || [],
        // Location data
        address: initialData?.address || '',
        lat: initialData?.lat || null,
        lng: initialData?.lng || null
    });
    // Tags input state
    const [tagInput, setTagInput] = useState('');
    // --- Validation ---
    const canProceed = useMemo(() => {
        return (
            formData.name.trim() &&
            formData.description.trim() &&
            formData.productPrice.trim() &&
            formData.categoryId
        );
    }, [formData]);
    const isPrimaryDisabled = useMemo(() =>
            isLoading || !isUserLoggedIn || !canProceed,
        [isLoading, isUserLoggedIn, canProceed]
    );
    // --- Handlers ---
    const updateFormData = useCallback((field, value) => {
        setFormData(prev => ({...prev, [field]: value}));
    }, []);
    const handleCategorySelect = useCallback((category) => {
        if (category) {
            updateFormData('categoryId', category.id);
            updateFormData('categorySlug', category.slug || '');
            updateFormData('categoryName', category.name);
        }
    }, [updateFormData]);
    const handleCategoryClear = useCallback(() => {
        updateFormData('categoryId', '');
        updateFormData('categorySlug', '');
        updateFormData('categoryName', '');
    }, [updateFormData]);
    // Tags handling
    const handleTagInputKeyDown = useCallback((e) => {
        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault();
            const newTag = tagInput.trim();
            if (newTag && !formData.tags.includes(newTag)) {
                updateFormData('tags', [...formData.tags, newTag]);
                setTagInput('');
            }
        }
    }, [tagInput, formData.tags, updateFormData]);
    const removeTag = useCallback((tagToRemove) => {
        updateFormData('tags', formData.tags.filter(tag => tag !== tagToRemove));
    }, [formData.tags, updateFormData]);
    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        if (!canProceed) return;
        onComplete(formData);
    }, [canProceed, formData, onComplete]);
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>Basic Information</h2>
            <p className={styles.formDescription}>
                Provide essential details about your product. All fields marked with * are required.
            </p>
            {/* Submit Error */}
            {errors?.submit && (
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {errors.submit}
                </div>
            )}
            <form className={styles.form} onSubmit={handleSubmit}>
                {/* Product Name & Description */}
                <NameDescriptionInputs
                    nameId="product-name"
                    descriptionId="description"
                    nameLabel="Product Name"
                    descriptionLabel="Description"
                    nameValue={formData.name}
                    descriptionValue={formData.description}
                    onNameChange={(value) => updateFormData('name', value)}
                    onDescriptionChange={(value) => updateFormData('description', value)}
                    namePlaceholder="e.g. Apple iPhone 14 Pro - 256GB - Deep Purple"
                    descriptionPlaceholder="Describe item condition, features, included accessories, etc."
                    nameError={errors?.name}
                    descriptionError={errors?.description}
                    styles={styles}
                    nameMaxLength={150}
                    descriptionRows={4}
                />
                {/* Price Row - Product Price + Base Price */}
                <div className={styles.formRow}>
                    <PriceInput
                        id="product-price"
                        label="Product Price"
                        value={formData.productPrice}
                        onChange={(value) => updateFormData('productPrice', value)}
                        error={errors?.productPrice}
                        styles={styles}
                        required={true}
                    />
                    <PriceInput
                        id="base-price"
                        label="Regular Price (Optional)"
                        value={formData.basePrice}
                        onChange={(value) => updateFormData('basePrice', value)}
                        error={errors?.basePrice}
                        styles={styles}
                    />
                </div>
                {/* Category Selection - Full Width */}
                <CategorySelector
                    categories={categories}
                    selectedCategoryId={formData.categoryId}
                    selectedCategoryName={formData.categoryName}
                    onCategorySelect={handleCategorySelect}
                    onCategoryClear={handleCategoryClear}
                    categoryType="marketplace"
                    error={errors?.categoryId}
                    styles={styles}
                    placeholder="Select a category for your product"
                />
                {/* Condition + Brand Row */}
                <div className={styles.formRow}>
                    <ConditionSelect
                        id="condition"
                        label="Condition"
                        value={formData.condition}
                        onChange={(value) => updateFormData('condition', value)}
                        conditions={CONDITION_SETS.PRODUCT}
                        styles={styles}
                        error={errors?.condition}
                        required={true}
                    />
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="brand">
                            Brand
                        </label>
                        <input
                            id="brand"
                            className={`${styles.formInput} ${errors?.brand ? styles.inputError : ''}`}
                            value={formData.brand}
                            onChange={(e) => updateFormData('brand', e.target.value)}
                            placeholder="e.g. Apple, Samsung, Nike"
                            aria-invalid={!!errors?.brand}
                        />
                        {errors?.brand && (
                            <div className={styles.fieldError}>{errors.brand}</div>
                        )}
                    </div>
                </div>
                {/* Model + SKU Row */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="model">
                            Model
                        </label>
                        <input
                            id="model"
                            className={`${styles.formInput} ${errors?.model ? styles.inputError : ''}`}
                            value={formData.model}
                            onChange={(e) => updateFormData('model', e.target.value)}
                            placeholder="e.g. iPhone 14 Pro, Galaxy S23"
                            aria-invalid={!!errors?.model}
                        />
                        {errors?.model && (
                            <div className={styles.fieldError}>{errors.model}</div>
                        )}
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="sku">
                            SKU (Optional)
                        </label>
                        <input
                            id="sku"
                            className={`${styles.formInput} ${errors?.sku ? styles.inputError : ''}`}
                            value={formData.sku}
                            onChange={(e) => updateFormData('sku', e.target.value)}
                            placeholder="Product SKU/ID"
                            aria-invalid={!!errors?.sku}
                        />
                        {errors?.sku && (
                            <div className={styles.fieldError}>{errors.sku}</div>
                        )}
                    </div>
                </div>
                {/* Tags - Full Width */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel}>
                        Tags (Optional)
                    </label>
                    <div className={styles.tagsContainer}>
                        {formData.tags.map((tag, index) => (
                            <span key={index} className={styles.tag}>
                                {tag}
                                <button
                                    type="button"
                                    className={styles.tagRemove}
                                    onClick={() => removeTag(tag)}
                                    aria-label={`Remove tag ${tag}`}
                                >
                                    ×
                                </button>
                            </span>
                        ))}
                        <input
                            className={styles.formInput}
                            value={tagInput}
                            onChange={(e) => setTagInput(e.target.value)}
                            onKeyDown={handleTagInputKeyDown}
                            placeholder={formData.tags.length === 0 ? "Add tags (press Enter or comma to add)" : "Add more tags..."}
                        />
                    </div>
                    <div style={{fontSize: '12px', color: '#6b7280', marginTop: '4px'}}>
                        Press Enter or comma to add tags. Maximum 10 tags.
                    </div>
                </div>
                {/* Seller Type + Negotiable Row */}
                <div className={styles.formRow}>
                    <UserTypeSelect
                        id="product-seller-type"
                        value={formData.userType}
                        onChange={(value) => updateFormData('userType', value)}
                        styles={styles}
                        error={errors?.userType}
                    />
                    <div className={styles.formGroup}>
                        <div className={styles.toggleGroup}>
                            <label className={styles.formLabel}>Price Negotiable</label>
                            <div
                                className={`${styles.toggleSwitch} ${formData.negotiable ? styles.checked : ''}`}
                                onClick={() => updateFormData('negotiable', !formData.negotiable)}
                            >
                                <div className={styles.toggleSlider}></div>
                            </div>
                        </div>
                    </div>
                </div>
                {/* Variants Toggle - Full Width */}
                <div className={styles.formGroup}>
                    <div className={styles.toggleGroup}>
                        <div>
                            <label className={styles.formLabel}>Product has variants (sizes, colors, etc.)</label>
                            <div style={{fontSize: '12px', color: '#6b7280', marginTop: '2px'}}>
                                Enable if your product comes in different sizes, colors, or configurations
                            </div>
                        </div>
                        <div
                            className={`${styles.toggleSwitch} ${formData.hasVariants ? styles.checked : ''}`}
                            onClick={() => updateFormData('hasVariants', !formData.hasVariants)}
                        >
                            <div className={styles.toggleSlider}></div>
                        </div>
                    </div>
                </div>
                {/* Location Input */}
                <LocationInput
                    value={formData.address}
                    latitude={formData.lat}
                    longitude={formData.lng}
                    onLocationChange={(address) => updateFormData('address', address)}
                    onCoordinatesChange={(lat, lng) => {
                        updateFormData('lat', lat);
                        updateFormData('lng', lng);
                    }}
                    placeholder="Enter zipcode or city name"
                    label="Product Location (Optional)"
                    error={errors?.location}
                    showCurrentLocationButton={true}
                />
                {/* Form Actions */}
                <FormActions
                    onCancel={onBack}
                    isSubmitting={isLoading}
                    isDisabled={isPrimaryDisabled}
                    submitLabel={isLoading ? 'Saving...' : 'Next'}
                />
            </form>
        </div>
    );
}
BasicInfoStep.propTypes = {
    initialData: PropTypes.object,
    onComplete: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    isUserLoggedIn: PropTypes.bool,
    categories: PropTypes.array
};

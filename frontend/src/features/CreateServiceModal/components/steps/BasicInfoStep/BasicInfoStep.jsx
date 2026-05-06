// src/features/CreateServiceModal/components/steps/BasicInfoStep.jsx
import React, {useState, useMemo, useCallback} from 'react'; // Added useCallback
import PropTypes from 'prop-types';
import {Tag} from "@/icons"; // Assuming Tag icon for tags input
// Assuming shared styles or specific ones for ServiceModal
import styles from '../../../ServiceModal.module.css';
import {FormActions} from "../../../../../common/components/FormActions";
// Import necessary components (adjust paths as needed)

// Import shared components
import { UserTypeSelect } from "../../../../shared/components/UserTypeSelect";
import { NameDescriptionInputs } from "../../../../shared/components/NameDescriptionInputs";
import { PriceInput } from "../../../../shared/components/PriceInput/PriceInput";
import { CategorySelector } from "../../../../shared/components/CategorySelector/CategorySelector";
import { LocationInput } from "../../../../shared/components/LocationInput";

// Define constants locally or import them
const PRODUCT_CONDITIONS = ["new", "like-new", "excellent", "good", "fair"];
const DEAL_TYPES = ["discount", "coupon", "sale", "clearance", "rebate", "bundle", "other"];

export function BasicInfoStep({
                                  initialData,
                                  onComplete,
                                  onBack,
                                  isLoading,
                                  errors,
                                  isUserLoggedIn,
                                  categories
                              }) {
    // Debug categories prop

    // --- Local State for Step 1 Fields ---
    // Initialize state using optional chaining for safety
    const [dealName, setServiceName] = useState(initialData?.name || "");
    const [dealDescription, setServiceDescription] = useState(initialData?.description || "");
    // Initialize basePrice as string to match input type="number" behavior
    const [basePrice, setBasePrice] = useState(initialData?.basePrice?.toString() || "");
    const [dealPrice, setServicePrice] = useState(initialData?.dealPrice?.toString() || "");
    const [dealUrl, setServiceUrl] = useState(initialData?.dealUrl || "");
    const [dealDuration, setServiceDuration] = useState(initialData?.dealDuration?.toString() || "");
    const [dealType, setServiceType] = useState(initialData?.dealType || "discount");
    const [condition, setCondition] = useState(initialData?.condition || "new");
    const [categoryId, setCategoryId] = useState(initialData?.categoryId || "");
    const [categorySlug, setCategorySlug] = useState(initialData?.categorySlug || "");
    const [categoryName, setCategoryName] = useState(initialData?.categoryName || ""); // Store name for display
    // Ensure tags is initialized as a string for the input
    const [tags, setTags] = useState(Array.isArray(initialData?.tags) ? initialData.tags.join(", ") : initialData?.tags || "");
    const [brand, setBrand] = useState(initialData?.brand || "");
    const [model, setModel] = useState(initialData?.model || "");
    const [negotiable, setNegotiable] = useState(initialData?.negotiable || false);
    const [userType, setUserType] = useState(initialData?.userType || "private");
    const [sku, setSku] = useState(initialData?.sku || "");
    
    // New service-specific fields
    const [availability, setAvailability] = useState(initialData?.availability || "");
    const [providerName, setProviderName] = useState(initialData?.providerName || "");
    const [descriptionShort, setDescriptionShort] = useState(initialData?.descriptionShort || "");
    const [descriptionLong, setDescriptionLong] = useState(initialData?.descriptionLong || "");
    const [qualifications, setQualifications] = useState(Array.isArray(initialData?.qualifications) ? initialData.qualifications.join(", ") : initialData?.qualifications || "");
    const [contact, setContact] = useState(initialData?.contact || "");
    const [faq, setFaq] = useState(initialData?.faq || "");
    const [status, setStatus] = useState(initialData?.status || "draft");
    const [hasVariants, setHasVariants] = useState(initialData?.hasVariants || false);
    
    // Location state
    const [locationAddress, setLocationAddress] = useState(initialData?.address || "");
    const [latitude, setLatitude] = useState(initialData?.lat || null);
    const [longitude, setLongitude] = useState(initialData?.lng || null);

    // --- Derived State & Validation ---
    // Basic client-side check for enabling the submit button
    const canProceed = useMemo(() => {
        // Check required fields have non-empty, trimmed values
        return (
            dealName.trim() &&
            dealDescription.trim() &&
            basePrice.trim() && // Check price string isn't empty
            categoryId // Check if a category is selected
        );
    }, [dealName, dealDescription, basePrice, categoryId]);

    // Determine if the primary action button should be disabled
    const isPrimaryDisabled = useMemo(() => isLoading || !isUserLoggedIn || !canProceed, [
        isLoading,
        isUserLoggedIn,
        canProceed
    ]);

    // Determine the label for the primary action button
    const primaryLabel = useMemo(() => {
        if (isLoading) return "Saving...";
        return "Next";
    }, [isLoading]);

    // --- Handlers ---
    const handleCategorySelect = useCallback((category) => {
        if (category) {
            setCategoryId(category.id);
            setCategorySlug(category.slug || "");
            setCategoryName(category.name); // Update name when category is selected
        }
    }, []); // Empty dependency array as setters are stable

    const handleCategoryClear = useCallback(() => {
        setCategoryId("");
        setCategorySlug("");
        setCategoryName("");
    }, []); // Empty dependency array

    const handleSubmit = (e) => {
        e.preventDefault();
        // Basic check before submitting, although parent handles main validation/logic
        if (!canProceed) {
            // Please fill all required fields
            // Optionally set a local error state here if needed
            // setLocalError("Please fill all required fields.");
            return;
        }
        // Pass all local state data up to the parent handler
        onComplete({
            name: dealName,
            description: dealDescription,
            basePrice: basePrice, // Pass as string, parent handles conversion
            serviceType: dealType, // Map dealType to serviceType for services
            categoryId,
            categorySlug,
            categoryName, // Pass name along too
            tags, // Pass as comma-separated string, parent handles split
            negotiable,
            userType,
            // New service-specific fields
            availability,
            providerName,
            descriptionShort,
            descriptionLong,
            qualifications, // Pass as comma-separated string, parent handles split
            contact,
            faq,
            status,
            hasVariants,
            // Location data
            address: locationAddress,
            lat: latitude,
            lng: longitude
        });
    };

    // --- Render Logic ---
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>Basic Information</h2>
            <p className={styles.formDescription}>
                Provide essential details about your item. Fields marked with * are required.
            </p>

            {/* Display Submit Error Message Passed From Parent */}
            {errors?.submit && (
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {errors.submit}
                </div>
            )}

            <form className={styles.form} onSubmit={handleSubmit}>
                {/* Service Name & Description */}
                <NameDescriptionInputs
                    nameId="service-name"
                    descriptionId="service-description"
                    nameLabel="Service Name"
                    descriptionLabel="Description"
                    nameValue={dealName}
                    descriptionValue={dealDescription}
                    onNameChange={setServiceName}
                    onDescriptionChange={setServiceDescription}
                    namePlaceholder="e.g. Professional Website Design Service"
                    descriptionPlaceholder="Describe your service, what's included, delivery time, etc."
                    nameError={errors?.name}
                    descriptionError={errors?.description}
                    styles={styles}
                    nameMaxLength={150}
                    descriptionRows={6}
                />

                {/* Price & Negotiable Row */}
                <div className={styles.formRow}>
                    <PriceInput
                        id="service-price"
                        label="Service Price"
                        value={basePrice}
                        onChange={setBasePrice}
                        error={errors?.basePrice}
                        styles={styles}
                        required={true}
                    />

                    <PriceInput
                        id="service-discount-price"
                        label="Discounted Price (Optional)"
                        value={dealPrice}
                        onChange={setServicePrice}
                        styles={styles}
                    />
                </div>

                {/* Service Type & Provider Name */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="service-type">
                            Service Type <span className={styles.requiredMark}>*</span>
                        </label>
                        <select
                            id="service-type"
                            className={styles.formSelect}
                            value={dealType}
                            onChange={(e) => setServiceType(e.target.value)}
                        >
                            <option value="consultation">Consultation</option>
                            <option value="design">Design</option>
                            <option value="development">Development</option>
                            <option value="maintenance">Maintenance</option>
                            <option value="marketing">Marketing</option>
                            <option value="coaching">Coaching</option>
                            <option value="repair">Repair</option>
                            <option value="installation">Installation</option>
                            <option value="cleaning">Cleaning</option>
                            <option value="other">Other</option>
                        </select>
                    </div>
                    
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="provider-name">
                            Provider Name
                        </label>
                        <input
                            id="provider-name"
                            className={styles.formInput}
                            type="text"
                            value={providerName}
                            onChange={(e) => setProviderName(e.target.value)}
                            placeholder="Your business or professional name"
                        />
                    </div>
                </div>

                {/* Availability & Short Description */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="availability">
                            Availability
                        </label>
                        <input
                            id="availability"
                            className={styles.formInput}
                            type="text"
                            value={availability}
                            onChange={(e) => setAvailability(e.target.value)}
                            placeholder="e.g. Monday-Friday 9AM-5PM, Weekends"
                        />
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="description-short">
                            Short Description
                        </label>
                        <input
                            id="description-short"
                            className={styles.formInput}
                            type="text"
                            value={descriptionShort}
                            onChange={(e) => setDescriptionShort(e.target.value)}
                            placeholder="Brief summary of your service"
                            maxLength={100}
                        />
                    </div>
                </div>

                {/* Contact & Status */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="contact">
                            Contact Information
                        </label>
                        <input
                            id="contact"
                            className={styles.formInput}
                            type="text"
                            value={contact}
                            onChange={(e) => setContact(e.target.value)}
                            placeholder="Email, phone, or preferred contact method"
                        />
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="status">
                            Status <span className={styles.requiredMark}>*</span>
                        </label>
                        <select
                            id="status"
                            className={styles.formSelect}
                            value={status}
                            onChange={(e) => setStatus(e.target.value)}
                        >
                            <option value="draft">Draft</option>
                            <option value="active">Active</option>
                        </select>
                    </div>
                </div>

                {/* Qualifications */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="qualifications">
                        Qualifications (comma separated)
                    </label>
                    <input
                        id="qualifications"
                        className={styles.formInput}
                        type="text"
                        value={qualifications}
                        onChange={(e) => setQualifications(e.target.value)}
                        placeholder="e.g. Certified Web Developer, 5 years experience, Adobe Certified"
                    />
                </div>

                {/* Long Description */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="description-long">
                        Detailed Description
                    </label>
                    <textarea
                        id="description-long"
                        className={styles.formTextarea}
                        value={descriptionLong}
                        onChange={(e) => setDescriptionLong(e.target.value)}
                        placeholder="Provide a detailed description of your service, what's included, process, deliverables, etc."
                        rows={4}
                    />
                </div>

                {/* FAQ */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="faq">
                        Frequently Asked Questions
                    </label>
                    <textarea
                        id="faq"
                        className={styles.formTextarea}
                        value={faq}
                        onChange={(e) => setFaq(e.target.value)}
                        placeholder="Common questions and answers about your service"
                        rows={3}
                    />
                </div>

                {/* Toggle Switches */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup} style={{alignItems: 'flex-start'}}>
                        <div style={{display: 'flex', alignItems: 'center', gap: '0.75rem'}}>
                            <div className={styles.toggleSwitch}>
                                <input
                                    type="checkbox"
                                    id="negotiable-toggle"
                                    checked={negotiable}
                                    onChange={(e) => setNegotiable(e.target.checked)}
                                    className={styles.toggleInput}
                                />
                                <label htmlFor="negotiable-toggle" className={styles.toggleLabel}></label>
                            </div>
                            <label className={styles.formLabel} htmlFor="negotiable-toggle">
                                Price is negotiable
                            </label>
                        </div>
                    </div>
                    
                    <div className={styles.formGroup} style={{alignItems: 'flex-start'}}>
                        <div style={{display: 'flex', alignItems: 'center', gap: '0.75rem'}}>
                            <div className={styles.toggleSwitch}>
                                <input
                                    type="checkbox"
                                    id="has-variants-toggle"
                                    checked={hasVariants}
                                    onChange={(e) => setHasVariants(e.target.checked)}
                                    className={styles.toggleInput}
                                />
                                <label htmlFor="has-variants-toggle" className={styles.toggleLabel}></label>
                            </div>
                            <label className={styles.formLabel} htmlFor="has-variants-toggle">
                                Service has variants/options
                            </label>
                        </div>
                    </div>
                </div>

                {/* Experience Level & Provider Type Row */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="experience-level">
                            Experience Level <span className={styles.requiredMark}>*</span>
                        </label>
                        <select
                            id="experience-level"
                            className={`${styles.formSelect} ${errors?.condition ? styles.inputError : ""}`}
                            value={condition}
                            onChange={(e) => setCondition(e.target.value)}
                            aria-invalid={!!errors?.condition}
                            aria-describedby={errors?.condition ? "condition-error" : undefined}
                        >
                            <option value="beginner">Beginner (0-2 years)</option>
                            <option value="intermediate">Intermediate (2-5 years)</option>
                            <option value="experienced">Experienced (5+ years)</option>
                            <option value="expert">Expert (10+ years)</option>
                        </select>
                        {errors?.condition && (
                            <div className={styles.fieldError} id="condition-error" role="alert">
                                {errors.condition}
                            </div>
                        )}
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="provider-type">Provider Type</label>
                        <select
                            id="provider-type"
                            className={styles.formSelect}
                            value={userType}
                            onChange={(e) => setUserType(e.target.value)}
                        >
                            <option value="freelancer">Freelancer</option>
                            <option value="agency">Agency</option>
                            <option value="consultant">Consultant</option>
                            <option value="company">Company</option>
                        </select>
                    </div>
                </div>

                {/* Category Selection */}
                <CategorySelector
                    categories={categories}
                    selectedCategoryId={categoryId}
                    selectedCategoryName={categoryName}
                    onCategorySelect={handleCategorySelect}
                    onCategoryClear={handleCategoryClear}
                    categoryType="services"
                    error={errors?.categoryId}
                    styles={styles}
                    placeholder="Select a category for your service"
                />

                {/* Duration & Location */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="service-duration">Delivery Time</label>
                        <input
                            id="service-duration"
                            className={styles.formInput}
                            type="text"
                            value={brand}
                            onChange={(e) => setBrand(e.target.value)}
                            placeholder="e.g. 1-2 weeks, 3 days, 1 month"
                        />
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="service-location">Service Location</label>
                        <input
                            id="service-location"
                            className={styles.formInput}
                            type="text"
                            value={model}
                            onChange={(e) => setModel(e.target.value)}
                            placeholder="Remote, On-site, Hybrid"
                        />
                    </div>
                </div>

                {/* Tags & Skills */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="service-tags">
                            <Tag size={16} className={styles.labelIcon}/> Tags (comma separated)
                        </label>
                        <input
                            id="service-tags"
                            className={styles.formInput}
                            type="text"
                            value={tags}
                            onChange={(e) => setTags(e.target.value)}
                            placeholder="e.g. wordpress, react, seo, branding"
                        />
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="service-skills">Skills/Tools</label>
                        <input
                            id="service-skills"
                            className={styles.formInput}
                            type="text"
                            value={sku}
                            onChange={(e) => setSku(e.target.value)}
                            placeholder="e.g. Photoshop, Figma, HTML/CSS"
                        />
                    </div>
                </div>

                {/* Location Input */}
                <LocationInput
                    value={locationAddress}
                    latitude={latitude}
                    longitude={longitude}
                    onLocationChange={setLocationAddress}
                    onCoordinatesChange={(lat, lng) => {
                        setLatitude(lat);
                        setLongitude(lng);
                    }}
                    placeholder="Enter zipcode or city name"
                    label="Service Location (Optional)"
                    error={errors?.location}
                    showCurrentLocationButton={true}
                />

                {/* Form Actions */}
                <FormActions
                    onCancel={onBack}
                    isSubmitting={isLoading}
                    isDisabled={isPrimaryDisabled}
                    primaryLabel={primaryLabel}
                />
            </form>
        </div>
    );
}

// PropTypes
BasicInfoStep.propTypes = {
    initialData: PropTypes.shape({
        name: PropTypes.string,
        description: PropTypes.string,
        basePrice: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        serviceType: PropTypes.string,
        categoryId: PropTypes.string,
        categorySlug: PropTypes.string,
        categoryName: PropTypes.string,
        tags: PropTypes.oneOfType([
            PropTypes.string,
            PropTypes.arrayOf(PropTypes.string)
        ]),
        negotiable: PropTypes.bool,
        userType: PropTypes.string,
        // New service-specific fields
        availability: PropTypes.string,
        providerName: PropTypes.string,
        descriptionShort: PropTypes.string,
        descriptionLong: PropTypes.string,
        qualifications: PropTypes.oneOfType([
            PropTypes.string,
            PropTypes.arrayOf(PropTypes.string)
        ]),
        contact: PropTypes.string,
        faq: PropTypes.string,
        status: PropTypes.string,
        hasVariants: PropTypes.bool,
        // Location fields
        address: PropTypes.string,
        lat: PropTypes.number,
        lng: PropTypes.number,
    }),
    onComplete: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    isUserLoggedIn: PropTypes.bool,
    categories: PropTypes.array,
};

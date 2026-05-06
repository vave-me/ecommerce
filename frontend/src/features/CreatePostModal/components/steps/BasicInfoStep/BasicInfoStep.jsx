// CreatePostModal/components/steps/BasicInfoStep/BasicInfoStep.jsx
import React, {useMemo, useState} from 'react';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl'; //  Import hook
import {Tag} from "@/icons";
import styles from '../../../CreatePostModal.module.css'; // Use shared styles
import RichTextEditor from "../../../../TextEditor"; // Ensure path is correct
import {FormActions} from "../../../../../common/components/FormActions"; // Ensure path is correct - Assume this handles its own internal translations like "Cancel"
import {LocationInput} from "../../../../shared/components/LocationInput";
export function BasicInfoStep({
                                  initialData,
                                  onSubmit,
                                  onCancel,
                                  isLoading,
                                  errors, // Expects errors object with translation keys, e.g., { name: 'errorTitleRequired', submit: 'errorLoginRequired' }
                                  isUserLoggedIn,
                                  categories
                              }) {
    const t = useTranslations('CreatePostModal'); //  Use shared namespace
    // Local state for form fields (remains the same)
    const [postName, setPostName] = useState(initialData?.name || "");
    const [postDescription, setPostDescription] = useState(initialData?.description || "");
    const [tags, setTags] = useState(initialData?.tags || "");
    // New API fields
    const [typeOfPost, setTypeOfPost] = useState(initialData?.typeOfPost || "general");
    const [userType, setUserType] = useState(initialData?.userType || "private");
    const [categoryId, setCategoryId] = useState(initialData?.categoryId || "");
    const [status, setStatus] = useState(initialData?.status || "active");
    // Location state
    const [locationAddress, setLocationAddress] = useState(initialData?.address || "");
    const [latitude, setLatitude] = useState(initialData?.lat || null);
    const [longitude, setLongitude] = useState(initialData?.lng || null);
    // Memoize derived state for button logic
    const isPrimaryDisabled = useMemo(() => isLoading || !isUserLoggedIn, [isLoading, isUserLoggedIn]);
    // Determine primary button label using translations
    const primaryLabel = useMemo(() => {
        if (isLoading) return "Saving...";
        return "Next";
    }, [isLoading]);
    // Handle form submission (passes data up, parent validates/submits)
    const handleSubmit = (e) => {
        e.preventDefault();
        // Parent component (CreatePostModal) handles actual validation and submission logic
        onSubmit({
            name: postName,
            description: postDescription,
            tags,
            // New API fields
            typeOfPost,
            userType,
            categoryId,
            status,
            // Location data
            address: locationAddress,
            lat: latitude,
            lng: longitude
        });
    };
    // Helper to translate error keys, potentially with values
    const translateError = (errorKey) => {
        if (!errorKey) return null;
        // Define known interpolation values based on keys (example)
        const values = {};
        if (errorKey === 'errorTitleMinLength') values.minLength = 5;
        if (errorKey === 'errorContentMinLength') values.minLength = 20;
        return t(errorKey, values);
    };
    return (
        <div className={styles.formContainer}>
            {/*   Use translation */}
            <h2 className={styles.formTitle}>{t('basicInfoTitle')}</h2>
            {/*   Use translation */}
            <p className={styles.formDescription}>{t('basicInfoDescription')}</p>
            {/* --- Display Translated Submit Error Message --- */}
            {errors?.submit && (
                <div className={`${styles.fieldError} ${styles.submitError}`} role="alert">
                    {translateError(errors.submit)} {/*  Translate error key */}
                </div>
            )}
            {/* ------------------------------------- */}
            <form className={styles.form} onSubmit={handleSubmit} noValidate>
                {/* Title field */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="post-name">
                        {/*   Use translation */}
                        {t('basicInfoTitleLabel')} <span className={styles.requiredMark}>*</span>
                    </label>
                    <input
                        id="post-name"
                        className={`${styles.formInput} ${errors?.name ? styles.inputError : ""}`}
                        value={postName}
                        onChange={(e) => setPostName(e.target.value)}
                        //   Use translation
                        placeholder={t('basicInfoTitlePlaceholder')}
                        maxLength={100}
                        aria-invalid={!!errors?.name}
                        aria-describedby={errors?.name ? "name-error" : undefined}
                    />
                    {/* Display translated name error */}
                    {errors?.name && (
                        <div className={styles.fieldError} id="name-error" role="alert">
                            {translateError(errors.name)} {/*  Translate error key */}
                        </div>
                    )}
                    {/*   Use translation with interpolation */}
                    <div className={styles.charCount} aria-live="polite">
                        {t('basicInfoCharCount', {count: postName.length})}
                    </div>
                </div>
                {/* Content/Description field */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="post-description">
                        {/*   Use translation */}
                        {t('basicInfoContentLabel')} <span className={styles.requiredMark}>*</span>
                    </label>
                    <div
                        className={`${styles.editorContainer} ${errors?.description ? styles.editorError : ""}`}
                    >
                        <RichTextEditor
                            value={postDescription}
                            onChange={setPostDescription}
                            //   Use translation
                            placeholder={t('basicInfoContentPlaceholder')}
                            // Pass mediaId known at initialization
                            mediaId={initialData?.mediaId}
                            id="post-description" // Ensure RichTextEditor applies this ID
                            // RichTextEditor should ideally handle its own aria-invalid/describedby
                        />
                    </div>
                    {/* Display translated description error */}
                    {errors?.description && (
                        <div className={styles.fieldError} id="description-error" role="alert">
                            {translateError(errors.description)} {/*  Translate error key */}
                        </div>
                    )}
                </div>
                {/* Tags field */}
                <div className={styles.formGroup}>
                    <label className={styles.formLabel} htmlFor="post-tags">
                        Tags
                    </label>
                    <div className={styles.tagsInputWrapper}>
                        <Tag size={16} className={styles.tagIcon} aria-hidden="true"/>
                        <input
                            id="post-tags"
                            className={styles.formInput}
                            value={tags}
                            onChange={(e) => setTags(e.target.value)}
                            placeholder="Enter tags separated by commas"
                            maxLength={200}
                        />
                    </div>
                    <div className={styles.inputHint}>Separate tags with commas to help people find your post</div>
                </div>
                {/* Post Type and User Type Row */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="post-type">
                            Post Type
                        </label>
                        <select
                            id="post-type"
                            className={styles.formInput}
                            value={typeOfPost}
                            onChange={(e) => setTypeOfPost(e.target.value)}
                        >
                            <option value="general">General</option>
                            <option value="question">Question</option>
                            <option value="discussion">Discussion</option>
                            <option value="announcement">Announcement</option>
                            <option value="event">Event</option>
                            <option value="other">Other</option>
                        </select>
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="user-type">
                            User Type
                        </label>
                        <select
                            id="user-type"
                            className={styles.formInput}
                            value={userType}
                            onChange={(e) => setUserType(e.target.value)}
                        >
                            <option value="private">Private</option>
                            <option value="business">Business</option>
                        </select>
                    </div>
                </div>
                {/* Category and Status Row */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="category">
                            Category
                        </label>
                        <select
                            id="category"
                            className={styles.formInput}
                            value={categoryId}
                            onChange={(e) => setCategoryId(e.target.value)}
                        >
                            <option value="">Select Category</option>
                            {categories?.map((category) => (
                                <option key={category.id} value={category.id}>
                                    {category.name}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.formLabel} htmlFor="status">
                            Status
                        </label>
                        <select
                            id="status"
                            className={styles.formInput}
                            value={status}
                            onChange={(e) => setStatus(e.target.value)}
                        >
                            <option value="active">Active</option>
                            <option value="draft">Draft</option>
                            <option value="archived">Archived</option>
                        </select>
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
                    label="Post Location (Optional)"
                    error={errors?.location}
                    showCurrentLocationButton={true}
                />
                {/* Use FormActions, passing translated label */}
                {/* Assume FormActions handles translation of its "Cancel" button internally */}
                <FormActions
                    primaryLabel={primaryLabel} // Pass translated label
                    primaryIcon="arrow-right"
                    onCancel={onCancel}
                    isPrimaryDisabled={isPrimaryDisabled}
                    // Assuming primary button inside FormActions is type="submit"
                />
                {/* Display translated login requirement message */}
                {!isUserLoggedIn && !errors?.submit && (
                    <p className={`${styles.fieldError} ${styles.loginRequirement}`} role="status">
                        {/*   Use translation */}
                        {t('basicInfoLoginRequired')}
                    </p>
                )}
            </form>
        </div>
    );
}
// PropTypes remain the same
BasicInfoStep.propTypes = {
    initialData: PropTypes.shape({
        name: PropTypes.string,
        description: PropTypes.string,
        tags: PropTypes.string,
        mediaId: PropTypes.string,
        // New API fields
        typeOfPost: PropTypes.string,
        userType: PropTypes.string,
        categoryId: PropTypes.string,
        status: PropTypes.string,
        // Location fields
        address: PropTypes.string,
        lat: PropTypes.number,
        lng: PropTypes.number,
    }),
    onSubmit: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object, // Should contain error keys
    isUserLoggedIn: PropTypes.bool.isRequired,
    categories: PropTypes.array,
};
// Default props if needed
BasicInfoStep.defaultProps = {
    initialData: {
        name: '', 
        description: '', 
        tags: '', 
        mediaId: null,
        typeOfPost: 'general',
        userType: 'private',
        categoryId: '',
        status: 'active'
    },
    isLoading: false,
    errors: {},
    categories: [],
};
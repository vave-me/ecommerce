// src/components/ListingManagement/ListingForm.jsx
"use client"
import { FaExclamationTriangle } from '../../utils/iconImports';
import React, { useState, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import { Save, X, AlertCircle } from '@/icons';
import styles from './ListingForm.module.css';
/**
 * ListingForm Component
 * Form for creating and editing listings.
 */
const ListingForm = memo(({ listing, onSubmit, onCancel }) => {
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        price: '',
        category: '',
        condition: 'new',
        location: '',
        images: [],
        tags: ''
    });
    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    useEffect(() => {
        if (listing) {
            setFormData({
                title: listing.title || '',
                description: listing.description || '',
                price: listing.price || '',
                category: listing.category || '',
                condition: listing.condition || 'new',
                location: listing.location || '',
                images: listing.images || [],
                tags: Array.isArray(listing.tags) ? listing.tags.join(', ') : (listing.tags || '')
            });
        }
    }, [listing]);
    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
        // Clear error when user starts typing
        if (errors[name]) {
            setErrors(prev => ({
                ...prev,
                [name]: ''
            }));
        }
    };
    const validateForm = () => {
        const newErrors = {};
        if (!formData.title.trim()) {
            newErrors.title = 'Title is required';
        }
        if (!formData.description.trim()) {
            newErrors.description = 'Description is required';
        }
        if (!formData.price || parseFloat(formData.price) <= 0) {
            newErrors.price = 'Valid price is required';
        }
        if (!formData.category) {
            newErrors.category = 'Category is required';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };
    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!validateForm()) {
            return;
        }
        setIsSubmitting(true);
        try {
            const submitData = {
                ...formData,
                price: parseFloat(formData.price),
                tags: formData.tags.split(',').map(tag => tag.trim()).filter(Boolean)
            };
            await onSubmit(submitData);
        } catch (error) {
            setErrors({
                submit: 'Failed to save listing. Please try again.'
            });
        } finally {
            setIsSubmitting(false);
        }
    };
    return (
        <div className={styles.formContainer}>
            <div className={styles.formHeader}>
                <h2 className={styles.formTitle}>
                    {listing ? 'Edit Listing' : 'Create New Listing'}
                </h2>
                <button
                    className={styles.closeButton}
                    onClick={onCancel}
                    aria-label="Close form"
                >
                    <X size={20} />
                </button>
            </div>
            <form onSubmit={handleSubmit} className={styles.form}>
                {errors.submit && (
                    <div className={styles.errorAlert}>
                        <AlertCircle size={16} />
                        {errors.submit}
                    </div>
                )}
                <div className={styles.formGroup}>
                    <label htmlFor="title" className={styles.label}>
                        Title *
                    </label>
                    <input
                        type="text"
                        id="title"
                        name="title"
                        value={formData.title}
                        onChange={handleInputChange}
                        className={`${styles.input} ${errors.title ? styles.inputError : ''}`}
                        placeholder="Enter listing title"
                        disabled={isSubmitting}
                    />
                    {errors.title && (
                        <span className={styles.errorText}>{errors.title}</span>
                    )}
                </div>
                <div className={styles.formGroup}>
                    <label htmlFor="description" className={styles.label}>
                        Description *
                    </label>
                    <textarea
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                        className={`${styles.textarea} ${errors.description ? styles.inputError : ''}`}
                        placeholder="Describe your listing"
                        rows={4}
                        disabled={isSubmitting}
                    />
                    {errors.description && (
                        <span className={styles.errorText}>{errors.description}</span>
                    )}
                </div>
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label htmlFor="price" className={styles.label}>
                            Price *
                        </label>
                        <input
                            type="number"
                            id="price"
                            name="price"
                            value={formData.price}
                            onChange={handleInputChange}
                            className={`${styles.input} ${errors.price ? styles.inputError : ''}`}
                            placeholder="0.00"
                            min="0"
                            step="0.01"
                            disabled={isSubmitting}
                        />
                        {errors.price && (
                            <span className={styles.errorText}>{errors.price}</span>
                        )}
                    </div>
                    <div className={styles.formGroup}>
                        <label htmlFor="category" className={styles.label}>
                            Category *
                        </label>
                        <select
                            id="category"
                            name="category"
                            value={formData.category}
                            onChange={handleInputChange}
                            className={`${styles.select} ${errors.category ? styles.inputError : ''}`}
                            disabled={isSubmitting}
                        >
                            <option value="">Select category</option>
                            <option value="electronics">Electronics</option>
                            <option value="fashion">Fashion</option>
                            <option value="home">Home & Garden</option>
                            <option value="vehicles">Vehicles</option>
                            <option value="services">Services</option>
                        </select>
                        {errors.category && (
                            <span className={styles.errorText}>{errors.category}</span>
                        )}
                    </div>
                </div>
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label htmlFor="condition" className={styles.label}>
                            Condition
                        </label>
                        <select
                            id="condition"
                            name="condition"
                            value={formData.condition}
                            onChange={handleInputChange}
                            className={styles.select}
                            disabled={isSubmitting}
                        >
                            <option value="new">New</option>
                            <option value="like-new">Like New</option>
                            <option value="good">Good</option>
                            <option value="fair">Fair</option>
                            <option value="poor">Poor</option>
                        </select>
                    </div>
                    <div className={styles.formGroup}>
                        <label htmlFor="location" className={styles.label}>
                            Location
                        </label>
                        <input
                            type="text"
                            id="location"
                            name="location"
                            value={formData.location}
                            onChange={handleInputChange}
                            className={styles.input}
                            placeholder="City, State"
                            disabled={isSubmitting}
                        />
                    </div>
                </div>
                <div className={styles.formGroup}>
                    <label htmlFor="tags" className={styles.label}>
                        Tags
                    </label>
                    <input
                        type="text"
                        id="tags"
                        name="tags"
                        value={formData.tags}
                        onChange={handleInputChange}
                        className={styles.input}
                        placeholder="Enter tags separated by commas"
                        disabled={isSubmitting}
                    />
                </div>
                <div className={styles.formActions}>
                    <button
                        type="button"
                        className={styles.cancelButton}
                        onClick={onCancel}
                        disabled={isSubmitting}
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className={styles.submitButton}
                        disabled={isSubmitting}
                    >
                        <Save size={16} />
                        {isSubmitting ? 'Saving...' : (listing ? 'Update Listing' : 'Create Listing')}
                    </button>
                </div>
            </form>
        </div>
    );
});
ListingForm.displayName = 'ListingForm';
ListingForm.propTypes = {
    listing: PropTypes.object,
    onSubmit: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
};
export default ListingForm;

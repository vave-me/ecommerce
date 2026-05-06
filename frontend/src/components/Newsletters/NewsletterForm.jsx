// src/components/NewsletterForm.jsx
import React, { useState, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import { Save, X, AlertCircle } from '@/icons';
import styles from './NewsletterForm.module.css';
/**
 * NewsletterForm Component
 * Form for creating and editing newsletters.
 */
const NewsletterForm = memo(({ onClose, onSubmit, existingData }) => {
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        category: '',
        frequency: 'weekly',
        templateId: '',
        isActive: true
    });
    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    useEffect(() => {
        if (existingData) {
            setFormData({
                name: existingData.name || '',
                description: existingData.description || '',
                category: existingData.category || '',
                frequency: existingData.frequency || 'weekly',
                templateId: existingData.template_id || '',
                isActive: existingData.is_active !== undefined ? existingData.is_active : true
            });
        }
    }, [existingData]);
    const handleInputChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
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
        if (!formData.name.trim()) {
            newErrors.name = 'Newsletter name is required';
        }
        if (!formData.description.trim()) {
            newErrors.description = 'Description is required';
        }
        if (!formData.category.trim()) {
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
            await onSubmit(formData);
            onClose();
        } catch (error) {
            setErrors({
                submit: 'Failed to save newsletter. Please try again.'
            });
        } finally {
            setIsSubmitting(false);
        }
    };
    return (
        <div className={styles.modalOverlay}>
            <div className={styles.modalContainer}>
                <div className={styles.modalHeader}>
                    <h2 className={styles.modalTitle}>
                        {existingData ? 'Edit Newsletter' : 'Create Newsletter'}
                    </h2>
                    <button
                        className={styles.closeButton}
                        onClick={onClose}
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
                        <label htmlFor="name" className={styles.label}>
                            Newsletter Name *
                        </label>
                        <input
                            type="text"
                            id="name"
                            name="name"
                            value={formData.name}
                            onChange={handleInputChange}
                            className={`${styles.input} ${errors.name ? styles.inputError : ''}`}
                            placeholder="Enter newsletter name"
                            disabled={isSubmitting}
                        />
                        {errors.name && (
                            <span className={styles.errorText}>{errors.name}</span>
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
                            placeholder="Brief description of the newsletter"
                            rows={3}
                            disabled={isSubmitting}
                        />
                        {errors.description && (
                            <span className={styles.errorText}>{errors.description}</span>
                        )}
                    </div>
                    <div className={styles.formGroup}>
                        <label htmlFor="frequency" className={styles.label}>
                            Frequency
                        </label>
                        <select
                            id="frequency"
                            name="frequency"
                            value={formData.frequency}
                            onChange={handleInputChange}
                            className={styles.select}
                            disabled={isSubmitting}
                        >
                            <option value="daily">Daily</option>
                            <option value="weekly">Weekly</option>
                            <option value="monthly">Monthly</option>
                            <option value="quarterly">Quarterly</option>
                        </select>
                    </div>
                    <div className={styles.formGroup}>
                        <label htmlFor="category" className={styles.label}>
                            Category *
                        </label>
                        <input
                            type="text"
                            id="category"
                            name="category"
                            value={formData.category}
                            onChange={handleInputChange}
                            className={`${styles.input} ${errors.category ? styles.inputError : ''}`}
                            placeholder="e.g., Technology, Business, Health"
                            disabled={isSubmitting}
                        />
                        {errors.category && (
                            <span className={styles.errorText}>{errors.category}</span>
                        )}
                    </div>
                    <div className={styles.formGroup}>
                        <label htmlFor="templateId" className={styles.label}>
                            Template ID (Optional)
                        </label>
                        <input
                            type="text"
                            id="templateId"
                            name="templateId"
                            value={formData.templateId}
                            onChange={handleInputChange}
                            className={styles.input}
                            placeholder="Select a template"
                            disabled={isSubmitting}
                        />
                    </div>
                    <div className={styles.checkboxGroup}>
                        <label className={styles.checkboxLabel}>
                            <input
                                type="checkbox"
                                name="isActive"
                                checked={formData.isActive}
                                onChange={handleInputChange}
                                disabled={isSubmitting}
                            />
                            Active newsletter
                        </label>
                    </div>
                    <div className={styles.formActions}>
                        <button
                            type="button"
                            className={styles.cancelButton}
                            onClick={onClose}
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
                            {isSubmitting ? 'Saving...' : 'Save Newsletter'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
});
NewsletterForm.displayName = 'NewsletterForm';
NewsletterForm.propTypes = {
    onClose: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
    existingData: PropTypes.object,
};
export default NewsletterForm;
// File: src/components/Support/SupportTicketForm.jsx
"use client";
import React, { useState, useCallback, useMemo, useRef, useEffect, memo } from "react";
import PropTypes from "prop-types";
import { useTranslations } from 'next-intl';
// Icons with Header-inspired import pattern
import { 
    CheckCircle, 
    Send, 
    Loader2, 
    X, 
    AlertCircle,
    FileText,
    MessageSquare,
    Type,
    Edit3
} from '@/icons';
// Custom hooks (Header pattern)
import { useIsMobile } from "../../hooks/useMobileDetection";
// Styles
import styles from "./SupportTicketForm.module.css";
// Form validation patterns (Header-inspired)
const VALIDATION_RULES = {
    title: {
        required: true,
        minLength: 3,
        maxLength: 100
    },
    description: {
        required: true,
        minLength: 20,
        maxLength: 2000
    }
};
// Priority levels - matching the protobuf enum
const PRIORITY_LEVELS = [
    { value: 'LOW', label: 'Low', color: '#10b981' },
    { value: 'MEDIUM', label: 'Medium', color: '#3b82f6' },
    { value: 'HIGH', label: 'High', color: '#f59e0b' },
    { value: 'URGENT', label: 'Urgent', color: '#ef4444' },
    { value: 'CRITICAL', label: 'Critical', color: '#dc2626' }
];
// Category options - matching the protobuf enum
const CATEGORY_OPTIONS = [
    { value: 'GENERAL_INQUIRY', label: 'General Inquiry' },
    { value: 'TECHNICAL_ISSUE', label: 'Technical Issue' },
    { value: 'BILLING_ISSUE', label: 'Billing & Payment' },
    { value: 'ACCOUNT_ISSUE', label: 'Account Management' },
    { value: 'PRODUCT_QUESTION', label: 'Product Question' },
    { value: 'FEATURE_REQUEST', label: 'Feature Request' },
    { value: 'COMPLAINT', label: 'Complaint' },
    { value: 'REFUND_REQUEST', label: 'Refund Request' },
    { value: 'ORDER_ISSUE', label: 'Order Issue' },
    { value: 'SHIPPING_ISSUE', label: 'Shipping Issue' }
];
// Form state reducer (Header pattern)
const formReducer = (state, action) => {
    switch (action.type) {
        case 'SET_FIELD':
            return {
                ...state,
                formData: {
                    ...state.formData,
                    [action.field]: action.value
                },
                errors: {
                    ...state.errors,
                    [action.field]: null // Clear error when user types
                }
            };
        case 'SET_ERRORS':
            return {
                ...state,
                errors: action.payload
            };
        case 'SET_SUBMITTING':
            return {
                ...state,
                isSubmitting: action.payload
            };
        case 'SET_SUCCESS':
            return {
                ...state,
                showSuccess: action.payload
            };
        case 'RESET_FORM':
            return {
                ...state,
                formData: {
                    title: '',
                    description: '',
                    category: 'GENERAL_INQUIRY',
                    priority: 'MEDIUM',
                    tags: [],
                    attachments: []
                },
                errors: {},
                showSuccess: false
            };
        default:
            return state;
    }
};
// Initial form state
const initialFormState = {
    formData: {
        title: '',
        description: '',
        category: 'GENERAL_INQUIRY',
        priority: 'MEDIUM',
        tags: [],
        attachments: []
    },
    errors: {},
    isSubmitting: false,
    showSuccess: false
};
/**
 * Enhanced SupportTicketForm Component with Header-inspired patterns
 * Features advanced validation, performance optimizations, and sophisticated UX
 */
const SupportTicketForm = memo(({ onCreate, onCancel, isCreating = false, isMobile: propIsMobile = false }) => {
    const t = useTranslations('SupportTicketForm');
    const isMobile = useIsMobile() || propIsMobile;
    // Enhanced state management with reducer (Header pattern)
    const [state, dispatch] = React.useReducer(formReducer, initialFormState);
    const { formData, errors, isSubmitting, showSuccess } = state;
    // Refs for focus management (Header pattern)
    const titleRef = useRef(null);
    const descriptionRef = useRef(null);
    const formRef = useRef(null);
    // Auto-focus on mount (Header pattern)
    useEffect(() => {
        if (!isMobile && titleRef.current) {
            titleRef.current.focus();
        }
    }, [isMobile]);
    // Enhanced validation function (Header-inspired)
    const validateField = useCallback((field, value) => {
        const rules = VALIDATION_RULES[field];
        if (!rules) return null;
        if (rules.required && !value.trim()) {
            return t(`error_${field}_required`);
        }
        if (rules.minLength && value.trim().length < rules.minLength) {
            return t(`error_${field}_minLength`, { min: rules.minLength });
        }
        if (rules.maxLength && value.length > rules.maxLength) {
            return t(`error_${field}_maxLength`, { max: rules.maxLength });
        }
        return null;
    }, [t]);
    // Validate entire form (Header pattern)
    const validateForm = useCallback(() => {
        const newErrors = {};
        Object.keys(VALIDATION_RULES).forEach(field => {
            const error = validateField(field, formData[field]);
            if (error) {
                newErrors[field] = error;
            }
        });
        return newErrors;
    }, [formData, validateField]);
    // Handle field changes with real-time validation (Header pattern)
    const handleFieldChange = useCallback((field, value) => {
        dispatch({ type: 'SET_FIELD', field, value });
        // Real-time validation for better UX
        const error = validateField(field, value);
        if (error) {
            dispatch({ 
                type: 'SET_ERRORS', 
                payload: { ...errors, [field]: error } 
            });
        }
    }, [errors, validateField]);
    // Handle form submission (Enhanced with Header patterns)
    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        // Validate form
        const validationErrors = validateForm();
        if (Object.keys(validationErrors).length > 0) {
            dispatch({ type: 'SET_ERRORS', payload: validationErrors });
            // Focus first field with error (Header pattern)
            const firstErrorField = Object.keys(validationErrors)[0];
            const fieldRef = firstErrorField === 'title' ? titleRef : descriptionRef;
            if (fieldRef.current) {
                fieldRef.current.focus();
            }
            return;
        }
        dispatch({ type: 'SET_SUBMITTING', payload: true });
        try {
            await onCreate(formData);
            // Show success state
            dispatch({ type: 'SET_SUCCESS', payload: true });
            // Reset form after delay
            setTimeout(() => {
                dispatch({ type: 'RESET_FORM' });
                if (onCancel) onCancel();
            }, 1500);
        } catch (error) {
            dispatch({ 
                type: 'SET_ERRORS', 
                payload: { form: t('error_form_submit') } 
            });
        } finally {
            dispatch({ type: 'SET_SUBMITTING', payload: false });
        }
    }, [formData, validateForm, onCreate, onCancel, t]);
    // Handle cancel (Header pattern)
    const handleCancel = useCallback(() => {
        if (onCancel) {
            onCancel();
        }
    }, [onCancel]);
    // Memoized character counts (Header performance pattern)
    const characterCounts = useMemo(() => ({
        title: formData.title.length,
        description: formData.description.length
    }), [formData.title.length, formData.description.length]);
    // Memoized form validation state (Header pattern)
    const isFormValid = useMemo(() => {
        return Object.keys(validateForm()).length === 0 && 
               formData.title.trim() && 
               formData.description.trim();
    }, [formData, validateForm]);
    // Form classes (Header pattern)
    const formClasses = [
        styles.formContainer,
        isMobile ? styles.mobile : styles.desktop,
        isSubmitting ? styles.submitting : '',
        showSuccess ? styles.success : ''
    ].filter(Boolean).join(' ');
    return (
        <div className={formClasses}>
            <div className={styles.formHeader}>
                <div className={styles.headerLeft}>
                    <MessageSquare size={24} className={styles.headerIcon} />
                    <div className={styles.headerText}>
                        <h3 className={styles.formTitle}>{t('formTitle')}</h3>
                        <p className={styles.formSubtitle}>{t('formSubtitle')}</p>
                    </div>
                </div>
                {onCancel && (
                    <button
                        type="button"
                        className={styles.closeButton}
                        onClick={handleCancel}
                        aria-label={t('closeAriaLabel')}
                    >
                        <X size={20} />
                    </button>
                )}
            </div>
            <form ref={formRef} onSubmit={handleSubmit} noValidate className={styles.form}>
                {/* Form-level error */}
                {errors.form && (
                    <div className={styles.formError} role="alert">
                        <AlertCircle size={16} />
                        <span>{errors.form}</span>
                    </div>
                )}
                {/* Category & Priority Row */}
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.label} htmlFor="category">
                            {t('categoryLabel')}
                        </label>
                        <div className={styles.selectWrapper}>
                            <select
                                id="category"
                                className={styles.select}
                                value={formData.category}
                                onChange={(e) => handleFieldChange('category', e.target.value)}
                            >
                                {CATEGORY_OPTIONS.map(option => (
                                    <option key={option.value} value={option.value}>
                                        {t(`category_${option.value}`, option.label)}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>
                    <div className={styles.formGroup}>
                        <label className={styles.label} htmlFor="priority">
                            {t('priorityLabel')}
                        </label>
                        <div className={styles.selectWrapper}>
                            <select
                                id="priority"
                                className={styles.select}
                                value={formData.priority}
                                onChange={(e) => handleFieldChange('priority', e.target.value)}
                            >
                                {PRIORITY_LEVELS.map(level => (
                                    <option key={level.value} value={level.value}>
                                        {t(`priority_${level.value}`, level.label)}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>
                </div>
                {/* Title Input */}
                <div className={styles.formGroup}>
                    <label className={styles.label} htmlFor="title">
                        {t('titleLabel')}
                        <span className={styles.required}>*</span>
                    </label>
                    <div className={styles.inputWrapper}>
                        <Type size={16} className={styles.inputIcon} />
                        <input
                            ref={titleRef}
                            type="text"
                            id="title"
                            className={`${styles.input} ${errors.title ? styles.inputError : ''}`}
                            value={formData.title}
                            onChange={(e) => handleFieldChange('title', e.target.value)}
                            placeholder={t('titlePlaceholder')}
                            aria-invalid={errors.title ? "true" : "false"}
                            aria-describedby={errors.title ? "title-error" : "title-help"}
                            maxLength={VALIDATION_RULES.title.maxLength}
                        />
                        <div className={styles.characterCount}>
                            <span className={characterCounts.title > VALIDATION_RULES.title.maxLength * 0.9 ? styles.warning : ''}>
                                {characterCounts.title}/{VALIDATION_RULES.title.maxLength}
                            </span>
                        </div>
                    </div>
                    {errors.title ? (
                        <div id="title-error" className={styles.errorMessage} role="alert">
                            <AlertCircle size={14} />
                            <span>{errors.title}</span>
                        </div>
                    ) : (
                        <div id="title-help" className={styles.helpText}>
                            {t('titleHelp')}
                        </div>
                    )}
                </div>
                {/* Description Textarea */}
                <div className={styles.formGroup}>
                    <label className={styles.label} htmlFor="description">
                        {t('descriptionLabel')}
                        <span className={styles.required}>*</span>
                    </label>
                    <div className={styles.textareaWrapper}>
                        <Edit3 size={16} className={styles.textareaIcon} />
                        <textarea
                            ref={descriptionRef}
                            id="description"
                            className={`${styles.textarea} ${errors.description ? styles.inputError : ''}`}
                            value={formData.description}
                            onChange={(e) => handleFieldChange('description', e.target.value)}
                            placeholder={t('descriptionPlaceholder')}
                            aria-invalid={errors.description ? "true" : "false"}
                            aria-describedby={errors.description ? "description-error" : "description-help"}
                            rows={isMobile ? 4 : 6}
                            maxLength={VALIDATION_RULES.description.maxLength}
                        />
                        <div className={styles.characterCount}>
                            <span className={characterCounts.description > VALIDATION_RULES.description.maxLength * 0.9 ? styles.warning : ''}>
                                {characterCounts.description}/{VALIDATION_RULES.description.maxLength}
                            </span>
                        </div>
                    </div>
                    {errors.description ? (
                        <div id="description-error" className={styles.errorMessage} role="alert">
                            <AlertCircle size={14} />
                            <span>{errors.description}</span>
                        </div>
                    ) : (
                        <div id="description-help" className={styles.helpText}>
                            {t('descriptionHelp')}
                        </div>
                    )}
                </div>
                {/* Form Actions */}
                <div className={styles.formActions}>
                    {onCancel && (
                        <button
                            type="button"
                            className={styles.cancelButton}
                            onClick={handleCancel}
                            disabled={isSubmitting}
                        >
                            {t('cancelButton')}
                        </button>
                    )}
                    <button
                        type="submit"
                        className={`${styles.submitButton} ${!isFormValid ? styles.disabled : ''}`}
                        disabled={isSubmitting || !isFormValid}
                    >
                        {isSubmitting || isCreating ? (
                            <>
                                <Loader2 size={16} className={styles.spinning} />
                                {t('submittingText')}
                            </>
                        ) : showSuccess ? (
                            <>
                                <CheckCircle size={16} />
                                {t('successText')}
                            </>
                        ) : (
                            <>
                                <Send size={16} />
                                {t('submitButton')}
                            </>
                        )}
                    </button>
                </div>
                {/* Success Message */}
                {showSuccess && (
                    <div className={styles.successMessage} role="alert">
                        <CheckCircle size={20} />
                        <div className={styles.successContent}>
                            <h4>{t('successTitle')}</h4>
                            <p>{t('successMessage')}</p>
                        </div>
                    </div>
                )}
            </form>
        </div>
    );
});
SupportTicketForm.displayName = 'SupportTicketForm';
SupportTicketForm.propTypes = {
    onCreate: PropTypes.func.isRequired,
    onCancel: PropTypes.func,
    isCreating: PropTypes.bool,
    isMobile: PropTypes.bool
};
export default SupportTicketForm;
/**
 * MobileForm - Enhanced Mobile-First Form Component
 *
 * Features:
 * - Intelligent keyboard type detection
 * - Dynamic viewport adjustment for keyboard
 * - Enhanced validation with real-time feedback
 * - Touch-optimized input interactions
 * - Auto-scroll to active fields
 * - Haptic feedback for form interactions
 *
 * Designed to work with existing form patterns while providing mobile enhancements
 */
"use client";
import React, {
    useRef,
    useCallback,
    useEffect,
    useState,
    forwardRef,
    useImperativeHandle,
    Children,
    cloneElement,
    memo
} from 'react';
import PropTypes from 'prop-types';
import { FaExclamationTriangle, FaCheckCircle } from '../../utils/iconImports';
// Mobile keyboard configuration
const KEYBOARD_CONFIG = {
    // Keyboard types for better mobile UX
    KEYBOARDS: {
        text: 'text',
        email: 'email',
        tel: 'tel',
        number: 'numeric',
        url: 'url',
        search: 'search',
        password: 'text'
    },
    // Input modes for modern browsers
    INPUT_MODES: {
        text: 'text',
        email: 'email',
        tel: 'tel',
        number: 'numeric',
        decimal: 'decimal',
        url: 'url',
        search: 'search'
    },
    // Auto-complete suggestions
    AUTOCOMPLETE: {
        name: 'name',
        email: 'email',
        tel: 'tel',
        'street-address': 'street-address',
        'postal-code': 'postal-code',
        'current-password': 'current-password',
        'new-password': 'new-password'
    }
};
// Validation patterns
const VALIDATION_PATTERNS = {
    email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    tel: /^[\+]?[\s\-\(\)]?[\d\s\-\(\)]{10,}$/,
    url: /^https?:\/\/.+\..+/,
    zipcode: /^\d{5}(-\d{4})?$/,
    creditcard: /^\d{4}\s?\d{4}\s?\d{4}\s?\d{4}$/
};
const MobileForm = memo(forwardRef(({
                                   children,
                                   onSubmit,
                                   validation = {},
                                   autoScrollToError = true,
                                   hapticFeedback = true,
                                   keyboardAdjustment = true,
                                   className = '',
                                   style = {},
                                   ...props
                               }, ref) => {
    // Refs and state
    const formRef = useRef(null);
    const fieldsRef = useRef(new Map());
    const [formData, setFormData] = useState({});
    const [errors, setErrors] = useState({});
    const [touched, setTouched] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [keyboardVisible, setKeyboardVisible] = useState(false);
    const [activeField, setActiveField] = useState(null);
    // Expose form methods
    useImperativeHandle(ref, () => ({
        submit: () => handleSubmit(),
        reset: () => handleReset(),
        setFieldValue: (name, value) => setFormData(prev => ({ ...prev, [name]: value })),
        setFieldError: (name, error) => setErrors(prev => ({ ...prev, [name]: error })),
        getFormData: () => formData,
        getErrors: () => errors,
        isValid: () => Object.keys(errors).length === 0,
        scrollToField: (fieldName) => scrollToField(fieldName)
    }), [formData, errors]);
    // Keyboard visibility detection
    useEffect(() => {
        if (!keyboardAdjustment) return;
        let initialHeight = window.innerHeight;
        const handleResize = () => {
            const currentHeight = window.innerHeight;
            const heightDifference = initialHeight - currentHeight;
            // Threshold for keyboard detection (mobile-specific)
            if (heightDifference > 150) {
                setKeyboardVisible(true);
            } else {
                setKeyboardVisible(false);
            }
        };
        const handleOrientationChange = () => {
            setTimeout(() => {
                initialHeight = window.innerHeight;
                setKeyboardVisible(false);
            }, 500);
        };
        window.addEventListener('resize', handleResize);
        window.addEventListener('orientationchange', handleOrientationChange);
        return () => {
            window.removeEventListener('resize', handleResize);
            window.removeEventListener('orientationchange', handleOrientationChange);
        };
    }, [keyboardAdjustment]);
    // Haptic feedback helper
    const triggerHaptic = useCallback((type = 'light') => {
        if (!hapticFeedback || !navigator.vibrate) return;
        try {
            switch (type) {
                case 'light':
                    navigator.vibrate(10);
                    break;
                case 'error':
                    navigator.vibrate([10, 10, 10]);
                    break;
                case 'success':
                    navigator.vibrate([10, 100, 10]);
                    break;
                default:
                    navigator.vibrate(10);
            }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, [hapticFeedback]);
    // Validate field
    const validateField = useCallback((name, value, rules) => {
        if (!rules) return null;
        // Required validation
        if (rules.required && (!value || value.toString().trim() === '')) {
            return rules.required === true ? `${name} is required` : rules.required;
        }
        // Skip other validations if field is empty and not required
        if (!value || value.toString().trim() === '') {
            return null;
        }
        // Pattern validation
        if (rules.pattern) {
            const pattern = typeof rules.pattern === 'string'
                ? VALIDATION_PATTERNS[rules.pattern] || new RegExp(rules.pattern)
                : rules.pattern;
            if (!pattern.test(value)) {
                return rules.patternMessage || `Invalid ${name} format`;
            }
        }
        // Length validation
        if (rules.minLength && value.length < rules.minLength) {
            return `${name} must be at least ${rules.minLength} characters`;
        }
        if (rules.maxLength && value.length > rules.maxLength) {
            return `${name} must be no more than ${rules.maxLength} characters`;
        }
        // Custom validation
        if (rules.validate && typeof rules.validate === 'function') {
            return rules.validate(value, formData);
        }
        return null;
    }, [formData]);
    // Handle field change
    const handleFieldChange = useCallback((name, value) => {
        setFormData(prev => ({ ...prev, [name]: value }));
        // Clear error when user starts typing
        if (errors[name]) {
            setErrors(prev => {
                const newErrors = { ...prev };
                delete newErrors[name];
                return newErrors;
            });
        }
        // Real-time validation for touched fields
        if (touched[name] && validation[name]) {
            const error = validateField(name, value, validation[name]);
            if (error) {
                setErrors(prev => ({ ...prev, [name]: error }));
            }
        }
    }, [errors, touched, validation, validateField]);
    // Handle field focus
    const handleFieldFocus = useCallback((name) => {
        setActiveField(name);
        triggerHaptic('light');
        // Auto-scroll to field on mobile
        if (keyboardAdjustment) {
            setTimeout(() => {
                scrollToField(name);
            }, 300); // Wait for keyboard animation
        }
    }, [keyboardAdjustment, triggerHaptic]);
    // Handle field blur
    const handleFieldBlur = useCallback((name, value) => {
        setTouched(prev => ({ ...prev, [name]: true }));
        setActiveField(null);
        // Validate on blur
        if (validation[name]) {
            const error = validateField(name, value, validation[name]);
            if (error) {
                setErrors(prev => ({ ...prev, [name]: error }));
                triggerHaptic('error');
            }
        }
    }, [validation, validateField, triggerHaptic]);
    // Scroll to field helper
    const scrollToField = useCallback((fieldName) => {
        const fieldElement = fieldsRef.current.get(fieldName);
        if (!fieldElement) return;
        const rect = fieldElement.getBoundingClientRect();
        const viewportHeight = window.innerHeight;
        const keyboardHeight = keyboardVisible ? viewportHeight * 0.4 : 0;
        const availableHeight = viewportHeight - keyboardHeight;
        // Check if field is in viewport
        if (rect.top < 0 || rect.bottom > availableHeight) {
            const scrollOffset = keyboardHeight / 2;
            fieldElement.scrollIntoView({
                behavior: 'smooth',
                block: 'center',
                inline: 'nearest'
            });
            // Additional adjustment for keyboard
            if (keyboardVisible) {
                setTimeout(() => {
                    window.scrollBy(0, -scrollOffset);
                }, 300);
            }
        }
    }, [keyboardVisible]);
    // Handle form submission
    const handleSubmit = useCallback((event) => {
        if (event) {
            event.preventDefault();
        }
        setIsSubmitting(true);
        // Validate all fields
        const newErrors = {};
        Object.keys(validation).forEach(fieldName => {
            const value = formData[fieldName];
            const error = validateField(fieldName, value, validation[fieldName]);
            if (error) {
                newErrors[fieldName] = error;
            }
        });
        setErrors(newErrors);
        if (Object.keys(newErrors).length > 0) {
            // Scroll to first error
            if (autoScrollToError) {
                const firstErrorField = Object.keys(newErrors)[0];
                scrollToField(firstErrorField);
            }
            triggerHaptic('error');
            setIsSubmitting(false);
            return;
        }
        // Trigger success haptic
        triggerHaptic('success');
        // Call onSubmit with form data
        if (onSubmit) {
            try {
                const result = onSubmit(formData);
                // Handle async onSubmit
                if (result && typeof result.then === 'function') {
                    result
                        .then(() => setIsSubmitting(false))
                        .catch(() => setIsSubmitting(false));
                } else {
                    setIsSubmitting(false);
                }
            } catch (error) {
                setIsSubmitting(false);
                triggerHaptic('error');
            }
        } else {
            setIsSubmitting(false);
        }
    }, [formData, validation, validateField, autoScrollToError, scrollToField, triggerHaptic, onSubmit]);
    // Handle form reset
    const handleReset = useCallback(() => {
        setFormData({});
        setErrors({});
        setTouched({});
        setIsSubmitting(false);
        setActiveField(null);
        triggerHaptic('light');
    }, [triggerHaptic]);
    // Auto-configure input based on name/type
    const getInputConfig = useCallback((name, type) => {
        const config = {
            autoComplete: 'off',
            inputMode: 'text',
            enterKeyHint: 'next'
        };
        // Auto-configure based on field name
        if (name.includes('email')) {
            config.type = 'email';
            config.inputMode = 'email';
            config.autoComplete = 'email';
            config.enterKeyHint = 'next';
        } else if (name.includes('phone') || name.includes('tel')) {
            config.type = 'tel';
            config.inputMode = 'tel';
            config.autoComplete = 'tel';
        } else if (name.includes('url') || name.includes('website')) {
            config.type = 'url';
            config.inputMode = 'url';
            config.autoComplete = 'url';
        } else if (name.includes('password')) {
            config.type = 'password';
            config.autoComplete = name.includes('new') ? 'new-password' : 'current-password';
        } else if (name.includes('price') || name.includes('amount') || name.includes('number')) {
            config.type = 'number';
            config.inputMode = 'numeric';
            config.pattern = '[0-9]*';
        }
        // Override with explicit type
        if (type) {
            config.type = type;
            if (KEYBOARD_CONFIG.INPUT_MODES[type]) {
                config.inputMode = KEYBOARD_CONFIG.INPUT_MODES[type];
            }
        }
        return config;
    }, []);
    // Enhanced children with form context
    const enhancedChildren = Children.map(children, (child) => {
        if (!child || typeof child !== 'object') return child;
        // Only enhance input elements
        if (child.type === 'input' || child.type === 'textarea' || child.type === 'select') {
            const { name, type, onChange, onFocus, onBlur, ...childProps } = child.props;
            if (!name) return child;
            const inputConfig = getInputConfig(name, type);
            const hasError = !!errors[name];
            const value = formData[name] || '';
            return cloneElement(child, {
                ...inputConfig,
                ...childProps,
                value,
                ref: (el) => {
                    if (el) fieldsRef.current.set(name, el);
                },
                onChange: (e) => {
                    handleFieldChange(name, e.target.value);
                    if (onChange) onChange(e);
                },
                onFocus: (e) => {
                    handleFieldFocus(name);
                    if (onFocus) onFocus(e);
                },
                onBlur: (e) => {
                    handleFieldBlur(name, e.target.value);
                    if (onBlur) onBlur(e);
                },
                'aria-invalid': hasError,
                'aria-describedby': hasError ? `${name}-error` : undefined,
                style: {
                    ...child.props.style,
                    borderColor: hasError ? '#ef4444' : undefined,
                    ...(activeField === name && {
                        borderColor: '#2980b9',
                        boxShadow: '0 0 0 3px rgba(79, 70, 229, 0.1)'
                    })
                }
            });
        }
        return child;
    });
    // Render field errors
    const renderFieldErrors = () => {
        return Object.entries(errors).map(([fieldName, error]) => (
            <div
                key={`${fieldName}-error`}
                id={`${fieldName}-error`}
                style={{
                    color: '#ef4444',
                    fontSize: '14px',
                    marginTop: '4px',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px'
                }}
                role="alert"
                aria-live="polite"
            >
                <FaExclamationTriangle size={14} />
                <span>{error}</span>
            </div>
        ));
    };
    return (
        <form
            ref={formRef}
            onSubmit={handleSubmit}
            className={`mobile-form ${className}`}
            style={{
                position: 'relative',
                paddingBottom: keyboardVisible ? '20px' : '0',
                transition: 'padding-bottom 0.3s ease',
                ...style
            }}
            noValidate
            {...props}
        >
            {enhancedChildren}
            {renderFieldErrors()}
            {/* Success indicators */}
            {Object.keys(formData).length > 0 && Object.keys(errors).length === 0 && (
                <div
                    style={{
                        position: 'fixed',
                        top: '20px',
                        right: '20px',
                        backgroundColor: '#10b981',
                        color: 'white',
                        padding: '8px 12px',
                        borderRadius: '8px',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        fontSize: '14px',
                        fontWeight: '500',
                        zIndex: 1000,
                        opacity: Object.keys(touched).length > 0 ? 1 : 0,
                        transition: 'opacity 0.3s ease'
                    }}
                >
                    <FaCheckCircle size={16} />
                    Form looks good!
                </div>
            )}
        </form>
    );
}));
MobileForm.displayName = 'MobileForm';
MobileForm.propTypes = {
    children: PropTypes.node.isRequired,
    onSubmit: PropTypes.func,
    validation: PropTypes.object,
    autoScrollToError: PropTypes.bool,
    hapticFeedback: PropTypes.bool,
    keyboardAdjustment: PropTypes.bool,
    className: PropTypes.string,
    style: PropTypes.object,
};
export default MobileForm;
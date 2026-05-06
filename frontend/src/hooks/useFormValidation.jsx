// Consolidated form validation hook
import { useState, useCallback } from 'react';
// Standard validation constants
const VALIDATION_CONSTANTS = {
  MIN_TITLE_LENGTH: 5,
  MIN_CONTENT_LENGTH: 20,
  MIN_PASSWORD_LENGTH: 8,
  EMAIL_REGEX: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
};
/**
 * Enhanced form validation hook
 * - Validates common form fields
 * - Customizable validation rules
 * - Real-time or on-submit validation
 * 
 * @param {Object} options - Configuration options
 * @returns {Object} Validation utilities and state
 */
export function useFormValidation(options = {}) {
  const {
    initialErrors = {},
    validateOnChange = false,
    constants = VALIDATION_CONSTANTS
  } = options;
  const [errors, setErrors] = useState(initialErrors);
  const [touched, setTouched] = useState({});
  // Clear all errors
  const clearErrors = useCallback(() => {
    setErrors({});
  }, []);
  // Set a specific field as touched
  const markAsTouched = useCallback((field) => {
    setTouched(prev => ({ ...prev, [field]: true }));
  }, []);
  // Common validation functions
  const validators = {
    required: (value, field, message = 'This field is required') => {
      if (!value || (typeof value === 'string' && !value.trim())) {
        return { [field]: message };
      }
      return {};
    },
    minLength: (value, field, minLength, message) => {
      if (!value) return {};
      const defaultMessage = `Must be at least ${minLength} characters`;
      if (value.trim().length < minLength) {
        return { [field]: message || defaultMessage };
      }
      return {};
    },
    email: (value, field, message = 'Invalid email address') => {
      if (!value) return {};
      if (!constants.EMAIL_REGEX.test(value)) {
        return { [field]: message };
      }
      return {};
    },
    htmlContent: (value, field, minLength = constants.MIN_CONTENT_LENGTH, message) => {
      if (!value) return {};
      const strippedContent = value.replace(/<[^>]*>/g, "").trim();
      const defaultMessage = `Content must be at least ${minLength} characters`;
      if (!strippedContent) {
        return { [field]: message || 'Content is required' };
      }
      if (strippedContent.length < minLength) {
        return { [field]: message || defaultMessage };
      }
      return {};
    },
    // Add more validators as needed
  };
  // Validate a specific field
  const validateField = useCallback((field, value, rules = []) => {
    let fieldErrors = {};
    rules.forEach(rule => {
      if (typeof rule === 'function') {
        // Custom validator function
        const result = rule(value, field);
        if (result) {
          fieldErrors = { ...fieldErrors, ...result };
        }
      } else if (rule.type === 'required') {
        const result = validators.required(value, field, rule.message);
        fieldErrors = { ...fieldErrors, ...result };
      } else if (rule.type === 'minLength') {
        const result = validators.minLength(value, field, rule.length, rule.message);
        fieldErrors = { ...fieldErrors, ...result };
      } else if (rule.type === 'email') {
        const result = validators.email(value, field, rule.message);
        fieldErrors = { ...fieldErrors, ...result };
      } else if (rule.type === 'htmlContent') {
        const result = validators.htmlContent(value, field, rule.minLength, rule.message);
        fieldErrors = { ...fieldErrors, ...result };
      }
    });
    // Update the errors state
    setErrors(prev => ({ 
      ...prev,
      ...fieldErrors
    }));
    return Object.keys(fieldErrors).length === 0;
  }, [validators]);
  // Legacy support for validateForm
  const validateForm = useCallback((formData) => {
    const { name, description, ...rest } = formData;
    const newErrors = {};
    // Validate title
    if (!name?.trim()) {
      newErrors.name = "Title is required";
    } else if (name.trim().length < constants.MIN_TITLE_LENGTH) {
      newErrors.name = `Title must be at least ${constants.MIN_TITLE_LENGTH} characters`;
    }
    // Validate description
    if (description) {
      const strippedContent = description?.replace(/<[^>]*>/g, "").trim();
      if (!strippedContent) {
        newErrors.description = "Content is required";
      } else if (strippedContent.length < constants.MIN_CONTENT_LENGTH) {
        newErrors.description = `Content must be at least ${constants.MIN_CONTENT_LENGTH} characters`;
      }
    }
    // Set errors
    setErrors(newErrors);
    return {
      isValid: Object.keys(newErrors).length === 0,
      errors: newErrors
    };
  }, [constants]);
  return {
    // Form state
    errors,
    setErrors,
    touched,
    // Validation methods
    validateField,
    validateForm,
    clearErrors,
    markAsTouched,
    // Helper data
    constants,
    validators,
    // Check if form is valid
    isValid: Object.keys(errors).length === 0
  };
}
// Default export
export default useFormValidation; 
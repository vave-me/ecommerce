import { useState, useEffect, useCallback } from 'react';
/**
 * Custom hook for debounced form validation
 * 
 * @param {Object} formData - The form data to validate
 * @param {Function} validationFunction - Function that validates the form data and returns errors
 * @param {number} delay - Debounce delay in milliseconds
 * @returns {Object} Object containing validation errors, isValidating status, and validateNow function
 */
export function useDebouncedValidation(formData, validationFunction, delay = 500) {
  const [errors, setErrors] = useState({});
  const [debouncedFormData, setDebouncedFormData] = useState(formData);
  const [isValidating, setIsValidating] = useState(false);
  // Update debounced value after specified delay
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedFormData(formData);
    }, delay);
    // Clean up timeout on component unmount or when formData changes
    return () => {
      clearTimeout(handler);
    };
  }, [formData, delay]);
  // Run validation when debounced form data changes
  useEffect(() => {
    const validateForm = async () => {
      try {
        setIsValidating(true);
        const validationErrors = await validationFunction(debouncedFormData);
        setErrors(validationErrors || {});
      } catch (error) {
        setErrors({ form: 'Validation failed. Please try again.' });
      } finally {
        setIsValidating(false);
      }
    };
    validateForm();
  }, [debouncedFormData, validationFunction]);
  // Function to force immediate validation
  const validateNow = useCallback(async () => {
    try {
      setIsValidating(true);
      const validationErrors = await validationFunction(formData);
      setErrors(validationErrors || {});
      return validationErrors || {};
    } catch (error) {
      const formError = { form: 'Validation failed. Please try again.' };
      setErrors(formError);
      return formError;
    } finally {
      setIsValidating(false);
    }
  }, [formData, validationFunction]);
  return { errors, isValidating, validateNow };
} 
/**
 * Utility functions for form validation
 */
/**
 * Validates a product form's data
 * @param {Object} formData - The form data to validate
 * @returns {Object|null} - Validation errors or null if valid
 */
export const validateProductForm = async (formData) => {
  const errors = {};
  // Validate required fields
  if (!formData.name?.trim()) {
    errors.name = 'Product name is required';
  } else if (formData.name.length < 3) {
    errors.name = 'Product name must be at least 3 characters';
  }
  if (!formData.description?.trim()) {
    errors.description = 'Description is required';
  } else if (formData.description.length < 10) {
    errors.description = 'Description must be at least 10 characters';
  }
  if (!formData.basePrice) {
    errors.basePrice = 'Price is required';
  } else if (isNaN(Number(formData.basePrice)) || Number(formData.basePrice) <= 0) {
    errors.basePrice = 'Price must be a positive number';
  }
  if (!formData.categoryId) {
    errors.categoryId = 'Category is required';
  }
  // Return null if no errors
  return Object.keys(errors).length ? errors : null;
};
/**
 * Validates a service form's data
 * @param {Object} formData - The form data to validate
 * @returns {Object|null} - Validation errors or null if valid
 */
export const validateServiceForm = async (formData) => {
  const errors = {};
  // Validate required fields
  if (!formData.name?.trim()) {
    errors.name = 'Service name is required';
  } else if (formData.name.length < 3) {
    errors.name = 'Service name must be at least 3 characters';
  }
  if (!formData.description?.trim()) {
    errors.description = 'Description is required';
  } else if (formData.description.length < 10) {
    errors.description = 'Description must be at least 10 characters';
  }
  if (!formData.basePrice) {
    errors.basePrice = 'Price is required';
  } else if (isNaN(Number(formData.basePrice)) || Number(formData.basePrice) <= 0) {
    errors.basePrice = 'Price must be a positive number';
  }
  if (!formData.categoryId) {
    errors.categoryId = 'Category is required';
  }
  if (!formData.serviceType) {
    errors.serviceType = 'Service type is required';
  }
  // Return null if no errors
  return Object.keys(errors).length ? errors : null;
};
/**
 * Validates a deal form's data
 * @param {Object} formData - The form data to validate
 * @returns {Object|null} - Validation errors or null if valid
 */
export const validateDealForm = async (formData) => {
  const errors = {};
  // Validate required fields
  if (!formData.name?.trim()) {
    errors.name = 'Deal name is required';
  } else if (formData.name.length < 3) {
    errors.name = 'Deal name must be at least 3 characters';
  }
  if (!formData.description?.trim()) {
    errors.description = 'Description is required';
  } else if (formData.description.length < 10) {
    errors.description = 'Description must be at least 10 characters';
  }
  if (!formData.basePrice) {
    errors.basePrice = 'Original price is required';
  } else if (isNaN(Number(formData.basePrice)) || Number(formData.basePrice) <= 0) {
    errors.basePrice = 'Original price must be a positive number';
  }
  if (!formData.dealPrice) {
    errors.dealPrice = 'Deal price is required';
  } else if (isNaN(Number(formData.dealPrice)) || Number(formData.dealPrice) <= 0) {
    errors.dealPrice = 'Deal price must be a positive number';
  } else if (Number(formData.dealPrice) >= Number(formData.basePrice)) {
    errors.dealPrice = 'Deal price must be less than original price';
  }
  if (!formData.categoryId) {
    errors.categoryId = 'Category is required';
  }
  if (!formData.dealUrl?.trim()) {
    errors.dealUrl = 'Deal URL is required';
  } else if (!/^https?:\/\/.+/.test(formData.dealUrl)) {
    errors.dealUrl = 'Please enter a valid URL starting with http:// or https://';
  }
  // Return null if no errors
  return Object.keys(errors).length ? errors : null;
}; 
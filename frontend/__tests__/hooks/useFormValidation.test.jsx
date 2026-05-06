import { renderHook, act } from '@testing-library/react';
import { useFormValidation } from '@/hooks/useFormValidation.jsx';

describe('useFormValidation hook', () => {
  test('should initialize with empty errors', () => {
    const { result } = renderHook(() => useFormValidation());
    
    expect(result.current.errors).toEqual({});
    expect(result.current.isValid).toBe(true);
    expect(result.current.touched).toEqual({});
  });
  
  test('should validate required fields', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: '',
        description: 'Some content'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title is required');
  });
  
  test('should validate minimum length requirements', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Hi',
        description: 'Some longer content here to meet requirements'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title must be at least 5 characters');
  });
  
  test('should validate content field only when provided', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      // Only validate description if it's provided
      validationResult = result.current.validateForm({
        name: 'Valid Title',
        // description not provided - should pass
      });
    });
    
    expect(validationResult.isValid).toBe(true);
  });
  
  test('should validate content minimum length', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Valid Title',
        description: 'Short'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.description).toBe('Content must be at least 20 characters');
  });
  
  test('should strip HTML tags when validating content length', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Valid Title',
        description: '<p>Short</p>'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.description).toBe('Content must be at least 20 characters');
  });
  
  test('should accept valid HTML content with sufficient length', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Valid Title',
        description: '<p>This is a longer description that meets the minimum character requirements for validation</p>'
      });
    });
    
    expect(validationResult.isValid).toBe(true);
  });
  
  test('should allow valid form data', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Valid Product Title',
        description: 'This is a comprehensive description that exceeds the minimum character requirements'
      });
    });
    
    expect(validationResult.isValid).toBe(true);
    expect(validationResult.errors).toEqual({});
  });
  
  test('should validate multiple fields with multiple errors', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Hi',
        description: 'Short'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title must be at least 5 characters');
    expect(validationResult.errors.description).toBe('Content must be at least 20 characters');
  });
  
  test('should validate individual fields with validateField', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let isValid;
    act(() => {
      isValid = result.current.validateField('name', 'Test', [
        { type: 'required' },
        { type: 'minLength', length: 5 }
      ]);
    });
    
    expect(isValid).toBe(false);
    expect(result.current.errors.name).toBe('Must be at least 5 characters');
  });
  
  test('should validate email fields when present', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let isValid;
    act(() => {
      isValid = result.current.validateField('email', 'invalid-email', [
        { type: 'email' }
      ]);
    });
    
    expect(isValid).toBe(false);
    expect(result.current.errors.email).toBe('Invalid email address');
  });
  
  test('should accept valid email addresses', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let isValid;
    act(() => {
      isValid = result.current.validateField('email', 'redacted-email@example.com', [
        { type: 'email' }
      ]);
    });
    
    expect(isValid).toBe(true);
    expect(result.current.errors.email).toBeUndefined();
  });
  
  test('should allow setting custom errors manually', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      result.current.setErrors({
        name: 'Custom name error',
        description: 'Custom description error'
      });
    });
    
    expect(result.current.errors.name).toBe('Custom name error');
    expect(result.current.errors.description).toBe('Custom description error');
  });
  
  test('should clear all errors', () => {
    const { result } = renderHook(() => useFormValidation());
    
    // Set some errors first
    act(() => {
      result.current.setErrors({
        name: 'Name error',
        description: 'Description error'
      });
    });
    
    expect(result.current.errors.name).toBe('Name error');
    
    // Clear all errors
    act(() => {
      result.current.clearErrors();
    });
    
    expect(result.current.errors).toEqual({});
  });
  
  test('should clear specific field error via validateField', () => {
    const { result } = renderHook(() => useFormValidation());
    
    // Set initial errors
    act(() => {
      result.current.setErrors({
        name: 'Name error',
        description: 'Description error'
      });
    });
    
    // Validate field with valid value should clear its error
    act(() => {
      result.current.validateField('name', 'Valid Name', [
        { type: 'required' },
        { type: 'minLength', length: 5 }
      ]);
    });
    
    expect(result.current.errors.name).toBeUndefined();
    expect(result.current.errors.description).toBe('Description error');
  });
  
  test('should handle custom validation options', () => {
    const customConstants = {
      MIN_TITLE_LENGTH: 3,
      MIN_CONTENT_LENGTH: 15,
    };

    const { result } = renderHook(() => 
      useFormValidation({ constants: customConstants })
    );
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Hi',
        description: 'Short content'
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title must be at least 3 characters');
    expect(validationResult.errors.description).toBe('Content must be at least 15 characters');
  });
  
  test('should handle edge cases with null/undefined values', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: null,
        description: undefined
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title is required');
    // description is undefined, so it should not be validated
  });
  
  test('should handle complex HTML content validation', () => {
    const { result } = renderHook(() => useFormValidation());
    
    const complexHTML = `
      <div>
        <h1>Title</h1>
        <p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
        <ul>
          <li>Item 1</li>
          <li>Item 2</li>
        </ul>
      </div>
    `;
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: 'Valid Title',
        description: complexHTML
      });
    });
    
    expect(validationResult.isValid).toBe(true);
  });
  
  test('should validate form with missing fields gracefully', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let validationResult;
    act(() => {
      validationResult = result.current.validateForm({
        name: undefined
        // description missing entirely
      });
    });
    
    expect(validationResult.isValid).toBe(false);
    expect(validationResult.errors.name).toBe('Title is required');
    // description should not have an error since it's not provided
  });
  
  test('should maintain error state across multiple validations', () => {
    const { result } = renderHook(() => useFormValidation());
    
    // First validation with errors
    let validationResult1;
    act(() => {
      validationResult1 = result.current.validateForm({
        name: 'Hi',
      });
    });
    
    expect(validationResult1.isValid).toBe(false);
    expect(result.current.errors.name).toBe('Title must be at least 5 characters');
    
    // Second validation that clears errors
    let validationResult2;
    act(() => {
      validationResult2 = result.current.validateForm({
        name: 'Valid Title',
      });
    });
    
    expect(validationResult2.isValid).toBe(true);
    expect(result.current.errors).toEqual({});
  });
  
  test('should handle custom validator functions', () => {
    const { result } = renderHook(() => useFormValidation());
    
    const customValidator = (value, field) => {
      if (value && value.includes('banned')) {
        return { [field]: 'Contains banned content' };
      }
      return {};
    };
    
    let isValid;
    act(() => {
      isValid = result.current.validateField('name', 'banned word', [
        customValidator
      ]);
    });
    
    expect(isValid).toBe(false);
    expect(result.current.errors.name).toBe('Contains banned content');
  });
  
  test('should support HTML content validation with validateField', () => {
    const { result } = renderHook(() => useFormValidation());
    
    let isValid;
    act(() => {
      isValid = result.current.validateField('description', '<p>Short</p>', [
        { type: 'htmlContent', minLength: 20 }
      ]);
    });
    
    expect(isValid).toBe(false);
    expect(result.current.errors.description).toBe('Content must be at least 20 characters');
  });
  
  test('should mark fields as touched', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      result.current.markAsTouched('name');
    });
    
    expect(result.current.touched.name).toBe(true);
  });
  
  test('should provide access to validation constants', () => {
    const { result } = renderHook(() => useFormValidation());
    
    expect(result.current.constants).toEqual({
      MIN_TITLE_LENGTH: 5,
      MIN_CONTENT_LENGTH: 20,
      MIN_PASSWORD_LENGTH: 8,
      EMAIL_REGEX: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    });
  });
}); 
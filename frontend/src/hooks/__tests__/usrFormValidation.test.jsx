import { renderHook, act } from '@testing-library/react';
import { useFormValidation } from '../usrFormValidation';

describe('useFormValidation hook', () => {
  test('should initialize with empty errors', () => {
    const { result } = renderHook(() => useFormValidation());
    
    expect(result.current.errors).toEqual({});
  });
  
  test('should validate title - empty error', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: '',
        description: 'This is a valid description with more than 20 characters'
      });
      
      expect(validationResult.isValid).toBe(false);
      expect(validationResult.errors.name).toBe('Title is required');
    });
    
    // Errors should also be available in the state
    expect(result.current.errors.name).toBe('Title is required');
  });
  
  test('should validate title - too short error', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: 'Test', // Less than 5 characters
        description: 'This is a valid description with more than 20 characters'
      });
      
      expect(validationResult.isValid).toBe(false);
      expect(validationResult.errors.name).toBe('Title must be at least 5 characters');
    });
  });
  
  test('should validate content - empty error', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: 'This is a valid title',
        description: ''
      });
      
      expect(validationResult.isValid).toBe(false);
      expect(validationResult.errors.description).toBe('Content is required');
    });
    
    // Errors should also be available in the state
    expect(result.current.errors.description).toBe('Content is required');
  });
  
  test('should validate content - too short error', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: 'This is a valid title',
        description: 'Too short'
      });
      
      expect(validationResult.isValid).toBe(false);
      expect(validationResult.errors.description).toBe('Content must be at least 20 characters');
    });
  });
  
  test('should strip HTML tags when validating content', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: 'This is a valid title',
        description: '<p>Too short</p>' // HTML content that's below the minimum length
      });
      
      expect(validationResult.isValid).toBe(false);
      expect(validationResult.errors.description).toBe('Content must be at least 20 characters');
    });
  });
  
  test('should allow valid form data', () => {
    const { result } = renderHook(() => useFormValidation());
    
    act(() => {
      const validationResult = result.current.validateForm({
        name: 'This is a valid title',
        description: 'This is a valid description with more than 20 characters'
      });
      
      expect(validationResult.isValid).toBe(true);
      expect(validationResult.errors).toEqual({});
    });
    
    // Errors should be empty
    expect(result.current.errors).toEqual({});
  });
  
  test('should allow setting errors manually', () => {
    const { result } = renderHook(() => useFormValidation());
    
    const customErrors = {
      name: 'Custom error for name',
      description: 'Custom error for description'
    };
    
    act(() => {
      result.current.setErrors(customErrors);
    });
    
    expect(result.current.errors).toEqual(customErrors);
  });
}); 
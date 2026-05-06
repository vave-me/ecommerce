"use client";
import React, { createContext, useContext, useReducer } from 'react';

// Create the context
const FormStateContext = createContext();

// Define action types
const ACTION_TYPES = {
  UPDATE_FORM_DATA: 'UPDATE_FORM_DATA',
  RESET_FORM_DATA: 'RESET_FORM_DATA',
  SET_STEP: 'SET_STEP',
  SET_ERROR: 'SET_ERROR',
  CLEAR_ERROR: 'CLEAR_ERROR',
  SET_LOADING: 'SET_LOADING',
  SET_SUCCESS: 'SET_SUCCESS',
};

// Initial state
const initialState = {
  formData: {},
  currentStep: 1,
  lastCompletedStep: 0,
  errors: {},
  isLoading: false,
  isSuccess: false,
};

// Reducer function
function formReducer(state, action) {
  switch (action.type) {
    case ACTION_TYPES.UPDATE_FORM_DATA:
      return {
        ...state,
        formData: {
          ...state.formData,
          ...action.payload,
        },
      };
    case ACTION_TYPES.RESET_FORM_DATA:
      return {
        ...initialState,
        formData: action.payload || {},
      };
    case ACTION_TYPES.SET_STEP:
      return {
        ...state,
        currentStep: action.payload.step,
        lastCompletedStep: action.payload.isCompleted 
          ? Math.max(state.lastCompletedStep, action.payload.step)
          : state.lastCompletedStep,
      };
    case ACTION_TYPES.SET_ERROR:
      return {
        ...state,
        errors: {
          ...state.errors,
          ...action.payload,
        },
      };
    case ACTION_TYPES.CLEAR_ERROR:
      const newErrors = { ...state.errors };
      if (action.payload) {
        delete newErrors[action.payload];
      } else {
        return {
          ...state,
          errors: {},
        };
      }
      return {
        ...state,
        errors: newErrors,
      };
    case ACTION_TYPES.SET_LOADING:
      return {
        ...state,
        isLoading: action.payload,
      };
    case ACTION_TYPES.SET_SUCCESS:
      return {
        ...state,
        isSuccess: action.payload,
      };
    default:
      return state;
  }
}

// Provider component
export function FormStateProvider({ children, initialFormData = {} }) {
  const [state, dispatch] = useReducer(formReducer, {
    ...initialState,
    formData: initialFormData,
  });

  // Actions
  const updateFormData = (data) => {
    dispatch({ type: ACTION_TYPES.UPDATE_FORM_DATA, payload: data });
  };

  const resetFormData = (data) => {
    dispatch({ type: ACTION_TYPES.RESET_FORM_DATA, payload: data });
  };

  const setStep = (step, isCompleted = false) => {
    dispatch({ 
      type: ACTION_TYPES.SET_STEP, 
      payload: { step, isCompleted } 
    });
  };

  const setError = (error) => {
    dispatch({ type: ACTION_TYPES.SET_ERROR, payload: error });
  };

  const clearError = (errorKey) => {
    dispatch({ type: ACTION_TYPES.CLEAR_ERROR, payload: errorKey });
  };

  const setLoading = (isLoading) => {
    dispatch({ type: ACTION_TYPES.SET_LOADING, payload: isLoading });
  };

  const setSuccess = (isSuccess) => {
    dispatch({ type: ACTION_TYPES.SET_SUCCESS, payload: isSuccess });
  };

  // Next and previous step helpers
  const goToNextStep = (isCurrentStepComplete = true) => {
    setStep(state.currentStep + 1, isCurrentStepComplete);
  };

  const goToPreviousStep = () => {
    setStep(Math.max(1, state.currentStep - 1));
  };

  // Value to expose to consumers
  const value = {
    formData: state.formData,
    currentStep: state.currentStep,
    lastCompletedStep: state.lastCompletedStep,
    errors: state.errors,
    isLoading: state.isLoading,
    isSuccess: state.isSuccess,
    updateFormData,
    resetFormData,
    setStep,
    goToNextStep,
    goToPreviousStep,
    setError,
    clearError,
    setLoading,
    setSuccess,
  };

  return (
    <FormStateContext.Provider value={value}>
      {children}
    </FormStateContext.Provider>
  );
}

// Custom hook to use the context
export function useFormState() {
  const context = useContext(FormStateContext);
  if (context === undefined) {
    throw new Error('useFormState must be used within a FormStateProvider');
  }
  return context;
} 
import React from 'react';
import { render, screen, waitForElementToBeRemoved } from '@testing-library/react';
import { renderWithProviders } from './test-utils';

/**
 * Renders a component with all providers and waits for any loading states to resolve
 * This is particularly useful for modal tests where we need to wait for async data fetching
 * @param {React.ReactElement} ui - The component to render
 * @param {Object} options - Additional render options
 * @returns {Promise<Object>} - The render result
 */
export async function renderModalWithLoading(ui, options = {}) {
  const result = renderWithProviders(ui, options);
  
  // Wait for any loading spinners to disappear
  try {
    // Check if there's a loading state to wait for
    const loadingElement = screen.queryByText('Loading...');
    if (loadingElement) {
      await waitForElementToBeRemoved(() => screen.queryByText('Loading...'));
    }
  } catch (error) {
    console.warn('Loading state not found or failed to resolve:', error);
  }
  
  return result;
}

/**
 * Waits for a step to be visible in a multi-step form modal
 * @param {string} stepTestId - The test ID of the step to wait for
 * @returns {Promise<HTMLElement>} - The step element
 */
export async function waitForStep(stepTestId) {
  let element;
  try {
    // Wait for a short time to allow for any async state changes
    await new Promise(resolve => setTimeout(resolve, 0));
    
    // Then check for the element
    element = await screen.findByTestId(stepTestId);
  } catch (error) {
    throw new Error(`Step ${stepTestId} was not found: ${error.message}`);
  }
  
  return element;
} 
import React from 'react';
import { render, screen } from '@testing-library/react';
import ErrorBoundary from '../ErrorBoundary';

// Create a component that throws an error
const ErrorComponent = () => {
  throw new Error('Test error');
};

describe('ErrorBoundary', () => {
  // Properly mock console.error before all tests
  const originalConsoleError = console.error;

  beforeAll(() => {
    console.error = jest.fn();
  });

  afterAll(() => {
    console.error = originalConsoleError;
  });

  it('renders children when there is no error', () => {
    render(
      <ErrorBoundary>
        <div data-testid="child">Test Child</div>
      </ErrorBoundary>
    );

    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('renders fallback UI when children throw an error', () => {
    render(
      <ErrorBoundary>
        <ErrorComponent />
      </ErrorBoundary>
    );

    // Check for error message in the fallback UI
    expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();
  });

  it('displays the error message when debug mode is enabled', () => {
    render(
      <ErrorBoundary debug={true}>
        <ErrorComponent />
      </ErrorBoundary>
    );

    // Check for detailed error message
    expect(screen.getByText(/Test error/i)).toBeInTheDocument();
  });

  it('calls onError when provided', () => {
    const handleError = jest.fn();
    
    render(
      <ErrorBoundary onError={handleError}>
        <ErrorComponent />
      </ErrorBoundary>
    );

    // Check that handleError was called with the error
    expect(handleError).toHaveBeenCalledWith(expect.any(Error));
    expect(handleError.mock.calls[0][0].message).toBe('Test error');
  });
}); 
import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { useFocusTrap } from '../useFocusTrap';

// Test component that uses the hook
const TestComponent = ({ isActive }) => {
  const ref = useFocusTrap(isActive);
  
  return (
    <div ref={ref} data-testid="trap-container">
      <button data-testid="button-1">First Button</button>
      <input data-testid="input" type="text" />
      <button data-testid="button-2">Last Button</button>
    </div>
  );
};

describe('useFocusTrap hook', () => {
  // Store original implementation
  const originalFocus = HTMLElement.prototype.focus;
  
  beforeEach(() => {
    // Mock focus method for all HTML elements
    HTMLElement.prototype.focus = jest.fn();
    
    // Mock document.activeElement
    Object.defineProperty(document, 'activeElement', {
      writable: true,
      value: document.body
    });
  });
  
  afterEach(() => {
    // Restore original focus method
    HTMLElement.prototype.focus = originalFocus;
  });
  
  test('should not attempt to trap focus when inactive', () => {
    render(<TestComponent isActive={false} />);
    
    // Focus should not be called on inactive trap
    expect(HTMLElement.prototype.focus).not.toHaveBeenCalled();
  });
  
  test('should focus first focusable element when activated', () => {
    const { getByTestId } = render(<TestComponent isActive={true} />);
    const firstButton = getByTestId('button-1');
    
    // First element should receive focus
    expect(HTMLElement.prototype.focus).toHaveBeenCalled();
    expect(firstButton.focus).toHaveBeenCalled();
  });
  
  test('should trap Tab key press and move to first element when last element has focus', () => {
    const { getByTestId } = render(<TestComponent isActive={true} />);
    
    const lastButton = getByTestId('button-2');
    
    // Simulate last button having focus
    Object.defineProperty(document, 'activeElement', {
      writable: true,
      value: lastButton
    });
    
    // Simulate Tab key press on last element
    fireEvent.keyDown(getByTestId('trap-container'), { key: 'Tab', keyCode: 9 });
    
    const firstButton = getByTestId('button-1');
    expect(firstButton.focus).toHaveBeenCalled();
  });
  
  test('should trap Shift+Tab key press and move to last element when first element has focus', () => {
    const { getByTestId } = render(<TestComponent isActive={true} />);
    
    const firstButton = getByTestId('button-1');
    
    // Simulate first button having focus
    Object.defineProperty(document, 'activeElement', {
      writable: true,
      value: firstButton
    });
    
    // Simulate Shift+Tab key press on first element
    fireEvent.keyDown(getByTestId('trap-container'), { 
      key: 'Tab', 
      keyCode: 9,
      shiftKey: true 
    });
    
    const lastButton = getByTestId('button-2');
    expect(lastButton.focus).toHaveBeenCalled();
  });
  
  test('should not interfere with other keys', () => {
    const { getByTestId } = render(<TestComponent isActive={true} />);
    
    // Clear mock calls from initial focus
    jest.clearAllMocks();
    
    // Simulate pressing a different key
    fireEvent.keyDown(getByTestId('trap-container'), { key: 'Enter', keyCode: 13 });
    
    // No focus change should occur
    expect(HTMLElement.prototype.focus).not.toHaveBeenCalled();
  });
  
  test('should restore focus to previously focused element on cleanup', () => {
    // Create an element to have previous focus
    const previouslyFocused = document.createElement('button');
    document.body.appendChild(previouslyFocused);
    
    // Set it as the active element
    Object.defineProperty(document, 'activeElement', {
      writable: true,
      value: previouslyFocused
    });
    
    // Render component then unmount
    const { unmount } = render(<TestComponent isActive={true} />);
    
    // Clear focus mock calls
    jest.clearAllMocks();
    
    // Trigger cleanup
    unmount();
    
    // Previous element should have focus restored
    expect(previouslyFocused.focus).toHaveBeenCalled();
    
    // Clean up
    document.body.removeChild(previouslyFocused);
  });
}); 
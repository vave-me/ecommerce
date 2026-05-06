import { renderHook, act } from '@testing-library/react';
import { useMobileDetection } from '../useMobileDetection';

describe('useMobileDetection hook', () => {
  const originalInnerWidth = window.innerWidth;
  let resizeEventListener;

  // Mock window resize event
  beforeEach(() => {
    resizeEventListener = null;
    
    // Mock addEventListener to capture the resize handler
    jest.spyOn(window, 'addEventListener').mockImplementation((event, handler) => {
      if (event === 'resize') {
        resizeEventListener = handler;
      }
    });
    
    // Mock removeEventListener
    jest.spyOn(window, 'removeEventListener').mockImplementation(() => {});
  });

  afterEach(() => {
    // Restore original window properties
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: originalInnerWidth
    });
    
    jest.restoreAllMocks();
  });

  test('should detect mobile when screen width <= 768px', () => {
    // Set window width to mobile
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768
    });
    
    const { result } = renderHook(() => useMobileDetection());
    
    expect(result.current).toBe(true);
  });

  test('should detect desktop when screen width > 768px', () => {
    // Set window width to desktop
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024
    });
    
    const { result } = renderHook(() => useMobileDetection());
    
    expect(result.current).toBe(false);
  });

  test('should update when window is resized', () => {
    // Start with desktop size
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024
    });
    
    const { result } = renderHook(() => useMobileDetection());
    
    // Initially should be desktop
    expect(result.current).toBe(false);
    
    // Change to mobile size and trigger resize event
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 600
      });
      
      // Call the resize handler directly
      if (resizeEventListener) {
        resizeEventListener();
      }
    });
    
    // Should update to mobile
    expect(result.current).toBe(true);
  });

  test('should set up and clean up event listeners', () => {
    const { unmount } = renderHook(() => useMobileDetection());
    
    // Should add resize listener
    expect(window.addEventListener).toHaveBeenCalledWith('resize', expect.any(Function));
    
    // Unmount to trigger cleanup
    unmount();
    
    // Should remove resize listener
    expect(window.removeEventListener).toHaveBeenCalledWith('resize', expect.any(Function));
  });
}); 
import { renderHook, act } from '@testing-library/react';
import { useIsMobile } from '../useMobileDetection';

describe('useIsMobile hook', () => {
  const originalInnerWidth = window.innerWidth;
  let resizeEventListener;

  // Mock addEventListener to capture the resize listener
  beforeEach(() => {
    resizeEventListener = jest.fn();
    jest.spyOn(window, 'addEventListener').mockImplementation((event, listener) => {
      if (event === 'resize') {
        resizeEventListener = listener;
      }
    });
    
    // Mock removeEventListener
    jest.spyOn(window, 'removeEventListener').mockImplementation(() => {});
  });

  afterEach(() => {
    // Restore original window.innerWidth
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: originalInnerWidth
    });
    
    // Restore original addEventListener and removeEventListener
    jest.restoreAllMocks();
  });

  test('should default to non-mobile initially', () => {
    const { result } = renderHook(() => useIsMobile());
    expect(result.current).toBe(false);
  });

  test('should detect mobile when screen width is below breakpoint', () => {
    // Set window width to be below the default breakpoint (768px)
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 500
    });
    
    const { result } = renderHook(() => useIsMobile());
    
    // Force resize event to trigger
    act(() => {
      resizeEventListener();
    });
    
    expect(result.current).toBe(true);
  });

  test('should detect non-mobile when screen width is above breakpoint', () => {
    // Set window width to be above the default breakpoint (768px)
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1000
    });
    
    const { result } = renderHook(() => useIsMobile());
    
    // Force resize event to trigger
    act(() => {
      resizeEventListener();
    });
    
    expect(result.current).toBe(false);
  });

  test('should respect custom breakpoint', () => {
    // Set window width to 900px
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 900
    });
    
    // Use custom breakpoint of 1000px
    const { result } = renderHook(() => useIsMobile(1000));
    
    // Force resize event to trigger
    act(() => {
      resizeEventListener();
    });
    
    // 900px is less than 1000px, so should be considered mobile
    expect(result.current).toBe(true);
  });

  test('should update when window size changes', () => {
    // Start with desktop size
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1000
    });
    
    const { result } = renderHook(() => useIsMobile());
    
    // Force initial resize event
    act(() => {
      resizeEventListener();
    });
    
    expect(result.current).toBe(false);
    
    // Change to mobile size
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 500
    });
    
    // Trigger resize
    act(() => {
      resizeEventListener();
    });
    
    expect(result.current).toBe(true);
  });

  test('should clean up resize listener on unmount', () => {
    const { unmount } = renderHook(() => useIsMobile());
    
    unmount();
    
    // Should have called removeEventListener for 'resize'
    expect(window.removeEventListener).toHaveBeenCalledWith('resize', expect.any(Function));
  });
}); 
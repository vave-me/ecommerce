import { renderHook, act } from '@testing-library/react';
import { useMobileDetection } from '@/hooks/useMediaQuery.jsx';

describe('useMobileDetection hook', () => {
  beforeEach(() => {
    // Mock window dimensions
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    });
    
    // Mock addEventListener and removeEventListener  
    window.addEventListener = jest.fn();
    window.removeEventListener = jest.fn();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('should return false for desktop width', () => {
    const { result } = renderHook(() => useMobileDetection());
    expect(result.current).toBe(false);
  });

  test('should return true for mobile width', () => {
    // Set mobile width
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 600,
    });

    const { result } = renderHook(() => useMobileDetection());
    
    act(() => {
      // Trigger resize event to update state
      window.dispatchEvent(new Event('resize'));
    });

    expect(result.current).toBe(true);
  });

  test('should update when window is resized', () => {
    const { result } = renderHook(() => useMobileDetection());
    
    // Initially false (desktop)
    expect(result.current).toBe(false);
    
    // Change to mobile width
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 500,
      });
      window.dispatchEvent(new Event('resize'));
    });
    
    expect(result.current).toBe(true);
    
    // Change back to desktop width
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 1200,
      });
      window.dispatchEvent(new Event('resize'));
    });
    
    expect(result.current).toBe(false);
  });

  test('should handle boundary conditions', () => {
    // Test exactly at mobile breakpoint (768px)
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768,
    });

    const { result } = renderHook(() => useMobileDetection());
    
    act(() => {
      window.dispatchEvent(new Event('resize'));
    });

    expect(result.current).toBe(true); // 768 <= 768 should be mobile

    // Test just above mobile breakpoint
    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 769,
      });
      window.dispatchEvent(new Event('resize'));
    });

    expect(result.current).toBe(false); // 769 > 768 should not be mobile
  });

  test('should clean up event listeners on unmount', () => {
    const { unmount } = renderHook(() => useMobileDetection());
    
    expect(window.addEventListener).toHaveBeenCalledWith('resize', expect.any(Function), { passive: true });
    
    unmount();
    
    expect(window.removeEventListener).toHaveBeenCalledWith('resize', expect.any(Function));
  });
}); 
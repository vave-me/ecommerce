import { renderHook, act } from '@testing-library/react';
import { useHeaderScroll } from '../useHeaderScroll';

describe('useHeaderScroll hook', () => {
  let scrollEventListener;
  const originalScrollY = window.scrollY;

  // Mock window.scrollY and event listeners
  beforeEach(() => {
    // Set initial scroll position
    Object.defineProperty(window, 'scrollY', {
      value: 0,
      writable: true
    });

    // Capture the scroll event listener
    scrollEventListener = jest.fn();
    jest.spyOn(window, 'addEventListener').mockImplementation((event, listener, options) => {
      if (event === 'scroll') {
        scrollEventListener = listener;
      }
    });

    // Mock removeEventListener
    jest.spyOn(window, 'removeEventListener').mockImplementation(() => {});
  });

  afterEach(() => {
    // Restore original properties and mocks
    Object.defineProperty(window, 'scrollY', {
      value: originalScrollY
    });
    jest.restoreAllMocks();
  });

  test('should default to not scrolled initially with 0 scrollY', () => {
    // scrollY is 0 from beforeEach
    const { result } = renderHook(() => useHeaderScroll({ offset: 20 }));
    
    // Initial check should show not scrolled
    expect(result.current).toBe(false);
  });

  test('should set isScrolled to true when scrolled beyond offset', () => {
    // Set up initial render
    const { result } = renderHook(() => useHeaderScroll({ offset: 20 }));
    expect(result.current).toBe(false);
    
    // Simulate scrolling down to 30px (beyond the 20px offset)
    Object.defineProperty(window, 'scrollY', { value: 30 });
    
    // Trigger scroll event
    act(() => {
      scrollEventListener();
    });
    
    // Should now be considered scrolled
    expect(result.current).toBe(true);
  });

  test('should set isScrolled to false when scrolled back above offset', () => {
    // Start with scrolled position
    Object.defineProperty(window, 'scrollY', { value: 50 });
    
    const { result } = renderHook(() => useHeaderScroll({ offset: 20 }));
    expect(result.current).toBe(true);
    
    // Scroll back to top
    Object.defineProperty(window, 'scrollY', { value: 10 });
    
    // Trigger scroll event
    act(() => {
      scrollEventListener();
    });
    
    // Should no longer be considered scrolled
    expect(result.current).toBe(false);
  });

  test('should respect custom offset', () => {
    // Scroll to 40px
    Object.defineProperty(window, 'scrollY', { value: 40 });
    
    // Using higher offset of a custom 50px
    const { result } = renderHook(() => useHeaderScroll({ offset: 50 }));
    
    // Trigger initial scroll check
    act(() => {
      scrollEventListener();
    });
    
    // Should not be considered scrolled yet (40 < 50)
    expect(result.current).toBe(false);
    
    // Now scroll to 60px
    Object.defineProperty(window, 'scrollY', { value: 60 });
    
    // Trigger scroll event
    act(() => {
      scrollEventListener();
    });
    
    // Should now be considered scrolled (60 > 50)
    expect(result.current).toBe(true);
  });

  test('should use default offset of 20 when not specified', () => {
    // Set to 15px (below default 20px)
    Object.defineProperty(window, 'scrollY', { value: 15 });
    
    const { result } = renderHook(() => useHeaderScroll({}));
    
    // Trigger scroll check
    act(() => {
      scrollEventListener();
    });
    
    // Should not be scrolled yet
    expect(result.current).toBe(false);
    
    // Set to 25px (above default 20px)
    Object.defineProperty(window, 'scrollY', { value: 25 });
    
    // Trigger scroll event
    act(() => {
      scrollEventListener();
    });
    
    // Should now be scrolled
    expect(result.current).toBe(true);
  });

  test('should clean up scroll listener on unmount', () => {
    const { unmount } = renderHook(() => useHeaderScroll({ offset: 20 }));
    
    unmount();
    
    // Should have called removeEventListener for 'scroll'
    expect(window.removeEventListener).toHaveBeenCalledWith('scroll', expect.any(Function));
  });
}); 
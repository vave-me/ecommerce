import { renderHook, act } from '@testing-library/react';
import { useMediaQuery } from '../useMobileDetection';

describe('useMediaQuery hook', () => {
  // Store original implementations
  let originalMatchMedia;
  let originalAddEventListener;
  let originalRemoveEventListener;
  
  // Mock media query object
  const createMatchMediaMock = (matches) => {
    return jest.fn().mockImplementation(query => ({
      matches,
      media: query,
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn()
    }));
  };

  beforeEach(() => {
    // Store original implementations
    originalMatchMedia = window.matchMedia;
    originalAddEventListener = Element.prototype.addEventListener;
    originalRemoveEventListener = Element.prototype.removeEventListener;
    
    // Default mock - not matching
    window.matchMedia = createMatchMediaMock(false);
  });

  afterEach(() => {
    // Restore original implementations
    window.matchMedia = originalMatchMedia;
    Element.prototype.addEventListener = originalAddEventListener;
    Element.prototype.removeEventListener = originalRemoveEventListener;
  });

  test('should default to non-matching state', () => {
    window.matchMedia = createMatchMediaMock(false);
    
    const { result } = renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    expect(result.current).toBe(false);
  });

  test('should return true when media query matches', () => {
    window.matchMedia = createMatchMediaMock(true);
    
    const { result } = renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    expect(result.current).toBe(true);
  });

  test('should handle media query changes', () => {
    // Start with non-matching
    const matchMediaMock = {
      matches: false,
      media: '(max-width: 768px)',
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn()
    };
    
    window.matchMedia = jest.fn().mockReturnValue(matchMediaMock);
    
    const { result } = renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    // Initial state should be false
    expect(result.current).toBe(false);
    
    // Simulate media query match change
    act(() => {
      matchMediaMock.matches = true;
      // Call the registered event handler
      const eventHandler = matchMediaMock.addEventListener.mock.calls[0][1];
      eventHandler();
    });
    
    // Should update to true
    expect(result.current).toBe(true);
  });

  test('should use addEventListener when available', () => {
    const matchMediaMock = {
      matches: false,
      media: '(max-width: 768px)',
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn()
    };
    
    window.matchMedia = jest.fn().mockReturnValue(matchMediaMock);
    
    renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    // Should use addEventListener
    expect(matchMediaMock.addEventListener).toHaveBeenCalledWith('change', expect.any(Function));
    expect(matchMediaMock.addListener).not.toHaveBeenCalled();
  });

  test('should fall back to addListener when addEventListener throws', () => {
    const matchMediaMock = {
      matches: false,
      media: '(max-width: 768px)',
      addEventListener: jest.fn().mockImplementation(() => {
        throw new Error('addEventListener not supported');
      }),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn()
    };
    
    window.matchMedia = jest.fn().mockReturnValue(matchMediaMock);
    
    renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    // Should fall back to addListener when addEventListener throws
    expect(matchMediaMock.addListener).toHaveBeenCalledWith(expect.any(Function));
  });

  test('should clean up listeners on unmount with removeEventListener', () => {
    const matchMediaMock = {
      matches: false,
      media: '(max-width: 768px)',
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn()
    };
    
    window.matchMedia = jest.fn().mockReturnValue(matchMediaMock);
    
    const { unmount } = renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    // Unmount the hook
    unmount();
    
    // Should clean up with removeEventListener
    expect(matchMediaMock.removeEventListener).toHaveBeenCalledWith('change', expect.any(Function));
    expect(matchMediaMock.removeListener).not.toHaveBeenCalled();
  });

  test('should fall back to removeListener when removeEventListener throws', () => {
    const matchMediaMock = {
      matches: false,
      media: '(max-width: 768px)',
      addEventListener: jest.fn(),
      removeEventListener: jest.fn().mockImplementation(() => {
        throw new Error('removeEventListener not supported');
      }),
      addListener: jest.fn(),
      removeListener: jest.fn()
    };
    
    window.matchMedia = jest.fn().mockReturnValue(matchMediaMock);
    
    const { unmount } = renderHook(() => useMediaQuery('(max-width: 768px)'));
    
    // Unmount the hook
    unmount();
    
    // Should fall back to removeListener when removeEventListener throws
    expect(matchMediaMock.removeListener).toHaveBeenCalledWith(expect.any(Function));
  });
}); 
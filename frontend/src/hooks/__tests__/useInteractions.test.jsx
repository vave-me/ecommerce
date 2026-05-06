import { renderHook, act } from '@testing-library/react';
import { useInteractions } from '../useInteractions';

describe('useInteractions hook', () => {
  test('should initialize with default values', () => {
    const { result } = renderHook(() => useInteractions());
    
    expect(result.current.currentMediaIndex).toBe(0);
    expect(result.current.activeVideoIndex).toBeNull();
    expect(result.current.isProcessing).toBe(false);
  });
  
  test('should update currentMediaIndex', () => {
    const { result } = renderHook(() => useInteractions());
    
    act(() => {
      result.current.setCurrentMediaIndex(2);
    });
    
    expect(result.current.currentMediaIndex).toBe(2);
  });
  
  test('should update activeVideoIndex', () => {
    const { result } = renderHook(() => useInteractions());
    
    act(() => {
      result.current.setActiveVideoIndex(1);
    });
    
    expect(result.current.activeVideoIndex).toBe(1);
  });
  
  test('should update isProcessing', () => {
    const { result } = renderHook(() => useInteractions());
    
    act(() => {
      result.current.setIsProcessing(true);
    });
    
    expect(result.current.isProcessing).toBe(true);
  });
  
  test('handleMediaNavigation should navigate to next media', () => {
    const { result } = renderHook(() => useInteractions());
    const galleryLength = 3;
    
    // Navigate to next (0 -> 1)
    act(() => {
      result.current.handleMediaNavigation('next', galleryLength);
    });
    
    expect(result.current.currentMediaIndex).toBe(1);
    expect(result.current.activeVideoIndex).toBeNull();
    
    // Navigate to next again (1 -> 2)
    act(() => {
      result.current.handleMediaNavigation('next', galleryLength);
    });
    
    expect(result.current.currentMediaIndex).toBe(2);
    
    // Navigate to next again, should wrap around (2 -> 0)
    act(() => {
      result.current.handleMediaNavigation('next', galleryLength);
    });
    
    expect(result.current.currentMediaIndex).toBe(0);
  });
  
  test('handleMediaNavigation should navigate to previous media', () => {
    const { result } = renderHook(() => useInteractions());
    const galleryLength = 3;
    
    // Set initial index to 1
    act(() => {
      result.current.setCurrentMediaIndex(1);
    });
    
    // Navigate to previous (1 -> 0)
    act(() => {
      result.current.handleMediaNavigation('prev', galleryLength);
    });
    
    expect(result.current.currentMediaIndex).toBe(0);
    
    // Navigate to previous again, should wrap around (0 -> 2)
    act(() => {
      result.current.handleMediaNavigation('prev', galleryLength);
    });
    
    expect(result.current.currentMediaIndex).toBe(2);
  });
  
  test('handleMediaNavigation should reset active video index', () => {
    const { result } = renderHook(() => useInteractions());
    const galleryLength = 3;
    
    // Set active video index
    act(() => {
      result.current.setActiveVideoIndex(0);
    });
    
    expect(result.current.activeVideoIndex).toBe(0);
    
    // Navigate - should reset activeVideoIndex
    act(() => {
      result.current.handleMediaNavigation('next', galleryLength);
    });
    
    expect(result.current.activeVideoIndex).toBeNull();
  });
  
  test('handleMediaNavigation should do nothing if gallery length is 0', () => {
    const { result } = renderHook(() => useInteractions());
    
    // Try to navigate with empty gallery
    act(() => {
      result.current.handleMediaNavigation('next', 0);
    });
    
    // Index should remain unchanged
    expect(result.current.currentMediaIndex).toBe(0);
  });
}); 
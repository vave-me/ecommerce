import { renderHook } from '@testing-library/react';
import { useBodyLock } from '../useBodyLock';

describe('useBodyLock hook', () => {
  beforeEach(() => {
    // Reset body style before each test
    document.body.style.overflow = '';
  });

  test('should set overflow to hidden when isLocked is true', () => {
    renderHook(() => useBodyLock(true));
    expect(document.body.style.overflow).toBe('hidden');
  });

  test('should not set overflow when isLocked is false', () => {
    renderHook(() => useBodyLock(false));
    expect(document.body.style.overflow).toBe('');
  });

  test('should reset overflow when unmounted', () => {
    // First lock the body
    const { unmount } = renderHook(() => useBodyLock(true));
    expect(document.body.style.overflow).toBe('hidden');
    
    // Then unmount the component
    unmount();
    expect(document.body.style.overflow).toBe('');
  });

  test('should update overflow when isLocked changes', () => {
    const { rerender } = renderHook(({ isLocked }) => useBodyLock(isLocked), {
      initialProps: { isLocked: false }
    });
    
    // Initially not locked
    expect(document.body.style.overflow).toBe('');
    
    // Update to locked
    rerender({ isLocked: true });
    expect(document.body.style.overflow).toBe('hidden');
    
    // Update back to unlocked
    rerender({ isLocked: false });
    expect(document.body.style.overflow).toBe('');
  });
}); 
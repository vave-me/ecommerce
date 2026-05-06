import { debounce } from '../debounce.js';

describe('Debounce Utility', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  test('should debounce function calls', () => {
    const mockFn = jest.fn();
    const debouncedFn = debounce(mockFn, 500);

    // Call function multiple times
    debouncedFn('test1');
    debouncedFn('test2');
    debouncedFn('test3');

    // Function should not have been called yet
    expect(mockFn).not.toHaveBeenCalled();

    // Fast-forward time by 500ms
    jest.advanceTimersByTime(500);

    // Function should have been called exactly once with the last arguments
    expect(mockFn).toHaveBeenCalledTimes(1);
    expect(mockFn).toHaveBeenCalledWith('test3');
  });

  test('should reset timer on subsequent calls', () => {
    const mockFn = jest.fn();
    const debouncedFn = debounce(mockFn, 500);

    // Call once
    debouncedFn('test1');

    // Advance time but not fully
    jest.advanceTimersByTime(400);
    
    // Function should not have been called yet
    expect(mockFn).not.toHaveBeenCalled();

    // Call again, resetting the timer
    debouncedFn('test2');

    // Advance time but not enough for the second call
    jest.advanceTimersByTime(400);
    
    // Function should still not have been called
    expect(mockFn).not.toHaveBeenCalled();

    // Advance time enough to trigger the second call
    jest.advanceTimersByTime(100);
    
    // Function should have been called exactly once with the last arguments
    expect(mockFn).toHaveBeenCalledTimes(1);
    expect(mockFn).toHaveBeenCalledWith('test2');
  });

  test('should preserve the last call arguments', () => {
    const mockFn = jest.fn();
    const debouncedFn = debounce(mockFn, 500);

    // Call with multiple arguments
    debouncedFn('test1', 123, { key: 'value' });
    debouncedFn('test2', 456, { key: 'updated' });

    // Advance time to trigger the call
    jest.advanceTimersByTime(500);

    // Function should have been called with the latest arguments
    expect(mockFn).toHaveBeenCalledTimes(1);
    expect(mockFn).toHaveBeenCalledWith('test2', 456, { key: 'updated' });
  });

  test('should allow multiple independent debounced functions', () => {
    const mockFn1 = jest.fn();
    const mockFn2 = jest.fn();
    
    const debouncedFn1 = debounce(mockFn1, 500);
    const debouncedFn2 = debounce(mockFn2, 1000);

    // Call both functions
    debouncedFn1('fn1');
    debouncedFn2('fn2');

    // Advance time enough for first function only
    jest.advanceTimersByTime(500);

    // Only first function should have been called
    expect(mockFn1).toHaveBeenCalledTimes(1);
    expect(mockFn2).not.toHaveBeenCalled();

    // Advance time for second function
    jest.advanceTimersByTime(500);

    // Both functions should have been called
    expect(mockFn1).toHaveBeenCalledTimes(1);
    expect(mockFn2).toHaveBeenCalledTimes(1);
  });
}); 
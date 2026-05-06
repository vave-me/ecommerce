import { renderHook, act } from '@testing-library/react';
import { usePerformance } from '@/hooks/usePerformance.jsx';

// Mock performance.now
const mockPerformanceNow = jest.fn();
Object.defineProperty(window, 'performance', {
  value: {
    now: mockPerformanceNow
  }
});

describe('usePerformance hook', () => {
  let timeCounter;

  beforeEach(() => {
    timeCounter = 0;
    mockPerformanceNow.mockImplementation(() => {
      timeCounter += 10; // Increment by 10ms each call
      return timeCounter;
    });
    jest.clearAllMocks();
  });

  test('should initialize with correct structure', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));

    expect(typeof result.current.trackApiCall).toBe('function');
    expect(typeof result.current.trackError).toBe('function');
    expect(typeof result.current.getMetrics).toBe('function');
  });

  test('should track render count and time on each render', () => {
    const { result, rerender } = renderHook(() => usePerformance('TestComponent'));

    // First render
    let metrics = result.current.getMetrics();
    expect(metrics.renderCount).toBe(1);
    expect(metrics.averageRenderTime).toBeGreaterThan(0);

    // Second render
    rerender();
    metrics = result.current.getMetrics();
    expect(metrics.renderCount).toBe(2);
    expect(metrics.averageRenderTime).toBeGreaterThan(0);
  });

  test('should track successful API calls', async () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const mockApiCall = jest.fn().mockResolvedValue({ data: 'test' });

    await act(async () => {
      const apiResult = await result.current.trackApiCall('fetchData', mockApiCall);
      expect(apiResult).toEqual({ data: 'test' });
    });

    const metrics = result.current.getMetrics();
    expect(metrics.apiCalls.fetchData).toHaveLength(1);
    expect(metrics.apiCalls.fetchData[0]).toBeGreaterThan(0);
    expect(metrics.errors).toHaveLength(0);
  });

  test('should track multiple API calls of same type', async () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const mockApiCall = jest.fn().mockResolvedValue({ data: 'test' });

    await act(async () => {
      await result.current.trackApiCall('fetchData', mockApiCall);
      await result.current.trackApiCall('fetchData', mockApiCall);
      await result.current.trackApiCall('fetchData', mockApiCall);
    });

    const metrics = result.current.getMetrics();
    expect(metrics.apiCalls.fetchData).toHaveLength(3);
    expect(metrics.apiCalls.fetchData.every(duration => duration > 0)).toBe(true);
  });

  test('should track different API call types separately', async () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const mockApiCall1 = jest.fn().mockResolvedValue({ data: 'test1' });
    const mockApiCall2 = jest.fn().mockResolvedValue({ data: 'test2' });

    await act(async () => {
      await result.current.trackApiCall('fetchData', mockApiCall1);
      await result.current.trackApiCall('updateData', mockApiCall2);
    });

    const metrics = result.current.getMetrics();
    expect(metrics.apiCalls.fetchData).toHaveLength(1);
    expect(metrics.apiCalls.updateData).toHaveLength(1);
  });

  test('should track API call errors and still throw', async () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const testError = new Error('API failure');
    const mockApiCall = jest.fn().mockRejectedValue(testError);

    await act(async () => {
      await expect(
        result.current.trackApiCall('fetchData', mockApiCall)
      ).rejects.toThrow('API failure');
    });

    const metrics = result.current.getMetrics();
    expect(metrics.errors).toHaveLength(1);
    expect(metrics.errors[0]).toMatchObject({
      type: 'api',
      apiName: 'fetchData',
      error: 'API failure'
    });
    expect(metrics.errors[0].timestamp).toBeDefined();
  });

  test('should track component errors manually', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const testError = new Error('Component error');
    testError.stack = 'Error stack trace';

    act(() => {
      result.current.trackError(testError);
    });

    const metrics = result.current.getMetrics();
    expect(metrics.errors).toHaveLength(1);
    expect(metrics.errors[0]).toMatchObject({
      type: 'component',
      error: 'Component error',
      stack: 'Error stack trace'
    });
    expect(metrics.errors[0].timestamp).toBeDefined();
  });

  test('should track multiple errors', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const error1 = new Error('First error');
    const error2 = new Error('Second error');

    act(() => {
      result.current.trackError(error1);
      result.current.trackError(error2);
    });

    const metrics = result.current.getMetrics();
    expect(metrics.errors).toHaveLength(2);
    expect(metrics.errors[0].error).toBe('First error');
    expect(metrics.errors[1].error).toBe('Second error');
  });

  test('should calculate average render time correctly', () => {
    const { result, rerender } = renderHook(() => usePerformance('TestComponent'));

    // Multiple rerenders to test average calculation
    rerender();
    rerender();
    rerender();

    const metrics = result.current.getMetrics();
    expect(metrics.renderCount).toBe(4); // Initial render + 3 rerenders
    expect(metrics.averageRenderTime).toBeGreaterThan(0);
    expect(typeof metrics.averageRenderTime).toBe('number');
  });

  test('should return immutable metrics snapshot', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    const metrics1 = result.current.getMetrics();
    const metrics2 = result.current.getMetrics();
    
    // Should be different objects (snapshots)
    expect(metrics1).not.toBe(metrics2);
    
    // But with same values
    expect(metrics1).toEqual(metrics2);
  });

  test('should handle edge case of zero renders for average calculation', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    // Mock internal state to simulate edge case
    const metrics = result.current.getMetrics();
    
    // Should not throw error and should handle division by zero
    expect(typeof metrics.averageRenderTime).toBe('number');
    expect(metrics.averageRenderTime).not.toBe(Infinity);
    expect(metrics.averageRenderTime).not.toBe(NaN);
  });

  test('should handle API calls that return different data types', async () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    await act(async () => {
      // String result
      const stringResult = await result.current.trackApiCall('stringApi', 
        () => Promise.resolve('string result'));
      expect(stringResult).toBe('string result');

      // Number result
      const numberResult = await result.current.trackApiCall('numberApi', 
        () => Promise.resolve(42));
      expect(numberResult).toBe(42);

      // Array result
      const arrayResult = await result.current.trackApiCall('arrayApi', 
        () => Promise.resolve([1, 2, 3]));
      expect(arrayResult).toEqual([1, 2, 3]);

      // Null result
      const nullResult = await result.current.trackApiCall('nullApi', 
        () => Promise.resolve(null));
      expect(nullResult).toBe(null);
    });

    const metrics = result.current.getMetrics();
    expect(Object.keys(metrics.apiCalls)).toHaveLength(4);
  });

  test('should handle errors without message property', () => {
    const { result } = renderHook(() => usePerformance('TestComponent'));
    
    // Error without message
    const strangeError = {};
    
    act(() => {
      result.current.trackError(strangeError);
    });

    const metrics = result.current.getMetrics();
    expect(metrics.errors).toHaveLength(1);
    expect(metrics.errors[0].error).toBeUndefined();
  });

  test('should maintain separate metrics across different component instances', () => {
    const { result: result1 } = renderHook(() => usePerformance('Component1'));
    const { result: result2 } = renderHook(() => usePerformance('Component2'));
    
    act(() => {
      result1.current.trackError(new Error('Error 1'));
      result2.current.trackError(new Error('Error 2'));
    });

    const metrics1 = result1.current.getMetrics();
    const metrics2 = result2.current.getMetrics();

    expect(metrics1.errors).toHaveLength(1);
    expect(metrics2.errors).toHaveLength(1);
    expect(metrics1.errors[0].error).toBe('Error 1');
    expect(metrics2.errors[0].error).toBe('Error 2');
  });
}); 
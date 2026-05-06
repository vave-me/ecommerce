import { renderHook, act } from '@testing-library/react';
import useActivityApi from '@/hooks/useActivityApi.jsx';
import { addInteraction, createActivity, getActivity } from "@/api/client/activityApi.jsx";

// Mock dependencies
jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn().mockReturnValue({}),
    info: jest.fn().mockReturnValue({}),
    error: jest.fn().mockReturnValue({}),
    warn: jest.fn().mockReturnValue({})
  }
}));

jest.mock('@/api/client/activityApi.jsx', () => ({
  addInteraction: jest.fn(),
  createActivity: jest.fn(),
  getActivity: jest.fn()
}));

// Mock console methods to suppress expected warnings
console.warn = jest.fn();
console.error = jest.fn();

// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: jest.fn((key) => store[key] || null),
    setItem: jest.fn((key, value) => {
      store[key] = value;
    }),
    removeItem: jest.fn((key) => {
      delete store[key];
    }),
    clear: jest.fn(() => {
      store = {};
    })
  };
})();
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('useActivityApi hook', () => {
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    localStorage.clear();
    
    // Default mock implementations
    getActivity.mockResolvedValue({ activityId: 'existing-activity-id' });
    createActivity.mockResolvedValue({ id: 'new-activity-id' });
    addInteraction.mockResolvedValue({ success: true });
  });

  test('should handle like interactions when user is logged in', async () => {
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleLike('product-123', 'user-123');
    });
    
    // Should fetch activity ID
    expect(getActivity).toHaveBeenCalledWith('user-123');
    
    // Should add like interaction
    expect(addInteraction).toHaveBeenCalledWith(
      'existing-activity-id', 
      'product-123', 
      'product', 
      'like'
    );
  });

  test('should handle dislike interactions when user is logged in', async () => {
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleDislike('product-123', 'user-123');
    });
    
    // Should fetch activity ID
    expect(getActivity).toHaveBeenCalledWith('user-123');
    
    // Should add dislike interaction
    expect(addInteraction).toHaveBeenCalledWith(
      'existing-activity-id', 
      'product-123', 
      'product', 
      'dislike'
    );
  });

  test('should use cached activity ID from localStorage if available', async () => {
    // Setup localStorage with existing activity ID
    localStorage.setItem('activityId', 'stored-activity-id');
    
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleLike('product-123', 'user-123');
    });
    
    // Should NOT fetch activity ID from API
    expect(getActivity).not.toHaveBeenCalled();
    
    // Should use cached ID
    expect(addInteraction).toHaveBeenCalledWith(
      'stored-activity-id', 
      'product-123', 
      'product', 
      'like'
    );
  });

  test('should create new activity if none exists', async () => {
    // Mock getActivity to fail
    getActivity.mockRejectedValueOnce(new Error('Not found'));
    
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleLike('product-123', 'user-123');
    });
    
    // Should try to fetch activity ID
    expect(getActivity).toHaveBeenCalledWith('user-123');
    
    // Should create new activity
    expect(createActivity).toHaveBeenCalledWith('user-123');
    
    // Should add interaction with new ID
    expect(addInteraction).toHaveBeenCalledWith(
      'new-activity-id', 
      'product-123', 
      'product', 
      'like'
    );
  });

  test('should warn if user is not logged in', async () => {
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleLike('product-123', null);
    });
    
    // Should NOT make API calls
    expect(getActivity).not.toHaveBeenCalled();
    expect(addInteraction).not.toHaveBeenCalled();
  });

  test('should handle API errors gracefully', async () => {
    // Mock addInteraction to fail
    addInteraction.mockRejectedValueOnce(new Error('API error'));
    
    const { result } = renderHook(() => useActivityApi());
    
    await act(async () => {
      await result.current.handleLike('product-123', 'user-123');
    });
    
    // Should still complete without throwing (the error is caught internally)
    expect(getActivity).toHaveBeenCalled();
  });
}); 
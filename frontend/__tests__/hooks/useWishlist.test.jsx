import { renderHook, act } from '@testing-library/react';
import useWishlist from '@/hooks/useWishlist.jsx';
import {
  getWishlist,
  startWishlist,
  addItemToWishlist,
  removeItemFromWishlist,
  getWishlistId
} from '@/api/client/wishlistApi.jsx';
import { useAuth } from '../../context/AuthContext.jsx';
import { toast } from 'react-toastify';
import React from 'react';

// Mock the dependencies
jest.mock('@/api/client/wishlistApi.jsx', () => ({
  getWishlist: jest.fn(),
  startWishlist: jest.fn(),
  addItemToWishlist: jest.fn(),
  removeItemFromWishlist: jest.fn(),
  getWishlistId: jest.fn()
}));

jest.mock('../../context/AuthContext.jsx', () => ({
  useAuth: jest.fn()
}));

jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
    warn: jest.fn()
  }
}));

// Create a custom wrapper component to modify state
const createWrapper = (initialWishlistId = null, initialItems = []) => {
  // Mock useState to return our controlled values
  const originalUseState = React.useState;
  jest.spyOn(React, 'useState').mockImplementation((initialValue) => {
    if (initialValue === null && arguments.callee.caller.name === 'useWishlist') {
      return [initialWishlistId, jest.fn()];
    } else if (Array.isArray(initialValue) && arguments.callee.caller.name === 'useWishlist') {
      return [initialItems, jest.fn()];
    }
    return originalUseState(initialValue);
  });
};

describe('useWishlist hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Default auth mock
    useAuth.mockReturnValue({
      user: { userId: 'user123' }
    });

    // Default API responses
    getWishlist.mockResolvedValue({ items: [{ productId: 'item1' }] });
    getWishlistId.mockResolvedValue({ wishlistId: 'wish123' });
    startWishlist.mockResolvedValue({ id: 'new-wish123' });
    addItemToWishlist.mockResolvedValue({ success: true });
    removeItemFromWishlist.mockResolvedValue({ success: true });
  });
  
  test('should initialize with empty state', () => {
    const { result } = renderHook(() => useWishlist());
    
    expect(result.current.wishlistId).toBeNull();
    expect(result.current.items).toEqual([]);
  });
  
  test('loadWishlist should show warning when user is not logged in', async () => {
    useAuth.mockReturnValue({ user: null });
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    expect(toast.warn).toHaveBeenCalledWith('Please log in to view your wishlist.');
    expect(getWishlistId).not.toHaveBeenCalled();
  });
  
  test('loadWishlist should find and use existing wishlist ID', async () => {
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    expect(getWishlistId).toHaveBeenCalledWith('user123');
    expect(getWishlist).toHaveBeenCalledWith('wish123');
  });
  
  test('loadWishlist should do nothing when no wishlist exists', async () => {
    getWishlistId.mockResolvedValue(null);
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    expect(getWishlistId).toHaveBeenCalledWith('user123');
    expect(startWishlist).not.toHaveBeenCalled(); // Should only create on first add
  });
  
  test('addItem should create new wishlist when none exists', async () => {
    getWishlistId.mockResolvedValue(null);
    startWishlist.mockResolvedValue({ id: 'new-wish123' });
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.addItem('prod1');
    });
    
    expect(getWishlistId).toHaveBeenCalledWith('user123');
    expect(startWishlist).toHaveBeenCalledWith('user123');
    expect(addItemToWishlist).toHaveBeenCalledWith('new-wish123', 'prod1');
    expect(toast.success).toHaveBeenCalledWith('Item added to wishlist.');
  });
  
  test('addItem should use existing wishlist ID when available', async () => {
    getWishlistId.mockResolvedValue({ wishlistId: 'wish123' });
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.addItem('prod1');
    });
    
    expect(getWishlistId).toHaveBeenCalledWith('user123');
    expect(addItemToWishlist).toHaveBeenCalledWith('wish123', 'prod1');
    expect(toast.success).toHaveBeenCalledWith('Item added to wishlist.');
  });
  
  test('addItem should show warning when user is not logged in', async () => {
    useAuth.mockReturnValue({ user: null });
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.addItem('prod1');
    });
    
    expect(toast.warn).toHaveBeenCalledWith('Please log in to manage your wishlist.');
    expect(addItemToWishlist).not.toHaveBeenCalled();
  });
  
  test('toggleItem should call addItem when product is not in wishlist', async () => {
    // Setup a mock implementation with initial items
    getWishlist.mockResolvedValue({ items: [{ productId: 'prod2' }] });
    
    // Render and load initial data
    const { result } = renderHook(() => useWishlist());
    
    // Load the wishlist first to populate items
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    // Reset mocks for the next operation
    jest.clearAllMocks();
    
    // Now toggle an item that is not in the wishlist
    await act(async () => {
      await result.current.toggleItem('prod1');
    });
    
    // We expect addItem to be called since the product isn't in the list
    expect(addItemToWishlist).toHaveBeenCalled();
  });
  
  test('isInWishlist should check against current items', async () => {
    // Setup a mock implementation with initial items
    getWishlist.mockResolvedValue({ items: [
      { productId: 'prod1' },
      { productId: 'prod2' }
    ]});
    
    // Render and load initial data
    const { result } = renderHook(() => useWishlist());
    
    // Load the wishlist first to populate items
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    // Check if isInWishlist works with the loaded items
    expect(result.current.isInWishlist('prod1')).toBe(true);
    expect(result.current.isInWishlist('prod3')).toBe(false);
  });
  
  test('error handling when API calls fail', async () => {
    getWishlist.mockRejectedValue(new Error('API error'));
    getWishlistId.mockResolvedValue({ wishlistId: 'wish123' });
    
    const { result } = renderHook(() => useWishlist());
    
    await act(async () => {
      await result.current.loadWishlist();
    });
    
    expect(toast.error).toHaveBeenCalledWith('Could not load wishlist items.');
  });
}); 
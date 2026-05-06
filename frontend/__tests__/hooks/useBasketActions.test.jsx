import { renderHook, act } from '@testing-library/react';
import useBasketActions from '@/hooks/useBasketActions.jsx';
import { 
  getOrCreateBasket, 
  getBasket, 
  addItemToBasket, 
  removeItemFromBasket, 
  updateItemQuantity,
  checkoutBasket 
} from '@/api/client/basketApi.jsx';
import { useAuth } from '../../context/AuthContext.jsx';
import { toast } from 'react-toastify';

// Mock dependencies
jest.mock('@/api/client/basketApi.jsx', () => ({
  getOrCreateBasket: jest.fn(),
  getBasket: jest.fn(),
  addItemToBasket: jest.fn(),
  removeItemFromBasket: jest.fn(),
  updateItemQuantity: jest.fn(),
  checkoutBasket: jest.fn()
}));

jest.mock('../../context/AuthContext.jsx', () => ({
  useAuth: jest.fn()
}));

jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn(),
    info: jest.fn(),
    error: jest.fn(),
    warn: jest.fn()
  }
}));

describe('useBasketActions hook', () => {
  const mockUserId = 'user-123';
  const mockProductId = 'product-456';
  const mockBasketId = 'basket-789';
  const mockBasket = { 
    id: mockBasketId, 
    items: [
      { id: 'item-1', productId: 'product-1', quantity: 2 },
      { id: 'item-2', productId: mockProductId, quantity: 1 }
    ],
    basketStatus: 'OPEN'
  };

  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Default auth mock returning a logged-in user
    useAuth.mockReturnValue({
      user: { userId: mockUserId }
    });
    
    // Default API responses
    getOrCreateBasket.mockResolvedValue({ basketId: mockBasketId });
    getBasket.mockResolvedValue({ basket: mockBasket });
    addItemToBasket.mockResolvedValue({ success: true });
    removeItemFromBasket.mockResolvedValue({ success: true });
    updateItemQuantity.mockResolvedValue({ success: true });
    checkoutBasket.mockResolvedValue({ success: true, orderId: 'order-123' });
  });

  test('fetchBasket should load the basket', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Should call APIs in sequence
    expect(getOrCreateBasket).toHaveBeenCalledWith(mockUserId);
    expect(getBasket).toHaveBeenCalledWith(mockBasketId);
    
    // Should update state with basket
    expect(result.current.basket).toEqual(mockBasket);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBe('');
  });

  test('fetchBasket should do nothing if not logged in', async () => {
    // Mock user as not logged in
    useAuth.mockReturnValue({ user: null });
    
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Should not call any APIs
    expect(getOrCreateBasket).not.toHaveBeenCalled();
    expect(getBasket).not.toHaveBeenCalled();
    
    // Should clear basket state
    expect(result.current.basket).toBeNull();
  });

  test('handleAddToCart should add an item to the basket', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    await act(async () => {
      await result.current.handleAddToCart(2); // Add 2 quantity
    });
    
    // Should call APIs in sequence
    expect(getOrCreateBasket).toHaveBeenCalledWith(mockUserId);
    expect(addItemToBasket).toHaveBeenCalledWith(mockBasketId, mockProductId, 2);
    expect(getBasket).toHaveBeenCalledWith(mockBasketId);
    
    // Should show success toast
    expect(toast.success).toHaveBeenCalled();
    
    // Should update state with basket
    expect(result.current.basket).toEqual(mockBasket);
  });

  test('handleAddToCart should warn if not logged in', async () => {
    // Mock user as not logged in
    useAuth.mockReturnValue({ user: null });
    
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    await act(async () => {
      await result.current.handleAddToCart();
    });
    
    // Should show warning toast
    expect(toast.warn).toHaveBeenCalled();
    
    // Should not call any APIs
    expect(getOrCreateBasket).not.toHaveBeenCalled();
    expect(addItemToBasket).not.toHaveBeenCalled();
  });

  test('removeItem should remove an item from the basket', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    // First set the basket to have a valid state
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Now try to remove an item
    await act(async () => {
      await result.current.removeItem('item-1');
    });
    
    // Should call remove API
    expect(removeItemFromBasket).toHaveBeenCalledWith(mockBasketId, 'item-1');
    
    // Should update basket state
    expect(getBasket).toHaveBeenCalledTimes(2); // Once for fetch, once for update
    
    // Should show success toast
    expect(toast.success).toHaveBeenCalled();
  });

  test('updateQuantity should update item quantity', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    // First set the basket
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Now update quantity
    await act(async () => {
      await result.current.updateQuantity('item-2', 3);
    });
    
    // Should call update API
    expect(updateItemQuantity).toHaveBeenCalledWith(mockBasketId, 'item-2', 3);
    
    // Should update basket state
    expect(getBasket).toHaveBeenCalledTimes(2);
  });

  test('updateQuantity should warn if quantity < 1', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    // First set the basket
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Now try to update with invalid quantity
    await act(async () => {
      await result.current.updateQuantity('item-2', 0);
    });
    
    // Should show warning toast
    expect(toast.warn).toHaveBeenCalled();
    
    // Should not call update API
    expect(updateItemQuantity).not.toHaveBeenCalled();
  });

  test('should checkout the basket', async () => {
    const { result } = renderHook(() => useBasketActions(mockProductId));
    
    // First set the basket
    await act(async () => {
      await result.current.fetchBasket();
    });
    
    // Now checkout
    await act(async () => {
      await result.current.handleCheckout(mockUserId);
    });
    
    // Should call checkout API
    expect(checkoutBasket).toHaveBeenCalledWith(mockBasketId, mockUserId);
    
    // Should show success toast
    expect(toast.success).toHaveBeenCalled();
  });
}); 
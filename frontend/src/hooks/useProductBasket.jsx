/**
 * Product-specific Basket Hook
 * 
 * This hook provides basket operations for a specific product,
 * useful for product detail pages and product cards.
 */
import { useCallback, useMemo } from 'react';
import { useBasket } from './useBasket';

/**
 * Hook for product-specific basket operations
 * @param {string} productId - The product ID to operate on
 * @param {object} productData - Optional product data (price, name, etc.)
 */
export function useProductBasket(productId, productData = {}) {
  const basket = useBasket();
  
  // Check if this specific product is in basket
  const isInBasket = useMemo(() => {
    return basket.isInBasket(productId);
  }, [basket, productId]);
  
  // Get quantity of this product in basket
  const quantity = useMemo(() => {
    return basket.getItemQuantity(productId);
  }, [basket, productId]);
  
  // Add this product to basket
  const addToBasket = useCallback((qty = 1) => {
    if (!productId) {
      
      return;
    }
    
    return basket.addToBasket(
      productId, 
      qty, 
      productData.price || 0, 
      {
        name: productData.name,
        image: productData.image,
        ...productData
      }
    );
  }, [basket, productId, productData]);
  
  // Remove this product from basket
  const removeFromBasket = useCallback(() => {
    if (!productId) return;
    
    // Find the item ID for this product
    const item = basket.items.find(item => 
      item.product_id === productId || item.productId === productId
    );
    
    if (item) {
      const itemId = item.item_id || item.id || productId;
      return basket.removeItem(itemId);
    }
  }, [basket, productId]);
  
  // Update quantity of this product
  const updateQuantity = useCallback((newQuantity) => {
    if (!productId) return;
    
    // Find the item ID for this product
    const item = basket.items.find(item => 
      item.product_id === productId || item.productId === productId
    );
    
    if (item) {
      const itemId = item.item_id || item.id || productId;
      
      if (newQuantity <= 0) {
        return basket.removeItem(itemId);
      }
      
      return basket.updateQuantity(itemId, newQuantity);
    } else if (newQuantity > 0) {
      // If not in basket and quantity > 0, add it
      return addToBasket(newQuantity);
    }
  }, [basket, productId, addToBasket]);
  
  // Increment quantity
  const incrementQuantity = useCallback(() => {
    return updateQuantity(quantity + 1);
  }, [updateQuantity, quantity]);
  
  // Decrement quantity
  const decrementQuantity = useCallback(() => {
    return updateQuantity(Math.max(0, quantity - 1));
  }, [updateQuantity, quantity]);
  
  // Toggle product in basket
  const toggleBasket = useCallback(() => {
    if (isInBasket) {
      return removeFromBasket();
    } else {
      return addToBasket(1);
    }
  }, [isInBasket, addToBasket, removeFromBasket]);
  
  return {
    // State
    isInBasket,
    quantity,
    
    // Actions
    addToBasket,
    removeFromBasket,
    updateQuantity,
    incrementQuantity,
    decrementQuantity,
    toggleBasket,
    
    // Loading states from parent
    isLoading: basket.isAddingToBasket || basket.isUpdatingQuantity || basket.isRemovingItem,
    
    // Access to full basket data if needed
    basket: basket.basket,
    basketItemCount: basket.itemCount,
    basketTotalAmount: basket.totalAmount,
  };
}

// Export as default for backward compatibility
export default useProductBasket;
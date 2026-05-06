/**
 * Basket Hooks Exports
 * 
 * Use these hooks for basket/cart functionality:
 * - useBasket: Main hook for general basket operations (recommended)
 * - useProductBasket: Hook for product-specific basket operations
 */

// Main hooks
export { useBasket, default as useBasketDefault } from '../useBasket';
export { useProductBasket, default as useProductBasketDefault } from '../useProductBasket';
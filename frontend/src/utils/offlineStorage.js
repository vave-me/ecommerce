// Offline storage utilities for wishlist and basket
const OFFLINE_WISHLIST_KEY = 'offline_wishlist_items';
const OFFLINE_BASKET_KEY = 'offline_basket_items';

// Wishlist offline storage
export const offlineWishlistStorage = {
  getItems() {
    try {
      const items = localStorage.getItem(OFFLINE_WISHLIST_KEY);
      return items ? JSON.parse(items) : [];
    } catch (error) {
      // Error: 'Error reading offline wishlist:', error...
      return [];
    }
  },

  addItem(itemId, entityType = 'product', notes = '') {
    try {
      const items = this.getItems();
      const exists = items.some(item => 
        item.itemId === itemId && item.entityType === entityType
      );
      
      if (!exists) {
        items.push({
          itemId,
          entityType,
          notes,
          addedAt: new Date().toISOString()
        });
        localStorage.setItem(OFFLINE_WISHLIST_KEY, JSON.stringify(items));
      }
      return true;
    } catch (error) {
      // Error: 'Error adding to offline wishlist:', error...
      return false;
    }
  },

  removeItem(itemId, entityType = 'product') {
    try {
      const items = this.getItems();
      const filtered = items.filter(item => 
        !(item.itemId === itemId && item.entityType === entityType)
      );
      localStorage.setItem(OFFLINE_WISHLIST_KEY, JSON.stringify(filtered));
      return true;
    } catch (error) {
      // Error: 'Error removing from offline wishlist:', error...
      return false;
    }
  },

  hasItem(itemId, entityType = 'product') {
    const items = this.getItems();
    return items.some(item => 
      item.itemId === itemId && item.entityType === entityType
    );
  },

  clear() {
    try {
      localStorage.removeItem(OFFLINE_WISHLIST_KEY);
      return true;
    } catch (error) {
      // Error: 'Error clearing offline wishlist:', error...
      return false;
    }
  }
};

// Basket offline storage
export const offlineBasketStorage = {
  getItems() {
    try {
      const items = localStorage.getItem(OFFLINE_BASKET_KEY);
      return items ? JSON.parse(items) : [];
    } catch (error) {
      // Error: 'Error reading offline basket:', error...
      return [];
    }
  },

  addItem(productId, quantity = 1) {
    try {
      const items = this.getItems();
      const existingIndex = items.findIndex(item => item.productId === productId);
      
      if (existingIndex >= 0) {
        // Update quantity if item exists
        items[existingIndex].quantity += quantity;
      } else {
        // Add new item
        items.push({
          productId,
          quantity,
          addedAt: new Date().toISOString()
        });
      }
      
      localStorage.setItem(OFFLINE_BASKET_KEY, JSON.stringify(items));
      return true;
    } catch (error) {
      // Error: 'Error adding to offline basket:', error...
      return false;
    }
  },

  updateQuantity(productId, newQuantity) {
    try {
      const items = this.getItems();
      const index = items.findIndex(item => item.productId === productId);
      
      if (index >= 0) {
        if (newQuantity <= 0) {
          // Remove item if quantity is 0 or less
          items.splice(index, 1);
        } else {
          items[index].quantity = newQuantity;
        }
        localStorage.setItem(OFFLINE_BASKET_KEY, JSON.stringify(items));
      }
      return true;
    } catch (error) {
      // Error: 'Error updating offline basket quantity:', error...
      return false;
    }
  },

  removeItem(productId) {
    try {
      const items = this.getItems();
      const filtered = items.filter(item => item.productId !== productId);
      localStorage.setItem(OFFLINE_BASKET_KEY, JSON.stringify(filtered));
      return true;
    } catch (error) {
      // Error: 'Error removing from offline basket:', error...
      return false;
    }
  },

  getTotalQuantity() {
    const items = this.getItems();
    return items.reduce((total, item) => total + item.quantity, 0);
  },

  clear() {
    try {
      localStorage.removeItem(OFFLINE_BASKET_KEY);
      return true;
    } catch (error) {
      // Error: 'Error clearing offline basket:', error...
      return false;
    }
  }
};

// Check if we're in a browser environment
export const isClientSide = () => typeof window !== 'undefined' && typeof localStorage !== 'undefined';
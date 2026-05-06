import axiosInstance from '../../axiosInstance';

// Merchant Center Integration APIs

export const listMerchantProducts = async (params = {}) => {
  try {
    const response = await axiosInstance.get('/merchant/products', { params });
    return response.data;
  } catch (error) {
    // Error: 'Failed to list merchant products:', error...
    throw error;
  }
};

export const syncProductToMerchant = async (productData) => {
  try {
    const response = await axiosInstance.post('/merchant/products/sync', productData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to sync product to merchant:', error...
    throw error;
  }
};

export const batchSyncProductsToMerchant = async (products) => {
  try {
    const response = await axiosInstance.post('/merchant/products/batch-sync', { products });
    return response.data;
  } catch (error) {
    // Error: 'Failed to batch sync products:', error...
    throw error;
  }
};

export const removeProductFromMerchant = async (productId) => {
  try {
    const response = await axiosInstance.delete(`/merchant/products/${productId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to remove product from merchant:', error...
    throw error;
  }
};

export const getProductMerchantStatus = async (productId) => {
  try {
    const response = await axiosInstance.get(`/merchant/products/${productId}/status`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to get product merchant status:', error...
    throw error;
  }
};

// Bulk operations
export const bulkSyncProductsToMerchant = async (productIds) => {
  try {
    const products = await Promise.all(
      productIds.map(async (id) => {
        const productData = await axiosInstance.get(`/products/${id}`);
        return {
          id: productData.data.product.id,
          name: productData.data.product.name,
          description: productData.data.product.description,
          price: productData.data.product.basePrice,
          currency: 'USD',
          availability: productData.data.product.stock > 0 ? 'in stock' : 'out of stock',
          condition: productData.data.product.condition || 'new',
          brand: productData.data.product.brand,
          googleProductCategory: productData.data.product.categoryId,
          imageUrl: productData.data.product.thumbnail,
          link: `${window.location.origin}/products/${id}`,
          stock: productData.data.product.stock,
          sku: productData.data.product.sku
        };
      })
    );
    
    return batchSyncProductsToMerchant(products);
  } catch (error) {
    // Error: 'Failed to bulk sync products to merchant:', error...
    throw error;
  }
};

export const bulkRemoveProductsFromMerchant = async (productIds) => {
  try {
    const promises = productIds.map(id => removeProductFromMerchant(id));
    const results = await Promise.allSettled(promises);
    return {
      success: results.filter(r => r.status === 'fulfilled').length,
      failed: results.filter(r => r.status === 'rejected').length,
      results
    };
  } catch (error) {
    // Error: 'Failed to bulk remove products from merchant:', e...
    throw error;
  }
};

// Export alias for batchSyncProducts
export const batchSyncProducts = batchSyncProductsToMerchant;
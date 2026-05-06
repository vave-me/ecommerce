import axiosInstance from '../../axiosInstance';

// Helper functions
const cleanParams = (params) => {
  const cleaned = {};
  Object.keys(params).forEach(key => {
    if (params[key] !== undefined && params[key] !== null && params[key] !== '') {
      cleaned[key] = params[key];
    }
  });
  return cleaned;
};

// Admin Product Management APIs

export const getAdminProducts = async (params = {}) => {
  try {
    const cleanedParams = cleanParams(params);
    const response = await axiosInstance.get('/products', { params: cleanedParams });
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch admin products:', error...
    throw error;
  }
};

export const getAdminProductById = async (productId) => {
  try {
    const response = await axiosInstance.get(`/products/${productId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch product details:', error...
    throw error;
  }
};

export const getProductsWithFilters = async (filters = {}) => {
  try {
    const response = await axiosInstance.post('/products/filters', filters);
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch products with filters:', error...
    throw error;
  }
};

export const getProductsByCategory = async (categoryId, params = {}) => {
  try {
    const cleanedParams = cleanParams(params);
    const response = await axiosInstance.get(`/products/categories/${categoryId}`, { params: cleanedParams });
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch products by category:', error...
    throw error;
  }
};

export const updateProductDetails = async (productId, productData) => {
  try {
    const response = await axiosInstance.post(`/products/${productId}`, productData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to update product:', error...
    throw error;
  }
};

export const updateProductPrice = async (productId, priceData) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/price`, priceData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to update product price:', error...
    throw error;
  }
};

export const adjustProductStock = async (productId, stockData) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/stock`, stockData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to adjust product stock:', error...
    throw error;
  }
};

export const archiveProduct = async (productId) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/archive`, {});
    return response.data;
  } catch (error) {
    // Error: 'Failed to archive product:', error...
    throw error;
  }
};

export const markProductSold = async (productId) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/sold`, {});
    return response.data;
  } catch (error) {
    // Error: 'Failed to mark product as sold:', error...
    throw error;
  }
};

export const markProductLeased = async (productId, leaseData) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/lease`, leaseData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to mark product as leased:', error...
    throw error;
  }
};

export const markProductPawned = async (productId, pawnData) => {
  try {
    const response = await axiosInstance.patch(`/products/${productId}/pawn`, pawnData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to mark product as pawned:', error...
    throw error;
  }
};

export const deleteProduct = async (productId) => {
  try {
    const response = await axiosInstance.delete(`/products/${productId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to delete product:', error...
    throw error;
  }
};

// Product Variants
export const getProductVariants = async (productId, params = {}) => {
  try {
    const cleanedParams = cleanParams(params);
    const response = await axiosInstance.get(`/products/${productId}/variants`, { params: cleanedParams });
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch product variants:', error...
    throw error;
  }
};

export const addProductVariant = async (variantData) => {
  try {
    const response = await axiosInstance.post('/products/variants', variantData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to add product variant:', error...
    throw error;
  }
};

export const updateVariantPrice = async (variantId, priceData, increase = false) => {
  try {
    const endpoint = increase ? 
      `/products/variants/${variantId}/increasePrice` : 
      `/products/variants/${variantId}/decreasePrice`;
    const response = await axiosInstance.patch(endpoint, priceData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to update variant price:', error...
    throw error;
  }
};

export const adjustVariantStock = async (variantId, stockData) => {
  try {
    const response = await axiosInstance.patch(`/products/variants/${variantId}/stock`, stockData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to adjust variant stock:', error...
    throw error;
  }
};

export const archiveVariant = async (variantId) => {
  try {
    const response = await axiosInstance.patch(`/products/variants/${variantId}/archive`, {});
    return response.data;
  } catch (error) {
    // Error: 'Failed to archive variant:', error...
    throw error;
  }
};

export const deleteVariant = async (variantId) => {
  try {
    const response = await axiosInstance.delete(`/products/variants/${variantId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to delete variant:', error...
    throw error;
  }
};

// Product Images
export const updateProductThumbnail = async (productId, thumbnailData) => {
  try {
    const response = await axiosInstance.post(`/products/${productId}/thumbnail/update`, thumbnailData);
    return response.data;
  } catch (error) {
    // Error: 'Failed to update product thumbnail:', error...
    throw error;
  }
};

// Bulk Operations
export const bulkUpdateProductStatus = async (productIds, status, additionalData = {}) => {
  try {
    const promises = productIds.map(id => {
      switch (status) {
        case 'archived':
          return archiveProduct(id);
        case 'sold':
          return markProductSold(id);
        case 'leased':
          return markProductLeased(id, additionalData);
        case 'pawned':
          return markProductPawned(id, additionalData);
        default:
          throw new Error(`Unknown status: ${status}`);
      }
    });
    
    const results = await Promise.allSettled(promises);
    return {
      success: results.filter(r => r.status === 'fulfilled').length,
      failed: results.filter(r => r.status === 'rejected').length,
      results
    };
  } catch (error) {
    // Error: 'Failed to bulk update product status:', error...
    throw error;
  }
};

export const bulkDeleteProducts = async (productIds) => {
  try {
    const promises = productIds.map(id => deleteProduct(id));
    const results = await Promise.allSettled(promises);
    return {
      success: results.filter(r => r.status === 'fulfilled').length,
      failed: results.filter(r => r.status === 'rejected').length,
      results
    };
  } catch (error) {
    // Error: 'Failed to bulk delete products:', error...
    throw error;
  }
};

// Export products
export const exportProducts = async (params = {}) => {
  try {
    const cleanedParams = cleanParams(params);
    const response = await axiosInstance.get('/products/export', { 
      params: cleanedParams,
      responseType: 'blob'
    });
    
    // Create a blob URL and trigger download
    const blob = new Blob([response.data], { type: response.headers['content-type'] });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `products-export-${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
    
    return { success: true };
  } catch (error) {
    // Error: 'Failed to export products:', error...
    throw error;
  }
};
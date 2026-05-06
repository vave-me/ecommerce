// src/api/postApi.jsx
import axiosInstance from '../axiosInstance';

export const addPost = async (postData) => {
    try {

        const response = await axiosInstance.post('/posts', postData);

        return response.data;
    } catch (error) {
        // Error: 'Error adding post:', error...
        throw error;
    }
};

/**
 * RETRIEVE ALL PRODUCTS
 * GET /api/posts
 * Query params => page, pageSize, sortBy, sortOrder
 */
export const getPosts = async (filters = {}) => {
    try {
        const response = await axiosInstance.get('/posts', {params: filters});
        const data = response.data; // This should match postspbGetPostsResponse

        // Convert totalCount, totalPages, currentPage to numbers if needed:
        return {
            posts: data.posts || [],
            totalCount: parseInt(data.totalCount, 10) || 0,
            totalPages: parseInt(data.totalPages, 10) || 0,
            currentPage: parseInt(data.currentPage, 10) || 0
        };
    } catch (error) {
        // Error: 'Error fetching posts:', error...
        throw error;
    }
};

export const getUserPosts = async (userId, filters = {}) => {
    try {
        const response = await axiosInstance.get(`/posts/users/${userId}`, {params: filters});
        const data = response.data; // This should match postspbGetPostsResponse

        // Convert totalCount, totalPages, currentPage to numbers if needed:
        return {
            posts: data.posts || [],
            totalCount: parseInt(data.totalCount, 10) || 0,
            totalPages: parseInt(data.totalPages, 10) || 0,
            currentPage: parseInt(data.currentPage, 10) || 0
        };
    } catch (error) {
        // Error: 'Error fetching posts:', error...
        throw error;
    }
};

/**
 * RETRIEVE PRODUCTS BY CATEGORY
 * GET /api/posts/categories/{categoryId}
 * Accepts query params: page, pageSize, sortBy, sortOrder
 */
export const getPostsByCategory = async (categoryId, filters = {}) => {
    try {
        // Per swagger => /api/posts/categories/{categoryId}
        // + optional: ?page, ?pageSize, ?sortBy, ?sortOrder
        const response = await axiosInstance.get(
            `/posts/categories/${encodeURIComponent(categoryId)}`,
            {params: filters}
        );
        // response matches postspbGetPostsByCategoryResponse
        return response.data;
    } catch (error) {
        // Error: 'Error fetching posts by category:', error...
        throw error;
    }
};

/**
 * RETRIEVE A USER'S CATALOG
 * GET /api/posts/catalog
 * Accepts query params: userId, page, pageSize, sortBy, sortOrder
 *
 *  NOTE: The swagger now indicates the path is
 *        GET /api/posts/catalog
 *        with userId as a query param, not in path.
 */
export const getCatalog = async (userId, filters = {}) => {
    try {
        // Include userId in the query parameters
        const query = { ...filters, userId };
        const response = await axiosInstance.get('/posts/catalog', {
            params: query,
        });
        // Matches postspbGetCatalogResponse => { posts: [...], totalCount, ... }
        return response.data;
    } catch (error) {
        // Error: 'Error fetching user catalog:', error...
        throw error;
    }
};

/**
 * GET A SPECIFIC PRODUCT
 * GET /api/posts/{id}
 * Optionally can pass ?userId=...
 */
export const getPost = async (postId) => {
    try {
        // According to swagger => GET /api/posts/{id}
        const response = await axiosInstance.get(`/posts/${encodeURIComponent(postId)}`);

        return response.data;
    } catch (error) {
        // Error: 'Error fetching post by ID:', error...
        throw error;
    }
};

/**
 * REMOVE A PRODUCT
 * DELETE /api/posts/{id}
 * Optionally can pass ?userId=...
 */
export const removePost = async (postId, userId = '') => {
    try {

        // According to swagger => DELETE /api/posts/{id}
        const response = await axiosInstance.delete(
            `/posts/${encodeURIComponent(postId)}`);

        return response.data;
    } catch (error) {
        // Error: 'Error removing post:', error...
        throw error;
    }
};

/**
 * UPDATE A PRODUCT
 * PATCH /api/posts/{id}
 * Body -> PostsServiceUpdatePostBody
 *
 * => operationId: PostsService_UpdatePost
 */
export const updatePost = async (updateData) => {
    try {

        const response = await axiosInstance.post(
            `/posts/update`,
            updateData
        );
        // response => postspbUpdatePostResponse => { id: string }
        return response.data;
    } catch (error) {
        // Error: 'Error updating post:', error...
        throw error;
    }
};

/**
 * REBRAND A PRODUCT
 * PATCH /api/posts/{id}/rebrand
 * Body -> PostsServiceRebrandPostBody
 */
export const rebrandPost = async (postId, rebrandData) => {
    try {
        // According to swagger => PATCH /api/posts/{id}/rebrand
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/rebrand`,
            rebrandData
        );
        // response => postspbRebrandPostResponse => {}
        return response.data;
    } catch (error) {
        // Error: 'Error rebranding post:', error...
        throw error;
    }
};

/**
 * UPDATE PRODUCT PRICE
 * PATCH /api/posts/{id}/price
 * Body -> PostsServiceUpdatePostPriceBody
 */
export const updatePostPrice = async (postId, priceData) => {
    try {
        // According to swagger => PATCH /api/posts/{id}/price
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/price`,
            priceData
        );
        // response => postspbUpdatePostPriceResponse => {}
        return response.data;
    } catch (error) {
        // Error: 'Error updating post price:', error...
        throw error;
    }
};

/**
 * ARCHIVE A PRODUCT
 * PATCH /api/posts/{postId}/archive
 * Body -> PostsServiceArchivePostBody
 */
export const archivePost = async (postId, archiveData) => {
    try {
        // According to swagger => PATCH /api/posts/{postId}/archive
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/archive`,
            archiveData
        );
        // response => postspbArchivePostResponse => { postId, archived }
        return response.data;
    } catch (error) {
        // Error: 'Error archiving post:', error...
        throw error;
    }
};

/**
 * MARK PRODUCT SOLD
 * PATCH /api/posts/{postId}/sold
 * Body -> PostsServiceMarkPostSoldBody
 */
export const markPostSold = async (postId, soldData) => {
    try {
        // According to swagger => PATCH /api/posts/{postId}/sold
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/sold`,
            soldData
        );
        // response => postspbMarkPostSoldResponse => { postId, status }
        return response.data;
    } catch (error) {
        // Error: 'Error marking post sold:', error...
        throw error;
    }
};

/**
 * MARK PRODUCT LEASED
 * PATCH /api/posts/{postId}/lease
 * Body -> PostsServiceMarkPostLeasedBody
 */
export const markPostLeased = async (postId, leaseData) => {
    try {
        // According to swagger => PATCH /api/posts/{postId}/lease
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/lease`,
            leaseData
        );
        // response => postspbMarkPostLeasedResponse => { postId, status }
        return response.data;
    } catch (error) {
        // Error: 'Error marking post leased:', error...
        throw error;
    }
};

/**
 * MARK PRODUCT PAWNED
 * PATCH /api/posts/{postId}/pawn
 * Body -> PostsServiceMarkPostPawnedBody
 */
export const markPostPawned = async (postId, pawnData) => {
    try {
        // According to swagger => PATCH /api/posts/{postId}/pawn
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/pawn`,
            pawnData
        );
        // response => postspbMarkPostPawnedResponse => { postId, status }
        return response.data;
    } catch (error) {
        // Error: 'Error marking post pawned:', error...
        throw error;
    }
};

/**
 * ADJUST PRODUCT STOCK
 * PATCH /api/posts/{postId}/stock
 * Body -> PostsServiceAdjustPostStockBody
 */
export const adjustPostStock = async (postId, stockData) => {
    try {
        // According to swagger => PATCH /api/posts/{postId}/stock
        const response = await axiosInstance.patch(
            `/posts/${encodeURIComponent(postId)}/stock`,
            stockData
        );
        // response => postspbAdjustPostStockResponse => { postId, oldStock, newStock }
        return response.data;
    } catch (error) {
        // Error: 'Error adjusting post stock:', error...
        throw error;
    }
};

/**
 * GET VARIANTS for a Post
 * GET /api/posts/{postId}/variants
 * query: userId, page, pageSize, sortBy, sortOrder
 */
export const getVariants = async (postId, filters = {}) => {
    try {
        // According to swagger => GET /api/posts/{postId}/variants
        // optionally ?userId=..., ?page=..., ?pageSize=..., ...
        const response = await axiosInstance.get(
            `/posts/${encodeURIComponent(postId)}/variants`,
            {params: filters}
        );
        // response => postspbGetVariantsResponse => { variants: [...], totalCount, ... }
        return response.data;
    } catch (error) {
        // Error: 'Error fetching variants:', error...
        throw error;
    }
};

/**
 * ADD VARIANT
 * POST /api/posts/variants
 * Body -> postspbAddVariantRequest
 */
export const addVariant = async (variantData) => {
    try {
        // According to swagger => POST /api/posts/variants
        const response = await axiosInstance.post('/posts/variants', variantData);
        // response => postspbAddVariantResponse => { variantId }
        return response.data;
    } catch (error) {
        // Error: 'Error adding variant:', error...
        throw error;
    }
};

/**
 * GET SPECIFIC VARIANT
 * GET /api/posts/variants/{variantId}
 * optionally ?userId=...
 */
export const getVariant = async (variantId, userId = '') => {
    try {
        const params = userId ? {userId} : {};
        // According to swagger => GET /api/posts/variants/{variantId}
        const response = await axiosInstance.get(
            `/posts/variants/${encodeURIComponent(variantId)}`,
            {params}
        );
        // response => postspbGetVariantResponse => { variant: {...} }
        return response.data;
    } catch (error) {
        // Error: 'Error fetching variant:', error...
        throw error;
    }
};

/**
 * REMOVE VARIANT
 * DELETE /api/posts/variants/{variantId}
 * optionally ?userId=...
 */
export const removeVariant = async (variantId, userId = '') => {
    try {
        const params = userId ? {userId} : {};
        // According to swagger => DELETE /api/posts/variants/{variantId}
        const response = await axiosInstance.delete(
            `/posts/variants/${encodeURIComponent(variantId)}`,
            {params}
        );
        // response => postspbRemoveVariantResponse => { variantId }
        return response.data;
    } catch (error) {
        // Error: 'Error removing variant:', error...
        throw error;
    }
};

/**
 * ARCHIVE VARIANT
 * PATCH /api/posts/variants/{variantId}/archive
 * Body -> PostsServiceArchiveVariantBody
 */
export const archiveVariant = async (variantId, archiveData) => {
    try {
        // According to swagger => PATCH /api/posts/variants/{variantId}/archive
        const response = await axiosInstance.patch(
            `/posts/variants/${encodeURIComponent(variantId)}/archive`,
            archiveData
        );
        // response => postspbArchiveVariantResponse => { variantId, archived }
        return response.data;
    } catch (error) {
        // Error: 'Error archiving variant:', error...
        throw error;
    }
};

/**
 * DECREASE VARIANT PRICE
 * PATCH /api/posts/variants/{variantId}/decreasePrice
 * Body -> PostsServiceDecreaseVariantPriceBody
 */
export const decreaseVariantPrice = async (variantId, decreaseData) => {
    try {
        // According to swagger => PATCH /api/posts/variants/{variantId}/decreasePrice
        const response = await axiosInstance.patch(
            `/posts/variants/${encodeURIComponent(variantId)}/decreasePrice`,
            decreaseData
        );
        // response => postspbDecreaseVariantPriceResponse => { variantId, oldPrice, newPrice }
        return response.data;
    } catch (error) {
        // Error: 'Error decreasing variant price:', error...
        throw error;
    }
};

/**
 * INCREASE VARIANT PRICE
 * PATCH /api/posts/variants/{variantId}/increasePrice
 * Body -> PostsServiceIncreaseVariantPriceBody
 */
export const increaseVariantPrice = async (variantId, increaseData) => {
    try {
        // According to swagger => PATCH /api/posts/variants/{variantId}/increasePrice
        const response = await axiosInstance.patch(
            `/posts/variants/${encodeURIComponent(variantId)}/increasePrice`,
            increaseData
        );
        // response => postspbIncreaseVariantPriceResponse => { variantId, oldPrice, newPrice }
        return response.data;
    } catch (error) {
        // Error: 'Error increasing variant price:', error...
        throw error;
    }
};

/**
 * ADJUST VARIANT STOCK
 * PATCH /api/posts/variants/{variantId}/stock
 * Body -> PostsServiceAdjustVariantStockBody
 */
export const adjustVariantStock = async (variantId, stockData) => {
    try {
        // According to swagger => PATCH /api/posts/variants/{variantId}/stock
        const response = await axiosInstance.patch(
            `/posts/variants/${encodeURIComponent(variantId)}/stock`,
            stockData
        );
        // response => postspbAdjustVariantStockResponse => { variantId, oldStock, newStock }
        return response.data;
    } catch (error) {
        // Error: 'Error adjusting variant stock:', error...
        throw error;
    }
};

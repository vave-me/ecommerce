// src/api/basketApi.jsx
import axiosInstance from '../../axiosInstance';

/**
 * POST /baskets
 * body: { userId: "xxx" }
 * => { id: "some-new-basket-id" }
 */
export const startBasket = async (userId) => {
    const response = await axiosInstance.post('/baskets', {userId});
    return response.data; // { id: "..." }
};

/**
 * GET /baskets/{basketId}
 * => { basket: { id, items: [...], basketStatus: ... } }
 */
export const getBasket = async (basketId) => {
    try {
        const response = await axiosInstance.get(`/baskets/${basketId}`);
        // returns { basket: {...} }
        return response.data;
    } catch (err) {
        if (err.response && err.response.status === 404) {
            // No basket found
            return null;
        }
        throw err; // re-throw for other errors
    }
};

/**
 * GET /baskets/current
 * => { basketId: "...", basketStatus: "open", ... }
 *    or 404 if no current basket
 * Note: The backend uses the authenticated user's ID from the JWT token
 */
export const getCurrentBasket = async (userId) => {
    try {
        // Pass userId as userCustomerId query parameter if provided
        const params = userId ? { userCustomerId: userId } : {};
        const response = await axiosInstance.get('/baskets/current', { params });
        return response.data; // e.g. { basketId, basketStatus, ... }
    } catch (err) {
        if (err.response && err.response.status === 404) {
            // No current basket found
            return null;
        }
        throw err; // re-throw if other error
    }
};

/**
 * getOrCreateBasket(userId):
 *  1) Tries to getCurrentBasket.
 *  2) If that returns null or missing basketId => create a new basket.
 */
export const getOrCreateBasket = async (userId) => {
    try {
        const currentBasket = await getCurrentBasket(userId);

        if (!currentBasket) {
            const newBasket = await startBasket(userId);
            return {basketId: newBasket.id, basketStatus: 'open'};
        }

        // We have some "currentBasket" object. Check if it has basketId
        const {basketId, basketStatus} = currentBasket;

        if (basketId) {
            // Return the existing basket
            return {basketId, basketStatus};
        }

    } catch (error) {
        // If 404 => create a new one
        if (error.response?.status === 404) {
            const newBasket = await startBasket(userId);
            return {basketId: newBasket.id, basketStatus: 'open'};
        }
        // Otherwise re-throw
        throw error;
    }
};

/**
 * DELETE /baskets/{basketId}?reason=xxx
 */
export const cancelBasket = async (basketId, reason) => {
    const response = await axiosInstance.delete(`/baskets/${basketId}`, {
        params: {reason},
    });
    return response.data;
};

/**
 * POST /baskets/{basketId}/addItem
 */
export const addItemToBasket = async (basketId, productId, quantity) => {
    const data = JSON.stringify({productId, quantity});
    const response = await axiosInstance.post(`/baskets/${basketId}/addItem`, data);
    return response.data;
};

/**
 * DELETE /baskets/{basketId}/removeItem?itemId=xxx&quantity=yyy
 */
export const removeItemFromBasket = async (basketId, itemId, quantity) => {
    const response = await axiosInstance.delete(
        `/baskets/${basketId}/removeItem`,
        {params: {itemId, quantity}}
    );
    return response.data;
};

/**
 * PATCH /baskets/{basketId}/items/{itemId}
 * body: { newQuantity: "123" }
 */
export const updateItemQuantity = async (basketId, itemId, newQuantity) => {
    const response = await axiosInstance.patch(
        `/baskets/${basketId}/items/${itemId}`,
        {newQuantity}
    );
    return response.data;
};

/**
 * POST /baskets/{basketId}/checkout
 * body: { userCustomerId: "..." }
 */
export const checkoutBasket = async (basketId, userCustomerId) => {
    const response = await axiosInstance.post(`/baskets/${basketId}/checkout`, {
        userCustomerId,
    });
    return response.data;
};

/**
 * GET /api/baskets/{basketId}/total
 * => 200 -> { amount: "12345" }
 */
export async function getTotalBasketAmount(basketId) {
    const response = await axiosInstance.get(`/baskets/${basketId}/total`);
    return response.data; // e.g. { amount: "12345" }
}

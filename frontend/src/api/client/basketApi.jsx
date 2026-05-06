// src/api/basketApi.jsx
import axiosInstance from '../axiosInstance';

/**
 * POST /api/baskets
 * body: { user_id: "xxx" }
 * => { id: "some-new-basket-id" }
 */
export const startBasket = async (userId) => {
    const response = await axiosInstance.post('/baskets', { user_id: userId });
    return response.data; // { id: "..." }
};

/**
 * GET /api/baskets/{basketId}
 * => { basket: { id, items: [...], basket_status: ... } }
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
 * GET /api/baskets/current
 * => { basket_id: "...", basket_status: "open", ... }
 *    or 404 if no current basket
 * Note: The backend uses the authenticated user's ID from the JWT token
 */
export const getCurrentBasket = async (userCustomerId) => {
    try {
        // Pass userCustomerId as query parameter if provided
        const params = userCustomerId ? { userCustomerId } : {};
        const response = await axiosInstance.get('/baskets/current', { params });
        return response.data; // e.g. { basket_id, basket_status, ... }
    } catch (err) {
        if (err.response && err.response.status === 404) {
            // No current basket found
            return null;
        }
        throw err; // re-throw if other error
    }
};

/**
 * @deprecated Use getCurrentBasket() to check for existing basket, 
 * and startBasket() only when adding first item
 * 
 * getOrCreateBasket(userCustomerId):
 *  1) Tries to getCurrentBasket.
 *  2) If that returns null or missing basket_id => create a new basket.
 */
export const getOrCreateBasket = async (userCustomerId, userId) => {
    try {
        const currentBasket = await getCurrentBasket(userCustomerId);

        if (!currentBasket || !currentBasket.basket_id) {
            const newBasket = await startBasket(userId);
            return { basketId: newBasket.id, basketStatus: 'open' };
        }

        // We have a current basket with basket_id
        const { basket_id, basket_status } = currentBasket;
        return { basketId: basket_id, basketStatus: basket_status || 'open' };

    } catch (error) {
        // If 404 => create a new one
        if (error.response?.status === 404) {
            const newBasket = await startBasket(userId);
            return { basketId: newBasket.id, basketStatus: 'open' };
        }
        // Otherwise re-throw
        throw error;
    }
};

/**
 * DELETE /api/baskets/{basketId}?reason=xxx
 */
export const cancelBasket = async (basketId, reason) => {
    const response = await axiosInstance.delete(`/baskets/${basketId}`, {
        params: { reason },
    });
    return response.data;
};

/**
 * POST /api/baskets/{basketId}/addItem
 * body: { product_id: "xxx", quantity: 1 }
 * => { basket_id: "...", item_id: "..." }
 */
export const addItemToBasket = async (basketId, productId, quantity) => {
    const response = await axiosInstance.post(`/baskets/${basketId}/addItem`, {
        product_id: productId,
        quantity: quantity
    });
    return response.data;
};

/**
 * DELETE /api/baskets/{basketId}/removeItem?item_id=xxx&quantity=yyy
 */
export const removeItemFromBasket = async (basketId, itemId, quantity = 1) => {
    const response = await axiosInstance.delete(
        `/baskets/${basketId}/removeItem`,
        { params: { item_id: itemId, quantity } }
    );
    return response.data;
};

/**
 * PATCH /api/baskets/{basketId}/items/{itemId}
 * body: { new_quantity: 123 }
 */
export const updateItemQuantity = async (basketId, itemId, newQuantity) => {
    const response = await axiosInstance.patch(
        `/baskets/${basketId}/items/${itemId}`,
        { new_quantity: newQuantity }
    );
    return response.data;
};

/**
 * POST /api/baskets/{basketId}/checkout
 * body: { user_customer_id: "...", payment_intent_id: "..." }
 */
export const checkoutBasket = async (basketId, userCustomerId, paymentIntentId) => {
    const body = { 
        user_customer_id: userCustomerId,
        payment_intent_id: paymentIntentId 
    };
    const response = await axiosInstance.post(`/baskets/${basketId}/checkout`, body);
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

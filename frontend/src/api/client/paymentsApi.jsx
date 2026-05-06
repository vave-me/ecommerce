import axiosInstance from "../axiosInstance";

const PAYMENTS_API_BASE_URL = '/payments';

export const authorizePayment = async (userCustomerId, amount, basketId = null, orderId = null) => {
    // POST /api/payments/authorize
    // This returns { id, client_secret }
    const paymentData = {
        user_customer_id: userCustomerId,
        amount: amount,
    };
    
    // Add optional fields if provided
    if (basketId) {
        paymentData.basket_id = basketId;
    }
    if (orderId) {
        paymentData.order_id = orderId;
    }

    const response = await axiosInstance.post(`${PAYMENTS_API_BASE_URL}/authorize`, paymentData);
    return response.data;
};

export const payInvoice = async (invoiceId, paymentMethodId) => {
    // PUT /api/payments/invoices/{id}/pay
    const response = await axiosInstance.put(
        `${PAYMENTS_API_BASE_URL}/invoices/${invoiceId}/pay`,
        { payment_method_id: paymentMethodId }
    );
    return response.data;
};

export const confirmPayment = async (paymentId, paymentMethodId) => {
    // POST /api/payments/{id}/confirm
    const response = await axiosInstance.post(
        `${PAYMENTS_API_BASE_URL}/${paymentId}/confirm`,
        { payment_method_id: paymentMethodId }
    );
    return response.data;
};

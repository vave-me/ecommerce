import axiosInstance from "../axiosInstance";

/**
 * Fetches the list of offers from the API
 */
const fetchOffers = async (filters = {}) => {
    try {
        const response = await axiosInstance.get('/aoffers', {params: filters});
        return response.data.offers; // Return the list of offers
    } catch (error) {
        // Error: 'Error fetching offers:', error...
        throw error;
    }
};

const acceptOffer = async (offerId, message = '') => {
    try {
        const payload = {offerResponse: 'ACCEPTED', message};
        await axiosInstance.post(`/offers/${offerId}/respond`, payload);
    } catch (error) {
        // Error: 'Error accepting offer:', error...
        throw error;
    }
};

const rejectOffer = async (offerId, message = '') => {
    try {
        const payload = {offerResponse: 'REJECTED', message};
        await axiosInstance.post(`/offers/${offerId}/respond`, payload);
    } catch (error) {
        // Error: 'Error rejecting offer:', error...
        throw error;
    }
};

const createOffer = async (offerData) => {
    try {
        const response = await axiosInstance.post('/offers', offerData);
        return response.data; // Return the created offer details
    } catch (error) {
        // Error: 'Error creating offer:', error...
        throw error;
    }
};

const negotiateOffer = async (offerId, newPrice, message = '') => {
    try {
        const payload = {newPrice, message};
        const response = await axiosInstance.patch(`/offers/${offerId}/negotiate`, payload);
        return response.data; // Return the updated offer details
    } catch (error) {
        // Error: 'Error negotiating offer:', error...
        throw error;
    }
};

const getOffer = async (offerId) => {
    try {
        const response = await axiosInstance.get(`/offers/${offerId}`);
        return response.data.offer; // Return the offer details
    } catch (error) {
        // Error: 'Error fetching offer:', error...
        throw error;
    }
};

export {
    fetchOffers,
    acceptOffer,
    rejectOffer,
    createOffer,
    negotiateOffer,
    getOffer,
};

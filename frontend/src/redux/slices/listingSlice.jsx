// File: src/redux/slices/listingSlice.js
import { createSlice } from '@reduxjs/toolkit';
const initialListingState = {
    listingType: 'products',
};
const listingSlice = createSlice({
    name: 'listing',
    initialState: initialListingState,
    reducers: {
        // We now expect `action.payload` to be a plain string,
        // e.g. 'products', 'news', 'cars', etc.
        setListingType(state, action) {
            state.listingType = action.payload;
        },
    },
});
export const { setListingType } = listingSlice.actions;
export default listingSlice.reducer;

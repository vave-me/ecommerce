// File: src/redux/slices/modalsSlice.js
import {createSlice} from "@reduxjs/toolkit";
const initialState = {
    // Existing fields
    commentsFullModalOpen: false,
    commentsFullModalItemId: null,
    commentsFullItemType: null,
    commentsFullCategoryId: null,
    messageModalOpen: false,
    messageModalItemId: null,
    messageRecipientId: null,
    addProductModalOpen: false,
    addVehicleModalOpen: false,
    addDealModalOpen: false,
    addPropertyModalOpen: false,
    addJobModalOpen: false,
    addServiceModalOpen: false,
    addPostModalOpen: false,
    isProductModalOpen: false,
    selectedProduct: null,
    isVideoModalOpen: false,
    openModalsCount: 0,
};
const modalsSlice = createSlice({
    name: "modals",
    initialState,
    reducers: {
        /* -------------------------------------
           1) CommentsFull
        ------------------------------------- */
        openCommentsFullModal: (state, action) => {
            const {itemId, itemType, categoryId} = action.payload;
            state.commentsFullModalOpen = true;
            state.commentsFullModalItemId = itemId;
            state.commentsFullItemType = itemType;
            state.commentsFullCategoryId = categoryId;
            state.openModalsCount += 1;
        },
        closeCommentsFullModal: (state) => {
            state.commentsFullModalOpen = false;
            state.commentsFullModalItemId = null;
            state.commentsFullItemType = null;
            state.commentsFullCategoryId = null;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           2) Messages
        ------------------------------------- */
        openMessageModal: (state, action) => {
            const {itemId, recipientId} = action.payload;
            state.messageModalOpen = true;
            state.messageModalItemId = itemId;
            state.messageRecipientId = recipientId;
            state.openModalsCount += 1;
        },
        closeMessageModal: (state) => {
            state.messageModalOpen = false;
            state.messageModalItemId = null;
            state.messageRecipientId = null;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           3) Add Product
        ------------------------------------- */
        openAddProductModal: (state) => {
            state.addProductModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddProductModal: (state) => {
            state.addProductModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           4) Add Post
        ------------------------------------- */
        openAddPostModal: (state) => {
            state.addPostModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddPostModal: (state) => {
            state.addPostModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           5) Add Vehicle
        ------------------------------------- */
        openAddVehicleModal: (state) => {
            state.addVehicleModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddVehicleModal: (state) => {
            state.addVehicleModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           6) Product Detail Modal
        ------------------------------------- */
        openProductModal: (state, action) => {
            state.isProductModalOpen = true;
            state.selectedProduct = action.payload.product;
            state.openModalsCount += 1;
        },
        closeProductModal: (state) => {
            state.isProductModalOpen = false;
            state.selectedProduct = null;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
           7) Video Modal
        ------------------------------------- */
        openAddVideoModal: (state) => {
            state.isVideoModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddVideoModal: (state) => {
            state.isVideoModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        /* -------------------------------------
        8) Add Deal
     ------------------------------------- */
        openAddDealModal: (state) => {
            state.addDealModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddDealModal: (state) => {
            state.addDealModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        openAddPropertyModal: (state) => {
            state.addPropertyModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddPropertyModal: (state) => {
            state.addPropertyModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        openAddServiceModal: (state) => {
            state.addServiceModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddServiceModal: (state) => {
            state.addServiceModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
        openAddJobModal: (state) => {
            state.addJobModalOpen = true;
            state.openModalsCount += 1;
        },
        closeAddJobModal: (state) => {
            state.addJobModalOpen = false;
            state.openModalsCount = Math.max(0, state.openModalsCount - 1);
        },
    },
});
export const {
    // CommentsFull
    openCommentsFullModal,
    closeCommentsFullModal,
    openMessageModal,
    closeMessageModal,
    openAddProductModal,
    closeAddProductModal,
    openAddPostModal,
    closeAddPostModal,
    openAddVehicleModal,
    closeAddVehicleModal,
    openAddDealModal,
    closeAddDealModal,
    openAddPropertyModal,
    closeAddPropertyModal,
    openAddServiceModal,
    closeAddServiceModal,
    openAddJobModal,
    closeAddJobModal,
    openProductModal,
    closeProductModal,
    openAddVideoModal,
    closeAddVideoModal,
} = modalsSlice.actions;
export default modalsSlice.reducer;

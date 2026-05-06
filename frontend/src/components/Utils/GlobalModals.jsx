// File: src/components/GlobalModals.jsx
"use client";
import React, {Suspense, useEffect, memo} from "react";
import {useDispatch, useSelector} from "react-redux";
import {
    closeAddPostModal,
    closeAddProductModal,
    closeAddServiceModal,
    closeAddVideoModal,
    closeCommentsFullModal,
    closeMessageModal,
    // closeProductModal,
} from "../../redux/slices/modalsSlice";
// Import the CSS module
import styles from "./GlobalModals.module.css";
import MessageSectionMobileSetup from "../../features/Messages/MessageSectionMobileSetup";
import CommentsFullSetup from "../../features/CommentsFull/CommentsFullSetup";
// Lazy-load "Add Product" & "Add Post", etc.
// const CreateProductModal = React.lazy(() =>
//     import("../../features/AddProductForm/ProductModal")
// );
const CreateProductModal = React.lazy(() =>
    import("../../features/CreateProductModal").then(module => ({default: module.default}))
);
// const CreatePostModal = React.lazy(() =>
//     import("../../features/AddPostForm/CreatePostModal")
// );
const CreatePostModal = React.lazy(() =>
    import("../../features/CreatePostModal").then(module => ({default: module.default}))
);
const CreateServiceModal = React.lazy(() =>
    import("../../features/CreateServiceModal").then(module => ({default: module.default}))
);
/** NEW: Lazy-load the VideoModal */
const VideoModal = React.lazy(() =>
    import("../../features/CreateVideoModal/VideoModal").then(module => ({default: module.default}))
);
const GlobalModals = memo(function GlobalModals() {
    const dispatch = useDispatch();
    const {
        // CommentsFull
        commentsFullModalOpen,
        commentsFullModalItemId,
        commentsFullItemType,
        commentsFullCategoryId,
        messageModalOpen,
        messageModalItemId,
        messageRecipientId,
        addProductModalOpen,
        addPostModalOpen,
        addServiceModalOpen,
        isVideoModalOpen,
        openModalsCount,
    } = useSelector((state) => state.modals);
    // Instead of a boolean isAnyModalOpen, rely on the count:
    const modalIsActive = openModalsCount > 0;
    // Lock body scroll if ANY modals are open
    useEffect(() => {
        if (modalIsActive) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "";
        }
        return () => {
            // On unmount, reset
            document.body.style.overflow = "";
        };
    }, [modalIsActive]);
    // Close handlers
    const handleCloseCommentsFull = () => dispatch(closeCommentsFullModal());
    const handleCloseMessages = () => dispatch(closeMessageModal());
    const handleCloseProductModal = () => dispatch(closeAddProductModal());
    const handleClosePostModal = () => dispatch(closeAddPostModal());
    const handleCloseServiceModal = () => dispatch(closeAddServiceModal());
    const handleCloseVideoModal = () => dispatch(closeAddVideoModal());
    return (
        <>
            {commentsFullModalOpen && commentsFullModalItemId && commentsFullItemType && commentsFullCategoryId && (
                <Suspense fallback={<div>Loading CommentsFull...</div>}>
                    <CommentsFullSetup
                        itemId={commentsFullModalItemId}
                        itemType={commentsFullItemType}
                        categoryId={commentsFullCategoryId}
                        toggleCommentsFullList={handleCloseCommentsFull}
                    />
                </Suspense>
            )}
            {/* Messages */}
            {messageModalOpen && messageModalItemId && messageRecipientId && (
                <Suspense fallback={<div>Loading Messages...</div>}>
                    <MessageSectionMobileSetup
                        itemId={messageModalItemId}
                        onClose={handleCloseMessages}
                        isOpen={messageModalOpen}
                        recipientId={messageRecipientId}
                        recipient="some recipient name"
                    />
                </Suspense>
            )}
            {/* Add Product */}
            {addProductModalOpen && (
                <Suspense fallback={<div>Loading CreateProductModal...</div>}>
                    <CreateProductModal onClose={handleCloseProductModal}/>
                </Suspense>
            )}
            {/* Add Post */}
            {addPostModalOpen && (
                <Suspense fallback={<div>Loading CreatePostModal...</div>}>
                    <CreatePostModal onClose={handleClosePostModal}/>
                </Suspense>
            )}
            {/* Add Service */}
            {addServiceModalOpen && (
                <Suspense fallback={<div>Loading CreateServiceModal...</div>}>
                    <CreateServiceModal onClose={handleCloseServiceModal}/>
                </Suspense>
            )}
            {/* Video Modal */}
            {isVideoModalOpen && (
                <Suspense fallback={<div>Loading VideoModal...</div>}>
                    <VideoModal onClose={handleCloseVideoModal}/>
                </Suspense>
            )}
        </>
    );
});
GlobalModals.displayName = 'GlobalModals';
export default GlobalModals;

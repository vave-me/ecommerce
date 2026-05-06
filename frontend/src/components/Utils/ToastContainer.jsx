// src/components/ToastContainer.jsx
import React, {memo} from 'react';
import {ToastContainer as ReactToastifyContainer} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
const ToastContainer = memo(() => {
    return (
        <ReactToastifyContainer
            position="top-right"
            autoClose={3000}
            hideProgressBar
            newestOnTop
            closeOnClick
            rtl={false}
            pauseOnFocusLoss
            draggable
            pauseOnHover
        />
    );
});
ToastContainer.displayName = 'ToastContainer';
export default ToastContainer;

// CreatePostModal/components/modal/ModalContainer.jsx
import React, { forwardRef, memo } from 'react';
import PropTypes from 'prop-types';
import styles from '../../CreatePostModal/CreatePostModal.module.css';
export const ModalContainer = memo(forwardRef(function ModalContainer({ children }, ref) {
    return (
        <div className={styles.modalContainer} ref={ref}>
            {children}
        </div>
    );
}));
ModalContainer.propTypes = {
    children: PropTypes.node.isRequired
};
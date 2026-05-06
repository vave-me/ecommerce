// CreatePostModal/components/modal/ModalOverlay.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../CreatePostModal/CreatePostModal.module.css';
export function ModalOverlay({ children }) {
    return (
        <div className={styles.modalOverlay}>
            {children}
        </div>
    );
}
ModalOverlay.propTypes = {
    children: PropTypes.node.isRequired
};
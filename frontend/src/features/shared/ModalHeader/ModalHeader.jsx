// CreatePostModal/components/ModalHeader/ModalHeader.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { X } from "@/icons";
import styles from '../../CreatePostModal/CreatePostModal.module.css';
export function ModalHeader({ title, lastSaved, isSaving, onClose }) {
    return (
        <div className={styles.logoContainer}>
            <h2 className={styles.logoText}>{title}</h2>
            {lastSaved && (
                <div
                    className={`${styles.autosaveIndicator} ${
                        isSaving ? styles.saving : ""
                    }`}
                    aria-live="polite"
                >
                    {isSaving
                        ? "Saving..."
                        : `Last saved ${new Date(lastSaved).toLocaleTimeString()}`}
                </div>
            )}
            <button
                className={styles.closeButton}
                onClick={onClose}
                aria-label="Close Modal"
            >
                <X size={18}/>
            </button>
        </div>
    );
}
ModalHeader.propTypes = {
    title: PropTypes.string.isRequired,
    lastSaved: PropTypes.instanceOf(Date),
    isSaving: PropTypes.bool,
    onClose: PropTypes.func.isRequired
};
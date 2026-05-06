// File: src/components/Header.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { X } from '@/icons';
// Import the CSS module
import styles from './MessageHeader.module.css';
const Header = ({ title, onClose }) => (
    <header className={styles.messageHeader}>
        <h2 className={styles.headerTitle}>{title}</h2>
        <button
            type="button"
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close"
        >
            <X />
        </button>
    </header>
);
Header.propTypes = {
    title: PropTypes.string.isRequired,
    onClose: PropTypes.func.isRequired,
};
export default memo(Header);

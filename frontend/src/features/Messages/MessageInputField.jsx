// File: src/components/MessageInputField.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
// Import CSS module
import styles from './MessageInputField.module.css';
const MessageInputField = ({ value, onChange, onKeyDown }) => (
    <textarea
        className={styles.messageInput}
        placeholder="Send a message..."
        value={value}
        onChange={onChange}
        onKeyDown={onKeyDown}
        aria-label="Type your message"
    />
);
MessageInputField.propTypes = {
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    onKeyDown: PropTypes.func.isRequired,
};
export default memo(MessageInputField);

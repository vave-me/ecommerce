// File: src/components/Messaging/MessageInput.jsx
"use client";
import React, { useState, useCallback, memo } from 'react';
import { Send } from '@/icons';
import PropTypes from 'prop-types';
import { useNATS } from '../../context/NATSContext';
import styles from './MessageInput.module.css';
const MessageInput = memo(function MessageInput({ onSendMessage, disabled = false }) {
    const [text, setText] = useState('');
    const { connectIfNeeded, isConnected } = useNATS();
    const [hasAttemptedNats, setHasAttemptedNats] = useState(false);
    const [isConnecting, setIsConnecting] = useState(false);
    const [connectionError, setConnectionError] = useState(null);
    const handleSend = useCallback(async () => {
        const trimmedText = text.trim();
        if (!trimmedText) return;
        // Clear any previous connection errors
        setConnectionError(null);
        // Check connection status
        const connectionStatus = Boolean(isConnected);
        // If not connected, try connecting
        if (!connectionStatus) {
            setIsConnecting(true);
            try {
                await connectIfNeeded();
                // Small delay to ensure connection is established
                await new Promise(resolve => setTimeout(resolve, 500));
            } catch (error) {
                setConnectionError("Connection failed. Please check your internet connection.");
                setIsConnecting(false);
                return;
            } finally {
                setIsConnecting(false);
            }
            // If still not connected after attempt
            if (!Boolean(isConnected)) {
                setConnectionError("Unable to establish connection. Please try again.");
                return;
            }
        }
        // Send the message
        try {
            await onSendMessage(trimmedText);
            setText('');
            setConnectionError(null);
        } catch (error) {
            setConnectionError("Failed to send message. Please try again.");
        }
    }, [text, isConnected, connectIfNeeded, onSendMessage]);
    // Typing handler with background connection
    const handleTextChange = useCallback(
        async (e) => {
            const newVal = e.target.value;
            setText(newVal);
            // Auto-connect on first keystroke
            const connectionStatus = Boolean(isConnected);
            if (!hasAttemptedNats && newVal.length === 1 && !connectionStatus) {
                setHasAttemptedNats(true);
                setIsConnecting(true);
                try {
                    await connectIfNeeded();
                    setConnectionError(null);
                } catch (err) {
                    setConnectionError("Connection issue detected.");
                } finally {
                    setIsConnecting(false);
                }
            }
        },
        [hasAttemptedNats, connectIfNeeded, isConnected]
    );
    const handleKeyPress = useCallback(
        (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleSend();
            }
        },
        [handleSend]
    );
    const handleSubmit = useCallback(
        (e) => {
            e.preventDefault();
            handleSend();
        },
        [handleSend]
    );
    // Determine the actual disabled state and placeholder text
    const connectionStatus = Boolean(isConnected);
    const isInputDisabled = disabled || isConnecting;
    const placeholderText = isConnecting 
        ? "Connecting..." 
        : !connectionStatus 
            ? "Connecting to chat..." 
            : "Type a message...";
    return (
        <form className={styles.form} onSubmit={handleSubmit} role="form" aria-label="Send message">
            {/* Connection status indicator */}
            <div 
                className={`${styles.connectionStatus} ${(!connectionStatus || isConnecting || connectionError) ? styles.visible : ''} ${connectionStatus && !connectionError ? styles.connected : styles.disconnected}`}
                role="status"
                aria-live="polite"
            >
                {connectionError 
                    ? connectionError
                    : isConnecting 
                        ? "Connecting..." 
                        : connectionStatus 
                            ? "Connected" 
                            : "Disconnected"
                }
            </div>
            <div className={styles.inputContainer}>
                <textarea
                    className={styles.input}
                    value={text}
                    onChange={handleTextChange}
                    onKeyPress={handleKeyPress}
                    placeholder={placeholderText}
                    disabled={isInputDisabled}
                    rows="1"
                    aria-label="Message input"
                    aria-describedby="message-help"
                    maxLength={2000}
                />
                <div id="message-help" className={styles.visuallyHidden}>
                    Press Enter to send, Shift+Enter for new line
                </div>
                <button
                    type="submit"
                    className={styles.sendButton}
                    disabled={!text.trim() || isInputDisabled}
                    aria-label={text.trim() ? "Send message" : "Enter a message to send"}
                >
                    <Send aria-hidden="true" />
                </button>
            </div>
            {/* Character count indicator */}
            {text.length > 1500 && (
                <div 
                    className={`${styles.characterCount} ${text.length > 1800 ? styles.warning : ''} ${text.length >= 2000 ? styles.limit : ''}`}
                    aria-live="polite"
                >
                    {text.length}/2000
                </div>
            )}
        </form>
    );
});
MessageInput.propTypes = {
    onSendMessage: PropTypes.func.isRequired,
    disabled: PropTypes.bool,
};
MessageInput.defaultProps = {
    disabled: false,
};
export default MessageInput;

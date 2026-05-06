// src/comments/CommentsFullInput.jsx
import { FaPaperPlane, FaPaperclip, FaSmile } from '../../utils/iconImports';
import React, { useState, useCallback } from 'react';
import PropTypes from 'prop-types';
import EmojiPicker from 'emoji-picker-react';
import { useNATS } from "../../context/NATSContext";
import styles from './CommentsFullInput.module.css';
function CommentsFullInput({
                               onSubmitComment,
                               commentText,
                               setCommentText,
                               placeholder = 'Add a comment...',
                               disabled = false,
                               isConnected = true,
                           }) {
    const { connectIfNeeded, isConnected: natsConnected } = useNATS();
    const [hasAttemptedNats, setHasAttemptedNats] = useState(false);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const [fileName, setFileName] = useState(null);
    const handleSubmitComment = useCallback(async () => {
        if (disabled || !commentText?.trim()) return;
        if (!natsConnected) {
            await connectIfNeeded();
            if (!natsConnected) {
                alert('Still not connected to NATS, please try again...');
                return;
            }
        }
        onSubmitComment(commentText.trim());
        if (setCommentText) {
            setCommentText('');
        }
        setShowEmojiPicker(false);
        setFileName(null);
    }, [
        disabled,
        commentText,
        connectIfNeeded,
        natsConnected,
        onSubmitComment,
        setCommentText,
    ]);
    const handleKeyDown = useCallback(
        async (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                await handleSubmitComment();
            }
        },
        [handleSubmitComment]
    );
    const handleEmojiClick = (emojiObject) => {
        if (setCommentText) {
            setCommentText((prev) => (prev || '') + emojiObject.emoji);
        }
    };
    const toggleEmojiPicker = () => {
        setShowEmojiPicker((prev) => !prev);
    };
    const handleFileUpload = (e) => {
        const file = e.target.files[0];
        if (file) {
            setFileName(file.name);
        }
    };
    const handleTextChange = useCallback(
        async (e) => {
            const newValue = e.target.value;
            if (setCommentText) {
                setCommentText(newValue);
            }
            // If user typed the FIRST character & we haven't tried connecting, do so.
            if (!hasAttemptedNats && newValue.length === 1) {
                setHasAttemptedNats(true);
                try {
                    await connectIfNeeded();
                } catch (err) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', err);
        }
        // Continue with default behavior
    }
            }
        },
        [hasAttemptedNats, connectIfNeeded, setCommentText]
    );
    return (
        <div className={styles.inputSection}>
            <button
                className={styles.emojiButton}
                onClick={toggleEmojiPicker}
                aria-label="Add Emoji"
                title="Add Emoji"
                disabled={disabled}
            >
                <FaSmile />
            </button>
            {showEmojiPicker && (
                <div className={styles.emojiPickerWrapper}>
                    <EmojiPicker onEmojiClick={(_, emoji) => handleEmojiClick(emoji)} />
                </div>
            )}
            <label className={styles.attachmentLabel} htmlFor="file-upload" aria-label="Attach File" title="Attach File">
                <FaPaperclip />
            </label>
            <input
                className={styles.hiddenFileInput}
                type="file"
                id="file-upload"
                onChange={handleFileUpload}
                aria-label="Upload Attachment"
            />
            {fileName && <div className={styles.fileNameDisplay}>{fileName}</div>}
            <input
                className={styles.inputField}
                value={commentText || ''}
                onChange={handleTextChange}
                onKeyDown={handleKeyDown}
                aria-label={placeholder}
                disabled={disabled}
                placeholder={placeholder}
            />
            <button
                className={styles.submitButton}
                onClick={handleSubmitComment}
                aria-label="Submit Comment"
                title="Submit Comment"
                disabled={disabled || !commentText?.trim()}
            >
                <FaPaperPlane />
            </button>
        </div>
    );
}
CommentsFullInput.propTypes = {
    onSubmitComment: PropTypes.func.isRequired,
    commentText: PropTypes.string,
    setCommentText: PropTypes.func,
    placeholder: PropTypes.string,
    disabled: PropTypes.bool,
    isConnected: PropTypes.bool,
};
export default CommentsFullInput;
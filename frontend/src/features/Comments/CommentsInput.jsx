// src/comments/CommentsInput.jsx
import React, {useState, useCallback, memo} from 'react';
import PropTypes from 'prop-types';
import {PaperclipIcon, SmileIcon, SendIcon, PlaneIcon} from '@/icons';
import EmojiPicker from 'emoji-picker-react';
import {useNATS} from "../../context/NATSContext";
import styles from './Comments.module.css';
const CommentsInput = memo(function CommentsInput({
                           onSubmitComment,
                           commentText,
                           setCommentText,
                           placeholder = 'Add a comment...',
                           disabled = false,
                           isConnected = true,
                       }) {
    const {connectIfNeeded, isConnected: natsConnected} = useNATS();
    const [hasAttemptedNats, setHasAttemptedNats] = useState(false);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const [fileName, setFileName] = useState(null);
    const handleSubmitComment = useCallback(async () => {
        if (disabled || !commentText?.trim()) return;
        if (!natsConnected) {
            await connectIfNeeded();
            if (!natsConnected) {
                // TODO: Replace with toast notification
            if (typeof window !== 'undefined' && window.toast) {
                window.toast.error('Unable to connect to the messaging service. Please try again.');
            }
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
        setShowEmojiPicker(false);
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
                type="button"
                onClick={toggleEmojiPicker}
                aria-label="Add Emoji"
                title="Add Emoji"
                disabled={disabled}
                className={styles.emojiButton}
            >
                <SmileIcon size={18} />
            </button>
            {showEmojiPicker && (
                <div className={styles.emojiPickerWrapper}>
                    <EmojiPicker onEmojiClick={(emojiData) => handleEmojiClick(emojiData)} />
                </div>
            )}
            <label
                htmlFor="file-upload"
                aria-label="Attach File"
                title="Attach File"
                className={styles.attachmentLabel}
            >
                <PaperclipIcon size={18} />
            </label>
            <input
                type="file"
                id="file-upload"
                onChange={handleFileUpload}
                aria-label="Upload Attachment"
                style={{ display: 'none' }}
            />
            {fileName && <div className={styles.fileNameDisplay}>{fileName}</div>}
            <input
                type="text"
                className={styles.inputField}
                value={commentText || ''}
                onChange={handleTextChange}
                onKeyDown={handleKeyDown}
                aria-label={placeholder}
                disabled={disabled}
                placeholder={placeholder}
            />
            <button
                type="button"
                onClick={handleSubmitComment}
                aria-label="Submit Comment"
                title="Submit Comment"
                disabled={disabled || !commentText?.trim()}
                className={styles.submitButton}
            >
                <SendIcon size={16} />
            </button>
        </div>
    );
});
CommentsInput.displayName = 'CommentsInput';
CommentsInput.propTypes = {
    onSubmitComment: PropTypes.func.isRequired,
    commentText: PropTypes.string,
    setCommentText: PropTypes.func,
    placeholder: PropTypes.string,
    disabled: PropTypes.bool,
    isConnected: PropTypes.bool,
};
export default CommentsInput;
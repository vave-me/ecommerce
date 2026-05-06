"use client";
import { FaPaperPlane, FaPaperclip, FaSmile } from '../../utils/iconImports';
import React, {
    useCallback,
    useEffect,
    useMemo,
    useRef,
    useState,
    memo,
} from "react";
import PropTypes from "prop-types";
import {toast} from "react-toastify";
import {v4 as uuidv4} from "uuid";
import EmojiPicker from "emoji-picker-react";
import {useDispatch} from "react-redux";
import ReactDOM from "react-dom";
import { useSelector } from 'react-redux';
import { Send, Paperclip, X } from '@/icons';
import { debounce } from '../../utils/debounce.js';
import {mes as messages_api} from "../../generated_proto/messages_api_events_pb";
import {message_type} from "../../generated_proto/message_types_pb";
import {jetstream as message_api} from "../../generated_proto/message_api_pb";
import MessageInputField from "./MessageInputField";
import Header from "./MessageHeader";
import {closeMessageModal} from "../../redux/slices/modalsSlice";
import {useConversation} from "../../hooks/useConversation";
import {useNATS} from "../../context/NATSContext"; // includes connectIfNeeded()
import styles from "./MessageForm.module.css";
import {useAuth} from "../../context/AuthContext";
const MAX_FILE_SIZE_MB = 10;
const VALID_FILE_TYPES = ["image/jpeg", "image/png", "image/gif", "application/pdf"];
const DEBOUNCE_DELAY_MS = 300;
function MessageForm({
                         itemId,
                         recipientId,
                         handleSendMessage,
                         metadata,
                         onClose,
                     }) {
    const dispatch = useDispatch();
    const {connectIfNeeded, isConnected, publish} = useNATS();
    const {ensureConversationId} = useConversation(recipientId, itemId);
    const natsName = process.env.NEXT_NATS_SM_NAME || "messenger.SendMessage";
    const {user} = useAuth();
    const [conversationId, setConversationID] = useState("");
    const [messageText, setMessageText] = useState("");
    const [attachments, setAttachments] = useState([]);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const [hasAttemptedNats, setHasAttemptedNats] = useState(false);
    // 3) Refs
    const fileInputRef = useRef(null);
    const emojiPickerRef = useRef(null);
    // Optimize debounce with reference to prevent memory leaks
    const debouncedHandleSendRef = useRef(null);
    // ============= Emoji Picker Logic =============
    const onEmojiClick = useCallback((emoji) => {
        setMessageText((prev) => prev + emoji.emoji);
    }, []);
    const toggleEmojiPicker = useCallback(() => {
        setShowEmojiPicker((prev) => !prev);
    }, []);
    useEffect(() => {
        if (showEmojiPicker && emojiPickerRef.current) {
            emojiPickerRef.current.focus();
        }
    }, [showEmojiPicker]);
    // ============= Attachments =============
    const handleAttachment = useCallback((e) => {
        const rawFiles = Array.from(e.target.files);
        const newFiles = [];
        rawFiles.forEach((file) => {
            if (file.size > MAX_FILE_SIZE_MB * 1024 * 1024) {
                toast.warn(`${file.name} exceeds the ${MAX_FILE_SIZE_MB}MB limit.`);
                return;
            }
            if (!VALID_FILE_TYPES.includes(file.type)) {
                toast.error(`Invalid file type: ${file.name}`);
                return;
            }
            const preview = file.type.startsWith("image/")
                ? URL.createObjectURL(file)
                : null;
            newFiles.push({id: uuidv4(), file, preview});
        });
        setAttachments((prev) => [...prev, ...newFiles]);
        // reset the file input so the same file can be re-uploaded if needed
        e.target.value = null;
    }, []);
    const handleRemoveAttachment = useCallback((attachId) => {
        setAttachments((prev) => {
            const updated = prev.filter((att) => att.id !== attachId);
            const removed = prev.find((att) => att.id === attachId);
            if (removed?.preview) {
                URL.revokeObjectURL(removed.preview);
            }
            return updated;
        });
    }, []);
    // ============= Typing & NATS Connection =============
    const onTextChange = useCallback(
        async (e) => {
            const newValue = e.target.value;
            setMessageText(newValue);
            // Connect to NATS on first keystroke if not attempted yet
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
        [hasAttemptedNats, connectIfNeeded]
    );
    // ============= Sending the Message =============
    const handleSend = useCallback(async () => {
        const text = messageText.trim();
        // 1) Ensure user is logged in
        if (!user || !user.userId) {
            toast.error("You must be logged in to send messages.");
            return;
        }
        // 2) Prevent sending empty messages
        if (!text && attachments.length === 0) {
            toast.error("Cannot send empty message");
            return;
        }
        // 3) Attempt to connect to NATS
        if (!isConnected) {
            await connectIfNeeded();
            if (!isConnected) {
                toast.error("Still not connected to NATS, please wait a moment...");
                return;
            }
        }
        // 4) Obtain or create the conversation
        let finalConvoId = conversationId;
        if (!finalConvoId) {
            try {
                finalConvoId = await ensureConversationId();
                setConversationID(finalConvoId);
            } catch (err) {
                toast.error("Failed to find or create conversation");
                return;
            }
        }
        if (!finalConvoId) {
            toast.error("No conversation ID found");
            return;
        }
        // 5) Build and publish the message
        try {
            const messageId = uuidv4();
            const msgData = messages_api.SendMessage.create({
                id: messageId,
                conversationId: finalConvoId,
                senderId: user.userId,  // <--- from useAuth
                recipientId: recipientId,
                itemId: itemId,
                body: text,
                isRead: false,
            });
            const msgBytes = messages_api.SendMessage.encode(msgData).finish();
            const wsCommand = message_type.WebsocketMessageData.create({
                payload: msgBytes,
                occurred_at: {
                    seconds: Math.floor(Date.now() / 1000),
                    nanos: 0,
                },
            });
            const encCommand = message_type.WebsocketMessageData.encode(wsCommand).finish();
            const streamMessage = message_api.StreamMessage.create({
                id: uuidv4(),
                name: natsName,
                data: encCommand,
                metadata: metadata || {user: user.userId, role: "sender"},
                sent_at: {
                    seconds: Math.floor(Date.now() / 1000),
                    nanos: 0,
                },
            });
            const encodedStreamMessage = message_api.StreamMessage.encode(streamMessage).finish();
            // Publish to NATS
            const subject = `${natsName}.${finalConvoId}`;
            await publish(subject, encodedStreamMessage);
            toast.success("Message sent!");
            // Clear local state
            setMessageText("");
            attachments.forEach((att) => {
                if (att.preview) URL.revokeObjectURL(att.preview);
            });
            setAttachments([]);
            // Fire external callback
            if (handleSendMessage) {
                handleSendMessage({
                    conversationId: finalConvoId,
                    senderId: user.userId,
                    recipientId: recipientId,
                    itemId: itemId,
                    text,
                });
            }
        } catch (err) {
            toast.error("Failed to send message");
        }
    }, [
        user,
        messageText,
        attachments,
        conversationId,
        ensureConversationId,
        connectIfNeeded,
        isConnected,
        publish,
        itemId,
        recipientId,
        metadata,
        handleSendMessage,
    ]);
    // Debounce for hitting Enter quickly - with memoization and cleanup
    useEffect(() => {
        // Create a new debounced function when handleSend changes
        debouncedHandleSendRef.current = debounce(() => {
            handleSend();
        }, DEBOUNCE_DELAY_MS);
        // Cleanup on dependency change or unmount
        return () => {
            if (debouncedHandleSendRef.current) {
                debouncedHandleSendRef.current.cancel();
            }
        };
    }, [handleSend]);
    // Cleanup on unmount
    useEffect(() => {
        return () => {
            // Release object URLs to prevent memory leaks
            attachments.forEach((att) => {
                if (att.preview) URL.revokeObjectURL(att.preview);
            });
        };
    }, [attachments]);
    // ============= Closing the Modal =============
    const handleClose = useCallback(() => {
        attachments.forEach((att) => {
            if (att.preview) URL.revokeObjectURL(att.preview);
        });
        if (onClose) {
            onClose();
        } else {
            dispatch(closeMessageModal());
        }
    }, [attachments, onClose, dispatch]);
    // ============= Render =============
    return ReactDOM.createPortal(
        <div className={styles.modalOverlay} onClick={handleClose}>
            <div
                className={styles.modalContent}
                onClick={(e) => e.stopPropagation()}
            >
                <Header title="Send Message" onClose={handleClose}/>
                <div className={styles.contentWrapper}>
                    {/* MAIN INPUT AREA */}
                    <div className={styles.inputArea}>
                        <MessageInputField
                            value={messageText}
                            onChange={onTextChange}
                            onKeyDown={(e) => {
                                if (e.key === "Enter" && !e.shiftKey) {
                                    e.preventDefault();
                                    if (debouncedHandleSendRef.current) {
                                        debouncedHandleSendRef.current();
                                    }
                                }
                            }}
                        />
                        <button
                            type="button"
                            className={styles.sendButton}
                            onClick={handleSend}
                            aria-label="Submit Message"
                            disabled={!messageText.trim() && attachments.length === 0}
                        >
                            <FaPaperPlane/>
                        </button>
                    </div>
                    {/* TOOLBAR FOR EMOJIS & ATTACHMENTS */}
                    <div className={styles.toolbar}>
                        <button
                            type="button"
                            className={`${styles.iconButton} ${styles.emojiButton}`}
                            onClick={toggleEmojiPicker}
                            aria-label="Toggle Emoji Picker"
                            aria-expanded={showEmojiPicker}
                        >
                            <FaSmile/>
                        </button>
                        {showEmojiPicker && (
                            <div className={styles.emojiPickerContainer} ref={emojiPickerRef}>
                                <EmojiPicker onEmojiClick={onEmojiClick}/>
                            </div>
                        )}
                        <button
                            type="button"
                            className={`${styles.iconButton} ${styles.attachmentButton}`}
                            onClick={() => fileInputRef.current?.click()}
                            aria-label="Attach File"
                        >
                            <FaPaperclip/>
                        </button>
                        <input
                            type="file"
                            ref={fileInputRef}
                            className={styles.hiddenInput}
                            onChange={handleAttachment}
                            multiple
                            accept={VALID_FILE_TYPES.join(",")}
                            aria-label="File Attachment"
                        />
                    </div>
                    {/* ATTACHMENT PREVIEW LIST */}
                    {/*{attachments.length > 0 && (*/}
                    {/*    <AttachmentList*/}
                    {/*        attachments={attachments}*/}
                    {/*        onRemove={handleRemoveAttachment}*/}
                    {/*    />*/}
                    {/*)}*/}
                </div>
            </div>
        </div>,
        document.body
    );
}
MessageForm.propTypes = {
    itemId: PropTypes.string.isRequired,
    recipientId: PropTypes.string.isRequired,
    handleSendMessage: PropTypes.func,
    metadata: PropTypes.object,
    onClose: PropTypes.func,
};
MessageForm.defaultProps = {
    handleSendMessage: null,
    metadata: null,
    onClose: null,
};
export default memo(MessageForm);

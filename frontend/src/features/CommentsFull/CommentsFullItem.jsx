// File: src/comments/CommentItem.jsx
"use client";
import React, { useState, useCallback, memo } from "react";
import PropTypes from "prop-types";
import DOMPurify from "dompurify";
import { useCommentsActions } from "../../hooks/useCommentsActions";
import { useAuth } from "../../context/AuthContext";
import { useUserData } from "../../hooks/useUserData";
import UserAvatar from "../../components/common/UserAvatar";
import CommentsFullInput from "./CommentsFullInput";
// Import the CSS module
import styles from "./CommentFullItem.module.css";
function CommentsFullItem({ comment, depth = 0 }) {
    const { user } = useAuth() || {};
    const userId = user?.userId;
    const { createComment } = useCommentsActions(
        comment.itemId,
        userId,
        comment.categoryId
    );
    const [showReplyInput, setShowReplyInput] = useState(false);
    const [showReplies, setShowReplies] = useState(true);
    const [replyText, setReplyText] = useState("");
    // Fetch sender user data using custom hook
    const { userData: senderUserData, isLoading: isLoadingUser } = useUserData(comment.senderId);
    const handleReplyClick = useCallback(() => {
        if (!userId) {
            alert('Please log in to reply to comments.');
            return;
        }
        setShowReplyInput((prev) => !prev);
    }, [userId]);
    const handleReplySubmit = useCallback(
        (text, parentId) => {
            if (!userId) {
                alert('Please log in to reply to comments.');
                return;
            }
            createComment(text, parentId);
            // Clear local input & hide
            setReplyText("");
            setShowReplyInput(false);
        },
        [createComment, userId]
    );
    const toggleReplies = useCallback(() => {
        setShowReplies((prev) => !prev);
    }, []);
    const handleActionClick = useCallback(
        (actionType) => {
            if (!userId) {
                alert(`Please log in to ${actionType.toLowerCase()} comments.`);
                return;
            }
            // Debug log removed for production
        },
        [comment.id, userId]
    );
    const displayTime = comment.createdAt
        ? new Date(comment.createdAt).toLocaleString()
        : "N/A";
    // Sanitize HTML
    const safeContent = DOMPurify.sanitize(comment.content || "No content");
    const displayName = isLoadingUser ? 'Loading...' : (senderUserData?.userName || 'Anonymous');
    const avatarSrc = senderUserData?.avatar || '/images/user-user.webp';
    return (
        <div
            className={styles.commentContainer}
            style={{ paddingLeft: Math.min(depth * 15, 30) }}
        >
            <div className={styles.commentBox}>
                <div className={styles.commentHeader}>
                    <UserAvatar src={avatarSrc} alt={`${displayName}'s avatar`} />
                    <span className={styles.commentUser}>
                        {displayName}
                    </span>
                    <span className={styles.commentTime}>{displayTime}</span>
                </div>
                <div
                    className={styles.commentContent}
                    dangerouslySetInnerHTML={{ __html: safeContent }}
                />
                <div className={styles.actions}>
                    <button
                        type="button"
                        className={styles.replyButton}
                        onClick={handleReplyClick}
                        title={userId ? `Reply to ${displayName}` : 'Log in to reply'}
                    >
                        Reply
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Like")}
                        title={userId ? `Like ${displayName}'s comment` : 'Log in to like'}
                    >
                        Like
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Dislike")}
                        title={userId ? `Dislike ${displayName}'s comment` : 'Log in to dislike'}
                    >
                        Dislike
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Share")}
                        title={userId ? `Share ${displayName}'s comment` : 'Log in to share'}
                    >
                        Share
                    </button>
                    {comment.replies && comment.replies.length > 0 && (
                        <button
                            type="button"
                            className={`${styles.replyButton} ${styles.viewRepliesButton}`}
                            onClick={toggleReplies}
                        >
                            {showReplies
                                ? "Hide Replies"
                                : `View Replies (${comment.replies.length})`}
                        </button>
                    )}
                </div>
                {showReplyInput && userId && (
                    <div className={styles.replyInputWrapper}>
                        <CommentsFullInput
                            commentText={replyText}
                            setCommentText={setReplyText}
                            onSubmitComment={(text) => handleReplySubmit(text, comment.id)}
                            placeholder={`Reply to ${displayName}...`}
                        />
                    </div>
                )}
                {showReplies && comment.replies && comment.replies.length > 0 && (
                    <div className={styles.repliesList}>
                        {comment.replies.map((reply) => (
                            <CommentsFullItem key={reply.id} comment={reply} depth={depth + 1}/>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
CommentsFullItem.propTypes = {
    comment: PropTypes.shape({
        id: PropTypes.string.isRequired,
        senderId: PropTypes.string,
        itemId: PropTypes.string,
        content: PropTypes.string,
        createdAt: PropTypes.string,
        parentId: PropTypes.string,
        categoryId: PropTypes.string,
        replies: PropTypes.array,
    }).isRequired,
    depth: PropTypes.number,
};
CommentsFullItem.defaultProps = {
    depth: 0,
};
export default memo(CommentsFullItem);

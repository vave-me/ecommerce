// File: src/comments/CommentItem.jsx
"use client";
import React, { useState, useCallback, memo } from "react";
import PropTypes from "prop-types";
import DOMPurify from "dompurify";
import {useCommentsActions} from "../../hooks/useCommentsActions";
import {useAuth} from "../../context/AuthContext";
import {useUserData} from "../../hooks/useUserData";
import UserAvatar from "../../components/common/UserAvatar";
import CommentsInput from "./CommentsInput";
// Import the CSS module
import styles from "./CommentItem.module.css";
function CommentItem({comment, depth = 0}) {
    const {user} = useAuth() || {};
    const userId = user?.userId;
    const {createComment} = useCommentsActions(
        comment.itemId,
        userId,
        comment.categoryId
    );
    const [showReplyInput, setShowReplyInput] = useState(false);
    const [showReplies, setShowReplies] = useState(true);
    const [replyText, setReplyText] = useState("");
    // Fetch sender user data using custom hook
    const { userData: senderUserData, isLoading: isLoadingUser, error } = useUserData(comment.senderId);
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
            style={{'--depth': Math.min(depth, 4)}}
            role="article"
            aria-label={`Comment by ${displayName}`}
        >
            <div className={styles.commentBox}>
                <div className={styles.commentHeader}>
                    <UserAvatar src={avatarSrc} alt={`${displayName}'s avatar`}/>
                    <span 
                        className={styles.commentUser}
                        aria-busy={isLoadingUser}
                        title={senderUserData?.firstName && senderUserData?.lastName 
                            ? `${senderUserData.firstName} ${senderUserData.lastName}` 
                            : displayName}
                    >
                        {displayName}
                    </span>
                    <time 
                        className={styles.commentTime} 
                        dateTime={comment.createdAt}
                        title={displayTime}
                    >
                        {displayTime}
                    </time>
                </div>
                <div
                    className={styles.commentContent}
                    dangerouslySetInnerHTML={{__html: safeContent}}
                    role="region"
                    aria-label="Comment content"
                />
                <div className={styles.actions} role="group" aria-label="Comment actions">
                    <button
                        type="button"
                        className={styles.replyButton}
                        onClick={handleReplyClick}
                        aria-expanded={showReplyInput}
                        aria-label={`Reply to ${displayName}'s comment`}
                        title={userId ? `Reply to ${displayName}` : 'Log in to reply'}
                    >
                        Reply
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Like")}
                        aria-label={`Like ${displayName}'s comment`}
                        title={userId ? `Like ${displayName}'s comment` : 'Log in to like'}
                    >
                        Like
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Dislike")}
                        aria-label={`Dislike ${displayName}'s comment`}
                        title={userId ? `Dislike ${displayName}'s comment` : 'Log in to dislike'}
                    >
                        Dislike
                    </button>
                    <button
                        type="button"
                        className={`${styles.replyButton} ${styles.actionButton}`}
                        onClick={() => handleActionClick("Share")}
                        aria-label={`Share ${displayName}'s comment`}
                        title={userId ? `Share ${displayName}'s comment` : 'Log in to share'}
                    >
                        Share
                    </button>
                    {comment.replies && comment.replies.length > 0 && (
                        <button
                            type="button"
                            className={`${styles.replyButton} ${styles.viewRepliesButton}`}
                            onClick={toggleReplies}
                            aria-expanded={showReplies}
                            aria-label={`${showReplies ? 'Hide' : 'Show'} ${comment.replies.length} replies`}
                        >
                            {showReplies
                                ? "Hide Replies"
                                : `View Replies (${comment.replies.length})`}
                        </button>
                    )}
                </div>
                {showReplyInput && userId && (
                    <div className={styles.replyInputWrapper} role="region" aria-label="Reply input">
                        <CommentsInput
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
                            <CommentItem key={reply.id} comment={reply} depth={depth + 1}/>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
CommentItem.propTypes = {
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
CommentItem.defaultProps = {
    depth: 0,
};
export default memo(CommentItem);

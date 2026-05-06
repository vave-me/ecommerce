// src/components/UserProfile/CommentItem.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import Rating from './Rating';
import styles from './CommentItem.module.css';
import { MessageCircle } from '@/icons';
const CommentItem = memo(({ comment, username }) => {
    return (
        <article className={styles.commentContainer}>
            <header className={styles.header}>
                <h4 className={styles.commenterName}>{username}</h4>
                <time className={styles.commentDate} dateTime={comment.date}>
                    {new Date(comment.date).toLocaleString()}
                </time>
            </header>
            {typeof comment.rating === 'number' && (
                <div className={styles.ratingWrapper}>
                    <Rating value={comment.rating} />
                </div>
            )}
            <div className={styles.details}>
                <div className={styles.commentRow}>
                    <span className={styles.label}>Comment:</span>
                    <span className={styles.commentContent}>{comment.content}</span>
                </div>
                <div className={styles.commentRow}>
                    <span className={styles.label}>Item ID:</span>
                    <span className={styles.commentId}>{comment.itemId}</span>
                </div>
                <div className={styles.commentRow}>
                    <span className={styles.label}>Item Type:</span>
                    <span className={styles.commentType}>{comment.itemType}</span>
                </div>
            </div>
        </article>
    );
});
CommentItem.propTypes = {
    comment: PropTypes.shape({
        id: PropTypes.string.isRequired,
        commenter: PropTypes.string.isRequired,
        rating: PropTypes.number,
        content: PropTypes.string.isRequired,
        itemId: PropTypes.string,
        itemType: PropTypes.string,
        date: PropTypes.string.isRequired,
    }).isRequired,
    username: PropTypes.string.isRequired,
};
CommentItem.displayName = 'CommentItem';
export default CommentItem;

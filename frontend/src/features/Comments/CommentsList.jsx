// src/comments/CommentsList.jsx
import { FaRegCommentDots } from '../../utils/iconImports';
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import CommentItem from './CommentItem';
import styles from './CommentsList.module.css';
function CommentsList({ comments }) {
    if (!comments || comments.length === 0) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyState}>
                    <FaRegCommentDots className={styles.emptyIcon} />
                    <p className={styles.emptyText}>No comments yet. Be the first to comment!</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {comments.map((comment) => (
                <CommentItem key={comment.id} comment={comment} />
            ))}
        </div>
    );
}
CommentsList.propTypes = {
    comments: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            senderId: PropTypes.string,
            itemId: PropTypes.string,
            content: PropTypes.string,
            createdAt: PropTypes.string,
            parentId: PropTypes.string,
            categoryId: PropTypes.string,
            replies: PropTypes.array,
        })
    ),
};
CommentsList.defaultProps = {
    comments: [],
};
export default memo(CommentsList);
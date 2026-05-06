// src/components/UserProfile/CommentList.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import CommentItem from './CommentItem';
import { getCommentsBySender } from '../../api/commentsApi';
import styles from './CommentsList.module.css';
const CommentList = memo(({ userId, username }) => {
    const [comments, setComments] = useState([]);
    useEffect(() => {
        const fetchComments = async () => {
            try {
                const data = await getCommentsBySender(userId);
                setComments(data.comments);
            } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
        };
        fetchComments();
    }, [userId]);
    return (
        <div className={styles.listContainer}>
            {comments.map((comment) => (
                <CommentItem username={username} key={comment.id} comment={comment} />
            ))}
        </div>
    );
});
CommentList.displayName = 'CommentList';
CommentList.propTypes = {
    userId: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
};
export default CommentList;

// src/components/UserProfile/ReviewItem.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Star } from '@/icons';
import styles from './ReviewItem.module.css';
const ReviewItem = memo(({ review, currentUser, onEdit, onDelete }) => {
    return (
        <div className={styles.reviewItem}>
            <div className={styles.header}>
                <div className={styles.rating}>
                    {[...Array(5)].map((_, i) => (
                        <Star
                            key={i}
                            className={i < review.rating ? styles.starFilled : styles.starEmpty}
                        />
                    ))}
                </div>
                <span className={styles.date}>
                    {new Date(review.createdAt).toLocaleDateString()}
                </span>
            </div>
            <p className={styles.comment}>{review.comment}</p>
            <div className={styles.reviewer}>
                <span>by {review.reviewer}</span>
            </div>
        </div>
    );
});
ReviewItem.displayName = 'ReviewItem';
ReviewItem.propTypes = {
    review: PropTypes.shape({
        id: PropTypes.string.isRequired,
        reviewer: PropTypes.string.isRequired,
        rating: PropTypes.number.isRequired,
        comment: PropTypes.string.isRequired,
        date: PropTypes.string.isRequired,
    }).isRequired,
};
export default ReviewItem;

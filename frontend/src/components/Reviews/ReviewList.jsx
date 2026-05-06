// src/components/ReviewList.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import ReviewItem from './ReviewItem';
import styles from './ReviewList.module.css';
/**
 * ReviewList Component
 * Displays a list of reviews.
 * Memoized for performance optimization
 */
const ReviewList = memo(({ reviews, onApprove, onReject, onEdit, onDelete }) => {
    return (
        <ul className={styles.listContainer}>
            {reviews.map((review) => (
                <ReviewItem
                    key={review.id}
                    review={review}
                    onApprove={onApprove}
                    onReject={onReject}
                    onEdit={onEdit}
                    onDelete={onDelete}
                />
            ))}
        </ul>
    );
});
ReviewList.propTypes = {
    reviews: PropTypes.arrayOf(PropTypes.object).isRequired,
    onApprove: PropTypes.func.isRequired,
    onReject: PropTypes.func.isRequired,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default ReviewList;
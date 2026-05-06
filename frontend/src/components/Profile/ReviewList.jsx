// src/components/UserProfile/ReviewList.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import ReviewItem from './ReviewItem';
import { getReviewsByUserId } from "../../api/userApi";
import styles from './ReviewList.module.css';
const ReviewList = memo(({ userId }) => {
    const [reviews, setReviews] = useState([]);
    useEffect(() => {
        const fetchReviews = async () => {
            try {
                const data = await getReviewsByUserId(userId);
                setReviews(data.reviews);
            } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
        };
        fetchReviews();
    }, [userId]);
    return (
        <div className={styles.listContainer}>
            {reviews.map((review) => (
                <ReviewItem key={review.id} review={review} />
            ))}
        </div>
    );
});
ReviewList.displayName = 'ReviewList';
ReviewList.propTypes = {
    userId: PropTypes.string.isRequired,
};
export default ReviewList;

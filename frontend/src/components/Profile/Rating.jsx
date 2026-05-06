// src/components/UserProfile/Rating.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Star } from '@/icons';
import { FaStarHalfAlt, FaRegStar } from '../../utils/iconImports';
import styles from './Rating.module.css';
const Rating = memo(({ value, totalReviews }) => {
    const stars = [];
    const fullStars = Math.floor(value);
    const hasHalfStar = value - fullStars >= 0.5;
    const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);
    for (let i = 0; i < fullStars; i++) {
        stars.push(<Star key={`full-${i}`} />);
    }
    if (hasHalfStar) {
        stars.push(<FaStarHalfAlt key="half" />);
    }
    for (let i = 0; i < emptyStars; i++) {
        stars.push(<FaRegStar key={`empty-${i}`} />);
    }
    return (
        <div className={styles.ratingContainer} aria-label={`Rating: ${value} out of 5`}>
            <div className={styles.stars}>{stars}</div>
            {totalReviews !== null && (
                <span className={styles.reviewCount}>({totalReviews} reviews)</span>
            )}
        </div>
    );
});
Rating.displayName = 'Rating';
Rating.propTypes = {
    value: PropTypes.number.isRequired,
    totalReviews: PropTypes.number,
};
Rating.defaultProps = {
    totalReviews: null,
};
export default Rating;

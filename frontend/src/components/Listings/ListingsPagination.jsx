import React, { memo } from 'react';
import { ChevronLeft, ChevronRight, MoreHorizontal } from '@/icons';
import styles from './ListingsPagination.module.css';
/**
 * ListingsPagination - Atomic Design Component
 * Pagination component for listings
 * 
 * @param {Object} props - Component props
 * @param {number} props.currentPage - Current page number
 * @param {number} props.totalPages - Total number of pages
 * @param {Function} props.onPageChange - Page change handler
 * @param {number} props.totalItems - Total number of items
 * @param {number} props.itemsPerPage - Items per page
 * @returns {JSX.Element} Rendered pagination component
 */
const ListingsPagination = memo(({ 
    currentPage = 1, 
    totalPages = 1, 
    onPageChange = () => {},
    totalItems = 0,
    itemsPerPage = 20
}) => {
    if (totalPages <= 1) return null;
    const getVisiblePages = () => {
        const delta = 2;
        const range = [];
        const rangeWithDots = [];
        for (let i = Math.max(2, currentPage - delta); 
             i <= Math.min(totalPages - 1, currentPage + delta); 
             i++) {
            range.push(i);
        }
        if (currentPage - delta > 2) {
            rangeWithDots.push(1, '...');
        } else {
            rangeWithDots.push(1);
        }
        rangeWithDots.push(...range);
        if (currentPage + delta < totalPages - 1) {
            rangeWithDots.push('...', totalPages);
        } else {
            rangeWithDots.push(totalPages);
        }
        return rangeWithDots;
    };
    const handlePageClick = (page) => {
        if (page !== '...' && page !== currentPage) {
            onPageChange(page);
        }
    };
    const handlePrevious = () => {
        if (currentPage > 1) {
            onPageChange(currentPage - 1);
        }
    };
    const handleNext = () => {
        if (currentPage < totalPages) {
            onPageChange(currentPage + 1);
        }
    };
    const startItem = (currentPage - 1) * itemsPerPage + 1;
    const endItem = Math.min(currentPage * itemsPerPage, totalItems);
    return (
        <div className={styles.container}>
            {/* Results Info */}
            <div className={styles.info}>
                <span className={styles.infoText}>
                    Showing {startItem}-{endItem} of {totalItems} results
                </span>
            </div>
            {/* Pagination Controls */}
            <div className={styles.pagination}>
                {/* Previous Button */}
                <button
                    className={`${styles.pageButton} ${styles.navButton} ${
                        currentPage === 1 ? styles.disabled : ''
                    }`}
                    onClick={handlePrevious}
                    disabled={currentPage === 1}
                    aria-label="Previous page"
                >
                    <ChevronLeft size={16} />
                    <span className={styles.navText}>Previous</span>
                </button>
                {/* Page Numbers */}
                <div className={styles.pageNumbers}>
                    {getVisiblePages().map((page, index) => (
                        <button
                            key={index}
                            className={`${styles.pageButton} ${
                                page === currentPage ? styles.active : ''
                            } ${page === '...' ? styles.dots : ''}`}
                            onClick={() => handlePageClick(page)}
                            disabled={page === '...'}
                            aria-label={page === '...' ? 'More pages' : `Page ${page}`}
                            aria-current={page === currentPage ? 'page' : undefined}
                        >
                            {page === '...' ? <MoreHorizontal size={16} /> : page}
                        </button>
                    ))}
                </div>
                {/* Next Button */}
                <button
                    className={`${styles.pageButton} ${styles.navButton} ${
                        currentPage === totalPages ? styles.disabled : ''
                    }`}
                    onClick={handleNext}
                    disabled={currentPage === totalPages}
                    aria-label="Next page"
                >
                    <span className={styles.navText}>Next</span>
                    <ChevronRight size={16} />
                </button>
            </div>
        </div>
    );
});
ListingsPagination.displayName = 'ListingsPagination';
export default ListingsPagination; 
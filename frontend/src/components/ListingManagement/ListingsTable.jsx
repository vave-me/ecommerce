// src/components/ListingManagement/ListingsTable.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import ListingRow from './ListingRow';
import styles from './ListingsTable.module.css';
/**
 * ListingsTable - Atomic Design Component
 * Displays listings in a responsive table format using CSS modules
 * 
 * @param {Object} props - Component props
 * @param {Array} props.listings - Array of listing objects
 * @param {Function} props.onEdit - Edit handler function
 * @param {Function} props.onDelete - Delete handler function
 * @returns {JSX.Element} Rendered table component
 */
const ListingsTable = memo(({ listings, onEdit, onDelete }) => {
    return (
        <div className={styles.tableContainer}>
            {listings.length === 0 ? (
                <div className={styles.emptyState}>
                    <div className={styles.emptyIcon}>📋</div>
                    <h3 className={styles.emptyTitle}>No listings found</h3>
                    <p className={styles.emptyMessage}>
                        Start by creating your first listing to see it appear here.
                    </p>
                </div>
            ) : (
                <div className={styles.tableWrapper}>
                    <table className={styles.table}>
                        <thead className={styles.tableHead}>
                            <tr className={styles.headerRow}>
                                <th className={styles.headerCell}>Name</th>
                                <th className={styles.headerCell}>Description</th>
                                <th className={styles.headerCell}>Price</th>
                                <th className={styles.headerCell}>Category</th>
                                <th className={styles.headerCell}>Stock</th>
                                <th className={styles.headerCell}>Actions</th>
                            </tr>
                        </thead>
                        <tbody className={styles.tableBody}>
                            {listings.map((listing) => (
                                <ListingRow
                                    key={listing.productId}
                                    listing={listing}
                                    onEdit={onEdit}
                                    onDelete={onDelete}
                                />
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
});
ListingsTable.displayName = 'ListingsTable';
ListingsTable.propTypes = {
    listings: PropTypes.arrayOf(
        PropTypes.shape({
            productId: PropTypes.string.isRequired,
            name: PropTypes.string,
            description: PropTypes.string,
            price: PropTypes.string,
            category: PropTypes.string,
            stock: PropTypes.string,
        })
    ).isRequired,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default ListingsTable;

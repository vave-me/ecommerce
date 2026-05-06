// src/components/ListingManagement/ListingRow.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Edit, Trash2, Eye, DollarSign } from '@/icons';
import styles from './ListingRow.module.css';
/**
 * ListingRow - Atomic Design Component
 * Displays a single listing row in the table using CSS modules
 * 
 * @param {Object} props - Component props
 * @param {Object} props.listing - Listing object
 * @param {Function} props.onEdit - Edit handler function
 * @param {Function} props.onDelete - Delete handler function
 * @returns {JSX.Element} Rendered table row component
 */
const ListingRow = memo(({ listing, onEdit, onDelete }) => {
    const { id, title, price, status, views, createdAt, thumbnail } = listing;
    const formatPrice = (price) => {
        if (!price) return 'N/A';
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    };
    const formatDate = (date) => {
        return new Date(date).toLocaleDateString();
    };
    const getStatusBadge = (status) => {
        const statusClass = status === 'active' ? styles.statusActive : 
                           status === 'pending' ? styles.statusPending : styles.statusInactive;
        return <span className={`${styles.statusBadge} ${statusClass}`}>{status}</span>;
    };
    return (
        <tr className={styles.row}>
            <td className={styles.cell}>
                <div className={styles.titleCell}>
                    {thumbnail && (
                        <img 
                            src={thumbnail} 
                            alt={title}
                            className={styles.thumbnail}
                        />
                    )}
                    <span className={styles.title}>{title}</span>
                </div>
            </td>
            <td className={styles.cell}>
                <div className={styles.priceCell}>
                    <DollarSign size={16} />
                    {formatPrice(price)}
                </div>
            </td>
            <td className={styles.cell}>
                {getStatusBadge(status)}
            </td>
            <td className={styles.cell}>
                <div className={styles.viewsCell}>
                    <Eye size={16} />
                    {views || 0}
                </div>
            </td>
            <td className={styles.cell}>
                {formatDate(createdAt)}
            </td>
            <td className={styles.cell}>
                <div className={styles.actions}>
                    <button
                        className={styles.actionButton}
                        onClick={() => onEdit(id)}
                        aria-label="Edit listing"
                    >
                        <Edit size={16} />
                    </button>
                    <button
                        className={`${styles.actionButton} ${styles.deleteButton}`}
                        onClick={() => onDelete(id)}
                        aria-label="Delete listing"
                    >
                        <Trash2 size={16} />
                    </button>
                </div>
            </td>
        </tr>
    );
});
ListingRow.propTypes = {
    listing: PropTypes.shape({
        id: PropTypes.string.isRequired,
        title: PropTypes.string,
        price: PropTypes.string,
        status: PropTypes.string,
        views: PropTypes.string,
        createdAt: PropTypes.string,
        thumbnail: PropTypes.string,
    }).isRequired,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
ListingRow.displayName = 'ListingRow';
export default ListingRow;

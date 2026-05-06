// src/components/ReviewItem.jsx
"use client"
import { Edit, Trash2 } from '@/icons';
import { FaCheckCircle, FaSpinner, FaTimesCircle } from '../../utils/iconImports';
import React, { useState, memo } from 'react';
import PropTypes from 'prop-types';
import styles from './ReviewItem.module.css';
/**
 * ReviewItem Component
 * Displays individual review details with management actions.
 * Memoized for performance optimization
 */
const ReviewItem = memo(({ review, onApprove, onReject, onEdit, onDelete }) => {
    const { id, reviewer, rating, comment, createdAt, status } = review;
    const [isEditing, setIsEditing] = useState(false);
    const [editedComment, setEditedComment] = useState(comment);
    const [isProcessing, setIsProcessing] = useState(false);
    const handleEditSubmit = () => {
        if (editedComment.trim() === '') {
            alert('Comment cannot be empty.');
            return;
        }
        setIsProcessing(true);
        onEdit(id, editedComment);
        setIsEditing(false);
        setIsProcessing(false);
    };
    return (
        <li className={`${styles.itemContainer} ${styles[status]}`}>
            <img
                className={styles.avatar}
                src={reviewer.avatar}
                alt={`${reviewer.name}'s avatar`}
            />
            <div className={styles.content}>
                <div className={styles.header}>
                    <span className={styles.reviewerName}>{reviewer.name}</span>
                    <span className={styles.rating}>Rating: {rating} / 5</span>
                    <span className={`${styles.statusBadge} ${styles[`status${status.charAt(0).toUpperCase() + status.slice(1)}`]}`}>
                        {status === 'approved' && <FaCheckCircle />}
                        {status === 'approved' && 'Approved'}
                        {status === 'pending' && <FaSpinner className={styles.spinner} />}
                        {status === 'pending' && 'Pending'}
                        {status === 'rejected' && <FaTimesCircle />}
                        {status === 'rejected' && 'Rejected'}
                    </span>
                </div>
                <div className={styles.commentSection}>
                    {isEditing ? (
                        <textarea
                            className={styles.editTextarea}
                            value={editedComment}
                            onChange={(e) => setEditedComment(e.target.value)}
                            aria-label="Edit Review Comment"
                        />
                    ) : (
                        <p className={styles.comment}>{comment}</p>
                    )}
                </div>
                <div className={styles.footer}>
                    <span className={styles.reviewDate}>{new Date(createdAt).toLocaleDateString()}</span>
                    <div className={styles.actions}>
                        {status === 'pending' && (
                            <>
                                <button
                                    className={styles.actionButton}
                                    onClick={() => onApprove(id)}
                                    aria-label="Approve Review"
                                >
                                    <FaCheckCircle /> Approve
                                </button>
                                <button
                                    className={`${styles.actionButton} ${styles.dangerButton}`}
                                    onClick={() => onReject(id)}
                                    aria-label="Reject Review"
                                >
                                    <FaTimesCircle /> Reject
                                </button>
                            </>
                        )}
                        {status === 'approved' && !isEditing && (
                            <button
                                className={styles.actionButton}
                                onClick={() => setIsEditing(true)}
                                aria-label="Edit Review"
                            >
                                <Edit /> Edit
                            </button>
                        )}
                        <button
                            className={`${styles.actionButton} ${styles.deleteButton}`}
                            onClick={() => onDelete(id)}
                            aria-label="Delete Review"
                        >
                            <Trash2 /> Delete
                        </button>
                        {isEditing && (
                            <>
                                <button
                                    className={`${styles.actionButton} ${styles.saveButton}`}
                                    onClick={handleEditSubmit}
                                    aria-label="Save Edited Review"
                                    disabled={isProcessing}
                                >
                                    {isProcessing ? 'Saving...' : 'Save'}
                                </button>
                                <button
                                    className={`${styles.actionButton} ${styles.cancelButton}`}
                                    onClick={() => setIsEditing(false)}
                                    aria-label="Cancel Editing"
                                >
                                    Cancel
                                </button>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </li>
    );
});
ReviewItem.propTypes = {
    review: PropTypes.shape({
        id: PropTypes.string.isRequired,
        reviewer: PropTypes.shape({
            name: PropTypes.string.isRequired,
            avatar: PropTypes.string,
        }).isRequired,
        rating: PropTypes.number.isRequired,
        comment: PropTypes.string.isRequired,
        createdAt: PropTypes.string.isRequired,
        status: PropTypes.oneOf(['approved', 'pending', 'rejected']).isRequired,
    }).isRequired,
    onApprove: PropTypes.func.isRequired,
    onReject: PropTypes.func.isRequired,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default ReviewItem;
"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { createPortal } from 'react-dom';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Search,
  Filter,
  Star,
  AlertCircle,
  CheckCircle,
  XCircle,
  Flag,
  Eye,
  Trash2,
  RefreshCw,
  MessageCircle,
  User,
  Package,
  Calendar,
  ThumbsUp,
  ThumbsDown,
  Shield,
  TrendingUp,
  MoreVertical,
  Reply,
  Ban,
  Archive,
  Heart,
  ExternalLink
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getReviews,
  approveReview,
  rejectReview,
  deleteReview,
  flagReview,
  unflagReview,
  getReviewStats,
  respondToReview
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ReviewManagement.module.css';

// Dummy data for demonstration
const dummyReviews = [
  {
    id: 1,
    rating: 5,
    title: "Excellent product quality!",
    content: "This product exceeded my expectations. The build quality is outstanding and it works exactly as described. Highly recommended!",
    status: 'approved',
    createdAt: '2024-01-15T10:30:00Z',
    author: {
      id: 'user1',
      name: 'John Smith',
      email: 'redacted-email@example.com',
      avatar: null,
      verified: true
    },
    product: {
      id: 'prod1',
      name: 'Premium Wireless Headphones',
      category: 'Electronics',
      thumbnail: null
    },
    helpful: 23,
    reported: 0,
    responses: []
  },
  {
    id: 2,
    rating: 2,
    title: "Poor packaging and delivery",
    content: "The product arrived damaged due to poor packaging. The box was completely crushed and the item inside was broken. Very disappointed with the service.",
    status: 'flagged',
    createdAt: '2024-01-14T15:45:00Z',
    author: {
      id: 'user2',
      name: 'Sarah Johnson',
      email: 'redacted-email@example.com',
      avatar: null,
      verified: false
    },
    product: {
      id: 'prod2',
      name: 'Kitchen Appliance Set',
      category: 'Home & Garden',
      thumbnail: null
    },
    helpful: 5,
    reported: 3,
    responses: [
      {
        id: 'resp1',
        content: 'We apologize for this experience. Please contact our support team for a replacement.',
        author: 'Customer Support',
        createdAt: '2024-01-14T16:00:00Z'
      }
    ]
  },
  {
    id: 3,
    rating: 4,
    title: "Good value for money",
    content: "Overall satisfied with the purchase. The product quality is decent for the price point. Shipping was fast and packaging was secure.",
    status: 'pending',
    createdAt: '2024-01-13T09:20:00Z',
    author: {
      id: 'user3',
      name: 'Mike Chen',
      email: 'redacted-email@example.com',
      avatar: null,
      verified: true
    },
    product: {
      id: 'prod3',
      name: 'Sports Equipment Bundle',
      category: 'Sports & Outdoors',
      thumbnail: null
    },
    helpful: 12,
    reported: 0,
    responses: []
  },
  {
    id: 4,
    rating: 1,
    title: "Completely unusable",
    content: "This product doesn't work at all. I've tried everything but it seems to be defective. Waste of money and time. Would not recommend to anyone.",
    status: 'pending',
    createdAt: '2024-01-12T14:30:00Z',
    author: {
      id: 'user4',
      name: 'Emily Davis',
      email: 'redacted-email@example.com',
      avatar: null,
      verified: false
    },
    product: {
      id: 'prod4',
      name: 'Tech Gadget Pro',
      category: 'Electronics',
      thumbnail: null
    },
    helpful: 8,
    reported: 1,
    responses: []
  },
  {
    id: 5,
    rating: 5,
    title: "Amazing service and quality!",
    content: "Perfect experience from ordering to delivery. The product quality is exceptional and the customer service team was very helpful throughout the process.",
    status: 'approved',
    createdAt: '2024-01-11T11:15:00Z',
    author: {
      id: 'user5',
      name: 'David Wilson',
      email: 'redacted-email@example.com',
      avatar: null,
      verified: true
    },
    product: {
      id: 'prod5',
      name: 'Luxury Watch Collection',
      category: 'Fashion & Accessories',
      thumbnail: null
    },
    helpful: 45,
    reported: 0,
    responses: []
  }
];

const dummyStats = {
  total: 1247,
  approved: 985,
  pending: 156,
  rejected: 78,
  flagged: 28,
  averageRating: 4.2,
  totalHelpful: 8942
};

// Star rating component
const StarRating = ({ rating, size = 14, showValue = true }) => {
  return (
    <div className={styles.starRating}>
      {[1, 2, 3, 4, 5].map((star) => (
        <Star
          key={star}
          size={size}
          className={star <= rating ? styles.starFilled : styles.starEmpty}
          fill={star <= rating ? 'currentColor' : 'none'}
        />
      ))}
      {showValue && <span className={styles.ratingValue}>{rating.toFixed(1)}</span>}
    </div>
  );
};

// Status badge component
const StatusBadge = ({ status }) => {
  const getStatusStyle = () => {
    switch (status) {
      case 'approved':
        return styles.statusApproved;
      case 'pending':
        return styles.statusPending;
      case 'rejected':
        return styles.statusRejected;
      case 'flagged':
        return styles.statusFlagged;
      default:
        return styles.statusPending;
    }
  };

  const getStatusIcon = () => {
    switch (status) {
      case 'approved':
        return <CheckCircle size={12} />;
      case 'pending':
        return <AlertCircle size={12} />;
      case 'rejected':
        return <XCircle size={12} />;
      case 'flagged':
        return <Flag size={12} />;
      default:
        return <AlertCircle size={12} />;
    }
  };

  return (
    <span className={`${styles.statusBadge} ${getStatusStyle()}`}>
      {getStatusIcon()}
      {status.replace('_', ' ')}
    </span>
  );
};

// Review row component
const ReviewRow = ({ review, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, right: 0 });
  const [expanded, setExpanded] = useState(false);

  const formatDate = (date) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const truncateText = (text, maxLength = 120) => {
    if (!text) return '';
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
  };

  return (
    <tr className={styles.reviewRow}>
      <td className={styles.reviewCell}>
        <div className={styles.reviewInfo}>
          <div className={styles.reviewHeader}>
            <div className={styles.reviewTitle}>{review.title}</div>
            <StarRating rating={review.rating} size={12} showValue={false} />
          </div>
          <div className={styles.reviewContent}>
            {expanded ? review.content : truncateText(review.content)}
            {review.content.length > 120 && (
              <button 
                className={styles.expandButton}
                onClick={() => setExpanded(!expanded)}
              >
                {expanded ? 'Show less' : 'Show more'}
              </button>
            )}
          </div>
          <div className={styles.reviewMeta}>
            <span className={styles.metaItem}>
              <Calendar size={12} />
              {formatDate(review.createdAt)}
            </span>
            {review.responses.length > 0 && (
              <span className={styles.metaItem}>
                <Reply size={12} />
                {review.responses.length} response{review.responses.length !== 1 ? 's' : ''}
              </span>
            )}
          </div>
        </div>
      </td>
      <td className={styles.authorCell}>
        <div className={styles.authorInfo}>
          <div className={styles.authorAvatar}>
            {review.author.avatar ? (
              <img src={review.author.avatar} alt={review.author.name} />
            ) : (
              <User size={16} />
            )}
          </div>
          <div className={styles.authorDetails}>
            <div className={styles.authorName}>
              {review.author.name}
              {review.author.verified && (
                <Shield size={12} className={styles.verifiedIcon} />
              )}
            </div>
            <div className={styles.authorEmail}>{review.author.email}</div>
          </div>
        </div>
      </td>
      <td className={styles.productCell}>
        <div className={styles.productInfo}>
          <div className={styles.productThumbnail}>
            {review.product.thumbnail ? (
              <img src={review.product.thumbnail} alt={review.product.name} />
            ) : (
              <Package size={16} />
            )}
          </div>
          <div className={styles.productDetails}>
            <div className={styles.productName}>{review.product.name}</div>
            <div className={styles.productCategory}>{review.product.category}</div>
          </div>
        </div>
      </td>
      <td className={styles.statsCell}>
        <div className={styles.statsInfo}>
          <div className={styles.statItem}>
            <ThumbsUp size={12} />
            <span>{review.helpful}</span>
          </div>
          {review.reported > 0 && (
            <div className={styles.statItem}>
              <Flag size={12} />
              <span>{review.reported}</span>
            </div>
          )}
        </div>
      </td>
      <td className={styles.statusCell}>
        <StatusBadge status={review.status} />
      </td>
      <td className={styles.actionCell}>
        <div className={styles.actionMenu}>
          <button
            className={styles.actionButton}
            onClick={(e) => {
              if (!showMenu) {
                const rect = e.currentTarget.getBoundingClientRect();
                setMenuPosition({
                  top: rect.bottom + window.scrollY,
                  right: window.innerWidth - rect.right
                });
              }
              setShowMenu(!showMenu);
            }}
            aria-label="Review actions"
          >
            <MoreVertical size={14} />
          </button>
          {showMenu && typeof document !== 'undefined' && createPortal(
            <>
              <div
                className={styles.menuOverlay}
                onClick={() => setShowMenu(false)}
              />
              <div 
                className={styles.menuDropdown} 
                style={{ 
                  position: 'fixed',
                  top: `${menuPosition.top}px`,
                  right: `${menuPosition.right}px`,
                  zIndex: 1000
                }}
              >
                <button onClick={() => { onAction('view', review); setShowMenu(false); }}>
                  <Eye size={14} /> View Details
                </button>
                {review.status === 'pending' && (
                  <>
                    <button onClick={() => { onAction('approve', review); setShowMenu(false); }}>
                      <CheckCircle size={14} /> Approve
                    </button>
                    <button onClick={() => { onAction('reject', review); setShowMenu(false); }}>
                      <XCircle size={14} /> Reject
                    </button>
                  </>
                )}
                {review.status !== 'flagged' && (
                  <button onClick={() => { onAction('flag', review); setShowMenu(false); }}>
                    <Flag size={14} /> Flag Review
                  </button>
                )}
                {review.status === 'flagged' && (
                  <button onClick={() => { onAction('unflag', review); setShowMenu(false); }}>
                    <Shield size={14} /> Remove Flag
                  </button>
                )}
                <button onClick={() => { onAction('respond', review); setShowMenu(false); }}>
                  <Reply size={14} /> Respond
                </button>
                <button onClick={() => { onAction('viewProduct', review); setShowMenu(false); }}>
                  <ExternalLink size={14} /> View Product
                </button>
                <button 
                  onClick={() => { onAction('delete', review); setShowMenu(false); }}
                  className={styles.deleteButton}
                >
                  <Trash2 size={14} /> Delete
                </button>
              </div>
            </>,
            document.body
          )}
        </div>
      </td>
    </tr>
  );
};

// Main component
const ReviewManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('ReviewManagement');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }
  
  const router = useRouter();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [filterRating, setFilterRating] = useState('all');
  const [filterCategory, setFilterCategory] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);

  const itemsPerPage = 20;

  // Use dummy data for now
  const reviews = dummyReviews;
  const stats = dummyStats;
  const isLoading = false;
  const error = null;

  // Mutations (placeholder for real implementation)
  const approveMutation = useMutation({
    mutationFn: approveReview,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminReviews']);
    },
  });

  const rejectMutation = useMutation({
    mutationFn: rejectReview,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminReviews']);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteReview,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminReviews']);
    },
  });

  const flagMutation = useMutation({
    mutationFn: flagReview,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminReviews']);
    },
  });

  // Filter reviews
  const filteredReviews = useMemo(() => {
    let filtered = reviews;

    if (searchTerm) {
      filtered = filtered.filter(review =>
        review.title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        review.content?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        review.author?.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        review.product?.name?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter(review => review.status === filterStatus);
    }

    if (filterRating !== 'all') {
      const rating = parseInt(filterRating);
      filtered = filtered.filter(review => Math.floor(review.rating) === rating);
    }

    if (filterCategory !== 'all') {
      filtered = filtered.filter(review => review.product?.category === filterCategory);
    }

    return filtered;
  }, [reviews, searchTerm, filterStatus, filterRating, filterCategory]);

  // Handle actions
  const handleReviewAction = useCallback((action, review) => {
    switch (action) {
      case 'view':
        // Open review details modal
        break;
      case 'approve':
        if (confirm(`Approve review by ${review.author.name}?`)) {
          
          // approveMutation.mutate(review.id);
        }
        break;
      case 'reject':
        if (confirm(`Reject review by ${review.author.name}?`)) {
          
          // rejectMutation.mutate(review.id);
        }
        break;
      case 'flag':
        if (confirm(`Flag this review for moderation?`)) {
          
          // flagMutation.mutate(review.id);
        }
        break;
      case 'unflag':
        
        break;
      case 'respond':
        // Open response modal
        break;
      case 'viewProduct':
        window.open(`/products/${review.product.id}`, '_blank');
        break;
      case 'delete':
        if (confirm(`Are you sure you want to delete this review? This action cannot be undone.`)) {
          
          // deleteMutation.mutate(review.id);
        }
        break;
    }
  }, []);

  // Get unique categories
  const categories = useMemo(() => {
    const categorySet = new Set();
    reviews.forEach(review => {
      if (review.product?.category) {
        categorySet.add(review.product.category);
      }
    });
    return Array.from(categorySet);
  }, [reviews]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>Access Denied</h2>
          <p>You need admin privileges to access this page.</p>
          <p>Current role: {user?.role || 'Not logged in'}</p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>Failed to load reviews</h2>
          <p>{error.message || 'An error occurred while fetching reviews'}</p>
          <button className={styles.retryButton}>
            <RefreshCw size={16} />
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div>
            <h1 className={styles.title}>
              {t('title', { defaultValue: 'Reviews Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Moderate and manage product reviews' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.refreshButton}>
              <RefreshCw size={14} />
            </button>
          </div>
        </div>

        {/* Stats */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <MessageCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.total.toLocaleString()}</span>
              <span className={styles.statLabel}>Total Reviews</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.approved.toLocaleString()}</span>
              <span className={styles.statLabel}>Approved</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <AlertCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.pending.toLocaleString()}</span>
              <span className={styles.statLabel}>Pending</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <Flag size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.flagged.toLocaleString()}</span>
              <span className={styles.statLabel}>Flagged</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(168, 85, 247, 0.1)', color: '#a855f7' }}>
              <Star size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.averageRating.toFixed(1)}</span>
              <span className={styles.statLabel}>Avg Rating</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <ThumbsUp size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.totalHelpful.toLocaleString()}</span>
              <span className={styles.statLabel}>Helpful Votes</span>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.searchBox}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search reviews...' })}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <button
            className={`${styles.filterButton} ${showFilters ? styles.filterActive : ''}`}
            onClick={() => setShowFilters(!showFilters)}
          >
            <Filter size={16} />
            {t('filters', { defaultValue: 'Filters' })}
          </button>
        </div>

        {/* Filter Panel */}
        {showFilters && (
          <div className={styles.filterPanel}>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByStatus', { defaultValue: 'Status' })}</label>
              <select
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Statuses</option>
                <option value="approved">Approved</option>
                <option value="pending">Pending</option>
                <option value="rejected">Rejected</option>
                <option value="flagged">Flagged</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByRating', { defaultValue: 'Rating' })}</label>
              <select
                value={filterRating}
                onChange={(e) => setFilterRating(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Ratings</option>
                <option value="5">5 Stars</option>
                <option value="4">4 Stars</option>
                <option value="3">3 Stars</option>
                <option value="2">2 Stars</option>
                <option value="1">1 Star</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByCategory', { defaultValue: 'Category' })}</label>
              <select
                value={filterCategory}
                onChange={(e) => setFilterCategory(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Categories</option>
                {categories.map(category => (
                  <option key={category} value={category}>{category}</option>
                ))}
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('sortBy', { defaultValue: 'Sort By' })}</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="createdAt">Date Created</option>
                <option value="rating">Rating</option>
                <option value="helpful">Helpful Votes</option>
                <option value="author">Author Name</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('sortOrder', { defaultValue: 'Order' })}</label>
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="desc">Descending</option>
                <option value="asc">Ascending</option>
              </select>
            </div>
          </div>
        )}

        {/* Reviews Table */}
        <div className={styles.tableContainer}>
          <table className={styles.reviewsTable}>
            <thead>
              <tr>
                <th>Review</th>
                <th>Author</th>
                <th>Product</th>
                <th>Stats</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {filteredReviews.length === 0 ? (
                <tr>
                  <td colSpan="6" className={styles.emptyState}>
                    <MessageCircle size={48} />
                    <h3>No reviews found</h3>
                    <p>Try adjusting your filters to see more reviews</p>
                  </td>
                </tr>
              ) : (
                filteredReviews.map(review => (
                  <ReviewRow
                    key={review.id}
                    review={review}
                    onAction={handleReviewAction}
                  />
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className={styles.pagination}>
          <button className={styles.paginationButton} disabled={currentPage === 1}>
            Previous
          </button>
          <span className={styles.paginationInfo}>
            Page {currentPage} of {Math.ceil(filteredReviews.length / itemsPerPage)}
          </span>
          <button 
            className={styles.paginationButton} 
            disabled={currentPage >= Math.ceil(filteredReviews.length / itemsPerPage)}
          >
            Next
          </button>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default ReviewManagement;
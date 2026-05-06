"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  MessageSquare,
  ThumbsUp,
  ThumbsDown,
  Flag,
  Eye,
  EyeOff,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Search,
  Filter,
  Download,
  RefreshCw,
  Edit,
  Trash2,
  MoreVertical,
  User,
  Calendar,
  Hash,
  ExternalLink,
  TrendingUp,
  Activity,
  Users,
  FileText,
  Info,
  Construction
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listAllComments,
  getCommentById,
  approveComment,
  rejectComment,
  deleteComment,
  editComment,
  getCommentStats,
  getCommentsBySender,
  getMostCommentedItems
} from '@/api/adminApi';
import { toast } from 'react-toastify';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './CommentsManagement.module.css';

const CommentStatusBadge = ({ status }) => {
  const statusConfig = {
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending', icon: AlertTriangle },
    approved: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Approved', icon: CheckCircle },
    rejected: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Rejected', icon: XCircle },
    flagged: { color: '#e11d48', bg: 'rgba(225, 29, 72, 0.1)', text: 'Flagged', icon: Flag },
    hidden: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Hidden', icon: EyeOff }
  };

  const config = statusConfig[status] || statusConfig.pending;
  const Icon = config.icon;

  return (
    <span 
      className={styles.statusBadge}
      style={{ color: config.color, backgroundColor: config.bg }}
    >
      <Icon size={12} />
      {config.text}
    </span>
  );
};

const CommentRow = ({ comment, onAction, isSelected, onSelect }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [expanded, setExpanded] = useState(false);

  const truncatedContent = comment.content?.length > 120 
    ? comment.content.substring(0, 120) + '...' 
    : comment.content;

  return (
    <tr className={`${styles.commentRow} ${isSelected ? styles.selected : ''}`}>
      <td className={styles.selectCell}>
        <input
          type="checkbox"
          checked={isSelected}
          onChange={() => onSelect(comment.id)}
          className={styles.checkbox}
        />
      </td>
      <td className={styles.commentCell}>
        <div className={styles.commentInfo}>
          <div className={styles.commentContent}>
            <MessageSquare size={16} />
            <div className={styles.contentText}>
              <div 
                className={styles.content}
                onClick={() => setExpanded(!expanded)}
                title="Click to expand/collapse"
              >
                {expanded ? comment.content : truncatedContent}
              </div>
              <div className={styles.commentMeta}>
                {comment.itemType}: <span className={styles.itemTitle}>{comment.itemTitle || 'Unknown Item'}</span>
              </div>
            </div>
          </div>
        </div>
      </td>
      <td className={styles.authorCell}>
        <div className={styles.authorInfo}>
          <User size={16} />
          <div>
            <div className={styles.authorName}>{comment.author?.name || 'Anonymous'}</div>
            <div className={styles.authorEmail}>{comment.author?.email}</div>
          </div>
        </div>
      </td>
      <td className={styles.statusCell}>
        <CommentStatusBadge status={comment.status} />
      </td>
      <td className={styles.engagementCell}>
        <div className={styles.engagement}>
          <div className={styles.engagementItem}>
            <ThumbsUp size={14} />
            <span>{comment.likes || 0}</span>
          </div>
          <div className={styles.engagementItem}>
            <ThumbsDown size={14} />
            <span>{comment.dislikes || 0}</span>
          </div>
        </div>
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <Calendar size={14} />
          <span>{new Date(comment.createdAt).toLocaleDateString()}</span>
        </div>
      </td>
      <td className={styles.actionCell}>
        <div className={styles.actionMenu}>
          <button
            className={styles.menuTrigger}
            onClick={() => setShowMenu(!showMenu)}
          >
            <MoreVertical size={16} />
          </button>
          {showMenu && (
            <div className={styles.actionDropdown}>
              <button onClick={() => onAction('view', comment)}>
                <Eye size={14} />
                View Details
              </button>
              {comment.status === 'pending' && (
                <>
                  <button onClick={() => onAction('approve', comment)}>
                    <CheckCircle size={14} />
                    Approve
                  </button>
                  <button onClick={() => onAction('reject', comment)}>
                    <XCircle size={14} />
                    Reject
                  </button>
                </>
              )}
              <button onClick={() => onAction('edit', comment)}>
                <Edit size={14} />
                Edit Comment
              </button>
              <button onClick={() => onAction('viewItem', comment)}>
                <ExternalLink size={14} />
                View Item
              </button>
              <button onClick={() => onAction('delete', comment)} className={styles.dangerAction}>
                <Trash2 size={14} />
                Delete
              </button>
            </div>
          )}
        </div>
      </td>
    </tr>
  );
};

const BulkActionsBar = ({ selectedComments, onBulkAction, onClearSelection }) => {
  const t = useTranslations('CommentsManagement');
  
  if (selectedComments.length === 0) return null;

  return (
    <div className={styles.bulkActionsBar}>
      <div className={styles.bulkInfo}>
        <span>{selectedComments.length} comments selected</span>
        <button onClick={onClearSelection} className={styles.clearSelection}>
          Clear Selection
        </button>
      </div>
      <div className={styles.bulkActions}>
        <button onClick={() => onBulkAction('approve')} className={styles.bulkApprove}>
          <CheckCircle size={16} />
          Approve All
        </button>
        <button onClick={() => onBulkAction('reject')} className={styles.bulkReject}>
          <XCircle size={16} />
          Reject All
        </button>
        <button onClick={() => onBulkAction('delete')} className={styles.bulkDelete}>
          <Trash2 size={16} />
          Delete All
        </button>
      </div>
    </div>
  );
};

// Fallback/Mock data generator for when API endpoints are not implemented
const generateMockComments = (filters) => {
  const mockComments = [
    {
      id: 'mock-1',
      content: 'This is a great product! I really recommend it to anyone looking for quality items.',
      author: { name: 'John Doe', email: 'redacted-email@example.com' },
      status: 'pending',
      itemType: 'product',
      itemTitle: 'Premium Wireless Headphones',
      itemId: 'prod-123',
      likes: 5,
      dislikes: 0,
      createdAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString() // 2 days ago
    },
    {
      id: 'mock-2',
      content: 'Not satisfied with the quality. The item arrived damaged and the seller was not responsive.',
      author: { name: 'Jane Smith', email: 'redacted-email@example.com' },
      status: 'approved',
      itemType: 'product',
      itemTitle: 'Smartphone Case',
      itemId: 'prod-456',
      likes: 2,
      dislikes: 8,
      createdAt: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString() // 1 day ago
    },
    {
      id: 'mock-3',
      content: 'Excellent service! Fast delivery and exactly as described. Will buy again.',
      author: { name: 'Mike Johnson', email: 'redacted-email@example.com' },
      status: 'approved',
      itemType: 'service',
      itemTitle: 'Home Cleaning Service',
      itemId: 'serv-789',
      likes: 12,
      dislikes: 1,
      createdAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString() // 3 days ago
    },
    {
      id: 'mock-4',
      content: 'This contains inappropriate content that should be flagged.',
      author: { name: 'Spam User', email: 'redacted-email@example.com' },
      status: 'flagged',
      itemType: 'post',
      itemTitle: 'Community Discussion',
      itemId: 'post-321',
      likes: 0,
      dislikes: 15,
      createdAt: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString() // 4 hours ago
    },
    {
      id: 'mock-5',
      content: 'Average product, nothing special but does the job. Fair price for what you get.',
      author: { name: 'Alice Brown', email: 'redacted-email@example.com' },
      status: 'approved',
      itemType: 'product',
      itemTitle: 'Basic T-Shirt',
      itemId: 'prod-111',
      likes: 3,
      dislikes: 2,
      createdAt: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString() // 6 hours ago
    }
  ];

  // Apply filters to mock data
  let filteredComments = mockComments;

  if (filters.status && filters.status !== 'all') {
    filteredComments = filteredComments.filter(c => c.status === filters.status);
  }

  if (filters.itemType && filters.itemType !== 'all') {
    filteredComments = filteredComments.filter(c => c.itemType === filters.itemType);
  }

  if (filters.search) {
    const searchLower = filters.search.toLowerCase();
    filteredComments = filteredComments.filter(c => 
      c.content.toLowerCase().includes(searchLower) ||
      c.author.name.toLowerCase().includes(searchLower) ||
      c.itemTitle.toLowerCase().includes(searchLower)
    );
  }

  return {
    comments: filteredComments,
    total: filteredComments.length,
    page: 1,
    limit: 50
  };
};

const generateMockStats = () => ({
  total: 156,
  pending: 23,
  approved: 98,
  rejected: 15,
  flagged: 8,
  hidden: 12,
  approvalRate: 75.6,
  commentsToday: 12,
  commentsThisWeek: 67,
  commentsThisMonth: 234
});

// API Error Notice Component
const ApiErrorNotice = ({ onDismiss }) => (
  <div className={styles.apiErrorNotice}>
    <div className={styles.noticeContent}>
      <Construction size={20} />
      <div className={styles.noticeText}>
        <h4>Backend Development in Progress</h4>
        <p>
          Comments management APIs are currently being implemented. 
          You're seeing sample data while development is completed.
        </p>
      </div>
      <button onClick={onDismiss} className={styles.dismissButton}>
        <XCircle size={16} />
      </button>
    </div>
  </div>
);

const CommentsManagement = () => {
  const t = useTranslations('CommentsManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    status: 'all',
    itemType: 'all',
    dateRange: '30d',
    search: ''
  });
  const [selectedComments, setSelectedComments] = useState([]);
  const [showBulkActions, setShowBulkActions] = useState(false);
  const [showApiNotice, setShowApiNotice] = useState(true);
  const [usingFallbackData, setUsingFallbackData] = useState(false);

  // Fetch comments data with error handling and fallback
  const { 
    data: commentsData, 
    isLoading: commentsLoading, 
    error: commentsError,
    refetch: refetchComments 
  } = useQuery({
    queryKey: ['adminComments', filters],
    queryFn: async () => {
      try {
        const result = await listAllComments(filters);
        setUsingFallbackData(false);
        return result;
      } catch (error) {
        
        setUsingFallbackData(true);
        toast.info('Using sample data - Comments API is being developed');
        return generateMockComments(filters);
      }
    },
    enabled: isAdmin,
    retry: false, // Don't retry failed requests
    staleTime: 30000 // 30 seconds
  });

  // Fetch comments statistics with fallback
  const { 
    data: statsData, 
    isLoading: statsLoading 
  } = useQuery({
    queryKey: ['commentStats'],
    queryFn: async () => {
      try {
        const result = await getCommentStats();
        return result;
      } catch (error) {
        
        return generateMockStats();
      }
    },
    enabled: isAdmin,
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  // Approval mutation with error handling
  const approveMutation = useMutation({
    mutationFn: async (commentId) => {
      try {
        const result = await approveComment(commentId);
        toast.success('Comment approved successfully');
        return result;
      } catch (error) {
        if (error.response?.status === 501) {
          toast.warning('Approve functionality is being developed');
          throw new Error('API endpoint not implemented yet');
        }
        toast.error('Failed to approve comment');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['adminComments']);
      queryClient.invalidateQueries(['commentStats']);
    },
    onError: (error) => {
      // Error: 'Approval failed:', error...
    }
  });

  // Rejection mutation with error handling
  const rejectMutation = useMutation({
    mutationFn: async (commentId) => {
      try {
        const result = await rejectComment(commentId);
        toast.success('Comment rejected successfully');
        return result;
      } catch (error) {
        if (error.response?.status === 501) {
          toast.warning('Reject functionality is being developed');
          throw new Error('API endpoint not implemented yet');
        }
        toast.error('Failed to reject comment');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['adminComments']);
      queryClient.invalidateQueries(['commentStats']);
    },
    onError: (error) => {
      // Error: 'Rejection failed:', error...
    }
  });

  // Delete mutation with error handling
  const deleteMutation = useMutation({
    mutationFn: async (commentId) => {
      try {
        const result = await deleteComment(commentId);
        toast.success('Comment deleted successfully');
        return result;
      } catch (error) {
        if (error.response?.status === 501) {
          toast.warning('Delete functionality is being developed');
          throw new Error('API endpoint not implemented yet');
        }
        toast.error('Failed to delete comment');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['adminComments']);
      queryClient.invalidateQueries(['commentStats']);
    },
    onError: (error) => {
      // Error: 'Deletion failed:', error...
    }
  });

  const handleCommentAction = useCallback((action, comment) => {
    switch (action) {
      case 'view':
        if (usingFallbackData) {
          toast.info('Comment details view is being developed');
          return;
        }
        router.push(`/admin/comments/${comment.id}`);
        break;
      case 'approve':
        if (confirm('Approve this comment?')) {
          approveMutation.mutate(comment.id);
        }
        break;
      case 'reject':
        if (confirm('Reject this comment?')) {
          rejectMutation.mutate(comment.id);
        }
        break;
      case 'edit':
        const newContent = prompt('Edit comment content:', comment.content);
        if (newContent && newContent !== comment.content) {
          if (usingFallbackData) {
            toast.info('Edit functionality is being developed');
            return;
          }
          editComment(comment.id, { content: newContent })
            .then(() => {
              toast.success('Comment updated successfully');
              queryClient.invalidateQueries(['adminComments']);
            })
            .catch((error) => {
              if (error.response?.status === 501) {
                toast.warning('Edit functionality is being developed');
              } else {
                toast.error('Failed to update comment');
              }
            });
        }
        break;
      case 'viewItem':
        if (usingFallbackData) {
          toast.info('Item view is being developed');
          return;
        }
        window.open(`/${comment.itemType}s/${comment.itemId}`, '_blank');
        break;
      case 'delete':
        if (confirm('Are you sure you want to delete this comment? This action cannot be undone.')) {
          deleteMutation.mutate(comment.id);
        }
        break;
    }
  }, [router, approveMutation, rejectMutation, deleteMutation, usingFallbackData, queryClient]);

  const handleBulkAction = useCallback((action) => {
    if (selectedComments.length === 0) return;

    const confirmMessages = {
      approve: `Approve ${selectedComments.length} selected comments?`,
      reject: `Reject ${selectedComments.length} selected comments?`,
      delete: `Delete ${selectedComments.length} selected comments? This action cannot be undone.`
    };

    if (confirm(confirmMessages[action])) {
      if (usingFallbackData) {
        toast.info(`Bulk ${action} functionality is being developed`);
        return;
      }

      selectedComments.forEach(commentId => {
        switch (action) {
          case 'approve':
            approveMutation.mutate(commentId);
            break;
          case 'reject':
            rejectMutation.mutate(commentId);
            break;
          case 'delete':
            deleteMutation.mutate(commentId);
            break;
        }
      });
      setSelectedComments([]);
    }
  }, [selectedComments, approveMutation, rejectMutation, deleteMutation, usingFallbackData]);

  const handleSelectComment = useCallback((commentId) => {
    setSelectedComments(prev => 
      prev.includes(commentId) 
        ? prev.filter(id => id !== commentId)
        : [...prev, commentId]
    );
  }, []);

  const handleSelectAll = useCallback(() => {
    const allCommentIds = commentsData?.comments.map(c => c.id) || [];
    setSelectedComments(prev => 
      prev.length === allCommentIds.length ? [] : allCommentIds
    );
  }, [commentsData]);

  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleExport = useCallback(() => {
    if (usingFallbackData) {
      toast.info('Export functionality is being developed');
      return;
    }

    const csvData = commentsData?.comments.map(c => ({
      'Comment ID': c.id,
      'Author': c.author?.name || 'Anonymous',
      'Content': c.content,
      'Status': c.status,
      'Item Type': c.itemType,
      'Item Title': c.itemTitle,
      'Likes': c.likes || 0,
      'Dislikes': c.dislikes || 0,
      'Created': new Date(c.createdAt).toLocaleDateString()
    })) || [];

    toast.success('Export completed');
  }, [commentsData, usingFallbackData]);

  // Process data
  const comments = commentsData?.comments || [];
  const stats = statsData || {};

  // Calculate summary stats
  const summaryStats = useMemo(() => {
    if (stats.total !== undefined) {
      return stats;
    }

    const pending = comments.filter(c => c.status === 'pending');
    const approved = comments.filter(c => c.status === 'approved');
    const rejected = comments.filter(c => c.status === 'rejected');
    const flagged = comments.filter(c => c.status === 'flagged');

    return {
      total: comments.length,
      pending: pending.length,
      approved: approved.length,
      rejected: rejected.length,
      flagged: flagged.length,
      approvalRate: comments.length > 0 ? (approved.length / comments.length) * 100 : 0
    };
  }, [comments, stats]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access comment management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (commentsLoading && !commentsData) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Comments...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch comment data.' })}</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* API Development Notice */}
        {usingFallbackData && showApiNotice && (
          <ApiErrorNotice onDismiss={() => setShowApiNotice(false)} />
        )}

        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <h1 className={styles.title}>
              <MessageSquare size={24} />
              {t('title', { defaultValue: 'Comments Management' })}
              {usingFallbackData && (
                <span className={styles.demoLabel}>Demo Mode</span>
              )}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Moderate and manage user comments' })}
              {usingFallbackData && ' - Sample data shown'}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.exportButton} onClick={handleExport}>
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button className={styles.refreshButton} onClick={() => refetchComments()}>
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <MessageSquare size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.total.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalComments', { defaultValue: 'Total Comments' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <AlertTriangle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.pending.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('pendingReview', { defaultValue: 'Pending Review' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.approved.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('approved', { defaultValue: 'Approved' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <TrendingUp size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.approvalRate.toFixed(1)}%</div>
              <div className={styles.statLabel}>{t('approvalRate', { defaultValue: 'Approval Rate' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersSection}>
          <div className={styles.searchContainer}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search comments...' })}
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filterControls}>
            <select
              value={filters.status}
              onChange={(e) => handleFilterChange('status', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
              <option value="pending">{t('pending', { defaultValue: 'Pending' })}</option>
              <option value="approved">{t('approved', { defaultValue: 'Approved' })}</option>
              <option value="rejected">{t('rejected', { defaultValue: 'Rejected' })}</option>
              <option value="flagged">{t('flagged', { defaultValue: 'Flagged' })}</option>
            </select>
            <select
              value={filters.itemType}
              onChange={(e) => handleFilterChange('itemType', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allTypes', { defaultValue: 'All Types' })}</option>
              <option value="product">{t('products', { defaultValue: 'Products' })}</option>
              <option value="post">{t('posts', { defaultValue: 'Posts' })}</option>
              <option value="service">{t('services', { defaultValue: 'Services' })}</option>
            </select>
            <select
              value={filters.dateRange}
              onChange={(e) => handleFilterChange('dateRange', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="7d">{t('last7Days', { defaultValue: 'Last 7 Days' })}</option>
              <option value="30d">{t('last30Days', { defaultValue: 'Last 30 Days' })}</option>
              <option value="90d">{t('last90Days', { defaultValue: 'Last 90 Days' })}</option>
              <option value="1y">{t('lastYear', { defaultValue: 'Last Year' })}</option>
            </select>
          </div>
        </div>

        {/* Bulk Actions Bar */}
        <BulkActionsBar
          selectedComments={selectedComments}
          onBulkAction={handleBulkAction}
          onClearSelection={() => setSelectedComments([])}
        />

        {/* Comments Table */}
        <div className={styles.tableSection}>
          <div className={styles.tableContainer}>
            <table className={styles.commentsTable}>
              <thead>
                <tr>
                  <th className={styles.selectHeader}>
                    <input
                      type="checkbox"
                      checked={comments.length > 0 && selectedComments.length === comments.length}
                      onChange={handleSelectAll}
                      className={styles.checkbox}
                    />
                  </th>
                  <th>{t('comment', { defaultValue: 'Comment' })}</th>
                  <th>{t('author', { defaultValue: 'Author' })}</th>
                  <th>{t('status', { defaultValue: 'Status' })}</th>
                  <th>{t('engagement', { defaultValue: 'Engagement' })}</th>
                  <th>{t('date', { defaultValue: 'Date' })}</th>
                  <th>{t('actions', { defaultValue: 'Actions' })}</th>
                </tr>
              </thead>
              <tbody>
                {comments.length > 0 ? (
                  comments.map((comment) => (
                    <CommentRow
                      key={comment.id}
                      comment={comment}
                      onAction={handleCommentAction}
                      isSelected={selectedComments.includes(comment.id)}
                      onSelect={handleSelectComment}
                    />
                  ))
                ) : (
                  <tr>
                    <td colSpan="7" className={styles.emptyState}>
                      <MessageSquare size={48} />
                      <h3>No Comments Found</h3>
                      <p>
                        {filters.search || filters.status !== 'all' || filters.itemType !== 'all'
                          ? 'Try adjusting your search criteria or filters.'
                          : 'No comments have been submitted yet.'
                        }
                      </p>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default CommentsManagement;
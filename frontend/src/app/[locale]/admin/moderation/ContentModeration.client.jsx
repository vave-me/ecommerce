"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  ArrowLeft,
  Shield,
  Eye,
  EyeOff,
  CheckCircle,
  XCircle,
  AlertTriangle,
  MessageSquare,
  Image,
  Video,
  FileText,
  User,
  Calendar,
  Filter,
  Search,
  MoreVertical,
  Flag,
  Trash2,
  Edit,
  RefreshCw
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getModerationQueue,
  approveContent,
  rejectContent,
  deleteContent,
  getModerationStats
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ContentModeration.module.css';

const ContentStatusBadge = ({ status }) => {
  const statusConfig = {
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending Review' },
    approved: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Approved' },
    rejected: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Rejected' },
    flagged: { color: '#e11d48', bg: 'rgba(225, 29, 72, 0.1)', text: 'Flagged' }
  };

  const config = statusConfig[status] || statusConfig.pending;

  return (
    <span 
      className={styles.statusBadge}
      style={{ color: config.color, backgroundColor: config.bg }}
    >
      {config.text}
    </span>
  );
};

const ContentTypeIcon = ({ type }) => {
  const icons = {
    post: FileText,
    comment: MessageSquare,
    image: Image,
    video: Video,
    review: MessageSquare
  };
  
  const Icon = icons[type] || FileText;
  return <Icon size={16} />;
};

const ContentItem = ({ content, onApprove, onReject, onDelete, loading }) => {
  const [showPreview, setShowPreview] = useState(false);

  return (
    <div className={styles.contentItem}>
      <div className={styles.contentHeader}>
        <div className={styles.contentType}>
          <ContentTypeIcon type={content.type} />
          <span>{content.type}</span>
        </div>
        <ContentStatusBadge status={content.status} />
      </div>

      <div className={styles.contentBody}>
        <div className={styles.contentMeta}>
          <div className={styles.author}>
            <User size={14} />
            <span>{content.author?.name || 'Anonymous'}</span>
          </div>
          <div className={styles.date}>
            <Calendar size={14} />
            <span>{new Date(content.createdAt).toLocaleDateString()}</span>
          </div>
          {content.reportCount > 0 && (
            <div className={styles.reports}>
              <Flag size={14} />
              <span>{content.reportCount} reports</span>
            </div>
          )}
        </div>

        <div className={styles.contentPreview}>
          <h4>{content.title || 'Untitled Content'}</h4>
          <p className={styles.contentText}>
            {showPreview ? content.content : content.content?.substring(0, 150) + '...'}
          </p>
          {content.content?.length > 150 && (
            <button 
              className={styles.togglePreview}
              onClick={() => setShowPreview(!showPreview)}
            >
              {showPreview ? (
                <>
                  <EyeOff size={14} />
                  Show Less
                </>
              ) : (
                <>
                  <Eye size={14} />
                  Show More
                </>
              )}
            </button>
          )}
        </div>

        {content.mediaUrls && content.mediaUrls.length > 0 && (
          <div className={styles.mediaPreview}>
            {content.mediaUrls.slice(0, 3).map((url, index) => (
              <img key={index} src={url} alt="Content media" className={styles.mediaThumbnail} />
            ))}
            {content.mediaUrls.length > 3 && (
              <div className={styles.moreMedia}>+{content.mediaUrls.length - 3} more</div>
            )}
          </div>
        )}
      </div>

      <div className={styles.contentActions}>
        {content.status === 'pending' && (
          <>
            <button
              className={`${styles.actionButton} ${styles.approve}`}
              onClick={() => onApprove(content.id)}
              disabled={loading}
            >
              <CheckCircle size={16} />
              Approve
            </button>
            <button
              className={`${styles.actionButton} ${styles.reject}`}
              onClick={() => onReject(content.id)}
              disabled={loading}
            >
              <XCircle size={16} />
              Reject
            </button>
          </>
        )}
        <button
          className={`${styles.actionButton} ${styles.delete}`}
          onClick={() => onDelete(content.id)}
          disabled={loading}
        >
          <Trash2 size={16} />
          Delete
        </button>
        <button className={styles.actionButton}>
          <MoreVertical size={16} />
        </button>
      </div>
    </div>
  );
};

const ContentModeration = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('ContentModeration');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    status: 'pending',
    type: 'all',
    searchTerm: '',
    reportedOnly: false
  });
  const [actionLoading, setActionLoading] = useState(null);

  // Fetch moderation data
  const { data: moderationStats, isLoading: statsLoading, refetch: refetchStats } = useQuery({
    queryKey: ['moderation', 'stats'],
    queryFn: getModerationStats,
    staleTime: 2 * 60 * 1000,
    enabled: isAdmin
  });

  const { data: contentQueue, isLoading: queueLoading, refetch: refetchQueue } = useQuery({
    queryKey: ['moderation', 'queue', filters],
    queryFn: () => getModerationQueue(filters),
    staleTime: 30 * 1000,
    enabled: isAdmin
  });

  // Mutations for content actions
  const approveContentMutation = useMutation({
    mutationFn: approveContent,
    onSuccess: () => {
      queryClient.invalidateQueries(['moderation']);
      setActionLoading(null);
    },
    onError: (error) => {
      // Error: 'Error approving content:', error...
      setActionLoading(null);
    }
  });

  const rejectContentMutation = useMutation({
    mutationFn: rejectContent,
    onSuccess: () => {
      queryClient.invalidateQueries(['moderation']);
      setActionLoading(null);
    },
    onError: (error) => {
      // Error: 'Error rejecting content:', error...
      setActionLoading(null);
    }
  });

  const deleteContentMutation = useMutation({
    mutationFn: deleteContent,
    onSuccess: () => {
      queryClient.invalidateQueries(['moderation']);
      setActionLoading(null);
    },
    onError: (error) => {
      // Error: 'Error deleting content:', error...
      setActionLoading(null);
    }
  });

  const handleApprove = useCallback((contentId) => {
    setActionLoading(contentId);
    approveContentMutation.mutate(contentId);
  }, [approveContentMutation]);

  const handleReject = useCallback((contentId) => {
    setActionLoading(contentId);
    rejectContentMutation.mutate(contentId);
  }, [rejectContentMutation]);

  const handleDelete = useCallback((contentId) => {
    if (window.confirm('Are you sure you want to permanently delete this content?')) {
      setActionLoading(contentId);
      deleteContentMutation.mutate(contentId);
    }
  }, [deleteContentMutation]);

  const handleRefresh = useCallback(() => {
    refetchStats();
    refetchQueue();
  }, [refetchStats, refetchQueue]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access content moderation.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  const isLoading = statsLoading || queueLoading;

  if (isLoading && !moderationStats && !contentQueue) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Content Moderation...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch content data.' })}</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <button 
              className={styles.backButton}
              onClick={() => router.back()}
            >
              <ArrowLeft size={16} />
              {t('backToDashboard', { defaultValue: 'Back to Dashboard' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <Shield size={24} />
                {t('title', { defaultValue: 'Content Moderation' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Review and moderate user-generated content' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button 
              className={styles.refreshButton}
              onClick={handleRefresh}
            >
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Grid */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <AlertTriangle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{moderationStats?.pendingReview || 0}</div>
              <div className={styles.statLabel}>{t('pendingReview', { defaultValue: 'Pending Review' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <Flag size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{moderationStats?.flaggedContent || 0}</div>
              <div className={styles.statLabel}>{t('flaggedContent', { defaultValue: 'Flagged Content' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{moderationStats?.approvedToday || 0}</div>
              <div className={styles.statLabel}>{t('approvedToday', { defaultValue: 'Approved Today' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <XCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{moderationStats?.rejectedToday || 0}</div>
              <div className={styles.statLabel}>{t('rejectedToday', { defaultValue: 'Rejected Today' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersPanel}>
          <div className={styles.searchBox}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search content...' })}
              value={filters.searchTerm}
              onChange={(e) => setFilters(prev => ({ ...prev, searchTerm: e.target.value }))}
            />
          </div>
          
          <div className={styles.filterControls}>
            <select
              value={filters.status}
              onChange={(e) => setFilters(prev => ({ ...prev, status: e.target.value }))}
            >
              <option value="pending">{t('pending', { defaultValue: 'Pending' })}</option>
              <option value="approved">{t('approved', { defaultValue: 'Approved' })}</option>
              <option value="rejected">{t('rejected', { defaultValue: 'Rejected' })}</option>
              <option value="flagged">{t('flagged', { defaultValue: 'Flagged' })}</option>
              <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
            </select>
            
            <select
              value={filters.type}
              onChange={(e) => setFilters(prev => ({ ...prev, type: e.target.value }))}
            >
              <option value="all">{t('allTypes', { defaultValue: 'All Types' })}</option>
              <option value="post">{t('posts', { defaultValue: 'Posts' })}</option>
              <option value="comment">{t('comments', { defaultValue: 'Comments' })}</option>
              <option value="review">{t('reviews', { defaultValue: 'Reviews' })}</option>
              <option value="image">{t('images', { defaultValue: 'Images' })}</option>
              <option value="video">{t('videos', { defaultValue: 'Videos' })}</option>
            </select>
            
            <label className={styles.checkboxFilter}>
              <input
                type="checkbox"
                checked={filters.reportedOnly}
                onChange={(e) => setFilters(prev => ({ ...prev, reportedOnly: e.target.checked }))}
              />
              {t('reportedOnly', { defaultValue: 'Reported Only' })}
            </label>
          </div>
        </div>

        {/* Content List */}
        <div className={styles.contentList}>
          {queueLoading ? (
            <div className={styles.loadingState}>
              <LoadingSpinner />
              <p>{t('loadingContent', { defaultValue: 'Loading content...' })}</p>
            </div>
          ) : contentQueue?.items?.length > 0 ? (
            contentQueue.items.map((content) => (
              <ContentItem
                key={content.id}
                content={content}
                onApprove={handleApprove}
                onReject={handleReject}
                onDelete={handleDelete}
                loading={actionLoading === content.id}
              />
            ))
          ) : (
            <div className={styles.emptyState}>
              <Shield size={48} />
              <h3>{t('noContent', { defaultValue: 'No content to moderate' })}</h3>
              <p>{t('noContentMessage', { defaultValue: 'All content has been reviewed or no content matches your filters.' })}</p>
            </div>
          )}
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default ContentModeration; 
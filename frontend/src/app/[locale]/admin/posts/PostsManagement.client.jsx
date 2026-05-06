"use client";

import React, { useState, useCallback, useMemo, lazy, Suspense } from 'react';
import { createPortal } from 'react-dom';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Search,
  Filter,
  Plus,
  Edit2,
  Trash2,
  Archive,
  Eye,
  EyeOff,
  FileText,
  Calendar,
  User,
  Tag,
  TrendingUp,
  RefreshCw,
  MoreVertical,
  AlertCircle,
  CheckCircle,
  Clock,
  Heart,
  MessageSquare,
  Share2,
  Copy,
  ExternalLink
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getPosts,
  getPostById,
  updatePost,
  deletePost,
  archivePost,
  getCategories
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './PostsManagement.module.css';

// Lazy load the CreatePostModal
const CreatePostModal = lazy(() => import('@/features/CreatePostModal'));

// Post status badge component
const StatusBadge = ({ status }) => {
  const getStatusStyle = () => {
    switch (status) {
      case 'published':
        return styles.statusPublished;
      case 'draft':
        return styles.statusDraft;
      case 'archived':
        return styles.statusArchived;
      case 'scheduled':
        return styles.statusScheduled;
      default:
        return styles.statusDraft;
    }
  };

  const getStatusIcon = () => {
    switch (status) {
      case 'published':
        return <CheckCircle size={14} />;
      case 'draft':
        return <Clock size={14} />;
      case 'archived':
        return <Archive size={14} />;
      case 'scheduled':
        return <Calendar size={14} />;
      default:
        return <AlertCircle size={14} />;
    }
  };

  return (
    <span className={`${styles.statusBadge} ${getStatusStyle()}`}>
      {getStatusIcon()}
      {status.replace('_', ' ')}
    </span>
  );
};

// Post row component
const PostRow = ({ post, onAction, categories }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, right: 0 });
  const menuButtonRef = React.useRef(null);
  const category = categories.find(c => c.id === post.categoryId);

  const formatDate = (date) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const truncateText = (text, maxLength = 100) => {
    if (!text) return '';
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
  };

  return (
    <tr className={styles.postRow}>
      <td className={styles.postCell}>
        <div className={styles.postInfo}>
          <div className={styles.postThumbnail}>
            {post.featuredImage ? (
              <img src={post.featuredImage} alt={post.title} />
            ) : (
              <FileText size={24} />
            )}
          </div>
          <div className={styles.postDetails}>
            <div className={styles.postTitle}>{post.title}</div>
            <div className={styles.postMeta}>
              <span className={styles.postSlug}>{post.slug}</span>
              {post.excerpt && (
                <span className={styles.postExcerpt}>{truncateText(post.excerpt, 80)}</span>
              )}
            </div>
          </div>
        </div>
      </td>
      <td className={styles.authorCell}>
        <div className={styles.authorInfo}>
          <div className={styles.authorAvatar}>
            {post.author?.avatar ? (
              <img src={post.author.avatar} alt={post.author.name} />
            ) : (
              <User size={16} />
            )}
          </div>
          <div className={styles.authorDetails}>
            <span className={styles.authorName}>{post.author?.name || 'Unknown'}</span>
            <span className={styles.authorRole}>{post.author?.role || 'Author'}</span>
          </div>
        </div>
      </td>
      <td className={styles.categoryCell}>
        <div className={styles.categoryInfo}>
          <span className={styles.categoryTag}>
            <Tag size={12} />
            {category?.name || 'Uncategorized'}
          </span>
          {post.tags && post.tags.length > 0 && (
            <span className={styles.tagsIndicator}>
              +{post.tags.length} tags
            </span>
          )}
        </div>
      </td>
      <td className={styles.statsCell}>
        <div className={styles.statsInfo}>
          <div className={styles.statItem}>
            <Eye size={12} />
            <span>{post.views || 0}</span>
          </div>
          <div className={styles.statItem}>
            <Heart size={12} />
            <span>{post.likes || 0}</span>
          </div>
          <div className={styles.statItem}>
            <MessageSquare size={12} />
            <span>{post.comments || 0}</span>
          </div>
        </div>
      </td>
      <td className={styles.statusCell}>
        <StatusBadge status={post.status || 'draft'} />
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <div className={styles.dateGroup}>
            <span className={styles.dateLabel}>Created</span>
            <span className={styles.dateValue}>{formatDate(post.createdAt)}</span>
          </div>
          {post.updatedAt !== post.createdAt && (
            <div className={styles.dateGroup}>
              <span className={styles.dateLabel}>Updated</span>
              <span className={styles.dateValue}>{formatDate(post.updatedAt)}</span>
            </div>
          )}
        </div>
      </td>
      <td className={styles.actionCell}>
        <div className={styles.actionMenu}>
          <button
            ref={menuButtonRef}
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
            aria-label="Post actions"
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
                <button onClick={() => { onAction('view', post); setShowMenu(false); }}>
                  <Eye size={14} /> View Post
                </button>
                <button onClick={() => { onAction('edit', post); setShowMenu(false); }}>
                  <Edit2 size={14} /> Edit Post
                </button>
                <button onClick={() => { onAction('copy', post); setShowMenu(false); }}>
                  <Copy size={14} /> Duplicate
                </button>
                <button onClick={() => { onAction('share', post); setShowMenu(false); }}>
                  <Share2 size={14} /> Share
                </button>
                <button onClick={() => { onAction('preview', post); setShowMenu(false); }}>
                  <ExternalLink size={14} /> Preview
                </button>
                {post.status === 'published' && (
                  <button onClick={() => { onAction('unpublish', post); setShowMenu(false); }}>
                    <EyeOff size={14} /> Unpublish
                  </button>
                )}
                {post.status === 'draft' && (
                  <button onClick={() => { onAction('publish', post); setShowMenu(false); }}>
                    <CheckCircle size={14} /> Publish
                  </button>
                )}
                <button onClick={() => { onAction('archive', post); setShowMenu(false); }}>
                  <Archive size={14} /> {post.status === 'archived' ? 'Unarchive' : 'Archive'}
                </button>
                <button 
                  onClick={() => { onAction('delete', post); setShowMenu(false); }}
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

// Quick action modals
const PublishModal = ({ post, onClose, onSave }) => {
  const [publishDate, setPublishDate] = useState(new Date().toISOString().slice(0, 16));
  const [publishNow, setPublishNow] = useState(true);

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({ 
      status: 'published',
      publishedAt: publishNow ? new Date() : new Date(publishDate)
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Publish Post - {post.title}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>
              <input
                type="radio"
                checked={publishNow}
                onChange={() => setPublishNow(true)}
              />
              Publish immediately
            </label>
          </div>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>
              <input
                type="radio"
                checked={!publishNow}
                onChange={() => setPublishNow(false)}
              />
              Schedule for later
            </label>
            {!publishNow && (
              <input
                type="datetime-local"
                value={publishDate}
                onChange={(e) => setPublishDate(e.target.value)}
                min={new Date().toISOString().slice(0, 16)}
              />
            )}
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              {publishNow ? 'Publish Now' : 'Schedule'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Main component
const PostsManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('PostsManagement');
  } catch (e) {
    // Fallback function for missing translations
    t = (key, options) => options?.defaultValue || key;
  }
  
  const router = useRouter();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  
  // Store auth context for debugging
  React.useEffect(() => {
    if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
      window.__REACT_AUTH_CONTEXT__ = { user };
    }
  }, [user]);

  const [searchTerm, setSearchTerm] = useState('');
  const [filterCategory, setFilterCategory] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [filterAuthor, setFilterAuthor] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);
  const [selectedPost, setSelectedPost] = useState(null);
  const [modalType, setModalType] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingPost, setEditingPost] = useState(null);

  const itemsPerPage = 20;

  // Fetch posts
  const { data: postsData, isLoading, error, refetch } = useQuery({
    queryKey: ['adminPosts', currentPage, sortBy, sortOrder],
    queryFn: () => getPosts({
      page: currentPage,
      pageSize: itemsPerPage,
      sortBy,
      sortOrder
    }),
    staleTime: 60000,
    onError: (error) => {
      // Error: 'Failed to fetch posts:', error...
    }
  });

  // Fetch categories
  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: getCategories,
    staleTime: 300000,
  });

  const categories = categoriesData?.categories || [];

  // Mutations
  const deleteMutation = useMutation({
    mutationFn: deletePost,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminPosts']);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, ...updateData }) => updatePost(id, updateData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminPosts']);
      setModalType(null);
      setSelectedPost(null);
    },
  });

  const archiveMutation = useMutation({
    mutationFn: archivePost,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminPosts']);
    },
  });

  // Filter posts
  const filteredPosts = useMemo(() => {
    let posts = postsData?.posts || postsData?.data || [];

    if (searchTerm) {
      posts = posts.filter(post =>
        post.title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        post.content?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        post.excerpt?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        post.author?.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        post.tags?.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
      );
    }

    if (filterCategory !== 'all') {
      posts = posts.filter(post => post.categoryId === filterCategory);
    }

    if (filterStatus !== 'all') {
      posts = posts.filter(post => post.status === filterStatus);
    }

    if (filterAuthor !== 'all') {
      posts = posts.filter(post => post.author?.id === filterAuthor);
    }

    return posts;
  }, [postsData, searchTerm, filterCategory, filterStatus, filterAuthor]);

  // Handle modal close
  const handleModalClose = useCallback(() => {
    setShowCreateModal(false);
    setEditingPost(null);
    queryClient.invalidateQueries(['adminPosts']);
  }, [queryClient]);

  // Handle actions
  const handlePostAction = useCallback((action, post) => {
    switch (action) {
      case 'view':
        window.open(`/posts/${post.slug}`, '_blank');
        break;
      case 'edit':
        setEditingPost(post);
        setShowCreateModal(true);
        break;
      case 'copy':
        setEditingPost({ ...post, id: null, title: `${post.title} (Copy)`, slug: null });
        setShowCreateModal(true);
        break;
      case 'share':
        if (navigator.share) {
          navigator.share({
            title: post.title,
            url: `${window.location.origin}/posts/${post.slug}`
          });
        } else {
          navigator.clipboard.writeText(`${window.location.origin}/posts/${post.slug}`);
        }
        break;
      case 'preview':
        window.open(`/posts/${post.slug}?preview=true`, '_blank');
        break;
      case 'publish':
        setSelectedPost(post);
        setModalType('publish');
        break;
      case 'unpublish':
        if (confirm(`Unpublish "${post.title}"?`)) {
          updateMutation.mutate({ id: post.id, status: 'draft' });
        }
        break;
      case 'archive':
        if (confirm(`Are you sure you want to ${post.status === 'archived' ? 'unarchive' : 'archive'} "${post.title}"?`)) {
          archiveMutation.mutate(post.id);
        }
        break;
      case 'delete':
        if (confirm(`Are you sure you want to delete "${post.title}"? This action cannot be undone.`)) {
          deleteMutation.mutate(post.id);
        }
        break;
    }
  }, [router, updateMutation, archiveMutation, deleteMutation]);

  // Calculate stats
  const stats = useMemo(() => {
    const posts = postsData?.posts || postsData?.data || [];
    return {
      total: postsData?.total || posts.length || 0,
      published: posts.filter(p => p.status === 'published').length,
      draft: posts.filter(p => p.status === 'draft').length,
      archived: posts.filter(p => p.status === 'archived').length,
      scheduled: posts.filter(p => p.status === 'scheduled').length,
      totalViews: posts.reduce((sum, p) => sum + (p.views || 0), 0),
      totalLikes: posts.reduce((sum, p) => sum + (p.likes || 0), 0)
    };
  }, [postsData]);

  // Get unique authors
  const authors = useMemo(() => {
    const posts = postsData?.posts || postsData?.data || [];
    const authorMap = new Map();
    posts.forEach(post => {
      if (post.author) {
        authorMap.set(post.author.id, post.author);
      }
    });
    return Array.from(authorMap.values());
  }, [postsData]);

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
    // Error: 'Posts loading error:', error...
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>Failed to load posts</h2>
          <p>{error.message || 'An error occurred while fetching posts'}</p>
          <button 
            className={styles.retryButton} 
            onClick={() => refetch()}
          >
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
              {t('title', { defaultValue: 'Posts Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage your blog posts and content' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button
              className={styles.addButton}
              onClick={() => setShowCreateModal(true)}
            >
              <Plus size={14} />
              {t('add', { defaultValue: 'New Post' })}
            </button>
            <button
              className={styles.refreshButton}
              onClick={() => refetch()}
            >
              <RefreshCw size={14} />
            </button>
          </div>
        </div>

        {/* Stats */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <FileText size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.total}</span>
              <span className={styles.statLabel}>Total Posts</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.published}</span>
              <span className={styles.statLabel}>Published</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Clock size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.draft}</span>
              <span className={styles.statLabel}>Drafts</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(147, 51, 234, 0.1)', color: '#9333ea' }}>
              <Calendar size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.scheduled}</span>
              <span className={styles.statLabel}>Scheduled</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <Eye size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.totalViews.toLocaleString()}</span>
              <span className={styles.statLabel}>Total Views</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(236, 72, 153, 0.1)', color: '#ec4899' }}>
              <Heart size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.totalLikes.toLocaleString()}</span>
              <span className={styles.statLabel}>Total Likes</span>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.searchBox}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search posts...' })}
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
              <label className={styles.filterLabel}>{t('filterByCategory', { defaultValue: 'Category' })}</label>
              <select
                value={filterCategory}
                onChange={(e) => setFilterCategory(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Categories</option>
                {categories.map(cat => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByStatus', { defaultValue: 'Status' })}</label>
              <select
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Statuses</option>
                <option value="published">Published</option>
                <option value="draft">Draft</option>
                <option value="scheduled">Scheduled</option>
                <option value="archived">Archived</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByAuthor', { defaultValue: 'Author' })}</label>
              <select
                value={filterAuthor}
                onChange={(e) => setFilterAuthor(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Authors</option>
                {authors.map(author => (
                  <option key={author.id} value={author.id}>{author.name}</option>
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
                <option value="updatedAt">Last Updated</option>
                <option value="title">Title</option>
                <option value="views">Views</option>
                <option value="likes">Likes</option>
                <option value="publishedAt">Published Date</option>
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

        {/* Posts Table */}
        <div className={styles.tableContainer}>
          <table className={styles.postsTable}>
            <thead>
              <tr>
                <th>Post</th>
                <th>Author</th>
                <th>Category</th>
                <th>Stats</th>
                <th>Status</th>
                <th>Dates</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {filteredPosts.length === 0 ? (
                <tr>
                  <td colSpan="7" className={styles.emptyState}>
                    <FileText size={48} />
                    <h3>No posts found</h3>
                    <p>Try adjusting your filters or create a new post</p>
                  </td>
                </tr>
              ) : (
                filteredPosts.map(post => (
                  <PostRow
                    key={post.id}
                    post={post}
                    onAction={handlePostAction}
                    categories={categories}
                  />
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {postsData && postsData.totalPages > 1 && (
          <div className={styles.pagination}>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
            >
              Previous
            </button>
            <span className={styles.paginationInfo}>
              Page {currentPage} of {postsData.totalPages}
            </span>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.min(postsData.totalPages, prev + 1))}
              disabled={currentPage === postsData.totalPages}
            >
              Next
            </button>
          </div>
        )}

        {/* Modals */}
        {modalType === 'publish' && selectedPost && (
          <PublishModal
            post={selectedPost}
            onClose={() => {
              setModalType(null);
              setSelectedPost(null);
            }}
            onSave={(publishData) => {
              updateMutation.mutate({ id: selectedPost.id, ...publishData });
            }}
          />
        )}

        {/* Create/Edit Post Modal */}
        {showCreateModal && (
          <Suspense fallback={<LoadingSpinner />}>
            <CreatePostModal
              onClose={handleModalClose}
              editMode={!!editingPost}
              initialPostData={editingPost}
            />
          </Suspense>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default PostsManagement;
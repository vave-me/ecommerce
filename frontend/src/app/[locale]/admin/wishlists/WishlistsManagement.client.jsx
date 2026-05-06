"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Search,
  Filter,
  MoreVertical,
  Heart,
  User,
  Calendar,
  Package,
  Trash2,
  RefreshCw,
  Download,
  Eye,
  TrendingUp,
  BarChart,
  Star
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getAllWishlists,
  getWishlistDetails,
  deleteWishlist
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './WishlistsManagement.module.css';

// Wishlist row component
const WishlistRow = ({ wishlist, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);

  return (
    <tr className={styles.wishlistRow}>
      <td className={styles.wishlistIdCell}>
        <div className={styles.wishlistId}>#{wishlist.id.slice(-8)}</div>
      </td>
      <td className={styles.userCell}>
        <div className={styles.userInfo}>
          <User size={16} />
          <div>
            <div className={styles.userName}>{wishlist.userName || 'Unknown User'}</div>
            <div className={styles.userEmail}>{wishlist.userEmail}</div>
          </div>
        </div>
      </td>
      <td className={styles.nameCell}>
        <div className={styles.wishlistName}>{wishlist.name || 'Untitled Wishlist'}</div>
      </td>
      <td className={styles.itemsCell}>
        <div className={styles.itemsInfo}>
          <Package size={16} />
          <span>{wishlist.itemCount || 0} items</span>
        </div>
      </td>
      <td className={styles.privacyCell}>
        <span className={`${styles.privacyBadge} ${wishlist.isPublic ? styles.privacyPublic : styles.privacyPrivate}`}>
          {wishlist.isPublic ? 'Public' : 'Private'}
        </span>
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <Calendar size={14} />
          <span>{new Date(wishlist.createdAt).toLocaleDateString()}</span>
        </div>
      </td>
      <td className={styles.actionCell}>
        <div className={styles.actionMenu}>
          <button
            className={styles.actionButton}
            onClick={() => setShowMenu(!showMenu)}
            aria-label="Wishlist actions"
          >
            <MoreVertical size={16} />
          </button>
          {showMenu && (
            <>
              <div
                className={styles.menuOverlay}
                onClick={() => setShowMenu(false)}
              />
              <div className={styles.menuDropdown}>
                <button onClick={() => onAction('view', wishlist)}>
                  <Eye size={16} /> View Details
                </button>
                <button onClick={() => onAction('analytics', wishlist)}>
                  <BarChart size={16} /> View Analytics
                </button>
                <button 
                  onClick={() => onAction('delete', wishlist)}
                  className={styles.deleteButton}
                >
                  <Trash2 size={16} /> Delete Wishlist
                </button>
              </div>
            </>
          )}
        </div>
      </td>
    </tr>
  );
};

// Main component
const WishlistsManagement = () => {
  const t = useTranslations('WishlistsManagement');
  const router = useRouter();
  const queryClient = useQueryClient();
  const { isAdmin } = useUserRole();
  
  // Debug: Check if token is available
  React.useEffect(() => {
    // Token check moved to secure storage
  }, []);

  const [searchTerm, setSearchTerm] = useState('');
  const [filterPrivacy, setFilterPrivacy] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);

  const itemsPerPage = 20;

  // Fetch wishlists
  const { data: wishlistsData, isLoading, refetch } = useQuery({
    queryKey: ['adminWishlists', currentPage, sortBy, sortOrder, filterPrivacy],
    queryFn: () => getAllWishlists({
      page: currentPage,
      limit: itemsPerPage,
      sortBy,
      sortOrder,
      isPublic: filterPrivacy === 'all' ? undefined : filterPrivacy === 'public'
    }),
    staleTime: 60000,
  });

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: deleteWishlist,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminWishlists']);
    },
  });

  // Filter data
  const filteredWishlists = useMemo(() => {
    let wishlists = wishlistsData?.wishlists || [];

    if (searchTerm) {
      wishlists = wishlists.filter(wishlist =>
        wishlist.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        wishlist.userName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        wishlist.userEmail?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        wishlist.id.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    return wishlists;
  }, [wishlistsData, searchTerm]);

  // Handle actions
  const handleAction = useCallback((action, wishlist) => {
    switch (action) {
      case 'view':
        router.push(`/admin/wishlists/${wishlist.id}`);
        break;
      case 'analytics':
        router.push(`/admin/wishlists/${wishlist.id}/analytics`);
        break;
      case 'delete':
        if (confirm(`Are you sure you want to delete wishlist "${wishlist.name || 'Untitled'}"?`)) {
          deleteMutation.mutate(wishlist.id);
        }
        break;
    }
  }, [router, deleteMutation]);

  // Calculate stats
  const stats = useMemo(() => {
    const wishlists = wishlistsData?.wishlists || [];
    return {
      totalWishlists: wishlistsData?.total || 0,
      publicWishlists: wishlists.filter(w => w.isPublic).length,
      privateWishlists: wishlists.filter(w => !w.isPublic).length,
      totalItems: wishlists.reduce((sum, w) => sum + (w.itemCount || 0), 0),
      avgItemsPerWishlist: wishlists.length > 0 
        ? Math.round(wishlists.reduce((sum, w) => sum + (w.itemCount || 0), 0) / wishlists.length)
        : 0
    };
  }, [wishlistsData]);

  if (!isAdmin) {
    return null;
  }

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
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
              {t('title', { defaultValue: 'Wishlists Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage and analyze user wishlists' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button
              className={styles.exportButton}
              onClick={() => {/* TODO: Export functionality */}}
            >
              <Download size={16} />
              Export
            </button>
            <button
              className={styles.refreshButton}
              onClick={() => refetch()}
            >
              <RefreshCw size={16} />
            </button>
          </div>
        </div>

        {/* Stats */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Heart size={20} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.totalWishlists}</span>
              <span className={styles.statLabel}>Total Wishlists</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <Eye size={20} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.publicWishlists}</span>
              <span className={styles.statLabel}>Public</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(107, 114, 128, 0.1)', color: '#6b7280' }}>
              <Eye size={20} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.privateWishlists}</span>
              <span className={styles.statLabel}>Private</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <Package size={20} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.totalItems}</span>
              <span className={styles.statLabel}>Total Items</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Star size={20} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.avgItemsPerWishlist}</span>
              <span className={styles.statLabel}>Avg Items/Wishlist</span>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.searchBox}>
            <Search size={20} />
            <input
              type="text"
              placeholder="Search wishlists..."
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
            Filters
          </button>
        </div>

        {/* Filter Panel */}
        {showFilters && (
          <div className={styles.filterPanel}>
            <div className={styles.filterGroup}>
              <label>Privacy</label>
              <select
                value={filterPrivacy}
                onChange={(e) => setFilterPrivacy(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All</option>
                <option value="public">Public Only</option>
                <option value="private">Private Only</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label>Sort By</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="createdAt">Date Created</option>
                <option value="itemCount">Item Count</option>
                <option value="name">Name</option>
                <option value="lastUpdated">Last Updated</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label>Order</label>
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

        {/* Data Table */}
        <div className={styles.tableContainer}>
          <table className={styles.wishlistsTable}>
            <thead>
              <tr>
                <th>Wishlist ID</th>
                <th>User</th>
                <th>Name</th>
                <th>Items</th>
                <th>Privacy</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredWishlists.length === 0 ? (
                <tr>
                  <td colSpan="7" className={styles.emptyState}>
                    <Heart size={48} />
                    <h3>No wishlists found</h3>
                    <p>Try adjusting your filters</p>
                  </td>
                </tr>
              ) : (
                filteredWishlists.map(wishlist => (
                  <WishlistRow
                    key={wishlist.id}
                    wishlist={wishlist}
                    onAction={handleAction}
                  />
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {wishlistsData?.totalPages > 1 && (
          <div className={styles.pagination}>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
            >
              Previous
            </button>
            <span className={styles.paginationInfo}>
              Page {currentPage} of {wishlistsData?.totalPages}
            </span>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.min(wishlistsData?.totalPages, prev + 1))}
              disabled={currentPage === wishlistsData?.totalPages}
            >
              Next
            </button>
          </div>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default WishlistsManagement;
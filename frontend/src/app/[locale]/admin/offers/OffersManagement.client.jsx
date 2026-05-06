"use client";

import React, { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  HandHeart, 
  ShoppingCart,
  Calendar,
  RefreshCw,
  DollarSign,
  TrendingUp,
  Users,
  Package,
  ArrowLeft,
  Plus,
  Filter,
  Search,
  Eye,
  Edit,
  Trash2,
  CheckCircle,
  XCircle,
  Clock,
  AlertCircle,
  Activity
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import LoadingSpinner from '@/components/common/LoadingSpinner';
import { 
  listOffers, 
  getOffersStats, 
  getRecentOffersActivity,
  activateOffer,
  closeOffer 
} from '@/api/offersApi';
import styles from './OffersManagement.module.css';

// Quick Stats Card Component
const StatsCard = ({ title, value, icon: Icon, trend, onClick, loading, className = '' }) => (
  <div 
    className={`${styles.statsCard} ${onClick ? styles.clickable : ''} ${className}`}
    onClick={onClick}
    role={onClick ? 'button' : undefined}
    tabIndex={onClick ? 0 : undefined}
    onKeyDown={onClick ? (e) => e.key === 'Enter' && onClick() : undefined}
  >
    <div className={styles.statsHeader}>
      <h3>{title}</h3>
      <Icon className={styles.statsIcon} aria-hidden="true" />
    </div>
    <div className={styles.statsValue}>
      {loading ? <div className={styles.skeleton}></div> : value}
    </div>
    {trend && (
      <div className={styles.statsTrend}>
        <TrendingUp className={styles.trendIcon} aria-hidden="true" />
        <span>{trend}</span>
      </div>
    )}
  </div>
);

// Offer Item Component
const OfferItem = ({ offer, onActivate, onClose, onView }) => {
  const getStatusColor = (status) => {
    switch (status) {
      case 'active': return styles.statusActive;
      case 'accepted': return styles.statusAccepted;
      case 'closed': return styles.statusClosed;
      case 'draft': return styles.statusDraft;
      default: return styles.statusDefault;
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'active': return CheckCircle;
      case 'accepted': return CheckCircle;
      case 'closed': return XCircle;
      case 'draft': return Clock;
      default: return AlertCircle;
    }
  };

  const StatusIcon = getStatusIcon(offer.offerStatus);

  return (
    <div className={styles.offerItem}>
      <div className={styles.offerMain}>
        <div className={styles.offerInfo}>
          <h4>Offer #{offer.id}</h4>
          <p>Product: {offer.productId}</p>
          <p>Seller: {offer.userSellerId}</p>
          {offer.userCustomerId && <p>Customer: {offer.userCustomerId}</p>}
          <p className={styles.price}>Price: ${(parseInt(offer.price) / 100).toLocaleString()}</p>
        </div>
        <div className={styles.offerStatus}>
          <span className={`${styles.statusBadge} ${getStatusColor(offer.offerStatus)}`}>
            <StatusIcon size={14} />
            {offer.offerStatus}
          </span>
        </div>
      </div>
      <div className={styles.offerActions}>
        <button
          className={styles.actionButton}
          onClick={() => onView(offer.id)}
          title="View Details"
        >
          <Eye size={16} />
        </button>
        {offer.offerStatus === 'draft' && (
          <button
            className={`${styles.actionButton} ${styles.activateButton}`}
            onClick={() => onActivate(offer.id)}
            title="Activate Offer"
          >
            <CheckCircle size={16} />
          </button>
        )}
        {(offer.offerStatus === 'active' || offer.offerStatus === 'draft') && (
          <button
            className={`${styles.actionButton} ${styles.closeButton}`}
            onClick={() => onClose(offer.id)}
            title="Close Offer"
          >
            <XCircle size={16} />
          </button>
        )}
      </div>
    </div>
  );
};

// Main Component
const OffersManagement = () => {
  const router = useRouter();
  const t = useTranslations('OffersManagement');
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [currentPage, setCurrentPage] = useState(1);

  // Data fetching
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['offers-stats'],
    queryFn: getOffersStats,
    refetchInterval: 30000,
  });

  const { data: offersData, isLoading: offersLoading, refetch: refetchOffers } = useQuery({
    queryKey: ['offers', { page: currentPage, status: statusFilter, search: searchTerm }],
    queryFn: () => listOffers({
      page: currentPage,
      limit: 10,
      offerStatus: statusFilter !== 'all' ? statusFilter : undefined,
    }),
    keepPreviousData: true,
  });

  const { data: recentActivity, isLoading: activityLoading } = useQuery({
    queryKey: ['offers-activity'],
    queryFn: getRecentOffersActivity,
    refetchInterval: 60000,
  });

  // Mutations
  const activateMutation = useMutation({
    mutationFn: activateOffer,
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      queryClient.invalidateQueries(['offers-stats']);
    },
  });

  const closeMutation = useMutation({
    mutationFn: (data) => closeOffer(data.offerId, data.reason),
    onSuccess: () => {
      queryClient.invalidateQueries(['offers']);
      queryClient.invalidateQueries(['offers-stats']);
    },
  });

  // Event handlers
  const handleActivateOffer = (offerId) => {
    activateMutation.mutate(offerId);
  };

  const handleCloseOffer = (offerId) => {
    const reason = prompt('Please provide a reason for closing this offer:');
    if (reason) {
      closeMutation.mutate({ offerId, reason });
    }
  };

  const handleViewOffer = (offerId) => {
    router.push(`/admin/offers/${offerId}`);
  };

  const handleRefresh = async () => {
    await Promise.all([
      queryClient.invalidateQueries(['offers-stats']),
      queryClient.invalidateQueries(['offers']),
      queryClient.invalidateQueries(['offers-activity'])
    ]);
  };

  // Filter offers based on search term
  const filteredOffers = useMemo(() => {
    if (!offersData?.offers) return [];
    
    return offersData.offers.filter(offer => {
      const matchesSearch = !searchTerm || 
        offer.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
        offer.productId.toLowerCase().includes(searchTerm.toLowerCase()) ||
        offer.userSellerId.toLowerCase().includes(searchTerm.toLowerCase());
      
      return matchesSearch;
    });
  }, [offersData?.offers, searchTerm]);

  const formatCurrency = (value) => {
    if (!value) return '$0';
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const formatNumber = (value) => {
    if (!value) return '0';
    return new Intl.NumberFormat('en-US').format(value);
  };

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.accessDenied}>
          <h2>Access Denied</h2>
          <p>You need admin privileges to access offers management.</p>
          <p>Current role: {user?.role || 'Not logged in'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerMain}>
          <button 
            onClick={() => router.push('/admin/dashboard')}
            className={styles.backButton}
            aria-label="Back to Dashboard"
          >
            <ArrowLeft size={20} />
          </button>
          <div className={styles.headerInfo}>
            <h1>Offers Management</h1>
            <p>Manage and monitor all offer types and transactions</p>
          </div>
        </div>
        <div className={styles.headerActions}>
          <button 
            onClick={handleRefresh}
            className={styles.refreshButton}
            aria-label="Refresh Data"
          >
            <RefreshCw size={18} />
            Refresh
          </button>
          <button 
            onClick={() => router.push('/admin/offers/create')}
            className={styles.createButton}
          >
            <Plus size={18} />
            Create Offer
          </button>
        </div>
      </header>

      {/* Quick Stats */}
      <section className={styles.statsGrid}>
        <StatsCard
          title="Total Offers"
          value={formatNumber(stats?.totalOffers)}
          icon={HandHeart}
          loading={statsLoading}
          onClick={() => setStatusFilter('all')}
        />
        <StatsCard
          title="Active Offers"
          value={formatNumber(stats?.activeOffers)}
          icon={CheckCircle}
          loading={statsLoading}
          onClick={() => setStatusFilter('active')}
        />
        <StatsCard
          title="Revenue Today"
          value={formatCurrency(stats?.revenueToday)}
          icon={DollarSign}
          loading={statsLoading}
        />
        <StatsCard
          title="Revenue This Month"
          value={formatCurrency(stats?.revenueThisMonth)}
          icon={TrendingUp}
          loading={statsLoading}
        />
      </section>

      {/* Navigation to Offer Types */}
      <section className={styles.offerTypes}>
        <h2>Offer Types</h2>
        <div className={styles.typeGrid}>
          <div 
            className={styles.typeCard}
            onClick={() => router.push('/admin/offers/buynow')}
          >
            <ShoppingCart className={styles.typeIcon} />
            <h3>Buy Now</h3>
            <p>{formatNumber(stats?.totalBuyNow)} transactions</p>
          </div>
          <div 
            className={styles.typeCard}
            onClick={() => router.push('/admin/offers/lease')}
          >
            <Calendar className={styles.typeIcon} />
            <h3>Leases</h3>
            <p>{formatNumber(stats?.totalLeases)} active</p>
          </div>
          <div 
            className={styles.typeCard}
            onClick={() => router.push('/admin/offers/reservation')}
          >
            <Clock className={styles.typeIcon} />
            <h3>Reservations</h3>
            <p>{formatNumber(stats?.totalReservations)} reserved</p>
          </div>
          <div 
            className={styles.typeCard}
            onClick={() => router.push('/admin/offers/buyback')}
          >
            <RefreshCw className={styles.typeIcon} />
            <h3>Buy Back</h3>
            <p>{formatNumber(stats?.totalBuyBacks)} agreements</p>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <div className={styles.mainContent}>
        {/* Filters and Search */}
        <div className={styles.filters}>
          <div className={styles.searchContainer}>
            <Search className={styles.searchIcon} />
            <input
              type="text"
              placeholder="Search offers..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className={styles.statusFilter}
          >
            <option value="all">All Statuses</option>
            <option value="draft">Draft</option>
            <option value="active">Active</option>
            <option value="accepted">Accepted</option>
            <option value="closed">Closed</option>
          </select>
        </div>

        {/* Offers List */}
        <div className={styles.offersSection}>
          <div className={styles.sectionHeader}>
            <h2>Recent Offers</h2>
            <span className={styles.offerCount}>
              {filteredOffers.length} offers
            </span>
          </div>
          
          {offersLoading ? (
            <div className={styles.loadingContainer}>
              <LoadingSpinner />
              <p>Loading offers...</p>
            </div>
          ) : filteredOffers.length === 0 ? (
            <div className={styles.emptyState}>
              <HandHeart className={styles.emptyIcon} />
              <h3>No offers found</h3>
              <p>No offers match your current filters.</p>
            </div>
          ) : (
            <div className={styles.offersList}>
              {filteredOffers.map((offer) => (
                <OfferItem
                  key={offer.id}
                  offer={offer}
                  onActivate={handleActivateOffer}
                  onClose={handleCloseOffer}
                  onView={handleViewOffer}
                />
              ))}
            </div>
          )}
        </div>

        {/* Recent Activity */}
        <aside className={styles.activityPanel}>
          <h3>Recent Activity</h3>
          {activityLoading ? (
            <div className={styles.activityLoading}>Loading activity...</div>
          ) : (
            <div className={styles.activityList}>
              {recentActivity?.activities?.slice(0, 10).map((activity) => (
                <div key={activity.id} className={styles.activityItem}>
                  <div className={styles.activityIcon}>
                    <Activity size={16} />
                  </div>
                  <div className={styles.activityDetails}>
                    <p>{activity.description}</p>
                    <span className={styles.activityTime}>
                      {new Date(activity.timestamp).toLocaleString()}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </aside>
      </div>
    </div>
  );
};

export default OffersManagement; 
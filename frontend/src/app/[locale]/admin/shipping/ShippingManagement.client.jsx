"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Truck,
  Package,
  MapPin,
  Clock,
  CheckCircle,
  AlertCircle,
  XCircle,
  Search,
  Filter,
  Download,
  RefreshCw,
  Eye,
  Edit,
  MoreVertical,
  Navigation,
  Calendar,
  DollarSign,
  User,
  Phone,
  Mail,
  Building,
  Globe,
  Zap,
  TrendingUp,
  Activity,
  Ban,
  UserCheck,
  FileText,
  History
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listShippings,
  getShippingDetails,
  updateShippingStatus,
  trackShipping,
  calculateShippingCost,
  getShippingAnalytics,
  cancelShipment,
  assignCarrier,
  schedulePickup,
  startShipment,
  markShipmentAsDelivered,
  downloadShippingLabel,
  getShipmentHistory
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ShippingManagement.module.css';

const ShippingStatusBadge = ({ status }) => {
  const statusConfig = {
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending', icon: Clock },
    processing: { color: '#3b82f6', bg: 'rgba(59, 130, 246, 0.1)', text: 'Processing', icon: Package },
    shipped: { color: '#8b5cf6', bg: 'rgba(139, 92, 246, 0.1)', text: 'Shipped', icon: Truck },
    in_transit: { color: '#06b6d4', bg: 'rgba(6, 182, 212, 0.1)', text: 'In Transit', icon: Navigation },
    out_for_delivery: { color: '#f97316', bg: 'rgba(249, 115, 22, 0.1)', text: 'Out for Delivery', icon: MapPin },
    delivered: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Delivered', icon: CheckCircle },
    failed: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Failed', icon: XCircle },
    returned: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Returned', icon: AlertCircle }
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

const ShippingRow = ({ shipment, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  return (
    <tr className={styles.shipmentRow}>
      <td className={styles.shipmentCell}>
        <div className={styles.shipmentInfo}>
          <div className={styles.trackingNumber}>
            <Package size={16} />
            <span>{shipment.trackingNumber}</span>
          </div>
          <div className={styles.shipmentMeta}>
            Order #{shipment.orderId}
          </div>
        </div>
      </td>
      <td className={styles.customerCell}>
        <div className={styles.customerInfo}>
          <User size={16} />
          <div>
            <div className={styles.customerName}>{shipment.customer?.name}</div>
            <div className={styles.customerEmail}>{shipment.customer?.email}</div>
          </div>
        </div>
      </td>
      <td className={styles.destinationCell}>
        <div className={styles.destination}>
          <MapPin size={16} />
          <div>
            <div className={styles.city}>{shipment.destination?.city}</div>
            <div className={styles.country}>{shipment.destination?.country}</div>
          </div>
        </div>
      </td>
      <td className={styles.statusCell}>
        <ShippingStatusBadge status={shipment.status} />
      </td>
      <td className={styles.carrierCell}>
        <span className={styles.carrier}>{shipment.carrier}</span>
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <Calendar size={14} />
          <span>{new Date(shipment.createdAt).toLocaleDateString()}</span>
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
              <button onClick={() => onAction('view', shipment)}>
                <Eye size={14} />
                View Details
              </button>
              <button onClick={() => onAction('track', shipment)}>
                <Navigation size={14} />
                Track Package
              </button>
              <button onClick={() => onAction('update', shipment)}>
                <Edit size={14} />
                Update Status
              </button>
              <button onClick={() => onAction('cancel', shipment)}>
                <Ban size={14} />
                Cancel Shipment
              </button>
              <button onClick={() => onAction('label', shipment)}>
                <FileText size={14} />
                Download Label
              </button>
              <button onClick={() => onAction('history', shipment)}>
                <History size={14} />
                View History
              </button>
              {shipment.status === 'out_for_delivery' && (
                <button onClick={() => onAction('deliver', shipment)}>
                  <CheckCircle size={14} />
                  Mark Delivered
                </button>
              )}
            </div>
          )}
        </div>
      </td>
    </tr>
  );
};

const TrackingModal = ({ shipment, onClose }) => {
  const t = useTranslations('ShippingManagement');
  const { data: trackingData, isLoading } = useQuery({
    queryKey: ['tracking', shipment?.trackingNumber],
    queryFn: () => trackShipping(shipment.trackingNumber),
    enabled: !!shipment
  });

  if (!shipment) return null;

  return (
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3>{t('trackingDetails', { defaultValue: 'Tracking Details' })}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <div className={styles.modalBody}>
          <div className={styles.trackingHeader}>
            <div className={styles.trackingNumber}>
              <strong>{t('trackingNumber', { defaultValue: 'Tracking Number' })}: </strong>
              {shipment.trackingNumber}
            </div>
            <div className={styles.currentStatus}>
              <ShippingStatusBadge status={shipment.status} />
            </div>
          </div>
          
          {isLoading ? (
            <div className={styles.trackingLoading}>
              <LoadingSpinner />
              <p>{t('loadingTracking', { defaultValue: 'Loading tracking information...' })}</p>
            </div>
          ) : (
            <div className={styles.trackingTimeline}>
              {trackingData?.events?.map((event, index) => (
                <div key={index} className={styles.timelineEvent}>
                  <div className={styles.eventDot}></div>
                  <div className={styles.eventContent}>
                    <div className={styles.eventTitle}>{event.description}</div>
                    <div className={styles.eventLocation}>{event.location}</div>
                    <div className={styles.eventTime}>
                      {new Date(event.timestamp).toLocaleString()}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const ShippingManagement = () => {
  const t = useTranslations('ShippingManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    status: '',
    carrierId: '',
    dateRange: '30d',
    search: '',
    orderId: '',
    productId: ''
  });
  const [selectedShipment, setSelectedShipment] = useState(null);
  const [showTrackingModal, setShowTrackingModal] = useState(false);
  const [showHistoryModal, setShowHistoryModal] = useState(false);
  const [showDeliveryModal, setShowDeliveryModal] = useState(false);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(20);

  // Fetch shipping data
  const { 
    data: shipmentsData, 
    isLoading: shipmentsLoading, 
    error: shipmentsError,
    refetch: refetchShipments 
  } = useQuery({
    queryKey: ['adminShipping', filters, page, pageSize],
    queryFn: () => listShippings({
      ...filters,
      limit: pageSize,
      offset: page * pageSize,
      status: filters.status || undefined,
      carrierId: filters.carrierId || undefined,
      orderId: filters.orderId || undefined,
      productId: filters.productId || undefined
    }),
    enabled: isAdmin,
    keepPreviousData: true
  });

  // Fetch shipping analytics
  const { 
    data: analyticsData, 
    isLoading: analyticsLoading 
  } = useQuery({
    queryKey: ['shippingAnalytics', filters.dateRange],
    queryFn: () => getShippingAnalytics({ dateRange: filters.dateRange }),
    enabled: isAdmin
  });

  // Update status mutation
  const updateStatusMutation = useMutation({
    mutationFn: ({ shippingId, status, location, notes }) => 
      updateShippingStatus(shippingId, { status, location, notes }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminShipping']);
      queryClient.invalidateQueries(['shippingAnalytics']);
      toast.success('Shipping status updated successfully');
    },
    onError: (error) => {
      
      toast.error(error.response?.data?.message || 'Failed to update shipping status');
    }
  });

  // Cancel shipment mutation
  const cancelShipmentMutation = useMutation({
    mutationFn: ({ shippingId, reason }) => cancelShipment(shippingId, { reason }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminShipping']);
      toast.success('Shipment cancelled successfully');
    },
    onError: (error) => {
      
      toast.error(error.response?.data?.message || 'Failed to cancel shipment');
    }
  });

  // Mark as delivered mutation
  const markDeliveredMutation = useMutation({
    mutationFn: (deliveryData) => markShipmentAsDelivered(deliveryData.shippingId, deliveryData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminShipping']);
      setShowDeliveryModal(false);
      toast.success('Shipment marked as delivered');
    },
    onError: (error) => {
      
      toast.error(error.response?.data?.message || 'Failed to mark as delivered');
    }
  });

  const handleShipmentAction = useCallback((action, shipment) => {
    switch (action) {
      case 'view':
        router.push(`/admin/shipping/${shipment.id}`);
        break;
      case 'track':
        setSelectedShipment(shipment);
        setShowTrackingModal(true);
        break;
      case 'update':
        // Show status update modal
        const statusOptions = ['pending', 'processing', 'shipped', 'in_transit', 'out_for_delivery', 'delivered', 'failed', 'returned'];
        const newStatus = prompt(`Enter new status (${statusOptions.join(', ')}):`, shipment.status);
        if (newStatus && statusOptions.includes(newStatus)) {
          const location = prompt('Enter current location (optional):', '');
          const notes = prompt('Enter notes (optional):', '');
          updateStatusMutation.mutate({
            shippingId: shipment.id,
            status: newStatus,
            location: location || undefined,
            notes: notes || undefined
          });
        } else if (newStatus) {
          alert('Invalid status. Please use one of: ' + statusOptions.join(', '));
        }
        break;
      case 'cancel':
        const reason = prompt('Enter cancellation reason:', '');
        if (reason) {
          if (confirm('Are you sure you want to cancel this shipment?')) {
            cancelShipmentMutation.mutate({
              shippingId: shipment.id,
              reason
            });
          }
        }
        break;
      case 'label':
        downloadShippingLabel(shipment.id, 'pdf');
        break;
      case 'history':
        setSelectedShipment(shipment);
        setShowHistoryModal(true);
        break;
      case 'deliver':
        setSelectedShipment(shipment);
        setShowDeliveryModal(true);
        break;
    }
  }, [router, updateStatusMutation, cancelShipmentMutation]);

  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleExport = useCallback(() => {
    // Export shipments to CSV
    const csvData = shipmentsData?.shipments || [];
    const csvContent = [
      ['Tracking Number', 'Order ID', 'Customer', 'Destination', 'Status', 'Carrier', 'Created'],
      ...csvData.map(s => [
        s.trackingNumber,
        s.orderId,
        s.customer?.name || 'Unknown',
        `${s.destination?.city}, ${s.destination?.country}`,
        s.status,
        s.carrier,
        new Date(s.createdAt).toLocaleDateString()
      ])
    ];

    const blob = new Blob([csvContent.map(row => row.join(',')).join('\n')], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'shipments.csv';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [shipmentsData]);

  // Process data
  const shipments = shipmentsData?.shipments || [];
  const analytics = analyticsData || {};

  // Calculate summary stats
  const stats = useMemo(() => {
    const delivered = shipments.filter(s => s.status === 'delivered');
    const inTransit = shipments.filter(s => ['shipped', 'in_transit', 'out_for_delivery'].includes(s.status));
    const failed = shipments.filter(s => s.status === 'failed');

    return {
      totalShipments: shipments.length,
      delivered: delivered.length,
      inTransit: inTransit.length,
      failed: failed.length,
      deliveryRate: shipments.length > 0 ? (delivered.length / shipments.length) * 100 : 0,
      averageDeliveryTime: analytics.averageDeliveryTime || 0
    };
  }, [shipments, analytics]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access shipping management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (shipmentsLoading && !shipmentsData) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Shipping Data...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch shipping information.' })}</p>
        </div>
      </div>
    );
  }

  if (shipmentsError) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>{t('errorTitle', { defaultValue: 'Failed to Load Shipping Data' })}</h2>
          <p>{shipmentsError.message || t('errorMessage', { defaultValue: 'An error occurred while fetching shipping data' })}</p>
          <button className={styles.retryButton} onClick={() => refetchShipments()}>
            <RefreshCw size={16} />
            {t('retry', { defaultValue: 'Try Again' })}
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
          <div className={styles.headerContent}>
            <h1 className={styles.title}>{t('title', { defaultValue: 'Shipping Management' })}</h1>
            <p className={styles.subtitle}>{t('subtitle', { defaultValue: 'Monitor and manage shipping logistics' })}</p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.exportButton} onClick={handleExport}>
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button className={styles.refreshButton} onClick={() => refetchShipments()}>
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Package size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.totalShipments.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalShipments', { defaultValue: 'Total Shipments' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.delivered.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('delivered', { defaultValue: 'Delivered' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Truck size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.inTransit.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('inTransit', { defaultValue: 'In Transit' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <TrendingUp size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.deliveryRate.toFixed(1)}%</div>
              <div className={styles.statLabel}>{t('deliveryRate', { defaultValue: 'Delivery Rate' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersSection}>
          <div className={styles.searchContainer}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search shipments...' })}
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
              <option value="">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
              <option value="pending">{t('pending', { defaultValue: 'Pending' })}</option>
              <option value="processing">{t('processing', { defaultValue: 'Processing' })}</option>
              <option value="shipped">{t('shipped', { defaultValue: 'Shipped' })}</option>
              <option value="in_transit">{t('inTransit', { defaultValue: 'In Transit' })}</option>
              <option value="delivered">{t('delivered', { defaultValue: 'Delivered' })}</option>
            </select>
            <select
              value={filters.carrierId}
              onChange={(e) => handleFilterChange('carrierId', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="">{t('allCarriers', { defaultValue: 'All Carriers' })}</option>
              <option value="fedex">{t('fedex', { defaultValue: 'FedEx' })}</option>
              <option value="ups">{t('ups', { defaultValue: 'UPS' })}</option>
              <option value="dhl">{t('dhl', { defaultValue: 'DHL' })}</option>
              <option value="usps">{t('usps', { defaultValue: 'USPS' })}</option>
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

        {/* Shipments Table */}
        <div className={styles.tableSection}>
          <div className={styles.tableContainer}>
            <table className={styles.shipmentsTable}>
              <thead>
                <tr>
                  <th>{t('tracking', { defaultValue: 'Tracking' })}</th>
                  <th>{t('customer', { defaultValue: 'Customer' })}</th>
                  <th>{t('destination', { defaultValue: 'Destination' })}</th>
                  <th>{t('status', { defaultValue: 'Status' })}</th>
                  <th>{t('carrier', { defaultValue: 'Carrier' })}</th>
                  <th>{t('date', { defaultValue: 'Date' })}</th>
                  <th>{t('actions', { defaultValue: 'Actions' })}</th>
                </tr>
              </thead>
              <tbody>
                {shipments.map((shipment) => (
                  <ShippingRow
                    key={shipment.id}
                    shipment={shipment}
                    onAction={handleShipmentAction}
                  />
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Pagination */}
        {shipmentsData && (
          <div className={styles.pagination}>
            <button
              onClick={() => setPage(Math.max(0, page - 1))}
              disabled={page === 0}
              className={styles.paginationButton}
            >
              {t('previous', { defaultValue: 'Previous' })}
            </button>
            <span className={styles.paginationInfo}>
              {t('pageInfo', { 
                defaultValue: 'Page {{current}} of {{total}}',
                current: page + 1,
                total: Math.ceil((shipmentsData.total || 0) / pageSize)
              })}
            </span>
            <button
              onClick={() => setPage(page + 1)}
              disabled={(page + 1) * pageSize >= (shipmentsData.total || 0)}
              className={styles.paginationButton}
            >
              {t('next', { defaultValue: 'Next' })}
            </button>
          </div>
        )}

        {/* Tracking Modal */}
        {showTrackingModal && selectedShipment && (
          <TrackingModal
            shipment={selectedShipment}
            onClose={() => {
              setShowTrackingModal(false);
              setSelectedShipment(null);
            }}
          />
        )}

        {/* History Modal */}
        {showHistoryModal && selectedShipment && (
          <HistoryModal
            shipment={selectedShipment}
            onClose={() => {
              setShowHistoryModal(false);
              setSelectedShipment(null);
            }}
          />
        )}

        {/* Delivery Modal */}
        {showDeliveryModal && selectedShipment && (
          <DeliveryModal
            shipment={selectedShipment}
            onDeliver={(data) => {
              markDeliveredMutation.mutate({
                shippingId: selectedShipment.id,
                ...data
              });
            }}
            onClose={() => {
              setShowDeliveryModal(false);
              setSelectedShipment(null);
            }}
          />
        )}
      </div>
    </ErrorBoundary>
  );
};

// History Modal Component
const HistoryModal = ({ shipment, onClose }) => {
  const t = useTranslations('ShippingManagement');
  const { data: historyData, isLoading } = useQuery({
    queryKey: ['shipmentHistory', shipment?.id],
    queryFn: () => getShipmentHistory(shipment.id),
    enabled: !!shipment
  });

  if (!shipment) return null;

  return (
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3>{t('shipmentHistory', { defaultValue: 'Shipment History' })}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <div className={styles.modalBody}>
          <div className={styles.historyHeader}>
            <strong>{t('trackingNumber', { defaultValue: 'Tracking Number' })}: </strong>
            {shipment.trackingNumber}
          </div>
          
          {isLoading ? (
            <div className={styles.loading}>
              <LoadingSpinner />
              <p>{t('loadingHistory', { defaultValue: 'Loading history...' })}</p>
            </div>
          ) : (
            <div className={styles.historyTimeline}>
              {historyData?.events?.map((event, index) => (
                <div key={event.id || index} className={styles.historyEvent}>
                  <div className={styles.eventDot}></div>
                  <div className={styles.eventContent}>
                    <div className={styles.eventTitle}>{event.eventType}</div>
                    <div className={styles.eventDescription}>{event.description}</div>
                    <div className={styles.eventMeta}>
                      <span>{event.location}</span>
                      <span>{new Date(event.timestamp).toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Delivery Modal Component
const DeliveryModal = ({ shipment, onDeliver, onClose }) => {
  const t = useTranslations('ShippingManagement');
  const [deliveryData, setDeliveryData] = useState({
    signedBy: '',
    deliveryTime: new Date().toISOString().slice(0, 16),
    proofOfDeliveryUrl: ''
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!deliveryData.signedBy) {
      alert('Please enter who signed for the delivery');
      return;
    }
    onDeliver(deliveryData);
  };

  if (!shipment) return null;

  return (
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3>{t('markAsDelivered', { defaultValue: 'Mark as Delivered' })}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalBody}>
          <div className={styles.formGroup}>
            <label>{t('signedBy', { defaultValue: 'Signed By' })} *</label>
            <input
              type="text"
              value={deliveryData.signedBy}
              onChange={(e) => setDeliveryData({ ...deliveryData, signedBy: e.target.value })}
              required
              placeholder={t('recipientName', { defaultValue: 'Recipient name' })}
            />
          </div>
          <div className={styles.formGroup}>
            <label>{t('deliveryTime', { defaultValue: 'Delivery Time' })}</label>
            <input
              type="datetime-local"
              value={deliveryData.deliveryTime}
              onChange={(e) => setDeliveryData({ ...deliveryData, deliveryTime: e.target.value })}
            />
          </div>
          <div className={styles.formGroup}>
            <label>{t('proofUrl', { defaultValue: 'Proof of Delivery URL' })}</label>
            <input
              type="url"
              value={deliveryData.proofOfDeliveryUrl}
              onChange={(e) => setDeliveryData({ ...deliveryData, proofOfDeliveryUrl: e.target.value })}
              placeholder={t('optionalUrl', { defaultValue: 'Optional URL' })}
            />
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              {t('cancel', { defaultValue: 'Cancel' })}
            </button>
            <button type="submit" className={styles.submitButton}>
              {t('confirmDelivery', { defaultValue: 'Confirm Delivery' })}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ShippingManagement; 
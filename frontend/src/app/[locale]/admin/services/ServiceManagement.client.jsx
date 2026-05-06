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
  Briefcase,
  DollarSign,
  TrendingUp,
  TrendingDown,
  AlertCircle,
  Download,
  Upload,
  RefreshCw,
  MoreVertical,
  Tag,
  Calendar,
  Image as ImageIcon,
  CheckCircle,
  XCircle,
  Clock,
  Users,
  Settings,
  Star,
  MapPin
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getEntities,
  getEntity,
  createEntity,
  updateEntity,
  deleteEntity
} from '@/api/client/entityApi';
import { getCategories } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ServiceManagement.module.css';

// Lazy load the CreateServiceModal
const CreateServiceModal = lazy(() => import('@/features/CreateServiceModal'));

// Service status badge component
const StatusBadge = ({ status }) => {
  const getStatusStyle = () => {
    switch (status) {
      case 'active':
        return styles.statusActive;
      case 'booked':
        return styles.statusBooked;
      case 'completed':
        return styles.statusCompleted;
      case 'cancelled':
        return styles.statusCancelled;
      case 'archived':
        return styles.statusArchived;
      case 'unavailable':
        return styles.statusUnavailable;
      default:
        return styles.statusInactive;
    }
  };

  const getStatusIcon = () => {
    switch (status) {
      case 'active':
        return <CheckCircle size={14} />;
      case 'booked':
        return <Calendar size={14} />;
      case 'completed':
        return <Star size={14} />;
      case 'cancelled':
        return <XCircle size={14} />;
      case 'archived':
        return <Archive size={14} />;
      case 'unavailable':
        return <AlertCircle size={14} />;
      default:
        return <Clock size={14} />;
    }
  };

  return (
    <span className={`${styles.statusBadge} ${getStatusStyle()}`}>
      {getStatusIcon()}
      {status.replace('_', ' ')}
    </span>
  );
};

// Service row component
const ServiceRow = ({ service, onAction, categories }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, right: 0 });
  const menuButtonRef = React.useRef(null);
  const category = categories.find(c => c.id === service.categoryId);
  
  const getAvailabilityStatus = () => {
    if (!service.availability) return 'not_set';
    if (service.availability === 'available') return 'available';
    if (service.availability === 'busy') return 'busy';
    return 'unavailable';
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: service.currency || 'USD'
    }).format(price);
  };

  return (
    <tr className={styles.serviceRow}>
      <td className={styles.serviceCell}>
        <div className={styles.serviceInfo}>
          <div className={styles.serviceImage}>
            {service.thumbnail ? (
              <img src={service.thumbnail} alt={service.name} />
            ) : (
              <Briefcase size={24} />
            )}
          </div>
          <div className={styles.serviceDetails}>
            <div className={styles.serviceName}>{service.name}</div>
            <div className={styles.serviceMeta}>
              <span className={styles.serviceType}>{service.serviceType || 'General'}</span>
              {service.providerName && <span className={styles.providerName}>{service.providerName}</span>}
              {service.userType && <span className={styles.userType}>{service.userType}</span>}
            </div>
          </div>
        </div>
      </td>
      <td className={styles.categoryCell}>
        <div className={styles.categoryInfo}>
          <span className={styles.categoryTag}>
            <Tag size={12} />
            {category?.name || 'Uncategorized'}
          </span>
          {service.tags && service.tags.length > 0 && (
            <span className={styles.tagsIndicator}>
              +{service.tags.length} tags
            </span>
          )}
        </div>
      </td>
      <td className={styles.priceCell}>
        <div className={styles.priceInfo}>
          <span className={styles.currentPrice}>
            {formatPrice(service.basePrice || service.price)}
            {service.negotiable && <span className={styles.negotiableTag} title="Negotiable">N</span>}
          </span>
          {service.pricing && service.pricing.length > 0 && (
            <span className={styles.pricingOptions}>+{service.pricing.length} options</span>
          )}
        </div>
      </td>
      <td className={styles.availabilityCell}>
        <div className={styles.availabilityInfo}>
          <Calendar size={12} />
          <span className={getAvailabilityStatus() === 'unavailable' ? styles.availabilityDanger : ''}>
            {service.availability || 'Not set'}
          </span>
          {service.hasVariants && (
            <span className={styles.variantIndicator} title="Has variants">
              +V
            </span>
          )}
        </div>
      </td>
      <td className={styles.statusCell}>
        <StatusBadge status={service.status || 'active'} />
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <Calendar size={12} />
          <span>{new Date(service.createdAt).toLocaleDateString()}</span>
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
            aria-label="Service actions"
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
                <button onClick={() => { onAction('view', service); setShowMenu(false); }}>
                  <Eye size={14} /> View
                </button>
                <button onClick={() => { onAction('edit', service); setShowMenu(false); }}>
                  <Edit2 size={14} /> Edit
                </button>
                {service.hasVariants && (
                  <button onClick={() => { onAction('variants', service); setShowMenu(false); }}>
                    <Settings size={14} /> Variants
                  </button>
                )}
                <button onClick={() => { onAction('price', service); setShowMenu(false); }}>
                  <DollarSign size={14} /> Price
                </button>
                <button onClick={() => { onAction('availability', service); setShowMenu(false); }}>
                  <Calendar size={14} /> Availability
                </button>
                {service.status !== 'completed' && service.status !== 'cancelled' && (
                  <>
                    <button onClick={() => { onAction('book', service); setShowMenu(false); }}>
                      <Users size={14} /> Book
                    </button>
                    <button onClick={() => { onAction('complete', service); setShowMenu(false); }}>
                      <Star size={14} /> Complete
                    </button>
                    <button onClick={() => { onAction('cancel', service); setShowMenu(false); }}>
                      <XCircle size={14} /> Cancel
                    </button>
                  </>
                )}
                <button onClick={() => { onAction('archive', service); setShowMenu(false); }}>
                  <Archive size={14} /> {service.status === 'archived' ? 'Unarchive' : 'Archive'}
                </button>
                {service.middlemanService && (
                  <button onClick={() => { onAction('middleman', service); setShowMenu(false); }}>
                    <CheckCircle size={14} /> Middleman
                  </button>
                )}
                <button 
                  onClick={() => { onAction('delete', service); setShowMenu(false); }}
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
const PriceUpdateModal = ({ service, onClose, onSave }) => {
  const [basePrice, setBasePrice] = useState(service.basePrice || service.price);
  const [negotiable, setNegotiable] = useState(service.negotiable || false);

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({ 
      basePrice: parseFloat(basePrice),
      negotiable 
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Update Pricing - {service.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Base Price *</label>
            <input
              type="number"
              value={basePrice}
              onChange={(e) => setBasePrice(e.target.value)}
              min="0"
              step="0.01"
              required
            />
          </div>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>
              <input
                type="checkbox"
                checked={negotiable}
                onChange={(e) => setNegotiable(e.target.checked)}
              />
              Negotiable
            </label>
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              Update Price
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const AvailabilityModal = ({ service, onClose, onSave }) => {
  const [availability, setAvailability] = useState(service.availability || 'available');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({ availability });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Update Availability - {service.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Availability Status</label>
            <select value={availability} onChange={(e) => setAvailability(e.target.value)}>
              <option value="available">Available</option>
              <option value="busy">Busy</option>
              <option value="unavailable">Unavailable</option>
            </select>
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              Update Availability
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Main component
const ServiceManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('ServiceManagement');
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
  const [filterServiceType, setFilterServiceType] = useState('all');
  const [filterUserType, setFilterUserType] = useState('all');
  const [filterNegotiable, setFilterNegotiable] = useState('all');
  const [filterMiddleman, setFilterMiddleman] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);
  const [selectedService, setSelectedService] = useState(null);
  const [modalType, setModalType] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingService, setEditingService] = useState(null);

  const itemsPerPage = 20;

  // Fetch services
  const { data: servicesData, isLoading, error, refetch } = useQuery({
    queryKey: ['adminServices', currentPage, sortBy, sortOrder],
    queryFn: () => getEntities('service', {
      page: currentPage,
      pageSize: itemsPerPage,
      sortBy,
      sortOrder
    }),
    staleTime: 60000,
    onError: (error) => {
      // Error: 'Failed to fetch services:', error...
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
    mutationFn: (serviceId) => deleteEntity('service', serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, ...updateData }) => updateEntity('service', id, updateData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
      setModalType(null);
      setSelectedService(null);
    },
  });

  const bookMutation = useMutation({
    mutationFn: (serviceId) => updateEntity('service', serviceId, { status: 'booked' }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
    },
  });

  const completeMutation = useMutation({
    mutationFn: (serviceId) => updateEntity('service', serviceId, { status: 'completed' }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
    },
  });

  const cancelMutation = useMutation({
    mutationFn: (serviceId) => updateEntity('service', serviceId, { status: 'cancelled' }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
    },
  });

  const archiveMutation = useMutation({
    mutationFn: (serviceId) => updateEntity('service', serviceId, { 
      status: 'archived',
      archived: true 
    }),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminServices']);
    },
  });

  // Filter services
  const filteredServices = useMemo(() => {
    let services = servicesData?.services || servicesData?.data || [];

    if (searchTerm) {
      services = services.filter(service =>
        service.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        service.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        service.serviceType?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        service.providerName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        service.tags?.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
      );
    }

    if (filterCategory !== 'all') {
      services = services.filter(service => service.categoryId === filterCategory);
    }

    if (filterStatus !== 'all') {
      services = services.filter(service => service.status === filterStatus);
    }

    if (filterServiceType !== 'all') {
      services = services.filter(service => service.serviceType === filterServiceType);
    }

    if (filterUserType !== 'all') {
      services = services.filter(service => service.userType === filterUserType);
    }

    if (filterNegotiable !== 'all') {
      services = services.filter(service => service.negotiable === (filterNegotiable === 'yes'));
    }

    if (filterMiddleman !== 'all') {
      services = services.filter(service => service.middlemanService === (filterMiddleman === 'yes'));
    }

    return services;
  }, [servicesData, searchTerm, filterCategory, filterStatus, filterServiceType, filterUserType, filterNegotiable, filterMiddleman]);

  // Handle modal close
  const handleModalClose = useCallback(() => {
    setShowCreateModal(false);
    setEditingService(null);
    queryClient.invalidateQueries(['adminServices']);
  }, [queryClient]);

  // Handle actions
  const handleServiceAction = useCallback((action, service) => {
    switch (action) {
      case 'view':
        router.push(`/admin/services/${service.id}`);
        break;
      case 'edit':
        setEditingService(service);
        setShowCreateModal(true);
        break;
      case 'variants':
        router.push(`/admin/services/${service.id}/variants`);
        break;
      case 'price':
        setSelectedService(service);
        setModalType('price');
        break;
      case 'availability':
        setSelectedService(service);
        setModalType('availability');
        break;
      case 'archive':
        if (confirm(`Are you sure you want to ${service.status === 'archived' ? 'unarchive' : 'archive'} "${service.name}"?`)) {
          archiveMutation.mutate(service.id);
        }
        break;
      case 'book':
        if (confirm(`Mark "${service.name}" as booked?`)) {
          bookMutation.mutate(service.id);
        }
        break;
      case 'complete':
        if (confirm(`Mark "${service.name}" as completed?`)) {
          completeMutation.mutate(service.id);
        }
        break;
      case 'cancel':
        if (confirm(`Cancel service "${service.name}"?`)) {
          cancelMutation.mutate(service.id);
        }
        break;
      case 'middleman':
        router.push(`/admin/middleman/services/${service.id}`);
        break;
      case 'delete':
        if (confirm(`Are you sure you want to delete "${service.name}"? This action cannot be undone.`)) {
          deleteMutation.mutate(service.id);
        }
        break;
    }
  }, [router, archiveMutation, bookMutation, completeMutation, cancelMutation, deleteMutation]);

  // Calculate stats
  const stats = useMemo(() => {
    const services = servicesData?.services || servicesData?.data || [];
    return {
      total: servicesData?.total || services.length || 0,
      active: services.filter(s => s.status === 'active').length,
      booked: services.filter(s => s.status === 'booked').length,
      completed: services.filter(s => s.status === 'completed').length,
      archived: services.filter(s => s.status === 'archived').length,
      totalValue: services.reduce((sum, s) => sum + (s.basePrice || s.price || 0), 0)
    };
  }, [servicesData]);

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
    // Error: 'Service loading error:', error...
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>Failed to load services</h2>
          <p>{error.message || 'An error occurred while fetching services'}</p>
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
              {t('title', { defaultValue: 'Service Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage your services, providers, and availability' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button
              className={styles.addButton}
              onClick={() => setShowCreateModal(true)}
            >
              <Plus size={14} />
              {t('add', { defaultValue: 'Add Service' })}
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
              <Briefcase size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.total}</span>
              <span className={styles.statLabel}>Total Services</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.active}</span>
              <span className={styles.statLabel}>Active</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <Users size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.booked}</span>
              <span className={styles.statLabel}>Booked</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Star size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.completed}</span>
              <span className={styles.statLabel}>Completed</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <DollarSign size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>
                ${stats.totalValue.toLocaleString()}
              </span>
              <span className={styles.statLabel}>Total Value</span>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.searchBox}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search services...' })}
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
                <option value="active">Active</option>
                <option value="booked">Booked</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
                <option value="archived">Archived</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByServiceType', { defaultValue: 'Service Type' })}</label>
              <select
                value={filterServiceType}
                onChange={(e) => setFilterServiceType(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Types</option>
                <option value="consultation">Consultation</option>
                <option value="maintenance">Maintenance</option>
                <option value="design">Design</option>
                <option value="development">Development</option>
                <option value="repair">Repair</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByUserType', { defaultValue: 'Provider Type' })}</label>
              <select
                value={filterUserType}
                onChange={(e) => setFilterUserType(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Types</option>
                <option value="private">Private</option>
                <option value="business">Business</option>
                <option value="freelancer">Freelancer</option>
                <option value="agency">Agency</option>
                <option value="consultant">Consultant</option>
                <option value="company">Company</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('negotiable', { defaultValue: 'Negotiable' })}</label>
              <select
                value={filterNegotiable}
                onChange={(e) => setFilterNegotiable(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All</option>
                <option value="yes">Yes</option>
                <option value="no">No</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('middleman', { defaultValue: 'Middleman' })}</label>
              <select
                value={filterMiddleman}
                onChange={(e) => setFilterMiddleman(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All</option>
                <option value="yes">Available</option>
                <option value="no">Not Available</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('sortBy', { defaultValue: 'Sort By' })}</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="createdAt">Date Added</option>
                <option value="name">Name</option>
                <option value="basePrice">Price</option>
                <option value="serviceType">Service Type</option>
                <option value="status">Status</option>
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

        {/* Services Table */}
        <div className={styles.tableContainer}>
          <table className={styles.servicesTable}>
            <thead>
              <tr>
                <th>Service</th>
                <th>Category</th>
                <th>Price</th>
                <th>Availability</th>
                <th>Status</th>
                <th>Date</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {filteredServices.length === 0 ? (
                <tr>
                  <td colSpan="7" className={styles.emptyState}>
                    <Briefcase size={48} />
                    <h3>No services found</h3>
                    <p>Try adjusting your filters or add a new service</p>
                  </td>
                </tr>
              ) : (
                filteredServices.map(service => (
                  <ServiceRow
                    key={service.id}
                    service={service}
                    onAction={handleServiceAction}
                    categories={categories}
                  />
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {servicesData && servicesData.totalPages > 1 && (
          <div className={styles.pagination}>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
            >
              Previous
            </button>
            <span className={styles.paginationInfo}>
              Page {currentPage} of {servicesData.totalPages}
            </span>
            <button
              className={styles.paginationButton}
              onClick={() => setCurrentPage(prev => Math.min(servicesData.totalPages, prev + 1))}
              disabled={currentPage === servicesData.totalPages}
            >
              Next
            </button>
          </div>
        )}

        {/* Modals */}
        {modalType === 'price' && selectedService && (
          <PriceUpdateModal
            service={selectedService}
            onClose={() => {
              setModalType(null);
              setSelectedService(null);
            }}
            onSave={(priceData) => {
              updateMutation.mutate({ id: selectedService.id, ...priceData });
            }}
          />
        )}

        {modalType === 'availability' && selectedService && (
          <AvailabilityModal
            service={selectedService}
            onClose={() => {
              setModalType(null);
              setSelectedService(null);
            }}
            onSave={(availabilityData) => {
              updateMutation.mutate({ id: selectedService.id, ...availabilityData });
            }}
          />
        )}

        {/* Create/Edit Service Modal */}
        {showCreateModal && (
          <Suspense fallback={<LoadingSpinner />}>
            <CreateServiceModal
              onClose={handleModalClose}
              editMode={!!editingService}
              initialServiceData={editingService}
            />
          </Suspense>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default ServiceManagement;
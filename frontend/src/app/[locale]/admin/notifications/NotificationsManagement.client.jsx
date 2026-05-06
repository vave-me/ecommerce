"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Bell,
  BellRing,
  AlertCircle,
  CheckCircle,
  Info,
  Megaphone,
  Users,
  Eye,
  Send,
  Calendar,
  Clock,
  Target,
  Search,
  Filter,
  Download,
  RefreshCw,
  Edit,
  Trash2,
  MoreVertical,
  User,
  Plus,
  TrendingUp,
  Activity,
  Zap,
  Mail,
  MessageSquare
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listNotifications,
  createNotification,
  updateNotification,
  deleteNotification,
  markNotificationAsRead,
  markAllNotificationsAsRead,
  getNotificationStats,
  sendBulkNotification,
  listNotificationTemplates,
  createNotificationTemplate
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './NotificationsManagement.module.css';

const NotificationTypeIcon = ({ type }) => {
  const icons = {
    system: AlertCircle,
    promotional: Megaphone,
    announcement: Bell,
    alert: BellRing,
    info: Info
  };
  
  const Icon = icons[type] || Bell;
  return <Icon size={16} />;
};

const NotificationStatusBadge = ({ status, deliveryRate }) => {
  const statusConfig = {
    draft: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Draft' },
    scheduled: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Scheduled' },
    sent: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Sent' },
    delivered: { color: '#3b82f6', bg: 'rgba(59, 130, 246, 0.1)', text: 'Delivered' },
    failed: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Failed' }
  };

  const config = statusConfig[status] || statusConfig.draft;

  return (
    <div className={styles.statusContainer}>
      <span 
        className={styles.statusBadge}
        style={{ color: config.color, backgroundColor: config.bg }}
      >
        {config.text}
      </span>
      {deliveryRate !== undefined && (
        <span className={styles.deliveryRate}>
          {deliveryRate}% delivered
        </span>
      )}
    </div>
  );
};

const NotificationRow = ({ notification, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);

  return (
    <tr className={styles.notificationRow}>
      <td className={styles.notificationCell}>
        <div className={styles.notificationInfo}>
          <div className={styles.notificationHeader}>
            <NotificationTypeIcon type={notification.type} />
            <span className={styles.notificationTitle}>{notification.title}</span>
          </div>
          <div className={styles.notificationPreview}>
            {notification.message?.substring(0, 80)}...
          </div>
        </div>
      </td>
      <td className={styles.typeCell}>
        <span className={styles.notificationType}>
          {notification.type?.toUpperCase()}
        </span>
      </td>
      <td className={styles.audienceCell}>
        <div className={styles.audience}>
          <Users size={14} />
          <span>{notification.recipientCount?.toLocaleString() || 0}</span>
        </div>
      </td>
      <td className={styles.statusCell}>
        <NotificationStatusBadge 
          status={notification.status} 
          deliveryRate={notification.deliveryRate}
        />
      </td>
      <td className={styles.metricsCell}>
        <div className={styles.metrics}>
          <div className={styles.metric}>
            <Send size={12} />
            <span>{notification.sentCount || 0}</span>
          </div>
          <div className={styles.metric}>
            <Eye size={12} />
            <span>{notification.openedCount || 0}</span>
          </div>
        </div>
      </td>
      <td className={styles.dateCell}>
        <div className={styles.dateInfo}>
          <Calendar size={14} />
          <span>{new Date(notification.createdAt).toLocaleDateString()}</span>
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
              <button onClick={() => onAction('view', notification)}>
                <Eye size={14} />
                View Details
              </button>
              {notification.status === 'draft' && (
                <>
                  <button onClick={() => onAction('edit', notification)}>
                    <Edit size={14} />
                    Edit
                  </button>
                  <button onClick={() => onAction('send', notification)}>
                    <Send size={14} />
                    Send Now
                  </button>
                </>
              )}
              <button onClick={() => onAction('duplicate', notification)}>
                <Plus size={14} />
                Duplicate
              </button>
              <button onClick={() => onAction('delete', notification)} className={styles.dangerAction}>
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

const CreateNotificationModal = ({ isOpen, onClose, onSave }) => {
  const t = useTranslations('NotificationsManagement');
  const [formData, setFormData] = useState({
    title: '',
    message: '',
    type: 'announcement',
    audience: 'all',
    scheduleDate: '',
    priority: 'normal'
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(formData);
  };

  if (!isOpen) return null;

  return (
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3>Create Notification</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label>Title</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({...formData, title: e.target.value})}
              required
              className={styles.formInput}
            />
          </div>
          
          <div className={styles.formGroup}>
            <label>Message</label>
            <textarea
              value={formData.message}
              onChange={(e) => setFormData({...formData, message: e.target.value})}
              required
              rows={4}
              className={styles.formTextarea}
            />
          </div>

          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label>Type</label>
              <select
                value={formData.type}
                onChange={(e) => setFormData({...formData, type: e.target.value})}
                className={styles.formSelect}
              >
                <option value="announcement">Announcement</option>
                <option value="promotional">Promotional</option>
                <option value="system">System</option>
                <option value="alert">Alert</option>
              </select>
            </div>

            <div className={styles.formGroup}>
              <label>Audience</label>
              <select
                value={formData.audience}
                onChange={(e) => setFormData({...formData, audience: e.target.value})}
                className={styles.formSelect}
              >
                <option value="all">All Users</option>
                <option value="customers">Customers Only</option>
                <option value="business">Business Users</option>
                <option value="admins">Administrators</option>
              </select>
            </div>
          </div>

          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label>Schedule (Optional)</label>
              <input
                type="datetime-local"
                value={formData.scheduleDate}
                onChange={(e) => setFormData({...formData, scheduleDate: e.target.value})}
                className={styles.formInput}
              />
            </div>

            <div className={styles.formGroup}>
              <label>Priority</label>
              <select
                value={formData.priority}
                onChange={(e) => setFormData({...formData, priority: e.target.value})}
                className={styles.formSelect}
              >
                <option value="low">Low</option>
                <option value="normal">Normal</option>
                <option value="high">High</option>
                <option value="urgent">Urgent</option>
              </select>
            </div>
          </div>

          <div className={styles.modalActions}>
            <button type="button" className={styles.cancelButton} onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className={styles.submitButton}>
              Create Notification
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const NotificationsManagement = () => {
  const t = useTranslations('NotificationsManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    type: 'all',
    status: 'all',
    dateRange: '30d',
    search: ''
  });
  const [showCreateModal, setShowCreateModal] = useState(false);

  // Fetch notifications data
  const { 
    data: notificationsData, 
    isLoading: notificationsLoading, 
    error: notificationsError,
    refetch: refetchNotifications 
  } = useQuery({
    queryKey: ['adminNotifications', filters],
    queryFn: () => listNotifications(filters),
    enabled: isAdmin
  });

  // Fetch notification statistics
  const { 
    data: statsData, 
    isLoading: statsLoading 
  } = useQuery({
    queryKey: ['notificationStats'],
    queryFn: () => getNotificationStats(),
    enabled: isAdmin
  });

  // Create notification mutation
  const createMutation = useMutation({
    mutationFn: createNotification,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminNotifications']);
      queryClient.invalidateQueries(['notificationStats']);
      setShowCreateModal(false);
    },
    onError: (error) => {
      
      alert('Failed to create notification. Please try again.');
    }
  });

  // Delete notification mutation
  const deleteMutation = useMutation({
    mutationFn: deleteNotification,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminNotifications']);
      queryClient.invalidateQueries(['notificationStats']);
    },
    onError: (error) => {
      
      alert('Failed to delete notification. Please try again.');
    }
  });

  const handleNotificationAction = useCallback((action, notification) => {
    switch (action) {
      case 'view':
        router.push(`/admin/notifications/${notification.id}`);
        break;
      case 'edit':
        router.push(`/admin/notifications/${notification.id}/edit`);
        break;
      case 'send':
        if (confirm('Send this notification now?')) {
          // Implement send functionality
          
        }
        break;
      case 'duplicate':
        // Implement duplicate functionality
        
        break;
      case 'delete':
        if (confirm('Are you sure you want to delete this notification?')) {
          deleteMutation.mutate(notification.id);
        }
        break;
    }
  }, [router, deleteMutation]);

  const handleCreateNotification = useCallback((formData) => {
    createMutation.mutate(formData);
  }, [createMutation]);

  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleExport = useCallback(() => {
    const csvData = notificationsData?.notifications.map(n => ({
      'Notification ID': n.id,
      'Title': n.title,
      'Type': n.type,
      'Status': n.status,
      'Recipients': n.recipientCount || 0,
      'Sent': n.sentCount || 0,
      'Opened': n.openedCount || 0,
      'Created': new Date(n.createdAt).toLocaleDateString()
    }));

  }, [notificationsData]);

  // Process data
  const notifications = notificationsData?.notifications || [];
  const stats = statsData || {};

  // Calculate summary stats
  const summaryStats = useMemo(() => {
    const today = new Date();
    const todaysNotifications = notifications.filter(n => 
      new Date(n.createdAt).toDateString() === today.toDateString()
    );
    
    const totalSent = notifications.reduce((sum, n) => sum + (n.sentCount || 0), 0);
    const totalOpened = notifications.reduce((sum, n) => sum + (n.openedCount || 0), 0);
    const totalDelivered = notifications.reduce((sum, n) => sum + (n.deliveredCount || 0), 0);

    return {
      totalNotifications: notifications.length,
      sentToday: todaysNotifications.reduce((sum, n) => sum + (n.sentCount || 0), 0),
      deliveryRate: totalSent > 0 ? (totalDelivered / totalSent) * 100 : 0,
      openRate: totalSent > 0 ? (totalOpened / totalSent) * 100 : 0
    };
  }, [notifications]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access notification management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (notificationsLoading && !notificationsData) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Notifications...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch notification data.' })}</p>
        </div>
      </div>
    );
  }

  if (notificationsError) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>{t('errorTitle', { defaultValue: 'Failed to Load Notifications' })}</h2>
          <p>{notificationsError.message || t('errorMessage', { defaultValue: 'An error occurred while fetching notification data' })}</p>
          <button className={styles.retryButton} onClick={() => refetchNotifications()}>
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
            <h1 className={styles.title}>{t('title', { defaultValue: 'Notifications Management' })}</h1>
            <p className={styles.subtitle}>{t('subtitle', { defaultValue: 'Manage system notifications and alerts' })}</p>
          </div>
          <div className={styles.headerActions}>
            <button 
              className={styles.createButton} 
              onClick={() => setShowCreateModal(true)}
            >
              <Plus size={16} />
              {t('createNotification', { defaultValue: 'Create Notification' })}
            </button>
            <button className={styles.exportButton} onClick={handleExport}>
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button className={styles.refreshButton} onClick={() => refetchNotifications()}>
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Bell size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.totalNotifications.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalNotifications', { defaultValue: 'Total Notifications' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Send size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.sentToday.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('sentToday', { defaultValue: 'Sent Today' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.deliveryRate.toFixed(1)}%</div>
              <div className={styles.statLabel}>{t('deliveryRate', { defaultValue: 'Delivery Rate' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Eye size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.openRate.toFixed(1)}%</div>
              <div className={styles.statLabel}>{t('openRate', { defaultValue: 'Open Rate' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersSection}>
          <div className={styles.searchContainer}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search notifications...' })}
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filterControls}>
            <select
              value={filters.type}
              onChange={(e) => handleFilterChange('type', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allTypes', { defaultValue: 'All Types' })}</option>
              <option value="system">{t('system', { defaultValue: 'System' })}</option>
              <option value="promotional">{t('promotional', { defaultValue: 'Promotional' })}</option>
              <option value="announcement">{t('announcement', { defaultValue: 'Announcement' })}</option>
              <option value="alert">{t('alert', { defaultValue: 'Alert' })}</option>
            </select>
            <select
              value={filters.status}
              onChange={(e) => handleFilterChange('status', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
              <option value="draft">{t('draft', { defaultValue: 'Draft' })}</option>
              <option value="scheduled">{t('scheduled', { defaultValue: 'Scheduled' })}</option>
              <option value="sent">{t('sent', { defaultValue: 'Sent' })}</option>
              <option value="delivered">{t('delivered', { defaultValue: 'Delivered' })}</option>
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

        {/* Notifications Table */}
        <div className={styles.tableSection}>
          <div className={styles.tableContainer}>
            <table className={styles.notificationsTable}>
              <thead>
                <tr>
                  <th>{t('title', { defaultValue: 'Title' })}</th>
                  <th>{t('type', { defaultValue: 'Type' })}</th>
                  <th>{t('recipients', { defaultValue: 'Recipients' })}</th>
                  <th>Status</th>
                  <th>Performance</th>
                  <th>{t('date', { defaultValue: 'Date' })}</th>
                  <th>{t('actions', { defaultValue: 'Actions' })}</th>
                </tr>
              </thead>
              <tbody>
                {notifications.map((notification) => (
                  <NotificationRow
                    key={notification.id}
                    notification={notification}
                    onAction={handleNotificationAction}
                  />
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Create Notification Modal */}
        <CreateNotificationModal
          isOpen={showCreateModal}
          onClose={() => setShowCreateModal(false)}
          onSave={handleCreateNotification}
        />
      </div>
    </ErrorBoundary>
  );
};

export default NotificationsManagement; 
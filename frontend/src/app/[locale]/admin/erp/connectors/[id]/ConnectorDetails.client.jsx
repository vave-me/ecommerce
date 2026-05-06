"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  ArrowLeft,
  RefreshCw,
  Settings,
  Activity,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Clock,
  Database,
  Link2,
  Power,
  Edit,
  Trash2,
  Package,
  Users,
  DollarSign,
  BarChart3,
  Calendar,
  Shield,
  Zap
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import {
  getConnectorStatus,
  getSyncHistory,
  removeConnector,
  toggleConnector,
  syncProducts,
  syncStock,
  syncPrices,
  syncCustomers,
  formatDate
} from '@/api/client/admin/erpApi';
import styles from './ConnectorDetails.module.css';

const StatusBadge = ({ status }) => {
  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'active':
      case 'healthy': return 'success';
      case 'inactive': return 'neutral';
      case 'syncing':
      case 'synchronizing': return 'info';
      case 'error':
      case 'unhealthy': return 'error';
      case 'maintenance': return 'warning';
      default: return 'neutral';
    }
  };

  const getStatusIcon = (status) => {
    switch (status?.toLowerCase()) {
      case 'active':
      case 'healthy': return CheckCircle;
      case 'inactive': return XCircle;
      case 'syncing':
      case 'synchronizing': return RefreshCw;
      case 'error':
      case 'unhealthy': return AlertTriangle;
      case 'maintenance': return Clock;
      default: return Database;
    }
  };

  const Icon = getStatusIcon(status);
  const color = getStatusColor(status);

  return (
    <div className={`${styles.statusBadge} ${styles[color]}`}>
      <Icon size={16} />
      <span>{(status || 'Unknown').toUpperCase()}</span>
    </div>
  );
};

const SyncCard = ({ title, icon: Icon, lastSync, status, onSync, loading }) => {
  return (
    <div className={styles.syncCard}>
      <div className={styles.syncCardHeader}>
        <Icon size={20} />
        <h4>{title}</h4>
      </div>
      <div className={styles.syncCardBody}>
        <div className={styles.syncInfo}>
          <span className={styles.syncLabel}>Last Sync:</span>
          <span className={styles.syncValue}>
            {lastSync ? formatDate(lastSync) : 'Never'}
          </span>
        </div>
        <div className={styles.syncInfo}>
          <span className={styles.syncLabel}>Status:</span>
          <StatusBadge status={status || 'pending'} />
        </div>
      </div>
      <button 
        className={styles.syncButton}
        onClick={onSync}
        disabled={loading}
      >
        <RefreshCw size={14} className={loading ? styles.spinning : ''} />
        Sync Now
      </button>
    </div>
  );
};

const ConnectorDetails = ({ connectorId }) => {
  const t = useTranslations('ERPConnectors');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [activeTab, setActiveTab] = useState('overview');

  // Fetch connector status
  const { 
    data: statusData, 
    isLoading: statusLoading, 
    error: statusError,
    refetch: refetchStatus 
  } = useQuery({
    queryKey: ['erp-connector-status', connectorId],
    queryFn: async () => {
      const response = await getConnectorStatus(connectorId);
      
      return response;
    },
    staleTime: 30000,
    refetchInterval: 60000, // Refresh every minute
  });

  // Fetch sync history
  const { 
    data: syncHistoryData, 
    isLoading: syncHistoryLoading 
  } = useQuery({
    queryKey: ['erp-sync-history', connectorId],
    queryFn: async () => {
      const response = await getSyncHistory(connectorId, { pageSize: 20 });
      
      return response;
    },
    staleTime: 60000,
  });

  // Delete connector mutation
  const deleteMutation = useMutation({
    mutationFn: async () => {
      return await removeConnector(connectorId, true);
    },
    onSuccess: () => {
      
      router.push('/admin/erp/connectors');
    },
    onError: (error) => {
      
      alert('Failed to delete connector: ' + error.message);
    }
  });

  // Toggle connector mutation
  const toggleMutation = useMutation({
    mutationFn: async () => {
      const isActive = statusData?.connector?.status === 'active';
      return await toggleConnector(connectorId, {
        activate: !isActive,
        reason: isActive ? 'Manual deactivation' : 'Manual activation'
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-connector-status', connectorId]);
      
    },
    onError: (error) => {
      
      alert('Failed to toggle connector: ' + error.message);
    }
  });

  // Sync mutations
  const syncProductsMutation = useMutation({
    mutationFn: () => syncProducts(connectorId),
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-sync-history', connectorId]);
      
    }
  });

  const syncStockMutation = useMutation({
    mutationFn: () => syncStock(connectorId),
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-sync-history', connectorId]);
      
    }
  });

  const syncPricesMutation = useMutation({
    mutationFn: () => syncPrices(connectorId),
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-sync-history', connectorId]);
      
    }
  });

  const syncCustomersMutation = useMutation({
    mutationFn: () => syncCustomers(connectorId),
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-sync-history', connectorId]);
      
    }
  });

  const handleDelete = () => {
    if (window.confirm(`Are you sure you want to delete this connector? This action cannot be undone.`)) {
      deleteMutation.mutate();
    }
  };

  const handleEdit = () => {
    router.push(`/admin/erp/connectors/${connectorId}/edit`);
  };

  if (!isAdmin) {
    return (
      <div className={styles.accessDenied}>
        <AlertTriangle size={48} />
        <h2>Access Denied</h2>
        <p>You don't have permission to view connector details.</p>
      </div>
    );
  }

  if (statusLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <p>Loading connector details...</p>
      </div>
    );
  }

  if (statusError) {
    return (
      <div className={styles.errorContainer}>
        <AlertTriangle size={24} />
        <p>Error loading connector: {statusError.message}</p>
        <button onClick={() => refetchStatus()}>Try Again</button>
      </div>
    );
  }

  const connector = statusData?.connector;
  const syncSummary = statusData?.lastSync;
  const webhookInfo = statusData?.webhookInfo;

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <button 
          className={styles.backButton}
          onClick={() => router.push('/admin/erp/connectors')}
        >
          <ArrowLeft size={20} />
          Back to Connectors
        </button>
        
        <div className={styles.headerActions}>
          <button 
            className={styles.actionButton}
            onClick={() => toggleMutation.mutate()}
            disabled={toggleMutation.isPending}
          >
            <Power size={16} />
            {connector?.status === 'active' ? 'Deactivate' : 'Activate'}
          </button>
          
          <button 
            className={styles.actionButton}
            onClick={handleEdit}
          >
            <Edit size={16} />
            Edit
          </button>
          
          <button 
            className={`${styles.actionButton} ${styles.danger}`}
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
          >
            <Trash2 size={16} />
            Delete
          </button>
        </div>
      </div>

      {/* Connector Info */}
      <div className={styles.connectorInfo}>
        <div className={styles.connectorHeader}>
          <div>
            <h1 className={styles.connectorName}>{connector?.name || 'Unknown Connector'}</h1>
            <p className={styles.connectorMeta}>
              <span>{connector?.type?.toUpperCase()}</span>
              <span className={styles.separator}>•</span>
              <span>{connector?.environment || 'production'}</span>
              <span className={styles.separator}>•</span>
              <span>ID: {connectorId}</span>
            </p>
          </div>
          <StatusBadge status={statusData?.healthStatus} />
        </div>
      </div>

      {/* Tabs */}
      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${activeTab === 'overview' ? styles.active : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          <Database size={16} />
          Overview
        </button>
        <button
          className={`${styles.tab} ${activeTab === 'sync' ? styles.active : ''}`}
          onClick={() => setActiveTab('sync')}
        >
          <RefreshCw size={16} />
          Synchronization
        </button>
        <button
          className={`${styles.tab} ${activeTab === 'history' ? styles.active : ''}`}
          onClick={() => setActiveTab('history')}
        >
          <Clock size={16} />
          History
        </button>
        <button
          className={`${styles.tab} ${activeTab === 'settings' ? styles.active : ''}`}
          onClick={() => setActiveTab('settings')}
        >
          <Settings size={16} />
          Settings
        </button>
      </div>

      {/* Tab Content */}
      <div className={styles.tabContent}>
        {activeTab === 'overview' && (
          <div className={styles.overviewTab}>
            {/* Connection Details */}
            <div className={styles.section}>
              <h2>Connection Details</h2>
              <div className={styles.detailsGrid}>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Base URL</span>
                  <span className={styles.detailValue}>{connector?.baseUrl || 'Not configured'}</span>
                </div>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Created</span>
                  <span className={styles.detailValue}>
                    {connector?.createdAt ? formatDate(connector.createdAt) : 'Unknown'}
                  </span>
                </div>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Last Updated</span>
                  <span className={styles.detailValue}>
                    {connector?.updatedAt ? formatDate(connector.updatedAt) : 'Unknown'}
                  </span>
                </div>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Health Check</span>
                  <span className={styles.detailValue}>
                    {connector?.lastHealthCheckAt ? formatDate(connector.lastHealthCheckAt) : 'Never'}
                  </span>
                </div>
              </div>
            </div>

            {/* Sync Summary */}
            {syncSummary && (
              <div className={styles.section}>
                <h2>Sync Summary</h2>
                <div className={styles.statsGrid}>
                  <div className={styles.statCard}>
                    <span className={styles.statValue}>{syncSummary.totalSyncs || 0}</span>
                    <span className={styles.statLabel}>Total Syncs</span>
                  </div>
                  <div className={styles.statCard}>
                    <span className={styles.statValue}>{syncSummary.failedSyncs || 0}</span>
                    <span className={styles.statLabel}>Failed Syncs</span>
                  </div>
                  <div className={styles.statCard}>
                    <span className={styles.statValue}>
                      {syncSummary.totalSyncs > 0 
                        ? `${Math.round(((syncSummary.totalSyncs - syncSummary.failedSyncs) / syncSummary.totalSyncs) * 100)}%`
                        : '0%'
                      }
                    </span>
                    <span className={styles.statLabel}>Success Rate</span>
                  </div>
                </div>
              </div>
            )}

            {/* Webhook Info */}
            {webhookInfo && (
              <div className={styles.section}>
                <h2>Webhook Configuration</h2>
                <div className={styles.webhookInfo}>
                  <div className={styles.webhookStatus}>
                    <Zap size={20} />
                    <span>{webhookInfo.enabled ? 'Webhooks Enabled' : 'Webhooks Disabled'}</span>
                  </div>
                  {webhookInfo.enabled && webhookInfo.url && (
                    <div className={styles.webhookUrl}>
                      <Link2 size={16} />
                      <code>{webhookInfo.url}</code>
                    </div>
                  )}
                  {webhookInfo.lastWebhookAt && (
                    <div className={styles.webhookStats}>
                      <span>Last webhook: {formatDate(webhookInfo.lastWebhookAt)}</span>
                      <span>Total: {webhookInfo.totalWebhooks || 0}</span>
                      <span>Failed: {webhookInfo.failedWebhooks || 0}</span>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'sync' && (
          <div className={styles.syncTab}>
            <div className={styles.syncActions}>
              <SyncCard
                title="Products"
                icon={Package}
                lastSync={syncSummary?.lastProductSync}
                status={syncProductsMutation.isPending ? 'syncing' : 'active'}
                onSync={() => syncProductsMutation.mutate()}
                loading={syncProductsMutation.isPending}
              />
              <SyncCard
                title="Stock Levels"
                icon={BarChart3}
                lastSync={syncSummary?.lastStockSync}
                status={syncStockMutation.isPending ? 'syncing' : 'active'}
                onSync={() => syncStockMutation.mutate()}
                loading={syncStockMutation.isPending}
              />
              <SyncCard
                title="Prices"
                icon={DollarSign}
                lastSync={syncSummary?.lastPriceSync}
                status={syncPricesMutation.isPending ? 'syncing' : 'active'}
                onSync={() => syncPricesMutation.mutate()}
                loading={syncPricesMutation.isPending}
              />
              <SyncCard
                title="Customers"
                icon={Users}
                lastSync={syncSummary?.lastOrderSync}
                status={syncCustomersMutation.isPending ? 'syncing' : 'active'}
                onSync={() => syncCustomersMutation.mutate()}
                loading={syncCustomersMutation.isPending}
              />
            </div>
          </div>
        )}

        {activeTab === 'history' && (
          <div className={styles.historyTab}>
            {syncHistoryLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p>Loading sync history...</p>
              </div>
            ) : !syncHistoryData?.syncOperations?.length ? (
              <div className={styles.emptyState}>
                <Clock size={48} />
                <h3>No Sync History</h3>
                <p>No synchronization operations have been performed yet.</p>
              </div>
            ) : (
              <div className={styles.historyList}>
                {syncHistoryData.syncOperations.map((operation, index) => (
                  <div key={operation.id || index} className={styles.historyItem}>
                    <div className={styles.historyIcon}>
                      {operation.status === 'completed' ? (
                        <CheckCircle size={20} className={styles.success} />
                      ) : operation.status === 'failed' ? (
                        <XCircle size={20} className={styles.error} />
                      ) : (
                        <Clock size={20} className={styles.pending} />
                      )}
                    </div>
                    <div className={styles.historyContent}>
                      <div className={styles.historyHeader}>
                        <span className={styles.historyOperation}>
                          {operation.entityType || operation.operation || 'Unknown Operation'}
                        </span>
                        <span className={styles.historyTime}>
                          {operation.completedAt ? formatDate(operation.completedAt) : 'In Progress'}
                        </span>
                      </div>
                      {operation.error && (
                        <div className={styles.historyError}>{operation.error}</div>
                      )}
                      <div className={styles.historyMeta}>
                        <span>Records: {operation.recordsProcessed || 0}</span>
                        {operation.duration && (
                          <span>Duration: {operation.duration}</span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {activeTab === 'settings' && (
          <div className={styles.settingsTab}>
            <div className={styles.settingsInfo}>
              <Shield size={48} />
              <h3>Connector Settings</h3>
              <p>
                To modify connector settings such as authentication, sync intervals, or rate limits,
                please use the edit function.
              </p>
              <button 
                className={styles.editSettingsButton}
                onClick={handleEdit}
              >
                <Edit size={16} />
                Edit Connector Settings
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConnectorDetails;
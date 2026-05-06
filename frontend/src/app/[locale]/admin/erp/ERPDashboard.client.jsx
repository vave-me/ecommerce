"use client";

import React, { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { 
  Database, 
  RefreshCw, 
  Settings, 
  Activity,
  Package,
  FileText,
  RotateCcw,
  Zap,
  Users,
  DollarSign,
  Truck,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Clock,
  BarChart3,
  ArrowRight,
  Plus,
  Link2
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
// REAL ERP API IMPORTS - NO MORE MOCK DATA
import {
  listConnectors,
  getConnectorStatus,
  getSyncHistory,
  syncCustomers,
  syncProducts,
  syncPrices,
  syncStock,
  sendOrder,
  createInventoryReservation,
  createInvoice,
  createReturn,
  getConnectorTypes,
  getConnectorStatuses,
  getInvoiceStatuses,
  getReturnStatuses
} from '@/api/client/admin/erpApi';
import styles from './ERPDashboard.module.css';

// Metric Card Component
const MetricCard = ({ title, value, icon: Icon, trend, onClick, loading }) => (
  <div 
    className={styles.metricCard}
    onClick={onClick}
    role={onClick ? 'button' : undefined}
    tabIndex={onClick ? 0 : undefined}
    onKeyDown={onClick ? (e) => e.key === 'Enter' && onClick() : undefined}
  >
    <Icon className={styles.metricIcon} />
    <div className={styles.metricValue}>
      {loading ? '...' : value}
    </div>
    <div className={styles.metricLabel}>{title}</div>
    {trend && (
      <div className={`${styles.metricTrend} ${styles[trend.direction]}`}>
        {trend.value}
      </div>
    )}
  </div>
);

// Module Card Component
const ModuleCard = ({ title, description, icon: Icon, href, status = 'active' }) => {
  const router = useRouter();
  
  const handleClick = () => {
    router.push(href);
  };

  return (
    <div 
      className={styles.moduleCard}
      onClick={handleClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === 'Enter' && handleClick()}
    >
      <div className={styles.moduleHeader}>
        <Icon className={styles.moduleIcon} />
        <div className={styles.moduleInfo}>
          <h3 className={styles.moduleTitle}>{title}</h3>
          <p className={styles.moduleDescription}>{description}</p>
        </div>
      </div>
      <div className={styles.moduleActions}>
        <button className={styles.moduleButton}>
          MANAGE
          <ArrowRight size={14} />
        </button>
        <span className={`${styles.moduleStatus} ${styles[`status${status.charAt(0).toUpperCase() + status.slice(1)}`]}`}>
          {status.toUpperCase()}
        </span>
      </div>
    </div>
  );
};

// Connector Item Component
const ConnectorItem = ({ connector, onClick, onSync, onViewDetails }) => {
  const getStatusClass = (status) => {
    switch (status?.toLowerCase()) {
      case 'connected':
      case 'active': return 'active';
      case 'syncing':
      case 'synchronizing': return 'syncing';
      case 'error':
      case 'failed':
      case 'disconnected': return 'error';
      case 'pending':
      case 'initializing': return 'pending';
      default: return 'active';
    }
  };

  const formatLastSync = (lastSyncTime) => {
    if (!lastSyncTime) return 'Never';
    try {
      return new Date(lastSyncTime).toLocaleString();
    } catch {
      return lastSyncTime;
    }
  };

  return (
    <div 
      className={styles.connectorItem}
      onClick={() => onClick?.(connector)}
    >
      <div className={styles.connectorIcon}>
        <Database size={20} />
      </div>
      <div className={styles.connectorInfo}>
        <div className={styles.connectorName}>{connector.name || connector.id}</div>
        <div className={styles.connectorType}>{connector.type || 'Unknown'}</div>
      </div>
      <div className={styles.connectorMeta}>
        <div className={styles.connectorStatus}>
          <div className={`${styles.statusIndicator} ${styles[getStatusClass(connector.status)]}`}></div>
          {(connector.status || 'Unknown').toUpperCase()}
        </div>
        <div>Last sync: {formatLastSync(connector.lastSyncTime)}</div>
      </div>
      <div className={styles.connectorActions}>
        <button 
          className={styles.syncButton}
          onClick={(e) => {
            e.stopPropagation();
            onSync?.(connector);
          }}
          disabled={connector.status === 'syncing'}
        >
          <RefreshCw size={14} />
          Sync
        </button>
        <button 
          className={styles.detailsButton}
          onClick={(e) => {
            e.stopPropagation();
            onViewDetails?.(connector);
          }}
        >
          Details
        </button>
      </div>
    </div>
  );
};

const ERPDashboard = () => {
  const t = useTranslations('ERP');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  // REAL API CALL FOR CONNECTORS - NO MORE MOCK DATA
  const { 
    data: connectorsData, 
    isLoading: connectorsLoading, 
    error: connectorsError,
    refetch: refetchConnectors 
  } = useQuery({
    queryKey: ['erp-connectors'],
    queryFn: async () => {
      try {
        const response = await listConnectors();
        
        return response;
      } catch (error) {
        // Error: '❌ Error loading ERP connectors:', error...
        throw error;
      }
    },
    staleTime: 30000,
    retry: 3,
  });

  // Get first connector for sync history (if any)
  const firstConnectorId = connectorsData?.connectors?.[0]?.id;
  
  // REAL API CALL FOR SYNC HISTORY
  const { 
    data: syncHistoryData, 
    isLoading: syncHistoryLoading 
  } = useQuery({
    queryKey: ['erp-sync-history', firstConnectorId],
    queryFn: async () => {
      if (!firstConnectorId) return { syncOperations: [] };
      try {
        const response = await getSyncHistory(firstConnectorId, { pageSize: 10 });
        
        return response;
      } catch (error) {
        // Error: '❌ Error loading sync history:', error...
        return { syncOperations: [] };
      }
    },
    enabled: !!firstConnectorId,
    staleTime: 60000,
  });

  // REAL API CALL FOR CONNECTOR TYPES
  const { 
    data: connectorTypesData 
  } = useQuery({
    queryKey: ['erp-connector-types'],
    queryFn: async () => {
      try {
        const response = await getConnectorTypes();
        
        return response;
      } catch (error) {
        // Error: '❌ Error loading connector types:', error...
        throw error;
      }
    },
    staleTime: 300000, // 5 minutes
  });

  // Calculate stats from REAL data
  const stats = useMemo(() => {
    if (!connectorsData?.connectors) {
      return {
        totalConnectors: 0,
        activeConnectors: 0,
        syncingConnectors: 0,
        errorConnectors: 0,
        healthRate: 0
      };
    }
    
    const connectors = connectorsData.connectors;
    const active = connectors.filter(c => 
      c.status?.toLowerCase() === 'connected' || 
      c.status?.toLowerCase() === 'active'
    ).length;
    const syncing = connectors.filter(c => 
      c.status?.toLowerCase() === 'syncing' || 
      c.status?.toLowerCase() === 'synchronizing'
    ).length;
    const errors = connectors.filter(c => 
      c.status?.toLowerCase() === 'error' || 
      c.status?.toLowerCase() === 'failed' || 
      c.status?.toLowerCase() === 'disconnected'
    ).length;
    
    return {
      totalConnectors: connectors.length,
      activeConnectors: active,
      syncingConnectors: syncing,
      errorConnectors: errors,
      healthRate: connectors.length > 0 ? Math.round((active / connectors.length) * 100) : 0
    };
  }, [connectorsData]);

  // REAL SYNC FUNCTIONS - sync all connectors
  const handleSyncCustomers = async () => {
    if (!connectorsData?.connectors?.length) {
      
      return;
    }
    
    try {
      const promises = connectorsData.connectors
        .filter(c => c.status === 'active')
        .map(connector => syncCustomers(connector.id));
      
      await Promise.allSettled(promises);
      
      refetchConnectors();
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  };

  const handleSyncProducts = async () => {
    if (!connectorsData?.connectors?.length) {
      
      return;
    }
    
    try {
      const promises = connectorsData.connectors
        .filter(c => c.status === 'active')
        .map(connector => syncProducts(connector.id));
      
      await Promise.allSettled(promises);
      
      refetchConnectors();
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  };

  const handleSyncPrices = async () => {
    if (!connectorsData?.connectors?.length) {
      
      return;
    }
    
    try {
      const promises = connectorsData.connectors
        .filter(c => c.status === 'active')
        .map(connector => syncPrices(connector.id));
      
      await Promise.allSettled(promises);
      
      refetchConnectors();
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  };

  const handleSyncStock = async () => {
    if (!connectorsData?.connectors?.length) {
      
      return;
    }
    
    try {
      const promises = connectorsData.connectors
        .filter(c => c.status === 'active')
        .map(connector => syncStock(connector.id));
      
      await Promise.allSettled(promises);
      
      refetchConnectors();
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  };

  const handleConnectorSync = async (connector) => {
    try {
      // Start multiple sync operations for the specific connector
      await Promise.allSettled([
        syncCustomers(connector.id),
        syncProducts(connector.id),
        syncPrices(connector.id),
        syncStock(connector.id)
      ]);
      
      refetchConnectors();
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  };

  const handleViewConnectorDetails = async (connector) => {
    try {
      const status = await getConnectorStatus(connector.id);
      
      // Navigate to detailed view or show modal
      router.push(`/admin/erp/connectors/${connector.id}`);
    } catch (error) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', error);
        }
        // Continue with default behavior
    }
  };

  if (!isAdmin) {
    return (
      <div className={styles.accessDenied}>
        <AlertTriangle size={48} />
        <h2>Access Denied</h2>
        <p>You don't have permission to access the ERP dashboard.</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerText}>
            <h1 className={styles.title}>ERP Integration Dashboard</h1>
            <p className={styles.subtitle}>
              Manage enterprise resource planning integrations and data synchronization
            </p>
          </div>
          <div className={styles.headerActions}>
            <button 
              className={styles.refreshButton}
              onClick={() => refetchConnectors()}
              disabled={connectorsLoading}
            >
              <RefreshCw size={16} className={connectorsLoading ? styles.spinning : ''} />
              Refresh
            </button>
            <button 
              className={styles.addConnectorButton}
              onClick={() => router.push('/admin/erp/connectors/new')}
            >
              <Plus size={16} />
              Add Connector
            </button>
          </div>
        </div>
      </div>

      {/* Metrics */}
      <div className={styles.metricsGrid}>
        <MetricCard
          title="Total Connectors"
          value={stats.totalConnectors}
          icon={Database}
          loading={connectorsLoading}
          onClick={() => router.push('/admin/erp/connectors')}
        />
        <MetricCard
          title="Active Connections"
          value={stats.activeConnectors}
          icon={CheckCircle}
          loading={connectorsLoading}
          trend={{ value: `${stats.healthRate}%`, direction: stats.healthRate > 80 ? 'positive' : 'neutral' }}
        />
        <MetricCard
          title="Syncing Now"
          value={stats.syncingConnectors}
          icon={RefreshCw}
          loading={connectorsLoading}
        />
        <MetricCard
          title="Error Status"
          value={stats.errorConnectors}
          icon={AlertTriangle}
          loading={connectorsLoading}
          trend={{ value: stats.errorConnectors > 0 ? 'Issues' : 'Healthy', direction: stats.errorConnectors > 0 ? 'negative' : 'positive' }}
        />
      </div>

      {/* Quick Actions */}
      <div className={styles.quickActionsSection}>
        <h2>Quick Sync Actions</h2>
        <div className={styles.quickActionsGrid}>
          <button className={styles.quickActionButton} onClick={handleSyncCustomers}>
            <Users size={20} />
            <span>Sync Customers</span>
          </button>
          <button className={styles.quickActionButton} onClick={handleSyncProducts}>
            <Package size={20} />
            <span>Sync Products</span>
          </button>
          <button className={styles.quickActionButton} onClick={handleSyncPrices}>
            <DollarSign size={20} />
            <span>Sync Prices</span>
          </button>
          <button className={styles.quickActionButton} onClick={handleSyncStock}>
            <BarChart3 size={20} />
            <span>Sync Stock</span>
          </button>
        </div>
      </div>

      {/* Connectors List */}
      <div className={styles.connectorsSection}>
        <div className={styles.sectionHeader}>
          <h2>ERP Connectors</h2>
          <span className={styles.connectorCount}>
            {connectorsData?.connectors?.length || 0} Connected
          </span>
        </div>
        
        {connectorsLoading ? (
          <div className={styles.loadingContainer}>
            <LoadingSpinner />
            <p>Loading ERP connectors...</p>
          </div>
        ) : connectorsError ? (
          <div className={styles.errorContainer}>
            <AlertTriangle size={24} />
            <p>Error loading connectors: {connectorsError.message}</p>
            <button onClick={() => refetchConnectors()}>Try Again</button>
          </div>
        ) : !connectorsData?.connectors?.length ? (
          <div className={styles.emptyState}>
            <Database size={48} />
            <h3>No ERP Connectors</h3>
            <p>Add your first ERP system integration to get started with data synchronization.</p>
            <button 
              className={styles.addFirstConnectorButton}
              onClick={() => router.push('/admin/erp/connectors/new')}
            >
              <Plus size={16} />
              Add First Connector
            </button>
          </div>
        ) : (
          <div className={styles.connectorsList}>
            {connectorsData.connectors.map((connector) => (
              <ConnectorItem
                key={connector.id}
                connector={connector}
                onClick={handleViewConnectorDetails}
                onSync={handleConnectorSync}
                onViewDetails={handleViewConnectorDetails}
              />
            ))}
          </div>
        )}
      </div>

      {/* Recent Sync History */}
      {syncHistoryData?.syncOperations?.length > 0 && (
        <div className={styles.syncHistorySection}>
          <h2>Recent Sync Operations</h2>
          <div className={styles.syncHistoryList}>
            {syncHistoryData.syncOperations.slice(0, 5).map((operation, index) => (
              <div key={operation.id || index} className={styles.syncHistoryItem}>
                <div className={styles.syncOperation}>{operation.operation || operation.type}</div>
                <div className={styles.syncStatus}>{operation.status}</div>
                <div className={styles.syncTime}>
                  {operation.completedAt ? new Date(operation.completedAt).toLocaleString() : 'In Progress'}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default ERPDashboard; 
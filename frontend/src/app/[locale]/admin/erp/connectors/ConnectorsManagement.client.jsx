"use client";

import React, { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Search,
  Filter,
  Plus,
  Settings,
  Activity,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Clock,
  RefreshCw,
  Eye,
  Trash2,
  ArrowLeft,
  Database,
  Link2,
  Power,
  Edit
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import {
  listConnectors,
  getConnectorStatus,
  removeConnector,
  toggleConnector,
  syncProducts,
  syncStock,
  syncPrices,
  syncCustomers,
  getConnectorTypes,
  getConnectorStatuses
} from '@/api/client/admin/erpApi';
import styles from './ConnectorsManagement.module.css';

const ConnectorCard = ({ connector, onClick, onEdit, onDelete, onToggle, onSync }) => {
  const t = useTranslations('ERPConnectors');
  const [isDeleting, setIsDeleting] = useState(false);
  const [isToggling, setIsToggling] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  
  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'active': return 'success';
      case 'inactive': return 'neutral';
      case 'syncing':
      case 'synchronizing': return 'info';
      case 'error': return 'error';
      case 'maintenance': return 'warning';
      default: return 'neutral';
    }
  };

  const getStatusIcon = (status) => {
    switch (status?.toLowerCase()) {
      case 'active': return CheckCircle;
      case 'inactive': return XCircle;
      case 'syncing':
      case 'synchronizing': return RefreshCw;
      case 'error': return AlertTriangle;
      case 'maintenance': return Clock;
      default: return Database;
    }
  };

  const getTypeLabel = (type) => {
    switch (type?.toLowerCase()) {
      case 'sap': return 'SAP ERP';
      case 'netsuite': return 'NetSuite';
      case 'odoo': return 'Odoo';
      case 'dynamics365': return 'Dynamics 365';
      case 'erpnext': return 'ERPNext';
      case 'frappe': return 'Frappe';
      default: return type?.toUpperCase() || 'Unknown';
    }
  };

  const StatusIcon = getStatusIcon(connector.status);

  const handleToggle = async (e) => {
    e.stopPropagation();
    setIsToggling(true);
    try {
      await onToggle(connector);
    } finally {
      setIsToggling(false);
    }
  };

  const handleDelete = async (e) => {
    e.stopPropagation();
    if (window.confirm(`Are you sure you want to delete connector "${connector.name}"?`)) {
      setIsDeleting(true);
      try {
        await onDelete(connector.id);
      } finally {
        setIsDeleting(false);
      }
    }
  };

  const handleSync = async (e) => {
    e.stopPropagation();
    setIsSyncing(true);
    try {
      await onSync(connector);
    } finally {
      setIsSyncing(false);
    }
  };

  const handleEdit = (e) => {
    e.stopPropagation();
    onEdit(connector);
  };

  return (
    <div 
      className={`${styles.connectorCard} ${styles[getStatusColor(connector.status)]}`}
      onClick={() => onClick(connector)}
    >
      <div className={styles.connectorHeader}>
        <div className={styles.connectorInfo}>
          <h4 className={styles.connectorName}>{connector.name}</h4>
          <span className={styles.connectorType}>{getTypeLabel(connector.type)}</span>
        </div>
        <div className={styles.connectorStatus}>
          <StatusIcon size={20} />
          <span>{(connector.status || 'Unknown').toUpperCase()}</span>
        </div>
      </div>
      
      <div className={styles.connectorDetails}>
        <div className={styles.detailItem}>
          <span className={styles.detailLabel}>Environment:</span>
          <span className={styles.detailValue}>{connector.environment || 'production'}</span>
        </div>
        <div className={styles.detailItem}>
          <span className={styles.detailLabel}>URL:</span>
          <span className={styles.detailValue}>{connector.baseUrl || 'Not configured'}</span>
        </div>
        <div className={styles.detailItem}>
          <span className={styles.detailLabel}>Last Sync:</span>
          <span className={styles.detailValue}>
            {connector.lastSyncTime ? new Date(connector.lastSyncTime).toLocaleString() : 'Never'}
          </span>
        </div>
      </div>

      <div className={styles.connectorActions}>
        <button 
          className={styles.actionButton}
          onClick={handleToggle}
          disabled={isToggling}
          title={connector.status === 'active' ? 'Deactivate' : 'Activate'}
        >
          <Power size={16} className={isToggling ? styles.spinning : ''} />
        </button>
        
        <button 
          className={styles.actionButton}
          onClick={handleSync}
          disabled={isSyncing || connector.status !== 'active'}
          title="Sync All Data"
        >
          <RefreshCw size={16} className={isSyncing ? styles.spinning : ''} />
        </button>
        
        <button 
          className={styles.actionButton}
          onClick={handleEdit}
          title="Edit Connector"
        >
          <Edit size={16} />
        </button>
        
        <button 
          className={`${styles.actionButton} ${styles.danger}`}
          onClick={handleDelete}
          disabled={isDeleting}
          title="Delete Connector"
        >
          <Trash2 size={16} />
        </button>
      </div>
    </div>
  );
};

const ConnectorsManagement = () => {
  const t = useTranslations('ERPConnectors');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');

  // Fetch connectors
  const { 
    data: connectorsData, 
    isLoading, 
    error,
    refetch 
  } = useQuery({
    queryKey: ['erp-connectors', filterType, filterStatus],
    queryFn: async () => {
      const params = {};
      if (filterType !== 'all') params.type = filterType;
      if (filterStatus !== 'all') params.status = filterStatus;
      
      const response = await listConnectors(params);
      
      return response;
    },
    staleTime: 30000,
  });

  // Delete connector mutation
  const deleteMutation = useMutation({
    mutationFn: async (connectorId) => {
      return await removeConnector(connectorId, true);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-connectors']);
      
    },
    onError: (error) => {
      
      alert('Failed to delete connector: ' + error.message);
    }
  });

  // Toggle connector mutation
  const toggleMutation = useMutation({
    mutationFn: async (connector) => {
      const isActive = connector.status === 'active';
      return await toggleConnector(connector.id, {
        activate: !isActive,
        reason: isActive ? 'Manual deactivation' : 'Manual activation'
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-connectors']);
      
    },
    onError: (error) => {
      
      alert('Failed to toggle connector: ' + error.message);
    }
  });

  // Sync connector mutation
  const syncMutation = useMutation({
    mutationFn: async (connector) => {
      // Sync all data types for the connector
      const results = await Promise.allSettled([
        syncProducts(connector.id),
        syncStock(connector.id),
        syncPrices(connector.id),
        syncCustomers(connector.id)
      ]);
      
      const failures = results.filter(r => r.status === 'rejected');
      if (failures.length > 0) {
        
      }
      
      return results;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['erp-connectors']);
      
    },
    onError: (error) => {
      
      alert('Failed to sync connector: ' + error.message);
    }
  });

  // Filter connectors
  const filteredConnectors = useMemo(() => {
    if (!connectorsData?.connectors) return [];
    
    return connectorsData.connectors.filter(connector => {
      const matchesSearch = searchTerm === '' || 
        connector.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        connector.type.toLowerCase().includes(searchTerm.toLowerCase());
      
      return matchesSearch;
    });
  }, [connectorsData, searchTerm]);

  const handleConnectorClick = (connector) => {
    router.push(`/admin/erp/connectors/${connector.id}`);
  };

  const handleEditConnector = (connector) => {
    router.push(`/admin/erp/connectors/${connector.id}/edit`);
  };

  const handleAddConnector = () => {
    router.push('/admin/erp/connectors/new');
  };

  if (!isAdmin) {
    return (
      <div className={styles.accessDenied}>
        <AlertTriangle size={48} />
        <h2>Access Denied</h2>
        <p>You don't have permission to manage ERP connectors.</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <button 
          className={styles.backButton}
          onClick={() => router.push('/admin/erp')}
        >
          <ArrowLeft size={20} />
          Back to ERP Dashboard
        </button>
        
        <div className={styles.headerActions}>
          <button 
            className={styles.refreshButton}
            onClick={() => refetch()}
            disabled={isLoading}
          >
            <RefreshCw size={16} className={isLoading ? styles.spinning : ''} />
            Refresh
          </button>
          
          <button 
            className={styles.addButton}
            onClick={handleAddConnector}
          >
            <Plus size={16} />
            Add Connector
          </button>
        </div>
      </div>

      {/* Title */}
      <div className={styles.titleSection}>
        <h1 className={styles.title}>ERP Connectors</h1>
        <p className={styles.subtitle}>
          Manage your enterprise resource planning system integrations
        </p>
      </div>

      {/* Filters */}
      <div className={styles.filters}>
        <div className={styles.searchBox}>
          <Search size={20} />
          <input
            type="text"
            placeholder="Search connectors..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        
        <div className={styles.filterGroup}>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className={styles.filterSelect}
          >
            <option value="all">All Types</option>
            <option value="sap">SAP</option>
            <option value="netsuite">NetSuite</option>
            <option value="odoo">Odoo</option>
            <option value="dynamics365">Dynamics 365</option>
            <option value="erpnext">ERPNext</option>
            <option value="frappe">Frappe</option>
          </select>
          
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className={styles.filterSelect}
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="error">Error</option>
            <option value="maintenance">Maintenance</option>
          </select>
        </div>
      </div>

      {/* Connectors Grid */}
      {isLoading ? (
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <p>Loading connectors...</p>
        </div>
      ) : error ? (
        <div className={styles.errorContainer}>
          <AlertTriangle size={24} />
          <p>Error loading connectors: {error.message}</p>
          <button onClick={() => refetch()}>Try Again</button>
        </div>
      ) : filteredConnectors.length === 0 ? (
        <div className={styles.emptyState}>
          <Database size={48} />
          <h3>No Connectors Found</h3>
          <p>
            {searchTerm || filterType !== 'all' || filterStatus !== 'all' 
              ? 'No connectors match your filters. Try adjusting your search criteria.'
              : 'Add your first ERP connector to start integrating with your business systems.'}
          </p>
          {!searchTerm && filterType === 'all' && filterStatus === 'all' && (
            <button 
              className={styles.addFirstButton}
              onClick={handleAddConnector}
            >
              <Plus size={16} />
              Add Your First Connector
            </button>
          )}
        </div>
      ) : (
        <div className={styles.connectorsGrid}>
          {filteredConnectors.map((connector) => (
            <ConnectorCard
              key={connector.id}
              connector={connector}
              onClick={handleConnectorClick}
              onEdit={handleEditConnector}
              onDelete={(id) => deleteMutation.mutate(id)}
              onToggle={(connector) => toggleMutation.mutate(connector)}
              onSync={(connector) => syncMutation.mutate(connector)}
            />
          ))}
        </div>
      )}

      {/* Summary */}
      {connectorsData && (
        <div className={styles.summary}>
          <span>Total: {connectorsData.totalCount || filteredConnectors.length} connectors</span>
          <span>Showing: {filteredConnectors.length} connectors</span>
        </div>
      )}
    </div>
  );
};

export default ConnectorsManagement;
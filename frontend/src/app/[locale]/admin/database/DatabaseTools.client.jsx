"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation } from '@tanstack/react-query';
import {
  ArrowLeft,
  Database,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  Activity,
  HardDrive,
  Zap,
  Download,
  Upload,
  Trash2,
  Settings,
  Monitor,
  BarChart3,
  Shield,
  Clock,
  Server,
  FileText,
  Play,
  Pause,
  XCircle
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getDatabaseStats,
  runDatabaseMaintenance,
  optimizeDatabase,
  backupDatabase,
  restoreDatabase,
  cleanupOldData,
  getSlowQueries,
  getDatabaseConnections
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './DatabaseTools.module.css';

const StatCard = ({ title, value, unit, icon: Icon, status, onClick, loading }) => (
  <div 
    className={`${styles.statCard} ${onClick ? styles.clickable : ''} ${status ? styles[status] : ''}`}
    onClick={onClick}
  >
    <div className={styles.statIcon}>
      <Icon size={20} />
    </div>
    <div className={styles.statContent}>
      <div className={styles.statValue}>
        {loading ? '...' : value}
        {unit && <span className={styles.unit}>{unit}</span>}
      </div>
      <div className={styles.statLabel}>{title}</div>
    </div>
  </div>
);

const ActionCard = ({ title, description, icon: Icon, onClick, variant = 'default', disabled, loading }) => (
  <button
    className={`${styles.actionCard} ${styles[variant]} ${disabled ? styles.disabled : ''}`}
    onClick={onClick}
    disabled={disabled || loading}
  >
    <div className={styles.actionIcon}>
      <Icon size={24} />
    </div>
    <div className={styles.actionContent}>
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
    {loading && (
      <div className={styles.actionLoader}>
        <RefreshCw size={16} className={styles.spinning} />
      </div>
    )}
  </button>
);

const DatabaseTools = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('DatabaseTools');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  const [activeOperation, setActiveOperation] = useState(null);

  // Fetch database statistics
  const { data: dbStats, isLoading: statsLoading, refetch: refetchStats } = useQuery({
    queryKey: ['database', 'stats'],
    queryFn: getDatabaseStats,
    staleTime: 30 * 1000,
    enabled: isAdmin
  });

  const { data: slowQueries, isLoading: queriesLoading, refetch: refetchQueries } = useQuery({
    queryKey: ['database', 'slow-queries'],
    queryFn: getSlowQueries,
    staleTime: 60 * 1000,
    enabled: isAdmin
  });

  const { data: connections, isLoading: connectionsLoading, refetch: refetchConnections } = useQuery({
    queryKey: ['database', 'connections'],
    queryFn: getDatabaseConnections,
    staleTime: 10 * 1000,
    enabled: isAdmin
  });

  // Database operation mutations
  const maintenanceMutation = useMutation({
    mutationFn: runDatabaseMaintenance,
    onMutate: () => setActiveOperation('maintenance'),
    onSuccess: () => {
      refetchStats();
      setActiveOperation(null);
    },
    onError: (error) => {
      // Error: 'Maintenance failed:', error...
      setActiveOperation(null);
    }
  });

  const optimizeMutation = useMutation({
    mutationFn: optimizeDatabase,
    onMutate: () => setActiveOperation('optimize'),
    onSuccess: () => {
      refetchStats();
      setActiveOperation(null);
    },
    onError: (error) => {
      // Error: 'Optimization failed:', error...
      setActiveOperation(null);
    }
  });

  const backupMutation = useMutation({
    mutationFn: backupDatabase,
    onMutate: () => setActiveOperation('backup'),
    onSuccess: () => {
      setActiveOperation(null);
    },
    onError: (error) => {
      // Error: 'Backup failed:', error...
      setActiveOperation(null);
    }
  });

  const cleanupMutation = useMutation({
    mutationFn: cleanupOldData,
    onMutate: () => setActiveOperation('cleanup'),
    onSuccess: () => {
      refetchStats();
      setActiveOperation(null);
    },
    onError: (error) => {
      // Error: 'Cleanup failed:', error...
      setActiveOperation(null);
    }
  });

  const handleMaintenance = useCallback(() => {
    if (window.confirm('This will temporarily impact database performance. Continue?')) {
      maintenanceMutation.mutate();
    }
  }, [maintenanceMutation]);

  const handleOptimize = useCallback(() => {
    if (window.confirm('Database optimization may take several minutes. Continue?')) {
      optimizeMutation.mutate();
    }
  }, [optimizeMutation]);

  const handleBackup = useCallback(() => {
    backupMutation.mutate();
  }, [backupMutation]);

  const handleCleanup = useCallback(() => {
    if (window.confirm('This will permanently delete old data. Are you sure?')) {
      cleanupMutation.mutate();
    }
  }, [cleanupMutation]);

  const handleRefresh = useCallback(() => {
    refetchStats();
    refetchQueries();
    refetchConnections();
  }, [refetchStats, refetchQueries, refetchConnections]);

  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
  };

  const getHealthStatus = (stats) => {
    if (!stats) return 'unknown';
    if (stats.diskUsage > 90 || stats.connectionUsage > 90) return 'critical';
    if (stats.diskUsage > 80 || stats.connectionUsage > 80) return 'warning';
    return 'healthy';
  };

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access database tools.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  const isLoading = statsLoading || queriesLoading || connectionsLoading;

  if (isLoading && !dbStats) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Database Tools...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch database information.' })}</p>
        </div>
      </div>
    );
  }

  const healthStatus = getHealthStatus(dbStats);

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
                <Database size={24} />
                {t('title', { defaultValue: 'Database Tools' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Database maintenance and optimization tools' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button 
              className={styles.refreshButton}
              onClick={handleRefresh}
              disabled={!!activeOperation}
            >
              <RefreshCw size={16} className={activeOperation ? styles.spinning : ''} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Health Status Banner */}
        <div className={`${styles.healthBanner} ${styles[healthStatus]}`}>
          <div className={styles.healthInfo}>
            {healthStatus === 'healthy' && <CheckCircle size={20} />}
            {healthStatus === 'warning' && <AlertTriangle size={20} />}
            {healthStatus === 'critical' && <XCircle size={20} />}
            <span>
              {healthStatus === 'healthy' && t('databaseHealthy', { defaultValue: 'Database is running smoothly' })}
              {healthStatus === 'warning' && t('databaseWarning', { defaultValue: 'Database requires attention' })}
              {healthStatus === 'critical' && t('databaseCritical', { defaultValue: 'Database needs immediate attention' })}
            </span>
          </div>
        </div>

        {/* Stats Grid */}
        <div className={styles.statsGrid}>
          <StatCard
            title={t('diskUsage', { defaultValue: 'Disk Usage' })}
            value={dbStats?.diskUsage || 0}
            unit="%"
            icon={HardDrive}
            status={dbStats?.diskUsage > 90 ? 'critical' : dbStats?.diskUsage > 80 ? 'warning' : 'healthy'}
            loading={statsLoading}
          />
          <StatCard
            title={t('connections', { defaultValue: 'Active Connections' })}
            value={`${connections?.active || 0}/${connections?.max || 100}`}
            icon={Activity}
            status={connections?.usage > 90 ? 'critical' : connections?.usage > 80 ? 'warning' : 'healthy'}
            loading={connectionsLoading}
          />
          <StatCard
            title={t('databaseSize', { defaultValue: 'Database Size' })}
            value={formatBytes(dbStats?.totalSize)}
            icon={Database}
            loading={statsLoading}
          />
          <StatCard
            title={t('slowQueries', { defaultValue: 'Slow Queries' })}
            value={slowQueries?.count || 0}
            icon={Clock}
            status={slowQueries?.count > 10 ? 'warning' : 'healthy'}
            loading={queriesLoading}
          />
          <StatCard
            title={t('uptime', { defaultValue: 'Uptime' })}
            value={dbStats?.uptime || '0h'}
            icon={Server}
            loading={statsLoading}
          />
          <StatCard
            title={t('queriesPerSecond', { defaultValue: 'Queries/sec' })}
            value={dbStats?.qps || 0}
            icon={Zap}
            loading={statsLoading}
          />
        </div>

        {/* Operations Grid */}
        <div className={styles.operationsSection}>
          <h2 className={styles.sectionTitle}>{t('operations', { defaultValue: 'Database Operations' })}</h2>
          <div className={styles.operationsGrid}>
            <ActionCard
              title={t('runMaintenance', { defaultValue: 'Run Maintenance' })}
              description={t('maintenanceDesc', { defaultValue: 'Analyze tables, rebuild indexes, and update statistics' })}
              icon={Settings}
              onClick={handleMaintenance}
              variant="primary"
              loading={activeOperation === 'maintenance'}
            />
            <ActionCard
              title={t('optimizeDatabase', { defaultValue: 'Optimize Database' })}
              description={t('optimizeDesc', { defaultValue: 'Defragment tables and optimize query performance' })}
              icon={Zap}
              onClick={handleOptimize}
              loading={activeOperation === 'optimize'}
            />
            <ActionCard
              title={t('createBackup', { defaultValue: 'Create Backup' })}
              description={t('backupDesc', { defaultValue: 'Create a full database backup for disaster recovery' })}
              icon={Download}
              onClick={handleBackup}
              loading={activeOperation === 'backup'}
            />
            <ActionCard
              title={t('cleanupOldData', { defaultValue: 'Cleanup Old Data' })}
              description={t('cleanupDesc', { defaultValue: 'Remove expired logs and temporary data' })}
              icon={Trash2}
              onClick={handleCleanup}
              variant="danger"
              loading={activeOperation === 'cleanup'}
            />
          </div>
        </div>

        {/* Monitoring Section */}
        <div className={styles.monitoringSection}>
          <div className={styles.monitoringCard}>
            <h3>{t('slowQueries', { defaultValue: 'Slow Queries' })}</h3>
            {queriesLoading ? (
              <div className={styles.loadingState}>
                <LoadingSpinner size={24} />
                <p>{t('loadingQueries', { defaultValue: 'Loading queries...' })}</p>
              </div>
            ) : slowQueries?.queries?.length > 0 ? (
              <div className={styles.queriesList}>
                {slowQueries.queries.slice(0, 5).map((query, index) => (
                  <div key={index} className={styles.queryItem}>
                    <div className={styles.queryInfo}>
                      <code className={styles.queryText}>{query.sql}</code>
                      <div className={styles.queryMeta}>
                        <span>{query.avgTime}ms avg</span>
                        <span>{query.execCount} executions</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className={styles.emptyState}>
                <CheckCircle size={32} />
                <p>{t('noSlowQueries', { defaultValue: 'No slow queries detected' })}</p>
              </div>
            )}
          </div>

          <div className={styles.monitoringCard}>
            <h3>{t('recentBackups', { defaultValue: 'Recent Backups' })}</h3>
            {dbStats?.recentBackups?.length > 0 ? (
              <div className={styles.backupsList}>
                {dbStats.recentBackups.map((backup, index) => (
                  <div key={index} className={styles.backupItem}>
                    <div className={styles.backupInfo}>
                      <div className={styles.backupName}>{backup.name}</div>
                      <div className={styles.backupMeta}>
                        <span>{formatBytes(backup.size)}</span>
                        <span>{new Date(backup.createdAt).toLocaleDateString()}</span>
                      </div>
                    </div>
                    <div className={styles.backupStatus}>
                      <CheckCircle size={16} />
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className={styles.emptyState}>
                <FileText size={32} />
                <p>{t('noBackups', { defaultValue: 'No recent backups found' })}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default DatabaseTools; 
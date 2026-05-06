'use client';

import { useState, useCallback, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { 
  Calendar, 
  Clock, 
  User, 
  Search, 
  Filter, 
  Plus, 
  Play, 
  Pause, 
  CheckCircle, 
  XCircle, 
  AlertCircle, 
  MoreHorizontal,
  Trash2,
  Eye,
  RefreshCw,
  ArrowLeft,
  Activity,
  Timer,
  Target,
  TrendingUp,
  Users,
  Zap,
  Archive,
  Edit,
  Send
} from 'lucide-react';
import LoadingSpinner from '@/components/common/LoadingSpinner';
import { schedulerApi } from '@/api/schedulerApi';
import styles from './SchedulerManagement.module.css';

export default function SchedulerManagement() {
  const t = useTranslations('SchedulerManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  // State management
  const [filters, setFilters] = useState({
    search: '',
    status: 'all',
    userId: '',
    sortBy: 'createdAt',
    sortOrder: 'desc'
  });
  const [selectedScheduler, setSelectedScheduler] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showActionModal, setShowActionModal] = useState(false);

  // Data fetching
  const { 
    data: schedulersData, 
    isLoading: schedulersLoading, 
    error: schedulersError,
    refetch: refetchSchedulers
  } = useQuery({
    queryKey: ['schedulers', filters],
    queryFn: () => schedulerApi.getSchedulers(filters.userId || undefined),
    staleTime: 30000,
  });

  const { 
    data: actionsData, 
    isLoading: actionsLoading 
  } = useQuery({
    queryKey: ['scheduler-actions', selectedScheduler?.id],
    queryFn: () => selectedScheduler ? schedulerApi.getActions(selectedScheduler.id) : null,
    enabled: !!selectedScheduler,
    staleTime: 15000,
  });

  // Mutations
  const createSchedulerMutation = useMutation({
    mutationFn: schedulerApi.createScheduler,
    onSuccess: () => {
      queryClient.invalidateQueries(['schedulers']);
      setShowCreateModal(false);
    },
  });

  const addActionMutation = useMutation({
    mutationFn: schedulerApi.addAction,
    onSuccess: () => {
      queryClient.invalidateQueries(['scheduler-actions']);
      setShowActionModal(false);
    },
  });

  const removeActionMutation = useMutation({
    mutationFn: schedulerApi.removeAction,
    onSuccess: () => {
      queryClient.invalidateQueries(['scheduler-actions']);
    },
  });

  // Computed values
  const schedulers = schedulersData?.activities || [];
  const actions = actionsData?.actions || [];
  
  const filteredSchedulers = useMemo(() => {
    return schedulers.filter(scheduler => {
      const matchesSearch = !filters.search || 
        scheduler.id.toLowerCase().includes(filters.search.toLowerCase()) ||
        scheduler.userId.toLowerCase().includes(filters.search.toLowerCase());
      
      const matchesStatus = filters.status === 'all' || 
        (filters.status === 'archived' && scheduler.archived) ||
        (filters.status === 'active' && !scheduler.archived);

      return matchesSearch && matchesStatus;
    });
  }, [schedulers, filters]);

  const stats = useMemo(() => {
    const totalSchedulers = schedulers.length;
    const activeSchedulers = schedulers.filter(s => !s.archived).length;
    const totalActions = actions.length;
    const completedActions = actions.filter(a => a.status === 'completed').length;
    const pendingActions = actions.filter(a => a.status === 'pending').length;
    const failedActions = actions.filter(a => a.status === 'failed').length;
    
    return {
      totalSchedulers,
      activeSchedulers,
      totalActions,
      completedActions,
      pendingActions,
      failedActions,
      successRate: totalActions > 0 ? Math.round((completedActions / totalActions) * 100) : 0
    };
  }, [schedulers, actions]);

  // Event handlers
  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleRefresh = useCallback(() => {
    refetchSchedulers();
    if (selectedScheduler) {
      queryClient.invalidateQueries(['scheduler-actions', selectedScheduler.id]);
    }
  }, [refetchSchedulers, queryClient, selectedScheduler]);

  const handleCreateScheduler = useCallback((userId) => {
    createSchedulerMutation.mutate({ userId });
  }, [createSchedulerMutation]);

  const handleAddAction = useCallback((data) => {
    addActionMutation.mutate(data);
  }, [addActionMutation]);

  const handleRemoveAction = useCallback((actionId) => {
    removeActionMutation.mutate(actionId);
  }, [removeActionMutation]);

  const getStatusIcon = (status) => {
    switch (status) {
      case 'completed': return <CheckCircle size={16} />;
      case 'failed': return <XCircle size={16} />;
      case 'executing': return <Play size={16} />;
      case 'pending': return <Clock size={16} />;
      default: return <AlertCircle size={16} />;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'completed': return 'success';
      case 'failed': return 'error';
      case 'executing': return 'info';
      case 'pending': return 'warning';
      default: return 'neutral';
    }
  };

  // Access control check - render different content but don't break hooks
  if (!isAdmin) {
    return (
      <div className={styles.accessDenied}>
        <div className={styles.accessDeniedContent}>
          <XCircle className={styles.accessDeniedIcon} />
          <h2 className={styles.accessDeniedTitle}>{t('accessDenied')}</h2>
          <p className={styles.accessDeniedText}>{t('adminRequired')}</p>
          <p className={styles.accessDeniedText}>
            {t('currentRole')}: {user?.role || t('notLoggedIn')}
          </p>
          <button 
            onClick={() => router.push('/admin')}
            className={styles.backButton}
          >
            <ArrowLeft size={16} />
            {t('backToDashboard')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerMain}>
          <button 
            onClick={() => router.push('/admin')}
            className={styles.backButton}
            aria-label={t('backToDashboard')}
          >
            <ArrowLeft size={20} />
          </button>
          <div className={styles.headerInfo}>
            <h1 className={styles.headerTitle}>{t('title')}</h1>
            <p className={styles.headerSubtitle}>{t('subtitle')}</p>
          </div>
        </div>
        <div className={styles.headerActions}>
          <button 
            onClick={handleRefresh}
            className={styles.refreshButton}
            disabled={schedulersLoading}
          >
            <RefreshCw size={16} className={schedulersLoading ? styles.spinning : ''} />
            {t('refresh')}
          </button>
          <button 
            onClick={() => setShowCreateModal(true)}
            className={styles.createButton}
          >
            <Plus size={16} />
            {t('createScheduler')}
          </button>
        </div>
      </div>

      {/* Stats Overview */}
      <div className={styles.statsGrid}>
        <div className={styles.statCard}>
          <div className={styles.statIcon}>
            <Users size={20} />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statValue}>{stats.totalSchedulers}</div>
            <div className={styles.statLabel}>{t('totalSchedulers')}</div>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIcon}>
            <Activity size={20} />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statValue}>{stats.activeSchedulers}</div>
            <div className={styles.statLabel}>{t('activeSchedulers')}</div>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIcon}>
            <Target size={20} />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statValue}>{stats.totalActions}</div>
            <div className={styles.statLabel}>{t('totalActions')}</div>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIcon}>
            <TrendingUp size={20} />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statValue}>{stats.successRate}%</div>
            <div className={styles.statLabel}>{t('successRate')}</div>
          </div>
        </div>
      </div>

      {/* Filters Bar */}
      <div className={styles.filtersBar}>
        <div className={styles.filtersLeft}>
          <div className={styles.searchContainer}>
            <Search className={styles.searchIcon} size={16} />
            <input
              type="text"
              placeholder={t('searchSchedulers')}
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              className={styles.searchInput}
            />
          </div>
          
          <select
            value={filters.status}
            onChange={(e) => handleFilterChange('status', e.target.value)}
            className={styles.filterSelect}
          >
            <option value="all">{t('allStatus')}</option>
            <option value="active">{t('active')}</option>
            <option value="archived">{t('archived')}</option>
          </select>

          <input
            type="text"
            placeholder={t('filterByUserId')}
            value={filters.userId}
            onChange={(e) => handleFilterChange('userId', e.target.value)}
            className={styles.userIdInput}
          />
        </div>

        <div className={styles.filtersRight}>
          <select
            value={`${filters.sortBy}-${filters.sortOrder}`}
            onChange={(e) => {
              const [sortBy, sortOrder] = e.target.value.split('-');
              handleFilterChange('sortBy', sortBy);
              handleFilterChange('sortOrder', sortOrder);
            }}
            className={styles.sortSelect}
          >
            <option value="createdAt-desc">{t('sortNewestFirst')}</option>
            <option value="createdAt-asc">{t('sortOldestFirst')}</option>
            <option value="userId-asc">{t('sortByUserId')}</option>
          </select>
        </div>
      </div>

      {/* Main Content */}
      <div className={styles.mainContent}>
        {/* Schedulers List */}
        <div className={styles.schedulersSection}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>
              {t('schedulers')} ({filteredSchedulers.length})
            </h2>
          </div>

          {schedulersLoading ? (
            <div className={styles.loadingContainer}>
              <LoadingSpinner />
              <p className={styles.loadingText}>{t('loadingSchedulers')}</p>
            </div>
          ) : schedulersError ? (
            <div className={styles.errorContainer}>
              <XCircle size={48} className={styles.errorIcon} />
              <h3 className={styles.errorTitle}>{t('loadingError')}</h3>
              <p className={styles.errorText}>{schedulersError.message}</p>
              <button onClick={handleRefresh} className={styles.retryButton}>
                <RefreshCw size={16} />
                {t('retry')}
              </button>
            </div>
          ) : filteredSchedulers.length === 0 ? (
            <div className={styles.emptyState}>
              <Calendar size={48} className={styles.emptyIcon} />
              <h3 className={styles.emptyTitle}>{t('noSchedulersFound')}</h3>
              <p className={styles.emptyText}>{t('noSchedulersMessage')}</p>
              <button 
                onClick={() => setShowCreateModal(true)}
                className={styles.createButton}
              >
                <Plus size={16} />
                {t('createFirstScheduler')}
              </button>
            </div>
          ) : (
            <div className={styles.schedulersList}>
              {filteredSchedulers.map((scheduler) => (
                <div 
                  key={scheduler.id} 
                  className={`${styles.schedulerCard} ${selectedScheduler?.id === scheduler.id ? styles.selected : ''}`}
                  onClick={() => setSelectedScheduler(scheduler)}
                >
                  <div className={styles.schedulerMain}>
                    <div className={styles.schedulerInfo}>
                      <div className={styles.schedulerHeader}>
                        <h4 className={styles.schedulerTitle}>
                          {t('scheduler')} #{scheduler.id}
                        </h4>
                                                 <div className={scheduler.archived ? styles.statusBadgeArchived : styles.statusBadgeActive}>
                           {scheduler.archived ? (
                             <><Archive size={12} /> {t('archived')}</>
                           ) : (
                             <><Zap size={12} /> {t('active')}</>
                           )}
                         </div>
                      </div>
                      <div className={styles.schedulerMeta}>
                        <div className={styles.metaItem}>
                          <User size={14} />
                          <span>{t('user')}: {scheduler.userId}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  <div className={styles.schedulerActions}>
                    <button 
                      className={styles.actionButton}
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedScheduler(scheduler);
                        setShowActionModal(true);
                      }}
                      title={t('addAction')}
                    >
                      <Plus size={16} />
                    </button>
                    <button 
                      className={styles.actionButton}
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedScheduler(scheduler);
                      }}
                      title={t('viewActions')}
                    >
                      <Eye size={16} />
                    </button>
                    <button 
                      className={`${styles.actionButton} ${styles.archiveButton}`}
                      onClick={(e) => {
                        e.stopPropagation();
                        // Handle archive/unarchive
                      }}
                      title={scheduler.archived ? t('unarchive') : t('archive')}
                    >
                      {scheduler.archived ? <Send size={16} /> : <Archive size={16} />}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Actions Panel */}
        {selectedScheduler && (
          <div className={styles.actionsPanel}>
            <div className={styles.sectionHeader}>
              <h2 className={styles.sectionTitle}>
                {t('actions')} - {t('scheduler')} #{selectedScheduler.id}
              </h2>
              <button 
                onClick={() => setShowActionModal(true)}
                className={styles.addActionButton}
              >
                <Plus size={16} />
                {t('addAction')}
              </button>
            </div>

            {actionsLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p className={styles.loadingText}>{t('loadingActions')}</p>
              </div>
            ) : actions.length === 0 ? (
              <div className={styles.emptyState}>
                <Timer size={48} className={styles.emptyIcon} />
                <h3 className={styles.emptyTitle}>{t('noActionsFound')}</h3>
                <p className={styles.emptyText}>{t('noActionsMessage')}</p>
                <button 
                  onClick={() => setShowActionModal(true)}
                  className={styles.createButton}
                >
                  <Plus size={16} />
                  {t('addFirstAction')}
                </button>
              </div>
            ) : (
              <div className={styles.actionsList}>
                {actions.map((action) => (
                  <div key={action.id} className={styles.actionCard}>
                    <div className={styles.actionMain}>
                      <div className={styles.actionInfo}>
                        <div className={styles.actionHeader}>
                          <h5 className={styles.actionTitle}>
                            {action.naturalLanguageTask}
                          </h5>
                                                     <div className={styles[`statusBadge${getStatusColor(action.status).charAt(0).toUpperCase() + getStatusColor(action.status).slice(1)}`]}>
                             {getStatusIcon(action.status)}
                             {action.status}
                           </div>
                        </div>
                        <div className={styles.actionMeta}>
                          <div className={styles.metaItem}>
                            <Clock size={14} />
                            <span>{t('execution')}: {new Date(action.executionTime).toLocaleString()}</span>
                          </div>
                          <div className={styles.metaItem}>
                            <Calendar size={14} />
                            <span>{t('created')}: {new Date(action.createdAt).toLocaleString()}</span>
                          </div>
                          {action.executedAt && (
                            <div className={styles.metaItem}>
                              <CheckCircle size={14} />
                              <span>{t('executed')}: {new Date(action.executedAt).toLocaleString()}</span>
                            </div>
                          )}
                        </div>
                        {action.result && (
                          <div className={styles.actionResult}>
                            <strong>{t('result')}:</strong> {action.result}
                          </div>
                        )}
                        {action.errorMessage && (
                          <div className={styles.actionError}>
                            <strong>{t('error')}:</strong> {action.errorMessage}
                          </div>
                        )}
                      </div>
                    </div>
                    
                    <div className={styles.actionActions}>
                      <button 
                        className={`${styles.actionButton} ${styles.deleteButton}`}
                        onClick={() => handleRemoveAction(action.id)}
                        title={t('removeAction')}
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Create Scheduler Modal */}
      {showCreateModal && (
        <CreateSchedulerModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateScheduler}
          isLoading={createSchedulerMutation.isLoading}
        />
      )}

      {/* Add Action Modal */}
      {showActionModal && selectedScheduler && (
        <AddActionModal
          schedulerId={selectedScheduler.id}
          onClose={() => setShowActionModal(false)}
          onSubmit={handleAddAction}
          isLoading={addActionMutation.isLoading}
        />
      )}
    </div>
  );
}

// Create Scheduler Modal Component
function CreateSchedulerModal({ onClose, onSubmit, isLoading }) {
  const t = useTranslations('SchedulerManagement');
  const [userId, setUserId] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (userId.trim()) {
      onSubmit({ userId: userId.trim() });
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3 className={styles.modalTitle}>{t('createScheduler')}</h3>
          <button onClick={onClose} className={styles.closeButton}>
            <XCircle size={20} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label htmlFor="userId" className={styles.formLabel}>
              {t('userId')}
            </label>
            <input
              id="userId"
              type="text"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              placeholder={t('enterUserId')}
              className={styles.formInput}
              required
            />
          </div>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={onClose}
              className={styles.cancelButton}
            >
              {t('cancel')}
            </button>
            <button
              type="submit"
              disabled={!userId.trim() || isLoading}
              className={styles.submitButton}
            >
              {isLoading ? <LoadingSpinner size="small" /> : <Plus size={16} />}
              {t('create')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// Add Action Modal Component
function AddActionModal({ schedulerId, onClose, onSubmit, isLoading }) {
  const t = useTranslations('SchedulerManagement');
  const [task, setTask] = useState('');
  const [executionTime, setExecutionTime] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (task.trim() && executionTime) {
      onSubmit({
        schedulerId,
        naturalLanguageTask: task.trim(),
        executionTime: new Date(executionTime).toISOString()
      });
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3 className={styles.modalTitle}>{t('addAction')}</h3>
          <button onClick={onClose} className={styles.closeButton}>
            <XCircle size={20} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label htmlFor="task" className={styles.formLabel}>
              {t('naturalLanguageTask')}
            </label>
            <textarea
              id="task"
              value={task}
              onChange={(e) => setTask(e.target.value)}
              placeholder={t('enterTaskDescription')}
              className={styles.formTextarea}
              rows={3}
              required
            />
          </div>
          <div className={styles.formGroup}>
            <label htmlFor="executionTime" className={styles.formLabel}>
              {t('executionTime')}
            </label>
            <input
              id="executionTime"
              type="datetime-local"
              value={executionTime}
              onChange={(e) => setExecutionTime(e.target.value)}
              className={styles.formInput}
              required
            />
          </div>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={onClose}
              className={styles.cancelButton}
            >
              {t('cancel')}
            </button>
            <button
              type="submit"
              disabled={!task.trim() || !executionTime || isLoading}
              className={styles.submitButton}
            >
              {isLoading ? <LoadingSpinner size="small" /> : <Plus size={16} />}
              {t('addAction')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
} 
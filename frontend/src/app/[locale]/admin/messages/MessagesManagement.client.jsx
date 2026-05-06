"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  MessageCircle,
  Users,
  Clock,
  CheckCircle,
  AlertCircle,
  Archive,
  Flag,
  Search,
  Filter,
  Download,
  RefreshCw,
  Eye,
  Trash2,
  MoreVertical,
  User,
  Calendar,
  Send,
  Reply,
  ExternalLink,
  TrendingUp,
  Activity,
  MessageSquare,
  Zap,
  Timer
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  listConversations,
  getConversation,
  getMessages,
  archiveConversation,
  restoreConversation,
  deleteMessage,
  markConversationAsRead,
  getMessagingStats,
  getAllActiveConversations,
  getMessageAnalytics
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './MessagesManagement.module.css';

const ConversationStatusBadge = ({ status }) => {
  const statusConfig = {
    active: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Active', icon: CheckCircle },
    archived: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Archived', icon: Archive },
    flagged: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Flagged', icon: Flag },
    unread: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Unread', icon: AlertCircle }
  };

  const config = statusConfig[status] || statusConfig.active;
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

const ConversationRow = ({ conversation, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);

  const formatTimeAgo = (date) => {
    const now = new Date();
    const messageDate = new Date(date);
    const diffInMinutes = Math.floor((now - messageDate) / (1000 * 60));
    
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
    return `${Math.floor(diffInMinutes / 1440)}d ago`;
  };

  return (
    <tr className={styles.conversationRow}>
      <td className={styles.conversationCell}>
        <div className={styles.conversationInfo}>
          <div className={styles.conversationHeader}>
            <MessageCircle size={16} />
            <span className={styles.conversationId}>#{conversation.id.slice(-8)}</span>
            {conversation.unreadCount > 0 && (
              <span className={styles.unreadBadge}>{conversation.unreadCount}</span>
            )}
          </div>
          <div className={styles.conversationPreview}>
            {conversation.lastMessage?.content?.substring(0, 60)}...
          </div>
        </div>
      </td>
      <td className={styles.participantsCell}>
        <div className={styles.participants}>
          <Users size={16} />
          <div className={styles.participantsList}>
            {conversation.participants?.map((participant, index) => (
              <div key={participant.id} className={styles.participant}>
                <span className={styles.participantName}>{participant.name}</span>
                {index < conversation.participants.length - 1 && (
                  <span className={styles.participantSeparator}>, </span>
                )}
              </div>
            ))}
          </div>
        </div>
      </td>
      <td className={styles.messagesCell}>
        <div className={styles.messageCount}>
          <MessageSquare size={14} />
          <span>{conversation.messageCount || 0}</span>
        </div>
      </td>
      <td className={styles.statusCell}>
        <ConversationStatusBadge status={conversation.status} />
      </td>
      <td className={styles.lastActivityCell}>
        <div className={styles.lastActivity}>
          <Clock size={14} />
          <span>{formatTimeAgo(conversation.lastActivity)}</span>
        </div>
      </td>
      <td className={styles.itemCell}>
        <div className={styles.itemInfo}>
          {conversation.itemId && (
            <>
              <span className={styles.itemType}>{conversation.itemType}</span>
              <span className={styles.itemTitle}>{conversation.itemTitle}</span>
            </>
          )}
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
              <button onClick={() => onAction('view', conversation)}>
                <Eye size={14} />
                View Messages
              </button>
              <button onClick={() => onAction('markRead', conversation)}>
                <CheckCircle size={14} />
                Mark as Read
              </button>
              {conversation.status === 'active' ? (
                <button onClick={() => onAction('archive', conversation)}>
                  <Archive size={14} />
                  Archive
                </button>
              ) : (
                <button onClick={() => onAction('restore', conversation)}>
                  <Activity size={14} />
                  Restore
                </button>
              )}
              {conversation.itemId && (
                <button onClick={() => onAction('viewItem', conversation)}>
                  <ExternalLink size={14} />
                  View Item
                </button>
              )}
              <button onClick={() => onAction('flag', conversation)} className={styles.dangerAction}>
                <Flag size={14} />
                Flag Conversation
              </button>
            </div>
          )}
        </div>
      </td>
    </tr>
  );
};

const ConversationModal = ({ conversation, onClose }) => {
  const t = useTranslations('MessagesManagement');
  const { data: messages, isLoading } = useQuery({
    queryKey: ['conversationMessages', conversation?.id],
    queryFn: () => getMessages(conversation.id),
    enabled: !!conversation
  });

  if (!conversation) return null;

  return (
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <div className={styles.modalTitle}>
            <MessageCircle size={20} />
            <h3>Conversation #{conversation.id.slice(-8)}</h3>
          </div>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <div className={styles.modalBody}>
          <div className={styles.conversationHeader}>
            <div className={styles.participantInfo}>
              <strong>Participants:</strong>
              <div className={styles.participantChips}>
                {conversation.participants?.map(participant => (
                  <span key={participant.id} className={styles.participantChip}>
                    <User size={14} />
                    {participant.name}
                  </span>
                ))}
              </div>
            </div>
            {conversation.itemId && (
              <div className={styles.itemContext}>
                <strong>Related Item:</strong>
                <span className={styles.itemLink}>
                  {conversation.itemType}: {conversation.itemTitle}
                </span>
              </div>
            )}
          </div>
          
          {isLoading ? (
            <div className={styles.messagesLoading}>
              <LoadingSpinner />
              <p>Loading messages...</p>
            </div>
          ) : (
            <div className={styles.messagesContainer}>
              {messages?.map((message, index) => (
                <div key={message.id} className={styles.messageItem}>
                  <div className={styles.messageHeader}>
                    <div className={styles.messageSender}>
                      <User size={14} />
                      <span>{message.sender?.name || 'Unknown User'}</span>
                    </div>
                    <div className={styles.messageTime}>
                      <Calendar size={12} />
                      <span>{new Date(message.createdAt).toLocaleString()}</span>
                    </div>
                  </div>
                  <div className={styles.messageContent}>
                    {message.content}
                  </div>
                  {message.attachments && message.attachments.length > 0 && (
                    <div className={styles.messageAttachments}>
                      <strong>Attachments:</strong>
                      {message.attachments.map(attachment => (
                        <span key={attachment.id} className={styles.attachment}>
                          {attachment.filename}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const MessagesManagement = () => {
  const t = useTranslations('MessagesManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    status: 'all',
    timeRange: '30d',
    search: ''
  });
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [showConversationModal, setShowConversationModal] = useState(false);

  // Fetch conversations data
  const { 
    data: conversationsData, 
    isLoading: conversationsLoading, 
    error: conversationsError,
    refetch: refetchConversations 
  } = useQuery({
    queryKey: ['adminConversations', filters],
    queryFn: () => listConversations(filters),
    enabled: isAdmin
  });

  // Fetch messaging statistics
  const { 
    data: statsData, 
    isLoading: statsLoading 
  } = useQuery({
    queryKey: ['messagingStats'],
    queryFn: () => getMessagingStats(),
    enabled: isAdmin
  });

  // Fetch analytics
  const { 
    data: analyticsData, 
    isLoading: analyticsLoading 
  } = useQuery({
    queryKey: ['messageAnalytics', filters.timeRange],
    queryFn: () => getMessageAnalytics({ timeRange: filters.timeRange }),
    enabled: isAdmin
  });

  // Archive mutation
  const archiveMutation = useMutation({
    mutationFn: archiveConversation,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminConversations']);
      queryClient.invalidateQueries(['messagingStats']);
    },
    onError: (error) => {
      
      alert('Failed to archive conversation. Please try again.');
    }
  });

  // Restore mutation
  const restoreMutation = useMutation({
    mutationFn: restoreConversation,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminConversations']);
      queryClient.invalidateQueries(['messagingStats']);
    },
    onError: (error) => {
      
      alert('Failed to restore conversation. Please try again.');
    }
  });

  // Mark as read mutation
  const markReadMutation = useMutation({
    mutationFn: markConversationAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminConversations']);
      queryClient.invalidateQueries(['messagingStats']);
    },
    onError: (error) => {
      
      alert('Failed to mark conversation as read. Please try again.');
    }
  });

  const handleConversationAction = useCallback((action, conversation) => {
    switch (action) {
      case 'view':
        setSelectedConversation(conversation);
        setShowConversationModal(true);
        break;
      case 'archive':
        if (confirm('Archive this conversation?')) {
          archiveMutation.mutate(conversation.id);
        }
        break;
      case 'restore':
        if (confirm('Restore this conversation?')) {
          restoreMutation.mutate(conversation.id);
        }
        break;
      case 'markRead':
        markReadMutation.mutate(conversation.id);
        break;
      case 'viewItem':
        if (conversation.itemId) {
          window.open(`/${conversation.itemType}s/${conversation.itemId}`, '_blank');
        }
        break;
      case 'flag':
        if (confirm('Flag this conversation for review?')) {
          // Implement flag functionality
          
        }
        break;
    }
  }, [archiveMutation, restoreMutation, markReadMutation]);

  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleExport = useCallback(() => {
    const csvData = conversationsData.map(c => ({
      'Conversation ID': c.id,
      'Participants': c.participants?.map(p => p.name).join(', ') || 'Unknown',
      'Messages': c.messageCount || 0,
      'Status': c.status,
      'Last Activity': new Date(c.lastActivity).toLocaleDateString(),
      'Item Type': c.itemType || 'N/A',
      'Item Title': c.itemTitle || 'N/A'
    }));

  }, [conversationsData]);

  // Process data
  const conversations = conversationsData?.conversations || [];
  const stats = statsData || {};
  const analytics = analyticsData || {};

  // Calculate summary stats
  const summaryStats = useMemo(() => {
    const active = conversations.filter(c => c.status === 'active');
    const archived = conversations.filter(c => c.status === 'archived');
    const unread = conversations.filter(c => c.unreadCount > 0);
    const totalMessages = conversations.reduce((sum, c) => sum + (c.messageCount || 0), 0);

    return {
      totalConversations: conversations.length,
      activeConversations: active.length,
      archivedConversations: archived.length,
      unreadConversations: unread.length,
      totalMessages,
      averageResponseTime: analytics.averageResponseTime || 0
    };
  }, [conversations, analytics]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access message management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (conversationsLoading && !conversationsData) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Messages...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch conversation data.' })}</p>
        </div>
      </div>
    );
  }

  if (conversationsError) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>{t('errorTitle', { defaultValue: 'Failed to Load Messages' })}</h2>
          <p>{conversationsError.message || t('errorMessage', { defaultValue: 'An error occurred while fetching message data' })}</p>
          <button className={styles.retryButton} onClick={() => refetchConversations()}>
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
            <h1 className={styles.title}>{t('title', { defaultValue: 'Messages Management' })}</h1>
            <p className={styles.subtitle}>{t('subtitle', { defaultValue: 'Monitor and manage user conversations' })}</p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.exportButton} onClick={handleExport}>
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button className={styles.refreshButton} onClick={() => refetchConversations()}>
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <MessageCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.totalConversations.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalConversations', { defaultValue: 'Total Conversations' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Activity size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.activeConversations.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('activeConversations', { defaultValue: 'Active Conversations' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <MessageSquare size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.totalMessages.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalMessages', { defaultValue: 'Total Messages' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Timer size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{summaryStats.averageResponseTime}h</div>
              <div className={styles.statLabel}>{t('averageResponseTime', { defaultValue: 'Avg Response Time' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersSection}>
          <div className={styles.searchContainer}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search conversations...' })}
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
              <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
              <option value="active">{t('active', { defaultValue: 'Active' })}</option>
              <option value="archived">{t('archived', { defaultValue: 'Archived' })}</option>
              <option value="flagged">{t('flagged', { defaultValue: 'Flagged' })}</option>
            </select>
            <select
              value={filters.timeRange}
              onChange={(e) => handleFilterChange('timeRange', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="7d">{t('last7Days', { defaultValue: 'Last 7 Days' })}</option>
              <option value="30d">{t('last30Days', { defaultValue: 'Last 30 Days' })}</option>
              <option value="90d">{t('last90Days', { defaultValue: 'Last 90 Days' })}</option>
              <option value="1y">{t('lastYear', { defaultValue: 'Last Year' })}</option>
            </select>
          </div>
        </div>

        {/* Conversations Table */}
        <div className={styles.tableSection}>
          <div className={styles.tableContainer}>
            <table className={styles.conversationsTable}>
              <thead>
                <tr>
                  <th>{t('conversation', { defaultValue: 'Conversation' })}</th>
                  <th>{t('participants', { defaultValue: 'Participants' })}</th>
                  <th>{t('messages', { defaultValue: 'Messages' })}</th>
                  <th>{t('status', { defaultValue: 'Status' })}</th>
                  <th>{t('lastActivity', { defaultValue: 'Last Activity' })}</th>
                  <th>Related Item</th>
                  <th>{t('actions', { defaultValue: 'Actions' })}</th>
                </tr>
              </thead>
              <tbody>
                {conversations.map((conversation) => (
                  <ConversationRow
                    key={conversation.id}
                    conversation={conversation}
                    onAction={handleConversationAction}
                  />
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Conversation Modal */}
        {showConversationModal && selectedConversation && (
          <ConversationModal
            conversation={selectedConversation}
            onClose={() => {
              setShowConversationModal(false);
              setSelectedConversation(null);
            }}
          />
        )}
      </div>
    </ErrorBoundary>
  );
};

export default MessagesManagement; 
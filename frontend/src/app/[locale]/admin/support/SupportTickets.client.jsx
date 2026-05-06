"use client";

import React, { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { 
  MessageSquare, 
  AlertCircle, 
  CheckCircle, 
  Clock, 
  Search, 
  Filter, 
  Plus, 
  Eye, 
  Edit, 
  MoreHorizontal,
  User,
  Calendar,
  Tag,
  AlertTriangle as PriorityHigh,
  AlertCircle as PriorityMedium,
  Info as PriorityLow,
  Info,
  Mail,
  Phone,
  MapPin,
  RefreshCw,
  Download,
  Archive,
  Star,
  StarOff,
  MessageCircle,
  Send,
  ArrowUpDown,
  Settings,
  Users,
  TrendingUp,
  BarChart3,
  ChevronUp,
  ChevronDown,
  X,
  Trash2
} from 'lucide-react';
import * as adminSupportApi from '../../../../api/client/admin/supportApi';
import * as supportApi from '../../../../api/client/supportApi';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'react-toastify';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './SupportTickets.module.css';

// Priority Badge Component
const PriorityBadge = ({ priority }) => {
  const priorityConfig = {
    low: { label: 'Low', icon: PriorityLow, color: 'success' },
    medium: { label: 'Medium', icon: PriorityMedium, color: 'warning' },
    high: { label: 'High', icon: PriorityHigh, color: 'danger' },
    urgent: { label: 'Urgent', icon: PriorityHigh, color: 'critical' }
  };

  const config = priorityConfig[priority] || priorityConfig.medium;
  const Icon = config.icon;

  return (
    <span className={`${styles.priorityBadge} ${styles[config.color]}`}>
      <Icon size={12} />
      {config.label}
    </span>
  );
};

// Status Badge Component
const StatusBadge = ({ status }) => {
  const statusConfig = {
    open: { label: 'Open', icon: AlertCircle, color: 'danger' },
    in_progress: { label: 'In Progress', icon: Clock, color: 'warning' },
    pending: { label: 'Pending Customer', icon: MessageSquare, color: 'info' },
    resolved: { label: 'Resolved', icon: CheckCircle, color: 'success' },
    closed: { label: 'Closed', icon: Archive, color: 'secondary' }
  };

  const config = statusConfig[status] || statusConfig.open;
  const Icon = config.icon;

  return (
    <span className={`${styles.statusBadge} ${styles[config.color]}`}>
      <Icon size={12} />
      {config.label}
    </span>
  );
};

// Ticket Detail Modal Component
const TicketDetailModal = ({ ticket, onClose, onUpdate }) => {
  const [activeTab, setActiveTab] = useState('conversation');
  const [replyContent, setReplyContent] = useState('');
  const [internalNote, setInternalNote] = useState('');
  const [isReplying, setIsReplying] = useState(false);

  const handleStatusUpdate = async (newStatus) => {
    try {
      await supportApi.updateTicket(ticket.id, { status: newStatus });
      onUpdate();
      toast.success(`Ticket status updated to ${newStatus}`);
    } catch (error) {
      toast.error('Failed to update status');
    }
  };

  const handlePriorityUpdate = async (newPriority) => {
    try {
      await supportApi.updateTicketPriority(ticket.id, newPriority, 'Admin priority adjustment');
      onUpdate();
      toast.success(`Priority updated to ${newPriority}`);
    } catch (error) {
      toast.error('Failed to update priority');
    }
  };

  const handleReply = async () => {
    if (!replyContent.trim()) return;
    
    setIsReplying(true);
    try {
      await supportApi.addTicketReply(ticket.id, {
        author_id: 'admin',
        author_type: 'AGENT',
        content: replyContent,
        is_public: true
      });
      setReplyContent('');
      toast.success('Reply sent successfully');
      onUpdate();
    } catch (error) {
      toast.error('Failed to send reply');
    } finally {
      setIsReplying(false);
    }
  };

  const handleInternalNote = async () => {
    if (!internalNote.trim()) return;
    
    try {
      await supportApi.addInternalNote(ticket.id, {
        author_id: 'admin',
        content: internalNote
      });
      setInternalNote('');
      toast.success('Internal note added');
    } catch (error) {
      toast.error('Failed to add note');
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modalContent} onClick={e => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <div className={styles.modalTitle}>
            <h2>Ticket #{ticket.id.slice(0, 8)}</h2>
            <StatusBadge status={ticket.status} />
            <PriorityBadge priority={ticket.priority} />
          </div>
          <button onClick={onClose} className={styles.modalClose}>
            <X size={20} />
          </button>
        </div>

        <div className={styles.modalBody}>
          <div className={styles.ticketInfo}>
            <h3>{ticket.title || ticket.subject}</h3>
            <p>{ticket.description}</p>
            
            <div className={styles.ticketMetadata}>
              <div className={styles.metaItem}>
                <User size={14} />
                <span>{ticket.created_by || 'Unknown User'}</span>
              </div>
              <div className={styles.metaItem}>
                <Calendar size={14} />
                <span>{new Date(ticket.created_at).toLocaleString()}</span>
              </div>
              {ticket.category && (
                <div className={styles.metaItem}>
                  <Tag size={14} />
                  <span>{ticket.category}</span>
                </div>
              )}
            </div>
          </div>

          <div className={styles.modalTabs}>
            <button 
              className={`${styles.tabButton} ${activeTab === 'conversation' ? styles.active : ''}`}
              onClick={() => setActiveTab('conversation')}
            >
              <MessageCircle size={16} />
              Conversation
            </button>
            <button 
              className={`${styles.tabButton} ${activeTab === 'details' ? styles.active : ''}`}
              onClick={() => setActiveTab('details')}
            >
              <Info size={16} />
              Details
            </button>
            <button 
              className={`${styles.tabButton} ${activeTab === 'actions' ? styles.active : ''}`}
              onClick={() => setActiveTab('actions')}
            >
              <Settings size={16} />
              Actions
            </button>
          </div>

          <div className={styles.tabContent}>
            {activeTab === 'conversation' && (
              <div className={styles.conversationTab}>
                <div className={styles.messageList}>
                  {/* Communications would be loaded here */}
                  <p className={styles.noMessages}>No messages yet</p>
                </div>
                
                <div className={styles.replySection}>
                  <textarea
                    value={replyContent}
                    onChange={(e) => setReplyContent(e.target.value)}
                    placeholder="Type your reply..."
                    className={styles.replyTextarea}
                  />
                  <div className={styles.replyActions}>
                    <button 
                      onClick={handleReply}
                      disabled={isReplying || !replyContent.trim()}
                      className={styles.sendButton}
                    >
                      <Send size={16} />
                      {isReplying ? 'Sending...' : 'Send Reply'}
                    </button>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'details' && (
              <div className={styles.detailsTab}>
                <div className={styles.detailGrid}>
                  <div className={styles.detailItem}>
                    <label>Status</label>
                    <select 
                      value={ticket.status}
                      onChange={(e) => handleStatusUpdate(e.target.value)}
                      className={styles.detailSelect}
                    >
                      <option value="SUBMITTED">Submitted</option>
                      <option value="ASSIGNED">Assigned</option>
                      <option value="IN_PROGRESS">In Progress</option>
                      <option value="PENDING_CUSTOMER">Pending Customer</option>
                      <option value="RESOLVED">Resolved</option>
                      <option value="CLOSED">Closed</option>
                    </select>
                  </div>
                  
                  <div className={styles.detailItem}>
                    <label>Priority</label>
                    <select 
                      value={ticket.priority}
                      onChange={(e) => handlePriorityUpdate(e.target.value)}
                      className={styles.detailSelect}
                    >
                      <option value="LOW">Low</option>
                      <option value="MEDIUM">Medium</option>
                      <option value="HIGH">High</option>
                      <option value="URGENT">Urgent</option>
                      <option value="CRITICAL">Critical</option>
                    </select>
                  </div>

                  <div className={styles.detailItem}>
                    <label>Category</label>
                    <span>{ticket.category || 'None'}</span>
                  </div>

                  <div className={styles.detailItem}>
                    <label>Channel ID</label>
                    <span className={styles.monoText}>{ticket.channel_id}</span>
                  </div>
                </div>

                <div className={styles.internalNotes}>
                  <h4>Internal Notes</h4>
                  <textarea
                    value={internalNote}
                    onChange={(e) => setInternalNote(e.target.value)}
                    placeholder="Add an internal note..."
                    className={styles.noteTextarea}
                  />
                  <button 
                    onClick={handleInternalNote}
                    disabled={!internalNote.trim()}
                    className={styles.noteButton}
                  >
                    Add Note
                  </button>
                </div>
              </div>
            )}

            {activeTab === 'actions' && (
              <div className={styles.actionsTab}>
                <button className={styles.actionButton}>
                  <Users size={16} />
                  Assign to Agent
                </button>
                <button className={styles.actionButton}>
                  <TrendingUp size={16} />
                  Escalate Ticket
                </button>
                <button className={styles.actionButton}>
                  <Archive size={16} />
                  Archive Ticket
                </button>
                <button className={`${styles.actionButton} ${styles.danger}`}>
                  <X size={16} />
                  Delete Ticket
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

// Ticket Card Component
const TicketCard = ({ ticket, onUpdate, onView }) => {
  const timeAgo = (date) => {
    const now = new Date();
    const diff = now - new Date(date);
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const days = Math.floor(hours / 24);
    
    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    return 'Recently';
  };

  return (
    <div className={`${styles.ticketCard} ${ticket.priority === 'urgent' ? styles.urgent : ''}`}>
      <div className={styles.ticketHeader}>
        <div className={styles.ticketMeta}>
          <span className={styles.ticketId}>#{ticket.id.slice(0, 8)}</span>
          <StatusBadge status={ticket.status?.toLowerCase() || 'open'} />
          <PriorityBadge priority={ticket.priority?.toLowerCase() || 'medium'} />
          {ticket.starred && <Star size={14} className={styles.starIcon} />}
        </div>
        <div className={styles.ticketActions}>
          <button 
            onClick={() => onView(ticket)}
            className={styles.actionButton}
            title="View Details"
          >
            <Eye size={14} />
          </button>
          <button className={styles.actionButton} title="More Actions">
            <MoreHorizontal size={14} />
          </button>
        </div>
      </div>

      <div className={styles.ticketContent}>
        <h3 className={styles.ticketSubject}>{ticket.title || ticket.subject || 'No subject'}</h3>
        <p className={styles.ticketDescription}>{ticket.description || 'No description'}</p>
        
        <div className={styles.ticketDetails}>
          <div className={styles.customerInfo}>
            <div className={styles.customerAvatar}>
              {ticket.customer?.avatar ? (
                <img src={ticket.customer.avatar} alt={ticket.customer.name} />
              ) : (
                <User size={16} />
              )}
            </div>
            <div className={styles.customerData}>
              <span className={styles.customerName}>{ticket.created_by || ticket.customer?.name || 'Anonymous'}</span>
              <span className={styles.customerEmail}>{ticket.customer?.email || ''}</span>
            </div>
          </div>
          
          <div className={styles.ticketTimeline}>
            <div className={styles.timelineItem}>
              <Calendar size={12} />
              <span>Created {timeAgo(ticket.created_at || ticket.createdAt)}</span>
            </div>
            {ticket.assigned_to && (
              <div className={styles.timelineItem}>
                <User size={12} />
                <span>Assigned to {ticket.assigned_to.name || ticket.assigned_to}</span>
              </div>
            )}
            {ticket.lastReply && (
              <div className={styles.timelineItem}>
                <MessageCircle size={12} />
                <span>Last reply {timeAgo(ticket.lastReply)}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className={styles.ticketFooter}>
        <div className={styles.ticketTags}>
          {ticket.category && (
            <span className={styles.categoryTag}>
              <Tag size={10} />
              {ticket.category}
            </span>
          )}
        </div>
        <div className={styles.ticketStats}>
          {ticket.responseTime && (
            <span className={styles.responseStat}>
              Response: {ticket.responseTime}
            </span>
          )}
          {ticket.replies > 0 && (
            <span className={styles.replyStat}>
              <MessageCircle size={12} />
              {ticket.replies}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

const SupportTickets = () => {
  const t = useTranslations('AdminSupport');
  const { user } = useAuth();
  const queryClient = useQueryClient();

  // State
  const [filter, setFilter] = useState('all');
  const [priority, setPriority] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedTickets, setSelectedTickets] = useState([]);
  const [viewTicket, setViewTicket] = useState(null);
  const [activeTab, setActiveTab] = useState('tickets');
  const [showAgentModal, setShowAgentModal] = useState(false);
  const [showKnowledgeModal, setShowKnowledgeModal] = useState(false);

  const itemsPerPage = 20;

  // Mock data for development
  const mockTickets = [
    {
      id: 'TICK-001',
      title: 'Unable to complete checkout',
      subject: 'Unable to complete checkout',
      description: 'I am getting an error when trying to complete my purchase. The payment page keeps timing out.',
      status: 'OPEN',
      priority: 'HIGH',
      category: 'Payment',
      created_by: 'John Doe',
      created_at: new Date().toISOString(),
      assigned_to: { name: 'Agent Smith' },
      channel_id: 'CH-001',
      customer: { name: 'John Doe', email: 'redacted-email@example.com' },
      replies: 2,
      responseTime: '2h'
    },
    {
      id: 'TICK-002',
      title: 'Product image not loading',
      subject: 'Product image not loading',
      description: 'The images for several products are not displaying on the product pages.',
      status: 'IN_PROGRESS',
      priority: 'MEDIUM',
      category: 'Technical',
      created_by: 'Jane Smith',
      created_at: new Date(Date.now() - 86400000).toISOString(),
      assigned_to: { name: 'Agent Jones' },
      channel_id: 'CH-002',
      customer: { name: 'Jane Smith', email: 'redacted-email@example.com' },
      replies: 5,
      responseTime: '30m'
    }
  ];

  const mockMetrics = {
    total_tickets: 156,
    open_tickets: 23,
    avg_response_time: '1.5h',
    resolution_rate: 87
  };

  const mockAgents = {
    agents: [
      {
        id: 'AG-001',
        name: 'Agent Smith',
        email: 'redacted-email@example.com',
        status: 'ONLINE',
        active_tickets: 5,
        resolved_tickets: 124,
        avg_response_time: '45m'
      },
      {
        id: 'AG-002',
        name: 'Agent Jones',
        email: 'redacted-email@example.com',
        status: 'AWAY',
        active_tickets: 3,
        resolved_tickets: 98,
        avg_response_time: '1.2h'
      }
    ]
  };

  const mockArticles = {
    articles: [
      {
        id: 'KB-001',
        title: 'How to reset your password',
        summary: 'Step-by-step guide to resetting your account password',
        category: 'Account',
        status: 'PUBLISHED',
        views: 1234,
        helpful_count: 89
      },
      {
        id: 'KB-002',
        title: 'Shipping and delivery FAQ',
        summary: 'Common questions about shipping times and delivery options',
        category: 'Orders',
        status: 'PUBLISHED',
        views: 2456,
        helpful_count: 156
      }
    ]
  };

  // Fetch tickets
  const { data: ticketsData, isLoading, error, refetch } = useQuery({
    queryKey: ['support-tickets', filter, priority, searchTerm, sortBy, sortOrder, currentPage],
    queryFn: () => adminSupportApi.listTickets({ 
      status: filter === 'all' ? undefined : filter.toUpperCase(),
      priority: priority === 'all' ? undefined : priority.toUpperCase(),
      search: searchTerm,
      sort_by: sortBy,
      sort_order: sortOrder,
      page: currentPage,
      limit: itemsPerPage
    }),
    staleTime: 30000, // 30 seconds
    retry: false,
    onError: (error) => {
      
    }
  });

  // Fetch metrics
  const { data: metrics } = useQuery({
    queryKey: ['support-metrics'],
    queryFn: () => adminSupportApi.getTicketMetrics(),
    staleTime: 300000, // 5 minutes
    retry: false
  });

  // Fetch agents
  const { data: agents } = useQuery({
    queryKey: ['support-agents'],
    queryFn: () => adminSupportApi.listAgents(),
    staleTime: 300000, // 5 minutes
    retry: false
  });

  // Fetch knowledge articles
  const { data: knowledgeArticles } = useQuery({
    queryKey: ['knowledge-articles'],
    queryFn: () => adminSupportApi.listKnowledgeArticles(),
    staleTime: 300000, // 5 minutes
    retry: false
  });

  // Mutations
  const updateTicketMutation = useMutation({
    mutationFn: ({ ticketId, updates }) => supportApi.updateTicket(ticketId, updates),
    onSuccess: () => {
      toast.success('Ticket updated successfully');
      queryClient.invalidateQueries(['support-tickets']);
    },
    onError: (error) => {
      toast.error(`Update failed: ${error.message}`);
    }
  });

  // Computed values - Use mock data if API fails
  const tickets = ticketsData?.tickets || (error ? mockTickets : []);
  const totalTickets = ticketsData?.total || tickets.length;
  const totalPages = Math.ceil(totalTickets / itemsPerPage);
  const displayMetrics = metrics || (error ? mockMetrics : null);
  const displayAgents = agents || mockAgents;
  const displayArticles = knowledgeArticles || mockArticles;

  const filteredTickets = useMemo(() => {
    return tickets.filter(ticket => {
      const matchesSearch = searchTerm === '' || 
        ticket.subject?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ticket.customer?.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ticket.id.toString().includes(searchTerm);
      return matchesSearch;
    });
  }, [tickets, searchTerm]);

  // Event handlers
  const handleStatusChange = (ticketId, newStatus) => {
    updateTicketMutation.mutate({
      ticketId,
      updates: { status: newStatus.toUpperCase() }
    });
  };

  const handleBulkAction = async (action) => {
    if (selectedTickets.length === 0) {
      toast.error('Please select tickets first');
      return;
    }

    try {
      await adminSupportApi.bulkUpdateTickets({
        ticket_ids: selectedTickets,
        updates: { status: action.toUpperCase() }
      });
      toast.success(`${selectedTickets.length} tickets updated successfully`);
      queryClient.invalidateQueries(['support-tickets']);
      setSelectedTickets([]);
    } catch (error) {
      toast.error(`Bulk update failed: ${error.message}`);
    }
  };

  const handleSelectTicket = (ticketId) => {
    setSelectedTickets(prev => 
      prev.includes(ticketId) 
        ? prev.filter(id => id !== ticketId)
        : [...prev, ticketId]
    );
  };

  const handleSelectAll = () => {
    if (selectedTickets.length === filteredTickets.length) {
      setSelectedTickets([]);
    } else {
      setSelectedTickets(filteredTickets.map(t => t.id));
    }
  };

  const handleExport = async () => {
    try {
      const blob = await adminSupportApi.exportReport('tickets', {
        status: filter === 'all' ? undefined : filter.toUpperCase(),
        priority: priority === 'all' ? undefined : priority.toUpperCase(),
        format: 'csv'
      });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `support-tickets-${new Date().toISOString().split('T')[0]}.csv`;
      a.click();
      window.URL.revokeObjectURL(url);
      toast.success('Export completed');
    } catch (error) {
      toast.error('Export failed');
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <p>Loading support tickets...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>Error Loading Tickets</h2>
          <p>{error.message}</p>
          <button onClick={() => refetch()} className={styles.retryButton}>
            <RefreshCw size={16} />
            Retry
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
          <div className={styles.headerLeft}>
            <h1 className={styles.title}>
              <MessageSquare size={24} />
              {t('title', { defaultValue: 'Support Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage support tickets, agents, and knowledge base' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            {activeTab === 'tickets' && (
              <button 
                onClick={handleExport}
                className={styles.exportButton}
              >
                <Download size={16} />
                Export
              </button>
            )}
            {activeTab === 'agents' && (
              <button 
                onClick={() => setShowAgentModal(true)}
                className={styles.createButton}
              >
                <Plus size={16} />
                Add Agent
              </button>
            )}
            {activeTab === 'knowledge' && (
              <button 
                onClick={() => setShowKnowledgeModal(true)}
                className={styles.createButton}
              >
                <Plus size={16} />
                New Article
              </button>
            )}
            <button 
              onClick={() => refetch()}
              className={styles.iconButton}
              disabled={isLoading}
            >
              <RefreshCw size={16} className={isLoading ? styles.spinning : ''} />
            </button>
            <button className={styles.settingsButton}>
              <Settings size={16} />
            </button>
          </div>
        </div>

        {/* Metrics Cards */}
        {displayMetrics && (
          <div className={styles.metricsGrid}>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <MessageSquare size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{displayMetrics.total_tickets || 0}</div>
                <div className={styles.metricLabel}>Total Tickets</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <AlertCircle size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{displayMetrics.open_tickets || 0}</div>
                <div className={styles.metricLabel}>Open</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <Clock size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{displayMetrics.avg_response_time || '0h'}</div>
                <div className={styles.metricLabel}>Avg Response</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <TrendingUp size={20} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{displayMetrics.resolution_rate || 0}%</div>
                <div className={styles.metricLabel}>Resolution Rate</div>
              </div>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className={styles.tabs}>
          <button 
            className={`${styles.tab} ${activeTab === 'tickets' ? styles.active : ''}`}
            onClick={() => setActiveTab('tickets')}
          >
            <MessageSquare size={16} />
            Tickets
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'agents' ? styles.active : ''}`}
            onClick={() => setActiveTab('agents')}
          >
            <Users size={16} />
            Agents
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'knowledge' ? styles.active : ''}`}
            onClick={() => setActiveTab('knowledge')}
          >
            <Info size={16} />
            Knowledge Base
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'analytics' ? styles.active : ''}`}
            onClick={() => setActiveTab('analytics')}
          >
            <BarChart3 size={16} />
            Analytics
          </button>
        </div>

        {/* Content based on active tab */}
        {activeTab === 'tickets' && (
          <>
            {/* Filters and Search */}
            <div className={styles.filtersSection}>
          <div className={styles.searchBar}>
            <Search size={16} />
            <input
              type="text"
              placeholder="Search tickets, customers, or IDs..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filters}>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">All Status</option>
              <option value="open">Open</option>
              <option value="in_progress">In Progress</option>
              <option value="pending">Pending</option>
              <option value="resolved">Resolved</option>
              <option value="closed">Closed</option>
            </select>
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">All Priority</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="urgent">Urgent</option>
            </select>
            <select
              value={`${sortBy}-${sortOrder}`}
              onChange={(e) => {
                const [field, order] = e.target.value.split('-');
                setSortBy(field);
                setSortOrder(order);
              }}
              className={styles.filterSelect}
            >
              <option value="created_at-desc">Newest First</option>
              <option value="created_at-asc">Oldest First</option>
              <option value="priority-desc">Priority High-Low</option>
              <option value="updated_at-desc">Recently Updated</option>
              <option value="subject-asc">Subject A-Z</option>
            </select>
          </div>
        </div>

        {/* Bulk Actions */}
        {selectedTickets.length > 0 && (
          <div className={styles.bulkActions}>
            <div className={styles.bulkInfo}>
              <span>{selectedTickets.length} tickets selected</span>
            </div>
            <div className={styles.bulkButtons}>
              <button 
                onClick={() => handleBulkAction('in_progress')}
                className={styles.bulkButton}
              >
                Mark In Progress
              </button>
              <button 
                onClick={() => handleBulkAction('resolved')}
                className={styles.bulkButton}
              >
                Mark Resolved
              </button>
              <button 
                onClick={() => handleBulkAction('closed')}
                className={`${styles.bulkButton} ${styles.secondary}`}
              >
                Close Tickets
              </button>
            </div>
          </div>
        )}

        {/* Tickets Grid */}
        <div className={styles.ticketsContainer}>
          {filteredTickets.length > 0 ? (
            <div className={styles.ticketsGrid}>
              {filteredTickets.map((ticket) => (
                <TicketCard
                  key={ticket.id}
                  ticket={ticket}
                  onUpdate={handleStatusChange}
                  onView={setViewTicket}
                />
              ))}
            </div>
          ) : (
            <div className={styles.emptyState}>
              <MessageSquare size={48} />
              <h3>No Support Tickets Found</h3>
              <p>
                {searchTerm || filter !== 'all' || priority !== 'all'
                  ? 'Try adjusting your search or filters.'
                  : 'No support tickets have been submitted yet.'
                }
              </p>
            </div>
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className={styles.pagination}>
            <button
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
              className={styles.paginationButton}
            >
              Previous
            </button>
            <span className={styles.pageInfo}>
              Page {currentPage} of {totalPages}
            </span>
            <button
              onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
              disabled={currentPage === totalPages}
              className={styles.paginationButton}
            >
              Next
            </button>
          </div>
        )}

        {/* Ticket Detail Modal */}
        {viewTicket && (
          <TicketDetailModal
            ticket={viewTicket}
            onClose={() => setViewTicket(null)}
            onUpdate={() => {
              refetch();
              setViewTicket(null);
            }}
          />
        )}
          </>
        )}

        {/* Agents Tab */}
        {activeTab === 'agents' && (
          <div className={styles.agentsSection}>
            <div className={styles.agentsGrid}>
              {displayAgents?.agents?.map(agent => (
                <div key={agent.id} className={styles.agentCard}>
                  <div className={styles.agentHeader}>
                    <div className={styles.agentInfo}>
                      <h3>{agent.name}</h3>
                      <span className={styles.agentEmail}>{agent.email}</span>
                    </div>
                    <span className={`${styles.statusBadge} ${styles[agent.status?.toLowerCase() || 'offline']}`}>
                      {agent.status || 'Offline'}
                    </span>
                  </div>
                  <div className={styles.agentStats}>
                    <div className={styles.stat}>
                      <span className={styles.statValue}>{agent.active_tickets || 0}</span>
                      <span className={styles.statLabel}>Active</span>
                    </div>
                    <div className={styles.stat}>
                      <span className={styles.statValue}>{agent.resolved_tickets || 0}</span>
                      <span className={styles.statLabel}>Resolved</span>
                    </div>
                    <div className={styles.stat}>
                      <span className={styles.statValue}>{agent.avg_response_time || '0h'}</span>
                      <span className={styles.statLabel}>Avg Response</span>
                    </div>
                  </div>
                  <div className={styles.agentActions}>
                    <button className={styles.actionButton}>
                      <Edit size={14} />
                    </button>
                    <button className={styles.actionButton}>
                      <MoreHorizontal size={14} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Knowledge Base Tab */}
        {activeTab === 'knowledge' && (
          <div className={styles.knowledgeSection}>
            <div className={styles.articlesList}>
              {displayArticles?.articles?.map(article => (
                <div key={article.id} className={styles.articleItem}>
                  <div className={styles.articleContent}>
                    <h3>{article.title}</h3>
                    <p>{article.summary}</p>
                    <div className={styles.articleMeta}>
                      <span>{article.category}</span>
                      <span>{article.views || 0} views</span>
                      <span>{article.helpful_count || 0} helpful</span>
                    </div>
                  </div>
                  <div className={styles.articleActions}>
                    <span className={`${styles.statusBadge} ${styles[article.status?.toLowerCase() || 'draft']}`}>
                      {article.status || 'Draft'}
                    </span>
                    <button className={styles.actionButton}>
                      <Edit size={14} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Analytics Tab */}
        {activeTab === 'analytics' && (
          <div className={styles.analyticsSection}>
            <div className={styles.analyticsGrid}>
              <div className={styles.analyticsCard}>
                <h3>Ticket Trends</h3>
                <p className={styles.placeholder}>Chart placeholder - Ticket volume over time</p>
              </div>
              <div className={styles.analyticsCard}>
                <h3>Response Times</h3>
                <p className={styles.placeholder}>Chart placeholder - Average response times</p>
              </div>
              <div className={styles.analyticsCard}>
                <h3>Agent Performance</h3>
                <p className={styles.placeholder}>Chart placeholder - Agent metrics</p>
              </div>
              <div className={styles.analyticsCard}>
                <h3>Customer Satisfaction</h3>
                <p className={styles.placeholder}>Chart placeholder - CSAT scores</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default SupportTickets;
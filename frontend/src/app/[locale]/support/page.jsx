"use client";
export const dynamic = 'force-dynamic';
import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useTranslations } from "next-intl";
import { 
    Plus,
    Search,
    Filter,
    RefreshCw,
    MessageSquare,
    Clock,
    CheckCircle,
    AlertCircle,
    Archive,
    ChevronDown,
    Star,
    HelpCircle,
    BookOpen,
    Zap,
    TrendingUp,
    X
} from 'lucide-react';
import { useSupport } from '../../../hooks/useSupport';
import { useIsMobile } from '../../../hooks/useMobileDetection';
import TicketDetail from '../../../components/Support/TicketDetail';
import SupportTicketForm from '../../../components/Support/SupportTicketForm';
import styles from './Support.module.css';
const TicketCard = ({ ticket, onClick }) => {
    const getPriorityColor = (priority) => {
        const colors = {
            LOW: 'low',
            MEDIUM: 'medium',
            HIGH: 'high',
            URGENT: 'urgent',
            CRITICAL: 'critical'
        };
        return colors[priority] || 'medium';
    };

    const getStatusIcon = (status) => {
        const icons = {
            SUBMITTED: Clock,
            IN_PROGRESS: TrendingUp,
            RESOLVED: CheckCircle,
            CLOSED: Archive,
            PENDING_CUSTOMER: MessageSquare
        };
        const Icon = icons[status] || AlertCircle;
        return <Icon size={14} />;
    };

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
        <div className={styles.ticketCard} onClick={() => onClick(ticket)}>
            <div className={styles.ticketHeader}>
                <div className={styles.ticketStatus}>
                    {getStatusIcon(ticket.status)}
                    <span>{ticket.status.replace('_', ' ')}</span>
                </div>
                <span className={`${styles.ticketPriority} ${styles[getPriorityColor(ticket.priority)]}`}>
                    {ticket.priority}
                </span>
            </div>
            <h3 className={styles.ticketTitle}>{ticket.title}</h3>
            <p className={styles.ticketDescription}>{ticket.description}</p>
            <div className={styles.ticketFooter}>
                <span className={styles.ticketId}>#{ticket.id.slice(0, 8)}</span>
                <span className={styles.ticketTime}>{timeAgo(ticket.created_at)}</span>
                {ticket.response_count > 0 && (
                    <span className={styles.ticketReplies}>
                        <MessageSquare size={12} />
                        {ticket.response_count}
                    </span>
                )}
            </div>
        </div>
    );
};

const Support = () => {
    const t = useTranslations('Support');
    const isMobile = useIsMobile();
    const {
        activeChannelId,
        tickets,
        isLoading,
        initializeSupport,
        createTicket,
        refetchTickets,
        isCreatingTicket,
        searchKnowledgeBase
    } = useSupport();

    const [showForm, setShowForm] = useState(false);
    const [selectedTicket, setSelectedTicket] = useState(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [statusFilter, setStatusFilter] = useState('all');
    const [priorityFilter, setPriorityFilter] = useState('all');
    const [sortBy, setSortBy] = useState('created_at');
    const [showKnowledge, setShowKnowledge] = useState(false);
    const [knowledgeArticles, setKnowledgeArticles] = useState([]);

    // Initialize support on mount
    useEffect(() => {
        initializeSupport();
    }, [initializeSupport]);

    // Filter and sort tickets
    const filteredTickets = useMemo(() => {
        let filtered = [...tickets];

        // Search filter
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(ticket =>
                ticket.title?.toLowerCase().includes(query) ||
                ticket.description?.toLowerCase().includes(query) ||
                ticket.id?.toLowerCase().includes(query)
            );
        }

        // Status filter
        if (statusFilter !== 'all') {
            const statusMap = {
                open: ['SUBMITTED', 'ASSIGNED', 'IN_PROGRESS', 'REOPENED'],
                pending: ['PENDING_CUSTOMER'],
                resolved: ['RESOLVED'],
                closed: ['CLOSED']
            };
            filtered = filtered.filter(ticket => 
                statusMap[statusFilter]?.includes(ticket.status)
            );
        }

        // Priority filter
        if (priorityFilter !== 'all') {
            filtered = filtered.filter(ticket => 
                ticket.priority === priorityFilter.toUpperCase()
            );
        }

        // Sort
        filtered.sort((a, b) => {
            switch (sortBy) {
                case 'created_at':
                    return new Date(b.created_at) - new Date(a.created_at);
                case 'updated_at':
                    return new Date(b.updated_at) - new Date(a.updated_at);
                case 'priority':
                    const priorityOrder = { CRITICAL: 0, URGENT: 1, HIGH: 2, MEDIUM: 3, LOW: 4 };
                    return priorityOrder[a.priority] - priorityOrder[b.priority];
                default:
                    return 0;
            }
        });

        return filtered;
    }, [tickets, searchQuery, statusFilter, priorityFilter, sortBy]);

    // Statistics
    const stats = useMemo(() => {
        const openCount = tickets.filter(t => 
            ['SUBMITTED', 'ASSIGNED', 'IN_PROGRESS', 'REOPENED'].includes(t.status)
        ).length;
        const pendingCount = tickets.filter(t => t.status === 'PENDING_CUSTOMER').length;
        const resolvedCount = tickets.filter(t => 
            ['RESOLVED', 'CLOSED'].includes(t.status)
        ).length;

        return {
            total: tickets.length,
            open: openCount,
            pending: pendingCount,
            resolved: resolvedCount
        };
    }, [tickets]);

    const handleCreateTicket = async (ticketData) => {
        await createTicket(ticketData);
        setShowForm(false);
    };

    const handleSearchKnowledge = async () => {
        if (!searchQuery) return;
        
        try {
            const results = await searchKnowledgeBase(searchQuery);
            setKnowledgeArticles(results.articles || []);
            setShowKnowledge(true);
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    };
    return (
        <div className={styles.container}>
            {/* Header */}
            <div className={styles.header}>
                <div className={styles.headerContent}>
                    <div className={styles.titleSection}>
                        <h1 className={styles.title}>
                            <MessageSquare size={20} />
                            Support Center
                        </h1>
                        <p className={styles.subtitle}>
                            Get help with your account, orders, and more
                        </p>
                    </div>
                    <div className={styles.headerActions}>
                        <button 
                            onClick={handleSearchKnowledge}
                            className={styles.knowledgeButton}
                            disabled={!searchQuery}
                        >
                            <BookOpen size={16} />
                            Search Help
                        </button>
                        <button 
                            onClick={() => refetchTickets()}
                            className={styles.refreshButton}
                        >
                            <RefreshCw size={16} />
                        </button>
                        <button 
                            onClick={() => setShowForm(true)}
                            className={styles.createButton}
                        >
                            <Plus size={16} />
                            New Ticket
                        </button>
                    </div>
                </div>

                {/* Stats */}
                <div className={styles.statsBar}>
                    <div className={styles.statCard}>
                        <div className={styles.statValue}>{stats.total}</div>
                        <div className={styles.statLabel}>Total Tickets</div>
                    </div>
                    <div className={`${styles.statCard} ${styles.open}`}>
                        <div className={styles.statValue}>{stats.open}</div>
                        <div className={styles.statLabel}>Open</div>
                    </div>
                    <div className={`${styles.statCard} ${styles.pending}`}>
                        <div className={styles.statValue}>{stats.pending}</div>
                        <div className={styles.statLabel}>Pending</div>
                    </div>
                    <div className={`${styles.statCard} ${styles.resolved}`}>
                        <div className={styles.statValue}>{stats.resolved}</div>
                        <div className={styles.statLabel}>Resolved</div>
                    </div>
                </div>
            </div>

            {/* Filters */}
            <div className={styles.filters}>
                <div className={styles.searchBar}>
                    <Search size={16} />
                    <input
                        type="text"
                        placeholder="Search tickets..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className={styles.searchInput}
                    />
                </div>
                <div className={styles.filterGroup}>
                    <select
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">All Status</option>
                        <option value="open">Open</option>
                        <option value="pending">Pending</option>
                        <option value="resolved">Resolved</option>
                        <option value="closed">Closed</option>
                    </select>
                    <select
                        value={priorityFilter}
                        onChange={(e) => setPriorityFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">All Priority</option>
                        <option value="low">Low</option>
                        <option value="medium">Medium</option>
                        <option value="high">High</option>
                        <option value="urgent">Urgent</option>
                        <option value="critical">Critical</option>
                    </select>
                    <select
                        value={sortBy}
                        onChange={(e) => setSortBy(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="created_at">Newest First</option>
                        <option value="updated_at">Recently Updated</option>
                        <option value="priority">Priority</option>
                    </select>
                </div>
            </div>

            {/* Quick Actions */}
            <div className={styles.quickActions}>
                <button className={styles.quickAction}>
                    <HelpCircle size={16} />
                    FAQs
                </button>
                <button className={styles.quickAction}>
                    <Zap size={16} />
                    Quick Solutions
                </button>
                <button className={styles.quickAction}>
                    <Star size={16} />
                    Popular Articles
                </button>
            </div>

            {/* Main Content */}
            <div className={styles.content}>
                {isLoading ? (
                    <div className={styles.loading}>
                        <RefreshCw size={24} className={styles.spinner} />
                        <p>Loading your tickets...</p>
                    </div>
                ) : filteredTickets.length === 0 ? (
                    <div className={styles.emptyState}>
                        <MessageSquare size={48} />
                        <h3>No tickets found</h3>
                        <p>
                            {searchQuery || statusFilter !== 'all' || priorityFilter !== 'all'
                                ? 'Try adjusting your filters'
                                : 'Create your first support ticket to get help'}
                        </p>
                        {(!searchQuery && statusFilter === 'all' && priorityFilter === 'all') && (
                            <button 
                                onClick={() => setShowForm(true)}
                                className={styles.emptyStateButton}
                            >
                                <Plus size={16} />
                                Create Ticket
                            </button>
                        )}
                    </div>
                ) : (
                    <div className={styles.ticketsGrid}>
                        {filteredTickets.map(ticket => (
                            <TicketCard
                                key={ticket.id}
                                ticket={ticket}
                                onClick={setSelectedTicket}
                            />
                        ))}
                    </div>
                )}
            </div>

            {/* Modals */}
            {showForm && (
                <div className={styles.modalOverlay} onClick={() => setShowForm(false)}>
                    <div className={styles.modalContent} onClick={e => e.stopPropagation()}>
                        <SupportTicketForm
                            onCreate={handleCreateTicket}
                            onCancel={() => setShowForm(false)}
                            isCreating={isCreatingTicket}
                        />
                    </div>
                </div>
            )}

            {selectedTicket && (
                <TicketDetail
                    ticket={selectedTicket}
                    onClose={() => setSelectedTicket(null)}
                    onUpdate={() => {
                        refetchTickets();
                        setSelectedTicket(null);
                    }}
                />
            )}

            {/* Knowledge Base Results */}
            {showKnowledge && knowledgeArticles.length > 0 && (
                <div className={styles.knowledgeOverlay} onClick={() => setShowKnowledge(false)}>
                    <div className={styles.knowledgeContent} onClick={e => e.stopPropagation()}>
                        <div className={styles.knowledgeHeader}>
                            <h3>Help Articles</h3>
                            <button onClick={() => setShowKnowledge(false)}>
                                <X size={20} />
                            </button>
                        </div>
                        <div className={styles.knowledgeList}>
                            {knowledgeArticles.map(article => (
                                <div key={article.id} className={styles.knowledgeItem}>
                                    <h4>{article.title}</h4>
                                    <p>{article.content.substring(0, 150)}...</p>
                                    <div className={styles.articleMeta}>
                                        <span>{article.view_count} views</span>
                                        {article.average_rating > 0 && (
                                            <span>⭐ {article.average_rating.toFixed(1)}</span>
                                        )}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Support;
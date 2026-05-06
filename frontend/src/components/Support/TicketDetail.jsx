import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import dayjs from 'dayjs';
import {
    X,
    Send,
    Paperclip,
    Clock,
    User,
    MessageSquare,
    AlertCircle,
    CheckCircle,
    Star,
    MoreVertical,
    Reply,
    Archive,
    RefreshCw,
    Tag,
    ChevronDown,
    ChevronUp,
    File,
    Image,
    FileText,
    Download,
    AlertTriangle,
    Info
} from 'lucide-react';
import { useSupport } from '../../hooks/useSupport';
import styles from './TicketDetail.module.css';

const PriorityIcon = ({ priority }) => {
    const icons = {
        LOW: Info,
        MEDIUM: AlertCircle,
        HIGH: AlertTriangle,
        URGENT: AlertTriangle,
        CRITICAL: AlertTriangle
    };
    const Icon = icons[priority] || Info;
    return <Icon size={14} className={styles[`priority${priority}`]} />;
};

const StatusBadge = ({ status }) => {
    const statusConfig = {
        DRAFT: { label: 'Draft', color: 'secondary' },
        SUBMITTED: { label: 'Submitted', color: 'info' },
        ASSIGNED: { label: 'Assigned', color: 'primary' },
        IN_PROGRESS: { label: 'In Progress', color: 'warning' },
        PENDING_CUSTOMER: { label: 'Pending Response', color: 'warning' },
        RESOLVED: { label: 'Resolved', color: 'success' },
        CLOSED: { label: 'Closed', color: 'secondary' },
        ESCALATED: { label: 'Escalated', color: 'danger' },
        REOPENED: { label: 'Reopened', color: 'warning' }
    };

    const config = statusConfig[status] || statusConfig.SUBMITTED;
    
    return (
        <span className={`${styles.statusBadge} ${styles[config.color]}`}>
            {config.label}
        </span>
    );
};

const AttachmentItem = ({ attachment, onDownload }) => {
    const getFileIcon = (contentType) => {
        if (contentType.startsWith('image/')) return Image;
        if (contentType.includes('pdf')) return FileText;
        return File;
    };

    const FileIcon = getFileIcon(attachment.content_type);

    return (
        <div className={styles.attachmentItem}>
            <FileIcon size={16} />
            <span className={styles.attachmentName}>{attachment.filename}</span>
            <span className={styles.attachmentSize}>
                {(attachment.size_bytes / 1024).toFixed(1)}KB
            </span>
            <button 
                onClick={() => onDownload(attachment)}
                className={styles.downloadButton}
                title="Download"
            >
                <Download size={14} />
            </button>
        </div>
    );
};

const CommunicationItem = ({ communication, isUser }) => {
    const authorTypeIcons = {
        CUSTOMER: User,
        AGENT: User,
        AI: MessageSquare,
        SYSTEM: Info
    };

    const AuthorIcon = authorTypeIcons[communication.author_type] || User;

    return (
        <div className={`${styles.communicationItem} ${isUser ? styles.userMessage : styles.agentMessage}`}>
            <div className={styles.messageHeader}>
                <div className={styles.authorInfo}>
                    <div className={styles.authorAvatar}>
                        <AuthorIcon size={16} />
                    </div>
                    <span className={styles.authorName}>
                        {communication.author_type === 'CUSTOMER' ? 'You' : 
                         communication.author_type === 'AI' ? 'AI Assistant' :
                         communication.author_type === 'SYSTEM' ? 'System' :
                         'Support Agent'}
                    </span>
                </div>
                <span className={styles.messageTime}>
                    {dayjs(communication.created_at).format('MMM D, h:mm A')}
                </span>
            </div>
            <div className={styles.messageContent}>
                {communication.content}
            </div>
            {communication.attachments?.length > 0 && (
                <div className={styles.messageAttachments}>
                    {communication.attachments.map(attachment => (
                        <AttachmentItem 
                            key={attachment.id} 
                            attachment={attachment}
                            onDownload={() => window.open(attachment.url, '_blank')}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

const TicketDetail = ({ ticket, onClose, onUpdate }) => {
    const t = useTranslations('Support');
    const { 
        addReply, 
        resolveTicket, 
        closeTicket, 
        reopenTicket,
        getTicketCommunications,
        isAddingReply,
        isResolvingTicket,
        isClosingTicket,
        isReopeningTicket
    } = useSupport();

    const [communications, setCommunications] = useState([]);
    const [loading, setLoading] = useState(true);
    const [replyContent, setReplyContent] = useState('');
    const [showActions, setShowActions] = useState(false);
    const [showResolutionForm, setShowResolutionForm] = useState(false);
    const [resolution, setResolution] = useState('');
    const [satisfaction, setSatisfaction] = useState(null);
    const [attachments, setAttachments] = useState([]);

    // Load communications
    useEffect(() => {
        const loadCommunications = async () => {
            if (!ticket?.id) return;
            
            setLoading(true);
            try {
                const result = await getTicketCommunications(ticket.id);
                setCommunications(result.communications || []);
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
                setLoading(false);
            }
        };

        loadCommunications();
    }, [ticket?.id, getTicketCommunications]);

    const handleSendReply = useCallback(async () => {
        if (!replyContent.trim()) return;

        await addReply({
            ticketId: ticket.id,
            content: replyContent,
            attachments: attachments.map(a => ({ id: a.id })),
            is_public: true
        });

        setReplyContent('');
        setAttachments([]);
        
        // Reload communications
        const result = await getTicketCommunications(ticket.id);
        setCommunications(result.communications || []);
    }, [ticket.id, replyContent, attachments, addReply, getTicketCommunications]);

    const handleResolve = useCallback(async () => {
        if (!resolution.trim()) return;

        await resolveTicket({
            ticketId: ticket.id,
            resolution,
            applied_solutions: []
        });

        setShowResolutionForm(false);
        setResolution('');
        onUpdate();
    }, [ticket.id, resolution, resolveTicket, onUpdate]);

    const handleClose = useCallback(async () => {
        await closeTicket({
            ticketId: ticket.id,
            closure_notes: '',
            satisfaction
        });

        onUpdate();
        onClose();
    }, [ticket.id, satisfaction, closeTicket, onUpdate, onClose]);

    const handleReopen = useCallback(async () => {
        const reason = window.prompt('Please provide a reason for reopening this ticket:');
        if (!reason) return;

        await reopenTicket({
            ticketId: ticket.id,
            reason
        });

        onUpdate();
    }, [ticket.id, reopenTicket, onUpdate]);

    const isTicketOpen = ['SUBMITTED', 'ASSIGNED', 'IN_PROGRESS', 'PENDING_CUSTOMER', 'REOPENED'].includes(ticket.status);
    const isTicketResolved = ['RESOLVED', 'CLOSED'].includes(ticket.status);

    return (
        <div className={styles.ticketDetailModal}>
            <div className={styles.modalOverlay} onClick={onClose} />
            <div className={styles.modalContent}>
                {/* Header */}
                <div className={styles.modalHeader}>
                    <div className={styles.headerLeft}>
                        <h2 className={styles.ticketTitle}>{ticket.title}</h2>
                        <div className={styles.ticketMeta}>
                            <span className={styles.ticketId}>#{ticket.id}</span>
                            <StatusBadge status={ticket.status} />
                            <div className={styles.priorityBadge}>
                                <PriorityIcon priority={ticket.priority} />
                                <span>{ticket.priority}</span>
                            </div>
                            {ticket.category && (
                                <span className={styles.categoryTag}>
                                    <Tag size={12} />
                                    {ticket.category}
                                </span>
                            )}
                        </div>
                    </div>
                    <div className={styles.headerActions}>
                        <button 
                            className={styles.actionButton}
                            onClick={() => setShowActions(!showActions)}
                        >
                            <MoreVertical size={16} />
                        </button>
                        <button className={styles.closeButton} onClick={onClose}>
                            <X size={20} />
                        </button>
                    </div>
                </div>

                {/* Actions Dropdown */}
                {showActions && (
                    <div className={styles.actionsDropdown}>
                        {isTicketOpen && (
                            <>
                                <button 
                                    onClick={() => setShowResolutionForm(true)}
                                    className={styles.dropdownItem}
                                >
                                    <CheckCircle size={14} />
                                    Mark as Resolved
                                </button>
                                <button 
                                    onClick={handleClose}
                                    className={styles.dropdownItem}
                                >
                                    <Archive size={14} />
                                    Close Ticket
                                </button>
                            </>
                        )}
                        {isTicketResolved && (
                            <button 
                                onClick={handleReopen}
                                className={styles.dropdownItem}
                            >
                                <RefreshCw size={14} />
                                Reopen Ticket
                            </button>
                        )}
                    </div>
                )}

                {/* Ticket Info */}
                <div className={styles.ticketInfo}>
                    <div className={styles.infoSection}>
                        <h3>Description</h3>
                        <p>{ticket.description}</p>
                    </div>
                    <div className={styles.infoSidebar}>
                        <div className={styles.infoItem}>
                            <Clock size={14} />
                            <span>Created {dayjs(ticket.created_at).format('MMM D, YYYY')}</span>
                        </div>
                        {ticket.updated_at && (
                            <div className={styles.infoItem}>
                                <RefreshCw size={14} />
                                <span>Updated {dayjs(ticket.updated_at).format('MMM D, YYYY')}</span>
                            </div>
                        )}
                        {ticket.response_count > 0 && (
                            <div className={styles.infoItem}>
                                <MessageSquare size={14} />
                                <span>{ticket.response_count} responses</span>
                            </div>
                        )}
                    </div>
                </div>

                {/* Resolution Form */}
                {showResolutionForm && (
                    <div className={styles.resolutionForm}>
                        <h3>Resolve Ticket</h3>
                        <textarea
                            value={resolution}
                            onChange={(e) => setResolution(e.target.value)}
                            placeholder="Describe how this issue was resolved..."
                            className={styles.resolutionTextarea}
                            rows={4}
                        />
                        <div className={styles.resolutionActions}>
                            <button 
                                onClick={() => setShowResolutionForm(false)}
                                className={styles.cancelButton}
                            >
                                Cancel
                            </button>
                            <button 
                                onClick={handleResolve}
                                disabled={!resolution.trim() || isResolvingTicket}
                                className={styles.primaryButton}
                            >
                                {isResolvingTicket ? 'Resolving...' : 'Resolve Ticket'}
                            </button>
                        </div>
                    </div>
                )}

                {/* Communications */}
                <div className={styles.communicationsSection}>
                    <h3>Conversation</h3>
                    <div className={styles.communicationsList}>
                        {loading ? (
                            <div className={styles.loadingState}>Loading conversation...</div>
                        ) : communications.length === 0 ? (
                            <div className={styles.emptyState}>No messages yet</div>
                        ) : (
                            communications.map((comm) => (
                                <CommunicationItem
                                    key={comm.id}
                                    communication={comm}
                                    isUser={comm.author_type === 'CUSTOMER'}
                                />
                            ))
                        )}
                    </div>
                </div>

                {/* Reply Form */}
                {isTicketOpen && (
                    <div className={styles.replySection}>
                        <div className={styles.replyForm}>
                            <textarea
                                value={replyContent}
                                onChange={(e) => setReplyContent(e.target.value)}
                                placeholder="Type your message..."
                                className={styles.replyTextarea}
                                rows={3}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter' && e.ctrlKey) {
                                        handleSendReply();
                                    }
                                }}
                            />
                            <div className={styles.replyActions}>
                                <button className={styles.attachButton}>
                                    <Paperclip size={16} />
                                    Attach
                                </button>
                                <button 
                                    onClick={handleSendReply}
                                    disabled={!replyContent.trim() || isAddingReply}
                                    className={styles.sendButton}
                                >
                                    <Send size={16} />
                                    {isAddingReply ? 'Sending...' : 'Send'}
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default TicketDetail;
"use client";

import React, { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { 
  Mail, 
  Users, 
  Send, 
  Calendar, 
  Plus, 
  Search, 
  Filter, 
  Eye, 
  Edit, 
  Trash2, 
  MoreHorizontal,
  RefreshCw,
  Download,
  Upload,
  CheckCircle,
  Clock,
  AlertTriangle,
  TrendingUp,
  BarChart3,
  Target,
  Globe,
  FileText,
  Settings,
  Copy,
  Archive,
  Star,
  Zap,
  Image,
  Type,
  Layout,
  Palette,
  MousePointer,
  PieChart,
  Activity
} from 'lucide-react';
import { 
  fetchNewsletters,
  createNewsletter,
  updateNewsletter,
  deleteNewsletter,
  listSubscriptions,
  listEditions,
  createEdition,
  updateEdition,
  sendEdition,
  scheduleEdition,
  listTemplates,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  getNewsletterStats,
  getEditionStats
} from '../../../../api/client/newsletterApi';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'react-toastify';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import TemplateEditor from '@/components/NewsletterEditor/TemplateEditor';
import styles from './NewsletterManagement.module.css';

// Campaign Status Badge Component
const CampaignStatusBadge = ({ status }) => {
  const statusConfig = {
    draft: { label: 'Draft', icon: FileText, color: 'secondary' },
    scheduled: { label: 'Scheduled', icon: Clock, color: 'warning' },
    sending: { label: 'Sending', icon: Zap, color: 'info' },
    sent: { label: 'Sent', icon: CheckCircle, color: 'success' },
    failed: { label: 'Failed', icon: AlertTriangle, color: 'danger' }
  };

  const config = statusConfig[status] || statusConfig.draft;
  const Icon = config.icon;

  return (
    <span className={`${styles.statusBadge} ${styles[config.color]}`}>
      <Icon size={12} />
      {config.label}
    </span>
  );
};

// Quick Action Button Component
const QuickActionButton = ({ icon: Icon, label, onClick, variant = 'default', disabled = false }) => {
  return (
    <button
      className={`${styles.actionButton} ${variant === 'primary' ? styles.primaryButton : ''} ${variant === 'danger' ? styles.danger : ''}`}
      onClick={onClick}
      disabled={disabled}
      title={label}
    >
      <Icon size={16} />
      <span className={styles.buttonText}>{label}</span>
    </button>
  );
};

// Subscriber Stats Component
const SubscriberStats = ({ subscribers }) => {
  const stats = useMemo(() => {
    if (!subscribers) return { active: 0, unsubscribed: 0, bounced: 0 };
    
    return {
      active: subscribers.filter(s => s.status === 'active').length,
      unsubscribed: subscribers.filter(s => s.status === 'unsubscribed').length,
      bounced: subscribers.filter(s => s.status === 'bounced').length
    };
  }, [subscribers]);

  return (
    <div className={styles.subscriberStats}>
      <div className={styles.statItem}>
        <span className={styles.statValue}>{stats.active}</span>
        <span className={styles.statLabel}>Active</span>
      </div>
      <div className={styles.statItem}>
        <span className={styles.statValue}>{stats.unsubscribed}</span>
        <span className={styles.statLabel}>Unsubscribed</span>
      </div>
      <div className={styles.statItem}>
        <span className={styles.statValue}>{stats.bounced}</span>
        <span className={styles.statLabel}>Bounced</span>
      </div>
    </div>
  );
};

// Newsletter Card Component
const NewsletterCard = ({ newsletter, onSelect, onEdit, onDelete, onDuplicate, onToggleActive, isSelected, onCheckboxChange, isChecked }) => {
  const formatDate = (date) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const handleCheckboxClick = (e) => {
    e.stopPropagation();
    onCheckboxChange(newsletter.id);
  };

  return (
    <div className={`${styles.campaignCard} ${isSelected ? styles.selected : ''}`} onClick={() => onSelect(newsletter)}>
      <div className={styles.campaignHeader}>
        <div className={styles.campaignMeta}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <input
              type="checkbox"
              checked={isChecked}
              onChange={handleCheckboxClick}
              onClick={(e) => e.stopPropagation()}
              className={styles.checkbox}
            />
            <h3 className={styles.campaignTitle}>{newsletter.name}</h3>
          </div>
          <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
            <span className={`${styles.statusBadge} ${newsletter.is_active ? styles.success : styles.secondary}`}>
              {newsletter.is_active ? 'Active' : 'Inactive'}
            </span>
            {newsletter.subscriber_count > 100 && (
              <span className={`${styles.statusBadge} ${styles.info}`}>
                <TrendingUp size={12} />
                Popular
              </span>
            )}
          </div>
        </div>
        <div className={styles.campaignActions}>
          <button 
            onClick={(e) => { e.stopPropagation(); onToggleActive(newsletter); }}
            className={styles.actionButton}
            title={newsletter.is_active ? "Deactivate" : "Activate"}
          >
            {newsletter.is_active ? <Archive size={14} /> : <Zap size={14} />}
          </button>
          <button 
            onClick={(e) => { e.stopPropagation(); onDuplicate(newsletter); }}
            className={styles.actionButton}
            title="Duplicate Newsletter"
          >
            <Copy size={14} />
          </button>
          <button 
            onClick={(e) => { e.stopPropagation(); onEdit(newsletter); }}
            className={styles.actionButton}
            title="Edit Newsletter"
          >
            <Edit size={14} />
          </button>
          <button 
            onClick={(e) => { e.stopPropagation(); onDelete(newsletter.id); }}
            className={`${styles.actionButton} ${styles.danger}`}
            title="Delete Newsletter"
          >
            <Trash2 size={14} />
          </button>
        </div>
      </div>

      <div className={styles.campaignContent}>
        <p className={styles.campaignPreview}>{newsletter.description}</p>
        
        <div className={styles.campaignDetails}>
          <div className={styles.campaignTimeline}>
            <div className={styles.timelineItem}>
              <Calendar size={12} />
              <span>{newsletter.frequency}</span>
            </div>
            {newsletter.category && (
              <div className={styles.timelineItem}>
                <Globe size={12} />
                <span>{newsletter.category}</span>
              </div>
            )}
            <div className={styles.timelineItem}>
              <Clock size={12} />
              <span>Created {formatDate(newsletter.created_at)}</span>
            </div>
          </div>
        </div>
      </div>

      <div className={styles.campaignMetrics}>
        <div className={styles.metric}>
          <span className={styles.metricValue}>{newsletter.subscriber_count || 0}</span>
          <span className={styles.metricLabel}>Subscribers</span>
        </div>
        <div className={styles.metric}>
          <span className={styles.metricValue}>{newsletter.editions_count || 0}</span>
          <span className={styles.metricLabel}>Editions</span>
        </div>
        <div className={styles.metric}>
          <span className={styles.metricValue}>{newsletter.open_rate || 0}%</span>
          <span className={styles.metricLabel}>Open Rate</span>
        </div>
        <div className={styles.metric}>
          <span className={styles.metricValue}>{newsletter.click_rate || 0}%</span>
          <span className={styles.metricLabel}>Click Rate</span>
        </div>
      </div>
    </div>
  );
};

// Create Newsletter Modal Component
const CreateNewsletterModal = ({ isOpen, onClose, onSubmit, initialData = null, isEdit = false }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    frequency: 'weekly',
    category: '',
    templateId: '',
    is_active: true
  });

  React.useEffect(() => {
    if (initialData && isEdit) {
      setFormData({
        name: initialData.name || '',
        description: initialData.description || '',
        frequency: initialData.frequency || 'weekly',
        category: initialData.category || '',
        templateId: initialData.template_id || '',
        is_active: initialData.is_active !== undefined ? initialData.is_active : true
      });
    }
  }, [initialData, isEdit]);

  if (!isOpen) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
    if (!isEdit) {
      setFormData({
        name: '',
        description: '',
        frequency: 'weekly',
        category: '',
        templateId: '',
        is_active: true
      });
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h2>{isEdit ? 'Edit Newsletter' : 'Create New Newsletter'}</h2>
          <button onClick={onClose} className={styles.closeButton}>×</button>
        </div>
        
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label>Newsletter Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              placeholder="Enter newsletter name..."
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label>Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              placeholder="Describe your newsletter..."
              rows={3}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label>Frequency</label>
            <select
              value={formData.frequency}
              onChange={(e) => setFormData(prev => ({ ...prev, frequency: e.target.value }))}
            >
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
            </select>
          </div>

          <div className={styles.formGroup}>
            <label>Category</label>
            <input
              type="text"
              value={formData.category}
              onChange={(e) => setFormData(prev => ({ ...prev, category: e.target.value }))}
              placeholder="e.g., Technology, Business, Health"
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label className={styles.checkboxLabel}>
              <input
                type="checkbox"
                checked={formData.is_active}
                onChange={(e) => setFormData(prev => ({ ...prev, is_active: e.target.checked }))}
              />
              Active newsletter
            </label>
          </div>

          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.submitButton}>
              {isEdit ? 'Update Newsletter' : 'Create Newsletter'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Create Edition Modal Component
const CreateEditionModal = ({ isOpen, onClose, onSubmit, templates }) => {
  const [formData, setFormData] = useState({
    subject: '',
    contentHtml: '',
    contentText: '',
    templateData: {},
    scheduledAt: ''
  });

  if (!isOpen) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
    setFormData({
      subject: '',
      contentHtml: '',
      contentText: '',
      templateData: {},
      scheduledAt: ''
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h2>Create New Edition</h2>
          <button onClick={onClose} className={styles.closeButton}>×</button>
        </div>
        
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label>Subject Line</label>
            <input
              type="text"
              value={formData.subject}
              onChange={(e) => setFormData(prev => ({ ...prev, subject: e.target.value }))}
              placeholder="Enter email subject..."
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label>HTML Content</label>
            <textarea
              value={formData.contentHtml}
              onChange={(e) => setFormData(prev => ({ ...prev, contentHtml: e.target.value }))}
              placeholder="HTML email content..."
              rows={6}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label>Plain Text Content</label>
            <textarea
              value={formData.contentText}
              onChange={(e) => setFormData(prev => ({ ...prev, contentText: e.target.value }))}
              placeholder="Plain text version..."
              rows={4}
            />
          </div>

          <div className={styles.formGroup}>
            <label>Schedule For (Optional)</label>
            <input
              type="datetime-local"
              value={formData.scheduledAt}
              onChange={(e) => setFormData(prev => ({ ...prev, scheduledAt: e.target.value }))}
            />
          </div>

          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.submitButton}>
              Create Edition
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const NewsletterManagement = () => {
  const t = useTranslations('AdminNewsletters');
  const { user } = useAuth();
  const queryClient = useQueryClient();

  // State
  const [activeTab, setActiveTab] = useState('newsletters');
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [filterCategory, setFilterCategory] = useState('all');
  const [currentPage, setCurrentPage] = useState(1);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditionModal, setShowEditionModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedNewsletter, setSelectedNewsletter] = useState(null);
  const [editingNewsletter, setEditingNewsletter] = useState(null);
  const [selectedNewsletters, setSelectedNewsletters] = useState(new Set());
  const [bulkActionLoading, setBulkActionLoading] = useState(false);
  const [showTemplateEditor, setShowTemplateEditor] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState(null);

  const itemsPerPage = 12;

  // Fetch data
  const { data: newsletters, isLoading: newslettersLoading } = useQuery({
    queryKey: ['newsletters', searchTerm, currentPage],
    queryFn: () => fetchNewsletters({
      search: searchTerm,
      page: currentPage,
      limit: itemsPerPage
    }),
    staleTime: 60000 // 1 minute
  });

  const { data: subscriptions, isLoading: subscriptionsLoading } = useQuery({
    queryKey: ['newsletter-subscriptions'],
    queryFn: () => listSubscriptions({}),
    staleTime: 300000 // 5 minutes
  });

  const { data: editions, isLoading: editionsLoading } = useQuery({
    queryKey: ['newsletter-editions', selectedNewsletter?.id],
    queryFn: () => selectedNewsletter ? listEditions(selectedNewsletter.id) : null,
    enabled: !!selectedNewsletter,
    staleTime: 60000 // 1 minute
  });

  const { data: templates } = useQuery({
    queryKey: ['newsletter-templates'],
    queryFn: () => listTemplates({}),
    staleTime: 600000 // 10 minutes
  });

  const { data: stats } = useQuery({
    queryKey: ['newsletter-stats', selectedNewsletter?.id],
    queryFn: () => selectedNewsletter ? getNewsletterStats(selectedNewsletter.id) : null,
    enabled: !!selectedNewsletter,
    staleTime: 300000 // 5 minutes
  });

  // Mutations
  const createNewsletterMutation = useMutation({
    mutationFn: createNewsletter,
    onSuccess: () => {
      toast.success('Newsletter created successfully');
      queryClient.invalidateQueries(['newsletters']);
      setShowCreateModal(false);
    },
    onError: (error) => {
      toast.error(`Failed to create newsletter: ${error.message}`);
    }
  });

  const updateNewsletterMutation = useMutation({
    mutationFn: ({ id, data }) => updateNewsletter(id, data),
    onSuccess: () => {
      toast.success('Newsletter updated successfully');
      queryClient.invalidateQueries(['newsletters']);
      setShowEditModal(false);
      setEditingNewsletter(null);
    },
    onError: (error) => {
      toast.error(`Failed to update newsletter: ${error.message}`);
    }
  });

  const deleteNewsletterMutation = useMutation({
    mutationFn: deleteNewsletter,
    onSuccess: () => {
      toast.success('Newsletter deleted successfully');
      queryClient.invalidateQueries(['newsletters']);
    },
    onError: (error) => {
      toast.error(`Failed to delete newsletter: ${error.message}`);
    }
  });

  const createEditionMutation = useMutation({
    mutationFn: ({ newsletterId, editionData }) => createEdition(newsletterId, editionData),
    onSuccess: () => {
      toast.success('Edition created successfully');
      queryClient.invalidateQueries(['newsletter-editions']);
      setShowEditionModal(false);
    },
    onError: (error) => {
      toast.error(`Failed to create edition: ${error.message}`);
    }
  });

  const sendEditionMutation = useMutation({
    mutationFn: ({ editionId, testMode }) => sendEdition(editionId, testMode),
    onSuccess: () => {
      toast.success('Edition sent successfully');
      queryClient.invalidateQueries(['newsletter-editions']);
    },
    onError: (error) => {
      toast.error(`Failed to send edition: ${error.message}`);
    }
  });

  const saveTemplateMutation = useMutation({
    mutationFn: (templateData) => {
      if (editingTemplate?.id) {
        return updateTemplate(editingTemplate.id, templateData);
      }
      return createTemplate(templateData);
    },
    onSuccess: () => {
      toast.success(editingTemplate ? 'Template updated successfully' : 'Template created successfully');
      queryClient.invalidateQueries(['newsletter-templates']);
      setShowTemplateEditor(false);
      setEditingTemplate(null);
    },
    onError: (error) => {
      toast.error(`Failed to save template: ${error.message}`);
    }
  });

  const deleteTemplateMutation = useMutation({
    mutationFn: deleteTemplate,
    onSuccess: () => {
      toast.success('Template deleted successfully');
      queryClient.invalidateQueries(['newsletter-templates']);
    },
    onError: (error) => {
      toast.error(`Failed to delete template: ${error.message}`);
    }
  });

  // Event handlers
  const handleCreateNewsletter = (formData) => {
    createNewsletterMutation.mutate(formData);
  };

  const handleUpdateNewsletter = (formData) => {
    if (!editingNewsletter) return;
    updateNewsletterMutation.mutate({
      id: editingNewsletter.id,
      data: formData
    });
  };

  const handleDeleteNewsletter = (newsletterId) => {
    if (confirm('Are you sure you want to delete this newsletter?')) {
      deleteNewsletterMutation.mutate(newsletterId);
      setSelectedNewsletters(prev => {
        const newSet = new Set(prev);
        newSet.delete(newsletterId);
        return newSet;
      });
    }
  };

  const handleDuplicateNewsletter = async (newsletter) => {
    try {
      const duplicatedData = {
        name: `${newsletter.name} (Copy)`,
        description: newsletter.description,
        frequency: newsletter.frequency,
        category: newsletter.category,
        template_id: newsletter.template_id,
        is_active: false
      };
      await createNewsletterMutation.mutateAsync(duplicatedData);
      toast.success('Newsletter duplicated successfully');
    } catch (error) {
      toast.error(`Failed to duplicate newsletter: ${error.message}`);
    }
  };

  const handleToggleActive = async (newsletter) => {
    try {
      await updateNewsletterMutation.mutateAsync({
        id: newsletter.id,
        data: { is_active: !newsletter.is_active }
      });
      toast.success(`Newsletter ${newsletter.is_active ? 'deactivated' : 'activated'} successfully`);
    } catch (error) {
      toast.error(`Failed to update newsletter status: ${error.message}`);
    }
  };

  const handleBulkDelete = async () => {
    if (!selectedNewsletters.size) return;
    if (!confirm(`Are you sure you want to delete ${selectedNewsletters.size} newsletters?`)) return;
    
    setBulkActionLoading(true);
    try {
      await Promise.all(
        Array.from(selectedNewsletters).map(id => deleteNewsletter(id))
      );
      toast.success(`${selectedNewsletters.size} newsletters deleted successfully`);
      setSelectedNewsletters(new Set());
    } catch (error) {
      toast.error(`Failed to delete some newsletters: ${error.message}`);
    } finally {
      setBulkActionLoading(false);
    }
  };

  const handleBulkActivate = async (activate = true) => {
    if (!selectedNewsletters.size) return;
    
    setBulkActionLoading(true);
    try {
      await Promise.all(
        Array.from(selectedNewsletters).map(id => {
          const newsletter = newslettersList.find(n => n.id === id);
          if (newsletter) {
            return updateNewsletter(id, { is_active: activate });
          }
        })
      );
      toast.success(`${selectedNewsletters.size} newsletters ${activate ? 'activated' : 'deactivated'} successfully`);
      queryClient.invalidateQueries(['newsletters']);
      setSelectedNewsletters(new Set());
    } catch (error) {
      toast.error(`Failed to update some newsletters: ${error.message}`);
    } finally {
      setBulkActionLoading(false);
    }
  };

  const handleSelectAll = () => {
    if (selectedNewsletters.size === filteredNewsletters.length) {
      setSelectedNewsletters(new Set());
    } else {
      setSelectedNewsletters(new Set(filteredNewsletters.map(n => n.id)));
    }
  };

  const handleCreateEdition = (formData) => {
    if (!selectedNewsletter) return;
    createEditionMutation.mutate({
      newsletterId: selectedNewsletter.id,
      editionData: formData
    });
  };

  const handleSendEdition = (editionId, testMode = false) => {
    sendEditionMutation.mutate({ editionId, testMode });
  };

  const handleSaveTemplate = (templateData) => {
    saveTemplateMutation.mutate(templateData);
  };

  const handleDeleteTemplate = (templateId) => {
    if (confirm('Are you sure you want to delete this template?')) {
      deleteTemplateMutation.mutate(templateId);
    }
  };

  const handlePreviewTemplate = (template) => {
    // For now, just show the template in a new window
    const previewWindow = window.open('', '_blank');
    previewWindow.document.write(template.html_template || '<p>No preview available</p>');
    previewWindow.document.close();
  };

  const handleExportSubscribers = async () => {
    try {
      // TODO: Implement CSV export from subscriptions data
      const csvContent = 'email,status,subscribed_at\n';
      const blob = new Blob([csvContent], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `subscribers-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      toast.success('Subscribers exported successfully');
    } catch (error) {
      toast.error(`Export failed: ${error.message}`);
    }
  };

  // Computed values
  const newslettersList = newsletters?.newsletters || [];
  const totalNewsletters = newsletters?.total || 0;
  const totalPages = Math.ceil(totalNewsletters / itemsPerPage);

  const filteredNewsletters = useMemo(() => {
    return newslettersList.filter(newsletter => {
      const matchesSearch = searchTerm === '' || 
        newsletter.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        newsletter.description?.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesStatus = filterStatus === 'all' || 
        (filterStatus === 'active' ? newsletter.is_active : !newsletter.is_active);
      const matchesCategory = filterCategory === 'all' || 
        newsletter.category?.toLowerCase() === filterCategory.toLowerCase();
      return matchesSearch && matchesStatus && matchesCategory;
    });
  }, [newslettersList, searchTerm, filterStatus, filterCategory]);

  const uniqueCategories = useMemo(() => {
    const categories = new Set(
      newslettersList
        .map(n => n.category)
        .filter(Boolean)
    );
    return Array.from(categories).sort();
  }, [newslettersList]);

  // Calculate metrics
  const metrics = useMemo(() => {
    const totalSubscribers = subscriptions?.total || 0;
    const activeSubscribers = subscriptions?.subscriptions?.filter(s => s.status === 'active').length || 0;
    const totalEditionsSent = editions?.editions?.filter(e => e.status === 'sent').length || 0;
    
    return {
      total_subscribers: totalSubscribers,
      active_subscribers: activeSubscribers,
      editions_sent: totalEditionsSent,
      avg_open_rate: stats?.average_open_rate || 0,
      avg_click_rate: stats?.average_click_rate || 0
    };
  }, [subscriptions, editions, stats]);

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <h1 className={styles.title}>
              <Mail size={24} />
              {t('title', { defaultValue: 'Newsletter Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage subscribers, create campaigns, and track newsletter performance' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            <button 
              onClick={() => setShowCreateModal(true)}
              className={styles.createButton}
            >
              <Plus size={18} />
              Create Newsletter
            </button>
            <button 
              onClick={handleExportSubscribers}
              className={styles.exportButton}
              title="Export Subscribers"
            >
              <Download size={16} />
              Export
            </button>
            <button 
              onClick={() => window.location.reload()}
              className={styles.importButton}
              title="Refresh"
            >
              <RefreshCw size={16} />
              Refresh
            </button>
          </div>
        </div>

        {/* Metrics Cards */}
        <div className={styles.metricsGrid}>
          <div className={styles.metricCard}>
            <div className={styles.metricIcon}>
              <Users size={20} />
            </div>
            <div className={styles.metricContent}>
              <div className={styles.metricValue}>{metrics.total_subscribers}</div>
              <div className={styles.metricLabel}>Total Subscribers</div>
            </div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricIcon}>
              <CheckCircle size={20} />
            </div>
            <div className={styles.metricContent}>
              <div className={styles.metricValue}>{metrics.active_subscribers}</div>
              <div className={styles.metricLabel}>Active Subscribers</div>
            </div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricIcon}>
              <Mail size={20} />
            </div>
            <div className={styles.metricContent}>
              <div className={styles.metricValue}>{metrics.editions_sent}</div>
              <div className={styles.metricLabel}>Editions Sent</div>
            </div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricIcon}>
              <MousePointer size={20} />
            </div>
            <div className={styles.metricContent}>
              <div className={styles.metricValue}>{metrics.avg_open_rate.toFixed(1)}%</div>
              <div className={styles.metricLabel}>Avg Open Rate</div>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className={styles.tabs}>
          <button 
            className={`${styles.tab} ${activeTab === 'newsletters' ? styles.active : ''}`}
            onClick={() => setActiveTab('newsletters')}
          >
            <Mail size={20} />
            Newsletters
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'editions' ? styles.active : ''}`}
            onClick={() => setActiveTab('editions')}
            disabled={!selectedNewsletter}
          >
            <FileText size={20} />
            Editions
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'subscribers' ? styles.active : ''}`}
            onClick={() => setActiveTab('subscribers')}
          >
            <Users size={20} />
            Subscribers
          </button>
          <button 
            className={`${styles.tab} ${activeTab === 'templates' ? styles.active : ''}`}
            onClick={() => setActiveTab('templates')}
          >
            <Layout size={20} />
            Templates
          </button>
        </div>

        {/* Newsletters Tab */}
        {activeTab === 'newsletters' && (
          <div className={styles.campaignsSection}>
            {/* Filters */}
            <div className={styles.filtersSection}>
              <div className={styles.searchBar}>
                <Search size={16} />
                <input
                  type="text"
                  placeholder="Search newsletters..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className={styles.searchInput}
                />
              </div>
              <div className={styles.filters}>
                <select
                  value={filterStatus}
                  onChange={(e) => setFilterStatus(e.target.value)}
                  className={styles.filterSelect}
                >
                  <option value="all">All Status</option>
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                </select>
                <select
                  value={filterCategory}
                  onChange={(e) => setFilterCategory(e.target.value)}
                  className={styles.filterSelect}
                >
                  <option value="all">All Categories</option>
                  {uniqueCategories.map(cat => (
                    <option key={cat} value={cat}>{cat}</option>
                  ))}
                </select>
                {selectedNewsletters.size > 0 && (
                  <div style={{ display: 'flex', gap: '8px', marginLeft: '16px' }}>
                    <button
                      onClick={() => handleBulkActivate(true)}
                      className={styles.actionButton}
                      disabled={bulkActionLoading}
                    >
                      <CheckCircle size={14} />
                      Activate ({selectedNewsletters.size})
                    </button>
                    <button
                      onClick={() => handleBulkActivate(false)}
                      className={styles.actionButton}
                      disabled={bulkActionLoading}
                    >
                      <Archive size={14} />
                      Deactivate ({selectedNewsletters.size})
                    </button>
                    <button
                      onClick={handleBulkDelete}
                      className={`${styles.actionButton} ${styles.danger}`}
                      disabled={bulkActionLoading}
                    >
                      <Trash2 size={14} />
                      Delete ({selectedNewsletters.size})
                    </button>
                  </div>
                )}
              </div>
            </div>

            {/* Newsletters Grid */}
            {newslettersLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p>Loading newsletters...</p>
              </div>
            ) : filteredNewsletters.length > 0 ? (
              <div className={styles.campaignsGrid}>
                {filteredNewsletters.map((newsletter) => (
                  <NewsletterCard
                    key={newsletter.id}
                    newsletter={newsletter}
                    isSelected={selectedNewsletter?.id === newsletter.id}
                    isChecked={selectedNewsletters.has(newsletter.id)}
                    onSelect={setSelectedNewsletter}
                    onCheckboxChange={(id) => {
                      setSelectedNewsletters(prev => {
                        const newSet = new Set(prev);
                        if (newSet.has(id)) {
                          newSet.delete(id);
                        } else {
                          newSet.add(id);
                        }
                        return newSet;
                      });
                    }}
                    onEdit={(newsletter) => {
                      setEditingNewsletter(newsletter);
                      setShowEditModal(true);
                    }}
                    onDelete={handleDeleteNewsletter}
                    onDuplicate={handleDuplicateNewsletter}
                    onToggleActive={handleToggleActive}
                  />
                ))}
              </div>
            ) : (
              <div className={styles.emptyState}>
                <Mail size={48} />
                <h3>No Newsletters Found</h3>
                <p>
                  {searchTerm || filterStatus !== 'all'
                    ? 'Try adjusting your search or filters.'
                    : 'Create your first newsletter to get started.'
                  }
                </p>
                <button 
                  onClick={() => setShowCreateModal(true)}
                  className={styles.createButton}
                >
                  <Plus size={18} />
                  Create Your First Newsletter
                </button>
              </div>
            )}

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
          </div>
        )}

        {/* Editions Tab */}
        {activeTab === 'editions' && selectedNewsletter && (
          <div className={styles.editionsSection}>
            <div className={styles.editionsHeader}>
              <h2>Editions for {selectedNewsletter.name}</h2>
              <button 
                onClick={() => setShowEditionModal(true)}
                className={styles.createButton}
              >
                <Plus size={18} />
                Create Edition
              </button>
            </div>

            {editionsLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p>Loading editions...</p>
              </div>
            ) : editions?.editions?.length > 0 ? (
              <div className={styles.editionsList}>
                {editions.editions.map((edition) => (
                  <div key={edition.id} className={styles.editionItem}>
                    <div className={styles.editionHeader}>
                      <h4>{edition.subject}</h4>
                      <CampaignStatusBadge status={edition.status} />
                    </div>
                    <div className={styles.editionMeta}>
                      <span>Created: {new Date(edition.created_at).toLocaleDateString()}</span>
                      {edition.scheduled_at && (
                        <span>Scheduled: {new Date(edition.scheduled_at).toLocaleDateString()}</span>
                      )}
                      {edition.sent_at && (
                        <span>Sent: {new Date(edition.sent_at).toLocaleDateString()}</span>
                      )}
                    </div>
                    <div className={styles.editionActions}>
                      {edition.status === 'draft' && (
                        <>
                          <button 
                            onClick={() => handleSendEdition(edition.id, true)}
                            className={styles.actionButton}
                          >
                            Test Send
                          </button>
                          <button 
                            onClick={() => handleSendEdition(edition.id, false)}
                            className={styles.createButton}
                          >
                            <Send size={16} />
                            Send Now
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className={styles.emptyState}>
                <FileText size={48} />
                <h3>No Editions Found</h3>
                <p>Create your first edition to send to subscribers.</p>
                <button 
                  onClick={() => setShowEditionModal(true)}
                  className={styles.createButton}
                >
                  <Plus size={18} />
                  Create Edition
                </button>
              </div>
            )}
          </div>
        )}

        {/* Subscribers Tab */}
        {activeTab === 'subscribers' && (
          <div className={styles.subscribersSection}>
            <div className={styles.subscribersHeader}>
              <h2>Newsletter Subscribers</h2>
              <div className={styles.subscribersActions}>
                <button className={styles.importButton}>
                  <Upload size={16} />
                  Import CSV
                </button>
                <button onClick={handleExportSubscribers} className={styles.exportButton}>
                  <Download size={16} />
                  Export CSV
                </button>
              </div>
            </div>

            {subscriptionsLoading ? (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p>Loading subscribers...</p>
              </div>
            ) : (
              <>
                <SubscriberStats subscribers={subscriptions?.subscriptions} />
                
                {subscriptions?.subscriptions?.length > 0 ? (
                  <div className={styles.subscribersTable}>
                    <table className={styles.table}>
                      <thead>
                        <tr>
                          <th>User</th>
                          <th>Newsletter</th>
                          <th>Status</th>
                          <th>Subscribed</th>
                          <th>Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {subscriptions.subscriptions.map((sub) => (
                          <tr key={sub.id}>
                            <td>{sub.user_id}</td>
                            <td>{sub.newsletter?.name || sub.newsletter_id}</td>
                            <td>
                              <span className={`${styles.statusBadge} ${styles[sub.status]}`}>
                                {sub.status}
                              </span>
                            </td>
                            <td>{new Date(sub.subscribed_at).toLocaleDateString()}</td>
                            <td>
                              <button 
                                className={styles.actionButton}
                                onClick={() => {/* View subscription details */}}
                              >
                                <Eye size={14} />
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className={styles.emptyState}>
                    <Users size={48} />
                    <h3>No Subscribers Yet</h3>
                    <p>Subscribers will appear here once users subscribe to newsletters.</p>
                  </div>
                )}
              </>
            )}
          </div>
        )}

        {/* Templates Tab */}
        {activeTab === 'templates' && (
          <div className={styles.templatesSection}>
            {showTemplateEditor ? (
              <TemplateEditor
                initialTemplate={editingTemplate}
                onSave={handleSaveTemplate}
                onCancel={() => {
                  setShowTemplateEditor(false);
                  setEditingTemplate(null);
                }}
                isLoading={saveTemplateMutation.isLoading}
              />
            ) : (
              <>
                <div className={styles.templatesHeader}>
                  <h2>Email Templates</h2>
                  <button 
                    onClick={() => setShowTemplateEditor(true)}
                    className={styles.createButton}
                  >
                    <Plus size={18} />
                    Create Template
                  </button>
                </div>
                
                {templates?.templates?.length > 0 ? (
              <div className={styles.templatesGrid}>
                {templates.templates.map((template) => (
                  <div key={template.id} className={styles.templateCard}>
                    <div className={styles.templateHeader}>
                      <h3>{template.name}</h3>
                      {template.is_public && (
                        <span className={styles.publicBadge}>Public</span>
                      )}
                    </div>
                    <p className={styles.templateDescription}>
                      {template.description}
                    </p>
                    <div className={styles.templateMeta}>
                      <span>Created: {new Date(template.created_at).toLocaleDateString()}</span>
                      {template.user_id && <span>By: {template.user_id}</span>}
                    </div>
                    <div className={styles.templateActions}>
                      <button 
                        className={styles.actionButton}
                        onClick={() => handlePreviewTemplate(template)}
                      >
                        <Eye size={14} /> Preview
                      </button>
                      <button 
                        className={styles.actionButton}
                        onClick={() => {
                          setEditingTemplate(template);
                          setShowTemplateEditor(true);
                        }}
                      >
                        <Edit size={14} /> Edit
                      </button>
                      <button 
                        className={`${styles.actionButton} ${styles.danger}`}
                        onClick={() => handleDeleteTemplate(template.id)}
                      >
                        <Trash2 size={14} /> Delete
                      </button>
                    </div>
                  </div>
                ))}
                </div>
              ) : (
                <div className={styles.emptyState}>
                  <Layout size={48} />
                  <h3>No Templates Found</h3>
                  <p>Create templates to reuse in your newsletter editions.</p>
                  <button 
                    onClick={() => setShowTemplateEditor(true)}
                    className={styles.createButton}
                  >
                    <Plus size={18} />
                    Create Your First Template
                  </button>
                </div>
              )}
              </>
            )}
          </div>
        )}

        {/* Analytics Tab */}
        {activeTab === 'analytics' && (
          <div className={styles.analyticsSection}>
            <h2>Newsletter Analytics</h2>
            <p className={styles.placeholder}>
              Detailed analytics dashboard would be displayed here.
            </p>
          </div>
        )}

        {/* Create Newsletter Modal */}
        <CreateNewsletterModal
          isOpen={showCreateModal}
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateNewsletter}
        />

        {/* Create Edition Modal */}
        <CreateEditionModal
          isOpen={showEditionModal}
          onClose={() => setShowEditionModal(false)}
          onSubmit={handleCreateEdition}
          templates={templates?.templates}
        />

        {/* Edit Newsletter Modal */}
        <CreateNewsletterModal
          isOpen={showEditModal}
          onClose={() => {
            setShowEditModal(false);
            setEditingNewsletter(null);
          }}
          onSubmit={handleUpdateNewsletter}
          initialData={editingNewsletter}
          isEdit={true}
        />
      </div>
    </ErrorBoundary>
  );
};

export default NewsletterManagement;
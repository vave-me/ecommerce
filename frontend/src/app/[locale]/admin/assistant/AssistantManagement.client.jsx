'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Bot, 
  RefreshCw, 
  Plus, 
  Edit, 
  Trash2, 
  MessageSquare, 
  Activity,
  Settings,
  BarChart3,
  Brain,
  Zap,
  Clock,
  Users
} from 'lucide-react';
import styles from './AssistantManagement.module.css';
import { assistantsApi } from '@/api/client/assistantsApi';

export default function AssistantManagement() {
  const t = useTranslations('admin.assistant');
  const queryClient = useQueryClient();
  const [selectedAssistant, setSelectedAssistant] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  // Fetch assistants
  const { data: assistants = [], isLoading, error, refetch } = useQuery({
    queryKey: ['assistants'],
    queryFn: async () => {
      const response = await assistantsApi.getAssistants();
      return response.assistants || [];
    }
  });

  // Create assistant mutation - NOT AVAILABLE IN API
  const createMutation = useMutation({
    mutationFn: async (data) => {
      // API doesn't have create endpoint - only activate/deactivate
      throw new Error('Create assistant not available - use system admin');
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['assistants']);
      setIsModalOpen(false);
    }
  });

  // Update assistant mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, ...data }) => assistantsApi.updateAssistantConfiguration(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries(['assistants']);
      setIsModalOpen(false);
    }
  });

  // Deactivate assistant (no delete in API)
  const deleteMutation = useMutation({
    mutationFn: (id) => assistantsApi.deactivateAssistant(id),
    onSuccess: () => {
      queryClient.invalidateQueries(['assistants']);
      setIsDeleteModalOpen(false);
    }
  });

  const handleCreateAssistant = () => {
    setSelectedAssistant(null);
    setIsModalOpen(true);
  };

  const handleEditAssistant = (assistant) => {
    setSelectedAssistant(assistant);
    setIsModalOpen(true);
  };

  const handleDeleteAssistant = (assistant) => {
    setSelectedAssistant(assistant);
    setIsDeleteModalOpen(true);
  };

  const handleSubmit = (formData) => {
    if (selectedAssistant) {
      updateMutation.mutate({ id: selectedAssistant.id, ...formData });
    } else {
      createMutation.mutate(formData);
    }
  };

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <p>Error loading assistants: {error.message}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <h1 className={styles.title}>
            <Bot size={24} />
            {t('title', 'Assistant Management')}
          </h1>
          <p className={styles.subtitle}>
            {t('subtitle', 'Manage AI assistants and chatbots for your platform')}
          </p>
        </div>
        <div className={styles.headerActions}>
          <button 
            className={styles.primaryButton} 
            onClick={handleCreateAssistant}
          >
            <Plus size={20} />
            {t('createAssistant', 'Create Assistant')}
          </button>
          <button 
            className={styles.iconButton} 
            onClick={() => refetch()}
            disabled={isLoading}
          >
            <RefreshCw size={20} className={isLoading ? styles.spinning : ''} />
          </button>
        </div>
      </div>

      {/* Metrics */}
      <div className={styles.metricsGrid}>
        <div className={styles.metricCard}>
          <div className={styles.metricIcon}>
            <Bot size={20} />
          </div>
          <div className={styles.metricContent}>
            <div className={styles.metricValue}>{assistants.length}</div>
            <div className={styles.metricLabel}>{t('totalAssistants', 'Total Assistants')}</div>
          </div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricIcon}>
            <Activity size={20} />
          </div>
          <div className={styles.metricContent}>
            <div className={styles.metricValue}>
              {assistants.filter(a => a.status === 'active').length}
            </div>
            <div className={styles.metricLabel}>{t('activeAssistants', 'Active')}</div>
          </div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricIcon}>
            <MessageSquare size={20} />
          </div>
          <div className={styles.metricContent}>
            <div className={styles.metricValue}>
              {assistants.reduce((sum, a) => sum + (a.conversationCount || 0), 0)}
            </div>
            <div className={styles.metricLabel}>{t('totalConversations', 'Conversations')}</div>
          </div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricIcon}>
            <BarChart3 size={20} />
          </div>
          <div className={styles.metricContent}>
            <div className={styles.metricValue}>
              {(assistants.reduce((sum, a) => sum + (a.satisfactionScore || 0), 0) / assistants.length || 0).toFixed(1)}
            </div>
            <div className={styles.metricLabel}>{t('avgSatisfaction', 'Avg Satisfaction')}</div>
          </div>
        </div>
      </div>

      {/* Assistants Grid */}
      {isLoading ? (
        <div className={styles.loadingContainer}>
          <RefreshCw size={32} className={styles.spinning} />
          <p>{t('loading', 'Loading assistants...')}</p>
        </div>
      ) : assistants.length === 0 ? (
        <div className={styles.emptyState}>
          <Bot size={64} />
          <h3>{t('noAssistants', 'No assistants yet')}</h3>
          <p>{t('noAssistantsDesc', 'Create your first AI assistant to get started')}</p>
          <button className={styles.primaryButton} onClick={handleCreateAssistant}>
            <Plus size={20} />
            {t('createFirstAssistant', 'Create First Assistant')}
          </button>
        </div>
      ) : (
        <div className={styles.assistantsGrid}>
          {assistants.map((assistant) => (
            <div key={assistant.id} className={styles.assistantCard}>
              <div className={styles.assistantHeader}>
                <h3 className={styles.assistantName}>{assistant.name}</h3>
                <span className={`${styles.statusBadge} ${styles[assistant.status]}`}>
                  {assistant.status}
                </span>
              </div>
              
              <p className={styles.assistantDescription}>
                {assistant.description || t('noDescription', 'No description provided')}
              </p>

              <div className={styles.assistantStats}>
                <div className={styles.statItem}>
                  <MessageSquare size={16} />
                  <span>{assistant.conversationCount || 0} {t('conversations', 'conversations')}</span>
                </div>
                <div className={styles.statItem}>
                  <Users size={16} />
                  <span>{assistant.userCount || 0} {t('users', 'users')}</span>
                </div>
                <div className={styles.statItem}>
                  <Clock size={16} />
                  <span>{assistant.avgResponseTime || '0'}ms</span>
                </div>
              </div>

              <div className={styles.assistantFeatures}>
                {assistant.features?.includes('nlp') && (
                  <span className={styles.featureBadge}>
                    <Brain size={14} />
                    NLP
                  </span>
                )}
                {assistant.features?.includes('ml') && (
                  <span className={styles.featureBadge}>
                    <Zap size={14} />
                    ML
                  </span>
                )}
              </div>

              <div className={styles.assistantActions}>
                <button 
                  className={`${styles.actionButton} ${styles.editButton}`}
                  onClick={() => handleEditAssistant(assistant)}
                  title={t('edit', 'Edit')}
                >
                  <Edit size={16} />
                </button>
                <button 
                  className={`${styles.actionButton} ${styles.settingsButton}`}
                  onClick={() => {/* Handle settings */}}
                  title={t('settings', 'Settings')}
                >
                  <Settings size={16} />
                </button>
                <button 
                  className={`${styles.actionButton} ${styles.deleteButton}`}
                  onClick={() => handleDeleteAssistant(assistant)}
                  title={t('delete', 'Delete')}
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create/Edit Modal */}
      {isModalOpen && (
        <AssistantModal
          assistant={selectedAssistant}
          onSubmit={handleSubmit}
          onClose={() => setIsModalOpen(false)}
          isLoading={createMutation.isLoading || updateMutation.isLoading}
        />
      )}

      {/* Delete Confirmation Modal */}
      {isDeleteModalOpen && (
        <DeleteConfirmationModal
          assistant={selectedAssistant}
          onConfirm={() => deleteMutation.mutate(selectedAssistant.id)}
          onClose={() => setIsDeleteModalOpen(false)}
          isLoading={deleteMutation.isLoading}
        />
      )}
    </div>
  );
}

// Assistant Modal Component
function AssistantModal({ assistant, onSubmit, onClose, isLoading }) {
  const t = useTranslations('admin.assistant');
  const [formData, setFormData] = useState({
    name: assistant?.name || '',
    description: assistant?.description || '',
    status: assistant?.status || 'active',
    model: assistant?.model || 'gpt-3.5-turbo',
    features: assistant?.features || [],
    systemPrompt: assistant?.systemPrompt || '',
    temperature: assistant?.temperature || 0.7,
    maxTokens: assistant?.maxTokens || 1000
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3 className={styles.modalTitle}>
            {assistant ? t('editAssistant', 'Edit Assistant') : t('createAssistant', 'Create Assistant')}
          </h3>
          <button className={styles.closeButton} onClick={onClose}>
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.modalBody}>
          <div className={styles.formGroup}>
            <label className={styles.formLabel}>{t('name', 'Name')}</label>
            <input
              type="text"
              className={styles.formInput}
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label className={styles.formLabel}>{t('description', 'Description')}</label>
            <textarea
              className={styles.formTextarea}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={3}
            />
          </div>

          <div className={styles.formGroup}>
            <label className={styles.formLabel}>{t('model', 'Model')}</label>
            <select
              className={styles.formSelect}
              value={formData.model}
              onChange={(e) => setFormData({ ...formData, model: e.target.value })}
            >
              <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
              <option value="gpt-4">GPT-4</option>
              <option value="claude-2">Claude 2</option>
              <option value="llama-2">Llama 2</option>
            </select>
          </div>

          <div className={styles.formGroup}>
            <label className={styles.formLabel}>{t('systemPrompt', 'System Prompt')}</label>
            <textarea
              className={styles.formTextarea}
              value={formData.systemPrompt}
              onChange={(e) => setFormData({ ...formData, systemPrompt: e.target.value })}
              rows={4}
              placeholder={t('systemPromptPlaceholder', 'Enter the system prompt for the assistant...')}
            />
          </div>

          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label className={styles.formLabel}>{t('temperature', 'Temperature')}</label>
              <input
                type="number"
                className={styles.formInput}
                value={formData.temperature}
                onChange={(e) => setFormData({ ...formData, temperature: parseFloat(e.target.value) })}
                min="0"
                max="2"
                step="0.1"
              />
            </div>

            <div className={styles.formGroup}>
              <label className={styles.formLabel}>{t('maxTokens', 'Max Tokens')}</label>
              <input
                type="number"
                className={styles.formInput}
                value={formData.maxTokens}
                onChange={(e) => setFormData({ ...formData, maxTokens: parseInt(e.target.value) })}
                min="1"
                max="4096"
              />
            </div>
          </div>

          <div className={styles.formGroup}>
            <label className={styles.formLabel}>{t('features', 'Features')}</label>
            <div className={styles.checkboxGroup}>
              <label className={styles.checkboxLabel}>
                <input
                  type="checkbox"
                  checked={formData.features.includes('nlp')}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setFormData({ ...formData, features: [...formData.features, 'nlp'] });
                    } else {
                      setFormData({ ...formData, features: formData.features.filter(f => f !== 'nlp') });
                    }
                  }}
                />
                Natural Language Processing
              </label>
              <label className={styles.checkboxLabel}>
                <input
                  type="checkbox"
                  checked={formData.features.includes('ml')}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setFormData({ ...formData, features: [...formData.features, 'ml'] });
                    } else {
                      setFormData({ ...formData, features: formData.features.filter(f => f !== 'ml') });
                    }
                  }}
                />
                Machine Learning
              </label>
            </div>
          </div>

          <div className={styles.modalActions}>
            <button 
              type="button" 
              className={styles.cancelButton} 
              onClick={onClose}
              disabled={isLoading}
            >
              {t('cancel', 'Cancel')}
            </button>
            <button 
              type="submit" 
              className={styles.submitButton}
              disabled={isLoading}
            >
              {isLoading ? t('saving', 'Saving...') : t('save', 'Save')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// Delete Confirmation Modal
function DeleteConfirmationModal({ assistant, onConfirm, onClose, isLoading }) {
  const t = useTranslations('admin.assistant');

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h3 className={styles.modalTitle}>{t('deleteConfirmation', 'Delete Assistant')}</h3>
        </div>
        <div className={styles.modalBody}>
          <p>
            {t('deleteConfirmationText', 'Are you sure you want to delete')} <strong>{assistant?.name}</strong>?
          </p>
          <p className={styles.warningText}>
            {t('deleteWarning', 'This action cannot be undone.')}
          </p>
        </div>
        <div className={styles.modalActions}>
          <button 
            className={styles.cancelButton} 
            onClick={onClose}
            disabled={isLoading}
          >
            {t('cancel', 'Cancel')}
          </button>
          <button 
            className={styles.deleteConfirmButton} 
            onClick={onConfirm}
            disabled={isLoading}
          >
            {isLoading ? t('deleting', 'Deleting...') : t('delete', 'Delete')}
          </button>
        </div>
      </div>
    </div>
  );
}
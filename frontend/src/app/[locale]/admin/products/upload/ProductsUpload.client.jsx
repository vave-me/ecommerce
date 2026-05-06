"use client";

import React, { useState, useCallback, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  ArrowLeft,
  Upload,
  Download,
  FileText,
  CheckCircle,
  XCircle,
  AlertTriangle,
  RefreshCw,
  Eye,
  Trash2,
  Play,
  Pause,
  BarChart3,
  Package,
  Calendar,
  Clock,
  FileSpreadsheet,
  AlertCircle
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  uploadProductsFile,
  getUploadHistory,
  getUploadTemplate,
  processUpload,
  deleteUpload,
  getUploadStats
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ProductsUpload.module.css';

const UploadStatusBadge = ({ status }) => {
  const statusConfig = {
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending', icon: Clock },
    processing: { color: '#3b82f6', bg: 'rgba(59, 130, 246, 0.1)', text: 'Processing', icon: RefreshCw },
    completed: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Completed', icon: CheckCircle },
    failed: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Failed', icon: XCircle },
    cancelled: { color: '#6b7280', bg: 'rgba(107, 114, 128, 0.1)', text: 'Cancelled', icon: XCircle }
  };

  const config = statusConfig[status] || statusConfig.pending;
  const Icon = config.icon;

  return (
    <span 
      className={styles.statusBadge}
      style={{ color: config.color, backgroundColor: config.bg }}
    >
      <Icon size={14} />
      {config.text}
    </span>
  );
};

const UploadItem = ({ upload, onProcess, onDelete, onView, loading }) => {
  return (
    <div className={styles.uploadItem}>
      <div className={styles.uploadHeader}>
        <div className={styles.uploadInfo}>
          <div className={styles.fileName}>
            <FileSpreadsheet size={16} />
            {upload.fileName}
          </div>
          <div className={styles.uploadMeta}>
            <span>{upload.totalRows} rows</span>
            <span>{new Date(upload.createdAt).toLocaleDateString()}</span>
          </div>
        </div>
        <UploadStatusBadge status={upload.status} />
      </div>

      <div className={styles.uploadStats}>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Valid</span>
          <span className={styles.statValue}>{upload.validRows || 0}</span>
        </div>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Errors</span>
          <span className={styles.statValue}>{upload.errorRows || 0}</span>
        </div>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Progress</span>
          <span className={styles.statValue}>{upload.progress || 0}%</span>
        </div>
      </div>

      {upload.progress > 0 && upload.progress < 100 && (
        <div className={styles.progressBar}>
          <div 
            className={styles.progressFill}
            style={{ width: `${upload.progress}%` }}
          />
        </div>
      )}

      <div className={styles.uploadActions}>
        {upload.status === 'pending' && (
          <button
            className={`${styles.actionButton} ${styles.process}`}
            onClick={() => onProcess(upload.id)}
            disabled={loading === upload.id}
          >
            <Play size={16} />
            Process
          </button>
        )}
        {upload.status === 'processing' && (
          <button
            className={`${styles.actionButton} ${styles.pause}`}
            onClick={() => onProcess(upload.id, 'pause')}
            disabled={loading === upload.id}
          >
            <Pause size={16} />
            Pause
          </button>
        )}
        <button
          className={styles.actionButton}
          onClick={() => onView(upload.id)}
        >
          <Eye size={16} />
          View Details
        </button>
        {upload.status !== 'processing' && (
          <button
            className={`${styles.actionButton} ${styles.delete}`}
            onClick={() => onDelete(upload.id)}
            disabled={loading === upload.id}
          >
            <Trash2 size={16} />
            Delete
          </button>
        )}
      </div>
    </div>
  );
};

const ProductsUpload = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('ProductsUpload');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();
  const fileInputRef = useRef(null);

  const [dragOver, setDragOver] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [actionLoading, setActionLoading] = useState(null);

  // Fetch upload data
  const { data: uploadStats, isLoading: statsLoading, refetch: refetchStats } = useQuery({
    queryKey: ['uploads', 'stats'],
    queryFn: getUploadStats,
    staleTime: 30 * 1000,
    enabled: isAdmin
  });

  const { data: uploadHistory, isLoading: historyLoading, refetch: refetchHistory } = useQuery({
    queryKey: ['uploads', 'history'],
    queryFn: () => getUploadHistory({ limit: 20 }),
    staleTime: 10 * 1000,
    enabled: isAdmin
  });

  // Mutations
  const uploadMutation = useMutation({
    mutationFn: uploadProductsFile,
    onSuccess: () => {
      queryClient.invalidateQueries(['uploads']);
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    },
    onError: (error) => {
      // Error: 'Upload failed:', error...
      setUploading(false);
    }
  });

  const processMutation = useMutation({
    mutationFn: ({ uploadId, action }) => processUpload(uploadId, action),
    onSuccess: () => {
      queryClient.invalidateQueries(['uploads']);
      setActionLoading(null);
    },
    onError: (error) => {
      // Error: 'Process action failed:', error...
      setActionLoading(null);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteUpload,
    onSuccess: () => {
      queryClient.invalidateQueries(['uploads']);
      setActionLoading(null);
    },
    onError: (error) => {
      // Error: 'Delete failed:', error...
      setActionLoading(null);
    }
  });

  const handleFileSelect = useCallback((files) => {
    const file = files[0];
    if (!file) return;

    // Validate file type
    const validTypes = [
      'text/csv',
      'application/vnd.ms-excel',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    ];

    if (!validTypes.includes(file.type)) {
      alert('Please upload a CSV or Excel file');
      return;
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('File size must be less than 10MB');
      return;
    }

    setUploading(true);
    uploadMutation.mutate(file);
  }, [uploadMutation]);

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    setDragOver(false);
    handleFileSelect(e.dataTransfer.files);
  }, [handleFileSelect]);

  const handleDragOver = useCallback((e) => {
    e.preventDefault();
    setDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e) => {
    e.preventDefault();
    setDragOver(false);
  }, []);

  const handleFileInputChange = useCallback((e) => {
    handleFileSelect(e.target.files);
  }, [handleFileSelect]);

  const handleProcess = useCallback((uploadId, action = 'start') => {
    setActionLoading(uploadId);
    processMutation.mutate({ uploadId, action });
  }, [processMutation]);

  const handleDelete = useCallback((uploadId) => {
    if (window.confirm('Are you sure you want to delete this upload?')) {
      setActionLoading(uploadId);
      deleteMutation.mutate(uploadId);
    }
  }, [deleteMutation]);

  const handleView = useCallback((uploadId) => {
    router.push(`/admin/products/upload/${uploadId}`);
  }, [router]);

  const handleDownloadTemplate = useCallback(async () => {
    try {
      const template = await getUploadTemplate();
      const blob = new Blob([template], { type: 'text/csv' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'products-template.csv';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (error) {
        // File operation error
        if (process.env.NODE_ENV === 'development') {
            console.error('File operation error:', error);
        }
        return null; // Return null for failed file operations
    }
  }, []);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access bulk upload.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  const isLoading = statsLoading || historyLoading;

  if (isLoading && !uploadStats && !uploadHistory) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Upload Manager...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch upload data.' })}</p>
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
            <button 
              className={styles.backButton}
              onClick={() => router.back()}
            >
              <ArrowLeft size={16} />
              {t('backToProducts', { defaultValue: 'Back to Products' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <Upload size={24} />
                {t('title', { defaultValue: 'Bulk Upload Products' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Upload and manage products in bulk using CSV or Excel files' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <button 
              className={styles.templateButton}
              onClick={handleDownloadTemplate}
            >
              <Download size={16} />
              {t('downloadTemplate', { defaultValue: 'Download Template' })}
            </button>
          </div>
        </div>

        {/* Stats Grid */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <FileText size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{uploadStats?.totalUploads || 0}</div>
              <div className={styles.statLabel}>{t('totalUploads', { defaultValue: 'Total Uploads' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <Clock size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{uploadStats?.pendingUploads || 0}</div>
              <div className={styles.statLabel}>{t('pendingUploads', { defaultValue: 'Pending Processing' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <Package size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{uploadStats?.processedProducts || 0}</div>
              <div className={styles.statLabel}>{t('processedProducts', { defaultValue: 'Products Processed' })}</div>
            </div>
          </div>
          
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <XCircle size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{uploadStats?.errorCount || 0}</div>
              <div className={styles.statLabel}>{t('errorCount', { defaultValue: 'Processing Errors' })}</div>
            </div>
          </div>
        </div>

        {/* Upload Area */}
        <div className={styles.uploadSection}>
          <h2 className={styles.sectionTitle}>{t('uploadNew', { defaultValue: 'Upload New File' })}</h2>
          <div 
            className={`${styles.dropzone} ${dragOver ? styles.dragOver : ''} ${uploading ? styles.uploading : ''}`}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onClick={() => fileInputRef.current?.click()}
          >
            <input
              ref={fileInputRef}
              type="file"
              accept=".csv,.xlsx,.xls"
              onChange={handleFileInputChange}
              style={{ display: 'none' }}
              disabled={uploading}
            />
            
            {uploading ? (
              <div className={styles.uploadingState}>
                <LoadingSpinner size={32} />
                <h3>{t('uploading', { defaultValue: 'Uploading...' })}</h3>
                <p>{t('uploadingMessage', { defaultValue: 'Please wait while we process your file' })}</p>
              </div>
            ) : (
              <div className={styles.dropzoneContent}>
                <Upload size={48} />
                <h3>{t('dropzoneTitle', { defaultValue: 'Drop your file here or click to browse' })}</h3>
                <p>{t('dropzoneMessage', { defaultValue: 'Supports CSV and Excel files up to 10MB' })}</p>
                <div className={styles.formatInfo}>
                  <span>CSV</span>
                  <span>XLSX</span>
                  <span>XLS</span>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Upload History */}
        <div className={styles.historySection}>
          <h2 className={styles.sectionTitle}>{t('uploadHistory', { defaultValue: 'Upload History' })}</h2>
          <div className={styles.uploadList}>
            {historyLoading ? (
              <div className={styles.loadingState}>
                <LoadingSpinner />
                <p>{t('loadingHistory', { defaultValue: 'Loading upload history...' })}</p>
              </div>
            ) : uploadHistory?.uploads?.length > 0 ? (
              uploadHistory.uploads.map((upload) => (
                <UploadItem
                  key={upload.id}
                  upload={upload}
                  onProcess={handleProcess}
                  onDelete={handleDelete}
                  onView={handleView}
                  loading={actionLoading}
                />
              ))
            ) : (
              <div className={styles.emptyState}>
                <FileText size={48} />
                <h3>{t('noUploads', { defaultValue: 'No uploads yet' })}</h3>
                <p>{t('noUploadsMessage', { defaultValue: 'Upload your first file to get started' })}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default ProductsUpload; 
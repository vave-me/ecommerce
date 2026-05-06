"use client";

import React, { useState, useCallback, useMemo, useRef } from 'react';
import { useTranslations } from 'next-intl';
import { 
  Image, 
  Video, 
  FileText, 
  Download, 
  Trash2, 
  Eye, 
  Upload,
  Search,
  Grid,
  List,
  Plus,
  AlertCircle,
  CheckCircle,
  Clock
} from 'lucide-react';
import { toast } from 'react-toastify';
import { useAdminMediaApi } from '@/hooks/useAdminMediaApi';
import styles from './MediaManagement.module.css';

const MediaManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('AdminMedia');
  } catch (e) {
    // Fallback function for missing translations
    t = (key, options) => options?.defaultValue || key;
  }
  
  // State management
  const [viewMode, setViewMode] = useState('grid');
  const [selectedType, setSelectedType] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedItems, setSelectedItems] = useState([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [sortBy, setSortBy] = useState('uploadedAt');
  const [sortOrder] = useState('desc');
  const [showUploadModal, setShowUploadModal] = useState(false);
  const fileInputRef = useRef(null);

  // API integration
  const {
    useAdminMediaQuery,
    deleteMedia,
    deleteImage,
    deleteVideo,
    uploadImage,
    uploadVideo,
    createMedia,
    isDeleting,
    isUploading,
    uploadProgress,
    extractFileName,
    calculateFileSize
  } = useAdminMediaApi();

  // Query parameters
  const queryParams = useMemo(() => ({
    page,
    pageSize,
    type: selectedType,
    sortBy,
    sortOrder,
    searchQuery: searchQuery.trim()
  }), [page, pageSize, selectedType, sortBy, sortOrder, searchQuery]);

  // Fetch media data
  const { data: mediaData, isLoading, error, refetch } = useAdminMediaQuery(queryParams);

  // Handle file upload
  const handleFileUpload = useCallback(async (files) => {
    if (!files || files.length === 0) return;

    try {
      for (const file of files) {
        // Create media container first
        const mediaContainer = await createMedia({
          itemId: `admin-upload-${Date.now()}`,
          itemType: 'admin-media',
          status: 'active'
        });

        // Upload file based on type
        if (file.type.startsWith('image/')) {
          await uploadImage({
            file,
            mediaId: mediaContainer.id,
            displayOrder: 0,
            isMain: true
          });
        } else if (file.type.startsWith('video/')) {
          await uploadVideo({
            file,
            mediaId: mediaContainer.id,
            displayOrder: 0,
            isMain: true
          });
        } else {
          throw new Error(`Unsupported file type: ${file.type}`);
        }
      }

      toast.success(t('uploadSuccess', { count: files.length }));
      setShowUploadModal(false);
      refetch();
    } catch (error) {
      // Error: 'Upload error:', error...
      toast.error(t('uploadError'));
    }
  }, [createMedia, uploadImage, uploadVideo, t, refetch]);

  // Handle file input change
  const handleFileInputChange = useCallback((event) => {
    const files = Array.from(event.target.files || []);
    if (files.length > 0) {
      handleFileUpload(files);
    }
    // Reset input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  }, [handleFileUpload]);

  // Handle delete single item
  const handleDelete = useCallback(async (item) => {
    if (!window.confirm(t('confirmDelete'))) return;

    try {
      if (item.type === 'image') {
        await deleteImage({ mediaId: item.mediaId, imageId: item.id });
      } else if (item.type === 'video') {
        await deleteVideo({ mediaId: item.mediaId, videoId: item.id });
      } else {
        await deleteMedia(item.id);
      }

      toast.success(t('deleteSuccess'));
      refetch();
    } catch (error) {
      // Error: 'Delete error:', error...
      toast.error(t('deleteError'));
    }
  }, [deleteImage, deleteVideo, deleteMedia, t, refetch]);

  // Handle bulk delete
  const handleBulkDelete = useCallback(async () => {
    if (selectedItems.length === 0) return;
    if (!window.confirm(t('confirmBulkDelete', { count: selectedItems.length }))) return;

    try {
      const items = mediaData?.items?.filter(item => selectedItems.includes(item.id)) || [];
      
      for (const item of items) {
        if (item.type === 'image') {
          await deleteImage({ mediaId: item.mediaId, imageId: item.id });
        } else if (item.type === 'video') {
          await deleteVideo({ mediaId: item.mediaId, videoId: item.id });
        } else {
          await deleteMedia(item.id);
        }
      }

      toast.success(t('bulkDeleteSuccess', { count: selectedItems.length }));
      setSelectedItems([]);
      refetch();
    } catch (error) {
      // Error: 'Bulk delete error:', error...
      toast.error(t('bulkDeleteError'));
    }
  }, [selectedItems, mediaData, deleteImage, deleteVideo, deleteMedia, t, refetch]);

  // Handle view/download
  const handleView = useCallback((item) => {
    if (item.url) {
      window.open(item.url, '_blank');
    }
  }, []);

  const handleDownload = useCallback((item) => {
    if (item.url) {
      const link = document.createElement('a');
      link.href = item.url;
      link.download = item.name;
      link.click();
    }
  }, []);

  // Get media icon
  const getMediaIcon = (type) => {
    switch (type) {
      case 'image':
        return <Image className={styles.typeIcon} />;
      case 'video':
        return <Video className={styles.typeIcon} />;
      default:
        return <FileText className={styles.typeIcon} />;
    }
  };

  // Filter and sort items
  const filteredItems = useMemo(() => {
    return mediaData?.items || [];
  }, [mediaData]);

  // Handle item selection
  const handleItemSelect = useCallback((itemId, checked) => {
    if (checked) {
      setSelectedItems(prev => [...prev, itemId]);
    } else {
      setSelectedItems(prev => prev.filter(id => id !== itemId));
    }
  }, []);

  const handleSelectAll = useCallback((checked) => {
    if (checked) {
      setSelectedItems(filteredItems.map(item => item.id));
    } else {
      setSelectedItems([]);
    }
  }, [filteredItems]);

  // Render upload modal
  const renderUploadModal = () => {
    if (!showUploadModal) return null;

    return (
      <div className={styles.uploadModal} onClick={() => setShowUploadModal(false)}>
        <div className={styles.uploadModalContent} onClick={(e) => e.stopPropagation()}>
          <h3>{t('uploadFiles')}</h3>
          
          <div className={styles.uploadArea}>
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/*,video/*"
              onChange={handleFileInputChange}
              className={styles.fileInput}
            />
            
            <div className={styles.uploadInstructions}>
              <Upload size={48} />
              <p>{t('uploadInstructions')}</p>
              <button 
                className={styles.selectFilesButton}
                onClick={() => fileInputRef.current?.click()}
              >
                {t('selectFiles')}
              </button>
            </div>
          </div>

          {isUploading && (
            <div className={styles.uploadProgress}>
              <div className={styles.progressBar}>
                <div 
                  className={styles.progressFill} 
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <span>{uploadProgress}%</span>
            </div>
          )}

          <div className={styles.uploadModalActions}>
            <button 
              className={styles.cancelButton}
              onClick={() => setShowUploadModal(false)}
              disabled={isUploading}
            >
              {t('cancel')}
            </button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>{t('title')}</h1>
          <p className={styles.subtitle}>{t('subtitle')}</p>
        </div>
        <button 
          className={styles.uploadButton}
          onClick={() => setShowUploadModal(true)}
          disabled={isUploading}
        >
          {isUploading ? (
            <>
              <Clock size={20} />
              {t('uploading')}
            </>
          ) : (
            <>
              <Upload size={20} />
              {t('uploadNew')}
            </>
          )}
        </button>
      </div>

      {/* Controls */}
      <div className={styles.controls}>
        <div className={styles.searchBar}>
          <Search size={20} />
          <input
            type="text"
            placeholder={t('searchPlaceholder')}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className={styles.searchInput}
          />
        </div>

        <div className={styles.filters}>
          <select 
            value={selectedType} 
            onChange={(e) => setSelectedType(e.target.value)}
            className={styles.filterSelect}
          >
            <option value="all">{t('filterAll')}</option>
            <option value="image">{t('filterImages')}</option>
            <option value="video">{t('filterVideos')}</option>
            <option value="document">{t('filterDocuments')}</option>
          </select>

          <select 
            value={sortBy} 
            onChange={(e) => setSortBy(e.target.value)}
            className={styles.filterSelect}
          >
            <option value="uploadedAt">{t('sortByDate')}</option>
            <option value="name">{t('sortByName')}</option>
            <option value="size">{t('sortBySize')}</option>
            <option value="type">{t('sortByType')}</option>
          </select>

          <div className={styles.viewToggle}>
            <button
              className={viewMode === 'grid' ? styles.active : ''}
              onClick={() => setViewMode('grid')}
            >
              <Grid size={20} />
            </button>
            <button
              className={viewMode === 'list' ? styles.active : ''}
              onClick={() => setViewMode('list')}
            >
              <List size={20} />
            </button>
          </div>
        </div>
      </div>

      {/* Bulk Actions */}
      {selectedItems.length > 0 && (
        <div className={styles.bulkActions}>
          <div className={styles.bulkInfo}>
            <input
              type="checkbox"
              checked={selectedItems.length === filteredItems.length}
              onChange={(e) => handleSelectAll(e.target.checked)}
            />
            <span>{t('selected', { count: selectedItems.length })}</span>
          </div>
          <button 
            className={styles.bulkDeleteButton}
            onClick={handleBulkDelete}
            disabled={isDeleting}
          >
            <Trash2 size={16} />
            {isDeleting ? t('deleting') : t('deleteSelected')}
          </button>
        </div>
      )}

      {/* Loading State */}
      {isLoading && (
        <div className={styles.loading}>
          <div className={styles.spinner} />
          {t('loading')}
        </div>
      )}

      {/* Error State */}
      {error && (
        <div className={styles.errorState}>
          <AlertCircle size={48} className={styles.errorIcon} />
          <h2>{t('errorTitle')}</h2>
          <p>{error.message}</p>
          <button className={styles.retryButton} onClick={() => refetch()}>
            {t('retry')}
          </button>
        </div>
      )}

      {/* Media Grid/List */}
      {!isLoading && !error && (
        <div className={viewMode === 'grid' ? styles.gridView : styles.listView}>
          {filteredItems.length === 0 ? (
            <div className={styles.empty}>
              <FileText size={48} />
              <p>{t('noMedia')}</p>
              <button 
                className={styles.uploadButton}
                onClick={() => setShowUploadModal(true)}
              >
                <Plus size={16} />
                {t('uploadFirst')}
              </button>
            </div>
          ) : (
            filteredItems.map(item => (
              <div key={item.id} className={styles.mediaItem}>
                <input
                  type="checkbox"
                  checked={selectedItems.includes(item.id)}
                  onChange={(e) => handleItemSelect(item.id, e.target.checked)}
                  className={styles.checkbox}
                />
                
                <div className={styles.preview}>
                  {item.type === 'image' ? (
                    <img src={item.thumbnail || item.url} alt={item.name} />
                  ) : (
                    getMediaIcon(item.type)
                  )}
                </div>

                <div className={styles.info}>
                  <h3>{item.name}</h3>
                  <p>{item.size} • {new Date(item.uploadedAt).toLocaleDateString()}</p>
                  {item.dimensions && <p>{item.dimensions}</p>}
                  {item.duration && <p>{item.duration}</p>}
                </div>

                <div className={styles.actions}>
                  <button 
                    className={styles.actionButton} 
                    onClick={() => handleView(item)}
                    title={t('view')}
                  >
                    <Eye size={16} />
                  </button>
                  <button 
                    className={styles.actionButton} 
                    onClick={() => handleDownload(item)}
                    title={t('download')}
                  >
                    <Download size={16} />
                  </button>
                  <button 
                    className={styles.actionButton} 
                    onClick={() => handleDelete(item)}
                    title={t('delete')}
                    disabled={isDeleting}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* Pagination */}
      {mediaData && mediaData.totalPages > 1 && (
        <div className={styles.pagination}>
          <button 
            onClick={() => setPage(prev => Math.max(1, prev - 1))}
            disabled={page === 1}
            className={styles.paginationButton}
          >
            {t('previous')}
          </button>
          
          <span className={styles.paginationInfo}>
            {t('pageInfo', { 
              current: page, 
              total: mediaData.totalPages,
              items: mediaData.total 
            })}
          </span>
          
          <button 
            onClick={() => setPage(prev => Math.min(mediaData.totalPages, prev + 1))}
            disabled={page === mediaData.totalPages}
            className={styles.paginationButton}
          >
            {t('next')}
          </button>
        </div>
      )}

      {/* Upload Modal */}
      {renderUploadModal()}
    </div>
  );
};

export default MediaManagement;
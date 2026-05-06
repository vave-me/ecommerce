'use client';

import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/api/axiosInstance';

/**
 * Admin Media API Hook
 * Integrates with the media management Swagger API
 * Provides comprehensive media management functionality for admin interface
 */
export const useAdminMediaApi = () => {
  const queryClient = useQueryClient();
  const [uploadProgress, setUploadProgress] = useState(0);

  // ===== QUERY FUNCTIONS =====

  /**
   * Get all media files with filtering, pagination, and sorting
   * Integrates with multiple endpoints for comprehensive media listing
   */
  const fetchAdminMedia = useCallback(async ({ 
    page = 1, 
    pageSize = 20, 
    type = 'all', 
    sortBy = 'uploadedAt', 
    sortOrder = 'desc',
    userId = null,
    searchQuery = ''
  }) => {
    try {
      let allMedia = [];
      let totalCount = 0;

      // Based on the Swagger spec, we need to combine different endpoints
      if (type === 'all' || type === 'image') {
        // Get images by author if userId provided, otherwise get all images
        if (userId) {
          const imageResponse = await axiosInstance.get(`/api/media/image/author/${userId}`, {
            params: { page, pageSize, sortBy, sortOrder }
          });
          
          if (imageResponse.data.images) {
            const images = imageResponse.data.images.map(img => ({
              id: img.id,
              type: 'image',
              name: extractFileName(img.url),
              url: img.url,
              thumbnail: img.thumbnail,
              size: calculateFileSize(img.metadata),
              uploadedAt: img.metadata ? JSON.parse(img.metadata).uploadedAt : new Date().toISOString(),
              uploadedBy: img.userId,
              dimensions: img.metadata ? JSON.parse(img.metadata).dimensions : null,
              mediaId: img.mediaId,
              isMain: img.isMain,
              displayOrder: img.displayOrder
            }));
            allMedia = allMedia.concat(images);
            totalCount += imageResponse.data.totalCount || images.length;
          }
        }
      }

      if (type === 'all' || type === 'video') {
        // Get videos by author if userId provided, otherwise get all videos
        if (userId) {
          const videoResponse = await axiosInstance.get(`/api/media/video/author/${userId}`, {
            params: { page, pageSize, sortBy, sortOrder }
          });
          
          if (videoResponse.data.videos) {
            const videos = videoResponse.data.videos.map(video => ({
              id: video.id,
              type: 'video',
              name: extractFileName(video.url),
              url: video.url,
              thumbnail: video.thumbnail,
              size: calculateFileSize(video.metadata),
              uploadedAt: video.metadata ? JSON.parse(video.metadata).uploadedAt : new Date().toISOString(),
              uploadedBy: video.userId,
              duration: video.metadata ? JSON.parse(video.metadata).duration : null,
              mediaId: video.mediaId,
              isMain: video.isMain,
              displayOrder: video.displayOrder
            }));
            allMedia = allMedia.concat(videos);
            totalCount += videoResponse.data.totalCount || videos.length;
          }
        } else {
          // Get all videos
          const videoResponse = await axiosInstance.get('/api/media/videos', {
            params: { page, pageSize, sortBy, sortOrder }
          });
          
          if (videoResponse.data.videos) {
            const videos = videoResponse.data.videos.map(video => ({
              id: video.id,
              type: 'video',
              name: extractFileName(video.url),
              url: video.url,
              thumbnail: video.thumbnail,
              size: calculateFileSize(video.metadata),
              uploadedAt: video.metadata ? JSON.parse(video.metadata).uploadedAt : new Date().toISOString(),
              uploadedBy: video.userId,
              duration: video.metadata ? JSON.parse(video.metadata).duration : null,
              mediaId: video.mediaId,
              isMain: video.isMain,
              displayOrder: video.displayOrder
            }));
            allMedia = allMedia.concat(videos);
            totalCount += videoResponse.data.totalCount || videos.length;
          }
        }
      }

      // Apply search filter if provided
      if (searchQuery) {
        allMedia = allMedia.filter(item => 
          item.name.toLowerCase().includes(searchQuery.toLowerCase())
        );
      }

      // Sort the combined results
      allMedia.sort((a, b) => {
        const aValue = a[sortBy] || a.uploadedAt;
        const bValue = b[sortBy] || b.uploadedAt;
        
        if (sortOrder === 'desc') {
          return new Date(bValue) - new Date(aValue);
        }
        return new Date(aValue) - new Date(bValue);
      });

      return {
        items: allMedia,
        total: totalCount,
        page: parseInt(page),
        pageSize: parseInt(pageSize),
        totalPages: Math.ceil(totalCount / pageSize)
      };

    } catch (error) {
      // Error: 'Error fetching admin media:', error...
      throw error;
    }
  }, []);

  // ===== QUERIES =====

  const useAdminMediaQuery = (params = {}) => {
    return useQuery({
      queryKey: ['admin-media', params],
      queryFn: () => fetchAdminMedia(params),
      staleTime: 30000, // 30 seconds
      cacheTime: 300000, // 5 minutes
    });
  };

  // ===== MUTATIONS =====

  /**
   * Create new media container
   */
  const createMediaMutation = useMutation({
    mutationFn: async ({ itemId, itemType, status = 'active' }) => {
      const response = await axiosInstance.post('/api/media', {
        itemId,
        itemType,
        status
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
    }
  });

  /**
   * Upload image
   */
  const uploadImageMutation = useMutation({
    mutationFn: async ({ file, mediaId, displayOrder = 0, isMain = false }) => {
      const formData = new FormData();
      formData.append('file', file);
      
      // First upload the file to get URL, then add to media
      const uploadResponse = await axiosInstance.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          setUploadProgress(progress);
        }
      });

      // Add image to media record
      const addImageResponse = await axiosInstance.post('/api/media/image', {
        mediaId,
        displayOrder,
        isMain,
        url: uploadResponse.data.url,
        metadata: JSON.stringify({
          originalName: file.name,
          size: file.size,
          type: file.type,
          uploadedAt: new Date().toISOString(),
          dimensions: await getImageDimensions(file)
        }),
        fileType: file.type,
        thumbnail: uploadResponse.data.thumbnailUrl
      });

      return addImageResponse.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
      setUploadProgress(0);
    },
    onError: () => {
      setUploadProgress(0);
    }
  });

  /**
   * Upload video
   */
  const uploadVideoMutation = useMutation({
    mutationFn: async ({ file, mediaId, displayOrder = 0, isMain = false }) => {
      const formData = new FormData();
      formData.append('file', file);
      
      const uploadResponse = await axiosInstance.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          setUploadProgress(progress);
        }
      });

      const addVideoResponse = await axiosInstance.post('/api/media/video', {
        mediaId,
        displayOrder,
        isMain,
        url: uploadResponse.data.url,
        metadata: JSON.stringify({
          originalName: file.name,
          size: file.size,
          type: file.type,
          uploadedAt: new Date().toISOString(),
          duration: await getVideoDuration(file)
        }),
        fileType: file.type,
        thumbnail: uploadResponse.data.thumbnailUrl
      });

      return addVideoResponse.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
      setUploadProgress(0);
    },
    onError: () => {
      setUploadProgress(0);
    }
  });

  /**
   * Delete media
   */
  const deleteMediaMutation = useMutation({
    mutationFn: async (mediaId) => {
      const response = await axiosInstance.delete(`/api/media/${mediaId}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
    }
  });

  /**
   * Delete image
   */
  const deleteImageMutation = useMutation({
    mutationFn: async ({ mediaId, imageId }) => {
      const response = await axiosInstance.delete(`/api/media/${mediaId}/image/${imageId}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
    }
  });

  /**
   * Delete video
   */
  const deleteVideoMutation = useMutation({
    mutationFn: async ({ mediaId, videoId }) => {
      const response = await axiosInstance.delete(`/api/media/${mediaId}/video/${videoId}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['admin-media']);
    }
  });

  /**
   * Bulk import
   */
  const startBulkImportMutation = useMutation({
    mutationFn: async ({ externalSystemId, externalSystemType, estimatedCount, options = {} }) => {
      const response = await axiosInstance.post('/api/media/import', {
        externalSystemId,
        externalSystemType,
        estimatedCount,
        options
      });
      return response.data;
    }
  });

  // ===== UTILITY FUNCTIONS =====

  const extractFileName = (url) => {
    if (!url) return 'Unknown';
    return url.split('/').pop() || 'Unknown';
  };

  const calculateFileSize = (metadata) => {
    if (!metadata) return 'Unknown';
    try {
      const parsed = JSON.parse(metadata);
      const bytes = parsed.size || 0;
      if (bytes === 0) return 'Unknown';
      
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    } catch {
      return 'Unknown';
    }
  };

  const getImageDimensions = (file) => {
    return new Promise((resolve) => {
      if (!file.type.startsWith('image/')) {
        resolve(null);
        return;
      }

      const img = new Image();
      img.onload = () => {
        resolve(`${img.width}x${img.height}`);
      };
      img.onerror = () => resolve(null);
      img.src = URL.createObjectURL(file);
    });
  };

  const getVideoDuration = (file) => {
    return new Promise((resolve) => {
      if (!file.type.startsWith('video/')) {
        resolve(null);
        return;
      }

      const video = document.createElement('video');
      video.onloadedmetadata = () => {
        const duration = Math.floor(video.duration);
        const minutes = Math.floor(duration / 60);
        const seconds = duration % 60;
        resolve(`${minutes}:${seconds.toString().padStart(2, '0')}`);
      };
      video.onerror = () => resolve(null);
      video.src = URL.createObjectURL(file);
    });
  };

  return {
    // Queries
    useAdminMediaQuery,
    
    // Mutations
    createMedia: createMediaMutation.mutateAsync,
    uploadImage: uploadImageMutation.mutateAsync,
    uploadVideo: uploadVideoMutation.mutateAsync,
    deleteMedia: deleteMediaMutation.mutateAsync,
    deleteImage: deleteImageMutation.mutateAsync,
    deleteVideo: deleteVideoMutation.mutateAsync,
    startBulkImport: startBulkImportMutation.mutateAsync,
    
    // Loading states
    isCreating: createMediaMutation.isLoading,
    isUploading: uploadImageMutation.isLoading || uploadVideoMutation.isLoading,
    isDeleting: deleteMediaMutation.isLoading || deleteImageMutation.isLoading || deleteVideoMutation.isLoading,
    
    // Upload progress
    uploadProgress,
    
    // Utility functions
    extractFileName,
    calculateFileSize
  };
}; 
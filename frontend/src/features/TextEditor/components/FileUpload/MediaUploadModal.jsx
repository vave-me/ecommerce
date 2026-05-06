// TextEditor/components/FileUpload/MediaUploadModal.jsx
import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import { Upload, Image, Video } from '@/icons';
import { FaTimesCircle } from '../../../../utils/iconImports';
import styles from '../../TextEditor.module.css';
import { FileDropzone } from './FileDropzone';
import { FilePreview } from './FilePreview';
import { ProgressIndicator } from '../common/ProgressIndicator';
import { ErrorMessage } from '../common/ErrorMessage';
import {useAuth} from "../../../../context/AuthContext";
import {addImage, addVideo} from "../../../../api/client/mediaApi";
export function MediaUploadModal({ isOpen, onClose, onUploadSuccess, mediaId, type }) {
    const [file, setFile] = useState(null);
    const [previewUrl, setPreviewUrl] = useState('');
    const [uploading, setUploading] = useState(false);
    const [progress, setProgress] = useState(0);
    const [error, setError] = useState('');
    const { user } = useAuth();
    const modalRef = useRef(null);
    // Reset state when modal opens/closes
    useEffect(() => {
        if (isOpen) {
            setFile(null);
            setPreviewUrl('');
            setError('');
            setProgress(0);
        }
    }, [isOpen]);
    // Handle ESC key to close modal
    useEffect(() => {
        const handleKeyDown = (e) => {
            if (e.key === 'Escape' && isOpen) {
                onClose();
            }
        };
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [isOpen, onClose]);
    // Focus trap for accessibility
    useEffect(() => {
        if (isOpen && modalRef.current) {
            const focusableElements = modalRef.current.querySelectorAll(
                'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
            );
            if (focusableElements.length > 0) {
                const firstElement = focusableElements[0];
                const lastElement = focusableElements[focusableElements.length - 1];
                // Focus first element when modal opens
                firstElement.focus();
                // Set up focus trap
                const handleTabKey = (e) => {
                    if (e.key === 'Tab') {
                        if (e.shiftKey && document.activeElement === firstElement) {
                            e.preventDefault();
                            lastElement.focus();
                        } else if (!e.shiftKey && document.activeElement === lastElement) {
                            e.preventDefault();
                            firstElement.focus();
                        }
                    }
                };
                modalRef.current.addEventListener('keydown', handleTabKey);
                return () => modalRef.current?.removeEventListener('keydown', handleTabKey);
            }
        }
    }, [isOpen]);
    const handleFileChange = (e) => {
        const selectedFile = e.target?.files?.[0] || e;
        if (!selectedFile) return;
        setFile(selectedFile);
        setError('');
        // Create preview for image
        if (type === 'image' && selectedFile.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onload = () => setPreviewUrl(reader.result);
            reader.readAsDataURL(selectedFile);
        } else if (type === 'video' && selectedFile.type.startsWith('video/')) {
            // For video, we could create a preview, but for simplicity just setting the file
            setPreviewUrl('');
        }
    };
    const getFileExtension = (filename) => {
        return filename.split('.').pop().toLowerCase();
    };
    const validateFile = () => {
        if (!file) {
            setError('Please select a file');
            return false;
        }
        // Validate file type
        if (type === 'image' && !file.type.startsWith('image/')) {
            setError('Please select a valid image file');
            return false;
        }
        if (type === 'video' && !file.type.startsWith('video/')) {
            setError('Please select a valid video file');
            return false;
        }
        // Validate file size (10MB max)
        const maxSize = 10 * 1024 * 1024; // 10MB
        if (file.size > maxSize) {
            setError(`File size exceeds 10MB limit (${(file.size / (1024 * 1024)).toFixed(2)}MB)`);
            return false;
        }
        return true;
    };
    const handleUpload = async () => {
        if (!validateFile()) return;
        // Check if mediaId exists
        if (!mediaId) {
            setError('No media container available. Please save your content first before uploading media.');
            return;
        }
        setUploading(true);
        setProgress(0);
        try {
            const fileExt = getFileExtension(file.name);
            // Request presigned URL from backend
            let apiResponse;
            if (type === 'image') {
                apiResponse = await addImage({
                    mediaId: mediaId,
                    displayOrder: 1,
                    isMain: false,
                    fileType: fileExt,
                    userId: user?.userId
                });
            } else {
                apiResponse = await addVideo({
                    mediaId: mediaId,
                    displayOrder: 1,
                    isMain: false,
                    fileType: fileExt,
                    userId: user?.userId,
                });
            }
            const { url: presignedUrl, viewUrl } = apiResponse;
            // Upload to S3
            await axios.put(presignedUrl, file, {
                headers: {
                    'Content-Type': file.type,
                    'x-amz-acl': 'private',
                },
                onUploadProgress: (progressEvent) => {
                    const percentCompleted = Math.round(
                        (progressEvent.loaded * 100) / progressEvent.total
                    );
                    setProgress(percentCompleted);
                },
            });
            // Return the viewUrl to be inserted in the editor
            onUploadSuccess(viewUrl || presignedUrl);
            onClose();
        } catch (err) {
            setError('An error occurred during upload. Please try again.');
        } finally {
            setUploading(false);
        }
    };
    if (!isOpen) return null;
    return (
        <div className={styles.modalOverlay} role="dialog" aria-modal="true" aria-labelledby="media-upload-title">
            <div className={styles.modalContainer} ref={modalRef}>
                <div className={styles.modalHeader}>
                    <h3 id="media-upload-title">{type === 'image' ? 'Insert Image' : 'Insert Video'}</h3>
                    <button
                        className={styles.closeButton}
                        onClick={onClose}
                        aria-label="Close dialog"
                    >
                        <FaTimesCircle />
                    </button>
                </div>
                <div className={styles.modalBody}>
                    {!mediaId ? (
                        <ErrorMessage message="Media uploads are only available after saving your content. Please save the listing first." />
                    ) : (
                        <FileDropzone
                            onDrop={handleFileChange}
                            accept={type === 'image' ? ['image'] : ['video']}
                        >
                            <div className={styles.fileInputContainer}>
                                <label className={styles.fileInputLabel}>
                                    <input
                                        type="file"
                                        onChange={(e) => handleFileChange(e)}
                                        accept={type === 'image' ? 'image/*' : 'video/*'}
                                        className={styles.fileInput}
                                        aria-label={`Select ${type} file`}
                                    />
                                    <div className={styles.fileInputButton}>
                                        {type === 'image' ? <Image className={styles.fileInputIcon} /> :
                                            <Video className={styles.fileInputIcon} />}
                                        <span>Select {type === 'image' ? 'an image' : 'a video'}</span>
                                    </div>
                                </label>
                                <p className={styles.dropHint}>or drag and drop here</p>
                            </div>
                        </FileDropzone>
                    )}
                    <FilePreview
                        file={file}
                        previewUrl={previewUrl}
                        type={type}
                    />
                    {uploading && <ProgressIndicator progress={progress} />}
                    <ErrorMessage message={error} />
                    <div className={styles.modalActions}>
                        <button
                            className={styles.cancelButton}
                            onClick={onClose}
                            type="button"
                            disabled={uploading}
                        >
                            Cancel
                        </button>
                        <button
                            className={styles.uploadButton}
                            onClick={handleUpload}
                            disabled={!file || uploading || !mediaId}
                            type="button"
                        >
                            {uploading ? (
                                <>
                                    <span className={styles.spinnerIcon} aria-hidden="true" />
                                    <span>Uploading...</span>
                                </>
                            ) : (
                                <>
                                    <Upload /> Upload
                                </>
                            )}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
MediaUploadModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onUploadSuccess: PropTypes.func.isRequired,
    mediaId: PropTypes.string,
    type: PropTypes.oneOf(['image', 'video']).isRequired
};
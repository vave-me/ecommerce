"use client"
import { Upload } from '@/icons';
import { FaCheckCircle, FaExclamationCircle } from '../../utils/iconImports';
import React, { useState, memo } from 'react';
import axios from 'axios';
import { addImage } from "../../api/client/mediaApi";
import { useAuth } from "../../context/AuthContext";
import styles from './ImageUploader.module.css';
function getFileExtension(filename) {
    const parts = filename.split('.');
    return parts.length > 1 ? parts.pop().toLowerCase() : '';
}
const ImageUploaderHookForm = memo(({ mediaId, onUploadSuccess }) => {
    const [selectedFile, setSelectedFile] = useState(null);
    const [previewSrc, setPreviewSrc] = useState(null);
    const [uploadProgress, setUploadProgress] = useState(0);
    const [error, setError] = useState('');
    const [successMsg, setSuccessMsg] = useState('');
    const user = useAuth();
    const handleFileChange = (e) => {
        const file = e.target.files?.[0];
        setSelectedFile(file);
        setError('');
        setSuccessMsg('');
        setUploadProgress(0);
        setPreviewSrc(null);
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => setPreviewSrc(reader.result);
            reader.onerror = () => {
                setError('Failed to read the selected file.');
                setPreviewSrc(null);
            };
            reader.readAsDataURL(file);
        }
    };
    const handleUpload = async () => {
        if (!selectedFile) {
            setError('Please select an image.');
            return;
        }
        setError('');
        setSuccessMsg('');
        setUploadProgress(0);
        try {
            const fileExt = getFileExtension(selectedFile.name);
            // 1) Request presigned URL from your backend
            const resp = await addImage({
                mediaId: mediaId,
                displayOrder: 1,
                isMain: false,
                fileType: fileExt,
                userId: user.userId
            });
            const { url: presignedUrl, viewUrl } = resp; // rename finalUrl => viewUrl
            if (!presignedUrl) {
                throw new Error('No presigned URL returned from server.');
            }
            // 2) PUT to S3
            const uploadResp = await axios.put(presignedUrl, selectedFile, {
                headers: {
                    'Content-Type': selectedFile.type,
                    'x-amz-acl': 'private',
                },
                onUploadProgress: (progressEvent) => {
                    const percentCompleted = Math.round(
                        (progressEvent.loaded * 100) / progressEvent.total
                    );
                    setUploadProgress(percentCompleted);
                },
            });
            if (uploadResp.status >= 200 && uploadResp.status < 300) {
                setSuccessMsg('Image uploaded successfully!');
                if (onUploadSuccess) {
                    onUploadSuccess(viewUrl || presignedUrl);
                }
                // Clear local states
                setSelectedFile(null);
                setPreviewSrc(null);
                setUploadProgress(0);
            } else {
                setError('Upload failed. Please try again.');
            }
        } catch (err) {
            setError('An error occurred during upload.');
        }
    };
    return (
        <div className={styles.container}>
            <div className={styles.label}>Select Image:</div>
            <input
                className={styles.fileInput}
                type="file"
                accept="image/*"
                onChange={handleFileChange}
            />
            {previewSrc && <img className={styles.previewImage} src={previewSrc} alt="Preview" />}
            {uploadProgress > 0 && (
                <div className={styles.progressBarContainer}>
                    <progress className={styles.progressBar} value={uploadProgress} max="100" />
                    <span className={styles.progressText}>{uploadProgress}%</span>
                </div>
            )}
            {successMsg && (
                <div className={styles.successMsg}>
                    <FaCheckCircle style={{ marginRight: '5px' }} />
                    {successMsg}
                </div>
            )}
            {error && (
                <div className={styles.errorMsg}>
                    <FaExclamationCircle style={{ marginRight: '5px' }} />
                    {error}
                </div>
            )}
            <button className={styles.uploadButton} onClick={handleUpload}>
                <Upload style={{ marginRight: '5px' }} />
                {uploadProgress > 0 ? 'Uploading...' : 'Upload Image'}
            </button>
        </div>
    );
});
ImageUploaderHookForm.displayName = 'ImageUploaderHookForm';
export default ImageUploaderHookForm;
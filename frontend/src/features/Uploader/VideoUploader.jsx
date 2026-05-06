"use client"
import React, { useState, memo } from 'react';
import axios from 'axios';
import { Upload } from '@/icons';
import { FaCheckCircle, FaExclamationCircle } from '../../utils/iconImports';
// Optional: FilePond + plugins if you want drag-drop
import { FilePond, registerPlugin } from 'react-filepond';
import 'filepond/dist/filepond.min.css';
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';
import { addVideo } from "../../api/client/mediaApi";
import { useAuth } from "../../context/AuthContext";
import styles from './VideoUploader.module.css';
registerPlugin(FilePondPluginFileValidateType, FilePondPluginFileValidateSize);
function getFileExtension(filename) {
    return filename.split('.').pop();
}
const VideoUploaderHookForm = memo(({ mediaId, onUploadSuccess }) => {
    const [videoFile, setVideoFile] = useState(null);
    const [uploadProgress, setUploadProgress] = useState(0);
    const [uploadError, setUploadError] = useState('');
    const [uploadSuccess, setUploadSuccess] = useState('');
    const [isConverting, setIsConverting] = useState(false);
    const user = useAuth();
    const handleSubmitVideo = async () => {
        setUploadError('');
        setUploadProgress(0);
        setUploadSuccess('');
        if (!videoFile) {
            setUploadError('Please select a video file to upload.');
            return;
        }
        try {
            let file = videoFile;
            // 1) Ask your server for a presigned URL
            const fileExt = getFileExtension(file.name).toLowerCase();
            const response = await addVideo({
                mediaId: mediaId,
                displayOrder: 1,
                isMain: false,
                fileType: fileExt,
                userId: user.userId,
            });
            const { url: presignedUrl, viewUrl } = response;
            // 2) Upload to S3
            const uploadResp = await axios.put(presignedUrl, file, {
                headers: {
                    'Content-Type': 'video/mp4',
                    'x-amz-acl': 'private',
                },
                onUploadProgress: (progressEvent) => {
                    const percent = Math.round(
                        (progressEvent.loaded * 100) / progressEvent.total
                    );
                    setUploadProgress(percent);
                },
            });
            if (uploadResp.status === 200) {
                setUploadSuccess('Video uploaded successfully!');
                if (onUploadSuccess) {
                    onUploadSuccess(viewUrl || presignedUrl);
                }
                setVideoFile(null);
                setUploadProgress(0);
            } else {
                setUploadError('Upload failed. Please try again.');
            }
        } catch (error) {
            setUploadError('An error occurred during upload.');
        } finally {
            setIsConverting(false);
        }
    };
    return (
        <div className={styles.container}>
            <h3 className={styles.title}>Video Uploader</h3>
            {isConverting && (
                <p className={styles.convertingNotice}>Converting to MP4...</p>
            )}
            {/* Using FilePond or a simple <input type="file" /> */}
            <FilePond
                files={videoFile ? [videoFile] : []}
                allowMultiple={false}
                acceptedFileTypes={['video/*']}
                maxFileSize="500MB"
                onupdatefiles={(items) => {
                    setVideoFile(items.length > 0 ? items[0].file : null);
                    setUploadError('');
                    setUploadSuccess('');
                }}
                labelIdle='Drag & Drop or <span class="filepond--label-action">Browse</span> video'
            />
            {uploadProgress > 0 && (
                <div className={styles.progressRow}>
                    <progress className={styles.progressBar} value={uploadProgress} max="100" />
                    <span className={styles.progressText}>{uploadProgress}%</span>
                </div>
            )}
            {uploadSuccess && (
                <div className={styles.successNotice}>
                    <FaCheckCircle style={{ marginRight: '5px' }} />
                    {uploadSuccess}
                </div>
            )}
            {uploadError && (
                <div className={styles.errorNotice}>
                    <FaExclamationCircle style={{ marginRight: '5px' }} />
                    {uploadError}
                </div>
            )}
            <button
                className={styles.uploadButton}
                disabled={isConverting || uploadProgress > 0}
                onClick={handleSubmitVideo}
            >
                <Upload style={{ marginRight: '5px' }} />
                {uploadProgress > 0 ? 'Uploading...' : isConverting ? 'Converting...' : 'Upload Video'}
            </button>
        </div>
    );
});
VideoUploaderHookForm.displayName = 'VideoUploaderHookForm';
export default VideoUploaderHookForm;
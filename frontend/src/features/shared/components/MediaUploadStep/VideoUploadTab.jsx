// CreatePostModal/components/steps/MediaUploadStep/VideoUploadTab.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {X} from "@/icons";
import styles from './MediaUploadStep.module.css'
import VideoUploaderHookForm from "@/features/Uploader/VideoUploader";
export function VideoUploadTab({
                                   mediaId,
                                   uploadedVideoUrl,
                                   onUploadSuccess,
                                   onRemoveVideo
                               }) {
    return (
        <>
            <VideoUploaderHookForm
                mediaId={mediaId}
                onUploadSuccess={onUploadSuccess}
            />
            {uploadedVideoUrl && (
                <div className={styles.videoPreviewContainer}>
                    <h3 className={styles.sectionTitle}>
                        Uploaded Video
                    </h3>
                    <div className={styles.videoWrapper}>
                        <video
                            src={uploadedVideoUrl}
                            controls
                            className={styles.videoPreview}
                        />
                        <button
                            className={styles.removeButton}
                            onClick={onRemoveVideo}
                            aria-label="Remove video"
                            type="button"
                        >
                            <X size={12}/>
                        </button>
                    </div>
                </div>
            )}
        </>
    );
}
VideoUploadTab.propTypes = {
    mediaId: PropTypes.string.isRequired,
    uploadedVideoUrl: PropTypes.string,
    onUploadSuccess: PropTypes.func.isRequired,
    onRemoveVideo: PropTypes.func.isRequired
};
// CreatePostModal/components/steps/MediaUploadStep/ImageUploadTab.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {X} from "@/icons";
import styles from './MediaUploadStep.module.css'
import ImageUploaderHookForm from "@/features/Uploader/ImageUploader";
export function ImageUploadTab({
                                   mediaId,
                                   uploadedImages,
                                   onUploadSuccess,
                                   onRemoveImage
                               }) {
    return (
        <>
            <ImageUploaderHookForm
                mediaId={mediaId}
                onUploadSuccess={onUploadSuccess}
            />
            {uploadedImages.length > 0 && (
                <div className={styles.thumbnailSection}>
                    <h3 className={styles.sectionTitle}>
                        Uploaded Images
                    </h3>
                    <div className={styles.imageGrid}>
                        {uploadedImages.map((url, idx) => (
                            <div
                                key={idx}
                                className={styles.imagePreview}
                            >
                                <img
                                    src={url}
                                    alt={`Image ${idx + 1}`}
                                    className={styles.previewImg}
                                />
                                <button
                                    className={styles.removeButton}
                                    onClick={() => onRemoveImage(idx)}
                                    aria-label="Remove image"
                                    type="button"
                                >
                                    <X size={12}/>
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </>
    );
}
ImageUploadTab.propTypes = {
    mediaId: PropTypes.string.isRequired,
    uploadedImages: PropTypes.array.isRequired,
    onUploadSuccess: PropTypes.func.isRequired,
    onRemoveImage: PropTypes.func.isRequired
};
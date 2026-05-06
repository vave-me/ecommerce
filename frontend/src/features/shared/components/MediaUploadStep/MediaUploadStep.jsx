import React, {useState} from 'react';
import PropTypes from 'prop-types';
import {ImageUploadTab} from './ImageUploadTab';
import {VideoUploadTab} from './VideoUploadTab';
import {FormActions} from "@/common/components/FormActions";
import styles from './MediaUploadStep.module.css';
/**
 * Shared MediaUploadStep Component
 * Used across all creation modals for consistent media upload functionality
 */
export function MediaUploadStep({
                                    mediaId,
                                    postId,
                                    initialImages = [],
                                    initialVideoUrl = '',
                                    onComplete,
                                    onBack,
                                    isLoading = false,
                                    errors = {},
                                    title = "Media Upload",
                                    description = "Upload photos and videos to enhance your post. High-quality media helps increase engagement.",
                                    continueLabel = "Continue",
                                    helpText = "Adding relevant images or videos will make your post more engaging and increase visibility."
                                }) {
    const [activeTab, setActiveTab] = useState('images');
    const [uploadedImages, setUploadedImages] = useState(initialImages);
    const [uploadedVideoUrl, setUploadedVideoUrl] = useState(initialVideoUrl);
    const handleRemoveImage = (idx) => {
        setUploadedImages((prev) => {
            const updated = [...prev];
            updated.splice(idx, 1);
            return updated;
        });
    };
    const handleDone = () => {
        onComplete({
            images: uploadedImages,
            videoUrl: uploadedVideoUrl,
            thumbnail: uploadedImages.length > 0 ? uploadedImages[0] : ''
        });
    };
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>{title}</h2>
            <p className={styles.formDescription}>
                {description}
            </p>
            <div className={styles.tabContainer}>
                <button
                    className={`${styles.tabButton} ${activeTab === 'images' ? styles.activeTab : ''}`}
                    onClick={() => setActiveTab('images')}
                    type="button"
                >
                    Images
                </button>
                <button
                    className={`${styles.tabButton} ${activeTab === 'videos' ? styles.activeTab : ''}`}
                    onClick={() => setActiveTab('videos')}
                    type="button"
                >
                    Videos
                </button>
            </div>
            <div className={styles.mediaContent}>
                {!mediaId ? (
                    <div className={styles.loadingContainer}>
                        <div className={styles.loadingText}>
                            <span>Preparing upload...</span>
                        </div>
                    </div>
                ) : (
                    <>
                        {activeTab === 'images' && (
                            <ImageUploadTab
                                mediaId={mediaId}
                                uploadedImages={uploadedImages}
                                onUploadSuccess={(url) => setUploadedImages((prev) => [...prev, url])}
                                onRemoveImage={handleRemoveImage}
                            />
                        )}
                        {activeTab === 'videos' && (
                            <VideoUploadTab
                                mediaId={mediaId}
                                uploadedVideoUrl={uploadedVideoUrl}
                                onUploadSuccess={setUploadedVideoUrl}
                                onRemoveVideo={() => setUploadedVideoUrl(null)}
                            />
                        )}
                    </>
                )}
            </div>
            <div className={styles.mediaHelpText}>
                <p>
                    <strong>Tip:</strong> {helpText}
                </p>
            </div>
            <FormActions
                primaryLabel={isLoading ? "Processing..." : continueLabel}
                primaryIcon="arrow-right"
                secondaryLabel="Back"
                onCancel={onBack}
                onPrimaryAction={handleDone}
                isSubmitting={isLoading}
            />
        </div>
    );
}
MediaUploadStep.propTypes = {
    mediaId: PropTypes.string,
    postId: PropTypes.string,
    initialImages: PropTypes.array,
    initialVideoUrl: PropTypes.string,
    onComplete: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    errors: PropTypes.object,
    title: PropTypes.string,
    description: PropTypes.string,
    continueLabel: PropTypes.string,
    helpText: PropTypes.string
}; 
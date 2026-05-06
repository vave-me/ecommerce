// TextEditor/components/FileUpload/FilePreview.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function FilePreview({ file, previewUrl, type }) {
    if (!file) return null;
    if (type === 'image' && previewUrl) {
        return (
            <div className={styles.previewContainer}>
                <img src={previewUrl} alt="Preview" className={styles.preview} />
            </div>
        );
    } else if (type === 'video') {
        return (
            <div className={styles.fileInfoContainer}>
                <div className={styles.fileInfo}>
                    <span className={styles.fileName}>{file.name}</span>
                    <span className={styles.fileSize}>({(file.size / (1024 * 1024)).toFixed(2)} MB)</span>
                </div>
            </div>
        );
    }
    return null;
}
FilePreview.propTypes = {
    file: PropTypes.object,
    previewUrl: PropTypes.string,
    type: PropTypes.oneOf(['image', 'video']).isRequired
};
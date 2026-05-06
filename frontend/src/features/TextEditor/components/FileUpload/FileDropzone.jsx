// TextEditor/components/FileUpload/FileDropzone.jsx
import React, { useState, useEffect, useCallback, useRef } from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function FileDropzone({ onDrop, accept, children }) {
    const [isDragging, setIsDragging] = useState(false);
    const dropzoneRef = useRef(null);
    const handleDragOver = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(true);
    }, []);
    const handleDragLeave = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);
    }, []);
    const handleDrop = useCallback((e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);
        if (!e.dataTransfer.files || e.dataTransfer.files.length === 0) return;
        const file = e.dataTransfer.files[0];
        // Check if file type is accepted
        if (accept && !accept.includes(file.type.split('/')[0])) return;
        onDrop(file);
    }, [onDrop, accept]);
    useEffect(() => {
        const element = dropzoneRef.current;
        if (element) {
            element.addEventListener('dragover', handleDragOver);
            element.addEventListener('dragleave', handleDragLeave);
            element.addEventListener('drop', handleDrop);
            return () => {
                element.removeEventListener('dragover', handleDragOver);
                element.removeEventListener('dragleave', handleDragLeave);
                element.removeEventListener('drop', handleDrop);
            };
        }
    }, [handleDragOver, handleDragLeave, handleDrop]);
    return (
        <div
            ref={dropzoneRef}
            className={`${styles.dropzone} ${isDragging ? styles.dragging : ''}`}
        >
            {children}
            {isDragging && (
                <div className={styles.dropOverlay}>
                    <div className={styles.dropMessage}>
                        Drop to upload
                    </div>
                </div>
            )}
        </div>
    );
}
FileDropzone.propTypes = {
    onDrop: PropTypes.func.isRequired,
    accept: PropTypes.arrayOf(PropTypes.string),
    children: PropTypes.node.isRequired
};
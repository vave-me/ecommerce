// TextEditor/utils/mediaHelpers.js
export function getFileExtension(filename) {
    if (!filename) return '';
    return filename.split('.').pop().toLowerCase();
}
export function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}
export function validateFile(file, type, maxSizeMB = 10) {
    if (!file) {
        return { isValid: false, error: 'No file selected' };
    }
    // Validate file type
    if (type === 'image' && !file.type.startsWith('image/')) {
        return { isValid: false, error: 'Please select a valid image file' };
    }
    if (type === 'video' && !file.type.startsWith('video/')) {
        return { isValid: false, error: 'Please select a valid video file' };
    }
    // Validate file size
    const maxSize = maxSizeMB * 1024 * 1024;
    if (file.size > maxSize) {
        return {
            isValid: false,
            error: `File size exceeds ${maxSizeMB}MB limit (${formatFileSize(file.size)})`
        };
    }
    return { isValid: true, error: null };
}
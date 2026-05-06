// File: src/components/AttachmentList.jsx
import { FaPaperclip } from '../../utils/iconImports';
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import Image from 'next/image';
// Import the CSS module
import styles from './AttachmentList.module.css';
const AttachmentList = ({ attachments, onRemove }) => (
    <div className={styles.attachments}>
        {attachments.map((file, index) => (
            <div key={file.id} className={styles.attachment}>
                <div className={styles.icon}>
                    {file.type.startsWith('image/') ? (
                        <Image 
                            src={file.preview} 
                            alt={`Attachment ${index + 1}`}
                            width={40}
                            height={40}
                            style={{ objectFit: 'cover' }}
                        />
                    ) : (
                        <FaPaperclip />
                    )}
                </div>
                <span className={styles.name} title={file.name}>
          {file.name}
        </span>
                <button
                    type="button"
                    className={styles.removeButton}
                    onClick={() => onRemove(file.id)}
                    aria-label={`Remove attachment ${file.name}`}
                >
                    &times;
                </button>
            </div>
        ))}
    </div>
);
AttachmentList.propTypes = {
    attachments: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            name: PropTypes.string.isRequired,
            type: PropTypes.string.isRequired,
            preview: PropTypes.string,
        })
    ).isRequired,
    onRemove: PropTypes.func.isRequired,
};
export default memo(AttachmentList);

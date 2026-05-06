// TextEditor/components/MenuBar/LinkInput.jsx
import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { ExternalLink } from '@/icons';
import { FaUnlink } from '../../../../utils/iconImports';
import styles from '../../TextEditor.module.css';
export function LinkInput({ editor }) {
    const [url, setUrl] = useState('');
    const [showInput, setShowInput] = useState(false);
    const inputRef = useRef(null);
    const handleSetLink = () => {
        if (!url) return;
        // Update link
        editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
        setUrl('');
        setShowInput(false);
    };
    const handleRemoveLink = () => {
        editor.chain().focus().extendMarkRange('link').unsetLink().run();
        setShowInput(false);
    };
    useEffect(() => {
        if (showInput && inputRef.current) {
            inputRef.current.focus();
        }
    }, [showInput]);
    return (
        <div className={styles.linkInputWrapper}>
            <button
                type="button"
                onClick={() => setShowInput(!showInput)}
                // TextEditor/components/MenuBar/LinkInput.jsx (continued)
                className={editor.isActive('link') ? styles.isActive : ''}
                title="Insert Link"
                aria-label="Insert Link"
            >
                <ExternalLink/>
            </button>
            {editor.isActive('link') && (
                <button
                    type="button"
                    onClick={handleRemoveLink}
                    title="Remove Link"
                    aria-label="Remove Link"
                    className={styles.unlinkButton}
                >
                    <FaUnlink/>
                </button>
            )}
            {showInput && (
                <div className={styles.linkInputContainer}>
                    <input
                        ref={inputRef}
                        type="url"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="https://example.com"
                        className={styles.linkInput}
                        aria-label="Link URL"
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                e.preventDefault();
                                handleSetLink();
                            }
                            if (e.key === 'Escape') {
                                e.preventDefault();
                                setShowInput(false);
                            }
                        }}
                    />
                    <button
                        type="button"
                        onClick={handleSetLink}
                        className={styles.linkSubmitButton}
                        aria-label="Confirm Link"
                    >
                        Add
                    </button>
                </div>
            )}
        </div>
    );
}
LinkInput.propTypes = {
    editor: PropTypes.object.isRequired
};
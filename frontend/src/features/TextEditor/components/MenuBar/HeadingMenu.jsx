// TextEditor/components/MenuBar/HeadingMenu.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function HeadingMenu({ editor, isOpen, onClose }) {
    if (!isOpen) return null;
    return (
        <div className={styles.headingMenu}>
            <button
                onClick={() => {
                    editor.chain().focus().toggleHeading({ level: 1 }).run();
                    onClose();
                }}
                className={editor.isActive('heading', { level: 1 }) ? styles.isActive : ''}
                type="button"
            >
                Heading 1
            </button>
            <button
                onClick={() => {
                    editor.chain().focus().toggleHeading({ level: 2 }).run();
                    onClose();
                }}
                className={editor.isActive('heading', { level: 2 }) ? styles.isActive : ''}
                type="button"
            >
                Heading 2
            </button>
            <button
                onClick={() => {
                    editor.chain().focus().toggleHeading({ level: 3 }).run();
                    onClose();
                }}
                className={editor.isActive('heading', { level: 3 }) ? styles.isActive : ''}
                type="button"
            >
                Heading 3
            </button>
            <button
                onClick={() => {
                    editor.chain().focus().setParagraph().run();
                    onClose();
                }}
                className={editor.isActive('paragraph') ? styles.isActive : ''}
                type="button"
            >
                Paragraph
            </button>
        </div>
    );
}
HeadingMenu.propTypes = {
    editor: PropTypes.object,
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired
};
// TextEditor/components/MenuBar/MenuBar.jsx
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Image, Video } from '@/icons';
import { FaBold, FaItalic, FaUnderline, FaListUl, FaListOl, FaQuoteRight, FaAlignLeft, FaAlignCenter, FaAlignRight, FaTable, FaHeading, FaUndo, FaRedo, FaBars } from '../../../../utils/iconImports';
import styles from '../../TextEditor.module.css';
import { ButtonGroup } from './ButtonGroup';
import { HeadingMenu } from './HeadingMenu';
import { LinkInput } from './LinkInput';
export function MenuBar({ editor, onImageClick, onVideoClick, mediaId, isMobile }) {
    const [isExpanded, setIsExpanded] = useState(true);
    const [showHeadingMenu, setShowHeadingMenu] = useState(false);
    // Close heading menu when clicking outside
    useEffect(() => {
        if (showHeadingMenu) {
            const handleOutsideClick = (e) => {
                if (!e.target.closest(`.${styles.headingMenu}`)) {
                    setShowHeadingMenu(false);
                }
            };
            document.addEventListener('click', handleOutsideClick);
            return () => document.removeEventListener('click', handleOutsideClick);
        }
    }, [showHeadingMenu]);
    if (!editor) {
        return null;
    }
    return (
        <div className={styles.menuBarWrapper}>
            {isMobile && (
                <button
                    className={styles.menuToggle}
                    onClick={() => setIsExpanded(!isExpanded)}
                    aria-expanded={isExpanded}
                    aria-label="Toggle formatting toolbar"
                >
                    <FaBars/>
                </button>
            )}
            <div className={`${styles.menuBar} ${!isExpanded ? styles.menuBarCollapsed : ''}`}>
                {/* Basic formatting */}
                <ButtonGroup>
                    <button
                        onClick={() => editor.chain().focus().toggleBold().run()}
                        className={editor.isActive('bold') ? styles.isActive : ''}
                        type="button"
                        title="Bold"
                        aria-label="Bold"
                    >
                        <FaBold/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleItalic().run()}
                        className={editor.isActive('italic') ? styles.isActive : ''}
                        type="button"
                        title="Italic"
                        aria-label="Italic"
                    >
                        <FaItalic/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleUnderline().run()}
                        className={editor.isActive('underline') ? styles.isActive : ''}
                        type="button"
                        title="Underline"
                        aria-label="Underline"
                    >
                        <FaUnderline/>
                    </button>
                </ButtonGroup>
                {/* Heading controls */}
                <ButtonGroup>
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            setShowHeadingMenu(!showHeadingMenu);
                        }}
                        className={editor.isActive('heading') ? styles.isActive : ''}
                        type="button"
                        title="Heading Style"
                        aria-label="Heading Style"
                        aria-haspopup="true"
                        aria-expanded={showHeadingMenu}
                    >
                        <FaHeading/>
                    </button>
                    <HeadingMenu
                        editor={editor}
                        isOpen={showHeadingMenu}
                        onClose={() => setShowHeadingMenu(false)}
                    />
                </ButtonGroup>
                {/* Lists */}
                <ButtonGroup>
                    <button
                        onClick={() => editor.chain().focus().toggleBulletList().run()}
                        className={editor.isActive('bulletList') ? styles.isActive : ''}
                        type="button"
                        title="Bullet List"
                        aria-label="Bullet List"
                    >
                        <FaListUl/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleOrderedList().run()}
                        className={editor.isActive('orderedList') ? styles.isActive : ''}
                        type="button"
                        title="Ordered List"
                        aria-label="Ordered List"
                    >
                        <FaListOl/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().toggleBlockquote().run()}
                        className={editor.isActive('blockquote') ? styles.isActive : ''}
                        type="button"
                        title="Blockquote"
                        aria-label="Blockquote"
                    >
                        <FaQuoteRight/>
                    </button>
                </ButtonGroup>
                {/* Alignment */}
                <ButtonGroup>
                    <button
                        onClick={() => editor.chain().focus().setTextAlign('left').run()}
                        className={editor.isActive({textAlign: 'left'}) ? styles.isActive : ''}
                        type="button"
                        title="Align Left"
                        aria-label="Align Left"
                    >
                        <FaAlignLeft/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().setTextAlign('center').run()}
                        className={editor.isActive({textAlign: 'center'}) ? styles.isActive : ''}
                        type="button"
                        title="Align Center"
                        aria-label="Align Center"
                    >
                        <FaAlignCenter/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().setTextAlign('right').run()}
                        className={editor.isActive({textAlign: 'right'}) ? styles.isActive : ''}
                        type="button"
                        title="Align Right"
                        aria-label="Align Right"
                    >
                        <FaAlignRight/>
                    </button>
                </ButtonGroup>
                {/* Link */}
                <ButtonGroup>
                    <LinkInput editor={editor}/>
                </ButtonGroup>
                {/* Table */}
                <ButtonGroup>
                    <button
                        onClick={() => editor.chain().focus().insertTable({
                            rows: 3,
                            cols: 3,
                            withHeaderRow: true
                        }).run()}
                        type="button"
                        title="Insert Table"
                        aria-label="Insert Table"
                        disabled={editor.isActive('table')}
                    >
                        <FaTable/>
                    </button>
                </ButtonGroup>
                {/* Undo/Redo */}
                <ButtonGroup>
                    <button
                        onClick={() => editor.chain().focus().undo().run()}
                        type="button"
                        title="Undo"
                        aria-label="Undo"
                        disabled={!editor.can().undo()}
                    >
                        <FaUndo/>
                    </button>
                    <button
                        onClick={() => editor.chain().focus().redo().run()}
                        type="button"
                        title="Redo"
                        aria-label="Redo"
                        disabled={!editor.can().redo()}
                    >
                        <FaRedo/>
                    </button>
                </ButtonGroup>
                {/* Media buttons */}
                <ButtonGroup>
                    <button
                        onClick={onImageClick}
                        className={`${styles.mediaButton} ${!mediaId ? styles.mediaButtonDisabled : ''}`}
                        type="button"
                        title={mediaId ? "Insert Image" : "Save content first to enable image upload"}
                        aria-label={mediaId ? "Insert Image" : "Save content first to enable image upload"}
                        disabled={!mediaId}
                    >
                        <Image/>
                    </button>
                    <button
                        onClick={onVideoClick}
                        className={`${styles.mediaButton} ${!mediaId ? styles.mediaButtonDisabled : ''}`}
                        type="button"
                        title={mediaId ? "Insert Video" : "Save content first to enable video upload"}
                        aria-label={mediaId ? "Insert Video" : "Save content first to enable video upload"}
                        disabled={!mediaId}
                    >
                        <Video/>
                    </button>
                </ButtonGroup>
            </div>
        </div>
    );
}
MenuBar.propTypes = {
    editor: PropTypes.object,
    onImageClick: PropTypes.func.isRequired,
    onVideoClick: PropTypes.func.isRequired,
    mediaId: PropTypes.string,
    isMobile: PropTypes.bool
};
// TextEditor/index.jsx
import React, { useState, useEffect, useRef, memo } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Underline from '@tiptap/extension-underline';
import Link from '@tiptap/extension-link';
import Placeholder from '@tiptap/extension-placeholder';
import Image from '@tiptap/extension-image';
import TextAlign from '@tiptap/extension-text-align';
import Table from '@tiptap/extension-table';
import TableRow from '@tiptap/extension-table-row';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import styles from './TextEditor.module.css';
import { MenuBar } from './components/MenuBar/MenuBar';
import { MediaUploadModal } from './components/FileUpload/MediaUploadModal';
import {useIsMobile} from "../../hooks/useMobileDetection";
const TextEditor = memo(({ value, onChange, placeholder, mediaId, id }) => {
    const [isMounted, setIsMounted] = useState(false);
    const [showImageModal, setShowImageModal] = useState(false);
    const [showVideoModal, setShowVideoModal] = useState(false);
    const editorRef = useRef(null);
    const isMobile = useIsMobile();
    // Initialize the editor with enhanced features
    const editor = useEditor({
        extensions: [
            StarterKit,
            Underline,
            Link.configure({
                openOnClick: false,
                validate: href => /^https?:\/\//.test(href),
                HTMLAttributes: {
                    rel: 'noopener noreferrer nofollow',
                    class: styles.editorLink,
                },
            }),
            Placeholder.configure({
                placeholder: placeholder || 'Write something...',
                emptyEditorClass: styles.emptyEditor,
            }),
            Image.configure({
                allowBase64: true,
                inline: false,
                HTMLAttributes: {
                    class: styles.editorImage,
                },
            }),
            TextAlign.configure({
                types: ['heading', 'paragraph'],
            }),
            Table.configure({
                resizable: true,
                HTMLAttributes: {
                    class: styles.editorTable,
                },
            }),
            TableRow,
            TableHeader,
            TableCell,
        ],
        content: value || '',
        onUpdate: ({ editor }) => {
            const html = editor.getHTML();
            onChange(html);
        },
        editorProps: {
            attributes: {
                class: styles.editor,
                id: id || 'tiptap-editor',
                role: 'textbox',
                'aria-multiline': 'true',
                'aria-label': placeholder || 'Rich text editor',
            },
        },
    });
    // Handle client-side only rendering
    useEffect(() => {
        setIsMounted(true);
    }, []);
    // Keep focus in sync with external focus commands
    useEffect(() => {
        if (editor && editorRef.current) {
            const handleClick = () => {
                editor.commands.focus();
            };
            editorRef.current.addEventListener('click', handleClick);
            return () => {
                editorRef.current?.removeEventListener('click', handleClick);
            };
        }
    }, [editor, editorRef.current]);
    // Handlers for media insertion
    const handleImageInsert = (url) => {
        editor
            .chain()
            .focus()
            .setImage({ src: url, alt: 'Uploaded image' })
            .run();
    };
    const handleVideoInsert = (url) => {
        editor
            .chain()
            .focus()
            .insertContent(`<div class="video-embed"><video src="${url}" controls width="100%" preload="metadata"></video></div>`)
            .run();
    };
    // Handle media button click
    const handleMediaButtonClick = (type) => {
        if (!mediaId) {
            // TODO: Replace with toast notification
        if (typeof window !== 'undefined' && window.toast) {
            window.toast.warn(`Please save your content first before adding ${type}s.`);
        }
            return;
        }
        if (type === 'image') {
            setShowImageModal(true);
        } else if (type === 'video') {
            setShowVideoModal(true);
        }
    };
    // Show a loading state during SSR
    if (!isMounted) {
        return (
            <div
                className={styles.editorLoading}
                role="progressbar"
                aria-busy="true"
                aria-label="Loading editor"
            >
                <div className={styles.editorLoadingSpinner}></div>
                <span>Loading editor...</span>
            </div>
        );
    }
    return (
        <div className={styles.TextEditor} ref={editorRef}>
            <MenuBar
                editor={editor}
                onImageClick={() => handleMediaButtonClick('image')}
                onVideoClick={() => handleMediaButtonClick('video')}
                mediaId={mediaId}
                isMobile={isMobile}
            />
            <EditorContent editor={editor} className={styles.editorContent} />
            {/* Accessibility announcement for editor state */}
            <div className={styles.editorStatus} aria-live="polite">
                {editor?.isActive('bold') && <span className="sr-only">Bold is active. </span>}
                {editor?.isActive('italic') && <span className="sr-only">Italic is active. </span>}
                {editor?.isActive('underline') && <span className="sr-only">Underline is active. </span>}
            </div>
            {/* Media upload status */}
            {!mediaId && (
                <div className={styles.mediaStatusMessage}>
                    <p>Media uploads will be available after saving your content.</p>
                </div>
            )}
            {/* Image Upload Modal */}
            <MediaUploadModal
                isOpen={showImageModal}
                onClose={() => setShowImageModal(false)}
                onUploadSuccess={handleImageInsert}
                mediaId={mediaId}
                type="image"
            />
            {/* Video Upload Modal */}
            <MediaUploadModal
                isOpen={showVideoModal}
                onClose={() => setShowVideoModal(false)}
                onUploadSuccess={handleVideoInsert}
                mediaId={mediaId}
                type="video"
            />
        </div>
    );
});
TextEditor.displayName = 'TextEditor';
export default TextEditor;
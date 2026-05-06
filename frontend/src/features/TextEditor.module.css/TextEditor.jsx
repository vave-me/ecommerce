// src/features/AddPostForm/TextEditor.jsx
"use client";
import React, {useState, useEffect, memo} from 'react';
import {useEditor, EditorContent} from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Underline from '@tiptap/extension-underline';
import Link from '@tiptap/extension-link';
import Placeholder from '@tiptap/extension-placeholder';
// This is the toolbar component
const MenuBar = memo(({editor}) => {
    if (!editor) {
        return null;
    }
    return (
        <div className="menu-bar">
            <button
                onClick={() => editor.chain().focus().toggleBold().run()}
                className={editor.isActive('bold') ? 'is-active' : ''}
                type="button"
            >
                Bold
            </button>
            <button
                onClick={() => editor.chain().focus().toggleItalic().run()}
                className={editor.isActive('italic') ? 'is-active' : ''}
                type="button"
            >
                Italic
            </button>
            <button
                onClick={() => editor.chain().focus().toggleUnderline().run()}
                className={editor.isActive('underline') ? 'is-active' : ''}
                type="button"
            >
                Underline
            </button>
            <button
                onClick={() => editor.chain().focus().toggleHeading({level: 1}).run()}
                className={editor.isActive('heading', {level: 1}) ? 'is-active' : ''}
                type="button"
            >
                H1
            </button>
            <button
                onClick={() => editor.chain().focus().toggleHeading({level: 2}).run()}
                className={editor.isActive('heading', {level: 2}) ? 'is-active' : ''}
                type="button"
            >
                H2
            </button>
            <button
                onClick={() => editor.chain().focus().toggleHeading({level: 3}).run()}
                className={editor.isActive('heading', {level: 3}) ? 'is-active' : ''}
                type="button"
            >
                H3
            </button>
            <button
                onClick={() => editor.chain().focus().toggleBulletList().run()}
                className={editor.isActive('bulletList') ? 'is-active' : ''}
                type="button"
            >
                Bullet List
            </button>
            <button
                onClick={() => editor.chain().focus().toggleOrderedList().run()}
                className={editor.isActive('orderedList') ? 'is-active' : ''}
                type="button"
            >
                Ordered List
            </button>
            <button
                onClick={() => editor.chain().focus().toggleBlockquote().run()}
                className={editor.isActive('blockquote') ? 'is-active' : ''}
                type="button"
            >
                Blockquote
            </button>
        </div>
    );
});
MenuBar.displayName = 'MenuBar';
const TextEditor = memo(({value, onChange, placeholder}) => {
    const [isMounted, setIsMounted] = useState(false);
    // Initialize the editor
    const editor = useEditor({
        extensions: [
            StarterKit,
            Underline,
            Link,
            Placeholder.configure({
                placeholder: placeholder || 'Write something...',
            }),
        ],
        content: value || '',
        onUpdate: ({editor}) => {
            const html = editor.getHTML();
            onChange(html);
        },
    });
    // Handle client-side only rendering
    useEffect(() => {
        setIsMounted(true);
    }, []);
    // Show a loading state or nothing during SSR
    if (!isMounted) {
        return (
            <div
                style={{
                    height: '200px',
                    border: '1px solid #ccc',
                    borderRadius: '6px',
                    backgroundColor: '#f9f9f9',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                }}
            >
                Loading editor...
            </div>
        );
    }
    return (
        <div className="rich-text-editor">
            <MenuBar editor={editor}/>
            <EditorContent editor={editor}/>
        </div>
    );
});
TextEditor.displayName = 'TextEditor';
export default TextEditor;
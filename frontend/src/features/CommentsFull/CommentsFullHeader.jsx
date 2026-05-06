// File: src/comments/CommentsFullHeader.jsx
"use client";
import React, { memo } from 'react';
import PropTypes from "prop-types";
import { X } from '@/icons';
import Image from "next/image";
// Import the CSS module
import styles from "./CommentsFullHeader.module.css";
const CommentsFullHeader = memo(({ toggleCommentsFullList, productName, thumbnail }) => (
    <header className={styles.header} role="banner">
        <button
            type="button"
            className={styles.closeButton}
            onClick={toggleCommentsFullList}
            aria-label="Close CommentsFull"
            title="Close CommentsFull"
        >
            <X />
        </button>
        <h1 id="comments-heading" className={styles.heading}>
            {thumbnail && (
                <div className={styles.thumbnailWrapper}>
                    <Image 
                        src={thumbnail} 
                        alt="" 
                        width={24} 
                        height={24} 
                        style={{ objectFit: 'cover' }}
                    />
                </div>
            )}
            {productName} CommentsFull
        </h1>
    </header>
));
CommentsFullHeader.propTypes = {
    toggleCommentsFullList: PropTypes.func.isRequired,
    productName: PropTypes.string,
    thumbnail: PropTypes.string,
};
CommentsFullHeader.defaultProps = {
    productName: "",
    thumbnail: "",
};
export default CommentsFullHeader;

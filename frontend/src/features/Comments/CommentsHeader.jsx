// File: src/comments/CommentsHeader.jsx
"use client";
import React, { memo } from 'react';
import PropTypes from "prop-types";
import {X, MessageCircle} from "@/icons";
import Image from "next/image";
// Import the CSS module
import styles from "./Comments.module.css";
const CommentsHeader = memo(({toggleCommentsList, productName, thumbnail, commentCount = 0}) => (
    <header className={styles.header} role="banner">
        <div className={styles.headerContent}>
            <MessageCircle size={18} color="#1aa89e"/>
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
                {productName ? `${productName} Comments` : "Comments"}
                {commentCount > 0 && <span className={styles.commentCount}>({commentCount})</span>}
            </h1>
        </div>
        <button
            type="button"
            className={styles.closeButton}
            onClick={toggleCommentsList}
            aria-label="Close Comments"
            title="Close Comments"
        >
            <X size={18}/>
        </button>
    </header>
));
CommentsHeader.propTypes = {
    toggleCommentsList: PropTypes.func.isRequired,
    productName: PropTypes.string,
    thumbnail: PropTypes.string,
    commentCount: PropTypes.number,
};
CommentsHeader.defaultProps = {
    productName: "",
    thumbnail: "",
    commentCount: 0,
};
export default CommentsHeader;
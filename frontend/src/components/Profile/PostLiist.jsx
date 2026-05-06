// src/components/UserProfile/PostList.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import { getUserPosts } from "../../api/postsApi";
import PostListItem from "../Items/PostListItem";
import styles from './PostList.module.css';
const randIMG = [
    "/images/lenovo2.png",
    "/images/lenovo2.png",
    "/images/lenovo2.png",
    "/images/default-product.webp"
];
const PostList = memo(({ userId }) => {
    const [posts, setPosts] = useState([]);
    useEffect(() => {
        const fetchUserPosts = async () => {
            try {
                const data = await getUserPosts(userId);
                const userPosts = data.posts.map((post) => ({
                    id: post.id,
                    title: post.title,
                    content: post.content,
                    thumbnail: randIMG[3], // placeholder
                }));
                setPosts(userPosts);
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
        };
        fetchUserPosts();
    }, [userId]);
    return (
        <div className={styles.listContainer}>
            {posts.map((post) => (
                <PostListItem key={post.id} product={post} />
            ))}
        </div>
    );
});
PostList.displayName = 'PostList';
PostList.propTypes = {
    userId: PropTypes.string.isRequired,
};
export default PostList;

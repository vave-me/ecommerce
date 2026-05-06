// File : src/app/page.jsx
"use client";
import { useQuery } from '@tanstack/react-query';
import { useIsMobile } from "../../../hooks/useMobileDetection";
import { getPosts } from "../../../api/postsApi";
import ImprovedSocialCard from "../design/tweets/page";
import styles from './Social.module.css';
export default function HomePage() {
    const isMobile = useIsMobile();
    const {
        data: featuredData,
        isLoading,
        isError,
    } = useQuery({
        queryKey: ['featuredPosts'],
        queryFn: getPosts, // real function
    });
    if (isLoading) return <div>Loading homepage...</div>;
    if (isError) return <div>Failed to load posts.</div>;
    const posts = featuredData?.posts || [];
    return (
        <main className={styles.mainContent}>
            <div className={styles.container}>
                <ul>
                    {posts.map(post => (
                        <ImprovedSocialCard tweet={post} key={post.id}/>
                    ))}
                </ul>
            </div>
        </main>
    );
}
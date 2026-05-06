"use client";
import React, {memo} from 'react';
import Link from 'next/link';
import styles from './TrendingTags.module.css';
const TrendingTags = memo(() => {
    // Mock data for trending tags - in a real app, this would come from an API
    const trendingTags = [
        {id: 1, name: '#technology', count: 1250},
        {id: 2, name: '#design', count: 980},
        {id: 3, name: '#marketing', count: 750},
        {id: 4, name: '#development', count: 620},
        {id: 5, name: '#business', count: 580},
        {id: 6, name: '#startup', count: 450},
        {id: 7, name: '#innovation', count: 420},
        {id: 8, name: '#digital', count: 380},
        {id: 9, name: '#trending', count: 350},
        {id: 10, name: '#viral', count: 320}
    ];
    return (
        <div className={styles.trendingTags}>
            <h3 className={styles.title}>Trending Tags</h3>
            <div className={styles.tagsContainer}>
                {trendingTags.map(tag => (
                    <Link
                        key={tag.id}
                        href={`/search?q=${tag.name.slice(1)}`}
                        className={styles.tag}
                    >
                        <span className={styles.tagName}>{tag.name}</span>
                        <span className={styles.tagCount}>{tag.count.toLocaleString()}</span>
                    </Link>
                ))}
            </div>
        </div>
    );
});
TrendingTags.displayName = 'TrendingTags';
export default TrendingTags; 
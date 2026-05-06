"use client";
import React, { memo } from 'react';
import styles from './ProductCard.module.css';
import { ShoppingCart, Tag } from '@/icons';
import Link from 'next/link';

/**
 * Simplified Product Card for AI Results
 * A lightweight component to display products returned by the AI assistant
 */
const ProductCard = ({ product }) => {
    // Format price from cents to currency
    const formatPrice = (price) => {
        const numPrice = typeof price === 'string' ? parseFloat(price) : price;
        return (numPrice / 100).toLocaleString('de-DE', { 
            style: 'currency', 
            currency: 'EUR' 
        });
    };

    // Extract essential product data with fallbacks
    const {
        id = '',
        name = 'Unnamed Product',
        description = '',
        basePrice = 0,
        price = 0,
        thumbnail = '',
        image = '',
        categorySlug = '',
        status = 'active'
    } = product || {};

    const displayPrice = basePrice || price;
    const displayImage = thumbnail || image;

    return (
        <div className={styles.productCard}>
            <Link 
                href={`/products/${categorySlug}/${id}`}
                className={styles.productLink}
            >
                <div className={styles.imageContainer}>
                    {displayImage ? (
                        <img 
                            src={displayImage} 
                            alt={name}
                            className={styles.productImage}
                            loading="lazy"
                        />
                    ) : (
                        <div className={styles.imagePlaceholder}>
                            <Tag size={32} />
                        </div>
                    )}
                </div>
                
                <div className={styles.productInfo}>
                    <h3 className={styles.productName}>{name}</h3>
                    
                    {description && (
                        <p className={styles.productDescription}>
                            {description.length > 100 
                                ? `${description.substring(0, 97)}...` 
                                : description
                            }
                        </p>
                    )}
                    
                    <div className={styles.priceRow}>
                        <span className={styles.price}>
                            {formatPrice(displayPrice)}
                        </span>
                        
                        {status === 'active' && (
                            <button 
                                className={styles.cartButton}
                                onClick={(e) => {
                                    e.preventDefault();
                                    // TODO: Add to cart functionality
                                    
                                }}
                                aria-label="Add to cart"
                            >
                                <ShoppingCart size={18} />
                            </button>
                        )}
                    </div>
                </div>
            </Link>
        </div>
    );
};

export default memo(ProductCard);
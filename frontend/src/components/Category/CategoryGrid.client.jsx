"use client";
import React, { memo } from 'react';
// import Link from "next/link"; // ❗️ Consider using next-intl's Link for automatic locale handling
import {Link} from "../../i18n/navigation"; // Assuming you have this configured like in HomePage
import PropTypes from "prop-types";
import {useTranslations} from "next-intl"; //  Import hook
import Image from 'next/image';
import styles from "./CategoryGrid.module.css";
/**
 * CategoryGrid: Displays categories in a responsive grid with optional images, tags, etc.
 * - Uses next-intl for translations.
 *
 * Props:
 * categories: array of category objects, e.g. [{ id, name, slug, description, imageUrl, tags, googleCategoryId }, ...]
 * showGoogleCategoryId: boolean (whether to display the "Google Category" text)
 */
const CategoryGrid = memo(function CategoryGrid({
                                         categories,
                                         showGoogleCategoryId = false,
                                     }) {
    const t = useTranslations('CategoryGrid'); //  Instantiate hook with namespace
    if (!categories || categories.length === 0) {
        //   Use translation
        return <p>{t('empty')}</p>;
    }
    return (
        <div className={styles.gridWrapper}>
            {categories.map((cat) => {
                // Assume cat.description and cat.name are already appropriately translated if needed
                const altText = cat.description || cat.name || t('unnamedCategory');
                const categoryTitle = cat.description || cat.name || t('unnamedCategory'); //  Use translation fallback
                return (
                    <div key={cat.id} className={styles.categoryCard}>
                        {/* If you have a category image, display it */}
                        {cat.imageUrl && (
                            <div className={styles.imageWrapper}>
                                <Image 
                                    src={cat.imageUrl} 
                                    alt={altText}
                                    width={300}
                                    height={200}
                                    style={{ objectFit: 'cover' }}
                                />
                            </div>
                        )}
                        {/* ❗️ Ensure this Link component handles locale prefixes correctly */}
                        <Link href={`/category/${cat.slug}`} className={styles.categoryTitle}>
                            {categoryTitle}
                        </Link>
                        {/* Optional description or SEO text (assume pre-translated) */}
                        {cat.seoDesc && <p className={styles.categoryDesc}>{cat.seoDesc}</p>}
                        {/* Optional tags (assume pre-translated) */}
                        {cat.tags && cat.tags.length > 0 && (
                            <div className={styles.tagsWrapper}>
                                {cat.tags.map((tag) => (
                                    <span key={tag} className={styles.tagItem}>
                                        {tag}
                                    </span>
                                ))}
                            </div>
                        )}
                        {/* Conditionally show googleCategoryId if you want */}
                        {showGoogleCategoryId && cat.googleCategoryId && (
                            <small className={styles.googleCat}>
                                {/*   Use translation with interpolation */}
                                {t('googleCategoryPrefix', {googleCategoryId: cat.googleCategoryId})}
                            </small>
                        )}
                    </div>
                );
            })}
        </div>
    );
});
CategoryGrid.propTypes = {
    categories: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.oneOfType([PropTypes.number, PropTypes.string]).isRequired, // Added isRequired
            name: PropTypes.string,
            slug: PropTypes.string.isRequired, // Added isRequired
            description: PropTypes.string,
            seoDesc: PropTypes.string,
            tags: PropTypes.arrayOf(PropTypes.string),
            googleCategoryId: PropTypes.string,
            imageUrl: PropTypes.string,
        })
    ),
    showGoogleCategoryId: PropTypes.bool,
};
export default CategoryGrid;
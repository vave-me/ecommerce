'use client';
import React from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useDispatch } from 'react-redux';
import { useTranslations } from 'next-intl';
import { dismissVaveAd } from '../../redux/slices/adsSlice';
import styles from './CompactVaveAd.module.css';

/**
 * Compact sfx-markt.de Platform Ad - Server Component
 * Professional B2B/B2C marketplace platform advertisement
 */
const CompactVaveAd = () => {
    const dispatch = useDispatch();
    const router = useRouter();
    const t = useTranslations('Demo');
    const tShared = useTranslations('shared');
    
    const handleClose = (e) => {
        e.preventDefault();
        e.stopPropagation();
        dispatch(dismissVaveAd());
    };
    
    const handleMainClick = (e) => {
        // Only handle click if it's on the main container, not on module links
        if (e.target.closest(`.${styles.moduleItem}`)) {
            return;
        }
        window.open('https://sfx-markt.de/contact-sales', '_blank', 'noopener,noreferrer');
    };
    
    return (
        <div className={styles.adContainer}>
            <div 
                className={styles.adLink}
                onClick={handleMainClick}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        handleMainClick(e);
                    }
                }}
                aria-label="sfx-markt.de - Enterprise Marketplace Platform - Scale Your Digital Commerce"
            >
                <div className={styles.content}>
                    <div className={styles.logoSection}>
                        <Image 
                            src="/images/logo-vaveme.png" 
                            alt="sfx-markt.de" 
                            width={48} 
                            height={48}
                            className={styles.logo}
                            style={{ objectFit: 'contain' }}
                        />
                    </div>
                    
                    <div className={styles.infoSection}>
                        <div className={styles.headerInfo}>
                            <h3 className={styles.headline}>
                                {t('compactAd.headline')}
                            </h3>
                            <span className={styles.badge}>{t('compactAd.badge')}</span>
                        </div>
                        
                        <div className={styles.modulesRow}>
                            <Link href="/demo/commerce" className={styles.moduleItem}>
                                <strong>✓ {t('modules.commerce.title')}</strong>
                                <span>{t('compactAd.modules.commerce')}</span>
                            </Link>
                            <Link href="/demo/social" className={styles.moduleItem}>
                                <strong>✓ {t('modules.social.title')}</strong>
                                <span>{t('compactAd.modules.social')}</span>
                            </Link>
                            <Link href="/demo/verticals" className={styles.moduleItem}>
                                <strong>✓ {t('modules.verticals.title')}</strong>
                                <span>{t('compactAd.modules.classified')}</span>
                            </Link>
                            <Link href="/demo/intelligence" className={styles.moduleItem}>
                                <strong>✓ {t('modules.intelligence.title')}</strong>
                                <span>{t('compactAd.modules.intelligence')}</span>
                            </Link>
                            <Link href="/demo/infrastructure" className={styles.moduleItem}>
                                <strong>✓ {t('modules.infrastructure.title')}</strong>
                                <span>{t('compactAd.modules.infrastructure')}</span>
                            </Link>
                            
                        </div>
                    </div>
                    
                    <div className={styles.ctaSection}>
                        <span className={styles.cta}>
                            {tShared('startBuilding')}
                        </span>
                    </div>
                </div>
            </div>
            <button 
                className={styles.closeButton}
                onClick={handleClose}
                aria-label="Close advertisement"
                type="button"
            >
                ×
            </button>
        </div>
    );
};

export default CompactVaveAd;
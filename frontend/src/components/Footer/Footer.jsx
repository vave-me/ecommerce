import React from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import styles from './Footer.module.css';

const Footer = () => {
    const t = useTranslations('HomePage');
    
    return (
        <footer className={styles.footer}>
            <div className={styles.container}>
                <div className={styles.content}>
                    {/* Footer Links */}
                    <div className={styles.footerLinks}>
                        <Link href="/help" className={styles.footerLink}>
                            {t('help')}
                        </Link>
                        <span className={styles.separator}>|</span>
                        <Link href="/contact" className={styles.footerLink}>
                            {t('contact')}
                        </Link>
                        <span className={styles.separator}>|</span>
                        <Link href="/about" className={styles.footerLink}>
                            {t('about')}
                        </Link>
                        <span className={styles.separator}>|</span>
                        <Link href="/terms" className={styles.footerLink}>
                            {t('terms')}
                        </Link>
                        <span className={styles.separator}>|</span>
                        <Link href="/privacy" className={styles.footerLink}>
                            {t('privacy')}
                        </Link>
                    </div>
                    
                    {/* Copyright */}
                    <div className={styles.copyright}>
                        © 2025 sfx markt. {t('allRightsReserved')}
                    </div>
                    
                    {/* Powered By */}
                    <div className={styles.poweredBy}>
                        <span className={styles.poweredByText}>Powered by</span>
                        <a 
                            href="https://sfx-markt.de"
                            target="_blank" 
                            rel="noopener noreferrer"
                            className={styles.vaveLink}
                        >
                            vave me
                        </a>
                    </div>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
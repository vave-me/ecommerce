import React, { useState, useRef, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { createPortal } from 'react-dom';
import { MapPin, Mail, ChevronDown, Globe } from 'lucide-react';
import {Link} from '@/i18n/navigation';
import styles from './UtilityBar.module.css';

const languages = [
  { code: 'en', name: 'English', flag: '🇬🇧' },
  { code: 'de', name: 'Deutsch', flag: '🇩🇪' },
  { code: 'pl', name: 'Polski', flag: '🇵🇱' },
  { code: 'it', name: 'Italiano', flag: '🇮🇹' }
];

const UtilityBar = () => {
  const [isLanguageOpen, setIsLanguageOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const dropdownRef = useRef(null);
  const buttonRef = useRef(null);
  const router = useRouter();
  const pathname = usePathname();
  const locale = useLocale();
  const t = useTranslations('UtilityBar');

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target) && 
          buttonRef.current && !buttonRef.current.contains(event.target)) {
        setIsLanguageOpen(false);
      }
    };

    if (isLanguageOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isLanguageOpen]);

  const handleLanguageChange = (newLocale) => {
    const pathWithoutLocale = pathname.replace(/^\/[a-z]{2}/, '');
    const newPath = `/${newLocale}${pathWithoutLocale}`;
    router.push(newPath);
    setIsLanguageOpen(false);
  };

  const currentLanguage = languages.find(lang => lang.code === locale) || languages[0];

  return (
    <div className={styles.utilityBar}>
      <div className={styles.utilityContainer}>
        <div className={styles.utilityLeft}>
          <div className={styles.utilityItem}>
            <MapPin className={styles.icon} size={14} />
            <span>{t('location', 'Deutschland')}</span>
          </div>
          
          <div className={styles.utilityItem}>
            <Mail className={styles.icon} size={14} />
            <a href="mailto:redacted-email@example.com" className={styles.link}>
              redacted-email@example.com
            </a>
          </div>
        </div>
        
        <div className={styles.utilityRight}>
          <Link href="/sell" className={styles.utilityLink}>
            {t('sell', 'Verkaufen')}
          </Link>
          <div className={styles.divider}></div>
          <Link href="/sfx-market" className={styles.utilityLink}>
            {t('sfxSite', 'SFX Site')}
          </Link>
          <div className={styles.divider}></div>
          <Link href="/help" className={styles.utilityLink}>
            {t('help', 'Hilfe')}
          </Link>
          <div className={styles.divider}></div>
          
          <div className={styles.languageSelector}>
            <button 
              ref={buttonRef}
              className={styles.languageButton}
              onClick={() => setIsLanguageOpen(!isLanguageOpen)}
              aria-expanded={isLanguageOpen}
              aria-haspopup="true"
              aria-label={t('changeLanguage', 'Change language')}
              title={currentLanguage.name}
            >
              <Globe className={styles.icon} size={14} />
              <span className={styles.languageCode}>{currentLanguage.code.toUpperCase()}</span>
              <ChevronDown className={`${styles.chevron} ${isLanguageOpen ? styles.chevronUp : ''}`} size={12} />
            </button>
            
            {mounted && isLanguageOpen && createPortal(
              <div 
                ref={dropdownRef}
                className={styles.languageDropdown}
                style={{
                  position: 'absolute',
                  top: buttonRef.current?.getBoundingClientRect().bottom + 8,
                  right: window.innerWidth - buttonRef.current?.getBoundingClientRect().right
                }}
              >
                {languages.map((lang) => (
                  <button
                    key={lang.code}
                    className={`${styles.languageOption} ${lang.code === locale ? styles.active : ''}`}
                    onClick={() => handleLanguageChange(lang.code)}
                  >
                    <span className={styles.languageFlag}>{lang.flag}</span>
                    <span className={styles.languageName}>{lang.name}</span>
                  </button>
                ))}
              </div>,
              document.body
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default UtilityBar;

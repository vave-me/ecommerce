"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { useTranslations } from 'next-intl';
import styles from './DropdownHeader.module.css';
/**
 * DropdownHeader Component
 * Mobile sheet header with title and close button
 * Extracted from SelectTopic to improve component modularity
 */
const DropdownHeader = memo(({ 
    topicData, 
    selectedParentId, 
    onClose 
}) => {
    const t = useTranslations('Topics');
    return (
        <header className={styles.sheetHeader}>
            <span
                id={`topic-dropdown-title-${topicData.value}`}
                className={styles.sheetTitle}
            >
                {selectedParentId 
                    ? t('subcategoriesOf', {category: topicData.label}) 
                    : topicData.label
                }
            </span>
            <button
                onClick={onClose}
                className={styles.sheetCloseButton}
                aria-label={t('closeButtonAriaLabel')}
                type="button"
            >
                ✕
            </button>
        </header>
    );
});
DropdownHeader.displayName = 'DropdownHeader';
DropdownHeader.propTypes = {
    topicData: PropTypes.shape({
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
    }).isRequired,
    selectedParentId: PropTypes.string,
    onClose: PropTypes.func.isRequired,
};
export default DropdownHeader; 
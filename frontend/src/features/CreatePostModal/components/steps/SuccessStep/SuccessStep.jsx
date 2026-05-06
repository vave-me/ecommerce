// CreatePostModal/components/steps/SuccessStep/SuccessStep.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl'; // ⬅️ Import hook
import {Check} from "@/icons";
import styles from '../../../CreatePostModal.module.css'; // Use shared styles
export function SuccessStep({onViewDashboard, onClose}) {
    const t = useTranslations('CreatePostModal'); // ⬅️ Use shared namespace
    return (
        <div className={styles.successContainer} role="alertdialog" aria-labelledby="success-title"
             aria-describedby="success-message"> {/* Added roles */}
            <div className={styles.successIcon} aria-hidden="true">
                <Check size={48}/>
            </div>
            <h2 id="success-title" className={styles.successTitle}>{t('successPublishedTitle')}</h2>
            <p id="success-message" className={styles.successMessage}>
                {t('successPublishedMessage')}
            </p>
            {/* Next Steps Section */}
            <div className={styles.nextStepOptions}>
                <h4>{t('successNextStepsTitle')}</h4>
                <ul className={styles.nextStepsList}>
                    <li>
                        {t('successNextStep1')}
                    </li>
                    <li>
                        {t('successNextStep2')}
                    </li>
                    <li>
                        {t('successNextStep3')}
                    </li>
                </ul>
            </div>
            {/* Action Buttons */}
            <div className={styles.successActions}>
                <button
                    className={styles.secondaryButton} // Assuming secondary style for dashboard button
                    type="button"
                    onClick={onViewDashboard}
                >
                    {t('successGoDashboardButton')}
                </button>
                <button
                    className={styles.primaryButton} // Assuming primary style for close button
                    type="button"
                    onClick={onClose}
                >
                    {t('successCloseButton')}
                </button>
            </div>
        </div>
    );
}
// PropTypes remain the same
SuccessStep.propTypes = {
    onViewDashboard: PropTypes.func.isRequired,
    onClose: PropTypes.func.isRequired
};
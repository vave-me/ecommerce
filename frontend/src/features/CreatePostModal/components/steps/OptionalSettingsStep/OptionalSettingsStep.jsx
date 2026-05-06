// CreatePostModal/components/steps/OptionalSettingsStep/OptionalSettingsStep.jsx
import React, {useMemo, useState} from 'react'; // Added useMemo
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl';
import styles from '../../../CreatePostModal.module.css';
import {FormActions} from "../../../../../common/components/FormActions";
export function OptionalSettingsStep({onPublish, onBack, isLoading}) {
    const t = useTranslations('CreatePostModal');
    const [termsAccepted, setTermsAccepted] = useState(false);
    // Determine button labels using translations
    const primaryLabel = useMemo(() => (
        isLoading ? t('optionalPublishingButton') : t('optionalPublishButton')
    ), [isLoading, t]);
    const secondaryLabel = useMemo(() => t('optionalBackButton'), [t]);
    return (
        <div className={styles.formContainer}>
            <h2 className={styles.formTitle}>{t('optionalStepTitle')}</h2>
            <p className={styles.formDescription}>{t('optionalStepDescription')}</p>
            {/* Placeholder for future settings */}
            <div className={styles.placeholder}>
                <p>{t('optionalPlaceholderText')}</p>
            </div>
            {/* Terms Confirmation Section */}
            <div className={styles.publishSection}>
                <div className={styles.publishCheckbox}>
                    <input
                        type="checkbox"
                        id="terms-checkbox"
                        className={styles.checkbox} // Ensure styling targets this class
                        checked={termsAccepted}
                        onChange={(e) => setTermsAccepted(e.target.checked)}
                        required // Keep for basic browser validation hint
                        aria-required="true" // Explicitly state requirement
                    />
                    <label
                        htmlFor="terms-checkbox"
                        className={styles.checkboxLabel}
                    >
                        {t('optionalTermsLabel')}
                    </label>
                </div>
            </div>
            {/* Use FormActions, passing translated labels */}
            {/* Assumes FormActions handles its own internal translations (e.g., default Cancel text if different) */}
            <FormActions
                primaryLabel={primaryLabel} // Pass translated label
                primaryIcon="check" // Or pass icon component if needed
                secondaryLabel={secondaryLabel} // Pass translated label
                onCancel={onBack} // Secondary action is "Back"
                onPrimaryAction={onPublish} // Primary action is "Publish"
                // Disable primary button if loading OR terms not accepted
                isPrimaryDisabled={isLoading || !termsAccepted}
                // isSubmitting={isLoading} // Pass loading state if FormActions needs it
            />
        </div>
    );
}
// PropTypes remain the same
OptionalSettingsStep.propTypes = {
    onPublish: PropTypes.func.isRequired,
    onBack: PropTypes.func.isRequired,
    isLoading: PropTypes.bool
};
// Default props if needed
OptionalSettingsStep.defaultProps = {
    isLoading: false,
};
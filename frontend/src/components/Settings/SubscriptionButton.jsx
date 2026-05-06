// src/components/Settings/SubscriptionButton.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Check, Plus } from '@/icons';
import { FaSpinner } from '../../utils/iconImports';
import styles from './SubscriptionButton.module.css';
const SubscriptionButton = memo(({ isSubscribed, onClick, disabled }) => {
    const buttonClass = `${styles.button} ${
        disabled ? styles.disabled : isSubscribed ? styles.subscribed : styles.notSubscribed
    }`;
    return (
        <button
            className={buttonClass}
            onClick={onClick}
            disabled={disabled}
            aria-label={isSubscribed ? 'Unsubscribe' : 'Subscribe'}
        >
            {disabled ? (
                <FaSpinner className={styles.spinner} />
            ) : isSubscribed ? (
                <>
                    <Check /> Unsubscribe
                </>
            ) : (
                <>
                    <Plus /> Subscribe
                </>
            )}
        </button>
    );
});
SubscriptionButton.displayName = 'SubscriptionButton';
SubscriptionButton.propTypes = {
    isSubscribed: PropTypes.bool.isRequired,
    onClick: PropTypes.func.isRequired,
    disabled: PropTypes.bool,
};
SubscriptionButton.defaultProps = {
    disabled: false,
};
export default SubscriptionButton;
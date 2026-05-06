import { FaSpinner } from '../../utils/iconImports';
import React, { memo } from 'react';
import PropTypes from 'prop-types';
 // For loading indicator
import styles from './Button.module.css';
const Button = memo(({
                    children,
                    onClick,
                    primary,
                    type,
                    disabled,
                    fullWidth,
                    loading,
                    icon,
                    ...rest
                }) => {
    const buttonClasses = [
        styles.button,
        primary ? styles.primary : '',
        fullWidth ? styles.fullWidth : '',
        loading ? styles.loading : '',
    ].filter(Boolean).join(' ');
    return (
        <button
            onClick={onClick}
            type={type}
            disabled={disabled || loading}
            aria-disabled={disabled || loading}
            className={buttonClasses}
            {...rest}
        >
            {loading ? (
                <span className={styles.spinnerIcon}><FaSpinner /></span>
            ) : (
                <>
                    {icon && <span className={styles.iconWrapper}>{icon}</span>}
                    <span className={styles.buttonText}>{children}</span>
                </>
            )}
        </button>
    );
});
Button.displayName = 'Button';
Button.propTypes = {
    children: PropTypes.node.isRequired,
    onClick: PropTypes.func,
    primary: PropTypes.bool,
    type: PropTypes.oneOf(['button', 'submit', 'reset']),
    disabled: PropTypes.bool,
    fullWidth: PropTypes.bool,
    loading: PropTypes.bool,
    icon: PropTypes.element, // Optional: To include an icon
};
Button.defaultProps = {
    onClick: () => {},
    primary: false,
    type: 'button',
    disabled: false,
    fullWidth: false,
    loading: false,
    icon: null,
};
export default Button;

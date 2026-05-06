// common/components/FormActions.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import {ArrowRight, Check} from "@/icons";
import styles from './FormActions.module.css';
export const FormActions = memo(function FormActions({
                                                               primaryLabel,
                                                               primaryIcon = null,
                                                               secondaryLabel = "Cancel",
                                                               onCancel,
                                                               onPrimaryAction,
                                                               isSubmitting = false
                                                           }) {
    return (
        <div className={styles.formActions}>
            <button
                type="button"
                className={styles.cancelButton}
                onClick={onCancel}
                disabled={isSubmitting}
            >
                {secondaryLabel}
            </button>
            <button
                type={onPrimaryAction ? "button" : "submit"}
                className={styles.primaryButton}
                onClick={onPrimaryAction}
                disabled={isSubmitting}
            >
                {primaryLabel}
                {!isSubmitting && primaryIcon === "arrow-right" && <ArrowRight size={16}/>}
                {!isSubmitting && primaryIcon === "check" && <Check size={16}/>}
            </button>
        </div>
    );
});
FormActions.propTypes = {
    primaryLabel: PropTypes.string.isRequired,
    primaryIcon: PropTypes.oneOf([null, "arrow-right", "check"]),
    secondaryLabel: PropTypes.string,
    onCancel: PropTypes.func.isRequired,
    onPrimaryAction: PropTypes.func,
    isSubmitting: PropTypes.bool
};
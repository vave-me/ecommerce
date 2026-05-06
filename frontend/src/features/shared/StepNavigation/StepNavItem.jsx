// CreatePostModal/components/StepNavigation/StepNavItem.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { Check } from "@/icons";
import styles from '../../CreatePostModal/CreatePostModal.module.css';
export function StepNavItem({
                                stepNum,
                                label,
                                isActive,
                                isCompleted,
                                isDisabled,
                                onClick,
                            }) {
    const baseClass = styles.stepNavItem;
    const activeClass = isActive ? styles.stepNavActive : "";
    const completedClass = isCompleted ? styles.stepNavCompleted : "";
    const disabledClass = isDisabled ? styles.stepNavDisabled : "";
    const combinedClass = `${baseClass} ${activeClass} ${completedClass} ${disabledClass}`;
    return (
        <div
            className={combinedClass}
            onClick={isDisabled ? null : onClick}
            role={isDisabled ? undefined : "button"}
            tabIndex={isDisabled ? undefined : 0}
        >
            <div className={styles.stepNumCircle}>
                {isCompleted ? <Check size={16}/> : stepNum}
            </div>
            <span className={styles.stepLabel}>{label}</span>
        </div>
    );
}
StepNavItem.propTypes = {
    stepNum: PropTypes.number.isRequired,
    label: PropTypes.string.isRequired,
    isActive: PropTypes.bool,
    isCompleted: PropTypes.bool,
    isDisabled: PropTypes.bool,
    onClick: PropTypes.func,
};
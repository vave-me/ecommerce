// CreatePostModal/components/StepNavigation/StepNavigation.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl'; //  Import hook
import {StepNavItem} from './StepNavItem'; // Assume StepNavItem handles internal text if any
import styles from '../../CreatePostModal/CreatePostModal.module.css'; // Use shared styles
export function StepNavigation({currentStep, lastCompletedStep, onStepClick}) {
    const t = useTranslations('CreatePostModal'); //  Use shared namespace
    // Determine if a step can be clicked (must be completed or the next step)
    // Note: Original logic allowed clicking completed steps only if >= current step.
    // Adjusted to allow clicking any *previously completed* step for easier navigation back.
    // Also prevent clicking the final "Publish" step as it's not a distinct form step.
    const canClickStep = (stepNum) => {
        return stepNum <= lastCompletedStep + 1 && stepNum < 4; // Can click up to last completed + 1, but not step 4
    };
    return (
        <div className={styles.sidebar}>
            <div className={styles.stepsContainer} role="navigation"
                 aria-label={t('stepNavAriaLabel')}> {/* Added ARIA */}
                <StepNavItem
                    stepNum={1}
                    //   Use translation
                    label={t('step1Name')}
                    isActive={currentStep === 1}
                    isCompleted={lastCompletedStep >= 1}
                    // Allow clicking step 1 if it's reachable
                    onClick={canClickStep(1) ? () => onStepClick(1) : undefined}
                    // Disable button functionality if not clickable
                    isDisabled={!canClickStep(1) && currentStep !== 1}
                />
                <StepNavItem
                    stepNum={2}
                    //   Use translation
                    label={t('step2Name')}
                    isActive={currentStep === 2}
                    isCompleted={lastCompletedStep >= 2}
                    isDisabled={lastCompletedStep < 1} // Disabled if step 1 not done
                    onClick={canClickStep(2) ? () => onStepClick(2) : undefined}
                />
                <StepNavItem
                    stepNum={3}
                    //   Use translation (Corrected key based on code)
                    label={t('step3Name')}
                    isActive={currentStep === 3}
                    isCompleted={lastCompletedStep >= 3}
                    isDisabled={lastCompletedStep < 2} // Disabled if step 2 not done
                    onClick={canClickStep(3) ? () => onStepClick(3) : undefined}
                />
                {/* Step 4 is the success/final state, typically not navigable */}
                <StepNavItem
                    stepNum={4}
                    //   Use translation
                    label={t('step4Name')}
                    // Active only when success state is reached (assuming currentStep would be 4 or handled by success state)
                    isActive={currentStep === 4}
                    isCompleted={currentStep === 4} // Completed only when active
                    isDisabled={lastCompletedStep < 3} // Disabled if step 3 not done
                    // Usually not clickable
                />
            </div>
        </div>
    );
}
StepNavigation.propTypes = {
    currentStep: PropTypes.number.isRequired,
    lastCompletedStep: PropTypes.number.isRequired,
    onStepClick: PropTypes.func.isRequired,
};
import React from 'react';
import PropTypes from 'prop-types';
/**
 * Shared NameDescriptionInputs component for consistent name and description input fields across all modals
 */
export function NameDescriptionInputs({
    nameId = "name",
    descriptionId = "description",
    nameLabel = "Name",
    descriptionLabel = "Description",
    nameValue,
    descriptionValue,
    onNameChange,
    onDescriptionChange,
    namePlaceholder = "Enter name",
    descriptionPlaceholder = "Enter description",
    nameError,
    descriptionError,
    className = "",
    styles,
    nameRequired = true,
    descriptionRequired = true,
    disabled = false,
    nameMaxLength = 150,
    descriptionRows = 6,
    ...props
}) {
    const handleNameChange = (e) => {
        if (onNameChange) {
            onNameChange(e.target.value);
        }
    };
    const handleDescriptionChange = (e) => {
        if (onDescriptionChange) {
            onDescriptionChange(e.target.value);
        }
    };
    return (
        <div className={className} {...props}>
            {/* Name Input */}
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={nameId}
                >
                    {nameLabel}
                    {nameRequired && <span className={styles?.requiredMark || 'required-mark'}> *</span>}
                </label>
                <input
                    id={nameId}
                    className={`${styles?.formInput || 'form-input'} ${nameError ? styles?.inputError || 'input-error' : ''}`}
                    type="text"
                    value={nameValue || ''}
                    onChange={handleNameChange}
                    placeholder={namePlaceholder}
                    maxLength={nameMaxLength}
                    disabled={disabled}
                    aria-invalid={!!nameError}
                    aria-describedby={nameError ? `${nameId}-error` : undefined}
                />
                {nameError && (
                    <div 
                        className={styles?.fieldError || 'field-error'} 
                        id={`${nameId}-error`} 
                        role="alert"
                    >
                        {nameError}
                    </div>
                )}
            </div>
            {/* Description Textarea */}
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={descriptionId}
                >
                    {descriptionLabel}
                    {descriptionRequired && <span className={styles?.requiredMark || 'required-mark'}> *</span>}
                </label>
                <textarea
                    id={descriptionId}
                    className={`${styles?.formTextarea || 'form-textarea'} ${descriptionError ? styles?.inputError || 'input-error' : ''}`}
                    value={descriptionValue || ''}
                    onChange={handleDescriptionChange}
                    placeholder={descriptionPlaceholder}
                    rows={descriptionRows}
                    disabled={disabled}
                    aria-invalid={!!descriptionError}
                    aria-describedby={descriptionError ? `${descriptionId}-error` : undefined}
                />
                {descriptionError && (
                    <div 
                        className={styles?.fieldError || 'field-error'} 
                        id={`${descriptionId}-error`} 
                        role="alert"
                    >
                        {descriptionError}
                    </div>
                )}
            </div>
        </div>
    );
}
NameDescriptionInputs.propTypes = {
    nameId: PropTypes.string,
    descriptionId: PropTypes.string,
    nameLabel: PropTypes.string,
    descriptionLabel: PropTypes.string,
    nameValue: PropTypes.string,
    descriptionValue: PropTypes.string,
    onNameChange: PropTypes.func.isRequired,
    onDescriptionChange: PropTypes.func.isRequired,
    namePlaceholder: PropTypes.string,
    descriptionPlaceholder: PropTypes.string,
    nameError: PropTypes.string,
    descriptionError: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    nameRequired: PropTypes.bool,
    descriptionRequired: PropTypes.bool,
    disabled: PropTypes.bool,
    nameMaxLength: PropTypes.number,
    descriptionRows: PropTypes.number,
};
export default NameDescriptionInputs; 
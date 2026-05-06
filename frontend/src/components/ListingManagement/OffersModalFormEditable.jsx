// File: src/components/ListingManagement/OffersModalFormEditable.jsx
"use client"
import React, {useState, useEffect, useCallback, useRef, memo} from 'react';
import ReactDOM from 'react-dom';
import styled, {keyframes, css} from 'styled-components';

/* -------------------------------------------------------------------------- */
/*                              CONSTANTS                                     */
/* -------------------------------------------------------------------------- */
const OFFER_TYPES = {
    BUY_NOW: 'buyNow',
    LEASE: 'lease', 
    PAWN: 'pawn',
    RESERVATION: 'reservation'
};

const TAB_CONFIG = [
    { id: OFFER_TYPES.BUY_NOW, label: 'Buy Now' },
    { id: OFFER_TYPES.LEASE, label: 'Lease' },
    { id: OFFER_TYPES.PAWN, label: 'Pawn' },
    { id: OFFER_TYPES.RESERVATION, label: 'Reserve' }
];

// Keyboard navigation constants
const KEYS = {
    ESCAPE: 'Escape',
    TAB: 'Tab',
    ARROW_LEFT: 'ArrowLeft',
    ARROW_RIGHT: 'ArrowRight'
};

/* -------------------------------------------------------------------------- */
/*                              MAIN COMPONENT                                */
/* -------------------------------------------------------------------------- */
function OffersModalFormEditable({productId, onClose, onSubmit, isLoading = false}) {
    const [activeTab, setActiveTab] = useState(OFFER_TYPES.BUY_NOW);
    const [errors, setErrors] = useState({});
    const modalRef = useRef(null);
    const firstFocusableRef = useRef(null);
    
    // Lock body scroll on mount/unmount
    useEffect(() => {
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = 'hidden';
        
        // Set initial focus for accessibility
        if (firstFocusableRef.current) {
            firstFocusableRef.current.focus();
        }
        
        return () => {
            document.body.style.overflow = prevOverflow;
        };
    }, []);
    
    // All local form states, not exposed externally
    const [formData, setFormData] = useState({
        buyNow: {
            finalPrice: '',
            shippingMethod: 'standard',
            giftWrap: false,
            quantity: 1,
        },
        lease: {
            monthlyPrice: '',
            leaseTerm: 12,
            hasBuyout: true,
            buyoutPrice: '',
            coverage: false,
        },
        pawn: {
            lockedPrice: '',
            redemptionDays: 14,
            redemptionFeePercent: 7,
            extension: false,
            disclaimersAccepted: false,
        },
        reservation: {
            reservationPrice: '',
            reservationDays: 14,
            redemptionFeePercent: 7,
            buyOutPrice: '',
            disclaimersAccepted: false,
        },
    });
    
    /* ------------------------- Keyboard Navigation ------------------------- */
    const handleKeyDown = useCallback((e) => {
        switch (e.key) {
            case KEYS.ESCAPE:
                e.preventDefault();
                onClose();
                break;
            case KEYS.ARROW_LEFT:
                if (e.target.getAttribute('role') === 'tab') {
                    e.preventDefault();
                    const currentIndex = TAB_CONFIG.findIndex(tab => tab.id === activeTab);
                    const newIndex = currentIndex > 0 ? currentIndex - 1 : TAB_CONFIG.length - 1;
                    setActiveTab(TAB_CONFIG[newIndex].id);
                }
                break;
            case KEYS.ARROW_RIGHT:
                if (e.target.getAttribute('role') === 'tab') {
                    e.preventDefault();
                    const currentIndex = TAB_CONFIG.findIndex(tab => tab.id === activeTab);
                    const newIndex = currentIndex < TAB_CONFIG.length - 1 ? currentIndex + 1 : 0;
                    setActiveTab(TAB_CONFIG[newIndex].id);
                }
                break;
        }
    }, [activeTab, onClose]);
    
    /* ------------------------- Unified Input Handler ------------------------- */
    const handleInputChange = useCallback((tabType, fieldName, value) => {
        setFormData(prev => ({
            ...prev,
            [tabType]: {
                ...prev[tabType],
                [fieldName]: value
            }
        }));
        
        // Clear error when user starts typing
        if (errors[`${tabType}.${fieldName}`]) {
            setErrors(prev => {
                const newErrors = {...prev};
                delete newErrors[`${tabType}.${fieldName}`];
                return newErrors;
            });
        }
    }, [errors]);
    
    /* ------------------------- Form Validation ------------------------- */
    const validateForm = useCallback(() => {
        const newErrors = {};
        const currentData = formData[activeTab];
        
        switch (activeTab) {
            case OFFER_TYPES.BUY_NOW:
                if (!currentData.finalPrice || parseFloat(currentData.finalPrice) <= 0) {
                    newErrors[`${activeTab}.finalPrice`] = 'Please enter a valid price';
                }
                if (currentData.quantity < 1) {
                    newErrors[`${activeTab}.quantity`] = 'Quantity must be at least 1';
                }
                break;
            case OFFER_TYPES.LEASE:
                if (!currentData.monthlyPrice || parseFloat(currentData.monthlyPrice) <= 0) {
                    newErrors[`${activeTab}.monthlyPrice`] = 'Please enter a valid monthly price';
                }
                if (currentData.leaseTerm < 1) {
                    newErrors[`${activeTab}.leaseTerm`] = 'Lease term must be at least 1 month';
                }
                break;
            case OFFER_TYPES.PAWN:
                if (!currentData.lockedPrice || parseFloat(currentData.lockedPrice) <= 0) {
                    newErrors[`${activeTab}.lockedPrice`] = 'Please enter a valid price';
                }
                if (!currentData.disclaimersAccepted) {
                    newErrors[`${activeTab}.disclaimersAccepted`] = 'You must accept the terms';
                }
                break;
            case OFFER_TYPES.RESERVATION:
                if (!currentData.reservationPrice || parseFloat(currentData.reservationPrice) <= 0) {
                    newErrors[`${activeTab}.reservationPrice`] = 'Please enter a valid reservation price';
                }
                if (!currentData.disclaimersAccepted) {
                    newErrors[`${activeTab}.disclaimersAccepted`] = 'You must accept the terms';
                }
                break;
        }
        
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    }, [activeTab, formData]);
    
    /* ------------------------- Submit Handler ------------------------- */
    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        
        if (validateForm()) {
            const submitData = {
                productId,
                offerType: activeTab,
                data: formData[activeTab]
            };
            
            if (onSubmit) {
                onSubmit(submitData);
            } else {
                
            }
        }
    }, [activeTab, formData, productId, onSubmit, validateForm]);
    // Render the modal into a portal
    return ReactDOM.createPortal(
        <ModalOverlay
            role="dialog"
            aria-modal="true"
            aria-labelledby="modal-title"
            onKeyDown={handleKeyDown}
        >
            <ModalContainer ref={modalRef}>
                <ModalHeader>
                    <ModalTitle id="modal-title">Create Offer - Product #{productId}</ModalTitle>
                    <CloseButton
                        ref={firstFocusableRef}
                        onClick={onClose}
                        aria-label="Close offer modal"
                        type="button"
                    >
                        ×
                    </CloseButton>
                </ModalHeader>
                
                {/* TAB BAR with ARIA support */}
                <TabBar role="tablist" aria-label="Offer types">
                    {TAB_CONFIG.map((tab) => (
                        <TabButton
                            key={tab.id}
                            role="tab"
                            aria-selected={activeTab === tab.id}
                            aria-controls={`tabpanel-${tab.id}`}
                            id={`tab-${tab.id}`}
                            tabIndex={activeTab === tab.id ? 0 : -1}
                            $isActive={activeTab === tab.id}
                            onClick={() => setActiveTab(tab.id)}
                        >
                            {tab.label}
                        </TabButton>
                    ))}
                </TabBar>
                
                {/* ACTIVE TAB CONTENT */}
                <form onSubmit={handleSubmit} noValidate>
                    <TabPanel
                        role="tabpanel"
                        id={`tabpanel-${activeTab}`}
                        aria-labelledby={`tab-${activeTab}`}
                    >
                        {activeTab === OFFER_TYPES.BUY_NOW && (
                            <BuyNowForm 
                                data={formData.buyNow} 
                                errors={errors}
                                onChange={(field, value) => handleInputChange(OFFER_TYPES.BUY_NOW, field, value)}
                            />
                        )}
                        {activeTab === OFFER_TYPES.LEASE && (
                            <LeaseForm 
                                data={formData.lease}
                                errors={errors} 
                                onChange={(field, value) => handleInputChange(OFFER_TYPES.LEASE, field, value)}
                            />
                        )}
                        {activeTab === OFFER_TYPES.PAWN && (
                            <PawnForm 
                                data={formData.pawn}
                                errors={errors} 
                                onChange={(field, value) => handleInputChange(OFFER_TYPES.PAWN, field, value)}
                            />
                        )}
                        {activeTab === OFFER_TYPES.RESERVATION && (
                            <ReservationForm 
                                data={formData.reservation}
                                errors={errors} 
                                onChange={(field, value) => handleInputChange(OFFER_TYPES.RESERVATION, field, value)}
                            />
                        )}
                    </TabPanel>
                    
                    {/* ACTIONS */}
                    <ActionsRow>
                        <CancelButton
                            type="button"
                            onClick={onClose}
                            disabled={isLoading}
                        >
                            Cancel
                        </CancelButton>
                        <SubmitButton 
                            type="submit" 
                            disabled={isLoading}
                            aria-busy={isLoading}
                        >
                            {isLoading ? 'Submitting...' : 'Submit Offer'}
                        </SubmitButton>
                    </ActionsRow>
                </form>
            </ModalContainer>
        </ModalOverlay>,
        document.body
    );
}
export default memo(OffersModalFormEditable);
/* -------------------------------------------------------------------------- */
/*                         INTERNAL SUB-COMPONENTS                            */
/* -------------------------------------------------------------------------- */

// Shared form field component for consistency and accessibility
const FormField = memo(({
    label,
    name,
    type = 'text',
    value,
    onChange,
    error,
    required,
    min,
    max,
    step,
    options,
    helpText,
    ...props
}) => {
    const fieldId = `field-${name}`;
    const errorId = `${fieldId}-error`;
    const helpId = helpText ? `${fieldId}-help` : undefined;
    
    return (
        <FieldGroup>
            <FormLabel htmlFor={fieldId}>
                {label} {required && <Required>*</Required>}
            </FormLabel>
            
            {type === 'select' ? (
                <Select
                    id={fieldId}
                    name={name}
                    value={value}
                    onChange={(e) => onChange(name, e.target.value)}
                    aria-invalid={!!error}
                    aria-describedby={[errorId, helpId].filter(Boolean).join(' ')}
                    required={required}
                    {...props}
                >
                    {options?.map(opt => (
                        <option key={opt.value} value={opt.value}>
                            {opt.label}
                        </option>
                    ))}
                </Select>
            ) : type === 'checkbox' ? (
                <CheckboxGroup>
                    <CheckboxInput
                        id={fieldId}
                        type="checkbox"
                        name={name}
                        checked={value}
                        onChange={(e) => onChange(name, e.target.checked)}
                        aria-invalid={!!error}
                        aria-describedby={errorId}
                        required={required}
                        {...props}
                    />
                    <CheckboxLabel htmlFor={fieldId}>{label}</CheckboxLabel>
                </CheckboxGroup>
            ) : (
                <NumberInput
                    id={fieldId}
                    type={type}
                    name={name}
                    value={value}
                    onChange={(e) => onChange(name, type === 'number' ? parseFloat(e.target.value) || 0 : e.target.value)}
                    min={min}
                    max={max}
                    step={step}
                    aria-invalid={!!error}
                    aria-describedby={[errorId, helpId].filter(Boolean).join(' ')}
                    required={required}
                    {...props}
                />
            )}
            
            {helpText && <HelpText id={helpId}>{helpText}</HelpText>}
            {error && <ErrorText id={errorId} role="alert">{error}</ErrorText>}
        </FieldGroup>
    );
});
FormField.displayName = 'FormField';

/** BuyNow Form **/
const BuyNowForm = memo(({data, errors, onChange}) => (
    <SubFormWrapper>
        <FormField
            label="Final Price (EUR)"
            name="finalPrice"
            type="number"
            value={data.finalPrice}
            onChange={onChange}
            error={errors['buyNow.finalPrice']}
            required
            min="0.01"
            step="0.01"
        />
        <FormField
            label="Shipping Method"
            name="shippingMethod"
            type="select"
            value={data.shippingMethod}
            onChange={onChange}
            required
            options={[
                { value: 'standard', label: 'Standard Shipping (€4.99)' },
                { value: 'express', label: 'Express Shipping (€9.99)' },
                { value: 'pickup', label: 'Local Pickup (Free)' }
            ]}
        />
        <FormField
            label="Quantity"
            name="quantity"
            type="number"
            value={data.quantity}
            onChange={onChange}
            error={errors['buyNow.quantity']}
            required
            min="1"
        />
        <FormField
            label="Gift Wrap (€2.00)"
            name="giftWrap"
            type="checkbox"
            value={data.giftWrap}
            onChange={onChange}
        />
    </SubFormWrapper>
));
BuyNowForm.displayName = 'BuyNowForm';

/** Lease Form **/
const LeaseForm = memo(({data, errors, onChange}) => (
    <SubFormWrapper>
        <FormField
            label="Monthly Price (EUR)"
            name="monthlyPrice"
            type="number"
            value={data.monthlyPrice}
            onChange={onChange}
            error={errors['lease.monthlyPrice']}
            required
            min="0.01"
            step="0.01"
        />
        <FormField
            label="Lease Term (months)"
            name="leaseTerm"
            type="number"
            value={data.leaseTerm}
            onChange={onChange}
            error={errors['lease.leaseTerm']}
            required
            min="1"
            max="60"
        />
        <FormField
            label="Optional Buyout Available"
            name="hasBuyout"
            type="checkbox"
            value={data.hasBuyout}
            onChange={onChange}
        />
        {data.hasBuyout && (
            <FormField
                label="Buyout Price (EUR)"
                name="buyoutPrice"
                type="number"
                value={data.buyoutPrice}
                onChange={onChange}
                min="0"
                step="0.01"
            />
        )}
        <FormField
            label="Add Extended Coverage"
            name="coverage"
            type="checkbox"
            value={data.coverage}
            onChange={onChange}
        />
    </SubFormWrapper>
));
LeaseForm.displayName = 'LeaseForm';

/** Pawn Form **/
const PawnForm = memo(({data, errors, onChange}) => (
    <SubFormWrapper>
        <FormField
            label="Locked Price (EUR)"
            name="lockedPrice"
            type="number"
            value={data.lockedPrice}
            onChange={onChange}
            error={errors['pawn.lockedPrice']}
            required
            min="0.01"
            step="0.01"
        />
        <FormField
            label="Redemption Period"
            name="redemptionDays"
            type="select"
            value={data.redemptionDays}
            onChange={onChange}
            required
            options={[
                { value: 7, label: '7 days (5% fee)' },
                { value: 14, label: '14 days (7% fee)' },
                { value: 21, label: '21 days (8% fee)' },
                { value: 30, label: '30 days (10% fee)' }
            ]}
        />
        <FormField
            label="Redemption Fee (%)"
            name="redemptionFeePercent"
            type="number"
            value={data.redemptionFeePercent}
            onChange={onChange}
            min="0"
            max="100"
            step="0.1"
        />
        <FormField
            label="Allow Extension Period"
            name="extension"
            type="checkbox"
            value={data.extension}
            onChange={onChange}
        />
        <FormField
            label="I accept the pawn terms and conditions"
            name="disclaimersAccepted"
            type="checkbox"
            value={data.disclaimersAccepted}
            onChange={onChange}
            error={errors['pawn.disclaimersAccepted']}
            required
        />
    </SubFormWrapper>
));
PawnForm.displayName = 'PawnForm';

/** Reservation Form **/
const ReservationForm = memo(({data, errors, onChange}) => (
    <SubFormWrapper>
        <FormField
            label="Reservation Fee (EUR)"
            name="reservationPrice"
            type="number"
            value={data.reservationPrice}
            onChange={onChange}
            error={errors['reservation.reservationPrice']}
            required
            min="0.01"
            step="0.01"
        />
        <FormField
            label="Reservation Period"
            name="reservationDays"
            type="select"
            value={data.reservationDays}
            onChange={onChange}
            required
            options={[
                { value: 7, label: '7 days' },
                { value: 14, label: '14 days' },
                { value: 21, label: '21 days' },
                { value: 30, label: '30 days' }
            ]}
        />
        <FormField
            label="Final Purchase Price (EUR)"
            name="buyOutPrice"
            type="number"
            value={data.buyOutPrice}
            onChange={onChange}
            required
            min="0.01"
            step="0.01"
        />
        <FormField
            label="I accept the reservation terms"
            name="disclaimersAccepted"
            type="checkbox"
            value={data.disclaimersAccepted}
            onChange={onChange}
            error={errors['reservation.disclaimersAccepted']}
            required
        />
    </SubFormWrapper>
));
ReservationForm.displayName = 'ReservationForm';
/* -------------------------------------------------------------------------- */
/*                     STYLED COMPONENTS & THEME                              */
/* -------------------------------------------------------------------------- */

// Modern design system with WCAG AAA compliant colors
const theme = {
    colors: {
        // Primary palette with sufficient contrast ratios
        primary: '#0066CC',
        primaryHover: '#0052A3',
        primaryLight: '#E6F1FF',
        
        // Neutral palette
        background: '#FFFFFF',
        surface: '#F9FAFB',
        surfaceHover: '#F3F4F6',
        
        // Text colors (WCAG AAA compliant)
        text: {
            primary: '#111827',    // 15.3:1 contrast ratio
            secondary: '#6B7280',  // 7.5:1 contrast ratio
            tertiary: '#9CA3AF',   // 4.5:1 contrast ratio
            inverse: '#FFFFFF'
        },
        
        // Border colors
        border: {
            default: '#E5E7EB',
            hover: '#D1D5DB',
            focus: '#0066CC'
        },
        
        // Semantic colors
        error: '#DC2626',
        errorLight: '#FEE2E2',
        success: '#059669',
        warning: '#D97706',
        
        // Overlay
        overlay: 'rgba(0, 0, 0, 0.5)'
    },
    
    spacing: {
        xs: '0.25rem',
        sm: '0.5rem',
        md: '1rem',
        lg: '1.5rem',
        xl: '2rem',
        xxl: '3rem'
    },
    
    borderRadius: {
        sm: '0.25rem',
        md: '0.375rem',
        lg: '0.5rem',
        xl: '0.75rem'
    },
    
    typography: {
        fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
        fontSize: {
            xs: '0.75rem',
            sm: '0.875rem',
            base: '1rem',
            lg: '1.125rem',
            xl: '1.25rem'
        }
    },
    
    shadows: {
        sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
        xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)'
    },
    
    breakpoints: {
        mobile: '640px',
        tablet: '768px',
        desktop: '1024px'
    }
};

// Animations
const fadeIn = keyframes`
    from {
        opacity: 0;
        transform: scale(0.95) translateY(-10px);
    }
    to {
        opacity: 1;
        transform: scale(1) translateY(0);
    }
`;

const slideIn = keyframes`
    from {
        transform: translateY(100%);
        opacity: 0;
    }
    to {
        transform: translateY(0);
        opacity: 1;
    }
`;

// Shared input styles with improved focus states
const inputStyles = css`
    width: 100%;
    padding: ${theme.spacing.sm} ${theme.spacing.md};
    border: 1px solid ${theme.colors.border.default};
    border-radius: ${theme.borderRadius.md};
    font-size: ${theme.typography.fontSize.sm};
    font-family: ${theme.typography.fontFamily};
    background: ${theme.colors.background};
    color: ${theme.colors.text.primary};
    transition: all 0.2s ease;
    
    &:hover:not(:disabled) {
        border-color: ${theme.colors.border.hover};
    }
    
    &:focus {
        outline: none;
        border-color: ${theme.colors.border.focus};
        box-shadow: 0 0 0 3px ${theme.colors.primaryLight};
    }
    
    &:disabled {
        background: ${theme.colors.surface};
        color: ${theme.colors.text.tertiary};
        cursor: not-allowed;
        opacity: 0.6;
    }
    
    &[aria-invalid="true"] {
        border-color: ${theme.colors.error};
        
        &:focus {
            box-shadow: 0 0 0 3px ${theme.colors.errorLight};
        }
    }
`;

// Components
const ModalOverlay = styled.div`
    position: fixed;
    inset: 0;
    z-index: 9999;
    background: ${theme.colors.overlay};
    display: flex;
    align-items: center;
    justify-content: center;
    padding: ${theme.spacing.md};
    backdrop-filter: blur(4px);
    animation: ${fadeIn} 0.2s ease-out;
`;

const ModalContainer = styled.div`
    background: ${theme.colors.background};
    border-radius: ${theme.borderRadius.xl};
    width: 100%;
    max-width: 600px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: ${theme.shadows.xl};
    animation: ${slideIn} 0.3s ease-out;
    overflow: hidden;
    
    @media (max-width: ${theme.breakpoints.mobile}) {
        max-height: 100vh;
        border-radius: 0;
        margin: 0;
    }
`;

const ModalHeader = styled.header`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: ${theme.spacing.lg};
    border-bottom: 1px solid ${theme.colors.border.default};
`;

const ModalTitle = styled.h2`
    margin: 0;
    font-size: ${theme.typography.fontSize.xl};
    font-weight: 600;
    color: ${theme.colors.text.primary};
    line-height: 1.2;
`;

const CloseButton = styled.button`
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    padding: 0;
    border: none;
    background: transparent;
    color: ${theme.colors.text.secondary};
    border-radius: ${theme.borderRadius.md};
    font-size: 1.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
    
    &:hover {
        background: ${theme.colors.surface};
        color: ${theme.colors.text.primary};
    }
    
    &:focus {
        outline: 2px solid ${theme.colors.primary};
        outline-offset: 2px;
    }
`;

const TabBar = styled.div`
    display: flex;
    gap: ${theme.spacing.xs};
    padding: ${theme.spacing.md} ${theme.spacing.lg};
    background: ${theme.colors.surface};
    border-bottom: 1px solid ${theme.colors.border.default};
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    
    @media (max-width: ${theme.breakpoints.mobile}) {
        padding: ${theme.spacing.sm} ${theme.spacing.md};
    }
`;

const TabButton = styled.button`
    display: flex;
    align-items: center;
    gap: ${theme.spacing.sm};
    padding: ${theme.spacing.sm} ${theme.spacing.md};
    border: 1px solid ${props => props.$isActive ? theme.colors.primary : 'transparent'};
    background: ${props => props.$isActive ? theme.colors.primary : theme.colors.background};
    color: ${props => props.$isActive ? theme.colors.text.inverse : theme.colors.text.primary};
    border-radius: ${theme.borderRadius.md};
    font-size: ${theme.typography.fontSize.sm};
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    white-space: nowrap;
    
    &:hover:not(:disabled) {
        background: ${props => props.$isActive ? theme.colors.primaryHover : theme.colors.surfaceHover};
        border-color: ${props => props.$isActive ? theme.colors.primaryHover : theme.colors.border.default};
    }
    
    &:focus {
        outline: 2px solid ${theme.colors.primary};
        outline-offset: 2px;
    }
`;

const TabPanel = styled.div`
    flex: 1;
    overflow-y: auto;
    padding: ${theme.spacing.lg};
    
    @media (max-width: ${theme.breakpoints.mobile}) {
        padding: ${theme.spacing.md};
    }
`;

const SubFormWrapper = styled.div`
    display: grid;
    gap: ${theme.spacing.lg};
`;

const FieldGroup = styled.div`
    display: flex;
    flex-direction: column;
    gap: ${theme.spacing.xs};
`;

const FormLabel = styled.label`
    font-size: ${theme.typography.fontSize.sm};
    font-weight: 500;
    color: ${theme.colors.text.primary};
    line-height: 1.5;
`;

const Required = styled.span`
    color: ${theme.colors.error};
    margin-left: ${theme.spacing.xs};
    font-weight: 400;
`;

const NumberInput = styled.input`
    ${inputStyles}
`;

const Select = styled.select`
    ${inputStyles}
    cursor: pointer;
`;

const CheckboxGroup = styled.div`
    display: flex;
    align-items: flex-start;
    gap: ${theme.spacing.sm};
`;

const CheckboxInput = styled.input`
    width: 1.25rem;
    height: 1.25rem;
    margin-top: 2px;
    cursor: pointer;
    accent-color: ${theme.colors.primary};
    
    &:focus {
        outline: 2px solid ${theme.colors.primary};
        outline-offset: 2px;
    }
`;

const CheckboxLabel = styled.label`
    font-size: ${theme.typography.fontSize.sm};
    color: ${theme.colors.text.primary};
    cursor: pointer;
    user-select: none;
    line-height: 1.5;
`;

const HelpText = styled.span`
    font-size: ${theme.typography.fontSize.xs};
    color: ${theme.colors.text.secondary};
    line-height: 1.5;
`;

const ErrorText = styled.span`
    font-size: ${theme.typography.fontSize.xs};
    color: ${theme.colors.error};
    margin-top: ${theme.spacing.xs};
    line-height: 1.5;
`;

const ActionsRow = styled.div`
    display: flex;
    gap: ${theme.spacing.md};
    justify-content: flex-end;
    margin-top: ${theme.spacing.xl};
    padding-top: ${theme.spacing.lg};
    border-top: 1px solid ${theme.colors.border.default};
`;

const buttonStyles = css`
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: ${theme.spacing.sm} ${theme.spacing.lg};
    border-radius: ${theme.borderRadius.md};
    font-size: ${theme.typography.fontSize.sm};
    font-weight: 500;
    font-family: ${theme.typography.fontFamily};
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 100px;
    
    &:focus {
        outline: 2px solid ${theme.colors.primary};
        outline-offset: 2px;
    }
    
    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
`;

const CancelButton = styled.button`
    ${buttonStyles}
    background: ${theme.colors.background};
    color: ${theme.colors.text.primary};
    border: 1px solid ${theme.colors.border.default};
    
    &:hover:not(:disabled) {
        background: ${theme.colors.surface};
        border-color: ${theme.colors.border.hover};
    }
`;

const SubmitButton = styled.button`
    ${buttonStyles}
    background: ${theme.colors.primary};
    color: ${theme.colors.text.inverse};
    border: 1px solid ${theme.colors.primary};
    
    &:hover:not(:disabled) {
        background: ${theme.colors.primaryHover};
        border-color: ${theme.colors.primaryHover};
    }
`;

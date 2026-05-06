"use client";
import React, { useState, useCallback, useMemo, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { X, CreditCard, Calendar, Clock, RotateCcw, ChevronRight, Info, Euro, CheckCircle } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { useOffersApi } from '../../hooks/useOffersApi';
import { useTranslations } from 'next-intl';
import { toast } from 'react-toastify';
import styles from './OfferModal.module.css';

/**
 * Premium OfferModal - Elegant, compact, and professional offer creation
 */
const OfferModal = ({ isOpen, onClose, service, onSuccess }) => {
    const t = useTranslations('Offers');
    const { user } = useAuth();
    const {
        createOffer,
        createBuyNow,
        createLease,
        createReservation,
        createBuyBack,
        calculateSuggestedPrices,
        isCreatingOffer,
        processingState
    } = useOffersApi();

    const [selectedType, setSelectedType] = useState(null);
    const [formData, setFormData] = useState({});
    const [step, setStep] = useState('type'); // 'type' | 'details' | 'confirm'
    const [mounted, setMounted] = useState(false);

    // Handle client-side mounting for portal
    useEffect(() => {
        setMounted(true);
        return () => setMounted(false);
    }, []);

    // Premium offer type configurations - Clean, accessible design
    const offerTypes = useMemo(() => [
        {
            type: 'BuyNow',
            title: t('buyNow', 'Sofort kaufen'),
            subtitle: t('buyNowDesc', 'Sofort zum angebotenen Preis kaufen'),
            icon: CreditCard,
            color: '#059669',
            bgColor: '#f0fdf4',
            borderColor: '#bbf7d0',
            recommended: true
        },
        {
            type: 'Reservation',
            title: t('reservation', 'Reservierung'),
            subtitle: t('reservationDesc', 'Den Service für später reservieren'),
            icon: Clock,
            color: '#d97706',
            bgColor: '#fffbeb',
            borderColor: '#fed7aa'
        },
        {
            type: 'Lease',
            title: t('lease', 'Mieten'),
            subtitle: t('leaseDesc', 'Den Service für einen bestimmten Zeitraum mieten'),
            icon: Calendar,
            color: '#2563eb',
            bgColor: '#eff6ff',
            borderColor: '#bfdbfe'
        },
        {
            type: 'BuyBack',
            title: t('buyBack', 'Rückkauf'),
            subtitle: t('buyBackDesc', 'Angebot zum Rückkauf nach einer Periode'),
            icon: RotateCcw,
            color: '#7c3aed',
            bgColor: '#f5f3ff',
            borderColor: '#c4b5fd'
        }
    ], [t]);

    // Get current offer type config
    const currentTypeConfig = useMemo(() => 
        offerTypes.find(ot => ot.type === selectedType),
        [offerTypes, selectedType]
    );

    // Calculate suggested prices
    const suggestedPrices = useMemo(() => {
        const basePrice = parseFloat(service?.basePrice || service?.hourlyRate || 0);
        return calculateSuggestedPrices(basePrice, selectedType);
    }, [service, selectedType, calculateSuggestedPrices]);

    // Reset form when modal opens/closes
    useEffect(() => {
        if (!isOpen) {
            setSelectedType(null);
            setFormData({});
            setStep('type');
        }
    }, [isOpen]);

    // Handle type selection
    const handleTypeSelect = useCallback((type) => {
        setSelectedType(type);
        setStep('details');
        
        // Pre-fill with suggested values
        const suggestions = calculateSuggestedPrices(
            parseFloat(service?.basePrice || service?.hourlyRate || 0), 
            type
        );
        setFormData(suggestions);
    }, [service, calculateSuggestedPrices]);

    // Handle form data update
    const updateFormData = useCallback((field, value) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    }, []);

    // Handle form submission
    const handleSubmit = useCallback(async () => {
        if (!user?.userId) {
            toast.error(t('loginRequired', 'Bitte loggen Sie sich ein, um ein Angebot zu senden'));
            return;
        }

        try {
            // First create the base offer
            const baseOffer = await createOffer({
                userSellerId: user.userId,
                productId: service.id,
                price: formData.finalPrice || formData.lockedPrice || formData.monthlyPrice || 0
            });

            // Then create the specific offer type
            let result;
            switch (selectedType) {
                case 'BuyNow':
                    result = await createBuyNow({
                        offerId: baseOffer.id,
                        finalPrice: formData.finalPrice
                    });
                    break;
                case 'Lease':
                    result = await createLease({
                        offerId: baseOffer.id,
                        monthlyPrice: formData.monthlyPrice,
                        leaseTermMonths: formData.leaseTermMonths || 12,
                        hasBuyout: formData.hasBuyout || false,
                        buyoutPrice: formData.buyoutPrice
                    });
                    break;
                case 'Reservation':
                    result = await createReservation({
                        offerId: baseOffer.id,
                        lockedPrice: formData.lockedPrice,
                        reservationFee: formData.reservationFee,
                        lockDurationDays: formData.lockDurationDays || 14,
                        lockBuyerId: user.userId
                    });
                    break;
                case 'BuyBack':
                    result = await createBuyBack({
                        offerId: baseOffer.id,
                        lockedPrice: formData.lockedPrice,
                        redemptionFee: formData.redemptionFee,
                        lockDurationDays: formData.lockDurationDays || 30,
                        lockBuyerId: user.userId
                    });
                    break;
                default:
                    throw new Error('Invalid offer type');
            }

            toast.success(t('offerSent', 'Angebot erfolgreich gesendet!'));
            
            if (onSuccess) {
                onSuccess({
                    type: selectedType,
                    result,
                    formData
                });
            }

            onClose();
        } catch (error) {
            // Error: 'Error creating offer:', error...
            toast.error(t('offerError', 'Fehler beim Senden des Angebots'));
        }
    }, [user, service, selectedType, formData, createOffer, createBuyNow, createLease, createReservation, createBuyBack, onSuccess, onClose, t]);

    // Don't render on server or when not open
    if (!mounted || !isOpen) return null;

    // Ensure document.body exists before creating portal
    const portalTarget = typeof document !== 'undefined' && document.body;
    if (!portalTarget) return null;

    const renderTypeSelection = () => (
        <div className={styles.stepContent}>
            <div className={styles.stepHeader}>
                <h3 className={styles.stepTitle}>{t('selectOfferType', 'Wählen Sie Ihre Angebotsart:')}</h3>
                <p className={styles.stepSubtitle}>{t('typeDescription', 'Wählen Sie die für Sie passende Option')}</p>
            </div>
            
            <div className={styles.typeGrid}>
                {offerTypes.map(({ type, title, subtitle, icon: Icon, color, bgColor, borderColor, recommended }) => (
                    <button
                        key={type}
                        className={`${styles.typeCard} ${recommended ? styles.recommended : ''}`}
                        onClick={() => handleTypeSelect(type)}
                        style={{ 
                            '--accent-color': color, 
                            '--bg-color': bgColor,
                            '--border-color': borderColor
                        }}
                        aria-label={`${title} - ${subtitle}`}
                        role="option"
                        aria-selected="false"
                    >
                        {recommended && (
                            <div className={styles.recommendedBadge} aria-label={t('recommended', 'Empfohlen')}>
                                <CheckCircle size={14} aria-hidden="true" />
                                <span>{t('recommended', 'Empfohlen')}</span>
                            </div>
                        )}
                        
                        <div className={styles.typeIcon} aria-hidden="true">
                            <Icon size={24} />
                        </div>
                        
                        <div className={styles.typeContent}>
                            <h4 className={styles.typeTitle}>{title}</h4>
                            <p className={styles.typeSubtitle}>{subtitle}</p>
                        </div>
                        
                        <ChevronRight size={20} className={styles.typeArrow} aria-hidden="true" />
                    </button>
                ))}
            </div>
        </div>
    );

    const renderDetailsForm = () => (
        <div className={styles.stepContent}>
            <div className={styles.stepHeader}>
                <div className={styles.selectedTypeInfo}>
                    <div 
                        className={styles.selectedIcon} 
                        style={{ 
                            backgroundColor: currentTypeConfig?.bgColor,
                            color: currentTypeConfig?.color,
                            border: `2px solid ${currentTypeConfig?.borderColor}`
                        }}
                        aria-hidden="true"
                    >
                        <currentTypeConfig.icon size={20} />
                    </div>
                    <div>
                        <h3 className={styles.stepTitle} id="selected-offer-type">{currentTypeConfig?.title}</h3>
                        <p className={styles.stepSubtitle}>{currentTypeConfig?.subtitle}</p>
                    </div>
                </div>
                <button 
                    className={styles.changeButton}
                    onClick={() => setStep('type')}
                    aria-describedby="selected-offer-type"
                    type="button"
                >
                    {t('change', 'Ändern')}
                </button>
            </div>

            <div className={styles.formGrid}>
                {selectedType === 'BuyNow' && (
                    <div className={styles.formGroup}>
                        <label className={styles.label}>
                            <Euro size={16} />
                            {t('finalPrice', 'Endpreis')}
                        </label>
                        <input
                            type="number"
                            step="0.01"
                            value={formData.finalPrice || ''}
                            onChange={(e) => updateFormData('finalPrice', parseFloat(e.target.value))}
                            className={styles.input}
                            placeholder="0.00"
                        />
                        {suggestedPrices.finalPrice && (
                            <div className={styles.suggestion}>
                                <Info size={14} />
                                Vorschlag: €{suggestedPrices.finalPrice}
                            </div>
                        )}
                    </div>
                )}

                {selectedType === 'Lease' && (
                    <>
                        <div className={styles.formGroup}>
                            <label className={styles.label}>
                                <Euro size={16} />
                                {t('monthlyPrice', 'Monatspreis')}
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                value={formData.monthlyPrice || ''}
                                onChange={(e) => updateFormData('monthlyPrice', parseFloat(e.target.value))}
                                className={styles.input}
                                placeholder="0.00"
                            />
                        </div>
                        
                        <div className={styles.formGroup}>
                            <label className={styles.label}>
                                <Calendar size={16} />
                                {t('leaseTerm', 'Laufzeit (Monate)')}
                            </label>
                            <select
                                value={formData.leaseTermMonths || 12}
                                onChange={(e) => updateFormData('leaseTermMonths', parseInt(e.target.value))}
                                className={styles.select}
                            >
                                <option value={6}>6 Monate</option>
                                <option value={12}>12 Monate</option>
                                <option value={24}>24 Monate</option>
                                <option value={36}>36 Monate</option>
                            </select>
                        </div>
                        
                        <div className={styles.checkboxGroup}>
                            <label className={styles.checkbox}>
                                <input
                                    type="checkbox"
                                    checked={formData.hasBuyout || false}
                                    onChange={(e) => updateFormData('hasBuyout', e.target.checked)}
                                />
                                <span className={styles.checkboxLabel}>{t('includeBuyout', 'Kaufoption einschließen')}</span>
                            </label>
                        </div>
                        
                        {formData.hasBuyout && (
                            <div className={styles.formGroup}>
                                <label className={styles.label}>
                                    <Euro size={16} />
                                    {t('buyoutPrice', 'Kaufpreis')}
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.buyoutPrice || ''}
                                    onChange={(e) => updateFormData('buyoutPrice', parseFloat(e.target.value))}
                                    className={styles.input}
                                    placeholder="0.00"
                                />
                            </div>
                        )}
                    </>
                )}

                {(selectedType === 'Reservation' || selectedType === 'BuyBack') && (
                    <>
                        <div className={styles.formGroup}>
                            <label className={styles.label}>
                                <Euro size={16} />
                                {selectedType === 'Reservation' ? t('reservationFee', 'Reservierungsgebühr') : t('lockedPrice', 'Festpreis')}
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                value={selectedType === 'Reservation' ? formData.reservationFee || '' : formData.lockedPrice || ''}
                                onChange={(e) => updateFormData(
                                    selectedType === 'Reservation' ? 'reservationFee' : 'lockedPrice', 
                                    parseFloat(e.target.value)
                                )}
                                className={styles.input}
                                placeholder="0.00"
                            />
                        </div>
                        
                        <div className={styles.formGroup}>
                            <label className={styles.label}>
                                <Clock size={16} />
                                {t('duration', 'Dauer (Tage)')}
                            </label>
                            <select
                                value={formData.lockDurationDays || (selectedType === 'Reservation' ? 14 : 30)}
                                onChange={(e) => updateFormData('lockDurationDays', parseInt(e.target.value))}
                                className={styles.select}
                            >
                                <option value={7}>7 Tage</option>
                                <option value={14}>14 Tage</option>
                                <option value={30}>30 Tage</option>
                                <option value={60}>60 Tage</option>
                            </select>
                        </div>
                        
                        {selectedType === 'BuyBack' && (
                            <div className={styles.formGroup}>
                                <label className={styles.label}>
                                    <Euro size={16} />
                                    {t('redemptionFee', 'Rückkaufgebühr')}
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.redemptionFee || ''}
                                    onChange={(e) => updateFormData('redemptionFee', parseFloat(e.target.value))}
                                    className={styles.input}
                                    placeholder="0.00"
                                />
                            </div>
                        )}
                    </>
                )}
            </div>

            <div className={styles.stepActions}>
                <button
                    className={styles.backButton}
                    onClick={() => setStep('type')}
                >
                    {t('back', 'Zurück')}
                </button>
                <button
                    className={styles.continueButton}
                    onClick={() => setStep('confirm')}
                    disabled={!Object.keys(formData).length}
                >
                    {t('continue', 'Weiter')}
                    <ChevronRight size={16} />
                </button>
            </div>
        </div>
    );

    const renderConfirmation = () => (
        <div className={styles.stepContent}>
            <div className={styles.confirmHeader}>
                <div 
                    className={styles.confirmIcon} 
                    style={{ 
                        backgroundColor: currentTypeConfig?.bgColor,
                        color: currentTypeConfig?.color,
                        border: `2px solid ${currentTypeConfig?.borderColor}`
                    }}
                    aria-hidden="true"
                >
                    <currentTypeConfig.icon size={24} />
                </div>
                <h3 className={styles.confirmTitle} id="confirm-title">{t('confirmOffer', 'Angebot bestätigen')}</h3>
                <p className={styles.confirmSubtitle}>{t('reviewDetails', 'Überprüfen Sie Ihre Angaben')}</p>
            </div>

            <div className={styles.confirmDetails}>
                <div className={styles.serviceInfo}>
                    <h4>{service.name}</h4>
                    <p>Listed at: €{service.basePrice || service.hourlyRate || '0'}/hr</p>
                </div>

                <div className={styles.offerSummary}>
                    <div className={styles.summaryRow}>
                        <span>{t('offerType', 'Angebotsart')}</span>
                        <strong>{currentTypeConfig?.title}</strong>
                    </div>
                    
                    {Object.entries(formData)
                        .filter(([key]) => !key.startsWith('suggested')) // Filter out suggested values
                        .map(([key, value]) => (
                        <div key={key} className={styles.summaryRow}>
                            <span>{t(key, key)}</span>
                            <strong>
                                {typeof value === 'boolean' ? (value ? t('yes', 'Ja') : t('no', 'Nein')) :
                                 typeof value === 'number' ? `€${value.toFixed(2)}` : value}
                                {key.includes('Days') && ` ${t('days', 'Tage')}`}
                                {key.includes('Months') && ` ${t('months', 'Monate')}`}
                            </strong>
                        </div>
                    ))}
                </div>
            </div>

            <div className={styles.stepActions}>
                <button
                    className={styles.backButton}
                    onClick={() => setStep('details')}
                >
                    {t('back', 'Zurück')}
                </button>
                <button
                    className={styles.submitButton}
                    onClick={handleSubmit}
                    disabled={isCreatingOffer || processingState[service.id]}
                >
                    {isCreatingOffer ? t('sending', 'Wird gesendet...') : t('sendOffer', 'Angebot senden')}
                </button>
            </div>
        </div>
    );

    return createPortal(
        <div className={styles.modalOverlay} onClick={onClose}>
            <div className={styles.modalContent} onClick={e => e.stopPropagation()}>
                <div className={styles.modalHeader}>
                    <div className={styles.headerContent}>
                        <h2 className={styles.modalTitle}>{t('makeOffer', 'Angebot machen')}</h2>
                        <div className={styles.stepIndicator}>
                            <div className={`${styles.step} ${step === 'type' ? styles.active : styles.completed}`}>1</div>
                            <div className={styles.stepLine} />
                            <div className={`${styles.step} ${step === 'details' ? styles.active : step === 'confirm' ? styles.completed : ''}`}>2</div>
                            <div className={styles.stepLine} />
                            <div className={`${styles.step} ${step === 'confirm' ? styles.active : ''}`}>3</div>
                        </div>
                    </div>
                    <button className={styles.closeButton} onClick={onClose}>
                        <X size={20} />
                    </button>
                </div>

                {step === 'type' && renderTypeSelection()}
                {step === 'details' && renderDetailsForm()}
                {step === 'confirm' && renderConfirmation()}
            </div>
        </div>,
        document.body
    );
};

export default OfferModal;
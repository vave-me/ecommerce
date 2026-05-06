// src/features/CreateServiceModal/index.jsx
// REVERTED: API call logic and data types reverted to match the older monolithic version provided.
// NOTE: This might NOT match the actual API spec (e.g., numeric types).
// FIXED: Added direct null check for 'user' object before accessing userId
// FIXED: Corrected import paths for child/shared components

import React, {useCallback, useEffect, useState, useMemo, Suspense, memo} from 'react';
import PropTypes from 'prop-types';
import styles from './ServiceModal.module.css'; // Assuming shared CSS Module
import {useAuth} from "../../context/AuthContext";
// Assuming updateService exists, though calls might be commented out now
import {addService, updateService} from "../../api/client/servicesApi";
import {createMedia} from "../../api/client/mediaApi";
import {fetchMainCategories} from "../../api/categories";
import {useCategories} from "../../hooks/useCategories";

// Lazy load step components - using named exports
const BasicInfoStep = React.lazy(() => import('./components/steps/BasicInfoStep/BasicInfoStep').then(module => ({default: module.BasicInfoStep})));
const MediaUploadStep = React.lazy(() => import('../shared/components/MediaUploadStep/MediaUploadStep').then(module => ({default: module.MediaUploadStep})));
const OptionalInfoStep = React.lazy(() => import('./components/steps/OptionalSettingsStep/OptionalSettingsStep').then(module => ({default: module.OptionalInfoStep})));
const FinalizeStep = React.lazy(() => import('./components/steps/FinalizeStep/FinalizeStep').then(module => ({default: module.FinalizeStep})));

// --- Import Shared UI Components (Using correct paths from shared directory) ---
import {ModalOverlay} from "../shared/modal/ModalOverlay";
import {ModalContainer} from "../shared/modal/ModalContainer";
import {StepNavigation} from "../shared/StepNavigation/StepNavigation";
import {ModalHeader} from "../shared/ModalHeader/ModalHeader";
import {ErrorAlert} from "../../common/components/ErrorAlert";

// Import Hooks (ensure paths are correct)
import {useFocusTrap} from "../../hooks/useFocusTrap";
// Assuming useAutosave hook is available if needed
// import { useAutosave } from "../../hooks/useAutoSave";

// Loading placeholder for Suspense fallback
const StepLoadingFallback = memo(() => (
    <div className={styles.loadingContainer}>
        <div className={styles.loadingSpinner}></div>
        <p>Loading...</p>
    </div>
));

StepLoadingFallback.displayName = 'StepLoadingFallback';

// Step labels for navigation
// Step labels for navigation
const STEPS = [
    {label: 'Basic Info', key: 'basicInfo'},
    {label: 'Media', key: 'media'},
    {label: 'Details', key: 'optional'},
    {label: 'Publish', key: 'finalize'}
];

// Utility functions for type safety
const safeParseNumber = (value) => {
    const num = parseFloat(value);
    return isNaN(num) ? 0 : num;
};

const safeParseInt = (value) => {
    const num = parseInt(value, 10);
    return isNaN(num) ? 0 : num;
};

// Centralized payload building function matching API spec exactly
const buildServicePayload = (accumulatedData, user) => {
    const { basicInfo, media, optionalInfo, location } = accumulatedData;
    
    // Format arrays properly
    const tagsArray = Array.isArray(basicInfo.tags)
        ? basicInfo.tags
        : (typeof basicInfo.tags === 'string' && basicInfo.tags)
            ? basicInfo.tags.split(',').map(tag => tag.trim()).filter(Boolean)
            : [];
    
    const pricingArray = Array.isArray(basicInfo.pricing)
        ? basicInfo.pricing
        : (typeof basicInfo.pricing === 'string' && basicInfo.pricing)
            ? basicInfo.pricing.split(',').map(price => price.trim()).filter(Boolean)
            : [];
    
    const qualificationsArray = Array.isArray(basicInfo.qualifications)
        ? basicInfo.qualifications
        : (typeof basicInfo.qualifications === 'string' && basicInfo.qualifications)
            ? basicInfo.qualifications.split(',').map(qual => qual.trim()).filter(Boolean)
            : [];
    
    // Build payload matching API spec exactly (no userId field in API)
    const payload = {
        // Basic Information
        name: basicInfo.name || '',
        description: basicInfo.description || '',
        basePrice: basicInfo.basePrice?.toString() || '0', // API expects string with int64 format
        serviceType: basicInfo.serviceType || '',
        pricing: pricingArray,
        availability: basicInfo.availability || '',
        providerName: basicInfo.providerName || '',
        
        // Category and Classification
        categoryId: basicInfo.categoryId || '',
        categorySlug: basicInfo.categorySlug || '',
        
        // Service Details
        descriptionShort: basicInfo.descriptionShort || '',
        descriptionLong: basicInfo.descriptionLong || '',
        qualifications: qualificationsArray,
        contact: basicInfo.contact || '',
        faq: basicInfo.faq || '',
        tags: tagsArray,
        status: 'active', // Default status
        userType: basicInfo.userType || 'private',
        shippingCost: optionalInfo?.shippingCost?.toString() || '0', // API expects string with int64 format
        hasVariants: !!basicInfo.hasVariants,
        middlemanService: !!basicInfo.middlemanService,
        negotiable: !!basicInfo.negotiable,
        
        // Attributes and Options - matching API definitions
        attributes: optionalInfo?.attributes || [],
        options: optionalInfo?.options || [],
        
        // Media
        thumbnail: media?.thumbnail || '',
        
        // Location - API expects numbers with float format
        lat: safeParseNumber(basicInfo?.lat) || location?.lat || 0,
        lng: safeParseNumber(basicInfo?.lng) || location?.lng || 0
    };
    
    return payload;
};

const CreateServiceModal = memo(function CreateServiceModal({onClose, editMode = false, initialServiceData = null}) {
    // --- Core State ---
    const [currentStep, setCurrentStep] = useState(1);
    const [lastCompletedStep, setLastCompletedStep] = useState(0); // Start at 0 for new
    const [success, setSuccess] = useState(false);

    // --- Shared IDs & Accumulated Data ---
    const [serviceId, setServiceId] = useState(initialServiceData?.id || null);
    const [mediaId, setMediaId] = useState(initialServiceData?.mediaId || null);
    const [accumulatedData, setAccumulatedData] = useState({
        // Initialize based on edit mode and initial data presence
        basicInfo: editMode && initialServiceData ? {
            name: initialServiceData.name,
            description: initialServiceData.description,
            basePrice: initialServiceData.basePrice,
            serviceType: initialServiceData.serviceType,
            pricing: initialServiceData.pricing || [],
            availability: initialServiceData.availability,
            providerName: initialServiceData.providerName,
            categoryId: initialServiceData.categoryId,
            categorySlug: initialServiceData.categorySlug,
            descriptionShort: initialServiceData.descriptionShort,
            descriptionLong: initialServiceData.descriptionLong,
            qualifications: initialServiceData.qualifications || [],
            contact: initialServiceData.contact,
            faq: initialServiceData.faq,
            tags: Array.isArray(initialServiceData.tags) ? initialServiceData.tags : [],
            userType: initialServiceData.userType || "private",
            negotiable: initialServiceData.negotiable,
            middlemanService: initialServiceData.middlemanService || false,
            hasVariants: initialServiceData.hasVariants || false,
            attributes: initialServiceData.attributes || [],
            options: initialServiceData.options || [],
        } : {
            name: "",
            description: "",
            basePrice: "",
            serviceType: "",
            pricing: [],
            availability: "",
            providerName: "",
            categoryId: "",
            categorySlug: "",
            descriptionShort: "",
            descriptionLong: "",
            qualifications: [],
            contact: "",
            faq: "",
            tags: [],
            userType: "private",
            negotiable: false,
            middlemanService: false,
            hasVariants: false,
            attributes: [],
            options: [],
        },
        media: editMode && initialServiceData ? {
            images: initialServiceData.images || [],
            videoUrl: initialServiceData.videoUrl,
            thumbnail: initialServiceData.thumbnail,
        } : {
            images: [],
            videoUrl: "",
            thumbnail: "",
        },
        optionalInfo: editMode && initialServiceData ? {
            weight: initialServiceData.weight,
            height: initialServiceData.height,
            width: initialServiceData.width,
            depth: initialServiceData.depth,
            manageStocks: initialServiceData.manageStocks,
            stock: initialServiceData.stock,
            shippingCost: initialServiceData.shippingCost,
        } : {
            weight: "",
            height: "",
            width: "",
            depth: "",
            manageStocks: false,
            stock: "",
            shippingCost: "",
        },
        location: editMode && initialServiceData ? {
            lat: initialServiceData.lat,
            lng: initialServiceData.lng,
        } : {
            lat: null,
            lng: null,
        },
    });

    // --- UI State ---
    const [errors, setErrors] = useState({});
    const [isLoading, setIsLoading] = useState(false);
    const [categories, setCategories] = useState([]);

    // --- Hooks ---
    const {user} = useAuth();
    const focusTrapRef = useFocusTrap(true);
    const isUserLoggedIn = useMemo(() => !!user && !!user.userId, [user]);
    // const { lastSaved, isSaving } = useAutosave(user?.userId, accumulatedData, dealId);

    // Use the categories hook to fetch and cache categories
    const {
        data: categoriesData = [],
        isLoading: isCategoriesLoading,
        error: categoriesError
    } = useCategories('service');

    // Log any category fetch errors
    useEffect(() => {
        if (categoriesError) {
            // Error: "Error fetching service categories:", categoriesEr...
            setErrors(prev => ({...prev, categories: "Failed to load categories"}));
        }
    }, [categoriesError]);
    
    // --- Status for sidebar indicator ---
    const getStatusText = () => {
        if (success) return 'Published';
        if (isLoading) return 'Saving...';
        return 'Draft';
    };
    
    const getStatusClass = () => {
        if (success) return styles.published;
        if (isLoading) return styles.saving;
        return '';
    };

    // --- Step Navigation ---
    const handleNextStep = useCallback(() => {
        setCurrentStep(prev => Math.min(prev + 1, 4));
    }, []);
    const handlePrevStep = useCallback(() => {
        setCurrentStep(prev => Math.max(prev - 1, 1));
    }, []);
    const handleStepClick = useCallback((step) => {
        if (step <= Math.min(lastCompletedStep + 1, 4)) {
            setCurrentStep(step);
        }
    }, [lastCompletedStep]);

    // --- Step Completion Callbacks ---

    // Basic Info (Step 1 -> 2) - Following ProductModal pattern
    const handleBasicInfoComplete = useCallback(async (basicInfo) => {
        // 1. Check Login Status
        if (!isUserLoggedIn || !user) {
            setErrors({submit: 'You must be logged in to create or update a listing.'});
            return;
        }

        // 2. Start Loading State & Clear Errors
        setIsLoading(true);
        setErrors({});

        let currentServiceId = serviceId;
        let currentMediaId = mediaId;

        try {
            // 3. Store Step 1 Data Locally
            setAccumulatedData(prev => ({...prev, basicInfo}));

            // 4. Prepare API Payload using centralized builder
            const servicePayload = buildServicePayload({
                ...accumulatedData,
                basicInfo
            }, user);

            // 5. Call API: Create Service (following ProductModal pattern)
            const addResp = await addService(servicePayload);
            if (!addResp?.id) throw new Error("No service ID returned from backend");
            currentServiceId = addResp.id;
            setServiceId(currentServiceId);

            // 6. Create Media Container (if needed)
            if (currentServiceId && !currentMediaId && user?.userId) {
                const mediaResp = await createMedia({
                    itemId: currentServiceId,
                    itemType: "service",
                    userId: user.userId,
                });
                if (!mediaResp?.id) {
                    
                } else {
                    currentMediaId = mediaResp.id;
                    setMediaId(currentMediaId);
                }
            }

            // 7. Advance Step
            setLastCompletedStep(1);
            handleNextStep();

        } catch (err) {
            // Error: "Error during basic info step:", err...
            setErrors({submit: err.response?.data?.message || "Failed to save basic info. Please try again."});
        } finally {
            setIsLoading(false);
        }
    }, [user, serviceId, mediaId, editMode, handleNextStep, isUserLoggedIn, accumulatedData]);

    // Media Upload (Step 2 -> 3) - Following ProductModal pattern
    const handleMediaComplete = useCallback(async (mediaData) => {
        try {
            setIsLoading(true);
            setErrors({});
            
            if (!serviceId) {
                throw new Error("Service ID is required for media upload");
            }
            
            let updatedMediaData = {...mediaData};
            
            // Set thumbnail based on first image if none selected
            if (updatedMediaData?.images?.length > 0 && !updatedMediaData?.thumbnail) {
                updatedMediaData.thumbnail = updatedMediaData.images[0];
            }
            
            // Update service with media information
            const servicePayload = buildServicePayload({
                ...accumulatedData,
                media: updatedMediaData
            }, user);
            
            await updateService({...servicePayload, id: serviceId});
            
            // Update accumulated data
            setAccumulatedData(prev => ({
                ...prev,
                media: updatedMediaData
            }));
            
            // Mark step as completed and move to next
            setLastCompletedStep(prev => Math.max(prev, 2));
            handleNextStep();
            
        } catch (error) {
            // Error: "Error during media step:", error...
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to save media'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, serviceId, handleNextStep]);

    // Optional Settings Step (3 -> 4) - Following ProductModal pattern
    const handleOptionalComplete = useCallback(async (optionalData) => {
        try {
            setIsLoading(true);
            setErrors({});
            
            if (!serviceId) {
                throw new Error("Service ID is required");
            }
            
            // Update service with optional settings
            const servicePayload = buildServicePayload({
                ...accumulatedData,
                optionalInfo: optionalData
            }, user);
            
            await updateService({...servicePayload, id: serviceId});
            
            // Update accumulated data
            setAccumulatedData(prev => ({
                ...prev,
                optionalInfo: optionalData
            }));
            
            // Mark step as completed and move to next
            setLastCompletedStep(prev => Math.max(prev, 3));
            handleNextStep();
            
        } catch (error) {
            // Error: "Error during optional info step:", error...
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to save settings'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, serviceId, handleNextStep]);

    // Final Step - Publish (Following ProductModal pattern)
    const handleFinalizeComplete = useCallback(async (finalData) => {
        try {
            setIsLoading(true);
            setErrors({});
            
            if (!serviceId) {
                throw new Error("Service ID is required for publishing");
            }
            
            // Final update with publish status - build complete payload
            const servicePayload = buildServicePayload(accumulatedData, user);
            
            // Override status to active for publishing
            const finalPayload = {
                ...servicePayload,
                id: serviceId, // Add ID for update
                status: 'active' // Change from draft to active
            };
            
            await updateService(finalPayload);
            
            // Mark as completed and show success
            setLastCompletedStep(4);
            setSuccess(true);
            
        } catch (error) {
            // Error: "Error finalizing service:", error...
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to publish service'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, serviceId]);
    
    // Success state completion
    const handleSuccess = useCallback(() => {
        onClose();
    }, [onClose]);

    return (
        <div className={styles.modalOverlay} onClick={(e) => e.target === e.currentTarget && onClose()}>
            <div className={styles.modalContainer} ref={focusTrapRef}>
                {/* Professional Sidebar */}
                <div className={styles.sidebar}>
                    <div className={styles.logoContainer}>
                        <h1 className={styles.logoText}>Create Service</h1>
                        <div className={`${styles.autosaveIndicator} ${getStatusClass()}`}>
                            {getStatusText()}
                        </div>
                    </div>
                    {/* Step Navigation */}
                    <div className={styles.stepsContainer}>
                        {STEPS.map((step, index) => {
                            const stepNumber = index + 1;
                            const isActive = currentStep === stepNumber;
                            const isCompleted = lastCompletedStep >= stepNumber;
                            const isAccessible = stepNumber <= Math.min(lastCompletedStep + 1, 4);
                            return (
                                <div
                                    key={step.key}
                                    className={`
                                        ${styles.stepNavItem}
                                        ${isActive ? styles.stepNavActive : ''}
                                        ${isCompleted ? styles.stepNavCompleted : ''}
                                        ${!isAccessible ? styles.stepNavDisabled : ''}
                                    `}
                                    onClick={() => isAccessible && handleStepClick(stepNumber)}
                                >
                                    <div className={styles.stepNumCircle}>
                                        {isCompleted && !isActive ? '' : stepNumber}
                                    </div>
                                    <span className={styles.stepLabel}>{step.label}</span>
                                </div>
                            );
                        })}
                    </div>
                </div>
                {/* Main Content */}
                <div className={styles.content}>
                    <button 
                        className={styles.closeButton} 
                        onClick={onClose}
                        aria-label="Close modal"
                    >
                        ✕
                    </button>
                    {errors.submit && <ErrorAlert message={errors.submit}/>}

                    {/* Conditional Rendering of Step Components */}
                    {currentStep === 1 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <BasicInfoStep
                                initialData={accumulatedData.basicInfo}
                                onComplete={handleBasicInfoComplete}
                                onBack={onClose}
                                isLoading={isLoading}
                                errors={errors}
                                isUserLoggedIn={isUserLoggedIn}
                                categories={categoriesData}
                                isCategoriesLoading={isCategoriesLoading}
                            />
                        </Suspense>
                    )}
                    {/* Step 2: Media Upload */}
                    {currentStep === 2 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <MediaUploadStep
                                initialImages={accumulatedData.media?.images}
                                initialVideoUrl={accumulatedData.media?.videoUrl}
                                onComplete={handleMediaComplete}
                                onBack={handlePrevStep}
                                mediaId={mediaId}
                                isLoading={isLoading}
                                errors={errors}
                            />
                        </Suspense>
                    )}

                    {currentStep === 3 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <OptionalInfoStep
                                initialData={accumulatedData.optionalInfo}
                                onComplete={handleOptionalComplete}
                                onBack={handlePrevStep}
                                isLoading={isLoading}
                                errors={errors}
                                isUserLoggedIn={isUserLoggedIn}
                            />
                        </Suspense>
                    )}

                    {currentStep === 4 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <FinalizeStep
                                isSuccess={success}
                                onClose={handleSuccess}
                                onFinalize={handleFinalizeComplete}
                                isLoading={isLoading}
                                serviceData={accumulatedData}
                                styles={styles}
                            />
                        </Suspense>
                    )}
                </div>
            </div>
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for optimal performance
    return (
        prevProps.editMode === nextProps.editMode &&
        prevProps.initialServiceData?.id === nextProps.initialServiceData?.id &&
        prevProps.onClose === nextProps.onClose
    );
});

// --- PropTypes ---
CreateServiceModal.propTypes = {
    onClose: PropTypes.func.isRequired,
    editMode: PropTypes.bool,
    initialServiceData: PropTypes.shape({
        id: PropTypes.string,
        mediaId: PropTypes.string,
        name: PropTypes.string,
        description: PropTypes.string,
        basePrice: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        serviceType: PropTypes.string,
        pricing: PropTypes.arrayOf(PropTypes.string),
        availability: PropTypes.string,
        providerName: PropTypes.string,
        categoryId: PropTypes.string,
        categorySlug: PropTypes.string,
        descriptionShort: PropTypes.string,
        descriptionLong: PropTypes.string,
        qualifications: PropTypes.arrayOf(PropTypes.string),
        contact: PropTypes.string,
        faq: PropTypes.string,
        userType: PropTypes.oneOf(["private", "business", "freelancer", "agency", "consultant", "company"]),
        tags: PropTypes.arrayOf(PropTypes.string),
        images: PropTypes.arrayOf(PropTypes.string),
        videoUrl: PropTypes.string,
        thumbnail: PropTypes.string,
        shippingCost: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        middlemanService: PropTypes.bool,
        hasVariants: PropTypes.bool,
        negotiable: PropTypes.bool,
        attributes: PropTypes.arrayOf(PropTypes.shape({
            key: PropTypes.string,
            value: PropTypes.string
        })),
        options: PropTypes.arrayOf(PropTypes.shape({
            name: PropTypes.string,
            value: PropTypes.string,
            price: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
        })),
        lat: PropTypes.number,
        lng: PropTypes.number,
        status: PropTypes.oneOf(["draft", "active"]),
    }),
};

export default CreateServiceModal;


// src/features/CreateProductModal/index.jsx
// REVERTED: API call logic and data types reverted to match the older monolithic version provided.
// NOTE: This might NOT match the actual API spec (e.g., numeric types).
// FIXED: Added direct null check for 'user' object before accessing userId
// FIXED: Corrected import paths for child/shared components
import React, {useCallback, useEffect, useState, useMemo, Suspense, memo} from 'react';
import PropTypes from 'prop-types';
import styles from './ProductModal.module.css'; // Assuming shared CSS Module
import {useAuth} from "../../context/AuthContext";
import {addProduct, updateProduct} from "../../api/client/productsApi";
import {createMedia} from "../../api/client/mediaApi";
import {useCategories} from "../../hooks/useCategories";
// Lazy load step components for code splitting
const BasicInfoStep = React.lazy(() => 
    import('./components/steps/BasicInfoStep/BasicInfoStep').then(module => ({
        default: module.BasicInfoStep
    }))
);
const MediaUploadStep = React.lazy(() => 
    import('../shared/components/MediaUploadStep/MediaUploadStep').then(module => ({
        default: module.MediaUploadStep
    }))
);
const OptionalInfoStep = React.lazy(() => 
    import('./components/steps/OptionalSettingsStep/OptionalSettingsStep').then(module => ({
        default: module.OptionalInfoStep
    }))
);
const FinalizeStep = React.lazy(() => 
    import('./components/steps/FinalizeStep/FinalizeStep').then(module => ({
        default: module.FinalizeStep
    }))
);
// --- Import Shared UI Components (Using original paths from user context) ---
import {ModalOverlay} from "../shared/modal/ModalOverlay";
import {ModalContainer} from "../shared/modal/ModalContainer";
import {StepNavigation} from "../shared/StepNavigation/StepNavigation";
import {ModalHeader} from "../shared/ModalHeader/ModalHeader";
import {ErrorAlert} from "../../common/components/ErrorAlert";
// Import Hooks (ensure paths are correct)
import {useFocusTrap} from "../../hooks/useFocusTrap";
import {fetchMainCategories} from "../../api/categories";
// Assuming useAutosave hook is available if needed
// import { useAutosave } from "../../hooks/useAutoSave";
// Import shared FinalizeStep
import {FinalizeStep as SharedFinalizeStep} from "../shared/components/FinalizeStep/FinalizeStep";
// Loading placeholder for Suspense fallback
const StepLoadingFallback = memo(() => (
    <div className={styles.loadingContainer}>
        <div className={styles.loadingSpinner}></div>
        <p>Loading...</p>
    </div>
));
StepLoadingFallback.displayName = 'StepLoadingFallback';
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
// Centralized payload building function
const buildProductPayload = (accumulatedData, user) => {
    const { basicInfo, media, optional, location } = accumulatedData;
    // Console log for debugging location data
    const payload = {
        // Basic Information
        name: basicInfo.name || '',
        description: basicInfo.description || '',
        base_price: safeParseNumber(basicInfo.basePrice),
        product_price: safeParseNumber(basicInfo.productPrice),
        // Category and Classification
        category_id: basicInfo.categoryId || '',
        category_slug: basicInfo.categorySlug || '',
        category_name: basicInfo.categoryName || '',
        // Product Details
        condition: basicInfo.condition || 'new',
        brand: basicInfo.brand || '',
        model: basicInfo.model || '',
        sku: basicInfo.sku || '',
        negotiable: !!basicInfo.negotiable,
        user_type: basicInfo.userType || 'private',
        product_type: basicInfo.productType || '',
        // Tags handling
        tags: Array.isArray(basicInfo.tags)
            ? basicInfo.tags
            : (typeof basicInfo.tags === 'string' && basicInfo.tags)
                ? basicInfo.tags.split(',').map(tag => tag.trim()).filter(Boolean)
                : [],
        // Variants and Attributes
        has_variants: !!basicInfo.hasVariants,
        attributes: Array.isArray(optional?.attributes) ? optional.attributes : [],
        // Optional settings
        weight: safeParseNumber(optional?.weight),
        height: safeParseNumber(optional?.height),
        width: safeParseNumber(optional?.width),
        depth: safeParseNumber(optional?.depth),
        manage_stocks: !!optional?.manageStocks,
        stock: safeParseInt(optional?.stock),
        shipping_cost: safeParseNumber(optional?.shippingCost),
        middleman_service: !!optional?.middlemanService,
        // Media - Include ALL media fields
        images: Array.isArray(media?.images) ? media.images : [],
        thumbnail: media?.thumbnail || '',
        video_url: media?.videoUrl || '',
        // Location
        lat: safeParseNumber(basicInfo?.lat) || location?.lat || 0,
        lng: safeParseNumber(basicInfo?.lng) || location?.lng || 0,
        // User information
        user_id: user?.userId || user?.id,
        // Location data
        address: basicInfo?.address || '',
        // Status - ensure active/published status is included
        status: 'active', // Default to active
        published: false // Default to unpublished
    };
    // Console log the final payload with location data
    return payload;
};
const CreateProductModal = memo(function CreateProductModal({onClose, editMode = false, initialProductData = null}) {
    // --- Core State ---
    const [currentStep, setCurrentStep] = useState(1);
    const [lastCompletedStep, setLastCompletedStep] = useState(0); // Start at 0 for new
    const [success, setSuccess] = useState(false);
    // --- Shared IDs & Accumulated Data ---
    const [productId, setProductId] = useState(initialProductData?.id || null);
    const [mediaId, setMediaId] = useState(initialProductData?.mediaId || null);
    const [accumulatedData, setAccumulatedData] = useState({
        // Initialize based on edit mode and initial data presence
        basicInfo: editMode && initialProductData ? {
            name: initialProductData.name,
            description: initialProductData.description,
            basePrice: initialProductData.basePrice,
            productPrice: initialProductData.productPrice,
            categoryId: initialProductData.categoryId,
            categorySlug: initialProductData.categorySlug,
            categoryName: initialProductData.categoryName,
            condition: initialProductData.condition,
            brand: initialProductData.brand,
            model: initialProductData.model,
            negotiable: initialProductData.negotiable,
            userType: initialProductData.userType,
            sku: initialProductData.sku,
            tags: Array.isArray(initialProductData.tags) ? initialProductData.tags : [],
            hasVariants: initialProductData.hasVariants,
            attributes: initialProductData.attributes || [],
            productType: initialProductData.productType || '',
        } : {
            name: "",
            description: "",
            basePrice: "",
            productPrice: "",
            categoryId: "",
            categorySlug: "",
            categoryName: "",
            condition: "new",
            brand: "",
            model: "",
            negotiable: false,
            userType: "private",
            sku: "",
            tags: [],
            hasVariants: false,
            attributes: [],
            productType: "",
        },
        media: editMode && initialProductData ? {
            images: initialProductData.images || [],
            videoUrl: initialProductData.videoUrl,
            thumbnail: initialProductData.thumbnail,
        } : {
            images: [],
            videoUrl: "",
            thumbnail: "",
        },
        optional: editMode && initialProductData ? {
            weight: initialProductData.weight,
            height: initialProductData.height,
            width: initialProductData.width,
            depth: initialProductData.depth,
            manageStocks: initialProductData.manageStocks,
            stock: initialProductData.stock,
            shippingCost: initialProductData.shippingCost,
            middlemanService: initialProductData.middlemanService,
            attributes: initialProductData.attributes || [],
        } : {
            weight: "",
            height: "",
            width: "",
            depth: "",
            manageStocks: false,
            stock: "",
            shippingCost: "",
            middlemanService: false,
            attributes: [],
        },
        location: editMode && initialProductData ? {
            lat: initialProductData.lat,
            lng: initialProductData.lng,
        } : {
            lat: null,
            lng: null,
        },
    });
    // --- UI State ---
    const [errors, setErrors] = useState({});
    const [isLoading, setIsLoading] = useState(false);
    // --- Hooks ---
    const {user} = useAuth();
    const focusTrapRef = useFocusTrap(true);
    const isUserLoggedIn = useMemo(() => !!user && !!user.userId, [user]);
    // const { lastSaved, isSaving } = useAutosave(user?.userId, accumulatedData, productId);
    // Use the categories hook to fetch and cache categories
    const {
        data: categoriesData = [],
        isLoading: isCategoriesLoading,
        error: categoriesError
    } = useCategories('marketplace');
    // Log any category fetch errors
    useEffect(() => {
        if (categoriesError) {
            setErrors(prev => ({...prev, categories: "Failed to load categories"}));
        }
    }, [categoriesError]);
    // --- Status for sidebar indicator ---
    const getStatusText = () => {
        if (success) return 'Published';
        if (isLoading) return 'Saving...';
        return 'Active';
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
    // Basic Info (Step 1 -> 2) - Reverted to match old logic/types
    const handleBasicInfoComplete = useCallback(async (basicInfo) => {
        // Console log the incoming basicInfo data
        // 1. Check Login Status
        if (!isUserLoggedIn || !user) {
            setErrors({submit: 'You must be logged in to create or update a listing.'});
            return;
        }
        // 2. Start Loading State & Clear Errors
        setIsLoading(true);
        setErrors({});
        let currentProductId = productId;
        let currentMediaId = mediaId;
        try {
            // 3. Store Step 1 Data Locally
            setAccumulatedData(prev => ({...prev, basicInfo}));
            // 4. Prepare API Payload - MATCHING OLD MONOLITHIC VERSION's TYPES ---
            const productPayload = buildProductPayload({
                ...accumulatedData,
                basicInfo
            }, user);
            // --- End API Payload Prep ---
            // Debug log removed for production
            // 5. Call API: Create Product (Old version only showed addProduct here)
            // if (editMode && currentProductId) {
            // Old version didn't explicitly call update here
            // } else {
            const addResp = await addProduct(productPayload); // Call addProduct
            if (!addResp?.id) throw new Error("No product ID returned from backend");
            currentProductId = addResp.id;
            setProductId(currentProductId); // Update state with the new product ID
            // }
            // 6. Create Media Container (if needed) - Ensure user is still valid
            if (currentProductId && !currentMediaId && user?.userId) {
                const mediaResp = await createMedia({
                    itemId: currentProductId,
                    itemType: "product",
                    userId: user.userId,
                });
                if (!mediaResp?.id) {
                } else {
                    currentMediaId = mediaResp.id;
                    setMediaId(currentMediaId); // Update state with the new media ID
                }
            }
            // 7. Advance Step
            setLastCompletedStep(1);
            handleNextStep();
        } catch (err) {
            setErrors({submit: err.response?.data?.message || "Failed to save basic info. Please try again."});
        } finally {
            setIsLoading(false);
        }
        // Add user, accumulatedData.media dependencies
    }, [user, productId, mediaId, editMode, handleNextStep, isUserLoggedIn, accumulatedData.media]);
    // Media Upload (Step 2 -> 3)
    const handleMediaComplete = useCallback(async (mediaData) => {
        try {
            setIsLoading(true);
            setErrors({});
            if (!productId) {
                throw new Error("Product ID is required for media upload");
            }
            let updatedMediaData = {...mediaData};
            // Media upload is handled by the MediaUploadStep component itself
            // through ImageUploadTab and VideoUploadTab components using addImage/addVideo APIs
            // We just need to save the current state to accumulated data
            // Update product with media information
            const productPayload = buildProductPayload({
                ...accumulatedData,
                media: updatedMediaData
            }, user);
            await updateProduct(productId, {...productPayload, id: productId});
            // Update accumulated data
            setAccumulatedData(prev => ({
                ...prev,
                media: updatedMediaData
            }));
            // Mark step as completed and move to next
            setLastCompletedStep(prev => Math.max(prev, 2));
            handleNextStep();
        } catch (error) {
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to save media'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, productId, handleNextStep]);
    // Optional Settings Step (3 -> 4)
    const handleOptionalComplete = useCallback(async (optionalData) => {
        try {
            setIsLoading(true);
            setErrors({});
            if (!productId) {
                throw new Error("Product ID is required");
            }
            // Update product with optional settings
            const productPayload = buildProductPayload({
                ...accumulatedData,
                optional: optionalData
            }, user);
            await updateProduct(productId, {...productPayload, id: productId});
            // Update accumulated data
            setAccumulatedData(prev => ({
                ...prev,
                optional: optionalData
            }));
            // Mark step as completed and move to next
            setLastCompletedStep(prev => Math.max(prev, 3));
            handleNextStep();
        } catch (error) {
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to save settings'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, productId, handleNextStep]);
    // Final Step - Publish
    const handleFinalizeComplete = useCallback(async (finalData) => {
        try {
            setIsLoading(true);
            setErrors({});
            if (!productId) {
                throw new Error("Product ID is required for publishing");
            }
            // Final update with publish status (location already set in basicInfo)
            const productPayload = buildProductPayload(accumulatedData, user);
            await updateProduct(productId, {
                ...productPayload,
                id: productId,
                published: true
            });
            // Mark as completed and show success
            setLastCompletedStep(4);
            setSuccess(true);
        } catch (error) {
            setErrors({
                submit: error.response?.data?.message || error.message || 'Failed to publish product'
            });
        } finally {
            setIsLoading(false);
        }
    }, [accumulatedData, user, productId]);
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
                        <h1 className={styles.logoText}>Create Product</h1>
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
                    {/* Step 1: Basic Information */}
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
                                isLoading={isLoading}
                                mediaId={mediaId}
                                errors={errors}
                            />
                        </Suspense>
                    )}
                    {/* Step 3: Optional Information */}
                    {currentStep === 3 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <OptionalInfoStep
                                initialData={accumulatedData.optional}
                                onComplete={handleOptionalComplete}
                                onBack={handlePrevStep}
                                isLoading={isLoading}
                                errors={errors}
                                isUserLoggedIn={isUserLoggedIn}
                            />
                        </Suspense>
                    )}
                    {/* Step 4: Finalize */}
                    {currentStep === 4 && (
                        <Suspense fallback={<StepLoadingFallback/>}>
                            <FinalizeStep
                                isSuccess={success}
                                onClose={handleSuccess}
                                onFinalize={handleFinalizeComplete}
                                isLoading={isLoading}
                                productData={accumulatedData}
                                styles={styles}
                            />
                        </Suspense>
                    )}
                </div>
            </div>
        </div>
    );
});
// --- PropTypes ---
CreateProductModal.propTypes = {
    onClose: PropTypes.func.isRequired,
    editMode: PropTypes.bool,
    initialProductData: PropTypes.shape({
        id: PropTypes.string,
        mediaId: PropTypes.string,
        name: PropTypes.string,
        description: PropTypes.string,
        basePrice: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        productPrice: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        categoryId: PropTypes.string,
        categorySlug: PropTypes.string,
        categoryName: PropTypes.string,
        condition: PropTypes.oneOf(["new", "like-new", "excellent", "good", "fair"]),
        brand: PropTypes.string,
        model: PropTypes.string,
        negotiable: PropTypes.bool,
        userType: PropTypes.oneOf(["private", "business", "seller", "manufacturer"]),
        sku: PropTypes.string,
        tags: PropTypes.arrayOf(PropTypes.string),
        images: PropTypes.arrayOf(PropTypes.string),
        videoUrl: PropTypes.string,
        thumbnail: PropTypes.string,
        weight: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        height: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        width: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        depth: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        manageStocks: PropTypes.bool,
        stock: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        shippingCost: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        middlemanService: PropTypes.bool,
        hasVariants: PropTypes.bool,
        attributes: PropTypes.arrayOf(PropTypes.shape({
            key: PropTypes.string,
            value: PropTypes.string
        })),
        lat: PropTypes.number,
        lng: PropTypes.number,
        status: PropTypes.oneOf(["draft", "active"]),
        productType: PropTypes.string,
    }),
};
export default CreateProductModal;

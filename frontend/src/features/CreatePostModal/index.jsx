// src/features/CreatePostModal/index.jsx
import React, {useCallback, useEffect, useMemo, useState, memo} from 'react';
import PropTypes from 'prop-types'; // Import PropTypes
import {useTranslations} from 'next-intl'; //  Import hook
import styles from './CreatePostModal.module.css';
import {useAuth} from "../../context/AuthContext"; // Ensure this provides user object or null
import {addPost, updatePost} from "../../api/client/postsApi"; // Make sure updatePost exists if using editMode
import {createMedia} from "../../api/client/mediaApi";
import {useCategories} from "../../hooks/useCategories"; // Import the categories hook
// Child Components (ensure paths are correct - assume they use translations internally)
import {BasicInfoStep} from './components/steps/BasicInfoStep/BasicInfoStep'; // Assumes translation of labels/placeholders/buttons inside
import {MediaUploadStep} from "../shared/components/MediaUploadStep/MediaUploadStep"; // Assumes translation inside
import {OptionalSettingsStep} from './components/steps/OptionalSettingsStep/OptionalSettingsStep'; // Assumes translation inside
import {SuccessStep} from './components/steps/SuccessStep/SuccessStep'; // Assumes translation inside
// Hooks (ensure paths are correct)
import {useIsMobile} from "../../hooks/useMobileDetection";
import {useFocusTrap} from "../../hooks/useFocusTrap";
import {useAutoSave} from "../../hooks/useAutoSave"; // Assuming this is implemented
// Common Components (ensure path is correct)
import {ErrorAlert} from "../../common/components/ErrorAlert";
const CreatePostModal = memo(function CreatePostModal({onClose, editMode = false, initialData = null}) {
    const t = useTranslations('CreatePostModal'); //  Instantiate hook
    // --- Core State ---
    const [currentStep, setCurrentStep] = useState(1);
    const [lastCompletedStep, setLastCompletedStep] = useState(0);
    const [success, setSuccess] = useState(false);
    // --- Post Data State ---
    const [postId, setPostId] = useState(initialData?.id || null);
    const [postData, setPostData] = useState(() => ({
        name: initialData?.name || "",
        description: initialData?.description || "",
        tags: Array.isArray(initialData?.tags) ? initialData.tags.join(", ") : initialData?.tags || "",
        // New API fields
        typeOfPost: initialData?.typeOfPost || "general",
        userType: initialData?.userType || "private",
        categoryId: initialData?.categoryId || "",
        categorySlug: initialData?.categorySlug || "",
        status: initialData?.status || "active",
        mediaId: initialData?.mediaId || null,
        images: initialData?.images || [],
        videoUrl: initialData?.videoUrl || null,
        thumbnail: initialData?.thumbnail || "", // Include other fields if needed by API
        lat: initialData?.lat || 0,
        lng: initialData?.lng || 0,
    }));
    // --- UI State ---
    const [errors, setErrors] = useState({}); // Stores error keys or general submit/finalize/media keys
    const [isLoading, setIsLoading] = useState(false);
    // --- Hooks ---
    const {user} = useAuth();
    const focusTrapRef = useFocusTrap(true);
    const isMobile = useIsMobile();
    const {lastSaved, isSaving} = useAutoSave(user?.userId, postData, postId);
    // Include the consolidated categories hook with 'posts' type for post categories
    const {
        data: categories = [],
        isLoading: isCategoriesLoading,
        error: categoriesError
    } = useCategories('posts');
    // Log any category fetch errors
    useEffect(() => {
        if (categoriesError) {
            setErrors(prev => ({...prev, categories: "Failed to load categories"}));
        }
    }, [categoriesError]);
    const isUserLoggedIn = useMemo(() => !!user && !!user.userId, [user]);
    // --- Step Navigation Handlers ---
    const handleNextStep = useCallback(() => {
        setCurrentStep(prev => Math.min(prev + 1, 3));
    }, []);
    const handlePrevStep = useCallback(() => {
        setCurrentStep(prev => Math.max(prev - 1, 1));
    }, []);
    const handleStepClick = useCallback((step) => {
        if (step <= Math.min(lastCompletedStep + 1, 3)) {
            setCurrentStep(step);
        }
    }, [lastCompletedStep]);
    // --- Step Completion Handlers ---
    const handleBasicInfoComplete = useCallback(async (formData) => {
        if (!isUserLoggedIn) {
            //   Use translation key
            setErrors({submit: 'errorLoginRequired'});
            return;
        }
        const newErrors = {};
        const trimmedName = formData.name?.trim();
        const trimmedDesc = formData.description?.replace(/<[^>]*>/g, "").trim(); // Check content without HTML tags
        //   Use translation keys for validation
        if (!trimmedName) newErrors.name = "errorTitleRequired";
        else if (trimmedName.length < 5) newErrors.name = t("errorTitleMinLength", {minLength: 5}); // Use t() directly for interpolation message
        if (!trimmedDesc) newErrors.description = "errorContentRequired";
        else if (trimmedDesc.length < 20) newErrors.description = t("errorContentMinLength", {minLength: 20}); // Use t() directly
        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors); // Set specific field errors (keys or translated messages)
            return;
        }
        setIsLoading(true);
        setErrors({});
        let resultingPostId = postId;
        let resultingMediaId = postData.mediaId;
        try {
            // Find categorySlug from categories
            const selectedCategory = categories.find(cat => cat.id === formData.categoryId);
            const categorySlug = selectedCategory ? selectedCategory.slug : "";
            setPostData(prev => ({
                ...prev,
                name: formData.name,
                description: formData.description,
                tags: formData.tags,
                typeOfPost: formData.typeOfPost,
                userType: formData.userType,
                categoryId: formData.categoryId,
                categorySlug: categorySlug,
                status: formData.status
            }));
            const tagsArray = formData.tags.split(",").map((t) => t.trim()).filter(Boolean);
            const postApiData = {
                ...(editMode && postId && {id: postId}),
                name: formData.name,
                description: formData.description,
                typeOfPost: formData.typeOfPost,
                userId: user.userId,
                userType: formData.userType,
                categoryId: formData.categoryId,
                categorySlug: categorySlug,
                tags: tagsArray,
                status: formData.status,
                thumbnail: postData.thumbnail || "",
                lat: formData.lat || 0,
                lng: formData.lng || 0,
            };
            if (editMode && postId) {
                await updatePost(postApiData);
                resultingPostId = postId;
            } else {
                const postResp = await addPost(postApiData);
                if (!postResp?.id) throw new Error("No post ID returned"); // Keep internal error
                resultingPostId = postResp.id;
                setPostId(resultingPostId);
            }
            if (resultingPostId && !resultingMediaId) {
                const mediaResp = await createMedia({
                    itemId: resultingPostId,
                    itemType: "post",
                    userId: user.userId,
                });
                if (mediaResp?.id) {
                    resultingMediaId = mediaResp.id;
                    setPostData(prev => ({...prev, mediaId: resultingMediaId}));
                } else {
                    // Optionally set a non-blocking warning key: setErrors({ media: 'warningMediaCreate' })
                }
            }
            setLastCompletedStep(1);
            handleNextStep();
        } catch (err) {
            //   Use generic translated error key
            setErrors({submit: 'errorSubmitGeneric'});
        } finally {
            setIsLoading(false);
        }
    }, [user, postId, editMode, handleNextStep, postData.mediaId, isUserLoggedIn, postData.thumbnail, categories]);
    const handleMediaComplete = useCallback((mediaData) => {
        setPostData(prev => ({
            ...prev,
            images: mediaData.images ?? prev.images,
            videoUrl: mediaData.videoUrl !== undefined ? mediaData.videoUrl : prev.videoUrl
        }));
        setLastCompletedStep(2);
        handleNextStep();
    }, [handleNextStep]);
    const handleFinalizePost = useCallback(async () => {
        if (!isUserLoggedIn) {
            //   Use translation key
            setErrors({finalize: 'errorLoginRequired'});
            return;
        }
        if (!postId) {
            //   Use translation key
            setErrors({finalize: 'errorMissingPostId'});
            return;
        }
        setIsLoading(true);
        setErrors({});
        try {
            const finalTagsArray = postData.tags.split(",").map(t => t.trim()).filter(Boolean);
            // Prepare media array combining images and videos
            const allMedia = [];
            
            // Add images to media array
            if (postData.images?.length > 0) {
                postData.images.forEach(url => {
                    allMedia.push(url);
                });
            }
            
            // Add video to media array (not as thumbnail!)
            if (postData.videoUrl) {
                allMedia.push(postData.videoUrl);
            }
            
            await updatePost({
                id: postId,
                name: postData.name,
                description: postData.description,
                typeOfPost: postData.typeOfPost,
                userId: user.userId,
                userType: postData.userType,
                categoryId: postData.categoryId,
                categorySlug: postData.categorySlug,
                tags: finalTagsArray,
                status: postData.status,
                // Only use images as thumbnails, never videos
                thumbnail: postData.images?.[0] || "",
                // Store all media (images and videos) in images array
                images: allMedia,
                lat: postData.lat || 0,
                lng: postData.lng || 0,
            });
            setLastCompletedStep(3);
            setSuccess(true);
        } catch (err) {
            //   Use generic translated error key
            setErrors({finalize: 'errorFinalizeGeneric'});
        } finally {
            setIsLoading(false);
        }
    }, [postId, postData, user, isUserLoggedIn, t]); // Added t dependency
    // --- Effect for creating media container (Fallback - Remains the same logic, but use error key) ---
    useEffect(() => {
        if (currentStep === 2 && postId && !postData.mediaId && isUserLoggedIn && !isLoading) {
            let isMounted = true;
            const createMediaForPostFallback = async () => {
                setIsLoading(true);
                setErrors(prev => ({...prev, media: undefined})); // Clear previous media error key
                try {
                    const mediaResp = await createMedia({itemId: postId, itemType: "post", userId: user.userId});
                    if (mediaResp?.id && isMounted) {
                        setPostData(prev => ({...prev, mediaId: mediaResp.id}));
                    } else if (isMounted) {
                        throw new Error("No media ID returned"); // Internal error
                    }
                } catch (err) {
                    if (isMounted) {
                        //   Use translation key for fallback error
                        setErrors(prev => ({...prev, media: "errorMediaCreateFallback"}));
                    }
                } finally {
                    if (isMounted) setIsLoading(false);
                }
            };
            createMediaForPostFallback();
            return () => {
                isMounted = false;
            };
        }
    }, [currentStep, postId, postData.mediaId, isUserLoggedIn, user?.userId, isLoading, t]); // Added t dependency
    // --- Determine Translated Error Message ---
    const errorMessage = useMemo(() => {
        // Prioritize specific field errors if they are keys, otherwise use general errors
        const errorKey = errors.name || errors.description || errors.submit || errors.finalize || errors.media;
        if (errorKey) {
            // Check if errorKey is a direct message (from validation using t()) or a key
            // This simple check assumes keys contain underscores or periods
            if (typeof errorKey === 'string' && (errorKey.includes('_') || errorKey.includes('.'))) {
                return t(errorKey); // Translate if it looks like a key
            }
            return errorKey; // Otherwise, return the message directly (e.g., from t() in validation)
        }
        return null; // No error message to display
    }, [errors, t]);
    // --- Render Logic ---
    return (
        <div className={styles.modalOverlay}>
            <div className={styles.modalContainer} ref={focusTrapRef}>
                <button
                    className={styles.closeButton}
                    onClick={onClose}
                    aria-label="Close modal"
                    type="button"
                >
                    ✕
                </button>
                <div className={styles.sidebar} role="navigation" aria-label="Post creation steps">
                    <div className={styles.logoContainer}>
                        <h1 className={styles.logoText}>Create Post</h1>
                        <div className={`${styles.autosaveIndicator} ${isLoading ? styles.saving : ''}`}>
                            {isLoading ? 'Saving...' : success ? 'Published' : 'Draft'}
                        </div>
                    </div>
                    <nav className={styles.stepsContainer}>
                        {[
                            {step: 1, label: "Basic Info"},
                            {step: 2, label: "Media"},
                            {step: 3, label: "Publish"}
                        ].map(({step, label}) => (
                            <div
                                key={step}
                                className={`${styles.stepNavItem} ${
                                    currentStep === step ? styles.stepNavActive :
                                        step <= lastCompletedStep ? styles.stepNavCompleted :
                                            styles.stepNavDisabled
                                }`}
                                onClick={() => handleStepClick(step)}
                                role="button"
                                tabIndex={step <= Math.min(lastCompletedStep + 1, 3) ? 0 : -1}
                                aria-current={currentStep === step ? "step" : undefined}
                                aria-disabled={step > Math.min(lastCompletedStep + 1, 3)}
                            >
                                <div className={styles.stepNumCircle}>
                                    {step <= lastCompletedStep ? "✓" : step}
                                </div>
                                <span className={styles.stepLabel}>{label}</span>
                            </div>
                        ))}
                    </nav>
                </div>
                <div className={styles.content}>
                    {/* Display translated error message */}
                    {errorMessage && <ErrorAlert message={errorMessage}/>}
                    {/* Render steps based on currentStep */}
                    {/* Steps are assumed to use translations internally */}
                    {currentStep === 1 && (
                        <BasicInfoStep
                            initialData={{
                                name: postData.name, 
                                description: postData.description, 
                                tags: postData.tags,
                                typeOfPost: postData.typeOfPost,
                                userType: postData.userType,
                                categoryId: postData.categoryId,
                                status: postData.status
                            }}
                            onSubmit={handleBasicInfoComplete}
                            onCancel={onClose}
                            isLoading={isLoading}
                            // Pass error keys/messages for specific fields if needed by BasicInfoStep validation display
                            errors={{name: errors.name, description: errors.description}}
                            isUserLoggedIn={isUserLoggedIn}
                            categories={categories}
                        />
                    )}
                    {currentStep === 2 && (
                        <MediaUploadStep
                            mediaId={postData.mediaId}
                            postId={postId}
                            initialImages={postData.images}
                            initialVideoUrl={postData.videoUrl}
                            onComplete={handleMediaComplete}
                            onBack={handlePrevStep}
                            isLoading={isLoading || !postData.mediaId}
                            errors={{media: errors.media}}
                        />
                    )}
                    {currentStep === 3 && !success && (
                        <OptionalSettingsStep
                            onPublish={handleFinalizePost}
                            onBack={handlePrevStep}
                            isLoading={isLoading}
                            // Pass finalize error key if needed
                            errorKey={errors.finalize}
                            isUserLoggedIn={isUserLoggedIn}
                        />
                    )}
                    {success && (
                        <SuccessStep
                            // Pass translated texts/props if SuccessStep needs them
                            // Or SuccessStep uses t() internally
                            onViewDashboard={() => window.location.href = "/dashboard"} // Keep route for now
                            onClose={onClose}
                        />
                    )}
                </div>
            </div>
        </div>
    );
});
// PropTypes for the main component (remains the same)
CreatePostModal.propTypes = {
    onClose: PropTypes.func.isRequired,
    editMode: PropTypes.bool,
    initialData: PropTypes.shape({
        id: PropTypes.string,
        name: PropTypes.string,
        description: PropTypes.string,
        tags: PropTypes.oneOfType([PropTypes.string, PropTypes.array]),
        // New API fields
        typeOfPost: PropTypes.string,
        userType: PropTypes.string,
        categoryId: PropTypes.string,
        categorySlug: PropTypes.string,
        status: PropTypes.string,
        mediaId: PropTypes.string,
        images: PropTypes.array,
        videoUrl: PropTypes.string,
        thumbnail: PropTypes.string,
        lat: PropTypes.number,
        lng: PropTypes.number,
    }),
};
// Default props
CreatePostModal.defaultProps = {
    editMode: false,
    initialData: null,
};
export default CreatePostModal;

/*****************************************************************************
 * FILE: CreatePostModal.jsx (Refactored to use CSS Module)
 *****************************************************************************/
"use client";
import React, {useState, useEffect, useRef, memo} from "react";
import styles from "./PostModal.module.css"; // Import the CSS module
import {useAuth} from "../../context/AuthContext";
import {addPost} from "../../api/client/postsApi";
import {createMedia} from "../../api/client/mediaApi";
import {useCategory} from "../../context/CategoryContext";
import CategoryTree from "../../components/Category/CategoryTree";
import ImageUploaderHookForm from "../Uploader/ImageUploader";
import VideoUploaderHookForm from "../Uploader/VideoUploader";
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import RichTextEditor from "./TextEditor"; // Import Quill styles
const CreatePostModal = memo(({onClose}) => {
    const {user} = useAuth();
    const modalRef = useRef(null);
    const [currentStep, setCurrentStep] = useState(1);
    const [createdPostId, setCreatedPostId] = useState(null);
    const [createdMediaId, setCreatedMediaId] = useState(null);
    // Step 1: Basic Info
    const [postTitle, setPostName] = useState("");
    const [postContent, setPostDescription] = useState("");
    const [tags, setTags] = useState("");
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    // Step 2: Media
    const [uploadedImages, setUploadedImages] = useState([]);
    const [uploadedVideoUrl, setUploadedVideoUrl] = useState(null);
    // Optionally close on outside click:
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (modalRef.current && !modalRef.current.contains(event.target)) {
                // If you want to close on outside click => handleCancel()
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);
    // ---------------------
    // Step 1: Create Post & Media
    // ---------------------
    const handleSaveBaseInfo = async (e) => {
        e.preventDefault();
        setError("");
        if (!postTitle.trim() || !postContent.trim()) {
            setError("Please fill in required fields: Name, Content, Category.");
            return;
        }
        // Build post data
        const postData = {
            title: postTitle,
            content: postContent,
            userId: user?.userId,
            status: "active",
        };
        try {
            setIsLoading(true);
            // 1) create post
            const postResp = await addPost(postData);
            const newPostId = postResp?.id;
            if (!newPostId) throw new Error("No post ID returned from backend");
            // 2) create media record
            const mediaResp = await createMedia({
                itemId: newPostId,
                itemType: "post",
                userId: user?.userId,
            });
            const newMediaId = mediaResp?.id;
            if (!newMediaId) throw new Error("No media ID returned from backend");
            setCreatedPostId(newPostId);
            setCreatedMediaId(newMediaId);
            setCurrentStep(2);
        } catch (err) {
            setError("Error. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };
    // ---------------------
    // Step 2: Media Upload
    // ---------------------
    const handleMediaDone = () => {
        setCurrentStep(3);
    };
    // ---------------------
    // Step 3: Optional => finalize => Step4
    // ---------------------
    const handleFinalizePost = () => {
        // Possibly patch post => set status='active'
        setCurrentStep(4);
    };
    // ---------------------
    // Step 4: Done => close
    // ---------------------
    const handleCloseAll = () => {
        onClose();
    };
    // Cancel
    const handleCancel = () => {
        onClose();
    };
    return (
        <div className={styles.overlay}>
            <div className={styles.modalWrapper} ref={modalRef}>
                <div className={styles.modalHeader}>
                    <h2 className={styles.modalTitle}>Create / Edit Listing</h2>
                    <button
                        className={styles.closeButton}
                        onClick={handleCancel}
                        aria-label="Close Modal"
                    >
                        <img src="/images/close-icon.svg" alt="Close"/>
                    </button>
                </div>
                <div className={styles.stepsContainer}>
                    {/* Step 1 */}
                    <div
                        className={
                            currentStep === 1 ? styles.singleStepActive : styles.singleStep
                        }
                    >
                        1
                    </div>
                    <div className={styles.stepLine}/>
                    {/* Step 2 */}
                    <div
                        className={
                            currentStep === 2 ? styles.singleStepActive : styles.singleStep
                        }
                    >
                        2
                    </div>
                    <div className={styles.stepLine}/>
                    {/* Step 3 */}
                    <div
                        className={
                            currentStep === 3 ? styles.singleStepActive : styles.singleStep
                        }
                    >
                        3
                    </div>
                    <div className={styles.stepLine}/>
                    {/* Step 4 */}
                    <div
                        className={
                            currentStep === 4 ? styles.singleStepActive : styles.singleStep
                        }
                    >
                        4
                    </div>
                </div>
                <div className={styles.contentArea}>
                    {error && <div className={styles.errorBox}>{error}</div>}
                    {currentStep === 1 && (
                        <Step1BaseInfo
                            postTitle={postTitle}
                            setPostName={setPostName}
                            postContent={postContent}
                            setPostDescription={setPostDescription}
                            tags={tags}
                            setTags={setTags}
                            onSaveBaseInfo={handleSaveBaseInfo}
                            isLoading={isLoading}
                        />
                    )}
                    {currentStep === 2 && (
                        <Step2MediaUpload
                            mediaId={createdMediaId}
                            uploadedImages={uploadedImages}
                            setUploadedImages={setUploadedImages}
                            uploadedVideoUrl={uploadedVideoUrl}
                            setUploadedVideoUrl={setUploadedVideoUrl}
                            onDone={handleMediaDone}
                        />
                    )}
                    {currentStep === 3 && (
                        <Step3OptionalInfo onFinish={handleFinalizePost}/>
                    )}
                    {currentStep === 4 && <Step4Done onCloseAll={handleCloseAll}/>}
                </div>
            </div>
        </div>
    );
});
export default CreatePostModal;
function Step1BaseInfo({
                           postTitle,
                           setPostName,
                           postContent,
                           setPostDescription,
                           tags,
                           setTags,
                           onSaveBaseInfo,
                           isLoading,
                       }) {
    const { categories, loading, error } = useCategory();
    // Handle tag change
    const handleTagChange = (e) => {
        setTags(e.target.value);
    };
    return (
        <div className={styles.stepBox}>
            <h3 className={styles.stepHeading}>Step 1: Basic Information</h3>
            <p className={styles.instruction}>Provide essential details about your item.</p>
            {error && (
                <div className={styles.errorBox}>
                    Failed to load categories: {String(error)}
                </div>
            )}
            {loading && <p>Loading categories...</p>}
            <form className={styles.styledForm} onSubmit={onSaveBaseInfo}>
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.label}>
                            Title <span>*</span>
                        </label>
                        <input
                            className={styles.textInput}
                            value={postTitle}
                            onChange={(e) => setPostName(e.target.value)}
                            required
                            placeholder="e.g. iPhone 12 in great condition"
                        />
                    </div>
                </div>
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.label}>
                            Content <span>*</span>
                        </label>
                        <div className={styles.editorContainer}>
                            {/* Import the new component instead of directly using ReactQuill */}
                            <RichTextEditor
                                value={postContent}
                                onChange={setPostDescription}
                                placeholder="Describe the condition, highlights, etc."
                            />
                        </div>
                    </div>
                </div>
                <div className={styles.formRow}>
                    <div className={styles.formGroup}>
                        <label className={styles.label}>Tags (comma-separated)</label>
                        <input
                            className={styles.textInput}
                            placeholder="e.g. electronics, smartphone"
                            value={tags}
                            onChange={handleTagChange}
                        />
                    </div>
                </div>
                <div className={styles.buttonContainer}>
                    <button className={styles.actionButton} type="submit" disabled={isLoading}>
                        {isLoading ? "Saving..." : "Save & Continue"}
                    </button>
                </div>
            </form>
        </div>
    );
}
/* ------------------------------------------------------------------
   STEP 2: Media Upload
------------------------------------------------------------------ */
function Step2MediaUpload({
                              mediaId,
                              uploadedImages,
                              setUploadedImages,
                              uploadedVideoUrl,
                              setUploadedVideoUrl,
                              onDone,
                          }) {
    const [activeTab, setActiveTab] = useState("images");
    return (
        <div className={styles.stepBox}>
            <h3 className={styles.stepHeading}>Step 2: Upload Photos / Videos</h3>
            <p className={styles.instruction}>
                Attach photos and/or videos for your listing.
            </p>
            <div className={styles.toggleTabs}>
                <button
                    className={
                        activeTab === "images"
                            ? `${styles.tabButton} ${styles.tabButtonActive}`
                            : styles.tabButton
                    }
                    onClick={() => setActiveTab("images")}
                >
                    Images
                </button>
                <button
                    className={
                        activeTab === "videos"
                            ? `${styles.tabButton} ${styles.tabButtonActive}`
                            : styles.tabButton
                    }
                    onClick={() => setActiveTab("videos")}
                >
                    Videos
                </button>
            </div>
            {activeTab === "images" && (
                <>
                    <ImageUploaderHookForm
                        mediaId={mediaId}
                        onUploadSuccess={(viewUrl) => {
                            setUploadedImages((prev) => [...prev, viewUrl]);
                        }}
                    />
                    {uploadedImages.length > 0 && (
                        <div className={styles.imageGrid}>
                            {uploadedImages.map((viewUrl, idx) => (
                                <div className={styles.imagePreview} key={idx}>
                                    <img src={viewUrl} alt={`uploaded-${idx}`}/>
                                </div>
                            ))}
                        </div>
                    )}
                </>
            )}
            {activeTab === "videos" && (
                <>
                    <VideoUploaderHookForm
                        mediaId={mediaId}
                        onUploadSuccess={(viewUrl) => {
                            setUploadedVideoUrl(viewUrl);
                        }}
                    />
                    {uploadedVideoUrl && (
                        <div className={styles.videoPreview}>
                            <video src={uploadedVideoUrl} controls/>
                        </div>
                    )}
                </>
            )}
            <div className={styles.buttonContainer}>
                <button className={styles.actionButton} type="button" onClick={onDone}>
                    Proceed to Next Step
                </button>
            </div>
        </div>
    );
}
/* ------------------------------------------------------------------
   STEP 3: Optional Info
------------------------------------------------------------------ */
function Step3OptionalInfo({onFinish}) {
    return (
        <div className={styles.stepBox}>
            <h3 className={styles.stepHeading}>Step 3: Optional Details</h3>
            <p className={styles.instruction}>
                Add shipping methods, variants (like color, size), or any extra info.
            </p>
            <div className={styles.placeholder}>
                <p>[Additional configuration could go here]</p>
            </div>
            <div className={styles.buttonContainer}>
                <button className={styles.actionButton} type="button" onClick={onFinish}>
                    Publish / Finish
                </button>
            </div>
        </div>
    );
}
/* ------------------------------------------------------------------
   STEP 4: Done
------------------------------------------------------------------ */
function Step4Done({onCloseAll}) {
    return (
        <div className={styles.stepBox}>
            <h3 className={styles.stepHeading}>Step 4: Done</h3>
            <p className={styles.instruction}>Your draft is created successfully!</p>
            <div className={styles.placeholder}>
                <p>
                    You can finalize or edit more details. If ready, set status="active" on your
                    backend, or just close.
                </p>
            </div>
            <div className={styles.buttonContainer}>
                <button className={styles.actionButton} type="button" onClick={onCloseAll}>
                    Close
                </button>
            </div>
        </div>
    );
}

/*****************************************************************************
 * FILE: VideoModal.jsx
 * PURPOSE:
 *   1) Step 1 => create a Media record (with a newly generated itemId).
 *   2) Step 2 => upload a video for that media using VideoUploaderHookForm.
 *****************************************************************************/
"use client";
import React, {useState, useEffect, useRef, memo} from "react";
import {v4 as uuidv4} from "uuid";               // <-- For generating itemId
import {useAuth} from "../../context/AuthContext";
import {createMedia} from "../../api/client/mediaApi";
import VideoUploaderHookForm from "../Uploader/VideoUploader";
/** Style constants */
const ACCENT_COLOR = "#7dcfb6";
const TEXT_COLOR = "#2d3748";
const BORDER_COLOR = "#cbd5e0";
const BACKGROUND_LEFT = "#f6f6f6";
/** The main modal component */
const VideoModal = memo(function VideoModal({onClose}) {
    const {user} = useAuth();
    const modalRef = useRef(null);
    // Step management
    const [currentStep, setCurrentStep] = useState(1);
    // The newly created media ID (once step 1 is complete)
    const [createdMediaId, setCreatedMediaId] = useState(null);
    // If you want to show the uploaded video in a <video> preview
    const [uploadedVideoUrl, setUploadedVideoUrl] = useState(null);
    // If you are referencing images but not using them, at least define them to avoid errors
    const [uploadedImages, setUploadedImages] = useState([]);
    // UI states
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    // (Optional) Close on outside click
    useEffect(() => {
        function handleClickOutside(e) {
            if (modalRef.current && !modalRef.current.contains(e.target)) {
                // onClose();
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, [onClose]);
    /**
     * STEP 1 => Create Media
     * Generate a random itemId, call createMedia(...), store the mediaId.
     */
    async function handleCreateMedia() {
        setError("");
        setIsLoading(true);
        try {
            // 1. Generate a random itemId
            const newItemId = uuidv4();
            // 2. Create media on the server
            const mediaResp = await createMedia({
                itemId: newItemId,
                itemType: "video",
                userId: user?.userId,
            });
            if (!mediaResp?.id) {
                throw new Error("No media ID returned from createMedia");
            }
            setCreatedMediaId(mediaResp.id);
            setCurrentStep(2);
        } catch (err) {
            setError("Something went wrong. Please try again.");
        } finally {
            setIsLoading(false);
        }
    }
    function handleMediaDone() {
        setCurrentStep(3);
    }
    function handleFinalizeProduct() {
        setCurrentStep(4);
    }
    function handleCloseAll() {
        onClose();
    }
    function handleCancel() {
        onClose();
    }
    return (
        <div style={css(styles.overlay)}>
            <div style={css(styles.modalContainer)} ref={modalRef}>
                {/* LEFT NAV */}
                <div style={css(styles.leftNav)}>
                    <div style={css(styles.logoArea)}>
                        <h2 style={css(styles.logoText)}>Steps</h2>
                    </div>
                    {/* Step 1 Nav Item */}
                    <StepNavItem
                        stepNum={1}
                        label="Create Media"
                        active={currentStep === 1}
                        onClick={() => setCurrentStep(1)}
                    />
                    {/* Step 2 Nav Item */}
                    <StepNavItem
                        stepNum={2}
                        label="Upload Video"
                        active={currentStep === 2}
                        onClick={() => {
                            if (createdMediaId) {
                                setCurrentStep(2);
                            }
                        }}
                    />
                    <StepNavItem
                        stepNum={3}
                        label="Optional"
                        active={currentStep === 3}
                        onClick={() => {
                            if (createdMediaId) {
                                setCurrentStep(3);
                            }
                        }}
                    />
                    <StepNavItem
                        stepNum={4}
                        label="Done"
                        active={currentStep === 4}
                        onClick={() => setCurrentStep(4)}
                    />
                </div>
                {/* RIGHT CONTENT */}
                <div style={css(styles.rightContent)}>
                    {/* Header */}
                    <div style={css(styles.headerBar)}>
                        <h3 style={css(styles.headerTitle)}>Create Video</h3>
                        <button style={css(styles.closeX)} onClick={handleCancel} aria-label="Close Modal">
                            ✕
                        </button>
                    </div>
                    <div style={css(styles.innerContent)}>
                        {error && <div style={css(styles.errorBox)}>{error}</div>}
                        {/* Step 1 => Button to create Media */}
                        {currentStep === 1 && (
                            <div style={css(styles.stepBox)}>
                                <h2 style={css(styles.stepHeading)}>Step 1: Create Media</h2>
                                <p style={css(styles.instruction)}>
                                    Generate a unique itemId and create a media record on the server.
                                </p>
                                <button
                                    style={css(styles.primaryBtn)}
                                    onClick={handleCreateMedia}
                                    disabled={isLoading}
                                >
                                    {isLoading ? "Creating..." : "Create Media"}
                                </button>
                            </div>
                        )}
                        {/* Step 2 => Upload Video */}
                        {currentStep === 2 && (
                            <Step2MediaUpload
                                mediaId={createdMediaId}
                                uploadedVideoUrl={uploadedVideoUrl}
                                setUploadedVideoUrl={setUploadedVideoUrl}
                                onDone={handleMediaDone}
                            />
                        )}
                        {/* Step 3 => Optional Info */}
                        {currentStep === 3 && <Step3OptionalInfo onFinish={handleFinalizeProduct}/>}
                        {/* Step 4 => Done */}
                        {currentStep === 4 && <Step4Done onCloseAll={handleCloseAll}/>}
                    </div>
                </div>
            </div>
        </div>
    );
});
export default VideoModal;
/** Step nav item with circle + label */
function StepNavItem({stepNum, label, active, onClick}) {
    return (
        <div
            style={{
                ...styles.navItem,
                background: active ? ACCENT_COLOR : "transparent",
                color: active ? "#fff" : TEXT_COLOR,
            }}
            onClick={onClick}
        >
            <div
                style={{
                    ...styles.circleIcon,
                    background: active ? "#fff" : "#e2e8f0",
                    color: active ? ACCENT_COLOR : TEXT_COLOR,
                }}
            >
                {stepNum}
            </div>
            <span style={styles.navLabel}>{label}</span>
        </div>
    );
}
/* ---------------------------------------
   STEP 2: Media (Upload Video)
--------------------------------------- */
function Step2MediaUpload({
                              mediaId,
                              uploadedVideoUrl,
                              setUploadedVideoUrl,
                              onDone,
                          }) {
    // If you have a tab for images or other media, you can do so,
    // but here we show just the video upload part.
    const canProceed = true; // or implement your own logic
    return (
        <div style={css(styles.stepBox)}>
            <h2 style={css(styles.stepHeading)}>Step 2: Upload Video</h2>
            <p style={css(styles.instruction)}>Show off your item with a great video.</p>
            <VideoUploaderHookForm
                mediaId={mediaId}
                onUploadSuccess={(vidUrl) => setUploadedVideoUrl(vidUrl)}
            />
            {uploadedVideoUrl && (
                <div style={css(styles.videoPreview)}>
                    <video src={uploadedVideoUrl} controls style={{width: "100%"}}/>
                </div>
            )}
            <div style={css(styles.btnRow)}>
                <button
                    style={css(styles.primaryBtn)}
                    type="button"
                    onClick={() => canProceed && onDone()}
                >
                    Proceed to Next Step
                </button>
            </div>
        </div>
    );
}
/* ---------------------------------------
   STEP 3: Optional Info
--------------------------------------- */
function Step3OptionalInfo({onFinish}) {
    return (
        <div style={css(styles.stepBox)}>
            <h2 style={css(styles.stepHeading)}>Step 3: Optional Info</h2>
            <p style={css(styles.instruction)}>
                Add shipping methods, item variants, or any other optional details.
            </p>
            <div style={css(styles.placeholderBox)}>
                <p>[Potential extra fields could go here]</p>
            </div>
            <div style={css(styles.btnRow)}>
                <button style={css(styles.primaryBtn)} type="button" onClick={onFinish}>
                    Publish / Finish
                </button>
            </div>
        </div>
    );
}
/* ---------------------------------------
   STEP 4: Done
--------------------------------------- */
function Step4Done({onCloseAll}) {
    return (
        <div style={css(styles.stepBox)}>
            <h2 style={css(styles.stepHeading)}>Step 4: Done!</h2>
            <p style={css(styles.instruction)}>
                Your media record and video have been successfully created.
            </p>
            <div style={css(styles.placeholderBox)}>
                <p>Feel free to edit or share now.</p>
            </div>
            <div style={css(styles.btnRow)}>
                <button style={css(styles.primaryBtn)} type="button" onClick={onCloseAll}>
                    Close
                </button>
            </div>
        </div>
    );
}
/** Helper to apply inline styles easily */
function css(obj) {
    return obj;
}
/** Styles object */
const styles = {
    overlay: {
        position: "fixed",
        top: 0, left: 0,
        zIndex: 9999,
        width: "100%",
        height: "100%",
        background: "rgba(0,0,0,0.5)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "20px",
        fontFamily: '"Helvetica Neue", Arial, sans-serif',
    },
    modalContainer: {
        backgroundColor: "#fff",
        display: "flex",
        flexDirection: "row",
        width: "100%",
        maxWidth: "1080px",
        borderRadius: "10px",
        boxShadow: "0 8px 20px rgba(0,0,0,0.25)",
        maxHeight: "90vh",
        overflow: "hidden",
    },
    leftNav: {
        width: "220px",
        background: BACKGROUND_LEFT,
        borderRight: `1px solid ${BORDER_COLOR}`,
        display: "flex",
        flexDirection: "column",
        padding: "12px",
        gap: "12px",
    },
    logoArea: {
        textAlign: "center",
        marginBottom: "12px",
    },
    logoText: {
        margin: 0,
        fontSize: "18px",
        color: TEXT_COLOR,
    },
    navItem: {
        padding: "10px 8px",
        borderRadius: "6px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        gap: "8px",
        cursor: "pointer",
        transition: "background 0.2s",
    },
    circleIcon: {
        width: "32px",
        height: "32px",
        borderRadius: "50%",
        background: "#e2e8f0",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        fontWeight: 600,
    },
    navLabel: {
        fontSize: "14px",
        fontWeight: 600,
    },
    rightContent: {
        flex: 1,
        display: "flex",
        flexDirection: "column",
    },
    headerBar: {
        borderBottom: `1px solid ${BORDER_COLOR}`,
        padding: "12px 16px",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
    },
    headerTitle: {
        margin: 0,
        fontSize: "20px",
        fontWeight: 600,
        color: TEXT_COLOR,
    },
    closeX: {
        background: "transparent",
        border: "none",
        fontSize: "18px",
        cursor: "pointer",
        color: "#666",
    },
    innerContent: {
        padding: "20px",
        overflowY: "auto",
        flex: 1,
        display: "flex",
        flexDirection: "column",
    },
    stepBox: {
        background: "#fff",
        border: "1px solid #eee",
        borderRadius: "6px",
        padding: "20px",
        marginBottom: "24px",
        flex: 1,
    },
    stepHeading: {
        margin: 0,
        marginBottom: "8px",
        fontSize: "18px",
        fontWeight: 600,
        color: TEXT_COLOR,
    },
    instruction: {
        marginBottom: "16px",
        fontSize: "14px",
        color: "#555",
        lineHeight: 1.4,
    },
    btnRow: {
        display: "flex",
        justifyContent: "flex-end",
        marginTop: "12px",
    },
    primaryBtn: {
        background: ACCENT_COLOR,
        color: "#fff",
        fontSize: "15px",
        fontWeight: 600,
        border: "none",
        borderRadius: "24px",
        padding: "10px 20px",
        cursor: "pointer",
    },
    videoPreview: {
        marginTop: "16px",
        width: "320px",
        border: "1px solid #ddd",
        borderRadius: "6px",
        overflow: "hidden",
    },
    errorBox: {
        background: "#ffe6e6",
        color: "#b30000",
        padding: "12px 16px",
        borderRadius: "6px",
        border: "1px solid #f5c6c6",
        marginBottom: "16px",
        fontWeight: 500,
    },
    placeholderBox: {
        padding: "16px",
        background: "#fafafa",
        border: "2px dashed #ccc",
        borderRadius: "6px",
        textAlign: "center",
        color: "#888",
        flex: 1,
    },
};

// File: MediaGalleryManagerEditable.jsx
import React, { useState, useCallback, memo } from 'react';
import styled from 'styled-components';
import MediaGalleryMobileEditable from './MediaGalleryMobileEditable';
import ImageUploaderHookForm from "../../features/Uploader/ImageUploader";
import VideoUploaderHookForm from "../../features/Uploader/VideoUploader";
/**
 * OPTIMIZED: Memoized for better media gallery performance
 */
const MediaGalleryManagerEditable = memo(function MediaGalleryManagerEditable({
                                                        initialGallery = [],
                                                        mediaId,
                                                        onGalleryUpdate,
                                                    }) {
    const [gallery, setGallery] = useState(initialGallery);
    const [currentMediaIndex, setCurrentMediaIndex] = useState(0);
    const [activeVideoIndex, setActiveVideoIndex] = useState(null);
    const [isEditMode, setIsEditMode] = useState(false);
    const handleNewMedia = useCallback((viewUrl, type) => {
        const newItem = {
            type,         // 'image' or 'video'
            src: viewUrl, // The final S3 location or CloudFront URL
            alt: '',
            // If video, optionally set a default poster
            poster: (type === 'video') ? '/images/video-default-poster.png' : undefined,
        };
        setGallery((prev) => {
            const updated = [...prev, newItem];
            // If parent wants to track changes:
            onGalleryUpdate?.(updated);
            return updated;
        });
    }, [onGalleryUpdate]);
    // 5) Delete item
    const handleDeleteItem = useCallback((index) => {
        setGallery((prev) => {
            const copy = [...prev];
            const removedItem = copy.splice(index, 1)[0];
            // Optionally: call your backend to remove or "mark as deleted":
            // await axios.delete(`/api/media?url=${removedItem.src}`) ...
            // or remove from S3, etc.
            // If parent wants changes:
            onGalleryUpdate?.(copy);
            return copy;
        });
        if (currentMediaIndex >= index && currentMediaIndex > 0) {
            setCurrentMediaIndex((prevIdx) => prevIdx - 1);
        }
    }, [currentMediaIndex, onGalleryUpdate]);
    // 6) Next / prev
    const handleNext = useCallback(() => {
        setCurrentMediaIndex((prev) => (prev + 1) % gallery.length);
        setActiveVideoIndex(null);
    }, [gallery.length]);
    const handlePrev = useCallback(() => {
        setCurrentMediaIndex((prev) => (prev - 1 + gallery.length) % gallery.length);
        setActiveVideoIndex(null);
    }, [gallery.length]);
    return (
        <Container>
            <HeaderBar>
                <Title>Media Gallery (Editable)</Title>
                <EditToggleButton onClick={() => setIsEditMode((v) => !v)}>
                    {isEditMode ? 'Close Edit Mode' : 'Edit Gallery'}
                </EditToggleButton>
            </HeaderBar>
            {/*
        MAIN DISPLAY using MediaGalleryMobileEditable
        so user can see a "read-only" style gallery (swipe left/right).
      */}
            <MediaGalleryMobileEditable
                gallery={gallery}
                currentMediaIndex={currentMediaIndex}
                setCurrentMediaIndex={setCurrentMediaIndex}
                handleNext={handleNext}
                handlePrev={handlePrev}
                productTitle="My Product"
                activeVideoIndex={activeVideoIndex}
                setActiveVideoIndex={setActiveVideoIndex}
            />
            <ImageUploaderHookForm
                mediaId={mediaId}
                onUploadSuccess={(viewUrl) => handleNewMedia(viewUrl, 'image')}
            />
            <UploaderRow>
                {/* Single-video upload */}
                <UploaderBox>
                    <SubTitle>Single Video</SubTitle>
                    <VideoUploaderHookForm
                        mediaId={mediaId}
                        onUploadSuccess={(viewUrl) => handleNewMedia(viewUrl, 'video')}
                    />
                </UploaderBox>
            </UploaderRow>
        </Container>
    );
}, (prevProps, nextProps) => {
    // Only re-render if essential props changed
    return (
        prevProps.mediaId === nextProps.mediaId &&
        JSON.stringify(prevProps.initialGallery) === JSON.stringify(nextProps.initialGallery)
        // Skip onGalleryUpdate function comparison
    );
});
/* ~~~~~~~~~~~~~~~~~ STYLED COMPONENTS ~~~~~~~~~~~~~~~~~ */
const Container = styled.div`
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 12px;
`;
const HeaderBar = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
`;
const Title = styled.h2`
    margin: 0;
`;
const EditToggleButton = styled.button`
    background: #007bc7;
    color: #fff;
    font-size: 0.9rem;
    border: none;
    padding: 6px 12px;
    border-radius: 4px;
    cursor: pointer;
    &:hover {
        background: #005e99;
    }
`;
const EditPanel = styled.div`
    border: 1px solid #ccc;
    padding: 12px;
    border-radius: 6px;
    background: #f7f7f7;
    display: flex;
    flex-direction: column;
    gap: 14px;
`;
const SectionHeader = styled.h3`
    margin: 0;
    font-size: 1rem;
    color: #333;
`;
const UploaderRow = styled.div`
    display: flex;
    gap: 14px;
    flex-wrap: wrap;
`;
const UploaderBox = styled.div`
    flex: 1 1 240px;
    min-width: 240px;
    border: 1px dashed #aaa;
    padding: 10px;
    border-radius: 6px;
    background: #fff;
`;
const SubTitle = styled.h4`
    margin-top: 0;
    font-size: 0.9rem;
    color: #444;
`;
const NoItemsLabel = styled.div`
    font-size: 0.85rem;
    color: #999;
    font-style: italic;
`;
const GalleryItemRow = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 4px 0;
`;
const ItemInfo = styled.div`
    font-size: 0.85rem;
    color: #333;
`;
const DeleteButton = styled.button`
    background: #e74c3c;
    color: #fff;
    border: none;
    padding: 4px 8px;
    font-size: 0.8rem;
    border-radius: 4px;
    cursor: pointer;
    &:hover {
        background: #c0392b;
    }
`;
export default MediaGalleryManagerEditable;

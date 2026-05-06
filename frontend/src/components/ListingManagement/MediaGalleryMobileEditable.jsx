import { Plus, Trash2, Upload } from '@/icons';
import { FaGripVertical, FaPlay } from '../../utils/iconImports';
import React, {
    useState,
    useCallback,
    useMemo,
    Suspense,
    memo,
    lazy,
} from 'react';
import PropTypes from 'prop-types';
import styled, { css } from 'styled-components';
import { useSwipeable } from 'react-swipeable';
import Spinner from '../Utils/Spinner';
import MediaDisplayEditable from './MediaDisplayEditable';
const LightboxViewerEditable = React.lazy(() => import('./LightboxViewerEditable').then(module => ({ default: module.default })));
const LazyImageUploader = lazy(() =>
    import('../../features/Uploader/ImageUploader').then(module => ({ default: module.default }))
);
const LazyVideoUploader = lazy(() =>
    import('../../features/Uploader/VideoUploader').then(module => ({ default: module.default }))
);
function wrapIndex(index, length) {
    return ((index % length) + length) % length;
}
function MediaGalleryMobileEditable({
                                        gallery,
                                        setGallery,
                                        currentMediaIndex,
                                        setCurrentMediaIndex,
                                        handleNext,
                                        handlePrev,
                                        productTitle,
                                        activeVideoIndex,
                                        setActiveVideoIndex,
                                        mediaId
                                    }) {
    const [isLightboxOpen, setIsLightboxOpen] = useState(false);
    const [dragOverIndex, setDragOverIndex] = useState(null);
    // 1) Filter out invalid items
    const validGallery = useMemo(() => {
        if (!Array.isArray(gallery)) return [];
        return gallery.filter((item) => item && item.src);
    }, [gallery]);
    // 2) Current item
    const currentMedia = useMemo(() => {
        if (validGallery.length === 0) return null;
        const safeIndex = wrapIndex(currentMediaIndex, validGallery.length);
        return validGallery[safeIndex];
    }, [currentMediaIndex, validGallery]);
    const isSingleMedia = validGallery.length <= 1;
    const handleThumbnailClick = useCallback(
        (index) => {
            setCurrentMediaIndex(index);
            setActiveVideoIndex(null);
        },
        [setCurrentMediaIndex, setActiveVideoIndex]
    );
    // 5) Lightbox toggles
    const openLightbox = useCallback(() => {
        setIsLightboxOpen(true);
        setActiveVideoIndex(null);
    }, [setActiveVideoIndex]);
    const closeLightbox = useCallback(() => {
        setIsLightboxOpen(false);
    }, []);
    // 6) Swipe handlers
    const swipeLeft = useCallback(() => {
        handleNext();
        setActiveVideoIndex(null);
    }, [handleNext, setActiveVideoIndex]);
    const swipeRight = useCallback(() => {
        handlePrev();
        setActiveVideoIndex(null);
    }, [handlePrev, setActiveVideoIndex]);
    const swipeHandlers = useSwipeable({
        onSwipedLeft: isSingleMedia ? () => {} : swipeLeft,
        onSwipedRight: isSingleMedia ? () => {} : swipeRight,
        trackMouse: !isSingleMedia,
        preventDefaultTouchmoveEvent: !isSingleMedia,
    });
    // 7) Main image => open Lightbox if image
    const handleMainImageClick = useCallback(() => {
        if (currentMedia?.type === 'image') {
            openLightbox();
        }
    }, [currentMedia?.type, openLightbox]);
    // 8) Remove
    const handleRemoveItem = useCallback(
        (index) => {
            setGallery((prev) => {
                const copy = [...prev];
                copy.splice(index, 1);
                return copy;
            });
            if (currentMediaIndex === index && currentMediaIndex > 0) {
                setCurrentMediaIndex((prev) => prev - 1);
            }
        },
        [currentMediaIndex, setGallery, setCurrentMediaIndex]
    );
    // 9) Drag & drop reorder
    const handleDragStart = (e, dragIndex) => {
        e.dataTransfer.setData('application/index', String(dragIndex));
        e.dataTransfer.effectAllowed = 'move';
    };
    const handleDragOver = (e, overIndex) => {
        e.preventDefault();
        setDragOverIndex(overIndex);
    };
    const handleDrop = (e, dropIndex) => {
        e.preventDefault();
        const dragIndex = Number(e.dataTransfer.getData('application/index'));
        if (dragIndex === dropIndex) {
            setDragOverIndex(null);
            return;
        }
        setGallery((prev) => {
            const copy = [...prev];
            const [dragItem] = copy.splice(dragIndex, 1);
            copy.splice(dropIndex, 0, dragItem);
            return copy;
        });
        setDragOverIndex(null);
    };
    const handleDragLeave = () => setDragOverIndex(null);
    // 10) Show/hide modals for uploading
    const [showImageModal, setShowImageModal] = useState(false);
    const [showVideoModal, setShowVideoModal] = useState(false);
    const openImageModal = () => setShowImageModal(true);
    const closeImageModal = () => setShowImageModal(false);
    const openVideoModal = () => setShowVideoModal(true);
    const closeVideoModal = () => setShowVideoModal(false);
    const handleImageSuccess = (imageUrl) => {
        // Add new image to gallery
        setGallery((prev) => [...prev, { type: 'image', src: imageUrl }]);
        closeImageModal();
    };
    const handleVideoSuccess = (videoUrl) => {
        // Add new video to gallery
        setGallery((prev) => [
            ...prev,
            { type: 'video', src: videoUrl, poster: '/images/video-icon.webp' },
        ]);
        closeVideoModal();
    };
    return (
        <GalleryContainer>
            {/* Buttons that open modals to add images/videos */}
            <MainLayout>
                <ThumbnailsColumn>
                    {validGallery.map((item, index) => {
                        const isActive = index === currentMediaIndex;
                        const isVideo = item.type === 'video';
                        const thumbSrc = isVideo
                            ? item.poster || '/images/video-icon.webp'
                            : item.src;
                        const altText = item.alt
                            ? item.alt
                            : isVideo
                                ? `Video ${index + 1}`
                                : `Image ${index + 1}`;
                        const isDragOver = dragOverIndex === index;
                        return (
                            <ThumbButton
                                key={`thumb-${index}`}
                                type="button"
                                onClick={() => handleThumbnailClick(index)}
                                $active={isActive}
                                $isDragOver={isDragOver}
                                draggable
                                onDragStart={(e) => handleDragStart(e, index)}
                                onDragOver={(e) => handleDragOver(e, index)}
                                onDrop={(e) => handleDrop(e, index)}
                                onDragLeave={handleDragLeave}
                            >
                                <DragHandle>
                                    <FaGripVertical />
                                </DragHandle>
                                <ThumbImage
                                    src={thumbSrc}
                                    alt={altText}
                                    onError={(e) => {
                                        e.currentTarget.src = '/images/video-icon.webp';
                                    }}
                                />
                                {isVideo && (
                                    <VideoOverlay>
                                        <FaPlay />
                                    </VideoOverlay>
                                )}
                                <RemoveBtn
                                    type="button"
                                    onClick={(ev) => {
                                        ev.stopPropagation();
                                        handleRemoveItem(index);
                                    }}
                                >
                                    <Trash2 />
                                </RemoveBtn>
                            </ThumbButton>
                        );
                    })}
                </ThumbnailsColumn>
                {/* Main "swipeable" area */}
                <MainArea {...swipeHandlers}>
                    {currentMedia ? (
                        <MainDisplay>
                            <MediaDisplayEditable
                                media={currentMedia}
                                onImageClick={handleMainImageClick}
                                isPlaying={activeVideoIndex === currentMediaIndex}
                                setIsPlaying={(playing) =>
                                    setActiveVideoIndex(playing ? currentMediaIndex : null)
                                }
                            />
                        </MainDisplay>
                    ) : (
                        <NoMediaFallback>No media available</NoMediaFallback>
                    )}
                    {/* Lightbox for images */}
                    {isLightboxOpen && currentMedia?.type === 'image' && (
                        <Suspense fallback={<LightboxFallback><Spinner /></LightboxFallback>}>
                            <LightboxViewerEditable
                                slides={validGallery.map((mItem, i) => ({
                                    src: mItem.src,
                                    type: mItem.type,
                                    alt: mItem.alt || `${productTitle} media #${i + 1}`,
                                }))}
                                currentIndex={currentMediaIndex}
                                onClose={closeLightbox}
                                onNext={() => {
                                    handleNext();
                                    setActiveVideoIndex(null);
                                }}
                                onPrev={() => {
                                    handlePrev();
                                    setActiveVideoIndex(null);
                                }}
                            />
                        </Suspense>
                    )}
                </MainArea>
            </MainLayout>
            <AddButtonsRow>
                <AddButton type="button" onClick={openImageModal}>
                    <Plus />
                    Add Image
                </AddButton>
                <AddButton type="button" onClick={openVideoModal}>
                    <Plus />
                    Add Video
                </AddButton>
            </AddButtonsRow>
            {/* === Modals for uploading === */}
            {showImageModal && (
                <Overlay onClick={closeImageModal}>
                    <ModalContent
                        onClick={(e) => e.stopPropagation()}
                        aria-modal="true"
                        role="dialog"
                    >
                        <CloseButton onClick={closeImageModal}>×</CloseButton>
                        <ModalHeader>
                            <Upload />
                            <span style={{ marginLeft: '8px' }}>Upload Image</span>
                        </ModalHeader>
                        {/* Lazy-load the image uploader */}
                        <Suspense fallback={<ModalFallback>Loading...</ModalFallback>}>
                            <LazyImageUploader
                                mediaId={mediaId}
                                onUploadSuccess={handleImageSuccess}
                            />
                        </Suspense>
                    </ModalContent>
                </Overlay>
            )}
            {showVideoModal && (
                <Overlay onClick={closeVideoModal}>
                    <ModalContent
                        onClick={(e) => e.stopPropagation()}
                        aria-modal="true"
                        role="dialog"
                    >
                        <CloseButton onClick={closeVideoModal}>×</CloseButton>
                        <ModalHeader>
                            <Upload />
                            <span style={{ marginLeft: '8px' }}>Upload Video</span>
                        </ModalHeader>
                        {/* Lazy-load the video uploader */}
                        <Suspense fallback={<ModalFallback>Loading...</ModalFallback>}>
                            <LazyVideoUploader
                                mediaId={mediaId}
                                onUploadSuccess={handleVideoSuccess}
                            />
                        </Suspense>
                    </ModalContent>
                </Overlay>
            )}
        </GalleryContainer>
    );
}
/* --------------------------------------------------------
   PropTypes / Defaults
--------------------------------------------------------- */
MediaGalleryMobileEditable.propTypes = {
    gallery: PropTypes.arrayOf(
        PropTypes.shape({
            type: PropTypes.oneOf(['image', 'video']).isRequired,
            src: PropTypes.string.isRequired,
            alt: PropTypes.string,
            poster: PropTypes.string,
        })
    ).isRequired,
    setGallery: PropTypes.func.isRequired,
    currentMediaIndex: PropTypes.number,
    setCurrentMediaIndex: PropTypes.func.isRequired,
    handleNext: PropTypes.func,
    handlePrev: PropTypes.func,
    productTitle: PropTypes.string,
    activeVideoIndex: PropTypes.number,
    setActiveVideoIndex: PropTypes.func.isRequired,
};
MediaGalleryMobileEditable.defaultProps = {
    currentMediaIndex: 0,
    handleNext: () => {},
    handlePrev: () => {},
    productTitle: 'Untitled Product',
    activeVideoIndex: null,
};
export default memo(MediaGalleryMobileEditable);
/* --------------------------------------------------------
   STYLED COMPONENTS
--------------------------------------------------------- */
const GalleryContainer = styled.div`
    display: flex;
    flex-direction: column;
    gap: 12px;
    box-sizing: border-box;
    padding: 12px;
    background: #fafafa;
`;
/* Row of "Add Image" / "Add Video" buttons */
const AddButtonsRow = styled.div`
    display: flex;
    gap: 12px;
    margin-bottom: 8px;
`;
const AddButton = styled.button`
    background: #007bc7;
    color: #fff;
    border: none;
    padding: 6px 10px;
    border-radius: 4px;
    font-size: 0.85rem;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    &:hover {
        background: #005e99;
    }
`;
/* The actual modals + overlay */
const Overlay = styled.div`
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
`;
const ModalContent = styled.div`
    background: #fff;
    border-radius: 8px;
    min-width: 320px;
    max-width: 480px;
    padding: 16px;
    position: relative;
`;
const CloseButton = styled.button`
    position: absolute;
    top: 8px;
    right: 10px;
    border: none;
    background: none;
    color: #999;
    font-size: 1.2rem;
    cursor: pointer;
    &:hover {
        color: #333;
    }
`;
const ModalHeader = styled.h3`
    margin: 0;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 1rem;
    color: #333;
    margin-bottom: 12px;
`;
const ModalFallback = styled.div`
    text-align: center;
    padding: 30px;
    font-size: 1rem;
`;
/* The main layout of thumbnails (left) + main area (right) */
const MainLayout = styled.div`
    display: grid;
    grid-template-columns: 80px 1fr;
    gap: 8px;
    @media (max-width: 768px) {
        grid-template-columns: 70px 1fr;
        gap: 6px;
    }
`;
const ThumbnailsColumn = styled.div`
    display: flex;
    flex-direction: column;
    gap: 6px;
    overflow-y: auto;
    max-height: 320px;
    &::-webkit-scrollbar { width: 0; }
    scrollbar-width: none;
    -ms-overflow-style: none;
`;
const ThumbButton = styled.button`
    position: relative;
    width: 100%;
    height: 60px;
    border: none;
    background: ${({ $active }) => ($active ? 'rgba(0, 123, 255, 0.15)' : '#fff')};
    border-radius: 4px;
    cursor: pointer;
    overflow: hidden;
    display: flex;
    align-items: center;
    margin-bottom: 2px;
    ${({ $active }) =>
            $active &&
            css`
                &:hover {
                    background: rgba(0, 123, 255, 0.25);
                }
            `}
    ${({ $isDragOver }) =>
            $isDragOver &&
            css`
                outline: 2px dashed #00c9a7;
                background: #e8faf7;
            `}
    &:focus {
        outline: 2px solid #66afe9;
        outline-offset: 2px;
    }
`;
const DragHandle = styled.div`
    width: 24px;
    height: 100%;
    background: #f2f2f2;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #888;
    cursor: move;
`;
const ThumbImage = styled.img`
    flex: 1;
    height: 100%;
    object-fit: cover;
`;
const RemoveBtn = styled.button`
    position: absolute;
    top: 4px;
    right: 4px;
    border: none;
    background: rgba(255, 0, 0, 0.75);
    color: #fff;
    padding: 4px;
    border-radius: 4px;
    font-size: 0.75rem;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    &:hover {
        background: rgba(255, 0, 0, 0.9);
    }
`;
const VideoOverlay = styled.div`
    position: absolute;
    inset: 0;
    background: rgba(0,0,0,0.25);
    display: flex;
    align-items: center;
    justify-content: center;
    svg {
        color: #fff;
        width: 20px;
        height: 20px;
    }
`;
const MainArea = styled.div`
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
`;
const MainDisplay = styled.div`
    width: 100%;
    max-width: 600px;
    aspect-ratio: 4/3;
    background: #000;
    border-radius: 6px;
    overflow: hidden;
    position: relative;
    @media (max-width: 768px) {
        max-width: 100%;
        aspect-ratio: unset;
        height: auto;
    }
`;
const NoMediaFallback = styled.div`
    width: 100%;
    min-height: 220px;
    border-radius: 6px;
    background: #222;
    color: #bbb;
    display: flex;
    align-items: center;
    justify-content: center;
`;
const LightboxFallback = styled.div`
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.75);
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
`;

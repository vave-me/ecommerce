"use client";
import React, { useState, useRef, useCallback, useEffect, useMemo, memo } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useTranslations} from 'next-intl';
import {
    Send,
    Loader2,
    Mic,
    MicOff,
    Image as ImageIcon,
    Upload,
    X,
    Paperclip,
    Plus,
    AlertCircle,
    CheckCircle
} from '@/icons';
// Removed undefined imports from appModeSlice
import { getAssistantService } from '../../services/ai/AssistantServiceWithRetry';
import {useAuth} from '../../context/AuthContext';
import {redirectToLogin, getCurrentPageUrl} from '../../utils/redirectUtils';
import useMediaPermissions from '../../hooks/useMediaPermissions';
import styles from './PromptInput.module.css';
/**
 * Modern AI Chat Input - Inspired by ChatGPT/Claude
 * Clean, minimal design with excellent UX
 * Optimized with React.memo for performance
 */
const PromptInput = memo(({
                                    onResponse = () => {
                                    },
                                    onFocus = () => {
                                    },
                                    onBlur = () => {
                                    },
                                    placeholder = "Message AI assistant...",
                                    disabled = false,
                                    className = '',
                                    conversationId = null
                                }) => {
    const dispatch = useDispatch();
    const t = useTranslations('AI');
    const {user, signInWithGoogle} = useAuth();
    // const lastPrompt = useSelector(selectLastPrompt); // Removed - selector doesn't exist
    const {
        cameraPermission,
        micPermission,
        requestCameraPermission,
        requestMicPermission,
        getPermissionError,
        clearPermissionError
    } = useMediaPermissions();
    // Format recording duration
    const formatDuration = useCallback((seconds) => {
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }, []);
    // Draw real-time waveform visualization
    const drawWaveform = useCallback(() => {
        const canvas = canvasRef.current;
        if (!canvas || !waveformDataRef.current.length) return;
        const ctx = canvas.getContext('2d');
        const width = canvas.width;
        const height = canvas.height;
        // Clear canvas
        ctx.fillStyle = 'rgba(0, 0, 0, 0.1)';
        ctx.fillRect(0, 0, width, height);
        // Draw waveform
        ctx.lineWidth = 2;
        ctx.strokeStyle = '#0066cc';
        ctx.beginPath();
        const sliceWidth = width / waveformDataRef.current.length;
        let x = 0;
        for (let i = 0; i < waveformDataRef.current.length; i++) {
            const v = waveformDataRef.current[i] / 128.0;
            const y = (v * height) / 2;
            if (i === 0) {
                ctx.moveTo(x, y);
            } else {
                ctx.lineTo(x, y);
            }
            x += sliceWidth;
        }
        ctx.lineTo(width, height / 2);
        ctx.stroke();
    }, []);
    // Draw frequency spectrum bars
    const drawSpectrum = useCallback(() => {
        const canvas = spectrumCanvasRef.current;
        if (!canvas || !frequencyDataRef.current.length) return;
        const ctx = canvas.getContext('2d');
        const width = canvas.width;
        const height = canvas.height;
        // Clear canvas
        ctx.fillStyle = 'rgba(0, 0, 0, 0.1)';
        ctx.fillRect(0, 0, width, height);
        const barWidth = width / frequencyDataRef.current.length;
        let x = 0;
        for (let i = 0; i < frequencyDataRef.current.length; i++) {
            const barHeight = (frequencyDataRef.current[i] / 255) * height;
            // Create gradient for bars
            const gradient = ctx.createLinearGradient(0, height - barHeight, 0, height);
            gradient.addColorStop(0, '#0066cc');
            gradient.addColorStop(0.5, '#0088ff');
            gradient.addColorStop(1, '#00aaff');
            ctx.fillStyle = gradient;
            ctx.fillRect(x, height - barHeight, barWidth - 1, barHeight);
            x += barWidth;
        }
    }, []);
    // Enhanced state with performance optimizations
    const [prompt, setPrompt] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [successMessage, setSuccessMessage] = useState(null);
    const [attachedFiles, setAttachedFiles] = useState([]);
    const [isRecording, setIsRecording] = useState(false);
    const [isDragOver, setIsDragOver] = useState(false);
    const [recordingDuration, setRecordingDuration] = useState(0);
    const [audioLevel, setAudioLevel] = useState(0);
    const [recordingProgress, setRecordingProgress] = useState(0);
    const [isOptimisticUpdate, setIsOptimisticUpdate] = useState(false);
    const [retryCount, setRetryCount] = useState(0);
    // Voice recording state
    const [recordingStatus, setRecordingStatus] = useState(null); // 'processing', 'success', 'error'
    const [isProcessingAudio, setIsProcessingAudio] = useState(false);
    const [lastRecordingResult, setLastRecordingResult] = useState(null);
    // Animation loop for visualizations (optimized with performance check)
    const updateVisualizations = useCallback(() => {
        if (isRecording && isMountedRef.current) {
            drawWaveform();
            drawSpectrum();
            animationFrameRef.current = requestAnimationFrame(updateVisualizations);
        }
    }, [isRecording, drawWaveform, drawSpectrum]);
    // Refs
    const textareaRef = useRef(null);
    const fileInputRef = useRef(null);
    const mediaRecorderRef = useRef(null);
    const audioContextRef = useRef(null);
    const analyserRef = useRef(null);
    const recordingTimerRef = useRef(null);
    const animationFrameRef = useRef(null);
    const streamRef = useRef(null);
    const waveformDataRef = useRef([]);
    const frequencyDataRef = useRef([]);
    const canvasRef = useRef(null);
    const spectrumCanvasRef = useRef(null);
    const debounceTimerRef = useRef(null);
    const isMountedRef = useRef(true);
    // Auto-resize textarea
    const adjustTextareaHeight = useCallback(() => {
        const textarea = textareaRef.current;
        if (!textarea) return;
        textarea.style.height = 'auto';
        const newHeight = Math.min(textarea.scrollHeight, 120);
        textarea.style.height = `${newHeight}px`;
    }, []);
    // Debounced input handling for better performance
    const debouncedPromptUpdate = useCallback(
        (value) => {
            // dispatch(setLastPrompt(value)); // Removed - action doesn't exist
        },
        [dispatch]
    );
    // Handle input changes with debouncing
    const handleInputChange = useCallback((e) => {
        const value = e.target.value;
        setPrompt(value);
        adjustTextareaHeight();
        setError(null);
        setSuccessMessage(null);
        // Debounce Redux update to prevent excessive dispatches
        if (debounceTimerRef.current) {
            clearTimeout(debounceTimerRef.current);
        }
        debounceTimerRef.current = setTimeout(() => {
            debouncedPromptUpdate(value);
        }, 300);
    }, [adjustTextareaHeight, debouncedPromptUpdate]);
    // Handle keyboard shortcuts
    const handleKeyDown = useCallback((e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            if (user && (prompt.trim() || attachedFiles.length > 0)) {
                handleSubmit(e);
            }
        }
    }, [user, prompt, attachedFiles]);
    // Enhanced form submission with retry logic and optimistic updates
    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();
        if (!user) {
            setError('Please sign in to continue');
            return;
        }
        if (!prompt.trim() && attachedFiles.length === 0) return;
        const messageContent = prompt.trim();
        const currentFiles = [...attachedFiles];
        setIsLoading(true);
        setError(null);
        setSuccessMessage(null);
        setIsOptimisticUpdate(true);
        // Optimistic UI update - immediately clear input
        setPrompt('');
        setAttachedFiles([]);
        adjustTextareaHeight();
        const attemptSubmit = async (attempt = 1) => {
            try {
                await onResponse(messageContent, currentFiles);
                // Success handling
                // dispatch(setLastPrompt('')); // Removed - action doesn't exist
                setSuccessMessage('Message sent successfully');
                setRetryCount(0);
                // Clear success message after 3 seconds
                setTimeout(() => setSuccessMessage(null), 3000);
            } catch (error) {
                // Restore input on error for retry
                setPrompt(messageContent);
                setAttachedFiles(currentFiles);
                if (attempt < 3 && error.name !== 'AbortError') {
                    // Retry with exponential backoff
                    const delay = Math.min(1000 * Math.pow(2, attempt - 1), 5000);
                    setTimeout(() => {
                        setRetryCount(attempt);
                        attemptSubmit(attempt + 1);
                    }, delay);
                } else {
                    // Final failure
                    const errorMessage = error.response?.data?.message ||
                        error.message ||
                        'Failed to send message. Please try again.';
                    setError(errorMessage);
                    setRetryCount(0);
                }
            }
        };
        try {
            await attemptSubmit();
        } finally {
            setIsLoading(false);
            setIsOptimisticUpdate(false);
        }
    }, [user, prompt, attachedFiles, onResponse, dispatch, adjustTextareaHeight]);
    // Enhanced file handling with validation and compression
    const handleFileSelect = useCallback(async (files) => {
        const maxFileSize = 10 * 1024 * 1024; // 10MB
        const maxFiles = 5; // Limit concurrent uploads
        if (attachedFiles.length + files.length > maxFiles) {
            setError(`Maximum ${maxFiles} files allowed`);
            return;
        }
        const validFiles = Array.from(files).filter(file => {
            const isValidType = file.type.startsWith('image/') ||
                file.type === 'application/pdf' ||
                file.type.startsWith('text/');
            const isValidSize = file.size <= maxFileSize;
            if (!isValidType) {
                setError(`Unsupported file type: ${file.name}`);
                return false;
            }
            if (!isValidSize) {
                setError(`File too large: ${file.name} (max 10MB)`);
                return false;
            }
            return true;
        });
        if (validFiles.length === 0) return;
        // Process files with loading state
        setIsLoading(true);
        try {
            const processedFiles = await Promise.all(
                validFiles.map(async (file) => {
                    const preview = file.type.startsWith('image/')
                        ? URL.createObjectURL(file)
                        : null;
                    // Generate file hash for deduplication
                    const arrayBuffer = await file.arrayBuffer();
                    const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
                    const hashArray = Array.from(new Uint8Array(hashBuffer));
                    const hash = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
                    return {
                        file,
                        name: file.name,
                        type: file.type,
                        size: file.size,
                        preview,
                        hash,
                        uploadProgress: 0
                    };
                })
            );
            // Check for duplicates
            const existingHashes = new Set(attachedFiles.map(f => f.hash));
            const uniqueFiles = processedFiles.filter(f => !existingHashes.has(f.hash));
            if (uniqueFiles.length !== processedFiles.length) {
                setError('Some files were skipped (duplicates detected)');
            }
            setAttachedFiles(prev => [...prev, ...uniqueFiles]);
            setError(null);
        } catch (error) {
            setError('Failed to process files');
        } finally {
            setIsLoading(false);
        }
    }, [attachedFiles]);
    const handleFileInputChange = useCallback((e) => {
        if (e.target.files && e.target.files.length > 0) {
            handleFileSelect(e.target.files);
        }
    }, [handleFileSelect]);
    const removeAttachment = useCallback((index) => {
        setAttachedFiles(prev => {
            const newFiles = [...prev];
            // Revoke object URL to prevent memory leaks
            if (newFiles[index].preview) {
                URL.revokeObjectURL(newFiles[index].preview);
            }
            newFiles.splice(index, 1);
            return newFiles;
        });
    }, []);
    // Drag and drop handling
    const handleDragOver = useCallback((e) => {
        e.preventDefault();
        setIsDragOver(true);
    }, []);
    const handleDragLeave = useCallback((e) => {
        e.preventDefault();
        setIsDragOver(false);
    }, []);
    const handleDrop = useCallback((e) => {
        e.preventDefault();
        setIsDragOver(false);
        if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
            handleFileSelect(e.dataTransfer.files);
        }
    }, [handleFileSelect]);
    // Audio level visualization with enhanced frequency analysis
    const updateAudioLevel = useCallback(() => {
        if (!analyserRef.current || !isRecording) return;
        // Enhanced frequency analysis with multiple data arrays
        const frequencyDataArray = new Uint8Array(analyserRef.current.frequencyBinCount);
        const timeDomainArray = new Uint8Array(analyserRef.current.frequencyBinCount);
        analyserRef.current.getByteFrequencyData(frequencyDataArray);
        analyserRef.current.getByteTimeDomainData(timeDomainArray);
        // Calculate RMS (Root Mean Square) for audio level
        let frequencySum = 0;
        let timeDomainSum = 0;
        for (let i = 0; i < frequencyDataArray.length; i++) {
            frequencySum += frequencyDataArray[i] * frequencyDataArray[i];
        }
        for (let i = 0; i < timeDomainArray.length; i++) {
            const normalizedValue = (timeDomainArray[i] - 128) / 128;
            timeDomainSum += normalizedValue * normalizedValue;
        }
        const frequencyRMS = Math.sqrt(frequencySum / frequencyDataArray.length);
        const timeDomainRMS = Math.sqrt(timeDomainSum / timeDomainArray.length);
        // Combine frequency and time domain for better audio level detection
        const combinedLevel = Math.max(
            Math.min(100, (frequencyRMS / 128) * 100),
            Math.min(100, timeDomainRMS * 100)
        );
        setAudioLevel(combinedLevel);
        // Store waveform data for visualization
        if (waveformDataRef.current) {
            // Downsample for visualization (take every 4th sample for performance)
            const downsampledWaveform = [];
            for (let i = 0; i < timeDomainArray.length; i += 4) {
                downsampledWaveform.push(timeDomainArray[i]);
            }
            waveformDataRef.current = downsampledWaveform;
        }
        // Store frequency data for spectrum visualization
        if (frequencyDataRef.current) {
            // Take first 32 frequency bins for visualization bars
            frequencyDataRef.current = Array.from(frequencyDataArray.slice(0, 32));
        }
        if (isRecording) {
            animationFrameRef.current = requestAnimationFrame(updateAudioLevel);
        }
    }, [isRecording]);
    // Recording timer
    const startRecordingTimer = useCallback(() => {
        setRecordingDuration(0);
        setRecordingProgress(0);
        recordingTimerRef.current = setInterval(() => {
            setRecordingDuration(prev => {
                const newDuration = prev + 1;
                // Calculate progress (max 60 seconds)
                const progress = Math.min(100, (newDuration / 60) * 100);
                setRecordingProgress(progress);
                // Auto-stop at 60 seconds
                if (newDuration >= 60) {
                    stopRecording();
                }
                return newDuration;
            });
        }, 1000);
    }, []);
    const stopRecordingTimer = useCallback(() => {
        if (recordingTimerRef.current) {
            clearInterval(recordingTimerRef.current);
            recordingTimerRef.current = null;
        }
        if (animationFrameRef.current) {
            cancelAnimationFrame(animationFrameRef.current);
            animationFrameRef.current = null;
        }
    }, []);
    // Process audio recording with new Swagger-compliant API
    const processAudioRecording = useCallback(async (audioBlob, mimeType) => {
        setIsProcessingAudio(true);
        setRecordingStatus('processing');
        setError(null);
        try {
            // Show processing feedback
            setSuccessMessage(`Processing ${recordingDuration}s recording...`);
            // Create a more detailed recording result
            const recordingResult = {
                blob: audioBlob,
                mimeType,
                size: audioBlob.size,
                duration: recordingDuration,
                timestamp: new Date().toISOString()
            };
            // Try speech-to-text with new API client
            try {
                // Convert blob to base64 for API
                const arrayBuffer = await audioBlob.arrayBuffer();
                const base64Audio = btoa(String.fromCharCode(...new Uint8Array(arrayBuffer)));
                const assistantService = getAssistantService();
                const speechResult = await assistantService.processSpeechInput({
                    audioData: base64Audio,
                    audioFormat: mimeType.includes('webm') ? 'webm' : 'mp3',
                    language: 'en',
                    context: {
                        duration: recordingDuration,
                        timestamp: new Date().toISOString()
                    }
                });
                if (speechResult.success && speechResult.data?.transcription) {
                    // Auto-fill the input with transcribed text
                    const transcription = speechResult.data.transcription;
                    setPrompt(transcription);
                    // dispatch(setLastPrompt(transcription)); // Removed - action doesn't exist
                    adjustTextareaHeight();
                    setRecordingStatus('success');
                    setSuccessMessage(`✓ Voice recorded and transcribed (${recordingDuration}s)`);
                    setLastRecordingResult({
                        ...recordingResult,
                        transcription: transcription,
                        confidence: speechResult.data.confidence
                    });
                } else {
                    throw new Error(speechResult.error || 'Speech processing failed');
                }
            } catch (apiError) {
                // Fallback: Save as audio attachment
                const audioFile = new File([audioBlob], `voice_${Date.now()}.${mimeType.includes('webm') ? 'webm' : 'mp3'}`, {
                    type: mimeType
                });
                const audioAttachment = {
                    file: audioFile,
                    name: audioFile.name,
                    type: audioFile.type,
                    size: audioFile.size,
                    preview: null, // No preview for audio
                    hash: `voice_${Date.now()}`,
                    isVoiceRecording: true,
                    duration: recordingDuration
                };
                setAttachedFiles(prev => [...prev, audioAttachment]);
                setRecordingStatus('success');
                setSuccessMessage(`✓ Voice recorded as audio file (${recordingDuration}s)`);
                setLastRecordingResult({
                    ...recordingResult,
                    attachment: audioAttachment,
                    fallbackMode: true
                });
            }
        } catch (error) {
            setRecordingStatus('error');
            setError(`Failed to process recording: ${error.message}`);
            setLastRecordingResult({
                ...recordingResult,
                error: error.message
            });
        } finally {
            setIsProcessingAudio(false);
            // Clear success message after delay
            setTimeout(() => {
                setSuccessMessage(null);
                setRecordingStatus(null);
            }, 4000);
        }
    }, [recordingDuration, dispatch, adjustTextareaHeight]);
    // Voice recording
    const startRecording = useCallback(async () => {
        // Clear any previous permission errors
        if (clearPermissionError) {
            clearPermissionError('microphone');
        }
        setError(null);
        if (micPermission !== 'granted') {
            const permission = await requestMicPermission();
            if (permission !== 'granted') {
                const permissionError = getPermissionError('microphone');
                setError(permissionError || 'Microphone permission is required for voice recording');
                return;
            }
        }
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: {
                    echoCancellation: true,
                    noiseSuppression: true,
                    autoGainControl: true,
                    sampleRate: 48000, // Higher quality
                    channelCount: 1, // Mono
                    volume: 1.0
                }
            });
            streamRef.current = stream;
            // Set up audio analysis for visualization
            if (window.AudioContext || window.webkitAudioContext) {
                try {
                    const AudioContext = window.AudioContext || window.webkitAudioContext;
                    audioContextRef.current = new AudioContext();
                    analyserRef.current = audioContextRef.current.createAnalyser();
                    const source = audioContextRef.current.createMediaStreamSource(stream);
                    source.connect(analyserRef.current);
                    // Enhanced analyser settings for better visualization
                    analyserRef.current.fftSize = 512; // Higher resolution
                    analyserRef.current.smoothingTimeConstant = 0.3; // Less smoothing for more responsive visualization
                    analyserRef.current.minDecibels = -100;
                    analyserRef.current.maxDecibels = -10;
                    // Start audio level monitoring
                    updateAudioLevel();
                } catch (audioError) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', audioError);
        }
    }
            }
            // Determine best supported format
            let mimeType = 'audio/webm;codecs=opus';
            if (!MediaRecorder.isTypeSupported(mimeType)) {
                mimeType = 'audio/webm';
                if (!MediaRecorder.isTypeSupported(mimeType)) {
                    mimeType = 'audio/mp4';
                    if (!MediaRecorder.isTypeSupported(mimeType)) {
                        mimeType = ''; // Use default
                    }
                }
            }
            mediaRecorderRef.current = new MediaRecorder(stream, {
                mimeType: mimeType || undefined,
                audioBitsPerSecond: 128000 // 128 kbps for good quality
            });
            const chunks = [];
            mediaRecorderRef.current.ondataavailable = (e) => {
                if (e.data.size > 0) {
                    chunks.push(e.data);
                }
            };
            mediaRecorderRef.current.onstop = async () => {
                const finalMimeType = mediaRecorderRef.current.mimeType || 'audio/webm';
                const blob = new Blob(chunks, {type: finalMimeType});
                // Clean up resources
                stopRecordingTimer();
                setAudioLevel(0);
                if (streamRef.current) {
                    streamRef.current.getTracks().forEach(track => track.stop());
                    streamRef.current = null;
                }
                if (audioContextRef.current) {
                    audioContextRef.current.close();
                    audioContextRef.current = null;
                }
                // Process the audio recording
                await processAudioRecording(blob, finalMimeType);
            };
            mediaRecorderRef.current.onerror = (event) => {
                setError('Recording error: ' + event.error.message);
                stopRecording();
            };
            mediaRecorderRef.current.start(250); // Collect data every 250ms for better responsiveness
            setIsRecording(true);
            startRecordingTimer();
            // Start visualization
            updateVisualizations();
        } catch (error) {
            let errorMessage = 'Failed to start recording';
            switch (error.name) {
                case 'NotAllowedError':
                    errorMessage = 'Microphone permission denied. Please allow microphone access.';
                    break;
                case 'NotFoundError':
                    errorMessage = 'No microphone found. Please ensure a microphone is connected.';
                    break;
                case 'NotReadableError':
                    errorMessage = 'Microphone is already in use by another application.';
                    break;
                case 'NotSupportedError':
                    errorMessage = 'Voice recording is not supported by your browser.';
                    break;
                case 'OverconstrainedError':
                    errorMessage = 'Microphone settings could not be satisfied. Please try again.';
                    break;
                default:
                    errorMessage = error.message || 'Failed to start recording';
            }
            setError(errorMessage);
        }
    }, [micPermission, requestMicPermission, getPermissionError, clearPermissionError, updateAudioLevel, startRecordingTimer, stopRecordingTimer, recordingDuration]);
    const stopRecording = useCallback(() => {
        if (!isRecording) return;
        try {
            // Stop the MediaRecorder
            if (mediaRecorderRef.current && mediaRecorderRef.current.state === 'recording') {
                mediaRecorderRef.current.stop();
            }
            // Stop timers and animations
            stopRecordingTimer();
            // Stop audio stream
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => track.stop());
                streamRef.current = null;
            }
            // Close audio context
            if (audioContextRef.current) {
                audioContextRef.current.close();
                audioContextRef.current = null;
            }
            setIsRecording(false);
            setAudioLevel(0);
        } catch (error) {
            setError('Failed to stop recording');
            setIsRecording(false);
            setAudioLevel(0);
            stopRecordingTimer();
        }
    }, [isRecording, stopRecordingTimer]);
    // Focus textarea on mount
    useEffect(() => {
        if (textareaRef.current) {
            textareaRef.current.focus();
        }
    }, []);
    // Clean up object URLs on unmount
    useEffect(() => {
        return () => {
            attachedFiles.forEach(file => {
                if (file.preview) {
                    URL.revokeObjectURL(file.preview);
                }
            });
        };
    }, [attachedFiles]);
    // Clean up recording on unmount
    useEffect(() => {
        return () => {
            // Clean up recording
            if (mediaRecorderRef.current && mediaRecorderRef.current.state === 'recording') {
                try {
                    mediaRecorderRef.current.stop();
                } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
            }
            // Clean up timers
            if (recordingTimerRef.current) {
                clearInterval(recordingTimerRef.current);
            }
            if (animationFrameRef.current) {
                cancelAnimationFrame(animationFrameRef.current);
            }
            // Clean up stream
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => track.stop());
            }
            // Clean up audio context
            if (audioContextRef.current) {
                audioContextRef.current.close();
            }
            // Mark as unmounted
            isMountedRef.current = false;
        };
    }, []);
    // Memoized performance calculations
    const memoizedValues = useMemo(() => {
        const canSend = user && (prompt.trim() || attachedFiles.length > 0) && !isLoading && !isProcessingAudio;
        const hasAttachments = attachedFiles.length > 0;
        const totalFileSize = attachedFiles.reduce((sum, file) => sum + file.size, 0);
        const formattedFileSize = totalFileSize > 0 ?
            `${(totalFileSize / (1024 * 1024)).toFixed(1)}MB` : '';
        const hasVoiceAttachments = attachedFiles.some(file => file.isVoiceRecording);
        const isAnyProcessing = isLoading || isProcessingAudio;
        return {
            canSend,
            hasAttachments,
            hasVoiceAttachments,
            totalFileSize,
            formattedFileSize,
            isRecordingLongDuration: recordingDuration > 30,
            isAnyProcessing,
            shouldShowHelperText: prompt.length === 0 && !hasAttachments && !isRecording && !isProcessingAudio
        };
    }, [user, prompt, attachedFiles, isLoading, recordingDuration, isProcessingAudio]);
    if (!user) {
        return (
            <div className={styles.loginPrompt}>
                <p>Please sign in to chat with AI assistants</p>
                <button
                    onClick={() => signInWithGoogle()}
                    className={styles.loginButton}
                >
                    Sign in with Google
                </button>
            </div>
        );
    }
    return (
        <div className={`${styles.container} ${className}`}>
            {/* Enhanced error display with retry functionality */}
            {error && (
                <div className={styles.error} role="alert" aria-live="polite">
                    <AlertCircle size={16}/>
                    <span>{error}</span>
                    {retryCount > 0 && (
                        <span className={styles.retryInfo}>
                            (Retry {retryCount}/3)
                        </span>
                    )}
                    <button
                        onClick={() => setError(null)}
                        aria-label="Dismiss error"
                        className={styles.dismissButton}
                    >
                        ×
                    </button>
                </div>
            )}
            {/* Success message display */}
            {successMessage && (
                <div className={styles.success} role="status" aria-live="polite">
                    <CheckCircle size={16}/>
                    <span>{successMessage}</span>
                </div>
            )}
            {/* Enhanced attachments preview with file info */}
            {memoizedValues.hasAttachments && (
                <div className={styles.attachmentsPreview}>
                    <div className={styles.attachmentsHeader}>
                        <span className={styles.attachmentsCount}>
                            {attachedFiles.length} file{attachedFiles.length !== 1 ? 's' : ''}
                        </span>
                        {memoizedValues.formattedFileSize && (
                            <span className={styles.attachmentsSize}>
                                ({memoizedValues.formattedFileSize})
                            </span>
                        )}
                    </div>
                    <div className={styles.attachmentsList}>
                        {attachedFiles.map((file, index) => (
                            <div key={file.hash || index} className={styles.attachmentItem}>
                                {file.preview ? (
                                    <img
                                        src={file.preview}
                                        alt={file.name}
                                        className={styles.attachmentImage}
                                        loading="lazy"
                                    />
                                ) : (
                                    <div
                                        className={`${styles.attachmentIcon} ${file.isVoiceRecording ? styles.voiceIcon : ''}`}>
                                        {file.isVoiceRecording ? <Mic size={16}/> : <Paperclip size={16}/>}
                                    </div>
                                )}
                                <div className={styles.attachmentInfo}>
                                    <span className={styles.attachmentName} title={file.name}>
                                        {file.name}
                                    </span>
                                    <span className={styles.attachmentMeta}>
                                        {(file.size / (1024 * 1024)).toFixed(1)}MB
                                    </span>
                                </div>
                                <button
                                    onClick={() => removeAttachment(index)}
                                    className={styles.removeAttachment}
                                    title="Remove attachment"
                                    aria-label={`Remove ${file.name}`}
                                >
                                    <X size={12}/>
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            {/* Bright Theme Recording Bar */}
            {isRecording && (
                <div className={styles.recordingBar}>
                    <div className={styles.recordingIndicator}>
                        <div className={styles.recordingPulse}></div>
                        <span className={styles.recordingLabel}>REC</span>
                        <span className={styles.recordingTimer}>
                            {formatDuration(recordingDuration)}
                        </span>
                    </div>
                    <div className={styles.waveformDisplay}>
                        <canvas
                            ref={canvasRef}
                            className={styles.waveformCanvas}
                            width={200}
                            height={32}
                        />
                    </div>
                    <div className={styles.progressIndicator}>
                        <div className={styles.progressBar}>
                            <div
                                className={styles.progressFill}
                                style={{width: `${recordingProgress}%`}}
                            ></div>
                        </div>
                    </div>
                    <div className={styles.recordingActions}>
                        <button
                            type="button"
                            onClick={stopRecording}
                            className={styles.stopRecordingBtn}
                            title="Stop and send recording"
                        >
                            <div className={styles.stopIcon}></div>
                        </button>
                        <button
                            type="button"
                            onClick={() => {
                                stopRecording();
                                setPrompt('');
                            }}
                            className={styles.cancelRecordingBtn}
                            title="Cancel recording"
                        >
                            ×
                        </button>
                    </div>
                </div>
            )}
            <form onSubmit={handleSubmit} className={styles.form}>
                <div
                    className={`${styles.inputWrapper} ${isDragOver ? styles.dragOver : ''}`}
                    onDragOver={handleDragOver}
                    onDragLeave={handleDragLeave}
                    onDrop={handleDrop}
                >
                    {isDragOver && (
                        <div className={styles.dragOverlay}>
                            <Upload size={32}/>
                            <span>Drop images here</span>
                        </div>
                    )}
                    <textarea
                        ref={textareaRef}
                        value={prompt}
                        onChange={handleInputChange}
                        onKeyDown={handleKeyDown}
                        onFocus={onFocus}
                        onBlur={onBlur}
                        placeholder={placeholder}
                        disabled={disabled || isLoading}
                        className={styles.textarea}
                        rows={1}
                    />
                    <div className={styles.actionButtons}>
                        {/* File upload */}
                        <button
                            type="button"
                            onClick={() => fileInputRef.current?.click()}
                            className={styles.iconButton}
                            disabled={disabled || isLoading}
                            title="Attach image"
                        >
                            <Paperclip size={20}/>
                        </button>
                        {/* Voice recording */}
                        <button
                            type="button"
                            onClick={isRecording ? stopRecording : startRecording}
                            className={`${styles.iconButton} ${isRecording ? styles.recording : ''}`}
                            disabled={disabled || isLoading}
                            title={isRecording ? "Stop recording" : "Start voice recording"}
                        >
                            {isRecording ? <MicOff size={20}/> : <Mic size={20}/>}
                        </button>
                        {/* Send button with enhanced states */}
                        <button
                            type="submit"
                            disabled={!memoizedValues.canSend}
                            className={`${styles.sendButton} ${isOptimisticUpdate ? styles.optimistic : ''}`}
                            title={
                                isLoading ? 'Sending...' :
                                    retryCount > 0 ? `Retrying (${retryCount}/3)` :
                                        'Send message'
                            }
                            aria-label={
                                isLoading ? 'Sending message' :
                                    retryCount > 0 ? `Retrying message, attempt ${retryCount}` :
                                        'Send message'
                            }
                        >
                            {isLoading ? (
                                <Loader2 size={20} className={styles.spinner}/>
                            ) : (
                                <Send size={20}/>
                            )}
                        </button>
                    </div>
                    <input
                        ref={fileInputRef}
                        type="file"
                        accept="image/*"
                        multiple
                        onChange={handleFileInputChange}
                        className={styles.hiddenInput}
                    />
                </div>
            </form>
            {/* Compact helper text */}
            {memoizedValues.shouldShowHelperText && (
                <div className={styles.compactHelperText}>
                    <span>Enter to send</span>
                    <span>•</span>
                    <span>Shift+Enter for line</span>
                    <span>•</span>
                    <span>📎 Attach files</span>
                    {micPermission === 'granted' && (
                        <>
                            <span>•</span>
                            <span>🎤 Voice</span>
                        </>
                    )}
                </div>
            )}
            {/* Recording duration warning */}
            {memoizedValues.isRecordingLongDuration && (
                <div className={styles.recordingWarning} role="status">
                    <AlertCircle size={14}/>
                    <span>Recording will auto-stop at 60 seconds</span>
                </div>
            )}
        </div>
    );
});
PromptInput.displayName = 'PromptInput';
export default PromptInput; 
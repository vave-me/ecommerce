package streaming

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Transcoder handles video transcoding using FFmpeg
type Transcoder struct {
	profiles      []QualityProfile
	config        *StreamingConfig
	ffmpegPath    string
	activeJobs    map[string]*TranscodeJob
	mu            sync.RWMutex
}

// TranscodeJob represents an active transcoding job
type TranscodeJob struct {
	ID            string
	Profile       string
	InputURL      string
	OutputPath    string
	Cmd           *exec.Cmd
	StartTime     time.Time
	BytesProcessed int64
	FramesProcessed int64
	Cancel        context.CancelFunc
}

// TranscodeOptions contains transcoding parameters
type TranscodeOptions struct {
	InputProtocol   string
	InputURL        string
	OutputFormat    string
	SegmentDuration int
	SegmentPattern  string
	PlaylistPath    string
	VideoCodec      string
	VideoBitrate    int
	VideoWidth      int
	VideoHeight     int
	VideoFramerate  int
	VideoProfile    string
	VideoLevel      string
	AudioCodec      string
	AudioBitrate    int
	LowLatency      bool
	HardwareAccel   string // nvidia, vaapi, qsv, videotoolbox
}

// NewTranscoder creates a new transcoder instance
func NewTranscoder(profiles []QualityProfile, config *StreamingConfig) *Transcoder {
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // Assume it's in PATH
	}

	return &Transcoder{
		profiles:   profiles,
		config:     config,
		ffmpegPath: ffmpegPath,
		activeJobs: make(map[string]*TranscodeJob),
	}
}

// TranscodeToQuality starts transcoding to a specific quality profile
func (t *Transcoder) TranscodeToQuality(profileName string, opts TranscodeOptions) error {
	jobID := fmt.Sprintf("%s_%s_%d", opts.InputProtocol, profileName, time.Now().Unix())

	// Build FFmpeg command
	args := t.buildFFmpegArgs(opts)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)

	// Set up logging
	cmd.Stdout = &transcodeLogger{jobID: jobID, logType: "stdout"}
	cmd.Stderr = &transcodeLogger{jobID: jobID, logType: "stderr"}

	job := &TranscodeJob{
		ID:         jobID,
		Profile:    profileName,
		InputURL:   opts.InputURL,
		OutputPath: opts.PlaylistPath,
		Cmd:        cmd,
		StartTime:  time.Now(),
		Cancel:     cancel,
	}

	t.mu.Lock()
	t.activeJobs[jobID] = job
	t.mu.Unlock()

	// Start transcoding
	go func() {
		defer func() {
			t.mu.Lock()
			delete(t.activeJobs, jobID)
			t.mu.Unlock()
		}()

		if err := cmd.Run(); err != nil {
			fmt.Printf("Transcoding error for job %s: %v\n", jobID, err)
		}
	}()

	return nil
}

// buildFFmpegArgs constructs FFmpeg command arguments
func (t *Transcoder) buildFFmpegArgs(opts TranscodeOptions) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
	}

	// Hardware acceleration
	if opts.HardwareAccel != "" {
		switch opts.HardwareAccel {
		case "nvidia":
			args = append(args, "-hwaccel", "cuda", "-hwaccel_output_format", "cuda")
		case "vaapi":
			args = append(args, "-hwaccel", "vaapi", "-hwaccel_device", "/dev/dri/renderD128")
		case "qsv":
			args = append(args, "-hwaccel", "qsv")
		case "videotoolbox":
			args = append(args, "-hwaccel", "videotoolbox")
		}
	}

	// Input options
	if opts.InputProtocol == "rtmp" {
		args = append(args,
			"-re", // Read input at native frame rate
			"-i", opts.InputURL,
		)
	} else if opts.InputProtocol == "srt" {
		args = append(args,
			"-f", "mpegts",
			"-i", fmt.Sprintf("srt://%s", opts.InputURL),
		)
	}

	// Video encoding options
	if opts.VideoCodec != "" {
		// Video codec selection
		codec := opts.VideoCodec
		if opts.HardwareAccel == "nvidia" && codec == "h264" {
			codec = "h264_nvenc"
		} else if opts.HardwareAccel == "vaapi" && codec == "h264" {
			codec = "h264_vaapi"
		} else if opts.HardwareAccel == "qsv" && codec == "h264" {
			codec = "h264_qsv"
		} else if opts.HardwareAccel == "videotoolbox" && codec == "h264" {
			codec = "h264_videotoolbox"
		}

		args = append(args,
			"-c:v", codec,
			"-b:v", fmt.Sprintf("%dk", opts.VideoBitrate),
			"-maxrate", fmt.Sprintf("%dk", int(float64(opts.VideoBitrate)*1.5)),
			"-bufsize", fmt.Sprintf("%dk", opts.VideoBitrate*2),
		)

		// Resolution
		if opts.VideoWidth > 0 && opts.VideoHeight > 0 {
			args = append(args, "-s", fmt.Sprintf("%dx%d", opts.VideoWidth, opts.VideoHeight))
		}

		// Framerate
		if opts.VideoFramerate > 0 {
			args = append(args, "-r", fmt.Sprintf("%d", opts.VideoFramerate))
		}

		// H.264 specific options
		if strings.Contains(codec, "h264") {
			args = append(args,
				"-profile:v", opts.VideoProfile,
				"-level", opts.VideoLevel,
			)

			// Low latency options
			if opts.LowLatency {
				args = append(args,
					"-preset", "ultrafast",
					"-tune", "zerolatency",
					"-x264-params", "keyint=60:min-keyint=60:no-scenecut",
				)
			} else {
				args = append(args,
					"-preset", "medium",
					"-x264-params", fmt.Sprintf("keyint=%d:min-keyint=%d:no-scenecut", 
						opts.VideoFramerate*2, opts.VideoFramerate*2),
				)
			}
		}

		// Pixel format
		args = append(args, "-pix_fmt", "yuv420p")
	}

	// Audio encoding options
	args = append(args,
		"-c:a", "aac",
		"-b:a", "128k",
		"-ar", "48000",
		"-ac", "2",
	)

	// Output format specific options
	if opts.OutputFormat == "hls" {
		args = append(args,
			"-f", "hls",
			"-hls_time", fmt.Sprintf("%d", opts.SegmentDuration),
			"-hls_list_size", fmt.Sprintf("%d", t.config.PlaylistSize),
			"-hls_flags", "delete_segments+append_list",
			"-hls_segment_filename", opts.SegmentPattern,
		)

		if opts.LowLatency {
			args = append(args,
				"-hls_flags", "delete_segments+append_list+program_date_time+independent_segments",
				"-hls_segment_type", "fmp4",
				"-hls_fmp4_init_filename", "init.mp4",
			)
		}

		args = append(args, opts.PlaylistPath)

	} else if opts.OutputFormat == "dash" {
		args = append(args,
			"-f", "dash",
			"-seg_duration", fmt.Sprintf("%d", opts.SegmentDuration),
			"-window_size", fmt.Sprintf("%d", t.config.PlaylistSize),
			"-use_template", "1",
			"-use_timeline", "1",
			"-init_seg_name", "init-$RepresentationID$.m4s",
			"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s",
		)

		if opts.LowLatency {
			args = append(args,
				"-ldash", "1",
				"-streaming", "1",
				"-adaptation_sets", "id=0,streams=v id=1,streams=a",
			)
		}

		args = append(args, opts.PlaylistPath)
	}

	return args
}

// Stop stops the transcoder and all active jobs
func (t *Transcoder) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, job := range t.activeJobs {
		if job.Cancel != nil {
			job.Cancel()
		}
	}
}

// GetActiveJobs returns the list of active transcoding jobs
func (t *Transcoder) GetActiveJobs() []TranscodeJob {
	t.mu.RLock()
	defer t.mu.RUnlock()

	jobs := make([]TranscodeJob, 0, len(t.activeJobs))
	for _, job := range t.activeJobs {
		jobs = append(jobs, *job)
	}
	return jobs
}

// transcodeLogger implements io.Writer for FFmpeg output logging
type transcodeLogger struct {
	jobID   string
	logType string
}

func (tl *transcodeLogger) Write(p []byte) (n int, err error) {
	// Log to structured logger in production
	// For now, just print
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("[%s/%s] %s\n", tl.jobID, tl.logType, line)
		}
	}
	return len(p), nil
}
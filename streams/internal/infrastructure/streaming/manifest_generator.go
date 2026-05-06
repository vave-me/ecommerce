package streaming

import (
	"fmt"
	"strings"
	"time"
)

// ManifestGenerator creates HLS and DASH manifests
type ManifestGenerator struct {
	config *StreamingConfig
}

// NewManifestGenerator creates a new manifest generator
func NewManifestGenerator(config *StreamingConfig) *ManifestGenerator {
	return &ManifestGenerator{
		config: config,
	}
}

// GenerateMasterPlaylist creates HLS master playlist
func (mg *ManifestGenerator) GenerateMasterPlaylist(streamID string, profiles []QualityProfile) string {
	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:7\n")

	// Add each quality variant
	for _, profile := range profiles {
		bandwidth := profile.Bitrate * 1000 // Convert to bps
		resolution := fmt.Sprintf("%dx%d", profile.Width, profile.Height)
		
		sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s,CODECS=\"avc1.%s,mp4a.40.2\",FRAME-RATE=%.3f\n",
			bandwidth,
			resolution,
			getAVCCodecString(profile.Profile, profile.Level),
			float64(profile.Framerate),
		))
		sb.WriteString(fmt.Sprintf("/hls/%s/%s/playlist.m3u8\n", streamID, profile.Name))
	}

	// Add subtitles if available
	sb.WriteString("#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"subs\",NAME=\"English\",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,LANGUAGE=\"en\",URI=\"/hls/" + streamID + "/subtitles/en.m3u8\"\n")

	return sb.String()
}

// GenerateDASHManifest creates DASH MPD manifest
func (mg *ManifestGenerator) GenerateDASHManifest(streamID string, profiles []QualityProfile) string {
	var sb strings.Builder

	// MPD header
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" `)
	sb.WriteString(`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" `)
	sb.WriteString(`xsi:schemaLocation="urn:mpeg:dash:schema:mpd:2011 http://standards.iso.org/ittf/PubliclyAvailableStandards/MPEG-DASH_schema_files/DASH-MPD.xsd" `)
	
	if mg.config.LowLatencyMode {
		sb.WriteString(`type="dynamic" `)
		sb.WriteString(`minimumUpdatePeriod="PT1S" `)
		sb.WriteString(`availabilityStartTime="` + time.Now().UTC().Format(time.RFC3339) + `" `)
		sb.WriteString(`publishTime="` + time.Now().UTC().Format(time.RFC3339) + `" `)
		sb.WriteString(`timeShiftBufferDepth="PT` + fmt.Sprintf("%d", mg.config.DVRWindowMinutes*60) + `S" `)
		sb.WriteString(`suggestedPresentationDelay="PT2S" `)
	} else {
		sb.WriteString(`type="dynamic" `)
		sb.WriteString(`minimumUpdatePeriod="PT` + fmt.Sprintf("%d", mg.config.SegmentDuration) + `S" `)
		sb.WriteString(`timeShiftBufferDepth="PT` + fmt.Sprintf("%d", mg.config.DVRWindowMinutes*60) + `S" `)
	}

	sb.WriteString(`minBufferTime="PT` + fmt.Sprintf("%d", mg.config.SegmentDuration) + `S" `)
	sb.WriteString(`profiles="urn:mpeg:dash:profile:isoff-live:2011,http://dashif.org/guidelines/dash-if-simple">`)
	sb.WriteString("\n")

	// Period
	sb.WriteString(`<Period id="0" start="PT0S">`)
	sb.WriteString("\n")

	// Video AdaptationSet
	sb.WriteString(`  <AdaptationSet id="0" contentType="video" segmentAlignment="true" bitstreamSwitching="true" frameRate="60/1" maxWidth="1920" maxHeight="1080" par="16:9">`)
	sb.WriteString("\n")

	for _, profile := range profiles {
		sb.WriteString(fmt.Sprintf(`    <Representation id="%s" mimeType="video/mp4" codecs="avc1.%s" bandwidth="%d" width="%d" height="%d" frameRate="%d/1">`,
			profile.Name,
			getAVCCodecString(profile.Profile, profile.Level),
			profile.Bitrate*1000,
			profile.Width,
			profile.Height,
			profile.Framerate,
		))
		sb.WriteString("\n")

		if mg.config.LowLatencyMode {
			sb.WriteString(`      <SegmentTemplate timescale="1000" media="/dash/` + streamID + `/` + profile.Name + `/chunk-$Number%05d$.m4s" initialization="/dash/` + streamID + `/` + profile.Name + `/init.mp4" duration="` + fmt.Sprintf("%d", mg.config.SegmentDuration*1000) + `">`)
		} else {
			sb.WriteString(`      <SegmentTemplate timescale="90000" media="/dash/` + streamID + `/` + profile.Name + `/segment_$Number$.m4s" initialization="/dash/` + streamID + `/` + profile.Name + `/init.mp4">`)
			sb.WriteString("\n")
			sb.WriteString(`        <SegmentTimeline>`)
			sb.WriteString("\n")
			sb.WriteString(`          <S t="0" d="` + fmt.Sprintf("%d", mg.config.SegmentDuration*90000) + `" r="-1"/>`)
			sb.WriteString("\n")
			sb.WriteString(`        </SegmentTimeline>`)
		}
		sb.WriteString("\n")
		sb.WriteString(`      </SegmentTemplate>`)
		sb.WriteString("\n")
		sb.WriteString(`    </Representation>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`  </AdaptationSet>`)
	sb.WriteString("\n")

	// Audio AdaptationSet
	sb.WriteString(`  <AdaptationSet id="1" contentType="audio" segmentAlignment="true" bitstreamSwitching="true" lang="en">`)
	sb.WriteString("\n")
	sb.WriteString(`    <Representation id="audio" mimeType="audio/mp4" codecs="mp4a.40.2" bandwidth="128000" audioSamplingRate="48000">`)
	sb.WriteString("\n")
	sb.WriteString(`      <AudioChannelConfiguration schemeIdUri="urn:mpeg:dash:23003:3:audio_channel_configuration:2011" value="2"/>`)
	sb.WriteString("\n")

	if mg.config.LowLatencyMode {
		sb.WriteString(`      <SegmentTemplate timescale="48000" media="/dash/` + streamID + `/audio/chunk-$Number%05d$.m4s" initialization="/dash/` + streamID + `/audio/init.mp4" duration="` + fmt.Sprintf("%d", mg.config.SegmentDuration*48000) + `">`)
	} else {
		sb.WriteString(`      <SegmentTemplate timescale="48000" media="/dash/` + streamID + `/audio/segment_$Number$.m4s" initialization="/dash/` + streamID + `/audio/init.mp4">`)
		sb.WriteString("\n")
		sb.WriteString(`        <SegmentTimeline>`)
		sb.WriteString("\n")
		sb.WriteString(`          <S t="0" d="` + fmt.Sprintf("%d", mg.config.SegmentDuration*48000) + `" r="-1"/>`)
		sb.WriteString("\n")
		sb.WriteString(`        </SegmentTimeline>`)
	}
	sb.WriteString("\n")
	sb.WriteString(`      </SegmentTemplate>`)
	sb.WriteString("\n")
	sb.WriteString(`    </Representation>`)
	sb.WriteString("\n")
	sb.WriteString(`  </AdaptationSet>`)
	sb.WriteString("\n")

	// End Period and MPD
	sb.WriteString(`</Period>`)
	sb.WriteString("\n")
	sb.WriteString(`</MPD>`)

	return sb.String()
}

// GenerateHLSVariantPlaylist creates quality-specific HLS playlist
func (mg *ManifestGenerator) GenerateHLSVariantPlaylist(segments []string, targetDuration int) string {
	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:7\n")
	sb.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", targetDuration))
	sb.WriteString("#EXT-X-PLAYLIST-TYPE:EVENT\n")

	if mg.config.LowLatencyMode {
		sb.WriteString("#EXT-X-SERVER-CONTROL:CAN-SKIP-UNTIL=12.0\n")
		sb.WriteString("#EXT-X-PART-INF:PART-TARGET=0.33334\n")
	}

	mediaSequence := 0
	sb.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", mediaSequence))

	// Add segments
	for _, segment := range segments {
		sb.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n", float64(mg.config.SegmentDuration)))
		sb.WriteString(segment + "\n")
	}

	return sb.String()
}

// GenerateWebVTTPlaylist creates subtitle playlist
func (mg *ManifestGenerator) GenerateWebVTTPlaylist(streamID, language string) string {
	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:7\n")
	sb.WriteString("#EXT-X-TARGETDURATION:10\n")
	sb.WriteString("#EXT-X-PLAYLIST-TYPE:EVENT\n")
	sb.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")

	// Add subtitle segments
	// In production, this would be generated from actual subtitle data
	sb.WriteString("#EXTINF:10.0,\n")
	sb.WriteString(fmt.Sprintf("/hls/%s/subtitles/%s_0.vtt\n", streamID, language))

	return sb.String()
}

// getAVCCodecString converts H.264 profile and level to codec string
func getAVCCodecString(profile, level string) string {
	profileMap := map[string]string{
		"baseline": "42",
		"main":     "4D",
		"high":     "64",
	}

	levelMap := map[string]string{
		"3.0": "1E",
		"3.1": "1F",
		"4.0": "28",
		"4.1": "29",
		"4.2": "2A",
		"5.0": "32",
		"5.1": "33",
		"5.2": "34",
	}

	profileHex := profileMap[profile]
	levelHex := levelMap[level]

	if profileHex == "" {
		profileHex = "64" // Default to high
	}
	if levelHex == "" {
		levelHex = "28" // Default to 4.0
	}

	// Format: avc1.PPCCLL where PP=profile, CC=constraints, LL=level
	return fmt.Sprintf("%s00%s", profileHex, levelHex)
}
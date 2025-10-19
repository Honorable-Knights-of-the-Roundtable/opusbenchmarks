package main

import (
	"fmt"
	"hmcalister/opusbenchmark/encdec"
	"math/rand"
	"testing"
	"time"
)

const (
	ENCDEC_TYPE          encdec.EncDecType = encdec.EncDecTypeJJ11hh
	BUFFER_SAFETY_FACTOR int               = 16
)

var (
	sampleRates []int = []int{
		8000,
		12000,
		16000,
		24000,
		48000,
	}
	channels       []int                      = []int{1, 2}
	frameDurations []encdec.OPUSFrameDuration = []encdec.OPUSFrameDuration{
		encdec.OPUS_FRAME_DURATION_20_MS,
		encdec.OPUS_FRAME_DURATION_40_MS,
		encdec.OPUS_FRAME_DURATION_60_MS,
	}
)

func encodeAudio(b *testing.B, audio []encdec.PCMFrame, sampleRate int, numChannels int, frameDuration encdec.OPUSFrameDuration) {
	encdec, err := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels, frameDuration, BUFFER_SAFETY_FACTOR)
	if err != nil {
		b.Errorf("error when creating encoder decoder %s", ENCDEC_TYPE)
	}

	for b.Loop() {
		for _, frame := range audio {
			encdec.Encode(frame)
			// _, err := encdec.Encode(frame)
			// if err != nil {
			// 	b.Errorf("error while encoding frame %v", err)
			// }
		}
	}
}

func decodeAudio(b *testing.B, audio []encdec.EncodedFrame, sampleRate int, numChannels int, frameDuration encdec.OPUSFrameDuration) {
	encdec, err := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels, frameDuration, BUFFER_SAFETY_FACTOR)
	if err != nil {
		b.Errorf("error when creating encoder decoder %s", ENCDEC_TYPE)
	}

	for b.Loop() {
		for _, frame := range audio {
			encdec.Decode(frame)
			// _, err := encdec.Decode(frame)
			// if err != nil {
			// 	b.Errorf("error while decoding frame %v", err)
			// }
		}
	}
}

func BenchmarkSilentTrack(b *testing.B) {
	trackDuration := 10 * time.Second

	for _, sampleRate := range sampleRates {
		for _, numChannels := range channels {
			for _, frameDuration := range frameDurations {

				// Make a silent track
				audio := make([]encdec.PCMFrame, trackDuration/time.Duration(frameDuration))
				for i := range audio {
					audio[i] = make(encdec.PCMFrame, (sampleRate * numChannels * int(frameDuration) / int(time.Second)))
				}

				// Encode the track with benchmarks
				b.Run(fmt.Sprintf("Encode Silent Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Silence")
					b.Attr("TaskType", "Encoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					encodeAudio(b, audio, sampleRate, numChannels, frameDuration)
				})

				// Encode the silent track
				ed, _ := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels, frameDuration, BUFFER_SAFETY_FACTOR)
				encodedAudio := make([]encdec.EncodedFrame, len(audio))
				encodedFrameIndex := 0
				for rawFrameIndex, frame := range audio {
					encodedFrames, err := ed.Encode(frame)
					if err != nil {
						b.Logf("error while encoding frame, frameIndex %v frameLength %v", rawFrameIndex, len(frame))
						continue
					}
					for _, encodedFrame := range encodedFrames {
						encodedAudio[encodedFrameIndex] = encodedFrame
						encodedFrameIndex += 1
					}
				}

				// Decode the track with benchmarks
				b.Run(fmt.Sprintf("Decode Silent Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Silence")
					b.Attr("TaskType", "Decoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					decodeAudio(b, encodedAudio, sampleRate, numChannels, frameDuration)
				})

			}
		}
	}
}

func BenchmarkRandomTrack(b *testing.B) {
	trackDuration := 10 * time.Second

	for _, sampleRate := range sampleRates {
		for _, numChannels := range channels {
			for _, frameDuration := range frameDurations {

				// Make a random track
				audio := make([]encdec.PCMFrame, trackDuration/time.Duration(frameDuration))
				for i := range audio {
					audio[i] = make(encdec.PCMFrame, (sampleRate * numChannels * int(frameDuration) / int(time.Second)))
					for sampleIndex := range audio[i] {
						audio[i][sampleIndex] = rand.Float32()*2 - 1
					}
				}

				// Encode the track with benchmarks
				b.Run(fmt.Sprintf("Encode Random Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Random")
					b.Attr("TaskType", "Encoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					encodeAudio(b, audio, sampleRate, numChannels, frameDuration)
				})

				// Encode the random track
				ed, _ := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels, frameDuration, BUFFER_SAFETY_FACTOR)
				encodedAudio := make([]encdec.EncodedFrame, len(audio))
				encodedFrameIndex := 0
				for rawFrameIndex, frame := range audio {
					encodedFrames, err := ed.Encode(frame)
					if err != nil {
						b.Logf("error while encoding frame, frameIndex %v frameLength %v", rawFrameIndex, len(frame))
						continue
					}
					for _, encodedFrame := range encodedFrames {
						encodedAudio[encodedFrameIndex] = encodedFrame
						encodedFrameIndex += 1
					}
				}

				// Decode the track with benchmarks
				b.Run(fmt.Sprintf("Decode Random Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Random")
					b.Attr("TaskType", "Decoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					decodeAudio(b, encodedAudio, sampleRate, numChannels, frameDuration)
				})
			}
		}
	}
}

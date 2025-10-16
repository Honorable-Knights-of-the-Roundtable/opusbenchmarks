package main

import (
	"fmt"
	"hmcalister/opusbenchmark/encdec"
	"log/slog"
	"math/rand"
	"testing"
	"time"
)

const (
	ENCDEC_TYPE encdec.EncDecType = encdec.EncDecTypeHraban
)

var (
	sampleRates []int = []int{
		8000,
		12000,
		16000,
		24000,
		48000,
	}
	channels       []int           = []int{1, 2}
	frameDurations []time.Duration = []time.Duration{
		time.Microsecond * 2500,
		time.Millisecond * 5,
		time.Millisecond * 10,
		time.Millisecond * 20,
		time.Millisecond * 40,
	}
)

func encodeAudio(b *testing.B, audio []encdec.PCMFrame, sampleRate int, numChannels int) {
	encdec, err := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels)
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

func decodeAudio(b *testing.B, audio []encdec.EncodedFrame, sampleRate int, numChannels int) {
	encdec, err := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels)
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

func BenchmarkEncodeSilence(b *testing.B) {
	trackDuration := 10 * time.Second

	for _, sampleRate := range sampleRates {
		for _, numChannels := range channels {
			for _, frameDuration := range frameDurations {

				// Make a silent track
				audio := make([]encdec.PCMFrame, trackDuration/frameDuration)
				for i := range audio {
					audio[i] = make(encdec.PCMFrame, (sampleRate / int(trackDuration/time.Second) * numChannels))
				}

				// Encode the track with benchmarks
				b.Run(fmt.Sprintf("Encode Silent Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Silence")
					b.Attr("TaskType", "Encoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					encodeAudio(b, audio, sampleRate, numChannels)
				})

				ed, err := encdec.NewOpusEncoderDecoder(ENCDEC_TYPE, sampleRate, numChannels)
				if err != nil {
					b.Logf("error while creating OPUS encoder/decoder %v", err)
					continue
				}
				// Encode the silent track
				encodedAudio := make([]encdec.EncodedFrame, len(audio))
				for frameIndex, frame := range audio {
					encodedFrame, err := ed.Encode(frame)
					if err != nil {
						b.Logf("error while encoding frame, frameIndex %v frameLength %v", frameIndex, len(frame))
						continue
					}
					encodedAudio[frameIndex] = encodedFrame
				}

				// Decode the track with benchmarks
				b.Run(fmt.Sprintf("Decode Silent Audio: Sample Rate %d, Channels %d, Frame Duration %v", sampleRate, numChannels, frameDuration), func(b *testing.B) {
					b.Attr("AudioType", "Silence")
					b.Attr("TaskType", "Decoding")
					b.Attr("sampleRate", fmt.Sprint(sampleRate))
					b.Attr("numChannels", fmt.Sprint(numChannels))
					b.Attr("frameDuration", fmt.Sprint(frameDuration))
					decodeAudio(b, encodedAudio, sampleRate, numChannels)
				})

			}
		}
	}
}


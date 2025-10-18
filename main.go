package main

import (
	"hmcalister/opusbenchmark/encdec"
	"log/slog"
	"os"
	"time"
)

func init() {
	loggerOptions := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, loggerOptions)))
}

func main() {
	trackDuration := 3 * time.Second
	sampleRate := 48000
	numChannels := 2
	frameDuration := encdec.OPUS_FRAME_DURATION_20_MS
	bufferSafetyFactor := 16

	audio := make([]encdec.PCMFrame, int(trackDuration/time.Duration(frameDuration)))
	for i := range audio {
		audio[i] = make(encdec.PCMFrame, int(frameDuration)*sampleRate*numChannels/int(time.Second))
	}
	slog.Debug("audio", "num frames", len(audio), "first frame length", len(audio[0]), "last frame len", len(audio[len(audio)-1]))

	// for i := range 10 {
	// 	go func() {
	// 		ed, _ := encdec.NewOpusEncoderDecoder(encdec.EncDecTypeHraban, sampleRate, numChannels)

	// 		encodedAudio := make([]encdec.EncodedFrame, len(audio))
	// 		for frameIndex, frame := range audio {
	// 			time.Sleep(frameDuration)
	// 			encodedFrame, _ := ed.Encode(frame)
	// 			encodedAudio[frameIndex] = encodedFrame
	// 		}
	// 		slog.Debug("finished encoding", "goroutine index", i)
	// 	}()
	// }
	// select {}

	ed, _ := encdec.NewOpusEncoderDecoder(encdec.EncDecTypeJJ11h, sampleRate, numChannels, frameDuration, bufferSafetyFactor)
	encodedAudio := make([]encdec.EncodedFrame, len(audio))
	slog.Debug("encoding audio")
	for frameIndex, frame := range audio {
		encodedFrame, err := ed.Encode(frame)
		if err != nil {
			slog.Error("error while encoding frame", "frameIndex", frameIndex, "frameLength", len(frame), "err", err)
			break
		}
		for _, frame := range encodedFrame {
			encodedAudio[frameIndex] = frame
		}
	}

	for i := range 2 {
		go func() {
			decodedAudio := make([]encdec.PCMFrame, len(audio))
			for frameIndex, frame := range encodedAudio {
				time.Sleep(time.Duration(frameDuration))
				decodedFrame, _ := ed.Decode(frame)
				decodedAudio[frameIndex] = decodedFrame
				// slog.Debug("decoded frame", "frame index", frameIndex, "frame len", len(frame), "num samples", len(decodedFrame))
			}
			slog.Debug("finished decoding", "goroutine index", i)
		}()
	}
	select {}
}

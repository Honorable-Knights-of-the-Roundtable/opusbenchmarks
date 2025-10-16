package encdec

import (
	"errors"
	"fmt"

	"github.com/hraban/opus"
)

type HrabanOpusEncoderDecoder struct {
	sampleRate  int
	numChannels int

	encoder       *opus.Encoder
	encodingFrame EncodedFrame
	decoder       *opus.Decoder
	decodedFrame  PCMFrame
}

func newHrabanEncoderDecoder(sampleRate int, numChannels int) (HrabanOpusEncoderDecoder, error) {
	encoder, errEnc := opus.NewEncoder(sampleRate, numChannels, opus.Application(opus.AppVoIP))
	decoder, errDec := opus.NewDecoder(sampleRate, numChannels)
	if err := errors.Join(errEnc, errDec); err != nil {
		return HrabanOpusEncoderDecoder{}, err
	}

	return HrabanOpusEncoderDecoder{
		sampleRate:    sampleRate,
		numChannels:   numChannels,
		encoder:       encoder,
		encodingFrame: make(EncodedFrame, bufferSize),
		decoder:       decoder,
		decodedFrame:  make(PCMFrame, bufferSize),
	}, nil
}

func (encdec HrabanOpusEncoderDecoder) String() string {
	return fmt.Sprintf("Type: hraban, Sample Rate: %d, Num Channels: %d", encdec.sampleRate, encdec.numChannels)
}

func (encdec HrabanOpusEncoderDecoder) Encode(pcmData PCMFrame) (EncodedFrame, error) {
	encodedBytes, err := encdec.encoder.EncodeFloat32(pcmData, encdec.encodingFrame)
	if err != nil {
		return nil, err
	}
	return encdec.encodingFrame[:encodedBytes], nil
}

func (encdec HrabanOpusEncoderDecoder) Decode(encodedData EncodedFrame) (PCMFrame, error) {
	decodedBytes, err := encdec.decoder.DecodeFloat32(encodedData, encdec.decodedFrame)
	if err != nil {
		return nil, err
	}
	return encdec.decodedFrame[:decodedBytes*encdec.numChannels], nil
}

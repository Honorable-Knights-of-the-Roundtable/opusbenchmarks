package encdec

import (
	"errors"
	"fmt"

	"github.com/jj11hh/opus"
)

type JJ11hhOpusEncoderDecoder struct {
	sampleRate  int
	numChannels int

	encoder       *opus.Encoder
	encodingFrame EncodedFrame
	decoder       *opus.Decoder
	decodedFrame  PCMFrame
}

func newJJ11hhEncoderDecoder(sampleRate int, numChannels int) (JJ11hhOpusEncoderDecoder, error) {
	encoder, errEnc := opus.NewEncoder(sampleRate, numChannels, opus.Application(opus.AppVoIP))
	decoder, errDec := opus.NewDecoder(sampleRate, numChannels)
	if err := errors.Join(errEnc, errDec); err != nil {
		return JJ11hhOpusEncoderDecoder{}, err
	}

	return JJ11hhOpusEncoderDecoder{
		sampleRate:    sampleRate,
		numChannels:   numChannels,
		encoder:       encoder,
		encodingFrame: make(EncodedFrame, bufferSize),
		decoder:       decoder,
		decodedFrame:  make(PCMFrame, bufferSize),
	}, nil
}

func (encdec JJ11hhOpusEncoderDecoder) String() string {
	return fmt.Sprintf("Type: JJ11hh, Sample Rate: %d, Num Channels: %d", encdec.sampleRate, encdec.numChannels)
}

func (encdec JJ11hhOpusEncoderDecoder) Encode(pcmData PCMFrame) (EncodedFrame, error) {
	encodedBytes, err := encdec.encoder.EncodeFloat32(pcmData, encdec.encodingFrame)
	if err != nil {
		return nil, err
	}
	return encdec.encodingFrame[:encodedBytes], nil
}

func (encdec JJ11hhOpusEncoderDecoder) Decode(encodedData EncodedFrame) (PCMFrame, error) {
	decodedBytes, err := encdec.decoder.DecodeFloat32(encodedData, encdec.decodedFrame)
	if err != nil {
		return nil, err
	}
	return encdec.decodedFrame[:decodedBytes*encdec.numChannels], nil
}

package encdec

import "fmt"

type OpusEncoderDecoder interface {
	String() string
	Encode(pcmData PCMFrame) ([]EncodedFrame, error)
	Decode(encodedData EncodedFrame) (PCMFrame, error)
}

func NewOpusEncoderDecoder(
	encdecType EncDecType,
	sampleRate int,
	numChannels int,
	frameDuration OPUSFrameDuration,
	bufferSafetyFactor int,
) (OpusEncoderDecoder, error) {
	switch encdecType {
	case EncDecTypeHraban:
		return newHrabanEncoderDecoder(sampleRate, numChannels, frameDuration, bufferSafetyFactor)
	case EncDecTypeJJ11hh:
		return newJJ11hhEncoderDecoder(sampleRate, numChannels, frameDuration, bufferSafetyFactor)
	default:
		return nil, fmt.Errorf("no encoder decoder associated with type %s", encdecType)
	}
}

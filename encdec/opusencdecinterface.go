package encdec

import "fmt"

const (
	bufferSize = 16384
)

type OpusEncoderDecoder interface {
	String() string
	Encode(pcmData PCMFrame) (EncodedFrame, error)
	Decode(encodedData EncodedFrame) (PCMFrame, error)
}

func NewOpusEncoderDecoder(encdecType EncDecType, sampleRate int, numChannels int) (OpusEncoderDecoder, error) {
	switch encdecType {
	case EncDecTypeHraban:
		return newHrabanEncoderDecoder(sampleRate, numChannels)
	case EncDecTypeJJ11h:
		return newHrabanEncoderDecoder(sampleRate, numChannels)
	default:
		return nil, fmt.Errorf("no encoder decoder associated with type %s", encdecType)
	}
}

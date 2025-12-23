package img

type DecodeOption func(*Decoder)

func WithMeta() DecodeOption {
	return func(dec *Decoder) {
		dec.withMeta = true
	}
}

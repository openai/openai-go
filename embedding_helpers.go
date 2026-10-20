package openai

import (
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func (r *Embedding) unmarshalEmbedding(data []byte) error {
	if err := apijson.UnmarshalRoot(data, r); err != nil {
		return err
	}

	raw := r.JSON.Embedding.Raw()
	if len(raw) == 0 || raw[0] != '"' {
		return nil
	}

	// A JSON string unmarshals into []byte as base64; the bytes contain little-endian float32 values.
	var decoded []byte
	if json.Unmarshal([]byte(raw), &decoded) == nil && len(decoded)%4 == 0 {
		r.Embedding = make([]float64, len(decoded)/4)
		for i := range r.Embedding {
			r.Embedding[i] = float64(math.Float32frombits(binary.LittleEndian.Uint32(decoded[i*4:])))
		}
		r.JSON.Embedding = respjson.NewField(raw)
	}
	return nil
}

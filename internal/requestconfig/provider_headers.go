package requestconfig

import (
	"net/http"
	"slices"
)

type headerOperation struct {
	value string
	layer *optionLayerIdentity
	kind  headerOperationKind
}

type headerOperationKind byte

const (
	headerSet headerOperationKind = iota
	headerAdd
	headerDelete
)

type optionHeader struct {
	defaults   []string
	operations []headerOperation
}

// Record options before applying them so a provider switch can restore SDK
// defaults and replay only options belonging to the new provider's layers.
func (cfg *RequestConfig) recordOptionHeader(key, value string, kind headerOperationKind) {
	state := &cfg.authentication
	if state.optionHeaders == nil {
		state.optionHeaders = make(map[string]optionHeader)
	}
	key = http.CanonicalHeaderKey(key)
	header, ok := state.optionHeaders[key]
	if !ok {
		header.defaults = slices.Clone(cfg.Request.Header.Values(key))
	}
	header.operations = append(header.operations, headerOperation{value, state.currentLayer, kind})
	state.optionHeaders[key] = header
}

// ClearInheritedProviderHeaders removes headers from inherited option layers
// when selecting a different provider. The switching layer and its enclosing
// layers remain: generated operation headers are applied in the outer request
// layer before inherited client/service options. Call before selecting auth.
func (cfg *RequestConfig) ClearInheritedProviderHeaders(provider string) {
	state := &cfg.authentication
	if selected := state.selectedProvider; selected != nil && selected.provider == provider {
		return
	}
	for key, header := range state.optionHeaders {
		cfg.Request.Header.Del(key)
		for _, value := range header.defaults {
			cfg.Request.Header.Add(key, value)
		}
		kept := make([]headerOperation, 0, len(header.operations))
		for _, op := range header.operations {
			if !enclosingHeaderLayer(op.layer, state.currentLayer) {
				continue
			}
			switch op.kind {
			case headerSet:
				cfg.Request.Header.Set(key, op.value)
			case headerAdd:
				cfg.Request.Header.Add(key, op.value)
			case headerDelete:
				cfg.Request.Header.Del(key)
			}
			kept = append(kept, op)
		}
		if len(kept) == 0 {
			delete(state.optionHeaders, key)
		} else {
			header.operations = kept
			state.optionHeaders[key] = header
		}
		switch key {
		case "Openai-Organization":
			cfg.Organization = cfg.Request.Header.Get(key)
		case "Openai-Project":
			cfg.Project = cfg.Request.Header.Get(key)
		}
	}
}

func enclosingHeaderLayer(layer, current *optionLayerIdentity) bool {
	for {
		if layer == current {
			return true
		}
		if current == nil {
			return false
		}
		current = current.parent
	}
}

func inheritedOptionHeaders(headers map[string]optionHeader) map[string]optionHeader {
	if headers == nil {
		return nil
	}
	inherited := make(map[string]optionHeader, len(headers))
	layer := new(optionLayerIdentity)
	for key, header := range headers {
		header.operations = slices.Clone(header.operations)
		for i := range header.operations {
			header.operations[i].layer = layer
		}
		inherited[key] = header
	}
	return inherited
}

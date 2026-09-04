// EDIT(begin): add custom options for JSON encoding
package json

type Option func(*encOpts)

// Every time a sub-type of [json.Marshaler] is encountered,
// skip a redundant and costly compaction step, trust it to self-compact.
//
// This is a divergence from the standard library behavior, and is only guaranteed
// safe with SDK types.
func WithSkipCompaction(b bool) Option {
	return func(eos *encOpts) {
		eos.skipCompaction = true
	}
}

func (eos encOpts) apply(opts ...Option) encOpts {
	for _, opt := range opts {
		opt(&eos)
	}
	return eos
}

// EDIT(end)

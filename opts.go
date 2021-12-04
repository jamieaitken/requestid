package requestid

// WithTracerFunc allows for using a custom function instead.
func WithTracerFunc(tracerFunc TracerFunc) Option {
	return func(tracer *Tracer) {
		tracer.fn = tracerFunc
	}
}

// WithTracerKey allows for a custom key to be used.
func WithTracerKey(tracerKey Key) Option {
	return func(tracer *Tracer) {
		tracer.key = tracerKey
	}
}

// WithIDGenerator allows for a custom id generator to be used.
func WithIDGenerator(generatorFunc IDGeneratorFunc) Option {
	return func(tracer *Tracer) {
		tracer.idGeneratorFunc = generatorFunc
	}
}

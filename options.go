package magic

type options struct {
	converters []Converter
	mapping    map[string]string
}

// WithMapping add mapping to options
func WithMapping(mapping map[string]string) func(o *options) {
	return func(o *options) {
		o.mapping = mapping
	}
}

// WithConverters add converters to options
func WithConverters(converters ...Converter) func(o *options) {
	return func(o *options) {
		o.converters = converters
	}
}

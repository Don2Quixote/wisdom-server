package app

// appConfig is struct for parsing ENV configuration.
type appConfig struct {
	// Port is port to launch TCP server.
	Port int `config:"PORT,required"`
	// ComplexityFactor shows how fast does complexity grow.
	ComplexityFactor float64 `config:"COMPLEXITY_FACTOR,required"`
	// MaxComplexity limits the maximum complexity level.
	MaxComplexity int `config:"MAX_COMPLEXITY,required"`
	// ComplexityDuration is value that shows how much time should pass
	// to restore complexity points after grow.
	ComplexityDurationSeconds int `config:"COMPLEXITY_DURATION_SECONDS,required"`
}

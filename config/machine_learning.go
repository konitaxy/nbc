package config

type MachineLearning struct {
	ImageDir          string  `mapstructure:"image-dir" json:"image-dir" yaml:"image-dir"`
	HistogramFilename string  `mapstructure:"histogram-filename" json:"histogram-filename" yaml:"histogram-filename"`
	PHashFilename     string  `mapstructure:"phash-filename" json:"phash-filename" yaml:"phash-filename"`
	PHASHDistance     int     `mapstructure:"phash-distance" json:"phash-distance" yaml:"phash-distance"`
	HistogramSimilar  float64 `mapstructure:"histogram-similar" json:"histogram-similar" yaml:"histogram-similar"`
	MLSimilar         float64 `mapstructure:"ml-similar" json:"ml-similar" yaml:"ml-similar"`
	MLUrl             string  `mapstructure:"ml-url" json:"ml-url" yaml:"ml-url"`
	Mode              string  `mapstructure:"mode" json:"mode" yaml:"mode"`
}

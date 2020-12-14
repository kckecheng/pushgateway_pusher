package types

// MAXINTERVAL maximum of interval (in minute) between 2 x records
const MAXINTERVAL = 5

// MAXERROR maximum num. continuous error can be tolerated
const MAXERROR = 5

// DELAY start to count for errors after DELAY seconds
const DELAY = 30

// ERREXIT error out
const ERREXIT = 1

// Field metric fields, since it is used by parser, donot delete the tags
type Field struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Metric raw metrics
type Metric struct {
	Field
	Value  float64
	Labels []string
}

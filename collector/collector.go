package collector

import (
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// PusherCollector implement the Collector interface
type PusherCollector struct {
	Fields      map[string]*prometheus.Desc
	Application string
	Values      map[string]float64
	Mutex       sync.Mutex
}

// New init collector
func New(fields []map[string]string, application string, values map[string]float64) (*PusherCollector, error) {
	fieldMap := make(map[string]*prometheus.Desc)
	for _, field := range fields {
		name, ok := field["name"]
		if !ok {
			return nil, errors.New("Each field requires a name")
		}
		help, ok := field["help"]
		if !ok {
			return nil, errors.New("Each field requires a help")
		}
		desc := prometheus.NewDesc(name, help, []string{"application"}, nil)
		fieldMap[name] = desc
	}

	pc := PusherCollector{
		Fields:      fieldMap,
		Application: application,
		Values:      values,
		Mutex:       sync.Mutex{},
	}
	return &pc, nil
}

// Describe implement interface
func (pc *PusherCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range pc.Fields {
		ch <- v
	}
}

// Collect implement interface
func (pc *PusherCollector) Collect(ch chan<- prometheus.Metric) {
	pc.Mutex.Lock()
	defer pc.Mutex.Unlock()

	for k, v := range pc.Values {
		ch <- prometheus.MustNewConstMetric(
			pc.Fields[k],
			prometheus.GaugeValue,
			v,
			pc.Application,
		)
	}
}

// UpdateValues update metric values
func (pc *PusherCollector) UpdateValues(values map[string]float64) {
	pc.Mutex.Lock()
	defer pc.Mutex.Unlock()
	pc.Values = values
}

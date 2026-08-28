package prometheusutil

import (
	"net/http"

	"github.com/fasionchan/goutils/stl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collectors []prometheus.Collector

func NewCollectors(collectors ...prometheus.Collector) Collectors {
	return collectors
}

func (collectors Collectors) Append(more ...prometheus.Collector) Collectors {
	return append(collectors, more...)
}

func (collectors Collectors) Concat(others ...Collectors) Collectors {
	return stl.ConcatSlicesTo(collectors, others...)
}

func (collectors Collectors) BuildRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors...)
	return registry
}

func (collectors Collectors) BuildHandler() http.Handler {
	return promhttp.HandlerFor(collectors.BuildRegistry(), promhttp.HandlerOpts{})
}

type CollectorGroup = Collectors

func NewCollectorGroup(collectors ...prometheus.Collector) CollectorGroup {
	return collectors
}

type CollectorGroups []CollectorGroup

func NewCollectorGroups(groups ...CollectorGroup) CollectorGroups {
	return groups
}

func (groups CollectorGroups) BuildRegistry() *prometheus.Registry {
	return groups.Stack().BuildRegistry()
}

func (groups CollectorGroups) BuildHandler() http.Handler {
	return groups.Stack().BuildHandler()
}

func (groups CollectorGroups) Stack() CollectorGroup {
	return stl.ConcatSlices(groups...)
}

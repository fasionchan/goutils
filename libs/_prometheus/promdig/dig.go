package promdig

import (
	"github.com/fasionchan/goutils/libs/_prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/dig"
)

type CollectorGroupsIn struct {
	dig.In
	Groups _prometheus.CollectorGroups `name:"default"`
}

func (in CollectorGroupsIn) Collectors() _prometheus.Collectors {
	return in.Groups.Stack()
}

type CollectorGroupsOut struct {
	dig.Out
	Groups _prometheus.CollectorGroups `name:"default"`
}

type collectorGroupSliceIn struct {
	dig.In
	Groups []_prometheus.CollectorGroup `group:"default"`
}

func (in collectorGroupSliceIn) GroupsOut() CollectorGroupsOut {
	return CollectorGroupsOut{
		Groups: _prometheus.CollectorGroups(in.Groups),
	}
}

type CollectorGroupOut struct {
	dig.Out
	Group _prometheus.CollectorGroup `group:"default"`
}

func (out *CollectorGroupOut) Append(collectors ...prometheus.Collector) {
	out.Group = out.Group.Append(collectors...)
}

func (out *CollectorGroupOut) Concat(groups ..._prometheus.CollectorGroup) {
	out.Group = out.Group.Concat(groups...)
}

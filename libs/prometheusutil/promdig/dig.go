package promdig

import (
	"github.com/fasionchan/goutils/libs/prometheusutil"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/dig"
)

type CollectorGroupsIn struct {
	dig.In
	Groups prometheusutil.CollectorGroups `name:"default"`
}

func (in CollectorGroupsIn) Collectors() prometheusutil.Collectors {
	return in.Groups.Stack()
}

type CollectorGroupsOut struct {
	dig.Out
	Groups prometheusutil.CollectorGroups `name:"default"`
}

type collectorGroupSliceIn struct {
	dig.In
	Groups []prometheusutil.CollectorGroup `group:"default"`
}

func (in collectorGroupSliceIn) GroupsOut() CollectorGroupsOut {
	return CollectorGroupsOut{
		Groups: prometheusutil.CollectorGroups(in.Groups),
	}
}

type CollectorGroupOut struct {
	dig.Out
	Group prometheusutil.CollectorGroup `group:"default"`
}

func (out *CollectorGroupOut) Append(collectors ...prometheus.Collector) {
	out.Group = out.Group.Append(collectors...)
}

func (out *CollectorGroupOut) Concat(groups ...prometheusutil.CollectorGroup) {
	out.Group = out.Group.Concat(groups...)
}

func GetProviderFuncs() []any {
	return []any{
		collectorGroupSliceIn.GroupsOut,
		CollectorGroupsIn.Collectors,
		prometheusutil.Collectors.BuildRegistry,
	}
}
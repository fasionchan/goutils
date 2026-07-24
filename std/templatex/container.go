package templatex

import "github.com/fasionchan/goutils/stl"

// ✅ DataContainer is a container for data.
type DataContainer struct {
	d any
}

func NewDataContainer() *DataContainer {
	return &DataContainer{}
}

func (c *DataContainer) Get() any {
	return c.d
}

func (c *DataContainer) Clear() *DataContainer {
	c.d = nil
	return c
}

func (c *DataContainer) Set(d any) *DataContainer {
	c.d = d
	return c
}

// DataContainers is a slice of DataContainers.
type DataContainers []*DataContainer

func (containers DataContainers) Append(others ...*DataContainer) DataContainers {
	return append(containers, others...)
}

// DataContainerMappingByString is a mapping of DataContainers by string.
type DataContainerMappingByString map[string]*DataContainer

func (m DataContainerMappingByString) Clear() DataContainerMappingByString {
	for _, container := range m {
		container.Clear()
	}
	return m
}

func (m DataContainerMappingByString) GetAllDatas() map[string]any {
	return stl.MapMap[map[string]any](m, func(key string, container *DataContainer, _ DataContainerMappingByString) (string, any) {
		return key, container.Get()
	})
}

func (m DataContainerMappingByString) GetData(key string) any {
	if m == nil {
		return nil
	}

	container, ok := m[key]
	if !ok {
		return nil
	}

	return container.Get()
}

// SetFuncs returns a TemplateFuncMap with the Set functions for the DataContainers.
func (m DataContainerMappingByString) SetFuncs() TemplateFuncMap {
	return stl.MapMap[TemplateFuncMap](m, func(name string, container *DataContainer, _ DataContainerMappingByString) (string, any) {
		return name, container.Set
	})
}

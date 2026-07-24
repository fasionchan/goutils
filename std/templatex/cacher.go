package templatex

import (
	"context"
	"time"
)

type SmartTemplateFetcher interface {
	FindTemplateByKeyAndFuncMap(context.Context, string, TemplateFuncMap) (*SmartTemplate, error)
}

type SmartTemplateCacher struct {
	fetcher       SmartTemplateFetcher
	cacheDuration time.Duration

	pool       SmartTemplateMappingByString
	cachedTime map[string]time.Time
	funcMap    TemplateFuncMap
}

func NewSmartTemplateCacher(fetcher SmartTemplateFetcher, cacheDuration time.Duration) *SmartTemplateCacher {
	return &SmartTemplateCacher{
		fetcher:       fetcher,
		cacheDuration: cacheDuration,

		pool:       NewSmartTemplateMappingByString(),
		cachedTime: make(map[string]time.Time),
	}
}

func (cacher *SmartTemplateCacher) WithFuncMap(funcMap TemplateFuncMap) *SmartTemplateCacher {
	if cacher == nil {
		return nil
	}

	cacher.funcMap = funcMap
	return cacher
}

func (cacher *SmartTemplateCacher) FetchTemplateByKey(ctx context.Context, key string) (*SmartTemplate, error) {
	tpl, err := cacher.fetcher.FindTemplateByKeyAndFuncMap(ctx, key, cacher.funcMap)
	if err != nil {
		return nil, err
	}

	if tpl != nil {
		cacher.pool[key] = tpl

		if cacher.cacheDuration > 0 {
			cacher.cachedTime[key] = time.Now()
		}
	}

	return tpl, err
}

func (cacher *SmartTemplateCacher) GetTemplateByKey(ctx context.Context, key string, cacheOk, useCacheFirst bool) (*SmartTemplate, error) {
	now := time.Now()

	if useCacheFirst {
		if now.Sub(cacher.cachedTime[key]) < cacher.cacheDuration {
			if tpl, ok := cacher.pool[key]; ok && tpl != nil {
				return tpl, nil
			}
		}
	}

	tpl, err := cacher.FetchTemplateByKey(ctx, key)
	if err == nil {
		return tpl, nil
	}

	if cacheOk {
		if tpl, ok := cacher.pool[key]; ok && tpl != nil {
			return tpl, nil
		}
	}

	return nil, err
}
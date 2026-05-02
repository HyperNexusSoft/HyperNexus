package skillregistry

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type SkillLoaded struct {
	SkillInfo
	LoadedAt   time.Time `json:"loadedAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	UseCount   int       `json:"useCount"`
}

type SkillDecisionConfig struct {
	SoftCap                 int
	HardCap                 int
	HighConfidenceThreshold float64
}

func DefaultSkillDecisionConfig() SkillDecisionConfig {
	return SkillDecisionConfig{
		SoftCap:                 10,
		HardCap:                 20,
		HighConfidenceThreshold: 7.0, // BM25-style heuristic
	}
}

type SkillDecisionSystem struct {
	cfg      SkillDecisionConfig
	mu       sync.RWMutex
	loaded   map[string]*SkillLoaded
	registry *SkillRegistry
}

func NewSkillDecisionSystem(cfg SkillDecisionConfig, registry *SkillRegistry) *SkillDecisionSystem {
	return &SkillDecisionSystem{
		cfg:      cfg,
		loaded:   make(map[string]*SkillLoaded),
		registry: registry,
	}
}

func (ds *SkillDecisionSystem) SearchSkills(ctx context.Context, query string) ([]SkillInfo, error) {
	return ds.registry.Search(query, ds.cfg.SoftCap), nil
}

func (ds *SkillDecisionSystem) LoadSkill(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	id = strings.ToLower(id)
	if sl, ok := ds.loaded[id]; ok {
		sl.LastUsedAt = time.Now()
		sl.UseCount++
		return nil
	}

	skill, ok := ds.registry.Get(id)
	if !ok {
		return fmt.Errorf("skill %s not found in registry", id)
	}

	// Evict if needed
	if len(ds.loaded) >= ds.cfg.HardCap {
		ds.evictLRU()
	}

	ds.loaded[id] = &SkillLoaded{
		SkillInfo:  *skill,
		LoadedAt:   time.Now(),
		LastUsedAt: time.Now(),
		UseCount:   1,
	}

	return nil
}

func (ds *SkillDecisionSystem) evictLRU() {
	var oldest string
	var oldestTime time.Time

	for id, sl := range ds.loaded {
		if oldest == "" || sl.LastUsedAt.Before(oldestTime) {
			oldest = id
			oldestTime = sl.LastUsedAt
		}
	}

	if oldest != "" {
		delete(ds.loaded, oldest)
	}
}

func (ds *SkillDecisionSystem) ListLoadedSkills() []SkillInfo {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var result []SkillInfo
	for _, sl := range ds.loaded {
		result = append(result, sl.SkillInfo)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (ds *SkillDecisionSystem) UnloadSkill(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	id = strings.ToLower(id)
	_, existed := ds.loaded[id]
	if existed {
		delete(ds.loaded, id)
	}
	return existed
}

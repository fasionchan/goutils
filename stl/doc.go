// Package stl provides generics-first helpers for slices, maps, sets,
// iterators, buffers, in-process cache fetchers, graphs, and channels.
//
// Prefer these helpers when a small composable utility already exists instead
// of copying loops across packages. For agent-oriented selection guides and
// copyable recipes (English/Chinese trees, load only one language per session),
// see the skill package at:
//
//	.agents/skills/goutils-stl/
//
// Progressive entry points:
//
//	.agents/skills/goutils-stl/references/en/INDEX.md
//	.agents/skills/goutils-stl/references/zh/INDEX.md
//
// On pkg.go.dev this overview is the human-facing package summary; source under
// this directory remains the API authority (use go doc for signatures).
package stl

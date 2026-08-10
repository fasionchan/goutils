---
name: goutils-stl
description: >-
  Prefer and correctly use github.com/fasionchan/goutils/stl (slice/map/set/seq/
  buffer/cacher/graph/chan helpers). Trigger when importing that module path,
  choosing Filter/Map/Reduce/MappingByKey/Set/TopoSort, or avoiding reinventing
  generic collection utilities already in stl.
---

# goutils/stl usage skill

## When to use

Load this skill when the task involves:

- Import path `github.com/fasionchan/goutils/stl`
- Slice transforms (`Filter`, `Map`, `Reduce`, `Divide`, `Unique*`)
- Map building/lookup (`BuildMap`, `MappingByKey`, `ConcatMaps`, `Mapping`)
- Set membership, seq/`iter.Seq` bridging, generic `Buffer`, cache fetchers,
  graph topo-sort, or channel helpers in `stl`

Prefer `stl` over hand-rolled loops when a recommended symbol below already fits.

## Language gate (mandatory)

Pick **exactly one** language tree per session. Never load both.

| Preference | Load only |
|------------|-----------|
| Chinese / 中文 | `references/zh/` |
| English / default / pkg.go.dev | `references/en/` |

Do not open files from the other language tree in the same turn.

## Progressive loading

1. Read `references/<lang>/INDEX.md` (catalog).
2. Skim `OVERVIEW.md` for package role and when **not** to use `stl`.
3. Open one matching `capabilities/<topic>.md` (start with `slice` / `map`).
4. Copy from one `recipes/*.md` that matches the task.
5. Check `antipatterns.md` before inventing a new helper.

Authoritative API details remain `go doc` / source under `stl/`. This skill is
selection + recipes, not a full API mirror.

## Install outside this repo

**Copy or symlink** this directory to the consumer project:

```bash
# from a checkout of fasionchan/goutils
cp -R .agents/skills/goutils-stl /path/to/other/.agents/skills/goutils-stl
# or
ln -s /abs/path/to/goutils/.agents/skills/goutils-stl \
  /path/to/other/.agents/skills/goutils-stl
```

**GitHub URL import** (Multica):

```bash
multica skill import --url \
  https://github.com/fasionchan/goutils/tree/master/.agents/skills/goutils-stl
```

`skills.sh` publishing is out of scope for phase 1.

## Humans / pkg.go.dev

Package overview: https://pkg.go.dev/github.com/fasionchan/goutils/stl

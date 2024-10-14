# Project Files

Pavement uses project files to provide configuration options for a pavement workspace. The project file uses a '.pvproj' extension. A default project file is assumed if you do not specify one. The default property values are shown below.

## Schema

A pvproj file contains the following fields. 

| Property | Type         | Description                | Default |
| -------- | ------------ | -------------------------- | ------- |
| include  | list(string) | include files with pattern | `["**/*.pv.hcl", "**/*.pv.yml", "**/*.pv.yaml", "**/*.pv.toml", "**/*.pv.json"]` |
| exclude  | list(string) | exclude files with pattern | `[]` |

## Examples

The default file

> default.pvproj

```
include:
- **/*.pv.hcl
- **/*.pv.yaml
- **/*.pv.toml
- **/*.pv.json
exclude: []
```
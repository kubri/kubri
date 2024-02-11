# Version Constraints

Version constraints in Kubri use [Semantic Versioning](https://semver.org/) with an optional prefix
of `v`.

You can use constraints to specify what versions are matched.

| Constraint | Description                                                                                                       |
| ---------- | ----------------------------------------------------------------------------------------------------------------- |
| `1`        | Any version that matches major version e.g. `>=1.0.0, <2.0.0`.                                                    |
| `1.2`      | Any version that matches minor version e.g. `>=1.2.0, <1.3.0`.                                                    |
| `1.2.3`    | Exactly specified version.                                                                                        |
| `>1.2`     | Any version that is above specified version.                                                                      |
| `>=1.2`    | Any version that is above or equal to specified version.                                                          |
| `<1.2`     | Any version that is below specified version.                                                                      |
| `<=1.2`    | Any version that is below or equal to specified version.                                                          |
| `~1.2`     | Any version that matches major version and is above or equal to minor version e.g. `>=1.2.0, <2.0.0`.             |
| `~1.2.3`   | Any version that matches major and minor version and is above or equal to patch version e.g. `>=1.2.3, <1.3.0`.   |
| `^1.2`     | Any version that matches major version and is above or equal to minor version e.g. `>=1.2.0, <2.0.0`.             |
| `^1.2.3`   | Any version that is above or equal to specified version and lower than next major version e.g. `>=1.2.3, <2.0.0`. |

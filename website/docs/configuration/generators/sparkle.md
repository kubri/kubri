---
sidebar_position: 10
---

# Sparkle Framework

Generate and publish a Sparkle appcast file from your `.dmg`, `.pkg`, `.mpkg`, `.msi` and `.exe`
files.

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable Sparkle.

### `folder`

- Type: `string`
- Default: `'appcast'`

Path to the directory on your target.

### `title`

- Type: `string`

Title of your appcast feed.

### `description`

- Type: `string`

Description of your appcast feed.

### `detect-os`

- Type: `map['macos'|'windows'|'windows-x64'|'windows-x86']string`

A map of globs to override detecting what file is what OS. Set this if your update packages are
`.zip` files or you use non-standard naming to differentiate between Windows 64-bit and 32-bit
installers.

:::info Default behaviour

- Files with the extensions `.dmg` and `.pkg` are picked up as `macos`.
- Files with the extensions `.exe` and `.msi` are picked up as `windows`.
  - Files matching `amd64|x64|x86[\W_]?64|64[\W_]?bit` are picked up as `windows-x64`.
  - Files matching `386|x86(?![\W_]?64)|ia32|32[\W_]?bit` are picked up as `windows-x86`.
  - Files matching both regular expressions are picked up as `windows`.

:::

#### Example

```yaml
sparkle:
  detect-os:
    macos: '*_MacOS.zip'
    windows: '*_Windows.zip'
    windows_x64: '*_Windows_x64.zip'
    windows_x86: '*_Windows_x86.zip'
```

### `params`

Set attributes on your appcast feed based on the OS & version of your release.

#### Example

```yaml
sparkle:
  - os: windows
    installer-arguments: /passive
  - os: macos
    minimum-system-version: '10.13.0'
```

### `params[*].os`

- Type: `'macos'|'windows'|'windows-x64'|'windows-x86'`

Apply these parameters only if the OS matches. Must be either `macos`, `windows`, `windows-x64` or
`windows-x86`.

:::info

Setting `os` to `windows` will also apply the parameters to `windows-x64` and `windows-x86`
releases.

:::

### `params[*].version`

- Type: `string`

A version constraint to limit what what releases these parameters should be applied to.  
See [Version Constraints](../../guides/version-constrains.md) for more information.

### `params[*].installer-arguments`

:::info Windows only

This parameter is only supported by WinSparkle, not Sparkle.

:::

- Type: `string`

On Windows, the enclosure is typically some kind of installer — an MSI, InnoSetup etc. It is often
useful to pass additional arguments to the installer when launching it, e.g. to force
non-interactive installation with reduce UI (notice that the installer shouldn’t be completely
invisible, because neither WinSparkle nor the hosting application is showing any UI at the time).

Useful values for common installers are listed below:

| Installer | Argument                | Description                                          |
| --------- | ----------------------- | ---------------------------------------------------- |
| InnoSetup | `/SILENT /SP- /NOICONS` | Shows only progress and errors, no startup prompt.   |
| MSI       | `/passive`              | Unattended mode, shows progress bar only.            |
| NSIS      | `/S`                    | Silent mode. No standard prompts or pages are shown. |

### `params[*].minimum-system-version`

- Type: `string`

The required minimum system operating version string for this update if provided. This version
string should contain three period-separated components e.g `10.13.0`.

### `params[*].minimum-autoupdate-version`

:::info MacOS only

This parameter is only supported by Sparkle, not WinSparkle.

:::

- Type: `string`

The minimum bundle version string this update requires for automatically downloading and installing
updates if provided. If an application’s bundle version meets this version requirement, it can
install the new update item in the background automatically. If the requirement is not met, the user
is always prompted to install the update.

### `params[*].ignore-skipped-upgrades-below-version`

:::info MacOS only

This parameter is only supported by Sparkle, not WinSparkle.

:::

- Type: `string`

Previously skipped upgrades by the user will be ignored if they skipped an update whose version
precedes this version.

### `params[*].critical-update`

- Type: `boolean`

Indicates whether or not the update item is critical. Critical updates are shown to the user more
promptly. Sparkle’s standard user interface also does not allow them to be skipped.

### `params[*].critical-update-below-version`

:::info MacOS only

This parameter is only supported by Sparkle, not WinSparkle.

:::

- Type: `string`

Indicates whether or not the update item is critical based on the version that is currently
installed.

## Example configuration

```yaml
sparkle:
  folder: appcast
  title: My app feed title
  description: My app feed description
  params:
    - os: windows
      installer-arguments: /passive
    - os: macos
      minimum-system-version: '10.13.0'
    - version: '1.0.0'
      critical-update: true
    - version: '> 1.0.0'
      critical-update-below-version: '1.0.0'
      minimum-autoupdate-version: '1.0.0'
    - version: '1.1.0'
      ignore-skipped-upgrades-below-version: '1.1.0'
```

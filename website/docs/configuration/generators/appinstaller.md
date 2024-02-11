---
sidebar_position: 20
---

# Windows App Installer

Generate and publish a Windows App Installer file from your `.msix`, `.msixbundle`, `.appx` and
`.appxbundle` files.

For more information see
https://learn.microsoft.com/en-us/uwp/schemas/appinstallerschema/element-update-settings

## Configuration

### `disabled`

- Type: `boolean`
- Default: `false`

Disable Windows App Installer.

### `folder`

- Type: `string`
- Default: `'appinstaller'`

Path to the directory on your target.

### `upload-packages`

- Type: `boolean`
- Default: `false`

Defines whether to upload packages to your target or reference them from your source.

### `on-launch`

Use `on-launch` to configure checking for updates on launch. This type of update can show UI.

### `on-launch.hours-between-update-checks`

- Type: `integer`
- Default: `24`

An integer that indicates how often (in how many hours) the system will check for updates to the
app. `0` to `255` inclusive. The default value is `24` (if this value is not specified). For example
if `hours-between-update-checks` is `3` then when the user launches the app, if the system has not
checked for updates within the past 3 hours, it will check for updates now.

### `on-launch.show-prompt`

- Type: `boolean`
- Default: `false`

A boolean that determines if UI will be shown to the user. This value is supported on Windows 10,
version 1903 and later.

### `on-launch.update-blocks-activation`

- Type: `boolean`
- Default: `false`

A boolean that determines if the UI shown to the user allows the user to launch the app without
taking the update, or if the user must take the update before launching the app. This attribute can
be set to `true` only if ShowPrompt is set to `true`. If set to `true` this means the UI the user
will see, allows the user to take the update or close the app. If set to `false` this means the UI
the user will see, allows the user to take the update or start the app without updating. In the
latter case, the update will be applied silently at an opportune time. This value is supported on
Windows 10, version 1903 and later.

:::info

`show-prompt` needs to be set to `true` if `update-blocks-activation` is set to `true`.

:::

### `automatic-background-task`

- Type: `boolean`
- Default: `false`

Checks for updates in the background every 8 hours independently of whether the user launched the
app. This type of update cannot show UI.

### `force-update-from-any-version`

- Type: `boolean`
- Default: `false`

Allows the app to update from version x to version x++ or to downgrade from version x to version
x--. Without this element, the app can only move to a higher version.

## Example

```yaml
appinstaller:
  folder: windows
  upload-packages: true
  on-launch:
    hours-between-update-checks: 12
    show-prompt: true
    update-blocks-activation: true
  automatic-background-task: true
  force-update-from-any-version: true
```

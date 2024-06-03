<div align="center">
  <h1>shazam.sh</h1>

  <p>Dotfiles 🗃️ manager on steroids ⚡.</p>
  <p>Makes managing 🤖 your dotfiles a breeze 💪.</p>

  <a href="#why">Why?</a>
  <span> • </span>
  <a href="#install">Install</a>
  <span> • </span>
  <a href="#configuration">Configuration</a>
  <span> • </span>
  <a href="#example">Example</a>
  <span> • </span>
  <a href="#advanced-configuration">Advanced</a>
  <span> • </span>
  <a href="#example-dotfiles-repository">Example Dotfiles Repo</a>

  <p></p>

[![Go](https://img.shields.io/badge/Made%20with%20Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=ffffff)](https://golang.org)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/mistweaverco/shazam.sh?style=for-the-badge)](https://github.com/mistweaverco/shazam.sh/releases/latest)

  <p></p>

[![GitHub Workflow Build Status](https://img.shields.io/github/actions/workflow/status/mistweaverco/shazam.sh/build.yml?style=for-the-badge)](https://github.com/mistweaverco/shazam.sh/actions/workflows/build.yml)
[![GitHub Workflow Lint Status](https://img.shields.io/github/actions/workflow/status/mistweaverco/shazam.sh/lint.yml?style=for-the-badge&label=lint)](https://github.com/mistweaverco/shazam.sh/actions/workflows/lint.yml)
[![GitHub Workflow Test Status](https://img.shields.io/github/actions/workflow/status/mistweaverco/shazam.sh/test.yml?style=for-the-badge&label=test)](https://github.com/mistweaverco/shazam.sh/actions/workflows/test.yml)

</div>

# Why

Why not simply use GNU stow?
GNU stow is great, but AFAIK it doesn't support symlinking files to
various different directories easily.

This is a feature that I wanted to have in my dotfiles manager.

Also it was just fun to write this.

# Install

Grab the latest, greatest and hottest release from the
[releases page](https://github.com/mistweaverco/shazam.sh/releases),

or use the following commands depending on your OS:

- [Linux](#install-on-linux)
- [macOS](#install-on-macos)
- [Windows](#install-on-windows)

## Install on Linux

On linux you can use `wget` to download the binary.

```sh
wget -qO- https://github.com/mistweaverco/shazam.sh/releases/latest/download/shazam-linux
```

Then you can move the binary to a directory in your `$PATH`.

```sh
sudo mv shazam-linux /usr/bin/shazam
```

## Install on macOS

On macOS you can use `wget` to download the binary.

```sh
wget -qO- https://github.com/mistweaverco/shazam.sh/releases/latest/download/shazam-macos
```

Then you can move the binary to a directory in your `$PATH`.

```sh
sudo mv shazam-macos /usr/local/bin/shazam
```

## Install on Windows

> [!NOTE]
> To be honest I don't know if dotfiles management is a thing on Windows,
> but if it is, you can use `shazam` to manage your dotfiles.

On Windows you can use the powershell `Invoke-WebRequest` cmdlet to download the binary.

```powershell
Invoke-WebRequest -Uri https://github.com/mistweaverco/shazam.sh/releases/latest/download/shazam-windows -OutFile shazam.exe
```

Then you can move the binary to a directory in your `$PATH`.

```powershell
Move-Item -Path .\shazam.exe -Destination C:\Windows\System32\shazam.exe
```
# Configuration

The configuration file is a YAML file that contains a list of symlinks
that you want to create.

The default configuration file is called `shazam.yml` and
should be placed in the root of your dotfiles directory.

Here is an example

`shazam.yml`

```yaml
configurations: # one of the root directories in your dotfiles directory
  - name: neovim # one subdirectory in the configurations directory
    files: # a list of files or directories that you want to symlink
      - source: nvim
        destination: $HOME/.config/nvim
  - name: tmux
    files:
      - source: .tmux
        destination: $HOME/.tmux
      - source: .tmux.conf
        destination: $HOME/.tmux.conf
  - name: editorconfig
    files:
      - source: .editorconfig
        destination: $HOME/.editorconfig
  - name: oh-my-zsh
    files:
      - source: .oh-my-zsh/custom
        destination: $HOME/.oh-my-zsh/custom
  - name: git
    files:
      - source: .gitconfig
        destination: $HOME/.gitconfig
  - name: WezTerm
    files:
      - source: .wezterm.lua
        destination: $HOME/.wezterm.lua
  - name: ssh
    files:
      - source: .ssh
        destination: $HOME/.ssh
```

The `configurations` key is corresponds to a
directory in your dotfiles directory.

It can be any name you want.

The `name` key corresponds to the name of a directory
in their respective root directory.

The `files` key is a list of files that you want to symlink.

Environment variables can be used in both the
`source` and `destination` keys.

## Example

You have a dotfiles directory with the following structure

```text
dotfiles
├── shazam.yml
├── configurations
│   ├── git
│   │   ├── gitconfig
│   ├── ssh
│   │   ├── config

```

You want to symlink the `gitconfig` file to `$HOME/.gitconfig`
and the `config` file to `$HOME/.ssh/config`.

Your `shazam.yml` file would look like this

```yaml
configurations:
  - name: git
    files:
      - source: gitconfig
        destination: $HOME/.gitconfig
  - name: ssh
    files:
      - source: config
        destination: $HOME/.ssh/config
```

You can then run `shazam` to symlink the files.

## Advanced configuration

You can have multiple root nodes in your configuration file,
each corresponding to a different directory in your dotfiles directory.

This can be useful if you want to structure your dotfiles directory
in a certain way.

> [!TIP]
> Here is an example

You have a dotfiles directory with the following structure:

```text
dotfiles
├── shazam.yml
├── configurations
│   ├── git
│   │   ├── gitconfig
│   ├── ssh
│   │   ├── config
├── neovimfiles
│   ├── neovim
│   │   ├── nvim
│   │   │   ├── ...
```
Then your `shazam.yml` file would look like this:

```yaml
configurations:
  - name: git
    files:
      - source: gitconfig
        destination: $HOME/.gitconfig
  - name: ssh
    files:
      - source: config
        destination: $HOME/.ssh/config
neovimfiles:
  - name: neovim
    files:
      - source: nvim
        destination: $HOME/.config/nvim
```

## `--config` flag

You can also have a custom configuration file name.

To specify a custom configuration file name, you can use the `--config` flag.

```sh
shazam --config my-others-machine-config.yml
```

## `--dry-run` flag

You can also run `shazam` with the `--dry-run` flag to see what symlinks will be created.

```sh
shazam --dry-run
```

## `--root` flag

```sh
shazam --root configurations
```
This will run all symlink actions in the of the `configurations` `root`.

> [!TIP]
> If you want to create a certain symlink,
> you can use the `--root` and  `--only` flags in combination.

## `--only` flag

```sh
shazam --only git
```

This will only create the `git` symlink(s).

> [!TIP]
> If you want to create a certain symlink,
> you can use the `--root` and  `--only` flags in combination.

## `--pull-in-existing` flag

If the destination file already exists,
`shazam` will do nothing and skip the creation of the symlink.

But it will notice you that the file already exists.

If you want to pull in that file into your dotfiles directory,
and also symlink it, you can use the `--pull-in-existing` flag.

```sh
shazam --pull-in-existing
```

> [!WARNING]
> It's recommended to backup any existing files in your dotfiles directory,
> before running this command,
> because it will overwrite the existing files.

> [!TIP]
> If you have some sort of version control system in place,
> you can use that to backup your files.

## `--dotfiles-path` flag

You can pass a path to `shazam` to specify the path to your dotfiles directory
via the `--dotfiles-path` flag.

```sh
shazam --dotfiles-path /path/to/dotfiles
```

This will use the specified path as the root directory for your dotfiles.

## `--path` flag

You can pass a path to `shazam` to specify a single configuration file,
that you want shazam.sh to create symlinks for via the `--path` flag.

shazam.sh will try to figure out what files to symlink
based on the configuration file.

It can be used in conjunction with the `--dotfiles-path` flag and/or
the `--config` flag.

```sh
shazam --path configurations/git
```
> [!TIP]
> The `--path` flag works by checking if the path starts
> with the root directory name + the configuration name + the file or directory name.
>
> So you could also run `shazam --path configurations/git/gitconfig`

> [!NOTE]
> The `--path` flag is not meant to be used in conjunction
> with the `--root` and `--only` flags.

## Example dotfiles repository

An example dotfiles repository can be found [here](https://github.com/gorillamoe/dotfiles).

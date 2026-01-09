# sysup

**sysup** is a lightweight, configurable system bootstrap tool written in Go. It helps you automate the installation of packages and applications on Fedora and Debian-based systems (Ubuntu, Pop!_OS, etc.). It supports system package managers (`dnf`, `apt`) and Flatpak.

## Features

*   **Multi-Distribution Support:** dedicated commands for Fedora and Debian/Ubuntu.
*   **Package Management:** Automates `dnf` and `apt` package installations.
*   **Flatpak Integration:** Installs Flatpak apps, manages remotes (like Flathub), and prevents redundant re-installations.
*   **Repository Management:** configure custom repositories (for Fedora).
*   **Post-Install Scripts:** Run custom shell scripts after installation (e.g., changing shell, setting up dotfiles).
*   **Idempotent:** safe to run multiple times; skips already installed repositories and Flatpaks.

## Installation

### Binary Release (Recommended)
Download the latest `.rpm` (for Fedora/RHEL) or `.deb` (for Debian/Ubuntu) from the [Releases](https://github.com/disturb16/sysup/releases) page.

**Fedora:**
```bash
sudo dnf install ./sysup_*.rpm
```

**Debian/Ubuntu:**
```bash
sudo apt install ./sysup_*.deb
```

### Go Install
If you have Go installed, you can install `sysup` directly:

```bash
go install github.com/disturb16/sysup@latest
```

## Usage

1.  **Create a Configuration File:**
    Create a `config.yaml` file defining what you want to install.

    ```yaml
    # config.yaml example

    # [Fedora Only] Extra DNF repositories
    repositories:
      - https://brave-browser-rpm-release.s3.brave.com/brave-browser.repo
      - https://download.docker.com/linux/fedora/docker-ce.repo

    # [Fedora Only] DNF packages to install
    dnf:
      - git
      - vim
      - docker-ce
      - brave-browser

    # [Debian/Ubuntu Only] APT packages to install
    apt:
      - git
      - curl
      - build-essential

    # [Cross-Platform] Configure Flatpak Remotes (optional, defaults to Flathub)
    flatpak_remotes:
      - name: flathub
        url: https://flathub.org/repo/flathub.flatpakrepo
      - name: gnome-nightly
        url: https://nightly.gnome.org/gnome-nightly.flatpakrepo

    # [Cross-Platform] Flatpak applications to install
    flatpak:
      - com.spotify.Client
      - com.visualstudio.code
      - org.gnome.Calculator

    # Post-installation shell commands
    post_install:
      - echo "System setup complete!"
    ```

2.  **Run the Installer:**

    **Fedora:**
    ```bash
    # Run everything
    sudo sysup install fedora

    # Run only specific parts
    sudo sysup install fedora --dnf      # Only DNF packages
    sudo sysup install fedora --flatpak  # Only Flatpaks
    ```

    **Debian/Ubuntu:**
    ```bash
    # Run everything
    sudo sysup install debian

    # Run only specific parts
    sudo sysup install debian --apt      # Only APT packages
    sudo sysup install debian --flatpak  # Only Flatpaks
    ```

    *Note: `sysup` usually requires `sudo` privileges to install system packages.*

## Building from Source

```bash
git clone https://github.com/disturb16/sysup.git
cd sysup
go build -o sysup
```

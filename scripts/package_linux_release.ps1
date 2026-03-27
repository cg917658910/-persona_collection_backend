param(
    [string]$Goos = "linux",
    [string]$Goarch = "amd64",
    [string]$Version = "",
    [string]$AppName = "pm-backend",
    [string]$OutputRoot = "dist",
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendRoot = Split-Path -Parent $scriptDir
$repoRoot = Split-Path -Parent $backendRoot

if ([string]::IsNullOrWhiteSpace($Version)) {
    $Version = Get-Date -Format "yyyyMMdd-HHmmss"
}

$releaseName = "$AppName-$Version-$Goos-$Goarch"
$distRoot = Join-Path $backendRoot $OutputRoot
$stageDir = Join-Path $distRoot $releaseName
$deployDir = Join-Path $stageDir "deploy"
$systemdDir = Join-Path $deployDir "systemd"
$binaryName = $AppName
$archivePath = Join-Path $distRoot "$releaseName.tar.gz"

Write-Host "==> backend root: $backendRoot"
Write-Host "==> release name: $releaseName"

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "go is not available in PATH."
}

New-Item -ItemType Directory -Force -Path $distRoot | Out-Null
if (Test-Path $stageDir) {
    Remove-Item -Recurse -Force $stageDir
}
New-Item -ItemType Directory -Force -Path $stageDir | Out-Null
New-Item -ItemType Directory -Force -Path $deployDir | Out-Null
New-Item -ItemType Directory -Force -Path $systemdDir | Out-Null

Push-Location $backendRoot
try {
    if (-not $SkipTests) {
        Write-Host "==> running go test ./..."
        go test ./...
    }

    Write-Host "==> building $binaryName for $Goos/$Goarch"
    $env:CGO_ENABLED = "0"
    $env:GOOS = $Goos
    $env:GOARCH = $Goarch
    go build -trimpath -ldflags "-s -w -X main.version=$Version" -o (Join-Path $stageDir $binaryName) ./cmd/server

    Write-Host "==> copying runtime files"
    Copy-Item ".env.example" (Join-Path $stageDir ".env.example")
    if (Test-Path "public") {
        Copy-Item "public" (Join-Path $stageDir "public") -Recurse
    }
}
finally {
    Pop-Location
}

$runSh = @'
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [[ ! -f ".env" && -f ".env.example" ]]; then
  cp .env.example .env
  echo "Created .env from .env.example. Review it before exposing this service."
fi

chmod +x ./pm-backend
exec ./pm-backend
'@
Set-Content -Path (Join-Path $stageDir "run.sh") -Value $runSh -NoNewline

$installSh = @'
#!/usr/bin/env bash
set -euo pipefail

APP_NAME="${APP_NAME:-pm-backend}"
INSTALL_DIR="${INSTALL_DIR:-/opt/pm-backend}"
SERVICE_NAME="${SERVICE_NAME:-pm-backend}"
RUN_USER="${RUN_USER:-www-data}"
RUN_GROUP="${RUN_GROUP:-www-data}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

mkdir -p "$INSTALL_DIR"
if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete "$PACKAGE_ROOT/" "$INSTALL_DIR/"
else
  find "$INSTALL_DIR" -mindepth 1 -maxdepth 1 ! -name ".env" -exec rm -rf {} +
  cp -a "$PACKAGE_ROOT/." "$INSTALL_DIR/"
fi

if [[ ! -f "$INSTALL_DIR/.env" && -f "$INSTALL_DIR/.env.example" ]]; then
  cp "$INSTALL_DIR/.env.example" "$INSTALL_DIR/.env"
fi

chown -R "$RUN_USER:$RUN_GROUP" "$INSTALL_DIR"
chmod +x "$INSTALL_DIR/$APP_NAME" "$INSTALL_DIR/run.sh"

SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
sed \
  -e "s|__APP_NAME__|$APP_NAME|g" \
  -e "s|__INSTALL_DIR__|$INSTALL_DIR|g" \
  -e "s|__RUN_USER__|$RUN_USER|g" \
  -e "s|__RUN_GROUP__|$RUN_GROUP|g" \
  "$INSTALL_DIR/deploy/systemd/pm-backend.service" > "$SERVICE_FILE"

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"
systemctl status "$SERVICE_NAME" --no-pager
'@
Set-Content -Path (Join-Path $deployDir "install.sh") -Value $installSh -NoNewline

$serviceTpl = @'
[Unit]
Description=__APP_NAME__
After=network.target

[Service]
Type=simple
User=__RUN_USER__
Group=__RUN_GROUP__
WorkingDirectory=__INSTALL_DIR__
EnvironmentFile=__INSTALL_DIR__/.env
ExecStart=__INSTALL_DIR__/__APP_NAME__
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
'@
Set-Content -Path (Join-Path $systemdDir "pm-backend.service") -Value $serviceTpl -NoNewline

$deployReadme = @'
# Linux deploy

1. Upload this tar.gz to the Linux server.
2. Extract it:

   tar -xzf pm-backend-<version>-linux-<arch>.tar.gz

3. Enter the release directory and review `.env.example`.
4. Install:

   sudo RUN_USER=www-data RUN_GROUP=www-data INSTALL_DIR=/opt/pm-backend ./deploy/install.sh

5. Check service logs:

   sudo journalctl -u pm-backend -f
'@
Set-Content -Path (Join-Path $deployDir "README.md") -Value $deployReadme -NoNewline

Write-Host "==> writing metadata"
$manifest = @{
    app_name = $AppName
    version = $Version
    goos = $Goos
    goarch = $Goarch
    built_at = (Get-Date).ToString("s")
}
$manifest | ConvertTo-Json | Set-Content -Path (Join-Path $stageDir "manifest.json")

Push-Location $distRoot
try {
    if (Test-Path $archivePath) {
        Remove-Item -Force $archivePath
    }
    Write-Host "==> creating $archivePath"
    tar -czf $archivePath $releaseName
}
finally {
    Pop-Location
}

Write-Host "==> done"
Write-Host "Package: $archivePath"

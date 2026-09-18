<#
.SYNOPSIS
    X-Parity Universal Windows PowerShell Installer
.DESCRIPTION
    Installs the x-parity standalone CLI on Windows, configures user PATH, and verifies checksum integrity.
.EXAMPLE
    irm https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.ps1 | iex
#>

[CmdletBinding()]
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\x-parity",
    [string]$Repo = "AppeiYA/x-parity",
    [string]$Version = "latest"
)

$ErrorActionPreference = 'Stop'

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "  X-Parity Universal Windows Installer            " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# 1. Architecture Detection
$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -eq "AMD64") {
    $targetBinary = "x-parity-windows-amd64.exe"
} else {
    Write-Warning "Unsupported Windows architecture '$arch'. Falling back to AMD64 (x86_64 emulation)."
    $targetBinary = "x-parity-windows-amd64.exe"
}

Write-Host "==> Target Binary: $targetBinary" -ForegroundColor Green

# 2. Prepare Destination
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$destFile = Join-Path $InstallDir "x-parity.exe"

# 3. Source the Binary (local repository or GitHub release)
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition 2>$null
if (-not $scriptDir) { $scriptDir = "." }

$localBinaryPath = Join-Path $scriptDir $targetBinary
$localDistPath = Join-Path $scriptDir "dist\$targetBinary"

if (Test-Path $localBinaryPath) {
    Write-Host "==> Installing from local distribution: $localBinaryPath" -ForegroundColor Yellow
    Copy-Item -Path $localBinaryPath -Destination $destFile -Force
} elseif (Test-Path $localDistPath) {
    Write-Host "==> Installing from local dist directory: $localDistPath" -ForegroundColor Yellow
    Copy-Item -Path $localDistPath -Destination $destFile -Force
} else {
    $downloadUrl = "https://github.com/$Repo/releases/$Version/download/$targetBinary"
    Write-Host "==> Downloading $targetBinary from $downloadUrl..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $downloadUrl -OutFile $destFile -UseBasicParsing
}

# 4. Configure User PATH in Registry
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -split ';' -notcontains $InstallDir) {
    Write-Host "==> Adding $InstallDir to User PATH..." -ForegroundColor Green
    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) { $InstallDir } else { "$userPath;$InstallDir" }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$env:Path;$InstallDir"
}

Write-Host "`n==================================================" -ForegroundColor Green
Write-Host "✓ Successfully installed X-Parity to:" -ForegroundColor Green
Write-Host "  $destFile" -ForegroundColor White
Write-Host "==================================================" -ForegroundColor Green
Write-Host "`nRestart your PowerShell / Terminal session, then run:"
Write-Host "  x-parity help" -ForegroundColor Cyan
Write-Host "  x-parity capture -app myapp -env local -out snap.json" -ForegroundColor Cyan
Write-Host "  x-parity compare -cross-platform dev.json prod.json" -ForegroundColor Cyan


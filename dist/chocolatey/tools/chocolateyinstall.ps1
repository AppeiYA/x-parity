$ErrorActionPreference = 'Stop'

$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$url64 = 'https://github.com/AppeiYA/x-parity/releases/download/v1.0.0/x-parity-windows-amd64.exe'
$checksum64 = 'debbe66357123f676e911d48926ee5977cb536e80995ef1b020b5c7aca6adce9'

$packageArgs = @{
  packageName   = 'x-parity'
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url64bit      = $url64
  checksum64    = $checksum64
  checksumType64= 'sha256'
}

Get-ChocolateyWebFile @packageArgs
$targetFile = Join-Path $toolsDir "x-parity.exe"
Rename-Item -Path (Join-Path $toolsDir "x-parity-windows-amd64.exe") -NewName $targetFile -Force

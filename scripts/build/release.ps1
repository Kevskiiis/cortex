param(
    [Parameter(Mandatory = $true)]
    [ValidatePattern('^v[0-9]+\.[0-9]+\.[0-9]+$')]
    [string]$Version
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$dist = Join-Path $repoRoot 'dist'
$binary = Join-Path $dist 'cortex.exe'
$archive = Join-Path $repoRoot "cortex-$Version-windows-amd64.zip"

if (Test-Path $dist) {
    Remove-Item $dist -Recurse -Force
}
if (Test-Path $archive) {
    Remove-Item $archive -Force
}
if (Test-Path "$archive.sha256") {
    Remove-Item "$archive.sha256" -Force
}

New-Item -ItemType Directory -Force -Path $dist | Out-Null
Push-Location $repoRoot
try {
    go build -ldflags "-X main.version=$Version" -o $binary ./cmd/cortex
} finally {
    Pop-Location
}

Compress-Archive -Path $binary -DestinationPath $archive -Force
$hash = (Get-FileHash -Algorithm SHA256 $archive).Hash.ToLowerInvariant()
"$hash  $(Split-Path $archive -Leaf)" | Set-Content "$archive.sha256"

Write-Output "Created $archive"
Write-Output "Created $archive.sha256"
param(
    [string]$Output = ""
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($Output)) {
    $Output = Join-Path $Root "bin\runner-monitor.exe"
}

New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Output) | Out-Null
Push-Location $Root
try {
    go build -o $Output ./cmd/runner-monitor
    Write-Output $Output
}
finally {
    Pop-Location
}


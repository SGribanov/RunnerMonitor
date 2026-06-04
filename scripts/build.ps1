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
    $Config = Join-Path (Split-Path -Parent $Output) "runner-monitor.json"
    if (-not (Test-Path -LiteralPath $Config)) {
        $DefaultConfig = [ordered]@{
            projectsRoot = "C:\Repos"
            windowsRunnerRoots = @("C:\Runners")
            wslRunnerRoots = @("/home/gsv777/Runners")
            linuxRunnerRoots = @("/opt/Runners", "/srv/Runners")
            tuiRefreshIntervalSeconds = 5
            wslSudoPassword = ""
        }
        $DefaultConfig | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath $Config -Encoding UTF8
    }
    Write-Output $Output
}
finally {
    Pop-Location
}

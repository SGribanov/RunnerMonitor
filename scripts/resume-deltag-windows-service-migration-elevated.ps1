$ErrorActionPreference = "Stop"

$LogPath = "D:\Repos\RunnerMonitor\reports\windows-service-runner-migration-2026-06-04.log"
function Write-Log {
    param([string]$Message)
    $line = "[{0}] {1}" -f (Get-Date -Format "yyyy-MM-dd HH:mm:ss"), $Message
    $line | Tee-Object -FilePath $LogPath -Append
}

function Assert-Admin {
    $principal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        throw "This script must run from elevated PowerShell."
    }
}

function Invoke-CheckedNative {
    param(
        [string]$FilePath,
        [string[]]$ArgumentList,
        [int[]]$SuccessExitCodes = @(0)
    )
    & $FilePath @ArgumentList
    $code = $LASTEXITCODE
    if ($SuccessExitCodes -notcontains $code) {
        throw "$FilePath failed with exit code $code"
    }
}

function Resolve-VersionTarget {
    param([string]$Root, [string]$Name)
    $versionDirs = Get-ChildItem -LiteralPath $Root -Directory -Filter "$Name.*" | Sort-Object Name -Descending
    if (-not $versionDirs) {
        throw "No versioned directory found for $Name under $Root."
    }
    return $versionDirs[0].FullName
}

function Retarget-RunnerJunctions {
    param([string]$Root)
    foreach ($name in @("bin", "externals")) {
        $path = Join-Path $Root $name
        $target = Resolve-VersionTarget -Root $Root -Name $name
        if (Test-Path -LiteralPath $path) {
            Remove-Item -LiteralPath $path -Recurse -Force
        }
        New-Item -ItemType Junction -Path $path -Target $target | Out-Null
        Write-Log "Retargeted $path -> $target"
    }
}

function Wait-ServiceState {
    param([string]$Name, [string]$State, [int]$TimeoutSeconds = 90)
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    while ((Get-Date) -lt $deadline) {
        $svc = Get-Service -Name $Name -ErrorAction Stop
        if ($svc.Status.ToString() -eq $State) {
            return
        }
        Start-Sleep -Seconds 1
    }
    throw "Service $Name did not reach state $State within $TimeoutSeconds seconds."
}

function Get-GitHubRunner {
    $response = gh api "repos/SGribanov/DeltaG/actions/runners" | ConvertFrom-Json
    $runner = $response.runners | Where-Object { $_.name -eq "deltag-win" } | Select-Object -First 1
    if (-not $runner) {
        throw "GitHub runner SGribanov/DeltaG deltag-win was not found."
    }
    return $runner
}

function Wait-GitHubRunnerOnline {
    $deadline = (Get-Date).AddSeconds(90)
    while ((Get-Date) -lt $deadline) {
        $runner = Get-GitHubRunner
        if ($runner.status -eq "online" -and -not $runner.busy) {
            return
        }
        Start-Sleep -Seconds 3
    }
    throw "GitHub runner SGribanov/DeltaG deltag-win did not become online/busy=false."
}

function Assert-ProcessFromPath {
    param([string]$Root)
    $prefix = (Resolve-Path -LiteralPath $Root).Path.TrimEnd("\") + "\"
    $proc = Get-CimInstance Win32_Process | Where-Object {
        $_.Name -eq "Runner.Listener.exe" -and
        $_.ExecutablePath -and
        $_.ExecutablePath.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase)
    } | Select-Object -First 1
    if (-not $proc) {
        throw "Runner.Listener.exe was not found under $Root."
    }
    Write-Log "Validated DeltaG process: pid=$($proc.ProcessId), path=$($proc.ExecutablePath)"
}

Assert-Admin

$serviceName = "actions.runner.SGribanov-DeltaG.deltag-win"
$oldPath = "C:\github-runners\deltag"
$newPath = "C:\Runners\SGribanov-DeltaG\deltag-win"
$backupPath = "C:\Runners-backup\github-runners-deltag-deltag-win-move-2026-06-04.zip"
$serviceExe = Join-Path $newPath "bin\RunnerService.exe"

try {
    Write-Log "=== Resuming DeltaG Windows service migration via robocopy ==="
    $backup = Get-Item -LiteralPath $backupPath -ErrorAction Stop
    if ($backup.Length -le 0) {
        throw "Backup is empty: $backupPath"
    }
    Write-Log "Using existing backup: $backupPath ($($backup.Length) bytes)"

    $runner = Get-GitHubRunner
    if ($runner.busy) {
        throw "GitHub runner SGribanov/DeltaG/deltag-win is busy; aborting."
    }
    Write-Log "GitHub precheck: status=$($runner.status), busy=$($runner.busy), version=$($runner.version)"

    $svc = Get-CimInstance Win32_Service -Filter "Name='$serviceName'"
    if ($svc.State -ne "Stopped") {
        Stop-Service -Name $serviceName -Force
        Wait-ServiceState -Name $serviceName -State "Stopped"
    }

    if (-not (Test-Path -LiteralPath $newPath)) {
        New-Item -ItemType Directory -Path $newPath -Force | Out-Null
    }

    Invoke-CheckedNative -FilePath "robocopy.exe" -ArgumentList @(
        $oldPath,
        $newPath,
        "/MIR",
        "/COPY:DAT",
        "/DCOPY:DAT",
        "/R:2",
        "/W:2",
        "/XJ",
        "/NFL",
        "/NDL",
        "/NP"
    ) -SuccessExitCodes @(0, 1, 2, 3, 4, 5, 6, 7)
    Write-Log "Robocopy completed $oldPath -> $newPath"

    Retarget-RunnerJunctions -Root $newPath
    if (-not (Test-Path -LiteralPath $serviceExe)) {
        throw "RunnerService.exe not found at $serviceExe"
    }

    Invoke-CheckedNative -FilePath "sc.exe" -ArgumentList @("config", $serviceName, "binPath=", $serviceExe, "start=", "demand")
    Write-Log "Configured service path and manual startup."

    Start-Service -Name $serviceName
    Wait-ServiceState -Name $serviceName -State "Running"
    Write-Log "Started service $serviceName"

    Start-Sleep -Seconds 5
    Assert-ProcessFromPath -Root $newPath
    Wait-GitHubRunnerOnline
    Write-Log "GitHub validation passed for SGribanov/DeltaG deltag-win"

    $archivePath = "C:\github-runners\deltag-migrated-20260604"
    try {
        if (Test-Path -LiteralPath $oldPath) {
            Rename-Item -LiteralPath $oldPath -NewName "deltag-migrated-20260604" -ErrorAction Stop
            Write-Log "Renamed old folder to $archivePath"
        }
    } catch {
        Write-Log "Old folder rename failed: $($_.Exception.Message)"
        $oldRunnerFile = Join-Path $oldPath ".runner"
        if (Test-Path -LiteralPath $oldRunnerFile) {
            Rename-Item -LiteralPath $oldRunnerFile -NewName ".runner_migrated_20260604" -Force
            Write-Log "Renamed old .runner to disable duplicate discovery."
        }
    }

    Write-Log "DeltaG Windows service migration completed."
} catch {
    Write-Log "ERROR: $($_.Exception.Message)"
    throw
}

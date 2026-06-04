$ErrorActionPreference = "Stop"

$LogPath = "C:\Repos\RunnerMonitor\reports\windows-service-runner-migration-2026-06-04.log"
New-Item -ItemType Directory -Path (Split-Path -Parent $LogPath) -Force | Out-Null

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

function Get-GitHubRunner {
    param(
        [string]$Repo,
        [string]$Runner
    )
    $response = gh api "repos/$Repo/actions/runners" | ConvertFrom-Json
    $match = $response.runners | Where-Object { $_.name -eq $Runner } | Select-Object -First 1
    if (-not $match) {
        throw "GitHub runner $Repo/$Runner was not found."
    }
    return $match
}

function Wait-ServiceState {
    param(
        [string]$Name,
        [string]$State,
        [int]$TimeoutSeconds = 60
    )
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

function Resolve-VersionTarget {
    param(
        [string]$Root,
        [string]$Name
    )
    $versionDirs = Get-ChildItem -LiteralPath $Root -Directory -Filter "$Name.*" |
        Sort-Object Name -Descending
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
            $item = Get-Item -LiteralPath $path -Force
            if ($item.LinkType -ne "Junction") {
                throw "Expected $path to be a junction, got LinkType=$($item.LinkType)."
            }
            Remove-Item -LiteralPath $path -Force
        }
        New-Item -ItemType Junction -Path $path -Target $target | Out-Null
        Write-Log "Retargeted $path -> $target"
    }
}

function Wait-GitHubRunnerOnline {
    param(
        [string]$Repo,
        [string]$Runner,
        [int]$TimeoutSeconds = 90
    )
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    while ((Get-Date) -lt $deadline) {
        $ghRunner = Get-GitHubRunner -Repo $Repo -Runner $Runner
        if ($ghRunner.status -eq "online" -and -not $ghRunner.busy) {
            return
        }
        Start-Sleep -Seconds 3
    }
    throw "GitHub runner $Repo/$Runner did not become online/busy=false within $TimeoutSeconds seconds."
}

function Assert-ProcessFromPath {
    param(
        [string]$Runner,
        [string]$Root
    )
    $prefix = (Resolve-Path -LiteralPath $Root).Path.TrimEnd("\") + "\"
    $proc = Get-CimInstance Win32_Process | Where-Object {
        $_.Name -eq "Runner.Listener.exe" -and
        $_.ExecutablePath -and
        $_.ExecutablePath.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase)
    } | Select-Object -First 1
    if (-not $proc) {
        throw "Runner.Listener.exe for $Runner was not found under $Root."
    }
    Write-Log "Validated process for $($Runner): pid=$($proc.ProcessId), path=$($proc.ExecutablePath)"
}

function Invoke-RunnerMigration {
    param([pscustomobject]$Runner)

    Write-Log "=== Migrating $($Runner.Repo) $($Runner.Name) ==="
    $ghRunner = Get-GitHubRunner -Repo $Runner.Repo -Runner $Runner.Name
    if ($ghRunner.busy) {
        throw "GitHub runner $($Runner.Repo)/$($Runner.Name) is busy; aborting."
    }
    Write-Log "GitHub precheck: status=$($ghRunner.status), busy=$($ghRunner.busy), version=$($ghRunner.version)"

    $svc = Get-CimInstance Win32_Service -Filter ("Name='{0}'" -f $Runner.ServiceName)
    if (-not $svc) {
        throw "Service not found: $($Runner.ServiceName)"
    }
    Write-Log "Service precheck: state=$($svc.State), startMode=$($svc.StartMode), path=$($svc.PathName)"

    $oldExists = Test-Path -LiteralPath $Runner.OldPath
    $newExists = Test-Path -LiteralPath $Runner.NewPath
    if (-not $oldExists -and $newExists) {
        Write-Log "Old path is already gone and target exists; validating migrated runner."
        if ($svc.State -ne "Running") {
            Start-Service -Name $Runner.ServiceName
            Wait-ServiceState -Name $Runner.ServiceName -State "Running" -TimeoutSeconds 90
        }
        Assert-ProcessFromPath -Runner $Runner.Name -Root $Runner.NewPath
        Wait-GitHubRunnerOnline -Repo $Runner.Repo -Runner $Runner.Name
        Write-Log "Already migrated validation passed for $($Runner.Repo) $($Runner.Name)"
        return
    }
    if (-not $oldExists) {
        throw "Old path not found: $($Runner.OldPath)"
    }
    if ($newExists) {
        throw "Both old and new paths exist for $($Runner.Name); refusing to continue."
    }

    if ($svc.State -ne "Stopped") {
        Stop-Service -Name $Runner.ServiceName -Force
        Wait-ServiceState -Name $Runner.ServiceName -State "Stopped" -TimeoutSeconds 90
    }
    Write-Log "Stopped service $($Runner.ServiceName)"

    New-Item -ItemType Directory -Path (Split-Path -Parent $Runner.NewPath) -Force | Out-Null
    New-Item -ItemType Directory -Path (Split-Path -Parent $Runner.BackupPath) -Force | Out-Null
    if (Test-Path -LiteralPath $Runner.BackupPath) {
        $backup = Get-Item -LiteralPath $Runner.BackupPath
        if ($backup.Length -le 0) {
            throw "Existing backup is empty: $($Runner.BackupPath)"
        }
        Write-Log "Backup already exists, reusing: $($Runner.BackupPath) ($($backup.Length) bytes)"
    } else {
        Compress-Archive -LiteralPath $Runner.OldPath -DestinationPath $Runner.BackupPath -Force
        Write-Log "Backup created: $($Runner.BackupPath)"
    }

    Move-Item -LiteralPath $Runner.OldPath -Destination $Runner.NewPath
    Write-Log "Moved $($Runner.OldPath) -> $($Runner.NewPath)"

    Retarget-RunnerJunctions -Root $Runner.NewPath

    $serviceExe = Join-Path $Runner.NewPath "bin\RunnerService.exe"
    if (-not (Test-Path -LiteralPath $serviceExe)) {
        throw "RunnerService.exe not found at $serviceExe"
    }
    $configOutput = & sc.exe config $Runner.ServiceName binPath= $serviceExe start= demand
    Write-Log "sc config output: $($configOutput -join ' | ')"

    Start-Service -Name $Runner.ServiceName
    Wait-ServiceState -Name $Runner.ServiceName -State "Running" -TimeoutSeconds 90
    Write-Log "Started service $($Runner.ServiceName)"

    Start-Sleep -Seconds 5
    $svcAfter = Get-CimInstance Win32_Service -Filter ("Name='{0}'" -f $Runner.ServiceName)
    if ($svcAfter.PathName -ne $serviceExe) {
        throw "Service path mismatch. Expected $serviceExe, got $($svcAfter.PathName)"
    }
    if ($svcAfter.StartMode -ne "Manual") {
        throw "Service start mode mismatch. Expected Manual, got $($svcAfter.StartMode)"
    }

    Assert-ProcessFromPath -Runner $Runner.Name -Root $Runner.NewPath
    Wait-GitHubRunnerOnline -Repo $Runner.Repo -Runner $Runner.Name
    Write-Log "GitHub validation passed for $($Runner.Repo) $($Runner.Name)"
}

Assert-Admin
Write-Log "Starting elevated Windows service runner migration."

$runners = @(
    [pscustomobject]@{
        Repo = "SGribanov/IdeaBox"
        Name = "ideabox-runner"
        ServiceName = "actions.runner.SGribanov-IdeaBox.ideabox-runner"
        OldPath = "C:\actions-runner-ideabox"
        NewPath = "C:\Runners\SGribanov-IdeaBox\ideabox-runner"
        BackupPath = "C:\Runners-backup\actions-runner-ideabox-ideabox-runner-move-2026-06-04.zip"
    },
    [pscustomobject]@{
        Repo = "SGribanov/DeltaG"
        Name = "deltag-win"
        ServiceName = "actions.runner.SGribanov-DeltaG.deltag-win"
        OldPath = "C:\github-runners\deltag"
        NewPath = "C:\Runners\SGribanov-DeltaG\deltag-win"
        BackupPath = "C:\Runners-backup\github-runners-deltag-deltag-win-move-2026-06-04.zip"
    }
)

try {
    foreach ($runner in $runners) {
        Invoke-RunnerMigration -Runner $runner
    }
    Write-Log "All elevated Windows service runner migrations completed."
} catch {
    Write-Log "ERROR: $($_.Exception.Message)"
    throw
}

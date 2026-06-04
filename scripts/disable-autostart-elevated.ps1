#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

$windowsServices = @(
    "actions.runner.SGribanov-DeltaG.deltag-win",
    "actions.runner.SGribanov-IdeaBox.ideabox-runner"
)

foreach ($service in $windowsServices) {
    sc.exe config $service start= demand | Write-Output
}

$wslServices = @(
    "actions.runner.SGribanov-DeltaG.deltag-linux-wsl.service",
    "actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service"
)

foreach ($service in $wslServices) {
    wsl.exe sudo systemctl disable $service
}

Write-Output "Autostart disable commands completed. Recheck service startup state before reboot."


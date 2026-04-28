[CmdletBinding()]
param(
    [string]$WzDirectory,
    [string]$ConverterPath,
    [string]$OutputPath
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($OutputPath)) {
    $OutputPath = Join-Path $repoRoot 'nx'
}

function Select-Folder([string]$Description) {
    try {
        Add-Type -AssemblyName System.Windows.Forms | Out-Null
        $dialog = New-Object System.Windows.Forms.FolderBrowserDialog
        $dialog.Description = $Description
        $dialog.ShowNewFolderButton = $false
        if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
            return $dialog.SelectedPath
        }
    }
    catch {
        Write-Warning "Could not open the folder picker automatically: $($_.Exception.Message)"
    }

    return $null
}

function Resolve-ConverterPath([string]$RequestedPath) {
    if (-not [string]::IsNullOrWhiteSpace($RequestedPath)) {
        return (Resolve-Path $RequestedPath).Path
    }

    $candidates = @(
        (Join-Path $PSScriptRoot 'go-wztonx-converter.exe'),
        (Join-Path $repoRoot 'go-wztonx-converter.exe'),
        (Join-Path $PSScriptRoot 'wztonx-converter.exe'),
        (Join-Path $repoRoot 'wztonx-converter.exe')
    )

    foreach ($candidate in $candidates) {
        if (Test-Path $candidate) {
            return (Resolve-Path $candidate).Path
        }
    }

    foreach ($commandName in @('go-wztonx-converter.exe', 'wztonx-converter.exe', 'go-wztonx-converter')) {
        $command = Get-Command $commandName -ErrorAction SilentlyContinue
        if ($command) {
            return $command.Source
        }
    }

    throw 'Could not find go-wztonx-converter. Put the .exe in the setup folder, make sure it is on PATH, or pass -ConverterPath.'
}

if ([string]::IsNullOrWhiteSpace($WzDirectory)) {
    $WzDirectory = Select-Folder 'Select your MapleStory v48 folder that contains the WZ files.'
}

if ([string]::IsNullOrWhiteSpace($WzDirectory)) {
    $WzDirectory = Read-Host 'Enter the full path to your MapleStory v48 folder'
}

if (-not (Test-Path $WzDirectory)) {
    throw "The MapleStory folder was not found: $WzDirectory"
}

$ConverterPath = Resolve-ConverterPath $ConverterPath

$expectedFiles = @(
    'Base.wz',
    'Character.wz',
    'Effect.wz',
    'Etc.wz',
    'Item.wz',
    'Map.wz',
    'Mob.wz',
    'Npc.wz',
    'Quest.wz',
    'Reactor.wz',
    'Skill.wz',
    'Sound.wz',
    'String.wz',
    'TamingMob.wz',
    'UI.wz'
)

$missing = @($expectedFiles | Where-Object { -not (Test-Path (Join-Path $WzDirectory $_)) })
if ($missing.Count -gt 0) {
    Write-Warning 'Some expected v48 WZ files were not found in the selected folder.'
    Write-Warning ('Missing: ' + ($missing -join ', '))
    Write-Warning 'The converter may still work if your files are stored differently, but make sure you selected the correct MapleStory folder.'
}

if (Test-Path $OutputPath) {
    $existingNxFiles = @(Get-ChildItem -Path $OutputPath -Filter '*.nx' -File -ErrorAction SilentlyContinue)
    if ($existingNxFiles.Count -gt 0) {
        $confirmation = Read-Host "The output folder already contains NX files. Continue and let the converter reuse or overwrite them? (Y/N)"
        if ($confirmation -notmatch '^(y|yes)$') {
            throw 'Conversion cancelled by user.'
        }
    }
}
else {
    New-Item -ItemType Directory -Path $OutputPath | Out-Null
}

$OutputPath = (Resolve-Path $OutputPath).Path

Write-Host ''
Write-Host 'Converter:' $ConverterPath
Write-Host 'MapleStory folder:' $WzDirectory
Write-Host 'Output folder:' $OutputPath
Write-Host ''
Write-Host 'Starting conversion... This can take a while.' -ForegroundColor Cyan

Push-Location $OutputPath
try {
    & $ConverterPath --server $WzDirectory
    if ($LASTEXITCODE -ne 0) {
        throw "Conversion failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

Write-Host ''
Write-Host 'Done.' -ForegroundColor Green
Write-Host 'Next steps:' -ForegroundColor Green
Write-Host '1. Edit config_dev.toml and set your MySQL password.'
Write-Host '2. Run: .\Valhalla.exe -type dev -config config_dev.toml'
if ($OutputPath -ne (Join-Path $repoRoot 'nx')) {
    Write-Host "3. Because you used a custom output path, run: .\Valhalla.exe -type dev -config config_dev.toml -nx \"$OutputPath\""
}

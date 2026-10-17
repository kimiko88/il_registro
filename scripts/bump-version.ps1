<#
.SYNOPSIS
    PowerShell script to bump versions across frontend and backend in Registrov2.
.DESCRIPTION
    Wrapper around scripts/bump-version.mjs with PowerShell argument handling.
.PARAMETER BumpType
    One of: patch, minor, major, beta, or an explicit version (e.g. 1.0.1).
.PARAMETER Git
    If specified, automatically creates a git commit and git tag.
.PARAMETER Check
    If specified, verifies synchronization across files without making changes.
.EXAMPLE
    .\scripts\bump-version.ps1 patch
    .\scripts\bump-version.ps1 minor -Git
    .\scripts\bump-version.ps1 -Check
#>
param (
    [Parameter(Position = 0)]
    [string]$BumpType = "patch",

    [switch]$Git,
    [switch]$Check
)

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$mjsScript = Join-Path $scriptDir "bump-version.mjs"

$nodeArgs = @()
if ($Check) {
    $nodeArgs += "--check"
} else {
    $nodeArgs += $BumpType
    if ($Git) {
        $nodeArgs += "--git"
    }
}

node $mjsScript @nodeArgs

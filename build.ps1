param(
    [ValidateSet("all", "windows", "linux", "linux-arm64", "darwin", "darwin-arm64")]
    [string]$Target = "all"
)

$AppName = "gh-proxy-go"
$DistDir = "dist"
$LdFlags = "-s -w"

# Ensure dist directory exists
if (-not (Test-Path $DistDir)) {
    New-Item -ItemType Directory -Path $DistDir -Force | Out-Null
}

# Get version info from git if available
$Version = "dev"
try {
    $gitTag = git describe --tags --abbrev=0 2>$null
    if ($gitTag) { $Version = $gitTag }
} catch {}

Write-Host "Building $AppName v$Version..." -ForegroundColor Cyan

$builds = @()

switch ($Target) {
    "all" {
        $builds = @(
            @{OS="windows"; Arch="amd64"; Suffix=".exe"}
            @{OS="linux";   Arch="amd64"; Suffix=""}
            @{OS="linux";   Arch="arm64"; Suffix=""}
            @{OS="darwin";  Arch="amd64"; Suffix=""}
            @{OS="darwin";  Arch="arm64"; Suffix=""}
        )
    }
    "windows" {
        $builds = @(@{OS="windows"; Arch="amd64"; Suffix=".exe"})
    }
    "linux" {
        $builds = @(@{OS="linux"; Arch="amd64"; Suffix=""})
    }
    "linux-arm64" {
        $builds = @(@{OS="linux"; Arch="arm64"; Suffix=""})
    }
    "darwin" {
        $builds = @(@{OS="darwin"; Arch="amd64"; Suffix=""})
    }
    "darwin-arm64" {
        $builds = @(@{OS="darwin"; Arch="arm64"; Suffix=""})
    }
}

foreach ($b in $builds) {
    $outputName = "${AppName}_$($b.OS)_$($b.Arch)$($b.Suffix)"
    $outputPath = Join-Path $DistDir $outputName

    Write-Host "  -> $outputName ... " -NoNewline

    $env:GOOS = $b.OS
    $env:GOARCH = $b.Arch
    $env:CGO_ENABLED = "0"

    $result = go build -ldflags "$LdFlags" -o $outputPath .
    if ($?) {
        $size = (Get-Item $outputPath).Length / 1KB
        Write-Host "OK ($([math]::Round($size, 1)) KB)" -ForegroundColor Green
    } else {
        Write-Host "FAILED" -ForegroundColor Red
    }
}

Remove-Item Env:\GOOS, Env:\GOARCH, Env:\CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host "`nAll builds placed in $DistDir/" -ForegroundColor Cyan

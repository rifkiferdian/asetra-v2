$ErrorActionPreference = 'Stop'

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$releaseRoot = Join-Path $projectRoot 'dist\asetra-windows-amd64-portable'
$releaseWeb = Join-Path $releaseRoot 'web'

Set-Location -LiteralPath $projectRoot
npm ci
npm run build:css
go test ./...

if (Test-Path -LiteralPath $releaseRoot) {
    $resolvedRelease = (Resolve-Path -LiteralPath $releaseRoot).Path
    $resolvedDist = (Resolve-Path -LiteralPath (Join-Path $projectRoot 'dist')).Path
    if (-not $resolvedRelease.StartsWith($resolvedDist + [IO.Path]::DirectorySeparatorChar)) {
        throw "Refusing to clean unexpected release path: $resolvedRelease"
    }
    Remove-Item -LiteralPath $releaseRoot -Recurse -Force
}

New-Item -ItemType Directory -Path $releaseWeb -Force | Out-Null
Copy-Item -LiteralPath (Join-Path $projectRoot 'web\templates') -Destination $releaseWeb -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'web\static') -Destination $releaseWeb -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot '.env.example') -Destination (Join-Path $releaseRoot '.env.example')
Copy-Item -LiteralPath (Join-Path $projectRoot 'database') -Destination $releaseRoot -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'migrations') -Destination $releaseRoot -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'seeds') -Destination $releaseRoot -Recurse

$vendorRoot = Join-Path $releaseRoot 'node_modules'
New-Item -ItemType Directory -Path (Join-Path $vendorRoot '@fortawesome\fontawesome-free') -Force | Out-Null
New-Item -ItemType Directory -Path (Join-Path $vendorRoot '@fontsource\roboto') -Force | Out-Null
New-Item -ItemType Directory -Path (Join-Path $vendorRoot 'sweetalert2') -Force | Out-Null
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fortawesome\fontawesome-free\css') -Destination (Join-Path $vendorRoot '@fortawesome\fontawesome-free') -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fortawesome\fontawesome-free\webfonts') -Destination (Join-Path $vendorRoot '@fortawesome\fontawesome-free') -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fontsource\roboto\400.css') -Destination (Join-Path $vendorRoot '@fontsource\roboto\400.css')
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fontsource\roboto\500.css') -Destination (Join-Path $vendorRoot '@fontsource\roboto\500.css')
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fontsource\roboto\700.css') -Destination (Join-Path $vendorRoot '@fontsource\roboto\700.css')
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\@fontsource\roboto\files') -Destination (Join-Path $vendorRoot '@fontsource\roboto') -Recurse
Copy-Item -LiteralPath (Join-Path $projectRoot 'node_modules\sweetalert2\dist') -Destination (Join-Path $vendorRoot 'sweetalert2') -Recurse

$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
go build -trimpath -ldflags '-s -w' -o (Join-Path $releaseRoot 'asetra.exe') ./cmd/web

$instructions = @'
ASETRA PORTABLE - WINDOWS 64 BIT

1. Salin .env.example menjadi .env.
2. Isi SESSION_SECRET (minimal 32 karakter) dan koneksi MySQL.
3. Siapkan database menggunakan database\asetra_starter.sql.
4. Jalankan asetra.exe.
5. Buka http://127.0.0.1:8080.

Folder web dan node_modules di paket ini wajib tetap berada di samping asetra.exe.
'@
Set-Content -LiteralPath (Join-Path $releaseRoot 'CARA-MENJALANKAN.txt') -Value $instructions -Encoding UTF8

$zipPath = Join-Path $projectRoot 'dist\asetra-windows-amd64-portable.zip'
if (Test-Path -LiteralPath $zipPath) { Remove-Item -LiteralPath $zipPath -Force }
Compress-Archive -LiteralPath $releaseRoot -DestinationPath $zipPath -CompressionLevel Optimal
Write-Host "Release created: $zipPath"

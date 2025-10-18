param(
    [Parameter(Mandatory=$true)]
    [string]$BaselineFile,
    
    [Parameter(Mandatory=$true)]
    [string]$CurrentFile,
    
    [Parameter(Mandatory=$true)]
    [int]$Threshold
)

if (-not (Test-Path $BaselineFile)) {
    Write-Error "Baseline file '$BaselineFile' not found"
    exit 1
}

if (-not (Test-Path $CurrentFile)) {
    Write-Error "Current file '$CurrentFile' not found"
    exit 1
}

function Extract-Benchmarks {
    param([string]$File)
    
    $results = @{}
    Get-Content $File | ForEach-Object {
        if ($_ -match '^(Benchmark\S+)\s+\d+\s+([\d.]+)\s+ns/op') {
            $name = $Matches[1]
            $nsPerOp = [double]$Matches[2]
            $results[$name] = $nsPerOp
        }
    }
    return $results
}

$baseline = Extract-Benchmarks -File $BaselineFile
$current = Extract-Benchmarks -File $CurrentFile

if ($baseline.Count -eq 0) {
    Write-Warning "No benchmark results found in baseline file"
    exit 0
}

if ($current.Count -eq 0) {
    Write-Warning "No benchmark results found in current file"
    exit 0
}

Write-Host "Benchmark Comparison Results:" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

$regressionsFound = $false

foreach ($benchName in $current.Keys) {
    $currentNs = $current[$benchName]
    
    if ($baseline.ContainsKey($benchName)) {
        $baselineNs = $baseline[$benchName]
        
        if ($baselineNs -ne 0) {
            $increase = (($currentNs - $baselineNs) / $baselineNs) * 100
            
            if ($increase -gt $Threshold) {
                Write-Host "REGRESSION: $benchName" -ForegroundColor Red
                Write-Host "   Baseline: $baselineNs ns/op"
                Write-Host "   Current:  $currentNs ns/op"
                Write-Host "   Increase: $([math]::Round($increase, 2))% (threshold: $Threshold%)"
                Write-Host ""
                $regressionsFound = $true
            }
            else {
                $changeText = if ($increase -ge 0) { "+$([math]::Round($increase, 2))%" } else { "$([math]::Round($increase, 2))%" }
                Write-Host "OK: $benchName ($changeText change)" -ForegroundColor Green
            }
        }
    }
    else {
        Write-Host "NEW: $benchName ($currentNs ns/op)" -ForegroundColor Yellow
    }
}

Write-Host ""
if ($regressionsFound) {
    Write-Host "Performance regressions detected! Build should fail." -ForegroundColor Red
    exit 1
}
else {
    Write-Host "No performance regressions detected." -ForegroundColor Green
    exit 0
}

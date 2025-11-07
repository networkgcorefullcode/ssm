# SSM API Test Script for PowerShell
# Usage: .\api_test.ps1 -Endpoint <endpoint> -Method <method> [-JsonFile <file>] [-BaseUrl <url>]
# Example: .\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"

param(
    [Parameter(Mandatory=$false, HelpMessage="Display help information")]
    [Alias("h")]
    [switch]$Help,

    [Parameter(Mandatory=$false, Position=0, HelpMessage="API endpoint (e.g., /encrypt, /decrypt)")]
    [string]$Endpoint,

    [Parameter(Mandatory=$false, Position=1, HelpMessage="HTTP method (GET, POST, PUT, DELETE)")]
    [ValidateSet("GET", "POST", "PUT", "DELETE", "PATCH")]
    [string]$Method = "POST",

    [Parameter(Mandatory=$false, Position=2, HelpMessage="Path to JSON file with request body")]
    [string]$JsonFile,

    [Parameter(Mandatory=$false, HelpMessage="Base URL of the API")]
    [string]$BaseUrl = "http://localhost:8080"
)

function Show-Help {
    Write-Host "`nSSM API Test Script - PowerShell Version" -ForegroundColor Cyan
    Write-Host "`nUsage:" -ForegroundColor Yellow
    Write-Host "  .\api_test.ps1 -Endpoint <endpoint> -Method <method> [-JsonFile <file>] [-BaseUrl <url>]"
    
    Write-Host "`nParameters:" -ForegroundColor Yellow
    Write-Host "  -Endpoint   API endpoint (e.g., /encrypt, /decrypt, /generate-aes-key)"
    Write-Host "  -Method     HTTP method (GET, POST, PUT, DELETE) - Default: POST"
    Write-Host "  -JsonFile   Path to JSON file with request body (optional)"
    Write-Host "  -BaseUrl    Base URL of the API - Default: http://localhost:8080"
    Write-Host "  -Help, -h   Show this help message"
    
    Write-Host "`nExamples:" -ForegroundColor Yellow
    Write-Host '  .\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"'
    Write-Host '  .\api_test.ps1 -Endpoint "/health-check" -Method "GET"'
    Write-Host '  .\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json" -BaseUrl "http://localhost:9000"'
    
    Write-Host "`nAvailable endpoints:" -ForegroundColor Yellow
    Write-Host "  /health-check         - Health check (GET/POST)"
    Write-Host "  /encrypt              - Encrypt data (POST)"
    Write-Host "  /decrypt              - Decrypt data (POST)"
    Write-Host "  /generate-aes-key     - Generate AES key (POST)"
    Write-Host "  /generate-des-key     - Generate DES key (POST)"
    Write-Host "  /generate-des3-key    - Generate DES3 key (POST)"
    Write-Host "  /store-key            - Store key (POST)"
    Write-Host "  /get-data-keys        - Get multiple keys by label (POST)"
    Write-Host "  /get-key              - Get single key by label (POST)"
    Write-Host "  /get-all-keys         - Get all keys (POST)"
    Write-Host ""
}

# Show help if requested or no parameters provided
if ($Help -or (-not $Endpoint)) {
    Show-Help
    exit 0
}

# Build the full URL
$FullUrl = "$BaseUrl$Endpoint"

Write-Host "`n=== SSM API Request ===" -ForegroundColor Green
Write-Host "URL:     " -ForegroundColor Yellow -NoNewline
Write-Host $FullUrl
Write-Host "Method:  " -ForegroundColor Yellow -NoNewline
Write-Host $Method

# Prepare headers
$Headers = @{
    "Content-Type" = "application/json"
    "Accept" = "application/json"
}

# Prepare body
$Body = $null
if ($JsonFile) {
    if (-not (Test-Path $JsonFile)) {
        Write-Host "`nError: JSON file not found: $JsonFile" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Body:    " -ForegroundColor Yellow -NoNewline
    Write-Host $JsonFile
    Write-Host "Content:" -ForegroundColor Yellow
    
    $Body = Get-Content $JsonFile -Raw
    
    # Pretty print JSON
    try {
        $JsonObject = $Body | ConvertFrom-Json
        $JsonObject | ConvertTo-Json -Depth 10 | Write-Host
    } catch {
        Write-Host $Body
    }
} else {
    Write-Host "Body:    " -ForegroundColor Yellow -NoNewline
    Write-Host "(none)"
}

Write-Host "`n=== Response ===" -ForegroundColor Green

# Measure request time
$StartTime = Get-Date

try {
    # Make the request
    $Response = if ($Body) {
        Invoke-RestMethod -Uri $FullUrl -Method $Method -Headers $Headers -Body $Body -ErrorAction Stop
    } else {
        Invoke-RestMethod -Uri $FullUrl -Method $Method -Headers $Headers -ErrorAction Stop
    }
    
    $EndTime = Get-Date
    $Duration = ($EndTime - $StartTime).TotalSeconds
    
    # Display response
    $Response | ConvertTo-Json -Depth 10 | Write-Host
    
    Write-Host "`nHTTP Status: 200" -ForegroundColor Green
    Write-Host "Time: $($Duration)s" -ForegroundColor Green
    Write-Host "`n✓ Request successful" -ForegroundColor Green
    
    exit 0
    
} catch {
    $EndTime = Get-Date
    $Duration = ($EndTime - $StartTime).TotalSeconds
    
    $StatusCode = $_.Exception.Response.StatusCode.value__
    $StatusDescription = $_.Exception.Response.StatusDescription
    
    Write-Host "`nError Response:" -ForegroundColor Red
    
    # Try to get error response body
    try {
        $ErrorStream = $_.Exception.Response.GetResponseStream()
        $Reader = New-Object System.IO.StreamReader($ErrorStream)
        $ErrorBody = $Reader.ReadToEnd()
        $Reader.Close()
        
        try {
            $ErrorJson = $ErrorBody | ConvertFrom-Json
            $ErrorJson | ConvertTo-Json -Depth 10 | Write-Host
        } catch {
            Write-Host $ErrorBody
        }
    } catch {
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
    
    Write-Host "`nHTTP Status: $StatusCode - $StatusDescription" -ForegroundColor Red
    Write-Host "Time: $($Duration)s" -ForegroundColor Yellow
    
    if ($StatusCode -ge 400 -and $StatusCode -lt 500) {
        Write-Host "`n⚠ Client error" -ForegroundColor Yellow
        exit 1
    } elseif ($StatusCode -ge 500) {
        Write-Host "`n✗ Server error" -ForegroundColor Red
        exit 1
    } else {
        Write-Host "`n? Unknown error" -ForegroundColor Yellow
        exit 1
    }
}

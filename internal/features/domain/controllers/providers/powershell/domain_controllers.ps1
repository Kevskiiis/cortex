param(
    [string]$Domain,
    [string]$Server
)

$ErrorActionPreference = 'Stop'

function Write-Result {
    param([string]$Status, [object[]]$Controllers, [object[]]$Checks, [object]$Error)
    [ordered]@{ status = $Status; controllers = @($Controllers); checks = @($Checks); error = $Error } | ConvertTo-Json -Depth 8 -Compress
    exit 0
}

function Write-Failure {
    param([string]$Status, [string]$Code, [string]$Message, [string]$Next, [string]$Detail)
    Write-Result -Status $Status -Controllers @() -Checks @() -Error ([ordered]@{ code = $Code; message = $Message; detail = $Detail; next = $Next })
}

try {
    $computer = Get-CimInstance -ClassName Win32_ComputerSystem
    if (-not $computer.PartOfDomain) { Write-Failure -Status 'not_domain_joined' -Code 'not_domain_joined' -Message 'This computer is not joined to an Active Directory domain.' -Next 'Run Cortex on a domain-joined computer.' -Detail 'Win32_ComputerSystem.PartOfDomain is false.' }
    $module = Get-Module -ListAvailable -Name ActiveDirectory | Select-Object -First 1
    if ($null -eq $module) { Write-Failure -Status 'rsat_missing' -Code 'rsat_missing' -Message 'The Active Directory PowerShell module is not installed.' -Next 'Install the RSAT Active Directory module and try again.' -Detail 'Get-Module did not find ActiveDirectory.' }
    Import-Module ActiveDirectory
    $checks = @(
        [ordered]@{ name = 'Domain joined'; status = 'Ready'; detail = [string]$computer.Domain; ok = $true },
        [ordered]@{ name = 'RSAT'; status = 'Ready'; detail = $module.Version.ToString(); ok = $true }
    )
    $query = @{ Filter = '*' }
    if ($Domain) { $query.Server = $Domain } elseif ($Server) { $query.Server = $Server }
    $controllers = @(Get-ADDomainController @query | Sort-Object HostName)
    $records = @($controllers | ForEach-Object { [ordered]@{ name = [string]$_.HostName; ipv4_address = [string]$_.IPv4Address; operating_system = [string]$_.OperatingSystem; site = [string]$_.Site; enabled = [bool]$_.Enabled; writable = -not [bool]$_.IsReadOnly; global_catalog = [bool]$_.IsGlobalCatalog } })
    if ($records.Count -eq 0) { Write-Result -Status 'empty' -Controllers @() -Checks $checks }
    $checks += [ordered]@{ name = 'Directory'; status = 'Ready'; detail = "$($records.Count) domain controller(s) found"; ok = $true }
    Write-Result -Status 'ready' -Controllers $records -Checks $checks
} catch {
    $message = $_.Exception.Message
    $lower = $message.ToLowerInvariant()
    if ($lower.Contains('access is denied') -or $lower.Contains('unauthorized')) { Write-Failure -Status 'access_denied' -Code 'access_denied' -Message 'The current account cannot query domain controllers.' -Next 'Use a domain-authenticated account with directory read access.' -Detail $message }
    elseif ($lower.Contains('server is not operational') -or $lower.Contains('cannot contact') -or $lower.Contains('network')) { Write-Failure -Status 'domain_unreachable' -Code 'domain_unreachable' -Message 'No domain controller could be reached.' -Next 'Connect to the domain network or VPN and try again.' -Detail $message }
    elseif ($lower.Contains('trust relationship')) { Write-Failure -Status 'secure_channel_failed' -Code 'secure_channel_failed' -Message 'The computer trust relationship with the domain is not healthy.' -Next 'Ask an administrator to repair the computer trust relationship.' -Detail $message }
    else { Write-Failure -Status 'unexpected_error' -Code 'unexpected_error' -Message 'The domain-controller query returned an unexpected error.' -Next 'Run the command again and use --json for technical details.' -Detail $message }
}
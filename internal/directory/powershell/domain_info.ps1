param(
    [string]$Domain,
    [string]$Server
)

$ErrorActionPreference = 'Stop'

function Write-Result {
    param(
        [string]$Status,
        [object]$Info,
        [object[]]$Checks,
        [object]$Error
    )

    [ordered]@{
        status = $Status
        info = $Info
        checks = @($Checks)
        error = $Error
    } | ConvertTo-Json -Depth 8 -Compress
    exit 0
}

function Write-Failure {
    param([string]$Status, [string]$Code, [string]$Message, [string]$Next, [string]$Detail)

    $checks = @(
        [ordered]@{ name = 'Platform'; status = 'Windows'; detail = 'Windows PowerShell is running'; ok = $true }
    )
    Write-Result -Status $Status -Checks $checks -Error ([ordered]@{
        code = $Code
        message = $Message
        detail = $Detail
        next = $Next
    })
}

try {
    $computer = Get-CimInstance -ClassName Win32_ComputerSystem
    if (-not $computer.PartOfDomain) {
        Write-Failure -Status 'not_domain_joined' -Code 'not_domain_joined' `
            -Message 'This computer is not joined to an Active Directory domain.' `
            -Next 'Run Cortex on a domain-joined computer.' -Detail 'Win32_ComputerSystem.PartOfDomain is false.'
    }

    $identity = [System.Security.Principal.WindowsIdentity]::GetCurrent()
    $userName = $identity.Name
    $userDomain = ($userName -split '\\')[0]
    $computerDomain = [string]$computer.Domain
    if ([string]::IsNullOrWhiteSpace($userDomain) -or $userDomain -ieq $computer.Name) {
        Write-Failure -Status 'not_domain_authenticated' -Code 'not_domain_authenticated' `
            -Message 'The current Windows user is not authenticated to the domain.' `
            -Next 'Sign in with a domain user or domain service account and try again.' -Detail $userName
    }

    $checks = @(
        [ordered]@{ name = 'Platform'; status = 'Ready'; detail = 'Windows PowerShell is running'; ok = $true },
        [ordered]@{ name = 'Domain joined'; status = 'Ready'; detail = $computerDomain; ok = $true },
        [ordered]@{ name = 'Identity'; status = 'Ready'; detail = $userName; ok = $true }
    )

    $module = Get-Module -ListAvailable -Name ActiveDirectory | Select-Object -First 1
    if ($null -eq $module) {
        Write-Failure -Status 'rsat_missing' -Code 'rsat_missing' `
            -Message 'The Active Directory PowerShell module is not installed.' `
            -Next 'Install the RSAT Active Directory module and try again.' -Detail 'Get-Module did not find ActiveDirectory.'
    }
    Import-Module ActiveDirectory
    $checks += [ordered]@{ name = 'RSAT'; status = 'Ready'; detail = $module.Version.ToString(); ok = $true }

    $secureChannel = Test-ComputerSecureChannel
    if (-not $secureChannel) {
        Write-Failure -Status 'secure_channel_failed' -Code 'secure_channel_failed' `
            -Message 'The computer trust relationship with the domain is not healthy.' `
            -Next 'Ask an administrator to repair the computer trust relationship.' -Detail 'Test-ComputerSecureChannel returned false.'
    }
    $checks += [ordered]@{ name = 'Secure channel'; status = 'Ready'; detail = 'Trust relationship is healthy'; ok = $true }

    if ($Server) {
        $domainInfo = Get-ADDomain -Server $Server
        $target = $Server
    } elseif ($Domain) {
        $domainInfo = Get-ADDomain -Identity $Domain
        $target = $Domain
    } else {
        $domainInfo = Get-ADDomain
        $target = $computerDomain
    }

    $checks += [ordered]@{ name = 'Directory'; status = 'Ready'; detail = 'Domain query succeeded'; ok = $true }
    $info = [ordered]@{
        dns_name = $domainInfo.DNSRoot
        netbios_name = $domainInfo.NetBIOSName
        distinguished_name = $domainInfo.DistinguishedName
        forest_name = $domainInfo.Forest
        current_user = $userName
        computer_name = $computer.Name
        computer_domain = $computerDomain
        domain_controllers = @($domainInfo.ReplicaDirectoryServers)
        target = $target
    }
    Write-Result -Status 'ready' -Info $info -Checks $checks
} catch {
    $message = $_.Exception.Message
    $lower = $message.ToLowerInvariant()
    if ($lower.Contains('access is denied') -or $lower.Contains('unauthorized')) {
        Write-Failure -Status 'access_denied' -Code 'access_denied' -Message 'The current account cannot query Active Directory.' `
            -Next 'Use a domain-authenticated account with directory read access.' -Detail $message
    } elseif ($lower.Contains('server is not operational') -or $lower.Contains('cannot contact') -or $lower.Contains('network')) {
        Write-Failure -Status 'domain_unreachable' -Code 'domain_unreachable' -Message 'No domain controller could be reached.' `
            -Next 'Connect to the domain network or VPN and try again.' -Detail $message
    } else {
        Write-Failure -Status 'unexpected_error' -Code 'unexpected_error' -Message 'The domain diagnostic returned an unexpected error.' `
            -Next 'Run the command again and use --json for technical details.' -Detail $message
    }
}
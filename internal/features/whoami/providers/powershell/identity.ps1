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

    Write-Result -Status $Status -Checks @() -Error ([ordered]@{
        code = $Code
        message = $Message
        detail = $Detail
        next = $Next
    })
}

try {
    $identity = [System.Security.Principal.WindowsIdentity]::GetCurrent()
    if ($null -eq $identity -or $null -eq $identity.User) {
        Write-Failure -Status 'identity_unavailable' -Code 'identity_unavailable' `
            -Message 'Windows did not provide a current user identity.' `
            -Next 'Run Cortex from a normal Windows user or service context.' -Detail 'WindowsIdentity.GetCurrent returned no usable identity.'
    }

    try {
        $computer = Get-CimInstance -ClassName Win32_ComputerSystem
    } catch {
        Write-Failure -Status 'system_info_unavailable' -Code 'system_info_unavailable' `
            -Message 'Cortex found the current identity but could not read computer information.' `
            -Next 'Check Windows system-management access and run the command again.' -Detail $_.Exception.Message
    }

    $name = $identity.Name
    $authority = ($name -split '\\')[0]
    $localAuthorities = @($computer.Name, 'BUILTIN', 'NT AUTHORITY', 'Window Manager', 'Font Driver Host')
    $isLocal = $localAuthorities -contains $authority
    $isDomainIdentity = [bool]$computer.PartOfDomain -and -not $isLocal
    $classification = if ($isDomainIdentity) { 'domain_identity' } else { 'local_identity' }
    $authentication = if ($isDomainIdentity) { 'domain' } else { 'local' }
    $status = if ($isDomainIdentity) { 'domain_identity' } else { 'local_identity' }

    $checks = @(
        [ordered]@{ name = 'Identity'; status = 'Ready'; detail = $name; ok = $true },
        [ordered]@{ name = 'Computer'; status = 'Ready'; detail = $computer.Name; ok = $true }
    )
    [ordered]@{
        status = $status
        info = [ordered]@{
            name = $name
            sid = $identity.User.Value
            authentication = $authentication
            computer_name = $computer.Name
            domain = [string]$computer.Domain
            domain_joined = [bool]$computer.PartOfDomain
            classification = $classification
        }
        checks = $checks
        error = $null
    } | ConvertTo-Json -Depth 8 -Compress
} catch {
    Write-Failure -Status 'unexpected_error' -Code 'unexpected_error' `
        -Message 'Cortex encountered an unexpected identity error.' `
        -Next 'Run the command again and use --json for technical details.' -Detail $_.Exception.Message
}
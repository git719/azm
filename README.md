# zman
Azure resource RBAC and MS Graph management utility. Uses the same supporting module as the `zls` utility (<https://github.com/git719/zls>), but it allows creation, update, and removal of certain objects.

## Getting Started
...

## Usage
```
zman Azure Resource RBAC and MS Graph MANAGER v0.8.4
    MANAGER FUNCTIONS
    -rm UUID|Specfile|"role name"     Delete role definition or assignment based on specifier
    -up Specfile                      Create or update definition or assignment based on specfile (YAML or JSON)
    -kd[j]                            Create a skeleton role-definition.yaml specfile (JSON option)
    -ka[j]                            Create a skeleton role-assignment.yaml specfile (JSON option)
    -spas SpUUID Expiry "name"        Add new secret to Service Principal, Expiry in YYYY-MM-DD format
    -sprs SpUUID SecretID             Remove secret from Service Principal
    -apas AppUUID Expiry "name"       Add new secret to application, Expiry in YYYY-MM-DD format
    -aprs AppUUID SecretID            Remove secret from application

    READER FUNCTIONS
    UUID                              Show object for given UUID
    -vs Specfile                      Compare YAML or JSON specfile to what's in Azure (only for d and a objects)
    -X[j] [Specifier]                 List all X objects tersely, with option for JSON output and/or match on Specifier
    -Xx                               Delete X object local file cache

      Where 'X' can be any of these object types:
      d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions
      m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups
      sp = Service Principals      ap = Applications            ad = Azure AD Roles

    -xx                               Delete ALL cache local files
    -ar                               List all RBAC role assignments with resolved names
    -mt                               List Management Group and subscriptions tree
    -pags                             List all Azure AD Privileged Access Groups
    -st                               List local cache count and Azure count of all objects

    -z                                Dump important program variables
    -cr                               Dump values in credentials file
    -cr  TenantId ClientId Secret     Set up MSAL automated ClientId + Secret login
    -cri TenantId Username            Set up MSAL interactive browser popup login
    -tx                               Delete MSAL accessTokens cache file
    -tc "TokenString"                 Dump token claims
    -v                                Print this usage page
```

## azm
`azm` is a [CLI](https://en.wikipedia.org/wiki/Command-line_interface) utility for managing [Indentity and Access Management (IAM)](https://www.nist.gov/identity-access-management) related Azure objects.


### Why?
Microsoft already has an official [Azure CLI tool](https://learn.microsoft.com/en-us/cli/azure/) called `az`, so **why** this?
- Only focuses on the smaller set of Azure objects that are related to IAM 
- Do quick and dirty searches of any IAM related object types in the azure tenant
- It is written in [Go](https://go.dev/) so it's much faster than `az`, or Python or PowerShell scripts
- Compiles into a binary executable that can easily be used from Windows, Linux, or macOS operating systems
- Automatic Linux [releases](https://github.com/git719/azm/releases/tag/v2.4.8) can easily be used in Github workflows
- Supports leveraging OIDC Github Action workflows with no passwords for configured Service Principal
- Primarily developed as a __proof-of-concept__ to learn to develop Azure utilities in the Go language
- Developed as part of a framework library for acquiring Azure [JWT](https://jwt.io/) token using the [MSAL library for Go](https://github.com/AzureAD/microsoft-authentication-library-for-go) (this utility leverages that library, [located here](https://github.com/git719/maz/blob/main/README.md))
- Quickly get a token to access the tenant's **Resources** Services API via endpoint <https://management.azure.com> ([REST API](https://learn.microsoft.com/en-us/rest/api/azure/))
- Quickly get a token to access the tenant's **Security** Services API via endpoint <https://graph.microsoft.com> ([MS Graph](https://learn.microsoft.com/en-us/graph/overview))

**PS:** Every repo should have a **Why?** entry at the top 😊


### Capabilities
This is a little _Swiss Army knife_ that can very quickly perform the following functions:
1. Read-Only Functions
    > **Note**<br>
    Of course these Read-Only functions are *only* available if you setup the tool to logon with an account that has the respective Read-Only privileges
    - List the following [Azure Resources Services](https://que.tips/azure/#azure-resource-services) objects in your tenant:
        - RBAC Role Definitions
        - RBAC Role Assignments
        - Azure Subscriptions
        - Azure Management Groups
    - List the following [Azure Security Services](https://que.tips/azure/#azure-security-services) objects:
        - Azure AD Users
        - Azure AD Groups
        - Applications
        - Service Principals
        - Azure AD Roles that have been **activated**
        - Azure AD Roles standard definitions
    - Compare RBAC role definitions and assignments that are defined in a YAML __specification file__ to what that object currently looks like in the Azure tenant
    - Dump the current Resources or Security JWT tokens (see **pman** below)
    - Perform *many* other related listing functions
2. Read-Write Functions
    > **Note**<br>
    Again, these Read-Write functions are **only** available if you setup the tool to logon with an account that has the respective Read-Write privileges
    - Delete/Create/Update the following [Azure Resources Services](https://que.tips/azure/#azure-resource-services) objects in your tenant:
        - RBAC Role Definitions
        - RBAC Role Assignments
    - Can output a sample RBAC Role definition or assignment YAML __specification file__, that can then be used to create a new role or assignment
    - Update the following [Azure Security Services](https://que.tips/azure/#azure-security-services) objects:
        - Service Principals: Can only create or delete SP secrets (Cannot yet create SPs)
        - Applications: Can only create or delete App secrets (Cannot yet create Apps)
    - Create a UUID
    - Other functions may be added


### Getting Started
To compile `azm`, first make sure you have installed and set up the Go language on your system. You can do that by following [these instructions here](https://que.tips/golang/#install-go-on-macos) or by following other similar recommendations found across the web.

- Also ensure that `$GOPATH/bin/` is in your `$PATH`, since that's where the executable binary will be placed.
- Open a `bash` shell, clone this repo, then switch to the `azm` working directory
- Type `./build` to build the binary executable
- To build from a *regular* Windows Command Prompt, just run the corresponding line in the `build` file (`go build ...`)
- If there are no errors, you should now be able to type `azm` and see the usage screen for this utility.

This utility has been successfully tested on macOS, Ubuntu Linux, and Windows. In Windows it works from a regular CMD.EXE, or PowerShell prompts, as well as from a GitBASH prompt.

Below other sections in this README explain how to set up access and use the utility in your own Azure tenant. 


### Access Requirements
First and foremost you need to know the special **Tenant ID** for your tenant. This is the unique UUID that identifies your Microsoft Azure tenant.

Then, you need a User ID or a Service Principal with the appropriate access rights. Either one will need the necessary privileges to perform the functions described above. For Read-Only functions that typically means getting _Reader_ role access to read **resource** objects, and _Global Reader_ role access to read **security** objects. The higher the scope of these access assignments, the more you will be able to see with the utility. 

When you run `azm` without any arguments you will see the **usage** screen listed below in this README. As you can quickly surmise, the `-id` argument will allow you to set up these 2 optional ways to connect to your tenant; either interactively with a User ID, also known as a [User Principal Name (UPN)](https://learn.microsoft.com/en-us/entra/identity/hybrid/connect/plan-connect-userprincipalname), or using a Service Principal or SP with a secret.

Another way of connecting is to use access tokens that have been acquired in another manner, perhaps via an OIDC login. 


#### User Logon
For example, if your Tenant ID was **c44154ad-6b37-4972-8067-0ef1068079b2**, and your User ID UPN was __bob@contoso.com__, you would type:

```
$ azm -id c44154ad-6b37-4972-8067-0ef1068079b2 bob@contoso.com
Updated /Users/myuser/.maz/credentials.yaml file
```
`azm` responds that the special `credentials.yaml` file has been updated accordingly.

To view, dump all configured logon values type the following:

```
$ azm -id
config_dir: /Users/myuser/.maz  # Config and cache directory
config_env_variables:
  # 1. MS Graph and Azure ARM tokens can be supplied directly via MAZ_MG_TOKEN and
  #    MAZ_AZ_TOKEN environment variables, and they have the highest precedence.
  #    Note, MAZ_TENANT_ID is still required when using these 2.
  # 2. Credentials supplied via environment variables have precedence over those
  #    provided via credentials file.
  # 3. The MAZ_USERNAME + MAZ_INTERACTIVE combo have priority over the MAZ_CLIENT_ID
  #    + MAZ_CLIENT_SECRET combination.
  MAZ_TENANT_ID:
  MAZ_USERNAME:
  MAZ_INTERACTIVE:
  MAZ_CLIENT_ID:
  MAZ_CLIENT_SECRET:
  MAZ_MG_TOKEN:
  MAZ_AZ_TOKEN:
config_creds_file:
  file_path: /Users/myuser/.maz/credentials.yaml
  tenant_id: c44154ad-6b37-4972-8067-0ef1068079b2
  username: bob@contoso.com
  interactive: true
```

Above tells you that the utility has been configured to use Bob's UPN for access via the special credentials file. Note that above is only a configuration setup, it actually hasn't logged Bob onto the tenant yet. To logon as Bob you have have to run any command, and the logon will happen automatically, in this case it will be an interactive browser popup logon.

Note also, that instead of setting up Bob's login with the `-id` argument, you could have setup the special 3 operating system environment variables to achieve the same. Had you done that, running `azm -id` would have displayed below instead:

```
$ azm -id
config_dir: /Users/myuser/.maz  # Config and cache directory
config_env_variables:
  # 1. Credentials supplied via environment variables override values provided via credentials file
  # 2. MAZ_USERNAME+MAZ_INTERACTIVE login have priority over MAZ_CLIENT_ID+MAZ_CLIENT_SECRET login
  MAZ_TENANT_ID: c44154ad-6b37-4972-8067-0ef1068079b2
  MAZ_USERNAME: bob@contoso.com
  MAZ_INTERACTIVE: true
  MAZ_CLIENT_ID:
  MAZ_CLIENT_SECRET:
  MAZ_MG_TOKEN:
  MAZ_AZ_TOKEN:
config_creds_file:
  file_path: /Users/myuser/.maz/credentials.yaml
  tenant_id: 
  username: 
  interactive:
```

#### SP Logon
To use an SP logon it means you first have to set up a dedicated App Registrations, grant it the same Reader resource and Global Reader security access roles mentioned above. For how to do an Azure App Registration please reference some other sources on the Web. By the way, this method is NOT RECOMMENDED, as you would be exposing the secret as an environment variables, which is not very good security practice.

Once above is setup, you then follow the same logic as for User ID logon above, but specifying 3 instead of 2 values; or use the other environment variables (MAZ_CLIENT_ID and MAZ_CLIENT_SECRET). 

The utility ensures that the permissions for configuration directory where the `credentials.yaml` file is only accessible by the owning user. However, storing a secrets in a clear-text file is a very poor security practice and should __never__ be use other than for quick tests, etc. The environment variable options was developed pricisely for this SP logon pattern, where the utility could be setup to run from say a [Docker container](https://en.wikipedia.org/wiki/Docker_(software)) and the secret injected as an environment variable, and that would be a much more secure way to run the utility.

An even better security practive when using the SP logon method is to leverage any process that can acquire OIDC tokens and make them available to this utility via the `MAZ_MG_TOKEN` and `MAZ_AZ_TOKEN` environment variable. If using OIDC logon, say for instance within a Github Workflow Action, you need to specify **both** these tokens and also the `MAZ_TENANT_ID` one.

(TODO: Need OIDC setup example, and how configure the SP on the Azure side with federated login.)

These login methods and the environment variables are described in more length in the [maz](https://github.com/git719/maz) package README.


### Quick Examples
1. List any Azure RBAC role, like the Built-in "Billing Reader" role for example:

```
$ azm -d "Billing Reader"
id: fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64
properties:
  roleName: Billing Reader
  description: Allows read access to billing data
  assignableScopes:
    - /
  permissions:
    - actions:
        - Microsoft.Authorization/*/read
        - Microsoft.Billing/*/read
        - Microsoft.Commerce/*/read
        - Microsoft.Consumption/*/read
        - Microsoft.Management/managementGroups/read
        - Microsoft.CostManagement/*/read
        - Microsoft.Support/*
      notActions:
      dataActions:
      notDataActions:
```

- Another way of listing the same role is to call it by its UUID: `azm -d fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64`
- The YAML listing format is more human-friendly and easier to read, and only displays the attributes that are most relevant to Azure systems engineers
- You can also display it in JSON format by calling: `azm -dj fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64`
- One advantage of the JSON format is that it displays every single attribute in the Azure object

2. Add a secret to an Application object: 

```
$ azm -apas 51afab9e-0225-4c36-81f0-f42289c1a57a "My Secret"
App_Object_Id: 51afab9e-0225-4c36-81f0-f42289c1a57a
New_Secret_Id: 7c140771-c547-43f9-8525-d08bd234e267
New_Secret_Name: My Secret
New_Secret_Expiry: 2025-01-06
New_Secret_Text: 8p68Q~Ab7OxR2nj.YOrXtOLwq1BT4bDy6wNebaYn
```

As the **usage** section shows, the secret Expiry defaults to 366 days if none is given. 

- Note that you have to use the **Objectd ID**, not the App ID (Client ID) of the application
- The name could have been nulled with `""`
- To remove above secret, you can simply do: `azm -aprs 51afab9e-0225-4c36-81f0-f42289c1a57a 7c140771-c547-43f9-8525-d08bd234e267`

3. Generate a random [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier). This can be very handy sometimes. Simply do: `azm -uuid`


### Usage
```
azm Azure IAM utility v2.4.8
    Read-Only Functions
    UUID                              Show object for given UUID
    -vs Specfile                      Compare YAML specfile to what's in Azure (only for d and a objects)
    -X[j] [Specifier]                 List all X objects tersely, with option for JSON output and/or match on Specifier
    -Xx                               Delete X object local file cache

      Where 'X' can be any of these object types:
      d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions
      m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups
      sp = Service Principals      ap = Applications            ad = Azure AD Roles

    Read-Write Functions
    -rm[f] UUID|Specfile|"role name"  Delete role definition or assignment based on specifier; force option
    -up[f] Specfile                   Create or update definition or assignment based on YAML specfile; force option
    -kd[j]                            Create a skeleton role-definition.yaml specfile
    -ka[j]                            Create a skeleton role-assignment.yaml specfile
    -spas SP_UUID "name" [Expiry]     Add secret to SP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -sprs SP_UUID SecretID            Remove secret from Service Principal
    -apas APP_UUID "name" [Expiry]    Add secret to APP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -aprs APP_UUID SecretID           Remove secret from application
    -uuid                             Generate new UUID

    -xx                               Delete ALL cache local files
    -ar                               List all RBAC role assignments with resolved names
    -mt                               List Management Group and subscriptions tree
    -pags                             List all Azure AD Privileged Access Groups
    -st                               List local cache count and Azure count of all objects
    -tmg                              Dump current token string for MS Graph API
    -taz                              Dump current token string for Azure Resource API
    -tc "TokenString"                 Dump token claims

    -id                               Display configured login values
    -id TenantId Username             Set up user for interactive login
    -id TenantId ClientId Secret      Set up ID for automated login
    -tx                               Delete current configured login values and token
    -v                                Print this usage page
```
Instead of documenting individual examples of all of the above switches, it is best for you to play around with the utility to see the different listing functionality that it offers.




### pman
[pman](https://github.com/git719/azm/tree/main/pman) is a poor man's REST API Postman BASH script, which leverages `azm`'s `-tmg` and `-taz` arguments to get the current user's token to make other generic REST API calls against those 2 Microsoft APIs.


### Token Sample Code
Included in this repo are examples of how to get a Microsoft token with 3 different languages. The token can be for any API. The examples use Docker Compose.
1. [Python Example](https://github.com/git719/azm/tree/main/token-python)
2. [PowerShell Example](https://github.com/git719/azm/tree/main/token-powershell)
3. [Node Example](https://github.com/git719/azm/tree/main/token-node)


### Container
There is also a Docker Compose file and a Dockerfile for an example of how to use this program within a container.


### To-Do and Known Issues
The program is stable enough to be relied on as a small utility, but there are a number of little niggly things that could be improved. Will put a list together at some point.

At any rate, no matter how stable any code is, it is always worth remembering computer scientist [Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)'s famous quote:
> "Inside every large program is a small program struggling to get out."


### Coding Philosophy and Feedback
The primary goal of this utility is to serve as a study aid for coding Azure utilities in Go, as well as to serve as a quick _Swiss Army knife* utility for managin tenant IAM objects. If you look through the code I think you will find that is relatively straightforward. There is a deliberate effor to keep the code as clear as possible, and simple to understand and maintain.

Note that the bulk of the code is actually in the [maz](https://github.com/git719/maz) library, and other packages. Please visit that repo for more info.

This utility along with the required libraries are obviously very useful to me. I don't think I'm going to formalize the way to contribute to this project, but if you find it useful and have some improvement suggestion please let me know. Anyway, this is published as an open source project, so feel free to clone and use on your own, with proper attributions.

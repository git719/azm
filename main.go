package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/queone/maz"
	"github.com/queone/utl"
)

const (
	prgname = "azm"
	prgver  = "2.6.1"
)

// Prints program usage
func printUsage() {
	X := utl.Red("X")
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Azure IAM CLI utility - https://github.com/git719/azm\n" +
		"Usage: " + prgname + " [options]\n" +
		"  Read-Only Functions\n" +
		"  UUID                              Show object(s) related to given UUID\n" +
		"  -vs Specfile                      Compare YAML specfile to what's in Azure (only for d and a objects)\n" +
		"  -" + X + "[j] [Specifier]                 List all " + X + " objects tersely, with option for JSON output and/or match on Specifier\n\n" +
		"    Where '" + X + "' can be any of these object types:\n" +
		"    d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions  \n" +
		"    m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups      \n" +
		"    sp = Service Principals      ap = Applications            ad = Azure AD Roles\n" +
		"\n" +
		"  -ar                               List all RBAC role assignments with resolved names\n" +
		"  -mt                               List Management Group and subscriptions tree\n" +
		"  -pags                             List all Azure AD Privileged Access Groups\n" +
		"  -st                               List local cache count and Azure count of all objects\n" +
		"  -tmg                              Dump current token string for MS Graph API\n" +
		"  -taz                              Dump current token string for Azure Resource API\n" +
		"  -tc \"TokenString\"                 Dump token claims\n" +
		"\n" +
		"  Read-Write Functions\n" +
		"  -rm[f] UUID|Specfile|\"role name\"  Delete role definition or assignment based on specifier; force option\n" +
		"  -up[f] Specfile                   Create or update definition or assignment based on YAML specfile; force option\n" +
		"  -kd[j]                            Create a skeleton role-definition.yaml specfile\n" +
		"  -ka[j]                            Create a skeleton role-assignment.yaml specfile\n" +
		"  -spas SP_UUID \"name\" [Expiry]     Add secret to SP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)\n" +
		"  -sprs SP_UUID SecretID            Remove secret from Service Principal\n" +
		"  -apas APP_UUID \"name\" [Expiry]    Add secret to APP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)\n" +
		"  -aprs APP_UUID SecretID           Remove secret from application\n" +
		"  -uuid                             Generate new UUID\n" +
		"\n" +
		"  Housekeeping Functions\n" +
		"  -" + X + "x                               Delete " + X + " object local file cache\n" +
		"  -xx                               Delete ALL cache local files\n" +
		"  -id                               Display configured login values\n" +
		"  -id TenantId Username             Set up user for interactive login\n" +
		"  -id TenantId ClientId Secret      Set up ID for automated login\n" +
		"  -tx                               Delete current configured login values and token\n" +
		"  -?, -h, --help                    Print this usage page\n")
	os.Exit(0)
}

// Sets up variable bundle object struct
func setupVariables(z *maz.Bundle) maz.Bundle {
	*z = maz.Bundle{
		ConfDir:      "",
		CredsFile:    "credentials.yaml",
		TokenFile:    "accessTokens.json",
		TenantId:     "",
		ClientId:     "",
		ClientSecret: "",
		Interactive:  false,
		Username:     "",
		AuthorityUrl: "",
		MgToken:      "",
		MgHeaders:    map[string]string{},
		AzToken:      "",
		AzHeaders:    map[string]string{},
	}
	// Set up configuration directory
	z.ConfDir = filepath.Join(os.Getenv("HOME"), ".maz") // IMPORTANT: Setting config dir = "~/.maz"
	if utl.FileNotExist(z.ConfDir) {
		if err := os.Mkdir(z.ConfDir, 0700); err != nil {
			panic(err.Error())
		}
	}
	return *z
}

func main() {
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 4 {
		printUsage() // Don't accept less than 1 or more than 4 arguments
	}

	var z maz.Bundle
	z = setupVariables(&z)

	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		// This first set of 1-arg requests do not require API tokens to be set up
		switch arg1 {
		case "-?", "-h", "--help":
			printUsage()
		case "-id":
			maz.DumpLoginValues(z)
		case "-tx":
			maz.RemoveCacheFile("t", z)
			maz.RemoveCacheFile("id", z)
		case "-uuid":
			utl.Die(uuid.New().String() + "\n")
		}
		z = maz.SetupApiTokens(&z) // The remaining 1-arg requests DO required API tokens to be set up
		switch arg1 {
		case "-xx":
			maz.RemoveCacheFile("all", z)
		case "-dx", "-ax", "-sx", "-mx", "-ux", "-gx", "-spx", "-apx", "-adx":
			t := arg1[1 : len(arg1)-1] // Single out the object type (t, d, sp, etc)
			maz.RemoveCacheFile(t, z)
		case "-dj", "-aj", "-sj", "-mj", "-uj", "-gj", "-spj", "-apj", "-adj":
			t := arg1[1 : len(arg1)-1]
			all := maz.GetObjects(t, "", false, z) // false means do not force Azure call, ok to use cache
			utl.PrintJsonColor(all)                // Print entire set in JSON
		case "-d", "-a", "-s", "-m", "-u", "-g", "-sp", "-ap", "-ad":
			t := arg1[1:]
			all := maz.GetObjects(t, "", false, z)
			for _, i := range all { // Print entire set tersely
				maz.PrintTersely(t, i)
			}
		case "-ar":
			maz.PrintRoleAssignmentReport(z)
		case "-mt":
			maz.PrintMgTree(z)
		case "-pags":
			maz.PrintPags(z)
		case "-st":
			maz.PrintCountStatus(z)
		case "-kd", "-kdj", "-ka", "-kaj":
			t := arg1[2:] // Single out the type (d, dj, a, aj)
			maz.CreateSkeletonFile(t)
		case "-tmg":
			fmt.Println(z.MgToken)
		case "-taz":
			fmt.Println(z.AzToken)
		default:
			c := rune(arg1[0])                            // Grab 1st charater of string, to check if it's hex
			if utl.IsHexDigit(c) && utl.ValidUuid(arg1) { // If valid UUID, search/print matching object(s?)
				maz.PrintObjectByUuid(arg1, z)
			} else {
				printUsage()
			}
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		z = maz.SetupApiTokens(&z)
		switch arg1 {
		case "-rm", "-rmf":
			force := false
			if arg1 == "-rmf" {
				force = true
			}
			maz.DeleteAzObject(force, arg2, z)
		case "-up", "-upf":
			force := false
			if arg1 == "-upf" {
				force = true
			}
			maz.UpsertAzObject(force, arg2, z)
		case "-tc":
			maz.DecodeJwtToken(arg2)
		case "-vs":
			maz.CompareSpecfileToAzure(arg2, z)
		case "-dj", "-aj", "-sj", "-mj", "-uj", "-gj", "-spj", "-apj", "-adj":
			t := arg1[1 : len(arg1)-1] // Single out the object type
			maz.PrintMatching("json", t, arg2, z)
		case "-d", "-a", "-s", "-m", "-u", "-g", "-sp", "-ap", "-ad":
			t := arg1[1:] // Single out the object type
			maz.PrintMatching("reg", t, arg2, z)
		default:
			printUsage()
		}
	case 3: // Process 3-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		switch arg1 {
		case "-id":
			z.TenantId = arg2
			z.Username = arg3
			maz.SetupInterativeLogin(z)
		}
		z = maz.SetupApiTokens(&z) // The remaining 3-arg requests do API access
		switch arg1 {
		case "-sprs":
			maz.RemoveSpSecret(arg2, arg3, z)
		case "-aprs":
			maz.RemoveAppSecret(arg2, arg3, z)
		case "-spas":
			expiryTime := utl.GetDateInDays("365")    // Set expiry to one year from now
			expiry := expiryTime.Format("2006-01-02") // Convert to yyyy-mm-dd format
			maz.AddSpSecret(arg2, arg3, expiry, z)
		case "-apas":
			expiryTime := utl.GetDateInDays("365")
			expiry := expiryTime.Format("2006-01-02")
			maz.AddAppSecret(arg2, arg3, expiry, z)
		default:
			printUsage()
		}
	case 4: // Process 4-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		arg4 := os.Args[4]
		switch arg1 {
		case "-id":
			z.TenantId = arg2
			z.ClientId = arg3
			z.ClientSecret = arg4
			maz.SetupAutomatedLogin(z)
		}
		z = maz.SetupApiTokens(&z) // The remaining 4-arg requests do API access
		switch arg1 {
		case "-spas":
			maz.AddSpSecret(arg2, arg3, arg4, z)
		case "-apas":
			maz.AddAppSecret(arg2, arg3, arg4, z)
		default:
			printUsage()
		}
	}
	os.Exit(0)
}

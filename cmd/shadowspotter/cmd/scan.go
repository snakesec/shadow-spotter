package cmd

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
	"math/rand"
	
	"gitlab.com/snake-security/shadowspotter/internal/scan"
	"gitlab.com/snake-security/shadowspotter/pkg/context"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"
)

func getRandomUserAgent() string {
	userAgents := []string{
		"Chrome. Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
		"Firefox. Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2.1 Safari/605.1.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.",
		"Mozilla/5.0 (Linux; arm_64; Android 10; LRA-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.94 YaBrowser/23.9.3.94.00 SA/3 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 13; CPH2339) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.280 Mobile Safari/537.36 OPR/80.3.4244.77596",
		"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.6167.144 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/118.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Mobile/15E148 Snapchat/10.77.5.59 (like Safari/604.1)",
		"Mozilla/5.0 (Linux; Android 13; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36 uacq",
		"Mozilla/5.0 (Android 13; Mobile; rv:109.0) Gecko/114.0 Firefox/114.0",
	}

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(userAgents))
	return userAgents[randomIndex]
}

var (
	kitebuilderFiles    = []string{}
	kitebuilderFullScan = false
	headers             = []string{}

	failStatusCodes    = []int{}
	successStatusCodes = []int{100,200,300,301,302,401,403,405,500}
	lengthIgnoreRange  = []string{}

	progressBar               = true
	disablePrecheck           = false
	wildcardDetection         = true
	maxConnPerHost            = 3
	maxParallelHosts          = 50
	delay                     = 1 * time.Second
	userAgent                 = "" //"Chrome. Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36"
	quarantineThreshold int64 = 10

	timeout      = 3 * time.Second
	maxRedirects = 3

	preflightDepth   int64 = 1
	blacklistDomains       = []string{}
	filterAPIs             = []string{}

	forceMethod = ""

	userQuery   = ""

	allowUnsafe = false

	profileName = ""

	assetnoteWordlist = []string{}

)



// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan INPUT [ -w wordlist.kite ]",
	Short: "scan one or multiple hosts with a provided wordlist",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `this will perform a concurrent scan of one or multiple hosts
using a wordlist.
We will attempt to find a file matching the <input> you provide, and otherwise attempt to parse it as a URI. 
If protocol is missing, then we will assume from the port, if the port is missing, then we will try both http:80 and https:443

The kitebuilder file format is a modified OpenAPI schema that allows you to specify arguments, parameters, headers, methods and body structure for structured api calls.

By default, we perform a 2 phase scan. The first phase uses a single route for API. If any of the routes respond, we perform a second phase scan on the host where all the routes for an API are scanned.

usage: 
Shadow-Spotter scan <input> <flags>
Shadow-Spotter scan hosts.txt -A=apiroutes-210228:5000 
Shadow-Spotter scan domain.com -w wordlist.kite

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		opts := []scan.ScanOption{
			scan.MaxParallelHosts(maxParallelHosts),
			scan.MaxConnPerHost(maxConnPerHost),
			scan.MaxRedirects(maxRedirects),
			scan.ContentLengthIgnoreRanges(lengthIgnoreRange),
			scan.Timeout(timeout),
			scan.Delay(delay),
			scan.AddHeaders(headers),
			scan.AllowUnsafe(allowUnsafe),
			scan.ForceMethod(forceMethod),
			scan.UserQuery(userQuery),
			scan.UserAgent(userAgent),
			scan.SuccessStatusCodes(successStatusCodes),
			scan.FailStatusCodes(failStatusCodes),
			scan.BlacklistDomains(blacklistDomains),
			scan.FilterAPIs(filterAPIs),
			scan.WildcardDetection(wildcardDetection),
			scan.ProgressBarEnabled(progressBar),
			scan.QuarantineThreshold(quarantineThreshold),
			scan.PreflightDepth(preflightDepth),
			scan.Precheck(!disablePrecheck),
			scan.LoadKitebuilderFile(kitebuilderFiles),
			scan.KitebuilderFullScan(kitebuilderFullScan),
			scan.LoadAssetnoteWordlistKitebuilder(assetnoteWordlist),
		}

		if len(userAgent) == 0 {
			opts = append(opts, scan.UserAgent(getRandomUserAgent()))
		} else {
			opts = append(opts, scan.UserAgent(userAgent))
		}

		go func() {
			log.Debug().Err(http.ListenAndServe("localhost:6060", nil)).Msg("Started http profiler server")
		}()

		if profileName != "" {
			f, err := os.Create(profileName)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create profile")
			}

			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		if domain == "-" {
			if err := scan.ScanStdin(context.Context(), opts...); err != nil {
				log.Fatal().Err(err).Msg("failed to read from stdin")
			}
		} else {
			if err := scan.ScanDomainOrFile(context.Context(), domain, opts...); err != nil {
				log.Fatal().Err(err).Msg("failed to scan domain")
			}
		}
		
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringSliceVarP(&kitebuilderFiles, "kitebuilder-list", "w", kitebuilderFiles, "ogl wordlist to use for scanning")
	scanCmd.Flags().BoolVar(&kitebuilderFullScan, "kitebuilder-full-scan", kitebuilderFullScan, "perform a full scan without first performing a phase scan.")
	scanCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{"x-forwarded-for: 127.0.0.1"}, "headers to add to requests")

	scanCmd.Flags().BoolVar(&disablePrecheck, "disable-precheck", false, "whether to skip host discovery")

	scanCmd.Flags().IntVarP(&maxConnPerHost, "max-connection-per-host", "x", maxConnPerHost, "max connections to a single host")
	scanCmd.Flags().IntVarP(&maxParallelHosts, "max-parallel-hosts", "j", maxParallelHosts, "max number of concurrent hosts to scan at once")
	scanCmd.Flags().DurationVar(&delay, "delay", delay, "delay to place inbetween requests to a single host")
	scanCmd.Flags().StringVar(&userAgent, "user-agent", "", "user agent to use for requests (default Random)")
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "t", timeout, "timeout to use on all requests")
	scanCmd.Flags().IntVar(&maxRedirects, "max-redirects", maxRedirects, "maximum number of redirects to follow")
	scanCmd.Flags().StringVar(&forceMethod, "force-method", forceMethod, "whether to ignore the methods specified in the ogl file and force this method")
	scanCmd.Flags().StringVar(&userQuery, "custom-query", userQuery, "Forces a query to be added to the URL path")
	scanCmd.Flags().BoolVar(&allowUnsafe, "allow-unsafe", allowUnsafe, "allow run unsafe methods (PUT, DELETE, PATCH)")

	scanCmd.Flags().IntSliceVar(&successStatusCodes, "success-status-codes", successStatusCodes,
		"which status codes whitelist as success. this is the default mode")
	scanCmd.Flags().IntSliceVar(&failStatusCodes, "fail-status-codes", failStatusCodes,
		"which status codes blacklist as fail. if this is set, this will override success-status-codes")

	scanCmd.Flags().StringSliceVar(&blacklistDomains, "blacklist-domain", blacklistDomains, "domains that are blacklisted for redirects. We will not follow redirects to these domains")
	scanCmd.Flags().BoolVar(&wildcardDetection, "wildcard-detection", wildcardDetection, "can be set to false to disable wildcard redirect detection")
	scanCmd.Flags().StringSliceVar(&lengthIgnoreRange, "ignore-length", lengthIgnoreRange, "a range of content length bytes to ignore. you can have multiple. e.g. 100-105 or 1234 or 123,34-53. This is inclusive on both ends")

	scanCmd.Flags().BoolVar(&progressBar, "progress", progressBar, "a progress bar while scanning. by default enabled only on Stderr")
	scanCmd.Flags().Int64Var(&quarantineThreshold, "quarantine-threshold", quarantineThreshold, "if the host return N consecutive hits, we quarantine the host as wildcard. Set to 0 to disable")

	scanCmd.Flags().Int64VarP(&preflightDepth, "preflight-depth", "d", 1, "when performing preflight checks, what directory depth do we attempt to check. 0 means that only the docroot is checked")
	scanCmd.Flags().StringVar(&profileName, "profile-name", profileName, "name for profile output file")

	scanCmd.Flags().StringSliceVar(&filterAPIs, "filter-api", filterAPIs, "only scan apis matching this ksuid")

	scanCmd.Flags().StringSliceVarP(&assetnoteWordlist, "assetnote-wordlist", "A", assetnoteWordlist, "use the wordlists from wordlists.assetnote.io. specify the type/name to use, e.g. apiroutes-210228. You can specify an additional maxlength to use only the first N values in the wordlist, e.g. apiroutes-210228;20000 will only use the first 20000 lines in that wordlist")

}

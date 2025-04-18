/*
Shadow-Spotter Next Gen Content Discovery
Copyright (C) 2024  Weidsom Nascimento - SNAKE Security

Based on kiterunner from AssetNote

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"regexp"
	"time"
	"math/rand"

	"gitlab.com/snake-security/shadowspotter/internal/kitebuilder"
	"gitlab.com/snake-security/shadowspotter/pkg/context"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"
)

func getRandomUserAgentR() string {
	userAgentRs := []string{
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

	randomIndex := rand.Intn(len(userAgentRs))
	return userAgentRs[randomIndex]
}

var (
	kitebuilderFile = ""
	proxy = ""
	userHeaders = []string{}
	userAgentR = ""
)

// replayCmd represents the replay command
var replayCmd = &cobra.Command{
	Use:   `replay "GET [   5,   2,   1] /foo -> /bar 0cc39f80913b550423454792b47ba661e6724a59" -w routes.kite`,
	Short: "replay a kitebuilder request based on the input",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `replay an kitebuilder request based on the input

Provide the raw Shadow-Spotter input and we'll figure it out for you.
or 
Shadow-Spotter kb replay -w routes.kite <id> <method> <route> <host>
Shadow-Spotter kb replay -w routes.kite "<full line output>"

e.g.
Shadow-Spotter kb replay -w routes.kite 0cc39f80913b550423454792b47ba661e6724a59 GET /foo http://127.0.0.1:5000
Shadow-Spotter kb replay -w routes.kite "POST    400 [    138,    5,  11] https://ap-service-team-hm.services.atlassian.com/volumes/create 0cc39f830ee6b0093e824073fd086bbd7c34b631"
Shadow-Spotter kb replay -w routes.kite --proxy http://localhost:8080 "POST    403 [    126,   25,   6] https://artifactory.services.atlassian.com/REPORT_AVAILABLE 0cc39f84f1c74d86ceb7727823cd1cc0f996ea19" 
`,
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ksuid string
			method string
			path string
			host string
		)

		if len(userAgentR) == 0 {
			userAgentR = getRandomUserAgentR()
		}

		if len(args) == 1 {
			rg := regexp.MustCompile(`(\w+).*\[.*\] (https?://[^/]+)([^ ]+).* (.*)$`)
			matches := rg.FindStringSubmatch(args[0])
			ksuid = matches[4]
			method = matches[1]
			path = matches[3]
			host = matches[2]
		} else if len(args) == 4 {
			ksuid = args[0]
			method = args[1]
			path = args[2]
			host = args[3]
		}

		if err := kitebuilder.Replay(context.Context(), kitebuilderFile, ksuid, method, path, host, proxy, userHeaders, userAgentR); err != nil {
			log.Fatal().Err(err).Msg("failed to replay request")
		}
	},
}

func init() {
	kidebuilderCmd.AddCommand(replayCmd)

	replayCmd.Flags().StringVarP(&kitebuilderFile, "kitebuilder-list", "w", kitebuilderFile, "ogl wordlist to use for scanning")
	replayCmd.Flags().StringVarP(&proxy, "proxy", "p", proxy, "proxy to replay the request through")
	replayCmd.Flags().StringSliceVarP(&userHeaders, "header", "H", []string{"x-forwarded-for: 127.0.0.1"}, "headers to add to requests")
	replayCmd.Flags().StringVar(&userAgentR, "user-agent", "", "user agent to use for requests (default Random)")
}

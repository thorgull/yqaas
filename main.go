/*
YQ As A Service
Copyright (C) 2024 Thorgull

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorgull/yqaas/gen/api"
	"github.com/thorgull/yqaas/impl"
	"gopkg.in/op/go-logging.v1"
	"runtime/debug"
	"strings"

	//"log"
	"net/http"
)

var log = logging.MustGetLogger("yqaas")

func findYQVersion() (string, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Warningf("can not dertermine the version of the yq library, no build info available")
		return "", false
	}
	for _, bi := range bi.Deps {
		if strings.HasPrefix(bi.Path, "github.com/mikefarah/yq") {
			return bi.Version, true
		}
	}
	log.Warningf("can not dertermine the version of the yq library, module github.com/mikefarah/yq/* not found")
	return "", false
}

func getVersionInfo() map[string]string {
	var data = make(map[string]string)
	if yqVersion, ok := findYQVersion(); ok {
		data["yq"] = yqVersion
	}
	return data
}

type buildInfo struct {
	Versions map[string]string `json:"versions"`
}

func getBuildInfo() buildInfo {
	return buildInfo{
		Versions: getVersionInfo(),
	}
}

func main() {
	metrics := flag.Bool("prometheus", false, "Enabled /metrics endpoint")
	probes := flag.Bool("probes", false, "Enable /health/* endpoints")
	port := flag.Int("port", 8080, "Configure port")
	verbose := flag.Bool("verbose", false, "Show debug logs")
	openapi := flag.Bool("openapi", false, "Enable /openapi endpoint")
	flag.Parse()

	if *verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}

	log.Info("Server starting...")

	DefaultApiService := impl.NewDefaultAPIService()
	DefaultApiController := api.NewDefaultAPIController(DefaultApiService)

	router := api.NewRouter(DefaultApiController)

	log.Info("[✔️] Enable /info endpoint")
	buildInfoData := getBuildInfo()
	router.HandleFunc("/buildInfo", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := json.Marshal(buildInfoData)
		if err != nil {
			writer.WriteHeader(500)
		} else {
			writer.WriteHeader(200)
			_, err = writer.Write(bs)
		}
	})

	if *metrics {
		log.Info("[✔️] Enable /metrics endpoint")
		router.Handle("/metrics", promhttp.Handler())
	}
	if *probes {
		log.Info("[✔️] Enable /health/* endpoints")
		respondNoContent := func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNoContent)
		}
		router.HandleFunc("/health/live", respondNoContent)
		router.HandleFunc("/health/ready", respondNoContent)
	}
	if *openapi {
		log.Info("[✔️] Enable /openapi endpoint")

		router.HandleFunc("/openapi", func(writer http.ResponseWriter, request *http.Request) {
			http.ServeFile(writer, request, "yqaas.yaml")
		})
	}

	log.Infof("Listening on %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}

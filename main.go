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
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorgull/yqaas/gen/api"
	"github.com/thorgull/yqaas/impl"
	"log"
	"net/http"
)

func main() {
	metrics := flag.Bool("prometheus", false, "Enabled /metrics endpoint")
	probes := flag.Bool("probes", false, "Enable /health/* endpoints")
	flag.Parse()
	log.Printf("Server started")

	DefaultApiService := impl.NewDefaultAPIService()
	DefaultApiController := api.NewDefaultAPIController(DefaultApiService)

	router := api.NewRouter(DefaultApiController)
	if *metrics {
		log.Printf("-- Enable /metrics endpoint")
		router.Handle("/metrics", promhttp.Handler())
	}
	if *probes {
		log.Printf("-- Enable /health/* endpoints")
		respondNoContent := func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNoContent)
		}
		router.HandleFunc("/health/live", respondNoContent)
		router.HandleFunc("/health/ready", respondNoContent)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}

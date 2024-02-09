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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorgull/yqaas/gen/api"
	"github.com/thorgull/yqaas/impl"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	DefaultApiService := impl.NewDefaultAPIService()
	DefaultApiController := api.NewDefaultAPIController(DefaultApiService)

	router := api.NewRouter(DefaultApiController)
	router.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", router))
}

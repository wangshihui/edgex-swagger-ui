//
// Copyright (c) 2021 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"github.com/gorilla/mux"
	"thundersoft.com/edgex/swagger-ui/internal"
)

//
func main() {
	//_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	ctx, cancel := context.WithCancel(context.Background())
	internal.Main(ctx, cancel, mux.NewRouter())
	//
	//os.Mkdir("file", 0777)
	//http.Handle("/edgex-swagger-ui/", http.StripPrefix("/edgex-swagger-ui/", http.FileServer(http.Dir("E:\\ts\\project\\cmcc_sher\\second_phase\\edgex-swagger-ui\\swagger-ui"))))
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
}

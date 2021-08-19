package interx

import (
	"net/http"
	"sort"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/global"
	"github.com/KiraCore/sekai/INTERX/tasks"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterNodeListQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryPubP2PList, QueryPubP2PNodeList(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryPrivP2PList, QueryPrivP2PNodeList(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryInterxList, QueryInterxList(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QuerySnapList, QuerySnapList(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryPubP2PList, "This is an API to query pub node list.", true)
	common.AddRPCMethod("GET", config.QueryPrivP2PList, "This is an API to query priv node list.", true)
	common.AddRPCMethod("GET", config.QueryInterxList, "This is an API to query interx list.", true)
	common.AddRPCMethod("GET", config.QuerySnapList, "This is an API to query snap node list.", true)
}

func queryPubP2PNodeList(r *http.Request, rpcAddr string) (interface{}, interface{}, int) {
	global.Mutex.Lock()
	sort.Sort(types.P2PNodes(tasks.PubP2PNodeListResponse.NodeList))
	global.Mutex.Unlock()

	_ = r.ParseForm()
	connected := r.FormValue("connected") == "true"
	ip_only := r.FormValue("ip_only") == "true"

	if ip_only {
		ips := []string{}
		for _, node := range tasks.PubP2PNodeListResponse.NodeList {
			if connected == node.Connected {
				ips = append(ips, node.IP)
			}
		}

		return strings.Join(ips, ", "), nil, http.StatusOK
	}

	return tasks.PubP2PNodeListResponse, nil, http.StatusOK
}

// QueryNodeList is a function to query node list.
func QueryPubP2PNodeList(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-pub-node-list] Entering pub p2p node lists query")

		if !common.RPCMethods["GET"][config.QueryPubP2PList].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryPubP2PList].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-pub-node-list] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryPubP2PNodeList(r, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

func queryPrivP2PNodeList(r *http.Request, rpcAddr string) (interface{}, interface{}, int) {
	global.Mutex.Lock()
	sort.Sort(types.P2PNodes(tasks.PrivP2PNodeListResponse.NodeList))
	global.Mutex.Unlock()

	_ = r.ParseForm()
	connected := r.FormValue("connected") == "true"
	ip_only := r.FormValue("ip_only") == "true"

	if ip_only {
		ips := []string{}
		for _, node := range tasks.PrivP2PNodeListResponse.NodeList {
			if connected == node.Connected {
				ips = append(ips, node.IP)
			}
		}

		return strings.Join(ips, ", "), nil, http.StatusOK
	}

	return tasks.PrivP2PNodeListResponse, nil, http.StatusOK
}

// QueryNodeList is a function to query node list.
func QueryPrivP2PNodeList(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-priv-node-list] Entering priv p2p node lists query")

		if !common.RPCMethods["GET"][config.QueryPrivP2PList].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryPrivP2PList].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-priv-node-list] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryPrivP2PNodeList(r, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

func queryInterxList(r *http.Request, rpcAddr string) (interface{}, interface{}, int) {
	global.Mutex.Lock()
	sort.Sort(types.InterxNodes(tasks.InterxP2PNodeListResponse.NodeList))
	global.Mutex.Unlock()

	_ = r.ParseForm()
	ip_only := r.FormValue("ip_only") == "true"

	if ip_only {
		ips := []string{}
		for _, node := range tasks.InterxP2PNodeListResponse.NodeList {
			ips = append(ips, node.IP)
		}

		return strings.Join(ips, ", "), nil, http.StatusOK
	}

	return tasks.InterxP2PNodeListResponse, nil, http.StatusOK
}

// QueryNodeList is a function to query node list.
func QueryInterxList(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-interx-list] Entering interx lists query")

		if !common.RPCMethods["GET"][config.QueryInterxList].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryInterxList].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-interx-list] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryInterxList(r, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

func querySnapList(r *http.Request, rpcAddr string) (interface{}, interface{}, int) {
	global.Mutex.Lock()
	sort.Sort(types.SnapNodes(tasks.SnapNodeListResponse.NodeList))
	global.Mutex.Unlock()

	_ = r.ParseForm()
	ip_only := r.FormValue("ip_only") == "true"

	if ip_only {
		ips := []string{}
		for _, node := range tasks.SnapNodeListResponse.NodeList {
			ips = append(ips, node.IP)
		}

		return strings.Join(ips, ", "), nil, http.StatusOK
	}

	return tasks.SnapNodeListResponse, nil, http.StatusOK
}

// QueryNodeList is a function to query node list.
func QuerySnapList(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-snap-list] Entering snap lists query")

		if !common.RPCMethods["GET"][config.QuerySnapList].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QuerySnapList].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-snap-list] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = querySnapList(r, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

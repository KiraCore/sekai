package interx

import (
	"math/rand"
	"net/http"
	"sort"
	"strconv"
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
	response := tasks.PubP2PNodeListResponse

	_ = r.ParseForm()
	connected := r.FormValue("connected") == "true"
	ip_only := r.FormValue("ip_only") == "true"
	is_random := r.FormValue("order") == "random"
	is_format_simple := r.FormValue("format") == "simple"
	is_synced := r.FormValue("synced") == "true"

	if is_random {
		dest := make([]types.P2PNode, len(response.NodeList))
		perm := rand.Perm(len(response.NodeList))
		for i, v := range perm {
			dest[v] = response.NodeList[i]
		}
		response.NodeList = dest
	} else {
		sort.Sort(types.P2PNodes(response.NodeList))
	}

	if is_format_simple {
		indexOfPeer := make(map[string]string)
		for index, node := range response.NodeList {
			indexOfPeer[node.ID] = strconv.Itoa(index)
		}

		for nID, _ := range response.NodeList {
			for pIndex, _ := range response.NodeList[nID].Peers {
				if pid, isIn := indexOfPeer[response.NodeList[nID].Peers[pIndex]]; isIn {
					response.NodeList[nID].Peers[pIndex] = pid
				}
			}
		}
	}

	if is_synced {
		dest := make([]types.P2PNode, 0)
		for _, node := range response.NodeList {
			if node.Synced {
				dest = append(dest, node)
			}
		}
		response.NodeList = dest
	}

	global.Mutex.Unlock()

	if ip_only {
		ips := []string{}
		for _, node := range response.NodeList {
			if connected == node.Connected {
				ips = append(ips, node.IP)
			}
		}

		return strings.Join(ips, "\n"), nil, http.StatusOK
	}

	return response, nil, http.StatusOK
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
	response := tasks.PrivP2PNodeListResponse

	_ = r.ParseForm()
	connected := r.FormValue("connected") == "true"
	ip_only := r.FormValue("ip_only") == "true"
	is_random := r.FormValue("order") == "random"
	is_format_simple := r.FormValue("format") == "simple"
	is_synced := r.FormValue("synced") == "true"

	if is_random {
		dest := make([]types.P2PNode, len(response.NodeList))
		perm := rand.Perm(len(response.NodeList))
		for i, v := range perm {
			dest[v] = response.NodeList[i]
		}
		response.NodeList = dest
	} else {
		sort.Sort(types.P2PNodes(response.NodeList))
	}

	if is_format_simple {
		indexOfPeer := make(map[string]string)
		for index, node := range response.NodeList {
			indexOfPeer[node.ID] = strconv.Itoa(index)
		}

		for nID, _ := range response.NodeList {
			for pIndex, _ := range response.NodeList[nID].Peers {
				if pid, isIn := indexOfPeer[response.NodeList[nID].Peers[pIndex]]; isIn {
					response.NodeList[nID].Peers[pIndex] = pid
				}
			}
		}
	}

	if is_synced {
		dest := make([]types.P2PNode, 0)
		for _, node := range response.NodeList {
			if node.Synced {
				dest = append(dest, node)
			}
		}
		response.NodeList = dest
	}

	global.Mutex.Unlock()

	if ip_only {
		ips := []string{}
		for _, node := range response.NodeList {
			if connected == node.Connected {
				ips = append(ips, node.IP)
			}
		}

		return strings.Join(ips, "\n"), nil, http.StatusOK
	}

	return response, nil, http.StatusOK
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
	response := tasks.InterxP2PNodeListResponse

	_ = r.ParseForm()
	ip_only := r.FormValue("ip_only") == "true"
	is_random := r.FormValue("order") == "random"
	is_synced := r.FormValue("synced") == "true"

	if is_random {
		dest := make([]types.InterxNode, len(response.NodeList))
		perm := rand.Perm(len(response.NodeList))
		for i, v := range perm {
			dest[v] = response.NodeList[i]
		}
		response.NodeList = dest
	} else {
		sort.Sort(types.InterxNodes(response.NodeList))
	}

	if is_synced {
		dest := make([]types.InterxNode, 0)
		for _, node := range response.NodeList {
			if node.Synced {
				dest = append(dest, node)
			}
		}
		response.NodeList = dest
	}
	global.Mutex.Unlock()

	if ip_only {
		ips := []string{}
		for _, node := range response.NodeList {
			ips = append(ips, node.IP)
		}

		return strings.Join(ips, "\n"), nil, http.StatusOK
	}

	return response, nil, http.StatusOK
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
	response := tasks.SnapNodeListResponse

	_ = r.ParseForm()
	ip_only := r.FormValue("ip_only") == "true"
	is_random := r.FormValue("order") == "random"
	is_synced := r.FormValue("synced") == "true"

	if is_random {
		dest := make([]types.SnapNode, len(response.NodeList))
		perm := rand.Perm(len(response.NodeList))
		for i, v := range perm {
			dest[v] = response.NodeList[i]
		}
		response.NodeList = dest
	} else {
		sort.Sort(types.SnapNodes(response.NodeList))
	}

	if is_synced {
		dest := make([]types.SnapNode, 0)
		for _, node := range response.NodeList {
			if node.Synced {
				dest = append(dest, node)
			}
		}
		response.NodeList = dest
	}
	global.Mutex.Unlock()

	if ip_only {
		ips := []string{}
		for _, node := range response.NodeList {
			ips = append(ips, node.IP)
		}

		return strings.Join(ips, "\n"), nil, http.StatusOK
	}

	return response, nil, http.StatusOK
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

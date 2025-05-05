package routers

import (
	"vngom/fiber_wrapper"
	"vngom/routers/auth"
)

var Routes map[string]fiber_wrapper.Router = make(map[string]fiber_wrapper.Router)

func init() {
	Routes["/auth/login"] = fiber_wrapper.Router{
		Method:  "GET",
		Handler: auth.Login,
	}
	Routes["/auth/get-tenant"] = fiber_wrapper.Router{
		Method:  "GET",
		Handler: auth.GetTenant,
	}

}

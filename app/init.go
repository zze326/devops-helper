package app

import (
	_ "github.com/revel/modules"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/filters"
	_ "github.com/zze326/devops-helper/app/startups"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		filters.HeaderFilter,          // Add some security based headers
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		filters.NoAuthFilter,          // 部分请求不需认证
		filters.AuthenticationFilter,  // JWT 认证
		filters.AuthorizationFilter,   // 授权
		filters.TerminalAuthFilter,    // 主机终端授权
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.InterceptorFilter,       // Run interceptors around the action.
		//revel.CompressFilter,    // Compress the result.
		revel.BeforeAfterFilter, // Call the before and after filter functions
		filters.ActionInvoker,
	}
}

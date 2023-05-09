package startups

import "github.com/revel/revel"

func init() {
	revel.OnAppStart(func() {
		registerTag()
		initDB()
		initCasbin()
	})
}

package dbInstance

import "github.com/upbound/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_db_instance", func(r *config.Resource) {

		//r.ShortGroup = "db"
		r.Kind = "DbInstance"
		r.UseAsync = true
	})
}

package kvstoreInstance

import "github.com/crossplane/terrajet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_kvstore_instance", func(r *config.Resource) {

		r.Kind = "KvStoreInstance"
		r.UseAsync = true
	})
}

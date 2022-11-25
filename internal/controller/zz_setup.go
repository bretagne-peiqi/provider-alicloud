/*
Copyright 2021 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/upbound/upjet/pkg/controller"

	dbinstance "github.com/bretagne-peiqi/provider-alicloud/internal/controller/db/dbinstance"
	kvstoreinstance "github.com/bretagne-peiqi/provider-alicloud/internal/controller/kvstore/kvstoreinstance"
	providerconfig "github.com/bretagne-peiqi/provider-alicloud/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		dbinstance.Setup,
		kvstoreinstance.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

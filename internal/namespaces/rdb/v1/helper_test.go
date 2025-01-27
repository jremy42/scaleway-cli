package rdb

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	name     = "cli-test"
	user     = "foobar"
	password = "{4xdl*#QOoP+&3XRkGA)]"
	engine   = "PostgreSQL-12"
)

func createInstance(engine string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf(baseCommand, name, engine, user, password),
	)
}

func createInstanceWithPrivateNetworkAndLoadBalancer(engine string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf(baseCommand+privateNetworkStaticSpec+loadBalancerSpec, name, engine, user, password),
	)
}

func createPN() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		api := vpc.NewAPI(ctx.Client)
		pn, err := api.CreatePrivateNetwork(&vpc.CreatePrivateNetworkRequest{})
		if err != nil {
			return err
		}
		ctx.Meta["PN"] = pn
		if len(pn.Subnets) > 0 {
			ctx.Meta["IPNet"], err = getIPSubnet(pn.Subnets[0])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func getIPSubnet(ipNet scw.IPNet) (*string, error) {
	addr := ipNet.IP.To4()
	if addr == nil {
		return nil, fmt.Errorf("could get ip 4 bytes")
	}
	addr = addr.Mask(addr.DefaultMask())
	addr[3] = +3

	sz, _ := ipNet.Mask.Size()
	ipNetStr := fmt.Sprintf("%s/%d", addr.String(), sz)
	return &ipNetStr, nil
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}")
}

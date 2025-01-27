package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func Test_UpdateSnapshot(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		t.Run("Change tags", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createVolume("Volume", 10, instance.VolumeVolumeTypeBSSD),
				core.ExecStoreBeforeCmd("CreateSnapshot", "scw instance snapshot create volume-id={{ .Volume.ID }} name=cli-test-snapshot-update-tags tags.0=foo tags.1=bar"),
			),
			Cmd: "scw instance snapshot update -w {{ .CreateSnapshot.Snapshot.ID }} tags.0=bar tags.1=foo",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotNil(t, ctx.Result)
					snapshot := ctx.Result.(*instance.Snapshot)
					assert.Equal(t, snapshot.Name, "cli-test-snapshot-update-tags")
					assert.Len(t, snapshot.Tags, 2)
					assert.Equal(t, snapshot.Tags[0], "bar")
					assert.Equal(t, snapshot.Tags[1], "foo")
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSnapshot("CreateSnapshot"),
				deleteVolume("Volume"),
			),
		}))
		t.Run("Change name", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createVolume("Volume", 10, instance.VolumeVolumeTypeBSSD),
				core.ExecStoreBeforeCmd("CreateSnapshot", "scw instance snapshot create volume-id={{ .Volume.ID }} name=cli-test-snapshot-update-name tags.0=foo tags.1=bar"),
			),
			Cmd: "scw instance snapshot update -w {{ .CreateSnapshot.Snapshot.ID }} name=cli-test-snapshot-update-name-updated",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotNil(t, ctx.Result)
					snapshot := ctx.Result.(*instance.Snapshot)
					assert.Equal(t, snapshot.Name, "cli-test-snapshot-update-name-updated")
					assert.Len(t, snapshot.Tags, 2)
					assert.Equal(t, snapshot.Tags[0], "foo")
					assert.Equal(t, snapshot.Tags[1], "bar")
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteSnapshot("CreateSnapshot"),
				deleteVolume("Volume"),
			),
		}))
	})
}

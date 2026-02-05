/**
 * 功能：提供CLI命令处理器包装函数，将core.CliHandlerFunc转换为cli.ActionFunc
 * 关联：依赖github.com/urfave/cli/v3和github.com/mayswind/ezbookkeeping/pkg/core
 * 注意：用于简化CLI命令的绑定过程
 */
package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func bindAction(fn core.CliHandlerFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		c := core.WrapCilContext(ctx, cmd)
		return fn(c)
	}
}

package cmd

import (
	"fmt"

	"hehan.net/my/stockcmd/baostock"
	"hehan.net/my/stockcmd/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/stat"
	"hehan.net/my/stockcmd/store"
)

var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Dev command to manipulate data",
}

var RemoveKDataCmd = &cobra.Command{
	Use:     "remove_kdata [code]",
	Short:   "remove kdata of certain code",
	Long:    `remove kdata of certain code`,
	Example: ` remove_kdata sz.002475`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    removeKDataCmdF,
}

var ShowKDataCmd = &cobra.Command{
	Use:     "show_kdata [code]",
	Short:   "show kdata of certain code",
	Long:    `show kdata of certain code`,
	Example: ` show_kdata sz.002475`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    showKDataCmdf,
}

var RemoveAllKData = &cobra.Command{
	Use:   "remove_all_kdata",
	Short: "remove kdata of all codes",
	RunE:  removeAllKDataCmdF,
}

var RemoveGroupKDataCmd = &cobra.Command{
	Use:     "remove_group_kdata [group]",
	Short:   "remove kdata of certain group",
	Long:    `remove kdata of certain group`,
	Example: ` remove_group_kdata hold`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    removeGroupKDataCmdF,
}

var UpdateBasicData = &cobra.Command{
	Use:   "update_basic_data",
	Short: "update_basic_data",
	RunE:  updateBasicDataCmdF,
}

func updateBasicDataCmdF(cmd *cobra.Command, args []string) error {
	store.RecreateBasicBucket()

	baostock.BS.Login()
	defer baostock.BS.Logout()

	rs, err := baostock.BS.QueryAllStock(util.GetLastWorkDay())
	if err != nil {
		fmt.Println(err)
		return err
	}

	store.WriteBasics(rs.Data)
	return nil
}

func removeKDataCmdF(cmd *cobra.Command, args []string) error {
	code := args[0]
	store.DeleteCodeRecords(code)
	return nil
}

func showKDataCmdf(cmd *cobra.Command, args []string) error {
	code := args[0]

	df, err := stat.GetDataFrame(code)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(df.Table())
	return nil
}

func removeAllKDataCmdF(cmd *cobra.Command, args []string) error {
	store.RecreateDailyBucket()
	return nil
}

func removeGroupKDataCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	g := store.GetGroup(gName)
	if g == nil {
		return errors.Errorf("Group [%s] not exist", gName)
	}

	for code, _ := range g.Codes {
		store.DeleteCodeRecords(code)
	}
	return nil
}

func init() {
	DevCmd.AddCommand(RemoveKDataCmd, RemoveGroupKDataCmd, RemoveAllKData, ShowKDataCmd, UpdateBasicData)
	rootCmd.AddCommand(DevCmd)
}

package main

import (
	"fmt"
	"time"

	"github.com/coreservice-io/general-counter"
	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log"
	"gorm.io/gorm"
)

func main() {

	llog, err := logrus_log.New("./logs", 1, 20, 30)
	if err != nil {
		panic(err.Error())
	}

	llog.SetLevel(log.TraceLevel)

	gcounter, err := general_counter.NewGeneralCounter(&general_counter.GeneralCounterConfig{
		Project_name: "", // config your own
		Db_config: &general_counter.DBConfig{
			Host:     "", // config your own
			Port:     0,  // config your own
			DbName:   "", // config your own
			UserName: "", // config your own
			Password: "", // config your own
		},
		Ecs_config: &general_counter.EcsConfig{
			Address:  "", // config your own
			UserName: "", // config your own
			Password: "", // config your own
		},
		Redis_config: &general_counter.RedisConfig{
			Addr:     "",    // config your own
			Port:     0,     // config your own
			UserName: "",    // config your own
			Password: "",    // config your own
			Prefix:   "",    // config your own
			UseTLS:   false, // config your own
		},
	}, llog)

	if err != nil {
		panic(err)
	}

	commit_err := gcounter.CreateTx().AppendFunc(func(tx *gorm.DB) error {
		fmt.Println("this is first func")
		return nil
	}).AppendOp(&general_counter.GcOp{
		Gkey:   "userid13",
		Gtype:  "total_balance",
		Amount: 200,
		Total_config: &general_counter.GcOpTotalConfig{
			Enable:        true,
			AllowNegative: false,
		},
	}).AppendOp(&general_counter.GcOp{
		Gkey:   "userid13",
		Gtype:  "transfer_out",
		Amount: -200,
		Total_config: &general_counter.GcOpTotalConfig{
			Enable:        true,
			AllowNegative: true,
		},
		Detail_config: &general_counter.GcOpDetailConfig{
			Enable: true,
			Msg:    "to userid14",
		},
	}).AppendOp(&general_counter.GcOp{
		Gkey:   "userid14",
		Gtype:  "transfer_in",
		Amount: 200,
		Total_config: &general_counter.GcOpTotalConfig{
			Enable:        true,
			AllowNegative: false,
		},
		Detail_config: &general_counter.GcOpDetailConfig{
			Enable: true,
			Msg:    "from userid13",
		},
	}).AppendFunc(func(tx *gorm.DB) error {
		fmt.Println("this is second func")
		return nil
	}).Commit()

	fmt.Println(commit_err)

	result, err := gcounter.QueryTotal("userid1", "total_balance")
	if err != nil {
		fmt.Println("query err:", err)
	} else {
		fmt.Println(result)
	}

	result, err = gcounter.QueryTotal("userid3", "total_balance")
	if err != nil {
		fmt.Println("query err:", err)
	} else {
		fmt.Println(result)
	}

	// aggResult, err := gcounter.QueryAgg("userid1", "mining_reward", "2022-12-16", "2023-01-15")
	// if err != nil {
	//     fmt.Println("query err:", err)
	// } else {
	//     for _, v := range aggResult {
	//         fmt.Println(v)
	//     }
	// }
	//
	// detailResult, err := gcounter.QueryDetail("userid1", "transfer_out", "2022-12-16", "2023-01-15")
	// if err != nil {
	//     fmt.Println("query err:", err)
	// } else {
	//     for _, v := range detailResult {
	//         fmt.Println(v)
	//     }
	// }
	//
	// detailResult2, err := gcounter.QueryDetail("userid2", "transfer_in", "2022-12-16", "2023-01-15")
	// if err != nil {
	//     fmt.Println("query err:", err)
	// } else {
	//     for _, v := range detailResult2 {
	//         fmt.Println(v)
	//     }
	// }

	for {
		time.Sleep(1 * time.Hour)
	}

}

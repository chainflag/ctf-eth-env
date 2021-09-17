package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

func makeCliqueGenesis(sealer common.Address, chainID *big.Int, period uint64) *core.Genesis {
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(1),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			ChainID:             chainID,
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			Clique: &params.CliqueConfig{
				Period: period,
				Epoch:  30000,
			},
		},
	}
	genesis.ExtraData = make([]byte, 32+common.AddressLength+65)
	copy(genesis.ExtraData[32:], sealer[:])
	genesis.Alloc[sealer] = core.GenesisAccount{
		Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
	}
	return genesis
}

func fatalExit(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(1)
}

func main() {
	app := &cli.App{
		Name:  "genesis-builder",
		Usage: "Create a ethereum clique consensus genesis spec file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "config",
				Usage:   "output `directory`",
			},
			&cli.StringFlag{
				Name:     "address",
				Usage:    "account for seal and pre-funded",
				Required: true,
			},
			&cli.Uint64Flag{
				Name:        "chainid",
				Value:       0,
				Usage:       "chainid for the POA Network",
				DefaultText: "random",
			},
			&cli.Uint64Flag{
				Name:  "period",
				Value: 15,
				Usage: "seconds of block time",
			},
		},
		Action: func(c *cli.Context) error {
			address := c.String("address")
			if !common.IsHexAddress(address) {
				fatalExit(errors.New("invalid address"))
			}
			if c.Uint64("chainid") > 65535 {
				fatalExit(errors.New("invalid chainid"))
			}
			chainID := c.Uint64("chainid")
			if chainID == 0 {
				chainID = uint64(rand.Intn(65536))
			}
			genesisPath := filepath.Join(c.String("output"), "genesis.json")
			genesis, _ := json.MarshalIndent(makeCliqueGenesis(common.HexToAddress(address), new(big.Int).SetUint64(chainID), c.Uint64("period")), "", "  ")
			if err := os.WriteFile(genesisPath, genesis, 0644); err != nil {
				fatalExit(fmt.Errorf("failed to save genesis file: %v", err))
			}
			fmt.Printf("Path of the genesis file: %s\n", genesisPath)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fatalExit(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

type Keystore struct {
	Address common.Address
	Path    string
}

func createKeystore(dir, auth string) (*Keystore, error) {
	account, err := keystore.StoreKey(dir, auth, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}

	return &Keystore{
		Address: account.Address,
		Path:    account.URL.Path,
	}, nil
}

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

	if chainID == nil {
		genesis.Config.ChainID = new(big.Int).SetUint64(uint64(rand.Intn(65536)))
	}
	if period == 0 {
		genesis.Config.Clique.Period = 15
	}

	genesis.ExtraData = make([]byte, 32+common.AddressLength+65)
	copy(genesis.ExtraData[32:], sealer[:])
	genesis.Alloc[sealer] = core.GenesisAccount{
		Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
	}

	return genesis
}

func saveGenesis(folder, network string, genesis *core.Genesis) error {
	path := filepath.Join(folder, network+".json")
	out, _ := json.MarshalIndent(genesis, "", "  ")
	return ioutil.WriteFile(path, out, 0644)
}

func main() {
	app := &cli.App{
		Name:  "conf-gen",
		Usage: "generate config",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "path",
				Usage: "",
			},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "all",
			Usage: "",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "keystore",
			Usage: "",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "genesis",
			Usage: "",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

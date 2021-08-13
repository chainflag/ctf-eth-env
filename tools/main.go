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
		UseShortOptionHandling: true,
		Name:                   "conf-gen",
		Usage:                  "generate config",
		Action: func(c *cli.Context) error {
			fmt.Println("you can use 'conf-gen create -p 199924' to make what you want")
			fmt.Println("you can use 'conf-gen  -h ' to get some help")
			return nil
		},
		Flags: []cli.Flag{},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "create",
			Usage: "Create all of them",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "folder",
					Usage:    "your config address",
					Value:    "../config",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "your password",
					Required: true,
					Aliases:  []string{"p"}},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(filepath.Join(c.String("path"), "keystore"), c.String("password"))
				if err != nil {
					log.Fatal(err)
				}
				saveGenesis("../config", "genesis", makeCliqueGenesis(ks.Address, nil, 15))
				fmt.Println("Successful,Here is the public key address corresponding to your keystore:")
				fmt.Println(ks.Address)
				return nil
			},
		},
		{
			Name:  "keystore",
			Usage: "Create your keystore",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "folder",
					Value:    "../config/keystore",
					Usage:    "your Keystore folder",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "your address",
					Required: true,
					Aliases:  []string{"p"}},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(c.String("folder"), c.String("password"))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Successful,Here is the public key address corresponding to your keystore:")
				fmt.Println(ks.Address)
				return nil
			},
		},
		{
			Name:  "genesis",
			Usage: "Create genesis to generate your private chains..",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "genName",
					Value:    "genesis",
					Usage:    "your genesisName",
					Required: false,
					Aliases:  []string{"n"}},
				&cli.StringFlag{
					Name:     "folder",
					Value:    "../config",
					Usage:    "your folder",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "address",
					Usage:    "your address",
					Required: true,
					Aliases:  []string{"a"}},
				&cli.Int64Flag{
					Name:     "chainId",
					Value:    1,
					Usage:    "your ChainId",
					Required: false,
					Aliases:  []string{"i"}},
				&cli.Uint64Flag{
					Name:     "seconds",
					Value:    15,
					Usage:    "Your seconds",
					Required: false,
					Aliases:  []string{"s"}},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("hello world")
				var genesis *core.Genesis
				if c.Int64("chainId") == 1 {
					genesis = makeCliqueGenesis(common.HexToAddress(c.String("address")), nil, c.Uint64("seconds"))
				} else {
					genesis = makeCliqueGenesis(common.HexToAddress(c.String("address")), big.NewInt(c.Int64("chainId")), c.Uint64("seconds"))
				}
				saveGenesis(c.String("folder"), c.String("genName"), genesis)
				fmt.Println("Successful")
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

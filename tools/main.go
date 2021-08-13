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
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "Welcome to conf-Gen, a tool used to create everything from keystore to Genesis that prepares you for creating your private chain", 0x1B)
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "           You can use 'conf-gen create -p <Your password>' to create everything very easily", 0x1B)
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "           You can also use 'conf-gen -h' to see how to do more personalized configuration", 0x1B)
			return nil
		},
		Flags: []cli.Flag{},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "create",
			Usage: "To create everything you need to set up a private chain",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "folder",
					Usage:    "Path of the configuration file.",
					Value:    "../config",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "The key used to decrypt your <keystore>.[Must be given]",
					Required: true,
					Aliases:  []string{"p"}},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(filepath.Join(c.String("path"), "keystore"), c.String("password"))
				if err != nil {
					log.Fatal(err)
				}
				saveGenesis("../config", "genesis", makeCliqueGenesis(ks.Address, nil, 15))
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, " Successfully created the required profile.", 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "The path to <genesis> : "+c.String("folder"), 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "The path to <keystore> : "+c.String("folder")+"/keystore", 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "Here is the public key address corresponding to your <keystore>:", 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, ks.Address, 0x1B)
				return nil
			},
		},
		{
			Name:  "keystore",
			Usage: "Create your <keystore> more personalized.Use - h for more information",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "folder",
					Value:    "../config/keystore",
					Usage:    "Path of the <keystore> file.",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "password",
					Usage:    "The key used to decrypt your keystore.[Must be given]",
					Required: true,
					Aliases:  []string{"p"}},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(c.String("folder"), c.String("password"))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "The path to <keystore> : "+c.String("folder"), 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "Here is the public key address corresponding to your <keystore>:", 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, ks.Address, 0x1B)
				return nil
			},
		},
		{
			Name:  "genesis",
			Usage: "Create your <genesis> more personalized.Use - h for more information.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "genName",
					Value:    "genesis",
					Usage:    "The name of your <genesis>.",
					Required: false,
					Aliases:  []string{"n"}},
				&cli.StringFlag{
					Name:     "folder",
					Value:    "../config",
					Usage:    "Path of the configuration file.",
					Required: false,
					Aliases:  []string{"f", "v"}},
				&cli.StringFlag{
					Name:     "address",
					Usage:    "Your original address[Must be given]",
					Required: true,
					Aliases:  []string{"a"}},
				&cli.Int64Flag{
					Name:     "chainId",
					Value:    1,
					Usage:    "The Chainid of your private chain.<Please be careful not to use the popular chainid>",
					Required: false,
					Aliases:  []string{"i"}},
				&cli.Uint64Flag{
					Name:     "seconds",
					Value:    15,
					Usage:    "The time to generate a new block.",
					Required: false,
					Aliases:  []string{"s"}},
			},
			Action: func(c *cli.Context) error {
				var genesis *core.Genesis
				if c.Int64("chainId") == 1 {
					genesis = makeCliqueGenesis(common.HexToAddress(c.String("address")), nil, c.Uint64("seconds"))
				} else {
					genesis = makeCliqueGenesis(common.HexToAddress(c.String("address")), big.NewInt(c.Int64("chainId")), c.Uint64("seconds"))
				}
				saveGenesis(c.String("folder"), c.String("genName"), genesis)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "Genesis created successfully", 0x1B)
				fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "The path to <Genesis> : "+c.String("folder"), 0x1B)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

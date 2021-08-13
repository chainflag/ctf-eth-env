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
	if network == "" {
		network = "genesis"
	}
	path := filepath.Join(folder, network+".json")
	out, _ := json.MarshalIndent(genesis, "", "  ")
	return ioutil.WriteFile(path, out, 0644)
}

func main() {
	app := &cli.App{
		Name:  "conf-gen",
		Usage: "generate config",
		Action: func(c *cli.Context) error {
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "Welcome to conf-Gen, a tool used to create everything from keystore to Genesis that prepares you for creating your private chain", 0x1B)
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "        You can use 'conf-gen create -password <Your password>' to create everything very easily", 0x1B)
			fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "           You can also use 'conf-gen -h' to see how to do more personalized configuration", 0x1B)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "folder",
				Value:    "config",
				Usage:    "Path of the configuration file.",
				Required: false,
				Aliases:  []string{"f"}},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "all",
			Usage: "To create everything you need to set up the ctf eth env",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "password",
					Usage:    "Used to unlock your account at geth launch",
					Required: true},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(filepath.Join(c.String("folder"), "keystore"), c.String("password"))
				if err != nil {
					log.Fatalf("Failed to create account: %v", err)
				}
				if err := saveGenesis(c.String("folder"), "", makeCliqueGenesis(ks.Address, nil, 15)); err != nil {
					log.Fatalf("Failed to save genesis file: %v", err)
				}

				fmt.Printf("\nSuccessfully created the required config\n\n")
				fmt.Printf("Path of the secret key file: %s\n", ks.Path)
				fmt.Printf("Path of the genesis file:    %s\n\n", filepath.Join(c.String("folder"), "genesis.json"))
				return nil
			},
		},
		{
			Name:  "keystore",
			Usage: "Create a new account and save it in keystore",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "password",
					Usage:    "Your new account is locked with the password",
					Required: true},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(filepath.Join(c.String("folder"), "keystore"), c.String("password"))
				if err != nil {
					log.Fatalf("Failed to create account: %v", err)
				}
				fmt.Printf("\nYour new key was generated\n\n")
				fmt.Printf("Public address of the key:   %s\n", ks.Address.Hex())
				fmt.Printf("Path of the secret key file: %s\n\n", ks.Path)
				return nil
			},
		},
		{
			Name:  "genesis",
			Usage: "Create a Clique consensus genesis spec",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "address",
					Usage:    "Account for seal and pre-funded",
					Required: true},
				&cli.Int64Flag{
					Name:     "chainid",
					Value:    0,
					Usage:    "Chain ID for the POA Network",
					Required: false},
				&cli.Uint64Flag{
					Name:     "period",
					Value:    15,
					Usage:    "Seconds of block time",
					Required: false},
			},
			Action: func(c *cli.Context) error {
				var chainID *big.Int
				if c.Int64("chainid") != 0 {
					chainID = big.NewInt(c.Int64("chainid"))
				}
				genesis := makeCliqueGenesis(common.HexToAddress(c.String("address")), chainID, c.Uint64("period"))
				fmt.Printf("\nConfigured new genesis spec\n\n")
				if err := saveGenesis(c.String("folder"), "", genesis); err != nil {
					log.Fatalf("Failed to save genesis file: %v", err)
				}
				fmt.Printf("Path of the genesis file: %s\n\n", filepath.Join(c.String("folder"), "genesis.json"))
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
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

func createKeystore(keydir, auth string) (*Keystore, error) {
	account, err := keystore.StoreKey(keydir, auth, keystore.StandardScryptN, keystore.StandardScryptP)
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
	genesis.ExtraData = make([]byte, 32+common.AddressLength+65)
	copy(genesis.ExtraData[32:], sealer[:])
	genesis.Alloc[sealer] = core.GenesisAccount{
		Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
	}

	return genesis
}

func saveGenesis(genesisPath string, genesis *core.Genesis) error {
	path, _ := filepath.Abs(genesisPath)
	out, _ := json.MarshalIndent(genesis, "", "  ")
	return ioutil.WriteFile(path, out, 0644)
}

func fatalExit(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(1)
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	app := &cli.App{
		Name:  "conf-gen",
		Usage: "To create everything you need to set up the ctf eth env",
		Action: func(c *cli.Context) error {
			folder := c.String("folder")
			rand.Seed(time.Now().UnixNano())
			password := randSeq(20)
			ks, err := createKeystore(filepath.Join(folder, "keystore"), password)
			if err != nil {
				fatalExit(fmt.Errorf("failed to create account: %v", err))
			}
			passwordPath := filepath.Join(folder, "password.txt")
			if err := ioutil.WriteFile(passwordPath, []byte(password), 0644); err != nil {
				fatalExit(fmt.Errorf("failed to save keystore pass: %v", err))
			}
			genesisPath := filepath.Join(folder, "genesis.json")
			if err := saveGenesis(genesisPath, makeCliqueGenesis(ks.Address, new(big.Int).SetUint64(uint64(rand.Intn(65536))), 15)); err != nil {
				fatalExit(fmt.Errorf("failed to save genesis file: %v", err))
			}
			fmt.Printf("\nSuccessfully created the required config\n\n")
			fmt.Printf("Path of the secret key file:   %s\n", ks.Path)
			fmt.Printf("Path of the keystore passowrd: %s\n", passwordPath)
			fmt.Printf("Path of the genesis file:      %s\n\n", genesisPath)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "folder",
				Value:    "config",
				Usage:    "directory to store configuration files",
				Required: false},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "keystore",
			Usage: "Create a new account and save it in keystore",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "password",
					Usage:    "your new account is locked with the password",
					Required: true},
			},
			Action: func(c *cli.Context) error {
				ks, err := createKeystore(filepath.Join(c.String("folder"), "keystore"), c.String("password"))
				if err != nil {
					fatalExit(fmt.Errorf("failed to create account: %v", err))
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
					Usage:    "account for seal and pre-funded",
					Required: true},
				&cli.Int64Flag{
					Name:     "chainid",
					Value:    0,
					Usage:    "chainid for the POA Network",
					Required: true},
				&cli.Uint64Flag{
					Name:     "period",
					Value:    15,
					Usage:    "seconds of block time",
					Required: false},
			},
			Action: func(c *cli.Context) error {
				address := c.String("address")
				re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
				if !re.MatchString(address) {
					fatalExit(errors.New("invalid address"))
				}
				chainID := c.Int64("chainid")
				if chainID <= 0 || chainID > 65535 {
					fatalExit(errors.New("invalid chainid"))
				}
				genesis := makeCliqueGenesis(common.HexToAddress(address), big.NewInt(chainID), c.Uint64("period"))
				fmt.Printf("\nConfigured new genesis spec\n\n")
				genesisPath := filepath.Join(c.String("folder"), "genesis.json")
				if err := saveGenesis(genesisPath, genesis); err != nil {
					fatalExit(fmt.Errorf("failed to save genesis file: %v", err))
				}
				fmt.Printf("Path of the genesis file: %s\n\n", genesisPath)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fatalExit(err)
	}
}

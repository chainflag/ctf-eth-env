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
		Usage: "ctf-eth-env configuration generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "config",
				Usage:   "output `directory`",
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
			if c.Uint64("chainid") > 65535 {
				fatalExit(errors.New("invalid chainid"))
			}
			folder := c.String("output")
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
			chainID := c.Uint64("chainid")
			if chainID == 0 {
				chainID = uint64(rand.Intn(65536))
			}
			genesisPath := filepath.Join(folder, "genesis.json")
			genesis, _ := json.MarshalIndent(makeCliqueGenesis(ks.Address, new(big.Int).SetUint64(chainID), c.Uint64("period")), "", "  ")
			if err := ioutil.WriteFile(genesisPath, genesis, 0644); err != nil {
				fatalExit(fmt.Errorf("failed to save genesis file: %v", err))
			}
			fmt.Printf("\nSuccessfully created the required config\n\n")
			fmt.Printf("Public address of new key:     %s\n", ks.Address.Hex())
			fmt.Printf("Path of the secret key file:   %s\n", ks.Path)
			fmt.Printf("Path of the keystore passowrd: %s\n", passwordPath)
			fmt.Printf("Path of the genesis file:      %s\n", genesisPath)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fatalExit(err)
	}
}

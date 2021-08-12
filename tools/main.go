package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
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
	path := filepath.Join(folder, fmt.Sprintf("%s.json", network))
	out, _ := json.MarshalIndent(genesis, "", "  ")
	return ioutil.WriteFile(path, out, 0644)
}

func main() {
	ks, err := createKeystore("../config/keystore", "password")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ks)

	genesis := makeCliqueGenesis(ks.Address, nil, 15)
	fmt.Println(genesis)

	saveGenesis("../config", "genesis", genesis)
}

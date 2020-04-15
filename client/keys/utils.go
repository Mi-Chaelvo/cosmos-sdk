package keys

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"

	// defaultKeyDBName is the client's subdirectory where keys are stored.
	defaultKeyDBName = "keys"
)

type bechKeyOutFn func(keyInfo cryptokeyring.Info) (cryptokeyring.KeyOutput, error)

// NewLegacyKeyBaseFromDir initializes a legacy keybase at the rootDir directory. Keybase
// options can be applied when generating this new Keybase.
func NewLegacyKeyBaseFromDir(rootDir string, opts ...cryptokeyring.KeybaseOption) (cryptokeyring.LegacyKeybase, error) {
	return getLegacyKeyBaseFromDir(rootDir, opts...)
}

func getLegacyKeyBaseFromDir(rootDir string, opts ...cryptokeyring.KeybaseOption) (cryptokeyring.LegacyKeybase, error) {
	return cryptokeyring.NewLegacy(defaultKeyDBName, filepath.Join(rootDir, "keys"), opts...)
}

func printKeyInfo(w io.Writer, keyInfo cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}

	switch viper.Get(cli.OutputFlag) {
	case OutputFormatText:
		printTextInfos(w, []cryptokeyring.KeyOutput{ko})

	case OutputFormatJSON:
		var out []byte
		var err error
		if viper.GetBool(flags.FlagIndentResponse) {
			out, err = KeysCdc.MarshalJSONIndent(ko, "", "  ")
		} else {
			out, err = KeysCdc.MarshalJSON(ko)
		}
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, string(out))
	}
}

func printInfos(w io.Writer, infos []cryptokeyring.Info) {
	kos, err := cryptokeyring.Bech32KeysOutput(infos)
	if err != nil {
		panic(err)
	}

	switch viper.Get(cli.OutputFlag) {
	case OutputFormatText:
		printTextInfos(w, kos)

	case OutputFormatJSON:
		var out []byte
		var err error

		if viper.GetBool(flags.FlagIndentResponse) {
			out, err = KeysCdc.MarshalJSONIndent(kos, "", "  ")
		} else {
			out, err = KeysCdc.MarshalJSON(kos)
		}

		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s", out)
	}
}

func printTextInfos(w io.Writer, kos []cryptokeyring.KeyOutput) {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(out))
}

func printKeyAddress(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.Address)
}

func printPubKey(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.PubKey)
}

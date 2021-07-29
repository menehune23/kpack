package cosigner

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sigstore/cosign/cmd/cosign/cli"
)

type ImageFetcher interface {
	Fetch(keychain authn.Keychain, repoName string) (v1.Image, string, error)
}

type ImageSigner struct {
	Logger *log.Logger
	Client ImageFetcher
}

var cliSignCmd = cli.SignCmd

// Todo: user service account secrets
// Iteration 1: use secret
// Iteration 2?: use service account
// The funcion calling this will do:
// For every key in the SA:
//   sign(imageRef, key.Path)
//      GetPass function would use: client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{});

// Other keyops support: https://github.com/sigstore/cosign/blob/143e47a120702f175e68e0a04594cb87a4ce8e02/cmd/cosign/cli/sign.go#L167
// Todo: Annotation obtained from kpack config

func (s *ImageSigner) Sign(refImage, keyPath string) error {
	if refImage == "" {
		return fmt.Errorf("signing reference image is empty")
	}

	if keyPath == "" {
		return fmt.Errorf("signing key path is empty")
	}

	ctx := context.Background()
	ko := cli.KeyOpts{KeyRef: keyPath}

	if err := cliSignCmd(ctx, ko, nil, refImage, "", true, "", false, false); err != nil {
		return fmt.Errorf("signing: %v", err)
	}

	return nil
}

// cosign sign -key <> <image>

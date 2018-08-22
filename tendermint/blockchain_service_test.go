package tendermint

import (
	"path/filepath"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
)

func TestSave_then_retrieve_Success(t *testing.T) {
	//variables:
	namespace := "xsuite"
	name := "users"
	id := uuid.NewV4()
	rootPath := filepath.Join("./test_files")

	blkchain, blkchainErr := generateBlockchain(namespace, name, &id)
	if blkchainErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", blkchainErr.Error())
		return
	}

	//create the service:
	service := createBlockchainService(rootPath)

	//save the blockchain:
	saveErr := service.Save(blkchain)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	//retrieve the blockchain:
	retBlkchain, retBlkchainErr := service.Retrieve(blkchain.GetGenesis().GetPath())
	if retBlkchainErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retBlkchainErr.Error())
		return
	}

	if !reflect.DeepEqual(blkchain, retBlkchain) {
		t.Errorf("the saved blockchain is invalid")
		return
	}

	//load the keys in tendermint, to make sure it works:
	confRootPath := filepath.Join(rootPath, blkchain.GetGenesis().GetPath().String())
	conf := config.DefaultConfig().SetRoot(confRootPath)
	pv := privval.LoadFilePV(conf.PrivValidatorFile())

	if !pv.GetPubKey().Equals(blkchain.GetPK().PubKey()) {
		t.Errorf("the generated blockchain files were invalid")
		return
	}

	//delete the blockchain:
	rmErr := service.Delete(blkchain.GetGenesis().GetPath())
	if rmErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", rmErr.Error())
		return
	}
}

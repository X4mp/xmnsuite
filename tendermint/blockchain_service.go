package tendermint

import (
	"io/ioutil"
	"os"
	"path/filepath"

	config "github.com/tendermint/tendermint/config"
	crypto "github.com/tendermint/tendermint/crypto"
)

type containedPK struct {
	PK crypto.PrivKey `json:"priv_key"`
}

type blockchainService struct {
	rootPath string
	conf     *config.Config
}

func createBlockchainService(rootPath string) BlockchainService {
	out := blockchainService{
		rootPath: rootPath,
		conf:     config.DefaultConfig(),
	}

	return &out
}

// Retrieve retrieves a blockchain
func (obj *blockchainService) Retrieve(path Path) (Blockchain, error) {

	//create the path:
	blkChainPath := filepath.Join(obj.rootPath, path.String())

	//create the config:
	conf := obj.conf.SetRoot(blkChainPath)

	//retrieve the node key:
	pk, pkErr := obj.retrieveNodeKey(conf)
	if pkErr != nil {
		return nil, pkErr
	}

	//retrieve the genesis:
	gen, genErr := obj.retrieveGenesis(conf)
	if genErr != nil {
		return nil, genErr
	}

	//retrieve the private validator:
	pv, pvErr := obj.retrievePV(conf)
	if pvErr != nil {
		return nil, pvErr
	}

	//create the blockchain instance:
	blkchain := createBlockchain(gen, pv, pk)
	return blkchain, nil

}

// Save saves a blockchain
func (obj *blockchainService) Save(blkChain Blockchain) error {
	//create the dirPath:
	gen := blkChain.GetGenesis()
	dirPath := filepath.Join(obj.rootPath, gen.GetPath().String())

	//create the config:
	conf := obj.conf.SetRoot(dirPath)

	//save the node key:
	saveNodeKeyErr := obj.saveNodeKey(blkChain.GetPK(), conf)
	if saveNodeKeyErr != nil {
		return saveNodeKeyErr
	}

	//save the genesis:
	saveGenErr := obj.saveGenesis(gen, conf)
	if saveGenErr != nil {
		return saveGenErr
	}

	//save the private validator:
	savePrivValErr := obj.savePV(blkChain.GetPV(), conf)
	if savePrivValErr != nil {
		return savePrivValErr
	}

	return nil
}

// Delete deletes a blockchain
func (obj *blockchainService) Delete(path Path) error {
	dirPath := filepath.Join(obj.rootPath, path.String())
	_, blkChainErr := obj.Retrieve(path)
	if blkChainErr != nil {
		return blkChainErr
	}

	rmErr := obj.deleteParentIfEmpty(dirPath)
	return rmErr
}

func (obj *blockchainService) deleteParentIfEmpty(path string) error {
	dirPath := filepath.Dir(path)
	files, filesErr := ioutil.ReadDir(dirPath)
	if filesErr != nil {
		return filesErr
	}

	hasFiles := false
	for _, oneFile := range files {
		if oneFile.IsDir() {
			continue
		}

		hasFiles = true
		break
	}

	if hasFiles {
		return nil
	}

	rmErr := os.RemoveAll(dirPath)
	if rmErr != nil {
		return rmErr
	}

	return obj.deleteParentIfEmpty(dirPath)
}

func (obj *blockchainService) retrieveNodeKey(conf *config.Config) (crypto.PrivKey, error) {
	//read the file:
	js, jsErr := ioutil.ReadFile(conf.NodeKeyFile())
	if jsErr != nil {
		return nil, jsErr
	}

	//convert the json to an instance:
	containedPK := new(containedPK)
	convJSErr := cdc.UnmarshalJSON(js, containedPK)
	if convJSErr != nil {
		return nil, convJSErr
	}

	return containedPK.PK, nil
}

func (obj *blockchainService) saveNodeKey(pk crypto.PrivKey, conf *config.Config) error {

	//convert the pk to js:
	cpk := containedPK{
		PK: pk,
	}

	jsPK, jsPKErr := cdc.MarshalJSON(cpk)
	if jsPKErr != nil {
		return jsPKErr
	}

	//write the pk to file:
	pkWriteErr := ioutil.WriteFile(obj.mkdirIfNeeded(conf.NodeKeyFile()), jsPK, 0600)
	if pkWriteErr != nil {
		return pkWriteErr
	}

	return nil
}

func (obj *blockchainService) retrieveGenesis(conf *config.Config) (Genesis, error) {
	//read the file:
	js, jsErr := ioutil.ReadFile(conf.GenesisFile())
	if jsErr != nil {
		return nil, jsErr
	}

	//convert the json to an instance:
	gen := new(genesis)
	convJSErr := cdc.UnmarshalJSON(js, gen)
	if convJSErr != nil {
		return nil, convJSErr
	}

	return gen, nil

}

func (obj *blockchainService) saveGenesis(gen Genesis, conf *config.Config) error {
	//convert the genesis to json:
	jsGen, jsGenErr := cdc.MarshalJSON(gen)
	if jsGenErr != nil {
		return jsGenErr
	}

	//write the genesis to file:
	genWriteErr := ioutil.WriteFile(obj.mkdirIfNeeded(conf.GenesisFile()), jsGen, 0600)
	if genWriteErr != nil {
		return genWriteErr
	}

	return nil
}

func (obj *blockchainService) retrievePV(conf *config.Config) (PrivateValidator, error) {
	//read the file:
	js, jsErr := ioutil.ReadFile(conf.PrivValidatorFile())
	if jsErr != nil {
		return nil, jsErr
	}

	//convert the json to instance:
	pv := new(privateValidator)
	convJSErr := cdc.UnmarshalJSON(js, pv)
	if convJSErr != nil {
		return nil, convJSErr
	}

	return pv, nil
}

func (obj *blockchainService) savePV(pv PrivateValidator, conf *config.Config) error {
	//convert the private validator to json:
	privValJS, privValJSErr := cdc.MarshalJSON(pv)
	if privValJSErr != nil {
		return privValJSErr
	}

	//write the private validator to file:
	privValErr := ioutil.WriteFile(obj.mkdirIfNeeded(conf.PrivValidatorFile()), privValJS, 0600)
	if privValErr != nil {
		return privValErr
	}

	return nil
}

func (obj *blockchainService) mkdirIfNeeded(filePath string) string {
	os.MkdirAll(filepath.Dir(filePath), 0777)
	return filePath
}

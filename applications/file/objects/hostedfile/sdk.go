package hostedfile

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/file/objects/file"
	"github.com/xmnservices/xmnsuite/applications/file/objects/host"
)

// HostedFile represents an hosted file
type HostedFile interface {
	ID() *uuid.UUID
	File() file.File
	Host() host.Host
}

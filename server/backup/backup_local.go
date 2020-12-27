package backup

import (
	"errors"
	"github.com/apex/log"
	"github.com/pterodactyl/wings/server/filesystem"
	"os"
)

type LocalBackup struct {
	Backup
}

var _ BackupInterface = (*LocalBackup)(nil)

// Locates the backup for a server and returns the local path. This will obviously only
// work if the backup was created as a local backup.
func LocateLocal(uuid string) (*LocalBackup, os.FileInfo, error) {
	b := &LocalBackup{
		Backup{
			Uuid:   uuid,
			Ignore: "",
		},
	}

	st, err := os.Stat(b.Path())
	if err != nil {
		return nil, nil, err
	}

	if st.IsDir() {
		return nil, nil, errors.New("invalid archive, is directory")
	}

	return b, st, nil
}

// Removes a backup from the system.
func (b *LocalBackup) Remove() error {
	return os.Remove(b.Path())
}

// Generates a backup of the selected files and pushes it to the defined location
// for this instance.
func (b *LocalBackup) Generate(basePath, ignore string) (*ArchiveDetails, error) {
	a := &filesystem.Archive{
		BasePath: basePath,
		Ignore:   ignore,
	}

	l := log.WithFields(log.Fields{
		"backup_id": b.Uuid,
		"adapter":   "local",
	})

	l.Info("attempting to create backup..")
	if err := a.Create(b.Path()); err != nil {
		return nil, err
	}
	l.Info("created backup successfully.")

	return b.Details(), nil
}

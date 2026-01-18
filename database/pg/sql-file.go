package pg

import (
	"io"
	"io/fs"
	"strings"
)

type OptFunc func(*Opts)

func WithSkipNotExistingFile(opts *Opts) {
	opts.skipNotExistingFile = true
}

func WithQuerySeparator(querySeparator string) OptFunc {
	return func(opts *Opts) {
		opts.querySeparator = querySeparator
	}
}

func defaultOpts(filePath string) Opts {
	return Opts{
		filePath:            filePath,
		skipNotExistingFile: false,
		querySeparator:      ";",
	}
}

type Opts struct {
	filePath            string
	skipNotExistingFile bool
	querySeparator      string
}
type SqlFile struct {
	fs   fs.FS
	db   *PgDb
	opts Opts
}

func NewSqlFile(fs fs.FS, db *PgDb, path string, opts ...OptFunc) *SqlFile {
	o := defaultOpts(path)
	for _, fn := range opts {
		fn(&o)
	}
	return &SqlFile{fs: fs, db: db, opts: o}
}

func (es *SqlFile) Execute() error {
	open, err := es.fs.Open(es.opts.filePath)
	if err != nil {
		_, ok := err.(*fs.PathError)
		if ok && es.opts.skipNotExistingFile {
			return nil
		}
		return err
	}
	sqlStatements, err := io.ReadAll(open)
	if err != nil {
		return err
	}
	sqlArray := strings.Split(string(sqlStatements), es.opts.querySeparator)
	return es.db.ExecuteSqls(sqlArray)
}

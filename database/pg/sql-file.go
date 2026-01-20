package pg

import (
	"io"
	"io/fs"
	"strings"
)

type SqlFile struct {
	fs      fs.FS
	db      *PgDb
	path    string
	options SqlFileOptions
}

func NewSqlFile(fs fs.FS, db *PgDb, path string, opts ...SqlFileOptionsFunc) *SqlFile {
	o := defaultOptions()
	for _, fn := range opts {
		fn(&o)
	}
	return &SqlFile{fs: fs, db: db, path: path, options: o}
}

type SqlFileOptions struct {
	skipNotExistingFile bool
	querySeparator      string
}

type SqlFileOptionsFunc func(*SqlFileOptions)

func WithSkipNotExistingFile(opts *SqlFileOptions) {
	opts.skipNotExistingFile = true
}

func WithQuerySeparator(querySeparator string) SqlFileOptionsFunc {
	return func(opts *SqlFileOptions) {
		opts.querySeparator = querySeparator
	}
}

func defaultOptions() SqlFileOptions {
	return SqlFileOptions{
		skipNotExistingFile: false,
		querySeparator:      ";",
	}
}

func (es *SqlFile) Execute() error {
	open, err := es.fs.Open(es.path)
	if err != nil {
		_, ok := err.(*fs.PathError)
		if ok && es.options.skipNotExistingFile {
			return nil
		}
		return err
	}
	sqlStatements, err := io.ReadAll(open)
	if err != nil {
		return err
	}
	sqlArray := strings.Split(string(sqlStatements), es.options.querySeparator)
	return es.db.ExecuteSqls(sqlArray)
}

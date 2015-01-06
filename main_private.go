package gorm

import (
	"regexp"
	"time"
)

func (s *DB) clone() *DB {
	db := DB{db: s.db, parent: s.parent, logMode: s.logMode, Value: s.Value, Error: s.Error, values: map[string]interface{}{}}

	for key, value := range s.values {
		db.values[key] = value
	}

	if s.search == nil {
		db.search = &search{}
	} else {
		db.search = s.search.clone()
	}

	db.search.db = &db
	return &db
}

func (s *DB) err(err error) error {
	if err != nil {
		if err != RecordNotFound {
			if s.logMode == 0 {
				go s.print(fileWithLineNum(), err)
			} else {
				s.log(err)
			}
			if regexp.MustCompile(`^sql: Scan error on column index`).MatchString(err.Error()) {
				return nil
			}
		}
		s.Error = err
	}
	return err
}

func (s *DB) print(v ...interface{}) {
	s.parent.logger.(logger).Print(v...)
}

func (s *DB) log(v ...interface{}) {
	if s != nil && s.logMode == 2 {
		s.print(append([]interface{}{"log", fileWithLineNum()}, v...)...)
	}
}

func (s *DB) slog(sql string, t time.Time, vars ...interface{}) {
	if s.logMode == 2 {
		s.print("sql", fileWithLineNum(), NowFunc().Sub(t), sql, vars)
	}
}

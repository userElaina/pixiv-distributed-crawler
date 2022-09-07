package pic

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type datas struct {
	db *sql.DB
}

func newDatas(p string) *datas {
	db, e := sql.Open("sqlite3", p)
	if e != nil {
		fmt.Println(e)
	}
	return &datas{db}
}

func (me *datas) Close() {
	me.db.Close()
}

func (me *datas) exec0(s string) (e error) {
	stmt, e := me.db.Prepare(s)
	if e != nil {
		return
	}
	_, e = stmt.Exec()
	return
}

func (me *datas) ePic(pid int) (bool, error) {
	if pid <= 0 {
		return false, nil
	}
	rows, e := me.db.Query("SELECT 1 FROM pic WHERE pid = ?", pid)
	if e != nil {
		return false, e
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (me *datas) iPic(info *InfoPic) (e error) {
	if info.pid <= 0 {
		return
	}
	stmt, e := me.db.Prepare("INSERT INTO pic VALUES (?,?,?,?,?,?,?,?,?)")
	if e != nil {
		return
	}
	_, e = stmt.Exec(info.pid, info.uid, info.series, info.stat, info.views, info.bookmark, info.likes, info.timestamp, info.title)
	return
}

func (me *datas) eAlias(i, j string) (bool, error) {
	if i == "" || j == "" {
		return false, nil
	}
	rows, e := me.db.Query("SELECT 1 FROM alias WHERE i = ? AND j = ?", i, j)
	if e != nil {
		return false, e
	}
	if rows.Next() {
		rows.Close()
		return true, nil
	}
	rows, e = me.db.Query("SELECT 1 FROM alias WHERE i = ? AND j = ?", j, i)
	if e != nil {
		return false, e
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (me *datas) iAlias(i, j string) (e error) {
	if i == "" || j == "" {
		return
	}
	stmt, e := me.db.Prepare("INSERT INTO alias VALUES (?,?)")
	if e != nil {
		return
	}
	_, e = stmt.Exec(i, j)
	return
}

func (me *datas) eiAlias(i, j string) (e error) {
	if i == "" || j == "" {
		return
	}
	b, e := me.eAlias(i, j)
	if e != nil {
		return
	}
	if b {
		return
	}
	e = me.iAlias(i, j)
	return
}

func (me *datas) ePicTag(pid int, name string) (bool, error) {
	if pid <= 0 || name == "" {
		return false, nil
	}
	rows, e := me.db.Query("SELECT 1 FROM pictag WHERE pid = ? AND name = ?", pid, name)
	if e != nil {
		return false, e
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (me *datas) iPicTag(pid int, name string) (e error) {
	if pid <= 0 || name == "" {
		return
	}
	stmt, e := me.db.Prepare("INSERT INTO pictag VALUES (?,?)")
	if e != nil {
		return
	}
	_, e = stmt.Exec(pid, name)
	return
}

func (me *datas) eiPicTag(pid int, name string) (e error) {
	if pid <= 0 || name == "" {
		return
	}
	b, e := me.ePicTag(pid, name)
	if e != nil {
		return
	}
	if b {
		return
	}
	e = me.iPicTag(pid, name)
	return
}

func (me *datas) loadInfoPic(info *InfoPic) (e error) {
	b, e := me.ePic(info.pid)
	if e != nil {
		return
	}
	if b {
		b = true
		// return
	}

	e = me.iPic(info)
	if e != nil {
		return
	}

	for i, j := range info.maps {
		e = me.eiAlias(i, j)
		if e != nil {
			return
		}
		e = me.eiPicTag(info.pid, i)
		if e != nil {
			return
		}
	}

	return
}

func (me *datas) initDB() (e error) {
	e = me.exec0("PRAGMA FOREIGN_KEYS = ON")
	if e != nil {
		return
	}

	e = me.exec0(`CREATE TABLE IF NOT EXISTS alias(
		i VARCHAR(64) NOT NULL,
		j VARCHAR(64) NOT NULL,
		PRIMARY KEY(i,j),
		CHECK (i != ""),
		CHECK (j != "")
	)`)
	if e != nil {
		return
	}

	e = me.exec0(`CREATE TABLE IF NOT EXISTS pic(
		pid INTEGER PRIMARY KEY NOT NULL,
		uid INTEGER NOT NULL,
		series INTEGER NOT NULL,
		stat INTEGER NOT NULL,
		views INTEGER NOT NULL,
		bookmark INTEGER NOT NULL,
		likes INTEGER NOT NULL,
		timestamp INTEGER NOT NULL,
		title VARCHAR(64) NOT NULL,
		CHECK (pid > 0),
		CHECK (uid > 0),
		CHECK (views >= 0),
		CHECK (bookmark >= 0),
		CHECK (likes >= 0)
	)`)
	if e != nil {
		return
	}

	e = me.exec0(`CREATE TABLE IF NOT EXISTS pictag(
		pid INTEGER NOT NULL,
		name VARCHAR(64) NOT NULL,
		FOREIGN KEY(pid) REFERENCES pic(pid),
		PRIMARY KEY(pid,name),
		CHECK (pid > 0),
		CHECK (name != "")
	)`)
	if e != nil {
		return
	}

	fmt.Println("INITDB SUCC")
	return
}

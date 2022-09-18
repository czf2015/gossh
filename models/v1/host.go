package models

import (
	"gossh/libs/db"
)

type Host struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Port        int    `json:"port"`
	FontSize    int    `json:"font_size"`
	Background  string `json:"background"`
	Foreground  string `json:"foreground"`
	CursorColor string `json:"cursor_color"`
	FontFamily  string `json:"font_family"`
	CursorStyle string `json:"cursor_style"`
	Shell       string `json:"shell"`
}

func (host *Host) Select() ([]Host, error) {
	rows, err := db.GetDB().Query(`select Id, Name, Address, User, Pwd, Port,FontSize, Background, Foreground, CursorColor, FontFamily, CursorStyle, Shell from host`)
	var hostList []Host
	if err != nil {
		return hostList, err
	}
	for rows.Next() {
		var h = new(Host)
		err = rows.Scan(&h.Id, &h.Name, &h.Address, &h.User, &h.Pwd, &h.Port, &h.FontSize, &h.Background, &h.Foreground, &h.CursorColor, &h.FontFamily, &h.CursorStyle, &h.Shell)
		if err != nil {
			return hostList, err
		}
		hostList = append(hostList, *h)
	}
	_ = rows.Close()
	return hostList, nil
}

func (host *Host) Insert(name, address, user, pwd string, port, fontSize int, background, foreground, cursorColor, fontFamily, cursorStyle, shell string) (int64, error) {
	insertSql := `INSERT INTO host(Name, Address, User, Pwd, Port, FontSize, Background, Foreground, CursorColor, FontFamily, CursorStyle, Shell)  values(?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := db.GetDB().Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(name, address, user, pwd, port, fontSize, background, foreground, cursorColor, fontFamily, cursorStyle, shell)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId() //返回新增的id号
	if err != nil {
		return 0, err
	}
	return id, err
}

func (host *Host) Update(id int, name, address, user, pwd string, port, fontSize int, background, foreground, cursorColor, fontFamily, cursorStyle, shell string) (int64, error) {

	stmt, err := db.GetDB().Prepare(`update host set Name=?, Address=?, User=?, Pwd=?, Port=?, FontSize=?, Background=?, Foreground=?, CursorColor=?, FontFamily=?, CursorStyle=?, Shell=?  where id=?`)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(name, address, user, pwd, port, fontSize, background, foreground, cursorColor, fontFamily, cursorStyle, shell, id)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, err
}

func (host *Host) Delete(id int) (int64, error) {
	stmt, err := db.GetDB().Prepare(`delete from host where id=?`)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(id) //将想删除的id输入进去就可以删除输入的id
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected() //几条数据受影响：返回int64类型数据
	if err != nil {
		return 0, err
	}
	return affect, err
}

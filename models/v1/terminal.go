package models

import (
	"gossh/libs/db"
)

type Terminal struct {
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

func (terminal *Terminal) Select() ([]Terminal, error) {
	rows, err := db.GetDB().Query(`select Id, Name, Address, User, Pwd, Port,FontSize, Background, Foreground, CursorColor, FontFamily, CursorStyle, Shell from terminal`)
	var terminalList []Terminal
	if err != nil {
		return terminalList, err
	}
	for rows.Next() {
		var h = new(Terminal)
		err = rows.Scan(&h.Id, &h.Name, &h.Address, &h.User, &h.Pwd, &h.Port, &h.FontSize, &h.Background, &h.Foreground, &h.CursorColor, &h.FontFamily, &h.CursorStyle, &h.Shell)
		if err != nil {
			return terminalList, err
		}
		terminalList = append(terminalList, *h)
	}
	_ = rows.Close()
	return terminalList, nil
}

func (terminal *Terminal) Insert(name, address, user, pwd string, port, fontSize int, background, foreground, cursorColor, fontFamily, cursorStyle, shell string) (int64, error) {
	insertSql := `INSERT INTO terminal(Name, Address, User, Pwd, Port, FontSize, Background, Foreground, CursorColor, FontFamily, CursorStyle, Shell)  values(?,?,?,?,?,?,?,?,?,?,?,?)`
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

func (terminal *Terminal) Update(id int, name, address, user, pwd string, port, fontSize int, background, foreground, cursorColor, fontFamily, cursorStyle, shell string) (int64, error) {

	stmt, err := db.GetDB().Prepare(`update terminal set Name=?, Address=?, User=?, Pwd=?, Port=?, FontSize=?, Background=?, Foreground=?, CursorColor=?, FontFamily=?, CursorStyle=?, Shell=?  where id=?`)
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

func (terminal *Terminal) Delete(id int) (int64, error) {
	stmt, err := db.GetDB().Prepare(`delete from terminal where id=?`)
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

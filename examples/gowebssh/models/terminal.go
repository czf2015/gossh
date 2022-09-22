package models

type Terminal struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Port        int    `json:"port"`
	Shell       string `json:"shell"`
	FontFamily  string `json:"font_family"`
	FontSize    int    `json:"font_size"`
	Foreground  string `json:"foreground"`
	Background  string `json:"background"`
	CursorColor string `json:"cursor_color"`
	CursorStyle string `json:"cursor_style"`
}

package binder

import (
	"github.com/zan8in/gologger"
)

func ShowBanner(VERSION string, CONTENT string) {
	gologger.Print().Msgf("\n|||\tB I N D E R\t|||\t%s\n大白哥捆绑器二开版本\tauthor:v01d\n%s\n\n", VERSION, CONTENT)
}

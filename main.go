package main;

import (
  "parse"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
)

func main() {
  f, err := ioutil.ReadFile(os.Args[len(os.Args) - 1])
  if err != nil {
    fmt.Printf("File Error: %v", err.Error())
  }

  switch os.Args[1] {
    case "lex": lexOnly(f)
    default: compile(f)
  }
}

func lexOnly(f []byte) {
  for tok := range parse.Lex(string(f)) {
    fmt.Printf("%v\n", tok)
  }
}

func compile(f []byte) {
  toks := parse.Lex(string(f))
  ast, er := parse.Parse(toks)

  if er != nil {
    fmt.Println(er.Error())
  } else {
    fmt.Printf(newLineOn(ast.String(), ";", "{", "}"))
  }
}

func newLineOn(s string, breakChars ...string) string {
  for _, breakChar := range breakChars {
    s = strings.Replace(s, breakChar, breakChar + "\n", -1)
  }

  return s
}

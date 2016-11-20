package main;

import (
  "parse"
  "compiler"
  "fmt"
  "asm"
  "io/ioutil"
  "strings"
  "flag"
)

func main() {
  lexPt := flag.Bool("lex", false, "Only lex the file into tokens. Don't parse.")
  parsePt := flag.Bool("parse", false, "Only lex & parse the file into an AST. Don't compile.")
  strBfPt := flag.Bool("str", false, "Show as BF descriptors.")

  flag.Parse()
  tail := flag.Args()

  f, err := ioutil.ReadFile(tail[0])
  if err != nil {
    fmt.Printf("File Error: %v", err.Error())
  }

  switch {
  case *lexPt: printLexicons(f)
  case *parsePt: printAst(f)
  default: compile(f, *strBfPt)
  }
}

func printLexicons(f []byte) {
  for tok := range parse.Lex(string(f)) {
    fmt.Printf("%v\n", tok)
  }
}

func printAst(f []byte) {
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

func compile(f []byte, strBf bool) {
  toks := parse.Lex(string(f))
  ast, er := parse.Parse(toks)

  if er != nil {
    fmt.Println(er.Error())
    return
  }

  assembler, out := asm.New()
  go func() {
    compiler.Compile(assembler, ast)
    close(out)
  }()

  for node := range out {
    if strBf {
      fmt.Println(node.String())
    } else {
      fmt.Printf("%s", node.ToBF())
    }
  }
  fmt.Println("")
}

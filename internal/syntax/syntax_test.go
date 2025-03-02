package syntax_test

// TODO(@FollowTheProcess): Integration tests that test the parse pipeline in phases using txtar archives
// lex: Take raw (valid) Lox source code as an entire program in src.lox, then tokenise and validate against tokens.txt
// parse: Take the tokens from the lex phase, parse and produce an AST, serialise it and dump to parsed.txt
// eval: Take the AST from parse and run it through the interpreter, dumping the raw string output to output.txt

package frontend

import "testing"

import (
	"github.com/timtadh/lexmachine/machines"
	"github.com/timtadh/lexmachine/inst"
)

func TestParse(t *testing.T) {
	ast, err := Parse([]byte("ab(a|c|d)?we*\\\\\\[\\..[s-f]+|qyx"))
	if err != nil {
		t.Error(err)
	}
	parsed := "(Match (Alternation (Concat (Concat (Character a), (Character b)), (? (Alternation (Character a), (Alternation (Character c), (Character d)))), (* (Concat (Character w), (Character e))), (+ (Concat (Character \\), (Character [), (Character .), (Range 0 255), (Range 115 102)))), (Concat (Character q), (Character y), (Character x))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
}

func t_match(program inst.InstSlice, text string, t *testing.T) {
	expected := []machines.Match{machines.Match{len(program)-1, []byte(text)}}
	i := 0
	scan := machines.LexerEngine(program, []byte(text))
	for m, err, scan := scan(); scan != nil; m, err, scan = scan() {
		t.Log("match", m)
		if err != nil {
			t.Error("error", err)
		} else if !m.Equals(&expected[i]) {
			t.Error(m, expected[i])
		}
		i++
	}
	if i != len(expected) {
		t.Error("unconsumed matches", expected[i:])
	}
}

func TestParseConcatAlts(t *testing.T) {
	ast, err := Parse([]byte("A|((C|D|E)(F|G)(H|I)B)"))
	if err != nil {
		t.Error(err)
	}
	parsed := "(Match (Alternation (Character A), (Concat (Alternation (Character C), (Alternation (Character D), (Character E))), (Alternation (Character F), (Character G)), (Alternation (Character H), (Character I)), (Character B))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "A", t)
	t_match(program, "CFHB", t)
	t_match(program, "CFIB", t)
	t_match(program, "CGHB", t)
	t_match(program, "CGIB", t)
	t_match(program, "DFHB", t)
	t_match(program, "DFIB", t)
	t_match(program, "DGHB", t)
	t_match(program, "DGIB", t)
	t_match(program, "EFHB", t)
	t_match(program, "EFIB", t)
	t_match(program, "EGHB", t)
	t_match(program, "EGIB", t)
}

func TestParseConcatAltMaybes(t *testing.T) {
	ast, err := Parse([]byte("((A?)?|(B|C))(D|E?)"))
	if err != nil {
		t.Error(err)
	}
	parsed := "(Match (Concat (Alternation (? (? (Character A))), (Alternation (Character B), (Character C))), (Alternation (Character D), (? (Character E)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "", t)
	t_match(program, "E", t)
	t_match(program, "D", t)
	t_match(program, "A", t)
	t_match(program, "AE", t)
	t_match(program, "AD", t)
	t_match(program, "B", t)
	t_match(program, "BE", t)
	t_match(program, "BD", t)
	t_match(program, "C", t)
	t_match(program, "CE", t)
	t_match(program, "CD", t)
}


func TestParseConcatAltPlus(t *testing.T) {
	ast, err := Parse([]byte("(A|(B|C))+(D|E?)"))
	if err != nil {
		t.Error(err)
	}
	parsed := "(Match (Concat (+ (Alternation (Character A), (Alternation (Character B), (Character C)))), (Alternation (Character D), (? (Character E)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "A", t)
	t_match(program, "AAA", t)
	t_match(program, "AAABBCC", t)
	t_match(program, "AAABBCC", t)
	t_match(program, "AAABBCCD", t)
}

func TestParseConcatAltStar(t *testing.T) {
	ast, err := Parse([]byte("(A|[C-G])*(X|Y?)"))
	if err != nil {
		t.Error(err)
	}
	parsed := "(Match (Concat (* (Alternation (Character A), (Range 67 71))), (Alternation (Character X), (? (Character Y)))))"
	if ast.String() != parsed {
		t.Log(ast.String())
		t.Log(parsed)
		t.Error("Did not parse correctly")
	}
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "", t)
	t_match(program, "X", t)
	t_match(program, "Y", t)
	t_match(program, "A", t)
	t_match(program, "AAA", t)
	t_match(program, "AAACC", t)
	t_match(program, "AAACC", t)
	t_match(program, "AAACCFFF", t)
	t_match(program, "CAACCGEDFX", t)
}


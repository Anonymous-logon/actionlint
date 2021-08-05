package actionlint

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func TestExprSemanticsCheckRealWorld(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "bench", "expressions.txt"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		expr := s.Text()
		l := NewExprLexer(expr + "}}")
		p := NewExprParser()
		root, err := p.Parse(l)
		if err != nil {
			t.Errorf("%q caused parse error: %v", expr, err)
			continue
		}
		c := NewExprSemanticsChecker(true)
		c.Check(root)
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkExprRealWorld(b *testing.B) {
	f, err := os.Open(filepath.Join("testdata", "bench", "expressions.txt"))
	if err != nil {
		b.Fatal(err)
	}

	exprs := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		exprs = append(exprs, s.Text()+"}}")
	}
	f.Close()
	if err := s.Err(); err != nil {
		b.Fatal(err)
	}

	b.Run("Lexer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, expr := range exprs {
				l := NewExprLexer(expr)
				for {
					t := l.Next()
					if l.lexErr != nil {
						b.Fatal(l.lexErr)
					}
					if t.Kind == TokenKindEnd {
						break
					}
				}
			}
		}
	})

	b.Run("Parser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, expr := range exprs {
				if _, err := NewExprParser().Parse(NewExprLexer(expr)); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("SematnicsCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, expr := range exprs {
				root, err := NewExprParser().Parse(NewExprLexer(expr + "}}"))
				if err != nil {
					b.Fatalf("%q caused parse error: %v", expr, err)
				}
				NewExprSemanticsChecker(true).Check(root)
			}
		}
	})

	b.Run("SematnicsCheckWithoutUntrustedInputsCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, expr := range exprs {
				root, err := NewExprParser().Parse(NewExprLexer(expr + "}}"))
				if err != nil {
					b.Fatalf("%q caused parse error: %v", expr, err)
				}
				NewExprSemanticsChecker(false).Check(root)
			}
		}
	})
}

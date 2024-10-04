package token

const Collateral = "Collateral"

type Token string

const (
	A Token = "TokenA"
	B Token = "TokenB"
)

func (t Token) Complement() Token {
	if t == A {
		return B
	}
	return A
}

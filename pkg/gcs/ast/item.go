package ast

import "fmt"

// Token represents a token or text string returned from the scanner.
type Token struct {
	Typ  TokenType // The type of this item.
	Pos  Pos       // The starting position, in bytes, of this item in the input string.
	Val  string    // The value of this item.
	Line int       // The line number at the start of this item.
}

func (i Token) String() string {
	switch {
	case i.Typ == ItemEOF:
		return "EOF"
	case i.Typ == ItemError:
		return i.Val
	case i.Typ == ItemTerminateLine:
		return ";"
	case i.Typ > ItemTerminateLine && i.Typ < ItemKeyword:
		return i.Val
	case i.Typ > ItemKeyword:
		return fmt.Sprintf("<%s>", i.Val)
		// case len(i.val) > 10:
		// 	return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// TokenType identifies the type of lex items.
type TokenType int

const (
	ItemError TokenType = iota // error occurred; value is text of error

	ItemEOF
	ItemTerminateLine    // \n to denote end of a line
	ItemAssign           // equals ('=') introducing an assignment
	ItemComma            // coma (,) used to break up list of ident
	ItemLeftParen        // '('
	ItemRightParen       // ')'
	ItemLeftSquareParen  // '['
	ItemRightSquareParen // ']'
	ItemLeftBrace        // '{'
	ItemRightBrace       // '}'
	ItemColon            // ':'
	ItemPlus             // '+'
	ItemMinus            // '-'
	ItemAsterisk         // '*'
	ItemForwardSlash     // '/'
	// following is logic operator
	ItemLogicOP // used only to delimit logical operation
	LogicNot    // !
	LogicAnd    // && keyword
	LogicOr     // || keyword
	// following is comparison operator
	ItemCompareOp        // used only to delimi comparison operators
	OpEqual              // == keyword
	OpNotEqual           // != keyword
	OpGreaterThan        // > keyword
	OpGreaterThanOrEqual // >= keyword
	OpLessThan           // < keyword
	OpLessThanOrEqual    // <= keyword
	ItemDot              // the cursor, spelled '.'
	// item types
	ItemTypes
	ItemField      // alphanumeric identifier starting with '.'
	ItemIdentifier // alphanumeric identifier not starting with '.'
	ItemNumber     // simple number
	ItemBool       // boolean
	ItemString     // string, including quotes
	// Keywords appear after all the rest.
	ItemKeyword        // used only to delimit the keywords
	KeywordLet         // let
	KeywordWhile       // while
	KeywordIf          // if
	KeywordElse        // else
	KeywordFn          // fn
	KeywordSwitch      // switch
	KeywordCase        // case
	KeywordDefault     // default
	KeywordBreak       // break
	KeywordContinue    // continue
	KeywordFallthrough // fallthrough
	KeywordReturn      // return
	KeywordFor         // for
	// Keywords after this are specific to Genshin (i.e. not generic scripting keywords)
	// These are special char related keywords
	KeywordOptions           // options
	KeywordAdd               // add
	KeywordChar              // char
	KeywordStats             // stats
	KeywordWeapon            // weapon
	KeywordSet               // set
	KeywordLvl               // lvl
	KeywordRefine            // refine
	KeywordCons              // cons
	KeywordTalent            // talent
	KeywordStartHP           // start_hp
	KeywordCount             // count
	KeywordParams            // params
	KeywordLabel             // label
	KeywordUntil             // until
	KeywordActive            // active
	KeywordTarget            // target
	KeywordResist            // resist
	KeywordEnergy            // energy
	KeywordParticleThreshold // particle_threshold
	KeywordParticleDropCount // particle_drop_count

	// Keywords specific to gcsim appears after this
	ItemKeys
	ItemStatKey      // stats: def%, def, etc..
	ItemElementKey   // elements: pyro, hydro, etc..
	ItemCharacterKey // characters: albedo, amber, etc..
	ItemActionKey    // actions: skill, burst, attack, charge, etc...
)

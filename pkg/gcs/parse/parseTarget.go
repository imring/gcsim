package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r enemy.EnemyProfile
	r.Resist = make(map[attributes.Element]float64)
	r.ParticleElement = attributes.NoElement
	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemIdentifier:
			switch n.Val {
			case "pos": //pos will end up defaulting to 0,0 if not set
				//pos=1.00,2,00
				item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				x, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				item, err = p.acceptSeqReturnLast(ast.ItemComma, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				y, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.X = x
				r.Pos.Y = y
			case "radius":
				item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.R = amt
			default:
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
			}
		case ast.KeywordLvl:
			n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case ast.ItemStatKey:
			//should be hp
			if ast.StatKeys[n.Val] != attributes.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
			}
			n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				if err != nil {
					return nil, err
				}
				p.res.Settings.DamageMode = true
			}
		case ast.KeywordResist:
			//this sets all resistance
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			//TODO: make this more elegant...
			r.Resist[attributes.Electro] += amt
			r.Resist[attributes.Cryo] += amt
			r.Resist[attributes.Hydro] += amt
			r.Resist[attributes.Physical] += amt
			r.Resist[attributes.Pyro] += amt
			r.Resist[attributes.Geo] += amt
			r.Resist[attributes.Dendro] += amt
			r.Resist[attributes.Anemo] += amt
		case ast.KeywordParticleThreshold:
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropThreshold = amt
		case ast.KeywordParticleDropCount:
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropCount = amt
		case ast.ItemElementKey:
			s := n.Val
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			r.Resist[ast.EleKeys[s]] += amt
		case ast.ItemTerminateLine:
			p.res.Targets = append(p.res.Targets, r)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing target")
}

package ast

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r enemy.EnemyProfile
	r.Resist = make(map[attributes.Element]float64)
	r.ParticleElement = attributes.NoElement
	for n := p.next(); n.Typ != ItemEOF; n = p.next() {
		switch n.Typ {
		case ItemIdentifier:
			switch n.Val {
			case "pos": //pos will end up defaulting to 0,0 if not set
				//pos=1.00,2,00
				item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
				if err != nil {
					return nil, err
				}
				x, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				item, err = p.acceptSeqReturnLast(ItemComma, ItemNumber)
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
				item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.R = amt
			default:
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
		case KeywordLvl:
			n, err = p.acceptSeqReturnLast(ItemAssign, ItemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case ItemStatKey:
			//should be hp
			if statKeys[n.Val] != attributes.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
			n, err = p.acceptSeqReturnLast(ItemAssign, ItemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				if err != nil {
					return nil, err
				}
				p.res.Settings.DamageMode = true
			}
		case KeywordResist:
			//this sets all resistance
			item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
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
		case KeywordParticleThreshold:
			item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropThreshold = amt
		case KeywordParticleDropCount:
			item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropCount = amt
		case ItemElementKey:
			s := n.Val
			item, err := p.acceptSeqReturnLast(ItemAssign, ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			r.Resist[eleKeys[s]] += amt
		case ItemTerminateLine:
			p.res.Targets = append(p.res.Targets, r)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing target")
}

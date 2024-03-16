// Code generated by "pipeline"; DO NOT EDIT.
package bennett

import (
	_ "embed"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/prototext"
)

//go:embed data_gen.textproto
var pbData []byte
var base *model.AvatarData

func init() {
	base = &model.AvatarData{}
	err := prototext.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
}

func (x *char) Data() *model.AvatarData {
	return base
}

var (
	attack = [][]float64{
		attack_1,
		attack_2,
		attack_3,
		attack_4,
		attack_5,
	}
	skillHold = [][][]float64{
		skillHold_1,
		skillHold_2,
	}
)

var (
	// attack: attack_1 = [0]
	attack_1 = []float64{
		0.44547998905181885,
		0.48173999786376953,
		0.5180000066757202,
		0.5698000192642212,
		0.6060600280761719,
		0.6474999785423279,
		0.704479992389679,
		0.76146000623703,
		0.8184400200843811,
		0.8805999755859375,
		0.9427599906921387,
		1.0049200057983398,
		1.067080020904541,
		1.1292400360107422,
		1.1914000511169434,
	}
	// attack: attack_2 = [1]
	attack_2 = []float64{
		0.4274199903011322,
		0.4622099995613098,
		0.4970000088214874,
		0.5467000007629395,
		0.5814899802207947,
		0.6212499737739563,
		0.6759200096130371,
		0.7305899858474731,
		0.785260021686554,
		0.8449000120162964,
		0.9045400023460388,
		0.9641799926757812,
		1.0238200426101685,
		1.0834599733352661,
		1.1431000232696533,
	}
	// attack: attack_3 = [2]
	attack_3 = []float64{
		0.5461000204086304,
		0.5905500054359436,
		0.6349999904632568,
		0.6984999775886536,
		0.7429500222206116,
		0.793749988079071,
		0.8636000156402588,
		0.9334499835968018,
		1.0032999515533447,
		1.0794999599456787,
		1.1556999683380127,
		1.2318999767303467,
		1.3080999851226807,
		1.3842999935150146,
		1.4605000019073486,
	}
	// attack: attack_4 = [3]
	attack_4 = []float64{
		0.5968400239944458,
		0.6454200148582458,
		0.6940000057220459,
		0.7634000182151794,
		0.8119800090789795,
		0.8675000071525574,
		0.9438400268554688,
		1.0201799869537354,
		1.096519947052002,
		1.179800033569336,
		1.2630800008773804,
		1.3463599681854248,
		1.4296400547027588,
		1.5129200220108032,
		1.5961999893188477,
	}
	// attack: attack_5 = [4]
	attack_5 = []float64{
		0.7189599871635437,
		0.7774800062179565,
		0.8360000252723694,
		0.9196000099182129,
		0.9781200289726257,
		1.0449999570846558,
		1.1369600296020508,
		1.2289199829101562,
		1.3208800554275513,
		1.4212000370025635,
		1.5215200185775757,
		1.621840000152588,
		1.7221599817276,
		1.8224799633026123,
		1.9227999448776245,
	}
	// attack: charge = [5 6]
	charge = [][]float64{
		{
			0.5590000152587891,
			0.6044999957084656,
			0.6499999761581421,
			0.7149999737739563,
			0.7605000138282776,
			0.8125,
			0.8840000033378601,
			0.9555000066757202,
			1.0269999504089355,
			1.1050000190734863,
			1.1829999685287476,
			1.2610000371932983,
			1.3389999866485596,
			1.4170000553131104,
			1.4950000047683716,
		},
		{
			0.6071599721908569,
			0.6565799713134766,
			0.7059999704360962,
			0.7766000032424927,
			0.8260200023651123,
			0.8824999928474426,
			0.9601600170135498,
			1.0378199815750122,
			1.1154799461364746,
			1.2001999616622925,
			1.2849199771881104,
			1.3696399927139282,
			1.454360008239746,
			1.539080023765564,
			1.6238000392913818,
		},
	}
	// attack: collision = [8]
	collision = []float64{
		0.6393240094184875,
		0.6913620233535767,
		0.743399977684021,
		0.8177400231361389,
		0.8697779774665833,
		0.9292500019073486,
		1.011023998260498,
		1.0927979946136475,
		1.1745719909667969,
		1.2637799978256226,
		1.3529880046844482,
		1.442196011543274,
		1.5314040184020996,
		1.6206120252609253,
		1.709820032119751,
	}
	// attack: highPlunge = [10]
	highPlunge = []float64{
		1.59676194190979,
		1.7267309427261353,
		1.8566999435424805,
		2.042370080947876,
		2.1723389625549316,
		2.3208749294281006,
		2.5251119136810303,
		2.72934889793396,
		2.9335858821868896,
		3.1563899517059326,
		3.3791940212249756,
		3.6019980907440186,
		3.8248019218444824,
		4.047605991363525,
		4.270410060882568,
	}
	// attack: lowPlunge = [9]
	lowPlunge = []float64{
		1.2783770561218262,
		1.3824310302734375,
		1.4864850044250488,
		1.635133981704712,
		1.7391870021820068,
		1.858106017112732,
		2.021620035171509,
		2.1851329803466797,
		2.3486459255218506,
		2.527024984359741,
		2.7054030895233154,
		2.8837809562683105,
		3.0621590614318848,
		3.24053692817688,
		3.418915033340454,
	}
	// skill: explosion = [5]
	explosion = []float64{
		1.3200000524520874,
		1.4190000295639038,
		1.5180000066757202,
		1.649999976158142,
		1.7489999532699585,
		1.8480000495910645,
		1.9800000190734863,
		2.111999988555908,
		2.24399995803833,
		2.375999927520752,
		2.507999897003174,
		2.640000104904175,
		2.805000066757202,
		2.9700000286102295,
		3.134999990463257,
	}
	// skill: skill = [0]
	skill = []float64{
		1.3760000467300415,
		1.479200005531311,
		1.5823999643325806,
		1.7200000286102295,
		1.823199987411499,
		1.9263999462127686,
		2.063999891281128,
		2.2016000747680664,
		2.339200019836426,
		2.476799964904785,
		2.6143999099731445,
		2.752000093460083,
		2.9240000247955322,
		3.0959999561309814,
		3.2679998874664307,
	}
	// skill: skillHold_1 = [1 2]
	skillHold_1 = [][]float64{
		{
			0.8399999737739563,
			0.902999997138977,
			0.9660000205039978,
			1.0499999523162842,
			1.1130000352859497,
			1.1759999990463257,
			1.2599999904632568,
			1.343999981880188,
			1.4279999732971191,
			1.5119999647140503,
			1.5959999561309814,
			1.6799999475479126,
			1.784999966621399,
			1.8899999856948853,
			1.9950000047683716,
		},
		{
			0.9200000166893005,
			0.9890000224113464,
			1.0579999685287476,
			1.149999976158142,
			1.218999981880188,
			1.2879999876022339,
			1.3799999952316284,
			1.472000002861023,
			1.5640000104904175,
			1.656000018119812,
			1.7480000257492065,
			1.840000033378601,
			1.9550000429153442,
			2.069999933242798,
			2.184999942779541,
		},
	}
	// skill: skillHold_2 = [3 4]
	skillHold_2 = [][]float64{
		{
			0.8799999952316284,
			0.9459999799728394,
			1.0119999647140503,
			1.100000023841858,
			1.1660000085830688,
			1.2319999933242798,
			1.3200000524520874,
			1.4079999923706055,
			1.496000051498413,
			1.5839999914169312,
			1.6720000505447388,
			1.7599999904632568,
			1.8700000047683716,
			1.9800000190734863,
			2.0899999141693115,
		},
		{
			0.9599999785423279,
			1.031999945640564,
			1.1039999723434448,
			1.2000000476837158,
			1.2719999551773071,
			1.343999981880188,
			1.440000057220459,
			1.5360000133514404,
			1.6319999694824219,
			1.7280000448226929,
			1.8240000009536743,
			1.9199999570846558,
			2.0399999618530273,
			2.1600000858306885,
			2.2799999713897705,
		},
	}
	// burst: burst = [0]
	burst = []float64{
		2.328000068664551,
		2.5025999546051025,
		2.6772000789642334,
		2.9100000858306885,
		3.0845999717712402,
		3.259200096130371,
		3.492000102996826,
		3.7248001098632812,
		3.9576001167297363,
		4.190400123596191,
		4.4232001304626465,
		4.656000137329102,
		4.947000026702881,
		5.23799991607666,
		5.5289998054504395,
	}
	// burst: burstatk = [3]
	burstatk = []float64{
		0.5600000023841858,
		0.6019999980926514,
		0.6439999938011169,
		0.699999988079071,
		0.7419999837875366,
		0.7839999794960022,
		0.8399999737739563,
		0.8960000276565552,
		0.9520000219345093,
		1.0080000162124634,
		1.0640000104904175,
		1.1200000047683716,
		1.190000057220459,
		1.2599999904632568,
		1.3300000429153442,
	}
	// burst: bursthp = [2]
	bursthp = []float64{
		577.3388061523438,
		635.0806884765625,
		697.6344604492188,
		765,
		837.1773681640625,
		914.1665649414062,
		995.9676513671875,
		1082.5804443359375,
		1174.005126953125,
		1270.24169921875,
		1371.2900390625,
		1477.150146484375,
		1587.8221435546875,
		1703.305908203125,
		1823.6015625,
	}
	// burst: bursthpp = [1]
	bursthpp = []float64{
		0.05999999865889549,
		0.06449999660253525,
		0.0689999982714653,
		0.07500000298023224,
		0.0794999971985817,
		0.08399999886751175,
		0.09000000357627869,
		0.09600000083446503,
		0.10199999809265137,
		0.1080000028014183,
		0.11400000005960464,
		0.11999999731779099,
		0.1274999976158142,
		0.13500000536441803,
		0.14249999821186066,
	}
)

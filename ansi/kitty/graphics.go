package kitty

import "errors"

// ErrMissingFile 在文件路径缺失时返回。
var ErrMissingFile = errors.New("missing file path")

// MaxChunkSize 是图像数据的最大块大小。
const MaxChunkSize = 1024 * 4

// Placeholder 是一个特殊的 Unicode 字符，可以用作图像的占位符。
const Placeholder = '\U0010EEEE'

// 图形图像格式。
const (
	// 32 位 RGBA 格式。
	RGBA = 32

	// 24 位 RGB 格式。
	RGB = 24

	// PNG 格式。
	PNG = 100
)

// 压缩类型。
const (
	Zlib = 'z'
)

// 传输类型。
const (
	// 数据直接在转义序列中传输。
	Direct = 'd'

	// 数据在常规文件中传输。
	File = 'f'

	// 使用临时文件并在传输后删除。
	TempFile = 't'

	// 共享内存对象。
	// 对于 POSIX，请参阅 https://pubs.opengroup.org/onlinepubs/9699919799/functions/shm_open.html
	// 对于 Windows，请参阅 https://docs.microsoft.com/en-us/windows/win32/memory/creating-named-shared-memory
	SharedMemory = 's'
)

// 操作类型。
const (
	// 传输图像数据。
	Transmit = 't'
	// TransmitAndPut 传输图像数据并显示（放置）它。
	TransmitAndPut = 'T'
	// 查询终端获取图像信息。
	Query = 'q'
	// 放置（显示）先前传输的图像。
	Put = 'p'
	// 删除图像。
	Delete = 'd'
	// Frame 传输动画帧的数据。
	Frame = 'f'
	// Animate 控制动画。
	Animate = 'a'
	// Compose 组合动画帧。
	Compose = 'c'
)

// 删除类型。
const (
	// 删除屏幕上可见的所有放置。
	DeleteAll = 'a'
	// 删除具有指定 id 的所有图像，使用 i 键指定。
	// 如果您同时为放置 id 指定 p 键，则只会删除具有指定图像 id 和放置 id 的放置。
	DeleteID = 'i'
	// 删除具有指定编号的最新图像，使用 I 键指定。
	// 如果您同时为放置 id 指定 p 键，则只会删除具有指定编号和放置 id 的放置。
	DeleteNumber = 'n'
	// 删除与当前光标位置相交的所有放置。
	DeleteCursor = 'c'
	// 删除动画帧。
	DeleteFrames = 'f'
	// 删除与特定单元格相交的所有放置，单元格使用 x 和 y 键指定。
	DeleteCell = 'p'
	// 删除与具有特定 z-index 的特定单元格相交的所有放置。
	// 单元格和 z-index 使用 x、y 和 z 键指定。
	DeleteCellZ = 'q'
	// 删除 id 大于或等于 x 键值且小于或等于 y 键值的所有图像。
	DeleteRange = 'r'
	// 删除与指定列相交的所有放置，使用 x 键指定。
	DeleteColumn = 'x'
	// 删除与指定行相交的所有放置，使用 y 键指定。
	DeleteRow = 'y'
	// 删除具有指定 z-index 的所有放置，使用 z 键指定。
	DeleteZ = 'z'
)

// Diacritic 返回指定索引处的变音符号字符。如果索引超出范围，则返回第一个变音符号字符。
func Diacritic(i int) rune {
	if i < 0 || i >= len(diacritics) {
		return diacritics[0]
	}
	return diacritics[i]
}

// 来自 https://sw.kovidgoyal.net/kitty/_downloads/f0a0de9ec8d9ff4456206db8e0814937/rowcolumn-diacritics.txt
// 有关进一步解释，请参阅 https://sw.kovidgoyal.net/kitty/graphics-protocol/#unicode-placeholders。
var diacritics = []rune{
	'\u0305',
	'\u030D',
	'\u030E',
	'\u0310',
	'\u0312',
	'\u033D',
	'\u033E',
	'\u033F',
	'\u0346',
	'\u034A',
	'\u034B',
	'\u034C',
	'\u0350',
	'\u0351',
	'\u0352',
	'\u0357',
	'\u035B',
	'\u0363',
	'\u0364',
	'\u0365',
	'\u0366',
	'\u0367',
	'\u0368',
	'\u0369',
	'\u036A',
	'\u036B',
	'\u036C',
	'\u036D',
	'\u036E',
	'\u036F',
	'\u0483',
	'\u0484',
	'\u0485',
	'\u0486',
	'\u0487',
	'\u0592',
	'\u0593',
	'\u0594',
	'\u0595',
	'\u0597',
	'\u0598',
	'\u0599',
	'\u059C',
	'\u059D',
	'\u059E',
	'\u059F',
	'\u05A0',
	'\u05A1',
	'\u05A8',
	'\u05A9',
	'\u05AB',
	'\u05AC',
	'\u05AF',
	'\u05C4',
	'\u0610',
	'\u0611',
	'\u0612',
	'\u0613',
	'\u0614',
	'\u0615',
	'\u0616',
	'\u0617',
	'\u0657',
	'\u0658',
	'\u0659',
	'\u065A',
	'\u065B',
	'\u065D',
	'\u065E',
	'\u06D6',
	'\u06D7',
	'\u06D8',
	'\u06D9',
	'\u06DA',
	'\u06DB',
	'\u06DC',
	'\u06DF',
	'\u06E0',
	'\u06E1',
	'\u06E2',
	'\u06E4',
	'\u06E7',
	'\u06E8',
	'\u06EB',
	'\u06EC',
	'\u0730',
	'\u0732',
	'\u0733',
	'\u0735',
	'\u0736',
	'\u073A',
	'\u073D',
	'\u073F',
	'\u0740',
	'\u0741',
	'\u0743',
	'\u0745',
	'\u0747',
	'\u0749',
	'\u074A',
	'\u07EB',
	'\u07EC',
	'\u07ED',
	'\u07EE',
	'\u07EF',
	'\u07F0',
	'\u07F1',
	'\u07F3',
	'\u0816',
	'\u0817',
	'\u0818',
	'\u0819',
	'\u081B',
	'\u081C',
	'\u081D',
	'\u081E',
	'\u081F',
	'\u0820',
	'\u0821',
	'\u0822',
	'\u0823',
	'\u0825',
	'\u0826',
	'\u0827',
	'\u0829',
	'\u082A',
	'\u082B',
	'\u082C',
	'\u082D',
	'\u0951',
	'\u0953',
	'\u0954',
	'\u0F82',
	'\u0F83',
	'\u0F86',
	'\u0F87',
	'\u135D',
	'\u135E',
	'\u135F',
	'\u17DD',
	'\u193A',
	'\u1A17',
	'\u1A75',
	'\u1A76',
	'\u1A77',
	'\u1A78',
	'\u1A79',
	'\u1A7A',
	'\u1A7B',
	'\u1A7C',
	'\u1B6B',
	'\u1B6D',
	'\u1B6E',
	'\u1B6F',
	'\u1B70',
	'\u1B71',
	'\u1B72',
	'\u1B73',
	'\u1CD0',
	'\u1CD1',
	'\u1CD2',
	'\u1CDA',
	'\u1CDB',
	'\u1CE0',
	'\u1DC0',
	'\u1DC1',
	'\u1DC3',
	'\u1DC4',
	'\u1DC5',
	'\u1DC6',
	'\u1DC7',
	'\u1DC8',
	'\u1DC9',
	'\u1DCB',
	'\u1DCC',
	'\u1DD1',
	'\u1DD2',
	'\u1DD3',
	'\u1DD4',
	'\u1DD5',
	'\u1DD6',
	'\u1DD7',
	'\u1DD8',
	'\u1DD9',
	'\u1DDA',
	'\u1DDB',
	'\u1DDC',
	'\u1DDD',
	'\u1DDE',
	'\u1DDF',
	'\u1DE0',
	'\u1DE1',
	'\u1DE2',
	'\u1DE3',
	'\u1DE4',
	'\u1DE5',
	'\u1DE6',
	'\u1DFE',
	'\u20D0',
	'\u20D1',
	'\u20D4',
	'\u20D5',
	'\u20D6',
	'\u20D7',
	'\u20DB',
	'\u20DC',
	'\u20E1',
	'\u20E7',
	'\u20E9',
	'\u20F0',
	'\u2CEF',
	'\u2CF0',
	'\u2CF1',
	'\u2DE0',
	'\u2DE1',
	'\u2DE2',
	'\u2DE3',
	'\u2DE4',
	'\u2DE5',
	'\u2DE6',
	'\u2DE7',
	'\u2DE8',
	'\u2DE9',
	'\u2DEA',
	'\u2DEB',
	'\u2DEC',
	'\u2DED',
	'\u2DEE',
	'\u2DEF',
	'\u2DF0',
	'\u2DF1',
	'\u2DF2',
	'\u2DF3',
	'\u2DF4',
	'\u2DF5',
	'\u2DF6',
	'\u2DF7',
	'\u2DF8',
	'\u2DF9',
	'\u2DFA',
	'\u2DFB',
	'\u2DFC',
	'\u2DFD',
	'\u2DFE',
	'\u2DFF',
	'\uA66F',
	'\uA67C',
	'\uA67D',
	'\uA6F0',
	'\uA6F1',
	'\uA8E0',
	'\uA8E1',
	'\uA8E2',
	'\uA8E3',
	'\uA8E4',
	'\uA8E5',
	'\uA8E6',
	'\uA8E7',
	'\uA8E8',
	'\uA8E9',
	'\uA8EA',
	'\uA8EB',
	'\uA8EC',
	'\uA8ED',
	'\uA8EE',
	'\uA8EF',
	'\uA8F0',
	'\uA8F1',
	'\uAAB0',
	'\uAAB2',
	'\uAAB3',
	'\uAAB7',
	'\uAAB8',
	'\uAABE',
	'\uAABF',
	'\uAAC1',
	'\uFE20',
	'\uFE21',
	'\uFE22',
	'\uFE23',
	'\uFE24',
	'\uFE25',
	'\uFE26',
	'\U00010A0F',
	'\U00010A38',
	'\U0001D185',
	'\U0001D186',
	'\U0001D187',
	'\U0001D188',
	'\U0001D189',
	'\U0001D1AA',
	'\U0001D1AB',
	'\U0001D1AC',
	'\U0001D1AD',
	'\U0001D242',
	'\U0001D243',
	'\U0001D244',
}

# 通过go语言实现区块链--工作量证明
## 引言
在上一篇文章我们构建一个简单的数据结构，其实质上就是一个区块链数据库。并且我们使得他们有了链式关系：每一个区块指向上一个区块。遗憾的是我们实现的区块链还有一个非常大的缺陷：可以非常简答和廉价的添加区块。但这是区块链和比特币一个基本的原则那就是添加一个区块是非常困难的，要经过大量的算法。今天我们俩完善这个缺点。

## 工作量证明
区块链一个核心的点就是要通过进行大量的工作才能实现添加数据到链上，这些大量的工作来保证了区块链的安全和一致性。因此，会给做这些大量工作的人一些奖励（矿工就是这么获得币）。

这个机制非常像我们的真实生活：一个人努力工作的人获得相应的报酬来维持他们的生活。在区块链里面，矿工们通过提供算力来维持网络安全和添加新的区块到网络里面，相应的他们也会得到一些工作奖励。它们的工作来保证了一个区块可以安全的添加到区块链里面去来维护了一个稳定的区块链数据库。值得注意的是，完成这项工作的人必须证明这一点。

这种通过"努力工作去证明"的机制被称为：工作量证明。它之所以难是因为需需要大量的计算机算力：即使性能很强的电脑也不能很快的完成。而且，计算难度会随着算力的增加而增加来保证平均每个小时出6个区块。在比特币中，这种工作的目标是为一个块找到一个符合要求的Hash值。 这就是这个Hash值作为证明。 因此，找到一个证明就是实际的工作。

最后一个需要注意的事情。工作量证明的算法必须满足：工作难度大、但是验证简单。证明通常交给其他人，所以对他们来说，不需要太多时间来验证它。

## Hashing
在这个段落，我们将讨论hashing.如果你已经非常熟悉这个概念了，你可以跳过。

哈希是为指定数据获取散列的过程。 哈希是它计算数据的唯一表示。 散列函数是一种可以获取任意大小数据并生成固定大小散列的函数。 以下是哈希的一些关键特性：

1、原始数据无法从散列中恢复。 因此，哈希不是加密。
2、相同的数据只能有一个散列，散列是唯一的。
3、更改输入数据中的一个字节将导致完全不同的散列。
![](media/15236744175870/15236791451615.png)
哈希方法在验证数据一致性上已经得到了广泛的应用。一些软件会为他们的软件包提供一个公开的校验和，当用户下载完软件可以通过哈希方法计算出一个值和软件提供方公开的校验和对比来确认下载的软件是否正确。

在区块链中，散列用于保证块的一致性。 散列算法的输入数据包含前一个块的散列，因此不可能（或者至少非常困难）修改链中的块：必须重新计算其后的所有块的散列和散列。

## Hashcash
比特币使用的是Hashcash工作量证明算法，起初这个算法是为了防止垃圾邮件而开发的。算法实现过程可以拆分为一下几个步骤：
1、拿到一些大家都知道的数据（例如email例子里面的邮箱地址，比特币里面的区块头）；
2、添加一个计数器，这个计数器的起始值为0；
3、通过哈希方法计算数据+计数器组合后的数据得到哈希值；
4、校验这个哈希值是否等于特定的一个值：
  - 如果相等，完成
  - 如果不相等，给计数器+1然后重复步骤3和4.

因此，这是一个蛮力算法：改变计数器的值、计算新的哈希值、校验、计数器+1、计算哈希值，就这样循环。这也就是为什么这个非常消耗算力原因 。

现在我们来看看哈希必须满足的要求。在Hashcash的最初实现版本中，需求要求“必须满足哈希值的前20位必须为0”。在比特币中，要求是动态调整的，因为根据设计10min左右必须产生出一个快，因此随着加入计算的矿工越多计算难度也会提升。

为了演示这个算法，我通过之前例子的数据并且找到一个前三位为0hash值：
![](media/15236744175870/15236909477749.png)
`ca07ca` 是计数器的十六进制值， 转换成十进制的值为`13240266` 。

## 实现
好了，我们完成了理论部分，现在让我们开始编码吧！首先，我们定一个挖矿的难度：

```
const targetBits = 24
```
在比特币中，"target bits"存储在区块header里面表示这个区块挖矿难度，我们暂时先不实现难度调整算法，因此我们先可以定一个全局的难度常量。

24是一个随便写的一个数字，我们的目标是在内存中占用少于 256 位的目标。 我们希望难度足够大，但不要太大，因为难度越大，找到合适的散列越困难。

```
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}
```
这里我们创建了`ProofOfWork`的数据结构它里面包含了一个`block`和`target`指针。 “target” 是前一段描述的要求的另一个名称。 我们使用一个大整数，因为我们将哈希与目标进行比较：我们将哈希转换为大整数，并检查它是否小于目标。

在`NewProofOfWork`这个方法里面，我们初始化`big.Int`的值为1左移256位, 256的长度是SHA-256 hash的长度，SHA-256哈希算法就是我们要使用的。16进制显示出来`target`如下：

```
0x10000000000000000000000000000000000000000000000000000000000
```
它在内存中占用 29 个字节。 下面是它与前面例子中的哈希值的视觉比较：

```
0fac49161af82ed938add1d8725835cc123a1a87b1b196488360e58d4bfb51e3
0000010000000000000000000000000000000000000000000000000000000000
0000008b0f41ec78bab747864db66bcb9fb89920ee75f43fdaaeb5544f7f76ca
```
第一个散列（根据 “I like donuts” 计算）大于目标，因此它不是有效的工作证明。 第二个散列（根据 “I like donutsca07ca” 计算）小于目标，因此这是一个有效的证明。

你可以将目标视为范围的上限：如果数字（哈希）低于边界，则该数字有效，反之亦然。 降低边界将导致更少的有效数字，因此，找到一个有效数字需要更困难的工作。

现在，我们需要数据进行哈希，我们来创建准备数据的方法：

```
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
```
这件事很简单：我们只是将块字段与目标和随机数合并。 nonce 这里是来自上面 Hashcash 描述的计数器，这是一个加密术语。

好的，所有的准备工作都完成了，我们来实现 PoW 算法的核心：

```
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}
```
首先，我们初始化变量：hashInt 是散列的整数表示; nonce 是计数器。 接下来，我们运行一个 “无限” 循环：它受限于 maxNonce，它等于 math.MaxInt64; 这样做是为了避免 nonce 可能溢出。 尽管我们的 PoW 实现的难度太低而不能溢出，但为了以防万一，进行这种检查还是更好。

我们在循环中：
1. 准备数据
2. 把通过SHA-256进行哈希计算
3. 转换哈希值为大整数
4. 比较大整数和目标值

和前面解释的一样简单。 现在我们可以删除`Block` 的 `SetHash` 方法并修改 `NewBlock` 函数：

```
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
```
在这里你可以看到 nonce 被保存为 Block 属性。 这是必要的，因为 nonce 需要验证证据。 块结构现在看起来如此：

```
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}
```
好的！ 让我们运行该程序，看看是否一切正常:

```
Mining the block containing "Genesis Block"
00000041662c5fc2883535dc19ba8a33ac993b535da9899e593ff98e1eda56a1

Mining the block containing "Send 1 BTC to Ivan"
00000077a856e697c69833d9effb6bdad54c730a98d674f73c0b30020cc82804

Mining the block containing "Send 2 more BTC to Ivan"
000000b33185e927c9a989cc7d5aaaed739c56dad9fd9361dea558b9bfaf5fbe

Prev. hash:
Data: Genesis Block
Hash: 00000041662c5fc2883535dc19ba8a33ac993b535da9899e593ff98e1eda56a1

Prev. hash: 00000041662c5fc2883535dc19ba8a33ac993b535da9899e593ff98e1eda56a1
Data: Send 1 BTC to Ivan
Hash: 00000077a856e697c69833d9effb6bdad54c730a98d674f73c0b30020cc82804

Prev. hash: 00000077a856e697c69833d9effb6bdad54c730a98d674f73c0b30020cc82804
Data: Send 2 more BTC to Ivan
Hash: 000000b33185e927c9a989cc7d5aaaed739c56dad9fd9361dea558b9bfaf5fbe
```
好极了！ 您可以看到每个散列现在都以三个零字节开始，并且需要一些时间才能获得这些散列。

还有一件事要做：让我们可以验证作品的证明

```
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
```
这就是我们需要保存的随机数的地方。

我们再一次减产一切是否OK：

```
func main() {
	...

	for _, block := range bc.blocks {
		...
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

```
输出：

```
...

Prev. hash:
Data: Genesis Block
Hash: 00000093253acb814afb942e652a84a8f245069a67b5eaa709df8ac612075038
PoW: true

Prev. hash: 00000093253acb814afb942e652a84a8f245069a67b5eaa709df8ac612075038
Data: Send 1 BTC to Ivan
Hash: 0000003eeb3743ee42020e4a15262fd110a72823d804ce8e49643b5fd9d1062b
PoW: true

Prev. hash: 0000003eeb3743ee42020e4a15262fd110a72823d804ce8e49643b5fd9d1062b
Data: Send 2 more BTC to Ivan
Hash: 000000e42afddf57a3daa11b43b2e0923f23e894f96d1f24bfd9b8d2d494c57a
PoW: true
```

## 结论
我们的区块链更接近其实际架构：现在添加区块需要努力工作，因此可以进行挖掘。 但它仍然缺乏一些关键特征：区块链数据库不是持久的，没有钱包，地址，交易，并且没有共识机制。 所有这些我们将在未来的文章中实现，而现在，快乐的挖掘！



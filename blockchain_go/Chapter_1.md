# 通过go语言自己实现区块链--区块
## 引言
区块链是21世纪最具革命性的技术之一，并且在未来它的潜力将还无限大。其实，区块链就是一个分布式的数据库。但是它的数据不是私有的是公开的任何人都可以拥有部分或者全部的数据。而且新增一条记录需要数据库维护者统一。因此区块链使得数字货币和只能合约成为可能。

在接下来的系列文章里面，我们将通过区块链技术构建一个简单的加密数字货币

## 区块
我们从区块链的“区块”开始。在区块链“区块”是进行有效信息存储的，例如，比特币的区块用于存储存储交易，本质上所有的加密货币这是这样的。除此之外，一个区块包含了一些技术信息，例如版本号、当前时间戳、和上一个区块的hash值。

在这篇文章里面我们实现不了我们描述的一个完成的区块链或者类似比特币的功能，我们先通过一个简单的版本，在区块里面包含机械重要的信息，这个版本的区块信息如下：
```
type Block struct {
  // 当前时间戳 (区块创建时间)
	Timestamp     int64
	// 区块里面包含的真实有效的信息
	Data          []byte
	// 上一个区块的的hash值
	PrevBlockHash []byte
	// 当前区块头数据的hash值
	Hash          []byte
}
```
真实的区块数据要比这个复杂，这里我们做了一些简单话的处理。

要怎么出技术hash值了？`hash`值的计算在区块链里面非常重要，它的特性保证了区块链的安全。通过电脑等去大量的计算hash值是一个比较难的事情，即使使用最快的电脑也要计算一会(这也是为什么很多人通过高性能的`GPU`进行挖矿)。特意设计的架构，使得产生一个新的区块是困难的，由此来阻止人们去修改已经添加进去的数据。我们将在接下来的文章里面实现这些机制。

目前，我们将`Block`里面的所有数据拼接起来，然后通过通过SHA-256计算去Hash值，我们来实现  `SetHash`方法：
```
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}
```

接下来，我们按照go习惯，实现一个简单的创建区块的方法
```
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}
```
上面我们实现的就是区块

## 区块链
现在我们来实现区块链，本质上他就是一个包含了`Block`区块数据的有序列表。也就是说每次要增加一个新的区块时我们都是按照顺序添加到最前面。这个结构必须可以快速的找到最新的一个区块并且获取到它的`hash`值.

在golang里面我们可以通过数组和`map`来实现：数组保证顺序，map通过存储`hash->block`对。这里我们只使用数组来实现，因为我们现在还用不到通过hash获取block的功能。

```
type Blockchain struct {
	blocks []*Block
}
```
这就是我们实现的第一个区块链，我们从来没有想到会如此简单

现在我们来实现可以向区块链里面添加区块：
```
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
```
添加一个区块之前我们先需要有一个区块，但是目前我们的区块链里面还没有区块。因此，在任何区块链里面至少需要一个区块，第一个在链里面的块通常称之为：`genesis block`.我们实现一个方法来创建这样一个区块：
```
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
```
接下来，我们实现一个方法通过我们的genesis区块创建一个区块链：
```
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
```
通过下面我们来测试一下区块链工作工作是否正常：
```
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
```
终端输出:
```
Prev. hash:
Data: Genesis Block
Hash: aff955a50dc6cd2abfe81b8849eab15f99ed1dc333d38487024223b5fe0f1168

Prev. hash: aff955a50dc6cd2abfe81b8849eab15f99ed1dc333d38487024223b5fe0f1168
Data: Send 1 BTC to Ivan
Hash: d75ce22a840abb9b4e8fc3b60767c4ba3f46a0432d3ea15b71aef9fde6a314e1

Prev. hash: d75ce22a840abb9b4e8fc3b60767c4ba3f46a0432d3ea15b71aef9fde6a314e1
Data: Send 2 more BTC to Ivan
Hash: 561237522bb7fcfbccbc6fe0e98bbbde7427ffe01c6fb223f7562288ca2295d1
```

## 结论
我们构建了一个非常简单的区块链：它仅仅通过一个数组存储区块，让一个区块链接上一个区块。真实的区块链要比这个负责的多。在我们的程序里面添加区块是非常简单和快的，但是在正式的区块链添加一个区块是需要进行一些计算的：需要去执行大量的计算才能获得添加快的权限(这个机制成为工作量证明`Proof-of-Work`)。区块链是一个分布式数据库，并且没有单一决策者。因此，要加入一个新块，必须要被网络的其他参与者确认和同意（这个机制叫做共识（consensus））。还有一点，我们的区块链还没有任何的交易！

## 代码测试
```
$ cd code
$ go run *.go
```

## 参考
[Building Blockchain in Go. Part 1: Basic Prototype](https://jeiwan.cc/posts/building-blockchain-in-go-part-1/)
[Block hashing algorithm](https://en.bitcoin.it/wiki/Block_hashing_algorithm)









package main

// Definindo a estrutura de um nó da árvore trie
type TrieNode struct {
	children map[uint64]*TrieNode // Mapeia os filhos do nó atual
	isEnd    bool                 // Indica se o nó é o final de uma palavra
	hashes   []uint64             // Armazena todos os hashes no nó final
}

// Função para criar um novo nó da árvore trie
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[uint64]*TrieNode),
		isEnd:    false,
	}
}

// Função para inserir um hash na árvore trie
func (root *TrieNode) Insert(hash uint64) {
	node := root
	// Itera sobre os 64 bits do hash
	for i := 0; i < 64; i++ {
		bit := (hash >> i) & 1 // Obtém o i-ésimo bit do hash
		// Verifica se o filho correspondente ao bit atual não existe
		if _, ok := node.children[bit]; !ok {
			// Se não existe, cria um novo nó para o bit
			node.children[bit] = NewTrieNode()
		}
		node = node.children[bit] // Move para o próximo nível na árvore
	}
	node.isEnd = true                       // Marca o nó como o final de uma palavra
	node.hashes = append(node.hashes, hash) // Adiciona o hash à lista de hashes no nó final
}

// Função para buscar hashes com uma distância de Hamming limitada em relação a um hash de referência
func (root *TrieNode) Search(hash uint64, distance int) []uint64 {
	var result []uint64                           // Inicializa a lista de resultados vazia
	root.searchHelper(hash, 0, distance, &result) // Chama a função auxiliar de busca
	return result                                 // Retorna os hashes encontrados
}

// Função auxiliar para buscar hashes com uma distância de Hamming limitada de forma recursiva
func (node *TrieNode) searchHelper(hash uint64, index, distance int, result *[]uint64) {
	// Verifica se atingimos o final do hash (todos os bits foram processados)
	if index == 64 {
		*result = append(*result, node.hashes...) // Adiciona todos os hashes encontrados ao resultado
		return
	}

	// Itera sobre os filhos do nó atual
	for bit := range node.children {
		// Verifica se o bit do filho não corresponde ao bit do hash atual e ainda temos distância disponível
		if bit != ((hash>>index)&1) && distance > 0 {
			// Se sim, faz uma chamada recursiva com a distância reduzida
			node.children[bit].searchHelper(hash, index+1, distance-1, result)
			// Verifica se o bit do filho corresponde ao bit do hash atual
		} else if bit == ((hash >> index) & 1) {
			// Se sim, faz uma chamada recursiva sem reduzir a distância
			node.children[bit].searchHelper(hash, index+1, distance, result)
		}
	}
}

// func main() {
// 	tree := NewTrieNode()
// 	hashes := []uint64{12486923320150515769, 12254001406827791480, 12257299735051066979}
// 	for _, hash := range hashes {
// 		tree.Insert(hash)
// 	}

// 	// Buscando hashes com uma distância de hamming limitada em relação a um hash de referência
// 	referenceHash := uint64(12486922770394827888)
// 	distance := 16
// 	fmt.Printf("Buscando hashes com distância de Hamming %d em relação a %d:\n", distance, referenceHash)
// 	foundHashes := tree.Search(referenceHash, distance)
// 	if len(foundHashes) > 0 {
// 		fmt.Println("Encontrado! Hashes:")
// 		for _, h := range foundHashes {
// 			fmt.Printf("Hash: %d\n", h)
// 		}
// 	} else {
// 		fmt.Println("Não encontrado!")
// 	}
// }

package main


import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"github.com/1amkaizen/pdfconverter"

)




func main() {
	// Inisialisasi variabel
	var stopwordsFile string
	var kata int

	// Mengambil argumen dengan flag
	flag.StringVar(&stopwordsFile, "s", "", "File stopwords")
	flag.IntVar(&kata, "t", 0, "Jumlah kata")
	flag.Parse()

	// Memeriksa jumlah argumen yang diberikan
	if flag.NArg() != 1 {
		fmt.Println("Usage: ./GoTextMiner -s <stopwords_file> -t <jumlah kata> <namafile>")
		return
	}

	// File input yang akan diproses (PDF atau teks)
	filePath := flag.Arg(0)

	// Memeriksa apakah file PDF
	isPDF := strings.HasSuffix(strings.ToLower(filePath), ".pdf")

	// Konversi PDF ke teks jika diperlukan
	
	
	// ...
var textFilePath string
if isPDF {
    textFilePath = pdfToText(filePath)
    if textFilePath == "" {
        fmt.Println("\nGagal melakukan konversi PDF ke teks.")
        return
    }
} else {
    textFilePath = filePath
}
// ...


	// Memeriksa apakah file stopwords telah disediakan
	if stopwordsFile == "" {
		fmt.Println("Gunakan opsi -s untuk menyediakan file stopwords.")
		return
	}

	// Memeriksa apakah file stopwords tersedia
	if _, err := os.Stat(stopwordsFile); os.IsNotExist(err) {
		fmt.Printf("File stopwords '%s' tidak ditemukan.\n", stopwordsFile)
		return
	}

	// Memeriksa apakah jumlah kata telah disediakan
	if kata <= 0 {
		fmt.Println("Gunakan opsi -t untuk menyediakan jumlah kata.")
		return
	}


	// Baca isi file teks
	file, err := os.Open(textFilePath)
	if err != nil {
		fmt.Printf("Gagal membuka file: %v\n", err)
		return
	}
	defer file.Close()

	// Baca file stopwords
	stopwords, err := readStopwords(stopwordsFile)
	if err != nil {
		fmt.Printf("Gagal membaca file stopwords: %v\n", err)
		return
	}

	// Menghapus kata-kata yang ada di stopwords dari teks
	filteredWords := filterWords(file, stopwords)

	// Hitung frekuensi kata-kata
	wordCounts := make(map[string]int)
	for _, word := range filteredWords {
		wordCounts[word]++
	}

	// Urutkan kata-kata berdasarkan frekuensi
	sortedWords := sortWordsByFrequency(wordCounts)

	// Ambil jumlah kata yang diinginkan
	topWords := sortedWords[:kata]

	// Menampilkan hasil
	for _, word := range topWords {
		fmt.Printf("%s: %d\n", word.word, word.count)
	}

	// Hapus file teks sementara jika hasil dari konversi PDF
	if isPDF {
		err := os.Remove(textFilePath)
		if err != nil {
			fmt.Printf("Gagal menghapus file teks sementara: %v\n", err)
		}
	}

}



func pdfToText(pdfFilePath string) string {
    // Menggunakan pdfconverter.ConvertPDFToText untuk mengkonversi PDF ke teks
    outputFilePath := "temp.txt"
    err := pdfconverter.ConvertPDFToText(pdfFilePath, outputFilePath)
    if err != nil {
        fmt.Printf("Gagal melakukan konversi PDF ke teks: %v\n", err)
        return ""
    }
    return outputFilePath
}


// Fungsi lain sama seperti sebelumnya...
type wordFrequency struct {
	word  string
	count int
}

func readStopwords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stopwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopwords = append(stopwords, strings.ToLower(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopwords, nil
}

func filterWords(file *os.File, stopwords []string) []string {
	var filteredWords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		for _, word := range words {
			wordLower := strings.ToLower(word)
			if !contains(stopwords, wordLower) {
				filteredWords = append(filteredWords, wordLower)
			}
		}
	}

	return filteredWords
}

func sortWordsByFrequency(wordCounts map[string]int) []wordFrequency {
	var sortedWords []wordFrequency
	for word, count := range wordCounts {
		sortedWords = append(sortedWords, wordFrequency{word, count})
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].count > sortedWords[j].count
	})

	return sortedWords
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

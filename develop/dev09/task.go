package main

/*
=== Утилита wget ===
Реализовать утилиту wget с возможностью скачивать сайты целиком
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

var (
	errInvalidMaxDepth = errors.New("wget: max depth must be positive")
	//errInvalidUrl = errors.New("wget: url is invalid")
	errInvalidMediaType = errors.New("wget: invalid media type")
	maxDepth int // Максимальная "глубина" скачивания (меньше может быть (если на каком-то из уровней глубины не оказалось url'ов), больше не может)
)

func init() {
	flag.IntVar(&maxDepth, "d", 0, `Максимальная "глубина" скачивания`)
}
func main() {
	flag.Parse()
	var url, folderPath string

	url = os.Args[len(os.Args) - 1] // Получаем указатель ресурса из командной строки
	workDir, _ := os.Getwd()
	folderPath = fmt.Sprint(workDir, `\testDirectory\`)
	fmt.Println(folderPath)
	//maxDepth := 0

	err := Download(url, folderPath, maxDepth)
	if err != nil {
		fmt.Println(err)
	}
}

func Download(resourceLocator, folderPath string, maxDepth int) error {
	if maxDepth < 0 { // Определяется максимальная "глубина" скачивания (внутри скачанной страницы могут быть гиперссылки, ресурсы, на которые они указывают, тоже надо скачать)
		return errInvalidMaxDepth // maxDepth также используется для определения крайнего случая рекурсии
	}

	//if _, e := url.ParseRequestURI(resourceLocator); e != nil { // ParseRequestURI парсит необработанный URL-адрес в структуру URL, если не удалось распарсить, возвращаем ошибку
	//	return errInvalidUrl
	//}
	// Получение содержимого страницы
	// ------------------------------
	fmt.Println("Downloading: ", resourceLocator)
	// URL-адрес нужного сайта передается функции http.Get. Она возвращает объект http.Response, а также любые обнаруженные ошибки
	res, err := http.Get(resourceLocator)
	if err != nil {
		log.Println(err)
		return err
	}
	// Объект http.Response представляет собой структуру с полем Body, представляющим содержимое страницы
	defer res.Body.Close()
	// Метод Close используется для освобождения сетевого подключения при завершении работы
	// Вызов Close откладывается при помощи defer, чтобы подключение было освобождено после завершения чтения данных
	data, err := ioutil.ReadAll(res.Body) // Поскольку конкретное значение поля Body имеет методы Read и Close,
	// мы можем передать это значение в качестве аргумента функции ReadAll (её параметр - интерфейсного типа)
	// ReadAll читает все данные из ответа
	if err != nil {
		log.Println(err)
		return err
	}
	// Создание файла, в который помещается содержимое страницы
	// --------------------------------------------------------
	filename, _ := getFilename(resourceLocator, res.Header) // Создаёт и возвращает имя нового файла, включая его расширение
	outputFile, err := os.Create(folderPath + filename) // Create создает значение типа File. Его методы можно использовать для ввода/вывода (например, метод Write (запись в файл))
	if err != nil {
		return err
	}
	defer outputFile.Close() // Close закрывает файл, делая его недоступным для ввода/вывода
	// Вызов Close откладывается
	outputFile.Write(data)
	// Распараллелим нашу программу
	links := getParseLinks(data) // Находим все линки в полученном html
	var wg sync.WaitGroup
	for _, sublink := range links { // Под загрузку каждого саблинка будет создана горутина, в которой загрузка будет осуществляться
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			Download(link, folderPath, maxDepth-1) // Рекурсивно вызываем Download в новой горутине для загрузки страниц по саблинкам
		}(sublink)
	}
	wg.Wait() // Горутина main останется заблокированной, пока ВСЕ остальные горутины не вызовут метод Done значениия wg и счётчик не уменьшиться до 0
	fmt.Println(resourceLocator, "successfully downloaded")
	return nil
}
// getFilename создаёт и возвращает имя нового файла, включая его расширение
func getFilename(url string, header http.Header) (string, error) { // Значение типа Header является набором пар "ключ-значение" (map) из http-header'а
	contentType := header.Get("Content-Type") // Get возвращает значение, соответствующее переданному ключу
	mediaType, _, _ := mime.ParseMediaType(contentType) // MIME тип состоит из типа и подтипа — двух строк разделённых наклонной чертой "text/html"
	fmt.Println(mediaType)
	subType := Cut(mediaType, "/", 1) // Разделяем тип и подтип и возвращаем ПОДТИП
	baseName := path.Base(url) // Возвращает последний элемент пути

	if subType == "" {
		return "", errInvalidMediaType
	}

	newName := Cut(baseName, "?", 0) // Поскольку нельзя создать файл с именем, содержащим "?", если "?" содержится в строке, ф-ция вернет часть строки до "?"

	if path.Ext(newName) == "" && subType != "" { // Если Ext() не нашла расширения в имени (.[расширение]) и подтип - не пустая строка,
		return newName + "." + subType, nil // Конкатенируем эти две строки и получаем имя файла с расширением
	}

	return newName, nil // Иначе расширение в имени и так присутствует
}
// Cut разделяет MIME тип на тип и подтип, если в строке присутствует наклонная черта, возвращает указанное (тип (0) или подтип (1))
// Кроме того, ф-ция универсальная и используется не только для вышеописанного
func Cut(s, sep string, i int) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[i]
	}
	return s
}
// getParseLinks находит все линки в полученном html и возвращает слайс с ними
func getParseLinks(data []byte) []string {
	reg := regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`) // для регулярных выражений лучше использовать необработанные строки (raw strings, строки без интерпретации экранированных литералов)
	result := reg.FindAll(data, -1)

	subUrls := make([]string, len(result)) // Поскольку искали совпадения в слайсе байт, нужно привести слайс байтов с совпадениями к слайсу строк
	for i := 0; i < len(result); i++ {
		subUrls = append(subUrls, string(result[i]))
	}

	return subUrls
}
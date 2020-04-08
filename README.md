## Onerilen tool'lar:
    - golint: Formatlama icin

## Test icin denenecek komutlar:
#### Alt paketlerdeki testleri calistirmak icin:
- Tum testler:

        $ go test -v ./api/model/...
- Tek test dosyasi:

        $ go test -v ./api/model/user.go ./api/model/user_test.go
#### Coverage:
    $ go test -covermode count -coverprofile <coverage report name>
    $ go tool cover -html=<coverage report name>

## Doc icin referans:
- https://golang.org/src/go/doc/example.go
### Local ortamda doc sayfasina erismek icin:
    $ godoc -http=localhost:6060
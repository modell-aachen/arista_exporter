# arista_exporter

Prometheus exporter for Arista EOS devices

## Note

To ensure compatibility with Arista EOS API responses and to make BGP metric collection work correctly, you need to patch the `goeapi` vendor files after running `go mod vendor`:

```bash
go mod vendor
sed -i 's/\(Version[[:space:]]*\)int/\1string/' ./vendor/github.com/aristanetworks/goeapi/eapi.go
sed -i 's/p := Parameters{1, commands, encoding}/p := Parameters{"latest", commands, encoding}/' ./vendor/github.com/aristanetworks/goeapi/eapilib.go
```

## Credits

- https://github.com/aristanetworks/goeapi
- https://github.com/henrikvtcodes/eoxporter
- https://github.com/ubccr/arista_exporter

install cyclonedx:
	go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.4.0
sbom:
	cyclonedx-gomod mod -licenses -verbose=false -test -output file/bom.xml -output-version 1.4
swag-init:
	swag init -g cmd/main.go -o api/docs
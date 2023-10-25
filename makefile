install cyclonedx:
	
create sbom:
	cyclonedx-gomod mod -licenses -type library -json -output bom.json
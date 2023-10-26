install cyclonedx:
	go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.4.0
sbom:
	cyclonedx-gomod mod -licenses -verbose=false -test -output file/bom.xml -output-version 1.4
swag-init:
	swag init -g cmd/main.go -o api/docs
get-detect:
	curl https://detect.synopsys.com/detect8.sh --output detect8.sh
scan:
	/bin/bash ./detect8.sh --blackduck.signature.scanner.memory=4096 --detect.timeout=6000 --blackduck.trust.cert=true --logging.level.com.synopsys.integration=DEBUG --blackduck.url=https://sap-staging.app.blackduck.com --blackduck.api.token=****** "--detect.project.name=test_OS3_sbom_upload1" "--detect.project.version.name=localScan" --detect.policy.check.fail.on.severities=NONE "--detect.force.success.on.skip=true" --detect.source.path="path_to_project" --detect.tools=SIGNATURE_SCAN,DETECTOR

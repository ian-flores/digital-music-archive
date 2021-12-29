all: build_data-extracter

clean: 
	rm -r build/
	mkdir build/

build_data-extracter:
	chmod +x ./secrets/data-extracter.env
	./secrets/data-extracter.env
	go build ./cmd/data-extracter
	mv ./data-extracter build/